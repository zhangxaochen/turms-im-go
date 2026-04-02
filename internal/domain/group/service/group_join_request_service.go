package service

import (
	"context"
	"time"

	"im.turms/server/internal/domain/group/po"
	"im.turms/server/internal/domain/group/repository"
)

type GroupJoinRequestService struct {
	joinReqRepo repository.GroupJoinRequestRepository
}

func NewGroupJoinRequestService(joinReqRepo repository.GroupJoinRequestRepository) *GroupJoinRequestService {
	return &GroupJoinRequestService{
		joinReqRepo: joinReqRepo,
	}
}

func (s *GroupJoinRequestService) CreateJoinRequest(ctx context.Context, groupID int64, requesterID int64, content string) (*po.GroupJoinRequest, error) {
	now := time.Now()
	
	id := now.UnixNano()

	req := &po.GroupJoinRequest{
		ID:           id,
		GroupID:      groupID,
		RequesterID:  requesterID,
		Content:      content,
		Status:       po.RequestStatusPending,
		CreationDate: now,
	}
	err := s.joinReqRepo.Insert(ctx, req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (s *GroupJoinRequestService) RecallPendingJoinRequest(ctx context.Context, requestID int64, requesterID int64) (bool, error) {
	return s.joinReqRepo.UpdateStatusIfPending(ctx, requestID, requesterID, po.RequestStatusCanceled, nil, time.Now())
}

func (s *GroupJoinRequestService) ReplyToJoinRequest(ctx context.Context, requestID int64, responderID int64, accept bool) (bool, error) {
	status := po.RequestStatusDeclined
	if accept {
		status = po.RequestStatusAccepted
	}
	return s.joinReqRepo.UpdateStatusIfPending(ctx, requestID, responderID, status, nil, time.Now())
}
