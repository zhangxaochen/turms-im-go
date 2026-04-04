package tcp

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

// @MappedFrom ExtendedHAProxyMessageReader
// In the Go version using go-proxyproto, this interceptor-style reader is not strictly needed
// since the wrapping net.Listener parses the headers transparently, but we keep the stub shape
// to align with Java handler pipelines.
type ExtendedHAProxyMessageReader struct {
	OnRemoteAddressConfirmed func(net.Addr)
}

func NewExtendedHAProxyMessageReader(callback func(net.Addr)) *ExtendedHAProxyMessageReader {
	return &ExtendedHAProxyMessageReader{
		OnRemoteAddressConfirmed: callback,
	}
}

// Read triggers the callback immediately as the *proxyproto.Conn already parsed the header
// during Accept().
// @MappedFrom channelRead(ChannelHandlerContext ctx, Object msg)
func (r *ExtendedHAProxyMessageReader) Read(conn net.Conn) error {
	if r.OnRemoteAddressConfirmed != nil {
		r.OnRemoteAddressConfirmed(conn.RemoteAddr())
	}
	return nil
}

// @MappedFrom HAProxyUtil
type HAProxyUtil struct{}

// @MappedFrom addProxyProtocolHandlers
func AddProxyProtocolHandlers(addr net.Addr, callback func(net.Addr)) {
	if callback != nil {
		callback(addr)
	}
}

// @MappedFrom addProxyProtocolDetectorHandler
func AddProxyProtocolDetectorHandler(addr net.Addr, callback func(net.Addr)) {
	if callback != nil {
		callback(addr)
	}
}
