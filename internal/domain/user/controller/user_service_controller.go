package controller

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/protobuf/proto"

	commonservice "im.turms/server/internal/domain/common/service"
	"im.turms/server/internal/domain/gateway/access/router"
	"im.turms/server/internal/domain/gateway/session"
	"im.turms/server/internal/domain/user/service"
	"im.turms/server/internal/domain/user/service/onlineuser"
	"im.turms/server/pkg/protocol"
)

type UserServiceController struct {
	userService             service.UserService
	userRelationshipService service.UserRelationshipService
	outboundMessageService  commonservice.OutboundMessageService
	nearbyUserService       onlineuser.NearbyUserService
	sessionService          onlineuser.SessionService
	userStatusService       onlineuser.UserStatusService
	sessionLocationService  onlineuser.SessionLocationService
}

func NewUserServiceController(
	userService service.UserService,
	userRelationshipService service.UserRelationshipService,
	outboundMessageService commonservice.OutboundMessageService,
	nearbyUserService onlineuser.NearbyUserService,
	sessionService onlineuser.SessionService,
	userStatusService onlineuser.UserStatusService,
	sessionLocationService onlineuser.SessionLocationService,
) *UserServiceController {
	return &UserServiceController{
		userService:             userService,
		userRelationshipService: userRelationshipService,
		outboundMessageService:  outboundMessageService,
		nearbyUserService:       nearbyUserService,
		sessionService:          sessionService,
		userStatusService:       userStatusService,
		sessionLocationService:  sessionLocationService,
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
// @MappedFrom handleQueryUserProfilesRequest()
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
// @MappedFrom handleQueryNearbyUsersRequest()
func (c *UserServiceController) HandleQueryNearbyUsersRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	queryReq := req.GetQueryNearbyUsersRequest()

	// Read maxCount and maxDistance from request (matching Java: request.hasMaxCount() ? ... : null)
	var maxCount *int
	if queryReq.MaxCount != nil {
		mc := int(queryReq.GetMaxCount())
		maxCount = &mc
	}
	var maxDistance *float64
	if queryReq.MaxDistance != nil {
		md := float64(queryReq.GetMaxDistance())
		maxDistance = &md
	}

	nearbyUsers, err := c.nearbyUserService.QueryNearbyUsers(
		ctx,
		s.UserID,
		s.DeviceType,
		&queryReq.Longitude,
		&queryReq.Latitude,
		maxCount,
		maxDistance,
		queryReq.GetWithCoordinates(),
		queryReq.GetWithDistance(),
		queryReq.GetWithUserInfo(),
	)
	if err != nil {
		return nil, err
	}

	// Return NO_CONTENT when no nearby users found (matching Java: RequestHandlerResult.NO_CONTENT)
	if len(nearbyUsers) == 0 {
		return buildSuccessNotification(req.RequestId), nil
	}

	nearbyUserProtos := make([]*protocol.NearbyUser, 0, len(nearbyUsers))
	for _, u := range nearbyUsers {
		nu := &protocol.NearbyUser{
			UserId: u.UserID,
		}
		if u.DeviceType != nil {
			nu.DeviceType = u.DeviceType.Enum()
		}
		if u.Longitude != nil || u.Latitude != nil {
			nu.Location = &protocol.UserLocation{}
			if u.Longitude != nil {
				nu.Location.Longitude = *u.Longitude
			}
			if u.Latitude != nil {
				nu.Location.Latitude = *u.Latitude
			}
		}
		if u.Distance != nil {
			nu.Distance = proto.Int32(int32(*u.Distance))
		}
		if u.User != nil {
			nu.Info = &protocol.UserInfo{
				Id:               proto.Int64(u.User.ID),
				Name:             proto.String(u.User.Name),
				Intro:            proto.String(u.User.Intro),
				ProfilePicture:   proto.String(u.User.ProfilePicture),
				RegistrationDate: proto.Int64(u.User.RegistrationDate.UnixMilli()),
				Active:           proto.Bool(u.User.IsActive),
			}
		}
		nearbyUserProtos = append(nearbyUserProtos, nu)
	}

	return &protocol.TurmsNotification{
		RequestId: req.RequestId,
		Code:      proto.Int32(1000),
		Data: &protocol.TurmsNotification_Data{
			Kind: &protocol.TurmsNotification_Data_NearbyUsers{
				NearbyUsers: &protocol.NearbyUsers{
					NearbyUsers: nearbyUserProtos,
				},
			},
		},
	}, nil
}

// HandleQueryUserOnlineStatusesRequest queries online statuses for a set of user IDs.
// @MappedFrom handleQueryUserOnlineStatusesRequest()
func (c *UserServiceController) HandleQueryUserOnlineStatusesRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	queryReq := req.GetQueryUserOnlineStatusesRequest()
	userIDs := queryReq.GetUserIds()

	// Early return for empty user IDs (matching Java: if (request.getUserIdsCount() == 0) return Mono.empty())
	if len(userIDs) == 0 {
		return nil, nil
	}

	sessions, err := c.sessionService.QueryUserSessions(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	statusProtos := make([]*protocol.UserOnlineStatus, 0, len(sessions))
	for _, sInfo := range sessions {
		deviceTypes := make([]protocol.DeviceType, 0, len(sInfo.Sessions))
		for _, sess := range sInfo.Sessions {
			deviceTypes = append(deviceTypes, sess.DeviceType)
		}
		statusProtos = append(statusProtos, &protocol.UserOnlineStatus{
			UserId:           sInfo.UserID,
			UserStatus:       sInfo.Status,
			UsingDeviceTypes: deviceTypes,
		})
	}

	return &protocol.TurmsNotification{
		RequestId: req.RequestId,
		Code:      proto.Int32(1000),
		Data: &protocol.TurmsNotification_Data{
			Kind: &protocol.TurmsNotification_Data_UserOnlineStatuses{
				UserOnlineStatuses: &protocol.UserOnlineStatuses{
					Statuses: statusProtos,
				},
			},
		},
	}, nil
}

// HandleUpdateUserLocationRequest updates the user's current location.
// @MappedFrom handleUpdateUserLocationRequest()
func (c *UserServiceController) HandleUpdateUserLocationRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	updateReq := req.GetUpdateUserLocationRequest()
	err := c.sessionLocationService.UpsertUserLocation(ctx, s.UserID, s.DeviceType, updateReq.Longitude, updateReq.Latitude)
	if err != nil {
		return nil, err
	}
	return buildSuccessNotification(req.RequestId), nil
}

// HandleUpdateUserOnlineStatusRequest updates the user's online status (invisible, busy, etc.).
// @MappedFrom handleUpdateUserOnlineStatusRequest()
func (c *UserServiceController) HandleUpdateUserOnlineStatusRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	updateReq := req.GetUpdateUserOnlineStatusRequest()
	userStatus := updateReq.GetUserStatus()

	// Validate UNRECOGNIZED status (matching Java: if (userStatus == UserStatus.UNRECOGNIZED) return ILLEGAL_ARGUMENT)
	_, known := protocol.UserStatus_name[int32(userStatus)]
	if !known {
		return nil, fmt.Errorf("ILLEGAL_ARGUMENT: unrecognized user status %v", userStatus)
	}

	// Handle OFFLINE status: disconnect the user (matching Java: if (userStatus == UserStatus.OFFLINE) sessionService.disconnect())
	if userStatus == protocol.UserStatus_OFFLINE {
		deviceTypes := updateReq.GetDeviceTypes()
		if len(deviceTypes) > 0 {
			// Disconnect specific device types
			for _, dt := range deviceTypes {
				_, err := c.sessionService.DisconnectWithDeviceType(ctx, s.UserID, int(dt), 0)
				if err != nil {
					return nil, err
				}
			}
		} else {
			// Disconnect all devices
			_, err := c.sessionService.Disconnect(ctx, s.UserID, 0)
			if err != nil {
				return nil, err
			}
		}
		return buildSuccessNotification(req.RequestId), nil
	}

	updated, err := c.userStatusService.UpdateStatus(ctx, s.UserID, userStatus)
	if err != nil {
		return nil, err
	}

	if updated && c.outboundMessageService != nil {
		// Broadcast the status update to friends
		isBlocked := false
		friendIDs, err := c.userRelationshipService.QueryRelatedUserIds(ctx, []int64{s.UserID}, nil, &isBlocked, nil, nil)
		if err == nil && len(friendIDs) > 0 {
			notification := &protocol.TurmsNotification{
				Data: &protocol.TurmsNotification_Data{
					Kind: &protocol.TurmsNotification_Data_UserOnlineStatuses{
						UserOnlineStatuses: &protocol.UserOnlineStatuses{
							Statuses: []*protocol.UserOnlineStatus{
								{
									UserId:     s.UserID,
									UserStatus: userStatus,
								},
							},
						},
					},
				},
			}
			_ = c.outboundMessageService.ForwardNotificationToMultiple(ctx, notification, friendIDs)
		}
	}

	return buildSuccessNotification(req.RequestId), nil
}

// HandleUpdateUserRequest updates user profile fields (name, intro, profilePicture, etc.).
// @MappedFrom handleUpdateUserRequest()
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

	if c.outboundMessageService != nil {
		// Broadcast profile update to friends
		isBlocked := false
		friendIDs, err := c.userRelationshipService.QueryRelatedUserIds(ctx, []int64{s.UserID}, nil, &isBlocked, nil, nil)
		if err == nil && len(friendIDs) > 0 {
			// Query the updated user profile to send in notification
			users, err := c.userService.QueryUsersProfile(ctx, []int64{s.UserID})
			if err == nil && len(users) > 0 {
				u := users[0]
				notification := &protocol.TurmsNotification{
					Data: &protocol.TurmsNotification_Data{
						Kind: &protocol.TurmsNotification_Data_UserInfosWithVersion{
							UserInfosWithVersion: &protocol.UserInfosWithVersion{
								UserInfos: []*protocol.UserInfo{
									{
										Id:               proto.Int64(u.ID),
										Name:             proto.String(u.Name),
										Intro:            proto.String(u.Intro),
										ProfilePicture:   proto.String(u.ProfilePicture),
										RegistrationDate: proto.Int64(u.RegistrationDate.UnixMilli()),
										Active:           proto.Bool(u.IsActive),
									},
								},
							},
						},
					},
				}
				_ = c.outboundMessageService.ForwardNotificationToMultiple(ctx, notification, friendIDs)
			}
		}
	}

	return buildSuccessNotification(req.RequestId), nil
}
