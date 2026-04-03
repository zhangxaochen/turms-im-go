package tcp

import (
	"context"
	"net"
	"time"

	"im.turms/server/internal/domain/gateway/access/client/common"
)

// @MappedFrom TcpConnection
type TcpConnection struct {
	*common.BaseNetConnection
	conn         net.Conn
	closeTimeout time.Duration
}

func NewTcpConnection(conn net.Conn, isConnected bool, closeTimeout time.Duration) *TcpConnection {
	return &TcpConnection{
		BaseNetConnection: common.NewBaseNetConnection(isConnected),
		conn:              conn,
		closeTimeout:      closeTimeout,
	}
}

// @MappedFrom getAddress()
func (c *TcpConnection) GetAddress() net.Addr {
	return c.conn.RemoteAddr()
}

// @MappedFrom send(ByteBuf buffer)
func (c *TcpConnection) Send(ctx context.Context, buffer []byte) error {
	_, err := c.conn.Write(buffer)
	return err
}

// @MappedFrom close(CloseReason closeReason)
func (c *TcpConnection) CloseWithReason(reason common.CloseReason) error {
	if !c.IsConnected() {
		return nil
	}
	c.BaseNetConnection.CloseWithReason(reason)
	// Pending logic to send notification before closing, similar to Java
	return c.conn.Close()
}

// @MappedFrom close()
func (c *TcpConnection) Close() error {
	if !c.IsConnected() {
		return nil
	}
	c.BaseNetConnection.Close()
	return c.conn.Close()
}

// @MappedFrom TcpServerFactory
type TcpServerFactory struct{}

// @MappedFrom create(...)
func (f *TcpServerFactory) Create() {
	// Pending implementation: net.Listen and accept loop
}

// @MappedFrom create(TcpProperties tcpProperties, BlocklistService blocklistService, ServerStatusManager serverStatusManager, SessionService sessionService, ConnectionListener connectionListener, int maxFrameLength)
func (f *TcpServerFactory) CreateWithArgs(tcpProperties any, blocklistService any, serverStatusManager any, sessionService any, connectionListener any, maxFrameLength int) {
}

// @MappedFrom TcpUserSessionAssembler
type TcpUserSessionAssembler struct {
	Enabled bool
	Host    string
	Port    int
	Server  net.Listener
}

func NewTcpUserSessionAssembler() *TcpUserSessionAssembler {
	return &TcpUserSessionAssembler{
		Enabled: false,
		Host:    "",
		Port:    -1,
	}
}

// @MappedFrom getHost()
func (a *TcpUserSessionAssembler) GetHost() string {
	if !a.Enabled {
		panic("TCP server is disabled")
	}
	return a.Host
}

// @MappedFrom getPort()
func (a *TcpUserSessionAssembler) GetPort() int {
	if !a.Enabled {
		panic("TCP server is disabled")
	}
	return a.Port
}

// @MappedFrom createConnection(Connection connection, Duration closeTimeout)
func (a *TcpUserSessionAssembler) CreateConnection(conn net.Conn, closeTimeout time.Duration) common.NetConnection {
	return NewTcpConnection(conn, true, closeTimeout)
}
