package e2e_test

import (
	"bufio"
	"context"
	"encoding/binary"
	"io"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"im.turms/server/internal/domain/common/infra/cluster/rpc/codec"
	"im.turms/server/internal/domain/common/infra/idgen"
	"im.turms/server/internal/domain/gateway/access/client/common"
	"im.turms/server/internal/domain/gateway/access/router"
	gatewayServer "im.turms/server/internal/domain/gateway/access/server"
	"im.turms/server/internal/domain/gateway/session"

	grouppo "im.turms/server/internal/domain/group/po"
	messagecontroller "im.turms/server/internal/domain/message/controller"
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
	br   *bufio.Reader
	mu   sync.Mutex
}

func newMockTurmsE2EClient(t *testing.T, addr string) *mockTurmsE2EClient {
	conn, err := net.Dial("tcp", addr)
	require.NoError(t, err)
	return &mockTurmsE2EClient{conn: conn, br: bufio.NewReader(conn)}
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

	payloadLen, err := binary.ReadUvarint(c.br)
	if err != nil {
		return nil, err
	}

	payload := make([]byte, payloadLen)
	_, err = io.ReadFull(c.br, payload)
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

	// 4. Gateway Dispatcher and Controllers
	msgController := messagecontroller.NewMessageController(msgSvc)
	
	r := router.NewRouter(sessionSvc)
	r.SetServiceAvailability(common.StatusRunning)
	r.RegisterController(&protocol.TurmsRequest_CreateMessageRequest{}, msgController.HandleCreateMessageRequest)

	tcpServer := gatewayServer.NewTCPServer("127.0.0.1:0", sessionSvc, r.HandleMessage)
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

	id := privMsgResp.GetData().GetLong()
	assert.Greater(t, id, int64(0), "Message ID generated should be valid")

	// Verify Message is saved to MongoDB
	savedMsg, err := msgRepo.FindByID(ctx, id)
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

	grpId := grpMsgResp.GetData().GetLong()
	assert.Greater(t, grpId, int64(0))

	savedGrpMsg, err := msgRepo.FindByID(ctx, grpId)
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
	assert.Equal(t, int32(1100), unauthResp.GetCode())
	assert.NotNil(t, unauthResp.Reason)
	assert.Contains(t, *unauthResp.Reason, "not a member of the target group")
	// Step 5. Rate Limiting Verification (High Frequency Traffic)
	// The default throttler is 100 req/s with a burst of 100.
	// We'll rapid-fire 110 requests and expect the last ones to be dropped with code 1400.
	var wg sync.WaitGroup
	var rateLimitedCount int32
	var mu sync.Mutex

	for i := 0; i < 110; i++ {
		wg.Add(1)
		go func(reqID int32) {
			defer wg.Done()
			
			// We need a separate connection per request otherwise we'll interleave writes
			// However, mockTurmsE2EClient is just writing, we can use the same client 
			// if we just want to hit the router, but its ReadTurmsNotification might get mixed up.
			// Instead of reading all 110 responses, let's just create a quick new client for the extra calls,
			// or just blast the current client.
			// Best approach: Use the same client (mockTurmsE2EClient is thread-safe on send),
			// and read 110 responses.
			
			client.SendTurmsRequest(t, reqID, &protocol.TurmsRequest{
				RequestId: proto.Int64(int64(reqID)),
				Kind: &protocol.TurmsRequest_CreateMessageRequest{
					CreateMessageRequest: &protocol.CreateMessageRequest{
						RecipientId:  proto.Int64(200),
						Text:         proto.String("spam"),
					},
				},
			})
		}(int32(100 + i))
	}
	wg.Wait()

	// Now read 110 responses
	for i := 0; i < 110; i++ {
		resp, err := client.ReadTurmsNotification(t)
		require.NoError(t, err)
		if resp.GetCode() == 1400 {
			mu.Lock()
			rateLimitedCount++
			mu.Unlock()
		}
	}

	assert.GreaterOrEqual(t, rateLimitedCount, int32(1), "Expected at least 1 request to be rate limited (code 1400)")
}
