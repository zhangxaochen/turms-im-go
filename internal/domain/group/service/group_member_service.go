package service

import (
	"context"
	"time"

	"im.turms/server/internal/domain/group/po"
	"im.turms/server/internal/domain/group/repository"
)

// GroupMemberService provides methods to check user roles inside groups.
type GroupMemberService interface {
	// IsGroupMember returns true if userID is a common member, admin, or owner of groupID.
	IsGroupMember(ctx context.Context, groupID int64, userID int64) (bool, error)
	// AddGroupMember inserts a new member into the group with a specified role.
	AddGroupMember(ctx context.Context, groupID int64, userID int64, role int32, name *string, muteEndDate *time.Time) (*po.GroupMember, error)
}

type groupMemberService struct {
	repo repository.GroupMemberRepository
}

func NewGroupMemberService(repo repository.GroupMemberRepository) GroupMemberService {
	return &groupMemberService{
		repo: repo,
	}
}

func (s *groupMemberService) IsGroupMember(ctx context.Context, groupID int64, userID int64) (bool, error) {
	return s.repo.IsGroupMember(ctx, groupID, userID)
}

func (s *groupMemberService) AddGroupMember(ctx context.Context, groupID int64, userID int64, role int32, name *string, muteEndDate *time.Time) (*po.GroupMember, error) {
	now := time.Now()
	member := &po.GroupMember{
		ID: po.GroupMemberKey{
			GroupID: groupID,
			UserID:  userID,
		},
		Name:        name,
		Role:        role,
		JoinDate:    &now,
		MuteEndDate: muteEndDate,
	}
	err := s.repo.Insert(ctx, member)
	if err != nil {
		return nil, err
	}
	return member, nil
}
