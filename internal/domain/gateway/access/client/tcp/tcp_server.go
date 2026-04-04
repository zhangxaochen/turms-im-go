package tcp

import (
	"bufio"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/pires/go-proxyproto"
	"im.turms/server/internal/domain/common/constant"
	"im.turms/server/internal/domain/gateway/access/client/common"
	"im.turms/server/internal/domain/gateway/access/client/udp"
	"im.turms/server/internal/infra/exception"
	"syscall"
)

// @MappedFrom TcpConnection
type TcpConnection struct {
	*common.BaseNetConnection
	conn           net.Conn
	closeTimeout   time.Duration
	maxFrameLength int
	onClose        chan struct{}
}

func NewTcpConnection(conn net.Conn, isConnected bool, closeTimeout time.Duration, maxFrameLength int, onClose chan struct{}) *TcpConnection {
	return &TcpConnection{
		BaseNetConnection: common.NewBaseNetConnection(isConnected),
		conn:              conn,
		closeTimeout:      closeTimeout,
		maxFrameLength:    maxFrameLength,
		onClose:           onClose,
	}
}

// @MappedFrom getAddress()
func (c *TcpConnection) GetAddress() net.Addr {
	return c.conn.RemoteAddr()
}

// @MappedFrom send(ByteBuf buffer)
func (c *TcpConnection) Send(ctx context.Context, buffer []byte) error {
	// Prepend varint length
	length := len(buffer)
	var varint [binary.MaxVarintLen64]byte
	n := binary.PutUvarint(varint[:], uint64(length))

	c.conn.SetWriteDeadline(time.Now().Add(5 * time.Second))

	// Combined write if possible, or just two writes.
	// For small buffers, combining is better.
	totalLen := n + length
	combined := make([]byte, totalLen)
	copy(combined, varint[:n])
	copy(combined[n:], buffer)

	_, err := c.conn.Write(combined)
	return err
}

func (c *TcpConnection) Start(onMessage func(common.NetConnection, []byte)) {
	defer c.Close()

	reader := bufio.NewReader(c.conn)
	for {
		// Read varint length
		length, err := binary.ReadUvarint(reader)
		if err != nil {
			if !exception.IsDisconnectedClientError(err) {
				log.Printf("Failed to read varint length: %v", err)
			}
			return
		}

		if length > uint64(c.maxFrameLength) {
			log.Printf("Frame size %d exceeds maxFrameLength %d", length, c.maxFrameLength)
			return
		}

		// Read payload
		payload := make([]byte, length)
		_, err = io.ReadFull(reader, payload)
		if err != nil {
			if !exception.IsDisconnectedClientError(err) {
				log.Printf("Failed to read payload: %v", err)
			}
			return
		}

		onMessage(c, payload)
	}
}

// @MappedFrom close(CloseReason closeReason)
func (c *TcpConnection) CloseWithReason(reason common.CloseReason) bool {
	if !c.BaseNetConnection.CloseWithReason(reason) {
		return false
	}

	// We don't use a separate goroutine if we want to ensure ordering, 
	// but Java's close(CloseReason) sends a notification and wait for closeTimeout.
	go func() {
		if reason.IsNotifyClient {
			// Try to send notification with backoff filter for disconnected-client errors
			for i := 0; i < 3; i++ {
				nf := common.NewNotificationFactory(nil)
				payload, err := nf.CreateCloseReasonBuffer(reason)
				if err != nil {
					log.Printf("Failed to marshal close notification: %v", err)
					break
				}

				err = c.Send(context.Background(), payload)
				if err == nil {
					break
				}
				// Java's RETRY_SEND_CLOSE_NOTIFICATION filters out disconnected client errors
				if exception.IsDisconnectedClientError(err) {
					break
				}
				if i < 2 {
					log.Printf("Failed to send the close notification attempt %d: %v", i+1, err)
					time.Sleep(3 * time.Second)
				}
			}
		}

		if c.closeTimeout == 0 {
			c.Close()
		} else if c.closeTimeout > 0 {
			select {
			case <-time.After(c.closeTimeout):
				c.Close()
			case <-c.onClose:
				// Peer closed first
				return
			}
		}
	}()
	return true
}

// @MappedFrom close()
func (c *TcpConnection) Close() error {
	// Call base close to update flags properly
	_ = c.BaseNetConnection.Close()
	
	err := c.conn.Close()
	if err != nil && !exception.IsDisconnectedClientError(err) {
		log.Printf("Failed to close the TCP connection %s: %v", c.GetAddress(), err)
	}
	return err
}

// @MappedFrom switchToUdp()
func (c *TcpConnection) SwitchToUdp() {
	c.CloseWithReason(common.NewCloseReason(constant.SessionCloseStatus_SWITCH))
}

// @MappedFrom TcpServerFactory
type TcpServerFactory struct{}

// @MappedFrom create(...)
func (f *TcpServerFactory) Create(
	props common.TcpProperties,
	blocklistService common.BlocklistService,
	serverStatusManager common.ServerStatusManager,
	sessionService common.SessionService,
	callback func(net.Conn),
) (net.Listener, error) {
	addr := fmt.Sprintf("%s:%d", props.Host, props.Port)
	
	// Use ListenConfig to set backlog and other options if possible
	lc := net.ListenConfig{
		Control: func(network, address string, c syscall.RawConn) error {
			return c.Control(func(fd uintptr) {
				// SO_REUSEADDR is usually default in Go for unix, 
				// but SO_KEEPALIVE or others can be set here if needed.
			})
		},
	}
	l, err := lc.Listen(context.Background(), "tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("Failed to bind the TCP server on: %s. Error: %w", addr, err)
	}

	if props.Ssl.Enabled {
		// Note: We'd need actual cert/key paths in SslProperties to implement this fully.
		// For now, aligning with the pattern but keeping it as a stub or placeholder.
		// tlsConfig := &tls.Config{}
		// l = tls.NewListener(l, tlsConfig)
	}

	if props.ProxyProtocolMode != common.ProxyProtocolMode_DISABLED {
		pL := &proxyproto.Listener{Listener: l}
		if props.ProxyProtocolMode == common.ProxyProtocolMode_REQUIRED {
			pL.Policy = func(upstream net.Addr) (proxyproto.Policy, error) {
				return proxyproto.REQUIRE, nil
			}
		} else {
			pL.Policy = func(upstream net.Addr) (proxyproto.Policy, error) {
				return proxyproto.USE, nil
			}
		}
		l = pL
	}

	availabilityHandler := common.NewServiceAvailabilityChannelHandler(blocklistService, serverStatusManager, sessionService)

	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				return
			}

			// Apply socket options on accepted connection
			if tcpConn, ok := conn.(*net.TCPConn); ok {
				// SO_LINGER=0 ensures the socket is reset immediately rather than staying in TIME_WAIT.
				// This matches Java's netty config for the child channel.
				tcpConn.SetNoDelay(true)
				tcpConn.SetLinger(0)
			}

			// Service availability and blocklist check
			if !availabilityHandler.HandleConnection(conn) {
				conn.Close()
				continue
			}

			go callback(conn)
		}
	}()
	return l, nil
}

// @MappedFrom TcpUserSessionAssembler
type TcpUserSessionAssembler struct {
	Host           string
	Port           int
	Server         net.Listener
	MaxFrameLength int
}

func NewTcpUserSessionAssembler() *TcpUserSessionAssembler {
	return &TcpUserSessionAssembler{
		Host: "",
		Port: -1,
	}
}

// @MappedFrom getHost()
func (a *TcpUserSessionAssembler) GetHost() (string, error) {
	if a.Server == nil {
		return "", &exception.FeatureDisabledError{Message: "TCP server is disabled"}
	}
	return a.Host, nil
}

// @MappedFrom getPort()
func (a *TcpUserSessionAssembler) GetPort() (int, error) {
	if a.Server == nil {
		return -1, &exception.FeatureDisabledError{Message: "TCP server is disabled"}
	}
	return a.Port, nil
}

// @MappedFrom createConnection(Connection connection, Duration closeTimeout)
func (a *TcpUserSessionAssembler) CreateConnection(conn net.Conn, closeTimeout time.Duration, onClose chan struct{}) common.NetConnection {
	c := NewTcpConnection(conn, true, closeTimeout, a.MaxFrameLength, onClose)
	c.SetUdpSignalDispatcher(func(addr *net.UDPAddr) {
		if udp.Instance != nil {
			udp.Instance.SendSignal(addr, udp.OpenConnectionNotification)
		}
	})
	return c
}
