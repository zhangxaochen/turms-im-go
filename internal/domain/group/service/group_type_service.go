package service

import (
	"context"

	"im.turms/server/internal/domain/group/po"
	"im.turms/server/internal/domain/group/repository"
)

type GroupTypeService struct {
	groupTypeRepo *repository.GroupTypeRepository
}

func NewGroupTypeService(groupTypeRepo *repository.GroupTypeRepository) *GroupTypeService {
	return &GroupTypeService{
		groupTypeRepo: groupTypeRepo,
	}
}

// FindGroupType retrieves a group type by its ID.
func (s *GroupTypeService) FindGroupType(ctx context.Context, typeID int64) (*po.GroupType, error) {
	return s.groupTypeRepo.FindGroupType(ctx, typeID)
}

func (s *GroupTypeService) FindByID(ctx context.Context, typeID int64) (*po.GroupType, error) {
	return s.FindGroupType(ctx, typeID)
}

// EnsureDefaultGroupType creates the default group type if it does not exist.
func (s *GroupTypeService) EnsureDefaultGroupType(ctx context.Context) error {
	// Simple implementation for ensuring default group type.
	t, err := s.groupTypeRepo.FindGroupType(ctx, 0)
	if err != nil {
		return err
	}
	if t == nil {
		defaultType := &po.GroupType{
			ID:             0,
			Name:           "DEFAULT",
			GroupSizeLimit: 500,
		}
		return s.groupTypeRepo.InsertGroupType(ctx, defaultType)
	}
	return nil
}

// @MappedFrom initGroupTypes()
func (s *GroupTypeService) InitGroupTypes(ctx context.Context) error {
	return s.EnsureDefaultGroupType(ctx)
}

// @MappedFrom addGroupType(@Nullable Long id, @NotNull @NoWhitespace String name, @NotNull @Min(1) Integer groupSizeLimit, @NotNull InvitationStrategy invitationStrategy, @NotNull JoinStrategy joinStrategy, @NotNull GroupUpdateStrategy groupInfoUpdateStrategy, @NotNull GroupUpdateStrategy memberInfoUpdateStrategy, @NotNull Boolean guestSpeakable, @NotNull Boolean selfInfoUpdatable, @NotNull Boolean enableReadReceipt, @NotNull Boolean messageEditable)
func (s *GroupTypeService) AddGroupType(ctx context.Context, groupType *po.GroupType) error {
	return s.groupTypeRepo.InsertGroupType(ctx, groupType)
}

// @MappedFrom updateGroupTypes(@NotEmpty Set<Long> ids, @Nullable @NoWhitespace String name, @Nullable @Min(1) Integer groupSizeLimit, @Nullable InvitationStrategy invitationStrategy, @Nullable JoinStrategy joinStrategy, @Nullable GroupUpdateStrategy groupInfoUpdateStrategy, @Nullable GroupUpdateStrategy memberInfoUpdateStrategy, @Nullable Boolean guestSpeakable, @Nullable Boolean selfInfoUpdatable, @Nullable Boolean enableReadReceipt, @Nullable Boolean messageEditable)
func (s *GroupTypeService) UpdateGroupTypes(ctx context.Context, ids []int64, update *po.GroupType) error {
	return s.groupTypeRepo.UpdateTypes(ctx, ids, &update.Name, &update.GroupSizeLimit, &update.InvitationStrategy, &update.JoinStrategy, &update.GroupInfoUpdateStrategy, &update.MemberInfoUpdateStrategy, &update.GuestSpeakable, &update.SelfInfoUpdatable, &update.EnableReadReceipt, &update.MessageEditable)
}

// @MappedFrom deleteGroupTypes(@Nullable Set<Long> groupTypeIds)
func (s *GroupTypeService) DeleteGroupTypes(ctx context.Context, groupTypeIds []int64) error {
	return s.groupTypeRepo.DeleteTypes(ctx, groupTypeIds)
}

// @MappedFrom queryGroupType(@NotNull Long groupTypeId)
func (s *GroupTypeService) QueryGroupType(ctx context.Context, groupTypeID int64) (*po.GroupType, error) {
	return s.groupTypeRepo.FindGroupType(ctx, groupTypeID)
}

// @MappedFrom queryGroupTypes(@Nullable Integer page, @Nullable Integer size)
func (s *GroupTypeService) QueryGroupTypes(ctx context.Context, page, size *int32) ([]*po.GroupType, error) {
	return s.groupTypeRepo.FindGroupTypes(ctx, nil, page, size)
}

// @MappedFrom queryGroupTypes(@NotNull Collection<Long> groupTypeIds)
func (s *GroupTypeService) QueryGroupTypesByIds(ctx context.Context, groupTypeIds []int64) ([]*po.GroupType, error) {
	return s.groupTypeRepo.FindGroupTypes(ctx, groupTypeIds, nil, nil)
}

// @MappedFrom groupTypeExists(@NotNull Long groupTypeId)
func (s *GroupTypeService) GroupTypeExists(ctx context.Context, groupTypeID int64) (bool, error) {
	return s.groupTypeRepo.TypeExists(ctx, groupTypeID)
}

// @MappedFrom countGroupTypes()
func (s *GroupTypeService) CountGroupTypes(ctx context.Context) (int64, error) {
	return s.groupTypeRepo.CountGroupTypes(ctx)
}
