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

type UserFriendRequestRepository interface {
	Insert(ctx context.Context, req *po.UserFriendRequest) error
	HasPendingFriendRequest(ctx context.Context, requesterID, recipientID int64) (bool, error)
	UpdateStatusIfPending(ctx context.Context, requestID, recipientID int64, newStatus po.RequestStatus, reason *string, responseDate time.Time) (bool, error)
	FindRequestsByRecipientID(ctx context.Context, recipientID int64) ([]po.UserFriendRequest, error)
	FindRequestsByRequesterID(ctx context.Context, requesterID int64) ([]po.UserFriendRequest, error)
}

type userFriendRequestRepository struct {
	coll *mongo.Collection
}

func NewUserFriendRequestRepository(client *turmsmongo.Client) UserFriendRequestRepository {
	return &userFriendRequestRepository{
		coll: client.Collection(po.CollectionNameUserFriendRequest),
	}
}

func (r *userFriendRequestRepository) Insert(ctx context.Context, req *po.UserFriendRequest) error {
	_, err := r.coll.InsertOne(ctx, req)
	return err
}

func (r *userFriendRequestRepository) HasPendingFriendRequest(ctx context.Context, requesterID, recipientID int64) (bool, error) {
	filter := bson.M{
		"rqid": requesterID,
		"rcid": recipientID,
		"s":    po.RequestStatusPending,
	}
	count, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *userFriendRequestRepository) UpdateStatusIfPending(ctx context.Context, requestID, recipientID int64, newStatus po.RequestStatus, reason *string, responseDate time.Time) (bool, error) {
	filter := bson.M{
		"_id":  requestID,
		"rcid": recipientID,
		"s":    po.RequestStatusPending,
	}
	updateOps := bson.M{
		"s":  newStatus,
		"rd": responseDate,
	}
	if reason != nil {
		updateOps["r"] = *reason
	}
	update := bson.M{
		"$set": updateOps,
	}

	res, err := r.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return false, err
	}
	return res.ModifiedCount > 0, nil
}

func (r *userFriendRequestRepository) FindRequestsByRecipientID(ctx context.Context, recipientID int64) ([]po.UserFriendRequest, error) {
	filter := bson.M{"rcid": recipientID}
	opts := options.Find().SetSort(bson.M{"cd": -1})
	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var reqs []po.UserFriendRequest
	if err := cursor.All(ctx, &reqs); err != nil {
		return nil, err
	}
	return reqs, nil
}

func (r *userFriendRequestRepository) FindRequestsByRequesterID(ctx context.Context, requesterID int64) ([]po.UserFriendRequest, error) {
	filter := bson.M{"rqid": requesterID}
	opts := options.Find().SetSort(bson.M{"cd": -1})
	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var reqs []po.UserFriendRequest
	if err := cursor.All(ctx, &reqs); err != nil {
		return nil, err
	}
	return reqs, nil
}
