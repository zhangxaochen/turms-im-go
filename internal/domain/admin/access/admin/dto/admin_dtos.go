package dto

import (
	"encoding/json"
	"fmt"

	"im.turms/server/internal/domain/admin/permission"
)

// AddAdminDTO maps to AddAdminDTO in Java.
// @MappedFrom AddAdminDTO
// Note: Java record fields are loginName, password, displayName, roleIds (no id field).
type AddAdminDTO struct {
	LoginName   string  `json:"loginName"`
	Password    string  `json:"password"`
	DisplayName *string `json:"displayName,omitempty"` // Java: displayName
	RoleIDs     []int64 `json:"roleIds,omitempty"`
}

// String masks password, matching Java's @SensitiveProperty toString().
// @MappedFrom AddAdminDTO.toString()
func (d AddAdminDTO) String() string {
	return fmt.Sprintf("AddAdminDTO[loginName=%s, password=***, displayName=%v, roleIds=%v]",
		d.LoginName, d.DisplayName, d.RoleIDs)
}

// AddAdminRoleDTO maps to AddAdminRoleDTO in Java.
// @MappedFrom AddAdminRoleDTO
// Note: rank is Integer (nullable) in Java → *int in Go.
type AddAdminRoleDTO struct {
	ID          *int64                       `json:"id,omitempty"`
	Name        *string                      `json:"name,omitempty"`
	Permissions []permission.AdminPermission `json:"permissions,omitempty"`
	Rank        *int                         `json:"rank,omitempty"` // nullable Integer in Java
}

// UpdateAdminDTO maps to UpdateAdminDTO in Java.
// @MappedFrom UpdateAdminDTO
// Password uses write-only semantics: can be deserialized but never serialized in responses.
type UpdateAdminDTO struct {
	// password is unexported from JSON responses via custom MarshalJSON.
	Password *string `json:"password,omitempty"`
	Name     *string `json:"name,omitempty"`
	RoleIDs  []int64 `json:"roleIds,omitempty"`
}

// MarshalJSON prevents password from being serialized in responses.
// @MappedFrom @SensitiveProperty(ALLOW_DESERIALIZATION) on password field.
func (d UpdateAdminDTO) MarshalJSON() ([]byte, error) {
	type Alias struct {
		Name    *string `json:"name,omitempty"`
		RoleIDs []int64 `json:"roleIds,omitempty"`
	}
	return json.Marshal(Alias{
		Name:    d.Name,
		RoleIDs: d.RoleIDs,
	})
}

// String masks password to avoid leaking it in logs.
// @MappedFrom UpdateAdminDTO.toString()
func (d UpdateAdminDTO) String() string {
	var passwordStr string
	if d.Password != nil {
		passwordStr = "***"
	} else {
		passwordStr = "null"
	}
	return fmt.Sprintf("UpdateAdminDTO[password=%s, name=%v, roleIds=%v]",
		passwordStr, d.Name, d.RoleIDs)
}

// UpdateAdminRoleDTO maps to UpdateAdminRoleDTO in Java.
// @MappedFrom UpdateAdminRoleDTO
type UpdateAdminRoleDTO struct {
	Name        *string                      `json:"name,omitempty"`
	Permissions []permission.AdminPermission `json:"permissions,omitempty"`
	Rank        *int                         `json:"rank,omitempty"`
}

// PermissionDTO maps to PermissionDTO in Java.
// @MappedFrom PermissionDTO
type PermissionDTO struct {
	Group      string                      `json:"group"`
	Permission permission.AdminPermission `json:"permission"`
}
