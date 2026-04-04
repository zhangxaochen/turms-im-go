package rpc

import (
	"encoding/json"

	"im.turms/server/internal/domain/common/access/servicerequest/dto"
	"im.turms/server/internal/domain/common/infra/cluster/rpc"
)

// HandleServiceRequest RPC request for forwarding gateway requests to backend services.
// @MappedFrom im.turms.server.common.access.servicerequest.rpc.HandleServiceRequest
type HandleServiceRequest struct {
	ServiceRequest *dto.ServiceRequest
}

func NewHandleServiceRequest(req *dto.ServiceRequest) *HandleServiceRequest {
	return &HandleServiceRequest{ServiceRequest: req}
}

func (r *HandleServiceRequest) Name() string { return "HandleServiceRequest" }

// CodecID is hardcoded to 100 for now to align with internal registry.
func (r *HandleServiceRequest) CodecID() uint16 { return 100 }

func (r *HandleServiceRequest) NodeTypeToRequest() rpc.NodeTypeToHandleRpc {
	return rpc.NodeTypeToHandleRpcService
}

func (r *HandleServiceRequest) NodeTypeToRespond() rpc.NodeTypeToHandleRpc {
	return rpc.NodeTypeToHandleRpcGateway
}

func (r *HandleServiceRequest) IsAsync() bool { return false }

// Payload uses JSON for simple serialization until binary codec is fully ready.
func (r *HandleServiceRequest) Payload() ([]byte, error) {
	return json.Marshal(r.ServiceRequest)
}
