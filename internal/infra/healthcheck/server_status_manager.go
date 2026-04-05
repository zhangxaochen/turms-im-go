package healthcheck

import (
	"fmt"

	"im.turms/server/internal/domain/gateway/access/client/common"
)

type ServerStatusManager struct {
	availabilityHandler *common.ServiceAvailabilityHandler
}

func NewServerStatusManager(availabilityHandler *common.ServiceAvailabilityHandler) *ServerStatusManager {
	return &ServerStatusManager{
		availabilityHandler: availabilityHandler,
	}
}

func (m *ServerStatusManager) GetServiceAvailability() common.ServiceAvailability {
	status := m.availabilityHandler.GetStatus()
	return common.ServiceAvailability{
		Available: status == common.StatusRunning,
		Reason:    fmt.Sprintf("The server is in the %s state", status),
	}
}
