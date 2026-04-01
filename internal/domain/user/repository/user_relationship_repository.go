package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"im.turms/server/internal/domain/user/po"
	turmsmongo "im.turms/server/internal/storage/mongo"
)

type UserRelationshipRepository interface {
	HasRelationshipAndNotBlocked(ctx context.Context, ownerID, relatedUserID int64) (bool, error)
	Insert(ctx context.Context, rel *po.UserRelationship) error
	UpdateBlockDate(ctx context.Context, ownerID, relatedUserID int64, blockDate *time.Time) error
}

type userRelationshipRepository struct {
	coll *mongo.Collection
}

func NewUserRelationshipRepository(client *turmsmongo.Client) UserRelationshipRepository {
	return &userRelationshipRepository{
		coll: client.Collection(po.CollectionNameUserRelationship),
	}
}

func (r *userRelationshipRepository) HasRelationshipAndNotBlocked(ctx context.Context, ownerID, relatedUserID int64) (bool, error) {
	filter := bson.M{
		"_id": bson.M{
			"oid": ownerID,
			"rid": relatedUserID,
		},
		"bd": nil, // Field is null or not exists -> not blocked
	}

	count, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *userRelationshipRepository) Insert(ctx context.Context, rel *po.UserRelationship) error {
	_, err := r.coll.InsertOne(ctx, rel)
	return err
}

func (r *userRelationshipRepository) UpdateBlockDate(ctx context.Context, ownerID, relatedUserID int64, blockDate *time.Time) error {
	filter := bson.M{
		"_id": bson.M{
			"oid": ownerID,
			"rid": relatedUserID,
		},
	}
	update := bson.M{
		"$set": bson.M{"bd": blockDate},
	}
	_, err := r.coll.UpdateOne(ctx, filter, update)
	return err
}
