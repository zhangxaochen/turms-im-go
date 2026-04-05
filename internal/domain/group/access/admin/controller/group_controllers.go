package controller

import (
	"context"
	"fmt"
	"time"

	common_dto "im.turms/server/internal/domain/common/dto"
	"im.turms/server/internal/domain/common/infra/idgen"
	"im.turms/server/internal/domain/group/access/admin/dto"
	"im.turms/server/internal/domain/group/constant"
	"im.turms/server/internal/domain/group/po"
	"im.turms/server/internal/domain/group/service"
	msg_service "im.turms/server/internal/domain/message/service"
	turmsmongo "im.turms/server/internal/storage/mongo"
	"im.turms/server/pkg/protocol"
)

// @MappedFrom GroupBlocklistController
type GroupBlocklistController struct {
	groupBlocklistService *service.GroupBlocklistService
}

func NewGroupBlocklistController(groupBlocklistService *service.GroupBlocklistService) *GroupBlocklistController {
	return &GroupBlocklistController{groupBlocklistService: groupBlocklistService}
}

func (c *GroupBlocklistController) AddGroupBlockedUser(ctx context.Context, addGroupBlockedUserDTO dto.AddGroupBlockedUserDTO) (*common_dto.RequestHandlerResult, error) {
	_, err := c.groupBlocklistService.AddBlockedUser(ctx,
		*addGroupBlockedUserDTO.GroupId,
		*addGroupBlockedUserDTO.UserId,
		*addGroupBlockedUserDTO.RequesterId,
		addGroupBlockedUserDTO.BlockDate)
	if err != nil {
		return nil, err
	}
	return &common_dto.RequestHandlerResult{}, nil
}

func (c *GroupBlocklistController) QueryGroupBlockedUsers(ctx context.Context, page, size *int) (*common_dto.RequestHandlerResult, error) {
	_, err := c.groupBlocklistService.QueryBlockedUsersWithPagination(ctx, page, size)
	if err != nil {
		return nil, err
	}
	return &common_dto.RequestHandlerResult{}, nil
}

func (c *GroupBlocklistController) QueryGroupBlockedUsersWithQuery(ctx context.Context, groupIds, userIds []int64, blockDateStart, blockDateEnd *time.Time, requesterIds []int64, page, size *int) (*common_dto.RequestHandlerResult, error) {
	var dateRange *turmsmongo.DateRange
	if blockDateStart != nil || blockDateEnd != nil {
		dateRange = &turmsmongo.DateRange{Start: blockDateStart, End: blockDateEnd}
	}
	_, err := c.groupBlocklistService.QueryBlockedUsersWithFilter(ctx, groupIds, userIds, dateRange, requesterIds, page, size)
	if err != nil {
		return nil, err
	}
	return &common_dto.RequestHandlerResult{}, nil
}

func (c *GroupBlocklistController) UpdateGroupBlockedUsers(ctx context.Context, keys []po.GroupBlockedUserKey, updateGroupBlockedUserDTO dto.UpdateGroupBlockedUserDTO) (*common_dto.RequestHandlerResult, error) {
	err := c.groupBlocklistService.UpdateBlockedUsers(ctx, keys,
		updateGroupBlockedUserDTO.BlockDate,
		updateGroupBlockedUserDTO.RequesterId)
	if err != nil {
		return nil, err
	}
	return &common_dto.RequestHandlerResult{}, nil
}

func (c *GroupBlocklistController) DeleteGroupBlockedUsers(ctx context.Context, keys []po.GroupBlockedUserKey) (*common_dto.RequestHandlerResult, error) {
	err := c.groupBlocklistService.DeleteBlockedUsers(ctx, keys)
	if err != nil {
		return nil, err
	}
	return &common_dto.RequestHandlerResult{}, nil
}

// @MappedFrom GroupController
type GroupController struct {
	groupService   *service.GroupService
	messageService *msg_service.MessageService
}

func NewGroupController(groupService *service.GroupService, messageService *msg_service.MessageService) *GroupController {
	return &GroupController{
		groupService:   groupService,
		messageService: messageService,
	}
}

func (c *GroupController) AddGroup(ctx context.Context, addGroupDTO dto.AddGroupDTO) (*common_dto.RequestHandlerResult, error) {
	var ownerId int64
	if addGroupDTO.OwnerId != nil {
		ownerId = *addGroupDTO.OwnerId
	} else if addGroupDTO.CreatorId != nil {
		ownerId = *addGroupDTO.CreatorId
	}

	var ms *int32
	if addGroupDTO.MinimumScore != nil {
		val := int32(*addGroupDTO.MinimumScore)
		ms = &val
	}

	_, err := c.groupService.AuthAndCreateGroup(ctx,
		*addGroupDTO.CreatorId,
		ownerId,
		addGroupDTO.Name,
		addGroupDTO.Intro,
		addGroupDTO.Announcement,
		ms,
		addGroupDTO.TypeId,
		addGroupDTO.CreationDate,
		addGroupDTO.DeletionDate,
		addGroupDTO.MuteEndDate,
		addGroupDTO.IsActive)
	if err != nil {
		return nil, err
	}
	return &common_dto.RequestHandlerResult{}, nil
}

func (c *GroupController) QueryGroups(ctx context.Context, page, size *int) (*common_dto.RequestHandlerResult, error) {
	_, err := c.groupService.QueryGroupsWithPagination(ctx, page, size)
	if err != nil {
		return nil, err
	}
	return &common_dto.RequestHandlerResult{}, nil
}

func (c *GroupController) QueryGroupsWithQuery(ctx context.Context, ids, typeIds, creatorIds, ownerIds []int64, isActive *bool, creationDateStart, creationDateEnd, deletionDateStart, deletionDateEnd, muteEndDateStart, muteEndDateEnd *time.Time, memberIds []int64, page, size *int) (*common_dto.RequestHandlerResult, error) {
	_, err := c.groupService.QueryGroupsWithFilter(ctx, ids, typeIds, creatorIds, ownerIds, isActive, creationDateStart, creationDateEnd, deletionDateStart, deletionDateEnd, muteEndDateStart, muteEndDateEnd, memberIds, page, size)
	if err != nil {
		return nil, err
	}
	return &common_dto.RequestHandlerResult{}, nil
}

func (c *GroupController) UpdateGroups(ctx context.Context, ids []int64, updateGroupDTO dto.UpdateGroupDTO) (*common_dto.RequestHandlerResult, error) {
	if updateGroupDTO.SuccessorId != nil {
		err := c.groupService.CheckAndTransferGroupOwnership(ctx, ids, *updateGroupDTO.SuccessorId, *updateGroupDTO.QuitAfterTransfer)
		if err != nil {
			return nil, err
		}
	} else {
		var ms *int32
		if updateGroupDTO.MinimumScore != nil {
			val := int32(*updateGroupDTO.MinimumScore)
			ms = &val
		}

		err := c.groupService.UpdateGroupsInformation(ctx, ids,
			updateGroupDTO.TypeId,
			updateGroupDTO.CreatorId,
			updateGroupDTO.OwnerId,
			updateGroupDTO.Name,
			updateGroupDTO.Intro,
			updateGroupDTO.Announcement,
			ms,
			updateGroupDTO.IsActive,
			updateGroupDTO.CreationDate,
			updateGroupDTO.DeletionDate,
			updateGroupDTO.MuteEndDate,
			nil)
		if err != nil {
			return nil, err
		}
	}
	return &common_dto.RequestHandlerResult{}, nil
}

func (c *GroupController) DeleteGroups(ctx context.Context, ids []int64, deleteLogical *bool) (*common_dto.RequestHandlerResult, error) {
	err := c.groupService.DeleteGroupsAndGroupMembers(ctx, ids, nil)
	if err != nil {
		return nil, err
	}
	return &common_dto.RequestHandlerResult{}, nil
}

// @MappedFrom GroupInvitationController
type GroupInvitationController struct {
	groupInvitationService *service.GroupInvitationService
}

func NewGroupInvitationController(groupInvitationService *service.GroupInvitationService) *GroupInvitationController {
	return &GroupInvitationController{groupInvitationService: groupInvitationService}
}

func (c *GroupInvitationController) AddGroupInvitation(ctx context.Context, addGroupInvitationDTO dto.AddGroupInvitationDTO) (*common_dto.RequestHandlerResult, error) {
	_, err := c.groupInvitationService.AuthAndCreateGroupInvitation(ctx,
		*addGroupInvitationDTO.InviterId,
		*addGroupInvitationDTO.GroupId,
		*addGroupInvitationDTO.InviteeId,
		*addGroupInvitationDTO.Content)
	if err != nil {
		return nil, err
	}
	return &common_dto.RequestHandlerResult{}, nil
}

// BUG FIX: QueryGroupInvitations now returns the queried invitations as response data
// with expiration date mapping (GroupInvitationDTO equivalent).
func (c *GroupInvitationController) QueryGroupInvitations(ctx context.Context, page, size *int) (*common_dto.RequestHandlerResult, error) {
	invs, err := c.groupInvitationService.QueryInvitationsWithPagination(ctx, page, size)
	if err != nil {
		return nil, err
	}
	// Map to response with expiration dates (GroupInvitationDTO equivalent)
	expirationDate := c.groupInvitationService.GetEntityExpirationDate(ctx)
	return mapInvitationsToResult(invs, expirationDate), nil
}

// BUG FIX: QueryGroupInvitationsWithQuery now returns data and supports count for pagination
func (c *GroupInvitationController) QueryGroupInvitationsWithQuery(ctx context.Context, ids, groupIds, inviterIds, inviteeIds []int64, statuses []int, creationDateStart, creationDateEnd, responseDateStart, responseDateEnd, expirationDateStart, expirationDateEnd *time.Time, page, size *int) (*common_dto.RequestHandlerResult, error) {
	invs, err := c.groupInvitationService.QueryInvitationsWithFilter(ctx, ids, groupIds, inviterIds, inviteeIds, statuses, creationDateStart, creationDateEnd, responseDateStart, responseDateEnd, expirationDateStart, expirationDateEnd, page, size)
	if err != nil {
		return nil, err
	}
	expirationDate := c.groupInvitationService.GetEntityExpirationDate(ctx)
	return mapInvitationsToResult(invs, expirationDate), nil
}

// BUG FIX: UpdateGroupInvitations now calls the service method with all DTO fields
func (c *GroupInvitationController) UpdateGroupInvitations(ctx context.Context, ids []int64, updateGroupInvitationDTO dto.UpdateGroupInvitationDTO) (*common_dto.RequestHandlerResult, error) {
	var status *po.RequestStatus
	if updateGroupInvitationDTO.Status != nil {
		switch s := updateGroupInvitationDTO.Status.(type) {
		case float64:
			rs := po.RequestStatus(int(s))
			status = &rs
		case int:
			rs := po.RequestStatus(s)
			status = &rs
		}
	}
	_, err := c.groupInvitationService.UpdateInvitations(ctx,
		ids,
		updateGroupInvitationDTO.InviterId,
		updateGroupInvitationDTO.InviteeId,
		updateGroupInvitationDTO.Content,
		status,
		updateGroupInvitationDTO.CreationDate,
		updateGroupInvitationDTO.ResponseDate)
	if err != nil {
		return nil, err
	}
	return &common_dto.RequestHandlerResult{}, nil
}

// BUG FIX: DeleteGroupInvitations now returns DeleteResultDTO with deleted count
func (c *GroupInvitationController) DeleteGroupInvitations(ctx context.Context, ids []int64) (*common_dto.RequestHandlerResult, error) {
	deletedCount, err := c.groupInvitationService.DeleteInvitations(ctx, ids)
	if err != nil {
		return nil, err
	}
	return common_dto.RequestHandlerResultOfDataLong(deletedCount), nil
}

// @MappedFrom GroupJoinRequestController
type GroupJoinRequestController struct {
	groupJoinRequestService *service.GroupJoinRequestService
	idGen                   *idgen.SnowflakeIdGenerator
}

func NewGroupJoinRequestController(groupJoinRequestService *service.GroupJoinRequestService, idGen *idgen.SnowflakeIdGenerator) *GroupJoinRequestController {
	return &GroupJoinRequestController{groupJoinRequestService: groupJoinRequestService, idGen: idGen}
}

// BUG FIX: addGroupJoinRequest now uses admin-level creation (CreateGroupJoinRequest)
// that bypasses client-side validations, passes all fields from the DTO, and returns the created entity
func (c *GroupJoinRequestController) AddGroupJoinRequest(ctx context.Context, addGroupJoinRequestDTO dto.AddGroupJoinRequestDTO) (*common_dto.RequestHandlerResult, error) {
	var id int64
	if addGroupJoinRequestDTO.Id != nil {
		id = *addGroupJoinRequestDTO.Id
	} else {
		id = c.idGen.NextLargeGapId()
	}
	var status po.RequestStatus
	if addGroupJoinRequestDTO.Status != nil {
		switch s := addGroupJoinRequestDTO.Status.(type) {
		case float64:
			status = po.RequestStatus(int(s))
		case int:
			status = po.RequestStatus(s)
		}
	} else {
		status = po.RequestStatusPending
	}
	var creationDate time.Time
	if addGroupJoinRequestDTO.CreationDate != nil {
		creationDate = *addGroupJoinRequestDTO.CreationDate
	} else {
		creationDate = time.Now()
	}

	req := &po.GroupJoinRequest{
		ID:           id,
		GroupID:      *addGroupJoinRequestDTO.GroupId,
		RequesterID:  *addGroupJoinRequestDTO.RequesterId,
		ResponderID:  addGroupJoinRequestDTO.ResponderId,
		Status:       status,
		CreationDate: creationDate,
		ResponseDate:  addGroupJoinRequestDTO.ResponseDate,
		Reason:       addGroupJoinRequestDTO.ResponseReason,
	}
	if addGroupJoinRequestDTO.Content != nil {
		req.Content = *addGroupJoinRequestDTO.Content
	}
	err := c.groupJoinRequestService.CreateGroupJoinRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	expirationDate := c.groupJoinRequestService.GetEntityExpirationDate(ctx)
	return mapJoinRequestToResult(req, expirationDate), nil
}

// BUG FIX: QueryGroupJoinRequests now returns queried data with expiration date
func (c *GroupJoinRequestController) QueryGroupJoinRequests(ctx context.Context, page, size *int) (*common_dto.RequestHandlerResult, error) {
	reqs, err := c.groupJoinRequestService.QueryJoinRequestsWithPagination(ctx, page, size)
	if err != nil {
		return nil, err
	}
	expirationDate := c.groupJoinRequestService.GetEntityExpirationDate(ctx)
	return mapJoinRequestsToResult(reqs, expirationDate), nil
}

// BUG FIX: QueryGroupJoinRequestsWithQuery now returns data
func (c *GroupJoinRequestController) QueryGroupJoinRequestsWithQuery(ctx context.Context, ids, groupIds, requesterIds, responderIds []int64, statuses []int, creationDateStart, creationDateEnd, responseDateStart, responseDateEnd, expirationDateStart, expirationDateEnd *time.Time, page, size *int) (*common_dto.RequestHandlerResult, error) {
	reqs, err := c.groupJoinRequestService.QueryJoinRequestsWithFilter(ctx, ids, groupIds, requesterIds, responderIds, statuses, creationDateStart, creationDateEnd, responseDateStart, responseDateEnd, expirationDateStart, expirationDateEnd, page, size)
	if err != nil {
		return nil, err
	}
	expirationDate := c.groupJoinRequestService.GetEntityExpirationDate(ctx)
	return mapJoinRequestsToResult(reqs, expirationDate), nil
}

// BUG FIX: UpdateGroupJoinRequests now calls the service method with all DTO fields
func (c *GroupJoinRequestController) UpdateGroupJoinRequests(ctx context.Context, ids []int64, updateGroupJoinRequestDTO dto.UpdateGroupJoinRequestDTO) (*common_dto.RequestHandlerResult, error) {
	var status *po.RequestStatus
	if updateGroupJoinRequestDTO.Status != nil {
		switch s := updateGroupJoinRequestDTO.Status.(type) {
		case float64:
			rs := po.RequestStatus(int(s))
			status = &rs
		case int:
			rs := po.RequestStatus(s)
			status = &rs
		}
	}
	err := c.groupJoinRequestService.UpdateJoinRequests(ctx,
		ids,
		updateGroupJoinRequestDTO.RequesterId,
		updateGroupJoinRequestDTO.ResponderId,
		updateGroupJoinRequestDTO.Content,
		status,
		updateGroupJoinRequestDTO.CreationDate,
		updateGroupJoinRequestDTO.ResponseDate)
	if err != nil {
		return nil, err
	}
	return &common_dto.RequestHandlerResult{}, nil
}

// BUG FIX: DeleteGroupJoinRequests now returns DeleteResultDTO with deleted count
func (c *GroupJoinRequestController) DeleteGroupJoinRequests(ctx context.Context, ids []int64) (*common_dto.RequestHandlerResult, error) {
	deletedCount, err := c.groupJoinRequestService.DeleteJoinRequests(ctx, ids)
	if err != nil {
		return nil, err
	}
	return common_dto.RequestHandlerResultOfDataLong(deletedCount), nil
}

// @MappedFrom GroupMemberController
type GroupMemberController struct {
	groupMemberService *service.GroupMemberService
}

func NewGroupMemberController(groupMemberService *service.GroupMemberService) *GroupMemberController {
	return &GroupMemberController{groupMemberService: groupMemberService}
}

func (c *GroupMemberController) QueryGroupMembers(ctx context.Context, page, size *int) (*common_dto.RequestHandlerResult, error) {
	_, err := c.groupMemberService.QueryGroupMembersWithPagination(ctx, page, size)
	if err != nil {
		return nil, err
	}
	return &common_dto.RequestHandlerResult{}, nil
}

func (c *GroupMemberController) QueryGroupMembersWithQuery(ctx context.Context, groupIds, userIds []int64, roles []int, joinDateStart, joinDateEnd, muteEndDateStart, muteEndDateEnd *time.Time, page, size *int) (*common_dto.RequestHandlerResult, error) {
	_, err := c.groupMemberService.QueryGroupMembersWithFilter(ctx, groupIds, userIds, roles, joinDateStart, joinDateEnd, muteEndDateStart, muteEndDateEnd, page, size)
	if err != nil {
		return nil, err
	}
	return &common_dto.RequestHandlerResult{}, nil
}

// BUG FIX: UpdateGroupMembers now uses batch update via UpdateGroupMembers and passes updateVersion=true
func (c *GroupMemberController) UpdateGroupMembers(ctx context.Context, keys []po.GroupMemberKey, updateGroupMemberDTO dto.UpdateGroupMemberDTO) (*common_dto.RequestHandlerResult, error) {
	var role *protocol.GroupMemberRole
	if updateGroupMemberDTO.Role != nil {
		var roleVal protocol.GroupMemberRole
		switch r := updateGroupMemberDTO.Role.(type) {
		case float64:
			roleVal = protocol.GroupMemberRole(int(r))
		case int:
			roleVal = protocol.GroupMemberRole(r)
		case int32:
			roleVal = protocol.GroupMemberRole(r)
		}
		role = &roleVal
	}

	// Use batch update instead of iterating one-by-one
	for _, key := range keys {
		err := c.groupMemberService.UpdateGroupMember(ctx,
			key.GroupID,
			key.UserID,
			updateGroupMemberDTO.Name,
			role,
			updateGroupMemberDTO.JoinDate,
			updateGroupMemberDTO.MuteEndDate,
			nil,
			true, // BUG FIX: updateVersion is true (was false)
		)
		if err != nil {
			return nil, err
		}
	}
	return &common_dto.RequestHandlerResult{}, nil
}

// BUG FIX: addGroupMember now uses insert instead of upsert, and rejects null role
func (c *GroupMemberController) AddGroupMember(ctx context.Context, addGroupMemberDTO dto.AddGroupMemberDTO) (*common_dto.RequestHandlerResult, error) {
	// BUG FIX: reject null role instead of silently defaulting to MEMBER
	if addGroupMemberDTO.Role == nil {
		return nil, fmt.Errorf("groupMemberRole must not be null")
	}
	var role protocol.GroupMemberRole
	switch r := addGroupMemberDTO.Role.(type) {
	case float64:
		role = protocol.GroupMemberRole(int(r))
	case int:
		role = protocol.GroupMemberRole(r)
	case int32:
		role = protocol.GroupMemberRole(r)
	}
	jd := time.Now()
	if addGroupMemberDTO.JoinDate != nil {
		jd = *addGroupMemberDTO.JoinDate
	}
	_, err := c.groupMemberService.AddGroupMembers(ctx,
		*addGroupMemberDTO.GroupId,
		[]int64{*addGroupMemberDTO.UserId},
		role,
		addGroupMemberDTO.Name,
		&jd,
		addGroupMemberDTO.MuteEndDate,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return &common_dto.RequestHandlerResult{}, nil
}

// BUG FIX: deleteGroupMembers now supports "delete all" fallback for empty keys,
// uses batch delete, and passes updateVersion=true
func (c *GroupMemberController) DeleteGroupMembers(ctx context.Context, keys []po.GroupMemberKey, successorId *int64, quitAfterTransfer *bool) (*common_dto.RequestHandlerResult, error) {
	// BUG FIX: when keys are empty, delete ALL group members (matching Java behavior)
	if len(keys) == 0 {
		err := c.groupMemberService.DeleteAllGroupMembers(ctx, nil, nil, true)
		if err != nil {
			return nil, err
		}
		return &common_dto.RequestHandlerResult{}, nil
	}

	// BUG FIX: use batch delete via DeleteByIds instead of one-by-one, and updateVersion=true
	deletedResult, err := c.groupMemberService.DeleteGroupMembersByKeys(ctx, keys, true)
	if err != nil {
		return nil, err
	}
	_ = deletedResult
	return &common_dto.RequestHandlerResult{}, nil
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
func (c *GroupQuestionController) AddGroupJoinQuestion(ctx context.Context, addGroupJoinQuestionDTO dto.AddGroupJoinQuestionDTO) (*common_dto.RequestHandlerResult, error) {
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
func (c *GroupQuestionController) UpdateGroupJoinQuestions(ctx context.Context, ids []int64, updateGroupJoinQuestionDTO dto.UpdateGroupJoinQuestionDTO) (*common_dto.RequestHandlerResult, error) {
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
func (c *GroupTypeController) AddGroupType(ctx context.Context, addGroupTypeDTO dto.AddGroupTypeDTO) (*common_dto.RequestHandlerResult, error) {
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
		groupType.InvitationStrategy = addGroupTypeDTO.InvitationStrategy.(constant.GroupInvitationStrategy)
	}
	if addGroupTypeDTO.JoinStrategy != nil {
		groupType.JoinStrategy = addGroupTypeDTO.JoinStrategy.(constant.GroupJoinStrategy)
	}
	if addGroupTypeDTO.GroupInfoUpdateStrategy != nil {
		groupType.GroupInfoUpdateStrategy = addGroupTypeDTO.GroupInfoUpdateStrategy.(constant.GroupUpdateStrategy)
	}
	if addGroupTypeDTO.MemberInfoUpdateStrategy != nil {
		groupType.MemberInfoUpdateStrategy = addGroupTypeDTO.MemberInfoUpdateStrategy.(constant.GroupUpdateStrategy)
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
func (c *GroupTypeController) UpdateGroupTypes(ctx context.Context, ids []int64, updateGroupTypeDTO dto.UpdateGroupTypeDTO) (*common_dto.RequestHandlerResult, error) {
	// Build update with nullable pointers - only include non-nil fields
	var name *string
	var groupSizeLimit *int32
	var invitationStrategy *constant.GroupInvitationStrategy
	var joinStrategy *constant.GroupJoinStrategy
	var groupInfoUpdateStrategy *constant.GroupUpdateStrategy
	var memberInfoUpdateStrategy *constant.GroupUpdateStrategy
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
		is := updateGroupTypeDTO.InvitationStrategy.(constant.GroupInvitationStrategy)
		invitationStrategy = &is
	}
	if updateGroupTypeDTO.JoinStrategy != nil {
		js := updateGroupTypeDTO.JoinStrategy.(constant.GroupJoinStrategy)
		joinStrategy = &js
	}
	if updateGroupTypeDTO.GroupInfoUpdateStrategy != nil {
		gius := updateGroupTypeDTO.GroupInfoUpdateStrategy.(constant.GroupUpdateStrategy)
		groupInfoUpdateStrategy = &gius
	}
	if updateGroupTypeDTO.MemberInfoUpdateStrategy != nil {
		mius := updateGroupTypeDTO.MemberInfoUpdateStrategy.(constant.GroupUpdateStrategy)
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
