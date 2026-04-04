package controller

import (
	"context"
	"time"

	common_dto "im.turms/server/internal/domain/common/dto"
	"im.turms/server/internal/domain/group/access/admin/dto"
	"im.turms/server/internal/domain/group/constant"
	"im.turms/server/internal/domain/group/po"
	"im.turms/server/internal/domain/group/service"
	msg_service "im.turms/server/internal/domain/message/service"
	turmsmongo "im.turms/server/internal/storage/mongo"
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
	// Java doesn't use the exact struct for keys, it maps directly. We map keys.
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
	// For parity, we don't have session or deleteLogical in DeleteGroupsAndGroupMembers yet
	// Let's call DeleteGroupsAndGroupMembers without deleteLogical
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
