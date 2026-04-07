package e2e_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	grouprepo "im.turms/server/internal/domain/group/repository"
	groupservice "im.turms/server/internal/domain/group/service"
	"im.turms/server/internal/testingutil"
	"im.turms/server/pkg/protocol"
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
	typeRepo := grouprepo.NewGroupTypeRepository(mongoClient)
	versionRepo := grouprepo.NewGroupVersionRepository(mongoClient)

	// 3. Services
	groupSvc := groupservice.NewGroupService(groupRepo)
	typeSvc := groupservice.NewGroupTypeService(typeRepo)
	versionSvc := groupservice.NewGroupVersionService(versionRepo)
	groupMemSvc := groupservice.NewGroupMemberService(groupRepo, groupMemRepo, versionSvc, typeSvc)

	// Break circular dependencies
	groupSvc.SetGroupMemberService(groupMemSvc)
	groupMemSvc.SetGroupService(groupSvc)

	creatorID := int64(100)
	groupName := "Turms Alpha Test Group"
	var groupID int64

	t.Run("Create Group", func(t *testing.T) {
		group, err := groupSvc.CreateGroup(ctx, creatorID, creatorID, &groupName, nil, nil, nil, nil, nil, nil, nil, nil)
		require.NoError(t, err)
		assert.NotNil(t, group)
		assert.Equal(t, creatorID, *group.OwnerID)
		groupID = group.ID
	})

	// Creator is already added as owner by CreateGroup in the service layer

	t.Run("Owner Adds Manager", func(t *testing.T) {
		userB := int64(200)
		err := groupMemSvc.AddGroupMember(ctx, groupID, userB, protocol.GroupMemberRole_MANAGER, &creatorID, nil)
		require.NoError(t, err)

		role, err := groupMemRepo.FindGroupMemberRole(ctx, groupID, userB)
		require.NoError(t, err)
		assert.Equal(t, protocol.GroupMemberRole_MANAGER, *role)
	})

	t.Run("Unauthorized User Cannot Add Members", func(t *testing.T) {
		userC := int64(300) // Not in group
		userD := int64(400)

		err := groupMemSvc.AddGroupMember(ctx, groupID, userD, protocol.GroupMemberRole_MEMBER, &userC, nil)
		assert.Error(t, err)
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
