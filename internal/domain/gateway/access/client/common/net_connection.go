package common

import (
	"context"
	"net"
	"sync"

	"im.turms/server/internal/domain/common/constant"
)

// CloseReason represents the reason for closing a connection.
type CloseReason struct {
	Status constant.SessionCloseStatus
	Reason string
}

func NewCloseReason(status constant.SessionCloseStatus) CloseReason {
	return CloseReason{Status: status}
}

// @MappedFrom NetConnection
type NetConnection interface {
	GetAddress() net.Addr
	Send(ctx context.Context, buffer []byte) error
	CloseWithReason(reason CloseReason) error
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
	mu                     sync.RWMutex
}

func NewBaseNetConnection(connected bool) *BaseNetConnection {
	return &BaseNetConnection{
		isConnected: connected,
	}
}

// @MappedFrom close(CloseReason closeReason)
func (b *BaseNetConnection) CloseWithReason(reason CloseReason) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.isConnected = false
	b.isConnectionRecovering = false
	b.isSwitchingToUdp = reason.Status == constant.SessionCloseStatus_SWITCH
	return nil
}

// @MappedFrom close()
func (b *BaseNetConnection) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.isConnected = false
	b.isConnectionRecovering = false
	b.isSwitchingToUdp = false
	return nil
}

// @MappedFrom switchToUdp()
func (b *BaseNetConnection) SwitchToUdp() {
	b.CloseWithReason(NewCloseReason(constant.SessionCloseStatus_SERVER_CLOSED))
}

// @MappedFrom tryNotifyClientToRecover()
func (b *BaseNetConnection) TryNotifyClientToRecover() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if !b.isConnected && !b.isConnectionRecovering && b.udpAddress != nil {
		// Mock logic: UdpRequestDispatcher.instance.sendSignal
		// udpdispatcher.SendSignal(b.udpAddress, udp.OPEN_CONNECTION)
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
