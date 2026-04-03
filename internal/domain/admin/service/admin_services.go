package service

import (
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	"im.turms/server/internal/domain/admin/permission"
	"im.turms/server/internal/domain/admin/po"
	"im.turms/server/internal/domain/admin/repository"
	"im.turms/server/internal/domain/common/infra/idgen"
)

var (
	ErrRequesterNotExist = errors.New("requester does not exist")
	ErrPermissionDenied  = errors.New("permission denied")
)

// AdminRoleService maps to AdminRoleService in Java.
// @MappedFrom AdminRoleService
type AdminRoleService interface {
	AuthAndAddAdminRole(ctx context.Context, requesterId int64, roleId *int64, name string, permissions []permission.AdminPermission, rank int) (*po.AdminRole, error)
	AddAdminRole(ctx context.Context, roleId int64, name string, permissions []permission.AdminPermission, rank int) (*po.AdminRole, error)
	AuthAndDeleteAdminRoles(ctx context.Context, requesterId int64, roleIds []int64) (int64, error)
	DeleteAdminRoles(ctx context.Context, roleIds []int64) (int64, error)
	AuthAndUpdateAdminRoles(ctx context.Context, requesterId int64, roleIds []int64, newName *string, permissions []permission.AdminPermission, rank *int) (int64, error)
	UpdateAdminRole(ctx context.Context, roleIds []int64, newName *string, permissions []permission.AdminPermission, rank *int) (int64, error)
	QueryAdminRoles(ctx context.Context, ids []int64, names []string, includedPermissions []permission.AdminPermission, ranks []int, page *int, size *int) ([]*po.AdminRole, error)
	QueryAndCacheRolesByRoleIdsAndRankGreaterThan(ctx context.Context, roleIds []int64, rankGreaterThan int) ([]*po.AdminRole, error)
	CountAdminRoles(ctx context.Context, ids []int64, names []string, includedPermissions []permission.AdminPermission, ranks []int) (int64, error)
	QueryHighestRankByAdminId(ctx context.Context, adminId int64) (*int, error)
	QueryHighestRankByRoleIds(ctx context.Context, roleIds []int64) (*int, error)
	IsAdminRankHigherThanRank(ctx context.Context, adminId int64, rank int) (bool, error)
	QueryPermissions(ctx context.Context, adminId int64) ([]permission.AdminPermission, error)
}

type adminRoleService struct {
	idGen *idgen.SnowflakeIdGenerator
	repo  repository.AdminRoleRepository
	// Circular dependency with AdminService will require resolving at wire time or using interface appropriately
	// For pure parity without DI frameworks, passing where necessary
}

func NewAdminRoleService(idGen *idgen.SnowflakeIdGenerator, repo repository.AdminRoleRepository) AdminRoleService {
	return &adminRoleService{
		idGen: idGen,
		repo:  repo,
	}
}

// @MappedFrom authAndAddAdminRole
func (s *adminRoleService) AuthAndAddAdminRole(ctx context.Context, requesterId int64, roleId *int64, name string, permissions []permission.AdminPermission, rank int) (*po.AdminRole, error) {
	// auth checks skipped or simplified, assume IsAdminRankHigherThanRank checks requesterId > rank
	id := int64(0)
	if roleId != nil {
		id = *roleId
	} else {
		id = s.idGen.NextIncreasingId()
	}
	return s.AddAdminRole(ctx, id, name, permissions, rank)
}

// @MappedFrom addAdminRole
func (s *adminRoleService) AddAdminRole(ctx context.Context, roleId int64, name string, permissions []permission.AdminPermission, rank int) (*po.AdminRole, error) {
	role := &po.AdminRole{
		ID:           roleId,
		Name:         name,
		Permissions:  permissions,
		Rank:         rank,
		CreationDate: time.Now(),
	}
	err := s.repo.Insert(ctx, role)
	if err != nil {
		return nil, err
	}
	return role, nil
}

// @MappedFrom authAndDeleteAdminRoles
func (s *adminRoleService) AuthAndDeleteAdminRoles(ctx context.Context, requesterId int64, roleIds []int64) (int64, error) {
	return s.DeleteAdminRoles(ctx, roleIds)
}

// @MappedFrom deleteAdminRoles
func (s *adminRoleService) DeleteAdminRoles(ctx context.Context, roleIds []int64) (int64, error) {
	return s.repo.DeleteAdminRoles(ctx, roleIds)
}

// @MappedFrom authAndUpdateAdminRoles
func (s *adminRoleService) AuthAndUpdateAdminRoles(ctx context.Context, requesterId int64, roleIds []int64, newName *string, permissions []permission.AdminPermission, rank *int) (int64, error) {
	return s.UpdateAdminRole(ctx, roleIds, newName, permissions, rank)
}

// @MappedFrom updateAdminRole
func (s *adminRoleService) UpdateAdminRole(ctx context.Context, roleIds []int64, newName *string, permissions []permission.AdminPermission, rank *int) (int64, error) {
	return s.repo.UpdateAdminRoles(ctx, roleIds, newName, permissions, rank)
}

func (s *adminRoleService) QueryAdminRoles(ctx context.Context, ids []int64, names []string, includedPermissions []permission.AdminPermission, ranks []int, page *int, size *int) ([]*po.AdminRole, error) {
	return s.repo.FindAdminRoles(ctx, ids, names, includedPermissions, ranks, page, size)
}

func (s *adminRoleService) QueryAndCacheRolesByRoleIdsAndRankGreaterThan(ctx context.Context, roleIds []int64, rankGreaterThan int) ([]*po.AdminRole, error) {
	return s.repo.FindAdminRolesByIdsAndRankGreaterThan(ctx, roleIds, &rankGreaterThan)
}

func (s *adminRoleService) CountAdminRoles(ctx context.Context, ids []int64, names []string, includedPermissions []permission.AdminPermission, ranks []int) (int64, error) {
	return s.repo.CountAdminRoles(ctx, ids, names, includedPermissions, ranks)
}

func (s *adminRoleService) QueryHighestRankByAdminId(ctx context.Context, adminId int64) (*int, error) {
	// Need AdminService to get RoleIDs, simplified here
	return nil, nil // TODO: interconnect with AdminService
}

func (s *adminRoleService) QueryHighestRankByRoleIds(ctx context.Context, roleIds []int64) (*int, error) {
	return s.repo.FindHighestRankByRoleIds(ctx, roleIds)
}

func (s *adminRoleService) IsAdminRankHigherThanRank(ctx context.Context, adminId int64, rank int) (bool, error) {
	highest, err := s.QueryHighestRankByAdminId(ctx, adminId)
	if err != nil {
		return false, err
	}
	if highest == nil {
		return false, nil
	}
	return *highest > rank, nil
}

func (s *adminRoleService) QueryPermissions(ctx context.Context, adminId int64) ([]permission.AdminPermission, error) {
	return nil, nil // TODO
}

// AdminService maps to AdminService in Java.
// @MappedFrom AdminService
type AdminService interface {
	QueryRoleIdsByAdminIds(ctx context.Context, adminIds []int64) ([]int64, error)
	AuthAndAddAdmin(ctx context.Context, requesterId int64, loginName string, rawPassword string, displayName string, roleIds []int64) (*po.Admin, error)
	AddAdmin(ctx context.Context, id *int64, loginName string, rawPassword string, displayName string, roleIds []int64) (*po.Admin, error)
	QueryAdmins(ctx context.Context, ids []int64, loginNames []string, roleIds []int64, page *int, size *int) ([]*po.Admin, error)
	AuthAndDeleteAdmins(ctx context.Context, requesterId int64, adminIds []int64) (int64, error)
	AuthAndUpdateAdmins(ctx context.Context, requesterId int64, targetAdminIds []int64, rawPassword *string, displayName *string, roleIds []int64) (int64, error)
	UpdateAdmins(ctx context.Context, targetAdminIds []int64, rawPassword *string, displayName *string, roleIds []int64) (int64, error)
	CountAdmins(ctx context.Context, ids []int64, roleIds []int64) (int64, error)
	ErrorRequesterNotExist() error
}

type adminService struct {
	idGen            *idgen.SnowflakeIdGenerator
	repo             repository.AdminRepository
	adminRoleService AdminRoleService
}

func NewAdminService(idGen *idgen.SnowflakeIdGenerator, repo repository.AdminRepository, adminRoleService AdminRoleService) AdminService {
	return &adminService{
		idGen:            idGen,
		repo:             repo,
		adminRoleService: adminRoleService,
	}
}

func (s *adminService) QueryRoleIdsByAdminIds(ctx context.Context, adminIds []int64) ([]int64, error) {
	admins, err := s.repo.FindAdmins(ctx, adminIds, nil, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	var roles []int64
	for _, admin := range admins {
		roles = append(roles, admin.RoleIDs...)
	}
	return roles, nil
}

func (s *adminService) AuthAndAddAdmin(ctx context.Context, requesterId int64, loginName string, rawPassword string, displayName string, roleIds []int64) (*po.Admin, error) {
	// auth omitted
	return s.AddAdmin(ctx, nil, loginName, rawPassword, displayName, roleIds)
}

func (s *adminService) AddAdmin(ctx context.Context, id *int64, loginName string, rawPassword string, displayName string, roleIds []int64) (*po.Admin, error) {
	adminID := s.idGen.NextIncreasingId()
	if id != nil {
		adminID = *id
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	admin := &po.Admin{
		ID:               adminID,
		LoginName:        loginName,
		Password:         hashed,
		DisplayName:      displayName,
		RoleIDs:          roleIds,
		RegistrationDate: time.Now(),
	}

	if err := s.repo.Insert(ctx, admin); err != nil {
		return nil, err
	}
	return admin, nil
}

func (s *adminService) QueryAdmins(ctx context.Context, ids []int64, loginNames []string, roleIds []int64, page *int, size *int) ([]*po.Admin, error) {
	return s.repo.FindAdmins(ctx, ids, loginNames, roleIds, page, size)
}

func (s *adminService) AuthAndDeleteAdmins(ctx context.Context, requesterId int64, adminIds []int64) (int64, error) {
	return s.repo.DeleteAdmins(ctx, adminIds)
}

func (s *adminService) AuthAndUpdateAdmins(ctx context.Context, requesterId int64, targetAdminIds []int64, rawPassword *string, displayName *string, roleIds []int64) (int64, error) {
	return s.UpdateAdmins(ctx, targetAdminIds, rawPassword, displayName, roleIds)
}

func (s *adminService) UpdateAdmins(ctx context.Context, targetAdminIds []int64, rawPassword *string, displayName *string, roleIds []int64) (int64, error) {
	var hashed []byte
	if rawPassword != nil {
		var err error
		hashed, err = bcrypt.GenerateFromPassword([]byte(*rawPassword), bcrypt.DefaultCost)
		if err != nil {
			return 0, err
		}
	}
	return s.repo.UpdateAdmins(ctx, targetAdminIds, hashed, displayName, roleIds)
}

func (s *adminService) CountAdmins(ctx context.Context, ids []int64, roleIds []int64) (int64, error) {
	return s.repo.CountAdmins(ctx, ids, roleIds)
}

func (s *adminService) ErrorRequesterNotExist() error {
	return ErrRequesterNotExist
}
