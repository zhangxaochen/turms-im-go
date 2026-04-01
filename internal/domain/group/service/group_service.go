package service

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"im.turms/server/internal/domain/common/infra/idgen"
	"im.turms/server/internal/domain/group/po"
	"im.turms/server/internal/domain/group/repository"
)

type GroupService interface {
	CreateGroup(ctx context.Context, creatorID int64, ownerID int64, name string, intro string, announcement string, minimumScore int32, isActive bool) (*po.Group, error)
	FindGroup(ctx context.Context, groupID int64) (*po.Group, error)
	UpdateGroup(ctx context.Context, groupID int64, update bson.M) error
}

type groupService struct {
	idGen *idgen.SnowflakeIdGenerator
	repo  *repository.GroupRepository
}

func NewGroupService(idGen *idgen.SnowflakeIdGenerator, repo *repository.GroupRepository) GroupService {
	return &groupService{
		idGen: idGen,
		repo:  repo,
	}
}

func (s *groupService) CreateGroup(ctx context.Context, creatorID int64, ownerID int64, name string, intro string, announcement string, minimumScore int32, isActive bool) (*po.Group, error) {
	groupID := s.idGen.NextIncreasingId()
	group := &po.Group{
		ID:           groupID,
		CreatorID:    creatorID,
		OwnerID:      ownerID,
		Name:         name,
		Intro:        intro,
		Announcement: announcement,
		MinimumScore: minimumScore,
		CreationDate: time.Now(),
		IsActive:     isActive,
	}
	err := s.repo.InsertGroup(ctx, group)
	if err != nil {
		return nil, err
	}
	return group, nil
}

func (s *groupService) FindGroup(ctx context.Context, groupID int64) (*po.Group, error) {
	return s.repo.FindGroup(ctx, groupID)
}

func (s *groupService) UpdateGroup(ctx context.Context, groupID int64, update bson.M) error {
	return s.repo.UpdateGroup(ctx, groupID, update)
}
