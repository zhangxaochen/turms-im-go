package common

import (
	"im.turms/server/internal/domain/common/constant"
	"im.turms/server/pkg/protocol"
)

// RequestHandlerResult carries the response code and optional reason of a handler execution.
type RequestHandlerResult struct {
	Code   constant.ResponseStatusCode
	Reason string
	Data   *protocol.TurmsNotification_Data
}

// NewRequestHandlerResult creates a new RequestHandlerResult.
// @MappedFrom RequestHandlerResult(ResponseStatusCode code, String reason, TurmsNotification.Data data)
func NewRequestHandlerResultWithData(code constant.ResponseStatusCode, reason string, data *protocol.TurmsNotification_Data) *RequestHandlerResult {
	return &RequestHandlerResult{
		Code:   code,
		Reason: reason,
		Data:   data,
	}
}

// @MappedFrom RequestHandlerResult(ResponseStatusCode code, String reason)
func NewRequestHandlerResult(code constant.ResponseStatusCode, reason string) *RequestHandlerResult {
	return &RequestHandlerResult{
		Code:   code,
		Reason: reason,
	}
}
