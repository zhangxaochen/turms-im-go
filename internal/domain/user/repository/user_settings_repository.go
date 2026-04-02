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
