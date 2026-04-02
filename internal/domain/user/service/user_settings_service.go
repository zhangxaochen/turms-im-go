package service

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"im.turms/server/internal/domain/user/po"
	"im.turms/server/internal/domain/user/repository"
	common "im.turms/server/internal/domain/common/service"
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
	settingsRepo           repository.UserSettingsRepository
	outboundMessageService common.OutboundMessageService
}

func NewUserSettingsService(
	settingsRepo repository.UserSettingsRepository,
	outboundMessageService common.OutboundMessageService,
) UserSettingsService {
	return &userSettingsService{
		settingsRepo:           settingsRepo,
		outboundMessageService: outboundMessageService,
	}
}

func (s *userSettingsService) UpsertSettings(ctx context.Context, userID int64, settings map[string]interface{}) error {
	if len(settings) == 0 {
		return nil
	}

	// Basic validation for immutable settings
	// In a real scenario, this would come from TurmsProperties
	immutableSettings := map[string]bool{
		"user_id": true,
	}
	for k := range settings {
		if immutableSettings[k] {
			return ErrImmutableSetting
		}
	}

	err := s.settingsRepo.UpsertSettings(ctx, userID, settings)
	if err != nil {
		return err
	}

	// Notify other devices
	// TODO: Enable notification when protocol supports UserSettings/Value
	/*
		if s.outboundMessageService != nil {
			notification := &protocol.TurmsNotification{
				Data: &protocol.TurmsNotification_Data{
					Kind: &protocol.TurmsNotification_Data_UserSettings{
						UserSettings: &protocol.UserSettings{
							Settings: make(map[string]*protocol.Value),
						},
					},
				},
			}
			// Convert map[string]interface{} to map[string]*protocol.Value (simplified)
			s.outboundMessageService.ForwardNotification(ctx, notification, userID)
		}
	*/

	return nil
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

