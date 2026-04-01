package repository

import (
	"context"

	"im.turms/server/internal/domain/group/po"
	turmsmongo "im.turms/server/internal/storage/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type GroupMemberRepository interface {
	IsGroupMember(ctx context.Context, groupID, userID int64) (bool, error)
	Insert(ctx context.Context, member *po.GroupMember) error
}

type groupMemberRepository struct {
	coll *mongo.Collection
}

func NewGroupMemberRepository(client *turmsmongo.Client) GroupMemberRepository {
	return &groupMemberRepository{
		coll: client.Collection(po.CollectionNameGroupMember),
	}
}

func (r *groupMemberRepository) IsGroupMember(ctx context.Context, groupID, userID int64) (bool, error) {
	filter := bson.M{
		"_id": bson.M{
			"gid": groupID,
			"uid": userID,
		},
	}

	count, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *groupMemberRepository) Insert(ctx context.Context, member *po.GroupMember) error {
	_, err := r.coll.InsertOne(ctx, member)
	return err
}
