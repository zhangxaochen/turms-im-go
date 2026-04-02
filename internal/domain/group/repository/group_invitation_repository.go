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

type GroupInvitationRepository interface {
	Insert(ctx context.Context, inv *po.GroupInvitation) error
	HasPendingInvitation(ctx context.Context, groupID, inviteeID int64) (bool, error)
	UpdateStatusIfPending(ctx context.Context, invitationID int64, newStatus po.RequestStatus, reason *string, responseDate time.Time) (bool, error)
	FindInvitationsByInviteeID(ctx context.Context, inviteeID int64) ([]po.GroupInvitation, error)
	FindInvitationsByGroupID(ctx context.Context, groupID int64) ([]po.GroupInvitation, error)
	FindByID(ctx context.Context, id int64) (*po.GroupInvitation, error)
	FindInvitations(ctx context.Context, groupID *int64, inviterID *int64, inviteeID *int64, status *po.RequestStatus, lastUpdatedDate *time.Time, page, size int) ([]*po.GroupInvitation, error)
	FindGroupIdAndInviterIdAndInviteeIdAndStatus(ctx context.Context, id int64) (int64, int64, int64, po.RequestStatus, error)
	FindGroupIdAndInviteeIdAndStatus(ctx context.Context, id int64) (int64, int64, po.RequestStatus, error)
	FindInviteeIdAndGroupIdAndCreationDateAndStatus(ctx context.Context, id int64) (int64, int64, time.Time, po.RequestStatus, error)
	UpdateInvitations(ctx context.Context, ids []int64, inviterID, inviteeID *int64, content *string, status *po.RequestStatus, creationDate, responseDate *time.Time) (int64, error)
	DeleteInvitations(ctx context.Context, ids []int64) (int64, error)
	CountInvitations(ctx context.Context, groupID, inviterID, inviteeID *int64, status *po.RequestStatus, lastUpdatedDate *time.Time) (int64, error)
}

type groupInvitationRepository struct {
	coll *mongo.Collection
}

func NewGroupInvitationRepository(client *turmsmongo.Client) GroupInvitationRepository {
	return &groupInvitationRepository{
		coll: client.Collection(po.CollectionNameGroupInvitation),
	}
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

// @MappedFrom findInvitationsByInviteeId(Long inviteeId)
func (r *groupInvitationRepository) FindInvitationsByInviteeID(ctx context.Context, inviteeID int64) ([]po.GroupInvitation, error) {
	filter := bson.M{"ieid": inviteeID}
	opts := options.Find().SetSort(bson.M{"cd": -1})
	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var invs []po.GroupInvitation
	if err := cursor.All(ctx, &invs); err != nil {
		return nil, err
	}
	return invs, nil
}

// @MappedFrom findInvitationsByGroupId(Long groupId)
func (r *groupInvitationRepository) FindInvitationsByGroupID(ctx context.Context, groupID int64) ([]po.GroupInvitation, error) {
	filter := bson.M{"gid": groupID}
	opts := options.Find().SetSort(bson.M{"cd": -1})
	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var invs []po.GroupInvitation
	if err := cursor.All(ctx, &invs); err != nil {
		return nil, err
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
// @MappedFrom findInvitations(Long groupId)
func (r *groupInvitationRepository) FindInvitations(ctx context.Context, groupID *int64, inviterID *int64, inviteeID *int64, status *po.RequestStatus, lastUpdatedDate *time.Time, page, size int) ([]*po.GroupInvitation, error) {
	filter := r.buildFilter(groupID, inviterID, inviteeID, status, lastUpdatedDate)

	opts := options.Find().
		SetSort(bson.M{"cd": -1}).
		SetSkip(int64(page * size)).
		SetLimit(int64(size))

	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var invs []*po.GroupInvitation
	if err := cursor.All(ctx, &invs); err != nil {
		return nil, err
	}
	return invs, nil
}

// @MappedFrom findGroupIdAndInviterIdAndInviteeIdAndStatus(Long invitationId)
func (r *groupInvitationRepository) FindGroupIdAndInviterIdAndInviteeIdAndStatus(ctx context.Context, id int64) (int64, int64, int64, po.RequestStatus, error) {
	var results struct {
		GroupID   int64            `bson:"gid"`
		InviterID int64            `bson:"irid"`
		InviteeID int64            `bson:"ieid"`
		Status    po.RequestStatus `bson:"stat"`
	}
	opts := options.FindOne().SetProjection(bson.M{"gid": 1, "irid": 1, "ieid": 1, "stat": 1})
	err := r.coll.FindOne(ctx, bson.M{"_id": id}, opts).Decode(&results)
	return results.GroupID, results.InviterID, results.InviteeID, results.Status, err
}

// @MappedFrom findGroupIdAndInviteeIdAndStatus(Long invitationId)
func (r *groupInvitationRepository) FindGroupIdAndInviteeIdAndStatus(ctx context.Context, id int64) (int64, int64, po.RequestStatus, error) {
	var results struct {
		GroupID   int64            `bson:"gid"`
		InviteeID int64            `bson:"ieid"`
		Status    po.RequestStatus `bson:"stat"`
	}
	opts := options.FindOne().SetProjection(bson.M{"gid": 1, "ieid": 1, "stat": 1})
	err := r.coll.FindOne(ctx, bson.M{"_id": id}, opts).Decode(&results)
	return results.GroupID, results.InviteeID, results.Status, err
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
	return results.InviteeID, results.GroupID, results.CreationDate, results.Status, err
}

// @MappedFrom updateInvitations(@NotEmpty Set<Long> invitationIds, @Nullable Long inviterId, @Nullable Long inviteeId, @Nullable String content, @Nullable @ValidRequestStatus RequestStatus status, @Nullable @PastOrPresent Date creationDate, @Nullable @PastOrPresent Date responseDate)
// @MappedFrom updateInvitations(Set<Long> invitationIds, @Nullable Long inviterId, @Nullable Long inviteeId, @Nullable String content, @Nullable RequestStatus status, @Nullable Date creationDate, @Nullable Date responseDate)
func (r *groupInvitationRepository) UpdateInvitations(ctx context.Context, ids []int64, inviterID, inviteeID *int64, content *string, status *po.RequestStatus, creationDate, responseDate *time.Time) (int64, error) {
	filter := bson.M{"_id": bson.M{"$in": ids}}
	update := bson.M{}
	set := bson.M{}
	if inviterID != nil {
		set["irid"] = *inviterID
	}
	if inviteeID != nil {
		set["ieid"] = *inviteeID
	}
	if content != nil {
		set["cnt"] = *content
	}
	if status != nil {
		set["stat"] = *status
	}
	if creationDate != nil {
		set["cd"] = *creationDate
	}
	if responseDate != nil {
		set["rd"] = *responseDate
	}
	if len(set) == 0 {
		return 0, nil
	}
	update["$set"] = set
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
func (r *groupInvitationRepository) CountInvitations(ctx context.Context, groupID, inviterID, inviteeID *int64, status *po.RequestStatus, lastUpdatedDate *time.Time) (int64, error) {
	filter := r.buildFilter(groupID, inviterID, inviteeID, status, lastUpdatedDate)
	return r.coll.CountDocuments(ctx, filter)
}

func (r *groupInvitationRepository) buildFilter(groupID, inviterID, inviteeID *int64, status *po.RequestStatus, lastUpdatedDate *time.Time) bson.M {
	filter := bson.M{}
	if groupID != nil {
		filter["gid"] = *groupID
	}
	if inviterID != nil {
		filter["irid"] = *inviterID
	}
	if inviteeID != nil {
		filter["ieid"] = *inviteeID
	}
	if status != nil {
		filter["stat"] = *status
	}
	if lastUpdatedDate != nil {
		filter["cd"] = bson.M{"$gt": *lastUpdatedDate}
	}
	return filter
}

