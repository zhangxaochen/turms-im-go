package e2e_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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

	relSvc := service.NewUserRelationshipService(relRepo, mongoClient)
	defer relSvc.Close()

	reqSvc := service.NewUserFriendRequestService(reqRepo, relSvc)

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

	// Send Request
	req, err := reqSvc.CreateFriendRequest(ctx, user1ID, user2ID, "Hello! Please add me.")
	require.NoError(t, err)
	assert.NotNil(t, req)
	assert.Equal(t, po.RequestStatusPending, req.Status)

	// Sending again should fail (spam prevention)
	_, err = reqSvc.CreateFriendRequest(ctx, user1ID, user2ID, "Hello again!")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already have a pending request")

	// User 2 accepts request
	accepted, err := reqSvc.HandleFriendRequest(ctx, req.ID, user1ID, user2ID, po.ResponseActionAccept, nil)
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

