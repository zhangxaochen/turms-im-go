package permission

import "im.turms/server/internal/domain/common/constant"

// ServicePermission maps to ServicePermission.java
// @MappedFrom ServicePermission
type ServicePermission struct {
	Code   constant.ResponseStatusCode
	Reason string
}

// @MappedFrom ServicePermission(ResponseStatusCode code, String reason)
func NewServicePermission(code constant.ResponseStatusCode, reason string) *ServicePermission {
	return &ServicePermission{
		Code:   code,
		Reason: reason,
	}
}

// OK defines the OK service permission constant.
var OK = NewServicePermission(constant.ResponseStatusCode_OK, "")

// Get returns a service permission for a code with no reason.
// @MappedFrom get(ResponseStatusCode code)
func Get(code constant.ResponseStatusCode) *ServicePermission {
	if code == constant.ResponseStatusCode_OK {
		return OK
	}
	return NewServicePermission(code, "")
}
