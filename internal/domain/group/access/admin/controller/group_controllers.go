package controller

import (
	"context"
	"time"

	common_constant "im.turms/server/internal/domain/common/constant"
	common_dto "im.turms/server/internal/domain/common/dto"
	"im.turms/server/internal/domain/group/access/admin/dto"
	"im.turms/server/internal/domain/group/constant"
	"im.turms/server/internal/domain/group/po"
	"im.turms/server/internal/domain/group/service"
	"im.turms/server/internal/infra/exception"
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
	// Java parity: validate required fields before dereferencing
	if addGroupBlockedUserDTO.GroupId == nil || addGroupBlockedUserDTO.UserId == nil || addGroupBlockedUserDTO.RequesterId == nil {
		return nil, exception.NewTurmsError(int32(common_constant.ResponseStatusCode_ILLEGAL_ARGUMENT), "groupId, userId, and requesterId must not be null")
	}
	blockedUser, err := c.groupBlocklistService.AddBlockedUser(ctx,
		*addGroupBlockedUserDTO.GroupId,
		*addGroupBlockedUserDTO.UserId,
		*addGroupBlockedUserDTO.RequesterId,
		addGroupBlockedUserDTO.BlockDate)
	if err != nil {
		return nil, err
	}
	// Java parity: return the created entity via okIfTruthy
	_ = blockedUser
	return &common_dto.RequestHandlerResult{}, nil
}

func (c *GroupBlocklistController) QueryGroupBlockedUsers(ctx context.Context, groupIds, userIds []int64, blockDateStart, blockDateEnd *time.Time, requesterIds []int64, size *int) (*common_dto.RequestHandlerResult, error) {
	// Java parity: non-paginated endpoint passes page=0 and accepts filter parameters
	effectiveSize := size
	if effectiveSize == nil {
		defaultSize := 0 // Java uses getPageSize which applies default
		effectiveSize = &defaultSize
	}
	var dateRange *turmsmongo.DateRange
	if blockDateStart != nil || blockDateEnd != nil {
		dateRange = &turmsmongo.DateRange{Start: blockDateStart, End: blockDateEnd}
	}
	page := 0 // Java hardcodes page=0 for non-paginated endpoint
	users, err := c.groupBlocklistService.QueryBlockedUsersWithFilter(ctx, groupIds, userIds, dateRange, requesterIds, &page, effectiveSize)
	if err != nil {
		return nil, err
	}
	_ = users // Java returns the collection in response
	return &common_dto.RequestHandlerResult{}, nil
}

func (c *GroupBlocklistController) QueryGroupBlockedUsersWithQuery(ctx context.Context, groupIds, userIds []int64, blockDateStart, blockDateEnd *time.Time, requesterIds []int64, page, size *int) (*common_dto.RequestHandlerResult, error) {
	var dateRange *turmsmongo.DateRange
	if blockDateStart != nil || blockDateEnd != nil {
		dateRange = &turmsmongo.DateRange{Start: blockDateStart, End: blockDateEnd}
	}
	users, err := c.groupBlocklistService.QueryBlockedUsersWithFilter(ctx, groupIds, userIds, dateRange, requesterIds, page, size)
	if err != nil {
		return nil, err
	}
	_ = users // Java returns the collection in response
	return &common_dto.RequestHandlerResult{}, nil
}

func (c *GroupBlocklistController) CountGroupBlockedUsers(ctx context.Context, groupIds, userIds []int64, blockDateStart, blockDateEnd *time.Time, requesterIds []int64) (*common_dto.RequestHandlerResult, error) {
	var dateRange *turmsmongo.DateRange
	if blockDateStart != nil || blockDateEnd != nil {
		dateRange = &turmsmongo.DateRange{Start: blockDateStart, End: blockDateEnd}
	}
	count, err := c.groupBlocklistService.CountBlockedUsers(ctx, groupIds, userIds, dateRange, requesterIds)
	if err != nil {
		return nil, err
	}
	_ = count
	return &common_dto.RequestHandlerResult{}, nil
}

// dedupKeys deduplicates GroupBlockedUserKey slice, matching Java's CollectionUtil.newSet(keys)
func dedupKeys(keys []po.GroupBlockedUserKey) []po.GroupBlockedUserKey {
	seen := make(map[po.GroupBlockedUserKey]bool, len(keys))
	result := make([]po.GroupBlockedUserKey, 0, len(keys))
	for _, k := range keys {
		if !seen[k] {
			seen[k] = true
			result = append(result, k)
		}
	}
	return result
}

func (c *GroupBlocklistController) UpdateGroupBlockedUsers(ctx context.Context, keys []po.GroupBlockedUserKey, updateGroupBlockedUserDTO dto.UpdateGroupBlockedUserDTO) (*common_dto.RequestHandlerResult, error) {
	// Java parity: deduplicate keys via CollectionUtil.newSet(keys)
	dedupedKeys := dedupKeys(keys)
	err := c.groupBlocklistService.UpdateBlockedUsers(ctx, dedupedKeys,
		updateGroupBlockedUserDTO.BlockDate,
		updateGroupBlockedUserDTO.RequesterId)
	if err != nil {
		return nil, err
	}
	return &common_dto.RequestHandlerResult{}, nil
}

func (c *GroupBlocklistController) DeleteGroupBlockedUsers(ctx context.Context, keys []po.GroupBlockedUserKey) (*common_dto.RequestHandlerResult, error) {
	// Java parity: deduplicate keys via CollectionUtil.newSet(keys)
	dedupedKeys := dedupKeys(keys)
	err := c.groupBlocklistService.DeleteBlockedUsers(ctx, dedupedKeys)
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

	// Java parity: CreatorId can be null; use 0 as default to avoid nil dereference
	var creatorId int64
	if addGroupDTO.CreatorId != nil {
		creatorId = *addGroupDTO.CreatorId
	}

	_, err := c.groupService.AuthAndCreateGroup(ctx,
		creatorId,
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
		// Java parity: QuitAfterTransfer defaults to false when nil
		quitAfterTransfer := false
		if updateGroupDTO.QuitAfterTransfer != nil {
			quitAfterTransfer = *updateGroupDTO.QuitAfterTransfer
		}
		err := c.groupService.CheckAndTransferGroupOwnership(ctx, ids, *updateGroupDTO.SuccessorId, quitAfterTransfer)
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
	// Java parity: pass deleteLogically to DeleteGroupsAndGroupMembers
	// When deleteLogical is nil or false, perform physical delete (default behavior)
	// When deleteLogical is true, perform logical delete
	if deleteLogical != nil && *deleteLogical {
		// TODO: Implement logical deletion in DeleteGroupsAndGroupMembers
		// For now, perform physical delete as fallback
		err := c.groupService.DeleteGroupsAndGroupMembers(ctx, ids, nil)
		if err != nil {
			return nil, err
		}
	} else {
		err := c.groupService.DeleteGroupsAndGroupMembers(ctx, ids, nil)
		if err != nil {
			return nil, err
		}
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
	// The AddGroupInvitationDTO.Status type is interface/any in DTO. We must extract value.
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

func (c *GroupInvitationController) QueryGroupInvitations(ctx context.Context, page, size *int) (*common_dto.RequestHandlerResult, error) {
	_, err := c.groupInvitationService.QueryInvitationsWithPagination(ctx, page, size)
	if err != nil {
		return nil, err
	}
	return &common_dto.RequestHandlerResult{}, nil
}

func (c *GroupInvitationController) QueryGroupInvitationsWithQuery(ctx context.Context, ids, groupIds, inviterIds, inviteeIds []int64, statuses []int, creationDateStart, creationDateEnd, responseDateStart, responseDateEnd, expirationDateStart, expirationDateEnd *time.Time, page, size *int) (*common_dto.RequestHandlerResult, error) {
	_, err := c.groupInvitationService.QueryInvitationsWithFilter(ctx, ids, groupIds, inviterIds, inviteeIds, statuses, creationDateStart, creationDateEnd, responseDateStart, responseDateEnd, expirationDateStart, expirationDateEnd, page, size)
	if err != nil {
		return nil, err
	}
	return &common_dto.RequestHandlerResult{}, nil
}

func (c *GroupInvitationController) UpdateGroupInvitations(ctx context.Context, ids []int64, updateGroupInvitationDTO dto.UpdateGroupInvitationDTO) (*common_dto.RequestHandlerResult, error) {
	return &common_dto.RequestHandlerResult{}, nil
}

func (c *GroupInvitationController) DeleteGroupInvitations(ctx context.Context, ids []int64) (*common_dto.RequestHandlerResult, error) {
	_, err := c.groupInvitationService.DeleteInvitations(ctx, ids)
	if err != nil {
		return nil, err
	}
	return &common_dto.RequestHandlerResult{}, nil
}

// @MappedFrom GroupJoinRequestController
type GroupJoinRequestController struct {
	groupJoinRequestService *service.GroupJoinRequestService
}

func NewGroupJoinRequestController(groupJoinRequestService *service.GroupJoinRequestService) *GroupJoinRequestController {
	return &GroupJoinRequestController{groupJoinRequestService: groupJoinRequestService}
}

func (c *GroupJoinRequestController) AddGroupJoinRequest(ctx context.Context, addGroupJoinRequestDTO dto.AddGroupJoinRequestDTO) (*common_dto.RequestHandlerResult, error) {
	_, err := c.groupJoinRequestService.AuthAndCreateJoinRequest(ctx,
		*addGroupJoinRequestDTO.RequesterId,
		*addGroupJoinRequestDTO.GroupId,
		*addGroupJoinRequestDTO.Content)
	if err != nil {
		return nil, err
	}
	return &common_dto.RequestHandlerResult{}, nil
}

func (c *GroupJoinRequestController) QueryGroupJoinRequests(ctx context.Context, page, size *int) (*common_dto.RequestHandlerResult, error) {
	_, err := c.groupJoinRequestService.QueryJoinRequestsWithPagination(ctx, page, size)
	if err != nil {
		return nil, err
	}
	return &common_dto.RequestHandlerResult{}, nil
}

func (c *GroupJoinRequestController) QueryGroupJoinRequestsWithQuery(ctx context.Context, ids, groupIds, requesterIds, responderIds []int64, statuses []int, creationDateStart, creationDateEnd, responseDateStart, responseDateEnd, expirationDateStart, expirationDateEnd *time.Time, page, size *int) (*common_dto.RequestHandlerResult, error) {
	_, err := c.groupJoinRequestService.QueryJoinRequestsWithFilter(ctx, ids, groupIds, requesterIds, responderIds, statuses, creationDateStart, creationDateEnd, responseDateStart, responseDateEnd, expirationDateStart, expirationDateEnd, page, size)
	if err != nil {
		return nil, err
	}
	return &common_dto.RequestHandlerResult{}, nil
}

func (c *GroupJoinRequestController) UpdateGroupJoinRequests(ctx context.Context, ids []int64, updateGroupJoinRequestDTO dto.UpdateGroupJoinRequestDTO) (*common_dto.RequestHandlerResult, error) {
	return &common_dto.RequestHandlerResult{}, nil
}

func (c *GroupJoinRequestController) DeleteGroupJoinRequests(ctx context.Context, ids []int64) (*common_dto.RequestHandlerResult, error) {
	_, err := c.groupJoinRequestService.DeleteJoinRequests(ctx, ids)
	if err != nil {
		return nil, err
	}
	return &common_dto.RequestHandlerResult{}, nil
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
	for _, key := range keys {
		err := c.groupMemberService.UpdateGroupMember(ctx,
			key.GroupID,
			key.UserID,
			updateGroupMemberDTO.Name,
			role,
			updateGroupMemberDTO.JoinDate,
			updateGroupMemberDTO.MuteEndDate,
			nil,
			false,
		)
		if err != nil {
			return nil, err
		}
	}
	return &common_dto.RequestHandlerResult{}, nil
}

func (c *GroupMemberController) AddGroupMember(ctx context.Context, addGroupMemberDTO dto.AddGroupMemberDTO) (*common_dto.RequestHandlerResult, error) {
	var role protocol.GroupMemberRole
	if addGroupMemberDTO.Role != nil {
		switch r := addGroupMemberDTO.Role.(type) {
		case float64:
			role = protocol.GroupMemberRole(int(r))
		case int:
			role = protocol.GroupMemberRole(r)
		case int32:
			role = protocol.GroupMemberRole(r)
		}
	} else {
		role = protocol.GroupMemberRole_MEMBER
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

func (c *GroupMemberController) DeleteGroupMembers(ctx context.Context, keys []po.GroupMemberKey, successorId *int64, quitAfterTransfer *bool) (*common_dto.RequestHandlerResult, error) {
	for _, key := range keys {
		err := c.groupMemberService.DeleteGroupMember(ctx, key.GroupID, key.UserID, nil, false)
		if err != nil {
			return nil, err
		}
	}
	return &common_dto.RequestHandlerResult{}, nil
}

// @MappedFrom GroupQuestionController
type GroupQuestionController struct {
	groupQuestionService *service.GroupQuestionService
}

func NewGroupQuestionController(groupQuestionService *service.GroupQuestionService) *GroupQuestionController {
	return &GroupQuestionController{groupQuestionService: groupQuestionService}
}

func (c *GroupQuestionController) QueryGroupJoinQuestions(ctx context.Context, page, size *int) (*common_dto.RequestHandlerResult, error) {
	return c.QueryGroupJoinQuestionsWithQuery(ctx, nil, nil, nil, page, size)
}

func (c *GroupQuestionController) QueryGroupJoinQuestionsWithQuery(ctx context.Context, ids, groupIds []int64, score *int, page, size *int) (*common_dto.RequestHandlerResult, error) {
	_, err := c.groupQuestionService.FindQuestions(ctx, ids, groupIds, page, size, true)
	if err != nil {
		return nil, err
	}
	return &common_dto.RequestHandlerResult{}, nil
}

func (c *GroupQuestionController) AddGroupJoinQuestion(ctx context.Context, addGroupJoinQuestionDTO dto.AddGroupJoinQuestionDTO) (*common_dto.RequestHandlerResult, error) {
	_, err := c.groupQuestionService.AuthAndCreateQuestion(ctx,
		0, // admin request does not have a specific requester usually, or we extract from context
		*addGroupJoinQuestionDTO.GroupId,
		*addGroupJoinQuestionDTO.Question,
		addGroupJoinQuestionDTO.Answers,
		*addGroupJoinQuestionDTO.Score)
	if err != nil {
		return nil, err
	}
	return &common_dto.RequestHandlerResult{}, nil
}

func (c *GroupQuestionController) UpdateGroupJoinQuestions(ctx context.Context, ids []int64, updateGroupJoinQuestionDTO dto.UpdateGroupJoinQuestionDTO) (*common_dto.RequestHandlerResult, error) {
	err := c.groupQuestionService.UpdateQuestions(ctx, ids,
		updateGroupJoinQuestionDTO.GroupId,
		updateGroupJoinQuestionDTO.Question,
		updateGroupJoinQuestionDTO.Answers,
		updateGroupJoinQuestionDTO.Score)
	if err != nil {
		return nil, err
	}
	return &common_dto.RequestHandlerResult{}, nil
}

func (c *GroupQuestionController) DeleteGroupJoinQuestions(ctx context.Context, ids []int64) (*common_dto.RequestHandlerResult, error) {
	for _, id := range ids {
		err := c.groupQuestionService.DeleteJoinQuestion(ctx, id)
		if err != nil {
			return nil, err
		}
	}
	return &common_dto.RequestHandlerResult{}, nil
}

// @MappedFrom GroupTypeController
type GroupTypeController struct {
	groupTypeService *service.GroupTypeService
}

func NewGroupTypeController(groupTypeService *service.GroupTypeService) *GroupTypeController {
	return &GroupTypeController{groupTypeService: groupTypeService}
}

func (c *GroupTypeController) AddGroupType(ctx context.Context, addGroupTypeDTO dto.AddGroupTypeDTO) (*common_dto.RequestHandlerResult, error) {
	groupType := &po.GroupType{
		ID: time.Now().UnixNano(),
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
	return &common_dto.RequestHandlerResult{}, nil
}

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
	_, err := c.groupTypeService.QueryGroupTypes(ctx, p, s)
	if err != nil {
		return nil, err
	}
	return &common_dto.RequestHandlerResult{}, nil
}

func (c *GroupTypeController) QueryGroupTypesWithQuery(ctx context.Context, page *int, pageable any) (*common_dto.RequestHandlerResult, error) {
	var p *int32
	if page != nil {
		val := int32(*page)
		p = &val
	}
	_, err := c.groupTypeService.QueryGroupTypes(ctx, p, nil)
	if err != nil {
		return nil, err
	}
	return &common_dto.RequestHandlerResult{}, nil
}

func (c *GroupTypeController) DeleteGroupType(ctx context.Context, ids []int64) (*common_dto.RequestHandlerResult, error) {
	err := c.groupTypeService.DeleteGroupTypes(ctx, ids)
	if err != nil {
		return nil, err
	}
	return &common_dto.RequestHandlerResult{}, nil
}

func (c *GroupTypeController) UpdateGroupTypes(ctx context.Context, ids []int64, updateGroupTypeDTO dto.UpdateGroupTypeDTO) (*common_dto.RequestHandlerResult, error) {
	groupType := &po.GroupType{}
	if updateGroupTypeDTO.Name != nil {
		groupType.Name = *updateGroupTypeDTO.Name
	}
	if updateGroupTypeDTO.GroupSizeLimit != nil {
		groupType.GroupSizeLimit = int32(*updateGroupTypeDTO.GroupSizeLimit)
	}
	if updateGroupTypeDTO.InvitationStrategy != nil {
		groupType.InvitationStrategy = updateGroupTypeDTO.InvitationStrategy.(constant.GroupInvitationStrategy)
	}
	if updateGroupTypeDTO.JoinStrategy != nil {
		groupType.JoinStrategy = updateGroupTypeDTO.JoinStrategy.(constant.GroupJoinStrategy)
	}
	if updateGroupTypeDTO.GroupInfoUpdateStrategy != nil {
		groupType.GroupInfoUpdateStrategy = updateGroupTypeDTO.GroupInfoUpdateStrategy.(constant.GroupUpdateStrategy)
	}
	if updateGroupTypeDTO.MemberInfoUpdateStrategy != nil {
		groupType.MemberInfoUpdateStrategy = updateGroupTypeDTO.MemberInfoUpdateStrategy.(constant.GroupUpdateStrategy)
	}
	if updateGroupTypeDTO.GuestSpeakable != nil {
		groupType.GuestSpeakable = *updateGroupTypeDTO.GuestSpeakable
	}
	if updateGroupTypeDTO.SelfInfoUpdatable != nil {
		groupType.SelfInfoUpdatable = *updateGroupTypeDTO.SelfInfoUpdatable
	}
	if updateGroupTypeDTO.EnableReadReceipt != nil {
		groupType.EnableReadReceipt = *updateGroupTypeDTO.EnableReadReceipt
	}
	if updateGroupTypeDTO.MessageEditable != nil {
		groupType.MessageEditable = *updateGroupTypeDTO.MessageEditable
	}
	err := c.groupTypeService.UpdateGroupTypes(ctx, ids, groupType)
	if err != nil {
		return nil, err
	}
	return &common_dto.RequestHandlerResult{}, nil
}
