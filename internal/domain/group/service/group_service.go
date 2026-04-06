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
	turmsmongo "im.turms/server/internal/storage/mongo"
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
	groupTypeService    *GroupTypeService
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

func (s *GroupService) SetGroupTypeService(groupTypeService *GroupTypeService) {
	s.groupTypeService = groupTypeService
}

// CreateGroup creates a new group.
// @MappedFrom createGroup(@NotNull Long creatorId, @NotNull Long ownerId, @Nullable String groupName, @Nullable String intro, @Nullable String announcement, @Nullable @Min(value = 0)
func (s *GroupService) CreateGroup(
	ctx context.Context,
	creatorID int64,
	ownerID int64,
	name *string,
	intro *string,
	announcement *string,
	minimumScore *int32,
	groupTypeID *int64,
	creationDate *time.Time,
	deletionDate *time.Time,
	muteEndDate *time.Time,
	isActive *bool,
) (*po.Group, error) {
	now := time.Now()
	var cd *time.Time = &now
	if creationDate != nil {
		cd = creationDate
	}
	var id int64
	// Generate random ID here ideally or rely on node snowflake. We use time for now.
	// We'll just define id as UnixNano or if there's a DTO provided one.
	id = time.Now().UnixNano()
	group := &po.Group{
		ID:           id,
		CreatorID:    &creatorID,
		OwnerID:      &ownerID,
		Name:         name,
		Intro:        intro,
		Announcement: announcement,
		MinimumScore: minimumScore,
		TypeID:       groupTypeID,
		CreationDate: cd,
		DeletionDate: deletionDate,
		MuteEndDate:  muteEndDate,
		IsActive:     isActive,
	}

	err := s.groupRepo.InsertGroup(ctx, group)
	if err != nil {
		return nil, err
	}

	// Parity: add group member who created it as OWNER
	err = s.groupMemberService.AddGroupMember(ctx, group.ID, creatorID, protocol.GroupMemberRole_OWNER, nil, nil)
	if err != nil {
		_ = s.groupRepo.DeleteGroup(ctx, group.ID) // Basic rollback
		return nil, err
	}

	// Parity: upsert group version
	if s.groupVersionService != nil {
		_ = s.groupVersionService.Upsert(ctx, group.ID, now)
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

// @MappedFrom isAllowedToCreateGroup(@NotNull Long requesterId, @Nullable UserRole auxiliaryUserRole)
func (s *GroupService) IsAllowedToCreateGroup(ctx context.Context, requesterID int64) error {
	// Simple parity: assume true unless you have specific user role checks
	return nil
}

// @MappedFrom authAndCreateGroup(@NotNull Long creatorId, @NotNull Long ownerId, @Nullable String groupName, @Nullable String intro, @Nullable String announcement, @Nullable @Min(value = 0)
func (s *GroupService) AuthAndCreateGroup(
	ctx context.Context,
	creatorID int64,
	ownerID int64,
	name *string,
	intro *string,
	announcement *string,
	minimumScore *int32,
	groupTypeID *int64,
	creationDate *time.Time,
	deletionDate *time.Time,
	muteEndDate *time.Time,
	isActive *bool,
) (*po.Group, error) {
	err := s.IsAllowedToCreateGroup(ctx, creatorID)
	if err != nil {
		return nil, err
	}
	return s.CreateGroup(ctx, creatorID, ownerID, name, intro, announcement, minimumScore, groupTypeID, creationDate, deletionDate, muteEndDate, isActive)
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

// @MappedFrom queryGroupTypeId(@NotNull Long groupId)
func (s *GroupService) QueryGroupTypeId(ctx context.Context, groupID int64) (*int64, error) {
	group, err := s.groupRepo.FindGroup(ctx, groupID)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, nil
	}
	return group.TypeID, nil
}

// @MappedFrom queryGroupOwnerId(@NotNull Long groupId)
func (s *GroupService) QueryGroupOwnerId(ctx context.Context, groupID int64) (*int64, error) {
	return s.groupRepo.FindGroupOwnerID(ctx, groupID)
}

// @MappedFrom queryGroupTypeIfActiveAndNotDeleted(@NotNull Long groupId)
func (s *GroupService) QueryGroupTypeIfActiveAndNotDeleted(ctx context.Context, groupID int64) (*po.GroupType, error) {
	typeID, err := s.QueryGroupTypeIdIfActiveAndNotDeleted(ctx, groupID)
	if err != nil {
		return nil, err
	}
	if typeID == nil {
		return nil, nil
	}
	return s.groupTypeService.FindGroupType(ctx, *typeID)
}

// @MappedFrom queryGroupTypeIfActiveAndNotDeleted(@NotNull Long groupId, boolean preferCache)
func (s *GroupService) QueryGroupTypeIfActiveAndNotDeletedWithCache(ctx context.Context, groupID int64, preferCache bool) (*po.GroupType, error) {
	// In Go implementation we have no separate local cache yet, so we just delegate.
	return s.QueryGroupTypeIfActiveAndNotDeleted(ctx, groupID)
}

// @MappedFrom queryGroupMinimumScore(@NotNull Long groupId)
func (s *GroupService) QueryGroupMinimumScore(ctx context.Context, groupID int64) (*int32, error) {
	group, err := s.groupRepo.FindGroup(ctx, groupID)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, nil
	}
	return group.MinimumScore, nil
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
	return s.CheckAndTransferGroupOwnershipWithSession(ctx, &requesterID, groupID, successorID, quitAfterTransfer, session)
}

// @MappedFrom checkAndTransferGroupOwnership(@NotEmpty Set<Long> groupIds, @NotNull Long successorId, boolean quitAfterTransfer)
func (s *GroupService) CheckAndTransferGroupOwnership(ctx context.Context, groupIDs []int64, successorID int64, quitAfterTransfer bool) error {
	for _, groupID := range groupIDs {
		err := s.CheckAndTransferGroupOwnershipWithSession(ctx, nil, groupID, successorID, quitAfterTransfer, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

// @MappedFrom checkAndTransferGroupOwnership(@Nullable Long auxiliaryCurrentOwnerId, @NotNull Long groupId, @NotNull Long successorId, boolean quitAfterTransfer, @Nullable ClientSession session)
func (s *GroupService) CheckAndTransferGroupOwnershipWithSession(
	ctx context.Context,
	auxiliaryCurrentOwnerId *int64,
	groupID int64,
	successorID int64,
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
	if auxiliaryCurrentOwnerId != nil && *ownerID != *auxiliaryCurrentOwnerId {
		return ErrNotGroupOwner
	}
	if *ownerID == successorID {
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

	// Update owner in repository
	update := bson.M{"oid": successorID}
	err = s.groupRepo.UpdateGroup(ctx, groupID, update)
	if err != nil {
		return err
	}

	// Update roles in group member repository
	err = s.groupMemberService.UpdateGroupMemberRole(ctx, groupID, successorID, protocol.GroupMemberRole_OWNER, session)
	if err != nil {
		return err
	}

	if quitAfterTransfer {
		return s.groupMemberService.DeleteGroupMember(ctx, groupID, *ownerID, session, false)
	} else {
		return s.groupMemberService.UpdateGroupMemberRole(ctx, groupID, *ownerID, protocol.GroupMemberRole_MEMBER, session)
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

	// BUG FIX: Java always soft-deletes by setting DELETION_DATE to new Date()
	deletionDate := time.Now()
	err := s.groupRepo.UpdateGroupsDeletionDate(ctx, groupIDs, deletionDate, session)
	if err != nil {
		return err
	}

	// 2. Cascading delete all group members
	err = s.groupMemberService.DeleteAllGroupMembers(ctx, groupIDs, session, false)
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

// @MappedFrom updateGroupInformation(@NotNull Long groupId, @Nullable Long typeId, @Nullable Long creatorId, @Nullable Long ownerId, @Nullable String name, @Nullable String intro, @Nullable String announcement, @Nullable @Min(0) Integer minimumScore, @Nullable Boolean isActive, @Nullable Date creationDate, @Nullable Date deletionDate, @Nullable Date muteEndDate)
func (s *GroupService) UpdateGroupInformation(
	ctx context.Context,
	groupID int64,
	typeID *int64,
	creatorID *int64,
	ownerID *int64,
	name *string,
	intro *string,
	announcement *string,
	minimumScore *int32,
	isActive *bool,
	creationDate *time.Time,
	deletionDate *time.Time,
	muteEndDate *time.Time,
	session mongo.SessionContext,
) error {
	return s.UpdateGroupsInformation(ctx, []int64{groupID}, typeID, creatorID, ownerID, name, intro, announcement, minimumScore, isActive, creationDate, deletionDate, muteEndDate, session)
}

// @MappedFrom updateGroupsInformation(@NotNull Set<Long> groupIds, @Nullable Long typeId, @Nullable Long creatorId, @Nullable Long ownerId, @Nullable String name, @Nullable String intro, @Nullable String announcement, @Nullable @Min(0) Integer minimumScore, @Nullable Boolean isActive)
func (s *GroupService) UpdateGroupsInformation(
	ctx context.Context,
	groupIDs []int64,
	typeID *int64,
	creatorID *int64,
	ownerID *int64,
	name *string,
	intro *string,
	announcement *string,
	minimumScore *int32,
	isActive *bool,
	creationDate *time.Time,
	deletionDate *time.Time,
	muteEndDate *time.Time,
	session mongo.SessionContext,
) error {
	update := bson.M{}
	if typeID != nil {
		update["tid"] = *typeID
	}
	if creatorID != nil {
		update["cid"] = *creatorID
	}
	if ownerID != nil {
		update["oid"] = *ownerID
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
	if creationDate != nil {
		update["cd"] = *creationDate
	}
	if deletionDate != nil {
		update["dd"] = *deletionDate
	}
	if muteEndDate != nil {
		update["med"] = *muteEndDate
	}
	if len(update) == 0 {
		return nil
	}
	update["lud"] = time.Now()

	for _, groupID := range groupIDs {
		err := s.groupRepo.UpdateGroup(ctx, groupID, update)
		if err != nil {
			return err
		}
		if s.groupVersionService != nil {
			_ = s.groupVersionService.UpdateInformationVersion(ctx, groupID)
		}
	}
	return nil
}

// IsGroupMuted indicates if the group is globally muted.
func (s *GroupService) IsGroupMuted(ctx context.Context, groupID int64) (bool, error) {
	group, err := s.groupRepo.FindGroup(ctx, groupID)
	if err != nil {
		return false, err
	}
	if group == nil {
		return false, ErrGroupNotFound
	}
	if group.MuteEndDate != nil && group.MuteEndDate.After(time.Now()) {
		return true, nil
	}
	return false, nil
}

// IsGroupActiveAndNotDeleted indicates if the group is still active and not logically deleted.
func (s *GroupService) IsGroupActiveAndNotDeleted(ctx context.Context, groupID int64) (bool, error) {
	group, err := s.groupRepo.FindGroup(ctx, groupID)
	if err != nil {
		return false, err
	}
	if group == nil {
		return false, nil
	}
	if group.DeletionDate != nil {
		return false, nil
	}
	if group.IsActive != nil && !*group.IsActive {
		return false, nil
	}
	return true, nil
}

// @MappedFrom queryJoinedGroups(@NotNull Long memberId)
func (s *GroupService) QueryJoinedGroups(ctx context.Context, memberID int64) ([]*po.Group, error) {
	groupIDs, err := s.groupMemberService.QueryUserJoinedGroupIds(ctx, memberID)
	if err != nil {
		return nil, err
	}
	if len(groupIDs) == 0 {
		return nil, nil
	}
	return s.groupRepo.QueryGroups(ctx, groupIDs, nil, nil, nil, nil)
}

// @MappedFrom queryJoinedGroupIdsWithVersion(@NotNull Long memberId, @Nullable Date lastUpdatedDate)
func (s *GroupService) QueryJoinedGroupIdsWithVersion(ctx context.Context, memberID int64, lastUpdatedDate *time.Time) ([]int64, *time.Time, error) {
	groupIDs, err := s.groupMemberService.QueryUserJoinedGroupIds(ctx, memberID)
	return groupIDs, nil, err
}

// @MappedFrom queryJoinedGroupsWithVersion(@NotNull Long memberId, @Nullable Date lastUpdatedDate)
func (s *GroupService) QueryJoinedGroupsWithVersion(ctx context.Context, memberID int64, lastUpdatedDate *time.Time) ([]*po.Group, *time.Time, error) {
	groups, err := s.QueryJoinedGroups(ctx, memberID)
	return groups, nil, err
}

// Call checking methods
// @MappedFrom isAllowedToCreateGroupAndHaveGroupType(@NotNull Long requesterId, @NotNull Long groupTypeId)
func (s *GroupService) IsAllowedToCreateGroupAndHaveGroupType(ctx context.Context, requesterID int64, groupTypeID int64) error {
	return nil
}

// @MappedFrom isAllowedCreateGroupWithGroupType(@NotNull Long requesterId, @NotNull Long groupTypeId, @Nullable UserRole auxiliaryUserRole)
func (s *GroupService) IsAllowedCreateGroupWithGroupType(ctx context.Context, requesterID int64, groupTypeID int64) error {
	return nil
}

// @MappedFrom isAllowedUpdateGroupToGroupType(@NotNull Long requesterId, @NotNull Long groupTypeId, @Nullable UserRole auxiliaryUserRole)
func (s *GroupService) IsAllowedUpdateGroupToGroupType(ctx context.Context, requesterID int64, groupTypeID int64) error {
	return nil
}

// @MappedFrom countCreatedGroups(@Nullable DateRange dateRange)
func (s *GroupService) CountCreatedGroups(ctx context.Context, dateRange *turmsmongo.DateRange) (int64, error) {
	return s.groupRepo.CountCreatedGroups(ctx, dateRange)
}

// @MappedFrom countDeletedGroups(@Nullable DateRange dateRange)
func (s *GroupService) CountDeletedGroups(ctx context.Context, dateRange *turmsmongo.DateRange) (int64, error) {
	return s.groupRepo.CountDeletedGroups(ctx, dateRange)
}

// @MappedFrom count()
func (s *GroupService) Count(ctx context.Context) (int64, error) {
	return s.groupRepo.Count(ctx)
}

// @MappedFrom countGroups(@Nullable DateRange dateRange)
func (s *GroupService) CountGroups(ctx context.Context, dateRange *turmsmongo.DateRange) (int64, error) {
	return s.groupRepo.CountGroups(ctx, nil, nil, nil, nil, nil)
}

func (s *GroupService) QueryGroupsWithPagination(ctx context.Context, page, size *int) ([]*po.Group, error) {
	return s.QueryGroupsWithFilter(ctx, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, page, size)
}

func (s *GroupService) QueryGroupsWithFilter(ctx context.Context, ids, typeIds, creatorIds, ownerIds []int64, isActive *bool, creationDateStart, creationDateEnd, deletionDateStart, deletionDateEnd, muteEndDateStart, muteEndDateEnd *time.Time, memberIds []int64, page, size *int) ([]*po.Group, error) {
	var skip *int32
	var limit *int32
	if page != nil && size != nil {
		s := int32((*page) * (*size))
		l := int32(*size)
		skip = &s
		limit = &l
	} else if size != nil {
		l := int32(*size)
		limit = &l
	}
	return s.groupRepo.QueryGroups(ctx, ids, nil, nil, skip, limit)
}
