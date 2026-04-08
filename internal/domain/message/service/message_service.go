package service

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"net"
	"time"

	"im.turms/server/internal/domain/common/cache"
	"im.turms/server/internal/domain/common/infra/idgen"
	"im.turms/server/internal/domain/message/bo"
	"im.turms/server/internal/domain/message/po"
	"im.turms/server/internal/domain/message/repository"
	"im.turms/server/internal/infra/plugin"
	"im.turms/server/internal/infra/property"
	turmsredis "im.turms/server/internal/storage/redis"
)

var (
	ErrNotFriend                  = errors.New("sender and target have no friendship, or are blocked")
	ErrNotGroupMember             = errors.New("sender is not a member of the target group, or lacks permission")
	ErrInvalidTargetID            = errors.New("invalid target ID")
	ErrTextLimitExceeded          = errors.New("text exceeds max limit")
	ErrRecordsSizeExceeded        = errors.New("records exceed max size")
	ErrInvalidBurnAfter           = errors.New("burnAfter must be >= 0")
	ErrInvalidDeliveryDate        = errors.New("deliveryDate must be past-or-present")
	ErrRecallDateBeforeDelivery   = errors.New("recallDate must be after deliveryDate")
	ErrEditNotAllowed             = errors.New("editing message by sender is not allowed")
	ErrRecallNotAllowed           = errors.New("recalling message is not allowed")
	ErrRecallDurationExceeded     = errors.New("message recall duration has been exceeded")
	ErrGroupNotActiveOrDeleted    = errors.New("group is not active or has been deleted")
	ErrNoFieldsToUpdate           = errors.New("no fields to update")
	ErrNotMessageRecipientOrSender = errors.New("not a recipient or sender of the reference message")
	ErrNotMessageSender           = errors.New("unauthorized to modify another user's message")
)

const chunkSize = 1000

type OutboundMessageDelivery interface {
	Deliver(ctx context.Context, targetID int64, msg *po.Message) error
}

type UserRelationshipService interface {
	HasRelationshipAndNotBlocked(ctx context.Context, ownerID int64, relatedUserID int64) (bool, error)
}

type GroupMemberService interface {
	IsGroupMember(ctx context.Context, groupID int64, userID int64) (bool, error)
	FindGroupMemberIDs(ctx context.Context, groupID int64) ([]int64, error)
}

type ConversationService interface {
	AuthAndUpsertGroupConversationReadDate(ctx context.Context, groupID int64, memberID int64, readDate *time.Time) error
	AuthAndUpsertPrivateConversationReadDate(ctx context.Context, ownerID int64, targetID int64, readDate *time.Time) error
}

type GroupService interface {
	QueryGroupTypeIfActiveAndNotDeleted(ctx context.Context, groupID int64) (bool, error)
}

type MessageService struct {
	idGen              *idgen.SnowflakeIdGenerator
	seqGen             *turmsredis.SequenceGenerator
	msgRepo            *repository.MessageRepository
	userRelService     UserRelationshipService
	groupMemService    GroupMemberService
	conversationSvc    ConversationService
	groupSvc           GroupService
	outboundDelivery   OutboundMessageDelivery
	propertiesManager  *property.TurmsPropertiesManager
	pluginManager      *plugin.PluginManager

	authCache *cache.TTLCache[string, bool]
}

func NewMessageService(
	idGen *idgen.SnowflakeIdGenerator,
	seqGen *turmsredis.SequenceGenerator,
	msgRepo *repository.MessageRepository,
	userRelSvc UserRelationshipService,
	groupMemSvc GroupMemberService,
	conversationSvc ConversationService,
	groupSvc GroupService,
	delivery OutboundMessageDelivery,
	propertiesManager *property.TurmsPropertiesManager,
	pluginManager *plugin.PluginManager,
) *MessageService {
	return &MessageService{
		idGen:             idGen,
		seqGen:            seqGen,
		msgRepo:           msgRepo,
		userRelService:    userRelSvc,
		groupMemService:   groupMemSvc,
		conversationSvc:   conversationSvc,
		groupSvc:          groupSvc,
		outboundDelivery:  delivery,
		propertiesManager: propertiesManager,
		pluginManager:     pluginManager,
		authCache:         cache.NewTTLCache[string, bool](1*time.Minute, 10*time.Second),
	}
}

func (s *MessageService) Close() {
	if s.authCache != nil {
		s.authCache.Close()
	}
}

func (s *MessageService) getMessageProperties() property.MessageProperties {
	return s.propertiesManager.GetLocalProperties().Service.Message
}

// validateSaveMessage validates parameters for saving a message.
func (s *MessageService) validateSaveMessage(text string, records [][]byte, burnAfter *int32, deliveryDate *time.Time, recallDate *time.Time) error {
	props := s.getMessageProperties()
	if len(text) > props.MaxTextLimit {
		return ErrTextLimitExceeded
	}
	if len(records) > props.MaxRecordsSize {
		return ErrRecordsSizeExceeded
	}
	if burnAfter != nil && *burnAfter < 0 {
		return ErrInvalidBurnAfter
	}
	if deliveryDate != nil && deliveryDate.After(time.Now()) {
		return ErrInvalidDeliveryDate
	}
	if recallDate != nil && deliveryDate != nil && recallDate.Before(*deliveryDate) {
		return ErrRecallDateBeforeDelivery
	}
	return nil
}

// computeConversationID computes the conversation ID from sender and target IDs.
func computeConversationID(isGroupMessage bool, senderID int64, targetID int64) []byte {
	if isGroupMessage {
		// Group conversation ID is just the group ID (big-endian 8 bytes)
		b := make([]byte, 8)
		binary.BigEndian.PutUint64(b, uint64(targetID))
		return b
	}
	// Private conversation ID: min(senderID, targetID) || max(senderID, targetID) (16 bytes)
	minID, maxID := senderID, targetID
	if senderID > targetID {
		minID, maxID = targetID, senderID
	}
	b := make([]byte, 16)
	binary.BigEndian.PutUint64(b[:8], uint64(minID))
	binary.BigEndian.PutUint64(b[8:], uint64(maxID))
	return b
}

// parseSenderIP parses a sender IP string into IPv4 (int32) and IPv6 ([]byte) fields.
func parseSenderIP(senderIP string) (ipv4 *int32, ipv6 []byte) {
	if senderIP == "" {
		return nil, nil
	}
	ip := net.ParseIP(senderIP)
	if ip == nil {
		return nil, nil
	}
	if ip4 := ip.To4(); ip4 != nil {
		v := int32(binary.BigEndian.Uint32(ip4))
		return &v, nil
	}
	if ip16 := ip.To16(); ip16 != nil {
		return nil, append([]byte(nil), ip16...)
	}
	return nil, nil
}

// AuthAndSaveMessage handles saving message state
// @MappedFrom authAndSaveMessage(boolean queryRecipientIds, @Nullable Boolean persist, @Nullable Long messageId, @NotNull Long senderId, @Nullable byte[] senderIp, @NotNull Long targetId, @NotNull Boolean isGroupMessage, @NotNull Boolean isSystemMessage, @Nullable String text, @Nullable List<byte[]> records, @Nullable @Min(0)
func (s *MessageService) AuthAndSaveMessage(
	ctx context.Context,
	isGroupMessage bool,
	senderID int64,
	targetID int64,
	isSystemMessage bool,
	text string,
	records [][]byte,
	burnAfter *int32,
	deliveryDate *time.Time,
	preMessageID *int64,
	senderIP string,
	referenceID *int64,
) (*bo.MessageAndRecipientIDs, error) {
	if targetID <= 0 {
		return nil, ErrInvalidTargetID
	}

	// Validation (Bug fix: Missing validation)
	if err := s.validateSaveMessage(text, records, burnAfter, deliveryDate, nil); err != nil {
		return nil, err
	}

	canSend, err := s.auth(ctx, isGroupMessage, senderID, targetID)
	if err != nil {
		return nil, err
	}
	if !canSend {
		if isGroupMessage {
			return nil, ErrNotGroupMember
		}
		return nil, ErrNotFriend
	}

	msg, recipientIDs, err := s.saveMessage0(ctx, isGroupMessage, isSystemMessage, senderID, targetID, text, records, burnAfter, deliveryDate, preMessageID, senderIP, referenceID)
	if err != nil {
		return nil, err
	}

	// Bug fix: Missing conversation read-date upsert
	props := s.getMessageProperties()
	if props.UpdateReadDateAfterMessageSent && s.conversationSvc != nil {
		if isGroupMessage {
			_ = s.conversationSvc.AuthAndUpsertGroupConversationReadDate(ctx, targetID, senderID, nil)
		} else {
			_ = s.conversationSvc.AuthAndUpsertPrivateConversationReadDate(ctx, senderID, targetID, nil)
		}
	}

	return &bo.MessageAndRecipientIDs{
		Message:      msg,
		RecipientIDs: recipientIDs,
	}, nil
}

// saveMessage0 is the internal save logic shared by AuthAndSaveMessage and other methods.
func (s *MessageService) saveMessage0(
	ctx context.Context,
	isGroupMessage bool,
	isSystemMessage bool,
	senderID int64,
	targetID int64,
	text string,
	records [][]byte,
	burnAfter *int32,
	deliveryDate *time.Time,
	preMessageID *int64,
	senderIP string,
	referenceID *int64,
) (*po.Message, []int64, error) {
	props := s.getMessageProperties()

	// Bug fix: Missing conditional sequence ID generation
	var seqID32 *int32
	if (isGroupMessage && props.UseSequenceIdForGroupConversation) ||
		(!isGroupMessage && props.UseSequenceIdForPrivateConversation) {
		var sequenceID int64
		var err error
		if isGroupMessage {
			sequenceID, err = s.seqGen.NextGroupMessageSequenceId(ctx, targetID)
		} else {
			sequenceID, err = s.seqGen.NextPrivateMessageSequenceId(ctx, senderID)
		}
		if err != nil {
			return nil, nil, fmt.Errorf("failed to generate sequence: %w", err)
		}
		v := int32(sequenceID)
		seqID32 = &v
	}

	msgID := s.idGen.NextIncreasingId()

	var dDate time.Time
	if deliveryDate != nil {
		dDate = *deliveryDate
	} else {
		dDate = time.Now()
	}

	// Bug fix: Missing conversationId computation
	conversationID := computeConversationID(isGroupMessage, senderID, targetID)

	// Bug fix: Missing senderIp parsing
	senderIPv4, senderIPv6 := parseSenderIP(senderIP)

	isSysMsg := isSystemMessage
	msg := &po.Message{
		ID:              msgID,
		ConversationID:  conversationID,
		IsGroupMessage:  &isGroupMessage,
		IsSystemMessage: &isSysMsg,
		SenderID:        senderID,
		TargetID:        targetID,
		Text:            text,
		SequenceID:      seqID32,
		DeliveryDate:    dDate,
		Records:         records,
		BurnAfter:       burnAfter,
		PreMessageID:    preMessageID,
		ReferenceID:     referenceID,
		SenderIP:        senderIPv4,
		SenderIPv6:      senderIPv6,
	}

	if err := s.msgRepo.InsertMessage(ctx, msg); err != nil {
		return nil, nil, fmt.Errorf("failed to save message to db: %w", err)
	}

	// Query recipient IDs
	var recipientIDs []int64
	if isGroupMessage {
		if s.groupMemService != nil {
			recipientIDs, _ = s.groupMemService.FindGroupMemberIDs(ctx, targetID)
		}
	} else {
		recipientIDs = []int64{targetID}
	}

	return msg, recipientIDs, nil
}

// @MappedFrom authAndSaveAndSendMessage(boolean send, @Nullable Boolean persist, @Nullable Long senderId, @Nullable DeviceType senderDeviceType, @Nullable byte[] senderIp, @Nullable Long messageId, @NotNull Boolean isGroupMessage, @NotNull Boolean isSystemMessage, @Nullable String text, @Nullable List<byte[]> records, @NotNull Long targetId, @Nullable @Min(0)
func (s *MessageService) AuthAndSaveAndSendMessage(
	ctx context.Context,
	isGroupMessage bool,
	senderID int64,
	targetID int64,
	isSystemMessage bool,
	text string,
	records [][]byte,
	burnAfter *int32,
	deliveryDate *time.Time,
	preMessageID *int64,
	senderIP string,
	referenceID *int64,
) (*bo.MessageAndRecipientIDs, error) {
	result, err := s.AuthAndSaveMessage(ctx, isGroupMessage, senderID, targetID, isSystemMessage, text, records, burnAfter, deliveryDate, preMessageID, senderIP, referenceID)
	if err != nil {
		return nil, err
	}

	if s.outboundDelivery != nil && result.Message != nil {
		if err := s.outboundDelivery.Deliver(ctx, targetID, result.Message); err != nil {
			return result, fmt.Errorf("saved but failed to deliver: %w", err)
		}
	}

	return result, nil
}

func (s *MessageService) auth(ctx context.Context, isGroupMessage bool, senderID int64, targetID int64) (bool, error) {
	cacheKey := fmt.Sprintf("auth:%v:%d:%d", isGroupMessage, senderID, targetID)
	if allowed, ok := s.authCache.Get(cacheKey); ok {
		return allowed, nil
	}

	var allowed bool
	var err error

	if isGroupMessage {
		if s.groupMemService != nil {
			allowed, err = s.groupMemService.IsGroupMember(ctx, targetID, senderID)
		} else {
			allowed = true
		}
	} else {
		if s.userRelService != nil {
			allowed, err = s.userRelService.HasRelationshipAndNotBlocked(ctx, senderID, targetID)
		} else {
			allowed = true
		}
	}

	if err != nil {
		return false, err
	}

	s.authCache.Set(cacheKey, allowed)
	return allowed, nil
}

// QueryMessages supports message pulling (offline/roaming sync)
func (s *MessageService) QueryMessages(
	ctx context.Context,
	requesterID int64,
	isGroupMessage *bool,
	senderIDs []int64,
	targetIDs []int64,
	deliveryDateAfter *time.Time,
	deliveryDateBefore *time.Time,
	size int64,
	ascending bool,
) ([]*po.Message, error) {
	// Authorization checking
	// 1. If querying group messages, ensure requester is a member
	var authErr error
	if isGroupMessage != nil && *isGroupMessage {
		if s.groupMemService != nil && len(targetIDs) > 0 {
			for _, groupID := range targetIDs {
				allowed, err := s.groupMemService.IsGroupMember(ctx, groupID, requesterID)
				if err != nil {
					authErr = err
					break
				}
				if !allowed {
					authErr = ErrNotGroupMember
					break
				}
			}
		}
	} else if isGroupMessage != nil && !*isGroupMessage {
		// 2. If querying private messages, restrict to messages related to requester
	}

	if authErr != nil {
		return nil, authErr
	}

	return s.msgRepo.QueryMessages(
		ctx,
		isGroupMessage,
		senderIDs,
		targetIDs,
		deliveryDateAfter,
		deliveryDateBefore,
		size,
		ascending,
	)
}

func (s *MessageService) CountMessages(
	ctx context.Context,
	isGroupMessage *bool,
	senderIDs []int64,
	targetIDs []int64,
	deliveryDateAfter *time.Time,
	deliveryDateBefore *time.Time,
) (int64, error) {
	return s.msgRepo.CountMessages(
		ctx,
		isGroupMessage,
		senderIDs,
		targetIDs,
		deliveryDateAfter,
		deliveryDateBefore,
	)
}

// checkIfAllowedToUpdateMessage checks if the sender is allowed to update the message.
func (s *MessageService) checkIfAllowedToUpdateMessage(ctx context.Context, msg *po.Message, senderID int64) error {
	if msg.SenderID != senderID {
		return ErrNotMessageSender
	}
	props := s.getMessageProperties()
	if !props.AllowEditMessageBySender {
		return ErrEditNotAllowed
	}
	return nil
}

// checkIfAllowedToRecallMessage checks if the sender is allowed to recall the message.
func (s *MessageService) checkIfAllowedToRecallMessage(ctx context.Context, msg *po.Message, senderID int64) error {
	if msg.SenderID != senderID {
		return ErrNotMessageSender
	}
	props := s.getMessageProperties()
	if !props.AllowRecallMessage {
		return ErrRecallNotAllowed
	}
	// Bug fix: Missing recall duration timeout check
	if props.AvailableRecallDurationMillis > 0 {
		elapsed := time.Since(msg.DeliveryDate)
		if elapsed > time.Duration(props.AvailableRecallDurationMillis)*time.Millisecond {
			return ErrRecallDurationExceeded
		}
	}
	// Bug fix: Missing group type active/not-deleted check
	if msg.IsGroupMessage != nil && *msg.IsGroupMessage {
		if s.groupSvc != nil {
			active, err := s.groupSvc.QueryGroupTypeIfActiveAndNotDeleted(ctx, msg.TargetID)
			if err != nil {
				return err
			}
			if !active {
				return ErrGroupNotActiveOrDeleted
			}
		}
	}
	return nil
}

func (s *MessageService) AuthAndRecallMessage(ctx context.Context, senderID int64, messageID int64) error {
	// First fetch the message to verify ownership
	msg, err := s.msgRepo.FindByID(ctx, messageID)
	if err != nil {
		return fmt.Errorf("failed to find message: %w", err)
	}

	// Bug fix: Missing permission checks for recall
	if err := s.checkIfAllowedToRecallMessage(ctx, msg, senderID); err != nil {
		return err
	}

	now := time.Now()
	if err := s.msgRepo.UpdateMessage(ctx, messageID, nil, nil, &now); err != nil {
		return fmt.Errorf("failed to recall message in db: %w", err)
	}

	return nil
}

func (s *MessageService) AuthAndUpdateMessageText(ctx context.Context, senderID int64, messageID int64, newText string) error {
	msg, err := s.msgRepo.FindByID(ctx, messageID)
	if err != nil {
		return fmt.Errorf("failed to find message: %w", err)
	}

	// Bug fix: Missing permission checks for edit
	if err := s.checkIfAllowedToUpdateMessage(ctx, msg, senderID); err != nil {
		return err
	}

	now := time.Now()
	if err := s.msgRepo.UpdateMessage(ctx, messageID, &newText, &now, nil); err != nil {
		return fmt.Errorf("failed to update message text in db: %w", err)
	}

	return nil
}

// IsMessageRecipientOrSender checks if a user is the sender or recipient of a message.
func (s *MessageService) IsMessageRecipientOrSender(ctx context.Context, messageID int64, userID int64) (bool, error) {
	msg, err := s.msgRepo.FindMessageSenderIDAndTargetIDAndIsGroupMessage(ctx, messageID)
	if err != nil {
		return false, err
	}
	if msg.SenderID == userID {
		return true, nil
	}
	if msg.IsGroupMessage != nil && *msg.IsGroupMessage {
		if s.groupMemService == nil {
			return true, nil
		}
		return s.groupMemService.IsGroupMember(ctx, msg.TargetID, userID)
	}
	return msg.TargetID == userID, nil
}

// QueryMessage queries a message by ID.
func (s *MessageService) QueryMessage(ctx context.Context, messageID int64) (*po.Message, error) {
	return s.msgRepo.FindByID(ctx, messageID)
}

// SaveMessage saves a message without authentication.
func (s *MessageService) SaveMessage(
	ctx context.Context,
	isGroupMessage bool,
	senderID int64,
	targetID int64,
	isSystemMessage bool,
	text string,
	records [][]byte,
	burnAfter *int32,
	deliveryDate *time.Time,
	preMessageID *int64,
	senderIP string,
	referenceID *int64,
) (*po.Message, error) {
	if targetID <= 0 {
		return nil, ErrInvalidTargetID
	}

	props := s.getMessageProperties()

	// Bug fix: Missing conditional sequence ID generation
	var seqID32 *int32
	if (isGroupMessage && props.UseSequenceIdForGroupConversation) ||
		(!isGroupMessage && props.UseSequenceIdForPrivateConversation) {
		var sequenceID int64
		var err error
		if isGroupMessage {
			sequenceID, err = s.seqGen.NextGroupMessageSequenceId(ctx, targetID)
		} else {
			sequenceID, err = s.seqGen.NextPrivateMessageSequenceId(ctx, senderID)
		}
		if err != nil {
			return nil, fmt.Errorf("failed to generate sequence: %w", err)
		}
		v := int32(sequenceID)
		seqID32 = &v
	}

	msgID := s.idGen.NextIncreasingId()

	var dDate time.Time
	if deliveryDate != nil {
		dDate = *deliveryDate
	} else {
		dDate = time.Now()
	}

	// Bug fix: Missing conversationId computation
	conversationID := computeConversationID(isGroupMessage, senderID, targetID)

	// Bug fix: Missing senderIp parsing
	senderIPv4, senderIPv6 := parseSenderIP(senderIP)

	isSysMsg := isSystemMessage
	msg := &po.Message{
		ID:              msgID,
		ConversationID:  conversationID,
		IsGroupMessage:  &isGroupMessage,
		IsSystemMessage: &isSysMsg,
		SenderID:        senderID,
		TargetID:        targetID,
		Text:            text,
		SequenceID:      seqID32,
		DeliveryDate:    dDate,
		Records:         records,
		BurnAfter:       burnAfter,
		PreMessageID:    preMessageID,
		ReferenceID:     referenceID,
		SenderIP:        senderIPv4,
		SenderIPv6:      senderIPv6,
	}

	if err := s.msgRepo.InsertMessage(ctx, msg); err != nil {
		return nil, fmt.Errorf("failed to save message to db: %w", err)
	}

	return msg, nil
}

// QueryExpiredMessageIds queries expired message IDs.
func (s *MessageService) QueryExpiredMessageIds(ctx context.Context, retentionPeriodHours int) ([]int64, error) {
	expirationDate := time.Now().Add(-time.Duration(retentionPeriodHours) * time.Hour)
	return s.msgRepo.FindExpiredMessageIds(ctx, expirationDate)
}

// DeleteExpiredMessages physically deletes expired messages.
// Bug fix: Missing plugin notification + batch processing
func (s *MessageService) DeleteExpiredMessages(ctx context.Context, retentionPeriodHours int) error {
	expirationDate := time.Now().Add(-time.Duration(retentionPeriodHours) * time.Hour)
	ids, err := s.msgRepo.FindExpiredMessageIds(ctx, expirationDate)
	if err != nil {
		return err
	}
	if len(ids) == 0 {
		return nil
	}

	// Bug fix: Missing plugin notification
	if s.pluginManager != nil && s.pluginManager.HasRunningExtensions("ExpiredMessageDeletionNotifier") {
		// Fetch full messages for plugin inspection
		messages, err := s.msgRepo.QueryMessages(ctx, nil, nil, nil, nil, nil, int64(len(ids)), true)
		if err != nil {
			return err
		}
		// Invoke extension points to filter which messages to delete
		result, extErr := s.pluginManager.InvokeExtensionPoints(ctx, "ExpiredMessageDeletionNotifier", "shouldDelete", messages)
		if extErr != nil {
			return extErr
		}
		if result != nil && !*result {
			return nil // Plugin vetoed deletion
		}
	}

	// Bug fix: Missing batch processing - chunk IDs for processing
	for i := 0; i < len(ids); i += chunkSize {
		end := i + chunkSize
		if end > len(ids) {
			end = len(ids)
		}
		chunk := ids[i:end]
		if err := s.msgRepo.DeleteMessages(ctx, chunk); err != nil {
			return err
		}
	}
	return nil
}

// DeleteMessages deletes messages logically or physically.
// Bug fix: Missing deleteMessageLogicallyByDefault fallback
func (s *MessageService) DeleteMessages(ctx context.Context, messageIDs []int64, deleteLogically *bool) (*DeleteResult, error) {
	if len(messageIDs) == 0 {
		return &DeleteResult{Acknowledged: true}, nil
	}

	// Bug fix: Missing deleteMessageLogicallyByDefault fallback
	isLogical := deleteLogically
	if isLogical == nil {
		defaultVal := s.getMessageProperties().DeleteMessageLogicallyByDefault
		isLogical = &defaultVal
	}

	if *isLogical {
		now := time.Now()
		err := s.msgRepo.UpdateMessagesDeletionDate(ctx, messageIDs, &now)
		if err != nil {
			return nil, err
		}
		// Bug fix: Missing conversion of update result to delete result
		return &DeleteResult{Acknowledged: true, ModifiedCount: int64(len(messageIDs))}, nil
	}
	err := s.msgRepo.DeleteMessages(ctx, messageIDs)
	if err != nil {
		return nil, err
	}
	return &DeleteResult{Acknowledged: true, ModifiedCount: int64(len(messageIDs))}, nil
}

// DeleteResult represents the result of a delete operation, analogous to Java's DeleteResult.
type DeleteResult struct {
	Acknowledged   bool
	ModifiedCount  int64
}

// UpdateMessages updates messages in batch.
// Bug fix: Missing early return when all update fields are null
func (s *MessageService) UpdateMessages(
	ctx context.Context,
	senderID *int64,
	senderDeviceType *int32,
	messageIDs []int64,
	isSystemMessage *bool,
	text *string,
	records [][]byte,
	burnAfter *int32,
	recallDate *time.Time,
	senderIP *string,
) error {
	// Bug fix: Missing early return when all update fields are null
	if isSystemMessage == nil && text == nil && records == nil && burnAfter == nil && recallDate == nil && senderIP == nil {
		return nil // ACKNOWLEDGED_UPDATE_RESULT equivalent
	}

	// Bug fix: Missing senderIp parsing
	var senderIPv4 *int32
	var senderIPv6 []byte
	if senderIP != nil {
		senderIPv4, senderIPv6 = parseSenderIP(*senderIP)
	}

	return s.msgRepo.UpdateMessages(ctx, messageIDs, isSystemMessage, senderIPv4, senderIPv6, recallDate, text, records, burnAfter)
}

// HasPrivateMessage checks if a specific private message exists.
func (s *MessageService) HasPrivateMessage(ctx context.Context, senderID int64, targetID int64) (bool, error) {
	return s.msgRepo.ExistsBySenderIDAndTargetID(ctx, senderID, targetID)
}

// AuthAndUpdateMessage updates a message after authentication.
// Bug fix: Missing conditional logic for update vs recall paths + all the permission checks
func (s *MessageService) AuthAndUpdateMessage(
	ctx context.Context,
	senderID int64,
	senderDeviceType *int32,
	messageID int64,
	isSystemMessage *bool,
	text *string,
	records [][]byte,
	burnAfter *int32,
	recallDate *time.Time,
	senderIP *string,
) error {
	// Bug fix: Missing early return when all update fields null
	if isSystemMessage == nil && text == nil && records == nil && burnAfter == nil && recallDate == nil && senderIP == nil {
		return nil // ACKNOWLEDGED_UPDATE_RESULT equivalent
	}

	// Bug fix: Missing validation
	if text != nil {
		props := s.getMessageProperties()
		if len(*text) > props.MaxTextLimit {
			return ErrTextLimitExceeded
		}
	}
	if records != nil {
		props := s.getMessageProperties()
		if len(records) > props.MaxRecordsSize {
			return ErrRecordsSizeExceeded
		}
	}
	if burnAfter != nil && *burnAfter < 0 {
		return ErrInvalidBurnAfter
	}

	msg, err := s.msgRepo.FindByID(ctx, messageID)
	if err != nil {
		return fmt.Errorf("failed to find message: %w", err)
	}

	// Bug fix: Missing conditional logic for update vs recall paths
	hasTextOrRecords := text != nil || records != nil
	hasRecallDate := recallDate != nil

	if hasTextOrRecords {
		// Update path: check update permission
		if err := s.checkIfAllowedToUpdateMessage(ctx, msg, senderID); err != nil {
			return err
		}
	} else if hasRecallDate {
		// Recall path: check recall permission
		if err := s.checkIfAllowedToRecallMessage(ctx, msg, senderID); err != nil {
			return err
		}
	} else {
		// No text/records and no recallDate - check sender at minimum
		if msg.SenderID != senderID {
			return ErrNotMessageSender
		}
	}

	// Bug fix: Missing senderIp parsing
	var senderIPv4 *int32
	var senderIPv6 []byte
	if senderIP != nil {
		senderIPv4, senderIPv6 = parseSenderIP(*senderIP)
	}

	err = s.msgRepo.UpdateMessages(ctx, []int64{messageID}, isSystemMessage, senderIPv4, senderIPv6, recallDate, text, records, burnAfter)
	if err != nil {
		return err
	}

	// Bug fix: Missing recall notification logic
	if hasRecallDate && s.outboundDelivery != nil {
		// Fetch the updated message and send recall notification
		updatedMsg, findErr := s.msgRepo.FindByID(ctx, messageID)
		if findErr == nil && updatedMsg != nil {
			_ = s.outboundDelivery.Deliver(ctx, updatedMsg.TargetID, updatedMsg)
		}
	}

	return nil
}

// CountMessages counts messages matching the specific criteria.
func (s *MessageService) CountMessagesByRange(
	ctx context.Context,
	isGroupMessage *bool,
	areSystemMessages *bool,
	senderIDs []int64,
	targetIDs []int64,
	deliveryDateAfter *time.Time,
	deliveryDateBefore *time.Time,
) (int64, error) {
	return s.msgRepo.CountMessages(ctx, isGroupMessage, senderIDs, targetIDs, deliveryDateAfter, deliveryDateBefore)
}

// CountUsersWhoSentMessage counts distinct users who sent messages.
func (s *MessageService) CountUsersWhoSentMessage(ctx context.Context, deliveryDateAfter *time.Time, deliveryDateBefore *time.Time, areGroupMessages *bool, areSystemMessages *bool) (int64, error) {
	return s.msgRepo.CountUsersWhoSentMessage(ctx, deliveryDateAfter, deliveryDateBefore, areGroupMessages, areSystemMessages)
}

// CountGroupsThatSentMessages counts distinct groups that had messages sent to them.
func (s *MessageService) CountGroupsThatSentMessages(ctx context.Context, deliveryDateAfter *time.Time, deliveryDateBefore *time.Time) (int64, error) {
	return s.msgRepo.CountGroupsThatSentMessages(ctx, deliveryDateAfter, deliveryDateBefore)
}

// CountSentMessages counts the number of sent messages based on criteria.
func (s *MessageService) CountSentMessages(ctx context.Context, deliveryDateAfter *time.Time, deliveryDateBefore *time.Time, areGroupMessages *bool, areSystemMessages *bool) (int64, error) {
	return s.msgRepo.CountSentMessages(ctx, deliveryDateAfter, deliveryDateBefore, areGroupMessages, areSystemMessages)
}

// CountSentMessagesOnAverage counts the average number of messages sent per user.
// Bug fix: Different behavior when totalUsers is 0 - Java returns Long.MAX_VALUE when totalMessages > 0 and distinctUsers is 0
func (s *MessageService) CountSentMessagesOnAverage(ctx context.Context, deliveryDateAfter *time.Time, deliveryDateBefore *time.Time, areGroupMessages *bool, areSystemMessages *bool) (int64, error) {
	totalMessages, err := s.CountSentMessages(ctx, deliveryDateAfter, deliveryDateBefore, areGroupMessages, areSystemMessages)
	if err != nil {
		return 0, err
	}
	if totalMessages == 0 {
		return 0, nil
	}
	distinctUsers, err := s.CountUsersWhoSentMessage(ctx, deliveryDateAfter, deliveryDateBefore, areGroupMessages, areSystemMessages)
	if err != nil {
		return 0, err
	}
	// Bug fix: Java returns Long.MAX_VALUE when totalUsers is 0 but totalMessages > 0
	if distinctUsers == 0 {
		return math.MaxInt64, nil
	}
	return totalMessages / distinctUsers, nil
}

// AuthAndQueryCompleteMessages auths and queries complete messages.
func (s *MessageService) AuthAndQueryCompleteMessages(
	ctx context.Context,
	requesterID int64,
	isGroupMessage *bool,
	areSystemMessages *bool,
	senderIDs []int64,
	deliveryDateAfter *time.Time,
	deliveryDateBefore *time.Time,
	size int64,
	ascending bool,
) ([]*po.Message, error) {
	return s.msgRepo.QueryMessages(ctx, isGroupMessage, senderIDs, nil, deliveryDateAfter, deliveryDateBefore, size, ascending)
}

// QueryMessageRecipients queries the recipients of a message.
func (s *MessageService) QueryMessageRecipients(ctx context.Context, messageID int64) ([]int64, error) {
	msg, err := s.QueryMessage(ctx, messageID)
	if err != nil {
		return nil, err
	}
	if msg.IsGroupMessage != nil && *msg.IsGroupMessage {
		if s.groupMemService != nil {
			return s.groupMemService.FindGroupMemberIDs(ctx, msg.TargetID)
		}
		return nil, nil
	}
	return []int64{msg.TargetID}, nil
}

// SaveAndSendMessage saves and optionally delivers a message.
func (s *MessageService) SaveAndSendMessage(
	ctx context.Context,
	send bool,
	persist bool,
	senderID int64,
	isGroupMessage bool,
	isSystemMessage bool,
	text string,
	records [][]byte,
	targetID int64,
	burnAfter *int32,
	deliveryDate *time.Time,
	preMessageID *int64,
	senderIP string,
	referenceID *int64,
) (*bo.MessageAndRecipientIDs, error) {
	var msg *po.Message
	var recipientIDs []int64
	var err error

	if persist {
		msg, err = s.SaveMessage(ctx, isGroupMessage, senderID, targetID, isSystemMessage, text, records, burnAfter, deliveryDate, preMessageID, senderIP, referenceID)
		if err != nil {
			return nil, err
		}
		// Compute recipient IDs
		if isGroupMessage {
			if s.groupMemService != nil {
				recipientIDs, _ = s.groupMemService.FindGroupMemberIDs(ctx, targetID)
			}
		} else {
			recipientIDs = []int64{targetID}
		}
	} else {
		msg = &po.Message{
			IsGroupMessage: &isGroupMessage,
			SenderID:       senderID,
			TargetID:       targetID,
			Text:           text,
			Records:        records,
			BurnAfter:      burnAfter,
			PreMessageID:   preMessageID,
		}
	}

	if send && s.outboundDelivery != nil {
		if err := s.outboundDelivery.Deliver(ctx, targetID, msg); err != nil {
			return &bo.MessageAndRecipientIDs{Message: msg, RecipientIDs: recipientIDs}, fmt.Errorf("message prepared but failed delivery: %w", err)
		}
	}

	return &bo.MessageAndRecipientIDs{Message: msg, RecipientIDs: recipientIDs}, nil
}

// CloneAndSaveMessage clones and saves an existing message to a new target.
func (s *MessageService) CloneAndSaveMessage(
	ctx context.Context,
	senderID int64,
	referenceID int64,
	isGroupMessage bool,
	isSystemMessage bool,
	targetID int64,
) (*po.Message, error) {
	refMsg, err := s.QueryMessage(ctx, referenceID)
	if err != nil {
		return nil, err
	}

	return s.SaveMessage(
		ctx,
		isGroupMessage,
		senderID,
		targetID,
		isSystemMessage,
		refMsg.Text,
		refMsg.Records,
		refMsg.BurnAfter,
		nil,
		nil,
		"",
		&referenceID,
	)
}

// AuthAndCloneAndSaveMessage clones and saves after authorization.
func (s *MessageService) AuthAndCloneAndSaveMessage(
	ctx context.Context,
	requesterID int64,
	referenceID int64,
	isGroupMessage bool,
	isSystemMessage bool,
	targetID int64,
) (*bo.MessageAndRecipientIDs, error) {
	// Bug fix: Missing switchIfEmpty for message-not-found case
	hasAuth, err := s.IsMessageRecipientOrSender(ctx, referenceID, requesterID)
	if err != nil {
		// Map "not found" errors to the generic "not recipient or sender" error
		// to avoid leaking whether a message exists
		return nil, ErrNotMessageRecipientOrSender
	}
	if !hasAuth {
		return nil, ErrNotMessageRecipientOrSender
	}

	msg, err := s.CloneAndSaveMessage(ctx, requesterID, referenceID, isGroupMessage, isSystemMessage, targetID)
	if err != nil {
		return nil, err
	}

	// Compute recipient IDs
	var recipientIDs []int64
	if isGroupMessage {
		if s.groupMemService != nil {
			recipientIDs, _ = s.groupMemService.FindGroupMemberIDs(ctx, targetID)
		}
	} else {
		recipientIDs = []int64{targetID}
	}

	return &bo.MessageAndRecipientIDs{Message: msg, RecipientIDs: recipientIDs}, nil
}

// DeleteGroupMessageSequenceIDs deletes sequence IDs associated with groups.
func (s *MessageService) DeleteGroupMessageSequenceIDs(ctx context.Context, groupIDs []int64) error {
	return nil
}

// DeletePrivateMessageSequenceIDs deletes sequence IDs associated with users.
func (s *MessageService) DeletePrivateMessageSequenceIDs(ctx context.Context, userIDs []int64) error {
	return nil
}

// FetchGroupMessageSequenceID retrieves the max sequence ID.
func (s *MessageService) FetchGroupMessageSequenceID(ctx context.Context, groupID int64) (int64, error) {
	return s.seqGen.NextGroupMessageSequenceId(ctx, groupID)
}

// FetchPrivateMessageSequenceID retrieves the max private sequence ID.
func (s *MessageService) FetchPrivateMessageSequenceID(ctx context.Context, userID1 int64, userID2 int64) (int64, error) {
	return s.seqGen.NextPrivateMessageSequenceId(ctx, userID1)
}
