package server

import (
	"bufio"
	"context"
	"encoding/binary"
	"log"
	"net"
	"sync"

	"github.com/pires/go-proxyproto"
	"im.turms/server/internal/domain/common/constant"
	"im.turms/server/internal/domain/common/infra/cluster/rpc/codec"
	"im.turms/server/internal/domain/gateway/session"
)

// TCPConnection wraps net.Conn to implement session.Connection
type TCPConnection struct {
	conn       net.Conn
	mu         sync.Mutex
	remoteAddr net.Addr
}

func (c *TCPConnection) Connect() error {
	return nil
}

func (c *TCPConnection) Send(payload []byte) error {
	return c.SendWithContext(context.Background(), payload)
}

func (c *TCPConnection) SendWithContext(ctx context.Context, payload []byte) error {
	// TCP requires Varint length prefix
	buf := make([]byte, binary.MaxVarintLen32+len(payload))
	n := binary.PutUvarint(buf, uint64(len(payload)))
	copy(buf[n:], payload)

	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if context is already done before sending
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	_, err := c.conn.Write(buf[:n+len(payload)])
	return err
}

func (c *TCPConnection) Close(reason constant.SessionCloseStatus) error {
	return c.conn.Close()
}

func (c *TCPConnection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *TCPConnection) TryNotifyClientToRecover() {}

func (c *TCPConnection) IsActive() bool {
	return c.conn != nil
}

// TCPServer listens for incoming TCP connections and handles them.
type TCPServer struct {
	addr           string
	listener       net.Listener
	handler        session.MessageHandler
	sessionService *session.SessionService

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewTCPServer(addr string, srv *session.SessionService, handler session.MessageHandler) *TCPServer {
	ctx, cancel := context.WithCancel(context.Background())
	return &TCPServer{
		addr:           addr,
		handler:        handler,
		sessionService: srv,
		ctx:            ctx,
		cancel:         cancel,
	}
}

func (s *TCPServer) ListenerAddr() string {
	if s.listener != nil {
		return s.listener.Addr().String()
	}
	return s.addr
}

func (s *TCPServer) Start() error {
	var err error
	l, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	s.listener = &proxyproto.Listener{Listener: l}

	go s.acceptLoop()
	return nil
}

func (s *TCPServer) acceptLoop() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.ctx.Done():
				return // Normal shutdown
			default:
				log.Printf("TCPServer accept error: %v", err)
				continue
			}
		}

		s.wg.Add(1)
		go s.handleConnection(conn)
	}
}

func (s *TCPServer) handleConnection(conn net.Conn) {
	defer s.wg.Done()

	tcpConn := &TCPConnection{conn: conn, remoteAddr: conn.RemoteAddr()}

	// Pre-create an unauthed session shell.
	// Actual details (UserID, DeviceType) will be populated upon first Auth request.
	userSession := &session.UserSession{
		Conn:      tcpConn,
		CloseChan: make(chan struct{}),
		IP:        conn.RemoteAddr().(*net.TCPAddr).IP,
	}

	// This goroutine ensures blocking reads are aborted if the server stops or session is closed.
	go func() {
		select {
		case <-s.ctx.Done():
			_ = conn.Close()
		case <-userSession.CloseChan:
			_ = conn.Close()
		}
	}()

	// This is the read loop
	br := bufio.NewReader(conn)
	for {
		// 1. We read the frame. Note: ReadFrame blocks!
		frame, err := codec.ReadFrame(br)
		if err != nil {
			// EOF or read error => disconnect
			_ = conn.Close()
			s.cleanup(userSession)
			return
		}

		// Update heartbeat
		userSession.SetLastHeartbeatRequestTimestampToNow()

		// Route it
		if s.handler != nil {
			go s.handler(s.ctx, userSession, frame.Payload)
		}
	}
}

func (s *TCPServer) cleanup(userSession *session.UserSession) {
	if userSession.UserID > 0 {
		s.sessionService.UnregisterSession(s.ctx, userSession.UserID, userSession.DeviceType, userSession.Conn, constant.SessionCloseStatus_CONNECTION_CLOSED)
	}
}

func (s *TCPServer) Stop() {
	s.cancel()
	if s.listener != nil {
		s.listener.Close()
	}
	s.wg.Wait()
}
