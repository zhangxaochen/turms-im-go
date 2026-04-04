package service

import (
	"context"
	"time"

	"im.turms/server/internal/domain/conversation/po"
	"im.turms/server/internal/domain/conversation/repository"
	group_service "im.turms/server/internal/domain/group/service"
	user_service "im.turms/server/internal/domain/user/service"
	"im.turms/server/internal/infra/exception"
	"im.turms/server/pkg/codes"
)

// ConversationSettingsService maps to ConversationSettingsService.java
// @MappedFrom ConversationSettingsService
type ConversationSettingsService struct {
	groupMemberService             *group_service.GroupMemberService
	userRelationshipService        user_service.UserRelationshipService
	conversationSettingsRepository *repository.ConversationSettingsRepository
}

func NewConversationSettingsService(
	groupMemberService *group_service.GroupMemberService,
	userRelationshipService user_service.UserRelationshipService,
	conversationSettingsRepository *repository.ConversationSettingsRepository,
) *ConversationSettingsService {
	return &ConversationSettingsService{
		groupMemberService:             groupMemberService,
		userRelationshipService:        userRelationshipService,
		conversationSettingsRepository: conversationSettingsRepository,
	}
}

// @MappedFrom upsertPrivateConversationSettings(Long ownerId, Long userId, Map<String, Value> settings)
func (s *ConversationSettingsService) UpsertPrivateConversationSettings(ctx context.Context, ownerId int64, userId int64, settings map[string]any) error {
	if len(settings) == 0 {
		return nil
	}
	related, err := s.userRelationshipService.HasRelationshipAndNotBlocked(ctx, ownerId, userId)
	if err != nil {
		return err
	}
	if !related {
		return exception.NewTurmsError(int32(codes.NotRelatedUserToUpdatePrivateConversationSetting), "not related user to update private conversation setting")
	}

	return s.conversationSettingsRepository.UpsertSettings(ctx, ownerId, userId, settings, time.Now())
}

// @MappedFrom upsertGroupConversationSettings(Long ownerId, Long groupId, Map<String, Value> settings)
func (s *ConversationSettingsService) UpsertGroupConversationSettings(ctx context.Context, ownerId int64, groupId int64, settings map[string]any) error {
	if len(settings) == 0 {
		return nil
	}
	isMember, err := s.groupMemberService.IsGroupMember(ctx, groupId, ownerId)
	if err != nil {
		return err
	}
	if !isMember {
		return exception.NewTurmsError(int32(codes.NotGroupMemberToUpdateGroupConversationSetting), "not group member to update group conversation setting")
	}

	return s.conversationSettingsRepository.UpsertSettings(ctx, ownerId, s.getTargetIdFromGroupId(groupId), settings, time.Now())
}

// @MappedFrom deleteSettings(Collection<Long> ownerIds, @Nullable ClientSession clientSession)
func (s *ConversationSettingsService) DeleteSettings(ctx context.Context, ownerIds []int64) (bool, error) {
	count, err := s.conversationSettingsRepository.DeleteByOwnerIds(ctx, ownerIds)
	return count > 0, err
}

// @MappedFrom unsetSettings(Long ownerId, @Nullable Set<Long> userIds, @Nullable Set<Long> groupIds, @Nullable Set<String> settingNames)
func (s *ConversationSettingsService) UnsetSettings(ctx context.Context, ownerId int64, userIds []int64, groupIds []int64, settingNames []string) error {
	targetIds := s.getTargetIds(userIds, groupIds)
	return s.conversationSettingsRepository.UnsetSettings(ctx, ownerId, targetIds, settingNames)
}

// @MappedFrom querySettings(Long ownerId, @Nullable Collection<Long> userIds, @Nullable Collection<Long> groupIds, @Nullable Set<String> settingNames, @Nullable Date lastUpdatedDateStart)
func (s *ConversationSettingsService) QuerySettings(ctx context.Context, ownerId int64, userIds []int64, groupIds []int64, settingNames []string, lastUpdatedDateStart *time.Time) ([]po.ConversationSettings, error) {
	var keys []po.ConversationSettingsKey
	if len(userIds) > 0 || len(groupIds) > 0 {
		keys = make([]po.ConversationSettingsKey, 0, len(userIds)+len(groupIds))
		for _, uid := range userIds {
			keys = append(keys, po.ConversationSettingsKey{OwnerId: ownerId, TargetId: uid})
		}
		for _, gid := range groupIds {
			keys = append(keys, po.ConversationSettingsKey{OwnerId: ownerId, TargetId: s.getTargetIdFromGroupId(gid)})
		}
		return s.conversationSettingsRepository.FindByIdAndSettingNamesWithKeys(ctx, keys, settingNames, lastUpdatedDateStart)
	}

	settings, err := s.conversationSettingsRepository.FindByKey(ctx, ownerId, 0, settingNames, lastUpdatedDateStart)
	if err != nil {
		return nil, err
	}
	if settings == nil {
		return nil, nil
	}
	// Note: the Java code returns multiple if ownerId or targetIds are collections.
	// Here we return a slice for consistency.
	return []po.ConversationSettings{*settings}, nil
}

func (s *ConversationSettingsService) getTargetIds(userIds []int64, groupIds []int64) []int64 {
	count := len(userIds) + len(groupIds)
	if count == 0 {
		return nil
	}
	targetIds := make([]int64, 0, count)
	targetIds = append(targetIds, userIds...)
	for _, gid := range groupIds {
		targetIds = append(targetIds, s.getTargetIdFromGroupId(gid))
	}
	return targetIds
}

func (s *ConversationSettingsService) getTargetIdFromGroupId(groupId int64) int64 {
	return -groupId
}
