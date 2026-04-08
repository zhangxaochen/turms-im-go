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
// Java order: isGroupMember → isBlocked → queryGroupTypeIdIfActiveAndNotDeleted → queryGroupType → create
func (s *GroupJoinRequestService) AuthAndCreateJoinRequest(ctx context.Context, requesterID int64, groupID int64, content string) (*po.GroupJoinRequest, error) {
	// 1. Check if requester is already a member
	isMember, err := s.groupMemberService.IsGroupMember(ctx, groupID, requesterID)
	if err != nil {
		return nil, err
	}
	if isMember {
		return nil, exception.NewTurmsError(int32(common_constant.ResponseStatusCode_USER_ALREADY_GROUP_MEMBER), "User is already a member of the group")
	}

	// 2. Check if requester is blocked
	isBlocked, err := s.groupBlocklistService.IsBlocked(ctx, groupID, requesterID)
	if err != nil {
		return nil, err
	}
	if isBlocked {
		return nil, exception.NewTurmsError(int32(common_constant.ResponseStatusCode_BLOCKED_USER_SEND_GROUP_JOIN_REQUEST), "User is blocked by group")
	}

	// 3. Check if group exists and is active
	typeID, err := s.groupService.QueryGroupTypeIdIfActiveAndNotDeleted(ctx, groupID)
	if err != nil {
		return nil, err
	}
	if typeID == nil {
		return nil, exception.NewTurmsError(int32(common_constant.ResponseStatusCode_ADD_USER_TO_INACTIVE_GROUP), "Group does not exist or is inactive")
	}

	// 4. Check join strategy and return specific error codes
	groupType, err := s.groupTypeService.FindByID(ctx, *typeID)
	if err != nil {
		return nil, err
	}
	if groupType == nil {
		return nil, exception.NewTurmsError(int32(common_constant.ResponseStatusCode_GROUP_JOIN_REQUEST_IS_DISABLED), "Group type not found")
	}

	switch groupType.JoinStrategy {
	case group_constant.GroupJoinStrategy_JOIN_REQUEST:
		// OK - proceed
	case group_constant.GroupJoinStrategy_MEMBERSHIP_REQUEST:
		return nil, exception.NewTurmsError(int32(common_constant.ResponseStatusCode_SEND_GROUP_JOIN_REQUEST_TO_GROUP_USING_MEMBERSHIP_REQUEST),
			"Please use membership request to join this group")
	case group_constant.GroupJoinStrategy_INVITATION:
		return nil, exception.NewTurmsError(int32(common_constant.ResponseStatusCode_SEND_GROUP_JOIN_REQUEST_TO_GROUP_USING_INVITATION),
			"Please use invitation to join this group")
	case group_constant.GroupJoinStrategy_QUESTION:
		return nil, exception.NewTurmsError(int32(common_constant.ResponseStatusCode_SEND_GROUP_JOIN_REQUEST_TO_GROUP_USING_QUESTION),
			"Please answer group question to join this group")
	default:
		return nil, exception.NewTurmsError(int32(common_constant.ResponseStatusCode_GROUP_JOIN_REQUEST_IS_DISABLED), "Group join request is disabled")
	}

	// 5. Create request - default content to empty string if nil
	finalContent := content
	if finalContent == "" {
		finalContent = ""
	}

	now := time.Now()
	id := now.UnixNano() // Simplified ID generation

	req := &po.GroupJoinRequest{
		ID:           id,
		GroupID:      groupID,
		RequesterID:  requesterID,
		Content:      finalContent,
		Status:       po.RequestStatusPending,
		CreationDate: now,
	}
	err = s.joinReqRepo.Insert(ctx, req)
	if err != nil {
		return nil, err
	}

	// 6. Update both group version and user version
	_ = s.groupVersionService.UpdateJoinRequestsVersion(ctx, groupID)
	_ = s.userVersionService.UpdateSentGroupJoinRequestsVersion(ctx, requesterID)

	return req, nil
}

// AuthAndRecallPendingJoinRequest recalls a pending join request.
// Java order: check allowRecallPendingJoinRequestBySender property → query requesterId/status/groupId →
// check authorization (requesterId match) → check status (Pending) → update → update versions
func (s *GroupJoinRequestService) AuthAndRecallPendingJoinRequest(ctx context.Context, requesterID int64, requestID int64) error {
	// 1. Use projected query (Java: queryRequesterIdAndStatusAndGroupId)
	reqRequesterID, status, groupID, err := s.joinReqRepo.FindRequesterIdAndStatusAndGroupId(ctx, requestID)
	if err != nil {
		return err
	}
	if reqRequesterID == nil {
		return exception.NewTurmsError(int32(common_constant.ResponseStatusCode_RECALL_NON_PENDING_GROUP_JOIN_REQUEST), "Join request not found")
	}

	// 2. Check authorization FIRST (before status, to avoid leaking status info)
	if *reqRequesterID != requesterID {
		return exception.NewTurmsError(int32(common_constant.ResponseStatusCode_NOT_SENDER_TO_RECALL_GROUP_JOIN_REQUEST), "Unauthorized to recall join request")
	}

	// 3. Check status is Pending
	if *status != po.RequestStatusPending {
		return exception.NewTurmsError(int32(common_constant.ResponseStatusCode_RECALL_NON_PENDING_GROUP_JOIN_REQUEST), "Join request is not pending")
	}

	// 4. Update status
	updated, err := s.joinReqRepo.UpdateStatusIfPending(ctx, requestID, requesterID, po.RequestStatusCanceled, nil, time.Now())
	if err != nil {
		return err
	}
	if !updated {
		return exception.NewTurmsError(int32(common_constant.ResponseStatusCode_RECALL_NON_PENDING_GROUP_JOIN_REQUEST), "Failed to recall join request")
	}

	// 5. Update both group version and user version
	_ = s.groupVersionService.UpdateJoinRequestsVersion(ctx, *groupID)
	_ = s.userVersionService.UpdateSentGroupJoinRequestsVersion(ctx, requesterID)

	return nil
}

// AuthAndHandleJoinRequest handles a join request (accept, decline, or ignore).
// Java order: findById → check authorization (isOwnerOrManager) → check status → check expiration → handle action → update versions
func (s *GroupJoinRequestService) AuthAndHandleJoinRequest(ctx context.Context, responderID int64, requestID int64, status po.RequestStatus, reason string) error {
	req, err := s.joinReqRepo.FindByID(ctx, requestID)
	if err != nil {
		return err
	}
	if req == nil {
		return exception.NewTurmsError(int32(common_constant.ResponseStatusCode_NOT_GROUP_OWNER_OR_MANAGER_TO_UPDATE_GROUP_JOIN_REQUEST), "Join request not found")
	}

	// 1. Check authorization FIRST (before status check, to avoid leaking status info)
	role, err := s.groupMemberService.QueryGroupMemberRole(ctx, req.GroupID, responderID)
	if err != nil {
		return err
	}
	if role == nil || (*role != protocol.GroupMemberRole_OWNER && *role != protocol.GroupMemberRole_MANAGER) {
		return exception.NewTurmsError(int32(common_constant.ResponseStatusCode_NOT_GROUP_OWNER_OR_MANAGER_TO_UPDATE_GROUP_JOIN_REQUEST), "Unauthorized to handle join request")
	}

	// 2. Check status is Pending
	if req.Status != po.RequestStatusPending {
		return exception.NewTurmsError(int32(common_constant.ResponseStatusCode_UPDATE_NON_PENDING_GROUP_JOIN_REQUEST), "Join request is not pending")
	}

	// 3. Check expiration for PENDING requests
	expireAfterSeconds := s.joinReqRepo.GetEntityExpireAfterSeconds()
	if expireAfterSeconds > 0 {
		expirationTime := req.CreationDate.Add(time.Duration(expireAfterSeconds) * time.Second)
		if time.Now().After(expirationTime) {
			return exception.NewTurmsError(int32(common_constant.ResponseStatusCode_EXPIRED_GROUP_JOIN_REQUEST), "Join request has expired")
		}
	}

	// 4. Handle IGNORE action - just update status, no member addition
	if status == po.RequestStatusIgnored {
		updated, err := s.joinReqRepo.UpdateStatusIfPending(ctx, requestID, responderID, po.RequestStatusIgnored, &reason, time.Now())
		if err != nil {
			return err
		}
		if !updated {
			return exception.NewTurmsError(int32(common_constant.ResponseStatusCode_UPDATE_NON_PENDING_GROUP_JOIN_REQUEST), "Failed to handle join request")
		}
		_ = s.groupVersionService.UpdateJoinRequestsVersion(ctx, req.GroupID)
		_ = s.userVersionService.UpdateSentGroupJoinRequestsVersion(ctx, req.RequesterID)
		return nil
	}

	// 5. Update status for ACCEPT/DECLINE
	updated, err := s.joinReqRepo.UpdateStatusIfPending(ctx, requestID, responderID, status, &reason, time.Now())
	if err != nil {
		return err
	}
	if !updated {
		return exception.NewTurmsError(int32(common_constant.ResponseStatusCode_UPDATE_NON_PENDING_GROUP_JOIN_REQUEST), "Failed to handle join request")
	}

	// 6. If ACCEPT, add user to group (handle DuplicateKey gracefully)
	if status == po.RequestStatusAccepted {
		err = s.groupMemberService.AddGroupMember(ctx, req.GroupID, req.RequesterID, protocol.GroupMemberRole_MEMBER, nil, nil)
		if err != nil {
			if !exception.IsDuplicateKeyError(err) {
				return err
			}
			// DuplicateKey: user is already a member, treat as success
		}
	}

	// 7. Update both group version and user version
	_ = s.groupVersionService.UpdateJoinRequestsVersion(ctx, req.GroupID)
	_ = s.userVersionService.UpdateSentGroupJoinRequestsVersion(ctx, req.RequesterID)

	return nil
}

// @MappedFrom queryJoinRequests(@Nullable Set<Long> ids, @Nullable Set<Long> groupIds, @Nullable Set<Long> requesterIds, @Nullable Set<Long> responderIds, @Nullable Set<RequestStatus> statuses, @Nullable DateRange creationDateRange, @Nullable DateRange responseDateRange, @Nullable DateRange expirationDateRange, @Nullable Integer page, @Nullable Integer size)
func (s *GroupJoinRequestService) QueryJoinRequests(ctx context.Context, groupID *int64, requesterID *int64, responderID *int64, status *po.RequestStatus, creationDate *time.Time, page int, size int) ([]*po.GroupJoinRequest, error) {
	var groupIds, requesterIds, responderIds []int64
	var statuses []po.RequestStatus
	if groupID != nil {
		groupIds = []int64{*groupID}
	}
	if requesterID != nil {
		requesterIds = []int64{*requesterID}
	}
	if responderID != nil {
		responderIds = []int64{*responderID}
	}
	if status != nil {
		statuses = []po.RequestStatus{*status}
	}
	var creationDateRange *turmsmongo.DateRange
	if creationDate != nil {
		creationDateRange = &turmsmongo.DateRange{Start: creationDate}
	}
	p, sz := page, size
	return s.joinReqRepo.FindRequests(ctx, nil, groupIds, requesterIds, responderIds, statuses, creationDateRange, nil, nil, &p, &sz)
}

func (s *GroupJoinRequestService) QueryJoinRequestsWithPagination(ctx context.Context, page, size *int) ([]*po.GroupJoinRequest, error) {
	return s.QueryJoinRequestsWithFilter(ctx, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, page, size)
}

func (s *GroupJoinRequestService) QueryJoinRequestsWithFilter(ctx context.Context, ids, groupIds, requesterIds, responderIds []int64, statuses []int, creationDateStart, creationDateEnd, responseDateStart, responseDateEnd, expirationDateStart, expirationDateEnd *time.Time, page, size *int) ([]*po.GroupJoinRequest, error) {
	var p, sz int
	if page != nil {
		p = *page
	}
	if size != nil {
		sz = *size
	}
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
	var reqStatuses []po.RequestStatus
	for _, st := range statuses {
		reqStatuses = append(reqStatuses, po.RequestStatus(st))
	}
	return s.joinReqRepo.FindRequests(ctx, ids, groupIds, requesterIds, responderIds, reqStatuses, creationDateRange, responseDateRange, expirationDateRange, &p, &sz)
}

// AuthAndQueryGroupJoinRequestsWithVersion queries the group join requests requested.
// Java has two branches: groupId != null → query by group (owner/manager check); groupId == null → query by requesterId via userVersion
func (s *GroupJoinRequestService) AuthAndQueryGroupJoinRequestsWithVersion(ctx context.Context, requesterID int64, groupID *int64, lastUpdatedDate *time.Time) (*po.GroupJoinRequestsWithVersion, error) {
	if groupID != nil {
		// Branch 1: Query by groupId (requires owner/manager role)
		role, err := s.groupMemberService.QueryGroupMemberRole(ctx, *groupID, requesterID)
		if err != nil {
			return nil, err
		}
		if role == nil || (*role != protocol.GroupMemberRole_OWNER && *role != protocol.GroupMemberRole_MANAGER) {
			return nil, exception.NewTurmsError(int32(common_constant.ResponseStatusCode_NOT_GROUP_OWNER_OR_MANAGER_TO_QUERY_GROUP_JOIN_REQUEST), "No permission to query group join requests")
		}

		version, err := s.groupVersionService.QueryGroupJoinRequestsVersion(ctx, *groupID)
		if err != nil {
			return nil, err
		}
		if lastUpdatedDate != nil && version != nil && !version.After(*lastUpdatedDate) {
			return &po.GroupJoinRequestsWithVersion{LastUpdatedDate: version}, nil
		}

		reqGroupIDs := []int64{*groupID}
		reqs, err := s.joinReqRepo.FindRequests(ctx, nil, reqGroupIDs, nil, nil, nil, nil, nil, nil, nil, nil)
		if err != nil {
			return nil, err
		}
		if len(reqs) == 0 {
			return nil, exception.NewTurmsError(int32(common_constant.ResponseStatusCode_NO_CONTENT), "No join requests found")
		}

		return &po.GroupJoinRequestsWithVersion{
			GroupJoinRequests: reqs,
			LastUpdatedDate:   version,
		}, nil
	}

	// Branch 2: Query by requesterId (user's own sent requests)
	version, err := s.userVersionService.QueryGroupJoinRequestsVersion(ctx, requesterID)
	if err != nil {
		return nil, err
	}
	if lastUpdatedDate != nil && version != nil && !version.After(*lastUpdatedDate) {
		return &po.GroupJoinRequestsWithVersion{LastUpdatedDate: version}, nil
	}

	reqRequesterIDs := []int64{requesterID}
	reqs, err := s.joinReqRepo.FindRequests(ctx, nil, nil, reqRequesterIDs, nil, nil, nil, nil, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	if len(reqs) == 0 {
		return nil, exception.NewTurmsError(int32(common_constant.ResponseStatusCode_NO_CONTENT), "No join requests found")
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

	reqRequesterIDs := []int64{requesterID}
	reqs, err := s.joinReqRepo.FindRequests(ctx, nil, nil, reqRequesterIDs, nil, nil, nil, nil, nil, nil, nil)
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
	if len(ids) == 0 {
		return 0, nil
	}
	return s.joinReqRepo.DeleteRequests(ctx, ids)
}

// UpdatePendingJoinRequestStatus updates the status of a pending join request.
// Java validates: groupId not null, joinRequestId not null, requestStatus not null, requestStatus != PENDING,
// responderId not null, responseReason maxLength. Updates groupVersion on success.
func (s *GroupJoinRequestService) UpdatePendingJoinRequestStatus(ctx context.Context, groupID int64, requestID int64, requestStatus po.RequestStatus, responderID int64, responseReason *string) (bool, error) {
	// Validations
	if requestStatus == po.RequestStatusPending {
		return false, exception.NewTurmsError(int32(common_constant.ResponseStatusCode_ILLEGAL_ARGUMENT), "Request status cannot be PENDING")
	}

	modified, err := s.joinReqRepo.UpdateStatusIfPending(ctx, requestID, responderID, requestStatus, responseReason, time.Now())
	if err != nil {
		return false, err
	}

	// Update group version on successful modification
	if modified {
		_ = s.groupVersionService.UpdateJoinRequestsVersion(ctx, groupID)
	}

	return modified, nil
}

// UpdateJoinRequests updates multiple join requests.
// Java validates: requestIds not empty, content maxLength, validRequestStatus, pastOrPresent dates.
// Returns early if all update fields are null.
func (s *GroupJoinRequestService) UpdateJoinRequests(ctx context.Context, requestIds []int64, requesterId, responderId *int64, content *string, status *po.RequestStatus, creationDate, responseDate *time.Time) error {
	// Validation: requestIds must not be empty
	if len(requestIds) == 0 {
		return exception.NewTurmsError(int32(common_constant.ResponseStatusCode_ILLEGAL_ARGUMENT), "Request IDs must not be empty")
	}

	// Early return if all update fields are null (Java: areAllNull check)
	if requesterId == nil && responderId == nil && content == nil && status == nil && creationDate == nil && responseDate == nil {
		return nil
	}

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

// CreateGroupJoinRequest creates a join request directly (admin-level, no auth checks).
// This is the admin version that creates a request with all fields, defaults null values,
// and updates both group and user versions.
func (s *GroupJoinRequestService) CreateGroupJoinRequest(ctx context.Context, id *int64, groupID int64, requesterID int64, responderID *int64, content *string, status *po.RequestStatus, creationDate, responseDate *time.Time, responseReason *string) (*po.GroupJoinRequest, error) {
	// Default ID if nil
	reqID := int64(0)
	if id != nil {
		reqID = *id
	} else {
		reqID = time.Now().UnixNano()
	}

	// Default content to empty string if nil
	finalContent := ""
	if content != nil {
		finalContent = *content
	}

	// Default creationDate to now if nil
	now := time.Now()
	finalCreationDate := now
	if creationDate != nil {
		finalCreationDate = *creationDate
	}

	// Default status to PENDING if nil
	finalStatus := po.RequestStatusPending
	if status != nil {
		finalStatus = *status
	}

	// Compute responseDate based on status for new record
	finalResponseDate := responseDate
	if finalStatus != po.RequestStatusPending {
		if finalResponseDate == nil {
			rd := now
			finalResponseDate = &rd
		}
	}

	req := &po.GroupJoinRequest{
		ID:           reqID,
		GroupID:      groupID,
		RequesterID:  requesterID,
		ResponderID:  responderID,
		Content:      finalContent,
		Status:       finalStatus,
		CreationDate: finalCreationDate,
		ResponseDate: finalResponseDate,
		Reason:       responseReason,
	}

	err := s.joinReqRepo.Insert(ctx, req)
	if err != nil {
		return nil, err
	}

	// Update both group and user versions
	_ = s.groupVersionService.UpdateJoinRequestsVersion(ctx, groupID)
	_ = s.userVersionService.UpdateSentGroupJoinRequestsVersion(ctx, requesterID)

	return req, nil
}

// GetEntityExpirationDate returns the entity expiration date for response wrapping
func (s *GroupJoinRequestService) GetEntityExpirationDate(ctx context.Context) *time.Time {
	return nil
}
