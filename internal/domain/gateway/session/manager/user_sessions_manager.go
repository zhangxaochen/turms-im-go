package manager

import (
	"sync"

	"im.turms/server/internal/domain/gateway/session"
	"im.turms/server/pkg/protocol"
)

// UserSessionsManager manages all sessions for a specific user.
type UserSessionsManager struct {
	UserID     int64
	UserStatus protocol.UserStatus

	// deviceTypeToSession maps DeviceType to UserSession
	deviceTypeToSession sync.Map // protocol.DeviceType -> *session.UserSession
}

func NewUserSessionsManager(userID int64, userStatus protocol.UserStatus) *UserSessionsManager {
	return &UserSessionsManager{
		UserID:     userID,
		UserStatus: userStatus,
	}
}

// AddSessionIfAbsent adds a new user session if one for the given device type does not already exist.
// @MappedFrom addSessionIfAbsent(int version, Set<TurmsRequest.KindCase> permissions, DeviceType loggingInDeviceType, Map<String, String> deviceDetails, Location location)
func (m *UserSessionsManager) AddSessionIfAbsent(
	version int,
	permissions map[interface{}]struct{},
	loggingInDeviceType protocol.DeviceType,
	deviceDetails map[string]string,
	// location *Location,
) *session.UserSession {
	userSession := &session.UserSession{
		// Map properties to session.UserSession here
		// In the actual system, we'd initialize the necessary structs
	}

	actual, loaded := m.deviceTypeToSession.LoadOrStore(loggingInDeviceType, userSession)
	if !loaded {
		return actual.(*session.UserSession)
	}
	return nil
}

// CloseSession closes the session for a specific device type with a given reason.
// @MappedFrom closeSession(DeviceType deviceType, CloseReason closeReason)
func (m *UserSessionsManager) CloseSession(deviceType protocol.DeviceType, closeReason interface{}) bool {
	sessionData, loaded := m.deviceTypeToSession.LoadAndDelete(deviceType)
	if loaded {
		sessionObj := sessionData.(*session.UserSession)
		// Assuming Close method exists on UserSession that takes a reason
		_ = sessionObj // sessionObj.Close(closeReason)
		return true
	}
	return false
}

// PushSessionNotification pushes a session notification to a specific device.
// @MappedFrom pushSessionNotification(DeviceType deviceType, String serverId)
func (m *UserSessionsManager) PushSessionNotification(deviceType protocol.DeviceType, serverID string) bool {
	sessionData, ok := m.deviceTypeToSession.Load(deviceType)
	if !ok {
		return false
	}
	sessionObj := sessionData.(*session.UserSession)

	// Stub mapping
	// In the real system it encodes user session notification and sends it over the connection
	_ = sessionObj

	return true
}

// GetSession retrieves the session for a given device type.
// @MappedFrom getSession(DeviceType deviceType)
func (m *UserSessionsManager) GetSession(deviceType protocol.DeviceType) *session.UserSession {
	sessionData, ok := m.deviceTypeToSession.Load(deviceType)
	if !ok {
		return nil
	}
	return sessionData.(*session.UserSession)
}

// CountSessions returns the total number of sessions for this user.
// @MappedFrom countSessions()
func (m *UserSessionsManager) CountSessions() int {
	count := 0
	m.deviceTypeToSession.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	return count
}

// GetLoggedInDeviceTypes returns the device types that have an active session.
// @MappedFrom getLoggedInDeviceTypes()
func (m *UserSessionsManager) GetLoggedInDeviceTypes() []protocol.DeviceType {
	var deviceTypes []protocol.DeviceType
	m.deviceTypeToSession.Range(func(key, value interface{}) bool {
		deviceTypes = append(deviceTypes, key.(protocol.DeviceType))
		return true
	})
	return deviceTypes
}
