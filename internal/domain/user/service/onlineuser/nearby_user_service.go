package onlineuser

import (
	"context"

	"fmt"

	"im.turms/server/internal/domain/user/po"
	"im.turms/server/internal/domain/user/service"
	"im.turms/server/internal/storage/redis"
	"im.turms/server/pkg/protocol"
	goredis "github.com/redis/go-redis/v9"
)

type NearbyUser struct {
	UserID     int64
	DeviceType *protocol.DeviceType
	Longitude  *float32
	Latitude   *float32
	Distance   *float64
	User       *po.User
}

type NearbyUserService interface {
	QueryNearbyUsers(ctx context.Context, userID int64, deviceType protocol.DeviceType, longitude *float32, latitude *float32, maxCount *int, maxDistance *float64, withCoordinates bool, withDistance bool, withUserInfo bool) ([]*NearbyUser, error)
}

type nearbyUserService struct {
	userService            service.UserService
	sessionLocationService SessionLocationService
	redisClient            *redis.Client
}

func NewNearbyUserService(userService service.UserService, sessionLocationService SessionLocationService, redisClient *redis.Client) NearbyUserService {
	return &nearbyUserService{
		userService:            userService,
		sessionLocationService: sessionLocationService,
		redisClient:            redisClient,
	}
}

func (s *nearbyUserService) QueryNearbyUsers(ctx context.Context, userID int64, deviceType protocol.DeviceType, longitude *float32, latitude *float32, maxCount *int, maxDistance *float64, withCoordinates bool, withDistance bool, withUserInfo bool) ([]*NearbyUser, error) {
	if longitude == nil || latitude == nil {
		return []*NearbyUser{}, nil
	}

	// Use go-redis GeoSearch
	radius := 1000.0
	if maxDistance != nil {
		radius = *maxDistance
	}
	limit := 10
	if maxCount != nil {
		limit = *maxCount
	}

	res, err := s.redisClient.RDB.GeoSearchLocation(ctx, redis.KeyLocation, &goredis.GeoSearchLocationQuery{
		GeoSearchQuery: goredis.GeoSearchQuery{
			Longitude:  float64(*longitude),
			Latitude:   float64(*latitude),
			Radius:     radius,
			RadiusUnit: "m",
			Count:      limit,
		},
		WithCoord: withCoordinates,
		WithDist:  withDistance,
	}).Result()

	if err != nil {
		return nil, err
	}

	nearbyUsers := make([]*NearbyUser, 0, len(res))
	userIDs := make([]int64, 0, len(res))
	for _, loc := range res {
		// member format: userID:deviceType
		var uid int64
		var dtype protocol.DeviceType
		_, parseErr := fmt.Sscanf(loc.Name, "%d:%d", &uid, &dtype)
		if parseErr != nil {
			continue
		}
		
		if uid == userID {
			continue
		}

		nu := &NearbyUser{
			UserID:     uid,
			DeviceType: &dtype,
		}
		if withCoordinates {
			lon := float32(loc.Longitude)
			lat := float32(loc.Latitude)
			nu.Longitude = &lon
			nu.Latitude = &lat
		}
		if withDistance {
			dist := loc.Dist
			nu.Distance = &dist
		}
		nearbyUsers = append(nearbyUsers, nu)
		if withUserInfo {
			userIDs = append(userIDs, uid)
		}
	}

	if withUserInfo && len(userIDs) > 0 {
		// Fix: QueryUsersProfile only takes (ctx, userIDs)
		users, profileErr := s.userService.QueryUsersProfile(ctx, userIDs)
		if profileErr == nil {
			userMap := make(map[int64]*po.User)
			for _, u := range users {
				userMap[u.ID] = u
			}
			for _, nu := range nearbyUsers {
				nu.User = userMap[nu.UserID]
			}
		}
	}

	return nearbyUsers, nil
}
