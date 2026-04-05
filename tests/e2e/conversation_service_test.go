package e2e

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"im.turms/server/internal/domain/conversation/repository"
	"im.turms/server/internal/domain/conversation/service"
	grouppo "im.turms/server/internal/domain/group/po"
	"im.turms/server/internal/infra/property"
	"im.turms/server/internal/testingutil"
)

type mockUserRelationshipService struct{}

func (m *mockUserRelationshipService) HasRelationshipAndNotBlocked(ctx context.Context, ownerID int64, relatedUserID int64) (bool, error) {
	return true, nil
}

type mockGroupService struct{}

func (m *mockGroupService) QueryGroupTypeIfActiveAndNotDeleted(ctx context.Context, groupID int64) (*grouppo.GroupType, error) {
	return &grouppo.GroupType{EnableReadReceipt: true}, nil
}

type mockGroupMemberService struct{}

func (m *mockGroupMemberService) IsGroupMember(ctx context.Context, groupID int64, userID int64, activeOnly ...bool) (bool, error) {
	return true, nil
}
func (m *mockGroupMemberService) IsGroupMemberActiveOnly(ctx context.Context, groupID int64, userID int64) (bool, error) {
	return true, nil
}
func (m *mockGroupMemberService) FindGroupMemberIDs(ctx context.Context, groupID int64) ([]int64, error) {
	return []int64{101, 102}, nil
}
func (m *mockGroupMemberService) FindActiveGroupMemberIDs(ctx context.Context, groupID int64) ([]int64, error) {
	return []int64{101, 102}, nil
}
func (m *mockGroupMemberService) QueryUserJoinedGroupIds(ctx context.Context, userID int64) ([]int64, error) {
	return []int64{1001}, nil
}

type mockMessageService struct{}

func (m *mockMessageService) HasPrivateMessage(ctx context.Context, senderID int64, targetID int64) (bool, error) {
	return true, nil
}

func TestConversationCore_E2E(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	db, cleanup := testingutil.SetupMongo(t, "turms_conversation_core_e2e_test")
	defer cleanup()

	privateConvRepo := repository.NewPrivateConversationRepository(db)
	groupConvRepo := repository.NewGroupConversationRepository(db)
	propsManager := property.NewTurmsPropertiesManager()
	localProps := propsManager.GetLocalProperties()
	localProps.Service.Conversation.ReadReceipt.Enabled = true
	localProps.Service.Conversation.ReadReceipt.UseServerTime = false
	localProps.Service.Conversation.ReadReceipt.AllowMoveReadDateForward = true

	convService := service.NewConversationService(
		privateConvRepo,
		groupConvRepo,
		&mockUserRelationshipService{},
		&mockGroupService{},
		&mockGroupMemberService{},
		&mockMessageService{},
		propsManager,
	)

	t.Run("PrivateConversation Lifecycle", func(t *testing.T) {
		ownerID := int64(101)
		targetID := int64(201)
		readDate := time.Now().UTC().Truncate(time.Millisecond)

		err := convService.AuthAndUpsertPrivateConversationReadDate(ctx, ownerID, targetID, &readDate)
		require.NoError(t, err)

		results, err := convService.QueryPrivateConversationsByOwnerIds(ctx, []int64{ownerID})
		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, ownerID, results[0].ID.OwnerID)
		assert.Equal(t, targetID, results[0].ID.TargetID)
		assert.Equal(t, readDate, results[0].ReadDate.UTC())

		// Update again to test UPSERT
		newReadDate := readDate.Add(5 * time.Minute)
		err = convService.AuthAndUpsertPrivateConversationReadDate(ctx, ownerID, targetID, &newReadDate)
		require.NoError(t, err)

		results, err = convService.QueryPrivateConversationsByOwnerIds(ctx, []int64{ownerID})
		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, newReadDate, results[0].ReadDate.UTC())
	})

	t.Run("GroupConversation Lifecycle", func(t *testing.T) {
		groupID := int64(1001)
		memberID1 := int64(101)
		memberID2 := int64(102)
		readDate1 := time.Now().UTC().Truncate(time.Millisecond)
		readDate2 := readDate1.Add(10 * time.Minute)

		err := convService.AuthAndUpsertGroupConversationReadDate(ctx, groupID, memberID1, &readDate1)
		require.NoError(t, err)

		err = convService.AuthAndUpsertGroupConversationReadDate(ctx, groupID, memberID2, &readDate2)
		require.NoError(t, err)

		results, err := convService.QueryGroupConversations(ctx, []int64{groupID})
		require.NoError(t, err)
		assert.Len(t, results, 1)

		groupConv := results[0]
		assert.Equal(t, groupID, groupConv.ID)
		require.NotNil(t, groupConv.MemberIDToReadDate)

		assert.Equal(t, readDate1, groupConv.MemberIDToReadDate["101"].UTC())
		assert.Equal(t, readDate2, groupConv.MemberIDToReadDate["102"].UTC())

		// Test Upsert for an existing member
		newReadDate1 := readDate1.Add(5 * time.Minute)
		err = convService.AuthAndUpsertGroupConversationReadDate(ctx, groupID, memberID1, &newReadDate1)
		require.NoError(t, err)

		results, err = convService.QueryGroupConversations(ctx, []int64{groupID})
		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, newReadDate1, results[0].MemberIDToReadDate["101"].UTC())
	})
}
