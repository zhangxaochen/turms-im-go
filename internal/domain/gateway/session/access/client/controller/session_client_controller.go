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

	// TODO: Java uses sessionService.closeLocalSession(userId, deviceType, SessionCloseStatus.DISCONNECTED_BY_CLIENT)
	// Currently SessionService methods aren't fully ported to match exactly.
	// We call UnregisterSession directly as a temporary stub.
	c.sessionService.UnregisterSession(userSession.UserID, userSession.DeviceType, userSession.Conn)
	return nil
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
	// If the value is somehow unrecognized, use a default/nil. Go protoc doesn't gen UNRECOGNIZED constants explicitly if skipping, but let's assume valid.
	if req.UserStatus != nil {
		userStatus = *req.UserStatus
	}

	deviceType := req.DeviceType
	// Go doesn't inherently have UNRECOGNIZED on generated pb unless using protoc-gen-go specific output.
	// The protoc generated enum unknown case usually maps to whatever default is or UNKNOWN.
	if deviceType == protocol.DeviceType(5) { // Assuming 5 is UNKNOWN
		deviceType = protocol.DeviceType_UNKNOWN
	}

	// TODO: Log deviceDetails in API logs (ClientApiLogging)
	deviceDetails := req.DeviceDetails

	// TODO: Map CreateSessionRequest.Location to the internal Location BO
	// location := req.Location

	// TODO: sessionService.handleLoginRequest needs a real implementation. Currently returning mock user session.
	session, err := c.sessionService.HandleLoginRequest(
		ctx,
		int(req.Version),
		[]byte(sessionWrapper.GetIP()),
		userID,
		password,
		deviceType,
		deviceDetails,
		userStatus,
		nil, // TODO: Map location correctly
		sessionWrapper.GetIPStr(),
	)
	if err != nil {
		return nil, err
	}

	// The sessionEstablishTimeout task cancellation logic from Java
	// (sessionEstablishTimeout == null || sessionEstablishTimeout.cancel())
	// In Go, timeouts are usually managed via Context or custom connection layer timers.
	// TODO: Check if session establish timeout expired via wrapper or connection layer context.
	isTimeout := false
	if isTimeout {
		// TODO: c.sessionService.CloseLocalSession(ctx, userID, []protocol.DeviceType{deviceType}, constant.SessionCloseStatus_LOGIN_TIMEOUT)
		return common.NewRequestHandlerResult(constant.ResponseStatusCode_LOGIN_TIMEOUT, ""), nil
	}

	// Ensure the network connection is still open before cementing the session
	// TODO: Fetch actual connected status via sessionWrapper.Connection().IsConnected()
	isConnectionAlive := true
	if isConnectionAlive {
		sessionWrapper.SetUserSession(session)

		// TODO: Complete `getUserSessionsManager` properly instead of `any` type constraint.
		userSessionsManager := c.sessionService.GetUserSessionsManager(ctx, userID)

		// Fire session established hooks
		c.sessionService.OnSessionEstablished(ctx, userSessionsManager, session.DeviceType)

		// Invoke online handlers (plugins)
		// TODO: Ensure we track and log the error via Logger.Errorf(ERROR_INVOKE_GO_ONLINE, err)
		c.sessionService.InvokeGoOnlineHandlers(ctx, userSessionsManager, session)

		return common.NewRequestHandlerResult(constant.ResponseStatusCode_OK, ""), nil
	}

	// If the connection dropped during the process, clean up
	// TODO: c.sessionService.CloseLocalSession(ctx, userID, []protocol.DeviceType{deviceType}, constant.SessionCloseStatus_LOGIN_TIMEOUT)
	return nil, nil
}
