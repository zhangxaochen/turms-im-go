package provider

import (
	"context"

	"im.turms/server/internal/domain/storage/constants"
)

type StorageProvider interface {
	GetPresignedDownloadURL(ctx context.Context, resourceType constants.StorageResourceType, keyStr string) (string, error)
	GetPresignedUploadURL(ctx context.Context, resourceType constants.StorageResourceType, keyStr string, contentType string, maxSize int64) (string, error)
	DeleteResource(ctx context.Context, resourceType constants.StorageResourceType, keyStr string) error
}
