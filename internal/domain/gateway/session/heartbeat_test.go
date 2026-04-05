package session

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"im.turms/server/pkg/protocol"
)

func TestHeartbeatManager_KickTimeout(t *testing.T) {
	svc := NewSessionService(nil, nil, nil, nil, "test-server-id", nil)

	conn := &MockConnection{}
	session := NewUserSession(1, nil, 123, protocol.DeviceType_ANDROID, nil, nil)
	session.IP = net.ParseIP("127.0.0.1")
	session.Conn = conn

	err := svc.RegisterSession(context.Background(), session)
	assert.NoError(t, err)

	// Set last heartbeat to 10 seconds ago
	session.lastHeartbeat = time.Now().Add(-10 * time.Second).UnixMilli()

	// Timeout is 5 seconds
	manager := NewHeartbeatManager(svc, 100*time.Millisecond, 5*time.Second)
	manager.Start()

	// Wait for cleanup loop
	time.Sleep(200 * time.Millisecond)
	manager.Stop()

	assert.Equal(t, 0, svc.CountOnlineUsers())
	assert.True(t, conn.Closed)
}
