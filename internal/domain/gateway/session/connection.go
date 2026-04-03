package session

import (
	"context"
	"net"
	"sync/atomic"
	"time"

	"im.turms/server/pkg/protocol"
)

// UserSession encapsulates the network connection and the user state.
type UserSession struct {
	ID         int
	Version    int
	UserID     int64
	DeviceType protocol.DeviceType

	// IP address of the client (for tracking/rate-limiting)
	IP net.IP

	LoginDate time.Time

	// Holds the actual transport wrapper (TCP/WebSocket)
	Conn Connection

	Permissions                 map[any]bool
	isDeleteSessionLockAcquired uint32

	lastHeartbeat int64 // Use atomic operations if accessed concurrently
	lastRequest   int64 // Use atomic operations if accessed concurrently

	// Channel to cleanly shutdown the connection loops
	CloseChan chan struct{}
}

// SetLastHeartbeatRequestTimestampToNow updates the heartbeat to now.
// @MappedFrom setLastHeartbeatRequestTimestampToNow()
func (s *UserSession) SetLastHeartbeatRequestTimestampToNow() {
	atomic.StoreInt64(&s.lastHeartbeat, time.Now().UnixMilli())
}

// GetLastHeartbeatRequestTimestamp returns the last heartbeat received in milliseconds.
func (s *UserSession) GetLastHeartbeatRequestTimestamp() int64 {
	return atomic.LoadInt64(&s.lastHeartbeat)
}

// SetLastRequestTimestampToNow updates the last request to now.
// @MappedFrom setLastRequestTimestampToNow()
func (s *UserSession) SetLastRequestTimestampToNow() {
	atomic.StoreInt64(&s.lastRequest, time.Now().UnixMilli())
}

// GetLastRequestTimestamp returns the last request timestamp in milliseconds.
func (s *UserSession) GetLastRequestTimestamp() int64 {
	return atomic.LoadInt64(&s.lastRequest)
}

// IsOpen returns whether the session's connection is active
// @MappedFrom isOpen()
func (s *UserSession) IsOpen() bool {
	return s.Conn != nil
}

// Connection abstracts away the difference between net.TCPConn and gorilla websocket.Conn
// It provides a unified way to send data to the client.
type Connection interface {
	// Send raw Protobuf notification to the client. The implementation handles
	// wrapping it in Varint length prefix (if TCP) or raw binary frame (if WS).
	WriteMessage(payload []byte) error

	// Close cleanly terminates the socket
	Close() error

	// RemoteAddr returns the client's network address
	RemoteAddr() net.Addr

	// TryNotifyClientToRecover stops the connection recovery status
	TryNotifyClientToRecover()

	// IsActive returns true if the connection is active
	IsActive() bool
}

// Handler defines how incoming raw payloads (post-frame decoding) are processed.
// This allows the gateway server to pass the routing logic into the listener loops.
type MessageHandler func(ctx context.Context, session *UserSession, payload []byte)

// @MappedFrom setConnection(NetConnection connection, ByteArrayWrapper ip)
func (s *UserSession) SetConnection(connection Connection, ip string) {
	s.Conn = connection
}

// @MappedFrom isConnected()
func (s *UserSession) IsConnected() bool {
	return s.Conn != nil
}

// @MappedFrom supportsSwitchingToUdp()
func (s *UserSession) SupportsSwitchingToUdp() bool {
	return s.DeviceType != protocol.DeviceType_BROWSER
}

// @MappedFrom toString()
func (s *UserSession) ToString() string {
	return "UserSession{}"
}

// @MappedFrom acquireDeleteSessionRequestLoggingLock()
func (s *UserSession) AcquireDeleteSessionRequestLoggingLock() bool {
	return atomic.CompareAndSwapUint32(&s.isDeleteSessionLockAcquired, 0, 1)
}

// @MappedFrom hasPermission(TurmsRequest.KindCase requestType)
func (s *UserSession) HasPermission(requestType any) bool {
	if s.Permissions == nil {
		return false
	}
	return s.Permissions[requestType]
}

// @MappedFrom close(@NotNull CloseReason closeReason)
func (s *UserSession) Close(closeReason any) {
	if s.Conn != nil {
		_ = s.Conn.Close()
	}
}
