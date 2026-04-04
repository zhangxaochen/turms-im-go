package logging

import (
	"fmt"
	"log/slog"
	"strings"
	"time"
)

// logFieldDelimiter maps to CommonLogger.LOG_FIELD_DELIMITER in Java.
const logFieldDelimiter = "|"

// clientAPILogger is the structured logger for client API requests.
// In Java this is CommonLogger.CLIENT_API_LOGGER.
var clientAPILogger = slog.Default()

// notificationLogger is the logger for notifications.
var notificationLogger = slog.Default()

// joinFields joins nullable fields with the log delimiter, rendering nil as "null".
func joinFields(fields ...interface{}) string {
	parts := make([]string, 0, len(fields))
	for _, f := range fields {
		if f == nil {
			parts = append(parts, "null")
		} else {
			parts = append(parts, fmt.Sprintf("%v", f))
		}
	}
	return strings.Join(parts, logFieldDelimiter)
}

// @MappedFrom ClientApiLogging
// LogRequest logs a client request with TurmsNotification response.
// Maps to:
//
//	log(@Nullable Integer sessionId, @Nullable Long userId, @Nullable DeviceType deviceType,
//	    @Nullable Integer version, String ip, long requestId, TurmsRequest.KindCase requestType,
//	    int requestSize, long requestTime, TurmsNotification response, long processingTime)
func LogRequest(
	sessionID *int32,
	userID *int64,
	deviceType *int32, // DeviceType enum value; nil means absent
	version *int32,
	ip string,
	requestID int64,
	requestType string, // TurmsRequest.KindCase.name()
	requestSize int,
	requestTime time.Time,
	responseCode int32,
	responseDataType *string, // nil if no data
	responseSize int,
	processingTime time.Duration,
) {
	msg := joinFields(
		sessionID,
		userID,
		deviceType,
		version,
		ip,
		requestID,
		requestType,
		requestSize,
		requestTime.UnixMilli(),
		responseCode,
		responseDataType,
		responseSize,
		processingTime.Milliseconds(),
	)
	clientAPILogger.Info(msg)
}

// LogRequestWithCode logs a client request with a numeric response code (no notification data).
// Maps to Java:
//
//	log(..., int responseCode, long processingTime)
func LogRequestWithCode(
	sessionID *int32,
	userID *int64,
	deviceType *int32,
	version *int32,
	ip string,
	requestID int64,
	requestType string,
	requestSize int,
	requestTime time.Time,
	responseCode int32,
	processingTime time.Duration,
) {
	msg := joinFields(
		sessionID,
		userID,
		deviceType,
		version,
		ip,
		requestID,
		requestType,
		requestSize,
		requestTime.UnixMilli(),
		responseCode,
		nil,  // responseDataType
		"0",  // responseSerializedSize
		processingTime.Milliseconds(),
	)
	clientAPILogger.Info(msg)
}
