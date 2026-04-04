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

	"im.turms/server/internal/domain/common/constant"
	"im.turms/server/internal/domain/gateway/access/client/common"
	"im.turms/server/internal/domain/gateway/access/client/udp"
	"im.turms/server/internal/infra/exception"
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

	go func() {
		if reason.IsNotifyClient {
			// Try to send notification up to 3 times (initial + 2 retries) with short backoff
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
				if exception.IsDisconnectedClientError(err) {
					break
				}
				if i < 2 {
					log.Printf("Failed to send the close notification attempt %d: %v", i+1, err)
					time.Sleep(3 * time.Second)
				} else {
					log.Printf("Failed to send the close notification after 2 retries: %v", err)
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
		// If c.closeTimeout < 0, we do not close the underlying connection forcefully.
	}()
	return true
}

// @MappedFrom close()
func (c *TcpConnection) Close() error {
	// Java doesn't check isConnected before disposing.
	// But in Go, multiple closes of a net.Conn are safe but return error.
	// We call c.conn.Close() directly to match Java's "disposeNow" logic.
	err := c.conn.Close()
	if err != nil && !exception.IsDisconnectedClientError(err) {
		log.Printf("Failed to close the TCP connection %s: %v", c.GetAddress(), err)
	}
	// Note: We don't call super.close() here because Java doesn't either.
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
	host string,
	port int,
	proxy bool,
	maxFrameLength int,
	blocklistService common.BlocklistService,
	serverStatusManager common.ServerStatusManager,
	sessionService common.SessionService,
	callback func(net.Conn),
) (net.Listener, error) {
	addr := fmt.Sprintf("%s:%d", host, port)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("Failed to bind the TCP server on: %s. Error: %w", addr, err)
	}

	if proxy {
		l = WrapWithProxyProtocol(l)
	}

	availabilityHandler := common.NewServiceAvailabilityChannelHandler(blocklistService, serverStatusManager, sessionService)

	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				return
			}

			// Apply socket options
			if tcpConn, ok := conn.(*net.TCPConn); ok {
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
