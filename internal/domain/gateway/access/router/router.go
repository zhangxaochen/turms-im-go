package router

import (
	"context"
	"net"
	"reflect"

	"google.golang.org/protobuf/proto"

	"im.turms/server/internal/domain/gateway/access/client/common"
	"im.turms/server/internal/domain/gateway/session"
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
	return &Router{
		sessionService:      sessionService,
		controllers:         make(ControllerMap),
		throttler:           common.DefaultIpRequestThrottler(),
		availabilityHandler: common.NewServiceAvailabilityHandler(),
		notifFactory:        common.NewNotificationFactory(),
	}
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
	if s.Conn != nil && s.Conn.RemoteAddr() != nil {
		if tcpAddr, ok := s.Conn.RemoteAddr().(*net.TCPAddr); ok {
			ipStr = tcpAddr.IP.String()
		}
	}
	if !r.throttler.TryAcquireToken(ipStr) {
		r.sendNotification(s, nil, 1400, "Too many requests") // e.g. Turms 1400/429
		return
	}

	req := &protocol.TurmsRequest{}
	if err := proto.Unmarshal(payload, req); err != nil {
		_ = s.Conn.WriteMessage([]byte("MALFORMED_PROTOBUF"))
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
		r.sessionService.UnregisterSession(s.UserID, s.DeviceType, s.Conn)
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
		_ = s.Conn.WriteMessage(respBytes)
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

func (r *Router) sendNotification(s *session.UserSession, requestID *int64, code int32, reason string) {
	buf, err := r.notifFactory.CreateBuffer(requestID, code, reason)
	if err == nil {
		_ = s.Conn.WriteMessage(buf)
	}
}
