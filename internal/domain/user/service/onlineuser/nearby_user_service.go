package onlineuser

import (
	"context"

	"im.turms/server/internal/domain/user/po"
	"im.turms/server/internal/domain/user/service"
)

type NearbyUser struct {
	UserID     int64
	DeviceType *int
	Longitude  *float32
	Latitude   *float32
	Distance   *int
	User       *po.User
}

type NearbyUserService struct {
	userService service.UserService
	// sessionLocationService *SessionLocationService
}

func NewNearbyUserService(userService service.UserService) *NearbyUserService {
	return &NearbyUserService{
		userService: userService,
	}
}

func (s *NearbyUserService) QueryNearbyUsers(ctx context.Context, userID int64, deviceType int, longitude *float32, latitude *float32, maxCount *int16, maxDistance *int, withCoordinates bool, withDistance bool, withUserInfo bool) ([]*NearbyUser, error) {
	// Stub implementation for compilation
	return []*NearbyUser{}, nil
}
