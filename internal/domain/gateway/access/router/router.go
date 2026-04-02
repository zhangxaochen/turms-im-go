package router

import (
	"context"
	"reflect"

	"google.golang.org/protobuf/proto"

	"im.turms/server/internal/domain/gateway/session"
	"im.turms/server/pkg/protocol"
)

type ControllerMap map[reflect.Type]func(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error)

type Router struct {
	sessionService *session.SessionService
	controllers    ControllerMap
}

func NewRouter(sessionService *session.SessionService) *Router {
	return &Router{
		sessionService: sessionService,
		controllers:    make(ControllerMap),
	}
}

// RegisterController takes a sample of the kind (e.g. &protocol.TurmsRequest_CreateMessageRequest{})
func (r *Router) RegisterController(kindSample interface{}, handler func(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error)) {
	r.controllers[reflect.TypeOf(kindSample)] = handler
}

func (r *Router) HandleMessage(ctx context.Context, s *session.UserSession, payload []byte) {
	req := &protocol.TurmsRequest{}
	if err := proto.Unmarshal(payload, req); err != nil {
		// Log or close connection on malformed protobuf
		_ = s.Conn.WriteMessage([]byte("MALFORMED_PROTOBUF"))
		return
	}

	// 1. Session Login Verification
	if req.GetCreateSessionRequest() != nil {
		r.handleCreateSession(ctx, s, req)
		return
	}

	if s.UserID == 0 {
		// Unauthed
		r.sendNotification(s, req.RequestId, 1200, "Unauthorized") // Using Turms code 1200 commonly
		return
	}

	// 2. Dispatch to controllers
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
	resp := &protocol.TurmsNotification{
		RequestId: requestID,
		Code:      proto.Int32(code),
		Reason:    proto.String(reason),
	}
	respBytes, _ := proto.Marshal(resp)
	_ = s.Conn.WriteMessage(respBytes)
}
