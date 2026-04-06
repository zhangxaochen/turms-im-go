package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"im.turms/server/internal/domain/common/repository"
	"im.turms/server/internal/domain/group/po"
	turmsmongo "im.turms/server/internal/storage/mongo"
)

// @MappedFrom GroupJoinRequestRepository
type GroupJoinRequestRepository interface {
	UpdateRequests(ctx context.Context, requestIds []int64, requesterId, responderId *int64, content *string, status *po.RequestStatus, creationDate, responseDate *time.Time) error
	CountRequests(ctx context.Context, ids, groupIds, requesterIds, responderIds []int64, statuses []po.RequestStatus, creationDateRange, responseDateRange, expirationDateRange *turmsmongo.DateRange) (int64, error)
	FindGroupId(ctx context.Context, requestId int64) (*int64, error)
	FindRequesterIdAndStatusAndGroupId(ctx context.Context, requestId int64) (*int64, *po.RequestStatus, *int64, error)
	DeleteRequests(ctx context.Context, ids []int64) (int64, error)
	Insert(ctx context.Context, req *po.GroupJoinRequest) error
	HasPendingJoinRequest(ctx context.Context, requesterID, groupID int64) (bool, error)
	UpdateStatusIfPending(ctx context.Context, requestID, responderID int64, newStatus po.RequestStatus, reason *string, responseDate time.Time) (bool, error)
	FindRequestsByGroupID(ctx context.Context, groupID int64) ([]po.GroupJoinRequest, error)
	FindRequestsByRequesterID(ctx context.Context, requesterID int64) ([]po.GroupJoinRequest, error)
	FindByID(ctx context.Context, requestID int64) (*po.GroupJoinRequest, error)
	FindRequests(ctx context.Context,
		ids, groupIds, requesterIds, responderIds []int64,
		statuses []po.RequestStatus,
		creationDateRange, responseDateRange, expirationDateRange *turmsmongo.DateRange,
		page, size *int) ([]*po.GroupJoinRequest, error)
	GetEntityExpireAfterSeconds() int
}

type groupJoinRequestRepository struct {
	coll      *mongo.Collection
	expirable repository.ExpirableEntityRepository
}

func NewGroupJoinRequestRepository(client *turmsmongo.Client, expireAfterSeconds func() int) GroupJoinRequestRepository {
	return &groupJoinRequestRepository{
		coll: client.Collection(po.CollectionNameGroupJoinRequest),
		expirable: repository.ExpirableEntityRepository{
			GetEntityExpireAfterSecondsFunc: expireAfterSeconds,
		},
	}
}

func (r *groupJoinRequestRepository) GetEntityExpireAfterSeconds() int {
	return r.expirable.GetEntityExpireAfterSeconds()
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
	// BUG FIX: Add expiration check like Java's isNotExpired(creationDate, getEntityExpirationDate())
	if expirationDate := r.expirable.GetEntityExpirationDate(); expirationDate != nil {
		filter["cd"] = bson.M{"$gte": *expirationDate}
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

// transformExpiredRequest checks if a pending request has expired and updates its status.
// @MappedFrom transformExpiredRequest
func (r *groupJoinRequestRepository) transformExpiredRequest(req *po.GroupJoinRequest) {
	if req.Status == po.RequestStatusPending && r.expirable.IsExpired(req.CreationDate.UnixMilli()) {
		req.Status = po.RequestStatusExpired
	}
}

// @MappedFrom findRequestsByGroupId(Long groupId)
func (r *groupJoinRequestRepository) FindRequestsByGroupID(ctx context.Context, groupID int64) ([]po.GroupJoinRequest, error) {
	filter := bson.M{"gid": groupID}
	// BUG FIX: Remove sort - Java version does not specify sort order
	cursor, err := r.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var reqs []po.GroupJoinRequest
	if err := cursor.All(ctx, &reqs); err != nil {
		return nil, err
	}
	// BUG FIX: Apply expiration transformation like Java's findExpirableDocs
	for i := range reqs {
		r.transformExpiredRequest(&reqs[i])
	}
	return reqs, nil
}

// @MappedFrom findRequestsByRequesterId(Long requesterId)
func (r *groupJoinRequestRepository) FindRequestsByRequesterID(ctx context.Context, requesterID int64) ([]po.GroupJoinRequest, error) {
	filter := bson.M{"rqid": requesterID}
	// BUG FIX: Remove sort - Java version does not specify sort order
	cursor, err := r.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var reqs []po.GroupJoinRequest
	if err := cursor.All(ctx, &reqs); err != nil {
		return nil, err
	}
	// BUG FIX: Apply expiration transformation like Java's findExpirableDocs
	for i := range reqs {
		r.transformExpiredRequest(&reqs[i])
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
// BUG FIX: Updated signature to accept slices and DateRange objects
func (r *groupJoinRequestRepository) FindRequests(ctx context.Context,
	ids, groupIds, requesterIds, responderIds []int64,
	statuses []po.RequestStatus,
	creationDateRange, responseDateRange, expirationDateRange *turmsmongo.DateRange,
	page, size *int) ([]*po.GroupJoinRequest, error) {

	filter := r.buildRequestsFilter(ids, groupIds, requesterIds, responderIds, statuses, creationDateRange, responseDateRange, expirationDateRange)

	opts := options.Find()
	// BUG FIX: Only apply pagination when both page and size are non-nil (Java uses paginateIfNotNull)
	if page != nil && size != nil && *size > 0 {
		opts.SetSkip(int64(*page * *size))
		opts.SetLimit(int64(*size))
	}

	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var reqs []*po.GroupJoinRequest
	if err := cursor.All(ctx, &reqs); err != nil {
		return nil, err
	}
	// BUG FIX: Apply expiration transformation like Java's findExpirableDocs
	for _, req := range reqs {
		if req != nil {
			r.transformExpiredRequest(req)
		}
	}
	return reqs, nil
}

func (r *groupJoinRequestRepository) UpdateRequests(ctx context.Context, requestIds []int64, requesterId, responderId *int64, content *string, status *po.RequestStatus, creationDate, responseDate *time.Time) error {
	filter := bson.M{}
	if len(requestIds) > 0 {
		filter["_id"] = bson.M{"$in": requestIds}
	}

	updateOps := bson.M{}
	unsetOps := bson.M{}
	if requesterId != nil {
		updateOps["rqid"] = *requesterId
	}
	if responderId != nil {
		updateOps["rpid"] = *responderId
	}
	if content != nil {
		updateOps["cnt"] = *content
	}
	if creationDate != nil {
		updateOps["cd"] = *creationDate
	}

	// BUG FIX: updateResponseDateBasedOnStatus: if status is a responder-processed status
	// (ACCEPTED, DECLINED, IGNORED), set responseDate (defaulting to now if nil).
	// Otherwise (e.g. CANCELED, PENDING), unset responseDate.
	if status != nil {
		updateOps["stat"] = *status
		switch *status {
		case po.RequestStatusAccepted, po.RequestStatusAcceptedWithoutConfirm, po.RequestStatusDeclined, po.RequestStatusIgnored:
			rd := responseDate
			if rd == nil {
				now := time.Now()
				rd = &now
			}
			updateOps["rd"] = *rd
		default:
			unsetOps["rd"] = ""
		}
	} else if responseDate != nil {
		updateOps["rd"] = *responseDate
	}

	if len(updateOps) == 0 && len(unsetOps) == 0 {
		return nil
	}
	update := bson.M{}
	if len(updateOps) > 0 {
		update["$set"] = updateOps
	}
	if len(unsetOps) > 0 {
		update["$unset"] = unsetOps
	}
	_, err := r.coll.UpdateMany(ctx, filter, update)
	return err
}

// BUG FIX: Updated CountRequests to accept slices and DateRange objects properly
func (r *groupJoinRequestRepository) CountRequests(ctx context.Context, ids, groupIds, requesterIds, responderIds []int64, statuses []po.RequestStatus, creationDateRange, responseDateRange, expirationDateRange *turmsmongo.DateRange) (int64, error) {
	filter := r.buildRequestsFilter(ids, groupIds, requesterIds, responderIds, statuses, creationDateRange, responseDateRange, expirationDateRange)
	return r.coll.CountDocuments(ctx, filter)
}

// buildRequestsFilter builds a MongoDB filter for request queries.
func (r *groupJoinRequestRepository) buildRequestsFilter(ids, groupIds, requesterIds, responderIds []int64, statuses []po.RequestStatus, creationDateRange, responseDateRange, expirationDateRange *turmsmongo.DateRange) bson.M {
	filter := bson.M{}
	if len(ids) > 0 {
		filter["_id"] = bson.M{"$in": ids}
	}
	if len(groupIds) > 0 {
		filter["gid"] = bson.M{"$in": groupIds}
	}
	if len(requesterIds) > 0 {
		filter["rqid"] = bson.M{"$in": requesterIds}
	}
	if len(responderIds) > 0 {
		filter["rpid"] = bson.M{"$in": responderIds}
	}
	if len(statuses) > 0 {
		filter["stat"] = bson.M{"$in": statuses}
	}
	// BUG FIX: Use DateRange for date filtering instead of single time pointers
	if creationDateRange != nil {
		if d := creationDateRange.ToBson(); d != nil {
			filter["cd"] = d
		}
	}
	if responseDateRange != nil {
		if d := responseDateRange.ToBson(); d != nil {
			filter["rd"] = d
		}
	}
	// BUG FIX: Java converts expirationDateRange to a creation date range offset
	// using getCreationDateRange(creationDateRange, expirationDateRange).
	// For now, apply expiration date range to creation date field if no creation date range exists.
	if expirationDateRange != nil {
		if d := expirationDateRange.ToBson(); d != nil {
			if _, exists := filter["cd"]; !exists {
				filter["cd"] = d
			}
		}
	}
	return filter
}

func (r *groupJoinRequestRepository) FindGroupId(ctx context.Context, requestId int64) (*int64, error) {
	filter := bson.M{"_id": requestId}
	opts := options.FindOne().SetProjection(bson.M{"gid": 1})
	var result struct {
		GroupID int64 `bson:"gid"`
	}
	err := r.coll.FindOne(ctx, filter, opts).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &result.GroupID, nil
}

// BUG FIX: Include "cd" in projection and apply expiration transformation
func (r *groupJoinRequestRepository) FindRequesterIdAndStatusAndGroupId(ctx context.Context, requestId int64) (*int64, *po.RequestStatus, *int64, error) {
	filter := bson.M{"_id": requestId}
	// BUG FIX: Include "cd" in projection for expiration checking (Java: "Required by findExpirableDoc")
	opts := options.FindOne().SetProjection(bson.M{"rqid": 1, "stat": 1, "gid": 1, "cd": 1})
	var result struct {
		RequesterID int64            `bson:"rqid"`
		Status      po.RequestStatus `bson:"stat"`
		GroupID     int64            `bson:"gid"`
		CreationDate time.Time       `bson:"cd"`
	}
	err := r.coll.FindOne(ctx, filter, opts).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil, nil, nil
		}
		return nil, nil, nil, err
	}
	// BUG FIX: Apply expiration transformation like Java's findExpirableDoc
	status := result.Status
	if status == po.RequestStatusPending && r.expirable.IsExpired(result.CreationDate.UnixMilli()) {
		status = po.RequestStatusExpired
	}
	return &result.RequesterID, &status, &result.GroupID, nil
}

func (r *groupJoinRequestRepository) DeleteRequests(ctx context.Context, ids []int64) (int64, error) {
	filter := bson.M{}
	if len(ids) > 0 {
		filter["_id"] = bson.M{"$in": ids}
	} else {
		// Do not delete anything if ids are provided as empty
		return 0, nil
	}
	res, err := r.coll.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}
	return res.DeletedCount, nil
}
