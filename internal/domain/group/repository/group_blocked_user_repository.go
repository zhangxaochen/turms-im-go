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
	UpdateBlockedUsers(ctx context.Context, keys []any, blockDate *time.Time, requesterId *int64) error
	FindBlockedUserIds(ctx context.Context, groupId int64) ([]int64, error)
	FindBlockedUsers(ctx context.Context, groupIds, userIds []int64, blockDateRange any, requesterIds []int64, page, size *int) ([]po.GroupBlockedUser, error)
	Insert(ctx context.Context, blockedUser *po.GroupBlockedUser) error
	Delete(ctx context.Context, groupID, userID int64) error
	Exists(ctx context.Context, groupID, userID int64) (bool, error)
	FindBlockedUsersByGroupID(ctx context.Context, groupID int64) ([]po.GroupBlockedUser, error)
	FilterBlockedUserIDs(ctx context.Context, groupID int64, userIDs []int64) ([]int64, error)
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

func (r *groupBlockedUserRepository) UpdateBlockedUsers(ctx context.Context, keys []any, blockDate *time.Time, requesterId *int64) error {
	return nil
}
func (r *groupBlockedUserRepository) FindBlockedUserIds(ctx context.Context, groupId int64) ([]int64, error) {
	return nil, nil
}
func (r *groupBlockedUserRepository) FindBlockedUsers(ctx context.Context, groupIds, userIds []int64, blockDateRange any, requesterIds []int64, page, size *int) ([]po.GroupBlockedUser, error) {
	return nil, nil
}
