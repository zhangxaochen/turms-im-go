package service

import (
	"context"
	"fmt"
	"time"

	"im.turms/server/internal/domain/user/po"
	"im.turms/server/internal/domain/user/repository"
)

type UserFriendRequestService interface {
	CreateFriendRequest(ctx context.Context, requesterID, recipientID int64, content string) (*po.UserFriendRequest, error)
	RecallPendingFriendRequest(ctx context.Context, requesterID, recipientID int64) (bool, error)
	HandleFriendRequest(ctx context.Context, requestID int64, requesterID, recipientID int64, action po.ResponseAction, reason *string) (bool, error)
}

type userFriendRequestService struct {
	repo                repository.UserFriendRequestRepository
	relationshipService UserRelationshipService
}

func NewUserFriendRequestService(repo repository.UserFriendRequestRepository, relService UserRelationshipService) UserFriendRequestService {
	return &userFriendRequestService{
		repo:                repo,
		relationshipService: relService,
	}
}

func (s *userFriendRequestService) CreateFriendRequest(ctx context.Context, requesterID, recipientID int64, content string) (*po.UserFriendRequest, error) {
	if requesterID == recipientID {
		return nil, fmt.Errorf("cannot send friend request to yourself")
	}

	// 1. Spam protection
	hasPending, err := s.repo.HasPendingFriendRequest(ctx, requesterID, recipientID)
	if err != nil {
		return nil, err
	}
	if hasPending {
		return nil, fmt.Errorf("already have a pending request")
	}

	// 2. Already friends or blocked?
	isBlocked, err := s.relationshipService.IsBlocked(ctx, recipientID, requesterID)
	if err != nil {
		return nil, err
	}
	if isBlocked {
		return nil, fmt.Errorf("you are blocked by the recipient")
	}

	hasRel, err := s.relationshipService.HasRelationshipAndNotBlocked(ctx, requesterID, recipientID)
	if err != nil {
		return nil, err
	}
	if hasRel {
		return nil, fmt.Errorf("already friends")
	}

	now := time.Now()
	req := &po.UserFriendRequest{
		ID:           now.UnixNano(), // using snowflake / unix nano for simplicity here
		Content:      content,
		Status:       po.RequestStatusPending,
		CreationDate: now,
		RequesterID:  requesterID,
		RecipientID:  recipientID,
	}

	if err := s.repo.Insert(ctx, req); err != nil {
		return nil, err
	}

	return req, nil
}

func (s *userFriendRequestService) RecallPendingFriendRequest(ctx context.Context, requesterID, recipientID int64) (bool, error) {
	// Not strictly atomic recall, but we can do a CAS update via status. Let's do a simple find/update.
	// We'd actually need a repository method to find request ID by requester/recipient.
	// To simplify, we'll mark this logically. Actually since Turms allows multiple, we should update all pending.
	// Let's assume we do this later if needed.
	return false, fmt.Errorf("not implemented")
}

func (s *userFriendRequestService) HandleFriendRequest(ctx context.Context, requestID int64, requesterID, recipientID int64, action po.ResponseAction, reason *string) (bool, error) {
	var targetStatus po.RequestStatus
	switch action {
	case po.ResponseActionAccept:
		targetStatus = po.RequestStatusAccepted
	case po.ResponseActionDecline:
		targetStatus = po.RequestStatusDeclined
	case po.ResponseActionIgnore:
		targetStatus = po.RequestStatusIgnored
	default:
		return false, fmt.Errorf("invalid response action")
	}

	now := time.Now()
	// CAS Update
	success, err := s.repo.UpdateStatusIfPending(ctx, requestID, recipientID, targetStatus, reason, now)
	if err != nil {
		return false, err
	}

	if success && targetStatus == po.RequestStatusAccepted {
		// Create bi-directional relationship!
		// In Java Turms, userVersion is incremented here.
		err = s.relationshipService.FriendTwoUsers(ctx, requesterID, recipientID)
		if err != nil {
			// This means friend request is marked accepted, but relationship failed to add.
			// In fully transactional systems this wouldn't happen, but since MongoDB cross-collection
			// transaction is sometimes heavy, Java Turms tries to do it inside.
			return false, err
		}
	}

	return success, nil
}
