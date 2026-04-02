package service

import (
	"context"

	"im.turms/server/internal/domain/group/repository"
)

type GroupVersionService struct {
	groupVersionRepo *repository.GroupVersionRepository
}

func NewGroupVersionService(groupVersionRepo *repository.GroupVersionRepository) *GroupVersionService {
	return &GroupVersionService{
		groupVersionRepo: groupVersionRepo,
	}
}

func (s *GroupVersionService) InitVersions(ctx context.Context, groupID int64) error {
	return s.groupVersionRepo.InsertVersion(ctx, groupID)
}

func (s *GroupVersionService) UpdateMembersVersion(ctx context.Context, groupID int64) error {
	return s.groupVersionRepo.UpdateMembersVersion(ctx, groupID)
}

func (s *GroupVersionService) UpdateBlocklistVersion(ctx context.Context, groupID int64) error {
	return s.groupVersionRepo.UpdateBlocklistVersion(ctx, groupID)
}
