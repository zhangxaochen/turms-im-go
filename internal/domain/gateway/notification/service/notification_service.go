package service

import (
	"context"
	"errors"

	"im.turms/server/internal/domain/gateway/session"
	"im.turms/server/pkg/protocol"
)

type NotificationService struct {
	apiLoggingContext          any
	sessionService             *session.SessionService
	notificationLoggingManager any
	pluginManager              any
	// isNotificationLoggingEnabled bool
}

func NewNotificationService(apiLoggingContext any, sessionService *session.SessionService, notificationLoggingManager any, pluginManager any, propertiesManager any) *NotificationService {
	return &NotificationService{
		apiLoggingContext:          apiLoggingContext,
		sessionService:             sessionService,
		notificationLoggingManager: notificationLoggingManager,
		pluginManager:              pluginManager,
		// isNotificationLoggingEnabled: propertiesManager.GetLocalProperties().GetGateway().GetNotificationLogging().IsEnabled(),
	}
}

// UserSessionID maps to im.turms.server.common.domain.session.bo.UserSessionId
type UserSessionID struct {
	UserID     int64
	DeviceType protocol.DeviceType
}

// SendNotificationToLocalClients sends notification to local clients.
// @MappedFrom sendNotificationToLocalClients(TracingContext tracingContext, ByteBuf notificationData, Set<Long> recipientIds, Set<UserSessionId> excludedUserSessionIds, @Nullable DeviceType excludedDeviceType)
func (s *NotificationService) SendNotificationToLocalClients(ctx context.Context, notificationData []byte, recipientIds []int64, excludedUserSessionIds map[UserSessionID]struct{}, excludedDeviceType *protocol.DeviceType) ([]int64, error) {
	if len(notificationData) == 0 {
		return nil, errors.New("notificationData cannot be empty")
	}
	if len(recipientIds) == 0 {
		return nil, errors.New("recipientIds cannot be empty")
	}

	var offlineRecipientIds []int64
	hasExcludedUserSessionIds := len(excludedUserSessionIds) > 0

	// We collect all the sending actions
	for _, recipientID := range recipientIds {
		sessions := s.sessionService.GetAllUserSessions(recipientID)
		if len(sessions) == 0 {
			offlineRecipientIds = append(offlineRecipientIds, recipientID)
		} else {
			for _, userSession := range sessions {
				if excludedDeviceType != nil && *excludedDeviceType == userSession.DeviceType {
					continue
				}
				if hasExcludedUserSessionIds {
					id := UserSessionID{UserID: userSession.UserID, DeviceType: userSession.DeviceType}
					if _, excluded := excludedUserSessionIds[id]; excluded {
						continue
					}
				}

				err := userSession.Conn.Send(notificationData)
				if err != nil {
					offlineRecipientIds = append(offlineRecipientIds, recipientID)
					if userSession.IsSessionOpen() {
						// log error
					}
				} else {
					userSession.Conn.TryNotifyClientToRecover()
				}
			}
		}
	}

	// Logging and plugin extension invocation would go here (omitted as per stubbing strategy)

	return offlineRecipientIds, nil
}
