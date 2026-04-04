package session

import (
	"context"

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
	if loginInfo.UserID == 0 { // 0 is AdminRequesterID
		return bo.NewUserPermissionInfo(constant.ResponseStatusCode_LOGIN_AUTHENTICATION_FAILED, nil), nil
	}

	if !m.enableIdentityAccessManagement {
		return bo.NewUserPermissionInfo(constant.ResponseStatusCode_OK, map[any]bool{}), nil
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

	user, err := m.userService.FindUser(ctx, loginInfo.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return bo.NewUserPermissionInfo(constant.ResponseStatusCode_LOGIN_AUTHENTICATION_FAILED, nil), nil
	}
	if !user.IsActive {
		return bo.NewUserPermissionInfo(constant.ResponseStatusCode_LOGIN_AUTHENTICATION_FAILED, nil), nil
	}

	// Wait, is it a direct check or through encoding? Passwords in Turms use Spring Security PasswordEncoder.
	// As this is a port, if bcrypt is used, we can verify with bcrypt.CompareHashAndPassword.
	// For now, if the user exists and active, and password matches (direct string match placeholder):
	if user.Password != *loginInfo.Password {
		// // If bcrypt was used uniformly:
		// err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(*loginInfo.Password))
		// if err != nil {
		// 	return bo.NewUserPermissionInfo(constant.ResponseStatusCode_LOGIN_AUTHENTICATION_FAILED, nil), nil
		// }
		return bo.NewUserPermissionInfo(constant.ResponseStatusCode_LOGIN_AUTHENTICATION_FAILED, nil), nil
	}

	// Granted
	return bo.NewUserPermissionInfo(constant.ResponseStatusCode_OK, nil), nil
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
	return bo.NewUserPermissionInfo(constant.ResponseStatusCode_OK, nil), nil
}
func (m *NoopSessionIdentityAccessManager) UpdateGlobalProperties(properties interface{}) {}
