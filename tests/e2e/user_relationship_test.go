package e2e_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"im.turms/server/internal/domain/common/infra/idgen"
	"im.turms/server/internal/domain/user/po"
	"im.turms/server/internal/domain/user/repository"
	"im.turms/server/internal/domain/user/service"
	"im.turms/server/internal/testingutil"
)

func TestUserRelationshipLifecycle(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test")
	}

	ctx := context.Background()

	// 1. Setup MongoDB
	mongoClient, mCleanup := testingutil.SetupMongo(t, "turms_user_e2e_test")
	defer mCleanup()

	// 2. Init Repositories & Services
	relRepo := repository.NewUserRelationshipRepository(mongoClient)
	reqRepo := repository.NewUserFriendRequestRepository(mongoClient)
	versionRepo := repository.NewUserVersionRepository(mongoClient)
	groupRepo := repository.NewUserRelationshipGroupRepository(mongoClient)
	groupMemberRepo := repository.NewUserRelationshipGroupMemberRepository(mongoClient)

	userVersionSvc := service.NewUserVersionService(versionRepo)
	groupSvc := service.NewUserRelationshipGroupService(groupRepo, groupMemberRepo, userVersionSvc)
	relSvc := service.NewUserRelationshipService(relRepo, groupSvc, userVersionSvc, mongoClient, nil)

	idGen, err := idgen.NewSnowflakeIdGenerator(1, 1)
	require.NoError(t, err)
	reqSvc := service.NewUserFriendRequestService(idGen, reqRepo, relSvc, userVersionSvc, nil)

	// 3. Test scenario: User 1 wants to add User 2
	var user1ID int64 = 1001
	var user2ID int64 = 1002

	// Initially no relationship
	hasRel, err := relSvc.HasRelationshipAndNotBlocked(ctx, user1ID, user2ID)
	require.NoError(t, err)
	assert.False(t, hasRel)

	isBlocked, err := relSvc.IsBlocked(ctx, user2ID, user1ID)
	require.NoError(t, err)
	assert.False(t, isBlocked)

	// Send Friend Request
	req, err := reqSvc.AuthAndCreateFriendRequest(ctx, user1ID, user2ID, "Hello! Please add me.", time.Now())
	require.NoError(t, err)
	assert.NotNil(t, req)
	assert.Equal(t, po.RequestStatusPending, req.Status)

	// User 2 accepts request
	accepted, err := reqSvc.AuthAndHandleFriendRequest(ctx, req.ID, user2ID, po.ResponseActionAccept, nil)
	require.NoError(t, err)
	assert.True(t, accepted)

	// 4. Verification
	hasRel1, err := relSvc.HasRelationshipAndNotBlocked(ctx, user1ID, user2ID)
	require.NoError(t, err)
	assert.True(t, hasRel1, "user1 should have relationship with user2")

	hasRel2, err := relSvc.HasRelationshipAndNotBlocked(ctx, user2ID, user1ID)
	require.NoError(t, err)
	assert.True(t, hasRel2, "user2 should have bidirectional relationship with user1")

	// User 1 blocks User 2
	err = relSvc.BlockUser(ctx, user1ID, user2ID)
	require.NoError(t, err)

	isBlockedNow, err := relSvc.IsBlocked(ctx, user1ID, user2ID)
	require.NoError(t, err)
	assert.True(t, isBlockedNow, "user1 blocked user2")

	// Check friend status after block
	hasRelAfterBlock, err := relSvc.HasRelationshipAndNotBlocked(ctx, user1ID, user2ID)
	require.NoError(t, err)
	assert.False(t, hasRelAfterBlock, "hasRelationshipAndNotBlocked should be false because user is blocked")

	// 5. Delete relationship
	err = relSvc.DeleteOneSidedRelationship(ctx, user1ID, user2ID)
	require.NoError(t, err)

	hasRelAfterDelete, err := relSvc.HasRelationshipAndNotBlocked(ctx, user1ID, user2ID)
	require.NoError(t, err)
	assert.False(t, hasRelAfterDelete)
}
