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
	sessionbo "im.turms/server/internal/domain/gateway/session/bo"
	"im.turms/server/internal/infra/exception"
	"syscall"
)

// @MappedFrom TcpConnection
type TcpConnection struct {
	*common.BaseNetConnection
	conn           net.Conn
	closeTimeout   time.Duration
	writeTimeout   time.Duration
	readTimeout    time.Duration
	maxFrameLength int
	onClose        chan struct{}
}

func NewTcpConnection(conn net.Conn, isConnected bool, closeTimeout time.Duration, writeTimeout time.Duration, readTimeout time.Duration, maxFrameLength int, onClose chan struct{}) *TcpConnection {
	return &TcpConnection{
		BaseNetConnection: common.NewBaseNetConnection(isConnected),
		conn:              conn,
		closeTimeout:      closeTimeout,
		writeTimeout:      writeTimeout,
		readTimeout:       readTimeout,
		maxFrameLength:    maxFrameLength,
		onClose:           onClose,
	}
}

// @MappedFrom getAddress()
func (c *TcpConnection) GetAddress() net.Addr {
	return c.conn.RemoteAddr()
}

// @MappedFrom send(ByteBuf buffer)
func (c *TcpConnection) Send(buffer []byte) error {
	return c.SendWithContext(context.Background(), buffer)
}

func (c *TcpConnection) SendWithContext(ctx context.Context, buffer []byte) error {
	if !c.IsConnected() {
		return nil
	}

	// Prepend varint length
	length := len(buffer)
	var varint [binary.MaxVarintLen64]byte
	n := binary.PutUvarint(varint[:], uint64(length))

	timeout := c.writeTimeout
	if timeout > 0 {
		_ = c.conn.SetWriteDeadline(time.Now().Add(timeout))
	} else {
		_ = c.conn.SetWriteDeadline(time.Time{})
	}

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
	defer func() {
		_ = c.Close()
		if c.onClose != nil {
			select {
			case <-c.onClose:
				// Already closed
			default:
				close(c.onClose)
			}
		}
	}()

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
		err = c.conn.SetReadDeadline(time.Now().Add(c.readTimeout))
		if err != nil {
			log.Printf("Failed to set read deadline: %v", err)
			return
		}
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
func (c *TcpConnection) CloseWithReason(reason sessionbo.CloseReason) bool {
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

				err = c.SendWithContext(context.Background(), payload)
				if err == nil {
					break
				}
				// Java's RETRY_SEND_CLOSE_NOTIFICATION filters out disconnected client errors
				if exception.IsDisconnectedClientError(err) {
					break
				}
				if i < 2 {
					if !exception.IsDisconnectedClientError(err) {
						log.Printf("Failed to send the close notification attempt %d: %v", i+1, err)
					}
					time.Sleep(3 * time.Second)
				} else {
					if !exception.IsDisconnectedClientError(err) {
						log.Printf("Failed to send the close notification after 3 attempts: %v", err)
					}
				}
			}
		}

		if c.closeTimeout == 0 {
			c.Close()
		} else if c.closeTimeout > 0 {
			// Bug 357: Wait for Peer terminate / onTerminate equivalent
			select {
			case <-time.After(c.closeTimeout):
				_ = c.Close()
			case <-c.onClose:
				// Peer closed first, cleanup via Close()
				_ = c.Close()
			}
		}
	}()
	return true
}

func (c *TcpConnection) IsDisposed() bool {
	return c.BaseNetConnection.IsDisposed()
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
	c.CloseWithReason(sessionbo.NewCloseReason(constant.SessionCloseStatus_SWITCH))
}

// @MappedFrom TcpServerFactory
type TcpServerFactory struct{}

// @MappedFrom create(...)
func (f *TcpServerFactory) Create(
	props common.TcpProperties,
	blocklistService common.BlocklistService,
	serverStatusManager common.ServerStatusManager,
	sessionService common.SessionService,
	callback func(net.Conn, int), // Pass maxFrameLength
) (net.Listener, error) {
	addr := fmt.Sprintf("%s:%d", props.Host, props.Port)
	
	lc := net.ListenConfig{
		Control: func(network, address string, c syscall.RawConn) error {
			return c.Control(func(fd uintptr) {
				// SO_REUSEADDR (Bug 873)
				_ = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
				// TCP Keep-alive
				if props.KeepAlive {
					_ = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_KEEPALIVE, 1)
				}
				// Bug 872: Set backlog to 4096 (matching Turms Java)
				_ = syscall.Listen(int(fd), 4096)
			})
		},
	}
	l, err := lc.Listen(context.Background(), "tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to bind the TCP server on %s: %w", addr, err)
	}

	if props.Ssl.Enabled {
		// Placeholder for SSL initialization to match Java parity (Bug 426)
		/*
			tlsConfig := &tls.Config{...}
			l = tls.NewListener(l, tlsConfig)
		*/
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

	if props.Wiretap {
		l = common.NewWiretapListener(l)
	}

	if props.MetricsEnabled {
		l = common.NewMetricsListener(l, "turms.gateway.server.tcp")
	}

	availabilityHandler := common.NewServiceAvailabilityChannelHandler(blocklistService, serverStatusManager, sessionService)

	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				return
			}

			if tcpConn, ok := conn.(*net.TCPConn); ok {
				// Bug 420: TCP_NODELAY parity
				_ = tcpConn.SetNoDelay(props.TcpNoDelay)
				// Bug 422: SO_LINGER parity - matched Java's default netty config (reset immediately)
				_ = tcpConn.SetLinger(0)
				
				if props.KeepAlive {
					_ = tcpConn.SetKeepAlive(true)
					_ = tcpConn.SetKeepAlivePeriod(1 * time.Minute) // Matching Java's default
				}
			}

			if props.ProxyProtocolMode != common.ProxyProtocolMode_DISABLED {
				AddProxyProtocolHandlers(conn.RemoteAddr(), func(addr net.Addr) {
					// Trigger any address confirmation logic if needed
					// (Bug 325/327: Confirming address after proxy resolution)
				})
			}

			if !availabilityHandler.HandleConnection(conn) {
				_ = conn.Close()
				continue
			}

			go callback(conn, props.MaxPayloadBytes)
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
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
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

func (a *TcpUserSessionAssembler) CreateConnection(conn net.Conn, closeTimeout time.Duration, onClose chan struct{}) common.NetConnection {
	c := NewTcpConnection(conn, true, closeTimeout, a.WriteTimeout, a.ReadTimeout, a.MaxFrameLength, onClose)
	c.SetUdpSignalDispatcher(func(addr *net.UDPAddr) {
		if udp.Instance != nil {
			udp.Instance.SendSignal(addr, udp.OpenConnectionNotification)
		}
	})
	return c
}
