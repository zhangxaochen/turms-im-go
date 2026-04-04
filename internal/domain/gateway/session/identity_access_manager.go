package session

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"im.turms/server/internal/domain/common/constant"
	"im.turms/server/internal/domain/gateway/session/bo"
	userservice "im.turms/server/internal/domain/user/service"
)

// SessionIdentityAccessManagementSupport maps to the Java support interface.
type SessionIdentityAccessManagementSupport interface {
	VerifyAndGrant(ctx context.Context, loginInfo *bo.UserLoginInfo) (*bo.UserPermissionInfo, error)
	UpdateGlobalProperties(properties interface{})
}

// SessionIdentityAccessManager orchestrates the authentication.
// @MappedFrom SessionIdentityAccessManager
type SessionIdentityAccessManager struct {
	enableIdentityAccessManagement bool
	support                        SessionIdentityAccessManagementSupport
	userService                    userservice.UserService
}

func NewSessionIdentityAccessManager(userService userservice.UserService) *SessionIdentityAccessManager {
	return &SessionIdentityAccessManager{
		enableIdentityAccessManagement: true, // Configurable via properties in a real implementation
		support:                        &PasswordSessionIdentityAccessManager{userService: userService},
		userService:                    userService,
	}
}

func (m *SessionIdentityAccessManager) UpdateGlobalProperties(properties interface{}) {
	// Update settings
	// m.enableIdentityAccessManagement = ...
	if m.support != nil {
		m.support.UpdateGlobalProperties(properties)
	}
}

// @MappedFrom verifyAndGrant(int version, Long userId, String password, DeviceType deviceType, Map<String, String> deviceDetails, UserStatus userStatus, Location location, String ip)
func (m *SessionIdentityAccessManager) VerifyAndGrant(ctx context.Context, loginInfo *bo.UserLoginInfo) (*bo.UserPermissionInfo, error) {
	if loginInfo.UserID == 0 { // 0 is AdminConst.ADMIN_REQUESTER_ID
		return bo.LoginAuthenticationFailed, nil
	}

	if !m.enableIdentityAccessManagement {
		return bo.GrantedWithAllPermissions, nil
	}

	// TODO: 插件系统尚未实现: 调用 PluginManager (UserAuthenticator) 的钩子
	// pluginManager := GetPluginManager()
	// if pluginManager.HasRunningExtensions(plugin.UserAuthenticator) {
	//     // return pluginManager.InvokeExtensionPointsSimultaneously(...)
	// }

	if m.support != nil {
		return m.support.VerifyAndGrant(ctx, loginInfo)
	}

	return bo.NewUserPermissionInfo(constant.ResponseStatusCode_LOGIN_AUTHENTICATION_FAILED, nil), nil
}

// PasswordSessionIdentityAccessManager maps to PasswordSessionIdentityAccessManager in Java.
type PasswordSessionIdentityAccessManager struct {
	userService userservice.UserService
}

// @MappedFrom verifyAndGrant(UserLoginInfo userLoginInfo)
func (m *PasswordSessionIdentityAccessManager) VerifyAndGrant(ctx context.Context, loginInfo *bo.UserLoginInfo) (*bo.UserPermissionInfo, error) {
	if loginInfo.Password == nil || *loginInfo.Password == "" {
		return bo.NewUserPermissionInfo(constant.ResponseStatusCode_LOGIN_AUTHENTICATION_FAILED, nil), nil
	}

	password, err := m.userService.FindPassword(ctx, loginInfo.UserID)
	if err != nil {
		return nil, err
	}
	if password == nil {
		return bo.NewUserPermissionInfo(constant.ResponseStatusCode_LOGIN_AUTHENTICATION_FAILED, nil), nil
	}

	active, err := m.userService.IsActiveAndNotDeleted(ctx, loginInfo.UserID)
	if err != nil {
		return nil, err
	}
	if !active {
		return bo.LoggingInUserNotActive, nil
	}

	err = bcrypt.CompareHashAndPassword([]byte(*password), []byte(*loginInfo.Password))
	if err != nil {
		return bo.NewUserPermissionInfo(constant.ResponseStatusCode_LOGIN_AUTHENTICATION_FAILED, nil), nil
	}

	// Granted
	return bo.GrantedWithAllPermissions, nil
}

func (m *PasswordSessionIdentityAccessManager) UpdateGlobalProperties(properties interface{}) {
}

// HttpSessionIdentityAccessManager maps to HttpSessionIdentityAccessManager in Java.
type HttpSessionIdentityAccessManager struct{}

func (m *HttpSessionIdentityAccessManager) VerifyAndGrant(ctx context.Context, loginInfo *bo.UserLoginInfo) (*bo.UserPermissionInfo, error) {
	return nil, nil
}
func (m *HttpSessionIdentityAccessManager) UpdateGlobalProperties(properties interface{}) {}

// JwtSessionIdentityAccessManager maps to JwtSessionIdentityAccessManager in Java.
type JwtSessionIdentityAccessManager struct{}

func (m *JwtSessionIdentityAccessManager) VerifyAndGrant(ctx context.Context, loginInfo *bo.UserLoginInfo) (*bo.UserPermissionInfo, error) {
	return nil, nil
}
func (m *JwtSessionIdentityAccessManager) UpdateGlobalProperties(properties interface{}) {}

// LdapSessionIdentityAccessManager maps to LdapSessionIdentityAccessManager in Java.
type LdapSessionIdentityAccessManager struct{}

func (m *LdapSessionIdentityAccessManager) VerifyAndGrant(ctx context.Context, loginInfo *bo.UserLoginInfo) (*bo.UserPermissionInfo, error) {
	return nil, nil
}
func (m *LdapSessionIdentityAccessManager) UpdateGlobalProperties(properties interface{}) {}

// NoopSessionIdentityAccessManager maps to NoopSessionIdentityAccessManager in Java.
type NoopSessionIdentityAccessManager struct{}

func (m *NoopSessionIdentityAccessManager) VerifyAndGrant(ctx context.Context, loginInfo *bo.UserLoginInfo) (*bo.UserPermissionInfo, error) {
	// Java: returns GRANTED_WITH_ALL_PERMISSIONS which contains TurmsRequestTypePool.ALL
	// In Go, nil Permissions == all permissions (see UserSession.HasPermission)
	return bo.GrantedWithAllPermissions, nil
}
func (m *NoopSessionIdentityAccessManager) UpdateGlobalProperties(properties interface{}) {}
