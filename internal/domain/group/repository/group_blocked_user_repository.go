package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"im.turms/server/internal/domain/group/po"
	turmsmongo "im.turms/server/internal/storage/mongo"
)

type GroupBlockedUserRepository interface {
	Insert(ctx context.Context, blockedUser *po.GroupBlockedUser) error
	Delete(ctx context.Context, groupID, userID int64) error
	Exists(ctx context.Context, groupID, userID int64) (bool, error)
	FindBlockedUsersByGroupID(ctx context.Context, groupID int64) ([]po.GroupBlockedUser, error)
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
