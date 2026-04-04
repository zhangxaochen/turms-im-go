package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"im.turms/server/internal/domain/conversation/po"
	turmsmongo "im.turms/server/internal/storage/mongo"
)

// ConversationSettingsRepository maps to ConversationSettingsRepository.java
// @MappedFrom ConversationSettingsRepository
type ConversationSettingsRepository struct {
	client *turmsmongo.Client
	col    *mongo.Collection
}

func NewConversationSettingsRepository(client *turmsmongo.Client) *ConversationSettingsRepository {
	return &ConversationSettingsRepository{
		client: client,
		col:    client.Collection(po.ConversationSettingsCollectionName),
	}
}

// @MappedFrom upsertSettings(Long ownerId, Long targetId, Map<String, Object> settings, Date lastUpdatedDate)
func (r *ConversationSettingsRepository) UpsertSettings(ctx context.Context, ownerId int64, targetId int64, settings map[string]any, lastUpdatedDate time.Time) error {
	if len(settings) == 0 {
		return nil
	}
	filter := bson.M{"_id": po.ConversationSettingsKey{OwnerId: ownerId, TargetId: targetId}}
	set := bson.M{po.ConversationSettingsFieldLastUpdatedDate: lastUpdatedDate}
	for k, v := range settings {
		set[po.ConversationSettingsFieldSettings+"."+k] = v
	}
	update := bson.M{"$set": set}
	opts := options.Update().SetUpsert(true)
	_, err := r.col.UpdateOne(ctx, filter, update, opts)
	return err
}

// @MappedFrom unsetSettings(Long ownerId, @Nullable Collection<Long> targetIds, @Nullable Collection<String> settingNames)
func (r *ConversationSettingsRepository) UnsetSettings(ctx context.Context, ownerId int64, targetIds []int64, settingNames []string) error {
	filter := bson.M{po.ConversationSettingsFieldIdOwnerId: ownerId}
	if len(targetIds) > 0 {
		filter["_id.tid"] = bson.M{"$in": targetIds}
	}

	var update bson.M
	if len(settingNames) > 0 {
		unset := bson.M{}
		for _, name := range settingNames {
			unset[po.ConversationSettingsFieldSettings+"."+name] = ""
		}
		update = bson.M{"$unset": unset}
	} else {
		update = bson.M{"$unset": bson.M{po.ConversationSettingsFieldSettings: ""}}
	}

	_, err := r.col.UpdateMany(ctx, filter, update)
	return err
}

// @MappedFrom unsetSettings(Collection<ConversationSettings.Key> keys, @Nullable Collection<String> settingNames)
func (r *ConversationSettingsRepository) UnsetSettingsWithKeys(ctx context.Context, keys []po.ConversationSettingsKey, settingNames []string) error {
	if len(keys) == 0 {
		return nil
	}
	filter := bson.M{"_id": bson.M{"$in": keys}}
	var update bson.M
	if len(settingNames) > 0 {
		unset := bson.M{}
		for _, name := range settingNames {
			unset[po.ConversationSettingsFieldSettings+"."+name] = ""
		}
		update = bson.M{"$unset": unset}
	} else {
		update = bson.M{"$unset": bson.M{po.ConversationSettingsFieldSettings: ""}}
	}
	_, err := r.col.UpdateMany(ctx, filter, update)
	return err
}

// @MappedFrom findByIdAndSettingNames(Long ownerId, Long targetId, @Nullable Collection<String> settingNames, @Nullable Date lastUpdatedDateStart)
func (r *ConversationSettingsRepository) FindByKey(ctx context.Context, ownerId int64, targetId int64, settingNames []string, lastUpdatedDateStart *time.Time) (*po.ConversationSettings, error) {
	filter := bson.M{"_id": po.ConversationSettingsKey{OwnerId: ownerId, TargetId: targetId}}
	if lastUpdatedDateStart != nil {
		filter[po.ConversationSettingsFieldLastUpdatedDate] = bson.M{"$gte": *lastUpdatedDateStart}
	}

	var projection bson.M
	if len(settingNames) > 0 {
		projection = bson.M{
			po.ConversationSettingsFieldLastUpdatedDate: 1,
		}
		for _, name := range settingNames {
			projection[po.ConversationSettingsFieldSettings+"."+name] = 1
		}
	}

	opts := options.FindOne()
	if projection != nil {
		opts.SetProjection(projection)
	}

	var settings po.ConversationSettings
	err := r.col.FindOne(ctx, filter, opts).Decode(&settings)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &settings, err
}

// @MappedFrom findByIdAndSettingNames(Long ownerId, @Nullable Collection<String> settingNames, @Nullable Date lastUpdatedDateStart)
func (r *ConversationSettingsRepository) FindByOwnerIdAndSettingNames(ctx context.Context, ownerId int64, settingNames []string, lastUpdatedDateStart *time.Time) ([]po.ConversationSettings, error) {
	filter := bson.M{po.ConversationSettingsFieldIdOwnerId: ownerId}
	if lastUpdatedDateStart != nil {
		filter[po.ConversationSettingsFieldLastUpdatedDate] = bson.M{"$gte": *lastUpdatedDateStart}
	}

	var projection bson.M
	if len(settingNames) > 0 {
		projection = bson.M{
			po.ConversationSettingsFieldLastUpdatedDate: 1,
		}
		for _, name := range settingNames {
			projection[po.ConversationSettingsFieldSettings+"."+name] = 1
		}
	}

	opts := options.Find()
	if projection != nil {
		opts.SetProjection(projection)
	}

	cursor, err := r.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var settingsList []po.ConversationSettings
	if err := cursor.All(ctx, &settingsList); err != nil {
		return nil, err
	}
	return settingsList, nil
}

// @MappedFrom findByIdAndSettingNames(Collection<ConversationSettings.Key> keys, @Nullable Collection<String> settingNames, @Nullable Date lastUpdatedDateStart)
func (r *ConversationSettingsRepository) FindByIdAndSettingNamesWithKeys(ctx context.Context, keys []po.ConversationSettingsKey, settingNames []string, lastUpdatedDateStart *time.Time) ([]po.ConversationSettings, error) {
	if len(keys) == 0 {
		return nil, nil
	}
	filter := bson.M{"_id": bson.M{"$in": keys}}
	if lastUpdatedDateStart != nil {
		filter[po.ConversationSettingsFieldLastUpdatedDate] = bson.M{"$gte": *lastUpdatedDateStart}
	}

	var projection bson.M
	if len(settingNames) > 0 {
		projection = bson.M{
			po.ConversationSettingsFieldLastUpdatedDate: 1,
		}
		for _, name := range settingNames {
			projection[po.ConversationSettingsFieldSettings+"."+name] = 1
		}
	}

	opts := options.Find()
	if projection != nil {
		opts.SetProjection(projection)
	}

	cursor, err := r.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var settingsList []po.ConversationSettings
	if err := cursor.All(ctx, &settingsList); err != nil {
		return nil, err
	}
	return settingsList, nil
}

// @MappedFrom findSettingFields(Long ownerId, Long targetId, Collection<String> includedFields)
func (r *ConversationSettingsRepository) FindSettingFields(ctx context.Context, ownerId int64, targetId int64, includedFields []string) (map[string]any, error) {
	filter := bson.M{"_id": po.ConversationSettingsKey{OwnerId: ownerId, TargetId: targetId}}
	projection := bson.M{}
	for _, field := range includedFields {
		projection[po.ConversationSettingsFieldSettings+"."+field] = 1
	}

	opts := options.FindOne().SetProjection(projection)
	var result struct {
		Settings map[string]any `bson:"s"`
	}
	err := r.col.FindOne(ctx, filter, opts).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return result.Settings, err
}

// @MappedFrom deleteByOwnerIds(Collection<Long> ownerIds, @Nullable ClientSession clientSession)
func (r *ConversationSettingsRepository) DeleteByOwnerIds(ctx context.Context, ownerIds []int64) (int64, error) {
	if len(ownerIds) == 0 {
		return 0, nil
	}
	filter := bson.M{po.ConversationSettingsFieldIdOwnerId: bson.M{"$in": ownerIds}}
	res, err := r.col.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}
	return res.DeletedCount, nil
}
