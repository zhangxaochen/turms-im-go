package controller

// MessageController maps to MessageController.java
// @MappedFrom MessageController
type MessageController struct {
}

// @MappedFrom createMessages(@QueryParam(defaultValue = "true")
func (c *MessageController) CreateMessages() {
	// TODO: implement
}

// @MappedFrom queryMessages(@QueryParam(required = false)
func (c *MessageController) QueryMessages() {
	// TODO: implement
}

// @MappedFrom countMessages(@QueryParam(required = false)
func (c *MessageController) CountMessages() {
	// TODO: implement
}

// @MappedFrom updateMessages(Set<Long> ids, @RequestBody UpdateMessageDTO updateMessageDTO)
func (c *MessageController) UpdateMessages() {
	// TODO: implement
}

// @MappedFrom deleteMessages(Set<Long> ids, @QueryParam(required = false)
func (c *MessageController) DeleteMessages() {
	// TODO: implement
}
