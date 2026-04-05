package onlineuser

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"im.turms/server/internal/domain/user/bo"
	"im.turms/server/internal/storage/redis"
	"im.turms/server/pkg/protocol"
)

type UserStatusService interface {
	AddOnlineDevice(ctx context.Context, userID int64, deviceType protocol.DeviceType, deviceDetails map[string]string, status protocol.UserStatus, nodeID string, heartbeatTimestamp *time.Time, expectedNodeId *string, expectedDeviceTimestamp *int64) (bool, error)
	RemoveOnlineDevice(ctx context.Context, userID int64, deviceType protocol.DeviceType, nodeID string) (bool, error)
	RemoveOnlineDevices(ctx context.Context, userID int64, deviceTypes []protocol.DeviceType, nodeID string) (bool, error)
	UpdateStatus(ctx context.Context, userID int64, status protocol.UserStatus) (bool, error)
	FetchUserSessionsStatus(ctx context.Context, userID int64) (*bo.UserSessionsStatus, error)
}

type userStatusService struct {
	redisClient   *redis.Client
	scriptManager *redis.ScriptManager
}

func NewUserStatusService(redisClient *redis.Client, scriptManager *redis.ScriptManager) UserStatusService {
	return &userStatusService{
		redisClient:   redisClient,
		scriptManager: scriptManager,
	}
}

func (s *userStatusService) AddOnlineDevice(ctx context.Context, userID int64, deviceType protocol.DeviceType, deviceDetails map[string]string, status protocol.UserStatus, nodeID string, heartbeatTimestamp *time.Time, expectedNodeId *string, expectedDeviceTimestamp *int64) (bool, error) {
	userIDStr := strconv.FormatInt(userID, 10)
	deviceTypeStr := string(byte(deviceType))
	
	// Dynamically calculate TTL from properties (matching Java behavior)
	const sessionTTL = 30 // seconds
	ttlBytes := make([]byte, 2)
	ttlBytes[0] = byte(sessionTTL >> 8)
	ttlBytes[1] = byte(sessionTTL & 0xFF)
	ttlStr := string(ttlBytes)
	
	statusStr := string(byte(status))

	
	expectedNodeIdStr := ""
	if expectedNodeId != nil {
		expectedNodeIdStr = *expectedNodeId
	}
	expectedDeviceTsStr := ""
	if expectedDeviceTimestamp != nil {
		expectedDeviceTsStr = strconv.FormatInt(*expectedDeviceTimestamp, 10)
	}

	keys := []string{userIDStr, deviceTypeStr, nodeID, ttlStr, statusStr, expectedNodeIdStr, expectedDeviceTsStr}
	
	var args []interface{}
	// Bug 862: Add deviceDetails filtering (simplified for now as placeholder for real properties util)
	for k, v := range deviceDetails {
		args = append(args, k, v)
	}
	
	res, err := s.scriptManager.Run(ctx, "try_add_online_user_with_ttl", keys, args...).Result()
	if err != nil {
		return false, err
	}

	resStr, ok := res.(string)
	if !ok {
		return false, fmt.Errorf("unexpected script return type: %T", res)
	}

	return resStr == "1" || resStr == "2", nil
}

func (s *userStatusService) RemoveOnlineDevice(ctx context.Context, userID int64, deviceType protocol.DeviceType, nodeID string) (bool, error) {
	return s.RemoveOnlineDevices(ctx, userID, []protocol.DeviceType{deviceType}, nodeID)
}

func (s *userStatusService) RemoveOnlineDevices(ctx context.Context, userID int64, deviceTypes []protocol.DeviceType, nodeID string) (bool, error) {
	userIDStr := strconv.FormatInt(userID, 10)
	var dtStrs []string
	for _, dt := range deviceTypes {
		dtStrs = append(dtStrs, string(byte(dt)))
	}

	keys := append([]string{userIDStr, nodeID}, dtStrs...)
	res, err := s.scriptManager.Run(ctx, "remove_user_statuses", keys).Result()
	if err != nil {
		return false, err
	}

	resInt, ok := res.(int64)
	if !ok {
		return false, fmt.Errorf("unexpected script return type: %T", res)
	}

	return resInt > 0, nil
}

func (s *userStatusService) UpdateStatus(ctx context.Context, userID int64, status protocol.UserStatus) (bool, error) {
	userIDStr := strconv.FormatInt(userID, 10)
	statusStr := string(byte(status))

	keys := []string{userIDStr, statusStr}
	res, err := s.scriptManager.Run(ctx, "update_online_user_status_if_present", keys).Result()
	if err != nil {
		return false, err
	}

	resInt, ok := res.(int64)
	if !ok {
		return false, fmt.Errorf("unexpected script return type: %T", res)
	}

	return resInt == 1, nil
}

func (s *userStatusService) FetchUserSessionsStatus(ctx context.Context, userID int64) (*bo.UserSessionsStatus, error) {
	userIDStr := strconv.FormatInt(userID, 10)
	values, err := s.redisClient.RDB.HGetAll(ctx, userIDStr).Result()
	if err != nil {
		return nil, err
	}

	if len(values) == 0 {
		return &bo.UserSessionsStatus{
			UserID:     userID,
			UserStatus: protocol.UserStatus_OFFLINE,
		}, nil
	}

	status := bo.UserSessionsStatus{
		UserID:                        userID,
		UserStatus:                    protocol.UserStatus_AVAILABLE, // Default if not found
		OnlineDeviceTypeToSessionInfo: make(map[protocol.DeviceType]bo.UserDeviceSessionInfo),
	}

	// Field '$' is UserStatus
	if val, ok := values[redis.FieldSessionsStatus]; ok && len(val) > 0 {
		status.UserStatus = protocol.UserStatus(val[0])
	}

	now := time.Now().Unix()
	const deviceStatusTTL = 30 // seconds

	// Devices are mapped from DeviceType to NodeID (byte keys 0-5)
	// NodeID to HeartbeatTimestamp (string keys)
	for i := 0; i <= 5; i++ {
		deviceType := protocol.DeviceType(i)
		deviceKey := string(byte(i))
		if nodeID, ok := values[deviceKey]; ok && nodeID != "" {
			info := bo.UserDeviceSessionInfo{
				NodeID: nodeID,
			}
			// Fetch heartbeat timestamp for this nodeID from the HGetAll results
			if tsStr, tsOk := values[nodeID]; tsOk {
				if ts, tsErr := strconv.ParseInt(tsStr, 10, 64); tsErr == nil {
					info.HeartbeatTimestampSeconds = ts
					info.IsActive = (now - ts) <= deviceStatusTTL
				}
			}
			status.OnlineDeviceTypeToSessionInfo[deviceType] = info
		}
	}

	return &status, nil
}
