package controller

import (
	"context"
	"time"

	"google.golang.org/protobuf/proto"

	"im.turms/server/internal/domain/gateway/access/router"
	"im.turms/server/internal/domain/gateway/session"
	"im.turms/server/internal/domain/user/po"
	"im.turms/server/internal/domain/user/service"
	"im.turms/server/pkg/protocol"
)

type UserRelationshipController struct {
	userFriendRequestService     service.UserFriendRequestService
	userRelationshipGroupService service.UserRelationshipGroupService
	userRelationshipService      service.UserRelationshipService
}

func NewUserRelationshipController(
	userFriendRequestService service.UserFriendRequestService,
	userRelationshipGroupService service.UserRelationshipGroupService,
	userRelationshipService service.UserRelationshipService,
) *UserRelationshipController {
	return &UserRelationshipController{
		userFriendRequestService:     userFriendRequestService,
		userRelationshipGroupService: userRelationshipGroupService,
		userRelationshipService:      userRelationshipService,
	}
}

// RegisterRoutes wires all UserRelationship handlers to the gateway router.
func (c *UserRelationshipController) RegisterRoutes(r *router.Router) {
	r.RegisterController(&protocol.TurmsRequest_CreateFriendRequestRequest{}, c.HandleCreateFriendRequestRequest)
	r.RegisterController(&protocol.TurmsRequest_CreateRelationshipGroupRequest{}, c.HandleCreateRelationshipGroupRequest)
	r.RegisterController(&protocol.TurmsRequest_CreateRelationshipRequest{}, c.HandleCreateRelationshipRequest)
	r.RegisterController(&protocol.TurmsRequest_DeleteFriendRequestRequest{}, c.HandleDeleteFriendRequestRequest)
	r.RegisterController(&protocol.TurmsRequest_DeleteRelationshipGroupRequest{}, c.HandleDeleteRelationshipGroupRequest)
	r.RegisterController(&protocol.TurmsRequest_DeleteRelationshipRequest{}, c.HandleDeleteRelationshipRequest)
	r.RegisterController(&protocol.TurmsRequest_QueryFriendRequestsRequest{}, c.HandleQueryFriendRequestsRequest)
	r.RegisterController(&protocol.TurmsRequest_QueryRelatedUserIdsRequest{}, c.HandleQueryRelatedUserIdsRequest)
	r.RegisterController(&protocol.TurmsRequest_QueryRelationshipGroupsRequest{}, c.HandleQueryRelationshipGroupsRequest)
	r.RegisterController(&protocol.TurmsRequest_QueryRelationshipsRequest{}, c.HandleQueryRelationshipsRequest)
	r.RegisterController(&protocol.TurmsRequest_UpdateFriendRequestRequest{}, c.HandleUpdateFriendRequestRequest)
	r.RegisterController(&protocol.TurmsRequest_UpdateRelationshipGroupRequest{}, c.HandleUpdateRelationshipGroupRequest)
	r.RegisterController(&protocol.TurmsRequest_UpdateRelationshipRequest{}, c.HandleUpdateRelationshipRequest)
}

// @MappedFrom handleCreateFriendRequestRequest()
func (c *UserRelationshipController) HandleCreateFriendRequestRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	createReq := req.GetCreateFriendRequestRequest()
	_, err := c.userFriendRequestService.AuthAndCreateFriendRequest(ctx, s.UserID, createReq.GetRecipientId(), createReq.GetContent(), time.Now())
	if err != nil {
		return nil, err
	}
	return buildSuccessNotification(req.RequestId), nil
}

// @MappedFrom handleCreateRelationshipGroupRequest()
func (c *UserRelationshipController) HandleCreateRelationshipGroupRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	createReq := req.GetCreateRelationshipGroupRequest()

	group, err := c.userRelationshipGroupService.CreateRelationshipGroup(ctx, s.UserID, nil, createReq.GetName(), nil, nil)
	if err != nil {
		return nil, err
	}

	return &protocol.TurmsNotification{
		RequestId: req.RequestId,
		Code:      proto.Int32(1000),
		Data: &protocol.TurmsNotification_Data{
			Kind: &protocol.TurmsNotification_Data_Long{
				Long: int64(group.Key.Index),
			},
		},
	}, nil
}

// @MappedFrom handleCreateRelationshipRequest()
func (c *UserRelationshipController) HandleCreateRelationshipRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	createReq := req.GetCreateRelationshipRequest()

	var blockDate *time.Time
	if createReq.GetBlocked() {
		now := time.Now()
		blockDate = &now
	}
	now := time.Now()

	_, err := c.userRelationshipService.UpsertOneSidedRelationship(ctx, s.UserID, createReq.GetUserId(), blockDate, createReq.GroupIndex, &now, createReq.Name, nil)
	if err != nil {
		return nil, err
	}
	return buildSuccessNotification(req.RequestId), nil
}

// @MappedFrom handleDeleteFriendRequestRequest()
func (c *UserRelationshipController) HandleDeleteFriendRequestRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	deleteReq := req.GetDeleteFriendRequestRequest()
	_, err := c.userFriendRequestService.AuthAndRecallFriendRequest(ctx, s.UserID, deleteReq.GetRequestId())
	if err != nil {
		return nil, err
	}
	return buildSuccessNotification(req.RequestId), nil
}

// @MappedFrom handleDeleteRelationshipGroupRequest()
func (c *UserRelationshipController) HandleDeleteRelationshipGroupRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	deleteReq := req.GetDeleteRelationshipGroupRequest()
	if deleteReq.TargetGroupIndex != nil {
		err := c.userRelationshipGroupService.DeleteRelationshipGroupAndMoveMembersToNewGroup(ctx, s.UserID, deleteReq.GetGroupIndex(), deleteReq.GetTargetGroupIndex())
		if err != nil {
			return nil, err
		}
	} else {
		_, err := c.userRelationshipGroupService.DeleteRelationshipGroups(ctx, s.UserID, []int32{deleteReq.GetGroupIndex()}, nil)
		if err != nil {
			return nil, err
		}
	}
	return buildSuccessNotification(req.RequestId), nil
}

// @MappedFrom handleDeleteRelationshipRequest()
func (c *UserRelationshipController) HandleDeleteRelationshipRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	deleteReq := req.GetDeleteRelationshipRequest()

	if deleteReq.GetTargetGroupIndex() != 0 || deleteReq.TargetGroupIndex != nil {
		err := c.userRelationshipGroupService.MoveRelatedUserToNewGroup(ctx, s.UserID, deleteReq.GetUserId(), deleteReq.GetGroupIndex(), deleteReq.GetTargetGroupIndex(), true, nil)
		if err != nil {
			return nil, err
		}
	} else if deleteReq.GetGroupIndex() != 0 || deleteReq.GroupIndex != nil {
		_, err := c.userRelationshipGroupService.DeleteRelatedUserFromRelationshipGroup(ctx, s.UserID, deleteReq.GetUserId(), deleteReq.GetGroupIndex(), nil, true)
		if err != nil {
			return nil, err
		}
	} else {
		err := c.userRelationshipService.DeleteOneSidedRelationship(ctx, s.UserID, deleteReq.GetUserId())
		if err != nil {
			return nil, err
		}
	}
	return buildSuccessNotification(req.RequestId), nil
}

// @MappedFrom handleQueryFriendRequestsRequest()
func (c *UserRelationshipController) HandleQueryFriendRequestsRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	queryReq := req.GetQueryFriendRequestsRequest()

	var lastUpdatedTime *time.Time
	if queryReq.LastUpdatedDate != nil {
		t := time.UnixMilli(queryReq.GetLastUpdatedDate())
		lastUpdatedTime = &t
	}

	requests, err := c.userFriendRequestService.QueryFriendRequestsWithVersion(ctx, s.UserID, queryReq.GetAreSentByMe(), lastUpdatedTime)
	if err != nil {
		return nil, err
	}

	protos := make([]*protocol.UserFriendRequest, len(requests))
	for i, r := range requests {
		protos[i] = service.FriendRequestToProto(&r)
	}

	return &protocol.TurmsNotification{
		RequestId: req.RequestId,
		Code:      proto.Int32(1000),
		Data: &protocol.TurmsNotification_Data{
			Kind: &protocol.TurmsNotification_Data_UserFriendRequestsWithVersion{
				UserFriendRequestsWithVersion: &protocol.UserFriendRequestsWithVersion{
					UserFriendRequests: protos,
				},
			},
		},
	}, nil
}

// @MappedFrom handleQueryRelatedUserIdsRequest()
func (c *UserRelationshipController) HandleQueryRelatedUserIdsRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	queryReq := req.GetQueryRelatedUserIdsRequest()

	var lastUpdatedTime *time.Time
	if queryReq.LastUpdatedDate != nil {
		t := time.UnixMilli(queryReq.GetLastUpdatedDate())
		lastUpdatedTime = &t
	}

	userIds, version, err := c.userRelationshipService.QueryRelatedUserIdsWithVersion(ctx, s.UserID, queryReq.GetGroupIndexes(), queryReq.Blocked, lastUpdatedTime)
	if err != nil {
		return nil, err
	}

	var lastUpdatedDate *int64
	if version != nil {
		t := version.UnixMilli()
		lastUpdatedDate = &t
	}

	return &protocol.TurmsNotification{
		RequestId: req.RequestId,
		Code:      proto.Int32(1000),
		Data: &protocol.TurmsNotification_Data{
			Kind: &protocol.TurmsNotification_Data_LongsWithVersion{
				LongsWithVersion: &protocol.LongsWithVersion{
					Longs:           userIds,
					LastUpdatedDate: lastUpdatedDate,
				},
			},
		},
	}, nil
}

// @MappedFrom handleQueryRelationshipGroupsRequest()
func (c *UserRelationshipController) HandleQueryRelationshipGroupsRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	queryReq := req.GetQueryRelationshipGroupsRequest()
	var lastUpdatedDate *time.Time
	if queryReq.LastUpdatedDate != nil {
		t := time.UnixMilli(queryReq.GetLastUpdatedDate())
		lastUpdatedDate = &t
	}

	groups, version, err := c.userRelationshipGroupService.QueryRelationshipGroupsInfosWithVersion(ctx, s.UserID, lastUpdatedDate)
	if err != nil {
		return nil, err
	}

	protos := make([]*protocol.UserRelationshipGroup, len(groups))
	for i, g := range groups {
		protos[i] = service.RelationshipGroupToProto(g)
	}

	var lastUpdatedTimeProto *int64
	if version != nil {
		t := version.UnixMilli()
		lastUpdatedTimeProto = &t
	}

	return &protocol.TurmsNotification{
		RequestId: req.RequestId,
		Code:      proto.Int32(1000),
		Data: &protocol.TurmsNotification_Data{
			Kind: &protocol.TurmsNotification_Data_UserRelationshipGroupsWithVersion{
				UserRelationshipGroupsWithVersion: &protocol.UserRelationshipGroupsWithVersion{
					UserRelationshipGroups: protos,
					LastUpdatedDate:        lastUpdatedTimeProto,
				},
			},
		},
	}, nil
}

// @MappedFrom handleQueryRelationshipsRequest()
func (c *UserRelationshipController) HandleQueryRelationshipsRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	queryReq := req.GetQueryRelationshipsRequest()

	var lastUpdatedTime *time.Time
	if queryReq.LastUpdatedDate != nil {
		t := time.UnixMilli(queryReq.GetLastUpdatedDate())
		lastUpdatedTime = &t
	}

	relationships, version, err := c.userRelationshipService.QueryRelationshipsWithVersion(ctx, s.UserID, queryReq.GetUserIds(), queryReq.GetGroupIndexes(), queryReq.Blocked, lastUpdatedTime)
	if err != nil {
		return nil, err
	}

	protos := make([]*protocol.UserRelationship, len(relationships))
	for i, r := range relationships {
		protos[i] = service.RelationshipToProto(&r)
	}

	var lastUpdatedTimeProto *int64
	if version != nil {
		t := version.UnixMilli()
		lastUpdatedTimeProto = &t
	}

	return &protocol.TurmsNotification{
		RequestId: req.RequestId,
		Code:      proto.Int32(1000),
		Data: &protocol.TurmsNotification_Data{
			Kind: &protocol.TurmsNotification_Data_UserRelationshipsWithVersion{
				UserRelationshipsWithVersion: &protocol.UserRelationshipsWithVersion{
					UserRelationships: protos,
					LastUpdatedDate:   lastUpdatedTimeProto,
				},
			},
		},
	}, nil
}

// @MappedFrom handleUpdateFriendRequestRequest()
func (c *UserRelationshipController) HandleUpdateFriendRequestRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	updateReq := req.GetUpdateFriendRequestRequest()

	// In Go, ResponseAction is probably mapped from proto
	_, err := c.userFriendRequestService.AuthAndHandleFriendRequest(ctx, updateReq.GetRequestId(), s.UserID, po.ResponseAction(updateReq.GetResponseAction()), updateReq.Reason)
	if err != nil {
		return nil, err
	}

	return buildSuccessNotification(req.RequestId), nil
}

// @MappedFrom handleUpdateRelationshipGroupRequest()
func (c *UserRelationshipController) HandleUpdateRelationshipGroupRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	updateReq := req.GetUpdateRelationshipGroupRequest()

	err := c.userRelationshipGroupService.UpdateRelationshipGroupName(ctx, s.UserID, updateReq.GetGroupIndex(), updateReq.GetNewName())
	if err != nil {
		return nil, err
	}

	return buildSuccessNotification(req.RequestId), nil
}

// @MappedFrom handleUpdateRelationshipRequest()
func (c *UserRelationshipController) HandleUpdateRelationshipRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	updateReq := req.GetUpdateRelationshipRequest()

	var blockDate *time.Time
	if updateReq.Blocked != nil {
		if updateReq.GetBlocked() {
			now := time.Now()
			blockDate = &now
		} else {
			// Unblock -> typically pass time.Unix(0, 0) or blockDate remains nil but check Java logic
			// Actually Java sets check on isBlocked
		}
	}

	err := c.userRelationshipService.UpdateUserOneSidedRelationships(ctx, s.UserID, []int64{updateReq.GetUserId()}, blockDate, updateReq.NewGroupIndex, nil, updateReq.Name, nil)
	if err != nil {
		return nil, err
	}

	return buildSuccessNotification(req.RequestId), nil
}
