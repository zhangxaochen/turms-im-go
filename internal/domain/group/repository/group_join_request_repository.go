package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"im.turms/server/internal/domain/group/po"
	turmsmongo "im.turms/server/internal/storage/mongo"
)

type GroupJoinRequestRepository interface {
	Insert(ctx context.Context, req *po.GroupJoinRequest) error
	HasPendingJoinRequest(ctx context.Context, requesterID, groupID int64) (bool, error)
	UpdateStatusIfPending(ctx context.Context, requestID, responderID int64, newStatus po.RequestStatus, reason *string, responseDate time.Time) (bool, error)
	FindRequestsByGroupID(ctx context.Context, groupID int64) ([]po.GroupJoinRequest, error)
	FindRequestsByRequesterID(ctx context.Context, requesterID int64) ([]po.GroupJoinRequest, error)
}

type groupJoinRequestRepository struct {
	coll *mongo.Collection
}

func NewGroupJoinRequestRepository(client *turmsmongo.Client) GroupJoinRequestRepository {
	return &groupJoinRequestRepository{
		coll: client.Collection(po.CollectionNameGroupJoinRequest),
	}
}

func (r *groupJoinRequestRepository) Insert(ctx context.Context, req *po.GroupJoinRequest) error {
	_, err := r.coll.InsertOne(ctx, req)
	return err
}

func (r *groupJoinRequestRepository) HasPendingJoinRequest(ctx context.Context, requesterID, groupID int64) (bool, error) {
	filter := bson.M{
		"rqid": requesterID,
		"gid":  groupID,
		"stat": po.RequestStatusPending,
	}
	count, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *groupJoinRequestRepository) UpdateStatusIfPending(ctx context.Context, requestID, responderID int64, newStatus po.RequestStatus, reason *string, responseDate time.Time) (bool, error) {
	filter := bson.M{
		"_id":  requestID,
		"stat": po.RequestStatusPending,
	}
	updateOps := bson.M{
		"stat": newStatus,
		"rd":   responseDate,
		"rpid": responderID,
	}
	if reason != nil {
		updateOps["rsn"] = *reason
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

func (r *groupJoinRequestRepository) FindRequestsByGroupID(ctx context.Context, groupID int64) ([]po.GroupJoinRequest, error) {
	filter := bson.M{"gid": groupID}
	opts := options.Find().SetSort(bson.M{"cd": -1})
	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var reqs []po.GroupJoinRequest
	if err := cursor.All(ctx, &reqs); err != nil {
		return nil, err
	}
	return reqs, nil
}

func (r *groupJoinRequestRepository) FindRequestsByRequesterID(ctx context.Context, requesterID int64) ([]po.GroupJoinRequest, error) {
	filter := bson.M{"rqid": requesterID}
	opts := options.Find().SetSort(bson.M{"cd": -1})
	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var reqs []po.GroupJoinRequest
	if err := cursor.All(ctx, &reqs); err != nil {
		return nil, err
	}
	return reqs, nil
}
