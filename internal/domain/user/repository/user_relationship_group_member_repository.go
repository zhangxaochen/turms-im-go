package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"im.turms/server/internal/domain/user/po"
	turmsmongo "im.turms/server/internal/storage/mongo"
)

type UserRelationshipGroupMemberRepository interface {
	Upsert(ctx context.Context, member *po.UserRelationshipGroupMember, session *mongo.Session) (*mongo.UpdateResult, error)
	Insert(ctx context.Context, member *po.UserRelationshipGroupMember, session *mongo.Session) error
	InsertAllOfSameType(ctx context.Context, members []*po.UserRelationshipGroupMember, session *mongo.Session) error
	DeleteById(ctx context.Context, ownerID int64, groupIndex int32, relatedUserID int64, session *mongo.Session) (int64, error)
	DeleteAllRelatedUserFromRelationshipGroup(ctx context.Context, ownerID int64, session *mongo.Session) (int64, error)
	DeleteRelatedUserFromRelationshipGroup(ctx context.Context, ownerID int64, relatedUserID int64, groupIndexes []int32, session *mongo.Session) (int64, error)
	DeleteRelatedUsersFromAllRelationshipGroups(ctx context.Context, ownerID int64, relatedUserIDs []int64, session *mongo.Session) (int64, error)
	CountGroups(ctx context.Context, ownerIDs []int64, relatedUserIDs []int64) (int64, error)
	CountMembers(ctx context.Context, ownerIDs []int64, groupIndexes []int32) (int64, error)
	FindGroupIndexes(ctx context.Context, ownerID int64, relatedUserID int64) ([]int32, error)
	FindRelationshipGroupMemberIds(ctx context.Context, ownerIDs []int64, groupIndexes []int32, page *int, size *int) ([]int64, error)
	FindRelationshipGroupMembers(ctx context.Context, ownerID int64, groupIndex int32) ([]*po.UserRelationshipGroupMember, error)
}

type userRelationshipGroupMemberRepository struct {
	collection *mongo.Collection
}

func NewUserRelationshipGroupMemberRepository(mongoClient *turmsmongo.Client) UserRelationshipGroupMemberRepository {
	return &userRelationshipGroupMemberRepository{
		collection: mongoClient.Collection(po.CollectionNameUserRelationshipGroupMember),
	}
}

func (r *userRelationshipGroupMemberRepository) Upsert(ctx context.Context, member *po.UserRelationshipGroupMember, session *mongo.Session) (*mongo.UpdateResult, error) {
	filter := bson.M{
		"_id.oid": member.Key.OwnerID,
		"_id.gi":  member.Key.GroupIndex,
		"_id.rid": member.Key.RelatedUserID,
	}
	update := bson.M{
		"$setOnInsert": bson.M{"jd": member.JoinDate},
	}
	opts := options.Update().SetUpsert(true)
	var res *mongo.UpdateResult
	var err error
	if session != nil {
		err = mongo.WithSession(ctx, *session, func(sc mongo.SessionContext) error {
			res, err = r.collection.UpdateOne(sc, filter, update, opts)
			return err
		})
	} else {
		res, err = r.collection.UpdateOne(ctx, filter, update, opts)
	}
	return res, err
}

func (r *userRelationshipGroupMemberRepository) Insert(ctx context.Context, member *po.UserRelationshipGroupMember, session *mongo.Session) error {
	var err error
	if session != nil {
		err = mongo.WithSession(ctx, *session, func(sc mongo.SessionContext) error {
			_, err = r.collection.InsertOne(sc, member)
			return err
		})
	} else {
		_, err = r.collection.InsertOne(ctx, member)
	}
	return err
}

func (r *userRelationshipGroupMemberRepository) InsertAllOfSameType(ctx context.Context, members []*po.UserRelationshipGroupMember, session *mongo.Session) error {
	if len(members) == 0 {
		return nil
	}
	docs := make([]interface{}, len(members))
	for i, m := range members {
		docs[i] = m
	}
	var err error
	if session != nil {
		err = mongo.WithSession(ctx, *session, func(sc mongo.SessionContext) error {
			_, err = r.collection.InsertMany(sc, docs)
			return err
		})
	} else {
		_, err = r.collection.InsertMany(ctx, docs)
	}
	return err
}

func (r *userRelationshipGroupMemberRepository) DeleteById(ctx context.Context, ownerID int64, groupIndex int32, relatedUserID int64, session *mongo.Session) (int64, error) {
	filter := bson.M{
		"_id.oid": ownerID,
		"_id.gi":  groupIndex,
		"_id.rid": relatedUserID,
	}
	var res *mongo.DeleteResult
	var err error
	if session != nil {
		err = mongo.WithSession(ctx, *session, func(sc mongo.SessionContext) error {
			res, err = r.collection.DeleteOne(sc, filter)
			return err
		})
	} else {
		res, err = r.collection.DeleteOne(ctx, filter)
	}
	if err != nil {
		return 0, err
	}
	return res.DeletedCount, nil
}

func (r *userRelationshipGroupMemberRepository) DeleteAllRelatedUserFromRelationshipGroup(ctx context.Context, ownerID int64, session *mongo.Session) (int64, error) {
	filter := bson.M{
		"_id.oid": ownerID,
	}
	var res *mongo.DeleteResult
	var err error
	if session != nil {
		err = mongo.WithSession(ctx, *session, func(sc mongo.SessionContext) error {
			res, err = r.collection.DeleteMany(sc, filter)
			return err
		})
	} else {
		res, err = r.collection.DeleteMany(ctx, filter)
	}
	if err != nil {
		return 0, err
	}
	return res.DeletedCount, nil
}

func (r *userRelationshipGroupMemberRepository) DeleteRelatedUserFromRelationshipGroup(ctx context.Context, ownerID int64, relatedUserID int64, groupIndexes []int32, session *mongo.Session) (int64, error) {
	filter := bson.M{
		"_id.oid": ownerID,
	}
	if relatedUserID > 0 {
		filter["_id.rid"] = relatedUserID
	}
	if len(groupIndexes) > 0 {
		filter["_id.gi"] = bson.M{"$in": groupIndexes}
	}
	var res *mongo.DeleteResult
	var err error
	if session != nil {
		err = mongo.WithSession(ctx, *session, func(sc mongo.SessionContext) error {
			res, err = r.collection.DeleteMany(sc, filter)
			return err
		})
	} else {
		res, err = r.collection.DeleteMany(ctx, filter)
	}
	if err != nil {
		return 0, err
	}
	return res.DeletedCount, nil
}

func (r *userRelationshipGroupMemberRepository) DeleteRelatedUsersFromAllRelationshipGroups(ctx context.Context, ownerID int64, relatedUserIDs []int64, session *mongo.Session) (int64, error) {
	if len(relatedUserIDs) == 0 {
		return 0, nil
	}
	filter := bson.M{
		"_id.oid": ownerID,
		"_id.rid": bson.M{"$in": relatedUserIDs},
	}
	var res *mongo.DeleteResult
	var err error
	if session != nil {
		err = mongo.WithSession(ctx, *session, func(sc mongo.SessionContext) error {
			res, err = r.collection.DeleteMany(sc, filter)
			return err
		})
	} else {
		res, err = r.collection.DeleteMany(ctx, filter)
	}
	if err != nil {
		return 0, err
	}
	return res.DeletedCount, nil
}

func (r *userRelationshipGroupMemberRepository) CountGroups(ctx context.Context, ownerIDs []int64, relatedUserIDs []int64) (int64, error) {
	filter := bson.M{}
	if len(ownerIDs) > 0 {
		filter["_id.oid"] = bson.M{"$in": ownerIDs}
	}
	if len(relatedUserIDs) > 0 {
		filter["_id.rid"] = bson.M{"$in": relatedUserIDs}
	}
	return r.collection.CountDocuments(ctx, filter)
}

func (r *userRelationshipGroupMemberRepository) CountMembers(ctx context.Context, ownerIDs []int64, groupIndexes []int32) (int64, error) {
	filter := bson.M{}
	if len(ownerIDs) > 0 {
		filter["_id.oid"] = bson.M{"$in": ownerIDs}
	}
	if len(groupIndexes) > 0 {
		filter["_id.gi"] = bson.M{"$in": groupIndexes}
	}
	return r.collection.CountDocuments(ctx, filter)
}

func (r *userRelationshipGroupMemberRepository) FindGroupIndexes(ctx context.Context, ownerID int64, relatedUserID int64) ([]int32, error) {
	filter := bson.M{
		"_id.oid": ownerID,
		"_id.rid": relatedUserID,
	}
	opts := options.Find().SetProjection(bson.M{"_id.gi": 1})
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var members []po.UserRelationshipGroupMember
	if err = cursor.All(ctx, &members); err != nil {
		return nil, err
	}
	var indexes []int32
	for _, m := range members {
		indexes = append(indexes, m.Key.GroupIndex)
	}
	return indexes, nil
}

func (r *userRelationshipGroupMemberRepository) FindRelationshipGroupMemberIds(ctx context.Context, ownerIDs []int64, groupIndexes []int32, page *int, size *int) ([]int64, error) {
	filter := bson.M{}
	if len(ownerIDs) > 0 {
		filter["_id.oid"] = bson.M{"$in": ownerIDs}
	}
	if len(groupIndexes) > 0 {
		filter["_id.gi"] = bson.M{"$in": groupIndexes}
	}

	opts := options.Find().SetProjection(bson.M{"_id.rid": 1})
	if page != nil && size != nil {
		opts.SetSkip(int64(*page * *size))
		opts.SetLimit(int64(*size))
	} else if size != nil {
		opts.SetLimit(int64(*size))
	}

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var members []po.UserRelationshipGroupMember
	if err = cursor.All(ctx, &members); err != nil {
		return nil, err
	}
	ids := make([]int64, len(members))
	for i, m := range members {
		ids[i] = m.Key.RelatedUserID
	}
	return ids, nil
}

func (r *userRelationshipGroupMemberRepository) FindRelationshipGroupMembers(ctx context.Context, ownerID int64, groupIndex int32) ([]*po.UserRelationshipGroupMember, error) {
	filter := bson.M{
		"_id.oid": ownerID,
		"_id.gi":  groupIndex,
	}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var members []*po.UserRelationshipGroupMember
	if err = cursor.All(ctx, &members); err != nil {
		return nil, err
	}
	return members, nil
}
