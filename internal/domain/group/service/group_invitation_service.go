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
	turmsmongo "im.turms/server/internal/storage/mongo"
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
// Bug fixes:
// - Strategy that doesn't require approval now REJECTS (Java parity), not auto-accepts
// - Added blocklist check for invitee (Java: isAllowedToBeInvited)
// - Added pending invitation duplicate check (Java parity)
// - Fixed error code for "invitee is member" (SendGroupInvitationToGroupMember)
// - Inviter must be a group member for ALL strategy (Java: queryGroupMemberRole empty → error)
// - Always creates with PENDING status (Java never auto-accepts in this method)
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

	// 2. Check invitation strategy
	groupType, err := s.groupTypeService.FindGroupType(ctx, *typeID)
	if err != nil {
		return nil, err
	}
	if groupType == nil {
		return nil, exception.NewTurmsError(codes.ServerInternalError, "Group type not found")
	}

	strategy := groupType.InvitationStrategy

	// 3. Java parity: if strategy doesn't require approval, REJECT the request
	// Java: if (!strategy.requiresApproval()) return error SEND_GROUP_INVITATION_TO_GROUP_NOT_REQUIRING_USERS_APPROVAL
	if !strategy.RequiresApproval() {
		return nil, exception.NewTurmsError(codes.SendGroupInvitationToGroupNotRequiringUsersApproval, "Cannot send group invitation for group not requiring user approval")
	}

	// 4. Check inviter role permissions
	// Java: queryGroupMemberRole → if empty → GROUP_INVITER_NOT_MEMBER (regardless of strategy)
	requesterRole, err := s.groupMemberService.QueryGroupMemberRole(ctx, groupID, requesterID)
	if err != nil {
		return nil, err
	}
	// Java: inviter must be a member (role != null) for all strategies
	if requesterRole == nil {
		return nil, exception.NewTurmsError(codes.NotGroupMemberToSendGroupInvitation, "Inviter is not a group member")
	}

	allowed := false
	switch strategy {
	case group_constant.GroupInvitationStrategy_ALL_REQUIRING_APPROVAL:
		allowed = true
	case group_constant.GroupInvitationStrategy_OWNER_MANAGER_MEMBER_REQUIRING_APPROVAL:
		allowed = true // already checked requesterRole != nil
	case group_constant.GroupInvitationStrategy_OWNER_MANAGER_REQUIRING_APPROVAL:
		allowed = *requesterRole == protocol.GroupMemberRole_OWNER || *requesterRole == protocol.GroupMemberRole_MANAGER
	case group_constant.GroupInvitationStrategy_OWNER_REQUIRING_APPROVAL:
		allowed = *requesterRole == protocol.GroupMemberRole_OWNER
	}

	if !allowed {
		return nil, exception.NewTurmsError(codes.NotGroupMemberToSendGroupInvitation, "Not allowed to send group invitation")
	}

	// 5. Check if invitee can be invited (not member + not blocked)
	// Java: groupMemberService.isAllowedToBeInvited(groupId, inviteeId)
	isMember, err := s.groupMemberService.IsGroupMember(ctx, groupID, inviteeID)
	if err != nil {
		return nil, err
	}
	if isMember {
		return nil, exception.NewTurmsError(codes.SendGroupInvitationToGroupMember, "Invitee is already a group member")
	}

	isBlocked, err := s.groupMemberService.IsBlocked(ctx, groupID, inviteeID)
	if err != nil {
		return nil, err
	}
	if isBlocked {
		return nil, exception.NewTurmsError(codes.SendGroupInvitationToBlockedUser, "Invitee has been blocked by the group")
	}

	// 6. Check for existing pending invitation
	hasPending, err := s.invRepo.HasPendingInvitation(ctx, groupID, inviteeID)
	if err != nil {
		return nil, err
	}
	if hasPending {
		return nil, exception.NewTurmsError(codes.RecordContainsDuplicateKey, "A pending invitation already exists")
	}

	// 7. Create invitation — always PENDING status (Java never auto-accepts)
	now := time.Now()
	id := s.idGen.NextLargeGapId()

	inv := &po.GroupInvitation{
		ID:           id,
		GroupID:      groupID,
		InviterID:    requesterID,
		InviteeID:    inviteeID,
		Content:      content,
		Status:       po.RequestStatusPending,
		CreationDate: now,
	}

	err = s.invRepo.Insert(ctx, inv)
	if err != nil {
		return nil, err
	}

	// 8. Update versions
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
	var groupIds, inviterIds, inviteeIds []int64
	var statuses []po.RequestStatus
	if groupID != nil {
		groupIds = []int64{*groupID}
	}
	if inviterID != nil {
		inviterIds = []int64{*inviterID}
	}
	if inviteeID != nil {
		inviteeIds = []int64{*inviteeID}
	}
	if status != nil {
		statuses = []po.RequestStatus{*status}
	}
	var creationDateRange *turmsmongo.DateRange
	if lastUpdatedDate != nil {
		creationDateRange = &turmsmongo.DateRange{Start: lastUpdatedDate}
	}
	p, sz := page, size
	return s.invRepo.FindInvitations(ctx, nil, groupIds, inviterIds, inviteeIds, statuses, creationDateRange, nil, nil, &p, &sz)
}

func (s *GroupInvitationService) QueryInvitationsWithPagination(ctx context.Context, page, size *int) ([]*po.GroupInvitation, error) {
	return s.QueryInvitationsWithFilter(ctx, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, page, size)
}

func (s *GroupInvitationService) QueryInvitationsWithFilter(ctx context.Context, ids, groupIds, inviterIds, inviteeIds []int64, statuses []int, creationDateStart, creationDateEnd, responseDateStart, responseDateEnd, expirationDateStart, expirationDateEnd *time.Time, page, size *int) ([]*po.GroupInvitation, error) {
	var p, sz int
	if page != nil {
		p = *page
	}
	if size != nil {
		sz = *size
	}
	// Build DateRange filters from start/end time pairs
	var creationDateRange, responseDateRange, expirationDateRange *turmsmongo.DateRange
	if creationDateStart != nil || creationDateEnd != nil {
		creationDateRange = &turmsmongo.DateRange{Start: creationDateStart, End: creationDateEnd}
	}
	if responseDateStart != nil || responseDateEnd != nil {
		responseDateRange = &turmsmongo.DateRange{Start: responseDateStart, End: responseDateEnd}
	}
	if expirationDateStart != nil || expirationDateEnd != nil {
		expirationDateRange = &turmsmongo.DateRange{Start: expirationDateStart, End: expirationDateEnd}
	}
	// Convert int statuses to po.RequestStatus
	var reqStatuses []po.RequestStatus
	for _, s := range statuses {
		reqStatuses = append(reqStatuses, po.RequestStatus(s))
	}
	return s.invRepo.FindInvitations(ctx, ids, groupIds, inviterIds, inviteeIds, reqStatuses, creationDateRange, responseDateRange, expirationDateRange, &p, &sz)
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
		invs, err = s.invRepo.FindInvitations(ctx, nil, []int64{userID}, nil, nil, nil, nil, nil, nil, nil, nil)
	} else {
		invs, err = s.invRepo.FindInvitations(ctx, nil, nil, nil, []int64{userID}, nil, nil, nil, nil, nil, nil)
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
	invGroupIDs := []int64{groupID}
	invs, err := s.invRepo.FindInvitations(ctx, nil, invGroupIDs, nil, nil, nil, nil, nil, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	return &po.GroupInvitationsWithVersion{
		GroupInvitations: invs,
		LastUpdatedDate:  version,
	}, nil
}

func (s *GroupInvitationService) CountInvitations(ctx context.Context, groupID, inviterID, inviteeID *int64, status *po.RequestStatus, lastUpdatedDate *time.Time) (int64, error) {
	var groupIds, inviterIds, inviteeIds []int64
	var statuses []po.RequestStatus
	if groupID != nil {
		groupIds = []int64{*groupID}
	}
	if inviterID != nil {
		inviterIds = []int64{*inviterID}
	}
	if inviteeID != nil {
		inviteeIds = []int64{*inviteeID}
	}
	if status != nil {
		statuses = []po.RequestStatus{*status}
	}
	var creationDateRange *turmsmongo.DateRange
	if lastUpdatedDate != nil {
		creationDateRange = &turmsmongo.DateRange{Start: lastUpdatedDate}
	}
	return s.invRepo.CountInvitations(ctx, nil, groupIds, inviterIds, inviteeIds, statuses, creationDateRange, nil, nil)
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
// @MappedFrom queryGroupInvitationsByInviterId(Long inviterId)
func (s *GroupInvitationService) QueryGroupInvitationsByInviterId(ctx context.Context, inviterID int64) ([]po.GroupInvitation, error) {
	// BUG FIX: Use FindInvitationsByInviterId which filters by "irid" (inviter ID).
	// Previously was passing inviterID as groupIds parameter (wrong field).
	return s.invRepo.FindInvitationsByInviterId(ctx, inviterID)
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

// CreateGroupInvitation creates a group invitation directly without auth checks (admin API).
// Java parity: createGroupInvitation(Long id, Long groupId, Long inviterId, Long inviteeId, String content, RequestStatus status, Date creationDate, Date responseDate)
func (s *GroupInvitationService) CreateGroupInvitation(
	ctx context.Context,
	id int64,
	groupID int64,
	inviterID int64,
	inviteeID int64,
	content string,
	status po.RequestStatus,
	creationDate *time.Time,
	responseDate *time.Time,
) (*po.GroupInvitation, error) {
	now := time.Now()
	cd := now
	if creationDate != nil {
		cd = *creationDate
	}
	inv := &po.GroupInvitation{
		ID:           id,
		GroupID:      groupID,
		InviterID:    inviterID,
		InviteeID:    inviteeID,
		Content:      content,
		Status:       status,
		CreationDate: cd,
		ResponseDate: responseDate,
	}
	err := s.invRepo.Insert(ctx, inv)
	if err != nil {
		return nil, err
	}

	// Update versions
	_ = s.groupVersionService.UpdateInvitationsVersion(ctx, groupID)
	_ = s.userVersionService.UpdateSentGroupInvitationsVersion(ctx, inviterID)
	_ = s.userVersionService.UpdateReceivedGroupInvitationsVersion(ctx, inviteeID)

	return inv, nil
}

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

// GenerateInvitationId generates a new unique invitation ID using the Snowflake ID generator.
func (s *GroupInvitationService) GenerateInvitationId() int64 {
	return s.idGen.NextLargeGapId()
}

// GetEntityExpirationDate returns the entity expiration date for response wrapping
func (s *GroupInvitationService) GetEntityExpirationDate(ctx context.Context) *time.Time {
	// Returns nil for now - the expiration date is managed by the turms-properties configuration
	return nil
}
