package repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"im.turms/server/internal/domain/user/po"
	turmsmongo "im.turms/server/internal/storage/mongo"
)

func TestUserRelationshipRepository_BDD(t *testing.T) {
	ctx := context.Background()

	mongodbContainer, err := mongodb.Run(ctx, "mongo:7.0")

	// Fallback/Graceful skip if Docker is not available in user terminal
	if err != nil {
		t.Skipf("Skipping BDD test, testcontainers failed to start: %v", err)
	}
	defer func() {
		_ = mongodbContainer.Terminate(ctx)
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

	repo := NewUserRelationshipRepository(client)

	// 1. Setup Data
	ownerID := int64(100)
	friendID := int64(200)
	blockedID := int64(300)
	now := time.Now()

	err = repo.Insert(ctx, &po.UserRelationship{
		ID:        po.UserRelationshipKey{OwnerID: ownerID, RelatedUserID: friendID},
		BlockDate: nil,
	})
	assert.NoError(t, err)

	err = repo.Insert(ctx, &po.UserRelationship{
		ID:        po.UserRelationshipKey{OwnerID: ownerID, RelatedUserID: blockedID},
		BlockDate: &now,
	})
	assert.NoError(t, err)

	// 2. Test HasRelationshipAndNotBlocked
	// Friend -> true
	ok, err := repo.HasRelationshipAndNotBlocked(ctx, ownerID, friendID)
	assert.NoError(t, err)
	assert.True(t, ok)

	// Blocked -> false
	ok, err = repo.HasRelationshipAndNotBlocked(ctx, ownerID, blockedID)
	assert.NoError(t, err)
	assert.False(t, ok)

	// Stranger -> false
	ok, err = repo.HasRelationshipAndNotBlocked(ctx, ownerID, 999)
	assert.NoError(t, err)
	assert.False(t, ok)
}
