package session

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"im.turms/server/pkg/protocol"
)

// MockConnection for testing
type MockConnection struct {
	Closed bool
}

func (m *MockConnection) WriteMessage(payload []byte) error {
	return nil
}

func (m *MockConnection) Close() error {
	m.Closed = true
	return nil
}

func (m *MockConnection) RemoteAddr() net.Addr {
	return &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 1234}
}

func TestSessionService_RegisterAndUnregister(t *testing.T) {
	svc := NewSessionService()

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
	svc.UnregisterSession(123, protocol.DeviceType_ANDROID, conn)
	assert.Equal(t, 0, svc.CountOnlineUsers())
	assert.True(t, conn.Closed)
}

func TestSessionService_ConflictKick(t *testing.T) {
	svc := NewSessionService()

	conn1 := &MockConnection{}
	session1 := &UserSession{
		UserID:     456,
		DeviceType: protocol.DeviceType_IOS,
		Conn:       conn1,
		CloseChan:  make(chan struct{}),
	}
	
	err := svc.RegisterSession(context.Background(), session1)
	assert.NoError(t, err)

	// Second device, same type, difference connection (e.g. reconnection)
	conn2 := &MockConnection{}
	session2 := &UserSession{
		UserID:     456,
		DeviceType: protocol.DeviceType_IOS,
		Conn:       conn2,
		CloseChan:  make(chan struct{}),
	}
	
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
	svc := NewSessionService()
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
