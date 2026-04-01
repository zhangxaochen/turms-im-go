package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"im.turms/server/internal/domain/group/po"
	turmsmongo "im.turms/server/internal/storage/mongo"
)

const GroupCollectionName = "group"

type GroupRepository struct {
	client *turmsmongo.Client
	col    *mongo.Collection
}

func NewGroupRepository(client *turmsmongo.Client) *GroupRepository {
	return &GroupRepository{
		client: client,
		col:    client.Collection(GroupCollectionName),
	}
}

// InsertGroup inserts a new group entity.
func (r *GroupRepository) InsertGroup(ctx context.Context, group *po.Group) error {
	_, err := r.col.InsertOne(ctx, group)
	return err
}

// FindGroup retrieves a group by ID.
func (r *GroupRepository) FindGroup(ctx context.Context, groupID int64) (*po.Group, error) {
	filter := bson.M{"_id": groupID}
	var group po.Group
	err := r.col.FindOne(ctx, filter).Decode(&group)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &group, nil
}

// UpdateGroup updates specified fields of a group.
func (r *GroupRepository) UpdateGroup(ctx context.Context, groupID int64, update bson.M) error {
	filter := bson.M{"_id": groupID}
	updateBson := bson.M{"$set": update}
	_, err := r.col.UpdateOne(ctx, filter, updateBson)
	return err
}
