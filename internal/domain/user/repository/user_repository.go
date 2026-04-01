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
	Update(ctx context.Context, userID int64, update bson.M) error
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

func (r *userRepository) Update(ctx context.Context, userID int64, update bson.M) error {
	filter := bson.M{"_id": userID}
	updateBson := bson.M{"$set": update}
	_, err := r.coll.UpdateOne(ctx, filter, updateBson)
	return err
}
