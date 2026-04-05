package common

import (
	"fmt"

	"im.turms/server/internal/domain/common/constant"
	"im.turms/server/pkg/protocol"
)

// RequestHandlerResult carries the response code and optional reason of a handler execution.
type RequestHandlerResult struct {
	Code   constant.ResponseStatusCode
	Reason *string
	Data   *protocol.TurmsNotification_Data
}

// NewRequestHandlerResult creates a new RequestHandlerResult.
// @MappedFrom RequestHandlerResult(ResponseStatusCode code, String reason, TurmsNotification.Data data)
func NewRequestHandlerResultWithData(code constant.ResponseStatusCode, reason *string, data *protocol.TurmsNotification_Data) *RequestHandlerResult {
	return &RequestHandlerResult{
		Code:   code,
		Reason: reason,
		Data:   data,
	}
}

func NewRequestHandlerResult(code constant.ResponseStatusCode, reason *string) *RequestHandlerResult {
	return &RequestHandlerResult{
		Code:   code,
		Reason: reason,
	}
}

// RequestHandlerResultOfCode creates a RequestHandlerResult with only a code
// @MappedFrom RequestHandlerResult(ResponseStatusCode code)
func RequestHandlerResultOfCode(code constant.ResponseStatusCode) *RequestHandlerResult {
	return &RequestHandlerResult{
		Code: code,
	}
}

// @MappedFrom toString()
func (r *RequestHandlerResult) String() string {
	reasonStr := "null"
	if r.Reason != nil {
		reasonStr = *r.Reason
	}
	return fmt.Sprintf("RequestHandlerResult[code=%v, reason='%s']", r.Code, reasonStr)
}
