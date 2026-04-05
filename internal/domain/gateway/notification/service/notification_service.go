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
	// Java: Validator.notNull(excludedUserSessionIds, "excludedUserSessionIds")
	if excludedUserSessionIds == nil {
		return nil, errors.New("excludedUserSessionIds cannot be nil")
	}

	offlineRecipientIdsSet := make(map[int64]struct{})
	var mu sync.Mutex

	// Bug 570: Collect errors as Java does with Mono.whenDelayError
	var sendErrors []error
	var errorsMu sync.Mutex

	hasExcludedUserSessionIds := len(excludedUserSessionIds) > 0

	var wg sync.WaitGroup

	for _, rid := range recipientIds {
		recipientID := rid
		manager := s.sessionService.GetUserSessionsManager(ctx, recipientID)
		if manager == nil {
			mu.Lock()
			offlineRecipientIdsSet[recipientID] = struct{}{}
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
			offlineRecipientIdsSet[recipientID] = struct{}{}
			mu.Unlock()
			continue
		}

		for _, userSession := range targetSessions {
			// Bug 8026: TryNotifyClientToRecover is called immediately/unconditionally in Java,
			// before the async send result is even evaluated.
			if userSession.Conn != nil {
				userSession.Conn.TryNotifyClientToRecover()
			} else {
				mu.Lock()
				offlineRecipientIdsSet[recipientID] = struct{}{}
				mu.Unlock()
				continue
			}

			wg.Add(1)
			go func(rid int64, sess *session.UserSession) {
				defer wg.Done()

				if sess.Conn == nil {
					mu.Lock()
					offlineRecipientIds = append(offlineRecipientIds, rid)
					mu.Unlock()
					return
				}

				// Java calls tryNotifyClientToRecover() unconditionally BEFORE the send,
				// in the for-loop body, not inside the send mono's error handler.
				sess.Conn.TryNotifyClientToRecover()

				err := sess.Conn.SendWithContext(ctx, notificationData)
				if err != nil {
					// Bug 8028: In Java, any error directly means the recipient is marked offline.
					mu.Lock()
					offlineRecipientIdsSet[rid] = struct{}{}
					mu.Unlock()

					if sess.IsOpen() {
						errorsMu.Lock()
						sendErrors = append(sendErrors, err)
						errorsMu.Unlock()
						log.Printf("Failed to send notification to session: user_id=%d, device_type=%s, error=%v", rid, sess.DeviceType, err)
					}
				}
			}(recipientID, userSession)
		}
	}

	wg.Wait()

	offlineRecipientIds := make([]int64, 0, len(offlineRecipientIdsSet))
	for id := range offlineRecipientIdsSet {
		offlineRecipientIds = append(offlineRecipientIds, id)
	}

	// Bug 8030: Missing logging of aggregated errors
	var finalErr error
	if len(sendErrors) > 0 {
		finalErr = errors.Join(sendErrors...)
		log.Printf("Caught an error while sending a notification to user sessions: %v", finalErr)
	}

	// Bug 558: Invoke NotificationHandler extension points
	if s.pluginManager != nil && s.pluginManager.HasRunningExtensions("NotificationHandler") {
		_, _ = s.pluginManager.InvokeExtensionPoints(ctx, "NotificationHandler", "HandleNotifications", notificationData, recipientIds, offlineRecipientIds)
	}

	return offlineRecipientIds, finalErr
}
