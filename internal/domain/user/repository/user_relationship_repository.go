package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"im.turms/server/internal/domain/user/po"
	turmsmongo "im.turms/server/internal/storage/mongo"
)

type UserRelationshipRepository interface {
	HasRelationshipAndNotBlocked(ctx context.Context, ownerID, relatedUserID int64) (bool, error)
	Insert(ctx context.Context, rel *po.UserRelationship) error
	UpdateBlockDate(ctx context.Context, ownerID, relatedUserID int64, blockDate *time.Time) error
	FindRelatedUserIDs(ctx context.Context, ownerID int64, isBlocked *bool) ([]int64, error)
	FindRelationships(ctx context.Context, ownerIDs []int64, relatedUserIDs []int64) ([]po.UserRelationship, error)
	DeleteById(ctx context.Context, ownerID, relatedUserID int64) error
	Upsert(ctx context.Context, ownerID, relatedUserID int64, blockDate *time.Time, groupIndex *int32, establishmentDate *time.Time, name *string, session mongo.SessionContext) error
	DeleteAllRelationships(ctx context.Context, ownerIDs []int64) error
	DeleteOneSidedRelationships(ctx context.Context, ownerIDs, relatedUserIDs []int64) error
	CountRelationships(ctx context.Context, ownerIDs, relatedUserIDs []int64, groupIndexes []int32, isBlocked *bool) (int64, error)
	QueryMembersRelationships(ctx context.Context, ownerID int64, groupIndexes []int32, page, size *int) ([]po.UserRelationship, error)
	HasOneSidedRelationship(ctx context.Context, ownerID, relatedUserID int64) (bool, error)
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
		"_id": po.UserRelationshipKey{
			OwnerID:       ownerID,
			RelatedUserID: relatedUserID,
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
		"_id": po.UserRelationshipKey{
			OwnerID:       ownerID,
			RelatedUserID: relatedUserID,
		},
	}
	update := bson.M{
		"$set": bson.M{"bd": blockDate},
	}
	_, err := r.coll.UpdateOne(ctx, filter, update)
	return err
}

func (r *userRelationshipRepository) FindRelatedUserIDs(ctx context.Context, ownerID int64, isBlocked *bool) ([]int64, error) {
	filter := bson.M{"_id.oid": ownerID}
	if isBlocked != nil {
		if *isBlocked {
			filter["bd"] = bson.M{"$ne": nil}
		} else {
			filter["bd"] = nil
		}
	}
	cursor, err := r.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var relatedIDs []int64
	for cursor.Next(ctx) {
		var rel po.UserRelationship
		if err := cursor.Decode(&rel); err != nil {
			return nil, err
		}
		relatedIDs = append(relatedIDs, rel.ID.RelatedUserID)
	}
	return relatedIDs, cursor.Err()
}

func (r *userRelationshipRepository) DeleteById(ctx context.Context, ownerID, relatedUserID int64) error {
	filter := bson.M{
		"_id": po.UserRelationshipKey{
			OwnerID:       ownerID,
			RelatedUserID: relatedUserID,
		},
	}
	_, err := r.coll.DeleteOne(ctx, filter)
	return err
}

func (r *userRelationshipRepository) Upsert(ctx context.Context, ownerID, relatedUserID int64, blockDate *time.Time, groupIndex *int32, establishmentDate *time.Time, name *string, session mongo.SessionContext) error {
	filter := bson.M{
		"_id": po.UserRelationshipKey{
			OwnerID:       ownerID,
			RelatedUserID: relatedUserID,
		},
	}
	update := bson.M{}
	setOps := bson.M{}
	if blockDate != nil {
		setOps["bd"] = blockDate
	}
	if groupIndex != nil {
		setOps["gi"] = groupIndex
	}
	if establishmentDate != nil {
		setOps["ed"] = establishmentDate
	}
	if name != nil {
		setOps["n"] = name
	}
	if len(setOps) > 0 {
		update["$set"] = setOps
	} else {
		update["$setOnInsert"] = bson.M{}
	}

	opts := options.Update().SetUpsert(true)
	_, err := r.coll.UpdateOne(ctx, filter, update, opts)
	return err
}

func (r *userRelationshipRepository) FindRelationships(ctx context.Context, ownerIDs []int64, relatedUserIDs []int64) ([]po.UserRelationship, error) {
	filter := bson.M{}
	if len(ownerIDs) > 0 {
		filter["_id.oid"] = bson.M{"$in": ownerIDs}
	}
	if len(relatedUserIDs) > 0 {
		filter["_id.rid"] = bson.M{"$in": relatedUserIDs}
	}

	cursor, err := r.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var rels []po.UserRelationship
	if err := cursor.All(ctx, &rels); err != nil {
		return nil, err
	}
	return rels, nil
}

func (r *userRelationshipRepository) DeleteAllRelationships(ctx context.Context, ownerIDs []int64) error {
	filter := bson.M{"_id.oid": bson.M{"$in": ownerIDs}}
	_, err := r.coll.DeleteMany(ctx, filter)
	return err
}

func (r *userRelationshipRepository) DeleteOneSidedRelationships(ctx context.Context, ownerIDs, relatedUserIDs []int64) error {
	filter := bson.M{}
	if len(ownerIDs) > 0 {
		filter["_id.oid"] = bson.M{"$in": ownerIDs}
	}
	if len(relatedUserIDs) > 0 {
		filter["_id.rid"] = bson.M{"$in": relatedUserIDs}
	}
	_, err := r.coll.DeleteMany(ctx, filter)
	return err
}

func (r *userRelationshipRepository) countOrFindFilter(ownerIDs, relatedUserIDs []int64, groupIndexes []int32, isBlocked *bool) bson.M {
	filter := bson.M{}
	if len(ownerIDs) > 0 {
		filter["_id.oid"] = bson.M{"$in": ownerIDs}
	}
	if len(relatedUserIDs) > 0 {
		filter["_id.rid"] = bson.M{"$in": relatedUserIDs}
	}
	if len(groupIndexes) > 0 {
		filter["gi"] = bson.M{"$in": groupIndexes}
	}
	if isBlocked != nil {
		if *isBlocked {
			filter["bd"] = bson.M{"$ne": nil}
		} else {
			filter["bd"] = nil
		}
	}
	return filter
}

func (r *userRelationshipRepository) CountRelationships(ctx context.Context, ownerIDs, relatedUserIDs []int64, groupIndexes []int32, isBlocked *bool) (int64, error) {
	filter := r.countOrFindFilter(ownerIDs, relatedUserIDs, groupIndexes, isBlocked)
	return r.coll.CountDocuments(ctx, filter)
}

func (r *userRelationshipRepository) QueryMembersRelationships(ctx context.Context, ownerID int64, groupIndexes []int32, page, size *int) ([]po.UserRelationship, error) {
	filter := r.countOrFindFilter([]int64{ownerID}, nil, groupIndexes, nil)
	opts := options.Find().SetSort(bson.M{"ed": -1})
	if page != nil && size != nil {
		opts.SetSkip(int64(*page * *size))
		opts.SetLimit(int64(*size))
	} else if size != nil {
		opts.SetLimit(int64(*size))
	}

	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var rels []po.UserRelationship
	if err := cursor.All(ctx, &rels); err != nil {
		return nil, err
	}
	return rels, nil
}

func (r *userRelationshipRepository) HasOneSidedRelationship(ctx context.Context, ownerID, relatedUserID int64) (bool, error) {
	filter := bson.M{
		"_id": po.UserRelationshipKey{
			OwnerID:       ownerID,
			RelatedUserID: relatedUserID,
		},
	}
	count, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

