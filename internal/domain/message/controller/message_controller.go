package controller

import (
	"context"

	"google.golang.org/protobuf/proto"

	"im.turms/server/internal/domain/gateway/session"
	"im.turms/server/internal/domain/message/service"
	"im.turms/server/pkg/protocol"
)

type MessageController struct {
	messageService *service.MessageService
}

func NewMessageController(messageService *service.MessageService) *MessageController {
	return &MessageController{
		messageService: messageService,
	}
}

// HandleCreateMessageRequest handles the creation of a message from the client.
// @MappedFrom handleCreateMessageRequest()
func (c *MessageController) HandleCreateMessageRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	createReq := req.GetCreateMessageRequest()

	var targetID int64
	var text string
	if createReq.Text != nil {
		text = *createReq.Text
	}

	isGroupMessage := createReq.GroupId != nil

	if isGroupMessage {
		targetID = *createReq.GroupId
	} else {
		if createReq.RecipientId != nil {
			targetID = *createReq.RecipientId
		}
	}

	createdMsg, err := c.messageService.AuthAndSaveAndSendMessage(ctx, isGroupMessage, s.UserID, targetID, text)
	if err != nil {
		return nil, err
	}

	return &protocol.TurmsNotification{
		RequestId: req.RequestId,
		Code:      proto.Int32(1000), // SUCCESS
		Data: &protocol.TurmsNotification_Data{
			Kind: &protocol.TurmsNotification_Data_Long{
				Long: createdMsg.ID,
			},
		},
	}, nil
}
