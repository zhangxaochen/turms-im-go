package e2e_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"im.turms/server/internal/domain/group/po"
	"im.turms/server/internal/domain/group/repository"
	"im.turms/server/internal/testingutil"
)

func TestGroupRelationshipPOValidation(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test")
	}

	ctx := context.Background()

	// 1. Setup MongoDB
	mongoClient, mCleanup := testingutil.SetupMongo(t, "turms_group_e2e_test")
	defer mCleanup()

	// 2. Initialize repositories
	joinReqRepo := repository.NewGroupJoinRequestRepository(mongoClient)
	invRepo := repository.NewGroupInvitationRepository(mongoClient)
	blockRepo := repository.NewGroupBlockedUserRepository(mongoClient)
	questionRepo := repository.NewGroupJoinQuestionRepository(mongoClient)

	groupID := int64(9001)
	userA := int64(1001)
	userB := int64(1002)

	// --- 3. Test GroupJoinRequest ---
	t.Run("GroupJoinRequest", func(t *testing.T) {
		req := &po.GroupJoinRequest{
			ID:           1,
			Content:      "Let me in",
			Status:       po.RequestStatusPending,
			CreationDate: time.Now(),
			GroupID:      groupID,
			RequesterID:  userA,
		}

		err := joinReqRepo.Insert(ctx, req)
		require.NoError(t, err)

		hasPending, err := joinReqRepo.HasPendingJoinRequest(ctx, userA, groupID)
		require.NoError(t, err)
		assert.True(t, hasPending)

		// Accept request
		now := time.Now()
		updated, err := joinReqRepo.UpdateStatusIfPending(ctx, 1, userB, po.RequestStatusAccepted, nil, now)
		require.NoError(t, err)
		assert.True(t, updated)

		hasPending2, err := joinReqRepo.HasPendingJoinRequest(ctx, userA, groupID)
		require.NoError(t, err)
		assert.False(t, hasPending2)
	})

	// --- 4. Test GroupInvitation ---
	t.Run("GroupInvitation", func(t *testing.T) {
		inv := &po.GroupInvitation{
			ID:           1,
			GroupID:      groupID,
			InviterID:    userA,
			InviteeID:    userB,
			Content:      "Join us!",
			Status:       po.RequestStatusPending,
			CreationDate: time.Now(),
		}

		err := invRepo.Insert(ctx, inv)
		require.NoError(t, err)

		hasPending, err := invRepo.HasPendingInvitation(ctx, groupID, userB)
		require.NoError(t, err)
		assert.True(t, hasPending)

		now := time.Now()
		updated, err := invRepo.UpdateStatusIfPending(ctx, 1, po.RequestStatusDeclined, nil, now)
		require.NoError(t, err)
		assert.True(t, updated)
	})

	// --- 5. Test GroupBlockedUser ---
	t.Run("GroupBlockedUser", func(t *testing.T) {
		now := time.Now()
		blocked := &po.GroupBlockedUser{
			ID: po.GroupBlockedUserKey{
				GroupID: groupID,
				UserID:  userA,
			},
			BlockDate:   &now,
			RequesterID: userB,
		}

		err := blockRepo.Insert(ctx, blocked)
		require.NoError(t, err)

		exists, err := blockRepo.Exists(ctx, groupID, userA)
		require.NoError(t, err)
		assert.True(t, exists)

		err = blockRepo.Delete(ctx, groupID, userA)
		require.NoError(t, err)

		exists2, err := blockRepo.Exists(ctx, groupID, userA)
		require.NoError(t, err)
		assert.False(t, exists2)
	})

	// --- 6. Test GroupJoinQuestion ---
	t.Run("GroupJoinQuestion", func(t *testing.T) {
		question := &po.GroupJoinQuestion{
			ID:       1,
			GroupID:  groupID,
			Question: "What is the meaning of life?",
			Answers:  []string{"42", "forty-two"},
			Score:    10,
		}

		err := questionRepo.Insert(ctx, question)
		require.NoError(t, err)

		qs, err := questionRepo.FindQuestionsByGroupID(ctx, groupID)
		require.NoError(t, err)
		assert.Len(t, qs, 1)
		assert.Equal(t, "42", qs[0].Answers[0])

		err = questionRepo.Delete(ctx, 1)
		require.NoError(t, err)
	})
}
