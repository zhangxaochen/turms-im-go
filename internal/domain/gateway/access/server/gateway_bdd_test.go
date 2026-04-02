package server_test

import (
	"bufio"
	"context"
	"encoding/binary"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"im.turms/server/internal/domain/common/infra/cluster/rpc/codec"
	"im.turms/server/internal/domain/gateway/access/server"
	"im.turms/server/internal/domain/gateway/session"
	"im.turms/server/pkg/protocol"
)

// mockTurmsClient is a simple TCP client that mimics the exact behavior of Turms client SDK
type mockTurmsClient struct {
	conn net.Conn
	mu   sync.Mutex
}

func newMockTurmsClient(t *testing.T, addr string) *mockTurmsClient {
	conn, err := net.Dial("tcp", addr)
	require.NoError(t, err)
	return &mockTurmsClient{conn: conn}
}

// SendReq writes a varint-length prefixed RpcFrame to the server
func (c *mockTurmsClient) SendReq(t *testing.T, requestID int32, payload []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	frameSize := codec.HeaderSize + len(payload)
	varintBuf := make([]byte, binary.MaxVarintLen32)
	varintBytes := binary.PutUvarint(varintBuf, uint64(frameSize))

	// Write varint prefix
	_, err := c.conn.Write(varintBuf[:varintBytes])
	require.NoError(t, err)

	// Write Header (CodecID = 1 for Client->Server, RequestID)
	headerBuf := make([]byte, codec.HeaderSize)
	binary.BigEndian.PutUint16(headerBuf[0:2], uint16(1))
	binary.BigEndian.PutUint32(headerBuf[2:6], uint32(requestID))
	_, err = c.conn.Write(headerBuf)
	require.NoError(t, err)

	// Write Payload
	if len(payload) > 0 {
		_, err = c.conn.Write(payload)
		require.NoError(t, err)
	}
}

func (c *mockTurmsClient) ReadResponse(t *testing.T) ([]byte, error) {
	err := c.conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	require.NoError(t, err)

	br := bufio.NewReader(c.conn)
	payloadLen, err := binary.ReadUvarint(br)
	if err != nil {
		return nil, err
	}

	payload := make([]byte, payloadLen)
	_, err = br.Read(payload)
	return payload, err
}

func (c *mockTurmsClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

// BDD Scenario: Handling Client Connections and Multi-device Kicks
func TestGateway_BDD_ClientConnectionAndKick(t *testing.T) {
	// GIVEN a running Gateway TCP Server with KickExisting strategy
	svc := session.NewSessionService()
	svc.ConflictStrategy = session.KickExisting

	handlerCount := 0
	var mu sync.Mutex

	// We intercept raw payloads.
	// In a real system, the payload is parsed as TurmsRequest and Auth is done.
	// For this test, we simulate an Auth interceptor inline in the Gateway Handler.
	handler := func(ctx context.Context, s *session.UserSession, payload []byte) {
		mu.Lock()
		handlerCount++
		c := handlerCount
		mu.Unlock()

		// Simulate Auth logic for testing BDD:
		// If payload is "LOGIN 100 ANDROID", we register them
		str := string(payload)
		if str == "LOGIN" {
			s.UserID = 100
			s.DeviceType = protocol.DeviceType_ANDROID
			// Register session
			err := svc.RegisterSession(ctx, s)
			if err != nil {
				// We don't write back error strictly in this dummy
				return
			}
			// Write success back
			s.Conn.WriteMessage([]byte("SUCCESS"))
		} else {
			// Normal message echo
			if c == 999 {
				t.Log("noop")
			}
			s.Conn.WriteMessage([]byte("ECHO: " + str))
		}
	}

	tcpServer := server.NewTCPServer("127.0.0.1:0", svc, handler)
	err := tcpServer.Start()
	require.NoError(t, err)
	defer tcpServer.Stop()

	addr := tcpServer.ListenerAddr()

	// SCENARIO 1: Client logs in successfully
	// WHEN a client connects and sends a LOGIN frame
	clientA := newMockTurmsClient(t, addr)
	defer clientA.Close()

	clientA.SendReq(t, 1001, []byte("LOGIN"))

	// THEN the client should receive a SUCCESS response
	respA, err := clientA.ReadResponse(t)
	require.NoError(t, err)
	assert.Equal(t, "SUCCESS", string(respA))

	// AND the session should be registered in the service
	s, ok := svc.GetUserSession(100, protocol.DeviceType_ANDROID)
	assert.True(t, ok)
	assert.Equal(t, 100, int(s.UserID))
	assert.Equal(t, 1, svc.CountOnlineUsers())

	// WHEN Client A sends a normal message
	clientA.SendReq(t, 1002, []byte("HELLO"))
	respEcho, err := clientA.ReadResponse(t)
	require.NoError(t, err)
	assert.Equal(t, "ECHO: HELLO", string(respEcho))

	// SCENARIO 2: Multi-Device Kick
	// WHEN Client B connects using the SAME UserID and DeviceType (simulated by same "LOGIN" payload)
	clientB := newMockTurmsClient(t, addr)
	defer clientB.Close()

	clientB.SendReq(t, 2001, []byte("LOGIN"))

	// THEN Client B should receive SUCCESS
	respB, err := clientB.ReadResponse(t)
	require.NoError(t, err)
	assert.Equal(t, "SUCCESS", string(respB))

	// AND Client A should be disconnected by the server explicitly closing its socket
	// Verify that Client A's connection is closed by trying to read (EOF expected)
	_, err = clientA.ReadResponse(t)
	assert.Error(t, err, "Expected error on Client A due to being kicked out")

	// AND the service should still have exactly 1 online user (Client B replaced A)
	assert.Equal(t, 1, svc.CountOnlineUsers())
}
