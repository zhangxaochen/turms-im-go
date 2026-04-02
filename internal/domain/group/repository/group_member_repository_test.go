package repository

import (
	"context"
	"github.com/stretchr/testify/assert"
	"im.turms/server/internal/domain/group/po"
	"im.turms/server/internal/testingutil"
	"testing"
)

func TestGroupMemberRepository_BDD(t *testing.T) {
	ctx := context.Background()

	if testing.Short() {
		t.Skip("Skipping testcases requiring Mongo in short mode")
	}

	client, cleanup := testingutil.SetupMongo(t, "turms_test")
	defer cleanup()

	repo := NewGroupMemberRepository(client)

	// 1. Setup Data
	groupID := int64(1)
	memberID := int64(100)
	nonMemberID := int64(200)

	err := repo.AddGroupMember(ctx, &po.GroupMember{
		ID:   po.GroupMemberKey{GroupID: groupID, UserID: memberID},
		Role: po.GroupMemberRole_MEMBER, // member
	})
	assert.NoError(t, err)

	// 2. Test FindGroupMemberRole
	role, err := repo.FindGroupMemberRole(ctx, groupID, memberID)
	assert.NoError(t, err)
	assert.NotNil(t, role)
	assert.Equal(t, po.GroupMemberRole_MEMBER, *role)

	role2, err := repo.FindGroupMemberRole(ctx, groupID, nonMemberID)
	assert.NoError(t, err)
	assert.Nil(t, role2)
}
