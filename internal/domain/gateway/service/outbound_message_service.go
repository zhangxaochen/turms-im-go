package service

import (
	"context"
	"errors"

	"google.golang.org/protobuf/proto"

	"im.turms/server/internal/domain/common/service"
	"im.turms/server/internal/domain/gateway/session"
	"im.turms/server/pkg/protocol"
)

type outboundMessageService struct {
	sessionService *session.SessionService
}

func NewOutboundMessageService(sessionService *session.SessionService) service.OutboundMessageService {
	return &outboundMessageService{
		sessionService: sessionService,
	}
}

func (s *outboundMessageService) ForwardNotification(ctx context.Context, notification *protocol.TurmsNotification, recipientID int64) error {
	if notification == nil {
		return errors.New("notification is nil")
	}

	sessions := s.sessionService.GetAllUserSessions(recipientID)
	if len(sessions) == 0 {
		// Recipient is offline on this node
		// (In a cluster, we would forward this via RPC to other nodes)
		return nil
	}

	data, err := proto.Marshal(notification)
	if err != nil {
		return err
	}

	for _, userSession := range sessions {
		if userSession.Conn != nil {
			_ = userSession.Conn.Send(data)
		}
	}

	return nil
}

func (s *outboundMessageService) ForwardNotificationToMultiple(ctx context.Context, notification *protocol.TurmsNotification, recipientIds []int64) error {
	if notification == nil || len(recipientIds) == 0 {
		return nil
	}

	data, err := proto.Marshal(notification)
	if err != nil {
		return err
	}

	for _, recipientID := range recipientIds {
		sessions := s.sessionService.GetAllUserSessions(recipientID)
		for _, userSession := range sessions {
			if userSession.Conn != nil {
				_ = userSession.Conn.Send(data)
			}
		}
	}

	return nil
}
