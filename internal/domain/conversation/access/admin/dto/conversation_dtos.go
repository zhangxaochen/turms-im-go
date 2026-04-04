package dto

import (
	"time"

	"im.turms/server/internal/domain/conversation/po"
)

// @MappedFrom UpdateConversationDTO(Date readDate)
type UpdateConversationDTO struct {
	ReadDate time.Time `json:"readDate"`
}

// @MappedFrom ConversationsDTO(List<PrivateConversation> privateConversations, List<GroupConversation> groupConversations)
type ConversationsDTO struct {
	PrivateConversations []*po.PrivateConversation `json:"privateConversations"`
	GroupConversations   []*po.GroupConversation   `json:"groupConversations"`
}
