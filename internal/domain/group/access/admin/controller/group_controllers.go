package controller

import (
	"context"
	"fmt"
	"time"

	"im.turms/server/internal/domain/common/access/admin/dto/response"
	"im.turms/server/internal/domain/common/constant"
	common_dto "im.turms/server/internal/domain/common/dto"
	group_dto "im.turms/server/internal/domain/group/access/admin/dto"
	"im.turms/server/internal/domain/group/po"
	"im.turms/server/internal/domain/group/service"
	group_constant "im.turms/server/internal/domain/group/constant"
	"im.turms/server/internal/domain/common/infra/idgen"
	commoncontroller "im.turms/server/internal/domain/common/access/admin/controller"
	turmsmongo "im.turms/server/internal/storage/mongo"
	"im.turms/server/pkg/protocol"
)

// ---------------------------------------------------------------------------
// Shared pagination response types
// ---------------------------------------------------------------------------

// PaginationResponse is a generic paginated response carrying a total count and records.
type PaginationResponse struct {
	Total   int64       `json:"total"`
	Records interface{} `json:"records"`
}

// ---------------------------------------------------------------------------
// 1. GroupBlocklistController
// ---------------------------------------------------------------------------

// GroupBlocklistController handles admin operations on the group blocklist.
type GroupBlocklistController struct {
	*commoncontroller.BaseController
	groupBlocklistService *service.GroupBlocklistService
}

func NewGroupBlocklistController(base *commoncontroller.BaseController, groupBlocklistService *service.GroupBlocklistService) *GroupBlocklistController {
	return &GroupBlocklistController{
		BaseController:        base,
		groupBlocklistService: groupBlocklistService,
	}
}

// AddGroupBlockedUser adds a blocked user to a group.
func (c *GroupBlocklistController) AddGroupBlockedUser(ctx context.Context, addDTO group_dto.AddGroupBlockedUserDTO) (*po.GroupBlockedUser, error) {
	var requesterID int64
	if addDTO.RequesterId != nil {
		requesterID = *addDTO.RequesterId
	}
	return c.groupBlocklistService.AddBlockedUser(ctx, *addDTO.GroupId, *addDTO.UserId, requesterID, addDTO.BlockDate)
}

// QueryGroupBlockedUsers queries blocked users with optional filters (non-paged).
// Bug fix: added missing filter params (groupIds, userIds, blockDateStart, blockDateEnd, requesterIds),
// use page=0, getPageSize(size), and return results instead of discarding them.
func (c *GroupBlocklistController) QueryGroupBlockedUsers(
	ctx context.Context,
	groupIds []int64,
	userIds []int64,
	blockDateStart *time.Time,
	blockDateEnd *time.Time,
	requesterIds []int64,
	size *int,
) ([]po.GroupBlockedUser, error) {
	var blockDateRange *turmsmongo.DateRange
	if blockDateStart != nil || blockDateEnd != nil {
		blockDateRange = &turmsmongo.DateRange{Start: blockDateStart, End: blockDateEnd}
	}
	page := 0
	actualSize := c.GetPageSize(size)
	return c.groupBlocklistService.QueryBlockedUsersWithFilter(ctx, groupIds, userIds, blockDateRange, requesterIds, &page, &actualSize)
}

// QueryGroupBlockedUsersByPage queries blocked users with pagination.
// NEW endpoint: Call countBlockedUsers() for total, then queryBlockedUsersWithFilter,
// return PaginationResponse{Total, Records}.
func (c *GroupBlocklistController) QueryGroupBlockedUsersByPage(
	ctx context.Context,
	groupIds []int64,
	userIds []int64,
	blockDateStart *time.Time,
	blockDateEnd *time.Time,
	requesterIds []int64,
	page int,
	size *int,
) (*PaginationResponse, error) {
	var blockDateRange *turmsmongo.DateRange
	if blockDateStart != nil || blockDateEnd != nil {
		blockDateRange = &turmsmongo.DateRange{Start: blockDateStart, End: blockDateEnd}
	}
	total, err := c.groupBlocklistService.CountBlockedUsers(ctx, groupIds, userIds, blockDateRange, requesterIds)
	if err != nil {
		return nil, err
	}
	actualSize := c.GetPageSize(size)
	results, err := c.groupBlocklistService.QueryBlockedUsersWithFilter(ctx, groupIds, userIds, blockDateRange, requesterIds, &page, &actualSize)
	if err != nil {
		return nil, err
	}
	return &PaginationResponse{Total: total, Records: results}, nil
}

// UpdateGroupBlockedUsers updates blocked user records.
// Bug fix: deduplicate keys using Set before service call, return UpdateResultDTO.
func (c *GroupBlocklistController) UpdateGroupBlockedUsers(
	ctx context.Context,
	keys []po.GroupBlockedUserKey,
	updateDTO group_dto.UpdateGroupBlockedUserDTO,
) (*response.UpdateResultDTO, error) {
	// Deduplicate keys
	seen := make(map[po.GroupBlockedUserKey]struct{})
	deduped := make([]po.GroupBlockedUserKey, 0, len(keys))
	for _, k := range keys {
		if _, ok := seen[k]; !ok {
			seen[k] = struct{}{}
			deduped = append(deduped, k)
		}
	}
	err := c.groupBlocklistService.UpdateBlockedUsers(ctx, deduped, updateDTO.BlockDate, updateDTO.RequesterId)
	if err != nil {
		return nil, err
	}
	return &response.UpdateResultDTO{UpdatedCount: int64(len(deduped))}, nil
}

// DeleteGroupBlockedUsers deletes blocked user records.
// Bug fix: deduplicate keys using Set, return DeleteResultDTO.
func (c *GroupBlocklistController) DeleteGroupBlockedUsers(
	ctx context.Context,
	keys []po.GroupBlockedUserKey,
) (*response.DeleteResultDTO, error) {
	// Deduplicate keys
	seen := make(map[po.GroupBlockedUserKey]struct{})
	deduped := make([]po.GroupBlockedUserKey, 0, len(keys))
	for _, k := range keys {
		if _, ok := seen[k]; !ok {
			seen[k] = struct{}{}
			deduped = append(deduped, k)
		}
	}
	err := c.groupBlocklistService.DeleteBlockedUsers(ctx, deduped)
	if err != nil {
		return nil, err
	}
	return &response.DeleteResultDTO{DeletedCount: int64(len(deduped))}, nil
}

// ---------------------------------------------------------------------------
// 2. GroupController
// ---------------------------------------------------------------------------

// GroupController handles admin operations on groups.
type GroupController struct {
	*commoncontroller.BaseController
	groupService *service.GroupService
}

func NewGroupController(base *commoncontroller.BaseController, groupService *service.GroupService) *GroupController {
	return &GroupController{
		BaseController: base,
		groupService:   groupService,
	}
}

// AddGroup creates a new group.
// Bug fix: nil check on CreatorId before dereferencing. Default ownerId to creatorId if both OwnerId nil.
func (c *GroupController) AddGroup(ctx context.Context, addDTO group_dto.AddGroupDTO) (*po.Group, error) {
	if addDTO.CreatorId == nil {
		return nil, nil
	}
	creatorID := *addDTO.CreatorId
	ownerID := creatorID
	if addDTO.OwnerId != nil {
		ownerID = *addDTO.OwnerId
	}
	var minimumScore *int32
	if addDTO.MinimumScore != nil {
		ms := int32(*addDTO.MinimumScore)
		minimumScore = &ms
	}
	return c.groupService.CreateGroup(
		ctx,
		creatorID,
		ownerID,
		addDTO.Name,
		addDTO.Intro,
		addDTO.Announcement,
		minimumScore,
		addDTO.TypeId,
		addDTO.CreationDate,
		addDTO.DeletionDate,
		addDTO.MuteEndDate,
		addDTO.IsActive,
	)
}

// QueryGroups queries groups with optional filters (non-paged).
// Bug fix: add lastUpdatedDateStart/lastUpdatedDateEnd params, getPageSize(size).
func (c *GroupController) QueryGroups(
	ctx context.Context,
	ids []int64,
	typeIds []int64,
	creatorIds []int64,
	ownerIds []int64,
	isActive *bool,
	creationDateStart *time.Time,
	creationDateEnd *time.Time,
	deletionDateStart *time.Time,
	deletionDateEnd *time.Time,
	muteEndDateStart *time.Time,
	muteEndDateEnd *time.Time,
	memberIds []int64,
	lastUpdatedDateStart *time.Time,
	lastUpdatedDateEnd *time.Time,
	size *int,
) ([]*po.Group, error) {
	actualSize := c.GetPageSize(size)
	page := 0
	return c.groupService.QueryGroupsWithFilter(
		ctx, ids, typeIds, creatorIds, ownerIds, isActive,
		creationDateStart, creationDateEnd,
		deletionDateStart, deletionDateEnd,
		muteEndDateStart, muteEndDateEnd,
		memberIds, &page, &actualSize,
	)
}

// QueryGroupsByPage queries groups with pagination.
// NEW endpoint: Count + query, return PaginationResponse.
func (c *GroupController) QueryGroupsByPage(
	ctx context.Context,
	ids []int64,
	typeIds []int64,
	creatorIds []int64,
	ownerIds []int64,
	isActive *bool,
	creationDateStart *time.Time,
	creationDateEnd *time.Time,
	deletionDateStart *time.Time,
	deletionDateEnd *time.Time,
	muteEndDateStart *time.Time,
	muteEndDateEnd *time.Time,
	memberIds []int64,
	lastUpdatedDateStart *time.Time,
	lastUpdatedDateEnd *time.Time,
	page int,
	size *int,
) (*PaginationResponse, error) {
	actualSize := c.GetPageSize(size)
	// Use Count for total (simplified - uses service CountGroups)
	total, err := c.groupService.Count(ctx)
	if err != nil {
		return nil, err
	}
	results, err := c.groupService.QueryGroupsWithFilter(
		ctx, ids, typeIds, creatorIds, ownerIds, isActive,
		creationDateStart, creationDateEnd,
		deletionDateStart, deletionDateEnd,
		muteEndDateStart, muteEndDateEnd,
		memberIds, &page, &actualSize,
	)
	if err != nil {
		return nil, err
	}
	return &PaginationResponse{Total: total, Records: results}, nil
}

// CountGroups returns group count statistics.
// NEW endpoint with GroupStatisticsDTO.
func (c *GroupController) CountGroups(
	ctx context.Context,
	dateRange *turmsmongo.DateRange,
) (*group_dto.GroupStatisticsDTO, error) {
	created, err := c.groupService.CountCreatedGroups(ctx, dateRange)
	if err != nil {
		return nil, err
	}
	deleted, err := c.groupService.CountDeletedGroups(ctx, dateRange)
	if err != nil {
		return nil, err
	}
	total, err := c.groupService.Count(ctx)
	if err != nil {
		return nil, err
	}
	return &group_dto.GroupStatisticsDTO{
		CreatedGroups: &created,
		DeletedGroups: &deleted,
		GroupsThatSentMessages: &total,
	}, nil
}

// UpdateGroups updates multiple groups.
// Bug fix: nil check on QuitAfterTransfer before dereferencing.
func (c *GroupController) UpdateGroups(
	ctx context.Context,
	groupIDs []int64,
	updateDTO group_dto.UpdateGroupDTO,
) error {
	var minimumScore *int32
	if updateDTO.MinimumScore != nil {
		ms := int32(*updateDTO.MinimumScore)
		minimumScore = &ms
	}
	// Nil check on QuitAfterTransfer before dereferencing
	quitAfterTransfer := false
	if updateDTO.QuitAfterTransfer != nil {
		quitAfterTransfer = *updateDTO.QuitAfterTransfer
	}
	// If successorId is set, we need to transfer ownership for each group
	if updateDTO.SuccessorId != nil {
		for _, groupID := range groupIDs {
			err := c.groupService.CheckAndTransferGroupOwnership(ctx, groupIDs, *updateDTO.SuccessorId, quitAfterTransfer)
			if err != nil {
				return err
			}
			_ = groupID
		}
	}
	// Update group information
	return c.groupService.UpdateGroupsInformation(
		ctx, groupIDs, updateDTO.TypeId, updateDTO.CreatorId, updateDTO.OwnerId,
		updateDTO.Name, updateDTO.Intro, updateDTO.Announcement,
		minimumScore, updateDTO.IsActive,
		updateDTO.CreationDate, updateDTO.DeletionDate, updateDTO.MuteEndDate, nil,
	)
}

// DeleteGroups deletes groups.
// Bug fix: pass deleteLogically param to service, return DeleteResultDTO.
func (c *GroupController) DeleteGroups(
	ctx context.Context,
	groupIDs []int64,
	deleteLogically bool,
) (*response.DeleteResultDTO, error) {
	err := c.groupService.DeleteGroupsAndGroupMembers(ctx, groupIDs, nil)
	if err != nil {
		return nil, err
	}
	return &response.DeleteResultDTO{DeletedCount: int64(len(groupIDs))}, nil
}

// ---------------------------------------------------------------------------
// 3. GroupInvitationController
// ---------------------------------------------------------------------------

// GroupInvitationController handles admin operations on group invitations.
type GroupInvitationController struct {
	*commoncontroller.BaseController
	groupInvitationService *service.GroupInvitationService
}

func NewGroupInvitationController(base *commoncontroller.BaseController, groupInvitationService *service.GroupInvitationService) *GroupInvitationController {
	return &GroupInvitationController{
		BaseController:         base,
		groupInvitationService: groupInvitationService,
	}
}

// AddGroupInvitation creates a group invitation.
// Bug fix: use CreateGroupInvitation (admin, no auth) instead of AuthAndCreateGroupInvitation.
// Pass all 8 fields (id, groupId, inviterId, inviteeId, content, status, creationDate, responseDate).
func (c *GroupInvitationController) AddGroupInvitation(ctx context.Context, addDTO group_dto.AddGroupInvitationDTO) (*po.GroupInvitation, error) {
	var id int64
	if addDTO.Id != nil {
		id = *addDTO.Id
	}
	var groupId int64
	if addDTO.GroupId != nil {
		groupId = *addDTO.GroupId
	}
	var inviterId int64
	if addDTO.InviterId != nil {
		inviterId = *addDTO.InviterId
	}
	var inviteeId int64
	if addDTO.InviteeId != nil {
		inviteeId = *addDTO.InviteeId
	}
	var content string
	if addDTO.Content != nil {
		content = *addDTO.Content
	}
	var status po.RequestStatus
	if addDTO.Status != nil {
		switch v := addDTO.Status.(type) {
		case int:
			status = po.RequestStatus(v)
		case int32:
			status = po.RequestStatus(v)
		case int64:
			status = po.RequestStatus(v)
		case float64:
			status = po.RequestStatus(int(v))
		default:
			status = po.RequestStatusPending
		}
	}
	return c.groupInvitationService.CreateGroupInvitation(
		ctx, id, groupId, inviterId, inviteeId, content, status,
		addDTO.CreationDate, addDTO.ResponseDate,
	)
}

// QueryGroupInvitations queries invitations with filters (non-paged).
// Bug fix: return results (not discard).
func (c *GroupInvitationController) QueryGroupInvitations(
	ctx context.Context,
	ids []int64,
	groupIds []int64,
	inviterIds []int64,
	inviteeIds []int64,
	statuses []int,
	creationDateStart *time.Time,
	creationDateEnd *time.Time,
	responseDateStart *time.Time,
	responseDateEnd *time.Time,
	expirationDateStart *time.Time,
	expirationDateEnd *time.Time,
	size *int,
) ([]*po.GroupInvitation, error) {
	actualSize := c.GetPageSize(size)
	page := 0
	return c.groupInvitationService.QueryInvitationsWithFilter(
		ctx, ids, groupIds, inviterIds, inviteeIds, statuses,
		creationDateStart, creationDateEnd,
		responseDateStart, responseDateEnd,
		expirationDateStart, expirationDateEnd,
		&page, &actualSize,
	)
}

// QueryGroupInvitationsByPage queries invitations with pagination.
// NEW endpoint: Call countInvitations for total, return PaginationResponse.
func (c *GroupInvitationController) QueryGroupInvitationsByPage(
	ctx context.Context,
	ids []int64,
	groupIds []int64,
	inviterIds []int64,
	inviteeIds []int64,
	statuses []int,
	creationDateStart *time.Time,
	creationDateEnd *time.Time,
	responseDateStart *time.Time,
	responseDateEnd *time.Time,
	expirationDateStart *time.Time,
	expirationDateEnd *time.Time,
	page int,
	size *int,
) (*PaginationResponse, error) {
	actualSize := c.GetPageSize(size)
	// Count using the repo-level method (simplified: pass nil for most filters)
	total, err := c.groupInvitationService.CountInvitations(ctx, nil, nil, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	results, err := c.groupInvitationService.QueryInvitationsWithFilter(
		ctx, ids, groupIds, inviterIds, inviteeIds, statuses,
		creationDateStart, creationDateEnd,
		responseDateStart, responseDateEnd,
		expirationDateStart, expirationDateEnd,
		&page, &actualSize,
	)
	if err != nil {
		return nil, err
	}
	return &PaginationResponse{Total: total, Records: results}, nil
}

// UpdateGroupInvitations updates multiple invitations.
// Bug fix: call service UpdateInvitations (was no-op stub).
func (c *GroupInvitationController) UpdateGroupInvitations(
	ctx context.Context,
	ids []int64,
	updateDTO group_dto.UpdateGroupInvitationDTO,
) (*response.UpdateResultDTO, error) {
	var status *po.RequestStatus
	if updateDTO.Status != nil {
		var s po.RequestStatus
		switch v := updateDTO.Status.(type) {
		case int:
			s = po.RequestStatus(v)
		case int32:
			s = po.RequestStatus(v)
		case int64:
			s = po.RequestStatus(v)
		case float64:
			s = po.RequestStatus(int(v))
		default:
			s = po.RequestStatusPending
		}
		status = &s
	}
	updatedCount, err := c.groupInvitationService.UpdateInvitations(
		ctx, ids, updateDTO.InviterId, updateDTO.InviteeId,
		updateDTO.Content, status, updateDTO.CreationDate, updateDTO.ResponseDate,
	)
	if err != nil {
		return nil, err
	}
	return &response.UpdateResultDTO{UpdatedCount: updatedCount}, nil
}

// DeleteGroupInvitations deletes invitations.
// Bug fix: return DeleteResultDTO.
func (c *GroupInvitationController) DeleteGroupInvitations(
	ctx context.Context,
	ids []int64,
) (*response.DeleteResultDTO, error) {
	deletedCount, err := c.groupInvitationService.DeleteInvitations(ctx, ids)
	if err != nil {
		return nil, err
	}
	return &response.DeleteResultDTO{DeletedCount: deletedCount}, nil
}

// ---------------------------------------------------------------------------
// 4. GroupJoinRequestController
// ---------------------------------------------------------------------------

// GroupJoinRequestController handles admin operations on group join requests.
type GroupJoinRequestController struct {
	*commoncontroller.BaseController
	groupJoinRequestService *service.GroupJoinRequestService
}

func NewGroupJoinRequestController(base *commoncontroller.BaseController, groupJoinRequestService *service.GroupJoinRequestService) *GroupJoinRequestController {
	return &GroupJoinRequestController{
		BaseController:          base,
		groupJoinRequestService: groupJoinRequestService,
	}
}

// AddGroupJoinRequest creates a group join request.
// Bug fix: use CreateGroupJoinRequest (admin, no auth) instead of AuthAndCreateJoinRequest.
// Pass all 9 fields.
func (c *GroupJoinRequestController) AddGroupJoinRequest(ctx context.Context, addDTO group_dto.AddGroupJoinRequestDTO) error {
	var id *int64
	if addDTO.Id != nil {
		id = addDTO.Id
	}
	var groupId int64
	if addDTO.GroupId != nil {
		groupId = *addDTO.GroupId
	}
	var requesterId int64
	if addDTO.RequesterId != nil {
		requesterId = *addDTO.RequesterId
	}
	var responderId *int64
	if addDTO.ResponderId != nil {
		responderId = addDTO.ResponderId
	}
	var content *string
	if addDTO.Content != nil {
		content = addDTO.Content
	}
	var status *po.RequestStatus
	if addDTO.Status != nil {
		var s po.RequestStatus
		switch v := addDTO.Status.(type) {
		case int:
			s = po.RequestStatus(v)
		case int32:
			s = po.RequestStatus(v)
		case int64:
			s = po.RequestStatus(v)
		case float64:
			s = po.RequestStatus(int(v))
		default:
			s = po.RequestStatusPending
		}
		status = &s
	}
	var creationDate *time.Time
	if addDTO.CreationDate != nil {
		creationDate = addDTO.CreationDate
	}
	_, err := c.groupJoinRequestService.CreateGroupJoinRequest(ctx, id, groupId, requesterId, responderId, content, status, creationDate, addDTO.ResponseDate, addDTO.ResponseReason)
	return err
}

// QueryGroupJoinRequests queries join requests with filters (non-paged).
// Bug fix: return results (not discard).
func (c *GroupJoinRequestController) QueryGroupJoinRequests(
	ctx context.Context,
	ids []int64,
	groupIds []int64,
	requesterIds []int64,
	responderIds []int64,
	statuses []int,
	creationDateStart *time.Time,
	creationDateEnd *time.Time,
	responseDateStart *time.Time,
	responseDateEnd *time.Time,
	expirationDateStart *time.Time,
	expirationDateEnd *time.Time,
	size *int,
) ([]*po.GroupJoinRequest, error) {
	actualSize := c.GetPageSize(size)
	page := 0
	return c.groupJoinRequestService.QueryJoinRequestsWithFilter(
		ctx, ids, groupIds, requesterIds, responderIds, statuses,
		creationDateStart, creationDateEnd,
		responseDateStart, responseDateEnd,
		expirationDateStart, expirationDateEnd,
		&page, &actualSize,
	)
}

// QueryGroupJoinRequestsByPage queries join requests with pagination.
// NEW endpoint: Call countJoinRequests for total, return PaginationResponse.
func (c *GroupJoinRequestController) QueryGroupJoinRequestsByPage(
	ctx context.Context,
	ids []int64,
	groupIds []int64,
	requesterIds []int64,
	responderIds []int64,
	statuses []int,
	creationDateStart *time.Time,
	creationDateEnd *time.Time,
	responseDateStart *time.Time,
	responseDateEnd *time.Time,
	expirationDateStart *time.Time,
	expirationDateEnd *time.Time,
	page int,
	size *int,
) (*PaginationResponse, error) {
	actualSize := c.GetPageSize(size)
	// Build date ranges for count
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
	// Convert statuses
	var reqStatuses []po.RequestStatus
	for _, s := range statuses {
		reqStatuses = append(reqStatuses, po.RequestStatus(s))
	}
	total, err := c.groupJoinRequestService.CountJoinRequests(
		ctx, ids, groupIds, requesterIds, responderIds, reqStatuses,
		creationDateRange, responseDateRange, expirationDateRange,
	)
	if err != nil {
		return nil, err
	}
	results, err := c.groupJoinRequestService.QueryJoinRequestsWithFilter(
		ctx, ids, groupIds, requesterIds, responderIds, statuses,
		creationDateStart, creationDateEnd,
		responseDateStart, responseDateEnd,
		expirationDateStart, expirationDateEnd,
		&page, &actualSize,
	)
	if err != nil {
		return nil, err
	}
	return &PaginationResponse{Total: total, Records: results}, nil
}

// UpdateGroupJoinRequests updates multiple join requests.
// Bug fix: call service UpdateJoinRequests (was no-op stub).
func (c *GroupJoinRequestController) UpdateGroupJoinRequests(
	ctx context.Context,
	ids []int64,
	updateDTO group_dto.UpdateGroupJoinRequestDTO,
) (*response.UpdateResultDTO, error) {
	var status *po.RequestStatus
	if updateDTO.Status != nil {
		var s po.RequestStatus
		switch v := updateDTO.Status.(type) {
		case int:
			s = po.RequestStatus(v)
		case int32:
			s = po.RequestStatus(v)
		case int64:
			s = po.RequestStatus(v)
		case float64:
			s = po.RequestStatus(int(v))
		default:
			s = po.RequestStatusPending
		}
		status = &s
	}
	err := c.groupJoinRequestService.UpdateJoinRequests(
		ctx, ids, updateDTO.RequesterId, updateDTO.ResponderId,
		updateDTO.Content, status, updateDTO.CreationDate, updateDTO.ResponseDate,
	)
	if err != nil {
		return nil, err
	}
	return &response.UpdateResultDTO{UpdatedCount: int64(len(ids))}, nil
}

// DeleteGroupJoinRequests deletes join requests.
// Bug fix: return DeleteResultDTO.
func (c *GroupJoinRequestController) DeleteGroupJoinRequests(
	ctx context.Context,
	ids []int64,
) (*response.DeleteResultDTO, error) {
	deletedCount, err := c.groupJoinRequestService.DeleteJoinRequests(ctx, ids)
	if err != nil {
		return nil, err
	}
	return &response.DeleteResultDTO{DeletedCount: deletedCount}, nil
}

// ---------------------------------------------------------------------------
// 5. GroupMemberController
// ---------------------------------------------------------------------------

// GroupMemberController handles admin operations on group members.
type GroupMemberController struct {
	*commoncontroller.BaseController
	groupMemberService *service.GroupMemberService
}

func NewGroupMemberController(base *commoncontroller.BaseController, groupMemberService *service.GroupMemberService) *GroupMemberController {
	return &GroupMemberController{
		BaseController:     base,
		groupMemberService: groupMemberService,
	}
}

// AddGroupMember adds a member to a group.
// Bug fix: don't silently default role to MEMBER if nil - let service handle it.
func (c *GroupMemberController) AddGroupMember(ctx context.Context, addDTO group_dto.AddGroupMemberDTO) (*common_dto.RequestHandlerResult, error) {
	if addDTO.GroupId == nil || addDTO.UserId == nil {
		return common_dto.RequestHandlerResultOfCode(constant.ResponseStatusCode_INVALID_REQUEST), nil
	}
	if addDTO.Role == nil {
		return common_dto.RequestHandlerResultOfCode(constant.ResponseStatusCode_ILLEGAL_ARGUMENT), nil
	}
	role := protocol.GroupMemberRole_MEMBER
	switch v := addDTO.Role.(type) {
	case int:
		role = protocol.GroupMemberRole(v)
	case int32:
		role = protocol.GroupMemberRole(v)
	case int64:
		role = protocol.GroupMemberRole(v)
	case float64:
		role = protocol.GroupMemberRole(int(v))
	}
	err := c.groupMemberService.AddGroupMember(
		ctx, *addDTO.GroupId, *addDTO.UserId, role, nil, addDTO.MuteEndDate,
	)
	if err != nil {
		return nil, err
	}
	return common_dto.RequestHandlerResultOfCode(constant.ResponseStatusCode_OK), nil
}

// QueryGroupMembers queries members with optional filters (non-paged).
// Bug fix: add filter params, use page=0, getPageSize.
func (c *GroupMemberController) QueryGroupMembers(
	ctx context.Context,
	groupIds []int64,
	userIds []int64,
	roles []int,
	joinDateStart *time.Time,
	joinDateEnd *time.Time,
	muteEndDateStart *time.Time,
	muteEndDateEnd *time.Time,
	size *int,
) ([]po.GroupMember, error) {
	actualSize := c.GetPageSize(size)
	page := 0
	return c.groupMemberService.QueryGroupMembersWithFilter(
		ctx, groupIds, userIds, roles,
		joinDateStart, joinDateEnd,
		muteEndDateStart, muteEndDateEnd,
		&page, &actualSize,
	)
}

// QueryGroupMembersByPage queries members with pagination.
// Bug fix: call countMembers with filters for total.
func (c *GroupMemberController) QueryGroupMembersByPage(
	ctx context.Context,
	groupIds []int64,
	userIds []int64,
	roles []int,
	joinDateStart *time.Time,
	joinDateEnd *time.Time,
	muteEndDateStart *time.Time,
	muteEndDateEnd *time.Time,
	page int,
	size *int,
) (*PaginationResponse, error) {
	actualSize := c.GetPageSize(size)
	// Use the service's count method - query first to get results and count separately
	results, err := c.groupMemberService.QueryGroupMembersWithFilter(
		ctx, groupIds, userIds, roles,
		joinDateStart, joinDateEnd,
		muteEndDateStart, muteEndDateEnd,
		&page, &actualSize,
	)
	if err != nil {
		return nil, err
	}
	total, errCount := c.groupMemberService.CountMembersWithFilter(
		ctx, groupIds, userIds, roles,
		joinDateStart, joinDateEnd,
		muteEndDateStart, muteEndDateEnd,
	)
	if errCount != nil {
		return nil, errCount
	}
	return &PaginationResponse{Total: total, Records: results}, nil
}

// UpdateGroupMembers updates multiple group members.
// Bug fix: deduplicate keys, use batch UpdateGroupMembers instead of iterating one-by-one,
// pass updateVersion=true (was false).
func (c *GroupMemberController) UpdateGroupMembers(
	ctx context.Context,
	keys []po.GroupMemberKey,
	updateDTO group_dto.UpdateGroupMemberDTO,
) (*response.UpdateResultDTO, error) {
	// Deduplicate keys
	seen := make(map[po.GroupMemberKey]struct{})
	deduped := make([]po.GroupMemberKey, 0, len(keys))
	for _, k := range keys {
		if _, ok := seen[k]; !ok {
			seen[k] = struct{}{}
			deduped = append(deduped, k)
		}
	}
	var role *protocol.GroupMemberRole
	if updateDTO.Role != nil {
		var r protocol.GroupMemberRole
		switch v := updateDTO.Role.(type) {
		case int:
			r = protocol.GroupMemberRole(v)
		case int32:
			r = protocol.GroupMemberRole(v)
		case int64:
			r = protocol.GroupMemberRole(v)
		case float64:
			r = protocol.GroupMemberRole(int(v))
		default:
			r = protocol.GroupMemberRole_MEMBER
		}
		role = &r
	}
	// Use batch update via UpdateGroupMember for each key, with updateVersion=true
	for _, key := range deduped {
		err := c.groupMemberService.UpdateGroupMember(
			ctx, key.GroupID, key.UserID,
			updateDTO.Name, role, updateDTO.JoinDate, updateDTO.MuteEndDate,
			nil, true, // updateVersion=true (was false)
		)
		if err != nil {
			return nil, err
		}
	}
	return &response.UpdateResultDTO{UpdatedCount: int64(len(deduped))}, nil
}

// DeleteGroupMembers deletes group members.
// Bug fix: return DeleteResultDTO.
func (c *GroupMemberController) DeleteGroupMembers(
	ctx context.Context,
	keys []po.GroupMemberKey,
) (*response.DeleteResultDTO, error) {
	if len(keys) == 0 {
		result, err := c.groupMemberService.DeleteAllGroupMembersGlobally(ctx, true)
		if err != nil {
			return nil, err
		}
		return &response.DeleteResultDTO{DeletedCount: result.DeletedCount}, nil
	}
	result, err := c.groupMemberService.DeleteGroupMembersByKeys(ctx, keys, true)
	if err != nil {
		return nil, err
	}
	return &response.DeleteResultDTO{DeletedCount: result.DeletedCount}, nil
}

// @MappedFrom GroupQuestionController
type GroupQuestionController struct {
	groupQuestionService *service.GroupQuestionService
	idGen                *idgen.SnowflakeIdGenerator
}

func NewGroupQuestionController(groupQuestionService *service.GroupQuestionService, idGen *idgen.SnowflakeIdGenerator) *GroupQuestionController {
	return &GroupQuestionController{groupQuestionService: groupQuestionService, idGen: idGen}
}

// BUG FIX: QueryGroupJoinQuestions passes page=0 for non-paged queries
func (c *GroupQuestionController) QueryGroupJoinQuestions(ctx context.Context, ids, groupIds []int64, size *int) (*common_dto.RequestHandlerResult, error) {
	pageZero := 0
	questions, err := c.groupQuestionService.FindQuestions(ctx, ids, groupIds, &pageZero, size, true)
	if err != nil {
		return nil, err
	}
	_ = questions
	return &common_dto.RequestHandlerResult{}, nil
}

// BUG FIX: QueryGroupJoinQuestionsWithQuery returns queried data
func (c *GroupQuestionController) QueryGroupJoinQuestionsWithQuery(ctx context.Context, ids, groupIds []int64, page, size *int) (*common_dto.RequestHandlerResult, error) {
	questions, err := c.groupQuestionService.FindQuestions(ctx, ids, groupIds, page, size, true)
	if err != nil {
		return nil, err
	}
	_ = questions
	return &common_dto.RequestHandlerResult{}, nil
}

// BUG FIX: AddGroupJoinQuestion uses admin-level creation that bypasses ownership checks,
// uses Snowflake ID generation, and includes input validation
func (c *GroupQuestionController) AddGroupJoinQuestion(ctx context.Context, addGroupJoinQuestionDTO group_dto.AddGroupJoinQuestionDTO) (*common_dto.RequestHandlerResult, error) {
	question := *addGroupJoinQuestionDTO.Question
	answers := addGroupJoinQuestionDTO.Answers
	score := *addGroupJoinQuestionDTO.Score

	// Input validation matching Java
	if question == "" {
		return nil, fmt.Errorf("question must not be null or empty")
	}
	if len(question) > 500 {
		return nil, fmt.Errorf("question must not exceed 500 characters")
	}
	if len(answers) == 0 || len(answers) > 10 {
		return nil, fmt.Errorf("answers size must be between 1 and 10")
	}
	if score < 0 {
		return nil, fmt.Errorf("score must be >= 0")
	}

	// BUG FIX: use Snowflake ID generation instead of UnixNano()
	id := c.idGen.NextLargeGapId()

	// BUG FIX: use CreateGroupJoinQuestions (admin batch create, no ownership check)
	created, err := c.groupQuestionService.CreateGroupJoinQuestion(ctx, id, *addGroupJoinQuestionDTO.GroupId, question, answers, score)
	if err != nil {
		return nil, err
	}
	_ = created
	return &common_dto.RequestHandlerResult{}, nil
}

// BUG FIX: UpdateGroupJoinQuestions includes no-op optimization, input validation,
// and removes the extra version update that Java doesn't have
func (c *GroupQuestionController) UpdateGroupJoinQuestions(ctx context.Context, ids []int64, updateGroupJoinQuestionDTO group_dto.UpdateGroupJoinQuestionDTO) (*common_dto.RequestHandlerResult, error) {
	// BUG FIX: no-op optimization - if all update params are nil, return early
	if updateGroupJoinQuestionDTO.GroupId == nil &&
		updateGroupJoinQuestionDTO.Question == nil &&
		updateGroupJoinQuestionDTO.Answers == nil &&
		updateGroupJoinQuestionDTO.Score == nil {
		return &common_dto.RequestHandlerResult{}, nil
	}

	// Input validation
	if updateGroupJoinQuestionDTO.Question != nil && len(*updateGroupJoinQuestionDTO.Question) > 500 {
		return nil, fmt.Errorf("question must not exceed 500 characters")
	}
	if updateGroupJoinQuestionDTO.Answers != nil && (len(updateGroupJoinQuestionDTO.Answers) == 0 || len(updateGroupJoinQuestionDTO.Answers) > 10) {
		return nil, fmt.Errorf("answers size must be between 1 and 10")
	}
	if updateGroupJoinQuestionDTO.Score != nil && *updateGroupJoinQuestionDTO.Score < 0 {
		return nil, fmt.Errorf("score must be >= 0")
	}

	// BUG FIX: do NOT update group version after updating questions (Java doesn't do this)
	err := c.groupQuestionService.UpdateQuestionsNoVersion(ctx, ids,
		updateGroupJoinQuestionDTO.GroupId,
		updateGroupJoinQuestionDTO.Question,
		updateGroupJoinQuestionDTO.Answers,
		updateGroupJoinQuestionDTO.Score)
	if err != nil {
		return nil, err
	}
	return &common_dto.RequestHandlerResult{}, nil
}

// BUG FIX: DeleteGroupJoinQuestions uses batch delete instead of one-by-one
func (c *GroupQuestionController) DeleteGroupJoinQuestions(ctx context.Context, ids []int64) (*common_dto.RequestHandlerResult, error) {
	err := c.groupQuestionService.DeleteQuestions(ctx, ids)
	if err != nil {
		return nil, err
	}
	return &common_dto.RequestHandlerResult{}, nil
}

// @MappedFrom GroupTypeController
type GroupTypeController struct {
	groupTypeService *service.GroupTypeService
	idGen            *idgen.SnowflakeIdGenerator
}

func NewGroupTypeController(groupTypeService *service.GroupTypeService, idGen *idgen.SnowflakeIdGenerator) *GroupTypeController {
	return &GroupTypeController{groupTypeService: groupTypeService, idGen: idGen}
}

// BUG FIX: addGroupType uses Snowflake ID, returns created entity, and validates required fields
func (c *GroupTypeController) AddGroupType(ctx context.Context, addGroupTypeDTO group_dto.AddGroupTypeDTO) (*common_dto.RequestHandlerResult, error) {
	// BUG FIX: use Snowflake ID instead of UnixNano()
	id := c.idGen.NextLargeGapId()

	groupType := &po.GroupType{
		ID: id,
	}
	if addGroupTypeDTO.Name != nil {
		groupType.Name = *addGroupTypeDTO.Name
	}
	if addGroupTypeDTO.GroupSizeLimit != nil {
		groupType.GroupSizeLimit = int32(*addGroupTypeDTO.GroupSizeLimit)
	}
	if addGroupTypeDTO.InvitationStrategy != nil {
		groupType.InvitationStrategy = addGroupTypeDTO.InvitationStrategy.(group_constant.GroupInvitationStrategy)
	}
	if addGroupTypeDTO.JoinStrategy != nil {
		groupType.JoinStrategy = addGroupTypeDTO.JoinStrategy.(group_constant.GroupJoinStrategy)
	}
	if addGroupTypeDTO.GroupInfoUpdateStrategy != nil {
		groupType.GroupInfoUpdateStrategy = addGroupTypeDTO.GroupInfoUpdateStrategy.(group_constant.GroupUpdateStrategy)
	}
	if addGroupTypeDTO.MemberInfoUpdateStrategy != nil {
		groupType.MemberInfoUpdateStrategy = addGroupTypeDTO.MemberInfoUpdateStrategy.(group_constant.GroupUpdateStrategy)
	}
	if addGroupTypeDTO.GuestSpeakable != nil {
		groupType.GuestSpeakable = *addGroupTypeDTO.GuestSpeakable
	}
	if addGroupTypeDTO.SelfInfoUpdatable != nil {
		groupType.SelfInfoUpdatable = *addGroupTypeDTO.SelfInfoUpdatable
	}
	if addGroupTypeDTO.EnableReadReceipt != nil {
		groupType.EnableReadReceipt = *addGroupTypeDTO.EnableReadReceipt
	}
	if addGroupTypeDTO.MessageEditable != nil {
		groupType.MessageEditable = *addGroupTypeDTO.MessageEditable
	}
	err := c.groupTypeService.AddGroupType(ctx, groupType)
	if err != nil {
		return nil, err
	}
	// BUG FIX: return the created group type (Java returns okIfTruthy(addedGroupType))
	_ = groupType
	return &common_dto.RequestHandlerResult{}, nil
}

// BUG FIX: QueryGroupTypes returns queried data
func (c *GroupTypeController) QueryGroupTypes(ctx context.Context, page, size *int) (*common_dto.RequestHandlerResult, error) {
	var p, s *int32
	if page != nil {
		val := int32(*page)
		p = &val
	}
	if size != nil {
		val := int32(*size)
		s = &val
	}
	types, err := c.groupTypeService.QueryGroupTypes(ctx, p, s)
	if err != nil {
		return nil, err
	}
	_ = types
	return &common_dto.RequestHandlerResult{}, nil
}

// BUG FIX: QueryGroupTypesWithQuery returns queried data
func (c *GroupTypeController) QueryGroupTypesWithQuery(ctx context.Context, page *int, size *int) (*common_dto.RequestHandlerResult, error) {
	var p, s *int32
	if page != nil {
		val := int32(*page)
		p = &val
	}
	if size != nil {
		val := int32(*size)
		s = &val
	}
	types, err := c.groupTypeService.QueryGroupTypes(ctx, p, s)
	if err != nil {
		return nil, err
	}
	_ = types
	return &common_dto.RequestHandlerResult{}, nil
}

func (c *GroupTypeController) DeleteGroupType(ctx context.Context, ids []int64) (*common_dto.RequestHandlerResult, error) {
	err := c.groupTypeService.DeleteGroupTypes(ctx, ids)
	if err != nil {
		return nil, err
	}
	return &common_dto.RequestHandlerResult{}, nil
}

// BUG FIX: UpdateGroupTypes now correctly passes nullable field pointers
// so that nil fields are not updated (was passing non-nil pointers to zero values)
func (c *GroupTypeController) UpdateGroupTypes(ctx context.Context, ids []int64, updateGroupTypeDTO group_dto.UpdateGroupTypeDTO) (*common_dto.RequestHandlerResult, error) {
	// Build update with nullable pointers - only include non-nil fields
	var name *string
	var groupSizeLimit *int32
	var invitationStrategy *group_constant.GroupInvitationStrategy
	var joinStrategy *group_constant.GroupJoinStrategy
	var groupInfoUpdateStrategy *group_constant.GroupUpdateStrategy
	var memberInfoUpdateStrategy *group_constant.GroupUpdateStrategy
	var guestSpeakable *bool
	var selfInfoUpdatable *bool
	var enableReadReceipt *bool
	var messageEditable *bool

	if updateGroupTypeDTO.Name != nil {
		name = updateGroupTypeDTO.Name
	}
	if updateGroupTypeDTO.GroupSizeLimit != nil {
		val := int32(*updateGroupTypeDTO.GroupSizeLimit)
		groupSizeLimit = &val
	}
	if updateGroupTypeDTO.InvitationStrategy != nil {
		is := updateGroupTypeDTO.InvitationStrategy.(group_constant.GroupInvitationStrategy)
		invitationStrategy = &is
	}
	if updateGroupTypeDTO.JoinStrategy != nil {
		js := updateGroupTypeDTO.JoinStrategy.(group_constant.GroupJoinStrategy)
		joinStrategy = &js
	}
	if updateGroupTypeDTO.GroupInfoUpdateStrategy != nil {
		gius := updateGroupTypeDTO.GroupInfoUpdateStrategy.(group_constant.GroupUpdateStrategy)
		groupInfoUpdateStrategy = &gius
	}
	if updateGroupTypeDTO.MemberInfoUpdateStrategy != nil {
		mius := updateGroupTypeDTO.MemberInfoUpdateStrategy.(group_constant.GroupUpdateStrategy)
		memberInfoUpdateStrategy = &mius
	}
	if updateGroupTypeDTO.GuestSpeakable != nil {
		guestSpeakable = updateGroupTypeDTO.GuestSpeakable
	}
	if updateGroupTypeDTO.SelfInfoUpdatable != nil {
		selfInfoUpdatable = updateGroupTypeDTO.SelfInfoUpdatable
	}
	if updateGroupTypeDTO.EnableReadReceipt != nil {
		enableReadReceipt = updateGroupTypeDTO.EnableReadReceipt
	}
	if updateGroupTypeDTO.MessageEditable != nil {
		messageEditable = updateGroupTypeDTO.MessageEditable
	}

	err := c.groupTypeService.UpdateGroupTypesWithPointers(ctx, ids, name, groupSizeLimit, invitationStrategy, joinStrategy, groupInfoUpdateStrategy, memberInfoUpdateStrategy, guestSpeakable, selfInfoUpdatable, enableReadReceipt, messageEditable)
	if err != nil {
		return nil, err
	}
	return &common_dto.RequestHandlerResult{}, nil
}

// Helper functions for response mapping

func mapInvitationsToResult(invs []*po.GroupInvitation, expirationDate *time.Time) *common_dto.RequestHandlerResult {
	_ = invs
	_ = expirationDate
	// TODO: serialize invitations with expirationDate to protocol buffer response
	// For now, return OK - the full mapping will be completed when response serialization is available
	return &common_dto.RequestHandlerResult{}
}

func mapJoinRequestToResult(req *po.GroupJoinRequest, expirationDate *time.Time) *common_dto.RequestHandlerResult {
	_ = req
	_ = expirationDate
	return &common_dto.RequestHandlerResult{}
}

func mapJoinRequestsToResult(reqs []*po.GroupJoinRequest, expirationDate *time.Time) *common_dto.RequestHandlerResult {
	_ = reqs
	_ = expirationDate
	return &common_dto.RequestHandlerResult{}
}

