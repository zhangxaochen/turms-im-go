package rpc

import (
	"encoding/json"

	"im.turms/server/pkg/protocol"
)

type SetUserOfflineRequest struct {
	UserID             int64
	DeviceTypes        []protocol.DeviceType
	SessionCloseStatus int
}

func (req *SetUserOfflineRequest) Name() string {
	return "SetUserOfflineRequest"
}

func (req *SetUserOfflineRequest) CodecID() uint16 {
	return 1 // Just a stub ID
}

func (req *SetUserOfflineRequest) NodeTypeToRequest() NodeTypeToHandleRpc {
	return NodeTypeToHandleRpcGateway
}

func (req *SetUserOfflineRequest) NodeTypeToRespond() NodeTypeToHandleRpc {
	return NodeTypeToHandleRpcGateway
}

func (req *SetUserOfflineRequest) IsAsync() bool {
	return false
}

func (req *SetUserOfflineRequest) Payload() ([]byte, error) {
	return json.Marshal(req)
}

// CountOnlineUsersRequest is sent to all other cluster members to aggregate
// the total online user count across the cluster.
// Java equivalent: CountOnlineUsersRequest
type CountOnlineUsersRequest struct{}

func (req *CountOnlineUsersRequest) Name() string {
	return "CountOnlineUsersRequest"
}

func (req *CountOnlineUsersRequest) CodecID() uint16 {
	return 2
}

func (req *CountOnlineUsersRequest) NodeTypeToRequest() NodeTypeToHandleRpc {
	return NodeTypeToHandleRpcGateway
}

func (req *CountOnlineUsersRequest) NodeTypeToRespond() NodeTypeToHandleRpc {
	return NodeTypeToHandleRpcGateway
}

func (req *CountOnlineUsersRequest) IsAsync() bool {
	return false
}

func (req *CountOnlineUsersRequest) Payload() ([]byte, error) {
	return json.Marshal(req)
}
