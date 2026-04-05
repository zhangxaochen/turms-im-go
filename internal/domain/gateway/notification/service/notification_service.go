package service

import (
	"context"
	"errors"
	"log"
	"sync"

	"im.turms/server/internal/domain/gateway/session"
	"im.turms/server/internal/infra/plugin"
	"im.turms/server/pkg/protocol"
)

type NotificationService struct {
	apiLoggingContext          any
	sessionService             *session.SessionService
	notificationLoggingManager any
	pluginManager              *plugin.PluginManager
	// isNotificationLoggingEnabled bool
}

func NewNotificationService(apiLoggingContext any, sessionService *session.SessionService, notificationLoggingManager any, pluginManager *plugin.PluginManager, propertiesManager any) *NotificationService {
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
	if notificationData == nil {
		return nil, errors.New("notificationData cannot be nil")
	}
	if len(recipientIds) == 0 {
		return nil, errors.New("recipientIds cannot be empty")
	}

	// Use a map as a set for offlineRecipientIds to prevent duplicates (Java uses Set<Long>)
	offlineSet := make(map[int64]struct{})
	var mu sync.Mutex

	// Collect errors as Java does with Mono.whenDelayError
	var sendErrors []error
	var errorsMu sync.Mutex

	hasExcludedUserSessionIds := len(excludedUserSessionIds) > 0

	var wg sync.WaitGroup

	for _, rid := range recipientIds {
		recipientID := rid
		manager := s.sessionService.GetUserSessionsManager(ctx, recipientID)
		if manager == nil {
			mu.Lock()
			offlineSet[recipientID] = struct{}{}
			mu.Unlock()
			continue
		}

		sessions := manager.GetAllSessions()

		// Filter sessions
		var targetSessions []*session.UserSession
		for _, sess := range sessions {
			if excludedDeviceType != nil && *excludedDeviceType == sess.DeviceType {
				continue
			}
			if hasExcludedUserSessionIds {
				id := UserSessionID{UserID: sess.UserID, DeviceType: sess.DeviceType}
				if _, excluded := excludedUserSessionIds[id]; excluded {
					continue
				}
			}
			targetSessions = append(targetSessions, sess)
		}

		if len(targetSessions) == 0 {
			mu.Lock()
			offlineSet[recipientID] = struct{}{}
			mu.Unlock()
			continue
		}

		for _, userSession := range targetSessions {
			wg.Add(1)
			go func(rid int64, sess *session.UserSession) {
				defer wg.Done()

				if sess.Conn == nil {
					// Java parity: onErrorResume immediately adds to offlineRecipientIds
					mu.Lock()
					offlineSet[rid] = struct{}{}
					mu.Unlock()
					return
				}

				err := sess.Conn.SendWithContext(ctx, notificationData)
				if err != nil {
					// Java parity: immediately add to offlineRecipientIds on any session send error
					mu.Lock()
					offlineSet[rid] = struct{}{}
					mu.Unlock()
					if sess.IsOpen() {
						errorsMu.Lock()
						sendErrors = append(sendErrors, err)
						errorsMu.Unlock()
						log.Printf("Failed to send notification to session: user_id=%d, device_type=%s, error=%v", rid, sess.DeviceType, err)
						sess.Conn.TryNotifyClientToRecover()
					}
				}
			}(recipientID, userSession)
		}
	}

	wg.Wait()

	// Log aggregated error message (matches Java: "Caught an error while sending a notification to user sessions")
	if len(sendErrors) > 0 {
		log.Printf("Caught %d error(s) while sending notifications to user sessions", len(sendErrors))
	}

	// Convert offlineSet to slice
	offlineRecipientIds := make([]int64, 0, len(offlineSet))
	for id := range offlineSet {
		offlineRecipientIds = append(offlineRecipientIds, id)
	}

	// Invoke NotificationHandler extension points
	if s.pluginManager != nil && s.pluginManager.HasRunningExtensions("NotificationHandler") {
		_, _ = s.pluginManager.InvokeExtensionPoints(ctx, "NotificationHandler", "HandleNotifications", notificationData, recipientIds, offlineRecipientIds)
	}

	return offlineRecipientIds, nil
}
