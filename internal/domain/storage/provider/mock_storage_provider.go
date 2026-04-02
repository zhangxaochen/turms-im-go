package provider

import (
	"context"
	"fmt"

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
