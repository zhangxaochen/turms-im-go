package address

import "fmt"

// ServiceAddressManager maps to ServiceAddressManager in Java.
// @MappedFrom ServiceAddressManager
type ServiceAddressManager struct {
	wsAddress  string
	tcpAddress string
	udpAddress string
}

func NewServiceAddressManager() *ServiceAddressManager {
	return &ServiceAddressManager{}
}

// @MappedFrom getWsAddress()
func (m *ServiceAddressManager) GetWsAddress() string {
	return m.wsAddress
}

// @MappedFrom getTcpAddress()
func (m *ServiceAddressManager) GetTcpAddress() string {
	return m.tcpAddress
}

// @MappedFrom getUdpAddress()
func (m *ServiceAddressManager) GetUdpAddress() string {
	return m.udpAddress
}

// @MappedFrom updateCustomAddresses()
func (m *ServiceAddressManager) UpdateCustomAddresses(
	advertiseStrategy string,
	advertiseHost string,
	attachPortToHost bool,
	wsHost string, wsPort int, wsSslEnabled bool,
	tcpHost string, tcpPort int,
	udpHost string, udpPort int,
) {
	// For WS
	resolvedWsHost := m.queryHost(advertiseStrategy, wsHost, advertiseHost)
	scheme := "ws://"
	if wsSslEnabled {
		scheme = "wss://"
	}
	if attachPortToHost {
		m.wsAddress = fmt.Sprintf("%s%s:%d", scheme, resolvedWsHost, wsPort)
	} else {
		m.wsAddress = fmt.Sprintf("%s%s", scheme, resolvedWsHost)
	}

	// For TCP
	resolvedTcpHost := m.queryHost(advertiseStrategy, tcpHost, advertiseHost)
	if attachPortToHost {
		m.tcpAddress = fmt.Sprintf("%s:%d", resolvedTcpHost, tcpPort)
	} else {
		m.tcpAddress = resolvedTcpHost
	}

	// For UDP
	resolvedUdpHost := m.queryHost(advertiseStrategy, udpHost, advertiseHost)
	if attachPortToHost {
		m.udpAddress = fmt.Sprintf("%s:%d", resolvedUdpHost, udpPort)
	} else {
		m.udpAddress = resolvedUdpHost
	}
}

func (m *ServiceAddressManager) queryHost(advertiseStrategy, host, advertiseHost string) string {
	// Basic simulation of Java's queryHost
	if advertiseStrategy == "ADVERTISE_ADDRESS" && advertiseHost != "" {
		return advertiseHost
	}
	if host != "" {
		return host
	}
	return "127.0.0.1" // Fallback local address
}
