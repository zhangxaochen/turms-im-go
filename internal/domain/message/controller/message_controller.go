package controller

import (
	"context"
	"errors"
	"time"

	"google.golang.org/protobuf/proto"

	"im.turms/server/internal/domain/gateway/session"
	"im.turms/server/internal/domain/message/po"
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

// ConvertMessageToProto converts po.Message to protocol.Message
func ConvertMessageToProto(msg *po.Message) *protocol.Message {
	if msg == nil {
		return nil
	}
	m := &protocol.Message{
		Id:              proto.Int64(msg.ID),
		DeliveryDate:    proto.Int64(msg.DeliveryDate.UnixMilli()),
		Text:            proto.String(msg.Text),
		SenderId:        proto.Int64(msg.SenderID),
		IsSystemMessage: proto.Bool(false),
		SequenceId:      msg.SequenceID,
		PreMessageId:    msg.PreMessageID,
		Records:         msg.Records,
	}
	if msg.IsGroupMessage != nil && *msg.IsGroupMessage {
		m.GroupId = proto.Int64(msg.TargetID)
	} else {
		m.RecipientId = proto.Int64(msg.TargetID)
	}
	if msg.IsSystemMessage != nil {
		m.IsSystemMessage = msg.IsSystemMessage
	}
	if msg.ModificationDate != nil {
		m.ModificationDate = proto.Int64(msg.ModificationDate.UnixMilli())
	}
	return m
}

// @MappedFrom handleQueryMessagesRequest()
func (c *MessageController) HandleQueryMessagesRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	queryReq := req.GetQueryMessagesRequest()

	var deliveryDateAfter, deliveryDateBefore *time.Time
	if queryReq.DeliveryDateStart != nil {
		t := time.UnixMilli(*queryReq.DeliveryDateStart)
		deliveryDateAfter = &t
	}
	if queryReq.DeliveryDateEnd != nil {
		t := time.UnixMilli(*queryReq.DeliveryDateEnd)
		deliveryDateBefore = &t
	}

	size := int64(50) // Default max count
	if queryReq.MaxCount != nil {
		size = int64(*queryReq.MaxCount)
	}
	ascending := false
	if queryReq.Descending != nil {
		ascending = !*queryReq.Descending
	}

	messages, err := c.messageService.QueryMessages(
		ctx,
		s.UserID,
		queryReq.AreGroupMessages,
		queryReq.FromIds,
		nil, // targetIDs logic can be complex, default to nil or infer from context
		deliveryDateAfter,
		deliveryDateBefore,
		size,
		ascending,
	)
	if err != nil {
		return nil, err
	}

	// NOTE: withTotal is ignored for now as we don't have a CountMessages method implemented.
	var protoMessages []*protocol.Message
	for _, m := range messages {
		protoMessages = append(protoMessages, ConvertMessageToProto(m))
	}

	return &protocol.TurmsNotification{
		RequestId: req.RequestId,
		Code:      proto.Int32(1000), // SUCCESS
		Data: &protocol.TurmsNotification_Data{
			Kind: &protocol.TurmsNotification_Data_Messages{
				Messages: &protocol.Messages{
					Messages: protoMessages,
				},
			},
		},
	}, nil
}

// @MappedFrom handleUpdateMessageRequest()
func (c *MessageController) HandleUpdateMessageRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	updateReq := req.GetUpdateMessageRequest()

	if updateReq.RecallDate != nil {
		err := c.messageService.AuthAndRecallMessage(ctx, s.UserID, updateReq.MessageId)
		if err != nil {
			return nil, err
		}
	} else if updateReq.Text != nil {
		err := c.messageService.AuthAndUpdateMessageText(ctx, s.UserID, updateReq.MessageId, *updateReq.Text)
		if err != nil {
			return nil, err
		}
	}

	return &protocol.TurmsNotification{
		RequestId: req.RequestId,
		Code:      proto.Int32(1000), // SUCCESS
	}, nil
}

// @MappedFrom handleCreateMessageReactionsRequest()
func (c *MessageController) HandleCreateMessageReactionsRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	// The feature hasn't published yet
	return nil, errors.New("ILLEGAL_ARGUMENT") // ILLEGAL_ARGUMENT equivalent
}

// @MappedFrom handleDeleteMessageReactionsRequest()
func (c *MessageController) HandleDeleteMessageReactionsRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	// The feature hasn't published yet
	return nil, errors.New("ILLEGAL_ARGUMENT") // ILLEGAL_ARGUMENT equivalent
}
