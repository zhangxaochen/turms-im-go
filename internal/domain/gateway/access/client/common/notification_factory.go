package common

import (
	"time"

	"google.golang.org/protobuf/proto"

	"im.turms/server/internal/domain/common/constant"
	"im.turms/server/internal/domain/gateway/config"
	sessionbo "im.turms/server/internal/domain/gateway/session/bo"
	"im.turms/server/internal/infra/exception"
	"im.turms/server/pkg/protocol"
)

type NotificationFactory struct {
	returnReasonForServerError bool
}

// NewNotificationFactory enforces configuration dependency injection.
// @MappedFrom init(TurmsPropertiesManager propertiesManager)
func NewNotificationFactory(props *config.GatewayProperties) *NotificationFactory {
	if props == nil {
		props = config.NewGatewayProperties()
	}
	f := &NotificationFactory{}
	f.UpdateGlobalProperties(props)
	return f
}

// UpdateGlobalProperties dynamically updates properties from the configuration.
func (f *NotificationFactory) UpdateGlobalProperties(props *config.GatewayProperties) {
	if props != nil && props.ClientAPI != nil {
		f.returnReasonForServerError = props.ClientAPI.ReturnReasonForServerError
	}
}

// Create generates a generic Notification payload.
// @MappedFrom create(ResponseStatusCode code, long requestId)
func (f *NotificationFactory) Create(requestID *int64, code constant.ResponseStatusCode) *protocol.TurmsNotification {
	return f.CreateWithReason(requestID, code, nil)
}

// CreateWithReason generates a payload allowing reason texts depending on config.
// @MappedFrom create(ResponseStatusCode code, @Nullable String reason, long requestId)
func (f *NotificationFactory) CreateWithReason(requestID *int64, code constant.ResponseStatusCode, reason *string) *protocol.TurmsNotification {
	notification := &protocol.TurmsNotification{
		Timestamp: time.Now().UnixMilli(),
		RequestId: requestID,
		Code:      proto.Int32(int32(code)),
	}

	var actualReason *string
	if reason == nil {
		r := code.Reason()
		actualReason = &r
	} else {
		actualReason = reason
	}

	f.trySetReason(notification, code, actualReason)
	return notification
}

// CreateFromError parses a typed TurmsError securely.
// @MappedFrom create(ThrowableInfo info, long requestId)
func (f *NotificationFactory) CreateFromError(err error, requestID *int64) *protocol.TurmsNotification {
	code := constant.ResponseStatusCode_SERVER_INTERNAL_ERROR
	var reason *string

	if te, ok := err.(*exception.TurmsError); ok {
		code = constant.ResponseStatusCode(te.Code)
		if te.Message != "" {
			reason = &te.Message
		} else {
			r := code.Reason()
			reason = &r
		}
	} else if err != nil {
		// Just map to internal server error if not TurmsError.
		code = constant.ResponseStatusCode_SERVER_INTERNAL_ERROR
		errStr := err.Error()
		reason = &errStr
	} else {
		r := code.Reason()
		reason = &r
	}

	notification := &protocol.TurmsNotification{
		Timestamp: time.Now().UnixMilli(),
		RequestId: requestID,
		Code:      proto.Int32(int32(code)),
	}

	f.trySetReason(notification, code, reason)
	return notification
}

// CreateBuffer generates the serialized protobuf bytes directly.
// @MappedFrom createBuffer(CloseReason closeReason)
func (f *NotificationFactory) CreateCloseReasonBuffer(reason sessionbo.CloseReason) ([]byte, error) {
	code := reason.BusinessStatusCode

	notification := &protocol.TurmsNotification{
		Timestamp:   time.Now().UnixMilli(),
		CloseStatus: proto.Int32(int32(reason.Status)),
	}

	if code != 0 {
		notification.Code = proto.Int32(int32(code))
		var r *string
		if reason.Reason != "" {
			r = &reason.Reason
		}
		f.trySetReason(notification, code, r)
	} else {
		if reason.Reason != "" {
			notification.Reason = &reason.Reason
		}
	}
	return proto.Marshal(notification)
}

func (f *NotificationFactory) CreateBuffer(requestID *int64, code constant.ResponseStatusCode, reason string) ([]byte, error) {
	notification := f.CreateWithReason(requestID, code, &reason)
	return proto.Marshal(notification)
}

// SessionClosed generates a specialized notification when the server forcefully kicks the client.
// @MappedFrom sessionClosed(long requestId)
func (f *NotificationFactory) SessionClosed(requestID *int64) *protocol.TurmsNotification {
	return &protocol.TurmsNotification{
		Timestamp: time.Now().UnixMilli(),
		RequestId: requestID,
		Code:      proto.Int32(int32(constant.ResponseStatusCode_SERVER_INTERNAL_ERROR)),
	}
}

func (f *NotificationFactory) trySetReason(notification *protocol.TurmsNotification, code constant.ResponseStatusCode, reason *string) {
	if reason == nil || *reason == "" {
		return
	}
	if constant.IsServerError(int32(code)) && !f.returnReasonForServerError {
		return
	}
	notification.Reason = reason
}
