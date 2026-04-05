package service

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"google.golang.org/protobuf/proto"
	"im.turms/server/internal/domain/common/access/servicerequest/dto"
	"im.turms/server/internal/domain/common/access/servicerequest/rpc"

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

	if rpcResp == nil || rpcResp.Payload == nil || len(rpcResp.Payload) == 0 {
		// Bug 8049: Missing defaultIfEmpty(REQUEST_RESPONSE_NO_CONTENT) equivalent
		return s.getNotificationFromResponse(&dto.ServiceResponse{
			Code: 1001, // 1001 corresponds to constant.ResponseStatusCode_NO_CONTENT business code
		}, serviceRequest.RequestId), nil
	}

	// Unmarshal search response from JSON (until binary codec is implemented)
	var serviceResp dto.ServiceResponse
	if err := json.Unmarshal(rpcResp.Payload, &serviceResp); err != nil {
		// Log error if backend returned something unparseable, but normally RPC layer handles this.
		return nil, err
	}

	return s.getNotificationFromResponse(&serviceResp, serviceRequest.RequestId), nil
}

// getNotificationFromResponse maps the backend ServiceResponse back to a TurmsNotification for the client.
// @MappedFrom getNotificationFromResponse(@NotNull ServiceResponse response, long requestId)
func (s *ServiceRequestService) getNotificationFromResponse(response *dto.ServiceResponse, requestId int64) *protocol.TurmsNotification {
	// Bug 8053: getNotificationFromResponse missing null code validation
	if response.Code == 0 {
		log.Printf("Received ServiceResponse with an unset/zero Code for request ID: %d", requestId)
		// Fallback to internal error if no valid code provided
		response.Code = 1002 // 1002 corresponds to constant.ResponseStatusCode_SERVER_INTERNAL_ERROR
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
