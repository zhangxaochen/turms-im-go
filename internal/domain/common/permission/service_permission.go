package permission

// ServicePermission maps to ServicePermission.java
// @MappedFrom ServicePermission
type ServicePermission struct {
}

// @MappedFrom ServicePermission(ResponseStatusCode code, String reason)
func NewServicePermission() *ServicePermission {
	return &ServicePermission{}
}
