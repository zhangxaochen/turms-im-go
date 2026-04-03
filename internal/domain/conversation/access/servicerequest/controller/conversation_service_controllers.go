package controller

// ConversationServiceController maps to ConversationServiceController.java
// @MappedFrom ConversationServiceController
type ConversationServiceController struct {
}

// @MappedFrom handleQueryConversationsRequest()
func (c *ConversationServiceController) HandleQueryConversationsRequest() {
}

// @MappedFrom handleUpdateTypingStatusRequest()
func (c *ConversationServiceController) HandleUpdateTypingStatusRequest() {
}

// @MappedFrom handleUpdateConversationRequest()
func (c *ConversationServiceController) HandleUpdateConversationRequest() {
}

// ConversationSettingsServiceController maps to ConversationSettingsServiceController.java
// @MappedFrom ConversationSettingsServiceController
type ConversationSettingsServiceController struct {
}

// @MappedFrom handleUpdateConversationSettingsRequest()
func (c *ConversationSettingsServiceController) HandleUpdateConversationSettingsRequest() {
}

// @MappedFrom handleDeleteConversationSettingsRequest()
func (c *ConversationSettingsServiceController) HandleDeleteConversationSettingsRequest() {
}

// @MappedFrom handleQueryConversationSettingsRequest()
func (c *ConversationSettingsServiceController) HandleQueryConversationSettingsRequest() {
}
