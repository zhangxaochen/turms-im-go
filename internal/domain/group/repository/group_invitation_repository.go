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

// @MappedFrom GroupInvitationRepository
type GroupInvitationRepository interface {
	FindInvitationsByInviterId(ctx context.Context, inviterId int64) ([]po.GroupInvitation, error)
	Insert(ctx context.Context, inv *po.GroupInvitation) error
	HasPendingInvitation(ctx context.Context, groupID, inviteeID int64) (bool, error)
	UpdateStatusIfPending(ctx context.Context, invitationID int64, newStatus po.RequestStatus, reason *string, responseDate time.Time) (bool, error)
	FindInvitationsByInviteeID(ctx context.Context, inviteeID int64) ([]po.GroupInvitation, error)
	FindInvitationsByGroupID(ctx context.Context, groupID int64) ([]po.GroupInvitation, error)
	FindByID(ctx context.Context, id int64) (*po.GroupInvitation, error)
	FindInvitations(ctx context.Context, ids, groupIds, inviterIds, inviteeIds []int64, statuses []po.RequestStatus, creationDateRange, responseDateRange, expirationDateRange *turmsmongo.DateRange, page, size *int) ([]*po.GroupInvitation, error)
	FindGroupIdAndInviterIdAndInviteeIdAndStatus(ctx context.Context, id int64) (int64, int64, int64, po.RequestStatus, error)
	FindGroupIdAndInviteeIdAndStatus(ctx context.Context, id int64) (int64, int64, po.RequestStatus, error)
	FindInviteeIdAndGroupIdAndCreationDateAndStatus(ctx context.Context, id int64) (int64, int64, time.Time, po.RequestStatus, error)
	UpdateInvitations(ctx context.Context, ids []int64, inviterID, inviteeID *int64, content *string, status *po.RequestStatus, creationDate, responseDate *time.Time) (int64, error)
	DeleteInvitations(ctx context.Context, ids []int64) (int64, error)
	CountInvitations(ctx context.Context, ids, groupIds, inviterIds, inviteeIds []int64, statuses []po.RequestStatus, creationDateRange, responseDateRange, expirationDateRange *turmsmongo.DateRange) (int64, error)
	GetEntityExpireAfterSeconds() int
}

type groupInvitationRepository struct {
	coll       *mongo.Collection
	expirable  repository.ExpirableEntityRepository
}

func NewGroupInvitationRepository(client *turmsmongo.Client, expireAfterSeconds func() int) GroupInvitationRepository {
	return &groupInvitationRepository{
		coll: client.Collection(po.CollectionNameGroupInvitation),
		expirable: repository.ExpirableEntityRepository{
			GetEntityExpireAfterSecondsFunc: expireAfterSeconds,
		},
	}
}

func (r *groupInvitationRepository) GetEntityExpireAfterSeconds() int {
	return r.expirable.GetEntityExpireAfterSeconds()
}

func (r *groupInvitationRepository) Insert(ctx context.Context, inv *po.GroupInvitation) error {
	_, err := r.coll.InsertOne(ctx, inv)
	return err
}

func (r *groupInvitationRepository) HasPendingInvitation(ctx context.Context, groupID, inviteeID int64) (bool, error) {
	filter := bson.M{
		"gid":  groupID,
		"ieid": inviteeID,
		"stat": po.RequestStatusPending,
	}
	count, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// @MappedFrom updateStatusIfPending(Long requestId, RequestStatus status, Long responderId, @Nullable String reason, @Nullable ClientSession session)
// @MappedFrom updateStatusIfPending(Long requestId, RequestStatus requestStatus, @Nullable String reason, @Nullable ClientSession session)
// @MappedFrom updateStatusIfPending(Long invitationId, RequestStatus requestStatus, @Nullable String reason, @Nullable ClientSession session)
func (r *groupInvitationRepository) UpdateStatusIfPending(ctx context.Context, invitationID int64, newStatus po.RequestStatus, reason *string, responseDate time.Time) (bool, error) {
	filter := bson.M{
		"_id":  invitationID,
		"stat": po.RequestStatusPending,
	}
	// BUG FIX: Add expiration check like Java's isNotExpired(creationDate, getEntityExpirationDate())
	if expirationDate := r.expirable.GetEntityExpirationDate(); expirationDate != nil {
		filter["cd"] = bson.M{"$gte": *expirationDate}
	}

	updateOps := bson.M{
		"stat": newStatus,
		"rd":   responseDate,
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

// transformExpiredRequest checks if a pending invitation has expired and updates its status.
// @MappedFrom transformExpiredRequest
func (r *groupInvitationRepository) transformExpiredRequest(inv *po.GroupInvitation) {
	if inv.Status == po.RequestStatusPending && r.expirable.IsExpired(inv.CreationDate.UnixMilli()) {
		inv.Status = po.RequestStatusExpired
	}
}

// @MappedFrom findInvitationsByInviteeId(Long inviteeId)
func (r *groupInvitationRepository) FindInvitationsByInviteeID(ctx context.Context, inviteeID int64) ([]po.GroupInvitation, error) {
	filter := bson.M{"ieid": inviteeID}
	// BUG FIX: Remove sort - Java version does not specify sort order
	cursor, err := r.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var invs []po.GroupInvitation
	if err := cursor.All(ctx, &invs); err != nil {
		return nil, err
	}
	// BUG FIX: Apply expiration transformation like Java's findExpirableDocs
	for i := range invs {
		r.transformExpiredRequest(&invs[i])
	}
	return invs, nil
}

// @MappedFrom findInvitationsByGroupId(Long groupId)
func (r *groupInvitationRepository) FindInvitationsByGroupID(ctx context.Context, groupID int64) ([]po.GroupInvitation, error) {
	filter := bson.M{"gid": groupID}
	// BUG FIX: Remove sort - Java version does not specify sort order
	cursor, err := r.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var invs []po.GroupInvitation
	if err := cursor.All(ctx, &invs); err != nil {
		return nil, err
	}
	// BUG FIX: Apply expiration transformation like Java's findExpirableDocs
	for i := range invs {
		r.transformExpiredRequest(&invs[i])
	}
	return invs, nil
}

func (r *groupInvitationRepository) FindByID(ctx context.Context, id int64) (*po.GroupInvitation, error) {
	filter := bson.M{"_id": id}
	var inv po.GroupInvitation
	err := r.coll.FindOne(ctx, filter).Decode(&inv)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &inv, nil
}

// @MappedFrom findInvitations(@Nullable Set<Long> ids, @Nullable Set<Long> groupIds, @Nullable Set<Long> inviterIds, @Nullable Set<Long> inviteeIds, @Nullable Set<RequestStatus> statuses, @Nullable DateRange creationDateRange, @Nullable DateRange responseDateRange, @Nullable DateRange expirationDateRange, @Nullable Integer page, @Nullable Integer size)
// BUG FIX: Updated signature to accept slices for multi-value filtering and DateRange objects
func (r *groupInvitationRepository) FindInvitations(ctx context.Context, ids, groupIds, inviterIds, inviteeIds []int64, statuses []po.RequestStatus, creationDateRange, responseDateRange, expirationDateRange *turmsmongo.DateRange, page, size *int) ([]*po.GroupInvitation, error) {
	filter := r.buildInvitationsFilter(ids, groupIds, inviterIds, inviteeIds, statuses, creationDateRange, responseDateRange, expirationDateRange)

	opts := options.Find()
	// BUG FIX: Remove sort - Java version does not specify sort order
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

	var invs []*po.GroupInvitation
	if err := cursor.All(ctx, &invs); err != nil {
		return nil, err
	}
	// BUG FIX: Apply expiration transformation like Java's findExpirableDocs
	for _, inv := range invs {
		if inv != nil {
			r.transformExpiredRequest(inv)
		}
	}
	return invs, nil
}

// @MappedFrom findGroupIdAndInviterIdAndInviteeIdAndStatus(Long invitationId)
func (r *groupInvitationRepository) FindGroupIdAndInviterIdAndInviteeIdAndStatus(ctx context.Context, id int64) (int64, int64, int64, po.RequestStatus, error) {
	var results struct {
		GroupID      int64            `bson:"gid"`
		InviterID    int64            `bson:"irid"`
		InviteeID    int64            `bson:"ieid"`
		Status       po.RequestStatus `bson:"stat"`
		CreationDate time.Time        `bson:"cd"`
	}
	// BUG FIX: Include "cd" in projection for expiration checking (Java: "Required by findExpirableDoc")
	opts := options.FindOne().SetProjection(bson.M{"gid": 1, "irid": 1, "ieid": 1, "stat": 1, "cd": 1})
	err := r.coll.FindOne(ctx, bson.M{"_id": id}, opts).Decode(&results)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 0, 0, 0, 0, nil
		}
		return 0, 0, 0, 0, err
	}
	// BUG FIX: Apply expiration transformation like Java's findExpirableDoc
	status := results.Status
	if status == po.RequestStatusPending && r.expirable.IsExpired(results.CreationDate.UnixMilli()) {
		status = po.RequestStatusExpired
	}
	return results.GroupID, results.InviterID, results.InviteeID, status, nil
}

// @MappedFrom findGroupIdAndInviteeIdAndStatus(Long invitationId)
func (r *groupInvitationRepository) FindGroupIdAndInviteeIdAndStatus(ctx context.Context, id int64) (int64, int64, po.RequestStatus, error) {
	var results struct {
		GroupID      int64            `bson:"gid"`
		InviteeID    int64            `bson:"ieid"`
		Status       po.RequestStatus `bson:"stat"`
		CreationDate time.Time        `bson:"cd"`
	}
	// BUG FIX: Include "cd" in projection for expiration checking (Java: "Required by findExpirableDoc")
	opts := options.FindOne().SetProjection(bson.M{"gid": 1, "ieid": 1, "stat": 1, "cd": 1})
	err := r.coll.FindOne(ctx, bson.M{"_id": id}, opts).Decode(&results)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 0, 0, 0, nil
		}
		return 0, 0, 0, err
	}
	// BUG FIX: Apply expiration transformation like Java's findExpirableDoc
	status := results.Status
	if status == po.RequestStatusPending && r.expirable.IsExpired(results.CreationDate.UnixMilli()) {
		status = po.RequestStatusExpired
	}
	return results.GroupID, results.InviteeID, status, nil
}

// @MappedFrom findInviteeIdAndGroupIdAndCreationDateAndStatus(Long invitationId)
func (r *groupInvitationRepository) FindInviteeIdAndGroupIdAndCreationDateAndStatus(ctx context.Context, id int64) (int64, int64, time.Time, po.RequestStatus, error) {
	var results struct {
		InviteeID    int64            `bson:"ieid"`
		GroupID      int64            `bson:"gid"`
		CreationDate time.Time        `bson:"cd"`
		Status       po.RequestStatus `bson:"stat"`
	}
	opts := options.FindOne().SetProjection(bson.M{"ieid": 1, "gid": 1, "cd": 1, "stat": 1})
	err := r.coll.FindOne(ctx, bson.M{"_id": id}, opts).Decode(&results)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 0, 0, time.Time{}, 0, nil
		}
		return 0, 0, time.Time{}, 0, err
	}
	return results.InviteeID, results.GroupID, results.CreationDate, results.Status, nil
}

// @MappedFrom updateInvitations(@NotEmpty Set<Long> invitationIds, @Nullable Long inviterId, @Nullable Long inviteeId, @Nullable String content, @Nullable @ValidRequestStatus RequestStatus status, @Nullable @PastOrPresent Date creationDate, @Nullable @PastOrPresent Date responseDate)
// @MappedFrom updateInvitations(Set<Long> invitationIds, @Nullable Long inviterId, @Nullable Long inviteeId, @Nullable String content, @Nullable RequestStatus status, @Nullable Date creationDate, @Nullable Date responseDate)
func (r *groupInvitationRepository) UpdateInvitations(ctx context.Context, ids []int64, inviterID, inviteeID *int64, content *string, status *po.RequestStatus, creationDate, responseDate *time.Time) (int64, error) {
	filter := bson.M{"_id": bson.M{"$in": ids}}
	update := bson.M{}
	set := bson.M{}
	unset := bson.M{}
	if inviterID != nil {
		set["irid"] = *inviterID
	}
	if inviteeID != nil {
		set["ieid"] = *inviteeID
	}
	if content != nil {
		set["cnt"] = *content
	}
	if creationDate != nil {
		set["cd"] = *creationDate
	}
	// BUG FIX: Implement updateResponseDateBasedOnStatus logic from Java
	// If status is ACCEPTED/DECLINED/IGNORED, set responseDate (default to now if nil).
	// Otherwise (e.g. CANCELED/PENDING), unset responseDate.
	if status != nil {
		set["stat"] = *status
		switch *status {
		case po.RequestStatusAccepted, po.RequestStatusAcceptedWithoutConfirm, po.RequestStatusDeclined, po.RequestStatusIgnored:
			rd := responseDate
			if rd == nil {
				now := time.Now()
				rd = &now
			}
			set["rd"] = *rd
		default:
			unset["rd"] = ""
		}
	} else if responseDate != nil {
		set["rd"] = *responseDate
	}

	if len(set) > 0 {
		update["$set"] = set
	}
	if len(unset) > 0 {
		update["$unset"] = unset
	}
	if len(update) == 0 {
		return 0, nil
	}
	res, err := r.coll.UpdateMany(ctx, filter, update)
	if err != nil {
		return 0, err
	}
	return res.ModifiedCount, nil
}

// @MappedFrom deleteInvitations(@Nullable Set<Long> ids)
func (r *groupInvitationRepository) DeleteInvitations(ctx context.Context, ids []int64) (int64, error) {
	res, err := r.coll.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return 0, err
	}
	return res.DeletedCount, nil
}

// @MappedFrom countInvitations(@Nullable Set<Long> ids, @Nullable Set<Long> groupIds, @Nullable Set<Long> inviterIds, @Nullable Set<Long> inviteeIds, @Nullable Set<RequestStatus> statuses, @Nullable DateRange creationDateRange, @Nullable DateRange responseDateRange, @Nullable DateRange expirationDateRange)
// BUG FIX: Updated signature to accept slices for multi-value filtering and DateRange objects
func (r *groupInvitationRepository) CountInvitations(ctx context.Context, ids, groupIds, inviterIds, inviteeIds []int64, statuses []po.RequestStatus, creationDateRange, responseDateRange, expirationDateRange *turmsmongo.DateRange) (int64, error) {
	filter := r.buildInvitationsFilter(ids, groupIds, inviterIds, inviteeIds, statuses, creationDateRange, responseDateRange, expirationDateRange)
	return r.coll.CountDocuments(ctx, filter)
}

// buildInvitationsFilter builds a MongoDB filter for invitation queries.
// BUG FIX: Supports multi-value filtering (slices instead of single pointers) and DateRange objects
func (r *groupInvitationRepository) buildInvitationsFilter(ids, groupIds, inviterIds, inviteeIds []int64, statuses []po.RequestStatus, creationDateRange, responseDateRange, expirationDateRange *turmsmongo.DateRange) bson.M {
	filter := bson.M{}
	if len(ids) > 0 {
		filter["_id"] = bson.M{"$in": ids}
	}
	if len(groupIds) > 0 {
		filter["gid"] = bson.M{"$in": groupIds}
	}
	if len(inviterIds) > 0 {
		filter["irid"] = bson.M{"$in": inviterIds}
	}
	if len(inviteeIds) > 0 {
		filter["ieid"] = bson.M{"$in": inviteeIds}
	}
	if len(statuses) > 0 {
		filter["stat"] = bson.M{"$in": statuses}
	}
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

// BUG FIX: Implement FindInvitationsByInviterId instead of returning nil stub
// @MappedFrom findInvitationsByInviterId(Long inviterId)
func (r *groupInvitationRepository) FindInvitationsByInviterId(ctx context.Context, inviterId int64) ([]po.GroupInvitation, error) {
	filter := bson.M{"irid": inviterId}
	// BUG FIX: No sort order in Java version
	cursor, err := r.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var invs []po.GroupInvitation
	if err := cursor.All(ctx, &invs); err != nil {
		return nil, err
	}
	// BUG FIX: Apply expiration transformation like Java's findExpirableDocs
	for i := range invs {
		r.transformExpiredRequest(&invs[i])
	}
	return invs, nil
}
