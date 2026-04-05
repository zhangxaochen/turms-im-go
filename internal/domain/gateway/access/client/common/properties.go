package common


type TcpProperties struct {
	Enabled               bool
	Host                  string
	Port                  int
	ConnectTimeoutMillis  int
	IdleTimeoutSeconds    int
	MaxPayloadBytes       int // maps to MaxFrameLength
	ProxyProtocolMode     RemoteAddressSourceProxyProtocolMode
	Ssl                   SslProperties
	KeepAlive             bool
	ReuseAddr             bool
	TcpNoDelay            bool
	Backlog               int
	WriteTimeoutMillis    int
	Wiretap               bool
	MetricsEnabled        bool
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
	Enabled               bool
	Host                  string
	Port                  int
	Backlog               int
	ConnectTimeoutMillis  int
	ProxyProtocolMode     RemoteAddressSourceProxyProtocolMode
	HttpHeaderMode        RemoteAddressSourceHttpHeaderMode
	MaxFramePayloadLength int
	WriteBufferSize       int
	ReadBufferSize        int
	IdleTimeoutSeconds    int
	WriteTimeoutMillis    int
	Wiretap               bool
	MetricsEnabled        bool
	Ssl                   SslProperties
}

type RemoteAddressSourceHttpHeaderMode int

const (
	HttpHeaderMode_DISABLED RemoteAddressSourceHttpHeaderMode = iota
	HttpHeaderMode_OPTIONAL
	HttpHeaderMode_REQUIRED
)
