package service

import (
	"context"
	"time"

	"im.turms/server/internal/domain/common/infra/idgen"
	group_constant "im.turms/server/internal/domain/group/constant"
	"im.turms/server/internal/domain/group/po"
	"im.turms/server/internal/domain/group/repository"
	user_service "im.turms/server/internal/domain/user/service"
	"im.turms/server/internal/infra/exception"
	"im.turms/server/pkg/codes"
	"im.turms/server/pkg/protocol"
)

type GroupInvitationService struct {
	invRepo             repository.GroupInvitationRepository
	groupMemberService  *GroupMemberService
	groupService        *GroupService
	groupTypeService    *GroupTypeService
	groupVersionService *GroupVersionService
	userVersionService  *user_service.UserVersionService
	idGen               *idgen.SnowflakeIdGenerator
}

func NewGroupInvitationService(
	invRepo repository.GroupInvitationRepository,
	groupMemberService *GroupMemberService,
	groupService *GroupService,
	groupTypeService *GroupTypeService,
	groupVersionService *GroupVersionService,
	userVersionService *user_service.UserVersionService,
	idGen *idgen.SnowflakeIdGenerator,
) *GroupInvitationService {
	return &GroupInvitationService{
		invRepo:             invRepo,
		groupMemberService:  groupMemberService,
		groupService:        groupService,
		groupTypeService:    groupTypeService,
		groupVersionService: groupVersionService,
		userVersionService:  userVersionService,
		idGen:               idGen,
	}
}

// @MappedFrom authAndCreateGroupInvitation(@NotNull Long groupId, @NotNull Long inviterId, @NotNull Long inviteeId, @Nullable String content)
func (s *GroupInvitationService) AuthAndCreateGroupInvitation(
	ctx context.Context,
	requesterID int64,
	groupID int64,
	inviteeID int64,
	content string,
) (*po.GroupInvitation, error) {
	// 1. Check if group is active
	typeID, err := s.groupService.QueryGroupTypeIdIfActiveAndNotDeleted(ctx, groupID)
	if err != nil {
		return nil, err
	}
	if typeID == nil {
		return nil, exception.NewTurmsError(codes.AddUserToInactiveGroup, "Group does not exist or is inactive")
	}

	// 2. Check if invitee is already a member
	isMember, err := s.groupMemberService.IsGroupMember(ctx, groupID, inviteeID)
	if err != nil {
		return nil, err
	}
	if isMember {
		return nil, exception.NewTurmsError(codes.AddUserToGroupWithSizeLimitReached, "Invitee is already a group member") // Actually should be a more specific code if exists
	}

	// 3. Check invitation strategy
	groupType, err := s.groupTypeService.FindGroupType(ctx, *typeID)
	if err != nil {
		return nil, err
	}
	if groupType == nil {
		return nil, exception.NewTurmsError(codes.ServerInternalError, "Group type not found")
	}

	strategy := groupType.InvitationStrategy
	requesterRole, err := s.groupMemberService.QueryGroupMemberRole(ctx, groupID, requesterID)
	if err != nil {
		return nil, err
	}

	allowed := false
	switch strategy {
	case group_constant.GroupInvitationStrategy_ALL, group_constant.GroupInvitationStrategy_ALL_REQUIRING_APPROVAL:
		allowed = true
	case group_constant.GroupInvitationStrategy_OWNER_MANAGER_MEMBER, group_constant.GroupInvitationStrategy_OWNER_MANAGER_MEMBER_REQUIRING_APPROVAL:
		allowed = requesterRole != nil
	case group_constant.GroupInvitationStrategy_OWNER_MANAGER, group_constant.GroupInvitationStrategy_OWNER_MANAGER_REQUIRING_APPROVAL:
		allowed = requesterRole != nil && (*requesterRole == protocol.GroupMemberRole_OWNER || *requesterRole == protocol.GroupMemberRole_MANAGER)
	case group_constant.GroupInvitationStrategy_OWNER, group_constant.GroupInvitationStrategy_OWNER_REQUIRING_APPROVAL:
		allowed = requesterRole != nil && *requesterRole == protocol.GroupMemberRole_OWNER
	}

	if !allowed {
		return nil, exception.NewTurmsError(codes.NotGroupMemberToSendGroupInvitation, "Not allowed to send group invitation")
	}

	// 4. Create invitation
	now := time.Now()
	id := s.idGen.NextLargeGapId()
	status := po.RequestStatusPending
	if !strategy.RequiresApproval() {
		status = po.RequestStatusAccepted
	}

	inv := &po.GroupInvitation{
		ID:           id,
		GroupID:      groupID,
		InviterID:    requesterID,
		InviteeID:    inviteeID,
		Content:      content,
		Status:       status,
		CreationDate: now,
	}

	err = s.invRepo.Insert(ctx, inv)
	if err != nil {
		return nil, err
	}

	// 5. If accepted automatically, add to members
	if status == po.RequestStatusAccepted {
		err = s.groupMemberService.AddGroupMember(ctx, groupID, inviteeID, protocol.GroupMemberRole_MEMBER, nil, nil)
		if err != nil {
			return nil, err
		}
	}

	// 6. Update versions
	_ = s.groupVersionService.UpdateInvitationsVersion(ctx, groupID)
	_ = s.userVersionService.UpdateSentGroupInvitationsVersion(ctx, requesterID)
	_ = s.userVersionService.UpdateReceivedGroupInvitationsVersion(ctx, inviteeID)

	return inv, nil
}

// @MappedFrom authAndRecallPendingGroupInvitation(@NotNull Long requesterId, @NotNull Long invitationId)
func (s *GroupInvitationService) AuthAndRecallPendingGroupInvitation(
	ctx context.Context,
	requesterID int64,
	invitationID int64,
) error {
	groupID, inviterID, inviteeID, status, err := s.invRepo.FindGroupIdAndInviterIdAndInviteeIdAndStatus(ctx, invitationID)
	if err != nil {
		return err
	}
	if status != po.RequestStatusPending {
		return exception.NewTurmsError(codes.RecallNonPendingGroupInvitation, "Cannot recall non-pending invitation")
	}

	// Permission check: requester must be inviter, owner, or manager
	if inviterID != requesterID {
		role, err := s.groupMemberService.QueryGroupMemberRole(ctx, groupID, requesterID)
		if err != nil {
			return err
		}
		if role == nil || (*role != protocol.GroupMemberRole_OWNER && *role != protocol.GroupMemberRole_MANAGER) {
			return exception.NewTurmsError(codes.NotGroupOwnerOrManagerOrSenderToRecallGroupInvitation, "No permission to recall")
		}
	}

	success, err := s.invRepo.UpdateStatusIfPending(ctx, invitationID, po.RequestStatusCanceled, nil, time.Now())
	if err != nil {
		return err
	}
	if !success {
		return exception.NewTurmsError(codes.RecallNonPendingGroupInvitation, "Recall failed")
	}

	_ = s.groupVersionService.UpdateInvitationsVersion(ctx, groupID)
	_ = s.userVersionService.UpdateSentGroupInvitationsVersion(ctx, inviterID)
	_ = s.userVersionService.UpdateReceivedGroupInvitationsVersion(ctx, inviteeID)

	return nil
}

// @MappedFrom authAndHandleInvitation(@NotNull Long requesterId, @NotNull Long invitationId, @NotNull @ValidResponseAction ResponseAction action, @Nullable String reason)
func (s *GroupInvitationService) AuthAndHandleInvitation(
	ctx context.Context,
	requesterID int64,
	invitationID int64,
	status po.RequestStatus,
	reason string,
) error {
	if status != po.RequestStatusAccepted && status != po.RequestStatusDeclined && status != po.RequestStatusIgnored {
		return exception.NewTurmsError(codes.IllegalArgument, "Invalid response status")
	}

	inviteeID, groupID, _, currentStatus, err := s.invRepo.FindInviteeIdAndGroupIdAndCreationDateAndStatus(ctx, invitationID)
	if err != nil {
		return err
	}
	if inviteeID != requesterID {
		return exception.NewTurmsError(codes.NotInviteeToUpdateGroupInvitation, "Only invitee can handle invitation")
	}
	if currentStatus != po.RequestStatusPending {
		return exception.NewTurmsError(codes.UpdateNonPendingGroupInvitation, "Invitation is already handled")
	}

	success, err := s.invRepo.UpdateStatusIfPending(ctx, invitationID, status, &reason, time.Now())
	if err != nil {
		return err
	}
	if !success {
		return exception.NewTurmsError(codes.UpdateNonPendingGroupInvitation, "Handle failed")
	}

	if status == po.RequestStatusAccepted {
		err = s.groupMemberService.AddGroupMember(ctx, groupID, inviteeID, protocol.GroupMemberRole_MEMBER, nil, nil)
		if err != nil {
			return err
		}
	}

	// Update versions
	inv, _ := s.invRepo.FindByID(ctx, invitationID)
	_ = s.groupVersionService.UpdateInvitationsVersion(ctx, groupID)
	if inv != nil {
		_ = s.userVersionService.UpdateSentGroupInvitationsVersion(ctx, inv.InviterID)
	}
	_ = s.userVersionService.UpdateReceivedGroupInvitationsVersion(ctx, inviteeID)

	return nil
}

// @MappedFrom queryInvitations(@Nullable Set<Long> ids, @Nullable Set<Long> groupIds, @Nullable Set<Long> inviterIds, @Nullable Set<Long> inviteeIds, @Nullable Set<RequestStatus> statuses, @Nullable DateRange creationDateRange, @Nullable DateRange responseDateRange, @Nullable DateRange expirationDateRange, @Nullable Integer page, @Nullable Integer size)
func (s *GroupInvitationService) QueryInvitations(ctx context.Context, groupID *int64, inviterID *int64, inviteeID *int64, status *po.RequestStatus, lastUpdatedDate *time.Time, page int, size int) ([]*po.GroupInvitation, error) {
	return s.invRepo.FindInvitations(ctx, groupID, inviterID, inviteeID, status, lastUpdatedDate, page, size)
}

// @MappedFrom queryUserGroupInvitationsWithVersion(@NotNull Long userId, boolean areSentByUser, @Nullable Date lastUpdatedDate)
func (s *GroupInvitationService) QueryUserGroupInvitationsWithVersion(ctx context.Context, userID int64, areSentInvitations bool, lastUpdatedDate *time.Time) (*po.GroupInvitationsWithVersion, error) {
	var version *time.Time
	var err error
	if areSentInvitations {
		version, err = s.userVersionService.QuerySentGroupInvitationsVersion(ctx, userID)
	} else {
		version, err = s.userVersionService.QueryReceivedGroupInvitationsVersion(ctx, userID)
	}
	if err != nil {
		return nil, err
	}
	if lastUpdatedDate != nil && version != nil && !version.After(*lastUpdatedDate) {
		return nil, exception.NewTurmsError(codes.AlreadyUpToDate, "Invitations are already up to date")
	}
	var invs []*po.GroupInvitation
	if areSentInvitations {
		invs, err = s.invRepo.FindInvitations(ctx, nil, &userID, nil, nil, nil, 0, 1000)
	} else {
		invs, err = s.invRepo.FindInvitations(ctx, nil, nil, &userID, nil, nil, 0, 1000)
	}
	if err != nil {
		return nil, err
	}
	return &po.GroupInvitationsWithVersion{
		GroupInvitations: invs,
		LastUpdatedDate:  version,
	}, nil
}

// @MappedFrom authAndQueryGroupInvitationsWithVersion(@NotNull Long userId, @NotNull Long groupId, @Nullable Date lastUpdatedDate)
func (s *GroupInvitationService) AuthAndQueryGroupInvitationsWithVersion(ctx context.Context, requesterID int64, groupID int64, lastUpdatedDate *time.Time) (*po.GroupInvitationsWithVersion, error) {
	role, err := s.groupMemberService.QueryGroupMemberRole(ctx, groupID, requesterID)
	if err != nil {
		return nil, err
	}
	if role == nil || (*role != protocol.GroupMemberRole_OWNER && *role != protocol.GroupMemberRole_MANAGER) {
		return nil, exception.NewTurmsError(codes.NotGroupOwnerOrManagerToQueryGroupInvitation, "No permission to query")
	}
	version, err := s.groupVersionService.QueryGroupInvitationsVersion(ctx, groupID)
	if err != nil {
		return nil, err
	}
	if lastUpdatedDate != nil && version != nil && !version.After(*lastUpdatedDate) {
		return nil, exception.NewTurmsError(codes.AlreadyUpToDate, "Invitations are already up to date")
	}
	invs, err := s.invRepo.FindInvitations(ctx, &groupID, nil, nil, nil, nil, 0, 1000)
	if err != nil {
		return nil, err
	}
	return &po.GroupInvitationsWithVersion{
		GroupInvitations: invs,
		LastUpdatedDate:  version,
	}, nil
}

func (s *GroupInvitationService) CountInvitations(ctx context.Context, groupID, inviterID, inviteeID *int64, status *po.RequestStatus, lastUpdatedDate *time.Time) (int64, error) {
	return s.invRepo.CountInvitations(ctx, groupID, inviterID, inviteeID, status, lastUpdatedDate)
}

func (s *GroupInvitationService) DeleteInvitations(ctx context.Context, ids []int64) (int64, error) {
	return s.invRepo.DeleteInvitations(ctx, ids)
}

func (s *GroupInvitationService) UpdateInvitations(ctx context.Context, ids []int64, inviterID, inviteeID *int64, content *string, status *po.RequestStatus, creationDate, responseDate *time.Time) (int64, error) {
	return s.invRepo.UpdateInvitations(ctx, ids, inviterID, inviteeID, content, status, creationDate, responseDate)
}

func (s *GroupInvitationService) QueryGroupIdAndInviterIdAndInviteeIdAndStatus(ctx context.Context, invitationID int64) (int64, int64, int64, po.RequestStatus, error) {
	return s.invRepo.FindGroupIdAndInviterIdAndInviteeIdAndStatus(ctx, invitationID)
}

func (s *GroupInvitationService) QueryGroupIdAndInviteeIdAndStatus(ctx context.Context, invitationID int64) (int64, int64, po.RequestStatus, error) {
	return s.invRepo.FindGroupIdAndInviteeIdAndStatus(ctx, invitationID)
}

func (s *GroupInvitationService) QueryGroupInvitationsByInviteeId(ctx context.Context, inviteeID int64) ([]po.GroupInvitation, error) {
	return s.invRepo.FindInvitationsByInviteeID(ctx, inviteeID)
}

// Java implementation queryGroupInvitationsByInviterId does the exact same thing
func (s *GroupInvitationService) QueryGroupInvitationsByInviterId(ctx context.Context, inviterID int64) ([]po.GroupInvitation, error) {
	// The repo method is not fully implemented but defined in interface. So we will defer to the generic FindInvitations.
	invs, err := s.invRepo.FindInvitations(ctx, nil, &inviterID, nil, nil, nil, 0, 1000)
	if err != nil {
		return nil, err
	}
	var res []po.GroupInvitation
	for _, v := range invs {
		res = append(res, *v)
	}
	return res, nil
}

func (s *GroupInvitationService) QueryGroupInvitationsByGroupId(ctx context.Context, groupID int64) ([]po.GroupInvitation, error) {
	return s.invRepo.FindInvitationsByGroupID(ctx, groupID)
}

func (s *GroupInvitationService) QueryInviteeIdAndGroupIdAndCreationDateAndStatusByInvitationId(ctx context.Context, invitationID int64) (int64, int64, time.Time, po.RequestStatus, error) {
	return s.invRepo.FindInviteeIdAndGroupIdAndCreationDateAndStatus(ctx, invitationID)
}

func (s *GroupInvitationService) UpdatePendingInvitationStatus(ctx context.Context, invitationID int64, requestStatus po.RequestStatus, reason *string) (bool, error) {
	return s.invRepo.UpdateStatusIfPending(ctx, invitationID, requestStatus, reason, time.Now())
}

// Backward Compatibility Aliases

func (s *GroupInvitationService) CreateInvitation(ctx context.Context, groupID int64, inviterID int64, inviteeID int64, content string) (*po.GroupInvitation, error) {
	return s.AuthAndCreateGroupInvitation(ctx, inviterID, groupID, inviteeID, content)
}

func (s *GroupInvitationService) RecallPendingInvitation(ctx context.Context, invitationID int64, inviterID int64) (bool, error) {
	err := s.AuthAndRecallPendingGroupInvitation(ctx, inviterID, invitationID)
	return err == nil, err
}

func (s *GroupInvitationService) ReplyToInvitation(ctx context.Context, invitationID int64, inviteeID int64, accept bool) (bool, error) {
	status := po.RequestStatusDeclined
	if accept {
		status = po.RequestStatusAccepted
	}
	err := s.AuthAndHandleInvitation(ctx, inviteeID, invitationID, status, "")
	return err == nil, err
}
