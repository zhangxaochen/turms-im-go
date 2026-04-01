package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	testredis "github.com/testcontainers/testcontainers-go/modules/redis"

	"im.turms/server/internal/domain/common/infra/idgen"
	"im.turms/server/internal/domain/message/po"
	"im.turms/server/internal/domain/message/repository"
	turmsmongo "im.turms/server/internal/storage/mongo"
	turmsredis "im.turms/server/internal/storage/redis"
)

type mockUserRelService struct {
	allowed bool
}

func (m *mockUserRelService) HasRelationshipAndNotBlocked(ctx context.Context, ownerID int64, relatedUserID int64) (bool, error) {
	return m.allowed, nil
}

type mockGroupMemService struct {
	allowed bool
}

func (m *mockGroupMemService) IsGroupMember(ctx context.Context, groupID int64, userID int64) (bool, error) {
	return m.allowed, nil
}

type mockDelivery struct {
	delivered map[int64]*po.Message
}

func (m *mockDelivery) Deliver(ctx context.Context, targetID int64, msg *po.Message) error {
	m.delivered[targetID] = msg
	return nil
}

func setupTestInfra(t *testing.T) (*turmsredis.Client, *turmsmongo.Client, func()) {
	ctx := context.Background()

	// Redis
	redisContainer, err := testredis.Run(ctx, "redis:7.0")
	require.NoError(t, err)

	redisURI, err := redisContainer.ConnectionString(ctx)
	require.NoError(t, err)

	redisClient, err := turmsredis.NewClient(ctx, turmsredis.Config{URI: redisURI})
	require.NoError(t, err)

	// MongoDB
	mongoContainer, err := mongodb.Run(ctx, "mongo:7.0")
	require.NoError(t, err)

	mongoURI, err := mongoContainer.ConnectionString(ctx)
	require.NoError(t, err)

	mongoClient, err := turmsmongo.NewClient(ctx, turmsmongo.Config{
		URI:            mongoURI,
		Database:       "turms_bdd_test",
		ConnectTimeout: 10 * time.Second,
	})
	require.NoError(t, err)

	cleanup := func() {
		redisClient.Close()
		mongoClient.Close(ctx)
		redisContainer.Terminate(ctx)
		mongoContainer.Terminate(ctx)
	}

	return redisClient, mongoClient, cleanup
}

func TestMessageService_AuthAndSaveAndSendMessage_BDD(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping BDD Testcontainers tests in short mode")
	}

	redisClient, mongoClient, cleanup := setupTestInfra(t)
	defer cleanup()

	idGen, err := idgen.NewSnowflakeIdGenerator(0, 0)
	require.NoError(t, err)

	seqGen := turmsredis.NewSequenceGenerator(redisClient)
	msgRepo := repository.NewMessageRepository(mongoClient)

	tests := []struct {
		name           string
		senderID       int64
		targetID       int64
		isGroup        bool
		text           string
		mockUserRel    bool
		mockGroup      bool
		wantErr        error
	}{
		{
			name:        "1. Private Message Success",
			senderID:    1,
			targetID:    2,
			isGroup:     false,
			text:        "Hello, my friend!",
			mockUserRel: true,
			wantErr:     nil,
		},
		{
			name:        "2. Private Message Blocked/Not Friend",
			senderID:    1,
			targetID:    3,
			isGroup:     false,
			text:        "Hello?",
			mockUserRel: false,
			wantErr:     ErrNotFriend,
		},
		{
			name:        "3. Group Message Success",
			senderID:    1,
			targetID:    1001,
			isGroup:     true,
			text:        "Hello, group!",
			mockGroup:   true,
			wantErr:     nil,
		},
		{
			name:        "4. Group Message Not Member",
			senderID:    1,
			targetID:    1002,
			isGroup:     true,
			text:        "Hello?",
			mockGroup:   false,
			wantErr:     ErrNotGroupMember,
		},
		{
			name:        "5. Invalid Target ID",
			senderID:    1,
			targetID:    0,
			isGroup:     false,
			text:        "Invalid Target",
			mockUserRel: true,
			wantErr:     ErrInvalidTargetID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			userSvc := &mockUserRelService{allowed: tt.mockUserRel}
			groupSvc := &mockGroupMemService{allowed: tt.mockGroup}
			delivery := &mockDelivery{delivered: make(map[int64]*po.Message)}

			svc := NewMessageService(idGen, seqGen, msgRepo, userSvc, groupSvc, delivery)
			defer svc.Close()

			msg, err := svc.AuthAndSaveAndSendMessage(ctx, tt.isGroup, tt.senderID, tt.targetID, tt.text)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, msg)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, msg)
				assert.Equal(t, tt.senderID, msg.SenderID)
				assert.Equal(t, tt.targetID, msg.TargetID)
				assert.Equal(t, tt.text, msg.Text)
				assert.NotNil(t, msg.SequenceID)
				assert.Condition(t, func() bool { return *msg.SequenceID > 0 })

				// Validate delivered
				assert.Contains(t, delivery.delivered, tt.targetID)
				
				// Validate persisted in MongoDB
				foundMsgs, err := msgRepo.FindMessagesByTarget(ctx, tt.targetID)
				require.NoError(t, err)
				require.NotEmpty(t, foundMsgs)
				
				// Ensure the message we just saved is found
				var found bool
				for _, m := range foundMsgs {
					if m.ID == msg.ID {
						found = true
						assert.Equal(t, tt.text, m.Text)
						break
					}
				}
				assert.True(t, found, "saved message not found in db")
			}
		})
	}
}
