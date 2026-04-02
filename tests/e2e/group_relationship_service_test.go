package e2e_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"im.turms/server/internal/domain/group/po"
	"im.turms/server/internal/domain/group/repository"
	"im.turms/server/internal/domain/group/service"
	"im.turms/server/internal/testingutil"
)

func TestGroupRelationshipService(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test")
	}

	ctx := context.Background()

	// 1. Setup MongoDB
	mongoClient, mCleanup := testingutil.SetupMongo(t, "turms_group_svc_e2e_test")
	defer mCleanup()

	// 2. Initialize repositories
	joinReqRepo := repository.NewGroupJoinRequestRepository(mongoClient)
	invRepo := repository.NewGroupInvitationRepository(mongoClient)
	blockRepo := repository.NewGroupBlockedUserRepository(mongoClient)
	questionRepo := repository.NewGroupJoinQuestionRepository(mongoClient)

	// 3. Initialize Services
	joinReqSvc := service.NewGroupJoinRequestService(joinReqRepo)
	invSvc := service.NewGroupInvitationService(invRepo)
	blockSvc := service.NewGroupBlocklistService(blockRepo)
	questionSvc := service.NewGroupQuestionService(questionRepo)

	groupID := int64(10001)
	adminID := int64(20001)
	userID := int64(30001)

	// --- 4. Test Group Blocklist ---
	t.Run("GroupBlocklistService", func(t *testing.T) {
		err := blockSvc.BlockUser(ctx, groupID, userID, adminID)
		require.NoError(t, err)

		isBlocked, err := blockSvc.IsBlocked(ctx, groupID, userID)
		require.NoError(t, err)
		assert.True(t, isBlocked)

		users, err := blockSvc.QueryBlockedUsers(ctx, groupID)
		require.NoError(t, err)
		assert.Len(t, users, 1)

		err = blockSvc.UnblockUser(ctx, groupID, userID)
		require.NoError(t, err)

		isBlocked, err = blockSvc.IsBlocked(ctx, groupID, userID)
		require.NoError(t, err)
		assert.False(t, isBlocked)
	})

	// --- 5. Test Group Join Request ---
	t.Run("GroupJoinRequestService", func(t *testing.T) {
		req, err := joinReqSvc.CreateJoinRequest(ctx, groupID, userID, "Hello")
		require.NoError(t, err)
		assert.Equal(t, po.RequestStatusPending, req.Status)

		updated, err := joinReqSvc.ReplyToJoinRequest(ctx, req.ID, adminID, true)
		require.NoError(t, err)
		assert.True(t, updated)

		// Create another request to test cancel
		req2, err := joinReqSvc.CreateJoinRequest(ctx, groupID, userID, "Hello again")
		require.NoError(t, err)

		updated, err = joinReqSvc.RecallPendingJoinRequest(ctx, req2.ID, userID)
		require.NoError(t, err)
		assert.True(t, updated)
	})

	// --- 6. Test Group Invitation ---
	t.Run("GroupInvitationService", func(t *testing.T) {
		inv, err := invSvc.CreateInvitation(ctx, groupID, adminID, userID, "Welcome")
		require.NoError(t, err)
		assert.Equal(t, po.RequestStatusPending, inv.Status)

		updated, err := invSvc.ReplyToInvitation(ctx, inv.ID, userID, false)
		require.NoError(t, err)
		assert.True(t, updated)

		// Create another invitation to test recall
		inv2, err := invSvc.CreateInvitation(ctx, groupID, adminID, userID, "Welcome 2")
		require.NoError(t, err)

		updated, err = invSvc.RecallPendingInvitation(ctx, inv2.ID, adminID)
		require.NoError(t, err)
		assert.True(t, updated)
	})

	// --- 7. Test Group Join Question ---
	t.Run("GroupQuestionService", func(t *testing.T) {
		answers := []string{"yes", "y"}
		q, err := questionSvc.CreateJoinQuestion(ctx, groupID, "Ready?", answers, 100)
		require.NoError(t, err)
		assert.NotZero(t, q.ID)

		questions, err := questionSvc.QueryJoinQuestions(ctx, groupID)
		require.NoError(t, err)
		assert.Len(t, questions, 1)

		q1Str := "Updated Question?"
		scoreStr := 200
		err = questionSvc.UpdateJoinQuestion(ctx, q.ID, groupID, &q1Str, answers, &scoreStr)
		require.NoError(t, err)

		questions, err = questionSvc.QueryJoinQuestions(ctx, groupID)
		require.NoError(t, err)
		assert.Equal(t, "Updated Question?", questions[0].Question)
		assert.Equal(t, 200, questions[0].Score)

		match, err := questionSvc.CheckGroupQuestionAnswerAndJoin(ctx, userID, q.ID, groupID, "yes")
		require.NoError(t, err)
		assert.True(t, match)

		err = questionSvc.DeleteJoinQuestion(ctx, q.ID)
		require.NoError(t, err)

		questions, err = questionSvc.QueryJoinQuestions(ctx, groupID)
		require.NoError(t, err)
		assert.Len(t, questions, 0)
	})
}
