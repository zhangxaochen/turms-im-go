package session

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync/atomic"
	"time"

	sessionbo "im.turms/server/internal/domain/gateway/session/bo"
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
	Location  *sessionbo.UserLocation

	// Holds the actual transport wrapper (TCP/WebSocket)
	Conn Connection

	Permissions                 map[any]bool
	isDeleteSessionLockAcquired uint32

	lastHeartbeat      int64 // Use atomic operations if accessed concurrently
	lastHeartbeatNanos int64
	lastRequest        int64 // Use atomic operations if accessed concurrently
	lastRequestNanos   int64

	isSessionOpen uint32 // 1 for open, 0 for closed

	// Channel to cleanly shutdown the connection loops
	CloseChan chan struct{}
}

// SetLastHeartbeatRequestTimestampToNow updates the heartbeat to now.
// @MappedFrom setLastHeartbeatRequestTimestampToNow()
func (s *UserSession) SetLastHeartbeatRequestTimestampToNow() {
	now := time.Now()
	atomic.StoreInt64(&s.lastHeartbeat, now.UnixMilli())
	atomic.StoreInt64(&s.lastHeartbeatNanos, now.UnixNano())
}

// GetLastHeartbeatRequestTimestamp returns the last heartbeat received in milliseconds.
func (s *UserSession) GetLastHeartbeatRequestTimestamp() int64 {
	return atomic.LoadInt64(&s.lastHeartbeat)
}

// SetLastRequestTimestampToNow updates the last request to now.
// @MappedFrom setLastRequestTimestampToNow()
func (s *UserSession) SetLastRequestTimestampToNow() {
	now := time.Now()
	atomic.StoreInt64(&s.lastRequest, now.UnixMilli())
	atomic.StoreInt64(&s.lastRequestNanos, now.UnixNano())
}

// GetLastRequestTimestamp returns the last request timestamp in milliseconds.
func (s *UserSession) GetLastRequestTimestamp() int64 {
	return atomic.LoadInt64(&s.lastRequest)
}

// IsOpen returns whether the session's connection is active
// @MappedFrom isOpen()
func (s *UserSession) IsOpen() bool {
	return atomic.LoadUint32(&s.isSessionOpen) == 1
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
func (s *UserSession) SetConnection(connection Connection, ip net.IP) {
	s.Conn = connection
	s.IP = ip
}

// @MappedFrom isConnected()
func (s *UserSession) IsConnected() bool {
	return s.Conn != nil && s.Conn.IsActive()
}

// @MappedFrom supportsSwitchingToUdp()
func (s *UserSession) SupportsSwitchingToUdp() bool {
	return s.DeviceType != protocol.DeviceType_BROWSER
}

// @MappedFrom toString()
func (s *UserSession) ToString() string {
	return fmt.Sprintf("UserSession{id=%d, version=%d, userId=%d, deviceType=%v, loginDate=%v, loginLocation=%v, isSessionOpen=%v, connection=%v}",
		s.ID, s.Version, s.UserID, s.DeviceType, s.LoginDate, s.Location, s.IsOpen(), s.Conn)
}

// @MappedFrom acquireDeleteSessionRequestLoggingLock()
func (s *UserSession) AcquireDeleteSessionRequestLoggingLock() bool {
	return atomic.CompareAndSwapUint32(&s.isDeleteSessionLockAcquired, 0, 1)
}

// @MappedFrom hasPermission(TurmsRequest.KindCase requestType)
func (s *UserSession) HasPermission(requestType any) bool {
	if s.Permissions == nil {
		// In Go refactor, a nil permissions map signifies TurmsRequestTypePool.ALL
		return true
	}
	return s.Permissions[requestType]
}

// @MappedFrom close(@NotNull CloseReason closeReason)
func (s *UserSession) Close(closeReason any) bool {
	if atomic.CompareAndSwapUint32(&s.isSessionOpen, 1, 0) {
		if s.Conn != nil {
			_ = s.Conn.Close() // In the future, pass closeReason if Conn supports it
		} else {
			log.Printf("The connection is missing for the user session: %v", s.ID)
		}
		return true
	}
	return false
}

// SendNotification forwards a notification payload to the client.
// @MappedFrom sendNotification(ByteBuf byteBuf)
func (s *UserSession) SendNotification(payload []byte) error {
	if s.Conn != nil {
		return s.Conn.WriteMessage(payload)
	}
	return nil
}
