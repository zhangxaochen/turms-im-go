package common

import (
	"time"

	"google.golang.org/protobuf/proto"

	"im.turms/server/internal/domain/common/constant"
	"im.turms/server/internal/domain/gateway/config"
	"im.turms/server/internal/infra/exception"
	"im.turms/server/pkg/protocol"
)

// NotificationFactory standardizes the creation of TurmsNotification objects
// sent down to clients containing success/failure outcomes.
type NotificationFactory struct {
	propsManager *config.GatewayProperties
}

// NewNotificationFactory enforces configuration dependency injection.
// @MappedFrom init(TurmsPropertiesManager propertiesManager)
func NewNotificationFactory(props *config.GatewayProperties) *NotificationFactory {
	if props == nil {
		props = config.NewGatewayProperties()
	}
	return &NotificationFactory{
		propsManager: props,
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

	f.trySetReason(notification, code, reason)
	return notification
}

// CreateFromError parses a typed TurmsError securely.
// @MappedFrom create(ThrowableInfo info, long requestId)
func (f *NotificationFactory) CreateFromError(err error, requestID *int64) *protocol.TurmsNotification {
	code := constant.ResponseStatusCode_SERVER_INTERNAL_ERROR
	var reason *string

	if te, ok := err.(*exception.TurmsError); ok {
		code = constant.ResponseStatusCode(te.Code)
		reason = &te.Message
	} else if err != nil {
		errStr := err.Error()
		reason = &errStr
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
	if reason != nil {
		if *reason == "" {
			return
		}
		if constant.IsServerError(int32(code)) {
			if f.propsManager.ClientAPI.ReturnReasonForServerError {
				notification.Reason = reason
			}
		} else {
			notification.Reason = reason
		}
	} else {
		// Fallback to default reason based on standard code values if missing
		// Since we don't have GetReason() on ResponseStatusCode mapped yet in Go,
		// we skip setting a reason if it's nil.
	}
}
