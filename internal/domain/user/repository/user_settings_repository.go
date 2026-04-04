package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"im.turms/server/internal/domain/user/po"
	turmsmongo "im.turms/server/internal/storage/mongo"
)

type UserSettingsRepository interface {
	UpsertSettings(ctx context.Context, userID int64, settings map[string]interface{}, lastUpdatedDate time.Time) error
	DeleteSettings(ctx context.Context, filter interface{}) (int64, error)
	FindSettings(ctx context.Context, filter interface{}) ([]*po.UserSettings, error)
	UnsetSettings(ctx context.Context, userID int64, settingsNames []string, lastUpdatedDate time.Time) error
	FindByIdAndSettingNames(ctx context.Context, userID int64, names []string, lastUpdatedDateStart *time.Time) (*po.UserSettings, error)
}

type userSettingsRepository struct {
	collection *mongo.Collection
}

func NewUserSettingsRepository(mongoClient *turmsmongo.Client) UserSettingsRepository {
	return &userSettingsRepository{
		collection: mongoClient.Collection(po.CollectionNameUserSettings),
	}
}

// @MappedFrom upsertSettings(Long userId, Map<String, Object> settings)
func (r *userSettingsRepository) UpsertSettings(ctx context.Context, userID int64, settings map[string]interface{}, lastUpdatedDate time.Time) error {
	setMap := make(map[string]interface{})
	setMap["lud"] = lastUpdatedDate
	for k, v := range settings {
		setMap["s."+k] = v
	}
	update := map[string]interface{}{
		"$set": setMap,
	}
	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateOne(ctx, map[string]interface{}{"_id": userID}, update, opts)
	return err
}

// @MappedFrom deleteSettings(Collection<Long> ownerIds, @Nullable ClientSession clientSession)
// @MappedFrom deleteSettings(Collection<Long> userIds, @Nullable ClientSession clientSession)
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

// @MappedFrom unsetSettings(Long userId, @Nullable Collection<String> settingNames)
func (r *userSettingsRepository) UnsetSettings(ctx context.Context, userID int64, settingsNames []string, lastUpdatedDate time.Time) error {
	update := bson.M{"$set": bson.M{"lud": lastUpdatedDate}}
	
	if len(settingsNames) > 0 {
		unsetMap := make(map[string]interface{})
		for _, name := range settingsNames {
			unsetMap["s."+name] = ""
		}
		update["$unset"] = unsetMap
	} else {
		update["$unset"] = bson.M{"s": ""}
	}
	
	_, err := r.collection.UpdateOne(ctx, map[string]interface{}{"_id": userID}, update)
	return err
}

// @MappedFrom findByIdAndSettingNames(Long userId, @Nullable Collection<String> settingNames, @Nullable Date lastUpdatedDateStart)
func (r *userSettingsRepository) FindByIdAndSettingNames(ctx context.Context, userID int64, names []string, lastUpdatedDateStart *time.Time) (*po.UserSettings, error) {
	filter := bson.M{"_id": userID}
	if lastUpdatedDateStart != nil {
		filter["lud"] = bson.M{"$gte": *lastUpdatedDateStart}
	}

	projection := map[string]interface{}{}
	if len(names) > 0 {
		projection["lud"] = 1
		for _, name := range names {
			projection["s."+name] = 1
		}
	}
	opts := options.FindOne().SetProjection(projection)
	var settings po.UserSettings
	err := r.collection.FindOne(ctx, filter, opts).Decode(&settings)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &settings, nil
}
