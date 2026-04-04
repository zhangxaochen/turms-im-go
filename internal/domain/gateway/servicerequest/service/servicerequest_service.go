package service

import (
	"context"
	"time"

	"google.golang.org/protobuf/proto"
	"im.turms/server/internal/domain/common/constant"
	"im.turms/server/internal/domain/gateway/session"
	"im.turms/server/pkg/protocol"
)

type ServiceRequestService struct {
	node any // placeholder for cluster.Node
}

func NewServiceRequestService(node any) *ServiceRequestService {
	return &ServiceRequestService{
		node: node,
	}
}

// HandleServiceRequest handles an incoming service request from a gateway context
// @MappedFrom handleServiceRequest(UserSession session, ServiceRequest serviceRequest)
func (s *ServiceRequestService) HandleServiceRequest(ctx context.Context, defaultSession *session.UserSession, serviceRequest any) (*protocol.TurmsNotification, error) {
	// Update request timestamp
	defaultSession.SetLastRequestTimestampToNow()

	// TODO: Obtain buffer from serviceRequest and retain it (like serviceRequest.getTurmsRequestBuffer().retain())

	// For now, return a basic notification with timestamp to satisfy minimal client requirements
	// until RPC forwarding is fully implemented.
	notification := &protocol.TurmsNotification{
		Timestamp: time.Now().UnixMilli(),
		Code:      proto.Int32(int32(constant.ResponseStatusCode_OK)),
	}
	return notification, nil
}

// TODO: Implement getNotificationFromResponse
// @MappedFrom getNotificationFromResponse(@NotNull ServiceResponse response, long requestId)
func (s *ServiceRequestService) getNotificationFromResponse(response any, requestId int64) *protocol.TurmsNotification {
	// TODO: Map the ServiceResponse to TurmsNotification builder (ClientMessagePool in Java)
	// Setting code, reason, data based on the business logic.
	return nil
}
