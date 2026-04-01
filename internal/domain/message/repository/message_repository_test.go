package repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"im.turms/server/internal/testingutil"

	"im.turms/server/internal/domain/message/po"
)

func TestMessageRepository_InsertAndFind(t *testing.T) {
	ctx := context.Background()

	client, cleanup := testingutil.SetupMongo(t, "turms_test")
	defer cleanup()

	repo := NewMessageRepository(client)

	// Create a dummy message
	now := time.Now().Truncate(time.Millisecond).UTC() // MongoDB truncates to ms
	isGroupMessage := false
	msg := &po.Message{
		ID:             1001,
		IsGroupMessage: &isGroupMessage,
		DeliveryDate:   now,
		Text:           "Hello Testcontainers",
		SenderID:       10,
		TargetID:       20,
		UserDefinedAttributes: map[string]any{
			"custom_key": "custom_value",
		},
	}

	// Test Insert
	err := repo.InsertMessage(ctx, msg)
	assert.NoError(t, err)

	// Test Find
	msgs, err := repo.FindMessagesByTarget(ctx, 20)
	assert.NoError(t, err)
	require.Len(t, msgs, 1)

	fetchedMsg := msgs[0]
	assert.Equal(t, int64(1001), fetchedMsg.ID)
	assert.Equal(t, "Hello Testcontainers", fetchedMsg.Text)
	assert.Equal(t, now, fetchedMsg.DeliveryDate.UTC())

	val, exists := fetchedMsg.UserDefinedAttributes["custom_key"]
	assert.True(t, exists)
	assert.Equal(t, "custom_value", val)
}
