package service

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"im.turms/server/internal/domain/user/po"
	"im.turms/server/internal/domain/user/repository"
)

var (
	ErrImmutableSetting = errors.New("cannot update immutable setting")
)

type UserSettingsService interface {
	UpsertSettings(ctx context.Context, userID int64, settings map[string]interface{}) error
	DeleteSettings(ctx context.Context, filter bson.M) (int64, error)
	UnsetSettings(ctx context.Context, userID int64, keys []string) error
	QuerySettings(ctx context.Context, filter bson.M) ([]*po.UserSettings, error)
	QuerySetting(ctx context.Context, userID int64, names []byte) (*po.UserSettings, error)
}

type userSettingsService struct {
	settingsRepo repository.UserSettingsRepository
}

func NewUserSettingsService(settingsRepo repository.UserSettingsRepository) UserSettingsService {
	return &userSettingsService{
		settingsRepo: settingsRepo,
	}
}

func (s *userSettingsService) UpsertSettings(ctx context.Context, userID int64, settings map[string]interface{}) error {
	if len(settings) == 0 {
		return nil
	}
	// TODO: Add validation from TurmsProperties if needed (e.g. check immutable settings)
	return s.settingsRepo.UpsertSettings(ctx, userID, settings)
}

func (s *userSettingsService) DeleteSettings(ctx context.Context, filter bson.M) (int64, error) {
	return s.settingsRepo.DeleteSettings(ctx, filter)
}

func (s *userSettingsService) UnsetSettings(ctx context.Context, userID int64, keys []string) error {
	if len(keys) == 0 {
		return nil
	}
	return s.settingsRepo.UnsetSettings(ctx, userID, keys)
}

func (s *userSettingsService) QuerySettings(ctx context.Context, filter bson.M) ([]*po.UserSettings, error) {
	return s.settingsRepo.FindSettings(ctx, filter)
}

func (s *userSettingsService) QuerySetting(ctx context.Context, userID int64, names []byte) (*po.UserSettings, error) {
	// Parse names if provided
	var nameStrs []string
	if len(names) > 0 {
		for _, b := range names {
			nameStrs = append(nameStrs, string(b))
		}
	}
	return s.settingsRepo.FindByIdAndSettingNames(ctx, userID, nameStrs)
}

