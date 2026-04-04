package session

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"im.turms/server/internal/domain/common/constant"
	"im.turms/server/internal/domain/gateway/session/bo"
	"im.turms/server/pkg/protocol"
)

// MockConnection for testing
type MockConnection struct {
	Closed bool
}

func (m *MockConnection) GetAddress() net.Addr {
	return &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 1234}
}

func (m *MockConnection) Send(payload []byte) error {
	return nil
}

func (m *MockConnection) SendWithContext(ctx context.Context, payload []byte) error {
	return nil
}

func (m *MockConnection) CloseWithReason(reason bo.CloseReason) bool {
	m.Closed = true
	return true
}

func (m *MockConnection) Close() error {
	m.Closed = true
	return nil
}

func (m *MockConnection) IsConnected() bool {
	return !m.Closed
}

func (m *MockConnection) IsSwitchingToUdp() bool {
	return false
}

func (m *MockConnection) IsConnectionRecovering() bool {
	return false
}

func (m *MockConnection) SwitchToUdp() {}

func (m *MockConnection) TryNotifyClientToRecover() {}

func (m *MockConnection) IsActive() bool {
	return !m.Closed
}

func TestSessionService_RegisterAndUnregister(t *testing.T) {
	svc := NewSessionService(nil, nil, nil, nil, "test-server-id", nil)

	conn := &MockConnection{}
	session := &UserSession{
		UserID:     123,
		DeviceType: protocol.DeviceType_ANDROID,
		Conn:       conn,
		CloseChan:  make(chan struct{}),
	}

	err := svc.RegisterSession(context.Background(), session)
	assert.NoError(t, err)

	assert.Equal(t, 1, svc.CountOnlineUsers())

	s, ok := svc.GetUserSession(123, protocol.DeviceType_ANDROID)
	assert.True(t, ok)
	assert.Equal(t, session, s)

	// Test Unregister
	svc.UnregisterSession(context.Background(), 123, protocol.DeviceType_ANDROID, conn, bo.NewCloseReason(constant.SessionCloseStatus_DISCONNECTED_BY_CLIENT))
	assert.Equal(t, 0, svc.CountOnlineUsers())
	assert.True(t, conn.Closed)
}

func TestSessionService_ConflictKick(t *testing.T) {
	svc := NewSessionService(nil, nil, nil, nil, "test-server-id", nil)

	conn1 := &MockConnection{}
	session1 := NewUserSession(1, nil, 456, protocol.DeviceType_IOS, nil, nil)
	session1.Conn = conn1

	err := svc.RegisterSession(context.Background(), session1)
	assert.NoError(t, err)

	// Second device, same type, difference connection (e.g. reconnection)
	conn2 := &MockConnection{}
	session2 := NewUserSession(1, nil, 456, protocol.DeviceType_IOS, nil, nil)
	session2.Conn = conn2

	err = svc.RegisterSession(context.Background(), session2)
	assert.NoError(t, err)

	// Should have kicked the first one
	assert.True(t, conn1.Closed)
	assert.False(t, conn2.Closed)

	s, ok := svc.GetUserSession(456, protocol.DeviceType_IOS)
	assert.True(t, ok)
	assert.Equal(t, session2, s) // should point to session2
}

func TestSessionService_DifferentDevices(t *testing.T) {
	svc := NewSessionService(nil, nil, nil, nil, "test-server-id", nil)
	svc.ConflictStrategy = KickExisting // Which only applies to the EXACT same device type currently

	conn1 := &MockConnection{}
	session1 := &UserSession{
		UserID:     456,
		DeviceType: protocol.DeviceType_IOS,
		Conn:       conn1,
		CloseChan:  make(chan struct{}),
	}

	conn2 := &MockConnection{}
	session2 := &UserSession{
		UserID:     456,
		DeviceType: protocol.DeviceType_DESKTOP, // Different
		Conn:       conn2,
		CloseChan:  make(chan struct{}),
	}

	svc.RegisterSession(context.Background(), session1)
	svc.RegisterSession(context.Background(), session2)

	assert.False(t, conn1.Closed)
	assert.False(t, conn2.Closed)

	sessions := svc.GetAllUserSessions(456)
	assert.Len(t, sessions, 2)
}
