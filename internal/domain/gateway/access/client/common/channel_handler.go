package common

import (
	"errors"
	"net"
	"sync/atomic"
)

var ErrOutOfMemory = errors.New("out of direct memory / fatal memory error")
var ErrCorruptedFrame = errors.New("corrupted frame")

// ServerStatus represents the lifecycle state of the Gateway server.
type ServerStatus int32

const (
	StatusStarting ServerStatus = iota
	StatusRunning
	StatusShuttingDown
)

func (s ServerStatus) String() string {
	switch s {
	case StatusStarting:
		return "STARTING"
	case StatusRunning:
		return "RUNNING"
	case StatusShuttingDown:
		return "SHUTTING_DOWN"
	default:
		return "UNKNOWN"
	}
}

// ServiceAvailabilityHandler maintains the global health and availability state of the nodes.
type ServiceAvailabilityHandler struct {
	status atomic.Int32
}

func NewServiceAvailabilityHandler() *ServiceAvailabilityHandler {
	sa := &ServiceAvailabilityHandler{}
	sa.SetStatus(StatusStarting)
	return sa
}

func (sa *ServiceAvailabilityHandler) SetStatus(status ServerStatus) {
	sa.status.Store(int32(status))
}

func (sa *ServiceAvailabilityHandler) GetStatus() ServerStatus {
	return ServerStatus(sa.status.Load())
}

// IsAvailable returns true if the server is in the RUNNING state.
func (sa *ServiceAvailabilityHandler) IsAvailable() bool {
	return sa.GetStatus() == StatusRunning
}

// @MappedFrom ServiceAvailabilityChannelHandler
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
	// IP Blocking Check
	addr := conn.RemoteAddr()
	if tcpAddr, ok := addr.(*net.TCPAddr); ok {
		if h.blocklistService != nil && h.blocklistService.IsIpBlocked(tcpAddr.IP) {
			return false
		}
	}

	// Server Status Check
	availability := h.serverStatusManager.GetServiceAvailability()
	if !availability.Available {
		return false
	}
	return true
}

// @MappedFrom exceptionCaught(ChannelHandlerContext ctx, Throwable cause)
// HandleException processes connection errors, mapping to Netty's exceptionCaught.
func (h *ServiceAvailabilityChannelHandler) HandleException(conn net.Conn, cause error) error {
	// If corrupted frame, block IP and session users
	if errors.Is(cause, ErrCorruptedFrame) {
		addr := conn.RemoteAddr()
		tcpAddr, ok := addr.(*net.TCPAddr)
		if !ok {
			return cause
		}
		ipBytes := tcpAddr.IP
		h.blocklistService.TryBlockIpForCorruptedFrame(ipBytes)

		sessions := h.sessionService.GetLocalUserSessionsByIp(ipBytes)
		for _, s := range sessions {
			h.blocklistService.TryBlockUserIdForCorruptedFrame(s.UserID)
		}
	} else if errors.Is(cause, ErrOutOfMemory) {
		_ = conn.Close()
	}

	// Unconditionally propagate the exception upstream
	return cause
}
