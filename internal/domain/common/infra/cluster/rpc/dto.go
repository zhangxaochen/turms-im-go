package rpc

import (
	"im.turms/server/internal/domain/common/infra/cluster/discovery"
)

// NodeTypeToHandleRpc represents the type of node that should handle the RPC request.
type NodeTypeToHandleRpc string

const (
	NodeTypeToHandleRpcBoth    NodeTypeToHandleRpc = "BOTH"
	NodeTypeToHandleRpcGateway NodeTypeToHandleRpc = "GATEWAY"
	NodeTypeToHandleRpcService NodeTypeToHandleRpc = "SERVICE"
)

// RpcRequest defines the interface for an RPC request.
type RpcRequest interface {
	Name() string
	CodecID() uint16
	NodeTypeToRequest() NodeTypeToHandleRpc
	NodeTypeToRespond() NodeTypeToHandleRpc
	IsAsync() bool
	Payload() ([]byte, error)
}

// RpcResponse defines a standard RPC response containing a payload or an error.
type RpcResponse struct {
	Payload []byte
	Err     error
}

// MapToDiscoveryNodeType maps the RPC node type to the discovery node type.
func (n NodeTypeToHandleRpc) MapToDiscoveryNodeType() discovery.NodeType {
	switch n {
	case NodeTypeToHandleRpcGateway:
		return discovery.NodeTypeGateway
	case NodeTypeToHandleRpcService:
		return discovery.NodeTypeService
	default:
		return ""
	}
}
