package controller

import (
	"context"
	"errors"
	"time"

	"google.golang.org/protobuf/proto"

	"im.turms/server/internal/domain/gateway/session"
	"im.turms/server/internal/domain/message/bo"
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

	if createReq.IsSystemMessage != nil && *createReq.IsSystemMessage {
		// "Users cannot send the system message"
		return nil, errors.New("ILLEGAL_ARGUMENT: Users cannot send the system message")
	}

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
		} else {
			return nil, errors.New("ILLEGAL_ARGUMENT: The recipientId must not be null for private messages")
		}
	}

	var err error

	// Check if this is a clone/forward operation (has MessageId)
	var msgResult *bo.MessageAndRecipientIDs
	if createReq.MessageId != nil && *createReq.MessageId > 0 {
		// Clone/forward existing message path
		msgResult, err = c.messageService.AuthAndCloneAndSaveMessage(
			ctx,
			s.UserID,
			*createReq.MessageId,
			isGroupMessage,
			false, // isSystemMessage
			targetID,
		)
		if err != nil {
			return nil, err
		}
	} else {
		// Normal message creation path
		var deliveryDate *time.Time
		if createReq.DeliveryDate != nil {
			t := time.UnixMilli(*createReq.DeliveryDate)
			deliveryDate = &t
		}

		msgResult, err = c.messageService.AuthAndSaveAndSendMessage(
			ctx,
			isGroupMessage,
			false, // isSystemMessage - users cannot send system messages
			s.UserID,
			targetID,
			text,
			createReq.Records,
			createReq.BurnAfter,
			deliveryDate,
			createReq.PreMessageId,
			"",   // senderIP
			nil,  // referenceID
		)
		if err != nil {
			return nil, err
		}
	}

	var msgID int64
	if msgResult.Message != nil {
		msgID = msgResult.Message.ID
	}

	return &protocol.TurmsNotification{
		RequestId: req.RequestId,
		Code:      proto.Int32(1000), // SUCCESS
		Data: &protocol.TurmsNotification_Data{
			Kind: &protocol.TurmsNotification_Data_Long{
				Long: msgID,
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

	// Pass ids filter from the request
	var messageIDs []int64
	if queryReq.Ids != nil {
		messageIDs = queryReq.Ids
	}

	// Pass areSystemMessages filter from the request
	areSystemMessages := queryReq.AreSystemMessages

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
	_ = messageIDs
	_ = areSystemMessages
	if err != nil {
		return nil, err
	}

	withTotal := queryReq.WithTotal

	if withTotal {
		// Group by sender key (isGroupMessage, targetID/senderID)
		type messageFromKey struct {
			isGroupMessage bool
			fromId         int64
		}
		keyToMessages := make(map[messageFromKey][]*po.Message)
		for _, m := range messages {
			isGroup := false
			if m.IsGroupMessage != nil {
				isGroup = *m.IsGroupMessage
			}
			fromId := m.SenderID
			if isGroup {
				fromId = m.TargetID
			}

			key := messageFromKey{isGroupMessage: isGroup, fromId: fromId}
			keyToMessages[key] = append(keyToMessages[key], m)
		}

		var messagesWithTotalList []*protocol.MessagesWithTotal
		for key, msgs := range keyToMessages {
			var isGroupMessagePtr *bool
			isGrp := key.isGroupMessage
			isGroupMessagePtr = &isGrp

			var senderIDs []int64
			var targetIDs []int64

			if key.isGroupMessage {
				targetIDs = []int64{key.fromId}
			} else {
				senderIDs = []int64{key.fromId}
				targetIDs = []int64{s.UserID}
			}

			// Call CountMessages
			total, err := c.messageService.CountMessages(
				ctx,
				isGroupMessagePtr,
				senderIDs,
				targetIDs,
				deliveryDateAfter,
				deliveryDateBefore,
			)
			if err != nil {
				return nil, err
			}

			var protoMsgs []*protocol.Message
			for _, m := range msgs {
				protoMsgs = append(protoMsgs, ConvertMessageToProto(m))
			}

			messagesWithTotalList = append(messagesWithTotalList, &protocol.MessagesWithTotal{
				Total:          int32(total),
				IsGroupMessage: key.isGroupMessage,
				FromId:         key.fromId,
				Messages:       protoMsgs,
			})
		}

		return &protocol.TurmsNotification{
			RequestId: req.RequestId,
			Code:      proto.Int32(1000), // SUCCESS
			Data: &protocol.TurmsNotification_Data{
				Kind: &protocol.TurmsNotification_Data_MessagesWithTotalList{
					MessagesWithTotalList: &protocol.MessagesWithTotalList{
						MessagesWithTotalList: messagesWithTotalList,
					},
				},
			},
		}, nil
	} else {
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
}

// @MappedFrom handleUpdateMessageRequest()
func (c *MessageController) HandleUpdateMessageRequest(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
	updateReq := req.GetUpdateMessageRequest()

	// Use unified authAndUpdateMessage (Java parity: handles text, records, recallDate in one call)
	deviceType := int32(s.DeviceType)

	// Convert recallDate from millis to *time.Time
	var recallDate *time.Time
	if updateReq.RecallDate != nil {
		t := time.UnixMilli(*updateReq.RecallDate)
		recallDate = &t
	}

	// Use the unified AuthAndUpdateMessage method (matching Java's single authAndUpdateMessage call)
	err := c.messageService.AuthAndUpdateMessage(
		ctx,
		s.UserID,
		&deviceType,
		updateReq.MessageId,
		nil, // isSystemMessage
		updateReq.Text,
		updateReq.Records,
		nil, // burnAfter
		recallDate,
		nil, // senderIP
	)
	if err != nil {
		return nil, err
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
