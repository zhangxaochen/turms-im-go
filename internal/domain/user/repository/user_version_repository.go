package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"im.turms/server/internal/domain/user/po"
	turmsmongo "im.turms/server/internal/storage/mongo"
)

type UserVersionRepository interface {
	UpsertEmptyUserVersion(ctx context.Context, userID int64) error
	UpdateUserVersion(ctx context.Context, userID int64, update interface{}) error
	UpdateUserVersions(ctx context.Context, userIDs []int64, update interface{}) error
	FindUserVersion(ctx context.Context, userID int64) (*po.UserVersion, error)
	DeleteUserVersion(ctx context.Context, userID int64) error
	DeleteUserVersions(ctx context.Context, userIDs []int64) error
}

type userVersionRepository struct {
	collection *mongo.Collection
}

func NewUserVersionRepository(mongoClient *turmsmongo.Client) UserVersionRepository {
	return &userVersionRepository{
		collection: mongoClient.Collection(po.CollectionNameUserVersion),
	}
}

func (r *userVersionRepository) UpsertEmptyUserVersion(ctx context.Context, userID int64) error {
	update := map[string]interface{}{
		"$setOnInsert": po.UserVersion{
			UserID: userID,
		},
	}
	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateOne(ctx, map[string]interface{}{"_id": userID}, update, opts)
	return err
}

func (r *userVersionRepository) UpdateUserVersion(ctx context.Context, userID int64, update interface{}) error {
	_, err := r.collection.UpdateOne(ctx, map[string]interface{}{"_id": userID}, update)
	return err
}

func (r *userVersionRepository) FindUserVersion(ctx context.Context, userID int64) (*po.UserVersion, error) {
	var version po.UserVersion
	err := r.collection.FindOne(ctx, map[string]interface{}{"_id": userID}).Decode(&version)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &version, nil
}

func (r *userVersionRepository) DeleteUserVersion(ctx context.Context, userID int64) error {
	_, err := r.collection.DeleteOne(ctx, map[string]interface{}{"_id": userID})
	return err
}

func (r *userVersionRepository) UpdateUserVersions(ctx context.Context, userIDs []int64, update interface{}) error {
	filter := map[string]interface{}{
		"_id": map[string]interface{}{
			"$in": userIDs,
		},
	}
	_, err := r.collection.UpdateMany(ctx, filter, update)
	return err
}

func (r *userVersionRepository) DeleteUserVersions(ctx context.Context, userIDs []int64) error {
	filter := map[string]interface{}{
		"_id": map[string]interface{}{
			"$in": userIDs,
		},
	}
	_, err := r.collection.DeleteMany(ctx, filter)
	return err
}
