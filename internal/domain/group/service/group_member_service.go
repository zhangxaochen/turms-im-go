package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"im.turms/server/internal/domain/common/cache"
	"im.turms/server/internal/domain/group/po"
	"im.turms/server/internal/domain/group/repository"
	"im.turms/server/pkg/protocol"
)

var (
	ErrUnauthorized   = errors.New("unauthorized role operation")
	ErrRoleNotAllowed = errors.New("target role not allowed")
)

type GroupMemberService struct {
	groupRepo       *repository.GroupRepository
	groupMemberRepo *repository.GroupMemberRepository
	memberCache     *cache.TTLCache[string, bool]
}

func NewGroupMemberService(groupRepo *repository.GroupRepository, groupMemberRepo *repository.GroupMemberRepository) *GroupMemberService {
	return &GroupMemberService{
		groupRepo:       groupRepo,
		groupMemberRepo: groupMemberRepo,
		memberCache:     cache.NewTTLCache[string, bool](1*time.Minute, 10*time.Second),
	}
}

func (s *GroupMemberService) Close() {
	if s.memberCache != nil {
		s.memberCache.Close()
	}
}

// AddGroupMember adds a new member with the given role.
// Only admins (Owner/Manager) can explicitly add someone without an invite.
func (s *GroupMemberService) AddGroupMember(ctx context.Context, requesterID, targetUserID, groupID int64, targetRole protocol.GroupMemberRole) error {
	// Simple RBAC: check if requester is Owner/Manager
	role, err := s.groupMemberRepo.FindGroupMemberRole(ctx, groupID, requesterID)
	if err != nil {
		return err
	}
	if role == nil || (*role != protocol.GroupMemberRole_OWNER && *role != protocol.GroupMemberRole_MANAGER) {
		return ErrUnauthorized
	}

	now := time.Now()
	member := &po.GroupMember{
		ID: po.GroupMemberKey{
			GroupID: groupID,
			UserID:  targetUserID,
		},
		Role:     targetRole,
		JoinDate: &now,
	}

	return s.groupMemberRepo.AddGroupMember(ctx, member)
}

// IsMemberMuted is intensely called during routing.
func (s *GroupMemberService) IsMemberMuted(ctx context.Context, groupID, userID int64) (bool, error) {
	cacheKey := fmt.Sprintf("muted:%d:%d", groupID, userID)
	if muted, ok := s.memberCache.Get(cacheKey); ok {
		return muted, nil
	}

	muted, err := s.groupMemberRepo.IsMemberMuted(ctx, groupID, userID)
	if err != nil {
		return false, err
	}

	s.memberCache.Set(cacheKey, muted)
	return muted, nil
}

func (s *GroupMemberService) IsGroupMember(ctx context.Context, groupID, userID int64) (bool, error) {
	cacheKey := fmt.Sprintf("ismember:%d:%d", groupID, userID)
	if isMember, ok := s.memberCache.Get(cacheKey); ok {
		return isMember, nil
	}

	role, err := s.groupMemberRepo.FindGroupMemberRole(ctx, groupID, userID)
	if err != nil {
		return false, err
	}

	isMember := role != nil
	s.memberCache.Set(cacheKey, isMember)
	return isMember, nil
}
