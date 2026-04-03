package session

import "im.turms/server/pkg/protocol"

type SimultaneousLoginStrategy int

const (
	SimultaneousLoginStrategy_ALLOW_ONE_DEVICE_OF_EACH_DEVICE_TYPE_ONLINE SimultaneousLoginStrategy = iota
	SimultaneousLoginStrategy_ALLOW_ONE_DEVICE_FOR_ALL_DEVICE_TYPES_ONLINE
	SimultaneousLoginStrategy_ALLOW_ONE_DEVICE_OF_DESKTOP_AND_ONE_DEVICE_OF_MOBILE_ONLINE
	SimultaneousLoginStrategy_ALLOW_ONE_DEVICE_OF_DESKTOP_OR_BROWSER_AND_ONE_DEVICE_OF_MOBILE_ONLINE
	SimultaneousLoginStrategy_ALLOW_ONE_DEVICE_OF_DESKTOP_AND_ONE_DEVICE_OF_BROWSER_AND_ONE_DEVICE_OF_MOBILE_ONLINE
	SimultaneousLoginStrategy_ALLOW_ONE_DEVICE_OF_DESKTOP_OR_MOBILE_ONLINE
	SimultaneousLoginStrategy_ALLOW_ONE_DEVICE_OF_DESKTOP_OR_BROWSER_OR_MOBILE_ONLINE
)

type LoginConflictStrategy int

const (
	LoginConflictStrategy_DISCONNECT_LOGGED_IN_DEVICES LoginConflictStrategy = iota
	LoginConflictStrategy_DISCONNECT_LOGGING_IN_DEVICE
)

var allAvailableDeviceTypes = []protocol.DeviceType{
	protocol.DeviceType_DESKTOP,
	protocol.DeviceType_BROWSER,
	protocol.DeviceType_IOS,
	protocol.DeviceType_ANDROID,
	protocol.DeviceType_OTHERS,
	protocol.DeviceType_UNKNOWN,
}

// UserSimultaneousLoginService maps to UserSimultaneousLoginService in Java.
// @MappedFrom UserSimultaneousLoginService
type UserSimultaneousLoginService struct {
	deviceTypeToExclusiveDeviceTypes map[protocol.DeviceType]map[protocol.DeviceType]struct{}
	forbiddenDeviceTypes             map[protocol.DeviceType]struct{}
	allowDeviceTypeUnknownLogin      bool
	allowDeviceTypeOthersLogin       bool
	loginConflictStrategy            LoginConflictStrategy
}

func NewUserSimultaneousLoginService() *UserSimultaneousLoginService {
	// For refactor parity, use the default properties:
	strategy := SimultaneousLoginStrategy_ALLOW_ONE_DEVICE_OF_EACH_DEVICE_TYPE_ONLINE
	loginConflictStrategy := LoginConflictStrategy_DISCONNECT_LOGGED_IN_DEVICES
	allowDeviceTypeUnknownLogin := true
	allowDeviceTypeOthersLogin := true

	s := &UserSimultaneousLoginService{
		allowDeviceTypeUnknownLogin: allowDeviceTypeUnknownLogin,
		allowDeviceTypeOthersLogin:  allowDeviceTypeOthersLogin,
		loginConflictStrategy:       loginConflictStrategy,
	}

	s.deviceTypeToExclusiveDeviceTypes = newExclusiveDeviceFromStrategy(strategy)
	s.forbiddenDeviceTypes = newForbiddenDeviceTypesFromStrategy(strategy, allowDeviceTypeUnknownLogin, allowDeviceTypeOthersLogin)
	return s
}

func newExclusiveDeviceFromStrategy(strategy SimultaneousLoginStrategy) map[protocol.DeviceType]map[protocol.DeviceType]struct{} {
	m := make(map[protocol.DeviceType]map[protocol.DeviceType]struct{})

	addConflicted := func(d1, d2 protocol.DeviceType) {
		if m[d1] == nil {
			m[d1] = make(map[protocol.DeviceType]struct{})
		}
		m[d1][d2] = struct{}{}

		if m[d2] == nil {
			m[d2] = make(map[protocol.DeviceType]struct{})
		}
		m[d2][d1] = struct{}{}
	}

	addConflictedWithAll := func(d1 protocol.DeviceType) {
		for _, dt := range allAvailableDeviceTypes {
			addConflicted(d1, dt)
		}
	}

	for _, dt := range allAvailableDeviceTypes {
		addConflicted(dt, dt)
	}

	switch strategy {
	case SimultaneousLoginStrategy_ALLOW_ONE_DEVICE_OF_EACH_DEVICE_TYPE_ONLINE:
		// Base loop already adds itself
	case SimultaneousLoginStrategy_ALLOW_ONE_DEVICE_FOR_ALL_DEVICE_TYPES_ONLINE:
		for _, typeOne := range allAvailableDeviceTypes {
			for _, typeTwo := range allAvailableDeviceTypes {
				if typeOne != typeTwo {
					addConflicted(typeOne, typeTwo)
				}
			}
		}
	case SimultaneousLoginStrategy_ALLOW_ONE_DEVICE_OF_DESKTOP_AND_ONE_DEVICE_OF_MOBILE_ONLINE,
		SimultaneousLoginStrategy_ALLOW_ONE_DEVICE_OF_DESKTOP_AND_ONE_DEVICE_OF_BROWSER_AND_ONE_DEVICE_OF_MOBILE_ONLINE:
		addConflicted(protocol.DeviceType_ANDROID, protocol.DeviceType_IOS)
	case SimultaneousLoginStrategy_ALLOW_ONE_DEVICE_OF_DESKTOP_OR_BROWSER_AND_ONE_DEVICE_OF_MOBILE_ONLINE:
		addConflicted(protocol.DeviceType_DESKTOP, protocol.DeviceType_BROWSER)
		addConflicted(protocol.DeviceType_ANDROID, protocol.DeviceType_IOS)
	case SimultaneousLoginStrategy_ALLOW_ONE_DEVICE_OF_DESKTOP_OR_MOBILE_ONLINE:
		addConflicted(protocol.DeviceType_DESKTOP, protocol.DeviceType_ANDROID)
		addConflicted(protocol.DeviceType_DESKTOP, protocol.DeviceType_IOS)
		addConflicted(protocol.DeviceType_ANDROID, protocol.DeviceType_IOS)
	case SimultaneousLoginStrategy_ALLOW_ONE_DEVICE_OF_DESKTOP_OR_BROWSER_OR_MOBILE_ONLINE:
		addConflicted(protocol.DeviceType_DESKTOP, protocol.DeviceType_BROWSER)
		addConflicted(protocol.DeviceType_DESKTOP, protocol.DeviceType_ANDROID)
		addConflicted(protocol.DeviceType_DESKTOP, protocol.DeviceType_IOS)
		addConflicted(protocol.DeviceType_BROWSER, protocol.DeviceType_ANDROID)
		addConflicted(protocol.DeviceType_BROWSER, protocol.DeviceType_IOS)
		addConflicted(protocol.DeviceType_ANDROID, protocol.DeviceType_IOS)
	}

	// Always conflict with unknown and others unless specifically allowed
	addConflictedWithAll(protocol.DeviceType_UNKNOWN)
	addConflictedWithAll(protocol.DeviceType_OTHERS)

	return m
}

func newForbiddenDeviceTypesFromStrategy(strategy SimultaneousLoginStrategy, allowDeviceTypeUnknownLogin bool, allowDeviceTypeOthersLogin bool) map[protocol.DeviceType]struct{} {
	m := make(map[protocol.DeviceType]struct{})

	if !allowDeviceTypeUnknownLogin {
		m[protocol.DeviceType_UNKNOWN] = struct{}{}
	}
	if !allowDeviceTypeOthersLogin {
		m[protocol.DeviceType_OTHERS] = struct{}{}
	}
	if strategy == SimultaneousLoginStrategy_ALLOW_ONE_DEVICE_OF_DESKTOP_AND_ONE_DEVICE_OF_MOBILE_ONLINE ||
		strategy == SimultaneousLoginStrategy_ALLOW_ONE_DEVICE_OF_DESKTOP_OR_MOBILE_ONLINE {
		m[protocol.DeviceType_BROWSER] = struct{}{}
	}

	return m
}

// @MappedFrom getConflictedDeviceTypes(@NotNull @ValidDeviceType DeviceType deviceType)
func (s *UserSimultaneousLoginService) GetConflictedDeviceTypes(deviceType protocol.DeviceType) []protocol.DeviceType {
	exclusiveTypes, ok := s.deviceTypeToExclusiveDeviceTypes[deviceType]
	if !ok || len(exclusiveTypes) == 0 {
		return nil
	}
	conflicted := make([]protocol.DeviceType, 0, len(exclusiveTypes))
	for dt := range exclusiveTypes {
		conflicted = append(conflicted, dt)
	}
	return conflicted
}

// @MappedFrom isForbiddenDeviceType(DeviceType deviceType)
func (s *UserSimultaneousLoginService) IsForbiddenDeviceType(deviceType protocol.DeviceType) bool {
	_, forbidden := s.forbiddenDeviceTypes[deviceType]
	return forbidden
}

// @MappedFrom shouldDisconnectLoggingInDeviceIfConflicts()
func (s *UserSimultaneousLoginService) ShouldDisconnectLoggingInDeviceIfConflicts() bool {
	return s.loginConflictStrategy == LoginConflictStrategy_DISCONNECT_LOGGING_IN_DEVICE
}
