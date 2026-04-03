package bo

import (
	"im.turms/server/internal/domain/common/constant"
)

// UserPermissionInfo represents the permissions granted to a user session.
type UserPermissionInfo struct {
	AuthenticationCode constant.ResponseStatusCode
	Permissions        map[interface{}]struct{} // Set of TurmsRequest.KindCase
}

// NewUserPermissionInfo creates a new UserPermissionInfo.
// @MappedFrom UserPermissionInfo(ResponseStatusCode authenticationCode, Set<TurmsRequest.KindCase> permissions)
func NewUserPermissionInfo(authenticationCode constant.ResponseStatusCode, permissions map[interface{}]struct{}) *UserPermissionInfo {
	return &UserPermissionInfo{
		AuthenticationCode: authenticationCode,
		Permissions:        permissions,
	}
}
