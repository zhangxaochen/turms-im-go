package address

import (
	"context"
	"fmt"
)

// ServiceAddressManager maps to ServiceAddressManager in Java.
// @MappedFrom ServiceAddressManager
type ServiceAddressManager struct {
	wsAddress  string
	tcpAddress string
	udpAddress string

	// Cached properties to simulate gatewayApiDiscoveryProperties in Java to avoid unnecessary updates
	lastAdvertiseStrategy string
	lastAdvertiseHost     string
	lastAttachPortToHost  *bool

	ipDetector *IpDetector
}

func NewServiceAddressManager() *ServiceAddressManager {
	return &ServiceAddressManager{
		ipDetector: NewIpDetector(),
	}
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

func (m *ServiceAddressManager) areAddressPropertiesChange(advertiseStrategy, advertiseHost string, attachPortToHost bool) bool {
	if m.lastAttachPortToHost == nil ||
		m.lastAdvertiseStrategy != advertiseStrategy ||
		m.lastAdvertiseHost != advertiseHost ||
		*m.lastAttachPortToHost != attachPortToHost {
		return true
	}
	return false
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
	if !m.areAddressPropertiesChange(advertiseStrategy, advertiseHost, attachPortToHost) {
		return
	}

	m.lastAdvertiseStrategy = advertiseStrategy
	m.lastAdvertiseHost = advertiseHost
	val := attachPortToHost
	m.lastAttachPortToHost = &val

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

func (m *ServiceAddressManager) queryHost(advertiseStrategy, bindHost, advertiseHost string) string {
	switch advertiseStrategy {
	case "ADVERTISE_ADDRESS":
		if advertiseHost != "" {
			return advertiseHost
		}
		panic("The advertised host is not specified")
	case "BIND_ADDRESS":
		if bindHost != "" {
			return bindHost
		}
		panic("The bind host is not specified")
	case "PRIVATE_ADDRESS":
		ip, err := m.ipDetector.QueryPrivateIp(0)
		if err == nil {
			return ip
		}
		panic("Failed to detect the local IP: " + err.Error())
	case "PUBLIC_ADDRESS":
		ip, err := m.ipDetector.QueryPublicIp(context.Background(), []string{"https://checkip.amazonaws.com"}, 0)
		if err == nil && ip != "" {
			return ip
		}
		panic("Failed to detect the public IP")
	default:
		return "127.0.0.1" // Fallback local address
	}
}
