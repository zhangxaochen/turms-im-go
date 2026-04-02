package service

import (
	"context"
	"time"

	"im.turms/server/internal/domain/group/po"
	"im.turms/server/internal/domain/group/repository"
)

type GroupInvitationService struct {
	invRepo repository.GroupInvitationRepository
}

func NewGroupInvitationService(invRepo repository.GroupInvitationRepository) *GroupInvitationService {
	return &GroupInvitationService{
		invRepo: invRepo,
	}
}

func (s *GroupInvitationService) CreateInvitation(ctx context.Context, groupID int64, inviterID int64, inviteeID int64, content string) (*po.GroupInvitation, error) {
	now := time.Now()
	
	// Create identity ID, could be randomly generated or generated via snowflake 
	// To simplify, we assume ID is populated externally or we leave it 0 if Mongo auto increments (Turms uses Snowflake)
	// We'll give it a timestamp-based ID proxy for now
	id := now.UnixNano()

	inv := &po.GroupInvitation{
		ID:           id,
		GroupID:      groupID,
		InviterID:    inviterID,
		InviteeID:    inviteeID,
		Content:      content,
		Status:       po.RequestStatusPending,
		CreationDate: now,
	}
	err := s.invRepo.Insert(ctx, inv)
	if err != nil {
		return nil, err
	}
	return inv, nil
}

func (s *GroupInvitationService) RecallPendingInvitation(ctx context.Context, invitationID int64, inviterID int64) (bool, error) {
	// auth check conceptually happens here or caller checks inviterID
	return s.invRepo.UpdateStatusIfPending(ctx, invitationID, po.RequestStatusCanceled, nil, time.Now())
}

func (s *GroupInvitationService) ReplyToInvitation(ctx context.Context, invitationID int64, inviteeID int64, accept bool) (bool, error) {
	status := po.RequestStatusDeclined
	if accept {
		status = po.RequestStatusAccepted
	}
	// auth check conceptually happens here or caller checks inviteeID
	return s.invRepo.UpdateStatusIfPending(ctx, invitationID, status, nil, time.Now())
}
