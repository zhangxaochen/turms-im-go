package controller

import (
	"context"
	"net"

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
func (c *SessionClientController) HandleDeleteSessionRequest(ctx context.Context, sessionWrapper *common.UserSessionWrapper) (*common.RequestHandlerResult, error) {
	if !sessionWrapper.HasUserSession() {
		return common.NewRequestHandlerResult(constant.ResponseStatusCode_OK, ""), nil
	}
	userSession := sessionWrapper.UserSession
	_, err := c.sessionService.CloseLocalSession(ctx, userSession.UserID, []protocol.DeviceType{userSession.DeviceType}, constant.SessionCloseStatus_DISCONNECTED_BY_CLIENT)
	if err != nil {
		// TODO: LOGGER.error("Caught an error while closing the session with the user ID: " + userId, t)
	}
	return common.NewRequestHandlerResult(constant.ResponseStatusCode_OK, ""), nil
}

// HandleCreateSessionRequest handles a client's request to create/login a session.
// @MappedFrom handleCreateSessionRequest(UserSessionWrapper sessionWrapper, CreateSessionRequest createSessionRequest)
func (c *SessionClientController) HandleCreateSessionRequest(ctx context.Context, sessionWrapper *common.UserSessionWrapper, req *protocol.CreateSessionRequest) (*common.RequestHandlerResult, error) {
	if sessionWrapper.HasUserSession() {
		return common.NewRequestHandlerResult(constant.ResponseStatusCode_CREATE_EXISTING_SESSION, ""), nil
	}

	userID := req.UserId
	var password string
	if req.Password != nil {
		password = *req.Password
	}

	var userStatus protocol.UserStatus
	if req.UserStatus != nil {
		userStatus = *req.UserStatus
	}

	deviceType := req.DeviceType
	if deviceType == protocol.DeviceType_UNKNOWN {
		deviceType = protocol.DeviceType_DESKTOP // Default or throw as per Java
	}

	deviceDetails := req.DeviceDetails
	var location *protocol.UserLocation
	if req.Location != nil {
		location = req.Location
	}

	session, err := c.sessionService.HandleLoginRequest(
		ctx,
		int(req.Version),
		net.ParseIP(sessionWrapper.GetIPStr()),
		userID,
		password,
		deviceType,
		deviceDetails,
		userStatus,
		location,
		sessionWrapper.GetIPStr(),
	)
	if err != nil {
		return nil, err
	}

	// Ensure the network connection is still open before cementing the session
	isConnectionAlive := false
	if conn := sessionWrapper.GetConnection(); conn != nil {
		isConnectionAlive = conn.IsActive()
	}

	if isConnectionAlive {
		if conn := sessionWrapper.GetConnection(); conn != nil {
			session.SetConnection(conn, net.ParseIP(sessionWrapper.GetIPStr()))
		}
		sessionWrapper.SetUserSession(session)

		userSessionsManager := c.sessionService.GetUserSessionsManager(ctx, userID)
		if userSessionsManager != nil {
			// Fire session established hooks
			c.sessionService.OnSessionEstablished(ctx, userSessionsManager, session.DeviceType)

			// Invoke online handlers (plugins)
			c.sessionService.InvokeGoOnlineHandlers(ctx, userSessionsManager, session)
		}

		return common.NewRequestHandlerResult(constant.ResponseStatusCode_OK, ""), nil
	}

	// If the connection dropped during the process, clean up
	c.sessionService.CloseLocalSession(ctx, userID, []protocol.DeviceType{deviceType}, constant.SessionCloseStatus_LOGIN_TIMEOUT)
	return common.NewRequestHandlerResult(constant.ResponseStatusCode_LOGIN_TIMEOUT, ""), nil
}
