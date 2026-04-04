package session

import (
	"context"
	"fmt"
	"net"
	"sync/atomic"
	"time"

	"google.golang.org/protobuf/proto"
	"im.turms/server/internal/domain/common/constant"
	sessionbo "im.turms/server/internal/domain/gateway/session/bo"
	"im.turms/server/pkg/protocol"
)

// MessageHandler defines the function signature for handling incoming client messages.
type MessageHandler func(ctx context.Context, userSession *UserSession, message []byte)

// UserSession encapsulates the network connection and the user state.
type UserSession struct {
	ID         int64
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

// GetLastHeartbeatRequestTimestamp returns the last heartbeat timestamp in ms.
// @MappedFrom getLastHeartbeatRequestTimestamp()
func (s *UserSession) GetLastHeartbeatRequestTimestamp() int64 {
	return atomic.LoadInt64(&s.lastHeartbeat)
}

// SetLastRequestTimestampToNow updates the last request time to now.
func (s *UserSession) SetLastRequestTimestampToNow() {
	now := time.Now()
	atomic.StoreInt64(&s.lastRequest, now.UnixMilli())
	atomic.StoreInt64(&s.lastRequestNanos, now.UnixNano())
}

// GetLastRequestTimestamp returns the last request timestamp in ms.
func (s *UserSession) GetLastRequestTimestamp() int64 {
	return atomic.LoadInt64(&s.lastRequest)
}

// IsSessionOpen returns true if the session is still open.
func (s *UserSession) IsSessionOpen() bool {
	return atomic.LoadUint32(&s.isSessionOpen) == 1
}

// Close marks the session as closed.
func (s *UserSession) Close() {
	atomic.StoreUint32(&s.isSessionOpen, 0)
}

func (s *UserSession) AcquireDeleteSessionLock() bool {
	return atomic.CompareAndSwapUint32(&s.isDeleteSessionLockAcquired, 0, 1)
}

// HasPermission checks if the user has permission for a specific request type.
func (s *UserSession) HasPermission(requestType any) bool {
	if s.Permissions == nil {
		// In Go refactor, a nil permissions map signifies TurmsRequestTypePool.ALL
		return true
	}
	return s.Permissions[requestType]
}

// AcquireDeleteSessionRequestLoggingLock acquires a lock ensuring that delete session requests
// are only logged once.
func (s *UserSession) AcquireDeleteSessionRequestLoggingLock() bool {
	return atomic.CompareAndSwapUint32(&s.isDeleteSessionLockAcquired, 0, 1)
}

// SetPermissions updates the user's permissions.
func (s *UserSession) SetPermissions(permissions map[any]bool) {
	s.Permissions = permissions
}

// SetConnection associates a network connection and IP with the user session.
func (s *UserSession) SetConnection(conn Connection, ip net.IP) {
	s.Conn = conn
	s.IP = ip
}

// SendMessage sends a protobuf notification to the user.
func (s *UserSession) SendMessage(notification *protocol.TurmsNotification) error {
	if !s.IsSessionOpen() {
		return fmt.Errorf("session is closed")
	}
	data, err := proto.Marshal(notification)
	if err != nil {
		return err
	}
	return s.Conn.Send(data)
}

func (s *UserSession) String() string {
	return fmt.Sprintf("UserSession{ID: %d, UserID: %d, DeviceType: %v, IP: %v}", s.ID, s.UserID, s.DeviceType, s.IP)
}

// Connection interface represents the underlying network transport.
type Connection interface {
	Connect() error
	Close(reason constant.SessionCloseStatus) error
	Send(data []byte) error
	RemoteAddr() net.Addr
	TryNotifyClientToRecover()
	IsActive() bool
}
