package router

import (
	"context"
	"net"
	"reflect"
	"time"

	"google.golang.org/protobuf/proto"

	"im.turms/server/internal/domain/common/constant"
	"im.turms/server/internal/domain/gateway/access/client/common"
	"im.turms/server/internal/domain/gateway/config"
	"im.turms/server/internal/domain/gateway/session"
	sessionbo "im.turms/server/internal/domain/gateway/session/bo"
	"im.turms/server/pkg/protocol"
)

type ControllerMap map[reflect.Type]func(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error)

type Router struct {
	sessionService      *session.SessionService
	controllers         ControllerMap
	throttler           *common.IpRequestThrottler
	availabilityHandler *common.ServiceAvailabilityHandler
	notifFactory        *common.NotificationFactory
}

// NewRouter initializes the gateway request router with core components.
func NewRouter(sessionService *session.SessionService) *Router {
	r := &Router{
		sessionService:      sessionService,
		controllers:         make(ControllerMap),
		throttler:           common.DefaultIpRequestThrottler(),
		availabilityHandler: common.NewServiceAvailabilityHandler(),
		notifFactory:        common.NewNotificationFactory(config.NewGatewayProperties()),
	}

	sessionService.AddOnSessionClosedListeners(context.Background(), func(s *session.UserSession) {
		if s.IP != nil {
			r.throttler.CleanupByIp(s.IP.String())
		}
	})

	return r
}

// SetServiceAvailability changes the availability status (e.g., shutting down).
func (r *Router) SetServiceAvailability(status common.ServerStatus) {
	r.availabilityHandler.SetStatus(status)
}

// RegisterController takes a sample of the kind (e.g. &protocol.TurmsRequest_CreateMessageRequest{})
func (r *Router) RegisterController(kindSample interface{}, handler func(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error)) {
	r.controllers[reflect.TypeOf(kindSample)] = handler
}

func (r *Router) HandleMessage(ctx context.Context, s *session.UserSession, payload []byte) {
	// 0. Service Availability Check
	if !r.availabilityHandler.IsAvailable() {
		// Drop silently or return server error
		return
	}

	// 1. IP Throttling Check
	ipStr := "<unknown>"
	if s.Conn != nil && s.Conn.GetAddress() != nil {
		if tcpAddr, ok := s.Conn.GetAddress().(*net.TCPAddr); ok {
			ipStr = tcpAddr.IP.String()
		}
	}
	timestamp := time.Now().UnixNano()
	if !r.throttler.TryAcquireToken(ipStr, timestamp) {
		r.sendNotification(s, nil, int32(constant.ResponseStatusCode_CLIENT_REQUESTS_TOO_FREQUENT), "Too many requests")
		return
	}

	req := &protocol.TurmsRequest{}
	if err := proto.Unmarshal(payload, req); err != nil {
		if s.Conn != nil {
			_ = s.Conn.Send([]byte("MALFORMED_PROTOBUF"))
		}
		return
	}

	// Always update heartbeat and activity timestamp if we parse a valid payload
	s.SetLastHeartbeatRequestTimestampToNow()
	s.SetLastRequestTimestampToNow()

	// 2. Session Login Verification
	if req.GetCreateSessionRequest() != nil {
		r.handleCreateSession(ctx, s, req)
		return
	}

	// Explicit Session Close by client
	if req.GetDeleteSessionRequest() != nil {
		r.sessionService.UnregisterSession(ctx, s.UserID, s.DeviceType, s.Conn, sessionbo.NewCloseReason(constant.SessionCloseStatus_DISCONNECTED_BY_CLIENT))
		r.sendNotification(s, req.RequestId, 1000, "OK")
		return
	}

	if s.UserID == 0 {
		// Unauthed
		r.sendNotification(s, req.RequestId, 1200, "Unauthorized")
		return
	}

	// 3. Dispatch to controllers
	kind := reflect.TypeOf(req.Kind)
	handler, exists := r.controllers[kind]
	if !exists {
		r.sendNotification(s, req.RequestId, 1100, "Not implemented")
		return
	}

	resp, err := handler(ctx, s, req)
	if err != nil {
		r.sendNotification(s, req.RequestId, 1100, err.Error())
		return
	}

	if resp != nil {
		if resp.RequestId == nil && req.RequestId != nil {
			resp.RequestId = req.RequestId
		}
		respBytes, _ := proto.Marshal(resp)
		if s.Conn != nil {
			_ = s.Conn.Send(respBytes)
		}
	}
}

func (r *Router) handleCreateSession(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) {
	loginReq := req.GetCreateSessionRequest()
	s.UserID = loginReq.UserId
	s.DeviceType = loginReq.DeviceType

	err := r.sessionService.RegisterSession(ctx, s)
	if err != nil {
		r.sendNotification(s, req.RequestId, 1100, err.Error())
		return
	}

	r.sendNotification(s, req.RequestId, 1000, "OK")
}

// @MappedFrom sendNotification(ByteBuf byteBuf, TracingContext tracingContext)
// @MappedFrom sendNotification(ByteBuf byteBuf)
func (r *Router) sendNotification(s *session.UserSession, requestID *int64, code int32, reason string) {
	buf, err := r.notifFactory.CreateBuffer(requestID, constant.ResponseStatusCode(code), reason)
	if err == nil && s.Conn != nil {
		_ = s.Conn.Send(buf)
	}
}
