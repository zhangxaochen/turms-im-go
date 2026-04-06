package bo

import (
	"im.turms/server/pkg/protocol"
)

// UserLoginInfo represents the login information provided by the user during connection establishment.
// Java parity: userId is Long (nullable boxed type) in Java, so we use *int64 to represent nil.
type UserLoginInfo struct {
	Version             int                   `json:"version"`
	UserID              *int64                `json:"userId"`
	Password            *string               `json:"password,omitempty"`
	LoggingInDeviceType protocol.DeviceType   `json:"deviceType"`
	DeviceDetails       map[string]string     `json:"deviceDetails,omitempty"`
	UserStatus          *protocol.UserStatus  `json:"userStatus,omitempty"`
	Location            *protocol.UserLocation `json:"location,omitempty"`
	IP                  string                `json:"ip"`
}

// NewUserLoginInfo creates a new UserLoginInfo.
// @MappedFrom UserLoginInfo(int version, Long userId, String password, DeviceType loggingInDeviceType, Map<String, String> deviceDetails, UserStatus userStatus, Location location, String ip)
func NewUserLoginInfo(version int, userID *int64, password *string, loggingInDeviceType protocol.DeviceType, deviceDetails map[string]string, userStatus *protocol.UserStatus, location *protocol.UserLocation, ip string) *UserLoginInfo {
	return &UserLoginInfo{
		Version:             version,
		UserID:              userID,
		Password:            password,
		LoggingInDeviceType: loggingInDeviceType,
		DeviceDetails:       deviceDetails,
		UserStatus:          userStatus,
		Location:            location,
		IP:                  ip,
	}
}
