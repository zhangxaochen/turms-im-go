package onlineuser

import (
	"context"
	"fmt"

	"im.turms/server/internal/storage/redis"
	"im.turms/server/pkg/protocol"
	goredis "github.com/redis/go-redis/v9"
)

type SessionLocationService interface {
	UpsertUserLocation(ctx context.Context, userID int64, deviceType protocol.DeviceType, longitude float32, latitude float32) error
	RemoveUserLocation(ctx context.Context, userID int64, deviceType protocol.DeviceType) error
	GetUserLocation(ctx context.Context, userID int64, deviceType protocol.DeviceType) (*protocol.UserLocation, error)
}

type sessionLocationService struct {
	redisClient *redis.Client
}

func NewSessionLocationService(redisClient *redis.Client) SessionLocationService {
	return &sessionLocationService{
		redisClient: redisClient,
	}
}

func (s *sessionLocationService) UpsertUserLocation(ctx context.Context, userID int64, deviceType protocol.DeviceType, longitude float32, latitude float32) error {
	// member: userID:deviceType
	member := fmt.Sprintf("%d:%d", userID, deviceType)
	
	err := s.redisClient.RDB.GeoAdd(ctx, redis.KeyLocation, &goredis.GeoLocation{
		Name:      member,
		Longitude: float64(longitude),
		Latitude:  float64(latitude),
	}).Err()
	
	return err
}

func (s *sessionLocationService) RemoveUserLocation(ctx context.Context, userID int64, deviceType protocol.DeviceType) error {
	member := fmt.Sprintf("%d:%d", userID, deviceType)
	err := s.redisClient.RDB.ZRem(ctx, redis.KeyLocation, member).Err()
	return err
}

func (s *sessionLocationService) GetUserLocation(ctx context.Context, userID int64, deviceType protocol.DeviceType) (*protocol.UserLocation, error) {
	member := fmt.Sprintf("%d:%d", userID, deviceType)
	pos, err := s.redisClient.RDB.GeoPos(ctx, redis.KeyLocation, member).Result()
	if err != nil {
		return nil, err
	}
	
	if len(pos) == 0 || pos[0] == nil {
		return nil, nil
	}
	
	return &protocol.UserLocation{
		Longitude: float32(pos[0].Longitude),
		Latitude:  float32(pos[0].Latitude),
	}, nil
}
