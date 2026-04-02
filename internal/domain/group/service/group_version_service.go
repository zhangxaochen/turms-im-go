package service

import (
	"context"
	"time"

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

func (s *GroupVersionService) UpdateInvitationsVersion(ctx context.Context, groupID int64) error {
	return s.groupVersionRepo.UpdateInvitationsVersion(ctx, groupID)
}

func (s *GroupVersionService) UpdateJoinRequestsVersion(ctx context.Context, groupID int64) error {
	return s.groupVersionRepo.UpdateJoinRequestsVersion(ctx, groupID)
}

func (s *GroupVersionService) UpdateJoinQuestionsVersion(ctx context.Context, groupID int64) error {
	return s.groupVersionRepo.UpdateJoinQuestionsVersion(ctx, groupID)
}

func (s *GroupVersionService) QueryGroupInvitationsVersion(ctx context.Context, groupID int64) (*time.Time, error) {
	v, err := s.groupVersionRepo.FindVersion(ctx, groupID)
	if err != nil || v == nil {
		return nil, err
	}
	return v.Invitations, nil
}
