package bo

import (
	"im.turms/server/internal/domain/common/constant"
)

// CloseReason represents the reason for closing a connection.
// @MappedFrom CloseReason.java
type CloseReason struct {
	Status             constant.SessionCloseStatus
	BusinessStatusCode constant.ResponseStatusCode
	Reason             string
	IsNotifyClient      bool
}

func NewCloseReason(status constant.SessionCloseStatus) CloseReason {
	// SessionCloseStatus.isNotifyClient() logic from Java
	isNotify := false
	if status == constant.SessionCloseStatus_DISCONNECTED_BY_OTHER_DEVICE ||
		status == constant.SessionCloseStatus_DISCONNECTED_BY_ADMIN ||
		status == constant.SessionCloseStatus_DISCONNECTED_BY_SERVER ||
		status == constant.SessionCloseStatus_SWITCH ||
		status == constant.SessionCloseStatus_HEARTBEAT_TIMEOUT {
		isNotify = true
	}
	return CloseReason{
		Status:         status,
		IsNotifyClient: isNotify,
	}
}

func CloseReasonFromError(err error) CloseReason {
	if err == nil {
		return NewCloseReason(constant.SessionCloseStatus_UNKNOWN_ERROR)
	}

	fromErr, ok := err.(interface {
		Code() constant.ResponseStatusCode
		Reason() string
	})

	if ok {
		code := fromErr.Code()
		status := constant.SessionCloseStatus_UNKNOWN_ERROR

		// Map some status codes to close status, simple for now
		if code >= constant.ResponseStatusCode_SERVER_INTERNAL_ERROR && code < 1300 {
			if code == constant.ResponseStatusCode_SERVER_UNAVAILABLE {
				status = constant.SessionCloseStatus_SERVER_UNAVAILABLE
			} else {
				status = constant.SessionCloseStatus_SERVER_ERROR
			}
		} else if code == constant.ResponseStatusCode_ILLEGAL_ARGUMENT || code == constant.ResponseStatusCode_INVALID_REQUEST {
			status = constant.SessionCloseStatus_ILLEGAL_REQUEST
		}

		return CloseReason{
			Status:             status,
			BusinessStatusCode: code,
			Reason:             fromErr.Reason(),
		}
	}

	return CloseReason{
		Status: constant.SessionCloseStatus_UNKNOWN_ERROR,
		Reason: err.Error(),
	}
}
