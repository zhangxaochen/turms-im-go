package repository

import (
	"context"
	"time"

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
	UpdateUsers(ctx context.Context, userIDs []int64, update bson.M) (int64, error)
	UpdateUsersDeletionDate(ctx context.Context, userIDs []int64) (int64, error)
	CountRegisteredUsers(ctx context.Context, startDate *time.Time, endDate *time.Time, queryDeletedRecords bool) (int64, error)
	CountDeletedUsers(ctx context.Context, startDate *time.Time, endDate *time.Time) (int64, error)
	CountUsers(ctx context.Context, startDate *time.Time, endDate *time.Time) (int64, error)
	CountAllUsers(ctx context.Context) (int64, error)
	FindName(ctx context.Context, userID int64) (string, error)
	FindAllNames(ctx context.Context) ([]string, error)
	FindProfileAccessIfNotDeleted(ctx context.Context, userID int64) (*int32, error)
	FindUsers(ctx context.Context, userIDs []int64) ([]*po.User, error)
	FindNotDeletedUserProfiles(ctx context.Context, userIDs []int64) ([]*po.User, error)
	FindUsersProfile(ctx context.Context, userIDs []int64) ([]*po.User, error)
	FindUserRoleID(ctx context.Context, userID int64) (*int64, error)
	IsActiveAndNotDeleted(ctx context.Context, userID int64) (bool, error)
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

func (r *userRepository) UpdateUsers(ctx context.Context, userIDs []int64, update bson.M) (int64, error) {
	if len(userIDs) == 0 {
		return 0, nil
	}
	filter := bson.M{"_id": bson.M{"$in": userIDs}}
	return r.UpdateMany(ctx, filter, update)
}

func (r *userRepository) UpdateUsersDeletionDate(ctx context.Context, userIDs []int64) (int64, error) {
	if len(userIDs) == 0 {
		return 0, nil
	}
	filter := bson.M{"_id": bson.M{"$in": userIDs}}
	return r.UpdateMany(ctx, filter, bson.M{"dd": time.Now()})
}

func (r *userRepository) CountRegisteredUsers(ctx context.Context, startDate *time.Time, endDate *time.Time, queryDeletedRecords bool) (int64, error) {
	filter := bson.M{}
	if startDate != nil || endDate != nil {
		cdFilter := bson.M{}
		if startDate != nil {
			cdFilter["$gte"] = startDate
		}
		if endDate != nil {
			cdFilter["$lt"] = endDate
		}
		filter["cd"] = cdFilter
	}
	if !queryDeletedRecords {
		filter["dd"] = nil
	}
	return r.Count(ctx, filter)
}

func (r *userRepository) CountDeletedUsers(ctx context.Context, startDate *time.Time, endDate *time.Time) (int64, error) {
	filter := bson.M{"dd": bson.M{"$ne": nil}}
	if startDate != nil || endDate != nil {
		ddFilter := bson.M{}
		if startDate != nil {
			ddFilter["$gte"] = startDate
		}
		if endDate != nil {
			ddFilter["$lt"] = endDate
		}
		filter["dd"] = ddFilter
	}
	return r.Count(ctx, filter)
}

func (r *userRepository) CountUsers(ctx context.Context, startDate *time.Time, endDate *time.Time) (int64, error) {
	filter := bson.M{}
	if startDate != nil || endDate != nil {
		cdFilter := bson.M{}
		if startDate != nil {
			cdFilter["$gte"] = startDate
		}
		if endDate != nil {
			cdFilter["$lt"] = endDate
		}
		filter["cd"] = cdFilter
	}
	return r.Count(ctx, filter)
}

func (r *userRepository) CountAllUsers(ctx context.Context) (int64, error) {
	return r.Count(ctx, bson.M{})
}

func (r *userRepository) FindName(ctx context.Context, userID int64) (string, error) {
	user, err := r.FindByID(ctx, userID)
	if err != nil || user == nil {
		return "", err
	}
	return user.Name, nil
}

func (r *userRepository) FindAllNames(ctx context.Context) ([]string, error) {
	users, err := r.FindMany(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(users))
	for _, user := range users {
		names = append(names, user.Name)
	}
	return names, nil
}

func (r *userRepository) FindProfileAccessIfNotDeleted(ctx context.Context, userID int64) (*int32, error) {
	user, err := r.FindByID(ctx, userID)
	if err != nil || user == nil || user.DeletionDate != nil {
		return nil, err
	}
	return &user.ProfileAccess, nil
}

func (r *userRepository) FindUsers(ctx context.Context, userIDs []int64) ([]*po.User, error) {
	if len(userIDs) == 0 {
		return []*po.User{}, nil
	}
	filter := bson.M{"_id": bson.M{"$in": userIDs}}
	return r.FindMany(ctx, filter)
}

func (r *userRepository) FindNotDeletedUserProfiles(ctx context.Context, userIDs []int64) ([]*po.User, error) {
	if len(userIDs) == 0 {
		return []*po.User{}, nil
	}
	filter := bson.M{"_id": bson.M{"$in": userIDs}, "dd": nil}
	return r.FindMany(ctx, filter)
}

func (r *userRepository) FindUsersProfile(ctx context.Context, userIDs []int64) ([]*po.User, error) {
	if len(userIDs) == 0 {
		return []*po.User{}, nil
	}
	filter := bson.M{"_id": bson.M{"$in": userIDs}}
	return r.FindMany(ctx, filter)
}

func (r *userRepository) FindUserRoleID(ctx context.Context, userID int64) (*int64, error) {
	user, err := r.FindByID(ctx, userID)
	if err != nil || user == nil {
		return nil, err
	}
	return &user.PermissionGroupID, nil
}

func (r *userRepository) IsActiveAndNotDeleted(ctx context.Context, userID int64) (bool, error) {
	user, err := r.FindByID(ctx, userID)
	if err != nil || user == nil {
		return false, err
	}
	return user.IsActive && user.DeletionDate == nil, nil
}

func (r *userRepository) Aggregate(ctx context.Context, pipeline mongo.Pipeline) (*mongo.Cursor, error) {
	return r.coll.Aggregate(ctx, pipeline)
}
