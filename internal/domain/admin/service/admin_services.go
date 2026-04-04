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
	"sync"
)

var (
	ErrRequesterNotExist = errors.New("requester does not exist")
	ErrPermissionDenied  = errors.New("permission denied")
)

const RootRoleID int64 = 0
const RootAdminID int64 = 0

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
	QueryRoleIdsByAdminId(ctx context.Context, adminId int64) ([]int64, error)
	SetAdminService(adminService AdminService)
}

type adminRoleService struct {
	idGen        *idgen.SnowflakeIdGenerator
	repo         repository.AdminRoleRepository
	adminService AdminService

	idToRole map[int64]*po.AdminRole
	mutex    sync.RWMutex
}

func NewAdminRoleService(idGen *idgen.SnowflakeIdGenerator, repo repository.AdminRoleRepository) AdminRoleService {
	return &adminRoleService{
		idGen:    idGen,
		repo:     repo,
		idToRole: make(map[int64]*po.AdminRole),
	}
}

func (s *adminRoleService) SetAdminService(adminService AdminService) {
	s.adminService = adminService
}

// @MappedFrom authAndAddAdminRole
func (s *adminRoleService) AuthAndAddAdminRole(ctx context.Context, requesterId int64, roleId *int64, name string, permissions []permission.AdminPermission, rank int) (*po.AdminRole, error) {
	if roleId != nil && *roleId == RootRoleID {
		return nil, errors.New("the root role cannot be created")
	}
	if name == "" {
		return nil, errors.New("name must not be blank")
	}
	higher, err := s.IsAdminRankHigherThanRank(ctx, requesterId, rank)
	if err != nil {
		return nil, err
	}
	if !higher {
		return nil, ErrPermissionDenied
	}
	// Verify that the requester has all the requested permissions
	requesterPermissions, err := s.QueryPermissions(ctx, requesterId)
	if err != nil {
		return nil, err
	}
	permMap := make(map[permission.AdminPermission]bool)
	for _, p := range requesterPermissions {
		permMap[p] = true
	}
	for _, p := range permissions {
		if !permMap[p] {
			return nil, ErrPermissionDenied
		}
	}

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
	if roleId == RootRoleID {
		return nil, errors.New("the root role already exists")
	}
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
	s.mutex.Lock()
	s.idToRole[roleId] = role
	s.mutex.Unlock()
	return role, nil
}

// @MappedFrom authAndDeleteAdminRoles
func (s *adminRoleService) AuthAndDeleteAdminRoles(ctx context.Context, requesterId int64, roleIds []int64) (int64, error) {
	for _, id := range roleIds {
		if id == RootRoleID {
			return 0, errors.New("the root role cannot be deleted")
		}
	}
	targetHighest, err := s.QueryHighestRankByRoleIds(ctx, roleIds)
	if err != nil {
		return 0, err
	}
	if targetHighest == nil {
		return 0, nil
	}
	higher, err := s.IsAdminRankHigherThanRank(ctx, requesterId, *targetHighest)
	if err != nil {
		return 0, err
	}
	if !higher {
		return 0, ErrPermissionDenied
	}
	return s.DeleteAdminRoles(ctx, roleIds)
}

// @MappedFrom deleteAdminRoles
func (s *adminRoleService) DeleteAdminRoles(ctx context.Context, roleIds []int64) (int64, error) {
	deleted, err := s.repo.DeleteAdminRoles(ctx, roleIds)
	if err == nil && deleted > 0 {
		s.mutex.Lock()
		for _, id := range roleIds {
			delete(s.idToRole, id)
		}
		s.mutex.Unlock()
	}
	return deleted, err
}

// @MappedFrom authAndUpdateAdminRoles
func (s *adminRoleService) AuthAndUpdateAdminRoles(ctx context.Context, requesterId int64, roleIds []int64, newName *string, permissions []permission.AdminPermission, rank *int) (int64, error) {
	for _, id := range roleIds {
		if id == RootRoleID {
			return 0, errors.New("the root role cannot be updated")
		}
	}
	if rank != nil {
		higher, err := s.IsAdminRankHigherThanRank(ctx, requesterId, *rank)
		if err != nil {
			return 0, err
		}
		if !higher {
			return 0, ErrPermissionDenied
		}
	}
	targetHighest, err := s.QueryHighestRankByRoleIds(ctx, roleIds)
	if err != nil {
		return 0, err
	}
	if targetHighest == nil {
		return 0, nil
	}
	higher, err := s.IsAdminRankHigherThanRank(ctx, requesterId, *targetHighest)
	if err != nil {
		return 0, err
	}
	if !higher {
		return 0, ErrPermissionDenied
	}
	return s.UpdateAdminRole(ctx, roleIds, newName, permissions, rank)
}

// @MappedFrom updateAdminRole
func (s *adminRoleService) UpdateAdminRole(ctx context.Context, roleIds []int64, newName *string, permissions []permission.AdminPermission, rank *int) (int64, error) {
	modified, err := s.repo.UpdateAdminRoles(ctx, roleIds, newName, permissions, rank)
	if err == nil && modified > 0 {
		s.mutex.Lock()
		for _, id := range roleIds {
			delete(s.idToRole, id) // Invalidate cache
		}
		s.mutex.Unlock()
	}
	return modified, err
}

func (s *adminRoleService) QueryRoleIdsByAdminId(ctx context.Context, adminId int64) ([]int64, error) {
	return s.adminService.QueryRoleIdsByAdminIds(ctx, []int64{adminId})
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
	roleIds, err := s.adminService.QueryRoleIdsByAdminIds(ctx, []int64{adminId})
	if err != nil {
		return nil, err
	}
	return s.QueryHighestRankByRoleIds(ctx, roleIds)
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
	roleIds, err := s.adminService.QueryRoleIdsByAdminIds(ctx, []int64{adminId})
	if err != nil {
		return nil, err
	}
	roles, err := s.repo.FindAdminRoles(ctx, roleIds, nil, nil, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	var permissions []permission.AdminPermission
	permMap := make(map[permission.AdminPermission]struct{})
	for _, role := range roles {
		for _, p := range role.Permissions {
			if _, ok := permMap[p]; !ok {
				permMap[p] = struct{}{}
				permissions = append(permissions, p)
			}
		}
	}
	return permissions, nil
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
	if len(roleIds) > 0 {
		targetHighest, err := s.adminRoleService.QueryHighestRankByRoleIds(ctx, roleIds)
		if err != nil {
			return nil, err
		}
		if targetHighest != nil {
			higher, err := s.adminRoleService.IsAdminRankHigherThanRank(ctx, requesterId, *targetHighest)
			if err != nil {
				return nil, err
			}
			if !higher {
				return nil, ErrPermissionDenied
			}
		}
	}
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
	for _, id := range adminIds {
		if id == RootAdminID {
			return 0, errors.New("the root admin cannot be deleted")
		}
	}
	targetRoleIds, err := s.QueryRoleIdsByAdminIds(ctx, adminIds)
	if err != nil {
		return 0, err
	}
	targetHighest, err := s.adminRoleService.QueryHighestRankByRoleIds(ctx, targetRoleIds)
	if err != nil {
		return 0, err
	}
	if targetHighest != nil {
		higher, err := s.adminRoleService.IsAdminRankHigherThanRank(ctx, requesterId, *targetHighest)
		if err != nil {
			return 0, err
		}
		if !higher {
			return 0, ErrPermissionDenied
		}
	}
	return s.repo.DeleteAdmins(ctx, adminIds)
}

func (s *adminService) AuthAndUpdateAdmins(ctx context.Context, requesterId int64, targetAdminIds []int64, rawPassword *string, displayName *string, roleIds []int64) (int64, error) {
	targetCurrentRoleIds, err := s.QueryRoleIdsByAdminIds(ctx, targetAdminIds)
	if err != nil {
		return 0, err
	}
	targetHighest, err := s.adminRoleService.QueryHighestRankByRoleIds(ctx, targetCurrentRoleIds)
	if err != nil {
		return 0, err
	}
	if targetHighest != nil {
		higher, err := s.adminRoleService.IsAdminRankHigherThanRank(ctx, requesterId, *targetHighest)
		if err != nil {
			return 0, err
		}
		if !higher {
			return 0, ErrPermissionDenied
		}
	}
	if len(roleIds) > 0 {
		newTargetHighest, err := s.adminRoleService.QueryHighestRankByRoleIds(ctx, roleIds)
		if err != nil {
			return 0, err
		}
		if newTargetHighest != nil {
			higher, err := s.adminRoleService.IsAdminRankHigherThanRank(ctx, requesterId, *newTargetHighest)
			if err != nil {
				return 0, err
			}
			if !higher {
				return 0, ErrPermissionDenied
			}
		}
	}
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
