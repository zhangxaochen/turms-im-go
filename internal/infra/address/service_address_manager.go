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
	wsHost string, wsPort int, wsSslEnabled bool, wsEnabled bool,
	tcpHost string, tcpPort int, tcpEnabled bool,
	udpHost string, udpPort int, udpEnabled bool,
) {
	// For WS - only update when WS is enabled (matches Java)
	if wsEnabled {
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
	}

	// For TCP - only update when TCP is enabled (matches Java)
	if tcpEnabled {
		resolvedTcpHost := m.queryHost(advertiseStrategy, tcpHost, advertiseHost)
		if attachPortToHost {
			m.tcpAddress = fmt.Sprintf("%s:%d", resolvedTcpHost, tcpPort)
		} else {
			m.tcpAddress = resolvedTcpHost
		}
	}

	// For UDP - only update when UDP is enabled (matches Java)
	if udpEnabled {
		resolvedUdpHost := m.queryHost(advertiseStrategy, udpHost, advertiseHost)
		if attachPortToHost {
			m.udpAddress = fmt.Sprintf("%s:%d", resolvedUdpHost, udpPort)
		} else {
			m.udpAddress = resolvedUdpHost
		}
	}
}

func (m *ServiceAddressManager) queryHost(advertiseStrategy, host, advertiseHost string) string {
	switch advertiseStrategy {
	case "ADVERTISE_ADDRESS":
		if advertiseHost == "" {
			panic("The advertised host is not specified")
		}
		return advertiseHost
	case "BIND_ADDRESS":
		if host == "" {
			panic("The bind host is not specified")
		}
		return host
	case "PRIVATE_ADDRESS":
		// TODO: Query local private IP via IpDetector.queryPrivateIp()
		// Falls back to host parameter for now
		if host != "" {
			return host
		}
		return "127.0.0.1"
	case "PUBLIC_ADDRESS":
		// TODO: Query public IP via IpDetector.queryPublicIp()
		// Falls back to host parameter for now
		if host != "" {
			return host
		}
		return "127.0.0.1"
	default:
		if host != "" {
			return host
		}
		return "127.0.0.1"
	}
}
