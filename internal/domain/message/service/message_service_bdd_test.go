package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"im.turms/server/internal/testingutil"

	"im.turms/server/internal/domain/common/infra/idgen"
	"im.turms/server/internal/infra/plugin"
	"im.turms/server/internal/infra/property"
	grouppo "im.turms/server/internal/domain/group/po"
	"im.turms/server/internal/domain/message/po"
	"im.turms/server/internal/domain/message/repository"
	userpo "im.turms/server/internal/domain/user/po"
	turmsmongo "im.turms/server/internal/storage/mongo"
	turmsredis "im.turms/server/internal/storage/redis"
)

type mockUserRelService struct {
	allowed bool
}

func (m *mockUserRelService) HasRelationshipAndNotBlocked(ctx context.Context, ownerID int64, relatedUserID int64) (bool, error) {
	return m.allowed, nil
}

func (m *mockUserRelService) AddFriend(ctx context.Context, ownerID int64, friendID int64) (*userpo.UserRelationship, error) {
	return nil, nil
}

func (m *mockUserRelService) BlockUser(ctx context.Context, ownerID int64, blockedID int64) error {
	return nil
}

type mockGroupMemService struct {
	allowed bool
}

func (m *mockGroupMemService) IsGroupMember(ctx context.Context, groupID int64, userID int64) (bool, error) {
	return m.allowed, nil
}

func (m *mockGroupMemService) FindGroupMemberIDs(ctx context.Context, groupID int64) ([]int64, error) {
	return []int64{1, 2, 3}, nil
}

func (m *mockGroupMemService) AddGroupMember(ctx context.Context, groupID int64, userID int64, role int32, name *string, muteEndDate *time.Time) (*grouppo.GroupMember, error) {
	return nil, nil
}

type mockDelivery struct {
	delivered map[int64]*po.Message
}

func (m *mockDelivery) Deliver(ctx context.Context, targetID int64, msg *po.Message) error {
	m.delivered[targetID] = msg
	return nil
}

func setupTestInfra(t *testing.T) (*turmsredis.Client, *turmsmongo.Client, func()) {
	redisClient, rCleanup := testingutil.SetupRedis(t)
	mongoClient, mCleanup := testingutil.SetupMongo(t, "turms_bdd_test")

	cleanup := func() {
		rCleanup()
		mCleanup()
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
		name        string
		senderID    int64
		targetID    int64
		isGroup     bool
		text        string
		mockUserRel bool
		mockGroup   bool
		wantErr     error
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
			name:      "3. Group Message Success",
			senderID:  1,
			targetID:  1001,
			isGroup:   true,
			text:      "Hello, group!",
			mockGroup: true,
			wantErr:   nil,
		},
		{
			name:      "4. Group Message Not Member",
			senderID:  1,
			targetID:  1002,
			isGroup:   true,
			text:      "Hello?",
			mockGroup: false,
			wantErr:   ErrNotGroupMember,
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
			propsMgr := property.NewTurmsPropertiesManager()
			plugMgr := plugin.NewPluginManager()

			svc := NewMessageService(idGen, seqGen, msgRepo, userSvc, groupSvc, nil, nil, delivery, propsMgr, plugMgr)
			defer svc.Close()

			result, err := svc.AuthAndSaveAndSendMessage(ctx, tt.isGroup, tt.senderID, tt.targetID, false, tt.text, nil, nil, nil, nil, "", nil)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				msg := result.Message
				assert.NotNil(t, msg)
				assert.Equal(t, tt.senderID, msg.SenderID)
				assert.Equal(t, tt.targetID, msg.TargetID)
				assert.Equal(t, tt.text, msg.Text)
				assert.Condition(t, func() bool { return msg.SequenceID != nil && *msg.SequenceID > 0 })

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
