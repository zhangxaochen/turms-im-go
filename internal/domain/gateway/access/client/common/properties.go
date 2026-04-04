package common


type TcpProperties struct {
	Enabled               bool
	Host                  string
	Port                  int
	ProxyProtocolMode     RemoteAddressSourceProxyProtocolMode
	MaxFrameLength        int
	ConnectTimeoutMillis int
	Backlog              int
	Wiretap              bool
	Ssl                   SslProperties
}

type RemoteAddressSourceProxyProtocolMode int

const (
	ProxyProtocolMode_DISABLED RemoteAddressSourceProxyProtocolMode = iota
	ProxyProtocolMode_OPTIONAL
	ProxyProtocolMode_REQUIRED
)

type SslProperties struct {
	Enabled bool
}

type WebSocketProperties struct {
	Enabled            bool
	Host               string
	Port               int
	ProxyProtocolMode  RemoteAddressSourceProxyProtocolMode
	HttpHeaderMode     RemoteAddressSourceHttpHeaderMode
	MaxFramePayloadLength int
}

type RemoteAddressSourceHttpHeaderMode int

const (
	HttpHeaderMode_DISABLED RemoteAddressSourceHttpHeaderMode = iota
	HttpHeaderMode_OPTIONAL
	HttpHeaderMode_REQUIRED
)
