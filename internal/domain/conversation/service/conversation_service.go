package service

import (
	"context"
	"time"

	"im.turms/server/internal/domain/common/constant"
	"im.turms/server/internal/domain/conversation/po"
	"im.turms/server/internal/domain/conversation/repository"
	grouppo "im.turms/server/internal/domain/group/po"
	"im.turms/server/internal/infra/exception"
	"im.turms/server/internal/infra/property"
)

type UserRelationshipService interface {
	HasRelationshipAndNotBlocked(ctx context.Context, ownerID int64, relatedUserID int64) (bool, error)
}

type GroupService interface {
	QueryGroupTypeIfActiveAndNotDeleted(ctx context.Context, groupID int64) (*grouppo.GroupType, error)
}

type GroupMemberService interface {
	IsGroupMember(ctx context.Context, groupID int64, userID int64) (bool, error)
	FindGroupMemberIDs(ctx context.Context, groupID int64) ([]int64, error)
	QueryUserJoinedGroupIds(ctx context.Context, userID int64) ([]int64, error)
}

type MessageService interface {
	HasPrivateMessage(ctx context.Context, senderID int64, targetID int64) (bool, error)
}

type ConversationService struct {
	privateConvRepo       *repository.PrivateConversationRepository
	groupConvRepo         *repository.GroupConversationRepository
	userRelationshipSvc   UserRelationshipService
	groupSvc              GroupService
	groupMemberSvc        GroupMemberService
	messageSvc            MessageService
	propertiesManager     *property.TurmsPropertiesManager
}

func NewConversationService(
	privateConvRepo *repository.PrivateConversationRepository,
	groupConvRepo *repository.GroupConversationRepository,
	userRelationshipSvc UserRelationshipService,
	groupSvc GroupService,
	groupMemberSvc GroupMemberService,
	messageSvc MessageService,
	propertiesManager *property.TurmsPropertiesManager,
) *ConversationService {
	return &ConversationService{
		privateConvRepo:     privateConvRepo,
		groupConvRepo:       groupConvRepo,
		userRelationshipSvc: userRelationshipSvc,
		groupSvc:            groupSvc,
		groupMemberSvc:      groupMemberSvc,
		messageSvc:          messageSvc,
		propertiesManager:   propertiesManager,
	}
}

// AuthAndUpsertGroupConversationReadDate updates the local high-water mark for a user in a group.
// @MappedFrom authAndUpsertGroupConversationReadDate(@NotNull Long groupId, @NotNull Long memberId, @Nullable @PastOrPresent Date readDate)
func (s *ConversationService) AuthAndUpsertGroupConversationReadDate(ctx context.Context, groupID int64, memberID int64, readDate *time.Time) error {
	props := s.propertiesManager.GetLocalProperties().Service.Conversation.ReadReceipt
	if !props.Enabled {
		return exception.NewTurmsError(int32(constant.ResponseStatusCode_UPDATING_READ_DATE_IS_DISABLED), "")
	}

	groupType, err := s.groupSvc.QueryGroupTypeIfActiveAndNotDeleted(ctx, groupID)
	if err != nil {
		return err
	}
	if groupType == nil {
		return exception.NewTurmsError(int32(constant.ResponseStatusCode_UPDATING_READ_DATE_OF_NONEXISTENT_GROUP_CONVERSATION), "")
	}

	isMember, err := s.groupMemberSvc.IsGroupMember(ctx, groupID, memberID)
	if err != nil {
		return err
	}
	if !isMember {
		return exception.NewTurmsError(int32(constant.ResponseStatusCode_NOT_GROUP_MEMBER_TO_UPDATE_READ_DATE_OF_GROUP_CONVERSATION), "")
	}

	if !groupType.EnableReadReceipt {
		return exception.NewTurmsError(int32(constant.ResponseStatusCode_UPDATING_READ_DATE_IS_DISABLED_BY_GROUP), "")
	}

	finalReadDate := time.Now()
	if !props.UseServerTime && readDate != nil {
		finalReadDate = *readDate
	}

	return s.UpsertGroupConversationReadDate(ctx, groupID, memberID, finalReadDate)
}

// AuthAndUpsertPrivateConversationReadDate updates the local high-water mark for reading private messages.
// @MappedFrom authAndUpsertPrivateConversationReadDate(@NotNull Long ownerId, @NotNull Long targetId, @Nullable @PastOrPresent Date readDate)
func (s *ConversationService) AuthAndUpsertPrivateConversationReadDate(ctx context.Context, ownerID int64, targetID int64, readDate *time.Time) error {
	props := s.propertiesManager.GetLocalProperties().Service.Conversation.ReadReceipt
	if !props.Enabled {
		return exception.NewTurmsError(int32(constant.ResponseStatusCode_UPDATING_READ_DATE_IS_DISABLED), "")
	}

	hasMessage, err := s.messageSvc.HasPrivateMessage(ctx, targetID, ownerID)
	if err != nil {
		return err
	}
	if !hasMessage {
		return nil
	}

	finalReadDate := time.Now()
	if !props.UseServerTime && readDate != nil {
		finalReadDate = *readDate
	}

	return s.UpsertPrivateConversationReadDate(ctx, ownerID, targetID, finalReadDate)
}

// UpsertGroupConversationReadDate
// @MappedFrom upsertGroupConversationReadDate(@NotNull Long groupId, @NotNull Long memberId, @Nullable @PastOrPresent Date readDate)
func (s *ConversationService) UpsertGroupConversationReadDate(ctx context.Context, groupID int64, memberID int64, readDate time.Time) error {
	allowMoveForward := s.propertiesManager.GetLocalProperties().Service.Conversation.ReadReceipt.AllowMoveReadDateForward
	err := s.groupConvRepo.Upsert(ctx, groupID, memberID, readDate, allowMoveForward)
	if err != nil && exception.IsDuplicateKeyError(err) {
		return exception.NewTurmsError(int32(constant.ResponseStatusCode_MOVING_READ_DATE_FORWARD_IS_DISABLED), "")
	}
	return err
}

// UpsertGroupConversationsReadDate
// @MappedFrom upsertGroupConversationsReadDate(@NotNull Set<GroupConversation.GroupConversionMemberKey> keys, @Nullable @PastOrPresent Date readDate)
func (s *ConversationService) UpsertGroupConversationsReadDate(ctx context.Context, groupIDToMemberIDs map[int64][]int64, readDate *time.Time) error {
	if len(groupIDToMemberIDs) == 0 {
		return nil
	}
	finalReadDate := time.Now()
	if readDate != nil {
		finalReadDate = *readDate
	}

	for groupID, memberIDs := range groupIDToMemberIDs {
		err := s.groupConvRepo.BulkUpsert(ctx, groupID, memberIDs, finalReadDate)
		if err != nil {
			return err
		}
	}
	return nil
}

// UpsertPrivateConversationReadDate
// @MappedFrom upsertPrivateConversationReadDate(@NotNull Long ownerId, @NotNull Long targetId, @Nullable @PastOrPresent Date readDate)
func (s *ConversationService) UpsertPrivateConversationReadDate(ctx context.Context, ownerID int64, targetID int64, readDate time.Time) error {
	return s.UpsertPrivateConversationsReadDate(ctx, []po.PrivateConversationKey{{OwnerID: ownerID, TargetID: targetID}}, &readDate)
}

// UpsertPrivateConversationsReadDate
// @MappedFrom upsertPrivateConversationsReadDate(@NotNull Set<PrivateConversation.Key> keys, @Nullable @PastOrPresent Date readDate)
func (s *ConversationService) UpsertPrivateConversationsReadDate(ctx context.Context, keys []po.PrivateConversationKey, readDate *time.Time) error {
	if len(keys) == 0 {
		return nil
	}
	finalReadDate := time.Now()
	if readDate != nil {
		finalReadDate = *readDate
	}
	allowMoveForward := s.propertiesManager.GetLocalProperties().Service.Conversation.ReadReceipt.AllowMoveReadDateForward
	err := s.privateConvRepo.Upsert(ctx, keys, finalReadDate, allowMoveForward)
	if err != nil && exception.IsDuplicateKeyError(err) {
		return exception.NewTurmsError(int32(constant.ResponseStatusCode_MOVING_READ_DATE_FORWARD_IS_DISABLED), "")
	}
	return err
}

// QueryGroupConversations fetches the read states for given group IDs.
// @MappedFrom queryGroupConversations(@NotNull Collection<Long> groupIds)
func (s *ConversationService) QueryGroupConversations(ctx context.Context, groupIDs []int64) ([]*po.GroupConversation, error) {
	if len(groupIDs) == 0 {
		return nil, nil
	}
	return s.groupConvRepo.FindByIds(ctx, groupIDs)
}

// QueryPrivateConversationsByOwnerIds fetches all private conversations for a given set of ownerIDs
// @MappedFrom queryPrivateConversationsByOwnerIds(@NotNull Set<Long> ownerIds)
func (s *ConversationService) QueryPrivateConversationsByOwnerIds(ctx context.Context, ownerIDs []int64) ([]*po.PrivateConversation, error) {
	if len(ownerIDs) == 0 {
		return nil, nil
	}
	return s.privateConvRepo.FindConversations(ctx, ownerIDs)
}

// QueryPrivateConversations
// @MappedFrom queryPrivateConversations(@NotNull Set<PrivateConversation.Key> keys)
func (s *ConversationService) QueryPrivateConversations(ctx context.Context, keys []po.PrivateConversationKey) ([]*po.PrivateConversation, error) {
	if len(keys) == 0 {
		return nil, nil
	}
	return s.privateConvRepo.FindByIds(ctx, keys)
}

// DeletePrivateConversationsByKeys
// @MappedFrom deletePrivateConversations(@NotNull Set<PrivateConversation.Key> keys)
func (s *ConversationService) DeletePrivateConversationsByKeys(ctx context.Context, keys []po.PrivateConversationKey) error {
	if len(keys) == 0 {
		return nil
	}
	return s.privateConvRepo.DeleteByIds(ctx, keys)
}

// DeletePrivateConversationsByUserIds
// @MappedFrom deletePrivateConversations(@NotNull Set<Long> userIds, @Nullable ClientSession session)
func (s *ConversationService) DeletePrivateConversationsByUserIds(ctx context.Context, userIDs []int64) error {
	if len(userIDs) == 0 {
		return nil
	}
	return s.privateConvRepo.DeleteConversationsByOwnerIds(ctx, userIDs)
}

// DeleteGroupConversations
// @MappedFrom deleteGroupConversations(@Nullable Set<Long> groupIds, @Nullable ClientSession session)
func (s *ConversationService) DeleteGroupConversations(ctx context.Context, groupIDs []int64) error {
	if len(groupIDs) == 0 {
		return nil
	}
	return s.groupConvRepo.DeleteByIds(ctx, groupIDs)
}

// DeleteGroupMemberConversations
// @MappedFrom deleteGroupMemberConversations(@NotNull Collection<Long> userIds, @Nullable ClientSession session)
func (s *ConversationService) DeleteGroupMemberConversations(ctx context.Context, userIDs []int64) error {
	if len(userIDs) == 0 {
		return nil
	}
	for _, userID := range userIDs {
		groupIDs, err := s.groupMemberSvc.QueryUserJoinedGroupIds(ctx, userID)
		if err != nil {
			return err
		}
		if len(groupIDs) > 0 {
			err = s.groupConvRepo.DeleteMemberConversations(ctx, groupIDs, userID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// AuthAndUpdateTypingStatus
// @MappedFrom authAndUpdateTypingStatus(@NotNull Long requesterId, boolean isGroupMessage, @NotNull Long toId)
func (s *ConversationService) AuthAndUpdateTypingStatus(ctx context.Context, requesterID int64, isGroupMessage bool, toID int64) ([]int64, error) {
	if !s.propertiesManager.GetLocalProperties().Service.Conversation.TypingStatus.Enabled {
		return nil, exception.NewTurmsError(int32(constant.ResponseStatusCode_UPDATING_TYPING_STATUS_IS_DISABLED), "")
	}

	if isGroupMessage {
		isMember, err := s.groupMemberSvc.IsGroupMember(ctx, toID, requesterID)
		if err != nil {
			return nil, err
		}
		if !isMember {
			return nil, exception.NewTurmsError(int32(constant.ResponseStatusCode_NOT_GROUP_MEMBER_TO_SEND_TYPING_STATUS), "")
		}
		return s.groupMemberSvc.FindGroupMemberIDs(ctx, toID)
	} else {
		canSend, err := s.userRelationshipSvc.HasRelationshipAndNotBlocked(ctx, toID, requesterID)
		if err != nil {
			return nil, err
		}
		if !canSend {
			return nil, exception.NewTurmsError(int32(constant.ResponseStatusCode_NOT_FRIEND_TO_SEND_TYPING_STATUS), "")
		}
		return []int64{toID}, nil
	}
}
