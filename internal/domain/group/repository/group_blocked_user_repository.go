package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"im.turms/server/internal/domain/group/po"
	turmsmongo "im.turms/server/internal/storage/mongo"
)

// @MappedFrom GroupBlocklistRepository
type GroupBlockedUserRepository interface {
	UpdateBlockedUsers(ctx context.Context, keys []po.GroupBlockedUserKey, blockDate *time.Time, requesterId *int64) error
	FindBlockedUserIds(ctx context.Context, groupId int64) ([]int64, error)
	FindBlockedUsers(ctx context.Context, groupIds, userIds []int64, blockDateRange *turmsmongo.DateRange, requesterIds []int64, page, size *int) ([]po.GroupBlockedUser, error)
	Insert(ctx context.Context, blockedUser *po.GroupBlockedUser) error
	Delete(ctx context.Context, groupID, userID int64) error
	Exists(ctx context.Context, groupID, userID int64) (bool, error)
	FindBlockedUsersByGroupID(ctx context.Context, groupID int64) ([]po.GroupBlockedUser, error)
	FilterBlockedUserIDs(ctx context.Context, groupID int64, userIDs []int64) ([]int64, error)
	CountBlockedUsers(ctx context.Context, groupIds []int64, userIds []int64, blockDateRange *turmsmongo.DateRange, requesterIds []int64) (int64, error)
	DeleteBlockedUsers(ctx context.Context, keys []po.GroupBlockedUserKey) error
}

type groupBlockedUserRepository struct {
	coll *mongo.Collection
}

func NewGroupBlockedUserRepository(client *turmsmongo.Client) GroupBlockedUserRepository {
	return &groupBlockedUserRepository{
		coll: client.Collection(po.CollectionNameGroupBlockedUser),
	}
}

func (r *groupBlockedUserRepository) Insert(ctx context.Context, blockedUser *po.GroupBlockedUser) error {
	_, err := r.coll.InsertOne(ctx, blockedUser)
	return err
}

func (r *groupBlockedUserRepository) Delete(ctx context.Context, groupID, userID int64) error {
	filter := bson.M{
		"_id.gid": groupID,
		"_id.uid": userID,
	}
	_, err := r.coll.DeleteOne(ctx, filter)
	return err
}

func (r *groupBlockedUserRepository) Exists(ctx context.Context, groupID, userID int64) (bool, error) {
	filter := bson.M{
		"_id.gid": groupID,
		"_id.uid": userID,
	}
	count, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *groupBlockedUserRepository) FindBlockedUsersByGroupID(ctx context.Context, groupID int64) ([]po.GroupBlockedUser, error) {
	filter := bson.M{"_id.gid": groupID}
	cursor, err := r.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []po.GroupBlockedUser
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *groupBlockedUserRepository) FilterBlockedUserIDs(ctx context.Context, groupID int64, userIDs []int64) ([]int64, error) {
	filter := bson.M{
		"_id.gid": groupID,
		"_id.uid": bson.M{"$in": userIDs},
	}
	opts := bson.M{"_id.uid": 1}
	cursor, err := r.coll.Find(ctx, filter, options.Find().SetProjection(opts))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []struct {
		ID struct {
			UserID int64 `bson:"uid"`
		} `bson:"_id"`
	}
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	blockedUserIDs := make([]int64, len(results))
	for i, res := range results {
		blockedUserIDs[i] = res.ID.UserID
	}
	return blockedUserIDs, nil
}

func (r *groupBlockedUserRepository) UpdateBlockedUsers(ctx context.Context, keys []po.GroupBlockedUserKey, blockDate *time.Time, requesterId *int64) error {
	if len(keys) == 0 {
		return nil
	}
	filter := bson.M{"_id": bson.M{"$in": keys}}
	updateFields := bson.M{}
	if blockDate != nil {
		updateFields["bd"] = blockDate
	}
	if requesterId != nil {
		updateFields["rid"] = requesterId
	}
	if len(updateFields) == 0 {
		return nil
	}
	update := bson.M{"$set": updateFields}
	_, err := r.coll.UpdateMany(ctx, filter, update)
	return err
}

func (r *groupBlockedUserRepository) FindBlockedUserIds(ctx context.Context, groupId int64) ([]int64, error) {
	filter := bson.M{"_id.gid": groupId}
	opts := options.Find().SetProjection(bson.M{"_id.uid": 1})
	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []struct {
		ID struct {
			UserID int64 `bson:"uid"`
		} `bson:"_id"`
	}
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	userIDs := make([]int64, len(results))
	for i, res := range results {
		userIDs[i] = res.ID.UserID
	}
	return userIDs, nil
}

func (r *groupBlockedUserRepository) FindBlockedUsers(
	ctx context.Context,
	groupIds []int64,
	userIds []int64,
	blockDateRange *turmsmongo.DateRange,
	requesterIds []int64,
	page *int,
	size *int,
) ([]po.GroupBlockedUser, error) {
	filter := bson.M{}
	if len(groupIds) > 0 && len(userIds) > 0 {
		// Optimization, though standard MongoDB driver might not need it as heavily, building an in query for keys
		var keys []bson.M
		for _, gid := range groupIds {
			for _, uid := range userIds {
				keys = append(keys, bson.M{"gid": gid, "uid": uid})
			}
		}
		filter["_id"] = bson.M{"$in": keys}
	} else if len(groupIds) > 0 {
		filter["_id.gid"] = bson.M{"$in": groupIds}
	} else if len(userIds) > 0 {
		filter["_id.uid"] = bson.M{"$in": userIds}
	}

	if blockDateRange != nil {
		if drBson := blockDateRange.ToBson(); drBson != nil {
			filter["bd"] = drBson
		}
	}
	if len(requesterIds) > 0 {
		filter["rid"] = bson.M{"$in": requesterIds}
	}

	opts := options.Find()
	if page != nil && size != nil {
		opts.SetSkip(int64(*page * *size))
		opts.SetLimit(int64(*size))
	} else if size != nil {
		opts.SetLimit(int64(*size))
	}

	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []po.GroupBlockedUser
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *groupBlockedUserRepository) CountBlockedUsers(
	ctx context.Context,
	groupIds []int64,
	userIds []int64,
	blockDateRange *turmsmongo.DateRange,
	requesterIds []int64,
) (int64, error) {
	filter := bson.M{}
	if len(groupIds) > 0 && len(userIds) > 0 {
		var keys []bson.M
		for _, gid := range groupIds {
			for _, uid := range userIds {
				keys = append(keys, bson.M{"gid": gid, "uid": uid})
			}
		}
		filter["_id"] = bson.M{"$in": keys}
	} else if len(groupIds) > 0 {
		filter["_id.gid"] = bson.M{"$in": groupIds}
	} else if len(userIds) > 0 {
		filter["_id.uid"] = bson.M{"$in": userIds}
	}

	if blockDateRange != nil {
		if drBson := blockDateRange.ToBson(); drBson != nil {
			filter["bd"] = drBson
		}
	}
	if len(requesterIds) > 0 {
		filter["rid"] = bson.M{"$in": requesterIds}
	}

	return r.coll.CountDocuments(ctx, filter)
}

func (r *groupBlockedUserRepository) DeleteBlockedUsers(ctx context.Context, keys []po.GroupBlockedUserKey) error {
	if len(keys) == 0 {
		return nil
	}
	filter := bson.M{"_id": bson.M{"$in": keys}}
	_, err := r.coll.DeleteMany(ctx, filter)
	return err
}
