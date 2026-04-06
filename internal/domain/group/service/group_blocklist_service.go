package service

import (
	"context"
	"time"

	common_constant "im.turms/server/internal/domain/common/constant"
	"im.turms/server/internal/domain/group/po"
	"im.turms/server/internal/domain/group/repository"
	"im.turms/server/internal/infra/exception"
	turmsmongo "im.turms/server/internal/storage/mongo"
	"im.turms/server/pkg/protocol"
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

func (s *GroupBlocklistService) AuthAndUnblockUser(ctx context.Context, requesterID int64, groupID int64, userID int64) error {	now := time.Now()
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
func (s *GroupBlocklistService) AuthAndBlockUser(
	ctx context.Context,
	requesterID int64,
	groupID int64,
	userID int64,
) error {
	// 1. Authorization check
	requesterRole, err := s.groupMemberService.FindGroupMemberRole(ctx, requesterID, groupID)
	if err != nil {
		return err
	}
	if requesterRole == nil || (*requesterRole != protocol.GroupMemberRole_OWNER && *requesterRole != protocol.GroupMemberRole_MANAGER) {
		return exception.NewTurmsError(int32(common_constant.ResponseStatusCode_NOT_GROUP_OWNER_OR_MANAGER_TO_ADD_BLOCKED_USER), "Only owner or manager can block users")
	}

	// 2. Target check
	targetRole, err := s.groupMemberService.FindGroupMemberRole(ctx, userID, groupID)
	if err != nil {
		return err
	}
	if targetRole != nil {
		if *requesterRole == protocol.GroupMemberRole_MANAGER && (*targetRole == protocol.GroupMemberRole_OWNER || *targetRole == protocol.GroupMemberRole_MANAGER) {
			return exception.NewTurmsError(int32(common_constant.ResponseStatusCode_NOT_GROUP_OWNER_TO_REMOVE_GROUP_OWNER_OR_MANAGER), "Manager cannot block owner or other managers")
		}
	}

	// 3. Block
	err = s.BlockUser(ctx, groupID, userID, requesterID)
	if err != nil {
		return err
	}

	// 4. Remove from group and update version
	if targetRole != nil {
		// No session context passed here, matching existing GroupMemberService.DeleteGroupMember calls
		_ = s.groupMemberService.DeleteGroupMember(ctx, groupID, userID, nil, false)
	}
	return s.groupVersionService.UpdateBlocklistVersion(ctx, groupID)
}

// AuthAndUnblockUser unblocks a user from a group after performing authorization checks.
func (s *GroupBlocklistService) AuthAndUnblockUser(
	ctx context.Context,
	requesterID int64,
	groupID int64,
	userID int64,
) error {
	// 1. Authorization check
	requesterRole, err := s.groupMemberService.FindGroupMemberRole(ctx, requesterID, groupID)
	if err != nil {
		return err
	}
	if requesterRole == nil || (*requesterRole != protocol.GroupMemberRole_OWNER && *requesterRole != protocol.GroupMemberRole_MANAGER) {
		return exception.NewTurmsError(int32(common_constant.ResponseStatusCode_NOT_GROUP_OWNER_OR_MANAGER_TO_REMOVE_BLOCKED_USER), "Only owner or manager can unblock users")
	}

	// 2. Unblock
	err = s.UnblockUser(ctx, groupID, userID)
	if err != nil {
		return err
	}

	return s.groupVersionService.UpdateBlocklistVersion(ctx, groupID)
}

// @MappedFrom unblockUser(@NotNull Long requesterId, @NotNull Long groupId, @NotNull Long userIdToUnblock, @Nullable ClientSession session, boolean updateBlocklistVersion)
func (s *GroupBlocklistService) UnblockUser(ctx context.Context, groupID int64, userID int64) error {
	return s.blockedUserRepo.Delete(ctx, groupID, userID)
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

	// 2. Version check
	var version *time.Time
	if lastUpdatedDate != nil {
		v, err := s.groupVersionService.QueryGroupBlocklistVersion(ctx, groupID)
		if err != nil {
			return nil, nil, err
		}
		if v == nil || v.Before(*lastUpdatedDate) || v.Equal(*lastUpdatedDate) {
			return nil, nil, nil
		}
		version = v
	}

	// 3. Query
	users, err := s.blockedUserRepo.FindBlockedUsersByGroupID(ctx, groupID)
	if err != nil {
		return nil, nil, err
	}

	var userIDs []int64
	for _, u := range users {
		userIDs = append(userIDs, u.ID.UserID)
	}

	return userIDs, version, nil
}

// AuthAndQueryGroupBlockedUserInfos queries blocked user infos with auth check.
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

	// 2. Version check
	var version *time.Time
	if lastUpdatedDate != nil {
		v, err := s.groupVersionService.QueryGroupBlocklistVersion(ctx, groupID)
		if err != nil {
			return nil, nil, err
		}
		if v == nil || v.Before(*lastUpdatedDate) || v.Equal(*lastUpdatedDate) {
			return nil, nil, nil
		}
		version = v
	}

	// 3. Query
	users, err := s.blockedUserRepo.FindBlockedUsersByGroupID(ctx, groupID)
	if err != nil {
		return nil, nil, err
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

func (s *GroupBlocklistService) QueryGroupBlockedUserIdsWithVersion(
	ctx context.Context,
	groupID int64,
	lastUpdatedDate *time.Time,
) ([]int64, *time.Time, error) {
	var version *time.Time
	if lastUpdatedDate != nil {
		v, err := s.groupVersionService.QueryGroupBlocklistVersion(ctx, groupID)
		if err != nil {
			return nil, nil, err
		}
		if v == nil || v.Before(*lastUpdatedDate) || v.Equal(*lastUpdatedDate) {
			return nil, nil, nil
		}
		version = v
	}

	userIDs, err := s.blockedUserRepo.FindBlockedUserIds(ctx, groupID)
	if err != nil {
		return nil, nil, err
	}
	return userIDs, version, nil
}

func (s *GroupBlocklistService) QueryGroupBlockedUserInfosWithVersion(
	ctx context.Context,
	groupID int64,
	lastUpdatedDate *time.Time,
) ([]po.GroupBlockedUser, *time.Time, error) {
	var version *time.Time
	if lastUpdatedDate != nil {
		v, err := s.groupVersionService.QueryGroupBlocklistVersion(ctx, groupID)
		if err != nil {
			return nil, nil, err
		}
		if v == nil || v.Before(*lastUpdatedDate) || v.Equal(*lastUpdatedDate) {
			return nil, nil, nil
		}
		version = v
	}

	users, err := s.blockedUserRepo.FindBlockedUsersByGroupID(ctx, groupID)
	if err != nil {
		return nil, nil, err
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

func (s *GroupBlocklistService) UpdateBlockedUsers(
	ctx context.Context,
	keys []po.GroupBlockedUserKey,
	blockDate *time.Time,
	requesterId *int64,
) error {
	return s.blockedUserRepo.UpdateBlockedUsers(ctx, keys, blockDate, requesterId)
}

func (s *GroupBlocklistService) DeleteBlockedUsers(ctx context.Context, keys []po.GroupBlockedUserKey) error {
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
