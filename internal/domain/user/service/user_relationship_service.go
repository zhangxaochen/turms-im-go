package service

import (
	"context"
	"time"

	"im.turms/server/internal/domain/user/po"
	"im.turms/server/internal/domain/user/repository"
)

// UserRelationshipService provides methods to check user relationships.
type UserRelationshipService interface {
	// HasRelationshipAndNotBlocked returns true if ownerID has a friend relationship
	// with relatedUserID and ownerID has NOT blocked relatedUserID.
	HasRelationshipAndNotBlocked(ctx context.Context, ownerID int64, relatedUserID int64) (bool, error)
	// AddFriend adds a friend relationship. If they are already friends, it could update fields.
	AddFriend(ctx context.Context, ownerID int64, relatedUserID int64) (*po.UserRelationship, error)
	// BlockUser blocks a related user for the owner.
	BlockUser(ctx context.Context, ownerID int64, relatedUserID int64) error
}

type userRelationshipService struct {
	repo repository.UserRelationshipRepository
}

func NewUserRelationshipService(repo repository.UserRelationshipRepository) UserRelationshipService {
	return &userRelationshipService{
		repo: repo,
	}
}

func (s *userRelationshipService) HasRelationshipAndNotBlocked(ctx context.Context, ownerID int64, relatedUserID int64) (bool, error) {
	return s.repo.HasRelationshipAndNotBlocked(ctx, ownerID, relatedUserID)
}

func (s *userRelationshipService) AddFriend(ctx context.Context, ownerID int64, relatedUserID int64) (*po.UserRelationship, error) {
	now := time.Now()
	rel := &po.UserRelationship{
		ID: po.UserRelationshipKey{
			OwnerID:       ownerID,
			RelatedUserID: relatedUserID,
		},
		EstablishmentDate: &now,
		// For a new friend, there is no BlockDate
	}
	err := s.repo.Insert(ctx, rel)
	if err != nil {
		return nil, err
	}
	return rel, nil
}

func (s *userRelationshipService) BlockUser(ctx context.Context, ownerID int64, relatedUserID int64) error {
	now := time.Now()
	// Update the BlockDate
	return s.repo.UpdateBlockDate(ctx, ownerID, relatedUserID, &now)
}
