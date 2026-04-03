package connection

import (
	"im.turms/server/internal/domain/common/infra/cluster/discovery"
)

// TurmsConnection represents an abstract connection to a remote cluster node.
type TurmsConnection interface {
	NodeID() string
	IsConnected() bool
	// SendRequest(...)
}

// ConnectionService tracks and manages connections to other cluster members.
type ConnectionService struct {
	discovery *discovery.DiscoveryService
}

func NewConnectionService(discovery *discovery.DiscoveryService) *ConnectionService {
	return &ConnectionService{
		discovery: discovery,
	}
}

func (s *ConnectionService) GetMemberConnection(nodeID string) TurmsConnection {
	// TODO: Return connection for the node when TCP connection pool is implemented
	return nil
}

func (s *ConnectionService) IsHasConnectedToAllMembers() bool {
	// TODO: Iterate through expected members and check connection statuses
	return false
}
