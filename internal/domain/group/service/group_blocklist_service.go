package service

import (
	"context"
	"time"

	common_constant "im.turms/server/internal/domain/common/constant"
	"im.turms/server/internal/domain/group/po"
	"im.turms/server/internal/domain/group/repository"
	"im.turms/server/internal/infra/exception"
	turmsmongo "im.turms/server/internal/storage/mongo"
	"im.turms/server/pkg/codes"
)

type GroupBlocklistService struct {
	blockedUserRepo     repository.GroupBlockedUserRepository
	groupMemberService  *GroupMemberService
	groupVersionService *GroupVersionService
}

func NewGroupBlocklistService(
	blockedUserRepo repository.GroupBlockedUserRepository,
	groupVersionService *GroupVersionService,
) *GroupBlocklistService {
	return &GroupBlocklistService{
		blockedUserRepo:     blockedUserRepo,
		groupVersionService: groupVersionService,
	}
}

func (s *GroupBlocklistService) SetGroupMemberService(groupMemberService *GroupMemberService) {
	s.groupMemberService = groupMemberService
}

// BlockUser blocks a user from a group without authorization checks.
// @MappedFrom blockUser(@NotNull Long requesterId, @NotNull Long groupId, @NotNull Long userIdToBlock, @Nullable ClientSession session)
func (s *GroupBlocklistService) BlockUser(ctx context.Context, groupID int64, userID int64, requesterID int64) error {
	now := time.Now()
	blockedUser := &po.GroupBlockedUser{
		ID: po.GroupBlockedUserKey{
			GroupID: groupID,
			UserID:  userID,
		},
		BlockDate:   &now,
		RequesterID: requesterID,
	}
	return s.blockedUserRepo.Insert(ctx, blockedUser)
}

// AuthAndBlockUser blocks a user from a group after performing authorization checks.
// @MappedFrom authAndBlockUser(@NotNull Long requesterId, @NotNull Long groupId, @NotNull Long userIdToBlock, @Nullable ClientSession session)
// Bug fixes:
// - Added "cannot block oneself" check (Java checks requesterId.equals(userIdToBlock))
// - Updated both members AND blocklist version when target is a member (Java: updateVersion(groupId, true, true, false, false))
// - Properly propagate errors from DeleteGroupMember instead of ignoring
// - Removed extra role hierarchy check not present in Java
func (s *GroupBlocklistService) AuthAndBlockUser(
	ctx context.Context,
	requesterID int64,
	groupID int64,
	userID int64,
) error {
	// 1. Cannot block oneself
	if requesterID == userID {
		return exception.NewTurmsError(codes.IllegalArgument, "Cannot block oneself")
	}

	// 2. Authorization check: must be owner or manager
	isOwnerOrManager, err := s.groupMemberService.IsOwnerOrManager(ctx, groupID, requesterID)
	if err != nil {
		return err
	}
	if !isOwnerOrManager {
		return exception.NewTurmsError(int32(common_constant.ResponseStatusCode_NOT_GROUP_OWNER_OR_MANAGER_TO_ADD_BLOCKED_USER), "Only owner or manager can block users")
	}

	// 3. Check if target is a group member
	isMember, err := s.groupMemberService.IsGroupMember(ctx, groupID, userID)
	if err != nil {
		return err
	}

	// 4. Block the user
	err = s.BlockUser(ctx, groupID, userID, requesterID)
	if err != nil {
		return err
	}

	// 5. If target is a member, remove from group and update both members + blocklist version
	if isMember {
		// Java: deleteGroupMember + insert(blockedUser) in a transaction
		// For now, sequential execution (transactions not yet implemented)
		err = s.groupMemberService.DeleteGroupMember(ctx, groupID, userID, nil, true)
		if err != nil {
			return err
		}
		// Java: updateVersion(groupId, true, true, false, false) — updates members AND blocklist
		_ = s.groupVersionService.UpdateVersionFields(ctx, groupID, true, true, false, false)
		return nil
	}

	// 6. If target is not a member, only update blocklist version
	return s.groupVersionService.UpdateBlocklistVersion(ctx, groupID)
}

// AuthAndUnblockUser unblocks a user from a group after performing authorization checks.
// @MappedFrom unblockUser(@NotNull Long requesterId, @NotNull Long groupId, @NotNull Long userIdToUnblock, @Nullable ClientSession session, boolean updateBlocklistVersion)
// Bug fixes:
// - Added updateBlocklistVersion parameter for caller control
// - Only update version if user was actually blocked (wasBlocked == true)
// - Returns (wasBlocked, error) - wasBlocked indicates whether the user was actually blocked before unblocking.
func (s *GroupBlocklistService) AuthAndUnblockUser(
	ctx context.Context,
	requesterID int64,
	groupID int64,
	userID int64,
	updateBlocklistVersion bool,
) (bool, error) {
	// 1. Authorization check
	isOwnerOrManager, err := s.groupMemberService.IsOwnerOrManager(ctx, groupID, requesterID)
	if err != nil {
		return false, err
	}
	if !isOwnerOrManager {
		return false, exception.NewTurmsError(int32(common_constant.ResponseStatusCode_NOT_GROUP_OWNER_OR_MANAGER_TO_REMOVE_BLOCKED_USER), "Only owner or manager can unblock users")
	}

	// 2. Unblock
	wasBlocked, err := s.UnblockUser(ctx, groupID, userID)
	if err != nil {
		return false, err
	}

	// 3. Only update version if user was actually blocked and updateBlocklistVersion is true
	if updateBlocklistVersion && wasBlocked {
		_ = s.groupVersionService.UpdateBlocklistVersion(ctx, groupID)
	}

	return wasBlocked, nil
}

// @MappedFrom unblockUser(@NotNull Long requesterId, @NotNull Long groupId, @NotNull Long userIdToUnblock, @Nullable ClientSession session, boolean updateBlocklistVersion)
// Bug fix: Now returns whether the user was actually blocked (deletedCount > 0) for caller control.
func (s *GroupBlocklistService) UnblockUser(ctx context.Context, groupID int64, userID int64) (bool, error) {
	deletedCount, err := s.blockedUserRepo.Delete(ctx, groupID, userID)
	if err != nil {
		return false, err
	}
	return deletedCount > 0, nil
}

// @MappedFrom queryBlockedUsers(int page, @QueryParam(required = false)
// @MappedFrom queryBlockedUsers(Set<Long> ids)
// @MappedFrom queryBlockedUsers(@Nullable Set<Long> groupIds, @Nullable Set<Long> userIds, @Nullable DateRange blockDateRange, @Nullable Set<Long> requesterIds, @Nullable Integer page, @Nullable Integer size)
// @MappedFrom isBlocked(Long ownerId, Long relatedUserId)
// @MappedFrom isBlocked(@NotNull Long groupId, @NotNull Long userId)
// @MappedFrom isBlocked(@NotNull Long ownerId, @NotNull Long relatedUserId, boolean preferCache)
func (s *GroupBlocklistService) IsBlocked(ctx context.Context, groupID int64, userID int64) (bool, error) {
	return s.blockedUserRepo.Exists(ctx, groupID, userID)
}

func (s *GroupBlocklistService) FilterBlockedUserIDs(ctx context.Context, groupID int64, userIDs []int64) ([]int64, error) {
	return s.blockedUserRepo.FilterBlockedUserIDs(ctx, groupID, userIDs)
}

// AuthAndQueryGroupBlockedUserIds queries blocked user IDs with auth check.
// Bug fixes:
// - Always query version (even when lastUpdatedDate is nil) to match Java behavior
// - Return NO_CONTENT error when blocked user IDs list is empty
func (s *GroupBlocklistService) AuthAndQueryGroupBlockedUserIds(
	ctx context.Context,
	requesterID int64,
	groupID int64,
	lastUpdatedDate *time.Time,
) ([]int64, *time.Time, error) {
	// 1. Authorization
	role, err := s.groupMemberService.FindGroupMemberRole(ctx, groupID, requesterID)
	if err != nil {
		return nil, nil, err
	}
	if role == nil {
		return nil, nil, exception.NewTurmsError(int32(common_constant.ResponseStatusCode_NOT_GROUP_MEMBER_TO_QUERY_GROUP_MEMBER_INFO), "Only group members can query blocked user info")
	}

	// 2. Version check — always query version (even when lastUpdatedDate is nil)
	version, err := s.groupVersionService.QueryGroupBlocklistVersion(ctx, groupID)
	if err != nil {
		return nil, nil, err
	}
	if version == nil {
		// Java: switchIfEmpty -> alreadyUpToUpdate()
		return nil, nil, nil
	}
	if lastUpdatedDate != nil && !version.After(*lastUpdatedDate) {
		// Already up to date
		return nil, nil, nil
	}

	// 3. Query
	userIDs, err := s.blockedUserRepo.FindBlockedUserIds(ctx, groupID)
	if err != nil {
		return nil, nil, err
	}

	if len(userIDs) == 0 {
		// Java: throw NO_CONTENT
		return nil, nil, exception.NewTurmsError(int32(common_constant.ResponseStatusCode_NO_CONTENT), "No blocked users found")
	}

	return userIDs, version, nil
}

// AuthAndQueryGroupBlockedUserInfos queries blocked user infos with auth check.
// Bug fixes:
// - Always query version (even when lastUpdatedDate is nil)
// - Return NO_CONTENT when empty results
func (s *GroupBlocklistService) AuthAndQueryGroupBlockedUserInfos(
	ctx context.Context,
	requesterID int64,
	groupID int64,
	lastUpdatedDate *time.Time,
) ([]po.GroupBlockedUser, *time.Time, error) {
	// 1. Authorization
	role, err := s.groupMemberService.FindGroupMemberRole(ctx, groupID, requesterID)
	if err != nil {
		return nil, nil, err
	}
	if role == nil {
		return nil, nil, exception.NewTurmsError(int32(common_constant.ResponseStatusCode_NOT_GROUP_MEMBER_TO_QUERY_GROUP_MEMBER_INFO), "Only group members can query blocked user info")
	}

	// 2. Version check — always query version (even when lastUpdatedDate is nil)
	version, err := s.groupVersionService.QueryGroupBlocklistVersion(ctx, groupID)
	if err != nil {
		return nil, nil, err
	}
	if version == nil {
		// Java: switchIfEmpty -> alreadyUpToUpdate()
		return nil, nil, nil
	}
	if lastUpdatedDate != nil && !version.After(*lastUpdatedDate) {
		return nil, nil, nil
	}

	// 3. Query
	users, err := s.blockedUserRepo.FindBlockedUsersByGroupID(ctx, groupID)
	if err != nil {
		return nil, nil, err
	}

	if len(users) == 0 {
		// Java: throw NO_CONTENT
		return nil, nil, exception.NewTurmsError(int32(common_constant.ResponseStatusCode_NO_CONTENT), "No blocked users found")
	}

	return users, version, nil
}

// QueryBlockedUsers returns blocked users for backward compatibility with tests
func (s *GroupBlocklistService) QueryBlockedUsers(ctx context.Context, groupID int64) ([]po.GroupBlockedUser, error) {
	return s.blockedUserRepo.FindBlockedUsersByGroupID(ctx, groupID)
}

func (s *GroupBlocklistService) QueryGroupBlockedUserIds(ctx context.Context, groupID int64) ([]int64, error) {
	return s.blockedUserRepo.FindBlockedUserIds(ctx, groupID)
}

func (s *GroupBlocklistService) CountBlockedUsers(
	ctx context.Context,
	groupIds []int64,
	userIds []int64,
	blockDateRange *turmsmongo.DateRange,
	requesterIds []int64,
) (int64, error) {
	return s.blockedUserRepo.CountBlockedUsers(ctx, groupIds, userIds, blockDateRange, requesterIds)
}

// QueryGroupBlockedUserIdsWithVersion queries blocked user IDs with version control.
// Bug fixes:
// - Always query version (even when lastUpdatedDate is nil) to match Java behavior
// - Return NO_CONTENT error when blocked user IDs list is empty
// - Return proper alreadyUpToUpdate semantics when version is nil
func (s *GroupBlocklistService) QueryGroupBlockedUserIdsWithVersion(
	ctx context.Context,
	groupID int64,
	lastUpdatedDate *time.Time,
) ([]int64, *time.Time, error) {
	// Always query version (Java always queries and uses it in response)
	version, err := s.groupVersionService.QueryGroupBlocklistVersion(ctx, groupID)
	if err != nil {
		return nil, nil, err
	}
	if version == nil {
		// Java: switchIfEmpty -> alreadyUpToUpdate()
		return nil, nil, nil
	}
	if lastUpdatedDate != nil && !version.After(*lastUpdatedDate) {
		// Already up to date
		return nil, nil, nil
	}

	userIDs, err := s.blockedUserRepo.FindBlockedUserIds(ctx, groupID)
	if err != nil {
		return nil, nil, err
	}

	if len(userIDs) == 0 {
		// Java: throw NO_CONTENT
		return nil, nil, exception.NewTurmsError(int32(common_constant.ResponseStatusCode_NO_CONTENT), "No blocked users found")
	}

	return userIDs, version, nil
}

// QueryGroupBlockedUserInfosWithVersion queries blocked user infos with version control.
// Bug fixes:
// - Always query version (even when lastUpdatedDate is nil)
// - Return NO_CONTENT error when empty results
func (s *GroupBlocklistService) QueryGroupBlockedUserInfosWithVersion(
	ctx context.Context,
	groupID int64,
	lastUpdatedDate *time.Time,
) ([]po.GroupBlockedUser, *time.Time, error) {
	// Always query version (Java always queries and uses it in response)
	version, err := s.groupVersionService.QueryGroupBlocklistVersion(ctx, groupID)
	if err != nil {
		return nil, nil, err
	}
	if version == nil {
		// Java: switchIfEmpty -> alreadyUpToUpdate()
		return nil, nil, nil
	}
	if lastUpdatedDate != nil && !version.After(*lastUpdatedDate) {
		return nil, nil, nil
	}

	users, err := s.blockedUserRepo.FindBlockedUsersByGroupID(ctx, groupID)
	if err != nil {
		return nil, nil, err
	}

	if len(users) == 0 {
		// Java: throw NO_CONTENT
		return nil, nil, exception.NewTurmsError(int32(common_constant.ResponseStatusCode_NO_CONTENT), "No blocked users found")
	}

	return users, version, nil
}

func (s *GroupBlocklistService) AddBlockedUser(
	ctx context.Context,
	groupID int64,
	userID int64,
	requesterID int64,
	blockDate *time.Time,
) (*po.GroupBlockedUser, error) {
	if blockDate == nil {
		now := time.Now()
		blockDate = &now
	}
	blockedUser := &po.GroupBlockedUser{
		ID: po.GroupBlockedUserKey{
			GroupID: groupID,
			UserID:  userID,
		},
		BlockDate:   blockDate,
		RequesterID: requesterID,
	}
	err := s.blockedUserRepo.Insert(ctx, blockedUser)
	if err != nil {
		return nil, err
	}
	return blockedUser, nil
}

// UpdateBlockedUsers updates blocked users.
// Bug fix: Added per-key validation.
func (s *GroupBlocklistService) UpdateBlockedUsers(
	ctx context.Context,
	keys []po.GroupBlockedUserKey,
	blockDate *time.Time,
	requesterId *int64,
) error {
	// Validate keys
	for _, key := range keys {
		if key.GroupID <= 0 || key.UserID <= 0 {
			return exception.NewTurmsError(codes.IllegalArgument, "Invalid blocked user key")
		}
	}
	return s.blockedUserRepo.UpdateBlockedUsers(ctx, keys, blockDate, requesterId)
}

// DeleteBlockedUsers deletes blocked users.
// Bug fix: Added per-key validation.
func (s *GroupBlocklistService) DeleteBlockedUsers(ctx context.Context, keys []po.GroupBlockedUserKey) error {
	// Validate keys
	for _, key := range keys {
		if key.GroupID <= 0 || key.UserID <= 0 {
			return exception.NewTurmsError(codes.IllegalArgument, "Invalid blocked user key")
		}
	}
	return s.blockedUserRepo.DeleteBlockedUsers(ctx, keys)
}

func (s *GroupBlocklistService) QueryBlockedUsersWithFilter(
	ctx context.Context,
	groupIds []int64,
	userIds []int64,
	blockDateRange *turmsmongo.DateRange,
	requesterIds []int64,
	page *int,
	size *int,
) ([]po.GroupBlockedUser, error) {
	return s.blockedUserRepo.FindBlockedUsers(ctx, groupIds, userIds, blockDateRange, requesterIds, page, size)
}

func (s *GroupBlocklistService) QueryBlockedUsersWithPagination(
	ctx context.Context,
	page *int,
	size *int,
) ([]po.GroupBlockedUser, error) {
	return s.QueryBlockedUsersWithFilter(ctx, nil, nil, nil, nil, page, size)
}
