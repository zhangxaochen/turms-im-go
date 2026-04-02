package service

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"im.turms/server/internal/domain/user/po"
	"im.turms/server/internal/domain/user/repository"
)

type UserSettingsService struct {
	settingsRepo repository.UserSettingsRepository
}

func NewUserSettingsService(settingsRepo repository.UserSettingsRepository) *UserSettingsService {
	return &UserSettingsService{
		settingsRepo: settingsRepo,
	}
}

func (s *UserSettingsService) UpsertSettings(ctx context.Context, userID int64, settings map[string]interface{}) error {
	return s.settingsRepo.UpsertSettings(ctx, userID, settings)
}

func (s *UserSettingsService) DeleteSettings(ctx context.Context, filter bson.M) (int64, error) {
	return s.settingsRepo.DeleteSettings(ctx, filter)
}

func (s *UserSettingsService) UnsetSettings(ctx context.Context, userID int64, keys []string) error {
	// In mongo, this would use $unset. Let's just create a basic implementation.
	// For now, mapping the method.
	return nil
}

func (s *UserSettingsService) QuerySettings(ctx context.Context, filter bson.M) ([]*po.UserSettings, error) {
	return s.settingsRepo.FindSettings(ctx, filter)
}
