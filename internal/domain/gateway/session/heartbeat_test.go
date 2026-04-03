package session

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"im.turms/server/pkg/protocol"
)

func TestHeartbeatManager_KickTimeout(t *testing.T) {
	svc := NewSessionService(nil, "test-server-id")

	conn := &MockConnection{}
	session := &UserSession{
		UserID:     123,
		DeviceType: protocol.DeviceType_ANDROID,
		Conn:       conn,
		CloseChan:  make(chan struct{}),
	}

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
