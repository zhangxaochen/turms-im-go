package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"im.turms/server/internal/domain/storage/constants"
	"im.turms/server/internal/domain/storage/provider"
	"im.turms/server/internal/domain/storage/service"
)

func TestStorageService(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mockProvider := provider.NewMockStorageProvider()
	storageService := service.NewStorageService(mockProvider)

	t.Run("Generate Upload URL", func(t *testing.T) {
		url, err := storageService.QueryResourceUploadInfo(
			ctx,
			1001,
			constants.StorageResourceTypeUserProfilePicture,
			nil,
			"1001.png",
			"image/png",
			1024*1024*5,
		)

		require.NoError(t, err)
		assert.Contains(t, url, "http://localhost:9000/mock/upload/")
		assert.Contains(t, url, "contentType=image/png")
	})

	t.Run("Generate Download URL", func(t *testing.T) {
		url, err := storageService.QueryResourceDownloadInfo(
			ctx,
			1001,
			constants.StorageResourceTypeGroupProfilePicture,
			nil,
			"group_5001.png",
		)

		require.NoError(t, err)
		assert.Contains(t, url, "http://localhost:9000/mock/")
		assert.Contains(t, url, "group_5001.png")
	})

	t.Run("Delete Resource", func(t *testing.T) {
		err := storageService.DeleteResource(
			ctx,
			1001,
			constants.StorageResourceTypeMessageAttachment,
			"msg_attach_999",
		)

		require.NoError(t, err)
	})
}
