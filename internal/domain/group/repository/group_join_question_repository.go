package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"im.turms/server/internal/domain/group/po"
	turmsmongo "im.turms/server/internal/storage/mongo"
)

type GroupJoinQuestionRepository interface {
	Insert(ctx context.Context, question *po.GroupJoinQuestion) error
	Delete(ctx context.Context, questionID int64) error
	FindQuestionsByGroupID(ctx context.Context, groupID int64) ([]po.GroupJoinQuestion, error)
	FindByID(ctx context.Context, questionID int64) (*po.GroupJoinQuestion, error)
	Update(ctx context.Context, questionID int64, question *string, answers []string, score *int) (bool, error)
	FindQuestions(ctx context.Context, ids []int64, groupIds []int64, page *int, size *int) ([]po.GroupJoinQuestion, error)
	CountQuestions(ctx context.Context, ids []int64, groupIds []int64) (int64, error)
	UpdateQuestions(ctx context.Context, ids []int64, groupID *int64, question *string, answers []string, score *int) error
	CheckQuestionAnswerAndGetScore(ctx context.Context, questionID int64, answer string, groupID *int64) (*int, error)
	FindGroupId(ctx context.Context, questionID int64) (*int64, error)
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

func (r *groupJoinQuestionRepository) FindQuestions(ctx context.Context, ids []int64, groupIds []int64, page *int, size *int) ([]po.GroupJoinQuestion, error) {
	filter := bson.M{}
	if len(ids) > 0 {
		filter["_id"] = bson.M{"$in": ids}
	}
	if len(groupIds) > 0 {
		filter["gid"] = bson.M{"$in": groupIds}
	}

	opts := options.Find()
	if page != nil {
		opts.SetSkip(int64(*page))
	}
	if size != nil {
		opts.SetLimit(int64(*size))
	}

	cursor, err := r.coll.Find(ctx, filter, opts)
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

func (r *groupJoinQuestionRepository) CountQuestions(ctx context.Context, ids []int64, groupIds []int64) (int64, error) {
	filter := bson.M{}
	if len(ids) > 0 {
		filter["_id"] = bson.M{"$in": ids}
	}
	if len(groupIds) > 0 {
		filter["gid"] = bson.M{"$in": groupIds}
	}

	return r.coll.CountDocuments(ctx, filter)
}

func (r *groupJoinQuestionRepository) UpdateQuestions(ctx context.Context, ids []int64, groupID *int64, question *string, answers []string, score *int) error {
	filter := bson.M{}
	if len(ids) > 0 {
		filter["_id"] = bson.M{"$in": ids}
	}

	updateOps := bson.M{}
	if groupID != nil {
		updateOps["gid"] = *groupID
	}
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
		return nil
	}

	update := bson.M{"$set": updateOps}
	_, err := r.coll.UpdateMany(ctx, filter, update)
	return err
}

func (r *groupJoinQuestionRepository) CheckQuestionAnswerAndGetScore(ctx context.Context, questionID int64, answer string, groupID *int64) (*int, error) {
	filter := bson.M{
		"_id": questionID,
		"ans": answer,
	}
	if groupID != nil {
		filter["gid"] = *groupID
	}

	opts := options.FindOne().SetProjection(bson.M{"score": 1})
	var result struct {
		Score int `bson:"score"`
	}
	err := r.coll.FindOne(ctx, filter, opts).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &result.Score, nil
}

func (r *groupJoinQuestionRepository) FindGroupId(ctx context.Context, questionID int64) (*int64, error) {
	filter := bson.M{"_id": questionID}
	opts := options.FindOne().SetProjection(bson.M{"gid": 1})
	var result struct {
		GroupID int64 `bson:"gid"`
	}
	err := r.coll.FindOne(ctx, filter, opts).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &result.GroupID, nil
}
