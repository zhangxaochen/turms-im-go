package service

import (
	"context"
	"errors"
	"log"
	"sync"

	"golang.org/x/sync/errgroup"
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
	// Bug 560: Java version checks non-null. In Go, we check the slice reference.
	if notificationData == nil {
		return nil, errors.New("notificationData cannot be nil")
	}
	if len(recipientIds) == 0 {
		return nil, errors.New("recipientIds cannot be empty")
	}

	hasExcludedUserSessionIds := len(excludedUserSessionIds) > 0
	var offlineRecipientIds []int64
	var mu sync.Mutex

	g, gCtx := errgroup.WithContext(ctx)

	for _, rid := range recipientIds {
		recipientID := rid
		g.Go(func() error {
			manager := s.sessionService.GetUserSessionsManager(gCtx, recipientID)
			if manager == nil {
				mu.Lock()
				offlineRecipientIds = append(offlineRecipientIds, recipientID)
				mu.Unlock()
				return nil
			}

			sessions := manager.GetAllSessions()
			isAnySent := false
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

				// Bug 572: Use gCtx (derived from parent ctx) to support cancellation propagation
				err := userSession.Conn.SendWithContext(gCtx, notificationData)
				if err != nil {
					// Bug 568: Log error as in Java
					if userSession.IsOpen() {
						log.Printf("Failed to send notification to session: user_id=%d, device_type=%s, error=%v", recipientID, userSession.DeviceType, err)
					}
				} else {
					isAnySent = true
				}
				// Bug 564: Unconditional recover notification
				userSession.Conn.TryNotifyClientToRecover()
			}

			if !isAnySent {
				mu.Lock()
				offlineRecipientIds = append(offlineRecipientIds, recipientID)
				mu.Unlock()
			}
			return nil
		})
	}

	// Wait for all sends to complete (Bug 576: parallel execution)
	_ = g.Wait()

	// Logging and plugin extension invocation would go here (omitted as per current stubbing strategy)

	return offlineRecipientIds, nil
}
