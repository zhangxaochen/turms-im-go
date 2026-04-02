package logging

import (
	"log"
)

// @MappedFrom log(@Nullable Integer sessionId, @Nullable Long userId, @Nullable DeviceType deviceType, @Nullable Integer version, String ip, long requestId, TurmsRequest.KindCase requestType, int requestSize, long requestTime, int responseCode, long processingTime)
// @MappedFrom log(@Nullable Integer sessionId, @Nullable Long userId, @Nullable DeviceType deviceType, @Nullable Integer version, String ip, long requestId, TurmsRequest.KindCase requestType, int requestSize, long requestTime, TurmsNotification response, long processingTime)
// @MappedFrom log(ClientRequest request, ServiceRequest serviceRequest, long requestSize, long requestTime, ServiceResponse response, long processingTime)
// @MappedFrom log(SimpleTurmsNotification notification, int notificationBytes, int recipientCount, int onlineRecipientCount)
// @MappedFrom log(@Nullable Integer sessionId, @Nullable Long userId, @Nullable DeviceType deviceType, @Nullable Integer version, String ip, long requestId, String requestType, int requestSize, long requestTime, int responseCode, @Nullable String responseDataType, int responseSize, long processingTime)
func Log(request interface{}, serviceRequest interface{}, requestSize int64, requestTime int64, response interface{}, processingTime int64) {
	// Stub implementation for ClientApiLogging.log
	log.Printf("mock client api log: %v", request)
}
