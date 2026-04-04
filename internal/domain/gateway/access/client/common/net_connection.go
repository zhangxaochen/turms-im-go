package common

import (
	"context"
	"net"
	"sync"

	"im.turms/server/internal/domain/common/constant"
)

// CloseReason represents the reason for closing a connection.
type CloseReason struct {
	Status             constant.SessionCloseStatus
	BusinessStatusCode constant.ResponseStatusCode
	Reason             string
	IsNotifyClient      bool
}

func NewCloseReason(status constant.SessionCloseStatus) CloseReason {
	// SessionCloseStatus.isNotifyClient() logic from Java
	isNotify := false
	if status == constant.SessionCloseStatus_DISCONNECTED_BY_OTHER_DEVICE ||
		status == constant.SessionCloseStatus_DISCONNECTED_BY_ADMIN ||
		status == constant.SessionCloseStatus_DISCONNECTED_BY_SERVER ||
		status == constant.SessionCloseStatus_SWITCH ||
		status == constant.SessionCloseStatus_HEARTBEAT_TIMEOUT {
		isNotify = true
	}
	return CloseReason{
		Status:         status,
		IsNotifyClient: isNotify,
	}
}

func CloseReasonFromError(err error) CloseReason {
	if err == nil {
		return NewCloseReason(constant.SessionCloseStatus_UNKNOWN_ERROR)
	}

	fromErr, ok := err.(interface {
		Code() constant.ResponseStatusCode
		Reason() string
	})

	if ok {
		code := fromErr.Code()
		status := constant.SessionCloseStatus_UNKNOWN_ERROR

		// Map some status codes to close status, simple for now
		if code >= constant.ResponseStatusCode_SERVER_INTERNAL_ERROR && code < 1300 {
			if code == constant.ResponseStatusCode_SERVER_UNAVAILABLE {
				status = constant.SessionCloseStatus_SERVER_UNAVAILABLE
			} else {
				status = constant.SessionCloseStatus_SERVER_ERROR
			}
		} else if code == constant.ResponseStatusCode_ILLEGAL_ARGUMENT || code == constant.ResponseStatusCode_INVALID_REQUEST {
			status = constant.SessionCloseStatus_ILLEGAL_REQUEST
		}

		return CloseReason{
			Status:             status,
			BusinessStatusCode: code,
			Reason:             fromErr.Reason(),
		}
	}

	return CloseReason{
		Status: constant.SessionCloseStatus_UNKNOWN_ERROR,
		Reason: err.Error(),
	}
}

// @MappedFrom NetConnection
type NetConnection interface {
	GetAddress() net.Addr
	Send(ctx context.Context, buffer []byte) error
	CloseWithReason(reason CloseReason) bool
	Close() error
	IsConnected() bool
	IsSwitchingToUdp() bool
	IsConnectionRecovering() bool
	SwitchToUdp()
	TryNotifyClientToRecover()
}

// BaseNetConnection provides common state for NetConnection implementations.
type BaseNetConnection struct {
	udpAddress             *net.UDPAddr
	isConnected            bool
	isSwitchingToUdp       bool
	isConnectionRecovering bool
	isDisposed             bool
	mu                     sync.RWMutex
	udpSignalDispatcher    func(*net.UDPAddr) // Injectable callback to notify via UDP
}

func NewBaseNetConnection(connected bool) *BaseNetConnection {
	return &BaseNetConnection{
		isConnected: connected,
	}
}

// SetUdpSignalDispatcher sets the callback to send recovery signals
func (b *BaseNetConnection) SetUdpSignalDispatcher(dispatcher func(*net.UDPAddr)) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.udpSignalDispatcher = dispatcher
}

// @MappedFrom close(CloseReason closeReason)
func (b *BaseNetConnection) CloseWithReason(reason CloseReason) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.isDisposed || !b.isConnected {
		return false
	}
	b.isConnected = false
	b.isConnectionRecovering = false
	b.isSwitchingToUdp = reason.Status == constant.SessionCloseStatus_SWITCH
	return true
}

// @MappedFrom close()
func (b *BaseNetConnection) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.isDisposed = true
	b.isConnected = false
	b.isConnectionRecovering = false
	b.isSwitchingToUdp = false
	return nil
}

// @MappedFrom tryNotifyClientToRecover()
func (b *BaseNetConnection) TryNotifyClientToRecover() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if !b.isConnected && !b.isConnectionRecovering && b.udpAddress != nil {
		if b.udpSignalDispatcher != nil {
			b.udpSignalDispatcher(b.udpAddress)
		}
		b.isConnectionRecovering = true
	}
}

// State accessors

func (b *BaseNetConnection) IsConnected() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.isConnected
}

func (b *BaseNetConnection) IsSwitchingToUdp() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.isSwitchingToUdp
}

func (b *BaseNetConnection) IsConnectionRecovering() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.isConnectionRecovering
}

func (b *BaseNetConnection) IsDisposed() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.isDisposed
}

func (b *BaseNetConnection) SetDisposed(disposed bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.isDisposed = disposed
}
