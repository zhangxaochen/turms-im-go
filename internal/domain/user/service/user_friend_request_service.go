package service

import (
	"context"
	"fmt"
	"time"

	"im.turms/server/internal/domain/common/infra/idgen"
	common "im.turms/server/internal/domain/common/service"
	"im.turms/server/internal/domain/user/po"
	"im.turms/server/internal/domain/user/repository"
	"im.turms/server/internal/infra/exception"
	"im.turms/server/internal/infra/validator"
	"im.turms/server/pkg/codes"
	"im.turms/server/pkg/protocol"
)

type UserFriendRequestService interface {
	RemoveAllExpiredFriendRequests(ctx context.Context, expirationDate time.Time) error
	HasPendingFriendRequest(ctx context.Context, requesterID, recipientID int64) (bool, error)
	CreateFriendRequest(ctx context.Context, requestID *int64, requesterID, recipientID int64, content string, status *po.RequestStatus, creationDate, responseDate *time.Time, reason *string) (*po.UserFriendRequest, error)
	AuthAndCreateFriendRequest(ctx context.Context, requesterID, recipientID int64, content string, creationDate time.Time) (*po.UserFriendRequest, error)
	AuthAndRecallFriendRequest(ctx context.Context, requesterID, requestID int64) (*po.UserFriendRequest, error)
	UpdatePendingFriendRequestStatus(ctx context.Context, requestID int64, targetStatus po.RequestStatus, reason *string) (bool, error)
	UpdateFriendRequests(ctx context.Context, requestIds []int64, requesterID, recipientID *int64, content *string, status *po.RequestStatus, reason *string, creationDate *time.Time, responseDate *time.Time) error
	QueryRecipientId(ctx context.Context, requestID int64) (int64, error)
	AuthAndHandleFriendRequest(ctx context.Context, friendRequestID int64, requesterID int64, action po.ResponseAction, reason *string) (bool, error)
	QueryFriendRequestsByRecipientId(ctx context.Context, recipientID int64) ([]po.UserFriendRequest, error)
	QueryFriendRequestsByRequesterId(ctx context.Context, requesterID int64) ([]po.UserFriendRequest, error)
	QueryFriendRequestsWithVersion(ctx context.Context, userID int64, isRecipient bool, lastUpdatedDate *time.Time) ([]po.UserFriendRequest, error)
	DeleteFriendRequests(ctx context.Context, ids []int64) error
	QueryFriendRequests(ctx context.Context, ids, requesterIds, recipientIds []int64, statuses []po.RequestStatus, creationDateStart, creationDateEnd, responseDateStart, responseDateEnd, expirationDateStart, expirationDateEnd *time.Time, page, size *int) ([]po.UserFriendRequest, error)
	CountFriendRequests(ctx context.Context, ids, requesterIds, recipientIds []int64, statuses []po.RequestStatus, creationDateStart, creationDateEnd, responseDateStart, responseDateEnd, expirationDateStart, expirationDateEnd *time.Time) (int64, error)
}

type userFriendRequestService struct {
	idGen                  *idgen.SnowflakeIdGenerator
	repo                   repository.UserFriendRequestRepository
	relationshipService    UserRelationshipService
	userVersionService     *UserVersionService
	outboundMessageService common.OutboundMessageService
}

func NewUserFriendRequestService(
	idGen *idgen.SnowflakeIdGenerator,
	repo repository.UserFriendRequestRepository,
	relService UserRelationshipService,
	userVersionService *UserVersionService,
	outboundMessageService common.OutboundMessageService,
) UserFriendRequestService {
	return &userFriendRequestService{
		idGen:                  idGen,
		repo:                   repo,
		relationshipService:    relService,
		userVersionService:     userVersionService,
		outboundMessageService: outboundMessageService,
	}
}

const (
	defaultMaxContentLength                        = 200
	defaultMaxResponseReasonLength                 = 200
	defaultAllowSendRequestAfterDeclinedOrIgnored  = true
	defaultAllowRecallPendingFriendRequestBySender = true
)

// @MappedFrom removeAllExpiredFriendRequests(Date expirationDate)
func (s *userFriendRequestService) RemoveAllExpiredFriendRequests(ctx context.Context, expirationDate time.Time) error {
	return s.repo.DeleteExpiredData(ctx, expirationDate)
}

func (s *userFriendRequestService) HasPendingFriendRequest(ctx context.Context, requesterID, recipientID int64) (bool, error) {
	return s.repo.HasPendingFriendRequest(ctx, requesterID, recipientID)
}

// @MappedFrom createFriendRequest(@RequestBody AddFriendRequestDTO addFriendRequestDTO)
// @MappedFrom createFriendRequest(@Nullable Long id, @NotNull Long requesterId, @NotNull Long recipientId, @NotNull String content, @Nullable @ValidRequestStatus RequestStatus status, @Nullable @PastOrPresent Date creationDate, @Nullable @PastOrPresent Date responseDate, @Nullable String reason)
func (s *userFriendRequestService) CreateFriendRequest(ctx context.Context, requestID *int64, requesterID, recipientID int64, content string, status *po.RequestStatus, creationDate, responseDate *time.Time, reason *string) (*po.UserFriendRequest, error) {
	if err := validator.NotNull(requesterID, "requesterID"); err != nil {
		return nil, err
	}
	if err := validator.NotNull(recipientID, "recipientID"); err != nil {
		return nil, err
	}
	if err := validator.MaxLength(&content, "content", defaultMaxContentLength); err != nil {
		return nil, err
	}
	if err := validator.NotEquals(requesterID, recipientID, "The requester ID must not be equal to the recipient ID"); err != nil {
		return nil, err
	}
	if err := validator.MaxLength(reason, "reason", defaultMaxResponseReasonLength); err != nil {
		return nil, err
	}

	id := int64(0)
	if requestID == nil {
		id = s.idGen.NextIncreasingId()
	} else {
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

	if responseDate == nil {
		if st != po.RequestStatusPending {
			responseDate = &now
		}
	}

	req := &po.UserFriendRequest{
		ID:           id,
		Content:      content,
		Status:       st,
		Reason:       reason,
		CreationDate: cd,
		ResponseDate: responseDate,
		RequesterID:  requesterID,
		RecipientID:  recipientID,
	}

	if err := s.repo.Insert(ctx, req); err != nil {
		return nil, err
	}

	// Update versions asynchronously
	go func() {
		bgCtx := context.Background()
		_ = s.userVersionService.UpdateReceivedFriendRequestsVersion(bgCtx, recipientID)
		_ = s.userVersionService.UpdateSentFriendRequestsVersion(bgCtx, requesterID)
	}()

	return req, nil
}

// @MappedFrom authAndCreateFriendRequest(@NotNull Long requesterId, @NotNull Long recipientId, @Nullable String content, @NotNull @PastOrPresent Date creationDate)
func (s *userFriendRequestService) AuthAndCreateFriendRequest(ctx context.Context, requesterID, recipientID int64, content string, creationDate time.Time) (*po.UserFriendRequest, error) {
	if err := validator.NotNull(requesterID, "requesterID"); err != nil {
		return nil, err
	}
	if err := validator.NotNull(recipientID, "recipientID"); err != nil {
		return nil, err
	}
	if err := validator.MaxLength(&content, "content", defaultMaxContentLength); err != nil {
		return nil, err
	}
	if err := validator.NotEquals(requesterID, recipientID, "The requester ID must not be equal to the recipient ID"); err != nil {
		return nil, err
	}

	isNotBlocked, err := s.relationshipService.IsNotBlocked(ctx, recipientID, requesterID)
	if err != nil {
		return nil, err
	}
	if !isNotBlocked {
		return nil, exception.NewTurmsError(int32(codes.BlockedUserToSendFriendRequest), "")
	}

	var requestExists bool
	if defaultAllowSendRequestAfterDeclinedOrIgnored {
		requestExists, err = s.HasPendingFriendRequest(ctx, requesterID, recipientID)
	} else {
		requestExists, err = s.repo.HasPendingOrDeclinedOrIgnoredOrExpiredRequest(ctx, requesterID, recipientID)
	}
	if err != nil {
		return nil, err
	}
	if requestExists {
		return nil, exception.NewTurmsError(int32(codes.CreateExistingFriendRequest), "")
	}

	req, err := s.CreateFriendRequest(ctx, nil, requesterID, recipientID, content, nil, &creationDate, nil, nil)
	if err != nil {
		return nil, err
	}

	// Send notification to recipient
	notification := &protocol.TurmsNotification{
		Data: &protocol.TurmsNotification_Data{
			Kind: &protocol.TurmsNotification_Data_UserFriendRequestsWithVersion{
				UserFriendRequestsWithVersion: &protocol.UserFriendRequestsWithVersion{
					UserFriendRequests: []*protocol.UserFriendRequest{
						FriendRequestToProto(req),
					},
				},
			},
		},
	}
	s.outboundMessageService.ForwardNotificationToMultiple(ctx, notification, []int64{recipientID})

	return req, nil
}

// @MappedFrom authAndRecallFriendRequest(@NotNull Long requesterId, @NotNull Long requestId)
func (s *userFriendRequestService) AuthAndRecallFriendRequest(ctx context.Context, requesterID, requestID int64) (*po.UserFriendRequest, error) {
	if err := validator.NotNull(requesterID, "requesterID"); err != nil {
		return nil, err
	}
	if err := validator.NotNull(requestID, "requestID"); err != nil {
		return nil, err
	}

	if !defaultAllowRecallPendingFriendRequestBySender {
		return nil, exception.NewTurmsError(int32(codes.RecallingFriendRequestIsDisabled), "")
	}

	req, err := s.repo.FindRequesterIdAndRecipientIdAndCreationDateAndStatus(ctx, requestID)
	if err != nil {
		return nil, err
	}
	// If the requester is not authorized to the request,
	// they should not know the status of the request from the error code.
	if req == nil || req.RequesterID != requesterID {
		return nil, exception.NewTurmsError(int32(codes.NotSenderToRecallFriendRequest), "")
	}
	if req.Status != po.RequestStatusPending {
		return nil, exception.NewTurmsError(int32(codes.RecallNonPendingFriendRequest), "")
	}

	success, err := s.repo.UpdateStatusIfPending(ctx, requestID, po.RequestStatusCanceled, nil, time.Now())
	if err != nil {
		return nil, err
	}
	if !success {
		return nil, exception.NewTurmsError(int32(codes.RecallNonPendingFriendRequest), "")
	}

	_ = s.userVersionService.UpdateReceivedFriendRequestsVersion(ctx, req.RecipientID)
	_ = s.userVersionService.UpdateSentFriendRequestsVersion(ctx, req.RequesterID)

	req.Status = po.RequestStatusCanceled
	return req, nil
}

func (s *userFriendRequestService) UpdateFriendRequests(ctx context.Context, requestIds []int64, requesterID, recipientID *int64, content *string, status *po.RequestStatus, reason *string, creationDate *time.Time, responseDate *time.Time) error {
	if err := validator.NotEmpty(requestIds, "requestIds"); err != nil {
		return err
	}
	if err := validator.MaxLength(content, "content", defaultMaxContentLength); err != nil {
		return err
	}
	if err := validator.PastOrPresent(creationDate, "creationDate"); err != nil {
		return err
	}
	if err := validator.PastOrPresent(responseDate, "responseDate"); err != nil {
		return err
	}
	if err := validator.MaxLength(reason, "reason", defaultMaxResponseReasonLength); err != nil {
		return err
	}
	if requesterID != nil && recipientID != nil && *requesterID == *recipientID {
		return exception.NewTurmsError(int32(codes.IllegalArgument), "The requester ID must not equal the recipient ID")
	}

	if validator.AreAllNull(requesterID, recipientID, content, status, reason, creationDate, responseDate) {
		return nil
	}

	return s.repo.UpdateFriendRequests(ctx, requestIds, requesterID, recipientID, content, status, reason, creationDate)
}

// @MappedFrom queryRecipientId(@NotNull Long requestId)
func (s *userFriendRequestService) QueryRecipientId(ctx context.Context, requestID int64) (int64, error) {
	return s.repo.FindRecipientId(ctx, requestID)
}

// @MappedFrom updatePendingFriendRequestStatus(@NotNull Long requestId, @NotNull @ValidRequestStatus RequestStatus requestStatus, @Nullable String reason, @Nullable ClientSession session)
func (s *userFriendRequestService) UpdatePendingFriendRequestStatus(ctx context.Context, requestID int64, targetStatus po.RequestStatus, reason *string) (bool, error) {
	if err := validator.NotNull(requestID, "requestID"); err != nil {
		return false, err
	}
	if err := validator.NotNull(targetStatus, "targetStatus"); err != nil {
		return false, err
	}
	if targetStatus == po.RequestStatusPending {
		return false, exception.NewTurmsError(int32(codes.IllegalArgument), "The request status must not be PENDING")
	}
	if err := validator.MaxLength(reason, "reason", defaultMaxResponseReasonLength); err != nil {
		return false, err
	}

	success, err := s.repo.UpdateStatusIfPending(ctx, requestID, targetStatus, reason, time.Now())
	if err != nil {
		return false, err
	}
	if success {
		recipientID, err := s.repo.FindRecipientId(ctx, requestID)
		if err == nil && recipientID != 0 {
			_ = s.userVersionService.UpdateReceivedFriendRequestsVersion(ctx, recipientID)
		}
	}
	return success, nil
}

// @MappedFrom authAndHandleFriendRequest(@NotNull Long friendRequestId, @NotNull Long requesterId, @NotNull @ValidResponseAction ResponseAction action, @Nullable String reason)
func (s *userFriendRequestService) AuthAndHandleFriendRequest(ctx context.Context, friendRequestID int64, requesterID int64, action po.ResponseAction, reason *string) (bool, error) {
	if friendRequestID <= 0 {
		return false, exception.NewTurmsError(int32(codes.IllegalArgument), "friendRequestID must be greater than 0")
	}
	if requesterID <= 0 { // In original Java, this is requesterId from session
		return false, exception.NewTurmsError(int32(codes.IllegalArgument), "requesterID must be greater than 0")
	}

	req, err := s.repo.FindRequesterIdAndRecipientIdAndCreationDateAndStatus(ctx, friendRequestID)
	if err != nil {
		return false, err
	}
	if req == nil || req.RecipientID != requesterID {
		return false, exception.NewTurmsError(int32(codes.NotRecipientToUpdateFriendRequest), "")
	}

	if req.Status != po.RequestStatusPending {
		return false, exception.NewTurmsError(int32(codes.UpdateNonPendingFriendRequest), "")
	}

	var status po.RequestStatus
	switch action {
	case po.ResponseActionAccept:
		status = po.RequestStatusAccepted
	case po.ResponseActionIgnore:
		status = po.RequestStatusIgnored
	case po.ResponseActionDecline:
		status = po.RequestStatusDeclined
	default:
		return false, exception.NewTurmsError(int32(codes.IllegalArgument), fmt.Sprintf("Illegal response action: %v", action))
	}

	success, err := s.UpdatePendingFriendRequestStatus(ctx, friendRequestID, status, reason)
	if err != nil {
		return false, err
	}
	if success && status == po.RequestStatusAccepted {
		err = s.relationshipService.FriendTwoUsers(ctx, req.RequesterID, requesterID)
	}

	if success {
		// Send notification to requester
		notification := &protocol.TurmsNotification{
			Data: &protocol.TurmsNotification_Data{
				Kind: &protocol.TurmsNotification_Data_UserFriendRequestsWithVersion{
					UserFriendRequestsWithVersion: &protocol.UserFriendRequestsWithVersion{
						UserFriendRequests: []*protocol.UserFriendRequest{
							FriendRequestToProto(req),
						},
					},
				},
			},
		}
		s.outboundMessageService.ForwardNotificationToMultiple(ctx, notification, []int64{req.RequesterID})
	}

	return success, err
}

// @MappedFrom queryFriendRequestsByRecipientId(@NotNull Long recipientId)
func (s *userFriendRequestService) QueryFriendRequestsByRecipientId(ctx context.Context, recipientID int64) ([]po.UserFriendRequest, error) {
	return s.repo.FindFriendRequestsByRecipientId(ctx, recipientID)
}

// @MappedFrom queryFriendRequestsByRequesterId(@NotNull Long requesterId)
func (s *userFriendRequestService) QueryFriendRequestsByRequesterId(ctx context.Context, requesterID int64) ([]po.UserFriendRequest, error) {
	return s.repo.FindFriendRequestsByRequesterId(ctx, requesterID)
}

// @MappedFrom queryFriendRequestsWithVersion(@NotNull Long userId, boolean areSentByUser, @Nullable Date lastUpdatedDate)
func (s *userFriendRequestService) QueryFriendRequestsWithVersion(ctx context.Context, userID int64, isRecipient bool, lastUpdatedDate *time.Time) ([]po.UserFriendRequest, error) {
	if isRecipient {
		return s.QueryFriendRequestsByRecipientId(ctx, userID)
	}
	return s.QueryFriendRequestsByRequesterId(ctx, userID)
}

// @MappedFrom deleteFriendRequests(@Nullable Set<Long> ids)
// @MappedFrom deleteFriendRequests(@QueryParam(required = false)
func (s *userFriendRequestService) DeleteFriendRequests(ctx context.Context, ids []int64) error {
	return s.repo.DeleteByIds(ctx, ids)
}

// @MappedFrom queryFriendRequests(@Nullable Set<Long> ids, @Nullable Set<Long> requesterIds, @Nullable Set<Long> recipientIds, @Nullable Set<RequestStatus> statuses, @Nullable DateRange creationDateRange, @Nullable DateRange responseDateRange, @Nullable DateRange expirationDateRange, @Nullable Integer page, @Nullable Integer size)
// @MappedFrom queryFriendRequests(@QueryParam(required = false)
func (s *userFriendRequestService) QueryFriendRequests(ctx context.Context, ids, requesterIds, recipientIds []int64, statuses []po.RequestStatus, creationDateStart, creationDateEnd, responseDateStart, responseDateEnd, expirationDateStart, expirationDateEnd *time.Time, page, size *int) ([]po.UserFriendRequest, error) {
	if err := validator.PastOrPresent(creationDateStart, "creationDateStart"); err != nil {
		return nil, err
	}
	if err := validator.PastOrPresent(creationDateEnd, "creationDateEnd"); err != nil {
		return nil, err
	}
	if err := validator.PastOrPresent(responseDateStart, "responseDateStart"); err != nil {
		return nil, err
	}
	if err := validator.PastOrPresent(responseDateEnd, "responseDateEnd"); err != nil {
		return nil, err
	}
	return s.repo.FindFriendRequests(ctx, ids, requesterIds, recipientIds, statuses, creationDateStart, creationDateEnd, responseDateStart, responseDateEnd, expirationDateStart, expirationDateEnd, page, size)
}

func (s *userFriendRequestService) CountFriendRequests(ctx context.Context, ids, requesterIds, recipientIds []int64, statuses []po.RequestStatus, creationDateStart, creationDateEnd, responseDateStart, responseDateEnd, expirationDateStart, expirationDateEnd *time.Time) (int64, error) {
	return s.repo.CountFriendRequests(ctx, ids, requesterIds, recipientIds, statuses, creationDateStart, creationDateEnd, responseDateStart, responseDateEnd, expirationDateStart, expirationDateEnd)
}
