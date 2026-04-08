package onlineuser

import (
	"context"
	"im.turms/server/pkg/protocol"
)

type UserSessionInfo struct {
	DeviceType protocol.DeviceType
}

type UserSessionsInfo struct {
	UserID   int64
	Status   protocol.UserStatus
	Sessions []UserSessionInfo
}

type SessionService interface {
	Disconnect(ctx context.Context, userID int64, closeStatus int) (bool, error)
	DisconnectWithDeviceTypes(ctx context.Context, userID int64, deviceTypes []int, closeStatus int) (bool, error)
	DisconnectWithDeviceType(ctx context.Context, userID int64, deviceType int, closeStatus int) (bool, error)
	DisconnectMultipleUsers(ctx context.Context, userIDs []int64, closeStatus int) (bool, error)
	DisconnectMultipleUsersWithDeviceTypes(ctx context.Context, userIDs []int64, deviceTypes []int, closeStatus int) (bool, error)
	QueryUserSessions(ctx context.Context, userIDs []int64) ([]*UserSessionsInfo, error)
}

type sessionService struct {
	userStatusService UserStatusService
	// rpcService        *RpcService
}

func NewSessionService(userStatusService UserStatusService) SessionService {
	return &sessionService{
		userStatusService: userStatusService,
	}
}

// @MappedFrom disconnect(@NotNull Long userId, @NotNull SessionCloseStatus closeStatus)
// @MappedFrom disconnect(@NotNull Set<Long> userIds, @Nullable Set<@ValidDeviceType DeviceType> deviceTypes, @NotNull SessionCloseStatus closeStatus)
// @MappedFrom disconnect(@NotNull Long userId, @NotNull @ValidDeviceType DeviceType deviceType, @NotNull SessionCloseStatus closeStatus)
// @MappedFrom disconnect(@NotNull Set<Long> userIds, @NotNull SessionCloseStatus closeStatus)
// @MappedFrom disconnect(@NotNull Long userId, @NotNull Set<@ValidDeviceType DeviceType> deviceTypes, @NotNull SessionCloseStatus closeStatus)
func (s *sessionService) Disconnect(ctx context.Context, userID int64, closeStatus int) (bool, error) {
	// TODO: Implement using UserStatusService and RPC to notify gateway
	return true, nil
}

func (s *sessionService) DisconnectWithDeviceTypes(ctx context.Context, userID int64, deviceTypes []int, closeStatus int) (bool, error) {
	return true, nil
}

func (s *sessionService) DisconnectWithDeviceType(ctx context.Context, userID int64, deviceType int, closeStatus int) (bool, error) {
	return true, nil
}

func (s *sessionService) DisconnectMultipleUsers(ctx context.Context, userIDs []int64, closeStatus int) (bool, error) {
	return true, nil
}

func (s *sessionService) DisconnectMultipleUsersWithDeviceTypes(ctx context.Context, userIDs []int64, deviceTypes []int, closeStatus int) (bool, error) {
	return true, nil
}

// @MappedFrom queryUserSessions(Set<Long> userIds)
// @MappedFrom queryUserSessions(Set<Long> ids, boolean returnNonExistingUsers)
func (s *sessionService) QueryUserSessions(ctx context.Context, userIDs []int64) ([]*UserSessionsInfo, error) {
	infos := make([]*UserSessionsInfo, 0, len(userIDs))
	for _, uid := range userIDs {
		status, err := s.userStatusService.FetchUserSessionsStatus(ctx, uid)
		if err != nil {
			return nil, err // Propagate error instead of silently continuing
		}

		var sessions []UserSessionInfo
		for dtype, sessInfo := range status.OnlineDeviceTypeToSessionInfo {
			if sessInfo.IsActive {
				sessions = append(sessions, UserSessionInfo{
					DeviceType: dtype,
				})
			}
		}

		infos = append(infos, &UserSessionsInfo{
			UserID:   uid,
			Status:   status.UserStatus,
			Sessions: sessions,
		})
	}
	return infos, nil
}
