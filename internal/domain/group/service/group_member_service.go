package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"im.turms/server/internal/domain/common/cache"
	common_constant "im.turms/server/internal/domain/common/constant"
	group_constant "im.turms/server/internal/domain/group/constant"
	"im.turms/server/internal/domain/group/po"
	"im.turms/server/internal/domain/group/repository"
	"im.turms/server/internal/infra/exception"
	"im.turms/server/pkg/protocol"
)

var (
	ErrUnauthorized   = errors.New("unauthorized role operation")
	ErrRoleNotAllowed = errors.New("target role not allowed")
)

type GroupMemberService struct {
	groupRepo             *repository.GroupRepository
	groupMemberRepo       *repository.GroupMemberRepository
	groupBlocklistService *GroupBlocklistService
	groupVersionService   *GroupVersionService
	groupTypeService      *GroupTypeService
	groupService          *GroupService // Circular dependency? Usually handled via interface or delayed assignment
	memberCache           *cache.TTLCache[string, bool]
}

func NewGroupMemberService(
	groupRepo *repository.GroupRepository,
	groupMemberRepo *repository.GroupMemberRepository,
	groupVersionService *GroupVersionService,
	groupTypeService *GroupTypeService,
) *GroupMemberService {
	return &GroupMemberService{
		groupRepo:           groupRepo,
		groupMemberRepo:     groupMemberRepo,
		groupVersionService: groupVersionService,
		groupTypeService:    groupTypeService,
		memberCache:         cache.NewTTLCache[string, bool](1*time.Minute, 10*time.Second),
	}
}

func (s *GroupMemberService) SetGroupBlocklistService(groupBlocklistService *GroupBlocklistService) {
	s.groupBlocklistService = groupBlocklistService
}

// SetGroupService is a helper to break circular dependency if any.
func (s *GroupMemberService) SetGroupService(groupService *GroupService) {
	s.groupService = groupService
}

func (s *GroupMemberService) Close() {
	if s.memberCache != nil {
		s.memberCache.Close()
	}
}

// AddGroupMember adds a new member with the given role.
// If requesterID is nil, it's considered a system operation (no RBAC check).
func (s *GroupMemberService) AddGroupMember(ctx context.Context, groupID, userID int64, role protocol.GroupMemberRole, requesterID *int64, muteEndDate *time.Time) error {
	if requesterID != nil {
		// RBAC: check if requester is Owner/Manager
		reqRole, err := s.groupMemberRepo.FindGroupMemberRole(ctx, groupID, *requesterID)
		if err != nil {
			return err
		}
		if reqRole == nil || (*reqRole != protocol.GroupMemberRole_OWNER && *reqRole != protocol.GroupMemberRole_MANAGER) {
			return exception.NewTurmsError(int32(common_constant.ResponseStatusCode_NOT_GROUP_OWNER_OR_MANAGER_TO_ADD_GROUP_MEMBER), "Unauthorized to add group member")
		}
	}

	now := time.Now()
	member := &po.GroupMember{
		ID: po.GroupMemberKey{
			GroupID: groupID,
			UserID:  userID,
		},
		Role:        role,
		JoinDate:    &now,
		MuteEndDate: muteEndDate,
	}

	return s.groupMemberRepo.AddGroupMember(ctx, member)
}

// UpdateGroupMember updates a group member's information.
func (s *GroupMemberService) UpdateGroupMember(
	ctx context.Context,
	groupID int64,
	memberID int64,
	name *string,
	role *protocol.GroupMemberRole,
	joinDate *time.Time,
	muteEndDate *time.Time,
	session mongo.SessionContext,
	updateVersion bool,
) error {
	keys := []po.GroupMemberKey{{GroupID: groupID, UserID: memberID}}
	_, err := s.groupMemberRepo.UpdateGroupMembers(ctx, keys, name, role, joinDate, muteEndDate)
	if err != nil {
		return err
	}
	s.memberCache.Delete(fmt.Sprintf("member:%d:%d", groupID, memberID))
	s.memberCache.Delete(fmt.Sprintf("muted:%d:%d", groupID, memberID))

	if updateVersion {
		return s.groupVersionService.UpdateMembersVersion(ctx, groupID)
	}
	return nil
}

// IsMemberMuted is intensely called during routing.
func (s *GroupMemberService) IsMemberMuted(ctx context.Context, groupID, userID int64) (bool, error) {
	cacheKey := fmt.Sprintf("muted:%d:%d", groupID, userID)
	if muted, ok := s.memberCache.Get(cacheKey); ok {
		return muted, nil
	}

	muted, err := s.groupMemberRepo.IsMemberMuted(ctx, groupID, userID)
	if err != nil {
		return false, err
	}

	s.memberCache.Set(cacheKey, muted)
	return muted, nil
}

// @MappedFrom isOwner(@NotNull Long userId, @NotNull Long groupId, boolean preferCache)
func (s *GroupMemberService) IsOwner(ctx context.Context, userID, groupID int64) (bool, error) {
	role, err := s.FindGroupMemberRole(ctx, userID, groupID)
	if err != nil {
		return false, err
	}
	if role == nil {
		return false, nil
	}
	return *role == protocol.GroupMemberRole_OWNER, nil
}

// @MappedFrom isOwnerOrManager(@NotNull Long userId, @NotNull Long groupId, boolean preferCache)
func (s *GroupMemberService) IsOwnerOrManager(ctx context.Context, groupID, userID int64) (bool, error) {
	role, err := s.FindGroupMemberRole(ctx, groupID, userID)
	if err != nil {
		return false, err
	}
	if role == nil {
		return false, nil
	}
	return *role == protocol.GroupMemberRole_OWNER || *role == protocol.GroupMemberRole_MANAGER, nil
}

func (s *GroupMemberService) FindGroupMemberRole(ctx context.Context, groupID, userID int64) (*protocol.GroupMemberRole, error) {
	return s.groupMemberRepo.FindGroupMemberRole(ctx, groupID, userID)
}

// @MappedFrom queryGroupMemberRole(@NotNull Long userId, @NotNull Long groupId, boolean preferCache)
func (s *GroupMemberService) QueryGroupMemberRole(ctx context.Context, groupID, userID int64) (*protocol.GroupMemberRole, error) {
	return s.groupMemberRepo.FindGroupMemberRole(ctx, groupID, userID)
}

func (s *GroupMemberService) FindGroupMemberIDs(ctx context.Context, groupID int64) ([]int64, error) {
	return s.groupMemberRepo.FindGroupMemberIDs(ctx, groupID)
}

// FindExistentMemberGroupIds returns a list of group IDs where the user is an existent member.
// @MappedFrom findExistentMemberGroupIds(@NotEmpty Set<Long> groupIds, @NotNull Long userId)
func (s *GroupMemberService) FindExistentMemberGroupIds(ctx context.Context, groupIDs []int64, userID int64) ([]int64, error) {
	return s.groupMemberRepo.FindExistentMemberGroupIds(ctx, groupIDs, userID)
}

// IsAllowedToInviteUser checks if the inviter has permission to invite users to the group.
// @MappedFrom isAllowedToInviteUser(@NotNull Long groupId, @NotNull Long inviterId)
func (s *GroupMemberService) IsAllowedToInviteUser(ctx context.Context, groupID, inviterID int64) (bool, error) {
	role, err := s.groupMemberRepo.FindGroupMemberRole(ctx, groupID, inviterID)
	if err != nil {
		return false, err
	}

	typeID, err := s.groupService.QueryGroupTypeIdIfActiveAndNotDeleted(ctx, groupID)
	if err != nil {
		return false, err
	}
	if typeID == nil {
		return false, nil // Inactive or deleted
	}

	groupType, err := s.groupTypeService.FindGroupType(ctx, *typeID)
	if err != nil {
		return false, err
	}
	if groupType == nil {
		return false, nil
	}

	strategy := groupType.InvitationStrategy
	if strategy == group_constant.GroupInvitationStrategy_ALL || strategy == group_constant.GroupInvitationStrategy_ALL_REQUIRING_APPROVAL {
		return true, nil
	}

	if role == nil {
		return false, nil
	}

	switch strategy {
	case group_constant.GroupInvitationStrategy_OWNER_MANAGER_MEMBER, group_constant.GroupInvitationStrategy_OWNER_MANAGER_MEMBER_REQUIRING_APPROVAL:
		return *role == protocol.GroupMemberRole_OWNER || *role == protocol.GroupMemberRole_MANAGER || *role == protocol.GroupMemberRole_MEMBER, nil
	case group_constant.GroupInvitationStrategy_OWNER_MANAGER, group_constant.GroupInvitationStrategy_OWNER_MANAGER_REQUIRING_APPROVAL:
		return *role == protocol.GroupMemberRole_OWNER || *role == protocol.GroupMemberRole_MANAGER, nil
	case group_constant.GroupInvitationStrategy_OWNER, group_constant.GroupInvitationStrategy_OWNER_REQUIRING_APPROVAL:
		return *role == protocol.GroupMemberRole_OWNER, nil
	default:
		return false, nil
	}
}

// IsAllowedToBeInvited checks if the user can be invited (not a member, not blocked).
// @MappedFrom isAllowedToBeInvited(@NotNull Long groupId, @NotNull Long inviteeId)
func (s *GroupMemberService) IsAllowedToBeInvited(ctx context.Context, groupID, inviteeID int64) (bool, error) {
	isMember, err := s.IsGroupMember(ctx, groupID, inviteeID)
	if err != nil {
		return false, err
	}
	if isMember {
		return false, exception.NewTurmsError(int32(common_constant.ResponseStatusCode_SEND_GROUP_INVITATION_TO_GROUP_MEMBER), "User is already a group member")
	}

	isBlocked, err := s.groupBlocklistService.IsBlocked(ctx, groupID, inviteeID)
	if err != nil {
		return false, err
	}
	if isBlocked {
		return false, exception.NewTurmsError(int32(common_constant.ResponseStatusCode_SEND_GROUP_INVITATION_TO_BLOCKED_USER), "User is blocked from the group")
	}

	return true, nil
}

// IsAllowedToSendMessage checks if the user is allowed to send a message to the group.
// @MappedFrom isAllowedToSendMessage(@NotNull Long groupId, @NotNull Long senderId)
func (s *GroupMemberService) IsAllowedToSendMessage(ctx context.Context, groupID, senderID int64) (bool, error) {
	isGroupMuted, err := s.groupService.IsGroupMuted(ctx, groupID)
	if err != nil {
		return false, err
	}
	if isGroupMuted {
		return false, exception.NewTurmsError(int32(common_constant.ResponseStatusCode_SEND_MESSAGE_TO_MUTED_GROUP), "Group is muted")
	}

	isActive, err := s.groupService.IsGroupActiveAndNotDeleted(ctx, groupID)
	if err != nil {
		return false, err
	}
	if !isActive {
		return false, exception.NewTurmsError(int32(common_constant.ResponseStatusCode_SEND_MESSAGE_TO_INACTIVE_GROUP), "Group is inactive or deleted")
	}

	isMember, err := s.IsGroupMember(ctx, groupID, senderID)
	if err != nil {
		return false, err
	}

	if isMember {
		isMuted, err := s.IsMemberMuted(ctx, groupID, senderID)
		if err != nil {
			return false, err
		}
		if isMuted {
			return false, exception.NewTurmsError(int32(common_constant.ResponseStatusCode_MUTED_GROUP_MEMBER_SEND_MESSAGE), "Member is muted")
		}
		return true, nil
	}

	typeID, err := s.groupService.QueryGroupTypeIdIfActiveAndNotDeleted(ctx, groupID)
	if err != nil || typeID == nil {
		return false, err
	}
	groupType, err := s.groupTypeService.FindGroupType(ctx, *typeID)
	if err != nil || groupType == nil {
		return false, err
	}

	if !groupType.GuestSpeakable {
		return false, exception.NewTurmsError(int32(common_constant.ResponseStatusCode_NOT_SPEAKABLE_GROUP_GUEST_TO_SEND_MESSAGE), "Guest speaking is not allowed")
	}

	isBlocked, err := s.groupBlocklistService.IsBlocked(ctx, groupID, senderID)
	if err != nil {
		return false, err
	}
	if isBlocked {
		return false, exception.NewTurmsError(int32(common_constant.ResponseStatusCode_BLOCKED_USER_SEND_GROUP_MESSAGE), "User is blocked from the group")
	}

	return true, nil
}

// QueryGroupMemberKeyAndRolePairs retrieves the ID and roles for a list of users in a group.
// @MappedFrom queryGroupMemberKeyAndRolePairs(@NotNull Set<Long> userIds, @NotNull Long groupId)
func (s *GroupMemberService) QueryGroupMemberKeyAndRolePairs(ctx context.Context, userIDs []int64, groupID int64) ([]po.GroupMember, error) {
	return s.groupMemberRepo.FindGroupMemberKeyAndRolePairs(ctx, groupID, userIDs)
}

// IsOwnerOrManagerOrMember checks if a user holds an owner, manager, or regular member role.
// @MappedFrom isOwnerOrManagerOrMember(@NotNull Long userId, @NotNull Long groupId, boolean preferCache)
func (s *GroupMemberService) IsOwnerOrManagerOrMember(ctx context.Context, userID, groupID int64) (bool, error) {
	role, err := s.FindGroupMemberRole(ctx, groupID, userID)
	if err != nil {
		return false, err
	}
	if role == nil {
		return false, nil
	}
	return *role == protocol.GroupMemberRole_OWNER || *role == protocol.GroupMemberRole_MANAGER || *role == protocol.GroupMemberRole_MEMBER, nil
}

func (s *GroupMemberService) IsGroupMember(ctx context.Context, groupID, userID int64) (bool, error) {
	cacheKey := fmt.Sprintf("member:%d:%d", groupID, userID)
	if isMember, ok := s.memberCache.Get(cacheKey); ok {
		return isMember, nil
	}

	isMember, err := s.groupMemberRepo.IsGroupMember(ctx, groupID, userID)
	if err != nil {
		return false, err
	}

	s.memberCache.Set(cacheKey, isMember)
	return isMember, nil
}

func (s *GroupMemberService) UpdateGroupMemberRole(
	ctx context.Context,
	groupID, userID int64,
	role protocol.GroupMemberRole,
	session mongo.SessionContext,
) error {
	keys := []po.GroupMemberKey{{GroupID: groupID, UserID: userID}}
	_, err := s.groupMemberRepo.UpdateGroupMembers(ctx, keys, nil, &role, nil, nil)
	// Clear cache
	s.memberCache.Delete(fmt.Sprintf("member:%d:%d", groupID, userID))
	return err
}

// @MappedFrom deleteGroupMember(@NotNull Long groupId, @NotNull Long memberId, @Nullable ClientSession session, boolean updateGroupMembersVersion)
func (s *GroupMemberService) DeleteGroupMember(
	ctx context.Context,
	groupID, userID int64,
	session mongo.SessionContext,
	updateVersion bool,
) error {
	err := s.groupMemberRepo.RemoveGroupMember(ctx, groupID, userID)
	if err != nil {
		return err
	}
	// Clear cache
	s.memberCache.Delete(fmt.Sprintf("member:%d:%d", groupID, userID))
	s.memberCache.Delete(fmt.Sprintf("muted:%d:%d", groupID, userID))

	if updateVersion {
		return s.groupVersionService.UpdateMembersVersion(ctx, groupID)
	}
	return nil
}

// DeleteAllGroupMembers deletes all members of multiple groups.
// @MappedFrom deleteAllGroupMembers(@Nullable Set<Long> groupIds, @Nullable ClientSession session, boolean updateMembersVersion)
// @MappedFrom deleteAllGroupMembers(@Nullable Set<Long> groupIds, @Nullable ClientSession session)
func (s *GroupMemberService) DeleteAllGroupMembers(ctx context.Context, groupIDs []int64, session mongo.SessionContext, updateVersion bool) error {
	if len(groupIDs) == 0 {
		return nil
	}
	_, err := s.groupMemberRepo.DeleteByGroupIDs(ctx, groupIDs)
	if err != nil {
		return err
	}
	if updateVersion {
		for _, groupID := range groupIDs {
			_ = s.groupVersionService.UpdateMembersVersion(ctx, groupID)
		}
	}
	return nil
}

// @MappedFrom addGroupMembers(@NotNull Long groupId, @NotNull Set<Long> userIds, @NotNull @ValidGroupMemberRole GroupMemberRole groupMemberRole, @Nullable String name, @Nullable @PastOrPresent Date joinDate, @Nullable Date muteEndDate, @Nullable ClientSession session)
func (s *GroupMemberService) AddGroupMembers(
	ctx context.Context,
	groupID int64,
	userIDs []int64,
	role protocol.GroupMemberRole,
	name *string,
	joinTime *time.Time,
	muteEndDate *time.Time,
	session mongo.SessionContext,
) ([]po.GroupMember, error) {
	if joinTime == nil {
		now := time.Now()
		joinTime = &now
	}
	members := make([]po.GroupMember, len(userIDs))
	for i, userID := range userIDs {
		members[i] = po.GroupMember{
			ID: po.GroupMemberKey{
				GroupID: groupID,
				UserID:  userID,
			},
			Role:        role,
			Name:        name,
			JoinDate:    joinTime,
			MuteEndDate: muteEndDate,
		}
		// In a real implementation, we would call repository bulk insert
		err := s.groupMemberRepo.AddGroupMember(ctx, &members[i])
		if err != nil {
			return nil, err
		}
		// Clear cache
		s.memberCache.Delete(fmt.Sprintf("member:%d:%d", groupID, userID))
	}
	_ = s.groupVersionService.UpdateMembersVersion(ctx, groupID)
	return members, nil
}

// AuthAndAddGroupMembers adds members to a group after performing authorization and strategy checks.
// @MappedFrom authAndAddGroupMembers(@NotNull Long requesterId, @NotNull Long groupId, @NotNull Set<Long> userIds, @Nullable @ValidGroupMemberRole GroupMemberRole groupMemberRole, @Nullable String name, @Nullable Date muteEndDate, @Nullable ClientSession session)
func (s *GroupMemberService) AuthAndAddGroupMembers(
	ctx context.Context,
	requesterID int64,
	groupID int64,
	userIDs []int64,
	role protocol.GroupMemberRole,
	muteEndDate *time.Time,
) ([]po.GroupMember, error) {
	if len(userIDs) == 0 {
		return nil, nil
	}

	// 1. Check if group exists and get its type
	typeID, err := s.groupService.QueryGroupTypeIdIfActiveAndNotDeleted(ctx, groupID)
	if err != nil {
		return nil, err
	}
	if typeID == nil {
		return nil, exception.NewTurmsError(int32(common_constant.ResponseStatusCode_UPDATE_INFO_OF_NONEXISTENT_GROUP), "Group does not exist or is deleted")
	}

	groupType, err := s.groupTypeService.FindGroupType(ctx, *typeID)
	if err != nil {
		return nil, err
	}
	if groupType == nil {
		return nil, exception.NewTurmsError(int32(common_constant.ResponseStatusCode_SERVER_INTERNAL_ERROR), "Group type does not exist")
	}

	// 2. Authorization check
	requesterRole, err := s.groupMemberRepo.FindGroupMemberRole(ctx, groupID, requesterID)
	if err != nil {
		return nil, err
	}

	isOwnerOrManager := requesterRole != nil && (*requesterRole == protocol.GroupMemberRole_OWNER || *requesterRole == protocol.GroupMemberRole_MANAGER)

	if !isOwnerOrManager {
		if groupType.JoinStrategy != group_constant.GroupJoinStrategy_MEMBERSHIP_REQUEST {
			return nil, exception.NewTurmsError(int32(common_constant.ResponseStatusCode_NOT_GROUP_OWNER_OR_MANAGER_TO_ADD_GROUP_MEMBER), "Only owner or manager can add members for this group type")
		}
		if len(userIDs) > 1 || userIDs[0] != requesterID {
			return nil, exception.NewTurmsError(int32(common_constant.ResponseStatusCode_NOT_GROUP_OWNER_OR_MANAGER_TO_ADD_GROUP_MEMBER), "Cannot add other users as a non-member")
		}
	}

	// 3. Filter blocked users
	validUserIDs, err := s.groupBlocklistService.FilterBlockedUserIDs(ctx, groupID, userIDs)
	if err != nil {
		return nil, err
	}
	if len(validUserIDs) == 0 {
		return nil, exception.NewTurmsError(int32(common_constant.ResponseStatusCode_ADD_BLOCKED_USER_TO_GROUP), "All users are blocked")
	}

	// 4. Check group size limit
	currentCount, err := s.groupMemberRepo.CountMembers(ctx, groupID)
	if err != nil {
		return nil, err
	}
	if int(currentCount)+len(validUserIDs) > int(groupType.GroupSizeLimit) {
		return nil, exception.NewTurmsError(int32(common_constant.ResponseStatusCode_ADD_USER_TO_GROUP_WITH_SIZE_LIMIT_REACHED), "Group size limit exceeded")
	}

	// 5. Add members
	return s.AddGroupMembers(ctx, groupID, validUserIDs, role, nil, nil, muteEndDate, nil)
}

// AuthAndDeleteGroupMembers deletes members from a group after performing authorization checks.
// @MappedFrom authAndDeleteGroupMembers(@NotNull Long requesterId, @NotNull Long groupId, @NotNull Set<Long> memberIdsToDelete, @Nullable Long successorId, @Nullable Boolean quitAfterTransfer)
func (s *GroupMemberService) AuthAndDeleteGroupMembers(
	ctx context.Context,
	requesterID int64,
	groupID int64,
	userIDs []int64,
	successorID *int64,
	quitAfterTransfer bool,
) error {
	if len(userIDs) == 0 {
		return nil
	}

	isQuitting := false
	for _, uid := range userIDs {
		if uid == requesterID {
			isQuitting = true
			break
		}
	}

	requesterRole, err := s.groupMemberRepo.FindGroupMemberRole(ctx, groupID, requesterID)
	if err != nil {
		return err
	}
	if requesterRole == nil {
		return exception.NewTurmsError(int32(common_constant.ResponseStatusCode_NOT_GROUP_OWNER_OR_MANAGER_TO_REMOVE_GROUP_MEMBER), "Requester is not a group member")
	}

	if isQuitting && *requesterRole == protocol.GroupMemberRole_OWNER {
		if successorID != nil {
			err := s.groupService.AuthAndTransferGroupOwnership(ctx, requesterID, groupID, *successorID, quitAfterTransfer, nil)
			if err != nil {
				return err
			}
			if !quitAfterTransfer {
				newUserIDs := make([]int64, 0, len(userIDs)-1)
				for _, uid := range userIDs {
					if uid != requesterID {
						newUserIDs = append(newUserIDs, uid)
					}
				}
				userIDs = newUserIDs
			}
		} else {
			return s.groupService.AuthAndDeleteGroup(ctx, requesterID, groupID)
		}
	}

	if len(userIDs) == 0 {
		return nil
	}

	if !isQuitting || len(userIDs) > 1 {
		if *requesterRole != protocol.GroupMemberRole_OWNER && *requesterRole != protocol.GroupMemberRole_MANAGER {
			return exception.NewTurmsError(int32(common_constant.ResponseStatusCode_NOT_GROUP_OWNER_OR_MANAGER_TO_REMOVE_GROUP_MEMBER), "Only owner or manager can remove other members")
		}

		if *requesterRole == protocol.GroupMemberRole_MANAGER {
			targetRoles, err := s.groupMemberRepo.FindGroupMemberKeyAndRolePairs(ctx, groupID, userIDs)
			if err != nil {
				return err
			}
			for _, m := range targetRoles {
				if m.ID.UserID != requesterID && (m.Role == protocol.GroupMemberRole_OWNER || m.Role == protocol.GroupMemberRole_MANAGER) {
					return exception.NewTurmsError(int32(common_constant.ResponseStatusCode_NOT_GROUP_OWNER_OR_MANAGER_TO_REMOVE_GROUP_MEMBER), "Manager cannot remove owner or other managers")
				}
			}
		}
	}

	// DeleteGroupMembers was intended to be bulk version of DeleteGroupMember
	keys := make([]po.GroupMemberKey, len(userIDs))
	for i, userID := range userIDs {
		keys[i] = po.GroupMemberKey{GroupID: groupID, UserID: userID}
	}
	_, err = s.groupMemberRepo.DeleteByIds(ctx, keys)
	if err != nil {
		return err
	}
	for _, userID := range userIDs {
		s.memberCache.Delete(fmt.Sprintf("member:%d:%d", groupID, userID))
		s.memberCache.Delete(fmt.Sprintf("muted:%d:%d", groupID, userID))
	}
	return s.groupVersionService.UpdateMembersVersion(ctx, groupID)
}

// AuthAndUpdateGroupMember updates a group member's info after performing authorization checks.
// @MappedFrom authAndUpdateGroupMember(@NotNull Long requesterId, @NotNull Long groupId, @NotNull Long memberId, @Nullable String name, @Nullable @ValidGroupMemberRole GroupMemberRole role, @Nullable Date muteEndDate)
func (s *GroupMemberService) AuthAndUpdateGroupMember(
	ctx context.Context,
	requesterID int64,
	groupID int64,
	memberID int64,
	name *string,
	role *protocol.GroupMemberRole,
	muteEndDate *time.Time,
) error {
	// 1. Authorization check
	requesterRole, err := s.groupMemberRepo.FindGroupMemberRole(ctx, groupID, requesterID)
	if err != nil {
		return err
	}
	if requesterRole == nil {
		return exception.NewTurmsError(int32(common_constant.ResponseStatusCode_NOT_GROUP_MEMBER_TO_UPDATE_GROUP_MEMBER_INFO), "Requester is not a group member")
	}

	// 2. Permission check
	if requesterID != memberID {
		if *requesterRole != protocol.GroupMemberRole_OWNER && *requesterRole != protocol.GroupMemberRole_MANAGER {
			return exception.NewTurmsError(int32(common_constant.ResponseStatusCode_NOT_GROUP_OWNER_OR_MANAGER_TO_UPDATE_GROUP_MEMBER_INFO), "Only owner or manager can update other members' info")
		}
		if role != nil {
			if *requesterRole != protocol.GroupMemberRole_OWNER {
				return exception.NewTurmsError(int32(common_constant.ResponseStatusCode_NOT_GROUP_OWNER_TO_UPDATE_GROUP_MEMBER_ROLE), "Only owner can update others' roles")
			}
			if *role == protocol.GroupMemberRole_OWNER {
				return exception.NewTurmsError(int32(common_constant.ResponseStatusCode_ILLEGAL_ARGUMENT), "Cannot update member role to OWNER")
			}
		}
	} else {
		if role != nil {
			return exception.NewTurmsError(int32(common_constant.ResponseStatusCode_ILLEGAL_ARGUMENT), "Cannot update your own role")
		}
	}

	// 3. Update
	keys := []po.GroupMemberKey{{GroupID: groupID, UserID: memberID}}
	_, err = s.groupMemberRepo.UpdateGroupMembers(ctx, keys, name, role, nil, muteEndDate)
	if err != nil {
		return err
	}

	// 4. Invalidate cache
	s.memberCache.Delete(fmt.Sprintf("member:%d:%d", groupID, memberID))
	s.memberCache.Delete(fmt.Sprintf("muted:%d:%d", groupID, memberID))

	return s.groupVersionService.UpdateMembersVersion(ctx, groupID)
}

// QueryUserJoinedGroupIds returns the group IDs the user has joined.
// @MappedFrom queryUserJoinedGroupIds(@NotNull Long userId)
func (s *GroupMemberService) QueryUserJoinedGroupIds(ctx context.Context, userID int64) ([]int64, error) {
	return s.groupMemberRepo.FindUserJoinedGroupIDs(ctx, userID)
}

// QueryUsersJoinedGroupIds returns the group IDs joined by any of the specified users.
// @MappedFrom queryUsersJoinedGroupIds(@Nullable Set<Long> groupIds, @NotEmpty Set<Long> userIds, @Nullable Integer page, @Nullable Integer size)
func (s *GroupMemberService) QueryUsersJoinedGroupIds(ctx context.Context, groupIDs []int64, userIDs []int64, page, size *int) ([]int64, error) {
	return s.groupMemberRepo.FindUsersJoinedGroupIds(ctx, groupIDs, userIDs, page, size)
}

// QueryMemberIdsInUsersJoinedGroups returns member IDs for all groups that the users have joined.
// @MappedFrom queryMemberIdsInUsersJoinedGroups(@NotEmpty Set<Long> userIds, boolean preferCache)
func (s *GroupMemberService) QueryMemberIdsInUsersJoinedGroups(ctx context.Context, userIDs []int64) ([]int64, error) {
	groupIDs, err := s.QueryUsersJoinedGroupIds(ctx, nil, userIDs, nil, nil)
	if err != nil {
		return nil, err
	}
	if len(groupIDs) == 0 {
		return []int64{}, nil
	}
	return s.QueryGroupMemberIds(ctx, groupIDs)
}

// QueryGroupMemberIds retrieves member IDs for multiple groups.
// @MappedFrom queryGroupMemberIds(@NotEmpty Set<Long> groupIds, boolean preferCache)
func (s *GroupMemberService) QueryGroupMemberIds(ctx context.Context, groupIDs []int64) ([]int64, error) {
	return s.groupMemberRepo.FindMemberIdsByGroupIds(ctx, groupIDs)
}

// AuthAndQueryGroupMembers queries members with access checks.
// @MappedFrom authAndQueryGroupMembers(@NotNull Long requesterId, @NotNull Long groupId, @NotEmpty Set<Long> memberIds, boolean withStatus)
func (s *GroupMemberService) AuthAndQueryGroupMembers(ctx context.Context, requesterID int64, groupID int64, memberIDs []int64, withStatus bool) ([]*po.GroupMember, error) {
	isMember, err := s.IsGroupMember(ctx, groupID, requesterID)
	if err != nil {
		return nil, err
	}
	if !isMember { // A non-member might try to access something
		return nil, exception.NewTurmsError(int32(common_constant.ResponseStatusCode_NOT_GROUP_MEMBER_TO_QUERY_GROUP_MEMBER_INFO), "Only group members can query member info")
	}

	var members []po.GroupMember
	if len(memberIDs) > 0 {
		members, err = s.groupMemberRepo.FindGroupMembersWithIds(ctx, groupID, memberIDs)
	} else {
		members, err = s.groupMemberRepo.FindGroupMembers(ctx, groupID)
	}
	if err != nil {
		return nil, err
	}

	result := make([]*po.GroupMember, len(members))
	for i := range members {
		result[i] = &members[i]
	}
	return result, nil
}

// QueryGroupManagersAndOwnerId retrieves managers and owner of a group.
// @MappedFrom queryGroupManagersAndOwnerId(@NotNull Long groupId)
func (s *GroupMemberService) QueryGroupManagersAndOwnerId(ctx context.Context, groupID int64) ([]po.GroupMember, error) {
	return s.groupMemberRepo.FindGroupManagersAndOwnerId(ctx, groupID)
}

func (s *GroupMemberService) QueryGroupMembersWithPagination(ctx context.Context, page, size *int) ([]po.GroupMember, error) {
	return s.QueryGroupMembersWithFilter(ctx, nil, nil, nil, nil, nil, nil, nil, page, size)
}

func (s *GroupMemberService) QueryGroupMembersWithFilter(ctx context.Context, groupIds, userIds []int64, roles []int, joinDateStart, joinDateEnd, muteEndDateStart, muteEndDateEnd *time.Time, page, size *int) ([]po.GroupMember, error) {
	return s.groupMemberRepo.FindGroupsMembers(ctx, groupIds, userIds, roles, joinDateStart, joinDateEnd, muteEndDateStart, muteEndDateEnd, page, size)
}

// AuthAndQueryGroupMembersWithVersion queries group members with version control and auth checks.
// @MappedFrom authAndQueryGroupMembersWithVersion
func (s *GroupMemberService) AuthAndQueryGroupMembersWithVersion(
	ctx context.Context,
	requesterID int64,
	groupID int64,
	memberIDs []int64,
	lastUpdatedDate *time.Time,
) ([]*po.GroupMember, *time.Time, error) {
	isMember, err := s.IsGroupMember(ctx, groupID, requesterID)
	if err != nil {
		return nil, nil, err
	}
	if !isMember {
		return nil, nil, exception.NewTurmsError(int32(common_constant.ResponseStatusCode_NOT_GROUP_MEMBER_TO_QUERY_GROUP_MEMBER_INFO), "Only group members can query member info")
	}

	var version *time.Time
	if lastUpdatedDate != nil {
		v, err := s.groupVersionService.QueryGroupMembersVersion(ctx, groupID)
		if err != nil {
			return nil, nil, err
		}
		if v == nil || v.Before(*lastUpdatedDate) || v.Equal(*lastUpdatedDate) {
			return nil, nil, nil // Not modified since lastUpdatedDate
		}
		version = v
	}

	var members []po.GroupMember
	if len(memberIDs) > 0 {
		members, err = s.groupMemberRepo.FindGroupMembersWithIds(ctx, groupID, memberIDs)
	} else {
		members, err = s.groupMemberRepo.FindGroupMembers(ctx, groupID)
	}
	if err != nil {
		return nil, nil, err
	}

	// Convert to pointers for response
	result := make([]*po.GroupMember, len(members))
	for i := range members {
		result[i] = &members[i]
	}

	return result, version, nil
}
