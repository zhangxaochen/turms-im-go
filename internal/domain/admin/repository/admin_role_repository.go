package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"im.turms/server/internal/domain/admin/permission"
	"im.turms/server/internal/domain/admin/po"
	turmsmongo "im.turms/server/internal/storage/mongo"
)

type AdminRoleRepository interface {
	Insert(ctx context.Context, adminRole *po.AdminRole) error
	UpdateAdminRoles(ctx context.Context, roleIds []int64, newName *string, permissions []permission.AdminPermission, rank *int) (int64, error)
	CountAdminRoles(ctx context.Context, ids []int64, names []string, includedPermissions []permission.AdminPermission, ranks []int) (int64, error)
	FindAdminRoles(ctx context.Context, roleIds []int64, names []string, includedPermissions []permission.AdminPermission, ranks []int, page *int, size *int) ([]*po.AdminRole, error)
	FindAdminRolesByIdsAndRankGreaterThan(ctx context.Context, roleIds []int64, rankGreaterThan *int) ([]*po.AdminRole, error)
	FindHighestRankByRoleIds(ctx context.Context, roleIds []int64) (*int, error)
	DeleteAdminRoles(ctx context.Context, ids []int64) (int64, error)
}

type adminRoleRepository struct {
	coll *mongo.Collection
}

func NewAdminRoleRepository(client *turmsmongo.Client) AdminRoleRepository {
	return &adminRoleRepository{
		coll: client.Collection(po.CollectionNameAdminRole),
	}
}

func (r *adminRoleRepository) Insert(ctx context.Context, adminRole *po.AdminRole) error {
	_, err := r.coll.InsertOne(ctx, adminRole)
	return err
}

func (r *adminRoleRepository) UpdateAdminRoles(ctx context.Context, roleIds []int64, newName *string, permissions []permission.AdminPermission, rank *int) (int64, error) {
	if len(roleIds) == 0 {
		return 0, nil
	}
	filter := bson.M{
		po.AdminRoleFieldID: bson.M{"$in": roleIds},
	}

	set := bson.M{}
	if newName != nil {
		set[po.AdminRoleFieldName] = *newName
	}
	if len(permissions) > 0 {
		set[po.AdminRoleFieldPermissions] = permissions
	}
	if rank != nil {
		set[po.AdminRoleFieldRank] = *rank
	}

	if len(set) == 0 {
		return 0, nil
	}
	update := bson.M{"$set": set}

	res, err := r.coll.UpdateMany(ctx, filter, update)
	if err != nil {
		return 0, err
	}
	return res.ModifiedCount, nil
}

func (r *adminRoleRepository) CountAdminRoles(ctx context.Context, ids []int64, names []string, includedPermissions []permission.AdminPermission, ranks []int) (int64, error) {
	filter := r.buildFilter(ids, names, includedPermissions, ranks)
	return r.coll.CountDocuments(ctx, filter)
}

func (r *adminRoleRepository) FindAdminRoles(ctx context.Context, roleIds []int64, names []string, includedPermissions []permission.AdminPermission, ranks []int, page *int, size *int) ([]*po.AdminRole, error) {
	filter := r.buildFilter(roleIds, names, includedPermissions, ranks)

	findOptions := options.Find()
	if size != nil {
		findOptions.SetLimit(int64(*size))
		if page != nil {
			findOptions.SetSkip(int64(*page * *size))
		}
	}

	cursor, err := r.coll.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var roles []*po.AdminRole
	if err = cursor.All(ctx, &roles); err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *adminRoleRepository) FindAdminRolesByIdsAndRankGreaterThan(ctx context.Context, roleIds []int64, rankGreaterThan *int) ([]*po.AdminRole, error) {
	if len(roleIds) == 0 {
		return nil, nil
	}
	filter := bson.M{
		po.AdminRoleFieldID: bson.M{"$in": roleIds},
	}
	if rankGreaterThan != nil {
		filter[po.AdminRoleFieldRank] = bson.M{"$gt": *rankGreaterThan}
	}

	cursor, err := r.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var roles []*po.AdminRole
	if err = cursor.All(ctx, &roles); err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *adminRoleRepository) FindHighestRankByRoleIds(ctx context.Context, roleIds []int64) (*int, error) {
	if len(roleIds) == 0 {
		return nil, nil
	}

	filter := bson.M{
		po.AdminRoleFieldID: bson.M{"$in": roleIds},
	}

	findOptions := options.FindOne().SetSort(bson.M{po.AdminRoleFieldRank: -1})
	var role po.AdminRole
	err := r.coll.FindOne(ctx, filter, findOptions).Decode(&role)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &role.Rank, nil
}

func (r *adminRoleRepository) DeleteAdminRoles(ctx context.Context, ids []int64) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	filter := bson.M{
		po.AdminRoleFieldID: bson.M{"$in": ids},
	}
	res, err := r.coll.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}
	return res.DeletedCount, nil
}

func (r *adminRoleRepository) buildFilter(ids []int64, names []string, includedPermissions []permission.AdminPermission, ranks []int) bson.M {
	filter := bson.M{}
	if len(ids) > 0 {
		filter[po.AdminRoleFieldID] = bson.M{"$in": ids}
	}
	if len(names) > 0 {
		filter[po.AdminRoleFieldName] = bson.M{"$in": names}
	}
	if len(includedPermissions) > 0 {
		filter[po.AdminRoleFieldPermissions] = bson.M{"$in": includedPermissions}
	}
	if len(ranks) > 0 {
		filter[po.AdminRoleFieldRank] = bson.M{"$in": ranks}
	}
	return filter
}
