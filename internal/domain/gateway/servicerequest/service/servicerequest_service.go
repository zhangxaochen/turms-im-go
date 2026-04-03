package service

import (
	"context"

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
func (s *ServiceRequestService) HandleServiceRequest(ctx context.Context, session *session.UserSession, serviceRequest any) (*protocol.TurmsNotification, error) {
	// 1. Update request timestamp
	session.SetLastRequestTimestampToNow()

	// 2. Forward request (Wait for cluster RPC implementation to be ready in Go)
	// We return an empty or dummy notification for now or forward to the proper Node RPC stub.
	var notification protocol.TurmsNotification
	return &notification, nil
}
