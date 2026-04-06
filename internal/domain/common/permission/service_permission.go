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

// Get returns a new service permission for a code with no reason.
// Bug fix: Always creates a new instance (Java always returns new ServicePermission(code, null)),
// never returns the shared OK singleton, to match Java behavior where get(OK) != OK.
// @MappedFrom get(ResponseStatusCode code)
func Get(code constant.ResponseStatusCode) *ServicePermission {
	return NewServicePermission(code, "")
}

// GetWithReason returns a new service permission for a code with a reason.
// Bug fix: Added missing two-parameter get method from Java.
// @MappedFrom get(ResponseStatusCode code, String reason)
func GetWithReason(code constant.ResponseStatusCode, reason string) *ServicePermission {
	return NewServicePermission(code, reason)
}
