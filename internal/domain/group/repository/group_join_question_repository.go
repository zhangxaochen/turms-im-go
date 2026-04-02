package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"im.turms/server/internal/domain/group/po"
	turmsmongo "im.turms/server/internal/storage/mongo"
)

type GroupJoinQuestionRepository interface {
	Insert(ctx context.Context, question *po.GroupJoinQuestion) error
	Delete(ctx context.Context, questionID int64) error
	FindQuestionsByGroupID(ctx context.Context, groupID int64) ([]po.GroupJoinQuestion, error)
	FindByID(ctx context.Context, questionID int64) (*po.GroupJoinQuestion, error)
	Update(ctx context.Context, questionID int64, question *string, answers []string, score *int) (bool, error)
}

type groupJoinQuestionRepository struct {
	coll *mongo.Collection
}

func NewGroupJoinQuestionRepository(client *turmsmongo.Client) GroupJoinQuestionRepository {
	return &groupJoinQuestionRepository{
		coll: client.Collection(po.CollectionNameGroupJoinQuestion),
	}
}

func (r *groupJoinQuestionRepository) Insert(ctx context.Context, question *po.GroupJoinQuestion) error {
	_, err := r.coll.InsertOne(ctx, question)
	return err
}

func (r *groupJoinQuestionRepository) Delete(ctx context.Context, questionID int64) error {
	filter := bson.M{"_id": questionID}
	_, err := r.coll.DeleteOne(ctx, filter)
	return err
}

func (r *groupJoinQuestionRepository) FindQuestionsByGroupID(ctx context.Context, groupID int64) ([]po.GroupJoinQuestion, error) {
	filter := bson.M{"gid": groupID}
	cursor, err := r.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var questions []po.GroupJoinQuestion
	if err := cursor.All(ctx, &questions); err != nil {
		return nil, err
	}
	return questions, nil
}
func (r *groupJoinQuestionRepository) FindByID(ctx context.Context, questionID int64) (*po.GroupJoinQuestion, error) {
	filter := bson.M{"_id": questionID}
	var res po.GroupJoinQuestion
	err := r.coll.FindOne(ctx, filter).Decode(&res)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &res, err
}

func (r *groupJoinQuestionRepository) Update(ctx context.Context, questionID int64, question *string, answers []string, score *int) (bool, error) {
	filter := bson.M{"_id": questionID}
	updateOps := bson.M{}
	if question != nil {
		updateOps["q"] = *question
	}
	if answers != nil {
		updateOps["ans"] = answers
	}
	if score != nil {
		updateOps["score"] = *score
	}
	if len(updateOps) == 0 {
		return true, nil
	}
	update := bson.M{"$set": updateOps}
	res, err := r.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return false, err
	}
	return res.ModifiedCount > 0, nil
}
