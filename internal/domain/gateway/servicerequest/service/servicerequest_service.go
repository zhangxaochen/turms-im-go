package service

import (
	"context"
	"encoding/json"
	"time"

	"google.golang.org/protobuf/proto"
	"im.turms/server/internal/domain/common/access/servicerequest/dto"
	"im.turms/server/internal/domain/common/access/servicerequest/rpc"
	"im.turms/server/internal/domain/common/constant"

	cluster "im.turms/server/internal/domain/common/infra/cluster/rpc"
	"im.turms/server/internal/domain/gateway/session"
	"im.turms/server/pkg/protocol"
)

type ServiceRequestService struct {
	rpcService *cluster.RpcService
}

func NewServiceRequestService(rpcService *cluster.RpcService) *ServiceRequestService {
	return &ServiceRequestService{
		rpcService: rpcService,
	}
}

// HandleServiceRequest handles an incoming service request from a gateway context by forwarding it via RPC.
// @MappedFrom handleServiceRequest(UserSession session, ServiceRequest serviceRequest)
func (s *ServiceRequestService) HandleServiceRequest(ctx context.Context, defaultSession *session.UserSession, serviceRequest *dto.ServiceRequest) (*protocol.TurmsNotification, error) {
	// Update request timestamp for heartbeat maintenance
	defaultSession.SetLastRequestTimestampToNow()

	// Forward the request to a backend service node via RPC
	rpcReq := rpc.NewHandleServiceRequest(serviceRequest)
	rpcResp, err := s.rpcService.RequestResponse(ctx, "", rpcReq)
	if err != nil {
		return nil, err
	}

	// Java parity: defaultIfEmpty(REQUEST_RESPONSE_NO_CONTENT)
	// If RPC returns nil/empty payload, return NO_CONTENT notification
	if rpcResp == nil || len(rpcResp.Payload) == 0 {
		return &protocol.TurmsNotification{
			Timestamp: time.Now().UnixMilli(),
			RequestId: proto.Int64(serviceRequest.RequestId),
			Code:      proto.Int32(int32(constant.ResponseStatusCode_NO_CONTENT)),
		}, nil
	}

	// Unmarshal search response from JSON (until binary codec is implemented)
	var serviceResp dto.ServiceResponse
	if err := json.Unmarshal(rpcResp.Payload, &serviceResp); err != nil {
		return nil, err
	}

	return s.getNotificationFromResponse(&serviceResp, serviceRequest.RequestId), nil
}

// getNotificationFromResponse maps the backend ServiceResponse back to a TurmsNotification for the client.
// @MappedFrom getNotificationFromResponse(@NotNull ServiceResponse response, long requestId)
func (s *ServiceRequestService) getNotificationFromResponse(response *dto.ServiceResponse, requestId int64) *protocol.TurmsNotification {
	// Java parity: validate that response.Code is not null/zero.
	// Java throws IllegalArgumentException if code is null.
	// In Go, Code=0 would mean the backend didn't set a code, which is an error.
	if response.Code == 0 {
		// Return NO_CONTENT as fallback for empty response code
		return &protocol.TurmsNotification{
			Timestamp: time.Now().UnixMilli(),
			RequestId: proto.Int64(requestId),
			Code:      proto.Int32(int32(constant.ResponseStatusCode_NO_CONTENT)),
		}
	}

	notification := &protocol.TurmsNotification{
		Timestamp: time.Now().UnixMilli(),
		RequestId: proto.Int64(requestId),
		Code:      proto.Int32(response.Code),
	}
	if response.Reason != "" {
		notification.Reason = proto.String(response.Reason)
	}
	if response.Data != nil {
		notification.Data = response.Data
	}
	return notification
}

// GetNotificationFromResponse is a helper to create a notification for a given status code and request ID.
// Used for creating NO_CONTENT and similar notifications externally.
func GetNotificationFromResponse(code constant.ResponseStatusCode, requestId int64) *protocol.TurmsNotification {
	return &protocol.TurmsNotification{
		Timestamp: time.Now().UnixMilli(),
		RequestId: proto.Int64(requestId),
		Code:      proto.Int32(int32(code)),
	}
}
