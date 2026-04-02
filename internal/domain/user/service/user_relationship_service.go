package service

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"im.turms/server/internal/domain/common/cache"
	"im.turms/server/internal/domain/user/bo"
	"im.turms/server/internal/domain/user/po"
	"im.turms/server/internal/domain/user/repository"
	"im.turms/server/internal/infra/exception"
	"im.turms/server/internal/infra/validator"
	turmsmongo "im.turms/server/internal/storage/mongo"
	common "im.turms/server/internal/domain/common/service"
	"im.turms/server/pkg/codes"
	"im.turms/server/pkg/protocol"
)

type UserRelationshipService interface {
	UpsertOneSidedRelationship(ctx context.Context, ownerID int64, relatedUserID int64, blockDate *time.Time, groupIndex *int32, establishmentDate *time.Time, name *string, session *mongo.Session) (bo.UpsertRelationshipResult, error)
	UpdateUserOneSidedRelationships(ctx context.Context, userID int64, relatedUserIDs []int64, blockDate *time.Time, groupIndex *int32, deleteGroupIndex *int32, name *string, lastUpdatedDate *time.Time) error
	BlockUser(ctx context.Context, ownerID, relatedUserID int64) error
	UnblockUser(ctx context.Context, ownerID, relatedUserID int64) error
	TryDeleteTwoSidedRelationships(ctx context.Context, user1ID int64, user2ID int64, session *mongo.Session) error
	DeleteAllRelationships(ctx context.Context, userIDs []int64, session *mongo.Session) error
	DeleteOneSidedRelationships(ctx context.Context, ownerID int64, relatedUserIDs []int64, session *mongo.Session) error
	DeleteOneSidedRelationship(ctx context.Context, ownerID int64, relatedUserID int64) error
	IsBlocked(ctx context.Context, ownerID, relatedUserID int64) (bool, error)
	IsNotBlocked(ctx context.Context, ownerID, relatedUserID int64) (bool, error)
	FriendTwoUsers(ctx context.Context, user1ID, user2ID int64) error
	QueryRelationships(ctx context.Context, ownerIDs []int64, relatedUserIDs []int64, groupIndexes []int32, isBlocked *bool, establishmentDateRange *turmsmongo.DateRange, page *int, size *int) ([]po.UserRelationship, error)
	QueryRelatedUserIds(ctx context.Context, ownerIDs []int64, groupIndexes []int32, isBlocked *bool, page *int, size *int) ([]int64, error)
	QueryRelationshipsWithVersion(ctx context.Context, ownerID int64, relatedUserIDs []int64, groupIndexes []int32, isBlocked *bool, lastUpdatedDate *time.Time) ([]po.UserRelationship, *time.Time, error)
	QueryRelatedUserIdsWithVersion(ctx context.Context, ownerID int64, groupIndexes []int32, isBlocked *bool, lastUpdatedDate *time.Time) ([]int64, *time.Time, error)
	CountRelationships(ctx context.Context, ownerIDs []int64, relatedUserIDs []int64, groupIndexes []int32, isBlocked *bool) (int64, error)
	HasRelationshipAndNotBlocked(ctx context.Context, ownerID, relatedUserID int64) (bool, error)
}

type userRelationshipService struct {
	repo               repository.UserRelationshipRepository
	groupService       UserRelationshipGroupService
	userVersionService *UserVersionService
	mongoClient        *turmsmongo.Client

	relCache               *cache.ShardedMap[string, any]
	blockedCache           *cache.ShardedMap[string, any]
	outboundMessageService common.OutboundMessageService
}

func NewUserRelationshipService(
	repo repository.UserRelationshipRepository,
	groupService UserRelationshipGroupService,
	userVersionService *UserVersionService,
	mongoClient *turmsmongo.Client,
	outboundMessageService common.OutboundMessageService,
) UserRelationshipService {
	return &userRelationshipService{
		repo:                   repo,
		groupService:           groupService,
		userVersionService:     userVersionService,
		mongoClient:            mongoClient,
		relCache:               cache.NewStringShardedMap[any](256),
		blockedCache:           cache.NewStringShardedMap[any](256),
		outboundMessageService: outboundMessageService,
	}
}

func (s *userRelationshipService) UpsertOneSidedRelationship(
	ctx context.Context,
	ownerID int64,
	relatedUserID int64,
	blockDate *time.Time,
	groupIndex *int32,
	establishmentDate *time.Time,
	name *string,
	session *mongo.Session,
) (bo.UpsertRelationshipResult, error) {
	if ownerID == relatedUserID {
		return bo.UpsertRelationshipResult{}, exception.NewTurmsError(int32(codes.IllegalArgument), "ownerID and relatedUserID must not be the same")
	}
	if name != nil {
		if err := validator.MaxLength(name, "name", 50); err != nil {
			return bo.UpsertRelationshipResult{}, err
		}
	}

	return turmsmongo.ExecuteWithSessionResult(ctx, s.mongoClient, session, func(sessCtx mongo.SessionContext, sess *mongo.Session) (bo.UpsertRelationshipResult, error) {
		return s.upsertOneSidedRelationship(sessCtx, ownerID, relatedUserID, blockDate, groupIndex, establishmentDate, name, sess)
	})
}

func (s *userRelationshipService) upsertOneSidedRelationship(
	ctx context.Context,
	ownerID int64,
	relatedUserID int64,
	blockDate *time.Time,
	groupIndex *int32,
	establishmentDate *time.Time,
	name *string,
	session *mongo.Session,
) (bo.UpsertRelationshipResult, error) {
	if establishmentDate == nil {
		now := time.Now()
		establishmentDate = &now
	}

	res, err := s.repo.Upsert(ctx, ownerID, relatedUserID, blockDate, establishmentDate, name, session)
	if err != nil {
		return bo.UpsertRelationshipResult{}, err
	}

	isCreated := res.UpsertedCount > 0
	isUpdated := res.MatchedCount > 0 || isCreated

	var finalGroupIndex *int32
	if isUpdated {
		finalGroupIndex, err = s.groupService.UpsertRelationshipGroupMember(ctx, ownerID, relatedUserID, groupIndex, nil, session)
		if err != nil {
			return bo.UpsertRelationshipResult{}, err
		}
		s.invalidMemberCache(ownerID, relatedUserID)
		_ = s.userVersionService.UpdateRelationshipsVersion(ctx, ownerID)
	}

	if err == nil && s.outboundMessageService != nil {
		s.sendRelationshipNotification(ctx, []int64{relatedUserID}, ownerID, relatedUserID)
	}

	return bo.UpsertRelationshipResult{
		IsCreated:  isCreated,
		GroupIndex: finalGroupIndex,
	}, err
}

func (s *userRelationshipService) UpdateUserOneSidedRelationships(
	ctx context.Context,
	userID int64,
	relatedUserIDs []int64,
	blockDate *time.Time,
	groupIndex *int32,
	deleteGroupIndex *int32,
	name *string,
	lastUpdatedDate *time.Time,
) error {
	if len(relatedUserIDs) == 0 {
		return nil
	}
	if name != nil {
		if err := validator.MaxLength(name, "name", 50); err != nil {
			return err
		}
	}

	return turmsmongo.ExecuteWithSession(ctx, s.mongoClient, nil, func(sessCtx mongo.SessionContext, sess *mongo.Session) error {
		count, err := s.repo.UpdateUserOneSidedRelationships(sessCtx, userID, relatedUserIDs, blockDate, lastUpdatedDate, name, sess)
		if err != nil {
			return err
		}
		if count > 0 {
			for _, relatedUserID := range relatedUserIDs {
				if groupIndex != nil || deleteGroupIndex != nil {
					_, err = s.groupService.UpsertRelationshipGroupMember(sessCtx, userID, relatedUserID, groupIndex, deleteGroupIndex, sess)
					if err != nil {
						return err
					}
				}
				s.invalidMemberCache(userID, relatedUserID)
			}
			return s.userVersionService.UpdateRelationshipsVersion(sessCtx, userID)
		}
		return nil
	})
}

func (s *userRelationshipService) BlockUser(ctx context.Context, ownerID, relatedUserID int64) error {
	now := time.Now()
	_, err := s.UpsertOneSidedRelationship(ctx, ownerID, relatedUserID, &now, nil, nil, nil, nil)
	return err
}

func (s *userRelationshipService) UnblockUser(ctx context.Context, ownerID, relatedUserID int64) error {
	return s.repo.UpdateBlockDate(ctx, ownerID, relatedUserID, nil, nil)
}

func (s *userRelationshipService) TryDeleteTwoSidedRelationships(
	ctx context.Context,
	user1ID int64,
	user2ID int64,
	session *mongo.Session,
) error {
	return turmsmongo.ExecuteWithSession(ctx, s.mongoClient, session, func(sessCtx mongo.SessionContext, sess *mongo.Session) error {
		isBlocked, err := s.IsBlocked(sessCtx, user1ID, user2ID)
		if err != nil {
			return err
		}
		if isBlocked {
			return exception.NewTurmsError(int32(codes.Unauthorized), "cannot delete relationship because you are blocked")
		}

		err = s.DeleteOneSidedRelationships(sessCtx, user1ID, []int64{user2ID}, sess)
		if err != nil {
			return err
		}
		return s.DeleteOneSidedRelationships(sessCtx, user2ID, []int64{user1ID}, sess)
	})
}

func (s *userRelationshipService) DeleteAllRelationships(
	ctx context.Context,
	userIDs []int64,
	session *mongo.Session,
) error {
	if len(userIDs) == 0 {
		return nil
	}
	return turmsmongo.ExecuteWithSession(ctx, s.mongoClient, session, func(sessCtx mongo.SessionContext, sess *mongo.Session) error {
		err := s.groupService.DeleteAllRelationshipGroups(sessCtx, userIDs, sess, false)
		if err != nil {
			return err
		}
		_, err = s.repo.DeleteAllRelationships(sessCtx, userIDs, sess)
		if err != nil {
			return err
		}
		for _, userID := range userIDs {
			_ = s.userVersionService.UpdateRelationshipsVersion(sessCtx, userID)
		}
		return nil
	})
}

func (s *userRelationshipService) DeleteOneSidedRelationships(
	ctx context.Context,
	ownerID int64,
	relatedUserIDs []int64,
	session *mongo.Session,
) error {
	if len(relatedUserIDs) == 0 {
		return nil
	}
	return turmsmongo.ExecuteWithSession(ctx, s.mongoClient, session, func(sessCtx mongo.SessionContext, sess *mongo.Session) error {
		_, err := s.groupService.DeleteRelatedUsersFromAllRelationshipGroups(sessCtx, ownerID, relatedUserIDs, sess, false)
		if err != nil {
			return err
		}
		res, err := s.repo.DeleteOneSidedRelationships(sessCtx, []int64{ownerID}, relatedUserIDs, sess)
		if err != nil {
			return err
		}
		if res.DeletedCount > 0 {
			for _, relatedUserID := range relatedUserIDs {
				s.invalidMemberCache(ownerID, relatedUserID)
				if s.outboundMessageService != nil {
					// Notify relevant users about the deletion
					s.sendRelationshipNotification(sessCtx, []int64{relatedUserID}, ownerID, relatedUserID)
				}
			}
			return s.userVersionService.UpdateRelationshipsVersion(sessCtx, ownerID)
		}
		return nil
	})
}

func (s *userRelationshipService) DeleteOneSidedRelationship(
	ctx context.Context,
	ownerID int64,
	relatedUserID int64,
) error {
	return s.DeleteOneSidedRelationships(ctx, ownerID, []int64{relatedUserID}, nil)
}

func (s *userRelationshipService) IsBlocked(ctx context.Context, ownerID, relatedUserID int64) (bool, error) {
	cacheKey := fmt.Sprintf("%d:%d", ownerID, relatedUserID)
	if val, ok := s.blockedCache.Get(cacheKey); ok {
		return val.(bool), nil
	}

	blocked, err := s.repo.IsBlocked(ctx, ownerID, relatedUserID, nil)
	if err != nil {
		return false, err
	}

	s.blockedCache.Set(cacheKey, blocked)
	return blocked, nil
}

func (s *userRelationshipService) IsNotBlocked(ctx context.Context, ownerID, relatedUserID int64) (bool, error) {
	blocked, err := s.IsBlocked(ctx, ownerID, relatedUserID)
	if err != nil {
		return false, err
	}
	return !blocked, nil
}

func (s *userRelationshipService) FriendTwoUsers(ctx context.Context, user1ID, user2ID int64) error {
	return turmsmongo.ExecuteWithSession(ctx, s.mongoClient, nil, func(sessCtx mongo.SessionContext, sess *mongo.Session) error {
		now := time.Now()
		_, err := s.upsertOneSidedRelationship(sessCtx, user1ID, user2ID, nil, nil, &now, nil, sess)
		if err != nil {
			return err
		}
		_, err = s.upsertOneSidedRelationship(sessCtx, user2ID, user1ID, nil, nil, &now, nil, sess)
		if err == nil && s.outboundMessageService != nil {
			s.sendRelationshipNotification(sessCtx, []int64{user1ID}, user2ID, user1ID)
			s.sendRelationshipNotification(sessCtx, []int64{user2ID}, user1ID, user2ID)
		}
		return err
	})
}

func (s *userRelationshipService) sendRelationshipNotification(ctx context.Context, targetUserIDs []int64, ownerID, relatedUserID int64) {
	relationships, _, err := s.QueryRelationshipsWithVersion(ctx, ownerID, []int64{relatedUserID}, nil, nil, nil)
	if err != nil || len(relationships) == 0 {
		return
	}
	notification := &protocol.TurmsNotification{
		Data: &protocol.TurmsNotification_Data{
			Kind: &protocol.TurmsNotification_Data_UserRelationshipsWithVersion{
				UserRelationshipsWithVersion: &protocol.UserRelationshipsWithVersion{
					UserRelationships: []*protocol.UserRelationship{
						RelationshipToProto(&relationships[0]),
					},
				},
			},
		},
	}
	s.outboundMessageService.ForwardNotificationToMultiple(ctx, notification, targetUserIDs)
}

func (s *userRelationshipService) QueryRelationships(
	ctx context.Context,
	ownerIDs []int64,
	relatedUserIDs []int64,
	groupIndexes []int32,
	isBlocked *bool,
	establishmentDateRange *turmsmongo.DateRange,
	page *int,
	size *int,
) ([]po.UserRelationship, error) {
	return s.repo.FindRelationships(ctx, ownerIDs, relatedUserIDs, groupIndexes, isBlocked, establishmentDateRange, page, size, nil)
}

func (s *userRelationshipService) QueryRelatedUserIds(
	ctx context.Context,
	ownerIDs []int64,
	groupIndexes []int32,
	isBlocked *bool,
	page *int,
	size *int,
) ([]int64, error) {
	return s.repo.FindRelatedUserIDs(ctx, ownerIDs, groupIndexes, isBlocked, page, size, nil)
}

func (s *userRelationshipService) QueryRelationshipsWithVersion(
	ctx context.Context,
	ownerID int64,
	relatedUserIDs []int64,
	groupIndexes []int32,
	isBlocked *bool,
	lastUpdatedDate *time.Time,
) ([]po.UserRelationship, *time.Time, error) {
	version, err := s.userVersionService.QueryRelationshipsLastUpdatedDate(ctx, ownerID)
	if err != nil {
		return nil, nil, err
	}
	if version != nil && lastUpdatedDate != nil && !version.After(*lastUpdatedDate) {
		return nil, nil, exception.NewTurmsError(int32(codes.AlreadyUpToDate), "already up to date")
	}
	rels, err := s.repo.FindRelationships(ctx, []int64{ownerID}, relatedUserIDs, groupIndexes, isBlocked, nil, nil, nil, nil)
	return rels, version, err
}

func (s *userRelationshipService) QueryRelatedUserIdsWithVersion(
	ctx context.Context,
	ownerID int64,
	groupIndexes []int32,
	isBlocked *bool,
	lastUpdatedDate *time.Time,
) ([]int64, *time.Time, error) {
	version, err := s.userVersionService.QueryRelationshipsLastUpdatedDate(ctx, ownerID)
	if err != nil {
		return nil, nil, err
	}
	if version != nil && lastUpdatedDate != nil && !version.After(*lastUpdatedDate) {
		return nil, nil, exception.NewTurmsError(int32(codes.AlreadyUpToDate), "already up to date")
	}
	ids, err := s.repo.FindRelatedUserIDs(ctx, []int64{ownerID}, groupIndexes, isBlocked, nil, nil, nil)
	return ids, version, err
}

func (s *userRelationshipService) CountRelationships(
	ctx context.Context,
	ownerIDs []int64,
	relatedUserIDs []int64,
	groupIndexes []int32,
	isBlocked *bool,
) (int64, error) {
	return s.repo.CountRelationships(ctx, ownerIDs, relatedUserIDs, groupIndexes, isBlocked, nil)
}

func (s *userRelationshipService) invalidMemberCache(ownerID int64, relatedUserID int64) {
	cacheKey := fmt.Sprintf("%d:%d", ownerID, relatedUserID)
	s.relCache.Delete(cacheKey)
	s.blockedCache.Delete(cacheKey)
}

func (s *userRelationshipService) HasRelationshipAndNotBlocked(ctx context.Context, ownerID, relatedUserID int64) (bool, error) {
	cacheKey := fmt.Sprintf("%d:%d", ownerID, relatedUserID)
	if val, ok := s.relCache.Get(cacheKey); ok {
		return val.(bool), nil
	}

	exists, err := s.repo.HasRelationshipAndNotBlocked(ctx, ownerID, relatedUserID, nil)
	if err != nil {
		return false, err
	}

	s.relCache.Set(cacheKey, exists)
	return exists, nil
}
