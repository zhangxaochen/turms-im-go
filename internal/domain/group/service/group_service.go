package service

import (
	"context"
	"errors"
	"time"

	"im.turms/server/internal/domain/group/po"
	"im.turms/server/internal/domain/group/repository"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	ErrGroupNotFound = errors.New("group not found")
	ErrNotGroupOwner = errors.New("not the group owner")
)

type GroupService struct {
	groupRepo *repository.GroupRepository
}

func NewGroupService(groupRepo *repository.GroupRepository) *GroupService {
	return &GroupService{
		groupRepo: groupRepo,
	}
}

// CreateGroup creates a new group.
func (s *GroupService) CreateGroup(ctx context.Context, creatorID, groupID int64, name, intro *string, minimumScore *int32) (*po.Group, error) {
	now := time.Now()
	group := &po.Group{
		ID:           groupID,
		CreatorID:    &creatorID,
		OwnerID:      &creatorID,
		Name:         name,
		Intro:        intro,
		MinimumScore: minimumScore,
		CreationDate: &now,
	}
	
	err := s.groupRepo.InsertGroup(ctx, group)
	if err != nil {
		return nil, err
	}
	return group, nil
}

// DeleteGroup performs a soft deletion of the group.
// Only the owner can delete the group.
func (s *GroupService) DeleteGroup(ctx context.Context, requesterID, groupID int64) error {
	ownerID, err := s.groupRepo.FindGroupOwnerID(ctx, groupID)
	if err != nil {
		return err
	}
	if ownerID == nil {
		return ErrGroupNotFound
	}

	if *ownerID != requesterID {
		return ErrNotGroupOwner
	}

	now := time.Now()
	update := bson.M{}
	update["dd"] = now

	return s.groupRepo.UpdateGroup(ctx, groupID, update)
}
