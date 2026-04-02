package service

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"im.turms/server/internal/domain/common/cache"
	"im.turms/server/internal/domain/user/repository"
	turmsmongo "im.turms/server/internal/storage/mongo"
)

// UserRelationshipService provides methods to check user relationships.
type UserRelationshipService interface {
	HasRelationshipAndNotBlocked(ctx context.Context, ownerID int64, relatedUserID int64) (bool, error)
	IsBlocked(ctx context.Context, ownerID int64, relatedUserID int64) (bool, error)
	FriendTwoUsers(ctx context.Context, user1ID int64, user2ID int64) error
	DeleteOneSidedRelationship(ctx context.Context, ownerID int64, relatedUserID int64) error
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
	return s.mongoClient.Client.UseSession(ctx, func(sessCtx mongo.SessionContext) error {
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
