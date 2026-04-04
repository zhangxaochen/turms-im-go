package provider

import (
	"context"
	"fmt"
	"time"

	"im.turms/server/internal/domain/storage/bo"
	"im.turms/server/internal/domain/storage/constants"
)

type MockStorageProvider struct{}

func NewMockStorageProvider() *MockStorageProvider {
	return &MockStorageProvider{}
}

func (p *MockStorageProvider) GetPresignedDownloadURL(ctx context.Context, resourceType constants.StorageResourceType, keyStr string) (string, error) {
	return fmt.Sprintf("http://localhost:9000/mock/%d/%s", resourceType, keyStr), nil
}

func (p *MockStorageProvider) GetPresignedUploadURL(ctx context.Context, resourceType constants.StorageResourceType, keyStr string, contentType string, maxSize int64) (string, error) {
	return fmt.Sprintf("http://localhost:9000/mock/upload/%d/%s?contentType=%s&maxSize=%d", resourceType, keyStr, contentType, maxSize), nil
}

// @MappedFrom deleteResource(Long requesterId, StorageResourceType resourceType, @Nullable Long resourceIdNum, @Nullable String resourceIdStr, List<Value> customAttributes)
func (p *MockStorageProvider) DeleteResource(ctx context.Context, resourceType constants.StorageResourceType, keyStr string) error {
	return nil
}

func (p *MockStorageProvider) ShareMessageAttachmentWithUser(ctx context.Context, requesterID int64, messageAttachmentIDNum *int64, messageAttachmentIDStr *string, userIDToShareWith int64) error {
	return nil
}

func (p *MockStorageProvider) ShareMessageAttachmentWithGroup(ctx context.Context, requesterID int64, messageAttachmentIDNum *int64, messageAttachmentIDStr *string, groupIDToShareWith int64) error {
	return nil
}

func (p *MockStorageProvider) UnshareMessageAttachmentWithUser(ctx context.Context, requesterID int64, messageAttachmentIDNum *int64, messageAttachmentIDStr *string, userIDToUnshareWith int64) error {
	return nil
}

func (p *MockStorageProvider) UnshareMessageAttachmentWithGroup(ctx context.Context, requesterID int64, messageAttachmentIDNum *int64, messageAttachmentIDStr *string, groupIDToUnshareWith int64) error {
	return nil
}

func (p *MockStorageProvider) QueryMessageAttachmentInfosUploadedByRequester(ctx context.Context, requesterID int64, creationDateStart *time.Time, creationDateEnd *time.Time) ([]bo.StorageResourceInfo, error) {
	return []bo.StorageResourceInfo{}, nil
}

func (p *MockStorageProvider) QueryMessageAttachmentInfosInPrivateConversations(ctx context.Context, requesterID int64, userIDs []int64, creationDateStart *time.Time, creationDateEnd *time.Time, areSharedByRequester *bool) ([]bo.StorageResourceInfo, error) {
	return []bo.StorageResourceInfo{}, nil
}

func (p *MockStorageProvider) QueryMessageAttachmentInfosInGroupConversations(ctx context.Context, requesterID int64, groupIDs []int64, userIDs []int64, creationDateStart *time.Time, creationDateEnd *time.Time) ([]bo.StorageResourceInfo, error) {
	return []bo.StorageResourceInfo{}, nil
}
