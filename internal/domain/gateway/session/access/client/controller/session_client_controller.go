package controller

import (
	"context"

	"im.turms/server/internal/domain/common/constant"
	"im.turms/server/internal/domain/gateway/access/client/common"
	"im.turms/server/internal/domain/gateway/session"
	"im.turms/server/pkg/protocol"
)

// SessionClientController handles incoming client requests related to session management.
type SessionClientController struct {
	sessionService *session.SessionService
}

func NewSessionClientController(sessionService *session.SessionService) *SessionClientController {
	return &SessionClientController{
		sessionService: sessionService,
	}
}

// HandleDeleteSessionRequest handles a client's request to delete/logout their session.
// @MappedFrom handleDeleteSessionRequest(UserSessionWrapper sessionWrapper)
func (c *SessionClientController) HandleDeleteSessionRequest(ctx context.Context, sessionWrapper *common.UserSessionWrapper) error {
	userSession := sessionWrapper.UserSession
	if userSession == nil {
		return nil
	}

	// This is a stub calling UnregisterSession. In the future this should map to CloseLocalSession with the
	// specific SessionCloseStatus.DISCONNECTED_BY_CLIENT
	c.sessionService.UnregisterSession(userSession.UserID, userSession.DeviceType, userSession.Conn)
	return nil
}

// HandleCreateSessionRequest handles a client's request to create/login a session.
// @MappedFrom handleCreateSessionRequest(UserSessionWrapper sessionWrapper, CreateSessionRequest createSessionRequest)
func (c *SessionClientController) HandleCreateSessionRequest(ctx context.Context, sessionWrapper *common.UserSessionWrapper, req *protocol.CreateSessionRequest) (*common.RequestHandlerResult, error) {
	if sessionWrapper.HasUserSession() {
		return common.NewRequestHandlerResult(constant.ResponseStatusCode_CREATE_EXISTING_SESSION, ""), nil
	}

	// Stub mapping implementation
	// We'd parse the req parameters and call SessionService's HandleLoginRequest
	// For now we map the basic flow

	userID := req.UserId
	deviceType := req.DeviceType
	if deviceType == protocol.DeviceType_UNKNOWN {
		// fallbacks
	}

	// Simulate successful login
	newSession := &session.UserSession{
		UserID:     userID,
		DeviceType: deviceType,
		Conn:       nil, // To be tied to the wrapper's network connection
	}

	err := c.sessionService.RegisterSession(ctx, newSession)
	if err != nil {
		return nil, err
	}

	sessionWrapper.SetUserSession(newSession)

	// In the real system, it'd check connection state and timeout,
	// and invoke GoOnlineHandlers here

	return common.NewRequestHandlerResult(constant.ResponseStatusCode_OK, ""), nil
}
