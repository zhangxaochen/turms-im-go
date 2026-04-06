package controller

import (
	"context"

	admindto "im.turms/server/internal/domain/admin/access/admin/dto"
	"im.turms/server/internal/domain/admin/permission"
	adminpo "im.turms/server/internal/domain/admin/po"
	"im.turms/server/internal/domain/admin/service"
	commoncontroller "im.turms/server/internal/domain/common/access/admin/controller"
	"im.turms/server/internal/infra/property"
)

// AdminController maps to AdminController Java.
// @MappedFrom AdminController
type AdminController struct {
	*commoncontroller.BaseController
	adminService service.AdminService
}

func NewAdminController(propertiesManager *property.TurmsPropertiesManager, adminService service.AdminService) *AdminController {
	return &AdminController{
		BaseController: commoncontroller.NewBaseController(propertiesManager),
		adminService:   adminService,
	}
}

// @MappedFrom checkLoginNameAndPassword
func (c *AdminController) CheckLoginNameAndPassword() error {
	return nil
}

// @MappedFrom addAdmin
func (c *AdminController) AddAdmin(ctx context.Context, requesterId int64, addAdminDTO admindto.AddAdminDTO) (*adminpo.Admin, error) {
	var displayName string
	if addAdminDTO.DisplayName != nil {
		displayName = *addAdminDTO.DisplayName
	}
	return c.adminService.AuthAndAddAdmin(
		ctx,
		requesterId,
		addAdminDTO.LoginName,
		addAdminDTO.Password,
		displayName,
		addAdminDTO.RoleIDs,
	)
}

// @MappedFrom queryAdmins
func (c *AdminController) QueryAdmins(ctx context.Context, ids []int64, loginNames []string, roleIds []int64, withPassword bool, size *int) ([]*adminpo.Admin, error) {
	actualSize := c.GetPageSize(size)
	page := 0
	admins, err := c.adminService.QueryAdmins(ctx, ids, loginNames, roleIds, &page, &actualSize)
	if err != nil {
		return nil, err
	}
	if !withPassword {
		for i := range admins {
			admins[i].Password = nil
		}
	}
	return admins, nil
}

// @MappedFrom queryAdminsByPage
func (c *AdminController) QueryAdminsByPage(ctx context.Context, ids []int64, loginNames []string, roleIds []int64, withPassword bool, page int, size *int) (int64, []*adminpo.Admin, error) {
	actualSize := c.GetPageSize(size)
	count, err := c.adminService.CountAdmins(ctx, ids, roleIds)
	if err != nil {
		return 0, nil, err
	}
	admins, err := c.adminService.QueryAdmins(ctx, ids, loginNames, roleIds, &page, &actualSize)
	if err != nil {
		return 0, nil, err
	}
	if !withPassword {
		for i := range admins {
			admins[i].Password = nil
		}
	}
	return count, admins, nil
}

// @MappedFrom updateAdmins
func (c *AdminController) UpdateAdmins(ctx context.Context, requesterId int64, ids []int64, updateAdminDTO admindto.UpdateAdminDTO) error {
	_, err := c.adminService.AuthAndUpdateAdmins(
		ctx,
		requesterId,
		ids,
		updateAdminDTO.Password,
		updateAdminDTO.DisplayName,
		updateAdminDTO.RoleIDs,
	)
	return err
}

// @MappedFrom deleteAdmins
func (c *AdminController) DeleteAdmins(ctx context.Context, requesterId int64, ids []int64) error {
	_, err := c.adminService.AuthAndDeleteAdmins(ctx, requesterId, ids)
	return err
}

// AdminPermissionController maps to AdminPermissionController Java.
// @MappedFrom AdminPermissionController
type AdminPermissionController struct {
	*commoncontroller.BaseController
}

func NewAdminPermissionController(propertiesManager *property.TurmsPropertiesManager) *AdminPermissionController {
	return &AdminPermissionController{
		BaseController: commoncontroller.NewBaseController(propertiesManager),
	}
}

// @MappedFrom queryAdminPermissions
func (c *AdminPermissionController) QueryAdminPermissions() []admindto.PermissionDTO {
	allPermissions := permission.AllAdminPermissions
	dtos := make([]admindto.PermissionDTO, len(allPermissions))
	for i, p := range allPermissions {
		dtos[i] = admindto.PermissionDTO{
			Group:      p.Group(),
			Permission: p,
		}
	}
	return dtos
}

// AdminRoleController maps to AdminRoleController Java.
// @MappedFrom AdminRoleController
type AdminRoleController struct {
	*commoncontroller.BaseController
	adminRoleService service.AdminRoleService
}

func NewAdminRoleController(propertiesManager *property.TurmsPropertiesManager, adminRoleService service.AdminRoleService) *AdminRoleController {
	return &AdminRoleController{
		BaseController:   commoncontroller.NewBaseController(propertiesManager),
		adminRoleService: adminRoleService,
	}
}

// @MappedFrom addAdminRole
func (c *AdminRoleController) AddAdminRole(ctx context.Context, requesterId int64, addAdminRoleDTO admindto.AddAdminRoleDTO) (*adminpo.AdminRole, error) {
	var name string
	if addAdminRoleDTO.Name != nil {
		name = *addAdminRoleDTO.Name
	}
	return c.adminRoleService.AuthAndAddAdminRole(
		ctx,
		requesterId,
		addAdminRoleDTO.ID,
		name,
		addAdminRoleDTO.Permissions,
		addAdminRoleDTO.Rank,
	)
}

// @MappedFrom queryAdminRoles
func (c *AdminRoleController) QueryAdminRoles(ctx context.Context, ids []int64, names []string, includedPermissions []permission.AdminPermission, ranks []int, size *int) ([]*adminpo.AdminRole, error) {
	actualSize := c.GetPageSize(size)
	page := 0
	return c.adminRoleService.QueryAdminRoles(ctx, ids, names, includedPermissions, ranks, &page, &actualSize)
}

// @MappedFrom queryAdminRolesByPage
func (c *AdminRoleController) QueryAdminRolesByPage(ctx context.Context, ids []int64, names []string, includedPermissions []permission.AdminPermission, ranks []int, page int, size *int) (int64, []*adminpo.AdminRole, error) {
	actualSize := c.GetPageSize(size)
	count, err := c.adminRoleService.CountAdminRoles(ctx, ids, names, includedPermissions, ranks)
	if err != nil {
		return 0, nil, err
	}
	roles, err := c.adminRoleService.QueryAdminRoles(ctx, ids, names, includedPermissions, ranks, &page, &actualSize)
	return count, roles, err
}

// @MappedFrom updateAdminRole
func (c *AdminRoleController) UpdateAdminRole(ctx context.Context, requesterId int64, ids []int64, updateAdminRoleDTO admindto.UpdateAdminRoleDTO) error {
	_, err := c.adminRoleService.AuthAndUpdateAdminRoles(
		ctx,
		requesterId,
		ids,
		updateAdminRoleDTO.Name,
		updateAdminRoleDTO.Permissions,
		updateAdminRoleDTO.Rank,
	)
	return err
}

// @MappedFrom deleteAdminRoles
func (c *AdminRoleController) DeleteAdminRoles(ctx context.Context, requesterId int64, ids []int64) error {
	_, err := c.adminRoleService.AuthAndDeleteAdminRoles(ctx, requesterId, ids)
	return err
}
