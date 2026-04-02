package server

import (
	"context"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"im.turms/server/internal/domain/gateway/session"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// In production, configure CheckOrigin properly
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WSConnection wraps websocket.Conn to implement session.Connection
type WSConnection struct {
	conn *websocket.Conn
	mu   sync.Mutex
}

func (c *WSConnection) WriteMessage(payload []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	// Turms uses BinaryMessage to carry Protobuf payloads
	return c.conn.WriteMessage(websocket.BinaryMessage, payload)
}

func (c *WSConnection) Close() error {
	return c.conn.Close()
}

func (c *WSConnection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

type WSServer struct {
	addr           string
	httpServer     *http.Server
	handler        session.MessageHandler
	sessionService *session.SessionService

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewWSServer(addr string, srv *session.SessionService, handler session.MessageHandler) *WSServer {
	ctx, cancel := context.WithCancel(context.Background())
	return &WSServer{
		addr:           addr,
		handler:        handler,
		sessionService: srv,
		ctx:            ctx,
		cancel:         cancel,
	}
}

func (s *WSServer) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleHTTPFunc)

	s.httpServer = &http.Server{
		Addr:    s.addr,
		Handler: mux,
	}

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("WSServer ListenAndServe error: %v", err)
		}
	}()

	return nil
}

func (s *WSServer) handleHTTPFunc(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WSServer upgrade error: %v", err)
		return
	}

	s.wg.Add(1)
	go s.handleConnection(conn, r)
}

func (s *WSServer) handleConnection(conn *websocket.Conn, r *http.Request) {
	defer s.wg.Done()

	wsConn := &WSConnection{conn: conn}

	ipStr, _, _ := net.SplitHostPort(r.RemoteAddr)

	userSession := &session.UserSession{
		Conn:      wsConn,
		CloseChan: make(chan struct{}),
		IP:        net.ParseIP(ipStr),
	}

	// Watcher for forced closes
	go func() {
		select {
		case <-s.ctx.Done():
			_ = conn.Close()
		case <-userSession.CloseChan:
			_ = conn.Close()
		}
	}()

	for {
		// Gorilla typically forces blocking read on ReadMessage
		messageType, payload, err := conn.ReadMessage()
		if err != nil {
			_ = conn.Close()
			s.cleanup(userSession)
			return
		}

		// Only process Binary Messages containing TurmsRequest Protobuf
		if messageType == websocket.BinaryMessage {
			userSession.SetLastHeartbeatRequestTimestampToNow()

			if s.handler != nil {
				go s.handler(s.ctx, userSession, payload)
			}
		}
	}
}

func (s *WSServer) cleanup(userSession *session.UserSession) {
	if userSession.UserID > 0 {
		s.sessionService.UnregisterSession(userSession.UserID, userSession.DeviceType, userSession.Conn)
	}
}

func (s *WSServer) Stop() {
	s.cancel()
	if s.httpServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = s.httpServer.Shutdown(ctx)
	}
	s.wg.Wait()
}
