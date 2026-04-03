package controller

import "im.turms/server/internal/domain/common/dto"

// @MappedFrom GroupBlocklistController
type GroupBlocklistController struct{}

func (c *GroupBlocklistController) AddGroupBlockedUser(addGroupBlockedUserDTO any) *dto.RequestHandlerResult {
	return nil
}
func (c *GroupBlocklistController) QueryGroupBlockedUsers(page, size *int) *dto.RequestHandlerResult {
	return nil
}
func (c *GroupBlocklistController) QueryGroupBlockedUsersWithQuery(groupIds, userIds []int64, blockDateStart, blockDateEnd *int64, requesterIds []int64, page, size *int) *dto.RequestHandlerResult {
	return nil
}
func (c *GroupBlocklistController) UpdateGroupBlockedUsers(keys []any, updateGroupBlockedUserDTO any) *dto.RequestHandlerResult {
	return nil
}
func (c *GroupBlocklistController) DeleteGroupBlockedUsers(keys []any) *dto.RequestHandlerResult {
	return nil
}

// @MappedFrom GroupController
type GroupController struct{}

func (c *GroupController) AddGroup(addGroupDTO any) *dto.RequestHandlerResult    { return nil }
func (c *GroupController) QueryGroups(page, size *int) *dto.RequestHandlerResult { return nil }
func (c *GroupController) QueryGroupsWithQuery(ids, typeIds, creatorIds, ownerIds []int64, isActive *bool, creationDateStart, creationDateEnd, deletionDateStart, deletionDateEnd, muteEndDateStart, muteEndDateEnd *int64, memberIds []int64, page, size *int) *dto.RequestHandlerResult {
	return nil
}
func (c *GroupController) UpdateGroups(ids []int64, updateGroupDTO any) *dto.RequestHandlerResult {
	return nil
}
func (c *GroupController) DeleteGroups(ids []int64, deleteLogical *bool) *dto.RequestHandlerResult {
	return nil
}

// @MappedFrom GroupInvitationController
type GroupInvitationController struct{}

func (c *GroupInvitationController) AddGroupInvitation(addGroupInvitationDTO any) *dto.RequestHandlerResult {
	return nil
}
func (c *GroupInvitationController) QueryGroupInvitations(page, size *int) *dto.RequestHandlerResult {
	return nil
}
func (c *GroupInvitationController) QueryGroupInvitationsWithQuery(ids, groupIds, inviterIds, inviteeIds []int64, statuses []int, creationDateStart, creationDateEnd, responseDateStart, responseDateEnd, expirationDateStart, expirationDateEnd *int64, page, size *int) *dto.RequestHandlerResult {
	return nil
}
func (c *GroupInvitationController) UpdateGroupInvitations(ids []int64, updateGroupInvitationDTO any) *dto.RequestHandlerResult {
	return nil
}
func (c *GroupInvitationController) DeleteGroupInvitations(ids []int64) *dto.RequestHandlerResult {
	return nil
}

// @MappedFrom GroupJoinRequestController
type GroupJoinRequestController struct{}

func (c *GroupJoinRequestController) AddGroupJoinRequest(addGroupJoinRequestDTO any) *dto.RequestHandlerResult {
	return nil
}
func (c *GroupJoinRequestController) QueryGroupJoinRequests(page, size *int) *dto.RequestHandlerResult {
	return nil
}
func (c *GroupJoinRequestController) QueryGroupJoinRequestsWithQuery(ids, groupIds, requesterIds, responderIds []int64, statuses []int, creationDateStart, creationDateEnd, responseDateStart, responseDateEnd, expirationDateStart, expirationDateEnd *int64, page, size *int) *dto.RequestHandlerResult {
	return nil
}
func (c *GroupJoinRequestController) UpdateGroupJoinRequests(ids []int64, updateGroupJoinRequestDTO any) *dto.RequestHandlerResult {
	return nil
}
func (c *GroupJoinRequestController) DeleteGroupJoinRequests(ids []int64) *dto.RequestHandlerResult {
	return nil
}

// @MappedFrom GroupMemberController
type GroupMemberController struct{}

func (c *GroupMemberController) QueryGroupMembers(page, size *int) *dto.RequestHandlerResult {
	return nil
}
func (c *GroupMemberController) QueryGroupMembersWithQuery(groupIds, userIds []int64, roles []int, joinDateStart, joinDateEnd, muteEndDateStart, muteEndDateEnd *int64, page, size *int) *dto.RequestHandlerResult {
	return nil
}
func (c *GroupMemberController) DeleteGroupMembers(keys []any, successorId *int64, quitAfterTransfer *bool) *dto.RequestHandlerResult {
	return nil
}

// @MappedFrom GroupQuestionController
type GroupQuestionController struct{}

func (c *GroupQuestionController) QueryGroupJoinQuestions(page, size *int) *dto.RequestHandlerResult {
	return nil
}
func (c *GroupQuestionController) QueryGroupJoinQuestionsWithQuery(ids, groupIds []int64, score *int, page, size *int) *dto.RequestHandlerResult {
	return nil
}
func (c *GroupQuestionController) AddGroupJoinQuestion(addGroupJoinQuestionDTO any) *dto.RequestHandlerResult {
	return nil
}
func (c *GroupQuestionController) UpdateGroupJoinQuestions(ids []int64, updateGroupJoinQuestionDTO any) *dto.RequestHandlerResult {
	return nil
}
func (c *GroupQuestionController) DeleteGroupJoinQuestions(ids []int64) *dto.RequestHandlerResult {
	return nil
}

// @MappedFrom GroupTypeController
type GroupTypeController struct{}

func (c *GroupTypeController) AddGroupType(addGroupTypeDTO any) *dto.RequestHandlerResult { return nil }
func (c *GroupTypeController) QueryGroupTypes(page, size *int) *dto.RequestHandlerResult  { return nil }
func (c *GroupTypeController) QueryGroupTypesWithQuery(page *int, pageable any) *dto.RequestHandlerResult {
	return nil
}                                                                                    // mapped parameter differently in Java
func (c *GroupTypeController) DeleteGroupType(ids []int64) *dto.RequestHandlerResult { return nil }
