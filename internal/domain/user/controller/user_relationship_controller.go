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

const defaultRelationshipGroupIndex int32 = 0

type UserRelationshipController struct {
	userFriendRequestService       service.UserFriendRequestService
	userRelationshipGroupService   service.UserRelationshipGroupService
	userRelationshipService        service.UserRelationshipService
	deleteTwoSidedRelationships    bool
	notifyFriendRequestRecipient   bool
	notifyRequesterOtherSessions   bool
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
	friendRequest, err := c.userFriendRequestService.AuthAndCreateFriendRequest(ctx, s.UserID, createReq.GetRecipientId(), createReq.GetContent(), time.Now())
	if err != nil {
		return nil, err
	}
	// Return the friend request ID as Long data (matching Java: RequestHandlerResult.ofDataLong(friendRequest.getId(), ...))
	return &protocol.TurmsNotification{
		RequestId: req.RequestId,
		Code:      proto.Int32(1000),
		Data: &protocol.TurmsNotification_Data{
			Kind: &protocol.TurmsNotification_Data_Long{
				Long: friendRequest.ID,
			},
		},
	}, nil
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

	// Default groupIndex to DEFAULT_RELATIONSHIP_GROUP_INDEX (0) when not set in the request
	groupIndex := int32(defaultRelationshipGroupIndex)
	if createReq.GroupIndex != nil {
		groupIndex = createReq.GetGroupIndex()
	}
	now := time.Now()

	_, err := c.userRelationshipService.UpsertOneSidedRelationship(ctx, s.UserID, createReq.GetUserId(), blockDate, &groupIndex, &now, createReq.Name, nil)
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

	// Java always calls deleteRelationshipGroupAndMoveMembersToNewGroup, defaulting targetGroupIndex
	// to DEFAULT_RELATIONSHIP_GROUP_INDEX (0) when not set in the request.
	targetGroupIndex := int32(defaultRelationshipGroupIndex)
	if deleteReq.TargetGroupIndex != nil {
		targetGroupIndex = deleteReq.GetTargetGroupIndex()
	}

	err := c.userRelationshipGroupService.DeleteRelationshipGroupAndMoveMembersToNewGroup(ctx, s.UserID, deleteReq.GetGroupIndex(), targetGroupIndex)
	if err != nil {
		return nil, err
	}

	return buildSuccessNotification(req.RequestId), nil
}

// @MappedFrom handleDeleteRelationshipRequest()
func (c *UserRelationshipController) HandleDeleteRelationshipRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	deleteReq := req.GetDeleteRelationshipRequest()

	if c.deleteTwoSidedRelationships {
		err := c.userRelationshipService.TryDeleteTwoSidedRelationships(ctx, s.UserID, deleteReq.GetUserId(), nil)
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

	// Convert empty groupIndexes list to nil (matching Java: groupIndexesList.isEmpty() ? null : groupIndexesList)
	groupIndexes := queryReq.GetGroupIndexes()
	if len(groupIndexes) == 0 {
		groupIndexes = nil
	}

	// Handle Blocked tri-state: nil when unset, false when explicitly false, true when explicitly true
	var isBlocked *bool
	if queryReq.Blocked != nil {
		isBlocked = queryReq.Blocked
	}

	userIds, version, err := c.userRelationshipService.QueryRelatedUserIdsWithVersion(ctx, s.UserID, groupIndexes, isBlocked, lastUpdatedTime)
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

	// Convert empty userIds list to nil (matching Java: userIdsList.isEmpty() ? null : userIdsList)
	userIDs := queryReq.GetUserIds()
	if len(userIDs) == 0 {
		userIDs = nil
	}

	// Convert empty groupIndexes list to nil
	groupIndexes := queryReq.GetGroupIndexes()
	if len(groupIndexes) == 0 {
		groupIndexes = nil
	}

	// Handle Blocked tri-state: nil when unset
	var isBlocked *bool
	if queryReq.Blocked != nil {
		isBlocked = queryReq.Blocked
	}

	relationships, version, err := c.userRelationshipService.QueryRelationshipsWithVersion(ctx, s.UserID, userIDs, groupIndexes, isBlocked, lastUpdatedTime)
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

	// Handle reason: nil when unset (matching Java: request.hasReason() ? request.getReason() : null)
	var reason *string
	if updateReq.Reason != nil {
		reason = updateReq.Reason
	}

	_, err := c.userFriendRequestService.AuthAndHandleFriendRequest(ctx, updateReq.GetRequestId(), s.UserID, po.ResponseAction(updateReq.GetResponseAction()), reason)
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
		}
	}

	// Match Java: use UpsertOneSidedRelationship instead of UpdateUserOneSidedRelationships
	// Java: upsertOneSidedRelationship(userId, request.getUserId(), name, blockDate, newGroupIndex, deleteGroupIndex, null, true, null)
	var newGroupIndex *int32
	if updateReq.NewGroupIndex != nil {
		newGroupIndex = updateReq.NewGroupIndex
	}

	var deleteGroupIndex *int32
	if updateReq.DeleteGroupIndex != nil {
		deleteGroupIndex = updateReq.DeleteGroupIndex
	}

	_, err := c.userRelationshipService.UpsertOneSidedRelationship(ctx, s.UserID, updateReq.GetUserId(), blockDate, newGroupIndex, nil, updateReq.Name, nil)
	if err != nil {
		return nil, err
	}

	// Handle deleteGroupIndex by moving to new group if specified
	if deleteGroupIndex != nil {
		err := c.userRelationshipGroupService.MoveRelatedUserToNewGroup(ctx, s.UserID, updateReq.GetUserId(), *deleteGroupIndex, defaultRelationshipGroupIndex, true, nil)
		if err != nil {
			return nil, err
		}
	}

	return buildSuccessNotification(req.RequestId), nil
}
