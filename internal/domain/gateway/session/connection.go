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
	userbo "im.turms/server/internal/domain/user/bo"
	"im.turms/server/internal/infra/tracing"
	"im.turms/server/pkg/protocol"
)

// MessageHandler defines the function signature for handling incoming client messages.
type MessageHandler func(ctx context.Context, userSession *UserSession, message []byte)

type SessionLocationService interface {
	UpsertUserLocation(ctx context.Context, userID int64, deviceType protocol.DeviceType, longitude float32, latitude float32) error
	RemoveUserLocation(ctx context.Context, userID int64, deviceType protocol.DeviceType) error
	RemoveUserLocations(ctx context.Context, userID int64, deviceTypes []protocol.DeviceType) error
	GetUserLocation(ctx context.Context, userID int64, deviceType protocol.DeviceType) (*protocol.UserLocation, error)
}

// UserSession encapsulates the network connection and the user state.
type UserSession struct {
	ID            int64
	Version       int
	UserID        int64
	DeviceType    protocol.DeviceType
	DeviceDetails map[string]string

	IP net.IP

	LoginDate time.Time
	Location  *sessionbo.UserLocation

	Conn Connection

	Permissions                           map[any]bool
	isDeleteSessionRequestLoggingAcquired uint32

	lastHeartbeat      int64
	lastHeartbeatNanos int64
	lastRequest        int64
	lastRequestNanos   int64

	isSessionOpen uint32 // 1 for open, 0 for closed

	// Channel to cleanly shutdown the connection loops
	CloseChan chan struct{}
}

func (s *UserSession) SetLastHeartbeatRequestTimestampToNow() {
	now := time.Now()
	atomic.StoreInt64(&s.lastHeartbeat, now.UnixMilli())
	atomic.StoreInt64(&s.lastHeartbeatNanos, now.UnixNano())
}

func (s *UserSession) GetLastHeartbeatRequestTimestamp() int64 {
	return atomic.LoadInt64(&s.lastHeartbeat)
}

func (s *UserSession) SetLastRequestTimestampToNow() {
	now := time.Now()
	atomic.StoreInt64(&s.lastRequest, now.UnixMilli())
	atomic.StoreInt64(&s.lastRequestNanos, now.UnixNano())
}

func (s *UserSession) GetLastRequestTimestamp() int64 {
	return atomic.LoadInt64(&s.lastRequest)
}

func (s *UserSession) IsOpen() bool {
	return atomic.LoadUint32(&s.isSessionOpen) == 1
}

func (s *UserSession) Close(closeReason constant.SessionCloseStatus) bool {
	if atomic.CompareAndSwapUint32(&s.isSessionOpen, 1, 0) {
		if s.Conn != nil {
			_ = s.Conn.Close(closeReason)
		} else {
			fmt.Printf("WARN: The connection is missing for the user session: %v\n", s)
		}
		if s.CloseChan != nil {
			close(s.CloseChan)
		}
		return true
	}
	return false
}

func (s *UserSession) AcquireDeleteSessionRequestLoggingLock() bool {
	return atomic.CompareAndSwapUint32(&s.isDeleteSessionRequestLoggingAcquired, 0, 1)
}

func (s *UserSession) SetConnection(conn Connection, ip net.IP) {
	s.Conn = conn
	s.IP = ip
}

func (s *UserSession) GetIPStr() string {
	if s.IP == nil {
		return ""
	}
	return s.IP.String()
}

func (s *UserSession) HasPermission(requestType any) bool {
	if s.Permissions == nil {
		return false
	}
	return s.Permissions[requestType]
}

func (s *UserSession) SetPermissions(permissions map[any]bool) {
	s.Permissions = permissions
}

func (s *UserSession) SendMessage(notification *protocol.TurmsNotification, tracingContext ...*tracing.TracingContext) error {
	data, err := proto.Marshal(notification)
	if err != nil {
		return err
	}
	if s.Conn != nil {
		return s.Conn.Send(data)
	}
	fmt.Printf("WARN: The connection is missing for the user session: %v\n", s)
	return nil
}

func (s *UserSession) String() string {
	return fmt.Sprintf("UserSession{ID: %d, Version: %d, UserID: %d, DeviceType: %v, LoginDate: %v, Location: %v, IsSessionOpen: %v, IP: %v, Conn: %v}",
		s.ID, s.Version, s.UserID, s.DeviceType, s.LoginDate, s.Location, s.IsOpen(), s.IP, s.Conn != nil)
}

func (s *UserSession) IsConnected() bool {
	return s.Conn != nil && s.Conn.IsActive()
}

func (s *UserSession) SupportsSwitchingToUdp() bool {
	return s.DeviceType != protocol.DeviceType_BROWSER
}

type UserStatusService interface {
	AddOnlineDevice(ctx context.Context, userID int64, deviceType protocol.DeviceType, status protocol.UserStatus, nodeID string, heartbeatTimestamp *time.Time) (bool, error)
	RemoveOnlineDevice(ctx context.Context, userID int64, deviceType protocol.DeviceType, nodeID string) (bool, error)
	RemoveOnlineDevices(ctx context.Context, userID int64, deviceTypes []protocol.DeviceType, nodeID string) (bool, error)
	UpdateStatus(ctx context.Context, userID int64, status protocol.UserStatus) (bool, error)
	FetchUserSessionsStatus(ctx context.Context, userID int64) (*userbo.UserSessionsStatus, error)
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
