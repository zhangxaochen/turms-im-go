package dto

import (
	"im.turms/server/pkg/protocol"
)

// ServiceRequest represents a client service request to be handled by turms-service.
// @MappedFrom im.turms.server.common.access.servicerequest.dto.ServiceRequest
type ServiceRequest struct {
	Ip           []byte
	UserId       int64
	DeviceType   protocol.DeviceType
	RequestId    int64
	TurmsRequest *protocol.TurmsRequest
	Type         any // In Java, this is an Object to support various request identifiers
	Buffer       []byte
}

// ServiceResponse represents a response from turms-service for a ServiceRequest.
// @MappedFrom im.turms.server.common.access.servicerequest.dto.ServiceResponse
type ServiceResponse struct {
	Code   int32
	Reason string
	Data   *protocol.TurmsNotification_Data
}
