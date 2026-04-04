package common

import (
	"context"
	"fmt"
	"net"
	"sync/atomic"
	"time"

	"google.golang.org/protobuf/proto"

	"im.turms/server/internal/domain/common/constant"
	"im.turms/server/internal/domain/common/access/servicerequest/dto"
	"im.turms/server/internal/domain/gateway/session"
	"im.turms/server/internal/infra/exception"
	"im.turms/server/internal/infra/metrics"
	"im.turms/server/internal/infra/tracing"
	"im.turms/server/pkg/protocol"
)

// Errors are centralized in channel_handler.go
var HeartbeatFailureRequestId int64 = -100

// Placeholder interfaces and structs pending migration
// Service definitions

type BlocklistService interface {
	TryBlockUserIdForCorruptedRequest(userId int64)
	TryBlockIpForCorruptedRequest(ip []byte)
	TryBlockIpForFrequentRequest(ip []byte)
	TryBlockUserIdForFrequentRequest(userId int64)

	IsIpBlocked(ip []byte) bool
	TryBlockIpForCorruptedFrame(ip []byte)
	TryBlockUserIdForCorruptedFrame(userId int64)
}

type SessionClientController interface {
	HandleCreateSessionRequest(ctx context.Context, wrapper *UserSessionWrapper, req *protocol.CreateSessionRequest) (*RequestHandlerResult, error)
	HandleDeleteSessionRequest(ctx context.Context, wrapper *UserSessionWrapper) (*RequestHandlerResult, error)
}

type SessionService interface {
	HandleHeartbeatUpdateRequest(session *session.UserSession)
	GetLocalUserSessionsByIp(ip []byte) []*session.UserSession
}

type ServiceRequestService interface {
	HandleServiceRequest(ctx context.Context, session *session.UserSession, req *dto.ServiceRequest) (*protocol.TurmsNotification, error)
}

type ServiceAvailability struct {
	Available bool
	Reason    string
}

type ServerStatusManager interface {
	GetServiceAvailability() ServiceAvailability
}

type ApiLoggingContext interface {
	ShouldLogHeartbeatRequest() bool
	ShouldLogRequest(requestType interface{}) bool
	LogRequest(sessionID *int32, userID *int64, deviceType *protocol.DeviceType, version *int32, ipStr string, requestID int64, requestType interface{}, requestSize int, requestTime int64, notification *protocol.TurmsNotification, processingTimeMilli int64)
}

// ClientRequestDispatcher routes incoming client requests.
type ClientRequestDispatcher struct {
	ApiLoggingContext     ApiLoggingContext
	BlocklistService      BlocklistService
	IpRequestThrottler    *IpRequestThrottler
	SessionController     SessionClientController
	SessionService        SessionService
	ServiceRequestService ServiceRequestService
	ServerStatusManager   ServerStatusManager
	MetricsService        metrics.MetricsService

	NotificationFactory *NotificationFactory

	PendingRequestCount  atomic.Int32
	OnAllRequestsHandled func()
}

// HandleRequest wraps handleRequest0 with error handling.
// @MappedFrom handleRequest(UserSessionWrapper sessionWrapper, ByteBuf serviceRequestBuffer)
func (d *ClientRequestDispatcher) HandleRequest(ctx context.Context, sessionWrapper *UserSessionWrapper, serviceRequestBuffer []byte) (buff []byte, err error) {
	d.PendingRequestCount.Add(1)
	defer func() {
		if d.PendingRequestCount.Add(-1) == 0 && d.OnAllRequestsHandled != nil {
			d.OnAllRequestsHandled()
		}
		if r := recover(); r != nil {
			err = fmt.Errorf("panic in HandleRequest: %v", r)
		}
	}()
	return d.HandleRequest0(ctx, sessionWrapper, serviceRequestBuffer)
}

// HandleRequest0 processes the parsed request.
// @MappedFrom handleRequest0(UserSessionWrapper sessionWrapper, ByteBuf serviceRequestBuffer)
func (d *ClientRequestDispatcher) HandleRequest0(ctx context.Context, sessionWrapper *UserSessionWrapper, serviceRequestBuffer []byte) ([]byte, error) {
	if len(serviceRequestBuffer) == 0 {
		availability := d.ServerStatusManager.GetServiceAvailability()
		if !availability.Available {
			notif := d.NotificationFactory.CreateWithReason(&HeartbeatFailureRequestId, constant.ResponseStatusCode_SERVER_UNAVAILABLE, &availability.Reason)
			return proto.Marshal(notif)
		}
		return d.handleHeartbeatRequest(sessionWrapper)
	}

	requestTime := time.Now().UnixMilli()
	startTime := time.Now().UnixNano()
	requestSize := len(serviceRequestBuffer)

	req := &protocol.TurmsRequest{}
	var notification *protocol.TurmsNotification
	var requestType interface{}
	var requestID int64

	err := proto.Unmarshal(serviceRequestBuffer, req)
	canLogRequest := true

	if err != nil {
		if sessionWrapper.HasUserSession() {
			d.BlocklistService.TryBlockUserIdForCorruptedRequest(sessionWrapper.UserSession.UserID)
		}
		d.BlocklistService.TryBlockIpForCorruptedRequest(net.ParseIP(sessionWrapper.GetIPStr()))
		if req.RequestId != nil {
			requestID = *req.RequestId
		}

		requestType = "UNRECOGNIZED_REQUEST"
		exc := exception.NewTurmsError(int32(constant.ResponseStatusCode_INVALID_REQUEST), err.Error())
		notification = d.NotificationFactory.CreateFromError(exc, &requestID)

		// Java logs the corrupted request error with trace context if possible
		fmt.Printf("ERROR: Failed to handle the service request with UNRECOGNIZED_REQUEST: %v\n", err)
	} else {
		if req.RequestId != nil {
			requestID = *req.RequestId
		}
		requestType = req.GetKind()

		if sessionWrapper.HasUserSession() {
			tag := int32(0)
			if req.Kind != nil {
				tag = int32(req.ProtoReflect().WhichOneof(req.ProtoReflect().Descriptor().Oneofs().ByName("kind")).Number())
			}
			if !sessionWrapper.UserSession.HasPermission(tag) {
				exc := exception.NewTurmsError(int32(constant.ResponseStatusCode_UNAUTHORIZED_REQUEST), "")
				notification = d.NotificationFactory.CreateFromError(exc, &requestID)
			} else if _, ok := req.Kind.(*protocol.TurmsRequest_DeleteSessionRequest); ok {
				canLogRequest = sessionWrapper.UserSession.AcquireDeleteSessionRequestLoggingLock()
			}
		}

		if notification == nil {
			var traceId string
			if sessionWrapper.HasUserSession() {
				traceId = fmt.Sprintf("%d-%d", sessionWrapper.UserSession.UserID, requestID)
			} else {
				traceId = fmt.Sprintf("anon-%d", requestID)
			}
			ctxWithTrace := ctx
			if supportsTracing(requestType) {
				tc := tracing.NewTracingContext(traceId)
				ctxWithTrace = tracing.WithTracingContext(ctx, tc)
			}

			notification, err = d.HandleServiceRequest(ctxWithTrace, sessionWrapper, req, serviceRequestBuffer)
			if err != nil {
				tc := tracing.FromContext(ctxWithTrace)
				traceId := ""
				if tc != nil {
					traceId = tc.TraceId
				}
				fmt.Printf("ERROR: [%s] Failed to handle the service request: %v, error: %v\n", traceId, req, err)
				notification = d.NotificationFactory.CreateFromError(err, &requestID)
			}
		}
	}

	finalCanLogRequest := canLogRequest
	isServerError := constant.IsServerError(notification.GetCode())

	if isServerError && err == nil { // err!=nil handled above for corrupted
		fmt.Printf("ERROR: Failed to handle the service request: type=%T, code=%d\n", requestType, notification.GetCode())
	}

	if isServerError || (d.ApiLoggingContext.ShouldLogRequest(requestType) && finalCanLogRequest) {
		var version *int32
		var userId *int64
		var sessionId *int32
		var deviceType *protocol.DeviceType

		if sessionWrapper.HasUserSession() {
			userId = &sessionWrapper.UserSession.UserID
			devType := sessionWrapper.UserSession.DeviceType
			deviceType = &devType

			ver := sessionWrapper.UserSession.Version.Load()
			version = &ver

			sid := int32(sessionWrapper.UserSession.ID)
			sessionId = &sid
		}

		processingTimeMilli := (time.Now().UnixNano() - startTime) / 1000000
		d.ApiLoggingContext.LogRequest(sessionId, userId, deviceType, version, sessionWrapper.GetIPStr(), requestID, requestType, requestSize, requestTime, notification, processingTimeMilli)
	}

	processingTimeMilli := (time.Now().UnixNano() - startTime) / 1000000
	if d.MetricsService != nil {
		d.MetricsService.RecordRequest(req, requestSize, processingTimeMilli)
	}

	return proto.Marshal(notification)
}

// HandleServiceRequest delegates typed requests to backend services.
// @MappedFrom handleServiceRequest(UserSessionWrapper sessionWrapper, SimpleTurmsRequest request, ByteBuf serviceRequestBuffer, TracingContext tracingContext)
func (d *ClientRequestDispatcher) HandleServiceRequest(ctx context.Context, sessionWrapper *UserSessionWrapper, request *protocol.TurmsRequest, serviceRequestBuffer []byte) (*protocol.TurmsNotification, error) {
	requestID := request.GetRequestId()
	if requestID <= 0 {
		reason := "The request ID must be greater than 0"
		return d.NotificationFactory.CreateWithReason(&requestID, constant.ResponseStatusCode_INVALID_REQUEST, &reason), nil
	}

	availability := d.ServerStatusManager.GetServiceAvailability()
	if !availability.Available {
		return d.NotificationFactory.CreateWithReason(&requestID, constant.ResponseStatusCode_SERVER_UNAVAILABLE, &availability.Reason), nil
	}

	// Rate limiting
	now := time.Now().UnixNano()
	if !d.IpRequestThrottler.TryAcquireToken(sessionWrapper.GetIPStr(), now) {
		d.BlocklistService.TryBlockIpForFrequentRequest(net.ParseIP(sessionWrapper.GetIPStr()))
		if sessionWrapper.HasUserSession() {
			d.BlocklistService.TryBlockUserIdForFrequentRequest(sessionWrapper.UserSession.UserID)
		}
		return d.NotificationFactory.Create(&requestID, constant.ResponseStatusCode_CLIENT_REQUESTS_TOO_FREQUENT), nil
	}

	switch kind := request.Kind.(type) {
	case *protocol.TurmsRequest_CreateSessionRequest:
		result, err := d.SessionController.HandleCreateSessionRequest(ctx, sessionWrapper, kind.CreateSessionRequest)
		if err != nil {
			return nil, err
		}
		return d.getNotificationFromHandlerResult(result, requestID), nil
	case *protocol.TurmsRequest_DeleteSessionRequest:
		result, err := d.SessionController.HandleDeleteSessionRequest(ctx, sessionWrapper)
		if err != nil {
			return nil, err
		}
		return d.getNotificationFromHandlerResult(result, requestID), nil
	default:
		return d.handleGenericServiceRequest(ctx, sessionWrapper, request, serviceRequestBuffer)
	}
}

func (d *ClientRequestDispatcher) handleGenericServiceRequest(ctx context.Context, sessionWrapper *UserSessionWrapper, request *protocol.TurmsRequest, serviceRequestBuffer []byte) (*protocol.TurmsNotification, error) {
	if !sessionWrapper.HasUserSession() || sessionWrapper.UserSession.Conn == nil {
		reqID := request.GetRequestId()
		return d.NotificationFactory.SessionClosed(&reqID), nil
	}

	session := sessionWrapper.UserSession
	svcReq := &dto.ServiceRequest{
		Ip:           []byte(sessionWrapper.GetIPStr()),
		UserId:       session.UserID,
		DeviceType:   session.DeviceType,
		RequestId:    request.GetRequestId(),
		TurmsRequest: request,
		Type:         request.GetKind(),
		Buffer:       serviceRequestBuffer,
	}
	return d.ServiceRequestService.HandleServiceRequest(ctx, session, svcReq)
}

func (d *ClientRequestDispatcher) handleHeartbeatRequest(sessionWrapper *UserSessionWrapper) ([]byte, error) {
	var data []byte
	var isSuccess bool

	if sessionWrapper.HasUserSession() && sessionWrapper.UserSession.Conn != nil {
		d.SessionService.HandleHeartbeatUpdateRequest(sessionWrapper.UserSession)
		data = []byte{}
		isSuccess = true
	} else {
		notif := d.NotificationFactory.Create(&HeartbeatFailureRequestId, constant.ResponseStatusCode_UPDATE_HEARTBEAT_OF_NONEXISTENT_SESSION)
		data, _ = proto.Marshal(notif)
		isSuccess = false
	}

	if d.ApiLoggingContext.ShouldLogHeartbeatRequest() {
		var version *int32
		var userId *int64
		var sessionId *int32
		var deviceType *protocol.DeviceType

		if sessionWrapper.HasUserSession() {
			userId = &sessionWrapper.UserSession.UserID
			devType := sessionWrapper.UserSession.DeviceType
			deviceType = &devType

			ver := sessionWrapper.UserSession.Version.Load()
			version = &ver

			sid := int32(sessionWrapper.UserSession.ID)
			sessionId = &sid
		}

		var notification *protocol.TurmsNotification
		if !isSuccess {
			notification = d.NotificationFactory.Create(&HeartbeatFailureRequestId, constant.ResponseStatusCode_UPDATE_HEARTBEAT_OF_NONEXISTENT_SESSION)
		}

		logSuccess := 0
		if isSuccess {
			logSuccess = 1
		}
		d.ApiLoggingContext.LogRequest(sessionId, userId, deviceType, version, sessionWrapper.GetIPStr(), 0, "HEARTBEAT", 0, time.Now().UnixMilli(), notification, int64(logSuccess))
	}
	return data, nil
}

func (d *ClientRequestDispatcher) getNotificationFromHandlerResult(result *RequestHandlerResult, reqId int64) *protocol.TurmsNotification {
	notif := &protocol.TurmsNotification{
		Timestamp: time.Now().UnixMilli(),
		RequestId: &reqId,
		Code:      proto.Int32(int32(result.Code)),
		Data:      result.Data,
	}
	if result.Reason != "" {
		notif.Reason = proto.String(result.Reason)
	}
	return notif
}

func supportsTracing(requestType interface{}) bool {
	switch requestType.(type) {
	case *protocol.TurmsRequest_CreateSessionRequest, *protocol.TurmsRequest_DeleteSessionRequest, string:
		return false
	default:
		return true
	}
}
