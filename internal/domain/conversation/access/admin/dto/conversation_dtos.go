package dto

import "time"

// @MappedFrom UpdateConversationDTO(Date readDate)
type UpdateConversationDTO struct {
	ReadDate time.Time `json:"readDate"`
}

// @MappedFrom ConversationsDTO(List<PrivateConversation> privateConversations, List<GroupConversation> groupConversations)
type ConversationsDTO struct {
	PrivateConversations []interface{} `json:"privateConversations"`
	GroupConversations   []interface{} `json:"groupConversations"`
}
