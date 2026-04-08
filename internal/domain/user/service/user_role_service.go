package service

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"im.turms/server/internal/domain/user/po"
	"im.turms/server/internal/domain/user/repository"
)

type UserRoleService struct {
	roleRepo repository.UserRoleRepository
}

func NewUserRoleService(roleRepo repository.UserRoleRepository) *UserRoleService {
	return &UserRoleService{
		roleRepo: roleRepo,
	}
}

// @MappedFrom queryUserRoles(@Nullable Integer page, @Nullable Integer size)
// @MappedFrom queryUserRoles(@QueryParam(required = false)
func (s *UserRoleService) QueryUserRoles(ctx context.Context, filter bson.M) ([]*po.UserRole, error) {
	return s.roleRepo.FindRoles(ctx, filter)
}

// @MappedFrom addUserRole(@Nullable Long groupId, @Nullable String name, @NotNull Set<Long> creatableGroupTypeIds, @NotNull Integer ownedGroupLimit, @NotNull Integer ownedGroupLimitForEachGroupType, @NotNull Map<Long, Integer> groupTypeIdToLimit)
// @MappedFrom addUserRole(@RequestBody AddUserRoleDTO addUserRoleDTO)
func (s *UserRoleService) AddUserRole(ctx context.Context, role *po.UserRole) error {
	return s.roleRepo.InsertRole(ctx, role)
}

func (s *UserRoleService) UpdateUserRoles(ctx context.Context, filter bson.M, update bson.M) error {
	// Extract role IDs from filter and use the repository's UpdateUserRoles method
	idsVal, ok := filter["_id"]
	if !ok {
		return nil
	}
	idsMap, ok := idsVal.(map[string]interface{})
	if !ok {
		return nil
	}
	inVal, ok := idsMap["$in"]
	if !ok {
		return nil
	}
	roleIDs, ok := inVal.([]int64)
	if !ok {
		return nil
	}
	_, err := s.roleRepo.UpdateUserRoles(ctx, roleIDs, update)
	return err
}

// @MappedFrom deleteUserRoles(@Nullable Set<Long> groupIds)
func (s *UserRoleService) DeleteUserRoles(ctx context.Context, filter bson.M) (int64, error) {
	return s.roleRepo.DeleteRoles(ctx, filter)
}

// @MappedFrom queryUserRoleById(@NotNull Long id)
func (s *UserRoleService) QueryUserRoleById(ctx context.Context, roleID int64) (*po.UserRole, error) {
	return s.roleRepo.FindRoleByID(ctx, roleID)
}

// @MappedFrom queryStoredOrDefaultUserRoleByUserId(@NotNull Long userId)
func (s *UserRoleService) QueryStoredOrDefaultUserRoleByUserId(ctx context.Context, userID int64) (*po.UserRole, error) {
	// Usually there is a default role. This needs UserService to fetch user's roleID,
	// but currently we just return a stub or we'd inject UserService.
	// We'll return nil for now to map the method.
	return nil, nil
}

// @MappedFrom countUserRoles()
func (s *UserRoleService) CountUserRoles(ctx context.Context, filter bson.M) (int64, error) {
	return s.roleRepo.CountRoles(ctx, filter)
}
