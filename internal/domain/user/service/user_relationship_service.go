package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"im.turms/server/internal/domain/common/cache"
	"im.turms/server/internal/domain/user/po"
	"im.turms/server/internal/domain/user/repository"
	turmsmongo "im.turms/server/internal/storage/mongo"
)

// UserRelationshipService provides methods to check user relationships.
type UserRelationshipService interface {
	HasRelationshipAndNotBlocked(ctx context.Context, ownerID int64, relatedUserID int64) (bool, error)
	IsBlocked(ctx context.Context, ownerID int64, relatedUserID int64) (bool, error)
	FriendTwoUsers(ctx context.Context, user1ID int64, user2ID int64) error
	DeleteOneSidedRelationship(ctx context.Context, ownerID int64, relatedUserID int64) error
	UpsertOneSidedRelationship(ctx context.Context, ownerID, relatedUserID int64, blockDate *time.Time, groupIndex *int32, establishmentDate *time.Time, name *string) error
	UpdateUserOneSidedRelationships(ctx context.Context, ownerID int64, relatedUserIDs []int64, blockDate *time.Time, groupIndex *int32, establishmentDate *time.Time, name *string) error
	TryDeleteTwoSidedRelationships(ctx context.Context, user1ID, user2ID int64) error
	DeleteAllRelationships(ctx context.Context, ids []int64) error
	DeleteOneSidedRelationships(ctx context.Context, ownerID int64, relatedUserIDs []int64) error
	QueryMembersRelationships(ctx context.Context, ownerID int64, groupIndexes []int32, page, size *int) ([]po.UserRelationship, error)
	CountRelationships(ctx context.Context, ownerIDs, relatedUserIDs []int64, groupIndexes []int32, isBlocked *bool) (int64, error)
	IsNotBlocked(ctx context.Context, ownerID, relatedUserID int64) (bool, error)
	HasOneSidedRelationship(ctx context.Context, ownerID, relatedUserID int64) (bool, error)
	QueryRelationships(ctx context.Context, ownerIDs []int64, relatedUserIDs []int64) ([]po.UserRelationship, error)
	QueryRelatedUserIds(ctx context.Context, ownerID int64, isBlocked *bool) ([]int64, error)
	QueryRelationshipsWithVersion(ctx context.Context, ownerID int64, relatedUserIDs []int64, groupIndexes []int32, isBlocked *bool, lastUpdatedDate *time.Time) ([]po.UserRelationship, error)
	QueryRelatedUserIdsWithVersion(ctx context.Context, ownerID int64, groupIndexes []int32, isBlocked *bool, lastUpdatedDate *time.Time) ([]int64, error)
	BlockUser(ctx context.Context, ownerID int64, relatedUserID int64) error
	Close()
}

type userRelationshipService struct {
	repo         repository.UserRelationshipRepository
	mongoClient  *turmsmongo.Client
	relCache     *cache.TTLCache[string, bool] // "ownerID:relatedUserID" -> hasRelationshipAndNotBlocked
	blockedCache *cache.TTLCache[string, bool] // "ownerID:relatedUserID" -> isBlocked
}

func NewUserRelationshipService(repo repository.UserRelationshipRepository, mongoClient *turmsmongo.Client) UserRelationshipService {
	return &userRelationshipService{
		repo:         repo,
		mongoClient:  mongoClient,
		relCache:     cache.NewTTLCache[string, bool](1*time.Minute, 10*time.Second),
		blockedCache: cache.NewTTLCache[string, bool](1*time.Minute, 10*time.Second),
	}
}

func (s *userRelationshipService) Close() {
	if s.relCache != nil {
		s.relCache.Close()
	}
	if s.blockedCache != nil {
		s.blockedCache.Close()
	}
}

func (s *userRelationshipService) HasRelationshipAndNotBlocked(ctx context.Context, ownerID int64, relatedUserID int64) (bool, error) {
	cacheKey := fmt.Sprintf("%d:%d", ownerID, relatedUserID)
	if hasRel, ok := s.relCache.Get(cacheKey); ok {
		return hasRel, nil
	}

	hasRel, err := s.repo.HasRelationshipAndNotBlocked(ctx, ownerID, relatedUserID)
	if err != nil {
		return false, err
	}
	s.relCache.Set(cacheKey, hasRel)
	return hasRel, nil
}

func (s *userRelationshipService) IsBlocked(ctx context.Context, ownerID int64, relatedUserID int64) (bool, error) {
	cacheKey := fmt.Sprintf("%d:%d", ownerID, relatedUserID)
	if isBlocked, ok := s.blockedCache.Get(cacheKey); ok {
		return isBlocked, nil
	}

	blockedValue := true
	relatedIDs, err := s.repo.FindRelatedUserIDs(ctx, ownerID, &blockedValue)
	if err != nil {
		return false, err
	}

	isBlocked := false
	for _, id := range relatedIDs {
		if id == relatedUserID {
			isBlocked = true
			break
		}
	}
	s.blockedCache.Set(cacheKey, isBlocked)
	return isBlocked, nil
}

func (s *userRelationshipService) FriendTwoUsers(ctx context.Context, user1ID int64, user2ID int64) error {
	if user1ID == user2ID {
		return nil
	}

	// Start a MongoDB session to insert transationally
	err := s.mongoClient.Client.UseSession(ctx, func(sessCtx mongo.SessionContext) error {
		err := sessCtx.StartTransaction()
		if err != nil {
			return err
		}

		now := time.Now()
		// Upsert user1 -> user2
		err = s.repo.Upsert(sessCtx, user1ID, user2ID, nil, nil, &now, nil, sessCtx)
		if err != nil {
			sessCtx.AbortTransaction(sessCtx)
			return err
		}

		// Upsert user2 -> user1
		err = s.repo.Upsert(sessCtx, user2ID, user1ID, nil, nil, &now, nil, sessCtx)
		if err != nil {
			sessCtx.AbortTransaction(sessCtx)
			return err
		}

		err = sessCtx.CommitTransaction(sessCtx)
		if err == nil {
			// Invalidate caches
			s.relCache.Delete(fmt.Sprintf("%d:%d", user1ID, user2ID))
			s.relCache.Delete(fmt.Sprintf("%d:%d", user2ID, user1ID))
		}
		return err
	})

	if err != nil && (strings.Contains(err.Error(), "Transaction numbers are only allowed on a replica set member") || strings.Contains(err.Error(), "Standalone")) {
		// Fallback to non-transactional execution
		now := time.Now()
		if upsertErr := s.repo.Upsert(ctx, user1ID, user2ID, nil, nil, &now, nil, nil); upsertErr != nil {
			return upsertErr
		}
		if upsertErr := s.repo.Upsert(ctx, user2ID, user1ID, nil, nil, &now, nil, nil); upsertErr != nil {
			return upsertErr
		}
		// Invalidate caches
		s.relCache.Delete(fmt.Sprintf("%d:%d", user1ID, user2ID))
		s.relCache.Delete(fmt.Sprintf("%d:%d", user2ID, user1ID))
		return nil
	}

	return err
}

func (s *userRelationshipService) DeleteOneSidedRelationship(ctx context.Context, ownerID int64, relatedUserID int64) error {
	err := s.repo.DeleteById(ctx, ownerID, relatedUserID)
	if err == nil {
		s.relCache.Delete(fmt.Sprintf("%d:%d", ownerID, relatedUserID))
	}
	return err
}

func (s *userRelationshipService) BlockUser(ctx context.Context, ownerID int64, relatedUserID int64) error {
	now := time.Now()
	err := s.repo.UpdateBlockDate(ctx, ownerID, relatedUserID, &now)
	if err == nil {
		s.relCache.Delete(fmt.Sprintf("%d:%d", ownerID, relatedUserID))
		s.blockedCache.Delete(fmt.Sprintf("%d:%d", ownerID, relatedUserID))
	}
	return err
}

func (s *userRelationshipService) UpsertOneSidedRelationship(ctx context.Context, ownerID, relatedUserID int64, blockDate *time.Time, groupIndex *int32, establishmentDate *time.Time, name *string) error {
	err := s.repo.Upsert(ctx, ownerID, relatedUserID, blockDate, groupIndex, establishmentDate, name, nil)
	if err == nil {
		s.relCache.Delete(fmt.Sprintf("%d:%d", ownerID, relatedUserID))
		s.blockedCache.Delete(fmt.Sprintf("%d:%d", ownerID, relatedUserID))
	}
	return err
}

func (s *userRelationshipService) UpdateUserOneSidedRelationships(ctx context.Context, ownerID int64, relatedUserIDs []int64, blockDate *time.Time, groupIndex *int32, establishmentDate *time.Time, name *string) error {
	// A real implementation might use a bulk update, but for now we loop
	for _, relatedUserID := range relatedUserIDs {
		err := s.UpsertOneSidedRelationship(ctx, ownerID, relatedUserID, blockDate, groupIndex, establishmentDate, name)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *userRelationshipService) TryDeleteTwoSidedRelationships(ctx context.Context, user1ID, user2ID int64) error {
	err1 := s.DeleteOneSidedRelationship(ctx, user1ID, user2ID)
	err2 := s.DeleteOneSidedRelationship(ctx, user2ID, user1ID)
	if err1 != nil {
		return err1
	}
	return err2
}

func (s *userRelationshipService) QueryRelationships(ctx context.Context, ownerIDs []int64, relatedUserIDs []int64) ([]po.UserRelationship, error) {
	return s.repo.FindRelationships(ctx, ownerIDs, relatedUserIDs)
}

func (s *userRelationshipService) QueryRelatedUserIds(ctx context.Context, ownerID int64, isBlocked *bool) ([]int64, error) {
	return s.repo.FindRelatedUserIDs(ctx, ownerID, isBlocked)
}

func (s *userRelationshipService) DeleteAllRelationships(ctx context.Context, ids []int64) error {
	return s.repo.DeleteAllRelationships(ctx, ids)
}

func (s *userRelationshipService) DeleteOneSidedRelationships(ctx context.Context, ownerID int64, relatedUserIDs []int64) error {
	return s.repo.DeleteOneSidedRelationships(ctx, []int64{ownerID}, relatedUserIDs)
}

func (s *userRelationshipService) QueryMembersRelationships(ctx context.Context, ownerID int64, groupIndexes []int32, page, size *int) ([]po.UserRelationship, error) {
	return s.repo.QueryMembersRelationships(ctx, ownerID, groupIndexes, page, size)
}

func (s *userRelationshipService) CountRelationships(ctx context.Context, ownerIDs, relatedUserIDs []int64, groupIndexes []int32, isBlocked *bool) (int64, error) {
	return s.repo.CountRelationships(ctx, ownerIDs, relatedUserIDs, groupIndexes, isBlocked)
}

func (s *userRelationshipService) IsNotBlocked(ctx context.Context, ownerID, relatedUserID int64) (bool, error) {
	isBlocked, err := s.IsBlocked(ctx, ownerID, relatedUserID)
	if err != nil {
		return false, err
	}
	return !isBlocked, nil
}

func (s *userRelationshipService) HasOneSidedRelationship(ctx context.Context, ownerID, relatedUserID int64) (bool, error) {
	return s.repo.HasOneSidedRelationship(ctx, ownerID, relatedUserID)
}

func (s *userRelationshipService) QueryRelationshipsWithVersion(ctx context.Context, ownerID int64, relatedUserIDs []int64, groupIndexes []int32, isBlocked *bool, lastUpdatedDate *time.Time) ([]po.UserRelationship, error) {
	return s.QueryRelationships(ctx, []int64{ownerID}, relatedUserIDs)
}

func (s *userRelationshipService) QueryRelatedUserIdsWithVersion(ctx context.Context, ownerID int64, groupIndexes []int32, isBlocked *bool, lastUpdatedDate *time.Time) ([]int64, error) {
	return s.QueryRelatedUserIds(ctx, ownerID, isBlocked)
}
