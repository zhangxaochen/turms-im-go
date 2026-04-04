package common_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"

	"im.turms/server/internal/domain/common/constant"
	"im.turms/server/internal/domain/gateway/access/client/common"
	"im.turms/server/internal/domain/gateway/config"
	"im.turms/server/internal/infra/exception"
)

func TestNotificationFactory_CreateWithReason(t *testing.T) {
	requestID := int64(12345)

	tests := []struct {
		name         string
		props        *config.GatewayProperties
		requestID    *int64
		code         constant.ResponseStatusCode
		reason       string
		expectReason bool
	}{
		{
			name:         "Client Error - Reason Included",
			props:        &config.GatewayProperties{ClientAPI: &config.ClientAPIProperties{ReturnReasonForServerError: false}},
			requestID:    &requestID,
			code:         constant.ResponseStatusCode_INVALID_REQUEST,
			reason:       "Bad parameter",
			expectReason: true,
		},
		{
			name:         "Server Error - Reason Masked",
			props:        &config.GatewayProperties{ClientAPI: &config.ClientAPIProperties{ReturnReasonForServerError: false}},
			requestID:    &requestID,
			code:         constant.ResponseStatusCode_SERVER_INTERNAL_ERROR,
			reason:       "DB connection failed",
			expectReason: false,
		},
		{
			name:         "Server Error - Reason Included When Configured",
			props:        &config.GatewayProperties{ClientAPI: &config.ClientAPIProperties{ReturnReasonForServerError: true}},
			requestID:    &requestID,
			code:         constant.ResponseStatusCode_SERVER_INTERNAL_ERROR,
			reason:       "DB connection failed",
			expectReason: true,
		},
		{
			name:         "No Reason Provided",
			props:        &config.GatewayProperties{ClientAPI: &config.ClientAPIProperties{ReturnReasonForServerError: true}},
			requestID:    &requestID,
			code:         constant.ResponseStatusCode_OK,
			reason:       "",
			expectReason: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory := common.NewNotificationFactory(tt.props)
			start := time.Now().UnixMilli()
			notification := factory.CreateWithReason(tt.requestID, tt.code, &tt.reason)

			assert.NotNil(t, notification)
			assert.Equal(t, tt.requestID, notification.RequestId)
			assert.Equal(t, int32(tt.code), *notification.Code)
			assert.GreaterOrEqual(t, notification.Timestamp, start)

			if tt.expectReason {
				assert.NotNil(t, notification.Reason)
				assert.Equal(t, tt.reason, *notification.Reason)
			} else {
				assert.Nil(t, notification.Reason)
			}
		})
	}
}

func TestNotificationFactory_CreateFromError(t *testing.T) {
	requestID := int64(67890)

	tests := []struct {
		name         string
		props        *config.GatewayProperties
		requestID    *int64
		err          error
		expectCode   constant.ResponseStatusCode
		expectReason string
		isServerErr  bool
	}{
		{
			name:         "TurmsError - Client Error",
			props:        &config.GatewayProperties{ClientAPI: &config.ClientAPIProperties{ReturnReasonForServerError: false}},
			requestID:    &requestID,
			err:          exception.NewTurmsError(int32(constant.ResponseStatusCode_INVALID_REQUEST), "Bad input"),
			expectCode:   constant.ResponseStatusCode_INVALID_REQUEST,
			expectReason: "Bad input",
			isServerErr:  false,
		},
		{
			name:         "TurmsError - Server Error - Masked",
			props:        &config.GatewayProperties{ClientAPI: &config.ClientAPIProperties{ReturnReasonForServerError: false}},
			requestID:    &requestID,
			err:          exception.NewTurmsError(int32(constant.ResponseStatusCode_SERVER_INTERNAL_ERROR), "Internal crash"),
			expectCode:   constant.ResponseStatusCode_SERVER_INTERNAL_ERROR,
			expectReason: "",
			isServerErr:  true,
		},
		{
			name:         "Native Error - Acts As Server Error - Masked",
			props:        &config.GatewayProperties{ClientAPI: &config.ClientAPIProperties{ReturnReasonForServerError: false}},
			requestID:    &requestID,
			err:          errors.New("Unknown panic"),
			expectCode:   constant.ResponseStatusCode_SERVER_INTERNAL_ERROR,
			expectReason: "",
			isServerErr:  true,
		},
		{
			name:         "Native Error - Acts As Server Error - Not Masked",
			props:        &config.GatewayProperties{ClientAPI: &config.ClientAPIProperties{ReturnReasonForServerError: true}},
			requestID:    &requestID,
			err:          errors.New("Unknown panic"),
			expectCode:   constant.ResponseStatusCode_SERVER_INTERNAL_ERROR,
			expectReason: "Unknown panic",
			isServerErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory := common.NewNotificationFactory(tt.props)
			notification := factory.CreateFromError(tt.err, tt.requestID)

			assert.NotNil(t, notification)
			assert.Equal(t, tt.requestID, notification.RequestId)
			assert.Equal(t, int32(tt.expectCode), *notification.Code)

			if tt.expectReason != "" {
				assert.NotNil(t, notification.Reason)
				assert.Equal(t, tt.expectReason, *notification.Reason)
			} else {
				assert.Nil(t, notification.Reason)
			}
		})
	}
}

func TestNotificationFactory_SessionClosed(t *testing.T) {
	requestID := int64(999)
	factory := common.NewNotificationFactory(nil)
	notification := factory.SessionClosed(&requestID)

	assert.NotNil(t, notification)
	assert.Equal(t, &requestID, notification.RequestId)
	assert.Equal(t, proto.Int32(int32(constant.ResponseStatusCode_SERVER_INTERNAL_ERROR)), notification.Code)
	assert.NotZero(t, notification.Timestamp)
	assert.Nil(t, notification.Reason)
}
