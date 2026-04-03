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
// @MappedFrom countOnlineUsers(boolean countByNodes)
// @MappedFrom countOnlineUsers()
// @MappedFrom countLocalOnlineUsers()
func (s *SessionService) CountOnlineUsers() int {
	return s.shardedMap.CountOnlineUsers()
}

// @MappedFrom destroy()
func (s *SessionService) Destroy(ctx context.Context) error {
	return nil
}

// @MappedFrom handleHeartbeatUpdateRequest(UserSession session)
func (s *SessionService) HandleHeartbeatUpdateRequest(session *UserSession) {
}

// @MappedFrom handleLoginRequest(int version, @NotNull ByteArrayWrapper ip, @NotNull Long userId, @Nullable String password, @NotNull DeviceType deviceType, @Nullable Map<String, String> deviceDetails, @Nullable UserStatus userStatus, @Nullable Location location, @Nullable String ipStr)
func (s *SessionService) HandleLoginRequest(ctx context.Context, version int, ip []byte, userId int64, password string, deviceType protocol.DeviceType, deviceDetails map[string]string, userStatus protocol.UserStatus, location any, ipStr string) (*UserSession, error) {
	return nil, nil
}

// @MappedFrom closeLocalSessions(@NotNull List<byte[]> ips, @NotNull CloseReason closeReason)
// @MappedFrom closeLocalSessions(@NotNull byte[] ip, @NotNull CloseReason closeReason)
func (s *SessionService) CloseLocalSessionsByIp(ctx context.Context, ips [][]byte, closeReason any) error {
	return nil
}

// @MappedFrom closeLocalSession(@NotNull Long userId, @NotNull @ValidDeviceType DeviceType deviceType, @NotNull SessionCloseStatus closeStatus)
// @MappedFrom closeLocalSession(@NotNull Long userId, @NotNull @ValidDeviceType DeviceType deviceType, @NotNull CloseReason closeReason)
// @MappedFrom closeLocalSession(@NotNull Long userId, @NotEmpty Set<@ValidDeviceType DeviceType> deviceTypes, @NotNull CloseReason closeReason)
// @MappedFrom closeLocalSession(Long userId, SessionCloseStatus closeStatus)
// @MappedFrom closeLocalSession(Long userId, CloseReason closeReason)
func (s *SessionService) CloseLocalSession(ctx context.Context, userId int64, deviceTypes []protocol.DeviceType, closeReason any) error {
	return nil
}

// @MappedFrom closeLocalSessions(@NotNull Set<Long> userIds, @NotNull CloseReason closeReason)
func (s *SessionService) CloseLocalSessionsByUserIds(ctx context.Context, userIds []int64, closeReason any) error {
	return nil
}

// @MappedFrom authAndCloseLocalSession(@NotNull Long userId, @NotNull DeviceType deviceType, @NotNull CloseReason closeReason, int sessionId)
func (s *SessionService) AuthAndCloseLocalSession(ctx context.Context, userId int64, deviceType protocol.DeviceType, closeReason any, sessionId int) error {
	return nil
}

// @MappedFrom closeAllLocalSessions(@NotNull CloseReason closeReason)
func (s *SessionService) CloseAllLocalSessions(ctx context.Context, closeReason any) error {
	return nil
}

// @MappedFrom getSessions(Set<Long> userIds)
func (s *SessionService) GetSessions(ctx context.Context, userIds []int64) []any {
	return nil
}

// @MappedFrom authAndUpdateHeartbeatTimestamp(long userId, @NotNull @ValidDeviceType DeviceType deviceType, int sessionId)
func (s *SessionService) AuthAndUpdateHeartbeatTimestamp(ctx context.Context, userId int64, deviceType protocol.DeviceType, sessionId int) *UserSession {
	return nil
}

// @MappedFrom tryRegisterOnlineUser(int version, @NotNull Set<TurmsRequest.KindCase> permissions, @NotNull ByteArrayWrapper ip, @NotNull Long userId, @NotNull DeviceType deviceType, @Nullable Map<String, String> deviceDetails, @Nullable UserStatus userStatus, @Nullable Location location)
func (s *SessionService) TryRegisterOnlineUser(ctx context.Context, version int, permissions any, ip []byte, userId int64, deviceType protocol.DeviceType, deviceDetails map[string]string, userStatus protocol.UserStatus, location any) (*UserSession, error) {
	return nil, nil
}

// @MappedFrom getUserSessionsManager(@NotNull Long userId)
func (s *SessionService) GetUserSessionsManager(ctx context.Context, userId int64) any {
	return nil
}

// @MappedFrom getLocalUserSession(@NotNull Long userId, @NotNull DeviceType deviceType)
// @MappedFrom getLocalUserSession(ByteArrayWrapper ip)
func (s *SessionService) GetLocalUserSession(ctx context.Context, userId int64, deviceType protocol.DeviceType) *UserSession {
	return nil
}

// @MappedFrom onSessionEstablished(@NotNull UserSessionsManager userSessionsManager, @NotNull @ValidDeviceType DeviceType deviceType)
func (s *SessionService) OnSessionEstablished(ctx context.Context, userSessionsManager any, deviceType protocol.DeviceType) {
}

// @MappedFrom addOnSessionClosedListeners(Consumer<UserSession> onSessionClosed)
func (s *SessionService) AddOnSessionClosedListeners(ctx context.Context, onSessionClosed func(*UserSession)) {
}

// @MappedFrom invokeGoOnlineHandlers(@NotNull UserSessionsManager userSessionsManager, @NotNull UserSession userSession)
func (s *SessionService) InvokeGoOnlineHandlers(ctx context.Context, userSessionsManager any, userSession *UserSession) {
}
