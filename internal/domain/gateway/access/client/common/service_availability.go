package common

import (
	"sync/atomic"
)

// ServerStatus represents the lifecycle state of the Gateway server.
type ServerStatus int32

const (
	StatusStarting ServerStatus = iota
	StatusRunning
	StatusShuttingDown
)

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
// Typically used by TCP/WS connection interceptors to immediately reject connections
// if the server is shutting down or not yet ready.
func (sa *ServiceAvailabilityHandler) IsAvailable() bool {
	return sa.GetStatus() == StatusRunning
}

// @MappedFrom channelRegistered(ChannelHandlerContext ctx)
func (sa *ServiceAvailabilityHandler) ChannelRegistered(isAvailable bool) bool {
	return isAvailable
}

// @MappedFrom exceptionCaught(ChannelHandlerContext ctx, Throwable cause)
func (sa *ServiceAvailabilityHandler) ExceptionCaught(err error) {
	// Log or handle the connection exception
}
