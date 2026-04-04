package dto

import "im.turms/server/internal/domain/admin/permission"

// AddAdminDTO maps to AddAdminDTO in Java.
// @MappedFrom AddAdminDTO
type AddAdminDTO struct {
	ID        *int64  `json:"id,omitempty"`
	LoginName string  `json:"loginName"`
	Password  string  `json:"password"`
	Name      string  `json:"name"`
	RoleIDs   []int64 `json:"roleIds,omitempty"`
}

// AddAdminRoleDTO maps to AddAdminRoleDTO in Java.
// @MappedFrom AddAdminRoleDTO
type AddAdminRoleDTO struct {
	ID          *int64                       `json:"id,omitempty"`
	Name        string                       `json:"name"`
	Permissions []permission.AdminPermission `json:"permissions,omitempty"`
	Rank        int                          `json:"rank"`
}

// UpdateAdminDTO maps to UpdateAdminDTO in Java.
// @MappedFrom UpdateAdminDTO
type UpdateAdminDTO struct {
	Password *string `json:"password,omitempty"`
	Name     *string `json:"name,omitempty"`
	RoleIDs  []int64 `json:"roleIds,omitempty"`
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
