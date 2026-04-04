package service

import (
	"context"
	"time"

	common_constant "im.turms/server/internal/domain/common/constant"
	group_constant "im.turms/server/internal/domain/group/constant"
	"im.turms/server/internal/domain/group/po"
	"im.turms/server/internal/domain/group/repository"
	user_service "im.turms/server/internal/domain/user/service"
	"im.turms/server/internal/infra/exception"
	turmsmongo "im.turms/server/internal/storage/mongo"
	"im.turms/server/pkg/protocol"
)

type GroupJoinRequestService struct {
	joinReqRepo           repository.GroupJoinRequestRepository
	groupMemberService    *GroupMemberService
	groupBlocklistService *GroupBlocklistService
	groupService          *GroupService
	groupTypeService      *GroupTypeService
	groupVersionService   *GroupVersionService
	userVersionService    *user_service.UserVersionService
}

func NewGroupJoinRequestService(
	joinReqRepo repository.GroupJoinRequestRepository,
	groupMemberService *GroupMemberService,
	groupBlocklistService *GroupBlocklistService,
	groupService *GroupService,
	groupTypeService *GroupTypeService,
	groupVersionService *GroupVersionService,
	userVersionService *user_service.UserVersionService,
) *GroupJoinRequestService {
	return &GroupJoinRequestService{
		joinReqRepo:           joinReqRepo,
		groupMemberService:    groupMemberService,
		groupBlocklistService: groupBlocklistService,
		groupService:          groupService,
		groupTypeService:      groupTypeService,
		groupVersionService:   groupVersionService,
		userVersionService:    userVersionService,
	}
}

// AuthAndCreateJoinRequest verifies permissions and creates a new join request.
func (s *GroupJoinRequestService) AuthAndCreateJoinRequest(ctx context.Context, requesterID int64, groupID int64, content string) (*po.GroupJoinRequest, error) {
	// 1. Check if group exists and is active
	typeID, err := s.groupService.QueryGroupTypeIdIfActiveAndNotDeleted(ctx, groupID)
	if err != nil {
		return nil, err
	}
	if typeID == nil {
		return nil, exception.NewTurmsError(int32(common_constant.ResponseStatusCode_ADD_USER_TO_INACTIVE_GROUP), "Group does not exist or is inactive")
	}

	// 2. Check if join request is disabled
	groupType, err := s.groupTypeService.FindByID(ctx, *typeID)
	if err != nil {
		return nil, err
	}
	if groupType == nil || groupType.JoinStrategy != group_constant.GroupJoinStrategy_JOIN_REQUEST {
		return nil, exception.NewTurmsError(int32(common_constant.ResponseStatusCode_GROUP_JOIN_REQUEST_IS_DISABLED), "Group join request is disabled")
	}

	// 3. Check if requester is blocked
	isBlocked, err := s.groupBlocklistService.IsBlocked(ctx, groupID, requesterID)
	if err != nil {
		return nil, err
	}
	if isBlocked {
		return nil, exception.NewTurmsError(int32(common_constant.ResponseStatusCode_BLOCKED_USER_SEND_GROUP_JOIN_REQUEST), "User is blocked by group")
	}

	// 4. Check if requester is already a member
	isMember, err := s.groupMemberService.IsGroupMember(ctx, groupID, requesterID)
	if err != nil {
		return nil, err
	}
	if isMember {
		return nil, exception.NewTurmsError(int32(common_constant.ResponseStatusCode_USER_ALREADY_GROUP_MEMBER), "User is already a member of the group")
	}

	// 5. Create request
	now := time.Now()
	id := now.UnixNano() // Simplified ID generation

	req := &po.GroupJoinRequest{
		ID:           id,
		GroupID:      groupID,
		RequesterID:  requesterID,
		Content:      content,
		Status:       po.RequestStatusPending,
		CreationDate: now,
	}
	err = s.joinReqRepo.Insert(ctx, req)
	if err != nil {
		return nil, err
	}

	// 6. Update version
	_ = s.groupVersionService.UpdateJoinRequestsVersion(ctx, groupID)

	return req, nil
}

// AuthAndRecallPendingJoinRequest recalls a pending join request.
func (s *GroupJoinRequestService) AuthAndRecallPendingJoinRequest(ctx context.Context, requesterID int64, requestID int64) error {
	req, err := s.joinReqRepo.FindByID(ctx, requestID)
	if err != nil {
		return err
	}
	if req == nil {
		return exception.NewTurmsError(int32(common_constant.ResponseStatusCode_RECALL_NON_PENDING_GROUP_JOIN_REQUEST), "Join request not found")
	}
	if req.Status != po.RequestStatusPending {
		return exception.NewTurmsError(int32(common_constant.ResponseStatusCode_RECALL_NON_PENDING_GROUP_JOIN_REQUEST), "Join request is not pending")
	}
	if req.RequesterID != requesterID {
		return exception.NewTurmsError(int32(common_constant.ResponseStatusCode_NOT_SENDER_TO_RECALL_GROUP_JOIN_REQUEST), "Unauthorized to recall join request")
	}

	updated, err := s.joinReqRepo.UpdateStatusIfPending(ctx, requestID, requesterID, po.RequestStatusCanceled, nil, time.Now())
	if err != nil {
		return err
	}
	if !updated {
		return exception.NewTurmsError(int32(common_constant.ResponseStatusCode_RECALL_NON_PENDING_GROUP_JOIN_REQUEST), "Failed to recall join request")
	}

	_ = s.groupVersionService.UpdateJoinRequestsVersion(ctx, req.GroupID)
	return nil
}

// AuthAndHandleJoinRequest handles a joint request (accept or decline).
// @MappedFrom authAndHandleJoinRequest(@NotNull Long requesterId, @NotNull Long joinRequestId, @NotNull @ValidResponseAction ResponseAction action, @Nullable String responseReason)
func (s *GroupJoinRequestService) AuthAndHandleJoinRequest(ctx context.Context, responderID int64, requestID int64, status po.RequestStatus, reason string) error {
	req, err := s.joinReqRepo.FindByID(ctx, requestID)
	if err != nil {
		return err
	}
	if req == nil {
		return exception.NewTurmsError(int32(common_constant.ResponseStatusCode_UPDATE_NON_PENDING_GROUP_JOIN_REQUEST), "Join request not found")
	}
	if req.Status != po.RequestStatusPending {
		return exception.NewTurmsError(int32(common_constant.ResponseStatusCode_UPDATE_NON_PENDING_GROUP_JOIN_REQUEST), "Join request is not pending")
	}

	// RBAC: Check if responder is Owner or Manager
	role, err := s.groupMemberService.QueryGroupMemberRole(ctx, req.GroupID, responderID)
	if err != nil {
		return err
	}
	if role == nil || (*role != protocol.GroupMemberRole_OWNER && *role != protocol.GroupMemberRole_MANAGER) {
		return exception.NewTurmsError(int32(common_constant.ResponseStatusCode_NOT_GROUP_OWNER_OR_MANAGER_TO_UPDATE_GROUP_JOIN_REQUEST), "Unauthorized to handle join request")
	}

	updated, err := s.joinReqRepo.UpdateStatusIfPending(ctx, requestID, responderID, status, &reason, time.Now())
	if err != nil {
		return err
	}
	if !updated {
		return exception.NewTurmsError(int32(common_constant.ResponseStatusCode_UPDATE_NON_PENDING_GROUP_JOIN_REQUEST), "Failed to handle join request")
	}

	if status == po.RequestStatusAccepted {
		// Add user to group
		err = s.groupMemberService.AddGroupMember(ctx, req.GroupID, req.RequesterID, protocol.GroupMemberRole_MEMBER, nil, nil)
		if err != nil {
			return err
		}
	}

	_ = s.groupVersionService.UpdateJoinRequestsVersion(ctx, req.GroupID)
	return nil
}

// @MappedFrom queryJoinRequests(@Nullable Set<Long> ids, @Nullable Set<Long> groupIds, @Nullable Set<Long> requesterIds, @Nullable Set<Long> responderIds, @Nullable Set<RequestStatus> statuses, @Nullable DateRange creationDateRange, @Nullable DateRange responseDateRange, @Nullable DateRange expirationDateRange, @Nullable Integer page, @Nullable Integer size)
func (s *GroupJoinRequestService) QueryJoinRequests(ctx context.Context, groupID *int64, requesterID *int64, responderID *int64, status *po.RequestStatus, creationDate *time.Time, page int, size int) ([]*po.GroupJoinRequest, error) {
	return s.joinReqRepo.FindRequests(ctx, groupID, requesterID, responderID, status, creationDate, nil, nil, page, size)
}

// AuthAndQueryGroupJoinRequestsWithVersion queries the group join requests requested.
func (s *GroupJoinRequestService) AuthAndQueryGroupJoinRequestsWithVersion(ctx context.Context, requesterID int64, groupID int64, lastUpdatedDate *time.Time) (*po.GroupJoinRequestsWithVersion, error) {
	role, err := s.groupMemberService.QueryGroupMemberRole(ctx, groupID, requesterID)
	if err != nil {
		return nil, err
	}
	if role == nil || (*role != protocol.GroupMemberRole_OWNER && *role != protocol.GroupMemberRole_MANAGER) {
		return nil, exception.NewTurmsError(int32(common_constant.ResponseStatusCode_NOT_GROUP_OWNER_OR_MANAGER_TO_QUERY_GROUP_JOIN_REQUEST), "No permission to query group join requests")
	}

	version, err := s.groupVersionService.QueryGroupJoinRequestsVersion(ctx, groupID)
	if err != nil {
		return nil, err
	}
	if lastUpdatedDate != nil && version != nil && !version.After(*lastUpdatedDate) {
		return &po.GroupJoinRequestsWithVersion{LastUpdatedDate: version}, nil
	}

	reqs, err := s.joinReqRepo.FindRequests(ctx, &groupID, nil, nil, nil, nil, nil, nil, 0, 1000)
	if err != nil {
		return nil, err
	}

	return &po.GroupJoinRequestsWithVersion{
		GroupJoinRequests: reqs,
		LastUpdatedDate:   version,
	}, nil
}

// QueryUserGroupJoinRequestsWithVersion queries the group join requests for a user.
func (s *GroupJoinRequestService) QueryUserGroupJoinRequestsWithVersion(ctx context.Context, requesterID int64, lastUpdatedDate *time.Time) (*po.GroupJoinRequestsWithVersion, error) {
	version, err := s.userVersionService.QueryGroupJoinRequestsVersion(ctx, requesterID)
	if err != nil {
		return nil, err
	}
	if lastUpdatedDate != nil && version != nil && !version.After(*lastUpdatedDate) {
		return &po.GroupJoinRequestsWithVersion{LastUpdatedDate: version}, nil
	}

	reqs, err := s.joinReqRepo.FindRequests(ctx, nil, &requesterID, nil, nil, nil, nil, nil, 0, 1000)
	if err != nil {
		return nil, err
	}

	return &po.GroupJoinRequestsWithVersion{
		GroupJoinRequests: reqs,
		LastUpdatedDate:   version,
	}, nil
}

func (s *GroupJoinRequestService) QueryGroupJoinRequestsByGroupId(ctx context.Context, groupID int64) ([]po.GroupJoinRequest, error) {
	return s.joinReqRepo.FindRequestsByGroupID(ctx, groupID)
}

func (s *GroupJoinRequestService) QueryGroupJoinRequestsByRequesterId(ctx context.Context, requesterID int64) ([]po.GroupJoinRequest, error) {
	return s.joinReqRepo.FindRequestsByRequesterID(ctx, requesterID)
}

func (s *GroupJoinRequestService) QueryGroupId(ctx context.Context, requestID int64) (*int64, error) {
	return s.joinReqRepo.FindGroupId(ctx, requestID)
}

func (s *GroupJoinRequestService) CountJoinRequests(ctx context.Context, ids, groupIds, requesterIds, responderIds []int64, statuses []po.RequestStatus, creationDateRange, responseDateRange, expirationDateRange *turmsmongo.DateRange) (int64, error) {
	return s.joinReqRepo.CountRequests(ctx, ids, groupIds, requesterIds, responderIds, statuses, creationDateRange, responseDateRange, expirationDateRange)
}

func (s *GroupJoinRequestService) DeleteJoinRequests(ctx context.Context, ids []int64) (int64, error) {
	return s.joinReqRepo.DeleteRequests(ctx, ids)
}

func (s *GroupJoinRequestService) UpdatePendingJoinRequestStatus(ctx context.Context, requestID int64, responderID int64, requestStatus po.RequestStatus, responseReason *string) (bool, error) {
	return s.joinReqRepo.UpdateStatusIfPending(ctx, requestID, responderID, requestStatus, responseReason, time.Now())
}

func (s *GroupJoinRequestService) UpdateJoinRequests(ctx context.Context, requestIds []int64, requesterId, responderId *int64, content *string, status *po.RequestStatus, creationDate, responseDate *time.Time) error {
	return s.joinReqRepo.UpdateRequests(ctx, requestIds, requesterId, responderId, content, status, creationDate, responseDate)
}

// Backward Compatibility Aliases

func (s *GroupJoinRequestService) CreateJoinRequest(ctx context.Context, groupID int64, requesterID int64, content string) (*po.GroupJoinRequest, error) {
	return s.AuthAndCreateJoinRequest(ctx, requesterID, groupID, content)
}

func (s *GroupJoinRequestService) RecallPendingJoinRequest(ctx context.Context, requestID int64, requesterID int64) (bool, error) {
	err := s.AuthAndRecallPendingJoinRequest(ctx, requesterID, requestID)
	return err == nil, err
}

func (s *GroupJoinRequestService) ReplyToJoinRequest(ctx context.Context, requestID int64, responderID int64, accept bool) (bool, error) {
	status := po.RequestStatusDeclined
	if accept {
		status = po.RequestStatusAccepted
	}
	err := s.AuthAndHandleJoinRequest(ctx, responderID, requestID, status, "")
	return err == nil, err
}
