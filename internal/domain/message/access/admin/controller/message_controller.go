package controller

import (
	"context"
	"time"

	"im.turms/server/internal/domain/common/access/admin/controller"
	"im.turms/server/internal/domain/common/access/admin/dto"
	"im.turms/server/internal/domain/common/access/admin/dto/response"
	message_dto "im.turms/server/internal/domain/message/access/admin/dto"
	"im.turms/server/internal/domain/message/service"
	turmsmongo "im.turms/server/internal/storage/mongo"
	timeutil "im.turms/server/internal/infra/time"
)

// MessageController maps to MessageController.java
// @MappedFrom MessageController
type MessageController struct {
	*controller.BaseController
	messageService *service.MessageService
}

func NewMessageController(base *controller.BaseController, messageService *service.MessageService) *MessageController {
	return &MessageController{
		BaseController: base,
		messageService: messageService,
	}
}

// @MappedFrom createMessages(@QueryParam(defaultValue = "true") Boolean send, @RequestBody CreateMessageDTO createMessageDTO)
func (c *MessageController) CreateMessages(ctx context.Context, send *bool, createDTO *message_dto.CreateMessageDTO) (*service.SaveResult, error) {
	isGroupMessage := false
	if createDTO.IsGroupMessage != nil {
		isGroupMessage = *createDTO.IsGroupMessage
	}
	isSystemMessage := false
	if createDTO.IsSystemMessage != nil {
		isSystemMessage = *createDTO.IsSystemMessage
	}

	var senderID int64
	if createDTO.SenderID != nil {
		senderID = *createDTO.SenderID
	}

	var targetID int64
	if createDTO.TargetID != nil {
		targetID = *createDTO.TargetID
	}

	var text string
	if createDTO.Text != nil {
		text = *createDTO.Text
	}

	var records [][]byte
	if createDTO.Records != nil {
		records = createDTO.Records
	}

	var burnAfter *int32
	if createDTO.BurnAfter != nil {
		ba := int32(*createDTO.BurnAfter)
		burnAfter = &ba
	}

	var preMessageID *int64
	if createDTO.PreMessageID != nil {
		preMessageID = createDTO.PreMessageID
	}

	var deliveryDate *time.Time
	if createDTO.PreMessageID != nil {
		// Use current time for delivery date
		now := time.Now()
		deliveryDate = &now
	}

	persist := true
	shouldSend := true
	if send != nil {
		shouldSend = *send
	}

	var msg *service.SaveResult
	if shouldSend {
		savedMsg, err := c.messageService.SaveAndSendMessage(ctx, shouldSend, persist, senderID, isGroupMessage, isSystemMessage, text, records, targetID, burnAfter, deliveryDate, preMessageID)
		if err != nil {
			return nil, err
		}
		msg = &service.SaveResult{Message: savedMsg}
	} else {
		savedMsg, err := c.messageService.SaveMessage(ctx, isGroupMessage, senderID, targetID, text, records, burnAfter, deliveryDate, preMessageID)
		if err != nil {
			return nil, err
		}
		msg = &service.SaveResult{Message: savedMsg}
	}

	_ = isSystemMessage
	return msg, nil
}

// @MappedFrom queryMessages(@QueryParam(required = false) Set<Long> ids, @QueryParam(required = false) Boolean areGroupMessages, @QueryParam(required = false) Boolean areSystemMessages, @QueryParam(required = false) Set<Long> senderIds, @QueryParam(required = false) Set<Long> targetIds, @QueryParam(required = false) Date deliveryDateStart, @QueryParam(required = false) Date deliveryDateEnd, @QueryParam(required = false) Date deletionDateStart, @QueryParam(required = false) Date deletionDateEnd, @QueryParam(required = false) Date recallDateStart, @QueryParam(required = false) Date recallDateEnd, @QueryParam(required = false) Integer size, @QueryParam(required = false) Boolean ascending)
func (c *MessageController) QueryMessages(
	ctx context.Context,
	ids []int64,
	areGroupMessages *bool,
	areSystemMessages *bool,
	senderIDs []int64,
	targetIDs []int64,
	deliveryDateStart *time.Time,
	deliveryDateEnd *time.Time,
	deletionDateStart *time.Time,
	deletionDateEnd *time.Time,
	recallDateStart *time.Time,
	recallDateEnd *time.Time,
	size *int,
	ascending *bool,
) ([]*service.QueryMessageResult, error) {
	actualSize := c.GetPageSize(size)

	var deliveryDateRange *turmsmongo.DateRange
	if deliveryDateStart != nil || deliveryDateEnd != nil {
		deliveryDateRange = &turmsmongo.DateRange{Start: deliveryDateStart, End: deliveryDateEnd}
	}

	var deletionDateRange *turmsmongo.DateRange
	if deletionDateStart != nil || deletionDateEnd != nil {
		deletionDateRange = &turmsmongo.DateRange{Start: deletionDateStart, End: deletionDateEnd}
	}

	var recallDateRange *turmsmongo.DateRange
	if recallDateStart != nil || recallDateEnd != nil {
		recallDateRange = &turmsmongo.DateRange{Start: recallDateStart, End: recallDateEnd}
	}

	asc := false
	if ascending != nil {
		asc = *ascending
	}

	_ = actualSize
	_ = deliveryDateRange
	_ = deletionDateRange
	_ = recallDateRange
	_ = areSystemMessages

	// Delegate to the message service's query
	messages, err := c.messageService.QueryMessagesForAdmin(
		ctx,
		ids,
		areGroupMessages,
		areSystemMessages,
		senderIDs,
		targetIDs,
		deliveryDateRange,
		deletionDateRange,
		recallDateRange,
		nil,
		&actualSize,
		asc,
	)
	if err != nil {
		return nil, err
	}

	var results []*service.QueryMessageResult
	for _, m := range messages {
		results = append(results, &service.QueryMessageResult{Message: m})
	}
	return results, nil
}

// @MappedFrom queryMessages(@QueryParam(required = false) Set<Long> ids, @QueryParam(required = false) Boolean areGroupMessages, @QueryParam(required = false) Boolean areSystemMessages, @QueryParam(required = false) Set<Long> senderIds, @QueryParam(required = false) Set<Long> targetIds, @QueryParam(required = false) Date deliveryDateStart, @QueryParam(required = false) Date deliveryDateEnd, @QueryParam(required = false) Date deletionDateStart, @QueryParam(required = false) Date deletionDateEnd, @QueryParam(required = false) Date recallDateStart, @QueryParam(required = false) Date recallDateEnd, @QueryParam(required = false) Integer page, @QueryParam(required = false) Integer size)
func (c *MessageController) QueryMessagesByPage(
	ctx context.Context,
	ids []int64,
	areGroupMessages *bool,
	areSystemMessages *bool,
	senderIDs []int64,
	targetIDs []int64,
	deliveryDateStart *time.Time,
	deliveryDateEnd *time.Time,
	deletionDateStart *time.Time,
	deletionDateEnd *time.Time,
	recallDateStart *time.Time,
	recallDateEnd *time.Time,
	page *int,
	size *int,
) (*PaginationResponse, error) {
	actualSize := c.GetPageSize(size)

	var deliveryDateRange *turmsmongo.DateRange
	if deliveryDateStart != nil || deliveryDateEnd != nil {
		deliveryDateRange = &turmsmongo.DateRange{Start: deliveryDateStart, End: deliveryDateEnd}
	}

	var deletionDateRange *turmsmongo.DateRange
	if deletionDateStart != nil || deletionDateEnd != nil {
		deletionDateRange = &turmsmongo.DateRange{Start: deletionDateStart, End: deletionDateEnd}
	}

	var recallDateRange *turmsmongo.DateRange
	if recallDateStart != nil || recallDateEnd != nil {
		recallDateRange = &turmsmongo.DateRange{Start: recallDateStart, End: recallDateEnd}
	}

	// Count messages (note: count does NOT include recallDateRange per Java spec)
	total, err := c.messageService.CountMessagesForAdmin(
		ctx,
		ids,
		areGroupMessages,
		areSystemMessages,
		senderIDs,
		targetIDs,
		deliveryDateRange,
		deletionDateRange,
	)
	if err != nil {
		return nil, err
	}

	// Query messages (includes recallDateRange per Java spec)
	messages, err := c.messageService.QueryMessagesForAdmin(
		ctx,
		ids,
		areGroupMessages,
		areSystemMessages,
		senderIDs,
		targetIDs,
		deliveryDateRange,
		deletionDateRange,
		recallDateRange,
		page,
		&actualSize,
		true,
	)
	if err != nil {
		return nil, err
	}

	return &PaginationResponse{Total: total, Records: messages}, nil
}

// @MappedFrom countMessages(@QueryParam(required = false) Set<Long> ids, @QueryParam(required = false) Boolean areGroupMessages, @QueryParam(required = false) Boolean areSystemMessages, @QueryParam(required = false) Set<Long> senderIds, @QueryParam(required = false) Set<Long> targetIds, @QueryParam(required = false) Date deliveryDateStart, @QueryParam(required = false) Date deliveryDateEnd, @QueryParam(required = false) Date deletionDateStart, @QueryParam(required = false) Date deletionDateEnd, @QueryParam(required = false) Integer page, @QueryParam(required = false) Integer size)
func (c *MessageController) CountMessages(
	ctx context.Context,
	ids []int64,
	areGroupMessages *bool,
	areSystemMessages *bool,
	senderIDs []int64,
	targetIDs []int64,
	deliveryDateStart *time.Time,
	deliveryDateEnd *time.Time,
	deletionDateStart *time.Time,
	deletionDateEnd *time.Time,
) (int64, error) {
	var deliveryDateRange *turmsmongo.DateRange
	if deliveryDateStart != nil || deliveryDateEnd != nil {
		deliveryDateRange = &turmsmongo.DateRange{Start: deliveryDateStart, End: deliveryDateEnd}
	}

	var deletionDateRange *turmsmongo.DateRange
	if deletionDateStart != nil || deletionDateEnd != nil {
		deletionDateRange = &turmsmongo.DateRange{Start: deletionDateStart, End: deletionDateEnd}
	}

	return c.messageService.CountMessagesForAdmin(
		ctx,
		ids,
		areGroupMessages,
		areSystemMessages,
		senderIDs,
		targetIDs,
		deliveryDateRange,
		deletionDateRange,
	)
}

// @MappedFrom queryMessageStatistics(@QueryParam(required = false) Date startDate, @QueryParam(required = false) Date endDate, @QueryParam(required = false) DivideBy divideBy, @QueryParam(required = false) Boolean areGroupMessages, @QueryParam(required = false) Boolean areSystemMessages)
func (c *MessageController) QueryMessageStatistics(
	ctx context.Context,
	startDate *time.Time,
	endDate *time.Time,
	divideBy *timeutil.DivideBy,
	areGroupMessages *bool,
	areSystemMessages *bool,
) (*message_dto.MessageStatisticsDTO, error) {
	if divideBy == nil {
		noop := timeutil.DivideBy_NOOP
		divideBy = &noop
	}

	var startVal, endVal time.Time
	if startDate != nil {
		startVal = *startDate
	}
	if endDate != nil {
		endVal = *endDate
	}
	dateRange := timeutil.DateRange{Start: startVal, End: endVal}

	sentMessagesOnAverage, err := c.messageService.CountSentMessagesOnAverage(ctx, startDate, endDate, areGroupMessages, areSystemMessages)
	if err != nil {
		return nil, err
	}

	sentMessages, err := c.messageService.CountSentMessages(ctx, startDate, endDate, areGroupMessages, areSystemMessages)
	if err != nil {
		return nil, err
	}

	// Build statistics DTO
	stats := &message_dto.MessageStatisticsDTO{
		SentMessagesOnAverage: &sentMessagesOnAverage,
		SentMessages:          &sentMessages,
	}

	// If divided query is needed (not NOOP), use CheckAndQueryBetweenDate
	if *divideBy != timeutil.DivideBy_NOOP {
		sentMessagesRecords, err := c.CheckAndQueryBetweenDate(
			dateRange, *divideBy,
			func(dr timeutil.DateRange, agm *bool, asm *bool) (int64, error) {
				return c.messageService.CountSentMessages(ctx, &dr.Start, &dr.End, agm, asm)
			},
			areGroupMessages,
			areSystemMessages,
		)
		if err != nil {
			return nil, err
		}
		var sentRecords []interface{}
		for _, r := range sentMessagesRecords {
			sentRecords = append(sentRecords, r)
		}
		stats.SentMessagesRecords = sentRecords

		avgRecords, err := c.CheckAndQueryBetweenDate(
			dateRange, *divideBy,
			func(dr timeutil.DateRange, agm *bool, asm *bool) (int64, error) {
				return c.messageService.CountSentMessagesOnAverage(ctx, &dr.Start, &dr.End, agm, asm)
			},
			areGroupMessages,
			areSystemMessages,
		)
		if err != nil {
			return nil, err
		}
		var avgRecs []interface{}
		for _, r := range avgRecords {
			avgRecs = append(avgRecs, r)
		}
		stats.SentMessagesOnAverageRecords = avgRecs
	}

	return stats, nil
}

// @MappedFrom updateMessages(Set<Long> ids, @RequestBody UpdateMessageDTO updateMessageDTO)
func (c *MessageController) UpdateMessages(
	ctx context.Context,
	ids []int64,
	updateDTO *message_dto.UpdateMessageDTO,
) (*response.UpdateResultDTO, error) {
	var burnAfter *int32
	if updateDTO.BurnAfter != nil {
		ba := int32(*updateDTO.BurnAfter)
		burnAfter = &ba
	}
	err := c.messageService.UpdateMessages(
		ctx,
		updateDTO.SenderID,
		nil, // senderDeviceType not available in admin context
		ids,
		updateDTO.IsSystemMessage,
		updateDTO.Text,
		updateDTO.Records,
		burnAfter,
	)
	if err != nil {
		return nil, err
	}
	return &response.UpdateResultDTO{UpdatedCount: int64(len(ids))}, nil
}

// @MappedFrom deleteMessages(Set<Long> ids, @QueryParam(required = false) Boolean deleteLogically)
func (c *MessageController) DeleteMessages(
	ctx context.Context,
	ids []int64,
	deleteLogically *bool,
) (*response.DeleteResultDTO, error) {
	err := c.messageService.DeleteMessages(ctx, ids, deleteLogically)
	if err != nil {
		return nil, err
	}
	return &response.DeleteResultDTO{DeletedCount: int64(len(ids))}, nil
}

// PaginationResponse is a generic paginated response carrying a total count and records.
type PaginationResponse struct {
	Total   int64       `json:"total"`
	Records interface{} `json:"records"`
}

// Ensure dto import is used
var _ dto.StatisticsRecordDTO
