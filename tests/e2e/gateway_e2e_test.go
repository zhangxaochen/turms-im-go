package e2e_test

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
	"google.golang.org/protobuf/proto"

	"im.turms/server/internal/domain/common/infra/cluster/rpc/codec"
	"im.turms/server/internal/domain/common/infra/idgen"
	gatewayServer "im.turms/server/internal/domain/gateway/access/server"
	"im.turms/server/internal/domain/gateway/session"

	grouppo "im.turms/server/internal/domain/group/po"
	messagepo "im.turms/server/internal/domain/message/po"
	messagerepo "im.turms/server/internal/domain/message/repository"
	messageservice "im.turms/server/internal/domain/message/service"

	userpo "im.turms/server/internal/domain/user/po"

	turmsredis "im.turms/server/internal/storage/redis"

	"im.turms/server/internal/testingutil"
	"im.turms/server/pkg/protocol"
)

// mockTurmsE2EClient is a simple TCP client that mimics the exact behavior of a Turms Protobuf client SDK
type mockTurmsE2EClient struct {
	conn net.Conn
	mu   sync.Mutex
}

func newMockTurmsE2EClient(t *testing.T, addr string) *mockTurmsE2EClient {
	conn, err := net.Dial("tcp", addr)
	require.NoError(t, err)
	return &mockTurmsE2EClient{conn: conn}
}

func (c *mockTurmsE2EClient) SendTurmsRequest(t *testing.T, requestID int32, turmsReq *protocol.TurmsRequest) {
	payload, err := proto.Marshal(turmsReq)
	require.NoError(t, err)

	c.mu.Lock()
	defer c.mu.Unlock()

	frameSize := codec.HeaderSize + len(payload)
	varintBuf := make([]byte, binary.MaxVarintLen32)
	varintBytes := binary.PutUvarint(varintBuf, uint64(frameSize))

	_, err = c.conn.Write(varintBuf[:varintBytes])
	require.NoError(t, err)

	headerBuf := make([]byte, codec.HeaderSize)
	binary.BigEndian.PutUint16(headerBuf[0:2], uint16(1)) // CodecID 1 = client->server
	binary.BigEndian.PutUint32(headerBuf[2:6], uint32(requestID))
	_, err = c.conn.Write(headerBuf)
	require.NoError(t, err)

	if len(payload) > 0 {
		_, err = c.conn.Write(payload)
		require.NoError(t, err)
	}
}

func (c *mockTurmsE2EClient) ReadTurmsNotification(t *testing.T) (*protocol.TurmsNotification, error) {
	err := c.conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		return nil, err
	}

	br := bufio.NewReader(c.conn)
	payloadLen, err := binary.ReadUvarint(br)
	if err != nil {
		return nil, err
	}

	payload := make([]byte, payloadLen)
	_, err = br.Read(payload)
	if err != nil {
		return nil, err
	}

	notification := &protocol.TurmsNotification{}
	err = proto.Unmarshal(payload, notification)
	return notification, err
}

func (c *mockTurmsE2EClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

// Minimal stub for E2E
type mockUserRelService struct{}

func (m *mockUserRelService) HasRelationshipAndNotBlocked(ctx context.Context, ownerID int64, relatedUserID int64) (bool, error) {
	return true, nil // Always true for test
}

func (m *mockUserRelService) AddFriend(ctx context.Context, ownerID int64, friendID int64) (*userpo.UserRelationship, error) {
	return nil, nil
}

func (m *mockUserRelService) BlockUser(ctx context.Context, ownerID int64, blockedID int64) error {
	return nil
}

type mockGroupMemService struct{}

func (m *mockGroupMemService) IsGroupMember(ctx context.Context, groupID int64, userID int64) (bool, error) {
	return groupID == 1, nil // Member of group 1, not of group 99
}

func (m *mockGroupMemService) AddGroupMember(ctx context.Context, groupID int64, userID int64, role int32, name *string, muteEndDate *time.Time) (*grouppo.GroupMember, error) {
	return nil, nil
}

type localDelivery struct {
	svc *session.SessionService
}

func (d *localDelivery) Deliver(ctx context.Context, targetID int64, msg *messagepo.Message) error {
	// Not full broadcast in E2E, just simulating delivery success
	return nil
}

// BDD E2E Protocol Scenario
func TestGateway_E2E_TCP_Lifecycle(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E tests requiring Redis and Mongo in short mode")
	}
	ctx := context.Background()

	// 1. Infra Boot
	redisClient, rCleanup := testingutil.SetupRedis(t)
	defer rCleanup()

	mongoClient, mCleanup := testingutil.SetupMongo(t, "turms_e2e_test")
	defer mCleanup()

	// 2. Repositories
	msgRepo := messagerepo.NewMessageRepository(mongoClient)

	// 3. Services
	sessionSvc := session.NewSessionService()
	sessionSvc.ConflictStrategy = session.KickExisting

	idGen, _ := idgen.NewSnowflakeIdGenerator(1, 1)
	seqGen := turmsredis.NewSequenceGenerator(redisClient)

	msgSvc := messageservice.NewMessageService(
		idGen,
		seqGen,
		msgRepo,
		&mockUserRelService{},
		&mockGroupMemService{},
		&localDelivery{svc: sessionSvc},
	)

	// 4. Gateway Handler (Simulating Turms-Gateway Dispatcher connecting to Turms-Service)
	handler := func(hCtx context.Context, s *session.UserSession, payload []byte) {
		req := &protocol.TurmsRequest{}
		if err := proto.Unmarshal(payload, req); err != nil {
			s.Conn.WriteMessage([]byte("MALFORMED"))
			return
		}

		if s.UserID == 0 && req.GetCreateSessionRequest() != nil {
			loginReq := req.GetCreateSessionRequest()
			s.UserID = loginReq.UserId
			s.DeviceType = loginReq.DeviceType
			err := sessionSvc.RegisterSession(hCtx, s)

			// Mock response TurmsNotification
			resp := &protocol.TurmsNotification{
				RequestId: req.RequestId,
				Code:      proto.Int32(1000),
			}
			if err != nil {
				resp.Code = proto.Int32(1100)
			}
			respBytes, _ := proto.Marshal(resp)
			s.Conn.WriteMessage(respBytes)
			return
		}

		// Require Auth
		if s.UserID == 0 {
			resp := &protocol.TurmsNotification{
				RequestId: req.RequestId,
				Code:      proto.Int32(1200), // unauthed
			}
			respBytes, _ := proto.Marshal(resp)
			s.Conn.WriteMessage(respBytes)
			return
		}

		// Message Request
		if msgReq := req.GetCreateMessageRequest(); msgReq != nil {
			var targetID int64
			var isGroupMessage bool
			if msgReq.RecipientId != nil {
				targetID = *msgReq.RecipientId
				isGroupMessage = false
			} else if msgReq.GroupId != nil {
				targetID = *msgReq.GroupId
				isGroupMessage = true
			}

			msg, err := msgSvc.AuthAndSaveAndSendMessage(hCtx, isGroupMessage, s.UserID, targetID, *msgReq.Text)
			resp := &protocol.TurmsNotification{
				RequestId: req.RequestId,
			}
			if err != nil {
				resp.Code = proto.Int32(1300)
				resp.Reason = proto.String(err.Error())
			} else {
				resp.Code = proto.Int32(1000)
				resp.Data = &protocol.TurmsNotification_Data{
					Kind: &protocol.TurmsNotification_Data_LongsWithVersion{
						LongsWithVersion: &protocol.LongsWithVersion{
							Longs: []int64{msg.ID},
						},
					},
				}
			}
			respBytes, _ := proto.Marshal(resp)
			s.Conn.WriteMessage(respBytes)
		}
	}

	tcpServer := gatewayServer.NewTCPServer("127.0.0.1:0", sessionSvc, handler)
	err := tcpServer.Start()
	require.NoError(t, err)
	defer tcpServer.Stop()

	addr := tcpServer.ListenerAddr()

	// ---- CLIENT E2E SCENARIO ----
	client := newMockTurmsE2EClient(t, addr)
	defer client.Close()

	// Step 1. Login Authentication
	client.SendTurmsRequest(t, 1, &protocol.TurmsRequest{
		RequestId: proto.Int64(1),
		Kind: &protocol.TurmsRequest_CreateSessionRequest{
			CreateSessionRequest: &protocol.CreateSessionRequest{
				Version:    1,
				UserId:     100,
				Password:   proto.String("password"),
				DeviceType: protocol.DeviceType_ANDROID,
			},
		},
	})

	loginResp, err := client.ReadTurmsNotification(t)
	require.NoError(t, err)
	assert.Equal(t, int32(1000), loginResp.GetCode())

	// Step 2. Send Private Message to a related friend (200)
	client.SendTurmsRequest(t, 2, &protocol.TurmsRequest{
		RequestId: proto.Int64(2),
		Kind: &protocol.TurmsRequest_CreateMessageRequest{
			CreateMessageRequest: &protocol.CreateMessageRequest{
				RecipientId:  proto.Int64(200),
				Text:         proto.String("Hello friend!"),
				DeliveryDate: proto.Int64(time.Now().UnixMilli()),
			},
		},
	})

	privMsgResp, err := client.ReadTurmsNotification(t)
	require.NoError(t, err)
	assert.Equal(t, int32(1000), privMsgResp.GetCode())

	ids := privMsgResp.GetData().GetLongsWithVersion().GetLongs()
	require.NotNil(t, ids)
	assert.Len(t, ids, 1)
	assert.Greater(t, ids[0], int64(0), "Message ID generated should be valid")

	// Verify Message is saved to MongoDB
	savedMsg, err := msgRepo.FindByID(ctx, ids[0])
	require.NoError(t, err)
	assert.NotNil(t, savedMsg)
	assert.Equal(t, "Hello friend!", savedMsg.Text)
	assert.Equal(t, int64(100), savedMsg.SenderID)
	assert.Equal(t, int64(200), savedMsg.TargetID)

	// Step 3. Try to send message to group where I am member (1)
	client.SendTurmsRequest(t, 3, &protocol.TurmsRequest{
		RequestId: proto.Int64(3),
		Kind: &protocol.TurmsRequest_CreateMessageRequest{
			CreateMessageRequest: &protocol.CreateMessageRequest{
				GroupId:      proto.Int64(1),
				Text:         proto.String("Hello group!"),
				DeliveryDate: proto.Int64(time.Now().UnixMilli()),
			},
		},
	})

	grpMsgResp, err := client.ReadTurmsNotification(t)
	require.NoError(t, err)
	assert.Equal(t, int32(1000), grpMsgResp.GetCode())

	grpIds := grpMsgResp.GetData().GetLongsWithVersion().GetLongs()
	require.NotNil(t, grpIds)
	assert.Len(t, grpIds, 1)

	savedGrpMsg, err := msgRepo.FindByID(ctx, grpIds[0])
	require.NoError(t, err)
	assert.NotNil(t, savedGrpMsg)
	assert.Equal(t, "Hello group!", savedGrpMsg.Text)
	assert.Equal(t, int64(1), savedGrpMsg.TargetID)

	// Step 4. Error Case: Try to send message to group I'm NOT in (99)
	client.SendTurmsRequest(t, 4, &protocol.TurmsRequest{
		RequestId: proto.Int64(4),
		Kind: &protocol.TurmsRequest_CreateMessageRequest{
			CreateMessageRequest: &protocol.CreateMessageRequest{
				GroupId:      proto.Int64(99),
				Text:         proto.String("Hello unauthorized group!"),
				DeliveryDate: proto.Int64(time.Now().UnixMilli()),
			},
		},
	})

	unauthResp, err := client.ReadTurmsNotification(t)
	require.NoError(t, err)
	assert.Equal(t, int32(1300), unauthResp.GetCode())
	assert.NotNil(t, unauthResp.Reason)
	assert.Contains(t, *unauthResp.Reason, "not a member of the target group")
}
