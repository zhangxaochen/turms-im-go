package session

import (
	"context"
	"errors"
	"net"
	"sync"
	"time"

	sessionbo "im.turms/server/internal/domain/gateway/session/bo"
	userbo "im.turms/server/internal/domain/user/bo"
	"im.turms/server/internal/domain/user/service/onlineuser"
	"im.turms/server/pkg/protocol"
	"im.turms/server/internal/domain/common/infra/cluster/rpc"
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

type SessionService struct {
	shardedMap *ShardedUserSessionsMap

	mu sync.RWMutex

	ConflictStrategy MultiDeviceStrategy

	userStatusService            onlineuser.UserStatusService
	sessionLocationService       onlineuser.SessionLocationService
	userSimultaneousLoginService *UserSimultaneousLoginService
	sessionAuthenticationManager *SessionIdentityAccessManager
	nodeID                       string
	ipToSessions                 sync.Map // ipStr -> *sync.Map (*UserSession -> struct{})
	onSessionClosedListeners     []func(*UserSession)
	rpcService                   *rpc.RpcService
}

func NewSessionService(
	userStatusService onlineuser.UserStatusService,
	sessionLocationService onlineuser.SessionLocationService,
	userSimultaneousLoginService *UserSimultaneousLoginService,
	sessionAuthenticationManager *SessionIdentityAccessManager,
	nodeID string,
	rpcService *rpc.RpcService,
) *SessionService {
	return &SessionService{
		shardedMap:                   NewShardedUserSessionsMap(256),
		ConflictStrategy:             KickExisting,
		userStatusService:            userStatusService,
		sessionLocationService:       sessionLocationService,
		userSimultaneousLoginService: userSimultaneousLoginService,
		sessionAuthenticationManager: sessionAuthenticationManager,
		nodeID:                       nodeID,
		rpcService:                   rpcService,
	}
}

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
			if existingSession.Conn != nil {
				existingSession.Conn.Close()
			}
			if existingSession.CloseChan != nil {
				close(existingSession.CloseChan)
			}
			delete(manager.Sessions, session.DeviceType)
		}
	}

	manager.Sessions[session.DeviceType] = session
	s.registerSessionIp(session)
	return nil
}

func (s *SessionService) registerSessionIp(session *UserSession) {
	ipStr := session.IP.String()
	var sessionMap *sync.Map
	v, ok := s.ipToSessions.Load(ipStr)
	if ok {
		sessionMap = v.(*sync.Map)
	} else {
		sessionMap = &sync.Map{}
		v, _ = s.ipToSessions.LoadOrStore(ipStr, sessionMap)
		sessionMap = v.(*sync.Map)
	}
	sessionMap.Store(session, struct{}{})
}

func (s *SessionService) unregisterSessionIp(session *UserSession) {
	ipStr := session.IP.String()
	v, ok := s.ipToSessions.Load(ipStr)
	if ok {
		sessionMap := v.(*sync.Map)
		sessionMap.Delete(session)
	}
}

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

	if existing.Conn == conn {
		delete(manager.Sessions, deviceType)

		if conn != nil {
			_ = conn.Close()
		}
		if existing.CloseChan != nil {
			close(existing.CloseChan)
		}
		s.unregisterSessionIp(existing)

		s.mu.RLock()
		listeners := s.onSessionClosedListeners
		s.mu.RUnlock()
		for _, listener := range listeners {
			listener(existing)
		}
	}
	manager.mu.Unlock()

	if manager := s.shardedMap.RemoveIfEmpty(userID); manager != nil {
		s.InvokeGoOfflineHandlers(context.Background(), manager, nil)
	}
}

func (s *SessionService) GetUserSession(userID int64, deviceType protocol.DeviceType) (*UserSession, bool) {
	manager, ok := s.shardedMap.Get(userID)
	if !ok {
		return nil, false
	}

	session := manager.GetSession(deviceType)
	return session, session != nil
}

func (s *SessionService) GetAllUserSessions(userID int64) []*UserSession {
	manager, ok := s.shardedMap.Get(userID)
	if !ok {
		return nil
	}
	return manager.GetAllSessions()
}

func (s *SessionService) CountOnlineUsers() int {
	return s.shardedMap.CountOnlineUsers()
}

func (s *SessionService) Destroy(ctx context.Context) error {
	s.CloseAllLocalSessions(ctx, nil)
	return nil
}

func (s *SessionService) HandleHeartbeatUpdateRequest(session *UserSession) {
	session.SetLastHeartbeatRequestTimestampToNow()
}

func (s *SessionService) HandleLoginRequest(ctx context.Context, version int, ip []byte, userId int64, password string, deviceType protocol.DeviceType, deviceDetails map[string]string, userStatus protocol.UserStatus, location any, ipStr string) (*UserSession, error) {
	if version != 1 {
		return nil, errors.New("unsupported client version")
	}
	if s.userSimultaneousLoginService.IsForbiddenDeviceType(deviceType) {
		return nil, errors.New("login from forbidden device type")
	}

	var passwordPtr *string
	if password != "" {
		passwordPtr = &password
	}
	loginInfo := sessionbo.NewUserLoginInfo(version, userId, passwordPtr, deviceType, deviceDetails, &userStatus, ipStr)
	permissionInfo, err := s.sessionAuthenticationManager.VerifyAndGrant(ctx, loginInfo)
	if err != nil {
		return nil, err
	}

	return s.TryRegisterOnlineUser(ctx, version, permissionInfo.Permissions, ip, userId, deviceType, deviceDetails, userStatus, location)
}

func (s *SessionService) CloseLocalSessionsByIp(ctx context.Context, ips [][]byte, closeReason any) error {
	for _, ip := range ips {
		ipStr := net.IP(ip).String()
		if v, ok := s.ipToSessions.Load(ipStr); ok {
			sessionMap := v.(*sync.Map)
			sessionMap.Range(func(key, value any) bool {
				session := key.(*UserSession)
				s.UnregisterSession(session.UserID, session.DeviceType, session.Conn)
				return true
			})
		}
	}
	return nil
}

func (s *SessionService) CloseLocalSession(ctx context.Context, userId int64, deviceTypes []protocol.DeviceType, closeReason any) error {
	manager, ok := s.shardedMap.Get(userId)
	if !ok {
		return nil
	}

	manager.mu.Lock()
	defer manager.mu.Unlock()

	for _, dt := range deviceTypes {
		if session, exists := manager.Sessions[dt]; exists {
			s.unregisterSessionIp(session)
			delete(manager.Sessions, dt)

			if session.Conn != nil {
				_ = session.Conn.Close()
			}
			if session.CloseChan != nil {
				close(session.CloseChan)
			}

			s.mu.RLock()
			listeners := s.onSessionClosedListeners
			s.mu.RUnlock()
			for _, listener := range listeners {
				listener(session)
			}
		}
	}
	
	if manager := s.shardedMap.RemoveIfEmpty(userId); manager != nil {
		s.InvokeGoOfflineHandlers(ctx, manager, closeReason)
	}
	return nil
}

func (s *SessionService) CloseLocalSessionsByUserIds(ctx context.Context, userIds []int64, closeReason any) error {
	for _, uid := range userIds {
		manager, ok := s.shardedMap.Get(uid)
		if ok {
			sessions := manager.GetAllSessions()
			for _, session := range sessions {
				s.UnregisterSession(uid, session.DeviceType, session.Conn)
			}
		}
	}
	return nil
}

func (s *SessionService) AuthAndCloseLocalSession(ctx context.Context, userId int64, deviceType protocol.DeviceType, closeReason any, sessionId int) error {
	manager, ok := s.shardedMap.Get(userId)
	if !ok {
		return nil
	}

	manager.mu.Lock()
	defer manager.mu.Unlock()

	if session, exists := manager.Sessions[deviceType]; exists {
		if session.ID == sessionId {
			s.unregisterSessionIp(session)
			delete(manager.Sessions, deviceType)

			if session.Conn != nil {
				_ = session.Conn.Close()
			}
			if session.CloseChan != nil {
				close(session.CloseChan)
			}

			s.mu.RLock()
			listeners := s.onSessionClosedListeners
			s.mu.RUnlock()
			for _, listener := range listeners {
				listener(session)
			}
		}
	}
	
	if manager := s.shardedMap.RemoveIfEmpty(userId); manager != nil {
		s.InvokeGoOfflineHandlers(ctx, manager, closeReason)
	}
	return nil
}

func (s *SessionService) CloseAllLocalSessions(ctx context.Context, closeReason any) error {
	s.shardedMap.Range(func(userId int64, manager *UserSessionsManager) bool {
		var toClose []protocol.DeviceType
		for _, sess := range manager.GetAllSessions() {
			toClose = append(toClose, sess.DeviceType)
		}
		if len(toClose) > 0 {
			_ = s.CloseLocalSession(ctx, userId, toClose, closeReason)
		}
		return true
	})
	return nil
}

func (s *SessionService) GetSessions(ctx context.Context, userIds []int64) []*sessionbo.UserSessionsInfo {
	if len(userIds) == 0 {
		return nil
	}

	var result []*sessionbo.UserSessionsInfo
	for _, uid := range userIds {
		manager, ok := s.shardedMap.Get(uid)
		if !ok {
			continue
		}

		sessions := manager.GetAllSessions()
		if len(sessions) == 0 {
			continue
		}

		var sessionInfos []sessionbo.UserSessionInfo
		for _, sess := range sessions {
			var loc *sessionbo.UserLocation
			if sess.Location != nil {
				loc = &sessionbo.UserLocation{
					Longitude: sess.Location.Longitude,
					Latitude:  sess.Location.Latitude,
				}
			}
			sessionInfos = append(sessionInfos, sessionbo.UserSessionInfo{
				ID:         sess.ID,
				Version:    sess.Version,
				DeviceType: sess.DeviceType,
				LoginDate:  sess.LoginDate.UnixMilli(),
				Location:   loc,
			})
		}

		result = append(result, &sessionbo.UserSessionsInfo{
			UserID:   uid,
			Sessions: sessionInfos,
		})
	}
	return result
}

func (s *SessionService) AuthAndUpdateHeartbeatTimestamp(ctx context.Context, userId int64, deviceType protocol.DeviceType, sessionId int) *UserSession {
	if session, ok := s.GetUserSession(userId, deviceType); ok {
		if session.ID == sessionId {
			s.HandleHeartbeatUpdateRequest(session)
			return session
		}
	}
	return nil
}

func (s *SessionService) TryRegisterOnlineUser(ctx context.Context, version int, permissions any, ip []byte, userId int64, deviceType protocol.DeviceType, deviceDetails map[string]string, userStatus protocol.UserStatus, location any) (*UserSession, error) {
	sessionsStatus, err := s.userStatusService.FetchUserSessionsStatus(ctx, userId)
	if err != nil {
		return nil, err
	}

	manager, ok := s.shardedMap.Get(userId)
	if ok {
		var toClose []protocol.DeviceType
		for _, sess := range manager.GetAllSessions() {
			if info, exists := sessionsStatus.OnlineDeviceTypeToSessionInfo[sess.DeviceType]; exists {
				if info.IsActive && info.NodeID != s.nodeID {
					toClose = append(toClose, sess.DeviceType)
				}
			}
		}
		if len(toClose) > 0 {
			_ = s.CloseLocalSession(ctx, userId, toClose, "disconnected_by_other_device")
		}
	}

	if sessionsStatus.UserStatus == protocol.UserStatus_OFFLINE {
		return s.addOnlineDeviceIfAbsent(ctx, version, permissions, ip, userId, deviceType, deviceDetails, userStatus, location, "", 0)
	}

	info, exists := sessionsStatus.OnlineDeviceTypeToSessionInfo[deviceType]
	if exists && info.IsActive {
		session, ok := s.GetUserSession(userId, deviceType)
		isClosedSessionOnLocal := ok && session.Conn != nil && !session.Conn.IsActive()
		if isClosedSessionOnLocal {
			if userStatus != 0 && sessionsStatus.UserStatus != userStatus {
				_, _ = s.userStatusService.UpdateStatus(ctx, userId, userStatus)
			}
			if location != nil {
				loc, ok := location.(*protocol.UserLocation)
				if ok && s.sessionLocationService != nil {
					_ = s.sessionLocationService.UpsertUserLocation(ctx, userId, deviceType, loc.Longitude, loc.Latitude)
				}
			}
			return session, nil
		} else if s.userSimultaneousLoginService.ShouldDisconnectLoggingInDeviceIfConflicts() {
			return nil, ErrSessionAlreadyExists
		}
	}

	var expectedNodeId string
	var expectedDeviceTimestamp int64
	if exists {
		expectedNodeId = info.NodeID
		expectedDeviceTimestamp = info.HeartbeatTimestampSeconds
	}

	wasSuccessful, err := s.closeSessionsWithConflictedDeviceTypes(ctx, userId, deviceType, sessionsStatus)
	if err != nil || !wasSuccessful {
		return nil, ErrSessionAlreadyExists
	}

	return s.addOnlineDeviceIfAbsent(ctx, version, permissions, ip, userId, deviceType, deviceDetails, userStatus, location, expectedNodeId, expectedDeviceTimestamp)
}

func (s *SessionService) closeSessionsWithConflictedDeviceTypes(ctx context.Context, userId int64, deviceType protocol.DeviceType, sessionsStatus *userbo.UserSessionsStatus) (bool, error) {
	if sessionsStatus == nil || len(sessionsStatus.OnlineDeviceTypeToSessionInfo) == 0 {
		return true, nil
	}

	for existingDeviceType, info := range sessionsStatus.OnlineDeviceTypeToSessionInfo {
		if !info.IsActive {
			continue
		}
		// Based on MultiDeviceStrategy, if we are KickExisting and there's a conflict
		if s.userSimultaneousLoginService.IsConflicted(deviceType, existingDeviceType) {
			if info.NodeID == s.nodeID {
				// Handled locally beforehand
				// SetLocalSessionOfflineByUserIdAndDeviceExists here if full logic was present.
				continue
			}
			
			// Send RPC to `info.NodeID` to kick `existingDeviceType` for `userId`
			// Use the new RpcService skeleton.
			// TODO: Actually implement SetSessionOfflineRequest in /dto
			if s.rpcService != nil {
				_, _ = s.rpcService.RequestResponse(ctx, info.NodeID, nil) // passing nil request struct temporarily
			}
		}
	}
	return true, nil
}

func (s *SessionService) addOnlineDeviceIfAbsent(ctx context.Context, version int, permissions any, ip []byte, userId int64, deviceType protocol.DeviceType, deviceDetails map[string]string, userStatus protocol.UserStatus, location any, expectedNodeId string, expectedDeviceTimestamp int64) (*UserSession, error) {
	now := time.Now()
	if s.userStatusService != nil {
		ok, err := s.userStatusService.AddOnlineDevice(ctx, userId, deviceType, userStatus, s.nodeID, &now)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, ErrSessionAlreadyExists
		}
	}

	session := &UserSession{
		UserID:     userId,
		DeviceType: deviceType,
		IP:         net.IP(ip),
		LoginDate:  now,
		CloseChan:  make(chan struct{}),
	}
	if location != nil {
		loc, ok := location.(*protocol.UserLocation)
		if ok {
			session.Location = &sessionbo.UserLocation{
				Longitude: loc.Longitude,
				Latitude:  loc.Latitude,
			}
		}
	}
	session.SetLastHeartbeatRequestTimestampToNow()

	err := s.RegisterSession(ctx, session)
	if err != nil {
		return nil, err
	}

	if location != nil {
		loc, ok := location.(*protocol.UserLocation)
		if ok && s.sessionLocationService != nil {
			_ = s.sessionLocationService.UpsertUserLocation(ctx, userId, deviceType, loc.Longitude, loc.Latitude)
		}
	}

	s.InvokeGoOnlineHandlers(ctx, nil, session)

	return session, nil
}

func (s *SessionService) GetUserSessionsManager(ctx context.Context, userId int64) any {
	manager, ok := s.shardedMap.Get(userId)
	if !ok {
		return nil
	}
	return manager
}

func (s *SessionService) GetLocalUserSession(ctx context.Context, userId int64, deviceType protocol.DeviceType) *UserSession {
	session, _ := s.GetUserSession(userId, deviceType)
	return session
}

func (s *SessionService) GetLocalUserSessionsByIp(ctx context.Context, ip []byte) []*UserSession {
	ipStr := net.IP(ip).String()
	var sessions []*UserSession
	if v, ok := s.ipToSessions.Load(ipStr); ok {
		sessionMap := v.(*sync.Map)
		sessionMap.Range(func(key, value any) bool {
			sessions = append(sessions, key.(*UserSession))
			return true
		})
	}
	return sessions
}

func (s *SessionService) OnSessionEstablished(ctx context.Context, userSessionsManager any, deviceType protocol.DeviceType) {
	// TODO: Increment metrics (e.g. LoggedInUsersCounter, OnlineUsersGauge)
	// TODO: Notify clients of session info if properties.Gateway.Session.NotifyClientsOfSessionInfoAfterConnected is true
}

func (s *SessionService) AddOnSessionClosedListeners(ctx context.Context, onSessionClosed func(*UserSession)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.onSessionClosedListeners = append(s.onSessionClosedListeners, onSessionClosed)
}

func (s *SessionService) InvokeGoOnlineHandlers(ctx context.Context, userSessionsManager any, userSession *UserSession) {
	// TODO: 插件系统尚未实现: 调用 PluginManager (UserOnlineStatusChangeHandler.goOnline)
}

func (s *SessionService) InvokeGoOfflineHandlers(ctx context.Context, userSessionsManager *UserSessionsManager, closeReason any) {
	// TODO: 插件系统尚未实现: 调用 PluginManager (UserOnlineStatusChangeHandler.goOffline)
}
