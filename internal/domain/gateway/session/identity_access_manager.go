package session

import (
	"context"

	"im.turms/server/internal/domain/gateway/session/bo"
)

// SessionIdentityAccessManager maps to SessionIdentityAccessManager in Java.
// @MappedFrom SessionIdentityAccessManager
type SessionIdentityAccessManager interface {
	// @MappedFrom verifyAndGrant(int version, Long userId, String password, DeviceType deviceType, Map<String, String> deviceDetails, UserStatus userStatus, Location location, String ip)
	VerifyAndGrant(ctx context.Context, loginInfo *bo.UserLoginInfo) (*bo.UserPermissionInfo, error)
}

// HttpSessionIdentityAccessManager maps to HttpSessionIdentityAccessManager in Java.
type HttpSessionIdentityAccessManager struct{}

// @MappedFrom verifyAndGrant(UserLoginInfo userLoginInfo)
// @MappedFrom verifyAndGrant(int version, Long userId, @Nullable String password, DeviceType deviceType, @Nullable Map<String, String> deviceDetails, @Nullable UserStatus userStatus, @Nullable Location location, @Nullable String ip)
func (m *HttpSessionIdentityAccessManager) VerifyAndGrant(ctx context.Context, loginInfo *bo.UserLoginInfo) (*bo.UserPermissionInfo, error) {
	// Stub implementation
	return nil, nil
}

// JwtSessionIdentityAccessManager maps to JwtSessionIdentityAccessManager in Java.
type JwtSessionIdentityAccessManager struct{}

// @MappedFrom verifyAndGrant(UserLoginInfo userLoginInfo)
// @MappedFrom verifyAndGrant(int version, Long userId, @Nullable String password, DeviceType deviceType, @Nullable Map<String, String> deviceDetails, @Nullable UserStatus userStatus, @Nullable Location location, @Nullable String ip)
func (m *JwtSessionIdentityAccessManager) VerifyAndGrant(ctx context.Context, loginInfo *bo.UserLoginInfo) (*bo.UserPermissionInfo, error) {
	// Stub implementation
	return nil, nil
}

// LdapSessionIdentityAccessManager maps to LdapSessionIdentityAccessManager in Java.
type LdapSessionIdentityAccessManager struct{}

// @MappedFrom verifyAndGrant(UserLoginInfo userLoginInfo)
// @MappedFrom verifyAndGrant(int version, Long userId, @Nullable String password, DeviceType deviceType, @Nullable Map<String, String> deviceDetails, @Nullable UserStatus userStatus, @Nullable Location location, @Nullable String ip)
func (m *LdapSessionIdentityAccessManager) VerifyAndGrant(ctx context.Context, loginInfo *bo.UserLoginInfo) (*bo.UserPermissionInfo, error) {
	// Stub implementation
	return nil, nil
}

// NoopSessionIdentityAccessManager maps to NoopSessionIdentityAccessManager in Java.
type NoopSessionIdentityAccessManager struct{}

// @MappedFrom verifyAndGrant(UserLoginInfo userLoginInfo)
// @MappedFrom verifyAndGrant(int version, Long userId, @Nullable String password, DeviceType deviceType, @Nullable Map<String, String> deviceDetails, @Nullable UserStatus userStatus, @Nullable Location location, @Nullable String ip)
func (m *NoopSessionIdentityAccessManager) VerifyAndGrant(ctx context.Context, loginInfo *bo.UserLoginInfo) (*bo.UserPermissionInfo, error) {
	// Stub implementation
	return nil, nil
}

// PasswordSessionIdentityAccessManager maps to PasswordSessionIdentityAccessManager in Java.
type PasswordSessionIdentityAccessManager struct{}

// @MappedFrom verifyAndGrant(UserLoginInfo userLoginInfo)
// @MappedFrom verifyAndGrant(int version, Long userId, @Nullable String password, DeviceType deviceType, @Nullable Map<String, String> deviceDetails, @Nullable UserStatus userStatus, @Nullable Location location, @Nullable String ip)
func (m *PasswordSessionIdentityAccessManager) VerifyAndGrant(ctx context.Context, loginInfo *bo.UserLoginInfo) (*bo.UserPermissionInfo, error) {
	// Stub implementation
	return nil, nil
}

// @MappedFrom updateGlobalProperties(TurmsProperties properties)
func (m *PasswordSessionIdentityAccessManager) UpdateGlobalProperties(properties interface{}) {
	// Stub implementation
}
