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
	HasPendingOrDeclinedOrIgnoredOrExpiredRequest(ctx context.Context, requesterID, recipientID int64) (bool, error)
	UpdateStatusIfPending(ctx context.Context, requestID int64, newStatus po.RequestStatus, reason *string, responseDate time.Time) (bool, error)
	UpdateFriendRequests(ctx context.Context, requestIds []int64, requesterID, recipientID *int64, content *string, status *po.RequestStatus, reason *string, creationDate *time.Time) error
	FindFriendRequestsByRecipientId(ctx context.Context, recipientID int64) ([]po.UserFriendRequest, error)
	FindFriendRequestsByRequesterId(ctx context.Context, requesterID int64) ([]po.UserFriendRequest, error)
	GetEntityExpireAfterSeconds() int
	FindRecipientId(ctx context.Context, requestID int64) (int64, error)
	FindRequesterIdAndRecipientIdAndStatus(ctx context.Context, requestID int64) (*po.UserFriendRequest, error)
	FindRequesterIdAndRecipientIdAndCreationDateAndStatus(ctx context.Context, requestID int64) (*po.UserFriendRequest, error)
	DeleteExpiredData(ctx context.Context, creationDateFieldName string, expirationDate time.Time) error
	DeleteByIds(ctx context.Context, ids []int64) error
	FindFriendRequests(ctx context.Context, ids, requesterIds, recipientIds []int64, statuses []po.RequestStatus, creationDateStart, creationDateEnd, responseDateStart, responseDateEnd, expirationDateStart, expirationDateEnd *time.Time, page, size *int) ([]po.UserFriendRequest, error)
	CountFriendRequests(ctx context.Context, ids, requesterIds, recipientIds []int64, statuses []po.RequestStatus, creationDateStart, creationDateEnd, responseDateStart, responseDateEnd, expirationDateStart, expirationDateEnd *time.Time) (int64, error)
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

// @MappedFrom hasPendingFriendRequest(@NotNull Long requesterId, @NotNull Long recipientId)
// @MappedFrom hasPendingFriendRequest(Long requesterId, Long recipientId)
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

// @MappedFrom hasPendingOrDeclinedOrIgnoredOrExpiredRequest(Long requesterId, Long recipientId)
func (r *userFriendRequestRepository) HasPendingOrDeclinedOrIgnoredOrExpiredRequest(ctx context.Context, requesterID, recipientID int64) (bool, error) {
	filter := bson.M{
		"rqid": requesterID,
		"rcid": recipientID,
		"s": bson.M{"$in": []po.RequestStatus{
			po.RequestStatusPending,
			po.RequestStatusDeclined,
			po.RequestStatusIgnored,
			po.RequestStatusExpired,
		}},
	}
	count, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *userFriendRequestRepository) UpdateStatusIfPending(ctx context.Context, requestID int64, newStatus po.RequestStatus, reason *string, responseDate time.Time) (bool, error) {
	filter := bson.M{
		"_id": requestID,
		"s":   po.RequestStatusPending,
	}
	updateOps := bson.M{
		"s":  newStatus,
		"rd": responseDate,
	}
	if reason != nil {
		updateOps["r"] = *reason
	}
	update := bson.M{"$set": updateOps}

	res, err := r.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return false, err
	}
	return res.ModifiedCount > 0, nil
}

// @MappedFrom updateFriendRequests(Set<Long> ids, @RequestBody UpdateFriendRequestDTO updateFriendRequestDTO)
// @MappedFrom updateFriendRequests(@NotEmpty Set<Long> requestIds, @Nullable Long requesterId, @Nullable Long recipientId, @Nullable String content, @Nullable @ValidRequestStatus RequestStatus status, @Nullable String reason, @Nullable @PastOrPresent Date creationDate, @Nullable @PastOrPresent Date responseDate)
// @MappedFrom updateFriendRequests(Set<Long> requestIds, @Nullable Long requesterId, @Nullable Long recipientId, @Nullable String content, @Nullable RequestStatus status, @Nullable String reason, @Nullable Date creationDate)
func (r *userFriendRequestRepository) UpdateFriendRequests(ctx context.Context, requestIds []int64, requesterID, recipientID *int64, content *string, status *po.RequestStatus, reason *string, creationDate *time.Time) error {
	filter := bson.M{"_id": bson.M{"$in": requestIds}}
	updateOps := bson.M{}
	if requesterID != nil {
		updateOps["rqid"] = *requesterID
	}
	if recipientID != nil {
		updateOps["rcid"] = *recipientID
	}
	if content != nil {
		updateOps["txt"] = *content
	}
	if status != nil {
		updateOps["s"] = *status
	}
	if reason != nil {
		updateOps["r"] = *reason
	}
	if creationDate != nil {
		updateOps["cd"] = *creationDate
	}
	if len(updateOps) == 0 {
		return nil
	}
	_, err := r.coll.UpdateMany(ctx, filter, bson.M{"$set": updateOps})
	return err
}

// @MappedFrom getEntityExpireAfterSeconds()
func (r *userFriendRequestRepository) GetEntityExpireAfterSeconds() int {
	return 0
}

// @MappedFrom findFriendRequestsByRecipientId(Long recipientId)
func (r *userFriendRequestRepository) FindFriendRequestsByRecipientId(ctx context.Context, recipientID int64) ([]po.UserFriendRequest, error) {
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

// @MappedFrom findFriendRequestsByRequesterId(Long requesterId)
func (r *userFriendRequestRepository) FindFriendRequestsByRequesterId(ctx context.Context, requesterID int64) ([]po.UserFriendRequest, error) {
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

// @MappedFrom findRecipientId(Long requestId)
func (r *userFriendRequestRepository) FindRecipientId(ctx context.Context, requestID int64) (int64, error) {
	filter := bson.M{"_id": requestID}
	opts := options.FindOne().SetProjection(bson.M{"rcid": 1})
	var result struct {
		RecipientID int64 `bson:"rcid"`
	}
	err := r.coll.FindOne(ctx, filter, opts).Decode(&result)
	return result.RecipientID, err
}

// @MappedFrom findRequesterIdAndRecipientIdAndStatus(Long requestId)
func (r *userFriendRequestRepository) FindRequesterIdAndRecipientIdAndStatus(ctx context.Context, requestID int64) (*po.UserFriendRequest, error) {
	filter := bson.M{"_id": requestID}
	opts := options.FindOne().SetProjection(bson.M{"rqid": 1, "rcid": 1, "s": 1})
	var result po.UserFriendRequest
	err := r.coll.FindOne(ctx, filter, opts).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

// @MappedFrom findRequesterIdAndRecipientIdAndCreationDateAndStatus(Long requestId)
func (r *userFriendRequestRepository) FindRequesterIdAndRecipientIdAndCreationDateAndStatus(ctx context.Context, requestID int64) (*po.UserFriendRequest, error) {
	filter := bson.M{"_id": requestID}
	opts := options.FindOne().SetProjection(bson.M{"rqid": 1, "rcid": 1, "cd": 1, "s": 1})
	var result po.UserFriendRequest
	err := r.coll.FindOne(ctx, filter, opts).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

// @MappedFrom deleteExpiredData(String creationDateFieldName, Date expirationDate)
func (r *userFriendRequestRepository) DeleteExpiredData(ctx context.Context, creationDateFieldName string, expirationDate time.Time) error {
	filter := bson.M{creationDateFieldName: bson.M{"$lt": expirationDate}}
	_, err := r.coll.DeleteMany(ctx, filter)
	return err
}

func (r *userFriendRequestRepository) DeleteByIds(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	filter := bson.M{"_id": bson.M{"$in": ids}}
	_, err := r.coll.DeleteMany(ctx, filter)
	return err
}

func (r *userFriendRequestRepository) countOrFind(ctx context.Context, ids, requesterIds, recipientIds []int64, statuses []po.RequestStatus, creationDateStart, creationDateEnd, responseDateStart, responseDateEnd, expirationDateStart, expirationDateEnd *time.Time) bson.M {
	filter := bson.M{}
	if len(ids) > 0 {
		filter["_id"] = bson.M{"$in": ids}
	}
	if len(requesterIds) > 0 {
		filter["rqid"] = bson.M{"$in": requesterIds}
	}
	if len(recipientIds) > 0 {
		filter["rcid"] = bson.M{"$in": recipientIds}
	}
	if len(statuses) > 0 {
		filter["s"] = bson.M{"$in": statuses}
	}
	if creationDateStart != nil || creationDateEnd != nil {
		cdFilter := bson.M{}
		if creationDateStart != nil {
			cdFilter["$gte"] = *creationDateStart
		}
		if creationDateEnd != nil {
			cdFilter["$lt"] = *creationDateEnd
		}
		filter["cd"] = cdFilter
	}
	if responseDateStart != nil || responseDateEnd != nil {
		rdFilter := bson.M{}
		if responseDateStart != nil {
			rdFilter["$gte"] = *responseDateStart
		}
		if responseDateEnd != nil {
			rdFilter["$lt"] = *responseDateEnd
		}
		filter["rd"] = rdFilter
	}
	return filter
}

// @MappedFrom findFriendRequests(@Nullable Set<Long> ids, @Nullable Set<Long> requesterIds, @Nullable Set<Long> recipientIds, @Nullable Set<RequestStatus> statuses, @Nullable DateRange creationDateRange, @Nullable DateRange responseDateRange, @Nullable DateRange expirationDateRange, @Nullable Integer page, @Nullable Integer size)
func (r *userFriendRequestRepository) FindFriendRequests(ctx context.Context, ids, requesterIds, recipientIds []int64, statuses []po.RequestStatus, creationDateStart, creationDateEnd, responseDateStart, responseDateEnd, expirationDateStart, expirationDateEnd *time.Time, page, size *int) ([]po.UserFriendRequest, error) {
	filter := r.countOrFind(ctx, ids, requesterIds, recipientIds, statuses, creationDateStart, creationDateEnd, responseDateStart, responseDateEnd, expirationDateStart, expirationDateEnd)
	opts := options.Find().SetSort(bson.M{"cd": -1})
	if page != nil && size != nil {
		opts.SetSkip(int64(*page * *size))
		opts.SetLimit(int64(*size))
	} else if size != nil {
		opts.SetLimit(int64(*size))
	}

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

// @MappedFrom countFriendRequests(@Nullable Set<Long> ids, @Nullable Set<Long> requesterIds, @Nullable Set<Long> recipientIds, @Nullable Set<RequestStatus> statuses, @Nullable DateRange creationDateRange, @Nullable DateRange responseDateRange, @Nullable DateRange expirationDateRange)
func (r *userFriendRequestRepository) CountFriendRequests(ctx context.Context, ids, requesterIds, recipientIds []int64, statuses []po.RequestStatus, creationDateStart, creationDateEnd, responseDateStart, responseDateEnd, expirationDateStart, expirationDateEnd *time.Time) (int64, error) {
	filter := r.countOrFind(ctx, ids, requesterIds, recipientIds, statuses, creationDateStart, creationDateEnd, responseDateStart, responseDateEnd, expirationDateStart, expirationDateEnd)
	return r.coll.CountDocuments(ctx, filter)
}
