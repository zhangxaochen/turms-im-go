package dto

import "im.turms/server/pkg/protocol"

// ClientRequest maps to ClientRequest in Java.
// @MappedFrom ClientRequest
type ClientRequest struct {
	turmsRequest *protocol.TurmsRequest
	userId       *int64
	deviceType   *protocol.DeviceType
	clientIp     *string
	requestId    *int64
}

// @MappedFrom turmsRequest()
func (c *ClientRequest) TurmsRequest() *protocol.TurmsRequest {
	return c.turmsRequest
}

// @MappedFrom userId()
func (c *ClientRequest) UserId() *int64 {
	return c.userId
}

// @MappedFrom deviceType()
func (c *ClientRequest) DeviceType() *protocol.DeviceType {
	return c.deviceType
}

// @MappedFrom clientIp()
func (c *ClientRequest) ClientIp() *string {
	return c.clientIp
}

// @MappedFrom requestId()
func (c *ClientRequest) RequestId() *int64 {
	return c.requestId
}

// @MappedFrom equals(Object obj)
func (c *ClientRequest) Equals(obj interface{}) bool {
	return false
}

// @MappedFrom hashCode()
func (c *ClientRequest) HashCode() int {
	return 0
}
