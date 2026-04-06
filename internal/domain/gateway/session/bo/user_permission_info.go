package bo

import (
	"im.turms/server/internal/domain/common/constant"
	"im.turms/server/internal/domain/gateway/access/client/common/authorization"
)

// UserPermissionInfo represents the permissions granted to a user session.
type UserPermissionInfo struct {
	AuthenticationCode constant.ResponseStatusCode
	Permissions        map[int32]bool // Set of TurmsRequest tag numbers
}

// NewUserPermissionInfo creates a new UserPermissionInfo.
// @MappedFrom UserPermissionInfo(ResponseStatusCode authenticationCode, Set<TurmsRequest.KindCase> permissions)
func NewUserPermissionInfo(authenticationCode constant.ResponseStatusCode, permissions map[int32]bool) *UserPermissionInfo {
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
		Permissions:        make(map[int32]bool),
	}
}

// allPermissionsMap holds a pre-built permissions map with all request types.
// Java parity: TurmsRequestTypePool.ALL (populated set of all request types).
var allPermissionsMap map[int32]bool

func init() {
	allPermissionsMap = make(map[int32]bool, len(authorization.ALL_REQUEST_TYPES))
	for _, rt := range authorization.ALL_REQUEST_TYPES {
		allPermissionsMap[rt] = true
	}
}

// GrantedWithAllPermissions returns a UserPermissionInfo with ALL request types granted.
// Java parity: GRANTED_WITH_ALL_PERMISSIONS uses TurmsRequestTypePool.ALL (populated set),
// not nil/empty. Downstream code checks permissions[requestType] and expects true for all.
var GrantedWithAllPermissions = NewUserPermissionInfo(constant.ResponseStatusCode_OK, allPermissionsMap)

var (
	LoginAuthenticationFailed = NewUserPermissionInfoCodeOnly(constant.ResponseStatusCode_LOGIN_AUTHENTICATION_FAILED)
	LoggingInUserNotActive    = NewUserPermissionInfoCodeOnly(constant.ResponseStatusCode_LOGGING_IN_USER_NOT_ACTIVE)
)
