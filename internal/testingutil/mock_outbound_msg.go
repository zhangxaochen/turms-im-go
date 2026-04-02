package testingutil

import (
	"context"
	"im.turms/server/pkg/protocol"
)

type MockOutboundMessageService struct{}

func NewMockOutboundMessageService() *MockOutboundMessageService {
	return &MockOutboundMessageService{}
}

func (m *MockOutboundMessageService) ForwardNotification(ctx context.Context, notification *protocol.TurmsNotification, recipientID int64) error {
	return nil
}

func (m *MockOutboundMessageService) ForwardNotificationToMultiple(ctx context.Context, notification *protocol.TurmsNotification, recipientIds []int64) error {
	return nil
}
