package repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"

	"im.turms/server/internal/domain/message/po"
	turmsmongo "im.turms/server/internal/storage/mongo"
)

func TestMessageRepository_InsertAndFind(t *testing.T) {
	ctx := context.Background()

	// Spin up MongoDB container
	mongodbContainer, err := mongodb.Run(ctx, "mongo:7.0")
	require.NoError(t, err)
	defer func() {
		if err := mongodbContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}()

	uri, err := mongodbContainer.ConnectionString(ctx)
	require.NoError(t, err)

	cfg := turmsmongo.Config{
		URI:            uri,
		Database:       "turms_test",
		ConnectTimeout: 10 * time.Second,
	}

	client, err := turmsmongo.NewClient(ctx, cfg)
	require.NoError(t, err)
	defer client.Close(ctx)

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
	err = repo.InsertMessage(ctx, msg)
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
