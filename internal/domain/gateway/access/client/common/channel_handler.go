package common

import (
	"errors"
	"net"
)

var ErrOutOfMemory = errors.New("out of direct memory / fatal memory error")

// @MappedFrom ServiceAvailabilityHandler
type ServiceAvailabilityChannelHandler struct {
	blocklistService    BlocklistService
	serverStatusManager ServerStatusManager
	sessionService      SessionService
}

func NewServiceAvailabilityChannelHandler(
	blocklistService BlocklistService,
	serverStatusManager ServerStatusManager,
	sessionService SessionService,
) *ServiceAvailabilityChannelHandler {
	return &ServiceAvailabilityChannelHandler{
		blocklistService:    blocklistService,
		serverStatusManager: serverStatusManager,
		sessionService:      sessionService,
	}
}

// @MappedFrom channelRegistered(ChannelHandlerContext ctx)
// HandleConnection answers whether the connection should be kept open based on service availability and blocklists.
func (h *ServiceAvailabilityChannelHandler) HandleConnection(conn net.Conn) bool {
	serviceAvailability := h.serverStatusManager.GetServiceAvailability()
	available := serviceAvailability.Available

	if available {
		addr := conn.RemoteAddr()
		if tcpAddr, ok := addr.(*net.TCPAddr); ok {
			ipStr := tcpAddr.IP.String()
			if h.blocklistService.IsIpBlocked(ipStr) {
				return false
			}
		}
		return true
	}
	return false
}

// @MappedFrom exceptionCaught(ChannelHandlerContext ctx, Throwable cause)
// HandleException processes connection errors, mapping to Netty's exceptionCaught.
func (h *ServiceAvailabilityChannelHandler) HandleException(conn net.Conn, cause error) error {
	// If corrupted frame, block IP and session users
	if errors.Is(cause, ErrCorruptedFrame) {
		addr := conn.RemoteAddr()
		// Replicate Java's implicit NullPointerException / ClassCastException behavior fail-fast
		tcpAddr := addr.(*net.TCPAddr)
		ipStr := tcpAddr.IP.String()
		h.blocklistService.TryBlockIpForCorruptedFrame(ipStr)

		sessions := h.sessionService.GetLocalUserSession(ipStr)
		for _, s := range sessions {
			h.blocklistService.TryBlockUserIdForCorruptedFrame(s.UserID)
		}
	} else if errors.Is(cause, ErrOutOfMemory) {
		conn.Close()
	}
	
	// Unconditionally propagate the exception upstream
	return cause
}
