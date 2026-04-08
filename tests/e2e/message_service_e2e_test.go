package e2e

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"im.turms/server/internal/domain/common/infra/idgen"
	"im.turms/server/internal/domain/message/repository"
	"im.turms/server/internal/domain/message/service"
	"im.turms/server/internal/infra/plugin"
	"im.turms/server/internal/infra/property"
	turmsredis "im.turms/server/internal/storage/redis"
	"im.turms/server/internal/testingutil"
)

func TestMessageCore_E2E(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	db, cleanupMongo := testingutil.SetupMongo(t, "turms_message_core_e2e_test")
	defer cleanupMongo()

	rdb, cleanupRedis := testingutil.SetupRedis(t)
	defer cleanupRedis()

	idGen, err := idgen.NewSnowflakeIdGenerator(1, 1)
	require.NoError(t, err)
	seqGen := turmsredis.NewSequenceGenerator(rdb)

	msgRepo := repository.NewMessageRepository(db)

	// Create MessageService (passing nil for user/group relation checks to isolate message testing)
	propsMgr := property.NewTurmsPropertiesManager()
	plugMgr := plugin.NewPluginManager()
	msgService := service.NewMessageService(idGen, seqGen, msgRepo, nil, nil, nil, nil, nil, propsMgr, plugMgr)
	defer msgService.Close()

	t.Run("Message Creation, Recall, and Modification", func(t *testing.T) {
		senderID := int64(101)
		targetID := int64(201)
		text := "Hello World"

		// 1. Create message
		result, err := msgService.AuthAndSaveMessage(ctx, false, senderID, targetID, false, text, nil, nil, nil, nil, "", nil)
		require.NoError(t, err)
		require.NotNil(t, result)
		msg := result.Message
		assert.NotNil(t, msg)
		assert.Equal(t, text, msg.Text)
		assert.Nil(t, msg.RecallDate)
		assert.Nil(t, msg.ModificationDate)

		msgID := msg.ID

		// 2. Modify message
		newText := "Hello Turms!"
		err = msgService.AuthAndUpdateMessageText(ctx, senderID, msgID, newText)
		require.NoError(t, err)

		// Verification
		fetchedMsg, err := msgRepo.FindByID(ctx, msgID)
		require.NoError(t, err)
		assert.Equal(t, newText, fetchedMsg.Text)
		assert.NotNil(t, fetchedMsg.ModificationDate)
		assert.Nil(t, fetchedMsg.RecallDate)

		// 3. Recall message
		err = msgService.AuthAndRecallMessage(ctx, senderID, msgID)
		require.NoError(t, err)

		// Verification
		fetchedMsg2, err := msgRepo.FindByID(ctx, msgID)
		require.NoError(t, err)
		assert.NotNil(t, fetchedMsg2.RecallDate)
		// Modification date remains
		assert.NotNil(t, fetchedMsg2.ModificationDate)
	})

	t.Run("Unauthorized Actions Should Fail", func(t *testing.T) {
		senderID := int64(101)
		targetID := int64(201)
		text := "A private message"

		result, err := msgService.AuthAndSaveMessage(ctx, false, senderID, targetID, false, text, nil, nil, nil, nil, "", nil)
		require.NoError(t, err)
		require.NotNil(t, result)

		wrongSenderID := int64(999)

		// Attempt to recall by wrong user
		err = msgService.AuthAndRecallMessage(ctx, wrongSenderID, result.Message.ID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unauthorized")

		// Attempt to modify by wrong user
		err = msgService.AuthAndUpdateMessageText(ctx, wrongSenderID, result.Message.ID, "hacked text")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unauthorized")
	})
}
