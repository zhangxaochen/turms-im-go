package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"im.turms/server/internal/domain/user/po"
	turmsmongo "im.turms/server/internal/storage/mongo"
)

type UserRoleRepository interface {
	InsertRole(ctx context.Context, role *po.UserRole) error
	FindRoles(ctx context.Context, filter interface{}) ([]*po.UserRole, error)
	FindRoleByID(ctx context.Context, roleID int64) (*po.UserRole, error)
	UpdateRole(ctx context.Context, roleID int64, update interface{}) error
	DeleteRoles(ctx context.Context, filter interface{}) (int64, error)
	CountRoles(ctx context.Context, filter interface{}) (int64, error)
	UpdateUserRoles(ctx context.Context, roleIDs []int64, update interface{}) (int64, error)
}

type userRoleRepository struct {
	collection *mongo.Collection
}

func NewUserRoleRepository(mongoClient *turmsmongo.Client) UserRoleRepository {
	return &userRoleRepository{
		collection: mongoClient.Collection(po.CollectionNameUserRole),
	}
}

func (r *userRoleRepository) InsertRole(ctx context.Context, role *po.UserRole) error {
	_, err := r.collection.InsertOne(ctx, role)
	return err
}

func (r *userRoleRepository) FindRoles(ctx context.Context, filter interface{}) ([]*po.UserRole, error) {
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var roles []*po.UserRole
	if err = cursor.All(ctx, &roles); err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *userRoleRepository) FindRoleByID(ctx context.Context, roleID int64) (*po.UserRole, error) {
	var role po.UserRole
	err := r.collection.FindOne(ctx, map[string]interface{}{"_id": roleID}).Decode(&role)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Return nil if not found
		}
		return nil, err
	}
	return &role, nil
}

func (r *userRoleRepository) UpdateRole(ctx context.Context, roleID int64, update interface{}) error {
	_, err := r.collection.UpdateOne(ctx, map[string]interface{}{"_id": roleID}, update)
	return err
}

func (r *userRoleRepository) DeleteRoles(ctx context.Context, filter interface{}) (int64, error) {
	res, err := r.collection.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}
	return res.DeletedCount, nil
}

func (r *userRoleRepository) CountRoles(ctx context.Context, filter interface{}) (int64, error) {
	return r.collection.CountDocuments(ctx, filter)
}

func (r *userRoleRepository) UpdateUserRoles(ctx context.Context, roleIDs []int64, update interface{}) (int64, error) {
	if len(roleIDs) == 0 {
		return 0, nil
	}
	filter := map[string]interface{}{"_id": map[string]interface{}{"$in": roleIDs}}
	res, err := r.collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return 0, err
	}
	return res.ModifiedCount, nil
}
