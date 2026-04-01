package repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"im.turms/server/internal/domain/group/po"
	turmsmongo "im.turms/server/internal/storage/mongo"
)

func TestGroupMemberRepository_BDD(t *testing.T) {
	ctx := context.Background()

	mongodbContainer, err := mongodb.Run(ctx, "mongo:7.0")
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

	repo := NewGroupMemberRepository(client)

	// 1. Setup Data
	groupID := int64(1)
	memberID := int64(100)
	nonMemberID := int64(200)

	err = repo.Insert(ctx, &po.GroupMember{
		ID:   po.GroupMemberKey{GroupID: groupID, UserID: memberID},
		Role: 1, // member
	})
	assert.NoError(t, err)

	// 2. Test IsGroupMember
	ok, err := repo.IsGroupMember(ctx, groupID, memberID)
	assert.NoError(t, err)
	assert.True(t, ok)

	ok, err = repo.IsGroupMember(ctx, groupID, nonMemberID)
	assert.NoError(t, err)
	assert.False(t, ok)
}
