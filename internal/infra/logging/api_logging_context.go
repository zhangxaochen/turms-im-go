package logging

// Stub implementation for ApiLoggingContext
type ApiLoggingContext struct {
}

func NewApiLoggingContext() *ApiLoggingContext {
	return &ApiLoggingContext{}
}

func (c *ApiLoggingContext) ShouldLogRequest(requestType int) bool {
	// Stub implementation
	return true
}

func (c *ApiLoggingContext) ShouldLogNotification(requestType int) bool {
	// Stub implementation
	return true
}
