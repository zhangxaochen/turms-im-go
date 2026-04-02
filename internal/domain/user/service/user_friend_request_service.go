package service

import (
	"context"
	"fmt"
	"time"

	"im.turms/server/internal/domain/user/po"
	"im.turms/server/internal/domain/user/repository"
)

type UserFriendRequestService interface {
	RemoveAllExpiredFriendRequests(ctx context.Context, expirationDate time.Time) error
	HasPendingFriendRequest(ctx context.Context, requesterID, recipientID int64) (bool, error)
	CreateFriendRequest(ctx context.Context, requestID *int64, requesterID, recipientID int64, content string, status *po.RequestStatus, creationDate, responseDate *time.Time, reason *string) (*po.UserFriendRequest, error)
	AuthAndCreateFriendRequest(ctx context.Context, requesterID, recipientID int64, content string, creationDate time.Time) (*po.UserFriendRequest, error)
	AuthAndRecallFriendRequest(ctx context.Context, requesterID, requestID int64) (*po.UserFriendRequest, error)
	UpdatePendingFriendRequestStatus(ctx context.Context, requestID, recipientID int64, targetStatus po.RequestStatus, reason *string) (bool, error)
	UpdateFriendRequests(ctx context.Context, requestIds []int64, requesterID, recipientID *int64, content *string, status *po.RequestStatus, reason *string, creationDate *time.Time, responseDate *time.Time) error
	QueryRecipientId(ctx context.Context, requestID int64) (int64, error)
	QueryRequesterIdAndRecipientIdAndStatus(ctx context.Context, requestID int64) (*po.UserFriendRequest, error)
	QueryRequesterIdAndRecipientIdAndCreationDateAndStatus(ctx context.Context, requestID int64) (*po.UserFriendRequest, error)
	AuthAndHandleFriendRequest(ctx context.Context, friendRequestID int64, requesterID int64, action po.ResponseAction, reason *string) (bool, error)
	QueryFriendRequestsByRecipientId(ctx context.Context, recipientID int64) ([]po.UserFriendRequest, error)
	QueryFriendRequestsByRequesterId(ctx context.Context, requesterID int64) ([]po.UserFriendRequest, error)
	DeleteFriendRequests(ctx context.Context, ids []int64) error
	QueryFriendRequests(ctx context.Context, ids, requesterIds, recipientIds []int64, statuses []po.RequestStatus, creationDateStart, creationDateEnd, responseDateStart, responseDateEnd, expirationDateStart, expirationDateEnd *time.Time, page, size *int) ([]po.UserFriendRequest, error)
	CountFriendRequests(ctx context.Context, ids, requesterIds, recipientIds []int64, statuses []po.RequestStatus, creationDateStart, creationDateEnd, responseDateStart, responseDateEnd, expirationDateStart, expirationDateEnd *time.Time) (int64, error)
}

type userFriendRequestService struct {
	repo                repository.UserFriendRequestRepository
	relationshipService UserRelationshipService
	userVersionService  UserVersionService
}

func NewUserFriendRequestService(repo repository.UserFriendRequestRepository, relService UserRelationshipService, userVersionService UserVersionService) UserFriendRequestService {
	return &userFriendRequestService{
		repo:                repo,
		relationshipService: relService,
		userVersionService:  userVersionService,
	}
}

func (s *userFriendRequestService) RemoveAllExpiredFriendRequests(ctx context.Context, expirationDate time.Time) error {
	return s.repo.DeleteExpiredData(ctx, expirationDate)
}

func (s *userFriendRequestService) HasPendingFriendRequest(ctx context.Context, requesterID, recipientID int64) (bool, error) {
	return s.repo.HasPendingFriendRequest(ctx, requesterID, recipientID)
}

func (s *userFriendRequestService) CreateFriendRequest(ctx context.Context, requestID *int64, requesterID, recipientID int64, content string, status *po.RequestStatus, creationDate, responseDate *time.Time, reason *string) (*po.UserFriendRequest, error) {
	if requesterID == recipientID {
		return nil, fmt.Errorf("requester == recipient")
	}

	id := time.Now().UnixNano()
	if requestID != nil {
		id = *requestID
	}
	
	now := time.Now()
	var cd time.Time
	if creationDate == nil {
		cd = now
	} else if creationDate.Before(now) {
		cd = *creationDate
	} else {
		cd = now
	}

	st := po.RequestStatusPending
	if status != nil {
		st = *status
	}

	req := &po.UserFriendRequest{
		ID:           id,
		Content:      content,
		Status:       st,
		CreationDate: cd,
		RequesterID:  requesterID,
		RecipientID:  recipientID,
	}
	if reason != nil {
		req.Reason = *reason
	}
	if responseDate != nil {
		req.ResponseDate = *responseDate
	}

	if err := s.repo.Insert(ctx, req); err != nil {
		return nil, err
	}

	_ = s.userVersionService.UpdateReceivedFriendRequestsVersion(ctx, recipientID)
	_ = s.userVersionService.UpdateSentFriendRequestsVersion(ctx, requesterID)

	return req, nil
}

func (s *userFriendRequestService) AuthAndCreateFriendRequest(ctx context.Context, requesterID, recipientID int64, content string, creationDate time.Time) (*po.UserFriendRequest, error) {
	isNotBlocked, err := s.relationshipService.IsBlocked(ctx, recipientID, requesterID)
	if err != nil {
		return nil, err
	}
	if isNotBlocked { // Note: original java name is isNotBlocked but here IsBlocked returns true if blocked
		return nil, fmt.Errorf("blocked user to send friend request")
	}

	requestExists, err := s.repo.HasPendingOrDeclinedOrIgnoredOrExpiredRequest(ctx, requesterID, recipientID)
	if err != nil {
		return nil, err
	}
	if requestExists {
		return nil, fmt.Errorf("create existing friend request")
	}

	return s.CreateFriendRequest(ctx, nil, requesterID, recipientID, content, nil, &creationDate, nil, nil)
}

func (s *userFriendRequestService) AuthAndRecallFriendRequest(ctx context.Context, requesterID, requestID int64) (*po.UserFriendRequest, error) {
	req, err := s.repo.FindRequesterIdAndRecipientIdAndCreationDateAndStatus(ctx, requestID)
	if err != nil {
		return nil, err
	}
	if req == nil {
		return nil, fmt.Errorf("not friend request sender to recall")
	}

	if requesterID != req.RequesterID {
		return nil, fmt.Errorf("not friend request sender to recall")
	}

	if req.Status != po.RequestStatusPending {
		return nil, fmt.Errorf("recall non-pending friend request")
	}

	success, err := s.repo.UpdateStatusIfPending(ctx, requestID, req.RecipientID, po.RequestStatusCanceled, nil, time.Now())
	if err != nil {
		return nil, err
	}
	if !success {
		return nil, fmt.Errorf("recall non-pending friend request")
	}

	_ = s.userVersionService.UpdateReceivedFriendRequestsVersion(ctx, req.RecipientID)
	_ = s.userVersionService.UpdateSentFriendRequestsVersion(ctx, req.RequesterID)

	return req, nil
}

func (s *userFriendRequestService) UpdatePendingFriendRequestStatus(ctx context.Context, requestID, recipientID int64, targetStatus po.RequestStatus, reason *string) (bool, error) {
	if targetStatus == po.RequestStatusPending {
		return false, fmt.Errorf("status must not be pending")
	}
	success, err := s.repo.UpdateStatusIfPending(ctx, requestID, recipientID, targetStatus, reason, time.Now())
	if err != nil {
		return false, err
	}
	if success {
		// Java: queryRecipientId then updateversion
		s.userVersionService.UpdateReceivedFriendRequestsVersion(ctx, recipientID)
	}
	return success, nil
}

func (s *userFriendRequestService) UpdateFriendRequests(ctx context.Context, requestIds []int64, requesterID, recipientID *int64, content *string, status *po.RequestStatus, reason *string, creationDate *time.Time, responseDate *time.Time) error {
	if len(requestIds) == 0 {
		return nil
	}
	return s.repo.UpdateFriendRequests(ctx, requestIds, requesterID, recipientID, content, status, reason, creationDate)
}

func (s *userFriendRequestService) QueryRecipientId(ctx context.Context, requestID int64) (int64, error) {
	return s.repo.FindRecipientId(ctx, requestID)
}

func (s *userFriendRequestService) QueryRequesterIdAndRecipientIdAndStatus(ctx context.Context, requestID int64) (*po.UserFriendRequest, error) {
	return s.repo.FindRequesterIdAndRecipientIdAndStatus(ctx, requestID)
}

func (s *userFriendRequestService) QueryRequesterIdAndRecipientIdAndCreationDateAndStatus(ctx context.Context, requestID int64) (*po.UserFriendRequest, error) {
	return s.repo.FindRequesterIdAndRecipientIdAndCreationDateAndStatus(ctx, requestID)
}

func (s *userFriendRequestService) AuthAndHandleFriendRequest(ctx context.Context, friendRequestID int64, requesterID int64, action po.ResponseAction, reason *string) (bool, error) {
	req, err := s.repo.FindRequesterIdAndRecipientIdAndCreationDateAndStatus(ctx, friendRequestID)
	if err != nil {
		return false, err
	}
	if req == nil {
		return false, fmt.Errorf("not recipient to update friend request")
	}

	if req.RecipientID != requesterID {
		return false, fmt.Errorf("not recipient to update friend request")
	}

	if req.Status != po.RequestStatusPending {
		return false, fmt.Errorf("update non-pending friend request")
	}

	switch action {
	case po.ResponseActionAccept:
		success, err := s.UpdatePendingFriendRequestStatus(ctx, friendRequestID, req.RecipientID, po.RequestStatusAccepted, reason)
		if err != nil {
			return false, err
		}
		if success {
			err = s.relationshipService.FriendTwoUsers(ctx, req.RequesterID, requesterID)
		}
		return success, err
	case po.ResponseActionIgnore:
		return s.UpdatePendingFriendRequestStatus(ctx, friendRequestID, req.RecipientID, po.RequestStatusIgnored, reason)
	case po.ResponseActionDecline:
		return s.UpdatePendingFriendRequestStatus(ctx, friendRequestID, req.RecipientID, po.RequestStatusDeclined, reason)
	default:
		return false, fmt.Errorf("illegal response action")
	}
}

func (s *userFriendRequestService) QueryFriendRequestsByRecipientId(ctx context.Context, recipientID int64) ([]po.UserFriendRequest, error) {
	return s.repo.FindRequestsByRecipientID(ctx, recipientID)
}

func (s *userFriendRequestService) QueryFriendRequestsByRequesterId(ctx context.Context, requesterID int64) ([]po.UserFriendRequest, error) {
	return s.repo.FindRequestsByRequesterID(ctx, requesterID)
}

func (s *userFriendRequestService) DeleteFriendRequests(ctx context.Context, ids []int64) error {
	return s.repo.DeleteByIds(ctx, ids)
}

func (s *userFriendRequestService) QueryFriendRequests(ctx context.Context, ids, requesterIds, recipientIds []int64, statuses []po.RequestStatus, creationDateStart, creationDateEnd, responseDateStart, responseDateEnd, expirationDateStart, expirationDateEnd *time.Time, page, size *int) ([]po.UserFriendRequest, error) {
	return s.repo.FindFriendRequests(ctx, ids, requesterIds, recipientIds, statuses, creationDateStart, creationDateEnd, responseDateStart, responseDateEnd, expirationDateStart, expirationDateEnd, page, size)
}

func (s *userFriendRequestService) CountFriendRequests(ctx context.Context, ids, requesterIds, recipientIds []int64, statuses []po.RequestStatus, creationDateStart, creationDateEnd, responseDateStart, responseDateEnd, expirationDateStart, expirationDateEnd *time.Time) (int64, error) {
	return s.repo.CountFriendRequests(ctx, ids, requesterIds, recipientIds, statuses, creationDateStart, creationDateEnd, responseDateStart, responseDateEnd, expirationDateStart, expirationDateEnd)
}
