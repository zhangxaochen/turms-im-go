package bo

import (
	"reflect"

	"im.turms/server/internal/domain/common/constant"
)

// UserPermissionInfo represents the permissions granted to a user session.
type UserPermissionInfo struct {
	AuthenticationCode constant.ResponseStatusCode
	Permissions        map[reflect.Type]struct{} // Set of TurmsRequest.KindCase
}

// NewUserPermissionInfo creates a new UserPermissionInfo.
// @MappedFrom UserPermissionInfo(ResponseStatusCode authenticationCode, Set<TurmsRequest.KindCase> permissions)
func NewUserPermissionInfo(authenticationCode constant.ResponseStatusCode, permissions map[reflect.Type]struct{}) *UserPermissionInfo {
	return &UserPermissionInfo{
		AuthenticationCode: authenticationCode,
		Permissions:        permissions,
	}
}

// NewUserPermissionInfoCodeOnly creates a new UserPermissionInfo with an empty set of permissions.
// @MappedFrom UserPermissionInfo(ResponseStatusCode authenticationCode)
func NewUserPermissionInfoCodeOnly(authenticationCode constant.ResponseStatusCode) *UserPermissionInfo {
	return &UserPermissionInfo{
		AuthenticationCode: authenticationCode,
		Permissions:        make(map[reflect.Type]struct{}),
	}
}

var (
	GrantedWithAllPermissions = NewUserPermissionInfoCodeOnly(constant.ResponseStatusCode_OK)
	LoginAuthenticationFailed = NewUserPermissionInfoCodeOnly(constant.ResponseStatusCode_LOGIN_AUTHENTICATION_FAILED)
	LoggingInUserNotActive    = NewUserPermissionInfoCodeOnly(constant.ResponseStatusCode_LOGGING_IN_USER_NOT_ACTIVE)
)
