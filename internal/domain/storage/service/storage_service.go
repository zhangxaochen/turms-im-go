package service

import (
	"context"
	"errors"
	"time"

	"im.turms/server/internal/domain/storage/bo"
	"im.turms/server/internal/domain/storage/constants"
	"im.turms/server/internal/domain/storage/provider"
)

type StorageService struct {
	provider provider.StorageProvider
}

func NewStorageService(provider provider.StorageProvider) *StorageService {
	return &StorageService{
		provider: provider,
	}
}

func (s *StorageService) DeleteResource(
	ctx context.Context,
	requesterID int64,
	resourceType constants.StorageResourceType,
	resourceIDStr string,
) error {
	if resourceType == 0 {
		return errors.New("unrecognized storage resource type")
	}

	return s.provider.DeleteResource(ctx, resourceType, resourceIDStr)
}

// @MappedFrom queryResourceUploadInfo(Long requesterId, StorageResourceType resourceType, @Nullable Long resourceIdNum, @Nullable String resourceName, @Nullable String resourceMediaType, List<Value> customAttributes)
func (s *StorageService) QueryResourceUploadInfo(
	ctx context.Context,
	requesterID int64,
	resourceType constants.StorageResourceType,
	resourceIDNum *int64,
	resourceName string,
	contentType string,
	maxSize int64,
	resourceIDKey string,
) (string, error) {
	if resourceType == 0 { // Treat 0 as unrecognized or default
		return "", errors.New("unrecognized storage resource type")
	}

	// Use resourceIDKey as keyStr if available, else fall back to resourceName
	keyStr := resourceName
	if resourceIDKey != "" {
		keyStr = resourceIDKey
	}

	return s.provider.GetPresignedUploadURL(ctx, resourceType, keyStr, contentType, maxSize)
}

// @MappedFrom queryResourceDownloadInfo(Long requesterId, StorageResourceType resourceType, @Nullable Long resourceIdNum, @Nullable String resourceIdStr, List<Value> customAttributes)
func (s *StorageService) QueryResourceDownloadInfo(
	ctx context.Context,
	requesterID int64,
	resourceType constants.StorageResourceType,
	resourceIDNum *int64,
	resourceIDStr string,
) (string, error) {
	if resourceType == 0 {
		return "", errors.New("unrecognized storage resource type")
	}

	// resourceIdNum is passed through for MESSAGE_ATTACHMENT type; provider may use it.
	_ = resourceIDNum

	return s.provider.GetPresignedDownloadURL(ctx, resourceType, resourceIDStr)
}

// @MappedFrom shareMessageAttachmentWithUser(Long requesterId, @Nullable Long messageAttachmentIdNum, @Nullable String messageAttachmentIdStr, Long userIdToShareWith)
func (s *StorageService) ShareMessageAttachmentWithUser(ctx context.Context, requesterID int64, messageAttachmentIDNum *int64, messageAttachmentIDStr *string, userIDToShareWith int64) error {
	if requesterID == 0 {
		return errors.New("requesterID must not be null")
	}
	if userIDToShareWith == 0 {
		return errors.New("userIDToShareWith must not be null")
	}
	return s.provider.ShareMessageAttachmentWithUser(ctx, requesterID, messageAttachmentIDNum, messageAttachmentIDStr, userIDToShareWith)
}

// @MappedFrom shareMessageAttachmentWithGroup(Long requesterId, @Nullable Long messageAttachmentIdNum, @Nullable String messageAttachmentIdStr, Long groupIdToShareWith)
func (s *StorageService) ShareMessageAttachmentWithGroup(ctx context.Context, requesterID int64, messageAttachmentIDNum *int64, messageAttachmentIDStr *string, groupIDToShareWith int64) error {
	if requesterID == 0 {
		return errors.New("requesterID must not be null")
	}
	if groupIDToShareWith == 0 {
		return errors.New("groupIDToShareWith must not be null")
	}
	return s.provider.ShareMessageAttachmentWithGroup(ctx, requesterID, messageAttachmentIDNum, messageAttachmentIDStr, groupIDToShareWith)
}

// @MappedFrom unshareMessageAttachmentWithUser(Long requesterId, @Nullable Long messageAttachmentIdNum, @Nullable String messageAttachmentIdStr, Long userIdToUnshareWith)
func (s *StorageService) UnshareMessageAttachmentWithUser(ctx context.Context, requesterID int64, messageAttachmentIDNum *int64, messageAttachmentIDStr *string, userIDToUnshareWith int64) error {
	if requesterID == 0 {
		return errors.New("requesterID must not be null")
	}
	if userIDToUnshareWith == 0 {
		return errors.New("userIDToUnshareWith must not be null")
	}
	return s.provider.UnshareMessageAttachmentWithUser(ctx, requesterID, messageAttachmentIDNum, messageAttachmentIDStr, userIDToUnshareWith)
}

// @MappedFrom unshareMessageAttachmentWithGroup(Long requesterId, @Nullable Long messageAttachmentIdNum, @Nullable String messageAttachmentIdStr, Long groupIdToUnshareWith)
func (s *StorageService) UnshareMessageAttachmentWithGroup(ctx context.Context, requesterID int64, messageAttachmentIDNum *int64, messageAttachmentIDStr *string, groupIDToUnshareWith int64) error {
	if requesterID == 0 {
		return errors.New("requesterID must not be null")
	}
	if groupIDToUnshareWith == 0 {
		return errors.New("groupIDToUnshareWith must not be null")
	}
	return s.provider.UnshareMessageAttachmentWithGroup(ctx, requesterID, messageAttachmentIDNum, messageAttachmentIDStr, groupIDToUnshareWith)
}

// @MappedFrom queryMessageAttachmentInfosUploadedByRequester(Long requesterId, @Nullable DateRange creationDateRange)
func (s *StorageService) QueryMessageAttachmentInfosUploadedByRequester(ctx context.Context, requesterID int64, creationDateStart *time.Time, creationDateEnd *time.Time) ([]bo.StorageResourceInfo, error) {
	return s.provider.QueryMessageAttachmentInfosUploadedByRequester(ctx, requesterID, creationDateStart, creationDateEnd)
}

// @MappedFrom queryMessageAttachmentInfosInPrivateConversations(Long requesterId, @Nullable Set<Long> userIds, @Nullable DateRange creationDateRange, @Nullable Boolean areSharedByRequester)
func (s *StorageService) QueryMessageAttachmentInfosInPrivateConversations(ctx context.Context, requesterID int64, userIDs []int64, creationDateStart *time.Time, creationDateEnd *time.Time, areSharedByRequester *bool) ([]bo.StorageResourceInfo, error) {
	return s.provider.QueryMessageAttachmentInfosInPrivateConversations(ctx, requesterID, userIDs, creationDateStart, creationDateEnd, areSharedByRequester)
}

// @MappedFrom queryMessageAttachmentInfosInGroupConversations(Long requesterId, @Nullable Set<Long> groupIds, @Nullable Set<Long> userIds, @Nullable DateRange creationDateRange)
func (s *StorageService) QueryMessageAttachmentInfosInGroupConversations(ctx context.Context, requesterID int64, groupIDs []int64, userIDs []int64, creationDateStart *time.Time, creationDateEnd *time.Time) ([]bo.StorageResourceInfo, error) {
	return s.provider.QueryMessageAttachmentInfosInGroupConversations(ctx, requesterID, groupIDs, userIDs, creationDateStart, creationDateEnd)
}
