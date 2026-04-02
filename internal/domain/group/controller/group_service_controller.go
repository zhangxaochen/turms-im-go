package controller

import (
	"context"
	"time"

	"google.golang.org/protobuf/proto"

	common_constant "im.turms/server/internal/domain/common/constant"
	"im.turms/server/internal/domain/gateway/access/router"
	"im.turms/server/internal/domain/gateway/session"
	"im.turms/server/internal/domain/group/service"
	"im.turms/server/pkg/protocol"
)

type GroupServiceController struct {
	groupService            *service.GroupService
	groupMemberService      *service.GroupMemberService
	groupBlocklistService   *service.GroupBlocklistService
	groupInvitationService  *service.GroupInvitationService
	groupJoinRequestService *service.GroupJoinRequestService
	groupQuestionService    *service.GroupQuestionService
}

func NewGroupServiceController(
	groupService *service.GroupService,
	groupMemberService *service.GroupMemberService,
	groupBlocklistService *service.GroupBlocklistService,
	groupInvitationService *service.GroupInvitationService,
	groupJoinRequestService *service.GroupJoinRequestService,
	groupQuestionService *service.GroupQuestionService,
) *GroupServiceController {
	return &GroupServiceController{
		groupService:            groupService,
		groupMemberService:      groupMemberService,
		groupBlocklistService:   groupBlocklistService,
		groupInvitationService:  groupInvitationService,
		groupJoinRequestService: groupJoinRequestService,
		groupQuestionService:    groupQuestionService,
	}
}

func (c *GroupServiceController) RegisterRoutes(r *router.Router) {
	// Group
	r.RegisterController(&protocol.TurmsRequest_CreateGroupRequest{}, c.HandleCreateGroupRequest)
	r.RegisterController(&protocol.TurmsRequest_DeleteGroupRequest{}, c.HandleDeleteGroupRequest)
	r.RegisterController(&protocol.TurmsRequest_QueryGroupsRequest{}, c.HandleQueryGroupsRequest)
	r.RegisterController(&protocol.TurmsRequest_QueryJoinedGroupIdsRequest{}, c.HandleQueryJoinedGroupIdsRequest)
	r.RegisterController(&protocol.TurmsRequest_QueryJoinedGroupInfosRequest{}, c.HandleQueryJoinedGroupInfosRequest)
	r.RegisterController(&protocol.TurmsRequest_UpdateGroupRequest{}, c.HandleUpdateGroupRequest)

	// Member
	r.RegisterController(&protocol.TurmsRequest_CreateGroupMembersRequest{}, c.HandleCreateGroupMembersRequest)
	r.RegisterController(&protocol.TurmsRequest_DeleteGroupMembersRequest{}, c.HandleDeleteGroupMembersRequest)
	r.RegisterController(&protocol.TurmsRequest_QueryGroupMembersRequest{}, c.HandleQueryGroupMembersRequest)
	r.RegisterController(&protocol.TurmsRequest_UpdateGroupMemberRequest{}, c.HandleUpdateGroupMemberRequest)

	// Blocklist
	r.RegisterController(&protocol.TurmsRequest_CreateGroupBlockedUserRequest{}, c.HandleCreateGroupBlockedUserRequest)
	r.RegisterController(&protocol.TurmsRequest_DeleteGroupBlockedUserRequest{}, c.HandleDeleteGroupBlockedUserRequest)
	r.RegisterController(&protocol.TurmsRequest_QueryGroupBlockedUserIdsRequest{}, c.HandleQueryGroupBlockedUserIdsRequest)
	r.RegisterController(&protocol.TurmsRequest_QueryGroupBlockedUserInfosRequest{}, c.HandleQueryGroupBlockedUserInfosRequest)

	// Invitation
	r.RegisterController(&protocol.TurmsRequest_CreateGroupInvitationRequest{}, c.HandleCreateGroupInvitationRequest)
	r.RegisterController(&protocol.TurmsRequest_DeleteGroupInvitationRequest{}, c.HandleDeleteGroupInvitationRequest)
	r.RegisterController(&protocol.TurmsRequest_QueryGroupInvitationsRequest{}, c.HandleQueryGroupInvitationsRequest)
	r.RegisterController(&protocol.TurmsRequest_UpdateGroupInvitationRequest{}, c.HandleUpdateGroupInvitationRequest)

	// Join Request
	r.RegisterController(&protocol.TurmsRequest_CreateGroupJoinRequestRequest{}, c.HandleCreateGroupJoinRequestRequest)
	r.RegisterController(&protocol.TurmsRequest_DeleteGroupJoinRequestRequest{}, c.HandleDeleteGroupJoinRequestRequest)
	r.RegisterController(&protocol.TurmsRequest_QueryGroupJoinRequestsRequest{}, c.HandleQueryGroupJoinRequestsRequest)
	r.RegisterController(&protocol.TurmsRequest_UpdateGroupJoinRequestRequest{}, c.HandleUpdateGroupJoinRequestRequest)

	// Question
	r.RegisterController(&protocol.TurmsRequest_CreateGroupJoinQuestionsRequest{}, c.HandleCreateGroupJoinQuestionsRequest)
	r.RegisterController(&protocol.TurmsRequest_DeleteGroupJoinQuestionsRequest{}, c.HandleDeleteGroupJoinQuestionsRequest)
	r.RegisterController(&protocol.TurmsRequest_QueryGroupJoinQuestionsRequest{}, c.HandleQueryGroupJoinQuestionsRequest)
	r.RegisterController(&protocol.TurmsRequest_UpdateGroupJoinQuestionRequest{}, c.HandleUpdateGroupJoinQuestionRequest)
	r.RegisterController(&protocol.TurmsRequest_CheckGroupJoinQuestionsAnswersRequest{}, c.HandleCheckGroupJoinQuestionsAnswersRequest)
}

func buildSuccessNotification(reqID *int64) *protocol.TurmsNotification {
	return &protocol.TurmsNotification{
		RequestId: reqID,
		Code:      proto.Int32(1000), // SUCCESS
	}
}

// Group Handlers

// @MappedFrom handleCreateGroupRequest()
func (c *GroupServiceController) HandleCreateGroupRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	createReq := req.GetCreateGroupRequest()
	group, err := c.groupService.CreateGroup(ctx, s.UserID, 0, &createReq.Name, createReq.Intro, createReq.MinScore)
	if err != nil {
		return nil, err
	}
	return &protocol.TurmsNotification{
		RequestId: req.RequestId,
		Code:      proto.Int32(1000),
		Data: &protocol.TurmsNotification_Data{
			Kind: &protocol.TurmsNotification_Data_Long{
				Long: group.ID,
			},
		},
	}, nil
}

// @MappedFrom handleDeleteGroupRequest()
func (c *GroupServiceController) HandleDeleteGroupRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	deleteReq := req.GetDeleteGroupRequest()
	err := c.groupService.DeleteGroup(ctx, s.UserID, deleteReq.GetGroupId())
	if err != nil {
		return nil, err
	}
	return buildSuccessNotification(req.RequestId), nil
}

// @MappedFrom handleQueryGroupsRequest()
func (c *GroupServiceController) HandleQueryGroupsRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	// TODO: Implement QueryGroupsRequest with versioning
	return buildSuccessNotification(req.RequestId), nil
}

// @MappedFrom handleQueryJoinedGroupIdsRequest()
func (c *GroupServiceController) HandleQueryJoinedGroupIdsRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	// TODO: Implement QueryJoinedGroupIdsRequest
	return buildSuccessNotification(req.RequestId), nil
}

func (c *GroupServiceController) HandleQueryJoinedGroupInfosRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	// TODO: Implement QueryJoinedGroupInfosRequest
	return buildSuccessNotification(req.RequestId), nil
}

// @MappedFrom handleUpdateGroupRequest()
func (c *GroupServiceController) HandleUpdateGroupRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	// TODO: Implement UpdateGroupRequest in service and controller
	return buildSuccessNotification(req.RequestId), nil
}

// Member Handlers

// @MappedFrom handleCreateGroupMembersRequest()
func (c *GroupServiceController) HandleCreateGroupMembersRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	createReq := req.GetCreateGroupMembersRequest()
	var muteEndDate *time.Time
	if createReq.MuteEndDate != nil {
		t := time.UnixMilli(*createReq.MuteEndDate)
		muteEndDate = &t
	}
	members, err := c.groupMemberService.AuthAndAddGroupMembers(
		ctx,
		s.UserID,
		createReq.GetGroupId(),
		createReq.GetUserIds(),
		createReq.GetRole(),
		muteEndDate,
	)
	if err != nil {
		return nil, err
	}
	ids := make([]int64, len(members))
	for i, m := range members {
		ids[i] = m.ID.UserID
	}
	return &protocol.TurmsNotification{
		RequestId: req.RequestId,
		Code:      proto.Int32(int32(common_constant.ResponseStatusCode_OK)),
		Data: &protocol.TurmsNotification_Data{
			Kind: &protocol.TurmsNotification_Data_LongsWithVersion{
				LongsWithVersion: &protocol.LongsWithVersion{
					Longs: ids,
				},
			},
		},
	}, nil
}

// @MappedFrom handleDeleteGroupMembersRequest()
func (c *GroupServiceController) HandleDeleteGroupMembersRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	deleteReq := req.GetDeleteGroupMembersRequest()
	err := c.groupMemberService.AuthAndDeleteGroupMembers(
		ctx,
		s.UserID,
		deleteReq.GetGroupId(),
		deleteReq.GetMemberIds(),
		deleteReq.SuccessorId,
		deleteReq.GetQuitAfterTransfer(),
	)
	if err != nil {
		return nil, err
	}
	return buildSuccessNotification(req.RequestId), nil
}

// @MappedFrom handleQueryGroupMembersRequest()
func (c *GroupServiceController) HandleQueryGroupMembersRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	// TODO: Implement QueryGroupMembersRequest
	return buildSuccessNotification(req.RequestId), nil
}

// @MappedFrom handleUpdateGroupMemberRequest()
func (c *GroupServiceController) HandleUpdateGroupMemberRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	updateReq := req.GetUpdateGroupMemberRequest()
	var muteEndDate *time.Time
	if updateReq.MuteEndDate != nil {
		t := time.UnixMilli(*updateReq.MuteEndDate)
		muteEndDate = &t
	}
	var role *protocol.GroupMemberRole
	if updateReq.Role != nil {
		r := updateReq.GetRole()
		role = &r
	}
	err := c.groupMemberService.AuthAndUpdateGroupMember(
		ctx,
		s.UserID,
		updateReq.GetGroupId(),
		updateReq.GetMemberId(),
		updateReq.Name,
		role,
		muteEndDate,
	)
	if err != nil {
		return nil, err
	}
	return buildSuccessNotification(req.RequestId), nil
}

// Blocklist Handlers

// @MappedFrom handleCreateGroupBlockedUserRequest()
func (c *GroupServiceController) HandleCreateGroupBlockedUserRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	createReq := req.GetCreateGroupBlockedUserRequest()
	err := c.groupBlocklistService.BlockUser(ctx, createReq.GetGroupId(), createReq.GetUserId(), s.UserID)
	if err != nil {
		return nil, err
	}
	return buildSuccessNotification(req.RequestId), nil
}

// @MappedFrom handleDeleteGroupBlockedUserRequest()
func (c *GroupServiceController) HandleDeleteGroupBlockedUserRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	deleteReq := req.GetDeleteGroupBlockedUserRequest()
	err := c.groupBlocklistService.UnblockUser(ctx, deleteReq.GetGroupId(), deleteReq.GetUserId())
	if err != nil {
		return nil, err
	}
	return buildSuccessNotification(req.RequestId), nil
}

// @MappedFrom handleQueryGroupBlockedUserIdsRequest()
func (c *GroupServiceController) HandleQueryGroupBlockedUserIdsRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	// TODO: Implement QueryGroupBlockedUserIdsRequest
	return buildSuccessNotification(req.RequestId), nil
}

func (c *GroupServiceController) HandleQueryGroupBlockedUserInfosRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	// TODO: Implement QueryGroupBlockedUserInfosRequest
	return buildSuccessNotification(req.RequestId), nil
}

// Invitation Handlers

func (c *GroupServiceController) HandleCreateGroupInvitationRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	createReq := req.GetCreateGroupInvitationRequest()
	invitation, err := c.groupInvitationService.CreateInvitation(ctx, createReq.GetGroupId(), s.UserID, createReq.GetInviteeId(), createReq.GetContent())
	if err != nil {
		return nil, err
	}
	return &protocol.TurmsNotification{
		RequestId: req.RequestId,
		Code:      proto.Int32(1000),
		Data: &protocol.TurmsNotification_Data{
			Kind: &protocol.TurmsNotification_Data_Long{
				Long: invitation.ID,
			},
		},
	}, nil
}

// @MappedFrom handleDeleteGroupInvitationRequest()
func (c *GroupServiceController) HandleDeleteGroupInvitationRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	deleteReq := req.GetDeleteGroupInvitationRequest()
	_, err := c.groupInvitationService.RecallPendingInvitation(ctx, deleteReq.GetInvitationId(), s.UserID)
	if err != nil {
		return nil, err
	}
	return buildSuccessNotification(req.RequestId), nil
}

// @MappedFrom handleQueryGroupInvitationsRequest()
func (c *GroupServiceController) HandleQueryGroupInvitationsRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	// TODO: Implement QueryGroupInvitationsRequest
	return buildSuccessNotification(req.RequestId), nil
}

// @MappedFrom handleUpdateGroupInvitationRequest()
func (c *GroupServiceController) HandleUpdateGroupInvitationRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	updateReq := req.GetUpdateGroupInvitationRequest()
	accept := protocol.ResponseAction_name[int32(updateReq.GetResponseAction())] == "ACCEPT" // Simplistic mapping
	_, err := c.groupInvitationService.ReplyToInvitation(ctx, updateReq.GetInvitationId(), s.UserID, accept)
	if err != nil {
		return nil, err
	}
	return buildSuccessNotification(req.RequestId), nil
}

// Join Request Handlers

// @MappedFrom handleCreateGroupJoinRequestRequest()
func (c *GroupServiceController) HandleCreateGroupJoinRequestRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	createReq := req.GetCreateGroupJoinRequestRequest()
	joinRequest, err := c.groupJoinRequestService.CreateJoinRequest(ctx, createReq.GetGroupId(), s.UserID, createReq.GetContent())
	if err != nil {
		return nil, err
	}
	return &protocol.TurmsNotification{
		RequestId: req.RequestId,
		Code:      proto.Int32(1000),
		Data: &protocol.TurmsNotification_Data{
			Kind: &protocol.TurmsNotification_Data_Long{
				Long: joinRequest.ID,
			},
		},
	}, nil
}

// @MappedFrom handleDeleteGroupJoinRequestRequest()
func (c *GroupServiceController) HandleDeleteGroupJoinRequestRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	deleteReq := req.GetDeleteGroupJoinRequestRequest()
	_, err := c.groupJoinRequestService.RecallPendingJoinRequest(ctx, deleteReq.GetRequestId(), s.UserID)
	if err != nil {
		return nil, err
	}
	return buildSuccessNotification(req.RequestId), nil
}

// @MappedFrom handleQueryGroupJoinRequestsRequest()
func (c *GroupServiceController) HandleQueryGroupJoinRequestsRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	// TODO: Implement QueryGroupJoinRequestsRequest
	return buildSuccessNotification(req.RequestId), nil
}

// @MappedFrom handleUpdateGroupJoinRequestRequest()
func (c *GroupServiceController) HandleUpdateGroupJoinRequestRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	updateReq := req.GetUpdateGroupJoinRequestRequest()
	accept := protocol.ResponseAction_name[int32(updateReq.GetResponseAction())] == "ACCEPT"
	_, err := c.groupJoinRequestService.ReplyToJoinRequest(ctx, updateReq.GetRequestId(), s.UserID, accept)
	if err != nil {
		return nil, err
	}
	return buildSuccessNotification(req.RequestId), nil
}

// Question Handlers

func (c *GroupServiceController) HandleCreateGroupJoinQuestionsRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	createReq := req.GetCreateGroupJoinQuestionsRequest()
	var ids []int64
	for _, question := range createReq.GetQuestions() {
		q, err := c.groupQuestionService.CreateJoinQuestion(ctx, createReq.GetGroupId(), question.GetQuestion(), question.GetAnswers(), int(question.GetScore()))
		if err != nil {
			return nil, err
		}
		ids = append(ids, q.ID)
	}
	return &protocol.TurmsNotification{
		RequestId: req.RequestId,
		Code:      proto.Int32(1000),
		Data: &protocol.TurmsNotification_Data{
			Kind: &protocol.TurmsNotification_Data_LongsWithVersion{
				LongsWithVersion: &protocol.LongsWithVersion{
					Longs: ids,
				},
			},
		},
	}, nil
}

// @MappedFrom handleDeleteGroupJoinQuestionsRequest()
func (c *GroupServiceController) HandleDeleteGroupJoinQuestionsRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	deleteReq := req.GetDeleteGroupJoinQuestionsRequest()
	for _, questionID := range deleteReq.GetQuestionIds() {
		err := c.groupQuestionService.DeleteJoinQuestion(ctx, questionID)
		if err != nil {
			return nil, err
		}
	}
	return buildSuccessNotification(req.RequestId), nil
}

// @MappedFrom handleQueryGroupJoinQuestionsRequest()
func (c *GroupServiceController) HandleQueryGroupJoinQuestionsRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	// TODO: Implement QueryGroupJoinQuestionsRequest
	return buildSuccessNotification(req.RequestId), nil
}

// @MappedFrom handleUpdateGroupJoinQuestionRequest()
func (c *GroupServiceController) HandleUpdateGroupJoinQuestionRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	updateReq := req.GetUpdateGroupJoinQuestionRequest()
	var score *int
	if updateReq.Score != nil {
		s := int(*updateReq.Score)
		score = &s
	}
	err := c.groupQuestionService.UpdateJoinQuestion(ctx, updateReq.GetQuestionId(), 0, updateReq.Question, updateReq.Answers, score)
	if err != nil {
		return nil, err
	}
	return buildSuccessNotification(req.RequestId), nil
}

func (c *GroupServiceController) HandleCheckGroupJoinQuestionsAnswersRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	// TODO: Implement CheckGroupJoinQuestionsAnswersRequest
	return buildSuccessNotification(req.RequestId), nil
}
