package service

import (
	"context"
	"time"

	"im.turms/server/internal/domain/group/po"
	"im.turms/server/internal/domain/group/repository"
)

type GroupBlocklistService struct {
	blockedUserRepo repository.GroupBlockedUserRepository
}

func NewGroupBlocklistService(blockedUserRepo repository.GroupBlockedUserRepository) *GroupBlocklistService {
	return &GroupBlocklistService{
		blockedUserRepo: blockedUserRepo,
	}
}

func (s *GroupBlocklistService) BlockUser(ctx context.Context, groupID int64, userID int64, requesterID int64) error {
	now := time.Now()
	blockedUser := &po.GroupBlockedUser{
		ID: po.GroupBlockedUserKey{
			GroupID: groupID,
			UserID:  userID,
		},
		BlockDate:   &now,
		RequesterID: requesterID,
	}
	return s.blockedUserRepo.Insert(ctx, blockedUser)
}

func (s *GroupBlocklistService) UnblockUser(ctx context.Context, groupID int64, userID int64) error {
	return s.blockedUserRepo.Delete(ctx, groupID, userID)
}

func (s *GroupBlocklistService) QueryBlockedUsers(ctx context.Context, groupID int64) ([]po.GroupBlockedUser, error) {
	return s.blockedUserRepo.FindBlockedUsersByGroupID(ctx, groupID)
}

func (s *GroupBlocklistService) IsBlocked(ctx context.Context, groupID int64, userID int64) (bool, error) {
	return s.blockedUserRepo.Exists(ctx, groupID, userID)
}
