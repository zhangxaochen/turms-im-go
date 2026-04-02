package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"im.turms/server/internal/domain/user/po"
	turmsmongo "im.turms/server/internal/storage/mongo"
)

type UserSettingsRepository interface {
	UpsertSettings(ctx context.Context, userID int64, settings map[string]interface{}) error
	DeleteSettings(ctx context.Context, filter interface{}) (int64, error)
	FindSettings(ctx context.Context, filter interface{}) ([]*po.UserSettings, error)
	UnsetSettings(ctx context.Context, userID int64, settingsNames []string) error
	FindByIdAndSettingNames(ctx context.Context, userID int64, names []string) (*po.UserSettings, error)
}

type userSettingsRepository struct {
	collection *mongo.Collection
}

func NewUserSettingsRepository(mongoClient *turmsmongo.Client) UserSettingsRepository {
	return &userSettingsRepository{
		collection: mongoClient.Collection(po.CollectionNameUserSettings),
	}
}

func (r *userSettingsRepository) UpsertSettings(ctx context.Context, userID int64, settings map[string]interface{}) error {
	update := map[string]interface{}{
		"$set": map[string]interface{}{
			"s": settings,
		},
	}
	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateOne(ctx, map[string]interface{}{"_id": userID}, update, opts)
	return err
}

func (r *userSettingsRepository) DeleteSettings(ctx context.Context, filter interface{}) (int64, error) {
	res, err := r.collection.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}
	return res.DeletedCount, nil
}

func (r *userSettingsRepository) FindSettings(ctx context.Context, filter interface{}) ([]*po.UserSettings, error) {
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var settings []*po.UserSettings
	if err = cursor.All(ctx, &settings); err != nil {
		return nil, err
	}
	return settings, nil
}

func (r *userSettingsRepository) UnsetSettings(ctx context.Context, userID int64, settingsNames []string) error {
	if len(settingsNames) == 0 {
		return nil
	}
	unsetMap := make(map[string]interface{})
	for _, name := range settingsNames {
		unsetMap["s."+name] = ""
	}
	update := map[string]interface{}{"$unset": unsetMap}
	_, err := r.collection.UpdateOne(ctx, map[string]interface{}{"_id": userID}, update)
	return err
}

func (r *userSettingsRepository) FindByIdAndSettingNames(ctx context.Context, userID int64, names []string) (*po.UserSettings, error) {
	projection := map[string]interface{}{}
	if len(names) > 0 {
		for _, name := range names {
			projection["s."+name] = 1
		}
	}
	opts := options.FindOne().SetProjection(projection)
	var settings po.UserSettings
	err := r.collection.FindOne(ctx, map[string]interface{}{"_id": userID}, opts).Decode(&settings)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &settings, nil
}
