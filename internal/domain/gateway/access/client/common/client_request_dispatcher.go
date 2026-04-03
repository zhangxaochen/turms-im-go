package common

import (
	"context"
	"errors"
	"fmt"
	"time"

	"google.golang.org/protobuf/proto"

	"im.turms/server/internal/domain/common/constant"
	"im.turms/server/internal/domain/gateway/session"
	"im.turms/server/internal/infra/exception"
	"im.turms/server/pkg/protocol"
)

var ErrCorruptedFrame = errors.New("corrupted frame")
var HeartbeatFailureRequestId int64 = -100

// Placeholder interfaces and structs pending migration
// Service definitions

type BlocklistService interface {
	TryBlockUserIdForCorruptedRequest(userId int64)
	TryBlockIpForCorruptedRequest(ip string)
	TryBlockIpForFrequentRequest(ip string)
	TryBlockUserIdForFrequentRequest(userId int64)

	IsIpBlocked(ip string) bool
	TryBlockIpForCorruptedFrame(ip string)
	TryBlockUserIdForCorruptedFrame(userId int64)
}

type SessionClientController interface {
	HandleCreateSessionRequest(wrapper *UserSessionWrapper, req *protocol.CreateSessionRequest) (*RequestHandlerResult, error)
	HandleDeleteSessionRequest(wrapper *UserSessionWrapper) (*RequestHandlerResult, error)
}

type SessionService interface {
	HandleHeartbeatUpdateRequest(session *session.UserSession)
	GetLocalUserSession(ip string) []*session.UserSession
}

type ServiceRequest struct {
	Ip         []byte
	UserId     int64
	DeviceType protocol.DeviceType
	RequestId  int64
	Type       interface{}
	Buffer     []byte
}

type ServiceRequestService interface {
	HandleServiceRequest(session *session.UserSession, req *ServiceRequest) (*protocol.TurmsNotification, error)
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

	NotificationFactory *NotificationFactory
}

// HandleRequest wraps handleRequest0 with error handling.
// @MappedFrom handleRequest(UserSessionWrapper sessionWrapper, ByteBuf serviceRequestBuffer)
func (d *ClientRequestDispatcher) HandleRequest(ctx context.Context, sessionWrapper *UserSessionWrapper, serviceRequestBuffer []byte) (buff []byte, err error) {
	defer func() {
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
			notif := d.NotificationFactory.CreateWithReason(&HeartbeatFailureRequestId, constant.ResponseStatusCode_SERVER_UNAVAILABLE, availability.Reason)
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
	if err != nil {
		if sessionWrapper.HasUserSession() {
			d.BlocklistService.TryBlockUserIdForCorruptedRequest(sessionWrapper.UserSession.UserID)
		}
		d.BlocklistService.TryBlockIpForCorruptedRequest(sessionWrapper.GetIPStr())
		if req.RequestId != nil {
			requestID = *req.RequestId
		}

		exc := exception.NewTurmsError(int32(constant.ResponseStatusCode_INVALID_REQUEST), err.Error())
		notification = d.NotificationFactory.CreateFromError(exc, &requestID)
	} else {
		if req.RequestId != nil {
			requestID = *req.RequestId
		}
		requestType = req.GetKind()

		canLogRequest := true
		if sessionWrapper.HasUserSession() {
			// if !sessionWrapper.UserSession.HasPermission(requestType) {
			//	 notification = d.NotificationFactory.Create(&requestID, constant.ResponseStatusCode_UNAUTHORIZED_REQUEST)
			// }
			if _, ok := req.Kind.(*protocol.TurmsRequest_DeleteSessionRequest); ok { // MOCK: userSession.AcquireDeleteSessionRequestLoggingLock()
				canLogRequest = true
			}
		}

		if notification == nil {
			notification, err = d.HandleServiceRequest(ctx, sessionWrapper, req, serviceRequestBuffer)
			if err != nil {
				notification = d.NotificationFactory.CreateFromError(err, &requestID)
			}
		}

		finalCanLogRequest := canLogRequest
		isServerError := constant.IsServerError(notification.GetCode())

		if isServerError || (d.ApiLoggingContext.ShouldLogRequest(requestType) && finalCanLogRequest) {
			var version *int32
			var userId *int64
			var sessionId *int32
			var deviceType *protocol.DeviceType

			if sessionWrapper.HasUserSession() {
				// MOCK parameters for userSession
				userId = &sessionWrapper.UserSession.UserID
				devType := sessionWrapper.UserSession.DeviceType
				deviceType = &devType
			}

			processingTimeMilli := (time.Now().UnixNano() - startTime) / 1000000
			d.ApiLoggingContext.LogRequest(sessionId, userId, deviceType, version, sessionWrapper.GetIPStr(), requestID, requestType, requestSize, requestTime, notification, processingTimeMilli)
		}
	}

	return proto.Marshal(notification)
}

// HandleServiceRequest delegates typed requests to backend services.
// @MappedFrom handleServiceRequest(UserSessionWrapper sessionWrapper, SimpleTurmsRequest request, ByteBuf serviceRequestBuffer, TracingContext tracingContext)
func (d *ClientRequestDispatcher) HandleServiceRequest(ctx context.Context, sessionWrapper *UserSessionWrapper, request *protocol.TurmsRequest, serviceRequestBuffer []byte) (*protocol.TurmsNotification, error) {
	requestID := request.GetRequestId()
	if requestID <= 0 {
		return d.NotificationFactory.CreateWithReason(&requestID, constant.ResponseStatusCode_INVALID_REQUEST, "The request ID must be greater than 0"), nil
	}

	availability := d.ServerStatusManager.GetServiceAvailability()
	if !availability.Available {
		return d.NotificationFactory.CreateWithReason(&requestID, constant.ResponseStatusCode_SERVER_UNAVAILABLE, availability.Reason), nil
	}

	// Rate limiting
	if !d.IpRequestThrottler.TryAcquireToken(sessionWrapper.GetIPStr()) {
		d.BlocklistService.TryBlockIpForFrequentRequest(sessionWrapper.GetIPStr())
		if sessionWrapper.HasUserSession() {
			d.BlocklistService.TryBlockUserIdForFrequentRequest(sessionWrapper.UserSession.UserID)
		}
		return d.NotificationFactory.Create(&requestID, constant.ResponseStatusCode_CLIENT_REQUESTS_TOO_FREQUENT), nil
	}

	switch kind := request.Kind.(type) {
	case *protocol.TurmsRequest_CreateSessionRequest:
		result, err := d.SessionController.HandleCreateSessionRequest(sessionWrapper, kind.CreateSessionRequest)
		if err != nil {
			return nil, err
		}
		return d.getNotificationFromHandlerResult(result, requestID), nil
	case *protocol.TurmsRequest_DeleteSessionRequest:
		result, err := d.SessionController.HandleDeleteSessionRequest(sessionWrapper)
		if err != nil {
			return nil, err
		}
		return d.getNotificationFromHandlerResult(result, requestID), nil
	default:
		return d.handleGenericServiceRequest(sessionWrapper, request, serviceRequestBuffer)
	}
}

func (d *ClientRequestDispatcher) handleGenericServiceRequest(sessionWrapper *UserSessionWrapper, request *protocol.TurmsRequest, serviceRequestBuffer []byte) (*protocol.TurmsNotification, error) {
	if !sessionWrapper.HasUserSession() || sessionWrapper.UserSession.Conn == nil {
		reqID := request.GetRequestId()
		return d.NotificationFactory.SessionClosed(&reqID), nil
	}

	session := sessionWrapper.UserSession
	svcReq := &ServiceRequest{
		Ip:         []byte(sessionWrapper.GetIPStr()),
		UserId:     session.UserID,
		DeviceType: session.DeviceType,
		RequestId:  request.GetRequestId(),
		Type:       request.GetKind(),
		Buffer:     serviceRequestBuffer,
	}
	return d.ServiceRequestService.HandleServiceRequest(session, svcReq)
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
		}

		var notification *protocol.TurmsNotification
		if !isSuccess {
			notification = d.NotificationFactory.Create(&HeartbeatFailureRequestId, constant.ResponseStatusCode_UPDATE_HEARTBEAT_OF_NONEXISTENT_SESSION)
		}

		d.ApiLoggingContext.LogRequest(sessionId, userId, deviceType, version, sessionWrapper.GetIPStr(), 0, "HEARTBEAT", 0, time.Now().UnixMilli(), notification, 0)
	}
	return data, nil
}

func (d *ClientRequestDispatcher) getNotificationFromHandlerResult(result *RequestHandlerResult, reqId int64) *protocol.TurmsNotification {
	notif := &protocol.TurmsNotification{
		Timestamp: time.Now().UnixMilli(),
		RequestId: &reqId,
		Code:      proto.Int32(int32(result.Code)),
	}
	if result.Reason != "" {
		notif.Reason = proto.String(result.Reason)
	}
	return notif
}
