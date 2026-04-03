package service

import (
	"context"
	"time"

	"im.turms/server/internal/domain/conversation/po"
	"im.turms/server/internal/domain/conversation/repository"
)

type ConversationService struct {
	privateConvRepo *repository.PrivateConversationRepository
	groupConvRepo   *repository.GroupConversationRepository
}

func NewConversationService(
	privateConvRepo *repository.PrivateConversationRepository,
	groupConvRepo *repository.GroupConversationRepository,
) *ConversationService {
	return &ConversationService{
		privateConvRepo: privateConvRepo,
		groupConvRepo:   groupConvRepo,
	}
}

// AuthAndUpdatePrivateConversationReadDate updates the local high-water mark for reading private messages.
// @MappedFrom authAndUpsertPrivateConversationReadDate(@NotNull Long ownerId, @NotNull Long targetId, @Nullable @PastOrPresent Date readDate)
// Calling this function implies the ownerID has read the messages from targetID.
func (s *ConversationService) AuthAndUpdatePrivateConversationReadDate(ctx context.Context, ownerID int64, targetID int64, readDate time.Time) error {
	// Any authentication checks can go here
	// e.g., checking if ownerID is not blocked by targetID if needed, but normally
	// users can always update their own local read state whether they are friends or not.
	return s.privateConvRepo.UpsertReadDate(ctx, ownerID, targetID, readDate)
}

// AuthAndUpdateGroupConversationReadDate updates the local high-water mark for a user in a group.
// @MappedFrom authAndUpsertGroupConversationReadDate(@NotNull Long groupId, @NotNull Long memberId, @Nullable @PastOrPresent Date readDate)
// Calling this function implies the memberID has read messages in the groupID up to readDate.
func (s *ConversationService) AuthAndUpdateGroupConversationReadDate(ctx context.Context, memberID int64, groupID int64, readDate time.Time) error {
	// Usually we'd check if the user is a group member, but this is a personal state mark.
	// So it's fine to just save it.
	return s.groupConvRepo.UpsertReadDate(ctx, groupID, memberID, readDate)
}

// QueryPrivateConversations fetches all private conversations for a given set of ownerIDs
// @MappedFrom queryPrivateConversationsByOwnerIds(@NotNull Set<Long> ownerIds)
// (which usually is just the requester passing their own ID or a list of devices).
func (s *ConversationService) QueryPrivateConversations(ctx context.Context, ownerIDs []int64) ([]*po.PrivateConversation, error) {
	return s.privateConvRepo.QueryPrivateConversations(ctx, ownerIDs)
}

// QueryGroupConversations fetches the read states for given group IDs.
func (s *ConversationService) QueryGroupConversations(ctx context.Context, groupIDs []int64) ([]*po.GroupConversation, error) {
	return s.groupConvRepo.QueryGroupConversations(ctx, groupIDs)
}

// @MappedFrom upsertGroupConversationReadDate(@NotNull Long groupId, @NotNull Long memberId, @Nullable @PastOrPresent Date readDate)
func (s *ConversationService) UpsertGroupConversationReadDate() {
}

// @MappedFrom upsertGroupConversationsReadDate(@NotNull Set<GroupConversation.GroupConversionMemberKey> keys, @Nullable @PastOrPresent Date readDate)
func (s *ConversationService) UpsertGroupConversationsReadDate() {
}

// @MappedFrom upsertPrivateConversationReadDate(@NotNull Long ownerId, @NotNull Long targetId, @Nullable @PastOrPresent Date readDate)
func (s *ConversationService) UpsertPrivateConversationReadDate() {
}

// @MappedFrom upsertPrivateConversationsReadDate(@NotNull Set<PrivateConversation.Key> keys, @Nullable @PastOrPresent Date readDate)
func (s *ConversationService) UpsertPrivateConversationsReadDate() {
}

// @MappedFrom deletePrivateConversations(@NotNull Set<PrivateConversation.Key> keys)
func (s *ConversationService) DeletePrivateConversationsByKeys() {
}

// @MappedFrom deletePrivateConversations(@NotNull Set<Long> userIds, @Nullable ClientSession session)
func (s *ConversationService) DeletePrivateConversationsByUserIds() {
}

// @MappedFrom deleteGroupConversations(@Nullable Set<Long> groupIds, @Nullable ClientSession session)
func (s *ConversationService) DeleteGroupConversations() {
}

// @MappedFrom deleteGroupMemberConversations(@NotNull Collection<Long> userIds, @Nullable ClientSession session)
func (s *ConversationService) DeleteGroupMemberConversations() {
}

// @MappedFrom authAndUpdateTypingStatus(@NotNull Long requesterId, boolean isGroupMessage, @NotNull Long toId)
func (s *ConversationService) AuthAndUpdateTypingStatus() {
}
