package common

import (
	"google.golang.org/protobuf/proto"

	"im.turms/server/pkg/protocol"
)

// NotificationFactory standardizes the creation of TurmsNotification objects
// sent down to clients containing success/failure outcomes.
type NotificationFactory struct{}

func NewNotificationFactory() *NotificationFactory {
	return &NotificationFactory{}
}

// Create generates a generic Notification payload.
func (f *NotificationFactory) Create(requestID *int64, code int32, reason string) *protocol.TurmsNotification {
	notification := &protocol.TurmsNotification{
		RequestId: requestID,
		Code:      proto.Int32(code),
	}

	if reason != "" {
		notification.Reason = proto.String(reason)
	}

	return notification
}

// CreateBuffer generates the serialized protobuf bytes directly.
func (f *NotificationFactory) CreateBuffer(requestID *int64, code int32, reason string) ([]byte, error) {
	notification := f.Create(requestID, code, reason)
	return proto.Marshal(notification)
}

// SessionClosed generates a specialized notification when the server forcefully kicks the client.
func (f *NotificationFactory) SessionClosed(reason string) *protocol.TurmsNotification {
	return &protocol.TurmsNotification{
		CloseStatus: proto.Int32(2000), // Standard Turms server-initiated close offset
		Reason:      proto.String(reason),
	}
}
