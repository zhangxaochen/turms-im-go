package logging

// Stub implementation for ApiLoggingContext
type ApiLoggingContext struct {
}

func NewApiLoggingContext() *ApiLoggingContext {
	return &ApiLoggingContext{}
}

// @MappedFrom shouldLogRequest(TurmsRequest.KindCase requestType)
func (c *ApiLoggingContext) ShouldLogRequest(requestType int) bool {
	// Stub implementation
	return true
}

// @MappedFrom shouldLogNotification(TurmsRequest.KindCase requestType)
func (c *ApiLoggingContext) ShouldLogNotification(requestType int) bool {
	// Stub implementation
	return true
}
