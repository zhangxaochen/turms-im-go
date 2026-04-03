package server

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"im.turms/server/internal/domain/gateway/session"
)

func TestTCPServer_Lifecycle(t *testing.T) {
	svc := session.NewSessionService(nil, "test-server-id")
	ts := NewTCPServer("127.0.0.1:0", svc, func(ctx context.Context, session *session.UserSession, payload []byte) {
		// handle
	})

	err := ts.Start()
	assert.NoError(t, err)

	addr := ts.listener.Addr().String()
	conn, err := net.Dial("tcp", addr)
	assert.NoError(t, err)

	// close conn
	conn.Close()

	// Normal Stop
	ts.Stop()
}

func TestWSServer_Lifecycle(t *testing.T) {
	svc := session.NewSessionService(nil, "test-server-id")

	ws := NewWSServer("127.0.0.1:0", svc, func(ctx context.Context, session *session.UserSession, payload []byte) {
		// handle
	})

	// Start without binding to port 0 properly, wait we bind to a static port for test to know the addr easily
	// because http.Server with:0 is tricky to get the actual port if it's not a proper listener
	// For testing, let's use a specific random port or start listener manually

	l, err := net.Listen("tcp", "127.0.0.1:0")
	assert.NoError(t, err)

	ws.addr = l.Addr().String()
	ws.httpServer = &http.Server{
		Handler: http.HandlerFunc(ws.handleHTTPFunc),
	}

	go func() {
		_ = ws.httpServer.Serve(l)
	}()

	u := url.URL{Scheme: "ws", Host: ws.addr, Path: "/"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	assert.NoError(t, err)
	defer c.Close()

	c.WriteMessage(websocket.BinaryMessage, []byte("hello"))

	ws.Stop()
}
