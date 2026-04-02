package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"im.turms/server/internal/domain/user/po"
	turmsmongo "im.turms/server/internal/storage/mongo"
)

type UserRelationshipGroupRepository interface {
	Insert(ctx context.Context, group *po.UserRelationshipGroup, session *mongo.Session) error
	DeleteByIds(ctx context.Context, keys []po.UserRelationshipGroupKey, session *mongo.Session) (int64, error)
	DeleteAllRelationshipGroups(ctx context.Context, ownerIDs []int64, session *mongo.Session) (int64, error)
	DeleteRelationshipGroups(ctx context.Context, ownerID int64, groupIndexes []int32, session *mongo.Session) (int64, error)
	UpdateRelationshipGroups(ctx context.Context, keys []po.UserRelationshipGroupKey, newName string, session *mongo.Session) (int64, error)
	UpdateRelationshipGroupName(ctx context.Context, ownerID int64, groupIndex int32, newName string, session *mongo.Session) (int64, error)
	CountRelationshipGroups(ctx context.Context, ownerIDs []int64, groupIndexes []int32) (int64, error)
	FindRelationshipGroups(ctx context.Context, ownerIDs []int64, groupIndexes []int32, page *int, size *int) ([]*po.UserRelationshipGroup, error)
	FindRelationshipGroupsInfos(ctx context.Context, ownerID int64) ([]*po.UserRelationshipGroup, error)
}

type userRelationshipGroupRepository struct {
	collection *mongo.Collection
}

func NewUserRelationshipGroupRepository(mongoClient *turmsmongo.Client) UserRelationshipGroupRepository {
	return &userRelationshipGroupRepository{
		collection: mongoClient.Collection(po.CollectionNameUserRelationshipGroup),
	}
}

func (r *userRelationshipGroupRepository) Insert(ctx context.Context, group *po.UserRelationshipGroup, session *mongo.Session) error {
	var err error
	if session != nil {
		err = mongo.WithSession(ctx, *session, func(sc mongo.SessionContext) error {
			_, err = r.collection.InsertOne(sc, group)
			return err
		})
	} else {
		_, err = r.collection.InsertOne(ctx, group)
	}
	return err
}

func (r *userRelationshipGroupRepository) DeleteByIds(ctx context.Context, keys []po.UserRelationshipGroupKey, session *mongo.Session) (int64, error) {
	if len(keys) == 0 {
		return 0, nil
	}
	filters := make([]bson.M, len(keys))
	for i, key := range keys {
		filters[i] = bson.M{
			"_id.oid":  key.OwnerID,
			"_id.gidx": key.Index,
		}
	}
	filter := bson.M{"$or": filters}
	var res *mongo.DeleteResult
	var err error
	if session != nil {
		err = mongo.WithSession(ctx, *session, func(sc mongo.SessionContext) error {
			res, err = r.collection.DeleteMany(sc, filter)
			return err
		})
	} else {
		res, err = r.collection.DeleteMany(ctx, filter)
	}
	if err != nil {
		return 0, err
	}
	return res.DeletedCount, nil
}

// @MappedFrom deleteAllRelationshipGroups(@NotEmpty Set<Long> ownerIds, @Nullable ClientSession session, boolean updateRelationshipGroupsVersion)
// @MappedFrom deleteAllRelationshipGroups(Set<Long> ownerIds, @Nullable ClientSession session)
func (r *userRelationshipGroupRepository) DeleteAllRelationshipGroups(ctx context.Context, ownerIDs []int64, session *mongo.Session) (int64, error) {
	if len(ownerIDs) == 0 {
		return 0, nil
	}
	filter := bson.M{
		"_id.oid": bson.M{"$in": ownerIDs},
	}
	var res *mongo.DeleteResult
	var err error
	if session != nil {
		err = mongo.WithSession(ctx, *session, func(sc mongo.SessionContext) error {
			res, err = r.collection.DeleteMany(sc, filter)
			return err
		})
	} else {
		res, err = r.collection.DeleteMany(ctx, filter)
	}
	if err != nil {
		return 0, err
	}
	return res.DeletedCount, nil
}

// @MappedFrom deleteRelationshipGroups()
// @MappedFrom deleteRelationshipGroups(@QueryParam(required = false)
// @MappedFrom deleteRelationshipGroups(@NotEmpty Set<UserRelationshipGroup.@ValidUserRelationshipGroupKey Key> keys)
func (r *userRelationshipGroupRepository) DeleteRelationshipGroups(ctx context.Context, ownerID int64, groupIndexes []int32, session *mongo.Session) (int64, error) {
	filter := bson.M{
		"_id.oid": ownerID,
	}
	if len(groupIndexes) > 0 {
		filter["_id.gidx"] = bson.M{"$in": groupIndexes}
	}
	var res *mongo.DeleteResult
	var err error
	if session != nil {
		err = mongo.WithSession(ctx, *session, func(sc mongo.SessionContext) error {
			res, err = r.collection.DeleteMany(sc, filter)
			return err
		})
	} else {
		res, err = r.collection.DeleteMany(ctx, filter)
	}
	if err != nil {
		return 0, err
	}
	return res.DeletedCount, nil
}

// @MappedFrom updateRelationshipGroups(@NotEmpty Set<UserRelationshipGroup.@ValidUserRelationshipGroupKey Key> keys, @Nullable String name, @Nullable @PastOrPresent Date creationDate)
// @MappedFrom updateRelationshipGroups(Set<UserRelationshipGroup.Key> keys, @Nullable String name, @Nullable Date creationDate)
// @MappedFrom updateRelationshipGroups(List<UserRelationshipGroup.Key> keys, @RequestBody UpdateRelationshipGroupDTO updateRelationshipGroupDTO)
func (r *userRelationshipGroupRepository) UpdateRelationshipGroups(ctx context.Context, keys []po.UserRelationshipGroupKey, newName string, session *mongo.Session) (int64, error) {
	if len(keys) == 0 {
		return 0, nil
	}
	filters := make([]bson.M, len(keys))
	for i, key := range keys {
		filters[i] = bson.M{
			"_id.oid":  key.OwnerID,
			"_id.gidx": key.Index,
		}
	}
	filter := bson.M{"$or": filters}
	update := bson.M{
		"$set": bson.M{"n": newName},
	}
	var res *mongo.UpdateResult
	var err error
	if session != nil {
		err = mongo.WithSession(ctx, *session, func(sc mongo.SessionContext) error {
			res, err = r.collection.UpdateMany(sc, filter, update)
			return err
		})
	} else {
		res, err = r.collection.UpdateMany(ctx, filter, update)
	}
	if err != nil {
		return 0, err
	}
	return res.ModifiedCount, nil
}

// @MappedFrom updateRelationshipGroupName(Long ownerId, Integer groupIndex, String newGroupName)
// @MappedFrom updateRelationshipGroupName(@NotNull Long ownerId, @NotNull Integer groupIndex, @NotNull String newGroupName)
func (r *userRelationshipGroupRepository) UpdateRelationshipGroupName(ctx context.Context, ownerID int64, groupIndex int32, newName string, session *mongo.Session) (int64, error) {
	filter := bson.M{
		"_id.oid":  ownerID,
		"_id.gidx": groupIndex,
	}
	update := bson.M{
		"$set": bson.M{"n": newName},
	}
	var res *mongo.UpdateResult
	var err error
	if session != nil {
		err = mongo.WithSession(ctx, *session, func(sc mongo.SessionContext) error {
			res, err = r.collection.UpdateOne(sc, filter, update)
			return err
		})
	} else {
		res, err = r.collection.UpdateOne(ctx, filter, update)
	}
	if err != nil {
		return 0, err
	}
	return res.ModifiedCount, nil
}

// @MappedFrom countRelationshipGroups(@Nullable Set<Long> ownerIds, @Nullable Set<Integer> indexes, @Nullable Set<String> names, @Nullable DateRange creationDateRange)
// @MappedFrom countRelationshipGroups(@Nullable Set<Long> ownerIds, @Nullable Set<Long> relatedUserIds)
func (r *userRelationshipGroupRepository) CountRelationshipGroups(ctx context.Context, ownerIDs []int64, groupIndexes []int32) (int64, error) {
	filter := bson.M{}
	if len(ownerIDs) > 0 {
		filter["_id.oid"] = bson.M{"$in": ownerIDs}
	}
	if len(groupIndexes) > 0 {
		filter["_id.gidx"] = bson.M{"$in": groupIndexes}
	}
	return r.collection.CountDocuments(ctx, filter)
}

// @MappedFrom findRelationshipGroups(Long userId)
// @MappedFrom findRelationshipGroups(@Nullable Set<Long> ownerIds, @Nullable Set<Integer> indexes, @Nullable Set<String> names, @Nullable DateRange creationDateRange, @Nullable Integer page, @Nullable Integer size)
func (r *userRelationshipGroupRepository) FindRelationshipGroups(ctx context.Context, ownerIDs []int64, groupIndexes []int32, page *int, size *int) ([]*po.UserRelationshipGroup, error) {
	filter := bson.M{}
	if len(ownerIDs) > 0 {
		filter["_id.oid"] = bson.M{"$in": ownerIDs}
	}
	if len(groupIndexes) > 0 {
		filter["_id.gidx"] = bson.M{"$in": groupIndexes}
	}

	opts := options.Find()
	if page != nil && size != nil {
		opts.SetSkip(int64(*page * *size))
		opts.SetLimit(int64(*size))
	} else if size != nil {
		opts.SetLimit(int64(*size))
	}

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	var groups []*po.UserRelationshipGroup
	if err = cursor.All(ctx, &groups); err != nil {
		return nil, err
	}
	return groups, nil
}

// @MappedFrom findRelationshipGroupsInfos(Long ownerId)
func (r *userRelationshipGroupRepository) FindRelationshipGroupsInfos(ctx context.Context, ownerID int64) ([]*po.UserRelationshipGroup, error) {
	filter := bson.M{
		"_id.oid": ownerID,
	}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var groups []*po.UserRelationshipGroup
	if err = cursor.All(ctx, &groups); err != nil {
		return nil, err
	}
	return groups, nil
}
