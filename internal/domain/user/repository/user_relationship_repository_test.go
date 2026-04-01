package repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"im.turms/server/internal/domain/user/po"
	"im.turms/server/internal/testingutil"
)

func TestUserRelationshipRepository_BDD(t *testing.T) {
	ctx := context.Background()

	client, cleanup := testingutil.SetupMongo(t, "turms_test")
	defer cleanup()

	repo := NewUserRelationshipRepository(client)

	// 1. Setup Data
	ownerID := int64(100)
	friendID := int64(200)
	blockedID := int64(300)
	now := time.Now()

	err := repo.Insert(ctx, &po.UserRelationship{
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
