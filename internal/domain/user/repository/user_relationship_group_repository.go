package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"im.turms/server/internal/domain/user/po"
	turmsmongo "im.turms/server/internal/storage/mongo"
)

type UserRelationshipGroupRepository interface {
	InsertGroup(ctx context.Context, group *po.UserRelationshipGroup) error
	FindGroups(ctx context.Context, filter interface{}) ([]*po.UserRelationshipGroup, error)
	DeleteGroups(ctx context.Context, filter interface{}) (int64, error)
	UpdateGroups(ctx context.Context, filter interface{}, update interface{}) (int64, error)
	CountGroups(ctx context.Context, filter interface{}) (int64, error)
}

type userRelationshipGroupRepository struct {
	collection *mongo.Collection
}

func NewUserRelationshipGroupRepository(mongoClient *turmsmongo.Client) UserRelationshipGroupRepository {
	return &userRelationshipGroupRepository{
		collection: mongoClient.Collection(po.CollectionNameUserRelationshipGroup),
	}
}

func (r *userRelationshipGroupRepository) InsertGroup(ctx context.Context, group *po.UserRelationshipGroup) error {
	_, err := r.collection.InsertOne(ctx, group)
	return err
}

func (r *userRelationshipGroupRepository) FindGroups(ctx context.Context, filter interface{}) ([]*po.UserRelationshipGroup, error) {
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

func (r *userRelationshipGroupRepository) DeleteGroups(ctx context.Context, filter interface{}) (int64, error) {
	res, err := r.collection.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}
	return res.DeletedCount, nil
}

func (r *userRelationshipGroupRepository) UpdateGroups(ctx context.Context, filter interface{}, update interface{}) (int64, error) {
	res, err := r.collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return 0, err
	}
	return res.ModifiedCount, nil
}

func (r *userRelationshipGroupRepository) CountGroups(ctx context.Context, filter interface{}) (int64, error) {
	return r.collection.CountDocuments(ctx, filter)
}
