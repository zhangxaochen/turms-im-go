package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"im.turms/server/internal/domain/group/po"
	turmsmongo "im.turms/server/internal/storage/mongo"
)

const GroupTypeCollectionName = "groupType"

type GroupTypeRepository struct {
	client *turmsmongo.Client
	col    *mongo.Collection
}

func NewGroupTypeRepository(client *turmsmongo.Client) *GroupTypeRepository {
	return &GroupTypeRepository{
		client: client,
		col:    client.Collection(GroupTypeCollectionName),
	}
}

// FindGroupType retrieves a group type by its ID.
func (r *GroupTypeRepository) FindGroupType(ctx context.Context, typeID int64) (*po.GroupType, error) {
	filter := bson.M{"_id": typeID}
	var groupType po.GroupType
	if err := r.col.FindOne(ctx, filter).Decode(&groupType); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &groupType, nil
}

// InsertGroupType inserts a new GroupType into MongoDB.
func (r *GroupTypeRepository) InsertGroupType(ctx context.Context, groupType *po.GroupType) error {
	_, err := r.col.InsertOne(ctx, groupType)
	return err
}

// UpdateGroupType modifies an existing GroupType.
// @MappedFrom updateGroupType(Set<Long> ids, @RequestBody UpdateGroupTypeDTO updateGroupTypeDTO)
func (r *GroupTypeRepository) UpdateGroupType(ctx context.Context, typeID int64, update bson.M) error {
	filter := bson.M{"_id": typeID}
	_, err := r.col.UpdateOne(ctx, filter, bson.M{"$set": update})
	return err
}
