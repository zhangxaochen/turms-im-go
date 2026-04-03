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

// @MappedFrom GroupJoinRequestRepository
type GroupJoinRequestRepository interface {
	UpdateRequests(ctx context.Context, requestIds []int64, requesterId, responderId *int64, content *string, status any, creationDate, responseDate *time.Time) error
	CountRequests(ctx context.Context, ids, groupIds, requesterIds, responderIds []int64, statuses []int, creationDateRange, responseDateRange, expirationDateRange any) (int64, error)
	FindGroupId(ctx context.Context, requestId int64) (int64, error)
	FindRequesterIdAndStatusAndGroupId(ctx context.Context, requestId int64) (any, error)
	Insert(ctx context.Context, req *po.GroupJoinRequest) error
	HasPendingJoinRequest(ctx context.Context, requesterID, groupID int64) (bool, error)
	UpdateStatusIfPending(ctx context.Context, requestID, responderID int64, newStatus po.RequestStatus, reason *string, responseDate time.Time) (bool, error)
	FindRequestsByGroupID(ctx context.Context, groupID int64) ([]po.GroupJoinRequest, error)
	FindRequestsByRequesterID(ctx context.Context, requesterID int64) ([]po.GroupJoinRequest, error)
	FindByID(ctx context.Context, requestID int64) (*po.GroupJoinRequest, error)
	FindRequests(ctx context.Context,
		groupID *int64,
		requesterID *int64,
		responderID *int64,
		status *po.RequestStatus,
		creationDate *time.Time,
		responseDate *time.Time,
		expirationDate *time.Time,
		page int,
		size int) ([]*po.GroupJoinRequest, error)
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

// @MappedFrom findRequestsByGroupId(Long groupId)
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

// @MappedFrom findRequestsByRequesterId(Long requesterId)
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
func (r *groupJoinRequestRepository) FindByID(ctx context.Context, requestID int64) (*po.GroupJoinRequest, error) {
	filter := bson.M{"_id": requestID}
	var res po.GroupJoinRequest
	err := r.coll.FindOne(ctx, filter).Decode(&res)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &res, err
}

// @MappedFrom findRequests(@Nullable Set<Long> ids, @Nullable Set<Long> groupIds, @Nullable Set<Long> requesterIds, @Nullable Set<Long> responderIds, @Nullable Set<RequestStatus> statuses, @Nullable DateRange creationDateRange, @Nullable DateRange responseDateRange, @Nullable DateRange expirationDateRange, @Nullable Integer page, @Nullable Integer size)
func (r *groupJoinRequestRepository) FindRequests(ctx context.Context,
	groupID *int64,
	requesterID *int64,
	responderID *int64,
	status *po.RequestStatus,
	creationDate *time.Time,
	responseDate *time.Time,
	expirationDate *time.Time,
	page int,
	size int) ([]*po.GroupJoinRequest, error) {
	filter := bson.M{}
	if groupID != nil {
		filter["gid"] = *groupID
	}
	if requesterID != nil {
		filter["rqid"] = *requesterID
	}
	if responderID != nil {
		filter["rpid"] = *responderID
	}
	if status != nil {
		filter["stat"] = *status
	}
	if creationDate != nil {
		filter["cd"] = bson.M{"$gte": *creationDate}
	}
	if responseDate != nil {
		filter["rd"] = bson.M{"$gte": *responseDate}
	}
	if expirationDate != nil {
		filter["ed"] = bson.M{"$lt": *expirationDate}
	}

	skip := int64(page * size)
	limit := int64(size)
	opts := options.Find().
		SetSort(bson.M{"cd": -1}).
		SetSkip(skip).
		SetLimit(limit)

	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var reqs []*po.GroupJoinRequest
	if err := cursor.All(ctx, &reqs); err != nil {
		return nil, err
	}
	return reqs, nil
}


func (r *groupJoinRequestRepository) UpdateRequests(ctx context.Context, requestIds []int64, requesterId, responderId *int64, content *string, status any, creationDate, responseDate *time.Time) error { return nil }
func (r *groupJoinRequestRepository) CountRequests(ctx context.Context, ids, groupIds, requesterIds, responderIds []int64, statuses []int, creationDateRange, responseDateRange, expirationDateRange any) (int64, error) { return 0, nil }
func (r *groupJoinRequestRepository) FindGroupId(ctx context.Context, requestId int64) (int64, error) { return 0, nil }
func (r *groupJoinRequestRepository) FindRequesterIdAndStatusAndGroupId(ctx context.Context, requestId int64) (any, error) { return nil, nil }
