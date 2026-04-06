package logging

import (
	"im.turms/server/internal/infra/proto"
)

// NotificationLoggingManager maps to NotificationLoggingManager in Java (turms-gateway).
// @MappedFrom NotificationLoggingManager
type NotificationLoggingManager struct{}

// SimpleNotificationFields holds the extracted fields from SimpleTurmsNotification
// used for logging, avoiding a circular dependency on the proto package.
type SimpleNotificationFields struct {
	RequesterID        int64
	CloseStatus        *int32
	RelayedRequestType proto.KindCase
}

// @MappedFrom log(SimpleTurmsNotification notification, int notificationBytes, int recipientCount, int onlineRecipientCount)
func (m *NotificationLoggingManager) Log(
	n *SimpleNotificationFields,
	notificationBytes int,
	recipientCount int,
	onlineRecipientCount int,
) {
	if n == nil {
		return
	}

	relayedStr := proto.KindCaseName(n.RelayedRequestType)

	msg := joinFields(
		n.RequesterID,
		recipientCount,
		onlineRecipientCount,
		n.CloseStatus,
		notificationBytes,
		relayedStr,
	)
	notificationLogger.Info(msg)
}
