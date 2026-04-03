package common

import (
	"context"
	"errors"

	"im.turms/server/internal/domain/gateway/session"
	"im.turms/server/pkg/protocol"
)

var ErrCorruptedFrame = errors.New("corrupted frame")

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
	UserId     int64
	DeviceType protocol.DeviceType
	RequestId  int64
	Type       int
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
	ShouldLogRequest(requestType int) bool
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
func (d *ClientRequestDispatcher) HandleRequest(ctx context.Context, sessionWrapper *UserSessionWrapper, serviceRequestBuffer []byte) ([]byte, error) {
	// Basic wrapper implementation
	return d.HandleRequest0(ctx, sessionWrapper, serviceRequestBuffer)
}

// HandleRequest0 processes the parsed request.
// @MappedFrom handleRequest0(UserSessionWrapper sessionWrapper, ByteBuf serviceRequestBuffer)
func (d *ClientRequestDispatcher) HandleRequest0(ctx context.Context, sessionWrapper *UserSessionWrapper, serviceRequestBuffer []byte) ([]byte, error) {
	// Minimal mock implementation pending full port
	if len(serviceRequestBuffer) == 0 {
		return nil, nil // Heartbeat
	}
	return nil, nil
}

// HandleServiceRequest delegates typed requests to backend services.
// @MappedFrom handleServiceRequest(UserSessionWrapper sessionWrapper, SimpleTurmsRequest request, ByteBuf serviceRequestBuffer, TracingContext tracingContext)
func (d *ClientRequestDispatcher) HandleServiceRequest(ctx context.Context, sessionWrapper *UserSessionWrapper, request *protocol.TurmsRequest, serviceRequestBuffer []byte) (*protocol.TurmsNotification, error) {
	// Minimal mock implementation
	return nil, nil
}
