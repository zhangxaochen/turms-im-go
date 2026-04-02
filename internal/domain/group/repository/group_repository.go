package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

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

// InsertGroup inserts a new group into MongoDB.
func (r *GroupRepository) InsertGroup(ctx context.Context, group *po.Group) error {
	_, err := r.col.InsertOne(ctx, group)
	return err
}

// FindGroups retrieves multiple groups by their IDs, filtering out deleted ones.
func (r *GroupRepository) FindGroups(ctx context.Context, groupIDs []int64) ([]*po.Group, error) {
	filter := bson.M{
		"_id": bson.M{"$in": groupIDs},
		"dd":  bson.M{"$exists": false}, // Ensure DeletionDate does not exist
	}

	cursor, err := r.col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var groups []*po.Group
	if err := cursor.All(ctx, &groups); err != nil {
		return nil, err
	}
	return groups, nil
}

// FindGroupOwnerID retrieves the owner ID of a specific group.
func (r *GroupRepository) FindGroupOwnerID(ctx context.Context, groupID int64) (*int64, error) {
	filter := bson.M{"_id": groupID}
	opts := options.FindOne().SetProjection(bson.M{"oid": 1})

	var group po.Group
	if err := r.col.FindOne(ctx, filter, opts).Decode(&group); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Group not found
		}
		return nil, err
	}
	return group.OwnerID, nil
}

func (r *GroupRepository) FindGroup(ctx context.Context, groupID int64) (*po.Group, error) {
	filter := bson.M{"_id": groupID}
	var group po.Group
	if err := r.col.FindOne(ctx, filter).Decode(&group); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &group, nil
}

// CountOwnedGroups counts the number of groups owned by a specific user.
func (r *GroupRepository) CountOwnedGroups(ctx context.Context, ownerID int64) (int64, error) {
	filter := bson.M{"oid": ownerID, "dd": bson.M{"$exists": false}}
	return r.col.CountDocuments(ctx, filter)
}

// UpdateGroup modifies specified fields of a group.
func (r *GroupRepository) UpdateGroup(ctx context.Context, groupID int64, update bson.M) error {
	filter := bson.M{"_id": groupID}
	_, err := r.col.UpdateOne(ctx, filter, bson.M{"$set": update})
	return err
}
// DeleteGroup removes a group from MongoDB.
func (r *GroupRepository) DeleteGroup(ctx context.Context, groupID int64) error {
	filter := bson.M{"_id": groupID}
	_, err := r.col.DeleteOne(ctx, filter)
	return err
}
