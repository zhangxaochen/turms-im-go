package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"im.turms/server/internal/domain/user/po"
	turmsmongo "im.turms/server/internal/storage/mongo"
)

const CollectionNameUser = "user"

type UserRepository interface {
	Insert(ctx context.Context, user *po.User) error
	FindByID(ctx context.Context, userID int64) (*po.User, error)
	FindMany(ctx context.Context, filter bson.M) ([]*po.User, error)
	Update(ctx context.Context, userID int64, update bson.M) error
	DeleteMany(ctx context.Context, filter bson.M) (int64, error)
	Count(ctx context.Context, filter bson.M) (int64, error)
	Exists(ctx context.Context, userID int64) (bool, error)
	UpdateMany(ctx context.Context, filter bson.M, update bson.M) (int64, error)
	Aggregate(ctx context.Context, pipeline mongo.Pipeline) (*mongo.Cursor, error)
}

type userRepository struct {
	coll *mongo.Collection
}

func NewUserRepository(client *turmsmongo.Client) UserRepository {
	return &userRepository{
		coll: client.Collection(CollectionNameUser),
	}
}

func (r *userRepository) Insert(ctx context.Context, user *po.User) error {
	_, err := r.coll.InsertOne(ctx, user)
	return err
}

func (r *userRepository) FindByID(ctx context.Context, userID int64) (*po.User, error) {
	filter := bson.M{"_id": userID}
	var user po.User
	err := r.coll.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindMany(ctx context.Context, filter bson.M) ([]*po.User, error) {
	cursor, err := r.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*po.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) Update(ctx context.Context, userID int64, update bson.M) error {
	filter := bson.M{"_id": userID}
	updateBson := bson.M{"$set": update}
	_, err := r.coll.UpdateOne(ctx, filter, updateBson)
	return err
}

func (r *userRepository) DeleteMany(ctx context.Context, filter bson.M) (int64, error) {
	result, err := r.coll.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}

func (r *userRepository) Count(ctx context.Context, filter bson.M) (int64, error) {
	return r.coll.CountDocuments(ctx, filter)
}

func (r *userRepository) Exists(ctx context.Context, userID int64) (bool, error) {
	filter := bson.M{"_id": userID}
	count, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *userRepository) UpdateMany(ctx context.Context, filter bson.M, update bson.M) (int64, error) {
	updateBson := bson.M{"$set": update}
	result, err := r.coll.UpdateMany(ctx, filter, updateBson)
	if err != nil {
		return 0, err
	}
	return result.ModifiedCount, nil
}

func (r *userRepository) Aggregate(ctx context.Context, pipeline mongo.Pipeline) (*mongo.Cursor, error) {
	return r.coll.Aggregate(ctx, pipeline)
}

