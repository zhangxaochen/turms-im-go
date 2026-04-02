package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"im.turms/server/internal/domain/user/po"
	turmsmongo "im.turms/server/internal/storage/mongo"
)

type UserRelationshipGroupMemberRepository interface {
	InsertMember(ctx context.Context, member *po.UserRelationshipGroupMember) error
	FindMembers(ctx context.Context, filter interface{}) ([]*po.UserRelationshipGroupMember, error)
	DeleteMembers(ctx context.Context, filter interface{}) (int64, error)
	UpsertMember(ctx context.Context, member *po.UserRelationshipGroupMember) error
	UpdateMembers(ctx context.Context, filter interface{}, update interface{}) (int64, error)
	CountMembers(ctx context.Context, filter interface{}) (int64, error)
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

func (r *userRelationshipGroupMemberRepository) FindMembers(ctx context.Context, filter interface{}) ([]*po.UserRelationshipGroupMember, error) {
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var members []*po.UserRelationshipGroupMember
	if err = cursor.All(ctx, &members); err != nil {
		return nil, err
	}
	return members, nil
}

func (r *userRelationshipGroupMemberRepository) DeleteMembers(ctx context.Context, filter interface{}) (int64, error) {
	res, err := r.collection.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}
	return res.DeletedCount, nil
}

func (r *userRelationshipGroupMemberRepository) UpsertMember(ctx context.Context, member *po.UserRelationshipGroupMember) error {
	opts := options.Update().SetUpsert(true)
	update := map[string]interface{}{
		"$set": member,
	}
	_, err := r.collection.UpdateOne(ctx, map[string]interface{}{"_id": member.Key}, update, opts)
	return err
}

func (r *userRelationshipGroupMemberRepository) UpdateMembers(ctx context.Context, filter interface{}, update interface{}) (int64, error) {
	res, err := r.collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return 0, err
	}
	return res.ModifiedCount, nil
}

func (r *userRelationshipGroupMemberRepository) CountMembers(ctx context.Context, filter interface{}) (int64, error) {
	return r.collection.CountDocuments(ctx, filter)
}
