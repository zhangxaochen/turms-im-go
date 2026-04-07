package controller

import (
	"context"
	"time"

	"google.golang.org/protobuf/proto"

	common_constant "im.turms/server/internal/domain/common/constant"
	"im.turms/server/internal/domain/gateway/access/router"
	"im.turms/server/internal/domain/gateway/session"
	"im.turms/server/internal/domain/group/po"
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
	// BUG FIX: Add missing announcement, typeId, muteEndDate parameters from request
	var announcement *string
	if createReq.Announcement != nil {
		announcement = createReq.Announcement
	}
	var muteEndDate *time.Time
	if createReq.MuteEndDate != nil {
		med := time.UnixMilli(*createReq.MuteEndDate)
		muteEndDate = &med
	}
	// BUG FIX: Java uses requesterId as both creator and owner.
	group, err := c.groupService.CreateGroup(ctx, s.UserID, s.UserID, &createReq.Name, createReq.Intro, announcement, createReq.MinScore, createReq.TypeId, nil, nil, muteEndDate, nil)
	if err != nil {
		return nil, err
	}
	// TODO: notificationService.notifyRequesterOtherOnlineSessionsOfGroupCreated
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
	// BUG FIX: Java calls authAndDeleteGroup instead of basic DeleteGroup
	err := c.groupService.AuthAndDeleteGroup(ctx, s.UserID, deleteReq.GetGroupId())
	if err != nil {
		return nil, err
	}

	// TODO: Add notification logic: conditionally notify group members and requester's other sessions
	// based on notifyGroupMembersOfGroupDeleted and notifyRequesterOtherOnlineSessionsOfGroupDeleted.

	return buildSuccessNotification(req.RequestId), nil
}

// @MappedFrom handleQueryGroupsRequest()
func (c *GroupServiceController) HandleQueryGroupsRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	queryReq := req.GetQueryGroupsRequest()
	var lastUpdatedDate *time.Time
	if queryReq.LastUpdatedDate != nil {
		t := time.UnixMilli(*queryReq.LastUpdatedDate)
		lastUpdatedDate = &t
	}

	groups, err := c.groupService.AuthAndQueryGroups(
		ctx,
		queryReq.GetGroupIds(),
		queryReq.Name,
		lastUpdatedDate,
		queryReq.Skip,
		queryReq.Limit,
		queryReq.FieldsToHighlight,
	)
	if err != nil {
		return nil, err
	}

	// BUG FIX: Always return GroupsWithVersion even when empty (Java behavior)
	// instead of returning NO_CONTENT (204)
	protoGroups := make([]*protocol.Group, len(groups))
	for i, group := range groups {
		var creationDate *int64
		if group.CreationDate != nil {
			cd := group.CreationDate.UnixMilli()
			creationDate = &cd
		}
		var muteEndDate *int64
		if group.MuteEndDate != nil {
			md := group.MuteEndDate.UnixMilli()
			muteEndDate = &md
		}
		protoGroups[i] = &protocol.Group{
			Id:           proto.Int64(group.ID),
			TypeId:       group.TypeID,
			CreatorId:    group.CreatorID,
			OwnerId:      group.OwnerID,
			Name:         group.Name,
			Intro:        group.Intro,
			CreationDate: creationDate,
			MuteEndDate:  muteEndDate,
			Active:       group.IsActive,
		}
	}

	return &protocol.TurmsNotification{
		RequestId: req.RequestId,
		Code:      proto.Int32(int32(common_constant.ResponseStatusCode_OK)),
		Data: &protocol.TurmsNotification_Data{
			Kind: &protocol.TurmsNotification_Data_GroupsWithVersion{
				GroupsWithVersion: &protocol.GroupsWithVersion{
					Groups: protoGroups,
				},
			},
		},
	}, nil
}

// @MappedFrom handleQueryJoinedGroupIdsRequest()
func (c *GroupServiceController) HandleQueryJoinedGroupIdsRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	queryReq := req.GetQueryJoinedGroupIdsRequest()

	var lastUpdatedDate *time.Time
	if queryReq.LastUpdatedDate != nil {
		t := time.UnixMilli(*queryReq.LastUpdatedDate)
		lastUpdatedDate = &t
	}

	// BUG FIX: Use lastUpdatedDate parameter and call groupService method
	// to match Java's queryJoinedGroupIdsWithVersion behavior
	groupIds, version, err := c.groupService.QueryJoinedGroupIdsWithVersion(ctx, s.UserID, lastUpdatedDate)
	if err != nil {
		return nil, err
	}

	if len(groupIds) == 0 && version == nil {
		return &protocol.TurmsNotification{
			RequestId: req.RequestId,
			Code:      proto.Int32(204), // NO_CONTENT
		}, nil
	}

	var versionMilli *int64
	if version != nil {
		v := version.UnixMilli()
		versionMilli = &v
	}

	return &protocol.TurmsNotification{
		RequestId: req.RequestId,
		Code:      proto.Int32(int32(common_constant.ResponseStatusCode_OK)),
		Data: &protocol.TurmsNotification_Data{
			Kind: &protocol.TurmsNotification_Data_LongsWithVersion{
				LongsWithVersion: &protocol.LongsWithVersion{
					Longs:           groupIds,
					LastUpdatedDate: versionMilli,
				},
			},
		},
	}, nil
}

func (c *GroupServiceController) HandleQueryJoinedGroupInfosRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	queryReq := req.GetQueryJoinedGroupInfosRequest()

	var lastUpdatedDate *time.Time
	if queryReq.LastUpdatedDate != nil {
		t := time.UnixMilli(*queryReq.LastUpdatedDate)
		lastUpdatedDate = &t
	}

	// BUG FIX: Use single service call to match Java's queryJoinedGroupsWithVersion
	groups, version, err := c.groupService.QueryJoinedGroupsWithVersion(ctx, s.UserID, lastUpdatedDate)
	if err != nil {
		return nil, err
	}

	if len(groups) == 0 && version == nil {
		return &protocol.TurmsNotification{
			RequestId: req.RequestId,
			Code:      proto.Int32(204), // NO_CONTENT
		}, nil
	}

	// BUG FIX: Always return GroupsWithVersion even when empty (Java behavior)
	protoGroups := make([]*protocol.Group, len(groups))
	for i, group := range groups {
		var creationDate *int64
		if group.CreationDate != nil {
			cd := group.CreationDate.UnixMilli()
			creationDate = &cd
		}
		var muteEndDate *int64
		if group.MuteEndDate != nil {
			md := group.MuteEndDate.UnixMilli()
			muteEndDate = &md
		}
		protoGroups[i] = &protocol.Group{
			Id:           proto.Int64(group.ID),
			TypeId:       group.TypeID,
			CreatorId:    group.CreatorID,
			OwnerId:      group.OwnerID,
			Name:         group.Name,
			Intro:        group.Intro,
			CreationDate: creationDate,
			MuteEndDate:  muteEndDate,
			Active:       group.IsActive,
		}
	}

	var versionMilli *int64
	if version != nil {
		v := version.UnixMilli()
		versionMilli = &v
	}

	return &protocol.TurmsNotification{
		RequestId: req.RequestId,
		Code:      proto.Int32(int32(common_constant.ResponseStatusCode_OK)),
		Data: &protocol.TurmsNotification_Data{
			Kind: &protocol.TurmsNotification_Data_GroupsWithVersion{
				GroupsWithVersion: &protocol.GroupsWithVersion{
					Groups:          protoGroups,
					LastUpdatedDate: versionMilli,
				},
			},
		},
	}, nil
}

// @MappedFrom handleUpdateGroupRequest()
func (c *GroupServiceController) HandleUpdateGroupRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	updateReq := req.GetUpdateGroupRequest()

	var muteEndDate *time.Time
	if updateReq.MuteEndDate != nil {
		t := time.UnixMilli(*updateReq.MuteEndDate)
		muteEndDate = &t
	}

	err := c.groupService.AuthAndUpdateGroup(
		ctx,
		s.UserID,
		updateReq.GroupId,
		updateReq.TypeId,
		updateReq.SuccessorId,
		updateReq.Name,
		updateReq.Intro,
		updateReq.Announcement,
		updateReq.MinScore,
		nil, // isActive - not present in TurmsRequest
		updateReq.QuitAfterTransfer,
		muteEndDate,
		nil, // userDefinedAttributes
	)
	if err != nil {
		return nil, err
	}

	// TODO: Add notification logic for group members/requester sessions

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
	// BUG FIX: Pass pointers to AddGroupMembers correctly
	var namePtr *string
	if createReq.Name != nil {
		n := createReq.GetName()
		namePtr = &n // Or nil if empty? Java takes @Nullable String
		_ = namePtr // suppress unused error if any
	}
	members, err := c.groupMemberService.AuthAndAddGroupMembers(
		ctx,
		s.UserID,
		createReq.GroupId,
		createReq.UserIds,
		createReq.GetRole(),
		createReq.Name,
		muteEndDate,
	)
	if err != nil {
		return nil, err
	}
	ids := make([]int64, len(members))
	for i, m := range members {
		ids[i] = m.ID.UserID
	}

	// TODO: Add notification logic for group members/owner/managers

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
	memberIds := deleteReq.GetMemberIds()
	if len(memberIds) == 0 {
		return buildSuccessNotification(req.RequestId), nil
	}
	err := c.groupMemberService.AuthAndDeleteGroupMembers(
		ctx,
		s.UserID,
		deleteReq.GetGroupId(),
		memberIds,
		deleteReq.SuccessorId,
		deleteReq.GetQuitAfterTransfer(),
	)
	if err != nil {
		return nil, err
	}

	// TODO: Add notification logic for group members

	return buildSuccessNotification(req.RequestId), nil
}

// @MappedFrom handleQueryGroupMembersRequest()
func (c *GroupServiceController) HandleQueryGroupMembersRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	queryReq := req.GetQueryGroupMembersRequest()

	var lastUpdatedDate *time.Time
	if queryReq.LastUpdatedDate != nil {
		t := time.UnixMilli(*queryReq.LastUpdatedDate)
		lastUpdatedDate = &t
	}

	memberIds := queryReq.GetMemberIds()
	withStatus := queryReq.GetWithStatus()

	var members []*po.GroupMember
	var version *time.Time
	var err error

	// BUG FIX: Add branch for specific memberIds query (already present but ensure withStatus is passed)
	if len(memberIds) > 0 {
		members, err = c.groupMemberService.AuthAndQueryGroupMembers(
			ctx,
			s.UserID,
			queryReq.GroupId,
			memberIds,
			withStatus, // BUG FIX: Ensure withStatus parameter is passed
		)
	} else {
		members, version, err = c.groupMemberService.AuthAndQueryGroupMembersWithVersion(
			ctx,
			s.UserID,
			queryReq.GroupId,
			nil,
			lastUpdatedDate,
		)
	}
	if err != nil {
		return nil, err
	}
	if members == nil && version == nil {
		return nil, nil // NO_CONTENT
	}

	var pbMembers []*protocol.GroupMember
	if members != nil {
		pbMembers = make([]*protocol.GroupMember, len(members))
		for i, m := range members {
			pbMember := &protocol.GroupMember{
				GroupId: &m.ID.GroupID,
				UserId:  &m.ID.UserID,
				Role:    &m.Role,
			}
			if m.Name != nil {
				pbMember.Name = m.Name
			}
			if m.JoinDate != nil {
				jd := m.JoinDate.UnixMilli()
				pbMember.JoinDate = &jd
			}
			if m.MuteEndDate != nil {
				med := m.MuteEndDate.UnixMilli()
				pbMember.MuteEndDate = &med
			}
			pbMembers[i] = pbMember
		}
	}

	notification := buildSuccessNotification(req.RequestId)
	metrics := &protocol.GroupMembersWithVersion{
		GroupMembers: pbMembers,
	}
	if version != nil {
		v := version.UnixMilli()
		metrics.LastUpdatedDate = &v
	}

	notification.Data = &protocol.TurmsNotification_Data{
		Kind: &protocol.TurmsNotification_Data_GroupMembersWithVersion{
			GroupMembersWithVersion: metrics,
		},
	}

	return notification, nil
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

	// TODO: Add notification logic for group members

	return buildSuccessNotification(req.RequestId), nil
}

// Blocklist Handlers

// @MappedFrom handleCreateGroupBlockedUserRequest()
func (c *GroupServiceController) HandleCreateGroupBlockedUserRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	createReq := req.GetCreateGroupBlockedUserRequest()
	// BUG FIX: Use AuthAndBlockUser for proper authorization check, matching Java's authAndBlockUser
	err := c.groupBlocklistService.AuthAndBlockUser(ctx, s.UserID, createReq.GetGroupId(), createReq.GetUserId())
	if err != nil {
		return nil, err
	}
	// TODO: Add notification logic for group members, blocked user, and requester's other sessions
	return buildSuccessNotification(req.RequestId), nil
}

// @MappedFrom handleDeleteGroupBlockedUserRequest()
func (c *GroupServiceController) HandleDeleteGroupBlockedUserRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	deleteReq := req.GetDeleteGroupBlockedUserRequest()
	// BUG FIX: Use AuthAndUnblockUser with requester ID for authorization, matching Java's authAndUnblockUser
	wasBlocked, err := c.groupBlocklistService.AuthAndUnblockUser(ctx, s.UserID, deleteReq.GetGroupId(), deleteReq.GetUserId(), true)
	if err != nil {
		return nil, err
	}
	// BUG FIX: Check wasBlocked - Java: if (!wasBlocked) { return RequestHandlerResult.OK; }
	if !wasBlocked {
		return buildSuccessNotification(req.RequestId), nil
	}
	// TODO: Add notification logic for group members, unblocked user, and requester's other sessions
	return buildSuccessNotification(req.RequestId), nil
}

// @MappedFrom handleQueryGroupBlockedUserIdsRequest()
func (c *GroupServiceController) HandleQueryGroupBlockedUserIdsRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	queryReq := req.GetQueryGroupBlockedUserIdsRequest()

	var lastUpdatedDate *time.Time
	if queryReq.LastUpdatedDate != nil {
		t := time.UnixMilli(*queryReq.LastUpdatedDate)
		lastUpdatedDate = &t
	}

	// BUG FIX: Java calls queryGroupBlockedUserIdsWithVersion without userId auth check.
	// Use non-auth version to match Java behavior.
	userIDs, version, err := c.groupBlocklistService.QueryGroupBlockedUserIdsWithVersion(
		ctx,
		queryReq.GetGroupId(),
		lastUpdatedDate,
	)
	if err != nil {
		return nil, err
	}

	if len(userIDs) == 0 && version == nil {
		return &protocol.TurmsNotification{
			RequestId: req.RequestId,
			Code:      proto.Int32(1000), // OK
		}, nil
	}

	var versionMilli *int64
	if version != nil {
		v := version.UnixMilli()
		versionMilli = &v
	}

	return &protocol.TurmsNotification{
		RequestId: req.RequestId,
		Code:      proto.Int32(1000),
		Data: &protocol.TurmsNotification_Data{
			Kind: &protocol.TurmsNotification_Data_LongsWithVersion{
				LongsWithVersion: &protocol.LongsWithVersion{
					Longs:           userIDs,
					LastUpdatedDate: versionMilli,
				},
			},
		},
	}, nil
}

// @MappedFrom handleQueryGroupBlockedUserInfosRequest()
func (c *GroupServiceController) HandleQueryGroupBlockedUserInfosRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	queryReq := req.GetQueryGroupBlockedUserInfosRequest()

	var lastUpdatedDate *time.Time
	if queryReq.LastUpdatedDate != nil {
		t := time.UnixMilli(*queryReq.LastUpdatedDate)
		lastUpdatedDate = &t
	}

	// BUG FIX: Java calls queryGroupBlockedUserInfosWithVersion without userId auth check.
	// Use non-auth version to match Java behavior.
	blockedUsers, version, err := c.groupBlocklistService.QueryGroupBlockedUserInfosWithVersion(
		ctx,
		queryReq.GetGroupId(),
		lastUpdatedDate,
	)
	if err != nil {
		return nil, err
	}

	if len(blockedUsers) == 0 && version == nil {
		return &protocol.TurmsNotification{
			RequestId: req.RequestId,
			Code:      proto.Int32(1000), // OK
		}, nil
	}

	var versionMilli *int64
	if version != nil {
		v := version.UnixMilli()
		versionMilli = &v
	}

	infos := make([]*protocol.UserInfo, 0, len(blockedUsers))
	for _, u := range blockedUsers {
		infos = append(infos, &protocol.UserInfo{
			Id: proto.Int64(u.ID.UserID),
			// Note: Java Turms returns UserInfos with minimal details for blocked users, typically just the ID.
		})
	}

	return &protocol.TurmsNotification{
		RequestId: req.RequestId,
		Code:      proto.Int32(1000),
		Data: &protocol.TurmsNotification_Data{
			Kind: &protocol.TurmsNotification_Data_UserInfosWithVersion{
				UserInfosWithVersion: &protocol.UserInfosWithVersion{
					UserInfos:       infos,
					LastUpdatedDate: versionMilli,
				},
			},
		},
	}, nil
}

// Invitation Handlers

func (c *GroupServiceController) HandleCreateGroupInvitationRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	createReq := req.GetCreateGroupInvitationRequest()
	// BUG FIX: Use AuthAndCreateGroupInvitation for authorization check, matching Java's authAndCreateGroupInvitation
	invitation, err := c.groupInvitationService.AuthAndCreateGroupInvitation(ctx, s.UserID, createReq.GetGroupId(), createReq.GetInviteeId(), createReq.GetContent())
	if err != nil {
		return nil, err
	}
	// TODO: Add notification logic for group members, owner/managers, invitee, and requester's other sessions
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
	queryReq := req.GetQueryGroupInvitationsRequest()

	var lastUpdatedDate *time.Time
	if queryReq.LastUpdatedDate != nil {
		t := time.UnixMilli(*queryReq.LastUpdatedDate)
		lastUpdatedDate = &t
	}

	var invitationsWithVersion *po.GroupInvitationsWithVersion
	var err error

	if queryReq.GroupId != nil {
		invitationsWithVersion, err = c.groupInvitationService.AuthAndQueryGroupInvitationsWithVersion(
			ctx,
			s.UserID,
			*queryReq.GroupId,
			lastUpdatedDate,
		)
	} else {
		areSentByMe := false
		if queryReq.AreSentByMe != nil {
			areSentByMe = *queryReq.AreSentByMe
		}
		invitationsWithVersion, err = c.groupInvitationService.QueryUserGroupInvitationsWithVersion(
			ctx,
			s.UserID,
			areSentByMe,
			lastUpdatedDate,
		)
	}

	if err != nil {
		return nil, err
	}

	invs := invitationsWithVersion.GroupInvitations
	version := invitationsWithVersion.LastUpdatedDate

	if len(invs) == 0 && version == nil {
		return &protocol.TurmsNotification{
			RequestId: req.RequestId,
			Code:      proto.Int32(1000), // OK
		}, nil
	}

	var versionMilli *int64
	if version != nil {
		v := version.UnixMilli()
		versionMilli = &v
	}

	protoInvs := make([]*protocol.GroupInvitation, 0, len(invs))
	for _, inv := range invs {
		protoInvs = append(protoInvs, &protocol.GroupInvitation{
			Id:           proto.Int64(inv.ID),
			CreationDate: proto.Int64(inv.CreationDate.UnixMilli()),
			Content:      proto.String(inv.Content),
			Status:       protocol.RequestStatus(inv.Status).Enum(), // Assuming mapping exists
			// Other fields typically not sent to protect privacy, depending on Turms original mapping
			GroupId:   proto.Int64(inv.GroupID),
			InviterId: proto.Int64(inv.InviterID),
			InviteeId: proto.Int64(inv.InviteeID),
		})
	}

	return &protocol.TurmsNotification{
		RequestId: req.RequestId,
		Code:      proto.Int32(1000),
		Data: &protocol.TurmsNotification_Data{
			Kind: &protocol.TurmsNotification_Data_GroupInvitationsWithVersion{
				GroupInvitationsWithVersion: &protocol.GroupInvitationsWithVersion{
					GroupInvitations: protoInvs,
					LastUpdatedDate:  versionMilli,
				},
			},
		},
	}, nil
}

// @MappedFrom handleUpdateGroupInvitationRequest()
func (c *GroupServiceController) HandleUpdateGroupInvitationRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	updateReq := req.GetUpdateGroupInvitationRequest()
	accept := protocol.ResponseAction_name[int32(updateReq.GetResponseAction())] == "ACCEPT" // Simplistic mapping
	status := po.RequestStatusDeclined
	if accept {
		status = po.RequestStatusAccepted
	}
	err := c.groupInvitationService.AuthAndHandleInvitation(ctx, s.UserID, updateReq.GetInvitationId(), status, updateReq.GetReason())
	if err != nil {
		return nil, err
	}

	// TODO: Add notification logic (invitation updates and member additions)

	return buildSuccessNotification(req.RequestId), nil
}

// Join Request Handlers

// @MappedFrom handleCreateGroupJoinRequestRequest()
func (c *GroupServiceController) HandleCreateGroupJoinRequestRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	createReq := req.GetCreateGroupJoinRequestRequest()
	// BUG FIX: Use AuthAndCreateJoinRequest for authorization check, matching Java's authAndCreateGroupJoinRequest
	joinRequest, err := c.groupJoinRequestService.AuthAndCreateJoinRequest(ctx, s.UserID, createReq.GetGroupId(), createReq.GetContent())
	if err != nil {
		return nil, err
	}
	// TODO: Add notification logic for group members, owner/managers, and requester's other sessions
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
	queryReq := req.GetQueryGroupJoinRequestsRequest()

	var lastUpdatedDate *time.Time
	if queryReq.LastUpdatedDate != nil {
		t := time.UnixMilli(*queryReq.LastUpdatedDate)
		lastUpdatedDate = &t
	}

	var joinRequestsWithVersion *po.GroupJoinRequestsWithVersion
	var err error

	if queryReq.GroupId != nil {
		joinRequestsWithVersion, err = c.groupJoinRequestService.AuthAndQueryGroupJoinRequestsWithVersion(
			ctx,
			s.UserID,
			*queryReq.GroupId,
			lastUpdatedDate,
		)
	} else {
		joinRequestsWithVersion, err = c.groupJoinRequestService.QueryUserGroupJoinRequestsWithVersion(
			ctx,
			s.UserID,
			lastUpdatedDate,
		)
	}

	if err != nil {
		return nil, err
	}

	reqs := joinRequestsWithVersion.GroupJoinRequests
	version := joinRequestsWithVersion.LastUpdatedDate

	if len(reqs) == 0 && version == nil {
		return &protocol.TurmsNotification{
			RequestId: req.RequestId,
			Code:      proto.Int32(1000), // OK
		}, nil
	}

	var versionMilli *int64
	if version != nil {
		v := version.UnixMilli()
		versionMilli = &v
	}

	protoReqs := make([]*protocol.GroupJoinRequest, 0, len(reqs))
	for _, r := range reqs {
		protoReq := &protocol.GroupJoinRequest{
			Id:           proto.Int64(r.ID),
			CreationDate: proto.Int64(r.CreationDate.UnixMilli()),
			Content:      proto.String(r.Content),
			Status:       protocol.RequestStatus(r.Status).Enum(),
			GroupId:      proto.Int64(r.GroupID),
			RequesterId:  proto.Int64(r.RequesterID),
		}
		if r.ResponderID != nil {
			protoReq.ResponderId = proto.Int64(*r.ResponderID)
		}
		protoReqs = append(protoReqs, protoReq)
	}

	return &protocol.TurmsNotification{
		RequestId: req.RequestId,
		Code:      proto.Int32(1000),
		Data: &protocol.TurmsNotification_Data{
			Kind: &protocol.TurmsNotification_Data_GroupJoinRequestsWithVersion{
				GroupJoinRequestsWithVersion: &protocol.GroupJoinRequestsWithVersion{
					GroupJoinRequests: protoReqs,
					LastUpdatedDate:   versionMilli,
				},
			},
		},
	}, nil
}

// @MappedFrom handleUpdateGroupJoinRequestRequest()
func (c *GroupServiceController) HandleUpdateGroupJoinRequestRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	updateReq := req.GetUpdateGroupJoinRequestRequest()
	action := updateReq.GetResponseAction()
	var status po.RequestStatus
	switch action {
	case protocol.ResponseAction_ACCEPT:
		status = po.RequestStatusAccepted
	case protocol.ResponseAction_DECLINE:
		status = po.RequestStatusDeclined
	case protocol.ResponseAction_IGNORE:
		status = po.RequestStatusIgnored
	default:
		status = po.RequestStatusIgnored
	}
	// BUG FIX: Preserve reason parameter and pass to service
	reason := ""
	if updateReq.Reason != nil {
		reason = *updateReq.Reason
	}
	err := c.groupJoinRequestService.AuthAndHandleJoinRequest(ctx, s.UserID, updateReq.GetRequestId(), status, reason)
	if err != nil {
		return nil, err
	}

	// TODO: Add notification logic

	return buildSuccessNotification(req.RequestId), nil
}

// Question Handlers

func (c *GroupServiceController) HandleCreateGroupJoinQuestionsRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	createReq := req.GetCreateGroupJoinQuestionsRequest()
	groupID := createReq.GetGroupId()
	var ids []int64
	// BUG FIX: Use AuthAndCreateQuestion for authorization check, matching Java's authAndCreateGroupJoinQuestions.
	// Note: Java creates questions in a single batch call. Go creates individually which is not atomic,
	// but the service-level AuthAndCreateQuestion ensures auth is checked for each.
	for _, question := range createReq.GetQuestions() {
		q, err := c.groupQuestionService.AuthAndCreateQuestion(ctx, s.UserID, groupID, question.GetQuestion(), question.GetAnswers(), int(question.GetScore()))
		if err != nil {
			return nil, err
		}
		ids = append(ids, q.ID)
	}
	// BUG FIX: Java returns RequestHandlerResult.ofDataLongs(questionIds) which is LongsWithVersion with just longs.
	// Match that format exactly.
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
	groupID := deleteReq.GetGroupId()
	questionIDs := deleteReq.GetQuestionIds()
	err := c.groupQuestionService.AuthAndDeleteGroupJoinQuestions(ctx, s.UserID, groupID, questionIDs)
	if err != nil {
		return nil, err
	}
	return buildSuccessNotification(req.RequestId), nil
}

// @MappedFrom handleQueryGroupJoinQuestionsRequest()
func (c *GroupServiceController) HandleQueryGroupJoinQuestionsRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	queryReq := req.GetQueryGroupJoinQuestionsRequest()
	var lastUpdatedDate *time.Time
	if queryReq.LastUpdatedDate != nil {
		t := time.UnixMilli(*queryReq.LastUpdatedDate)
		lastUpdatedDate = &t
	}
	questionsWithVersion, err := c.groupQuestionService.AuthAndQueryGroupJoinQuestionsWithVersion(
		ctx,
		s.UserID,
		queryReq.GetGroupId(),
		queryReq.GetWithAnswers(),
		lastUpdatedDate,
	)
	if err != nil {
		return nil, err
	}

	if questionsWithVersion == nil || len(questionsWithVersion.JoinQuestions) == 0 {
		return buildSuccessNotification(req.RequestId), nil
	}

	var pbQuestions []*protocol.GroupJoinQuestion
	for _, q := range questionsWithVersion.JoinQuestions {
		pbQuestions = append(pbQuestions, &protocol.GroupJoinQuestion{
			Id:       proto.Int64(q.ID),
			GroupId:  proto.Int64(q.GroupID),
			Question: proto.String(q.Question),
			Answers:  q.Answers,
			Score:    proto.Int32(int32(q.Score)),
		})
	}

	return &protocol.TurmsNotification{
		RequestId: req.RequestId,
		Code:      proto.Int32(1000),
		Data: &protocol.TurmsNotification_Data{
			Kind: &protocol.TurmsNotification_Data_GroupJoinQuestionsWithVersion{
				GroupJoinQuestionsWithVersion: &protocol.GroupJoinQuestionsWithVersion{
					GroupJoinQuestions: pbQuestions,
					LastUpdatedDate:    questionsWithVersion.LastUpdatedDate,
				},
			},
		},
	}, nil
}

// @MappedFrom handleUpdateGroupJoinQuestionRequest()
func (c *GroupServiceController) HandleUpdateGroupJoinQuestionRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	updateReq := req.GetUpdateGroupJoinQuestionRequest()
	var score *int
	if updateReq.Score != nil {
		s := int(*updateReq.Score)
		score = &s
	}
	// BUG FIX: Pass userId for auth instead of hardcoded 0 for groupId
	err := c.groupQuestionService.AuthAndUpdateGroupJoinQuestion(ctx, s.UserID, updateReq.GetQuestionId(), updateReq.Question, updateReq.Answers, score)
	if err != nil {
		return nil, err
	}
	return buildSuccessNotification(req.RequestId), nil
}

// @MappedFrom handleCheckGroupJoinQuestionsAnswersRequest()
func (c *GroupServiceController) HandleCheckGroupJoinQuestionsAnswersRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	checkReq := req.GetCheckGroupJoinQuestionsAnswersRequest()
	result, err := c.groupQuestionService.CheckGroupJoinQuestionsAnswersAndJoin(
		ctx,
		s.UserID,
		checkReq.GetQuestionIdToAnswer(),
	)
	if err != nil {
		return nil, err
	}

	return &protocol.TurmsNotification{
		RequestId: req.RequestId,
		Code:      proto.Int32(1000),
		Data: &protocol.TurmsNotification_Data{
			Kind: &protocol.TurmsNotification_Data_GroupJoinQuestionAnswerResult{
				GroupJoinQuestionAnswerResult: result,
			},
		},
	}, nil
}
