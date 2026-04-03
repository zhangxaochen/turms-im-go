package common

import "im.turms/server/internal/domain/common/constant"

// RequestHandlerResult carries the response code and optional reason of a handler execution.
type RequestHandlerResult struct {
	Code   constant.ResponseStatusCode
	Reason string
}

// NewRequestHandlerResult creates a new RequestHandlerResult.
// @MappedFrom RequestHandlerResult(ResponseStatusCode code, String reason)
func NewRequestHandlerResult(code constant.ResponseStatusCode, reason string) *RequestHandlerResult {
	return &RequestHandlerResult{
		Code:   code,
		Reason: reason,
	}
}
