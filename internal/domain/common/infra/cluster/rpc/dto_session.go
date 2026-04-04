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
	// Simple stub json serialization for now
	return json.Marshal(req)
}
