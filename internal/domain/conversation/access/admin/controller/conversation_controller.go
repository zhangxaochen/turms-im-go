package controller

import (
	"im.turms/server/internal/domain/common/access/admin/controller"
	"im.turms/server/internal/domain/conversation/service"
	"im.turms/server/internal/infra/property"
)

// ConversationController maps to ConversationController.java
// @MappedFrom ConversationController
type ConversationController struct {
	*controller.BaseController
	conversationService *service.ConversationService
}

func NewConversationController(
	propertiesManager *property.TurmsPropertiesManager,
	conversationService *service.ConversationService,
) *ConversationController {
	return &ConversationController{
		BaseController:      controller.NewBaseController(propertiesManager),
		conversationService: conversationService,
	}
}

// @MappedFrom queryConversations(@QueryParam(required = false)
func (c *ConversationController) QueryConversations() {
	// TODO: Implement with proper admin request handling pattern
}

// @MappedFrom deleteConversations(@QueryParam(required = false)
func (c *ConversationController) DeleteConversations() {
	// TODO: Implement with proper admin request handling pattern
}

// @MappedFrom updateConversations(@QueryParam(required = false)
func (c *ConversationController) UpdateConversations() {
	// TODO: Implement with proper admin request handling pattern
}
