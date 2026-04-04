package common

import (
	"github.com/pires/go-proxyproto"
	"net"
)

// WrapWithProxyProtocol wraps a net.Listener with proxyproto connection handling setup.
// This fully handles parsing the PROXY v1/v2 header transparently when Accept() is called.
func WrapWithProxyProtocol(l net.Listener) net.Listener {
	return &proxyproto.Listener{
		Listener: l,
	}
}
