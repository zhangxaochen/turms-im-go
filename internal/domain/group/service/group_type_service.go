package service

import (
	"context"
	"log"
	"strings"
	"sync"

	"im.turms/server/internal/domain/common/constant"
	group_constant "im.turms/server/internal/domain/group/constant"
	"im.turms/server/internal/domain/group/po"
	"im.turms/server/internal/domain/group/repository"
	"im.turms/server/internal/infra/exception"
)

const DefaultGroupTypeID int64 = 0

type GroupTypeService struct {
	groupTypeRepo *repository.GroupTypeRepository
	// Bug fix: in-memory cache matching Java's ConcurrentHashMap<Long, GroupType>
	idToGroupType sync.Map
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
	t, err := s.groupTypeRepo.FindGroupType(ctx, 0)
	if err != nil {
		return err
	}
	if t == nil {
		// Bug fix: populate all 11 fields matching Java's default group type.
		// Java: ID=0, name="DEFAULT", groupSizeLimit=500,
		// invitationStrategy=OWNER_MANAGER_MEMBER_REQUIRING_APPROVAL,
		// joinStrategy=INVITATION, groupInfoUpdateStrategy=OWNER_MANAGER,
		// memberInfoUpdateStrategy=OWNER_MANAGER, guestSpeakable=false,
		// selfInfoUpdatable=true, enableReadReceipt=true, messageEditable=true
		defaultType := &po.GroupType{
			ID:                       0,
			Name:                     "DEFAULT",
			GroupSizeLimit:           500,
			InvitationStrategy:       group_constant.GroupInvitationStrategy_OWNER_MANAGER_MEMBER_REQUIRING_APPROVAL,
			JoinStrategy:             group_constant.GroupJoinStrategy_INVITATION,
			GroupInfoUpdateStrategy:  group_constant.GroupUpdateStrategy_OWNER_MANAGER,
			MemberInfoUpdateStrategy: group_constant.GroupUpdateStrategy_OWNER_MANAGER,
			GuestSpeakable:           false,
			SelfInfoUpdatable:        true,
			EnableReadReceipt:        true,
			MessageEditable:          true,
		}
		err = s.groupTypeRepo.InsertGroupType(ctx, defaultType)
		// Bug fix: silently ignore duplicate key errors (Java parity).
		// Java calls addGroupType(...).onErrorComplete(DuplicateKeyException.class)
		if err != nil && !exception.IsDuplicateKeyError(err) {
			return err
		}
		// Cache the default type
		s.idToGroupType.Store(int64(0), defaultType)
	}
	return nil
}

// @MappedFrom initGroupTypes()
func (s *GroupTypeService) InitGroupTypes(ctx context.Context) error {
	err := s.EnsureDefaultGroupType(ctx)
	if err != nil {
		return err
	}
	// Bug fix: populate in-memory cache from database.
	// Java loads all group types into idToGroupType cache during init.
	types, err := s.groupTypeRepo.FindGroupTypes(ctx, nil, nil, nil)
	if err != nil {
		return err
	}
	for _, gt := range types {
		s.idToGroupType.Store(gt.ID, gt)
	}
	// Note: Java also sets up a MongoDB change stream watcher to keep cache in sync.
	// This is not implemented in Go yet. A TODO comment is left for future work.
	// TODO: implement MongoDB change stream watcher for cache invalidation.
	return nil
}

// @MappedFrom addGroupType(@Nullable Long id, @NotNull @NoWhitespace String name, @NotNull @Min(1) Integer groupSizeLimit, @NotNull InvitationStrategy invitationStrategy, @NotNull JoinStrategy joinStrategy, @NotNull GroupUpdateStrategy groupInfoUpdateStrategy, @NotNull GroupUpdateStrategy memberInfoUpdateStrategy, @NotNull Boolean guestSpeakable, @NotNull Boolean selfInfoUpdatable, @NotNull Boolean enableReadReceipt, @NotNull Boolean messageEditable)
func (s *GroupTypeService) AddGroupType(ctx context.Context, groupType *po.GroupType) error {
	// Bug fix: input validation matching Java.
	if strings.TrimSpace(groupType.Name) == "" {
		return exception.NewTurmsError(int32(constant.ResponseStatusCode_ILLEGAL_ARGUMENT), "name must not be null or contain only whitespace")
	}
	if groupType.GroupSizeLimit < 1 {
		return exception.NewTurmsError(int32(constant.ResponseStatusCode_ILLEGAL_ARGUMENT), "groupSizeLimit must be >= 1")
	}

	err := s.groupTypeRepo.InsertGroupType(ctx, groupType)
	if err != nil {
		return err
	}

	// Bug fix: update in-memory cache on success.
	s.idToGroupType.Store(groupType.ID, groupType)
	return nil
}

// @MappedFrom updateGroupTypes(@NotEmpty Set<Long> ids, @Nullable @NoWhitespace String name, @Nullable @Min(1) Integer groupSizeLimit, @Nullable InvitationStrategy invitationStrategy, @Nullable JoinStrategy joinStrategy, @Nullable GroupUpdateStrategy groupInfoUpdateStrategy, @Nullable GroupUpdateStrategy memberInfoUpdateStrategy, @Nullable Boolean guestSpeakable, @Nullable Boolean selfInfoUpdatable, @Nullable Boolean enableReadReceipt, @Nullable Boolean messageEditable)
func (s *GroupTypeService) UpdateGroupTypes(ctx context.Context, ids []int64, update *po.GroupType) error {
	return s.UpdateGroupTypesWithPointers(ctx, ids, &update.Name, &update.GroupSizeLimit, &update.InvitationStrategy, &update.JoinStrategy, &update.GroupInfoUpdateStrategy, &update.MemberInfoUpdateStrategy, &update.GuestSpeakable, &update.SelfInfoUpdatable, &update.EnableReadReceipt, &update.MessageEditable)
}

// @MappedFrom deleteGroupTypes(@Nullable Set<Long> groupTypeIds)
func (s *GroupTypeService) DeleteGroupTypes(ctx context.Context, groupTypeIds []int64) error {
	// Bug fix: default group type deletion protection.
	// Java explicitly checks if groupTypeIds contains DEFAULT_GROUP_TYPE_ID and throws an error.
	for _, id := range groupTypeIds {
		if id == DefaultGroupTypeID {
			return exception.NewTurmsError(int32(constant.ResponseStatusCode_ILLEGAL_ARGUMENT), "Cannot delete the default group type")
		}
	}

	// Bug fix: handle nil/null groupTypeIds.
	// Java handles null by deleting all non-default types.
	if len(groupTypeIds) == 0 {
		// Delete all types except the default
		allTypes, err := s.groupTypeRepo.FindGroupTypes(ctx, nil, nil, nil)
		if err != nil {
			return err
		}
		var nonDefaultIDs []int64
		for _, gt := range allTypes {
			if gt.ID != DefaultGroupTypeID {
				nonDefaultIDs = append(nonDefaultIDs, gt.ID)
			}
		}
		if len(nonDefaultIDs) == 0 {
			return nil
		}
		groupTypeIds = nonDefaultIDs
	}

	err := s.groupTypeRepo.DeleteTypes(ctx, groupTypeIds)
	if err != nil {
		return err
	}

	// Bug fix: invalidate in-memory cache for deleted types.
	for _, id := range groupTypeIds {
		s.idToGroupType.Delete(id)
	}
	return nil
}

// @MappedFrom queryGroupType(@NotNull Long groupTypeId)
func (s *GroupTypeService) QueryGroupType(ctx context.Context, groupTypeID int64) (*po.GroupType, error) {
	// Bug fix: check in-memory cache first, then fall back to database.
	if cached, ok := s.idToGroupType.Load(groupTypeID); ok {
		return cached.(*po.GroupType), nil
	}
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
	// Bug fix: check in-memory cache first.
	if _, ok := s.idToGroupType.Load(groupTypeID); ok {
		return true, nil
	}
	return s.groupTypeRepo.TypeExists(ctx, groupTypeID)
}

// @MappedFrom countGroupTypes()
func (s *GroupTypeService) CountGroupTypes(ctx context.Context) (int64, error) {
	return s.groupTypeRepo.CountGroupTypes(ctx)
}

// UpdateGroupTypesWithPointers passes nullable field pointers directly to the repository.
// Unlike UpdateGroupTypes which takes a *po.GroupType and always sets all fields (even zero values),
// this method only updates fields that are explicitly provided (non-nil).
func (s *GroupTypeService) UpdateGroupTypesWithPointers(ctx context.Context, ids []int64, name *string, groupSizeLimit *int32, invitationStrategy *group_constant.GroupInvitationStrategy, joinStrategy *group_constant.GroupJoinStrategy, groupInfoUpdateStrategy *group_constant.GroupUpdateStrategy, memberInfoUpdateStrategy *group_constant.GroupUpdateStrategy, guestSpeakable *bool, selfInfoUpdatable *bool, enableReadReceipt *bool, messageEditable *bool) error {
	// Bug fix: input validation matching Java.
	if len(ids) == 0 {
		return nil
	}
	if name != nil && strings.TrimSpace(*name) == "" {
		return exception.NewTurmsError(int32(constant.ResponseStatusCode_ILLEGAL_ARGUMENT), "name must not contain only whitespace")
	}
	if groupSizeLimit != nil && *groupSizeLimit < 1 {
		return exception.NewTurmsError(int32(constant.ResponseStatusCode_ILLEGAL_ARGUMENT), "groupSizeLimit must be >= 1")
	}

	// Bug fix: "all null" short-circuit matching Java.
	// Java checks if all update parameters are null and returns ACKNOWLEDGED_UPDATE_RESULT immediately.
	if name == nil && groupSizeLimit == nil && invitationStrategy == nil && joinStrategy == nil &&
		groupInfoUpdateStrategy == nil && memberInfoUpdateStrategy == nil &&
		guestSpeakable == nil && selfInfoUpdatable == nil && enableReadReceipt == nil && messageEditable == nil {
		return nil
	}

	err := s.groupTypeRepo.UpdateTypes(ctx, ids, name, groupSizeLimit, invitationStrategy, joinStrategy, groupInfoUpdateStrategy, memberInfoUpdateStrategy, guestSpeakable, selfInfoUpdatable, enableReadReceipt, messageEditable)
	if err != nil {
		return err
	}

	// Bug fix: invalidate in-memory cache for updated types.
	for _, id := range ids {
		s.idToGroupType.Delete(id)
	}

	// Log cache invalidation for debugging
	log.Printf("Invalidated cache for group type IDs: %v", ids)
	return nil
}
