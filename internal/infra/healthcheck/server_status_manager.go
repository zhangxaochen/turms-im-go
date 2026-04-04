package healthcheck

import (
	"im.turms/server/internal/domain/gateway/access/client/common"
)

type ServerStatusManager struct {
}

func NewServerStatusManager() *ServerStatusManager {
	return &ServerStatusManager{}
}

func (m *ServerStatusManager) GetServiceAvailability() common.ServiceAvailability {
	return common.ServiceAvailability{
		Available: true,
		Reason:    "",
	}
}
