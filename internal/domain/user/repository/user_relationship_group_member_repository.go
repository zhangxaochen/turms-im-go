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
	InsertMember(ctx context.Context, member *po.UserRelationshipGroupMember) error
	DeleteAllRelatedUserFromRelationshipGroup(ctx context.Context, ownerID int64, session *mongo.Session) (int64, error)
	DeleteRelatedUserFromRelationshipGroup(ctx context.Context, ownerID int64, relatedUserID int64, groupIndexes []int32, session *mongo.Session) (int64, error)
	DeleteRelatedUsersFromAllRelationshipGroups(ctx context.Context, ownerID int64, relatedUserIDs []int64, session *mongo.Session) (int64, error)
	CountGroups(ctx context.Context, ownerID int64, relatedUserID int64) (int64, error)
	CountMembers(ctx context.Context, ownerID int64, groupIndex int32) (int64, error)
	FindGroupIndexes(ctx context.Context, ownerID int64, relatedUserID int64) ([]int32, error)
	FindRelationshipGroupMemberIds(ctx context.Context, ownerID int64) ([]*po.UserRelationshipGroupMember, error)
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

func (r *userRelationshipGroupMemberRepository) InsertMember(ctx context.Context, member *po.UserRelationshipGroupMember) error {
	_, err := r.collection.InsertOne(ctx, member)
	return err
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
		"_id.oid":  ownerID,
		"_id.ruid": relatedUserID,
	}
	if len(groupIndexes) > 0 {
		filter["_id.gidx"] = bson.M{"$in": groupIndexes}
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
		"_id.oid":  ownerID,
		"_id.ruid": bson.M{"$in": relatedUserIDs},
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

func (r *userRelationshipGroupMemberRepository) CountGroups(ctx context.Context, ownerID int64, relatedUserID int64) (int64, error) {
	filter := bson.M{
		"_id.oid":  ownerID,
		"_id.ruid": relatedUserID,
	}
	return r.collection.CountDocuments(ctx, filter)
}

func (r *userRelationshipGroupMemberRepository) CountMembers(ctx context.Context, ownerID int64, groupIndex int32) (int64, error) {
	filter := bson.M{
		"_id.oid":  ownerID,
		"_id.gidx": groupIndex,
	}
	return r.collection.CountDocuments(ctx, filter)
}

func (r *userRelationshipGroupMemberRepository) FindGroupIndexes(ctx context.Context, ownerID int64, relatedUserID int64) ([]int32, error) {
	filter := bson.M{
		"_id.oid":  ownerID,
		"_id.ruid": relatedUserID,
	}
	opts := options.Find().SetProjection(bson.M{"_id.gidx": 1})
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

func (r *userRelationshipGroupMemberRepository) FindRelationshipGroupMemberIds(ctx context.Context, ownerID int64) ([]*po.UserRelationshipGroupMember, error) {
	filter := bson.M{
		"_id.oid": ownerID,
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

func (r *userRelationshipGroupMemberRepository) FindRelationshipGroupMembers(ctx context.Context, ownerID int64, groupIndex int32) ([]*po.UserRelationshipGroupMember, error) {
	filter := bson.M{
		"_id.oid":  ownerID,
		"_id.gidx": groupIndex,
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
