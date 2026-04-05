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
	Code    constant.ResponseStatusCode
	Message string
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

	// @MappedFrom Java: SessionService.notifyClientsOfSessionInfoAfterConnected
	// Default true: send session notification (sessionId + serverId) to the new device upon login.
	notifyClientsOfSessionInfoAfterConnected bool
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
		// Java default: notifyClientsOfSessionInfoAfterConnected = true
		notifyClientsOfSessionInfoAfterConnected: true,
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
			svc.CloseLocalSession(ctx, req.UserID, req.DeviceTypes, sessionbo.NewCloseReason(constant.SessionCloseStatus(req.SessionCloseStatus)))
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
			existingSession.Close(sessionbo.NewCloseReason(constant.SessionCloseStatus_DISCONNECTED_BY_OTHER_DEVICE))
			delete(manager.Sessions, session.DeviceType)
		}
	}

	manager.Sessions[session.DeviceType] = session
	session.SetLastHeartbeatRequestTimestampToNow()
	session.SetLastRequestTimestampToNow()
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
		
		isEmpty := true
		sessionMap.Range(func(key, value any) bool {
			isEmpty = false
			return false
		})
		if isEmpty {
			s.ipToSessions.Delete(ipStr)
		}
	}
}

func (s *SessionService) UnregisterSession(ctx context.Context, userID int64, deviceType protocol.DeviceType, conn Connection, closeReason sessionbo.CloseReason) {
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

	if (conn != nil && existing.Conn == conn) || conn == nil {
		delete(manager.Sessions, deviceType)

		existing.Close(closeReason)
		s.unregisterSessionIp(existing)
		s.notifySessionClosedListeners(existing)
	}
	manager.mu.Unlock()

	s.shardedMap.RemoveIfEmpty(userID)
	s.InvokeGoOfflineHandlers(ctx, manager, closeReason.Status)

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
	_, _ = s.CloseAllLocalSessions(ctx, sessionbo.NewCloseReason(constant.SessionCloseStatus_SERVER_CLOSED))
	return nil
}

func (s *SessionService) HandleHeartbeatUpdateRequest(session *UserSession) {
	session.SetLastHeartbeatRequestTimestampToNow()
}

func (s *SessionService) HandleLoginRequest(ctx context.Context, version int, ip []byte, userId int64, password string, deviceType protocol.DeviceType, deviceDetails map[string]string, userStatus protocol.UserStatus, location *protocol.UserLocation, ipStr string) (*UserSession, error) {
	if ip == nil {
		return nil, errors.New("ip cannot be nil")
	}
	if userId == 0 {
		return nil, errors.New("userId cannot be 0")
	}
	if version != 1 {
		return nil, &SessionAuthError{Code: constant.ResponseStatusCode_UNSUPPORTED_CLIENT_VERSION}
	}
	if s.userSimultaneousLoginService.IsForbiddenDeviceType(deviceType) {
		return nil, &SessionAuthError{Code: constant.ResponseStatusCode_LOGIN_FROM_FORBIDDEN_DEVICE_TYPE}
	}

	var passwordPtr *string
	if password != "" {
		passwordPtr = &password
	}
	loginInfo := sessionbo.NewUserLoginInfo(version, userId, passwordPtr, deviceType, deviceDetails, &userStatus, location, ipStr)
	permissionInfo, err := s.sessionAuthenticationManager.VerifyAndGrant(ctx, loginInfo)
	if err != nil {
		return nil, err
	}
	if permissionInfo.AuthenticationCode != constant.ResponseStatusCode_OK {
		return nil, &SessionAuthError{Code: permissionInfo.AuthenticationCode}
	}

	permissions := permissionInfo.Permissions

	return s.TryRegisterOnlineUser(ctx, version, permissions, ip, userId, deviceType, deviceDetails, userStatus, location)
}

func (s *SessionService) TryRegisterOnlineUser(ctx context.Context, version int, permissions map[int32]bool, ip []byte, userId int64, deviceType protocol.DeviceType, deviceDetails map[string]string, userStatus protocol.UserStatus, location *protocol.UserLocation) (*UserSession, error) {
	if ip == nil {
		return nil, &SessionAuthError{Code: constant.ResponseStatusCode_ILLEGAL_ARGUMENT, Message: "ip must not be null"}
	}
	if deviceType == protocol.DeviceType_UNKNOWN {
		return nil, &SessionAuthError{Code: constant.ResponseStatusCode_ILLEGAL_ARGUMENT, Message: "deviceType must not be UNKNOWN"}
	}
	if userStatus == protocol.UserStatus_OFFLINE {
		return nil, &SessionAuthError{Code: constant.ResponseStatusCode_ILLEGAL_ARGUMENT, Message: "userStatus must not be OFFLINE"}
	}
	if location != nil {
		if location.Longitude < -180 || location.Longitude > 180 { // MIN/MAX range approx
			return nil, &SessionAuthError{Code: constant.ResponseStatusCode_ILLEGAL_ARGUMENT, Message: "longitude out of range"}
		}
		if location.Latitude < -90 || location.Latitude > 90 {
			return nil, &SessionAuthError{Code: constant.ResponseStatusCode_ILLEGAL_ARGUMENT, Message: "latitude out of range"}
		}
	}
	if len(deviceDetails) > 0 {
		newDetails := make(map[string]string)
		for k, v := range deviceDetails {
			if v != "" {
				newDetails[k] = v
			}
		}
		deviceDetails = newDetails
	}

	sessionsStatus, err := s.userStatusService.FetchUserSessionsStatus(ctx, userId)
	if err != nil {
		return nil, err
	}

	manager, _ := s.shardedMap.Get(userId)
	if manager != nil {
		deviceTypeToSessionInfo := sessionsStatus.OnlineDeviceTypeToSessionInfo
		for _, typ := range manager.GetLoggedInDeviceTypes() {
			sessionInfo, exists := deviceTypeToSessionInfo[typ]
			if exists && sessionInfo.IsActive && s.nodeID != sessionInfo.NodeID {
				_, _ = s.CloseLocalSession(ctx, userId, []protocol.DeviceType{typ}, sessionbo.NewCloseReason(constant.SessionCloseStatus_DISCONNECTED_BY_OTHER_DEVICE))
			}
		}
	}

	existingUserStatus := sessionsStatus.UserStatus
	// If the user is offline, register the current session.
	if existingUserStatus == protocol.UserStatus_OFFLINE {
		return s.addOnlineDeviceIfAbsent(ctx, version, permissions, ip, userId, deviceType, deviceDetails, userStatus, location, nil, nil)
	}

	// If the user is already online, check if there is any device conflict.
	sessionInfo, ok := sessionsStatus.OnlineDeviceTypeToSessionInfo[deviceType]
	if ok && sessionInfo.IsActive {
		session, _ := s.GetUserSession(userId, deviceType)
		isClosedSessionOnLocal := session != nil && session.Conn != nil && !session.Conn.IsConnected()
		if isClosedSessionOnLocal {
			// Replace disconnected connection
			if existingUserStatus != userStatus || userStatus == 0 {
				// Java checks userStatus == null || existingUserStatus == userStatus
				// If userStatus == 0 (UNRECOGNIZED), don't update
			} else {
				_, err = s.userStatusService.UpdateStatus(ctx, userId, userStatus)
				if err != nil {
					fmt.Printf("failed to update online status for user %d: %v\n", userId, err)
				}
			}
			if location != nil {
				err = s.sessionLocationService.UpsertUserLocation(ctx, userId, deviceType, location.Longitude, location.Latitude)
				if err != nil {
					fmt.Printf("failed to upsert location for user %d: %v\n", userId, err)
				}
			}
			return session, nil
		} else if s.userSimultaneousLoginService.ShouldDisconnectLoggingInDeviceIfConflicts() {
			return nil, &SessionAuthError{Code: constant.ResponseStatusCode_SESSION_SIMULTANEOUS_CONFLICTS_DECLINE}
		}
	}

	var expectedNodeId *string
	var expectedDeviceTimestamp *int64
	if ok {
		expectedNodeId = &sessionInfo.NodeID
		ts := sessionInfo.HeartbeatTimestampSeconds
		expectedDeviceTimestamp = &ts
	}

	wasSuccessful, err := s.closeSessionsWithConflictedDeviceTypes(ctx, userId, deviceType, sessionsStatus)
	if err != nil {
		return nil, err
	}
	if wasSuccessful {
		return s.addOnlineDeviceIfAbsent(ctx, version, permissions, ip, userId, deviceType, deviceDetails, userStatus, location, expectedNodeId, expectedDeviceTimestamp)
	}
	return nil, &SessionAuthError{Code: constant.ResponseStatusCode_SESSION_SIMULTANEOUS_CONFLICTS_DECLINE}
}

func (s *SessionService) addOnlineDeviceIfAbsent(ctx context.Context, version int, permissions map[int32]bool, ip []byte, userId int64, deviceType protocol.DeviceType, deviceDetails map[string]string, userStatus protocol.UserStatus, location *protocol.UserLocation, expectedNodeId *string, expectedDeviceTimestamp *int64) (*UserSession, error) {
	now := time.Now()
	added, err := s.userStatusService.AddOnlineDevice(ctx, userId, deviceType, deviceDetails, userStatus, s.nodeID, &now, expectedNodeId, expectedDeviceTimestamp)
	if err != nil {
		return nil, err
	}
	if !added {
		return nil, &SessionAuthError{Code: constant.ResponseStatusCode_SESSION_SIMULTANEOUS_CONFLICTS_DECLINE}
	}

	var loc *sessionbo.UserLocation
	if location != nil {
		loc = &sessionbo.UserLocation{
			Longitude: location.Longitude,
			Latitude:  location.Latitude,
			Timestamp: location.Timestamp,
			Details:   location.Details,
		}
	}

	session := NewUserSession(
		int32(version),
		permissions,
		userId,
		deviceType,
		deviceDetails,
		loc,
	)
	session.ID = s.nextSessionID()
	session.SetConnection(nil, net.IP(ip))

	err = s.RegisterSession(ctx, session)
	if err != nil {
		_, _ = s.CloseLocalSession(ctx, userId, []protocol.DeviceType{deviceType}, sessionbo.NewCloseReason(constant.SessionCloseStatus_SERVER_ERROR))
		return nil, &SessionAuthError{Code: constant.ResponseStatusCode_SERVER_INTERNAL_ERROR, Message: "Caught an error while adding the user session: " + err.Error()}
	}

	if location != nil && deviceType != protocol.DeviceType_BROWSER {
		_ = s.sessionLocationService.UpsertUserLocation(ctx, userId, deviceType, location.Longitude, location.Latitude)
	}

	manager, _ := s.shardedMap.Get(userId)
	s.OnSessionEstablished(ctx, manager, deviceType)
	s.InvokeGoOnlineHandlers(ctx, manager, session)
	return session, nil
}

func (s *SessionService) closeSessionsWithConflictedDeviceTypes(ctx context.Context, userId int64, deviceType protocol.DeviceType, status *userbo.UserSessionsStatus) (bool, error) {
	conflictedDeviceTypes := s.userSimultaneousLoginService.GetConflictedDeviceTypes(deviceType)
	if len(conflictedDeviceTypes) == 0 {
		return true, nil
	}

	// Group by NodeID to minimize RPC calls (Bug 892)
	nodeToConflictedTypes := make(map[string][]protocol.DeviceType)
	for _, conflictedDT := range conflictedDeviceTypes {
		if info, exists := status.OnlineDeviceTypeToSessionInfo[conflictedDT]; exists && info.IsActive {
			if s.userSimultaneousLoginService.ShouldDisconnectLoggingInDeviceIfConflicts() {
				return true, nil
			}
			// Bug fix: Java adds deviceType (the logging-in device) to the node mapping,
			// not conflictedDT. Java tells remote nodes to disconnect the logging-in device type.
			nodeToConflictedTypes[info.NodeID] = append(nodeToConflictedTypes[info.NodeID], deviceType)
		}
	}

	if len(nodeToConflictedTypes) == 0 {
		return true, nil
	}

	for nodeID, dts := range nodeToConflictedTypes {
		if nodeID == s.nodeID {
			_, _ = s.CloseLocalSession(ctx, userId, dts, sessionbo.NewCloseReason(constant.SessionCloseStatus_DISCONNECTED_BY_CLIENT))
		} else if s.rpcService != nil {
			req := &rpc.SetUserOfflineRequest{
				UserID:             userId,
				DeviceTypes:        dts,
				SessionCloseStatus: int(constant.SessionCloseStatus_DISCONNECTED_BY_CLIENT),
			}
			_, err := s.rpcService.RequestResponse(ctx, nodeID, req)
			if err != nil {
				// Bug 860: Handle ConnectionNotFound with node discovery fallback
				if errors.Is(err, rpc.ErrConnectionNotFound) {
					discovery := s.rpcService.DiscoveryService()
					if discovery != nil && !discovery.IsKnownMember(nodeID) {
						continue // Consider offline, move to next
					}
				}
				return false, err
			}
		} else {
			return false, errors.New("rpc service is missing, cannot resolve distributed conflicts")
		}
	}
	return true, nil
}

func (s *SessionService) CloseLocalSession(ctx context.Context, userId int64, deviceTypes []protocol.DeviceType, closeReason sessionbo.CloseReason) (int, error) {
	if userId == 0 {
		return 0, errors.New("userId must not be 0")
	}
	if closeReason.Status == 0 { // Or whatever implies invalid or nil
		// actually closeReason is a struct, not a pointer in Go. It's essentially not-nil implicitly.
	}

	manager, ok := s.shardedMap.Get(userId)
	if !ok {
		return 0, nil
	}

	manager.mu.Lock()

	var toClose []protocol.DeviceType
	// Java: manager.getLoggedInDeviceTypes()
	if len(deviceTypes) == 0 {
		toClose = manager.GetLoggedInDeviceTypes() // Return a copy of logged-in device types
	} else {
		for _, dt := range deviceTypes {
			if _, exists := manager.Sessions[dt]; exists {
				toClose = append(toClose, dt)
			}
		}
	}

	if len(toClose) == 0 {
		manager.mu.Unlock()
		return 0, nil
	}

	// 1. Cleanup remote status BEFORE local closing (matches Java)
	_, _ = s.userStatusService.RemoveOnlineDevices(ctx, userId, toClose, s.nodeID)
	if s.sessionLocationService != nil && s.sessionLocationService.IsLocationEnabled() {
		_ = s.sessionLocationService.RemoveUserLocations(ctx, userId, toClose)
	}

	count := 0
	for _, dt := range toClose {
		if sess, exists := manager.Sessions[dt]; exists {
			wasOpen := sess.IsOpen()
			s.unregisterSessionIp(sess)
			delete(manager.Sessions, dt)
			sess.Close(closeReason)
			if wasOpen {
				s.notifySessionClosedListeners(sess)
			}
			count++
		}
	}

	manager.mu.Unlock()
	s.shardedMap.RemoveIfEmpty(userId)
	// Call InvokeGoOfflineHandlers outside the lock (matches Java behavior)
	s.InvokeGoOfflineHandlers(ctx, manager, closeReason.Status)
	return count, nil
}

// Java's removeSessionsManagerIfEmpty:
//  if manager.countSessions() == 0 { userIdToSessionsManager.remove(...) }
//  pluginManager.invokeGoOfflineHandlers(userId, closeReason)

func (s *SessionService) notifySessionClosedListeners(session *UserSession) {
	s.mu.RLock()
	listeners := s.onSessionClosedListeners
	s.mu.RUnlock()
	for _, listener := range listeners {
		// @MappedFrom Java: wraps each invocation in try-catch to prevent one bad listener from aborting others
		func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("WARN: panic in session closed listener: %v\n", r)
				}
			}()
			listener(session)
		}()
	}
}

func (s *SessionService) CloseLocalSessions(ctx context.Context, userIds []int64, ips [][]byte, closeReason sessionbo.CloseReason) (int, error) {
	if len(userIds) == 0 && len(ips) == 0 {
		return 0, nil
	}

	userIdSet := make(map[int64]struct{})
	for _, id := range userIds {
		userIdSet[id] = struct{}{}
	}
	ipSet := make(map[string]struct{})
	for _, ip := range ips {
		if ip != nil {
			ipSet[net.IP(ip).String()] = struct{}{}
		}
	}

	totalCount := 0
	checkIPs := len(ipSet) > 0
	checkIDs := len(userIdSet) > 0

	if !checkIDs && !checkIPs {
		return 0, nil
	}

	if checkIDs {
		for userId := range userIdSet {
			n, _ := s.CloseLocalSession(ctx, userId, nil, closeReason)
			totalCount += n
		}
	}

	if checkIPs {
		// Fix: Close ALL sessions matching the IP regardless of userIdSet.
		// Java has completely separate methods for userIds and IPs, always closing
		// all sessions for each IP. The userIdSet filtering was a logic bug.
		for ipStr := range ipSet {
			if v, ok := s.ipToSessions.Load(ipStr); ok {
				sessionMap := v.(*sync.Map)
				sessionMap.Range(func(key, value any) bool {
					sess := key.(*UserSession)
					s.CloseLocalSession(ctx, sess.UserID, []protocol.DeviceType{sess.DeviceType}, closeReason)
					totalCount++
					return true
				})
			}
		}
	}

	return totalCount, nil
}

func (s *SessionService) QuerySessions(ctx context.Context, userIds []int64) []*sessionbo.UserSessionsInfo {
	infos := make([]*sessionbo.UserSessionsInfo, 0, len(userIds))
	for _, userId := range userIds {
		infos = append(infos, s.getUserSessions(userId))
	}
	return infos
}

func (s *SessionService) getUserSessions(userID int64) *sessionbo.UserSessionsInfo {
	manager, ok := s.shardedMap.Get(userID)
	if !ok {
		return &sessionbo.UserSessionsInfo{
			UserID:   userID,
			Status:   protocol.UserStatus_OFFLINE,
			Sessions: nil,
		}
	}
	
	sessionList := manager.GetAllSessions()
	if len(sessionList) == 0 {
		return &sessionbo.UserSessionsInfo{
			UserID:   userID,
			Status:   protocol.UserStatus_OFFLINE,
			Sessions: nil,
		}
	}
	
	infos := make([]sessionbo.UserSessionInfo, 0, len(sessionList))
	for _, sess := range sessionList {
		var loc *sessionbo.UserLocation
		if sess.Location != nil {
			loc = &sessionbo.UserLocation{
				Longitude: sess.Location.Longitude,
				Latitude:  sess.Location.Latitude,
				Timestamp: sess.Location.Timestamp,
				Details:   sess.Location.Details,
			}
		}
		
		infos = append(infos, sessionbo.UserSessionInfo{
			ID:                                  sess.ID,
			Version:                             sess.Version.Load(),
			DeviceType:                          sess.DeviceType,
			DeviceDetails:                       sess.DeviceDetails,
			LoginDate:                           sess.LoginDate.UnixMilli(),
			LastHeartbeatRequestTimestampMillis: sess.GetLastHeartbeatRequestTimestamp(),
			LastRequestTimestampMillis:          sess.GetLastRequestTimestamp(),
			IsSessionOpen:                       sess.IsOpen(),
			Location:                            loc,
			IP:                                  sess.IP,
		})
	}

	return &sessionbo.UserSessionsInfo{
		UserID:   userID,
		Status:   manager.UserStatus,
		Sessions: infos,
	}
}

func (s *SessionService) QuerySessionsCount() int {
	return s.shardedMap.CountOnlineUsers()
}

func (s *SessionService) SetUserOffline(ctx context.Context, userID int64, closeReason sessionbo.CloseReason) (int, error) {
	return s.CloseLocalSession(ctx, userID, nil, closeReason)
}

func (s *SessionService) SetUsersOffline(ctx context.Context, userIDs []int64, closeReason sessionbo.CloseReason) (int, error) {
	total := 0
	for _, userID := range userIDs {
		count, _ := s.SetUserOffline(ctx, userID, closeReason)
		total += count
	}
	return total, nil
}

func (s *SessionService) CloseLocalSessionsByUserIds(ctx context.Context, userIds []int64, closeReason sessionbo.CloseReason) (int, error) {
	return s.CloseLocalSessions(ctx, userIds, nil, closeReason)
}

func (s *SessionService) AuthAndCloseLocalSession(ctx context.Context, userId int64, deviceType protocol.DeviceType, sessionId int64, closeReason sessionbo.CloseReason) (int, error) {
	if userId == 0 {
		return 0, errors.New("userId must not be 0")
	}
	manager, ok := s.shardedMap.Get(userId)
	if !ok {
		return 0, nil
	}
	session := manager.GetSession(deviceType)
	if session == nil || session.ID != sessionId {
		return 0, nil
	}
	return s.CloseLocalSession(ctx, userId, []protocol.DeviceType{deviceType}, closeReason)
}

func (s *SessionService) UpdateLocalSession(ctx context.Context, userId int64, deviceType protocol.DeviceType, userStatus protocol.UserStatus, location *protocol.UserLocation) error {
	if userStatus == protocol.UserStatus_OFFLINE && location == nil {
		return nil // Matches Java behavior or throw if required
	}
	manager, ok := s.shardedMap.Get(userId)
	if !ok {
		return &SessionAuthError{Code: constant.ResponseStatusCode_UPDATE_NON_EXISTING_SESSION_STATUS}
	}

	manager.mu.Lock()
	defer manager.mu.Unlock()

	session, exists := manager.Sessions[deviceType]
	if !exists {
		return &SessionAuthError{Code: constant.ResponseStatusCode_UPDATE_NON_EXISTING_SESSION_STATUS}
	}

	if userStatus != protocol.UserStatus_OFFLINE {
		manager.UserStatus = userStatus
	}
	if location != nil {
		session.Location = &sessionbo.UserLocation{
			Longitude: location.Longitude,
			Latitude:  location.Latitude,
		}
	}
	return nil
}

func (s *SessionService) CloseAllLocalSessions(ctx context.Context, closeReason sessionbo.CloseReason) (int, error) {
	count := 0
	s.shardedMap.Range(func(userId int64, manager *UserSessionsManager) bool {
		n, _ := s.CloseLocalSession(ctx, userId, nil, closeReason)
		count += n
		return true
	})
	return count, nil
}

func (s *SessionService) GetLocalUserSessionsInfo(ctx context.Context, userIds []int64) []*sessionbo.UserSessionsInfo {
	var result []*sessionbo.UserSessionsInfo
	for _, uid := range userIds {
		manager, ok := s.shardedMap.Get(uid)
		if !ok {
			// Java: Include offline users with OFFLINE status
			result = append(result, &sessionbo.UserSessionsInfo{
				UserID:   uid,
				Status:   protocol.UserStatus_OFFLINE,
				Sessions: []sessionbo.UserSessionInfo{},
			})
			continue
		}

		sessions := manager.GetAllSessions()
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
				ID:                                  sess.ID,
				Version:                             sess.Version.Load(),
				DeviceType:                          sess.DeviceType,
				DeviceDetails:                       sess.DeviceDetails,
				LoginDate:                           sess.LoginDate.UnixMilli(),
				LastHeartbeatRequestTimestampMillis: sess.GetLastHeartbeatRequestTimestamp(),
				LastRequestTimestampMillis:          sess.GetLastRequestTimestamp(),
				IsSessionOpen:                       sess.IsOpen(),
				Location:                            loc,
				IP:                                  []byte(sess.IP),
			})
		}

		result = append(result, &sessionbo.UserSessionsInfo{
			UserID:   uid,
			Status:   manager.UserStatus, // CORRECT: Includes manager status (Bug 877)
			Sessions: sessionInfos,
		})
	}
	return result
}

func (s *SessionService) AuthAndUpdateHeartbeatTimestamp(ctx context.Context, userId int64, deviceType protocol.DeviceType, sessionId int64) *UserSession {
	// Validate deviceType (matches Java: Validator.notNull + DeviceTypeUtil.validDeviceType)
	if deviceType == protocol.DeviceType_UNKNOWN {
		return nil
	}
	if session, ok := s.GetUserSession(userId, deviceType); ok {
		// Bug 882 & 849: Check if connection is active and NOT recovering (matches Java)
		if session.ID == sessionId && session.Conn != nil && session.Conn.IsActive() && !session.Conn.IsConnectionRecovering() {
			s.HandleHeartbeatUpdateRequest(session)
			return session
		}
	}
	return nil
}

func (s *SessionService) CloseLocalSessionsByIp(ctx context.Context, ips [][]byte, closeReason sessionbo.CloseReason) (int, error) {
	if len(ips) == 0 {
		return 0, nil
	}
	return s.CloseLocalSessions(ctx, nil, ips, closeReason)
}

func (s *SessionService) GetUserSessionsManager(ctx context.Context, userId int64) *UserSessionsManager {
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

func (s *SessionService) GetLocalUserSessionsByIp(ip []byte) []*UserSession {
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

func (s *SessionService) OnSessionEstablished(ctx context.Context, manager *UserSessionsManager, deviceType protocol.DeviceType) {
	// TODO: Increment metrics (e.g. LoggedInUsersCounter) (Bug 851)

	// @MappedFrom Java: if (notifyClientsOfSessionInfoAfterConnected) userSessionsManager.pushSessionNotification(deviceType, serverId)
	// Java sends to the device that just logged in (deviceType), NOT to other devices.
	// The notification includes the session's numeric ID and the server ID.
	if s.notifyClientsOfSessionInfoAfterConnected && manager != nil {
		manager.PushSessionNotification(ctx, deviceType, s.nodeID)
	}
}

func (s *SessionService) AddOnSessionClosedListeners(ctx context.Context, onSessionClosed func(*UserSession)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.onSessionClosedListeners = append(s.onSessionClosedListeners, onSessionClosed)
}

func (s *SessionService) InvokeGoOnlineHandlers(ctx context.Context, userSessionsManager *UserSessionsManager, userSession *UserSession) {
	// TODO: 插件系统尚未实现: 调用 PluginManager (UserOnlineStatusChangeHandler.goOnline)
}

func (s *SessionService) InvokeGoOfflineHandlers(ctx context.Context, userSessionsManager *UserSessionsManager, closeReason constant.SessionCloseStatus) {
	// TODO: plugin hooks (Bug 897)
}

func (s *SessionService) CloseLocalSessionByDeviceType(ctx context.Context, userId int64, deviceType protocol.DeviceType, closeReason sessionbo.CloseReason) bool {
	count, _ := s.CloseLocalSession(ctx, userId, []protocol.DeviceType{deviceType}, closeReason)
	return count > 0
}

// CloseLocalSessionByDeviceTypeWithStatus wraps a SessionCloseStatus into a CloseReason,
// matching Java's closeLocalSession(userId, deviceType, closeStatus) overload.
func (s *SessionService) CloseLocalSessionByDeviceTypeWithStatus(ctx context.Context, userId int64, deviceType protocol.DeviceType, closeStatus constant.SessionCloseStatus) bool {
	return s.CloseLocalSessionByDeviceType(ctx, userId, deviceType, sessionbo.NewCloseReason(closeStatus))
}
