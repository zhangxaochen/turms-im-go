package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"im.turms/server/internal/domain/admin/po"
	turmsmongo "im.turms/server/internal/storage/mongo"
)

type AdminRepository interface {
	Insert(ctx context.Context, admin *po.Admin) error
	UpdateAdmins(ctx context.Context, ids []int64, password []byte, displayName *string, roleIDs []int64) (int64, error)
	CountAdmins(ctx context.Context, ids []int64, roleIDs []int64) (int64, error)
	FindAdmins(ctx context.Context, ids []int64, loginNames []string, roleIDs []int64, page *int, size *int) ([]*po.Admin, error)
	DeleteAdmins(ctx context.Context, ids []int64) (int64, error)
}

type adminRepository struct {
	coll *mongo.Collection
}

func NewAdminRepository(client *turmsmongo.Client) AdminRepository {
	return &adminRepository{
		coll: client.Collection(po.CollectionNameAdmin),
	}
}

func (r *adminRepository) Insert(ctx context.Context, admin *po.Admin) error {
	_, err := r.coll.InsertOne(ctx, admin)
	return err
}

func (r *adminRepository) UpdateAdmins(ctx context.Context, ids []int64, password []byte, displayName *string, roleIDs []int64) (int64, error) {
	filter := bson.M{}
	if len(ids) > 0 {
		filter[po.AdminFieldID] = bson.M{"$in": ids}
	}

	update := bson.M{}
	set := bson.M{}
	unset := bson.M{}

	if len(password) > 0 {
		set[po.AdminFieldPassword] = password
	}
	if displayName != nil {
		set[po.AdminFieldDisplayName] = *displayName
	}
	if roleIDs != nil {
		if len(roleIDs) == 0 {
			unset[po.AdminFieldRoleIDs] = ""
		} else {
			set[po.AdminFieldRoleIDs] = roleIDs
		}
	}

	if len(set) > 0 {
		update["$set"] = set
	}
	if len(unset) > 0 {
		update["$unset"] = unset
	}

	if len(update) == 0 {
		return 0, nil
	}

	res, err := r.coll.UpdateMany(ctx, filter, update)
	if err != nil {
		return 0, err
	}
	return res.ModifiedCount, nil
}

func (r *adminRepository) CountAdmins(ctx context.Context, ids []int64, roleIDs []int64) (int64, error) {
	filter := bson.M{}
	if len(ids) > 0 {
		filter[po.AdminFieldID] = bson.M{"$in": ids}
	}
	if len(roleIDs) > 0 {
		filter[po.AdminFieldRoleIDs] = bson.M{"$in": roleIDs}
	}
	return r.coll.CountDocuments(ctx, filter)
}

func (r *adminRepository) FindAdmins(ctx context.Context, ids []int64, loginNames []string, roleIDs []int64, page *int, size *int) ([]*po.Admin, error) {
	filter := bson.M{}
	if len(ids) > 0 {
		filter[po.AdminFieldID] = bson.M{"$in": ids}
	}
	if len(loginNames) > 0 {
		filter[po.AdminFieldLoginName] = bson.M{"$in": loginNames}
	}
	if len(roleIDs) > 0 {
		filter[po.AdminFieldRoleIDs] = bson.M{"$in": roleIDs}
	}

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

	var admins []*po.Admin
	if err = cursor.All(ctx, &admins); err != nil {
		return nil, err
	}
	return admins, nil
}

func (r *adminRepository) DeleteAdmins(ctx context.Context, ids []int64) (int64, error) {
	filter := bson.M{}
	if len(ids) > 0 {
		filter[po.AdminFieldID] = bson.M{"$in": ids}
	}
	res, err := r.coll.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}
	return res.DeletedCount, nil
}
