package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"im.turms/server/internal/domain/conference/po"
	turmsmongo "im.turms/server/internal/storage/mongo"
)

// MeetingRepository maps to MeetingRepository.java
// @MappedFrom MeetingRepository
type MeetingRepository struct {
	client *turmsmongo.Client
	col    *mongo.Collection
}

func NewMeetingRepository(client *turmsmongo.Client) *MeetingRepository {
	return &MeetingRepository{
		client: client,
		col:    client.Collection(po.MeetingCollectionName),
	}
}

// Insert inserts a new meeting into MongoDB.
func (r *MeetingRepository) Insert(ctx context.Context, meeting *po.Meeting) error {
	_, err := r.col.InsertOne(ctx, meeting)
	return err
}

// FindByID retrieves a meeting by its ID.
func (r *MeetingRepository) FindByID(ctx context.Context, meetingID int64) (*po.Meeting, error) {
	var meeting po.Meeting
	err := r.col.FindOne(ctx, bson.M{"_id": meetingID}).Decode(&meeting)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &meeting, err
}

// @MappedFrom updateEndDate(Long meetingId, Date endDate)
func (r *MeetingRepository) UpdateEndDate(ctx context.Context, meetingID int64, endDate time.Time) error {
	filter := bson.M{"_id": meetingID}
	update := bson.M{"$set": bson.M{"ed": endDate}}
	_, err := r.col.UpdateOne(ctx, filter, update)
	return err
}

// @MappedFrom updateCancelDateIfNotCanceled(Long meetingId, Date cancelDate)
func (r *MeetingRepository) UpdateCancelDateIfNotCanceled(ctx context.Context, meetingID int64, cancelDate time.Time) (bool, error) {
	// Bug fix: Use $eq: nil instead of $exists: false to match Java's eq(CANCEL_DATE, null).
	// Java's eq(null) matches documents where the field is absent OR explicitly null.
	// $exists: false only matches documents where the field doesn't exist at all.
	filter := bson.M{
		"_id": meetingID,
		"ccd": nil,
	}
	update := bson.M{"$set": bson.M{"ccd": cancelDate}}
	res, err := r.col.UpdateOne(ctx, filter, update)
	if err != nil {
		return false, err
	}
	return res.ModifiedCount > 0, nil
}

// @MappedFrom updateMeeting(Long meetingId, @Nullable String name, @Nullable String intro, @Nullable String password)
func (r *MeetingRepository) UpdateMeeting(ctx context.Context, meetingID int64, name *string, intro *string, password *string) error {
	_, err := r.UpdateMeetingWithResult(ctx, meetingID, name, intro, password)
	return err
}

// UpdateMeetingWithResult updates a meeting and returns whether any document was modified.
// Bug fix: Java checks updateResult.getModifiedCount() > 0 and returns FAILED if 0 rows modified.
func (r *MeetingRepository) UpdateMeetingWithResult(ctx context.Context, meetingID int64, name *string, intro *string, password *string) (bool, error) {
	filter := bson.M{"_id": meetingID}
	set := bson.M{}
	if name != nil {
		set["n"] = *name
	}
	if intro != nil {
		set["intro"] = *intro
	}
	if password != nil {
		set["pw"] = *password
	}
	if len(set) == 0 {
		return false, nil
	}
	res, err := r.col.UpdateOne(ctx, filter, bson.M{"$set": set})
	if err != nil {
		return false, err
	}
	return res.ModifiedCount > 0, nil
}

// @MappedFrom find(@Nullable Collection<Long> ids, @Nullable Collection<Long> creatorIds, @Nullable Collection<Long> userIds, @Nullable Collection<Long> groupIds, @Nullable Date creationDateStart, @Nullable Date creationDateEnd, @Nullable Integer skip, @Nullable Integer limit)
func (r *MeetingRepository) Find(ctx context.Context,
	ids []int64,
	creatorIDs []int64,
	userIDs []int64,
	groupIDs []int64,
	creationDateStart *time.Time,
	creationDateEnd *time.Time,
	skip *int32,
	limit *int32,
) ([]*po.Meeting, error) {
	filter := bson.M{}
	if len(ids) > 0 {
		filter["_id"] = bson.M{"$in": ids}
	}
	if len(creatorIDs) > 0 {
		filter["cid"] = bson.M{"$in": creatorIDs}
	}
	if len(userIDs) > 0 {
		filter["uid"] = bson.M{"$in": userIDs}
	}
	if len(groupIDs) > 0 {
		filter["gid"] = bson.M{"$in": groupIDs}
	}
	if creationDateStart != nil || creationDateEnd != nil {
		dateFilter := bson.M{}
		if creationDateStart != nil {
			dateFilter["$gte"] = *creationDateStart
		}
		if creationDateEnd != nil {
			dateFilter["$lte"] = *creationDateEnd
		}
		filter["cd"] = dateFilter
	}

	opts := options.Find()
	if skip != nil {
		opts.SetSkip(int64(*skip))
	}
	if limit != nil {
		opts.SetLimit(int64(*limit))
	}

	cursor, err := r.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var meetings []*po.Meeting
	if err := cursor.All(ctx, &meetings); err != nil {
		return nil, err
	}
	return meetings, nil
}

// @MappedFrom find(@Nullable Collection<Long> ids, @NotNull Long creatorId, @NotNull Long userId, @Nullable Date creationDateStart, @Nullable Date creationDateEnd, @Nullable Integer skip, @Nullable Integer limit)
func (r *MeetingRepository) FindByCreatorAndUser(ctx context.Context,
	ids []int64,
	creatorID int64,
	userID int64,
	creationDateStart *time.Time,
	creationDateEnd *time.Time,
	skip *int32,
	limit *int32,
) ([]*po.Meeting, error) {
	// Bug fix: Java uses .or(Filter.newBuilder(2).eq(CREATOR_ID, creatorId).eq(USER_ID, userId))
	// which is actually AND inside a single OR clause (one group), equivalent to matching
	// documents where BOTH creatorId AND userId match. Go incorrectly used $or with two
	// separate conditions (either match). Fix: use top-level AND (both must match).
	filter := bson.M{
		"cid": creatorID,
		"uid": userID,
	}
	if len(ids) > 0 {
		filter["_id"] = bson.M{"$in": ids}
	}
	if creationDateStart != nil || creationDateEnd != nil {
		dateFilter := bson.M{}
		if creationDateStart != nil {
			dateFilter["$gte"] = *creationDateStart
		}
		if creationDateEnd != nil {
			dateFilter["$lte"] = *creationDateEnd
		}
		filter["cd"] = dateFilter
	}

	opts := options.Find()
	if skip != nil {
		opts.SetSkip(int64(*skip))
	}
	if limit != nil {
		opts.SetLimit(int64(*limit))
	}

	cursor, err := r.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var meetings []*po.Meeting
	if err := cursor.All(ctx, &meetings); err != nil {
		return nil, err
	}
	return meetings, nil
}
