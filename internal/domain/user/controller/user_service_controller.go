package controller

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/protobuf/proto"

	"im.turms/server/internal/domain/gateway/access/router"
	"im.turms/server/internal/domain/gateway/session"
	"im.turms/server/internal/domain/user/service"
	"im.turms/server/pkg/protocol"
)

type UserServiceController struct {
	userService service.UserService
}

func NewUserServiceController(
	userService service.UserService,
) *UserServiceController {
	return &UserServiceController{
		userService: userService,
	}
}

// RegisterRoutes wires all UserService handlers to the gateway router.
func (c *UserServiceController) RegisterRoutes(r *router.Router) {
	r.RegisterController(&protocol.TurmsRequest_QueryUserProfilesRequest{}, c.HandleQueryUserProfilesRequest)
	r.RegisterController(&protocol.TurmsRequest_QueryNearbyUsersRequest{}, c.HandleQueryNearbyUsersRequest)
	r.RegisterController(&protocol.TurmsRequest_QueryUserOnlineStatusesRequest{}, c.HandleQueryUserOnlineStatusesRequest)
	r.RegisterController(&protocol.TurmsRequest_UpdateUserLocationRequest{}, c.HandleUpdateUserLocationRequest)
	r.RegisterController(&protocol.TurmsRequest_UpdateUserOnlineStatusRequest{}, c.HandleUpdateUserOnlineStatusRequest)
	r.RegisterController(&protocol.TurmsRequest_UpdateUserRequest{}, c.HandleUpdateUserRequest)
}

// HandleQueryUserProfilesRequest queries user profiles by user IDs.
func (c *UserServiceController) HandleQueryUserProfilesRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	queryReq := req.GetQueryUserProfilesRequest()

	userIDs := queryReq.GetUserIds()

	var lastUpdatedDate *time.Time
	if queryReq.LastUpdatedDate != nil {
		t := time.UnixMilli(queryReq.GetLastUpdatedDate())
		lastUpdatedDate = &t
	}

	users, err := c.userService.AuthAndQueryUsersProfile(ctx, s.UserID, userIDs, "", lastUpdatedDate, 0, 0)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return buildSuccessNotification(req.RequestId), nil
	}

	// Build UserInfosWithVersion
	userInfos := make([]*protocol.UserInfo, 0, len(users))
	for _, u := range users {
		info := &protocol.UserInfo{
			Id:   proto.Int64(u.ID),
			Name: proto.String(u.Name),
		}
		if u.Intro != "" {
			info.Intro = proto.String(u.Intro)
		}
		if u.ProfilePicture != "" {
			info.ProfilePicture = proto.String(u.ProfilePicture)
		}
		info.RegistrationDate = proto.Int64(u.RegistrationDate.UnixMilli())
		info.Active = proto.Bool(u.IsActive)
		userInfos = append(userInfos, info)
	}

	return &protocol.TurmsNotification{
		RequestId: req.RequestId,
		Code:      proto.Int32(1000),
		Data: &protocol.TurmsNotification_Data{
			Kind: &protocol.TurmsNotification_Data_UserInfosWithVersion{
				UserInfosWithVersion: &protocol.UserInfosWithVersion{
					UserInfos: userInfos,
				},
			},
		},
	}, nil
}

// HandleQueryNearbyUsersRequest queries nearby users based on location.
// NOTE: This requires SessionLocationService/NearbyUserService which are not yet ported.
// Returns empty result for now.
func (c *UserServiceController) HandleQueryNearbyUsersRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	// TODO: Implement when NearbyUserService is ported
	_ = req.GetQueryNearbyUsersRequest()
	return buildSuccessNotification(req.RequestId), nil
}

// HandleQueryUserOnlineStatusesRequest queries online statuses for a set of user IDs.
// NOTE: This requires UserStatusService which is not yet ported.
// Returns empty result for now.
func (c *UserServiceController) HandleQueryUserOnlineStatusesRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	// TODO: Implement when UserStatusService is ported
	_ = req.GetQueryUserOnlineStatusesRequest()
	return buildSuccessNotification(req.RequestId), nil
}

// HandleUpdateUserLocationRequest updates the user's current location.
// NOTE: This requires SessionLocationService which is not yet ported.
// Returns OK for now.
func (c *UserServiceController) HandleUpdateUserLocationRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	// TODO: Implement when SessionLocationService is ported
	_ = req.GetUpdateUserLocationRequest()
	return buildSuccessNotification(req.RequestId), nil
}

// HandleUpdateUserOnlineStatusRequest updates the user's online status (invisible, busy, etc.).
// NOTE: This requires UserStatusService/SessionService which are not yet fully ported.
// Returns OK for now.
func (c *UserServiceController) HandleUpdateUserOnlineStatusRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	// TODO: Implement when UserStatusService is ported
	_ = req.GetUpdateUserOnlineStatusRequest()
	return buildSuccessNotification(req.RequestId), nil
}

// HandleUpdateUserRequest updates user profile fields (name, intro, profilePicture, etc.).
func (c *UserServiceController) HandleUpdateUserRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	updateReq := req.GetUpdateUserRequest()

	update := bson.M{}
	if updateReq.Password != nil {
		update["pw"] = updateReq.GetPassword()
	}
	if updateReq.Name != nil {
		update["n"] = updateReq.GetName()
	}
	if updateReq.Intro != nil {
		update["intro"] = updateReq.GetIntro()
	}
	if updateReq.ProfilePicture != nil {
		update["pp"] = updateReq.GetProfilePicture()
	}
	if updateReq.ProfileAccessStrategy != nil && updateReq.GetProfileAccessStrategy() != protocol.ProfileAccessStrategy_ALL {
		update["pas"] = int32(updateReq.GetProfileAccessStrategy())
	}

	if len(update) == 0 {
		return buildSuccessNotification(req.RequestId), nil
	}

	err := c.userService.UpdateUser(ctx, s.UserID, update)
	if err != nil {
		return nil, err
	}

	return buildSuccessNotification(req.RequestId), nil
}
