package service

import (
	"context"
	"im.turms/server/pkg/protocol"
)

type OutboundMessageService interface {
	// ForwardNotification sends a notification to target recipients.
	// recipientIds can be a single ID or multiple.
	ForwardNotification(ctx context.Context, notification *protocol.TurmsNotification, recipientID int64) error
	ForwardNotificationToMultiple(ctx context.Context, notification *protocol.TurmsNotification, recipientIds []int64) error
}
