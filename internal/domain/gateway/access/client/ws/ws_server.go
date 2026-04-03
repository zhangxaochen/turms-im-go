package ws

import (
	"im.turms/server/internal/domain/gateway/session"
)

// @MappedFrom HttpForwardedHeaderHandler
type HttpForwardedHeaderHandler struct{}

// @MappedFrom apply(ConnectionInfo connectionInfo, HttpRequest request)
func (h *HttpForwardedHeaderHandler) Apply(connectionInfo any, request any) any {
	return nil
}

// @MappedFrom WebSocketServerFactory
type WebSocketServerFactory struct{}

// @MappedFrom create(WebSocketProperties webSocketProperties, BlocklistService blocklistService, ServerStatusManager serverStatusManager, SessionService sessionService, ConnectionListener connectionListener, int maxFramePayloadLength)
func (f *WebSocketServerFactory) Create(webSocketProperties any, blocklistService any, serverStatusManager any, sessionService *session.SessionService, connectionListener any, maxFramePayloadLength int) {
}
