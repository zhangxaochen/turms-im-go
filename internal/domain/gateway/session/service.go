package session

import (
	"context"
	"errors"
	"sync"

	"im.turms/server/pkg/protocol"
)

var (
	ErrSessionAlreadyExists = errors.New("cluster session already exists but conflicting action denied")
)

type MultiDeviceStrategy int

const (
	KickExisting MultiDeviceStrategy = iota
	DenyNew
	AllowMultiple
)

// SessionService manages the active connections and handles the session lifecycles for the gateway.
type SessionService struct {
	shardedMap *ShardedUserSessionsMap

	mu sync.RWMutex

	// How the system handles multiple logins from the same user on the same/different device types.
	// We'll hardcode to KickExisting (Turms default for identical device types) for now
	ConflictStrategy MultiDeviceStrategy
}

func NewSessionService() *SessionService {
	// Let's assume 256 shards is sufficient for a moderately high traffic gateway node
	return &SessionService{
		shardedMap:       NewShardedUserSessionsMap(256),
		ConflictStrategy: KickExisting,
	}
}

// RegisterSession adds a new session connected to this gateway.
// Here we perform collision detection and apply kick logic.
func (s *SessionService) RegisterSession(ctx context.Context, session *UserSession) error {
	userID := session.UserID

	manager := s.shardedMap.GetOrAdd(userID)

	manager.mu.Lock()
	defer manager.mu.Unlock()

	existingSession := manager.Sessions[session.DeviceType]
	if existingSession != nil {
		if s.ConflictStrategy == DenyNew {
			return ErrSessionAlreadyExists
		} else if s.ConflictStrategy == KickExisting {
			// Disconnect existing
			existingSession.Conn.Close()
			close(existingSession.CloseChan)
			delete(manager.Sessions, session.DeviceType)
		}
	}

	// Wait! Even across different devices, there might be constraints,
	// but default IMs usually allow Desktop + Mobile.

	manager.Sessions[session.DeviceType] = session
	return nil
}

// UnregisterSession removes a session from internal management.
func (s *SessionService) UnregisterSession(userID int64, deviceType protocol.DeviceType, conn Connection) {
	manager, ok := s.shardedMap.Get(userID)
	if !ok {
		return
	}

	manager.mu.Lock()

	existing, exists := manager.Sessions[deviceType]
	if !exists {
		manager.mu.Unlock()
		return
	}

	// Ensure we only remove if it's the exact same connection object (prevent removing replaced sessions)
	if existing.Conn == conn {
		delete(manager.Sessions, deviceType)

		// Close connection if not closed yet
		_ = conn.Close()
		close(existing.CloseChan)
	}
	manager.mu.Unlock()

	// We can lazily garbage collect the manager here
	s.shardedMap.RemoveIfEmpty(userID)
}

// GetUserSession fetches a specific user's session by device type.
func (s *SessionService) GetUserSession(userID int64, deviceType protocol.DeviceType) (*UserSession, bool) {
	manager, ok := s.shardedMap.Get(userID)
	if !ok {
		return nil, false
	}

	session := manager.GetSession(deviceType)
	return session, session != nil
}

// GetAllUserSessions fetches all active devices for a user.
func (s *SessionService) GetAllUserSessions(userID int64) []*UserSession {
	manager, ok := s.shardedMap.Get(userID)
	if !ok {
		return nil
	}
	return manager.GetAllSessions()
}

// CountOnlineUsers returns the approximate number of active connections in this gateway node.
func (s *SessionService) CountOnlineUsers() int {
	return s.shardedMap.CountOnlineUsers()
}
