package controller

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/proto"

	"go.mongodb.org/mongo-driver/mongo"
	"im.turms/server/internal/domain/gateway/session"
	"im.turms/server/internal/domain/user/bo"
	"im.turms/server/internal/domain/user/po"
	turmsmongo "im.turms/server/internal/storage/mongo"
	"im.turms/server/pkg/protocol"
)

// Mock services
type mockUserRelationshipService struct {
	mock.Mock
}

func (m *mockUserRelationshipService) UpsertOneSidedRelationship(ctx context.Context, ownerID int64, relatedUserID int64, blockDate *time.Time, groupIndex *int32, establishmentDate *time.Time, name *string, session *mongo.Session) (bo.UpsertRelationshipResult, error) {
	args := m.Called(ctx, ownerID, relatedUserID, blockDate, groupIndex, establishmentDate, name, session)
	return args.Get(0).(bo.UpsertRelationshipResult), args.Error(1)
}

func (m *mockUserRelationshipService) UpdateUserOneSidedRelationships(ctx context.Context, userID int64, relatedUserIDs []int64, blockDate *time.Time, groupIndex *int32, deleteGroupIndex *int32, name *string, lastUpdatedDate *time.Time) error {
	args := m.Called(ctx, userID, relatedUserIDs, blockDate, groupIndex, deleteGroupIndex, name, lastUpdatedDate)
	return args.Error(0)
}

func (m *mockUserRelationshipService) BlockUser(ctx context.Context, ownerID, relatedUserID int64) error {
	args := m.Called(ctx, ownerID, relatedUserID)
	return args.Error(0)
}

func (m *mockUserRelationshipService) UnblockUser(ctx context.Context, ownerID, relatedUserID int64) error {
	args := m.Called(ctx, ownerID, relatedUserID)
	return args.Error(0)
}

func (m *mockUserRelationshipService) TryDeleteTwoSidedRelationships(ctx context.Context, user1ID int64, user2ID int64, session *mongo.Session) error {
	args := m.Called(ctx, user1ID, user2ID, session)
	return args.Error(0)
}

func (m *mockUserRelationshipService) DeleteAllRelationships(ctx context.Context, userIDs []int64, session *mongo.Session) error {
	args := m.Called(ctx, userIDs, session)
	return args.Error(0)
}

func (m *mockUserRelationshipService) DeleteOneSidedRelationships(ctx context.Context, ownerID int64, relatedUserIDs []int64, session *mongo.Session) error {
	args := m.Called(ctx, ownerID, relatedUserIDs, session)
	return args.Error(0)
}

func (m *mockUserRelationshipService) DeleteOneSidedRelationship(ctx context.Context, ownerID int64, relatedUserID int64) error {
	args := m.Called(ctx, ownerID, relatedUserID)
	return args.Error(0)
}

func (m *mockUserRelationshipService) IsBlocked(ctx context.Context, ownerID, relatedUserID int64) (bool, error) {
	args := m.Called(ctx, ownerID, relatedUserID)
	return args.Bool(0), args.Error(1)
}

func (m *mockUserRelationshipService) IsNotBlocked(ctx context.Context, ownerID, relatedUserID int64) (bool, error) {
	args := m.Called(ctx, ownerID, relatedUserID)
	return args.Bool(0), args.Error(1)
}

func (m *mockUserRelationshipService) FriendTwoUsers(ctx context.Context, user1ID, user2ID int64) error {
	args := m.Called(ctx, user1ID, user2ID)
	return args.Error(0)
}

func (m *mockUserRelationshipService) QueryRelationships(ctx context.Context, ownerIDs []int64, relatedUserIDs []int64, groupIndexes []int32, isBlocked *bool, establishmentDateRange *turmsmongo.DateRange, page *int, size *int) ([]po.UserRelationship, error) {
	args := m.Called(ctx, ownerIDs, relatedUserIDs, groupIndexes, isBlocked, establishmentDateRange, page, size)
	return args.Get(0).([]po.UserRelationship), args.Error(1)
}

func (m *mockUserRelationshipService) QueryRelatedUserIds(ctx context.Context, ownerIDs []int64, groupIndexes []int32, isBlocked *bool, page *int, size *int) ([]int64, error) {
	args := m.Called(ctx, ownerIDs, groupIndexes, isBlocked, page, size)
	return args.Get(0).([]int64), args.Error(1)
}

func (m *mockUserRelationshipService) QueryRelationshipsWithVersion(ctx context.Context, userID int64, relatedUserIDs []int64, groupIndexes []int32, blocked *bool, lastUpdatedDate *time.Time) ([]po.UserRelationship, *time.Time, error) {
	args := m.Called(ctx, userID, relatedUserIDs, groupIndexes, blocked, lastUpdatedDate)
	return args.Get(0).([]po.UserRelationship), args.Get(1).(*time.Time), args.Error(2)
}

func (m *mockUserRelationshipService) QueryRelatedUserIdsWithVersion(ctx context.Context, userID int64, groupIndexes []int32, blocked *bool, lastUpdatedDate *time.Time) ([]int64, *time.Time, error) {
	args := m.Called(ctx, userID, groupIndexes, blocked, lastUpdatedDate)
	return args.Get(0).([]int64), args.Get(1).(*time.Time), args.Error(2)
}

func (m *mockUserRelationshipService) CountRelationships(ctx context.Context, ownerIDs []int64, relatedUserIDs []int64, groupIndexes []int32, isBlocked *bool) (int64, error) {
	args := m.Called(ctx, ownerIDs, relatedUserIDs, groupIndexes, isBlocked)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockUserRelationshipService) HasRelationshipAndNotBlocked(ctx context.Context, ownerID, relatedUserID int64) (bool, error) {
	args := m.Called(ctx, ownerID, relatedUserID)
	return args.Bool(0), args.Error(1)
}

type mockUserFriendRequestService struct {
	mock.Mock
}

func (m *mockUserFriendRequestService) RemoveAllExpiredFriendRequests(ctx context.Context, expirationDate time.Time) error {
	args := m.Called(ctx, expirationDate)
	return args.Error(0)
}

func (m *mockUserFriendRequestService) HasPendingFriendRequest(ctx context.Context, requesterID, recipientID int64) (bool, error) {
	args := m.Called(ctx, requesterID, recipientID)
	return args.Bool(0), args.Error(1)
}

func (m *mockUserFriendRequestService) CreateFriendRequest(ctx context.Context, requestID *int64, requesterID, recipientID int64, content string, status *po.RequestStatus, creationDate, responseDate *time.Time, reason *string) (*po.UserFriendRequest, error) {
	args := m.Called(ctx, requestID, requesterID, recipientID, content, status, creationDate, responseDate, reason)
	return args.Get(0).(*po.UserFriendRequest), args.Error(1)
}

func (m *mockUserFriendRequestService) AuthAndCreateFriendRequest(ctx context.Context, requesterID int64, recipientID int64, content string, creationDate time.Time) (*po.UserFriendRequest, error) {
	args := m.Called(ctx, requesterID, recipientID, content, creationDate)
	return args.Get(0).(*po.UserFriendRequest), args.Error(1)
}

func (m *mockUserFriendRequestService) AuthAndRecallFriendRequest(ctx context.Context, requesterID, requestID int64) (*po.UserFriendRequest, error) {
	args := m.Called(ctx, requesterID, requestID)
	return args.Get(0).(*po.UserFriendRequest), args.Error(1)
}

func (m *mockUserFriendRequestService) UpdatePendingFriendRequestStatus(ctx context.Context, requestID int64, targetStatus po.RequestStatus, reason *string) (bool, error) {
	args := m.Called(ctx, requestID, targetStatus, reason)
	return args.Bool(0), args.Error(1)
}

func (m *mockUserFriendRequestService) UpdateFriendRequests(ctx context.Context, requestIds []int64, requesterID, recipientID *int64, content *string, status *po.RequestStatus, reason *string, creationDate *time.Time, responseDate *time.Time) error {
	args := m.Called(ctx, requestIds, requesterID, recipientID, content, status, reason, creationDate, responseDate)
	return args.Error(0)
}

func (m *mockUserFriendRequestService) QueryRecipientId(ctx context.Context, requestID int64) (int64, error) {
	args := m.Called(ctx, requestID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockUserFriendRequestService) QueryRequesterIdAndRecipientIdAndStatus(ctx context.Context, requestID int64) (*po.UserFriendRequest, error) {
	args := m.Called(ctx, requestID)
	return args.Get(0).(*po.UserFriendRequest), args.Error(1)
}

func (m *mockUserFriendRequestService) QueryRequesterIdAndRecipientIdAndCreationDateAndStatus(ctx context.Context, requestID int64) (*po.UserFriendRequest, error) {
	args := m.Called(ctx, requestID)
	return args.Get(0).(*po.UserFriendRequest), args.Error(1)
}

func (m *mockUserFriendRequestService) AuthAndHandleFriendRequest(ctx context.Context, friendRequestID int64, requesterID int64, action po.ResponseAction, reason *string) (bool, error) {
	args := m.Called(ctx, friendRequestID, requesterID, action, reason)
	return args.Bool(0), args.Error(1)
}

func (m *mockUserFriendRequestService) QueryFriendRequestsByRecipientId(ctx context.Context, recipientID int64) ([]po.UserFriendRequest, error) {
	args := m.Called(ctx, recipientID)
	return args.Get(0).([]po.UserFriendRequest), args.Error(1)
}

func (m *mockUserFriendRequestService) QueryFriendRequestsByRequesterId(ctx context.Context, requesterID int64) ([]po.UserFriendRequest, error) {
	args := m.Called(ctx, requesterID)
	return args.Get(0).([]po.UserFriendRequest), args.Error(1)
}

func (m *mockUserFriendRequestService) QueryFriendRequestsWithVersion(ctx context.Context, userID int64, isRecipient bool, lastUpdatedDate *time.Time) ([]po.UserFriendRequest, error) {
	args := m.Called(ctx, userID, isRecipient, lastUpdatedDate)
	return args.Get(0).([]po.UserFriendRequest), args.Error(1)
}

func (m *mockUserFriendRequestService) DeleteFriendRequests(ctx context.Context, ids []int64) error {
	args := m.Called(ctx, ids)
	return args.Error(0)
}

func (m *mockUserFriendRequestService) QueryFriendRequests(ctx context.Context, ids, requesterIds, recipientIds []int64, statuses []po.RequestStatus, creationDateStart, creationDateEnd, responseDateStart, responseDateEnd, expirationDateStart, expirationDateEnd *time.Time, page, size *int) ([]po.UserFriendRequest, error) {
	args := m.Called(ctx, ids, requesterIds, recipientIds, statuses, creationDateStart, creationDateEnd, responseDateStart, responseDateEnd, expirationDateStart, expirationDateEnd, page, size)
	return args.Get(0).([]po.UserFriendRequest), args.Error(1)
}

func (m *mockUserFriendRequestService) CountFriendRequests(ctx context.Context, ids, requesterIds, recipientIds []int64, statuses []po.RequestStatus, creationDateStart, creationDateEnd, responseDateStart, responseDateEnd, expirationDateStart, expirationDateEnd *time.Time) (int64, error) {
	args := m.Called(ctx, ids, requesterIds, recipientIds, statuses, creationDateStart, creationDateEnd, responseDateStart, responseDateEnd, expirationDateStart, expirationDateEnd)
	return args.Get(0).(int64), args.Error(1)
}

type mockUserRelationshipGroupService struct {
	mock.Mock
}

func (m *mockUserRelationshipGroupService) CreateRelationshipGroup(ctx context.Context, ownerID int64, groupIndex *int32, groupName string, creationDate *time.Time, session *mongo.Session) (*po.UserRelationshipGroup, error) {
	args := m.Called(ctx, ownerID, groupIndex, groupName, creationDate, session)
	return args.Get(0).(*po.UserRelationshipGroup), args.Error(1)
}

func (m *mockUserRelationshipGroupService) QueryRelationshipGroupsInfos(ctx context.Context, ownerID int64) ([]*po.UserRelationshipGroup, error) {
	args := m.Called(ctx, ownerID)
	return args.Get(0).([]*po.UserRelationshipGroup), args.Error(1)
}

func (m *mockUserRelationshipGroupService) QueryRelationshipGroupsInfosWithVersion(ctx context.Context, ownerID int64, lastUpdatedDate *time.Time) ([]*po.UserRelationshipGroup, *time.Time, error) {
	args := m.Called(ctx, ownerID, lastUpdatedDate)
	return args.Get(0).([]*po.UserRelationshipGroup), args.Get(1).(*time.Time), args.Error(2)
}

func (m *mockUserRelationshipGroupService) QueryGroupIndexes(ctx context.Context, ownerID int64, relatedUserID int64) ([]int32, error) {
	args := m.Called(ctx, ownerID, relatedUserID)
	return args.Get(0).([]int32), args.Error(1)
}

func (m *mockUserRelationshipGroupService) QueryRelationshipGroupMemberIds(ctx context.Context, ownerID int64, groupIndex int32) ([]int64, error) {
	args := m.Called(ctx, ownerID, groupIndex)
	return args.Get(0).([]int64), args.Error(1)
}

func (m *mockUserRelationshipGroupService) UpdateRelationshipGroupName(ctx context.Context, ownerID int64, groupIndex int32, newGroupName string) error {
	args := m.Called(ctx, ownerID, groupIndex, newGroupName)
	return args.Error(0)
}

func (m *mockUserRelationshipGroupService) UpsertRelationshipGroupMember(ctx context.Context, ownerID int64, relatedUserID int64, newGroupIndex *int32, deleteGroupIndex *int32, session *mongo.Session) (*int32, error) {
	args := m.Called(ctx, ownerID, relatedUserID, newGroupIndex, deleteGroupIndex, session)
	return args.Get(0).(*int32), args.Error(1)
}

func (m *mockUserRelationshipGroupService) UpdateRelationshipGroups(ctx context.Context, keys []po.UserRelationshipGroupKey, newName *string, creationDate *time.Time) error {
	args := m.Called(ctx, keys, newName, creationDate)
	return args.Error(0)
}

func (m *mockUserRelationshipGroupService) AddRelatedUserToRelationshipGroup(ctx context.Context, ownerID int64, groupIndex int32, relatedUserID int64, session *mongo.Session) (bool, error) {
	args := m.Called(ctx, ownerID, groupIndex, relatedUserID, session)
	return args.Bool(0), args.Error(1)
}

func (m *mockUserRelationshipGroupService) DeleteRelationshipGroups(ctx context.Context, ownerID int64, groupIndexes []int32, session *mongo.Session) (int64, error) {
	args := m.Called(ctx, ownerID, groupIndexes, session)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockUserRelationshipGroupService) DeleteRelationshipGroupAndMoveMembersToNewGroup(ctx context.Context, ownerID int64, deleteGroupIndex int32, newGroupIndex int32) error {
	args := m.Called(ctx, ownerID, deleteGroupIndex, newGroupIndex)
	return args.Error(0)
}

func (m *mockUserRelationshipGroupService) DeleteAllRelationshipGroups(ctx context.Context, ownerIDs []int64, session *mongo.Session, updateVersion bool) error {
	args := m.Called(ctx, ownerIDs, session, updateVersion)
	return args.Error(0)
}

func (m *mockUserRelationshipGroupService) DeleteRelatedUserFromRelationshipGroup(ctx context.Context, ownerID int64, relatedUserID int64, groupIndex int32, session *mongo.Session, updateVersion bool) (int64, error) {
	args := m.Called(ctx, ownerID, relatedUserID, groupIndex, session, updateVersion)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockUserRelationshipGroupService) DeleteRelatedUserFromAllRelationshipGroups(ctx context.Context, ownerID int64, relatedUserID int64, session *mongo.Session, updateVersion bool) (int64, error) {
	args := m.Called(ctx, ownerID, relatedUserID, session, updateVersion)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockUserRelationshipGroupService) DeleteRelatedUsersFromAllRelationshipGroups(ctx context.Context, ownerID int64, relatedUserIDs []int64, session *mongo.Session, updateVersion bool) (int64, error) {
	args := m.Called(ctx, ownerID, relatedUserIDs, session, updateVersion)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockUserRelationshipGroupService) MoveRelatedUserToNewGroup(ctx context.Context, ownerID int64, relatedUserID int64, currentGroupIndex int32, targetGroupIndex int32, suppressIfAlreadyExists bool, session *mongo.Session) error {
	args := m.Called(ctx, ownerID, relatedUserID, currentGroupIndex, targetGroupIndex, suppressIfAlreadyExists, session)
	return args.Error(0)
}

func (m *mockUserRelationshipGroupService) CountRelationshipGroups(ctx context.Context, ownerIDs []int64, groupIndexes []int32, names []string, creationDateStart *time.Time, creationDateEnd *time.Time) (int64, error) {
	args := m.Called(ctx, ownerIDs, groupIndexes, names, creationDateStart, creationDateEnd)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockUserRelationshipGroupService) CountRelationshipGroupMembers(ctx context.Context, ownerIDs []int64, groupIndexes []int32) (int64, error) {
	args := m.Called(ctx, ownerIDs, groupIndexes)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockUserRelationshipGroupService) QueryRelationshipGroups(ctx context.Context, ownerIDs []int64, groupIndexes []int32, names []string, creationDateStart *time.Time, creationDateEnd *time.Time, page *int, size *int) ([]*po.UserRelationshipGroup, error) {
	args := m.Called(ctx, ownerIDs, groupIndexes, names, creationDateStart, creationDateEnd, page, size)
	return args.Get(0).([]*po.UserRelationshipGroup), args.Error(1)
}

func TestUserRelationshipController_HandleQueryRelationshipsRequest(t *testing.T) {
	mockRelSvc := new(mockUserRelationshipService)
	mockFriendSvc := new(mockUserFriendRequestService)
	mockGroupSvc := new(mockUserRelationshipGroupService)

	c := &UserRelationshipController{
		userRelationshipService:      mockRelSvc,
		userFriendRequestService:     mockFriendSvc,
		userRelationshipGroupService: mockGroupSvc,
	}

	ctx := context.Background()
	s := &session.UserSession{UserID: 1}

	now := time.Now()
	mockRelSvc.On("QueryRelationshipsWithVersion", ctx, int64(1), []int64{2}, []int32(nil), (*bool)(nil), (*time.Time)(nil)).
		Return([]po.UserRelationship{
			{
				ID:                po.UserRelationshipKey{OwnerID: 1, RelatedUserID: 2},
				EstablishmentDate: &now,
			},
		}, &now, nil)

	req := &protocol.TurmsRequest{
		RequestId: proto.Int64(1),
		Kind: &protocol.TurmsRequest_QueryRelationshipsRequest{
			QueryRelationshipsRequest: &protocol.QueryRelationshipsRequest{
				UserIds: []int64{2},
			},
		},
	}

	resp, err := c.HandleQueryRelationshipsRequest(ctx, s, req)

	assert.NoError(t, err)
	assert.Equal(t, int32(1000), *resp.Code)

	data := resp.GetData().GetUserRelationshipsWithVersion()
	assert.NotNil(t, data)
	assert.Len(t, data.UserRelationships, 1)
	assert.Equal(t, int64(2), *data.UserRelationships[0].RelatedUserId)
}
