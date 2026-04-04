package session

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"im.turms/server/internal/domain/common/constant"
	"im.turms/server/internal/domain/common/infra/cluster/rpc"
	sessionbo "im.turms/server/internal/domain/gateway/session/bo"
	userbo "im.turms/server/internal/domain/user/bo"
	"im.turms/server/internal/domain/user/service/onlineuser"
	"im.turms/server/pkg/protocol"
)

var (
	ErrSessionAlreadyExists = errors.New("cluster session already exists but conflicting action denied")
)

// SessionAuthError is returned when authentication fails with a specific status code.
// @MappedFrom Java: ResponseException.get(statusCode)
type SessionAuthError struct {
	Code constant.ResponseStatusCode
}

func (e *SessionAuthError) Error() string {
	return fmt.Sprintf("session authentication failed with code: %d", e.Code)
}

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

	heartbeatManager *HeartbeatManager
	sessionIDCounter atomic.Int64
}

func NewSessionService(
	userStatusService onlineuser.UserStatusService,
	sessionLocationService onlineuser.SessionLocationService,
	userSimultaneousLoginService *UserSimultaneousLoginService,
	sessionAuthenticationManager *SessionIdentityAccessManager,
	nodeID string,
	rpcService *rpc.RpcService,
) *SessionService {
	svc := &SessionService{
		shardedMap:                   NewShardedUserSessionsMap(256),
		ConflictStrategy:             KickExisting,
		userStatusService:            userStatusService,
		sessionLocationService:       sessionLocationService,
		userSimultaneousLoginService: userSimultaneousLoginService,
		sessionAuthenticationManager: sessionAuthenticationManager,
		nodeID:                       nodeID,
		rpcService:                   rpcService,
	}

	// Initialize session ID counter with a random positive value
	svc.sessionIDCounter.Store(int64(rand.Int31n(1000000) + 1))

	// Default heartbeat timeouts (should be loadable from config in Phase 4)
	svc.heartbeatManager = NewHeartbeatManager(svc, 10*time.Second, 40*time.Second)
	svc.heartbeatManager.Start()

	if rpcService != nil {
		rpcService.Router().Register(1, func(ctx context.Context, payload []byte) ([]byte, error) {
			var req rpc.SetUserOfflineRequest
			if err := json.Unmarshal(payload, &req); err != nil {
				return nil, err
			}
			svc.CloseLocalSession(ctx, req.UserID, req.DeviceTypes, constant.SessionCloseStatus(req.SessionCloseStatus))
			return json.Marshal(true)
		})
	}

	return svc
}

func (s *SessionService) nextSessionID() int64 {
	return s.sessionIDCounter.Add(1)
}

func (s *SessionService) RegisterSession(ctx context.Context, session *UserSession) error {
	userID := session.UserID

	manager := s.shardedMap.GetOrAdd(userID, protocol.UserStatus_AVAILABLE)

	manager.mu.Lock()
	defer manager.mu.Unlock()

	existingSession := manager.Sessions[session.DeviceType]
	if existingSession != nil {
		if s.ConflictStrategy == DenyNew {
			return ErrSessionAlreadyExists
		} else if s.ConflictStrategy == KickExisting {
			if existingSession.Conn != nil {
				_ = existingSession.Conn.Close(constant.SessionCloseStatus_DISCONNECTED_BY_OTHER_DEVICE)
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

func (s *SessionService) UnregisterSession(ctx context.Context, userID int64, deviceType protocol.DeviceType, conn Connection, closeReason constant.SessionCloseStatus) {
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
			_ = conn.Close(closeReason)
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
		s.InvokeGoOfflineHandlers(ctx, manager, closeReason)
	}

	// Remove from distributed status
	if s.userStatusService != nil {
		_, _ = s.userStatusService.RemoveOnlineDevice(ctx, userID, deviceType, s.nodeID)
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
	if s.heartbeatManager != nil {
		s.heartbeatManager.Stop()
	}
	s.CloseAllLocalSessions(ctx, constant.SessionCloseStatus_SERVER_CLOSED)
	return nil
}

func (s *SessionService) HandleHeartbeatUpdateRequest(session *UserSession) {
	session.SetLastHeartbeatRequestTimestampToNow()
}

func (s *SessionService) HandleLoginRequest(ctx context.Context, version int, ip []byte, userId int64, password string, deviceType protocol.DeviceType, deviceDetails map[string]string, userStatus protocol.UserStatus, location *protocol.UserLocation, ipStr string) (*UserSession, error) {
	if version != 1 {
		return nil, errors.New("unsupported client version")
	}
	if s.userSimultaneousLoginService.IsForbiddenDeviceType(deviceType) {
		return nil, &SessionAuthError{Code: constant.ResponseStatusCode_LOGIN_FROM_FORBIDDEN_DEVICE_TYPE}
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
	if permissionInfo.AuthenticationCode != constant.ResponseStatusCode_OK {
		return nil, &SessionAuthError{Code: permissionInfo.AuthenticationCode}
	}

	return s.TryRegisterOnlineUser(ctx, version, permissionInfo.Permissions, ip, userId, deviceType, deviceDetails, userStatus, location)
}

func (s *SessionService) TryRegisterOnlineUser(ctx context.Context, version int, permissions any, ip []byte, userId int64, deviceType protocol.DeviceType, deviceDetails map[string]string, userStatus protocol.UserStatus, location *protocol.UserLocation) (*UserSession, error) {
	// 1. Fetch user status and handle conflicts
	sessionsStatus, err := s.userStatusService.FetchUserSessionsStatus(ctx, userId)
	if err != nil {
		return nil, err
	}

	// 2. Resolve conflicts
	conflictedByDenyNew, err := s.resolveConflicts(ctx, userId, deviceType, sessionsStatus)
	if err != nil {
		return nil, err
	}
	if conflictedByDenyNew {
		return nil, &SessionAuthError{Code: constant.ResponseStatusCode_SESSION_SIMULTANEOUS_CONFLICTS_DECLINE}
	}

	// 3. Register in UserStatusService (Redis)
	now := time.Now()
	added, err := s.userStatusService.AddOnlineDevice(ctx, userId, deviceType, userStatus, s.nodeID, &now)
	if err != nil {
		return nil, err
	}
	if !added {
		// This handles the race condition if another node won
		return nil, &SessionAuthError{Code: constant.ResponseStatusCode_SESSION_SIMULTANEOUS_CONFLICTS_DECLINE}
	}

	// 4. Create and local register session
	session := &UserSession{
		ID:         s.nextSessionID(),
		Version:    version,
		UserID:     userId,
		DeviceType: deviceType,
		IP:         net.IP(ip),
		LoginDate:  now,
		CloseChan:  make(chan struct{}),
	}
	if location != nil {
		session.Location = &sessionbo.UserLocation{
			Longitude: location.Longitude,
			Latitude:  location.Latitude,
		}
	}
	session.SetLastHeartbeatRequestTimestampToNow()

	err = s.RegisterSession(ctx, session)
	if err != nil {
		return nil, err
	}

	// 5. Update location and notify online
	if location != nil {
		_ = s.sessionLocationService.UpsertUserLocation(ctx, userId, deviceType, location.Longitude, location.Latitude)
	}

	s.InvokeGoOnlineHandlers(ctx, nil, session)
	return session, nil
}

func (s *SessionService) resolveConflicts(ctx context.Context, userId int64, deviceType protocol.DeviceType, status *userbo.UserSessionsStatus) (bool, error) {
	conflictedDeviceTypes := s.userSimultaneousLoginService.GetConflictedDeviceTypes(deviceType)
	if len(conflictedDeviceTypes) == 0 {
		return false, nil
	}

	for _, conflictedDT := range conflictedDeviceTypes {
		if info, exists := status.OnlineDeviceTypeToSessionInfo[conflictedDT]; exists && info.IsActive {
			if s.userSimultaneousLoginService.ShouldDisconnectLoggingInDeviceIfConflicts() {
				return true, nil
			}
			// Kick existing
			if info.NodeID == s.nodeID {
				_ = s.CloseLocalSession(ctx, userId, []protocol.DeviceType{conflictedDT}, constant.SessionCloseStatus_DISCONNECTED_BY_OTHER_DEVICE)
			} else if s.rpcService != nil {
				req := &rpc.SetUserOfflineRequest{
					UserID:             userId,
					DeviceTypes:        []protocol.DeviceType{conflictedDT},
					SessionCloseStatus: int(constant.SessionCloseStatus_DISCONNECTED_BY_OTHER_DEVICE),
				}
				_, _ = s.rpcService.RequestResponse(ctx, info.NodeID, req)
			}
		}
	}
	return false, nil
}

func (s *SessionService) CloseLocalSession(ctx context.Context, userId int64, deviceTypes []protocol.DeviceType, closeReason constant.SessionCloseStatus) error {
	manager, ok := s.shardedMap.Get(userId)
	if !ok {
		return nil
	}

	manager.mu.Lock()
	defer manager.mu.Unlock()

	if len(deviceTypes) == 0 {
		// Close all sessions
		for dt, session := range manager.Sessions {
			s.unregisterSessionIp(session)
			delete(manager.Sessions, dt)
			if session.Conn != nil {
				_ = session.Conn.Close(closeReason)
			}
			if session.CloseChan != nil {
				select {
				case <-session.CloseChan:
				default:
					close(session.CloseChan)
				}
			}
			s.notifySessionClosedListeners(session)
			// Cleanup remote
			_, _ = s.userStatusService.RemoveOnlineDevice(ctx, userId, dt, s.nodeID)
			_ = s.sessionLocationService.RemoveUserLocation(ctx, userId, dt)
		}
	} else {
		for _, dt := range deviceTypes {
			if session, exists := manager.Sessions[dt]; exists {
				s.unregisterSessionIp(session)
				delete(manager.Sessions, dt)
				if session.Conn != nil {
					_ = session.Conn.Close(closeReason)
				}
				if session.CloseChan != nil {
					select {
					case <-session.CloseChan:
					default:
						close(session.CloseChan)
					}
				}
				s.notifySessionClosedListeners(session)
				// Cleanup remote
				_, _ = s.userStatusService.RemoveOnlineDevice(ctx, userId, dt, s.nodeID)
				_ = s.sessionLocationService.RemoveUserLocation(ctx, userId, dt)
			}
		}
	}

	if manager := s.shardedMap.RemoveIfEmpty(userId); manager != nil {
		s.InvokeGoOfflineHandlers(ctx, manager, closeReason)
	}
	return nil
}

func (s *SessionService) notifySessionClosedListeners(session *UserSession) {
	s.mu.RLock()
	listeners := s.onSessionClosedListeners
	s.mu.RUnlock()
	for _, listener := range listeners {
		listener(session)
	}
}

func (s *SessionService) CloseLocalSessionsByUserIds(ctx context.Context, userIds []int64, closeReason constant.SessionCloseStatus) error {
	for _, uid := range userIds {
		manager, ok := s.shardedMap.Get(uid)
		if ok {
			sessions := manager.GetAllSessions()
			for _, session := range sessions {
				s.UnregisterSession(ctx, uid, session.DeviceType, session.Conn, closeReason)
			}
		}
	}
	return nil
}

func (s *SessionService) AuthAndCloseLocalSession(ctx context.Context, userId int64, deviceType protocol.DeviceType, password *string, sessionId int64) error {
	_, err := s.sessionAuthenticationManager.VerifyAndGrant(ctx, sessionbo.NewUserLoginInfo(0, userId, password, deviceType, nil, nil, ""))
	if err != nil {
		return err
	}
	return s.CloseLocalSession(ctx, userId, []protocol.DeviceType{deviceType}, constant.SessionCloseStatus_DISCONNECTED_BY_CLIENT_REDUNDANTLY)
}

func (s *SessionService) UpdateLocalSession(ctx context.Context, userId int64, deviceType protocol.DeviceType, userStatus protocol.UserStatus, location *protocol.UserLocation) error {
	manager, ok := s.shardedMap.Get(userId)
	if !ok {
		return &SessionAuthError{Code: constant.ResponseStatusCode_UPDATE_NON_EXISTING_SESSION_STATUS}
	}

	manager.mu.Lock()
	session, exists := manager.Sessions[deviceType]
	if !exists {
		manager.mu.Unlock()
		return &SessionAuthError{Code: constant.ResponseStatusCode_UPDATE_NON_EXISTING_SESSION_STATUS}
	}

	if userStatus != protocol.UserStatus_OFFLINE { // Use OFFLINE or some placeholder to check
		manager.UserStatus = userStatus
	}
	manager.mu.Unlock()

	if userStatus != protocol.UserStatus_OFFLINE {
		_, _ = s.userStatusService.UpdateStatus(ctx, userId, userStatus)
	}
	if location != nil {
		_ = s.sessionLocationService.UpsertUserLocation(ctx, userId, deviceType, location.Longitude, location.Latitude)
		session.Location = &sessionbo.UserLocation{
			Longitude: location.Longitude,
			Latitude:  location.Latitude,
		}
	}
	session.SetLastHeartbeatRequestTimestampToNow()
	return nil
}

func (s *SessionService) CloseAllLocalSessions(ctx context.Context, closeReason constant.SessionCloseStatus) error {
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

func (s *SessionService) AuthAndUpdateHeartbeatTimestamp(ctx context.Context, userId int64, deviceType protocol.DeviceType, sessionId int64) *UserSession {
	if session, ok := s.GetUserSession(userId, deviceType); ok {
		if session.ID == sessionId {
			s.HandleHeartbeatUpdateRequest(session)
			return session
		}
	}
	return nil
}

func (s *SessionService) CloseLocalSessionsByIp(ctx context.Context, ips [][]byte, closeReason constant.SessionCloseStatus) (int, error) {
	count := 0
	for _, ip := range ips {
		ipStr := net.IP(ip).String()
		if v, ok := s.ipToSessions.Load(ipStr); ok {
			sessionMap := v.(*sync.Map)
			sessionMap.Range(func(key, value any) bool {
				session := key.(*UserSession)
				s.UnregisterSession(ctx, session.UserID, session.DeviceType, session.Conn, closeReason)
				count++
				return true
			})
		}
	}
	return count, nil
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

func (s *SessionService) InvokeGoOfflineHandlers(ctx context.Context, userSessionsManager *UserSessionsManager, closeReason constant.SessionCloseStatus) {
	// TODO: 插件系统尚未实现: 调用 PluginManager (UserOnlineStatusChangeHandler.goOffline)
}
