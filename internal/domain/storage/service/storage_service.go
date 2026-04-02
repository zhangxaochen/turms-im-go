package service

import (
	"context"
	"errors"

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

func (s *StorageService) QueryResourceUploadInfo(
	ctx context.Context,
	requesterID int64,
	resourceType constants.StorageResourceType,
	resourceName string,
	contentType string,
	maxSize int64,
) (string, error) {
	if resourceType == 0 { // Treat 0 as unrecognized or default
		return "", errors.New("unrecognized storage resource type")
	}

	return s.provider.GetPresignedUploadURL(ctx, resourceType, resourceName, contentType, maxSize)
}

func (s *StorageService) QueryResourceDownloadInfo(
	ctx context.Context,
	requesterID int64,
	resourceType constants.StorageResourceType,
	resourceIDStr string,
) (string, error) {
	if resourceType == 0 {
		return "", errors.New("unrecognized storage resource type")
	}

	return s.provider.GetPresignedDownloadURL(ctx, resourceType, resourceIDStr)
}
