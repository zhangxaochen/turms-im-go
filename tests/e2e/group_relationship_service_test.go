package e2e_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"im.turms/server/internal/domain/common/infra/idgen"
	"im.turms/server/internal/domain/group/constant"
	"im.turms/server/internal/domain/group/po"
	"im.turms/server/internal/domain/group/repository"
	"im.turms/server/internal/domain/group/service"
	user_repository "im.turms/server/internal/domain/user/repository"
	user_service "im.turms/server/internal/domain/user/service"
	"im.turms/server/internal/testingutil"
	"im.turms/server/pkg/protocol"
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
	joinReqRepo := repository.NewGroupJoinRequestRepository(mongoClient, nil)
	invRepo := repository.NewGroupInvitationRepository(mongoClient, nil)
	blockRepo := repository.NewGroupBlockedUserRepository(mongoClient)
	questionRepo := repository.NewGroupJoinQuestionRepository(mongoClient)
	memberRepo := repository.NewGroupMemberRepository(mongoClient)
	groupRepo := repository.NewGroupRepository(mongoClient)
	typeRepo := repository.NewGroupTypeRepository(mongoClient)
	groupVerRepo := repository.NewGroupVersionRepository(mongoClient)
	userVerRepo := user_repository.NewUserVersionRepository(mongoClient)
	idGen, _ := idgen.NewSnowflakeIdGenerator(1, 1)

	// 3. Initialize Services
	typeSvc := service.NewGroupTypeService(typeRepo)
	groupVerSvc := service.NewGroupVersionService(groupVerRepo)
	userVerSvc := user_service.NewUserVersionService(userVerRepo)
	blockSvc := service.NewGroupBlocklistService(blockRepo, groupVerSvc)

	memberSvc := service.NewGroupMemberService(groupRepo, memberRepo, groupVerSvc, typeSvc)
	groupSvc := service.NewGroupService(groupRepo)

	// Break circular dependencies
	memberSvc.SetGroupService(groupSvc)
	memberSvc.SetGroupBlocklistService(blockSvc)
	blockSvc.SetGroupMemberService(memberSvc)
	groupSvc.SetGroupMemberService(memberSvc)

	joinReqSvc := service.NewGroupJoinRequestService(joinReqRepo, memberSvc, blockSvc, groupSvc, typeSvc, groupVerSvc, userVerSvc)
	invSvc := service.NewGroupInvitationService(invRepo, memberSvc, groupSvc, typeSvc, groupVerSvc, userVerSvc, idGen)
	questionSvc := service.NewGroupQuestionService(questionRepo, memberSvc, groupSvc, groupVerSvc)

	groupID := int64(10001)
	adminID := int64(20001)
	userID := int64(30001)

	var err error

	// --- 4. Prepare Test Data ---
	// Ensure Default Group Type
	err = typeRepo.InsertGroupType(ctx, &po.GroupType{
		ID:           0,
		Name:         "DEFAULT",
		JoinStrategy: constant.GroupJoinStrategy_JOIN_REQUEST,
	})
	require.NoError(t, err)

	// Insert Group
	isActive := true
	typeID := int64(0)
	err = groupRepo.InsertGroup(ctx, &po.Group{
		ID:       groupID,
		TypeID:   &typeID,
		OwnerID:  &adminID,
		IsActive: &isActive,
	})
	require.NoError(t, err)

	// Add Admin as Owner
	now := time.Now()
	err = memberRepo.AddGroupMember(ctx, &po.GroupMember{
		ID: po.GroupMemberKey{
			GroupID: groupID,
			UserID:  adminID,
		},
		Role:     protocol.GroupMemberRole_OWNER,
		JoinDate: &now,
	})
	require.NoError(t, err)

	// --- 5. Test Group Blocklist ---
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
		q, err := questionSvc.AuthAndCreateQuestion(ctx, adminID, groupID, "Ready?", answers, 100)
		require.NoError(t, err)
		assert.NotZero(t, q.ID)
		assert.Equal(t, groupID, q.GroupID)

		questions, err := questionSvc.QueryJoinQuestions(ctx, groupID)
		require.NoError(t, err)
		assert.Len(t, questions, 1)

		q1Str := "Updated Question?"
		scoreStr := 200
		err = questionSvc.AuthAndUpdateQuestion(ctx, adminID, groupID, q.ID, &q1Str, answers, &scoreStr)
		require.NoError(t, err)

		questions, err = questionSvc.QueryJoinQuestions(ctx, groupID)
		require.NoError(t, err)
		assert.Equal(t, "Updated Question?", questions[0].Question)
		assert.Equal(t, 200, questions[0].Score)

		result, err := questionSvc.CheckGroupJoinQuestionsAnswersAndJoin(ctx, userID, map[int64]string{
			q.ID: "yes",
		})
		require.NoError(t, err)
		assert.True(t, result.Joined)

		err = questionSvc.AuthAndDeleteQuestion(ctx, adminID, groupID, q.ID)
		require.NoError(t, err)

		questions, err = questionSvc.QueryJoinQuestions(ctx, groupID)
		require.NoError(t, err)
		assert.Len(t, questions, 0)
	})
}
