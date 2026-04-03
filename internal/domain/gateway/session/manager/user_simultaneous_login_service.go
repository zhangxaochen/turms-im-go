package manager

import "im.turms/server/pkg/protocol"

// UserSimultaneousLoginService maps to UserSimultaneousLoginService in Java.
// @MappedFrom UserSimultaneousLoginService
type UserSimultaneousLoginService struct {
}

// @MappedFrom getConflictedDeviceTypes(@NotNull @ValidDeviceType DeviceType deviceType)
func (s *UserSimultaneousLoginService) GetConflictedDeviceTypes(deviceType protocol.DeviceType) []protocol.DeviceType {
	// Stub implementation
	return nil
}

// @MappedFrom isForbiddenDeviceType(DeviceType deviceType)
func (s *UserSimultaneousLoginService) IsForbiddenDeviceType(deviceType protocol.DeviceType) bool {
	// Stub implementation
	return false
}

// @MappedFrom shouldDisconnectLoggingInDeviceIfConflicts()
func (s *UserSimultaneousLoginService) ShouldDisconnectLoggingInDeviceIfConflicts() bool {
	// Stub implementation
	return false
}
