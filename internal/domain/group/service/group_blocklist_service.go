package service

import (
	"context"
	"time"

	common_constant "im.turms/server/internal/domain/common/constant"
	"im.turms/server/internal/domain/group/po"
	"im.turms/server/internal/domain/group/repository"
	"im.turms/server/internal/infra/exception"
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
func (s *GroupBlocklistService) QueryBlockedUsers(ctx context.Context, groupID int64) ([]po.GroupBlockedUser, error) {
	return s.blockedUserRepo.FindBlockedUsersByGroupID(ctx, groupID)
}

// @MappedFrom isBlocked(Long ownerId, Long relatedUserId)
// @MappedFrom isBlocked(@NotNull Long groupId, @NotNull Long userId)
// @MappedFrom isBlocked(@NotNull Long ownerId, @NotNull Long relatedUserId, boolean preferCache)
func (s *GroupBlocklistService) IsBlocked(ctx context.Context, groupID int64, userID int64) (bool, error) {
	return s.blockedUserRepo.Exists(ctx, groupID, userID)
}

func (s *GroupBlocklistService) FilterBlockedUserIDs(ctx context.Context, groupID int64, userIDs []int64) ([]int64, error) {
	return s.blockedUserRepo.FilterBlockedUserIDs(ctx, groupID, userIDs)
}
