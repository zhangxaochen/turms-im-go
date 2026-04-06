package router_test

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"im.turms/server/internal/domain/gateway/access/client/common"
	"im.turms/server/internal/domain/gateway/access/router"
	"im.turms/server/internal/domain/gateway/session"
	"im.turms/server/internal/domain/gateway/session/bo"
	"im.turms/server/pkg/protocol"
)

type mockConnection struct {
	sent           [][]byte
	lastWrittenMsg []byte
	closed         bool
	closedReason   bo.CloseReason
}

func (m *mockConnection) GetAddress() net.Addr { return &net.IPAddr{} }
func (m *mockConnection) Send(data []byte) error {
	m.sent = append(m.sent, data)
	m.lastWrittenMsg = data
	return nil
}
func (m *mockConnection) SendWithContext(ctx context.Context, data []byte) error {
	return m.Send(data)
}
func (m *mockConnection) CloseWithReason(reason bo.CloseReason) bool {
	m.closed = true
	m.closedReason = reason
	return true
}
func (m *mockConnection) Close() error {
	m.closed = true
	return nil
}
func (m *mockConnection) IsConnected() bool            { return !m.closed }
func (m *mockConnection) IsActive() bool               { return !m.closed }
func (m *mockConnection) IsSwitchingToUdp() bool       { return false }
func (m *mockConnection) IsConnectionRecovering() bool { return false }
func (m *mockConnection) SwitchToUdp()                 {}
func (m *mockConnection) TryNotifyClientToRecover()    {}

func TestRouter_HandleMessage(t *testing.T) {
	ctx := context.Background()

	setupRouter := func() (*session.SessionService, *router.Router) {
		sessionSvc := session.NewSessionService(nil, nil, nil, nil, "test-server-id", nil, nil)
		r := router.NewRouter(sessionSvc)
		r.SetServiceAvailability(common.StatusRunning)

		// Register a simple echo controller
		r.RegisterController(&protocol.TurmsRequest_CreateMessageRequest{}, func(ctx context.Context, s *session.UserSession, req *protocol.TurmsRequest) (*protocol.TurmsNotification, error) {
			return &protocol.TurmsNotification{
				RequestId: req.RequestId,
				Code:      proto.Int32(1000),
				Data: &protocol.TurmsNotification_Data{
					Kind: &protocol.TurmsNotification_Data_Long{
						Long: 123,
					},
				},
			}, nil
		})
		return sessionSvc, r
	}

	t.Run("Requires Session Auth", func(t *testing.T) {
		_, r := setupRouter()
		mockConn := &mockConnection{}
		uSession := &session.UserSession{Conn: mockConn, CloseChan: make(chan struct{})}

		req := &protocol.TurmsRequest{
			RequestId: proto.Int64(1),
			Kind: &protocol.TurmsRequest_CreateMessageRequest{
				CreateMessageRequest: &protocol.CreateMessageRequest{
					Text: proto.String("hello"),
				},
			},
		}
		payload, _ := proto.Marshal(req)

		r.HandleMessage(ctx, uSession, payload)

		// Verification: should get unauthed (1200)
		resp := &protocol.TurmsNotification{}
		err := proto.Unmarshal(mockConn.lastWrittenMsg, resp)
		require.NoError(t, err)
		assert.Equal(t, int32(1200), resp.GetCode())
	})

	t.Run("Create Session", func(t *testing.T) {
		_, r := setupRouter()
		mockConn := &mockConnection{}
		uSession := &session.UserSession{Conn: mockConn, CloseChan: make(chan struct{})}

		req := &protocol.TurmsRequest{
			RequestId: proto.Int64(2),
			Kind: &protocol.TurmsRequest_CreateSessionRequest{
				CreateSessionRequest: &protocol.CreateSessionRequest{
					UserId:     101,
					DeviceType: protocol.DeviceType_DESKTOP,
				},
			},
		}
		payload, _ := proto.Marshal(req)

		r.HandleMessage(ctx, uSession, payload)

		// Verification
		resp := &protocol.TurmsNotification{}
		err := proto.Unmarshal(mockConn.lastWrittenMsg, resp)
		require.NoError(t, err)
		assert.Equal(t, int32(1000), resp.GetCode())
		assert.Equal(t, int64(101), uSession.UserID)
	})

	t.Run("Valid Routing After Auth", func(t *testing.T) {
		sessionSvc, r := setupRouter()
		mockConn := &mockConnection{}
		uSession := &session.UserSession{Conn: mockConn, UserID: 101, DeviceType: protocol.DeviceType_DESKTOP, CloseChan: make(chan struct{})}
		// Register it in session svc to make sure heartbeat updates work
		_ = sessionSvc.RegisterSession(ctx, uSession)

		req := &protocol.TurmsRequest{
			RequestId: proto.Int64(3),
			Kind: &protocol.TurmsRequest_CreateMessageRequest{
				CreateMessageRequest: &protocol.CreateMessageRequest{
					Text: proto.String("hello mapped"),
				},
			},
		}
		payload, _ := proto.Marshal(req)

		r.HandleMessage(ctx, uSession, payload)

		resp := &protocol.TurmsNotification{}
		err := proto.Unmarshal(mockConn.lastWrittenMsg, resp)
		require.NoError(t, err)
		assert.Equal(t, int32(1000), resp.GetCode())
		assert.Equal(t, int64(123), resp.GetData().GetLong())
	})
}
