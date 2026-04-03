package common

import (
	"errors"
	"log"
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
func (h *ServiceAvailabilityChannelHandler) HandleException(conn net.Conn, cause error) {
	// If corrupted frame, block IP and session users
	if errors.Is(cause, ErrCorruptedFrame) {
		addr := conn.RemoteAddr()
		if tcpAddr, ok := addr.(*net.TCPAddr); ok {
			ipStr := tcpAddr.IP.String()
			h.blocklistService.TryBlockIpForCorruptedFrame(ipStr)

			sessions := h.sessionService.GetLocalUserSession(ipStr)
			for _, s := range sessions {
				h.blocklistService.TryBlockUserIdForCorruptedFrame(s.UserID)
			}
		}
		// Java: ctx.fireExceptionCaught(...) followed by potential close upstream or here.
		// Corrupted frames should force close the connection because the stream is no longer valid.
		conn.Close()
	} else if errors.Is(cause, ErrOutOfMemory) {
		log.Printf("Fatal memory error caught on connection %v: %v", conn.RemoteAddr(), cause)
		conn.Close()
		// In Go, usually we'd allow panics to bubble up, but if we catch it, we must close.
	} else {
		// Log or handle the connection exception
		log.Printf("Connection exception caught on %v: %v", conn.RemoteAddr(), cause)
	}
}
