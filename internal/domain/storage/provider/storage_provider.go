package provider

import (
	"context"
	"time"

	"im.turms/server/internal/domain/storage/bo"
	"im.turms/server/internal/domain/storage/constants"
)

type StorageProvider interface {
	GetPresignedDownloadURL(ctx context.Context, resourceType constants.StorageResourceType, keyStr string) (string, error)
	GetPresignedUploadURL(ctx context.Context, resourceType constants.StorageResourceType, keyStr string, contentType string, maxSize int64) (string, error)
	DeleteResource(ctx context.Context, resourceType constants.StorageResourceType, keyStr string) error

	// Message-attachment-specific methods
	ShareMessageAttachmentWithUser(ctx context.Context, requesterID int64, messageAttachmentIDNum *int64, messageAttachmentIDStr *string, userIDToShareWith int64) error
	ShareMessageAttachmentWithGroup(ctx context.Context, requesterID int64, messageAttachmentIDNum *int64, messageAttachmentIDStr *string, groupIDToShareWith int64) error
	UnshareMessageAttachmentWithUser(ctx context.Context, requesterID int64, messageAttachmentIDNum *int64, messageAttachmentIDStr *string, userIDToUnshareWith int64) error
	UnshareMessageAttachmentWithGroup(ctx context.Context, requesterID int64, messageAttachmentIDNum *int64, messageAttachmentIDStr *string, groupIDToUnshareWith int64) error
	QueryMessageAttachmentInfosUploadedByRequester(ctx context.Context, requesterID int64, creationDateStart *time.Time, creationDateEnd *time.Time) ([]bo.StorageResourceInfo, error)
	QueryMessageAttachmentInfosInPrivateConversations(ctx context.Context, requesterID int64, userIDs []int64, creationDateStart *time.Time, creationDateEnd *time.Time, areSharedByRequester *bool) ([]bo.StorageResourceInfo, error)
	QueryMessageAttachmentInfosInGroupConversations(ctx context.Context, requesterID int64, groupIDs []int64, userIDs []int64, creationDateStart *time.Time, creationDateEnd *time.Time) ([]bo.StorageResourceInfo, error)
}
