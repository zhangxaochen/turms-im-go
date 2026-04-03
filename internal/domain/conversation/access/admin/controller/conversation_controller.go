package controller

// ConversationController maps to ConversationController.java
// @MappedFrom ConversationController
type ConversationController struct {
}

// @MappedFrom queryConversations(@QueryParam(required = false)
func (c *ConversationController) QueryConversations() {
}

// @MappedFrom deleteConversations(@QueryParam(required = false)
func (c *ConversationController) DeleteConversations() {
}

// @MappedFrom updateConversations(@QueryParam(required = false)
func (c *ConversationController) UpdateConversations() {
}
