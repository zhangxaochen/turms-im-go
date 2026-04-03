package ws

import (
	"im.turms/server/internal/domain/gateway/access/server"
	"im.turms/server/internal/domain/gateway/session"
)

// @MappedFrom HttpForwardedHeaderHandler
type HttpForwardedHeaderHandler struct{}

// @MappedFrom apply(ConnectionInfo connectionInfo, HttpRequest request)
func (h *HttpForwardedHeaderHandler) Apply(connectionInfo any, request any) any {
	// Pending implementation for Forwarded / X-Forwarded-For parsing
	return nil
}

// @MappedFrom WebSocketServerFactory
type WebSocketServerFactory struct{}

// @MappedFrom create(...)
func (f *WebSocketServerFactory) Create(addr string, handler session.MessageHandler, sessionService *session.SessionService) *server.WSServer {
	// Note: Replaced multiple Spring/Reactor specifics with actual Go WSServer setup
	wsServer := server.NewWSServer(addr, sessionService, handler)
	_ = wsServer.Start()
	return wsServer
}
