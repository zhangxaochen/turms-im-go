package ws

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"im.turms/server/internal/domain/common/constant"
	"im.turms/server/internal/domain/gateway/access/client/common"
	"im.turms/server/internal/domain/gateway/access/client/udp"
	"im.turms/server/internal/infra/exception"
)

// @MappedFrom HttpForwardedHeaderHandler
type HttpForwardedHeaderHandler struct {
	isForwardedIpRequired bool
}

func NewHttpForwardedHeaderHandler(isForwardedIpRequired bool) *HttpForwardedHeaderHandler {
	return &HttpForwardedHeaderHandler{
		isForwardedIpRequired: isForwardedIpRequired,
	}
}

var (
	forRegex   = regexp.MustCompile(`(?i)for="?([^";,]+)"?`)
	protoRegex = regexp.MustCompile(`(?i)proto="?([^";,]+)"?`)
	hostRegex  = regexp.MustCompile(`(?i)host="?([^";,]+)"?`)
)

func (h *HttpForwardedHeaderHandler) Apply(r *http.Request) error {
	// 1. Forwarded header
	if forwarded := r.Header.Get("Forwarded"); forwarded != "" {
		return h.parseForwardedInfo(r, forwarded)
	}
	// 2. X-Forwarded-* headers
	return h.parseXForwardedInfo(r)
}

func (h *HttpForwardedHeaderHandler) parseForwardedInfo(r *http.Request, forwarded string) error {
	// Take first entry
	part := strings.Split(forwarded, ",")[0]
	
	if match := forRegex.FindStringSubmatch(part); len(match) > 1 {
		r.RemoteAddr = match[1]
	} else if h.isForwardedIpRequired {
		return fmt.Errorf("forwarded IP is required but not found in Forwarded header")
	}

	if match := protoRegex.FindStringSubmatch(part); len(match) > 1 {
		r.URL.Scheme = strings.ToLower(match[1])
	}
	if match := hostRegex.FindStringSubmatch(part); len(match) > 1 {
		r.URL.Host = match[1]
	}
	return nil
}

func (h *HttpForwardedHeaderHandler) parseXForwardedInfo(r *http.Request) error {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		r.RemoteAddr = strings.TrimSpace(strings.Split(xff, ",")[0])
	} else if h.isForwardedIpRequired {
		return fmt.Errorf("forwarded IP is required but not found in X-Forwarded-For header")
	}

	if xfp := r.Header.Get("X-Forwarded-Proto"); xfp != "" {
		r.URL.Scheme = strings.ToLower(strings.TrimSpace(strings.Split(xfp, ",")[0]))
	}
	if xfh := r.Header.Get("X-Forwarded-Host"); xfh != "" {
		r.URL.Host = strings.TrimSpace(strings.Split(xfh, ",")[0])
	}
	if xfport := r.Header.Get("X-Forwarded-Port"); xfport != "" {
		portStr := strings.TrimSpace(strings.Split(xfport, ",")[0])
		if _, err := strconv.Atoi(portStr); err == nil {
			if !strings.Contains(r.URL.Host, ":") {
				r.URL.Host = fmt.Sprintf("%s:%s", r.URL.Host, portStr)
			}
		}
	}
	return nil
}

// @MappedFrom WSConnection (corresponds to WebSocketConnection.java)
type WSConnection struct {
	*common.BaseNetConnection
	conn         *websocket.Conn
	closeTimeout time.Duration
	onClose      chan struct{}
}

func NewWSConnection(conn *websocket.Conn, isConnected bool, closeTimeout time.Duration, onClose chan struct{}) *WSConnection {
	return &WSConnection{
		BaseNetConnection: common.NewBaseNetConnection(isConnected),
		conn:              conn,
		closeTimeout:      closeTimeout,
		onClose:           onClose,
	}
}

func (c *WSConnection) GetAddress() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *WSConnection) Send(ctx context.Context, buffer []byte) error {
	c.conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	return c.conn.WriteMessage(websocket.BinaryMessage, buffer)
}

func (c *WSConnection) CloseWithReason(reason common.CloseReason) bool {
	if !c.BaseNetConnection.CloseWithReason(reason) {
		return false
	}

	go func() {
		if reason.IsNotifyClient {
			for i := 0; i < 3; i++ {
				nf := common.NewNotificationFactory(nil)
				payload, err := nf.CreateCloseReasonBuffer(reason)
				if err == nil {
					err = c.Send(context.Background(), payload)
				}
				if err == nil || exception.IsDisconnectedClientError(err) {
					break
				}
				time.Sleep(3 * time.Second)
			}
		}

		if c.closeTimeout == 0 {
			c.Close()
		} else if c.closeTimeout > 0 {
			select {
			case <-time.After(c.closeTimeout):
				c.Close()
			case <-c.onClose:
				return
			}
		}
	}()
	return true
}

func (c *WSConnection) Close() error {
	_ = c.BaseNetConnection.Close()
	return c.conn.Close()
}

func (c *WSConnection) SwitchToUdp() {
	c.CloseWithReason(common.NewCloseReason(constant.SessionCloseStatus_SWITCH))
}

// @MappedFrom WebSocketServerFactory
type WebSocketServerFactory struct{}

func (f *WebSocketServerFactory) Create(
	props common.WebSocketProperties,
	blocklistService common.BlocklistService,
	serverStatusManager common.ServerStatusManager,
	sessionService common.SessionService,
	callback func(*websocket.Conn, http.Header),
) (*http.Server, error) {
	
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true // CORS loosely handled here, refined in handler
		},
	}

	forwardedHandler := NewHttpForwardedHeaderHandler(props.HttpHeaderMode == common.HttpHeaderMode_REQUIRED)
	availabilityHandler := common.NewServiceAvailabilityChannelHandler(blocklistService, serverStatusManager, sessionService)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// CORS Preflight
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Proxy headers
		if props.HttpHeaderMode != common.HttpHeaderMode_DISABLED {
			if err := forwardedHandler.Apply(r); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		// Service availability
		if !availabilityHandler.HandleConnection(&dummyNetConn{r}) {
			http.Error(w, "Service unavailable or IP blocked", http.StatusServiceUnavailable)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Upgrade error: %v", err)
			return
		}
		
		go callback(conn, r.Header)
	})

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", props.Host, props.Port),
		Handler: mux,
	}

	return server, nil
}

// @MappedFrom WebSocketUserSessionAssembler
type WebSocketUserSessionAssembler struct {
	Host           string
	Port           int
	Server         *http.Server
}

func NewWebSocketUserSessionAssembler() *WebSocketUserSessionAssembler {
	return &WebSocketUserSessionAssembler{
		Host: "",
		Port: -1,
	}
}

func (a *WebSocketUserSessionAssembler) GetHost() (string, error) {
	if a.Server == nil {
		return "", &exception.FeatureDisabledError{Message: "WebSocket server is disabled"}
	}
	return a.Host, nil
}

func (a *WebSocketUserSessionAssembler) GetPort() (int, error) {
	if a.Server == nil {
		return -1, &exception.FeatureDisabledError{Message: "WebSocket server is disabled"}
	}
	return a.Port, nil
}

func (a *WebSocketUserSessionAssembler) CreateConnection(conn *websocket.Conn, closeTimeout time.Duration, onClose chan struct{}) common.NetConnection {
	c := NewWSConnection(conn, true, closeTimeout, onClose)
	c.SetUdpSignalDispatcher(func(addr *net.UDPAddr) {
		if udp.Instance != nil {
			udp.Instance.SendSignal(addr, udp.OpenConnectionNotification)
		}
	})
	return c
}

type dummyNetConn struct {
	r *http.Request
}

func (d *dummyNetConn) Read(b []byte) (n int, err error)   { return 0, io.EOF }
func (d *dummyNetConn) Write(b []byte) (n int, err error)  { return 0, nil }
func (d *dummyNetConn) Close() error                       { return nil }
func (d *dummyNetConn) LocalAddr() net.Addr                { return nil }
func (d *dummyNetConn) RemoteAddr() net.Addr               { 
	// Basic parsing of r.RemoteAddr
	host, _, _ := net.SplitHostPort(d.r.RemoteAddr)
	if host == "" { host = d.r.RemoteAddr }
	return &net.TCPAddr{IP: net.ParseIP(host)}
}
func (d *dummyNetConn) SetDeadline(t time.Time) error      { return nil }
func (d *dummyNetConn) SetReadDeadline(t time.Time) error  { return nil }
func (d *dummyNetConn) SetWriteDeadline(t time.Time) error { return nil }
