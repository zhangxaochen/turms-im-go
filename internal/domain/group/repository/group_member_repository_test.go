package repository

import (
	"context"
	"testing"
	"github.com/stretchr/testify/assert"
	"im.turms/server/internal/domain/group/po"
	"im.turms/server/internal/testingutil"
)

func TestGroupMemberRepository_BDD(t *testing.T) {
	ctx := context.Background()

	client, cleanup := testingutil.SetupMongo(t, "turms_test")
	defer cleanup()

	repo := NewGroupMemberRepository(client)

	// 1. Setup Data
	groupID := int64(1)
	memberID := int64(100)
	nonMemberID := int64(200)

	err := repo.Insert(ctx, &po.GroupMember{
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
