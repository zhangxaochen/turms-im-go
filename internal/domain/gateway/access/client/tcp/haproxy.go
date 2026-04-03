package tcp

import (
	"net"
)

// @MappedFrom ExtendedHAProxyMessageReader
type ExtendedHAProxyMessageReader struct {
	OnRemoteAddressConfirmed func(net.Addr)
}

func NewExtendedHAProxyMessageReader(callback func(net.Addr)) *ExtendedHAProxyMessageReader {
	return &ExtendedHAProxyMessageReader{
		OnRemoteAddressConfirmed: callback,
	}
}

// Read processes PROXY protocol headers and triggers the callback.
// Maps to Netty channelRead handling HAProxyMessage.
// @MappedFrom channelRead(ChannelHandlerContext ctx, Object msg)
func (r *ExtendedHAProxyMessageReader) Read(conn net.Conn) error {
	// Pending implementation: Read initial bytes to parse PROXY v1/v2 headers.
	// For now, immediately confirm using the underlying connection's remote address.
	if r.OnRemoteAddressConfirmed != nil {
		r.OnRemoteAddressConfirmed(conn.RemoteAddr())
	}
	return nil
}

// @MappedFrom HAProxyUtil
type HAProxyUtil struct{}

// @MappedFrom addProxyProtocolHandlers(ChannelPipeline pipeline, Consumer<InetSocketAddress> onRemoteAddressConfirmed)
func AddProxyProtocolHandlers(callback func(net.Addr)) {
	// Pending implementation: Integrate with Go's net package or custom pipeline
}

// @MappedFrom addProxyProtocolDetectorHandler(ChannelPipeline pipeline, Consumer<InetSocketAddress> onRemoteAddressConfirmed)
func AddProxyProtocolDetectorHandler(callback func(net.Addr)) {
	// Pending implementation: Integrate with Go's net package or custom pipeline
}
