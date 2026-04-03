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
func (s *ServiceRequestService) HandleServiceRequest(ctx context.Context, defaultSession *session.UserSession, serviceRequest any) (*protocol.TurmsNotification, error) {
	// Update request timestamp
	defaultSession.SetLastRequestTimestampToNow()

	// TODO: Obtain buffer from serviceRequest and retain it (like serviceRequest.getTurmsRequestBuffer().retain())
	
	// TODO: Construct HandleServiceRequest wrapper
	// request := newHandleServiceRequest(serviceRequest)
	
	// TODO: Replace with proper RPC call when 'Node' cluster infra is implemented
	// Example port: 
	// response, err := s.node.GetRpcService().RequestResponse(request)
	// if err != nil { return nil, err }
	
	// TODO: defer release buffer
	// defer serviceRequest.getTurmsRequestBuffer().release()

	// TODO: Parse response properly using getNotificationFromResponse
	// For now, return a mocked NO_CONTENT or empty TurmsNotification
	var notification protocol.TurmsNotification
	return &notification, nil
}

// TODO: Implement getNotificationFromResponse
// @MappedFrom getNotificationFromResponse(@NotNull ServiceResponse response, long requestId)
func (s *ServiceRequestService) getNotificationFromResponse(response any, requestId int64) *protocol.TurmsNotification {
	// TODO: Map the ServiceResponse to TurmsNotification builder (ClientMessagePool in Java)
	// Setting code, reason, data based on the business logic.
	return nil
}
