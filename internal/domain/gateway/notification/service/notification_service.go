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

	var offlineRecipientIds []int64
	var mu sync.Mutex

	// Bug 570: Collect errors as Java does with Mono.whenDelayError
	var sendErrors []error
	var errorsMu sync.Mutex

	hasExcludedUserSessionIds := len(excludedUserSessionIds) > 0
	
	// Create a WaitGroup to handle concurrent sends for ALL sessions (Bug 576)
	var wg sync.WaitGroup

	// Track per-recipient success
	recipientSuccessfulCount := make(map[int64]int32)
	recipientSessionCount := make(map[int64]int32)

	for _, rid := range recipientIds {
		recipientID := rid
		manager := s.sessionService.GetUserSessionsManager(ctx, recipientID)
		if manager == nil {
			mu.Lock()
			offlineRecipientIds = append(offlineRecipientIds, recipientID)
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
			offlineRecipientIds = append(offlineRecipientIds, recipientID)
			mu.Unlock()
			continue
		}

		recipientSessionCount[recipientID] = int32(len(targetSessions))

		for _, userSession := range targetSessions {
			wg.Add(1)
			go func(rid int64, sess *session.UserSession) {
				defer wg.Done()
				
				err := sess.Conn.SendWithContext(ctx, notificationData)
				if err != nil {
					if sess.IsOpen() {
						errorsMu.Lock()
						sendErrors = append(sendErrors, err)
						errorsMu.Unlock()
						log.Printf("Failed to send notification to session: user_id=%d, device_type=%s, error=%v", rid, sess.DeviceType, err)
					}
				} else {
					mu.Lock()
					recipientSuccessfulCount[rid]++
					mu.Unlock()
				}
				sess.Conn.TryNotifyClientToRecover()
			}(recipientID, userSession)
		}
	}

	wg.Wait()

	// Identify recipients where NO session was successful
	for rid, total := range recipientSessionCount {
		if recipientSuccessfulCount[rid] == 0 && total > 0 {
			offlineRecipientIds = append(offlineRecipientIds, rid)
		}
	}

	// Bug 558: Invoke NotificationHandler extension points
	if s.pluginManager != nil && s.pluginManager.HasRunningExtensions("NotificationHandler") {
		_, _ = s.pluginManager.InvokeExtensionPoints(ctx, "NotificationHandler", "HandleNotifications", notificationData, recipientIds, offlineRecipientIds)
	}

	return offlineRecipientIds, nil
}
