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
	HasRelationshipAndNotBlocked(ctx context.Context, ownerID, relatedUserID int64, session *mongo.Session) (bool, error)
	Insert(ctx context.Context, rel *po.UserRelationship, session *mongo.Session) error
	UpdateBlockDate(ctx context.Context, ownerID, relatedUserID int64, blockDate *time.Time, session *mongo.Session) error
	FindRelatedUserIDs(ctx context.Context, ownerIDs []int64, groupIndexes []int32, isBlocked *bool, page, size *int, session *mongo.Session) ([]int64, error)
	FindRelationships(ctx context.Context, ownerIDs []int64, relatedUserIDs []int64, groupIndexes []int32, isBlocked *bool, establishmentDateRange *turmsmongo.DateRange, page, size *int, session *mongo.Session) ([]po.UserRelationship, error)
	DeleteById(ctx context.Context, ownerID, relatedUserID int64, session *mongo.Session) (*mongo.DeleteResult, error)
	DeleteByIds(ctx context.Context, keys []po.UserRelationshipKey, session *mongo.Session) (*mongo.DeleteResult, error)
	Upsert(ctx context.Context, ownerID, relatedUserID int64, blockDate, establishmentDate *time.Time, name *string, session *mongo.Session) (*mongo.UpdateResult, error)
	DeleteAllRelationships(ctx context.Context, ownerIDs []int64, session *mongo.Session) (*mongo.DeleteResult, error)
	DeleteOneSidedRelationships(ctx context.Context, ownerIDs, relatedUserIDs []int64, session *mongo.Session) (*mongo.DeleteResult, error)
	UpdateUserOneSidedRelationships(ctx context.Context, ownerID int64, relatedUserIDs []int64, blockDate *time.Time, establishmentDate *time.Time, name *string, session *mongo.Session) (int64, error)
	CountRelationships(ctx context.Context, ownerIDs, relatedUserIDs []int64, groupIndexes []int32, isBlocked *bool, session *mongo.Session) (int64, error)
	QueryMembersRelationships(ctx context.Context, ownerID int64, groupIndexes []int32, page, size *int, session *mongo.Session) ([]po.UserRelationship, error)
	HasOneSidedRelationship(ctx context.Context, ownerID, relatedUserID int64, session *mongo.Session) (bool, error)
	IsBlocked(ctx context.Context, ownerID, relatedUserID int64, session *mongo.Session) (bool, error)
}

type userRelationshipRepository struct {
	client *turmsmongo.Client
	coll   *mongo.Collection
}

func NewUserRelationshipRepository(client *turmsmongo.Client) UserRelationshipRepository {
	return &userRelationshipRepository{
		client: client,
		coll:   client.Collection(po.CollectionNameUserRelationship),
	}
}

// @MappedFrom hasRelationshipAndNotBlocked(Long ownerId, Long relatedUserId)
// @MappedFrom hasRelationshipAndNotBlocked(@NotNull Long ownerId, @NotNull Long relatedUserId)
// @MappedFrom hasRelationshipAndNotBlocked(@NotNull Long ownerId, @NotNull Long relatedUserId, boolean preferCache)
func (r *userRelationshipRepository) HasRelationshipAndNotBlocked(ctx context.Context, ownerID, relatedUserID int64, session *mongo.Session) (bool, error) {
	filter := bson.M{
		"_id": po.UserRelationshipKey{
			OwnerID:       ownerID,
			RelatedUserID: relatedUserID,
		},
		"bd": nil,
	}

	return turmsmongo.ExecuteWithSessionResult(ctx, r.client, session, func(sessCtx mongo.SessionContext, sess *mongo.Session) (bool, error) {
		count, err := r.coll.CountDocuments(sessCtx, filter)
		return count > 0, err
	})
}

func (r *userRelationshipRepository) IsBlocked(ctx context.Context, ownerID, relatedUserID int64, session *mongo.Session) (bool, error) {
	filter := bson.M{
		"_id": po.UserRelationshipKey{
			OwnerID:       ownerID,
			RelatedUserID: relatedUserID,
		},
		"bd": bson.M{"$ne": nil},
	}

	return turmsmongo.ExecuteWithSessionResult(ctx, r.client, session, func(sessCtx mongo.SessionContext, sess *mongo.Session) (bool, error) {
		count, err := r.coll.CountDocuments(sessCtx, filter)
		return count > 0, err
	})
}

func (r *userRelationshipRepository) Insert(ctx context.Context, rel *po.UserRelationship, session *mongo.Session) error {
	return turmsmongo.ExecuteWithSession(ctx, r.client, session, func(sessCtx mongo.SessionContext, sess *mongo.Session) error {
		_, err := r.coll.InsertOne(sessCtx, rel)
		return err
	})
}

func (r *userRelationshipRepository) UpdateBlockDate(ctx context.Context, ownerID, relatedUserID int64, blockDate *time.Time, session *mongo.Session) error {
	filter := bson.M{
		"_id": po.UserRelationshipKey{
			OwnerID:       ownerID,
			RelatedUserID: relatedUserID,
		},
	}
	update := bson.M{
		"$set": bson.M{"bd": blockDate},
	}
	return turmsmongo.ExecuteWithSession(ctx, r.client, session, func(sessCtx mongo.SessionContext, sess *mongo.Session) error {
		_, err := r.coll.UpdateOne(sessCtx, filter, update)
		return err
	})
}

// @MappedFrom findRelatedUserIds(@Nullable Set<Long> ownerIds, @Nullable Boolean isBlocked)
func (r *userRelationshipRepository) FindRelatedUserIDs(ctx context.Context, ownerIDs []int64, groupIndexes []int32, isBlocked *bool, page, size *int, session *mongo.Session) ([]int64, error) {
	filter := r.countOrFindFilter(ownerIDs, nil, groupIndexes, isBlocked)
	opts := options.Find().SetProjection(bson.M{"_id.rid": 1})
	if page != nil && size != nil {
		opts.SetSkip(int64(*page * *size)).SetLimit(int64(*size))
	} else if size != nil {
		opts.SetLimit(int64(*size))
	}

	return turmsmongo.ExecuteWithSessionResult(ctx, r.client, session, func(sessCtx mongo.SessionContext, sess *mongo.Session) ([]int64, error) {
		cursor, err := r.coll.Find(sessCtx, filter, opts)
		if err != nil {
			return nil, err
		}
		defer cursor.Close(sessCtx)

		var results []int64
		for cursor.Next(sessCtx) {
			var rel po.UserRelationship
			if err := cursor.Decode(&rel); err != nil {
				return nil, err
			}
			results = append(results, rel.ID.RelatedUserID)
		}
		return results, cursor.Err()
	})
}

func (r *userRelationshipRepository) DeleteById(ctx context.Context, ownerID, relatedUserID int64, session *mongo.Session) (*mongo.DeleteResult, error) {
	filter := bson.M{
		"_id": po.UserRelationshipKey{
			OwnerID:       ownerID,
			RelatedUserID: relatedUserID,
		},
	}
	return turmsmongo.ExecuteWithSessionResult(ctx, r.client, session, func(sessCtx mongo.SessionContext, sess *mongo.Session) (*mongo.DeleteResult, error) {
		return r.coll.DeleteOne(sessCtx, filter)
	})
}

func (r *userRelationshipRepository) DeleteByIds(ctx context.Context, keys []po.UserRelationshipKey, session *mongo.Session) (*mongo.DeleteResult, error) {
	if len(keys) == 0 {
		return &mongo.DeleteResult{}, nil
	}
	filter := bson.M{"_id": bson.M{"$in": keys}}
	return turmsmongo.ExecuteWithSessionResult(ctx, r.client, session, func(sessCtx mongo.SessionContext, sess *mongo.Session) (*mongo.DeleteResult, error) {
		return r.coll.DeleteMany(sessCtx, filter)
	})
}

func (r *userRelationshipRepository) Upsert(
	ctx context.Context,
	ownerID, relatedUserID int64,
	blockDate, establishmentDate *time.Time,
	name *string,
	session *mongo.Session,
) (*mongo.UpdateResult, error) {
	filter := bson.M{
		"_id": po.UserRelationshipKey{
			OwnerID:       ownerID,
			RelatedUserID: relatedUserID,
		},
	}
	update := bson.M{
		"$set": bson.M{
			"bd": blockDate,
			"ed": establishmentDate,
			"n":  name,
		},
	}
	opts := options.Update().SetUpsert(true)
	return turmsmongo.ExecuteWithSessionResult(ctx, r.client, session, func(sessCtx mongo.SessionContext, sess *mongo.Session) (*mongo.UpdateResult, error) {
		return r.coll.UpdateOne(sessCtx, filter, update, opts)
	})
}

// @MappedFrom findRelationships(@Nullable Set<Long> ownerIds, @Nullable Set<Long> relatedUserIds, @Nullable Integer page, @Nullable Integer size)
// @MappedFrom findRelationships(Long userId)
// @MappedFrom findRelationships(@Nullable Set<Long> ownerIds, @Nullable Set<Long> relatedUserIds, @Nullable Boolean isBlocked, @Nullable DateRange establishmentDateRange, @Nullable Integer page, @Nullable Integer size)
func (r *userRelationshipRepository) FindRelationships(
	ctx context.Context,
	ownerIDs []int64,
	relatedUserIDs []int64,
	groupIndexes []int32,
	isBlocked *bool,
	establishmentDateRange *turmsmongo.DateRange,
	page, size *int,
	session *mongo.Session,
) ([]po.UserRelationship, error) {
	filter := r.countOrFindFilter(ownerIDs, relatedUserIDs, groupIndexes, isBlocked)
	if establishmentDateRange != nil {
		if edFilter := establishmentDateRange.ToBson(); edFilter != nil {
			filter["ed"] = edFilter
		}
	}

	opts := options.Find()
	if page != nil && size != nil {
		opts.SetSkip(int64(*page * *size)).SetLimit(int64(*size))
	} else if size != nil {
		opts.SetLimit(int64(*size))
	}

	return turmsmongo.ExecuteWithSessionResult(ctx, r.client, session, func(sessCtx mongo.SessionContext, sess *mongo.Session) ([]po.UserRelationship, error) {
		cursor, err := r.coll.Find(sessCtx, filter, opts)
		if err != nil {
			return nil, err
		}
		defer cursor.Close(sessCtx)

		var rels []po.UserRelationship
		if err := cursor.All(sessCtx, &rels); err != nil {
			return nil, err
		}
		return rels, nil
	})
}

// @MappedFrom deleteAllRelationships(Set<Long> userIds, @Nullable ClientSession session)
// @MappedFrom deleteAllRelationships(@NotEmpty Set<Long> userIds, @Nullable ClientSession session, boolean updateRelationshipsVersion)
func (r *userRelationshipRepository) DeleteAllRelationships(ctx context.Context, ownerIDs []int64, session *mongo.Session) (*mongo.DeleteResult, error) {
	if len(ownerIDs) == 0 {
		return &mongo.DeleteResult{}, nil
	}
	filter := bson.M{
		"$or": []bson.M{
			{"_id.oid": bson.M{"$in": ownerIDs}},
			{"_id.rid": bson.M{"$in": ownerIDs}},
		},
	}
	return turmsmongo.ExecuteWithSessionResult(ctx, r.client, session, func(sessCtx mongo.SessionContext, sess *mongo.Session) (*mongo.DeleteResult, error) {
		return r.coll.DeleteMany(sessCtx, filter)
	})
}

// @MappedFrom deleteOneSidedRelationships(@NotEmpty Set<UserRelationship.@ValidUserRelationshipKey Key> keys)
func (r *userRelationshipRepository) DeleteOneSidedRelationships(ctx context.Context, ownerIDs, relatedUserIDs []int64, session *mongo.Session) (*mongo.DeleteResult, error) {
	filter := bson.M{}
	if len(ownerIDs) > 0 {
		filter["_id.oid"] = bson.M{"$in": ownerIDs}
	}
	if len(relatedUserIDs) > 0 {
		filter["_id.rid"] = bson.M{"$in": relatedUserIDs}
	}
	return turmsmongo.ExecuteWithSessionResult(ctx, r.client, session, func(sessCtx mongo.SessionContext, sess *mongo.Session) (*mongo.DeleteResult, error) {
		return r.coll.DeleteMany(sessCtx, filter)
	})
}

// @MappedFrom updateUserOneSidedRelationships(Set<UserRelationship.Key> keys, @Nullable String name, @Nullable Date blockDate, @Nullable Date establishmentDate)
// @MappedFrom updateUserOneSidedRelationships(@NotEmpty Set<UserRelationship.@ValidUserRelationshipKey Key> keys, @Nullable String name, @Nullable @PastOrPresent Date blockDate, @Nullable @PastOrPresent Date establishmentDate)
func (r *userRelationshipRepository) UpdateUserOneSidedRelationships(
	ctx context.Context,
	ownerID int64,
	relatedUserIDs []int64,
	blockDate *time.Time,
	establishmentDate *time.Time,
	name *string,
	session *mongo.Session,
) (int64, error) {
	if len(relatedUserIDs) == 0 {
		return 0, nil
	}
	filter := bson.M{
		"_id.oid": ownerID,
		"_id.rid": bson.M{"$in": relatedUserIDs},
	}
	setOps := bson.M{}
	unsetOps := bson.M{}
	if blockDate != nil {
		if blockDate.UnixMilli() <= 0 {
			unsetOps["bd"] = ""
		} else {
			setOps["bd"] = blockDate
		}
	}
	if establishmentDate != nil {
		setOps["ed"] = establishmentDate
	}
	if name != nil {
		setOps["n"] = name
	}

	if len(setOps) == 0 && len(unsetOps) == 0 {
		return 0, nil
	}

	update := bson.M{}
	if len(setOps) > 0 {
		update["$set"] = setOps
	}
	if len(unsetOps) > 0 {
		update["$unset"] = unsetOps
	}

	return turmsmongo.ExecuteWithSessionResult(ctx, r.client, session, func(sessCtx mongo.SessionContext, sess *mongo.Session) (int64, error) {
		res, err := r.coll.UpdateMany(sessCtx, filter, update)
		if err != nil {
			return 0, err
		}
		return res.ModifiedCount, nil
	})
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

// @MappedFrom countRelationships(@Nullable Set<Long> ownerIds, @Nullable Set<Long> relatedUserIds, @Nullable Set<Integer> groupIndexes, @Nullable Boolean isBlocked)
// @MappedFrom countRelationships(@Nullable Set<Long> ownerIds, @Nullable Set<Long> relatedUserIds, @Nullable Boolean isBlocked)
func (r *userRelationshipRepository) CountRelationships(ctx context.Context, ownerIDs, relatedUserIDs []int64, groupIndexes []int32, isBlocked *bool, session *mongo.Session) (int64, error) {
	filter := r.countOrFindFilter(ownerIDs, relatedUserIDs, groupIndexes, isBlocked)
	return turmsmongo.ExecuteWithSessionResult(ctx, r.client, session, func(sessCtx mongo.SessionContext, sess *mongo.Session) (int64, error) {
		return r.coll.CountDocuments(sessCtx, filter)
	})
}

// @MappedFrom queryMembersRelationships(@Nullable Set<Long> ownerIds, @Nullable Set<Integer> groupIndexes, @Nullable Integer page, @Nullable Integer size)
func (r *userRelationshipRepository) QueryMembersRelationships(ctx context.Context, ownerID int64, groupIndexes []int32, page, size *int, session *mongo.Session) ([]po.UserRelationship, error) {
	filter := r.countOrFindFilter([]int64{ownerID}, nil, groupIndexes, nil)
	opts := options.Find().SetSort(bson.M{"ed": -1})
	if page != nil && size != nil {
		opts.SetSkip(int64(*page * *size)).SetLimit(int64(*size))
	} else if size != nil {
		opts.SetLimit(int64(*size))
	}

	return turmsmongo.ExecuteWithSessionResult(ctx, r.client, session, func(sessCtx mongo.SessionContext, sess *mongo.Session) ([]po.UserRelationship, error) {
		cursor, err := r.coll.Find(sessCtx, filter, opts)
		if err != nil {
			return nil, err
		}
		defer cursor.Close(sessCtx)

		var rels []po.UserRelationship
		if err := cursor.All(sessCtx, &rels); err != nil {
			return nil, err
		}
		return rels, nil
	})
}

// @MappedFrom hasOneSidedRelationship(@NotNull Long ownerId, @NotNull Long relatedUserId)
func (r *userRelationshipRepository) HasOneSidedRelationship(ctx context.Context, ownerID, relatedUserID int64, session *mongo.Session) (bool, error) {
	filter := bson.M{
		"_id": po.UserRelationshipKey{
			OwnerID:       ownerID,
			RelatedUserID: relatedUserID,
		},
	}
	return turmsmongo.ExecuteWithSessionResult(ctx, r.client, session, func(sessCtx mongo.SessionContext, sess *mongo.Session) (bool, error) {
		count, err := r.coll.CountDocuments(sessCtx, filter)
		return count > 0, err
	})
}
