package e2e_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	grouppo "im.turms/server/internal/domain/group/po"
	grouprepo "im.turms/server/internal/domain/group/repository"
	groupservice "im.turms/server/internal/domain/group/service"
	"im.turms/server/internal/testingutil"
)

func TestGroup_E2E_Lifecycle(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E tests requiring Mongo in short mode")
	}
	ctx := context.Background()

	// 1. Infra Boot
	mongoClient, mCleanup := testingutil.SetupMongo(t, "turms_group_e2e_test")
	defer mCleanup()

	// 2. Repositories
	groupRepo := grouprepo.NewGroupRepository(mongoClient)
	groupMemRepo := grouprepo.NewGroupMemberRepository(mongoClient)

	// 3. Services
	groupSvc := groupservice.NewGroupService(groupRepo)
	groupMemSvc := groupservice.NewGroupMemberService(groupRepo, groupMemRepo)

	creatorID := int64(100)
	groupID := int64(1)
	groupName := "Turms Alpha Test Group"

	t.Run("Create Group", func(t *testing.T) {
		group, err := groupSvc.CreateGroup(ctx, creatorID, groupID, &groupName, nil, nil)
		require.NoError(t, err)
		assert.NotNil(t, group)
		assert.Equal(t, groupID, group.ID)
		assert.Equal(t, creatorID, *group.OwnerID)
	})

	t.Run("Set Creator as Owner", func(t *testing.T) {
		now := time.Now()
		err := groupMemRepo.AddGroupMember(ctx, &grouppo.GroupMember{
			ID:       grouppo.GroupMemberKey{GroupID: groupID, UserID: creatorID},
			Role:     grouppo.GroupMemberRole_OWNER,
			JoinDate: &now,
		})
		require.NoError(t, err)
	})

	t.Run("Owner Adds Manager", func(t *testing.T) {
		userB := int64(200)
		err := groupMemSvc.AddGroupMember(ctx, creatorID, userB, groupID, grouppo.GroupMemberRole_MANAGER)
		require.NoError(t, err)

		role, err := groupMemRepo.FindGroupMemberRole(ctx, groupID, userB)
		require.NoError(t, err)
		assert.Equal(t, grouppo.GroupMemberRole_MANAGER, *role)
	})

	t.Run("Unauthorized User Cannot Add Members", func(t *testing.T) {
		userC := int64(300) // Not in group
		userD := int64(400)

		err := groupMemSvc.AddGroupMember(ctx, userC, userD, groupID, grouppo.GroupMemberRole_MEMBER)
		assert.Error(t, err)
		assert.Equal(t, groupservice.ErrUnauthorized, err)
	})

	t.Run("Soft Delete Group by Owner", func(t *testing.T) {
		err := groupSvc.DeleteGroup(ctx, creatorID, groupID)
		require.NoError(t, err)

		// Verify Soft Delete
		groups, err := groupRepo.FindGroups(ctx, []int64{groupID})
		require.NoError(t, err)
		assert.Len(t, groups, 0, "Group should be filtered out because it is soft-deleted")
	})
}
