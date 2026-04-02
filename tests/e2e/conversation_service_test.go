package e2e

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"im.turms/server/internal/domain/conversation/repository"
	"im.turms/server/internal/domain/conversation/service"
	"im.turms/server/internal/testingutil"
)

func TestConversationCore_E2E(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	db, cleanup := testingutil.SetupMongo(t, "turms_conversation_core_e2e_test")
	defer cleanup()

	privateConvRepo := repository.NewPrivateConversationRepository(db)
	groupConvRepo := repository.NewGroupConversationRepository(db)

	convService := service.NewConversationService(privateConvRepo, groupConvRepo)

	t.Run("PrivateConversation Lifecycle", func(t *testing.T) {
		ownerID := int64(101)
		targetID := int64(201)
		readDate := time.Now().UTC().Truncate(time.Millisecond)

		err := convService.AuthAndUpdatePrivateConversationReadDate(ctx, ownerID, targetID, readDate)
		require.NoError(t, err)

		results, err := convService.QueryPrivateConversations(ctx, []int64{ownerID})
		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, ownerID, results[0].ID.OwnerID)
		assert.Equal(t, targetID, results[0].ID.TargetID)
		assert.Equal(t, readDate, results[0].ReadDate.UTC())

		// Update again to test UPSERT
		newReadDate := readDate.Add(5 * time.Minute)
		err = convService.AuthAndUpdatePrivateConversationReadDate(ctx, ownerID, targetID, newReadDate)
		require.NoError(t, err)

		results, err = convService.QueryPrivateConversations(ctx, []int64{ownerID})
		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, newReadDate, results[0].ReadDate.UTC())
	})

	t.Run("GroupConversation Lifecycle", func(t *testing.T) {
		groupID := int64(1001)
		memberID1 := int64(101)
		memberID2 := int64(102)
		readDate1 := time.Now().UTC().Truncate(time.Millisecond)
		readDate2 := readDate1.Add(10 * time.Minute)

		err := convService.AuthAndUpdateGroupConversationReadDate(ctx, memberID1, groupID, readDate1)
		require.NoError(t, err)

		err = convService.AuthAndUpdateGroupConversationReadDate(ctx, memberID2, groupID, readDate2)
		require.NoError(t, err)

		results, err := convService.QueryGroupConversations(ctx, []int64{groupID})
		require.NoError(t, err)
		assert.Len(t, results, 1)

		// Map keys are strings in MongoDB when using generic map[string]interface{} BSON tags
		// We verify the map deserializes properly into map[string]time.Time
		groupConv := results[0]
		assert.Equal(t, groupID, groupConv.ID)
		require.NotNil(t, groupConv.MemberIDToReadDate)

		assert.Equal(t, readDate1, groupConv.MemberIDToReadDate["101"].UTC())
		assert.Equal(t, readDate2, groupConv.MemberIDToReadDate["102"].UTC())

		// Test Upsert for an existing member
		newReadDate1 := readDate1.Add(5 * time.Minute)
		err = convService.AuthAndUpdateGroupConversationReadDate(ctx, memberID1, groupID, newReadDate1)
		require.NoError(t, err)

		results, err = convService.QueryGroupConversations(ctx, []int64{groupID})
		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, newReadDate1, results[0].MemberIDToReadDate["101"].UTC())
	})
}
