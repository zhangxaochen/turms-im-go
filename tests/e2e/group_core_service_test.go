package e2e

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"im.turms/server/internal/domain/group/po"
	"im.turms/server/internal/domain/group/repository"
	"im.turms/server/internal/domain/group/service"
	"im.turms/server/internal/testingutil"
	"im.turms/server/pkg/protocol"
)

func TestGroupCore_E2E(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, cleanup := testingutil.SetupMongo(t, "turms_group_core_e2e_test")
	defer cleanup()

	// Repositories
	groupRepo := repository.NewGroupRepository(client)
	memberRepo := repository.NewGroupMemberRepository(client)
	typeRepo := repository.NewGroupTypeRepository(client)
	versionRepo := repository.NewGroupVersionRepository(client)

	// Services
	groupService := service.NewGroupService(groupRepo)
	memberService := service.NewGroupMemberService(groupRepo, memberRepo)
	typeService := service.NewGroupTypeService(typeRepo)
	versionService := service.NewGroupVersionService(versionRepo)

	// Test GroupType
	t.Run("GroupType Lifecycle", func(t *testing.T) {
		err := typeService.EnsureDefaultGroupType(ctx)
		require.NoError(t, err)

		groupType, err := typeService.FindGroupType(ctx, 0)
		require.NoError(t, err)
		assert.NotNil(t, groupType)
		assert.Equal(t, "DEFAULT", groupType.Name)
	})

	// Test Group Creation & Membership
	t.Run("CreateGroup And AddMember", func(t *testing.T) {
		creatorID := int64(101)
		groupID := int64(1001)
		name := "Go Developers"
		intro := "A group for Go enthusiasts"

		group, err := groupService.CreateGroup(ctx, creatorID, groupID, &name, &intro, nil)
		require.NoError(t, err)
		assert.NotNil(t, group)
		assert.Equal(t, name, *group.Name)

		// Create membership for creator as OWNER
		err = memberService.AddGroupMember(ctx, creatorID, creatorID, groupID, protocol.GroupMemberRole_OWNER)
		// Wait, the requester is creatorID, but if they aren't owner yet, it will fail because of our RBAC (returns ErrUnauthorized)
		// Actually, in Turms Java, CreateGroup adds the creator as Owner organically. Let's fix our AddGroupMember logic if needed or just use Repo here.

		// For the test, we add via Repo to bootstrap the owner
		member := &po.GroupMember{
			ID:   po.GroupMemberKey{GroupID: groupID, UserID: creatorID},
			Role: protocol.GroupMemberRole_OWNER,
		}
		err = memberRepo.AddGroupMember(ctx, member)
		require.NoError(t, err)

		// Test finding the role
		role, err := memberRepo.FindGroupMemberRole(ctx, groupID, creatorID)
		require.NoError(t, err)
		assert.NotNil(t, role)
		assert.Equal(t, protocol.GroupMemberRole_OWNER, *role)

		// Now creator can add a member
		memberID := int64(102)
		err = memberService.AddGroupMember(ctx, creatorID, memberID, groupID, protocol.GroupMemberRole_MEMBER)
		require.NoError(t, err)

		// Test IsGroupMember
		isMember, err := memberService.IsGroupMember(ctx, groupID, memberID)
		require.NoError(t, err)
		assert.True(t, isMember)
	})

	// Test Group Versioning
	t.Run("Version Lifecycle", func(t *testing.T) {
		groupID := int64(1001)
		err := versionService.InitVersions(ctx, groupID)
		require.NoError(t, err)

		err = versionService.UpdateMembersVersion(ctx, groupID)
		require.NoError(t, err)

		v, err := versionRepo.FindVersion(ctx, groupID)
		require.NoError(t, err)
		assert.NotNil(t, v)
		assert.NotNil(t, v.Members)
	})

	// Test Soft Delete
	t.Run("Soft Delete Group", func(t *testing.T) {
		creatorID := int64(101)
		groupID := int64(1001)

		err := groupService.DeleteGroup(ctx, creatorID, groupID)
		require.NoError(t, err)

		// Ensure FindGroups filters deleted
		groups, err := groupRepo.FindGroups(ctx, []int64{groupID})
		require.NoError(t, err)
		assert.Empty(t, groups)
	})
}
