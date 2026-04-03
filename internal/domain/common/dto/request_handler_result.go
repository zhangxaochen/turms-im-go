package dto

import (
	"im.turms/server/internal/domain/common/constant"
	"im.turms/server/pkg/protocol"
)

// RequestHandlerResult maps to RequestHandlerResult in Java.
// @MappedFrom RequestHandlerResult
type RequestHandlerResult struct {
	Code          constant.ResponseStatusCode
	Reason        *string
	Response      *protocol.TurmsNotification_Data
	Notifications []any // Placeholder for List<Notification>
}

// @MappedFrom RequestHandlerResult(ResponseStatusCode code, @Nullable String reason, @Nullable TurmsNotification.Data response, List<Notification> notifications)
func NewRequestHandlerResult(code constant.ResponseStatusCode, reason *string, response *protocol.TurmsNotification_Data, notifications []any) *RequestHandlerResult {
	return &RequestHandlerResult{
		Code:          code,
		Reason:        reason,
		Response:      response,
		Notifications: notifications,
	}
}
