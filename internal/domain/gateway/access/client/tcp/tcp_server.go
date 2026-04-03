package tcp

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"im.turms/server/internal/domain/common/constant"
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
	c.conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	_, err := c.conn.Write(buffer)
	return err
}

// @MappedFrom close(CloseReason closeReason)
func (c *TcpConnection) CloseWithReason(reason common.CloseReason) error {
	if !c.IsConnected() {
		return nil
	}
	c.BaseNetConnection.CloseWithReason(reason)
	
	if reason.Status != constant.SessionCloseStatus_UNKNOWN_ERROR {
		// Try to send notification up to 3 times (initial + 2 retries) with short backoff, imitating Java behavior
		go func() {
			for i := 0; i < 3; i++ {
				// We could send real CloseReason status code here using NotificationFactory
				err := c.Send(context.Background(), []byte{byte(reason.Status)}) 
				if err == nil {
					break
				}
				log.Printf("Failed to send close notification attempt %d: %v", i+1, err)
				time.Sleep(3 * time.Second)
			}
			
			if c.closeTimeout > 0 {
				time.Sleep(c.closeTimeout)
			}
			
			err := c.conn.Close()
			if err != nil {
				log.Printf("Failed to close the TCP connection %s: %v", c.GetAddress(), err)
			}
		}()
		return nil
	}
	
	err := c.conn.Close()
	if err != nil {
		log.Printf("Failed to close the TCP connection %s: %v", c.GetAddress(), err)
	}
	return err
}

// @MappedFrom close()
func (c *TcpConnection) Close() error {
	if !c.IsConnected() {
		return nil
	}
	c.BaseNetConnection.Close()
	err := c.conn.Close()
	if err != nil {
		log.Printf("Failed to close the TCP connection %s: %v", c.GetAddress(), err)
	}
	return err
}

// @MappedFrom TcpServerFactory
type TcpServerFactory struct{}

// @MappedFrom create(...)
func (f *TcpServerFactory) Create(host string, port int, proxy bool, callback func(net.Conn)) (net.Listener, error) {
	addr := fmt.Sprintf("%s:%d", host, port)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	if proxy {
		// Wrap with go-proxyproto if proxy is enabled
		l = WrapWithProxyProtocol(l)
	}

	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				// Typically happens on listener close
				return
			}
			go callback(conn)
		}
	}()
	return l, nil
}

// @MappedFrom TcpUserSessionAssembler
type TcpUserSessionAssembler struct {
	Host    string
	Port    int
	Server  net.Listener
}

func NewTcpUserSessionAssembler() *TcpUserSessionAssembler {
	return &TcpUserSessionAssembler{
		Host:    "",
		Port:    -1,
	}
}

// @MappedFrom getHost()
func (a *TcpUserSessionAssembler) GetHost() (string, error) {
	if a.Server == nil {
		return "", fmt.Errorf("TCP server is disabled")
	}
	return a.Host, nil
}

// @MappedFrom getPort()
func (a *TcpUserSessionAssembler) GetPort() (int, error) {
	if a.Server == nil {
		return -1, fmt.Errorf("TCP server is disabled")
	}
	return a.Port, nil
}

// @MappedFrom createConnection(Connection connection, Duration closeTimeout)
func (a *TcpUserSessionAssembler) CreateConnection(conn net.Conn, closeTimeout time.Duration) common.NetConnection {
	return NewTcpConnection(conn, true, closeTimeout)
}
