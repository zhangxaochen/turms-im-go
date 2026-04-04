package service

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"im.turms/server/internal/domain/common/constant"
	"im.turms/server/internal/domain/group/po"
	"im.turms/server/internal/domain/group/repository"
	"im.turms/server/internal/infra/exception"
	"im.turms/server/pkg/protocol"
)

var (
	ErrGroupNotFound = errors.New("group not found")
	ErrNotGroupOwner = errors.New("not the group owner")
)

type GroupService struct {
	groupRepo           *repository.GroupRepository
	groupMemberService  *GroupMemberService
	groupVersionService *GroupVersionService
}

func NewGroupService(groupRepo *repository.GroupRepository) *GroupService {
	return &GroupService{
		groupRepo: groupRepo,
	}
}

func (s *GroupService) SetGroupMemberService(groupMemberService *GroupMemberService) {
	s.groupMemberService = groupMemberService
}

func (s *GroupService) SetGroupVersionService(groupVersionService *GroupVersionService) {
	s.groupVersionService = groupVersionService
}

// CreateGroup creates a new group.
// @MappedFrom createGroup(@NotNull Long creatorId, @NotNull Long ownerId, @Nullable String groupName, @Nullable String intro, @Nullable String announcement, @Nullable @Min(value = 0)
func (s *GroupService) CreateGroup(ctx context.Context, creatorID, groupID int64, name, intro *string, minimumScore *int32) (*po.Group, error) {
	now := time.Now()
	group := &po.Group{
		ID:           groupID,
		CreatorID:    &creatorID,
		OwnerID:      &creatorID,
		Name:         name,
		Intro:        intro,
		MinimumScore: minimumScore,
		CreationDate: &now,
	}

	err := s.groupRepo.InsertGroup(ctx, group)
	if err != nil {
		return nil, err
	}

	// Parity: add group member who created it as OWNER
	err = s.groupMemberService.AddGroupMember(ctx, groupID, creatorID, protocol.GroupMemberRole_OWNER, nil, nil)
	if err != nil {
		_ = s.groupRepo.DeleteGroup(ctx, groupID) // Basic rollback
		return nil, err
	}

	// Parity: upsert group version
	if s.groupVersionService != nil {
		_ = s.groupVersionService.Upsert(ctx, groupID, now)
	}

	// TODO: add metric increment (createdGroupsCounter.increment)
	// TODO: add Elasticsearch integration if supported

	return group, nil
}

// DeleteGroup performs a soft deletion of the group.
// Only the owner can delete the group.
func (s *GroupService) DeleteGroup(ctx context.Context, requesterID, groupID int64) error {
	ownerID, err := s.groupRepo.FindGroupOwnerID(ctx, groupID)
	if err != nil {
		return err
	}
	if ownerID == nil {
		return ErrGroupNotFound
	}

	if *ownerID != requesterID {
		return ErrNotGroupOwner
	}

	now := time.Now()
	update := bson.M{}
	update["dd"] = now

	return s.groupRepo.UpdateGroup(ctx, groupID, update)
}

// @MappedFrom queryGroupTypeIdIfActiveAndNotDeleted(@NotNull Long groupId)
func (s *GroupService) QueryGroupTypeIdIfActiveAndNotDeleted(ctx context.Context, groupID int64) (*int64, error) {
	group, err := s.groupRepo.FindGroup(ctx, groupID)
	if err != nil {
		return nil, err
	}
	if group == nil || group.DeletionDate != nil || (group.IsActive != nil && !*group.IsActive) {
		return nil, nil
	}
	return group.TypeID, nil
}

// @MappedFrom queryGroupMinimumScoreIfActiveAndNotDeleted(@NotNull Long groupId)
func (s *GroupService) QueryGroupMinimumScoreIfActiveAndNotDeleted(ctx context.Context, groupID int64) (*int32, error) {
	group, err := s.groupRepo.FindGroup(ctx, groupID)
	if err != nil {
		return nil, err
	}
	if group == nil || group.DeletionDate != nil || (group.IsActive != nil && !*group.IsActive) {
		return nil, nil
	}
	if group.MinimumScore == nil {
		var zero int32 = 0
		return &zero, nil
	}
	return group.MinimumScore, nil
}

// @MappedFrom authAndTransferGroupOwnership(@NotNull Long requesterId, @NotNull Long groupId, @NotNull Long successorId, boolean quitAfterTransfer, @Nullable ClientSession session)
func (s *GroupService) AuthAndTransferGroupOwnership(
	ctx context.Context,
	requesterID, groupID, successorID int64,
	quitAfterTransfer bool,
	session mongo.SessionContext,
) error {
	ownerID, err := s.groupRepo.FindGroupOwnerID(ctx, groupID)
	if err != nil {
		return err
	}
	if ownerID == nil {
		return ErrGroupNotFound
	}
	if *ownerID != requesterID {
		return ErrNotGroupOwner
	}
	if requesterID == successorID {
		return nil
	}

	// Parity: check if successor is a member
	isMember, err := s.groupMemberService.IsGroupMember(ctx, groupID, successorID)
	if err != nil {
		return err
	}
	if !isMember {
		return exception.NewTurmsError(int32(constant.ResponseStatusCode_GROUP_SUCCESSOR_NOT_GROUP_MEMBER), "Successor is not a member of the group")
	}
	// TODO: add isAllowedToCreateGroupAndHaveGroupType constraint check parity

	// Update owner in repository
	update := bson.M{"oid": successorID}
	err = s.groupRepo.UpdateGroup(ctx, groupID, update)
	if err != nil {
		return err
	}

	// Update roles in group member repository
	// Successor becomes Owner, Requester becomes Member (if quitAfterTransfer is false)
	err = s.groupMemberService.UpdateGroupMemberRole(ctx, groupID, successorID, protocol.GroupMemberRole_OWNER, session)
	if err != nil {
		return err
	}

	if quitAfterTransfer {
		return s.groupMemberService.DeleteGroupMember(ctx, groupID, requesterID, session, false)
	} else {
		return s.groupMemberService.UpdateGroupMemberRole(ctx, groupID, requesterID, protocol.GroupMemberRole_MEMBER, session)
	}
}

// AuthAndDeleteGroup deletes a group after authorization check.
// @MappedFrom authAndDeleteGroup(boolean queryGroupMemberIds, @NotNull Long requesterId, @NotNull Long groupId)
func (s *GroupService) AuthAndDeleteGroup(ctx context.Context, requesterID int64, groupID int64) error {
	ownerID, err := s.groupRepo.FindGroupOwnerID(ctx, groupID)
	if err != nil {
		return err
	}
	if ownerID == nil {
		return exception.NewTurmsError(int32(constant.ResponseStatusCode_UPDATE_INFO_OF_NONEXISTENT_GROUP), "Group does not exist")
	}
	if *ownerID != requesterID {
		return exception.NewTurmsError(int32(constant.ResponseStatusCode_NOT_GROUP_OWNER_TO_DELETE_GROUP), "Only the owner can delete the group")
	}

	// Call DeleteGroupsAndGroupMembers parity (cascading)
	return s.DeleteGroupsAndGroupMembers(ctx, []int64{groupID}, nil)
}

// DeleteGroupsAndGroupMembers performs cascading deletion parity.
// @MappedFrom deleteGroupsAndGroupMembers(@Nullable Set<Long> groupIds, @Nullable Boolean deleteLogically)
func (s *GroupService) DeleteGroupsAndGroupMembers(ctx context.Context, groupIDs []int64, session mongo.SessionContext) error {
	if len(groupIDs) == 0 {
		return nil
	}

	// 1. Delete groups physically/logically
	// Normally handled by groupRepo.DeleteByIds, here we loop delete for simplicity based on repo
	for _, groupID := range groupIDs {
		// Assuming hard delete conceptually, turms-orig toggles logic via props.
		// For now we use the existing DeleteGroup
		err := s.groupRepo.DeleteGroup(ctx, groupID)
		if err != nil {
			return err
		}
	}

	// 2. Cascading delete all group members
	err := s.groupMemberService.DeleteAllGroupMembers(ctx, groupIDs, session, false)
	if err != nil {
		return err
	}

	// 3. Cascading delete group versions
	if s.groupVersionService != nil {
		err = s.groupVersionService.Delete(ctx, groupIDs)
		if err != nil {
			return err
		}
	}

	// TODO: cascading message sequence IDs and conversations
	return nil
}

// AuthAndQueryGroups queries groups. In Java, this method is called on groupService.
// @MappedFrom authAndQueryGroups
func (s *GroupService) AuthAndQueryGroups(ctx context.Context, groupIDs []int64, name *string, lastUpdatedDate *time.Time, skip *int32, limit *int32, fieldsToHighlight []int32) ([]*po.Group, error) {
	// TODO: Add auth and highlights logic if necessary based on fieldsToHighlight
	// For basic parity, we just delegate to QueryGroups
	return s.groupRepo.QueryGroups(ctx, groupIDs, name, lastUpdatedDate, skip, limit)
}

// AuthAndUpdateGroup updates a group.
// @MappedFrom authAndUpdateGroup(@NotNull Long requesterId, @NotNull Long groupId, @Nullable Long typeId, @Nullable Long successorId, @Nullable String name, @Nullable String intro, @Nullable String announcement, @Nullable @Min(value = 0) Integer minimumScore, @Nullable @ValidGroupType GroupType groupType, @Nullable Boolean isActive, @Nullable Boolean quitAfterTransfer)
func (s *GroupService) AuthAndUpdateGroup(
	ctx context.Context,
	requesterID int64,
	groupID int64,
	typeID *int64,
	successorID *int64,
	name *string,
	intro *string,
	announcement *string,
	minimumScore *int32,
	isActive *bool,
	quitAfterTransfer *bool,
) error {
	ownerID, err := s.groupRepo.FindGroupOwnerID(ctx, groupID)
	if err != nil {
		return err
	}
	if ownerID == nil {
		return exception.NewTurmsError(int32(constant.ResponseStatusCode_UPDATE_INFO_OF_NONEXISTENT_GROUP), "Group does not exist")
	}

	isManager := false
	if *ownerID != requesterID {
		role, err := s.groupMemberService.FindGroupMemberRole(ctx, groupID, requesterID)
		if err != nil {
			return err
		}
		if role == nil || (*role != protocol.GroupMemberRole_MANAGER) {
			return exception.NewTurmsError(int32(constant.ResponseStatusCode_NOT_GROUP_OWNER_OR_MANAGER_TO_UPDATE_GROUP_INFO), "Only owner or manager can update group info")
		}
		isManager = true
	}

	if isManager && (typeID != nil || successorID != nil || isActive != nil) {
		return exception.NewTurmsError(int32(constant.ResponseStatusCode_NOT_GROUP_OWNER_TO_UPDATE_GROUP_INFO), "Only owner can update type, owner, or active status")
	}

	if successorID != nil {
		transferQuit := false
		if quitAfterTransfer != nil {
			transferQuit = *quitAfterTransfer
		}
		err := s.AuthAndTransferGroupOwnership(ctx, requesterID, groupID, *successorID, transferQuit, nil)
		if err != nil {
			return err
		}
		if name == nil && intro == nil && announcement == nil && minimumScore == nil && typeID == nil && isActive == nil {
			return nil
		}
	}

	update := bson.M{}
	if typeID != nil {
		update["tid"] = *typeID
	}
	if name != nil {
		update["n"] = *name
	}
	if intro != nil {
		update["intro"] = *intro
	}
	if announcement != nil {
		update["annc"] = *announcement
	}
	if minimumScore != nil {
		update["ms"] = *minimumScore
	}
	if isActive != nil {
		update["ac"] = *isActive
	}

	if len(update) == 0 {
		return nil
	}

	update["lud"] = time.Now()

	err = s.groupRepo.UpdateGroup(ctx, groupID, update)
	if err != nil {
		return err
	}

	if s.groupVersionService != nil {
		return s.groupVersionService.UpdateInformationVersion(ctx, groupID)
	}

	return nil
}
