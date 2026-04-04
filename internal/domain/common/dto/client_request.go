package dto

import (
	"fmt"
	"hash/fnv"

	"im.turms/server/pkg/protocol"
)

// ClientRequest maps to ClientRequest in Java.
// @MappedFrom ClientRequest
type ClientRequest struct {
	turmsRequest *protocol.TurmsRequest
	userId       *int64
	deviceType   *protocol.DeviceType
	clientIp     []byte // raw bytes like Java's byte[]
	requestId    *int64
}

// NewClientRequest constructs a ClientRequest.
func NewClientRequest(userId *int64, deviceType *protocol.DeviceType, clientIp []byte, requestId *int64, turmsRequest *protocol.TurmsRequest) *ClientRequest {
	return &ClientRequest{
		userId:       userId,
		deviceType:   deviceType,
		clientIp:     clientIp,
		requestId:    requestId,
		turmsRequest: turmsRequest,
	}
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
// Returns the raw IP bytes, matching Java's byte[] clientIp().
func (c *ClientRequest) ClientIp() []byte {
	return c.clientIp
}

// @MappedFrom requestId()
func (c *ClientRequest) RequestId() *int64 {
	return c.requestId
}

// String implements fmt.Stringer, matching Java's toString().
// Java: "ClientRequest[userId=..., deviceType=..., clientIp=..., requestId=..., turmsRequest=...]"
// @MappedFrom toString()
func (c *ClientRequest) String() string {
	var userIdStr, deviceTypeStr, ipStr, requestIdStr, requestStr string
	if c.userId != nil {
		userIdStr = fmt.Sprintf("%d", *c.userId)
	} else {
		userIdStr = "null"
	}
	if c.deviceType != nil {
		deviceTypeStr = c.deviceType.String()
	} else {
		deviceTypeStr = "null"
	}
	ipStr = fmt.Sprintf("%v", c.clientIp)
	if c.requestId != nil {
		requestIdStr = fmt.Sprintf("%d", *c.requestId)
	} else {
		requestIdStr = "null"
	}
	if c.turmsRequest != nil {
		requestStr = c.turmsRequest.String()
	} else {
		requestStr = "null"
	}
	return fmt.Sprintf("ClientRequest[userId=%s, deviceType=%s, clientIp=%s, requestId=%s, turmsRequest=%s]",
		userIdStr, deviceTypeStr, ipStr, requestIdStr, requestStr)
}

// Equals compares two ClientRequest instances, mirroring Java's equals().
// @MappedFrom equals(Object obj)
func (c *ClientRequest) Equals(other *ClientRequest) bool {
	if other == nil {
		return false
	}
	if c == other {
		return true
	}
	if !int64PtrEq(c.userId, other.userId) {
		return false
	}
	if !deviceTypeEq(c.deviceType, other.deviceType) {
		return false
	}
	if !bytesEq(c.clientIp, other.clientIp) {
		return false
	}
	if !int64PtrEq(c.requestId, other.requestId) {
		return false
	}
	// turmsRequest: proto messages don't implement ==; compare by pointer identity as Java does
	return c.turmsRequest == other.turmsRequest
}

// HashCode mirrors Java's Objects.hash + Arrays.hashCode pattern.
// @MappedFrom hashCode()
func (c *ClientRequest) HashCode() int {
	h := fnv.New32a()
	if c.userId != nil {
		fmt.Fprintf(h, "%d", *c.userId)
	}
	if c.deviceType != nil {
		fmt.Fprintf(h, "%d", int32(*c.deviceType))
	}
	h.Write(c.clientIp)
	if c.requestId != nil {
		fmt.Fprintf(h, "%d", *c.requestId)
	}
	return int(h.Sum32())
}

// --- helpers ---

func int64PtrEq(a, b *int64) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

func deviceTypeEq(a, b *protocol.DeviceType) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

func bytesEq(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
