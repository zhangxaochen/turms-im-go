package dto

import "time"

// @MappedFrom CreateMessageDTO.java
type CreateMessageDTO struct {
	ID               *int64   `json:"id,omitempty"`
	IsGroupMessage   *bool    `json:"isGroupMessage,omitempty"`
	IsSystemMessage  *bool    `json:"isSystemMessage,omitempty"`
	Text             *string  `json:"text,omitempty"`
	Records          [][]byte `json:"records,omitempty"`
	SenderID         *int64   `json:"senderId,omitempty"`
	SenderIP         *string  `json:"senderIp,omitempty"`
	SenderDeviceType *int     `json:"senderDeviceType,omitempty"`
	TargetID         *int64   `json:"targetId,omitempty"`
	BurnAfter        *int     `json:"burnAfter,omitempty"`
	ReferenceID      *int64   `json:"referenceId,omitempty"`
	PreMessageID     *int64   `json:"preMessageId,omitempty"`
}

// @MappedFrom UpdateMessageDTO.java
type UpdateMessageDTO struct {
	SenderID         *int64     `json:"senderId,omitempty"`
	SenderIP         *string    `json:"senderIp,omitempty"`
	SenderDeviceType *int       `json:"senderDeviceType,omitempty"`
	IsSystemMessage  *bool      `json:"isSystemMessage,omitempty"`
	Text             *string    `json:"text,omitempty"`
	Records          [][]byte   `json:"records,omitempty"`
	BurnAfter        *int       `json:"burnAfter,omitempty"`
	RecallDate       *time.Time `json:"recallDate,omitempty"`
}

// @MappedFrom MessageStatisticsDTO.java
type MessageStatisticsDTO struct {
	SentMessagesOnAverage               *int64        `json:"sentMessagesOnAverage,omitempty"`
	AcknowledgedMessages                *int64        `json:"acknowledgedMessages,omitempty"`
	AcknowledgedMessagesOnAverage       *int64        `json:"acknowledgedMessagesOnAverage,omitempty"`
	SentMessages                        *int64        `json:"sentMessages,omitempty"`
	SentMessagesOnAverageRecords        []interface{} `json:"sentMessagesOnAverageRecords,omitempty"` // placeholder for StatisticsRecordDTO
	AcknowledgedMessagesRecords         []interface{} `json:"acknowledgedMessagesRecords,omitempty"`  // placeholder for StatisticsRecordDTO
	AcknowledgedMessagesOnAverageRecord []interface{} `json:"acknowledgedMessagesOnAverageRecords,omitempty"`
	SentMessagesRecords                 []interface{} `json:"sentMessagesRecords,omitempty"`
}
