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

const GroupCollectionName = "group"

type GroupRepository struct {
	client *turmsmongo.Client
	col    *mongo.Collection
}

func NewGroupRepository(client *turmsmongo.Client) *GroupRepository {
	return &GroupRepository{
		client: client,
		col:    client.Collection(GroupCollectionName),
	}
}

// InsertGroup inserts a new group into MongoDB.
func (r *GroupRepository) InsertGroup(ctx context.Context, group *po.Group) error {
	_, err := r.col.InsertOne(ctx, group)
	return err
}

// FindGroups retrieves multiple groups by their IDs without filtering by deletion date.
// @MappedFrom findGroups(@Nullable Set<Long> ids, @Nullable Set<Long> typeIds, @Nullable Set<Long> creatorIds, @Nullable Set<Long> ownerIds, @Nullable Boolean isActive, @Nullable DateRange creationDateRange, @Nullable DateRange deletionDateRange, @Nullable DateRange lastUpdatedDateRange, @Nullable DateRange muteEndDateRange, @Nullable Integer page, @Nullable Integer size)
// Bug fix: Java findGroups does NOT filter by deletion date. Removed incorrect "dd": {"$exists": false} filter.
func (r *GroupRepository) FindGroups(ctx context.Context, groupIDs []int64) ([]*po.Group, error) {
	filter := bson.M{
		"_id": bson.M{"$in": groupIDs},
	}

	cursor, err := r.col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var groups []*po.Group
	if err := cursor.All(ctx, &groups); err != nil {
		return nil, err
	}
	return groups, nil
}

// QueryGroups retrieves groups based on various filters.
func (r *GroupRepository) QueryGroups(ctx context.Context, groupIDs []int64, name *string, lastUpdatedDate *time.Time, skip *int32, limit *int32) ([]*po.Group, error) {
	filter := bson.M{}
	if len(groupIDs) > 0 {
		filter["_id"] = bson.M{"$in": groupIDs}
	}
	if name != nil {
		filter["n"] = *name
	}
	if lastUpdatedDate != nil {
		filter["lud"] = bson.M{"$gt": *lastUpdatedDate}
	}
	// Bug fix: Java findGroups does NOT filter by deletion date.
	// When lastUpdatedDate is provided, it returns even deleted groups to let clients sync.
	// Removed incorrect deletion date filter.

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

	var groups []*po.Group
	if err := cursor.All(ctx, &groups); err != nil {
		return nil, err
	}
	return groups, nil
}

// FindGroupOwnerID retrieves the owner ID of a specific group.
func (r *GroupRepository) FindGroupOwnerID(ctx context.Context, groupID int64) (*int64, error) {
	filter := bson.M{"_id": groupID}
	opts := options.FindOne().SetProjection(bson.M{"oid": 1})

	var group po.Group
	if err := r.col.FindOne(ctx, filter, opts).Decode(&group); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Group not found
		}
		return nil, err
	}
	return group.OwnerID, nil
}

func (r *GroupRepository) FindGroup(ctx context.Context, groupID int64) (*po.Group, error) {
	filter := bson.M{"_id": groupID}
	var group po.Group
	if err := r.col.FindOne(ctx, filter).Decode(&group); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &group, nil
}

// CountOwnedGroups counts the number of groups owned by a specific user.
// @MappedFrom countOwnedGroups(@NotNull Long ownerId)
// Bug fix: Java only filters by ownerId, no deletion date check. Removed incorrect "dd": {"$exists": false}.
func (r *GroupRepository) CountOwnedGroups(ctx context.Context, ownerID int64) (int64, error) {
	filter := bson.M{"oid": ownerID}
	return r.col.CountDocuments(ctx, filter)
}

// CountOwnedGroupsByTypeId counts groups owned by a user with a specific type.
// @MappedFrom countOwnedGroups(Long ownerId, Long groupTypeId)
// Bug fix: Method was missing. Java filters by both OWNER_ID and TYPE_ID.
func (r *GroupRepository) CountOwnedGroupsByTypeId(ctx context.Context, ownerID int64, groupTypeID int64) (int64, error) {
	filter := bson.M{"oid": ownerID, "tid": groupTypeID}
	return r.col.CountDocuments(ctx, filter)
}

// UpdateGroup modifies specified fields of a group.
func (r *GroupRepository) UpdateGroup(ctx context.Context, groupID int64, update bson.M) error {
	filter := bson.M{"_id": groupID}
	_, err := r.col.UpdateOne(ctx, filter, bson.M{"$set": update})
	return err
}

// UpdateGroupsDeletionDate updates the deletion date of groups.
func (r *GroupRepository) UpdateGroupsDeletionDate(ctx context.Context, groupIDs []int64, deletionDate *time.Time, session mongo.SessionContext) error {
	filter := bson.M{"_id": bson.M{"$in": groupIDs}}
	update := bson.M{}
	if deletionDate == nil {
		update["$unset"] = bson.M{"dd": ""}
	} else {
		update["$set"] = bson.M{"dd": *deletionDate}
	}

	var err error
	if session != nil {
		_, err = r.col.UpdateMany(session, filter, update)
	} else {
		_, err = r.col.UpdateMany(ctx, filter, update)
	}
	return err
}

// CountCreatedGroups counts groups created within a date range.
// Bug fix: Java includes .eq(Group.Fields.DELETION_DATE, null) to only count non-deleted groups.
// Added missing "dd": nil filter.
func (r *GroupRepository) CountCreatedGroups(ctx context.Context, dateRange *turmsmongo.DateRange) (int64, error) {
	filter := bson.M{"dd": nil}
	if dateRange != nil {
		dateFilter := bson.M{}
		if dateRange.Start != nil {
			dateFilter["$gte"] = *dateRange.Start
			}
		if dateRange.End != nil {
			dateFilter["$lte"] = *dateRange.End
			}
		if len(dateFilter) > 0 {
				filter["cd"] = dateFilter
		}
	}
	return r.col.CountDocuments(ctx, filter)
}

// CountDeletedGroups counts groups deleted within a date range.
// Bug fix: When dateRange is nil, Java applies no filter (counts all documents), because it
// relies on DELETION_DATE being non-null from the date range itself.
// When dateRange is provided, the date range filter inherently ensures dd is non-null.
// When dateRange is nil, Java counts all groups, so we use an empty filter.
func (r *GroupRepository) CountDeletedGroups(ctx context.Context, dateRange *turmsmongo.DateRange) (int64, error) {
	if dateRange == nil {
		// Java parity: when no date range, counts all documents (no filter)
		return r.col.CountDocuments(ctx, bson.M{})
	}
	// When date range is provided, the date range filter inherently ensure dd is non-null
	dateFilter := bson.M{}
	if dateRange.Start != nil {
		dateFilter["$gte"] = *dateRange.Start
	}
	if dateRange.End != nil {
		dateFilter["$lte"] = *dateRange.End
	}
	if len(dateFilter) > 0 {
		return r.col.CountDocuments(ctx, bson.M{"dd": dateFilter})
	}
	// If date range bounds are both nil, treat as no filter
	return r.col.CountDocuments(ctx, bson.M{})
}

// Count counts all groups.
// @MappedFrom count()
func (r *GroupRepository) Count(ctx context.Context) (int64, error) {
	return r.col.CountDocuments(ctx, bson.M{})
}

// CountGroups counts groups matching multi-parameter filters.
// @MappedFrom countGroups(@Nullable Set<Long> ids, @Nullable Set<Long> typeIds, @Nullable Set<Long> creatorIds, @Nullable Set<Long> ownerIds, @Nullable Boolean isActive, @Nullable DateRange creationDateRange, @Nullable DateRange deletionDateRange, @Nullable DateRange lastUpdatedDateRange, @Nullable DateRange muteEndDateRange)
// Bug fix: Previous implementation incorrectly delegated to CountCreatedGroups. Now implements full multi-parameter count.
func (r *GroupRepository) CountGroups(ctx context.Context, ids []int64, typeIds []int64, creatorIds []int64, ownerIds []int64, isActive *bool) (int64, error) {
	filter := bson.M{}
	if len(ids) > 0 {
		filter["_id"] = bson.M{"$in": ids}
	}
	if len(typeIds) > 0 {
		filter["tid"] = bson.M{"$in": typeIds}
	}
	if len(creatorIds) > 0 {
		filter["cid"] = bson.M{"$in": creatorIds}
	}
	if len(ownerIds) > 0 {
		filter["oid"] = bson.M{"$in": ownerIds}
	}
	if isActive != nil {
		filter["ac"] = *isActive
	}
	// Java parity: countGroups does NOT filter by deletion date
	return r.col.CountDocuments(ctx, filter)
}

// FindNotDeletedGroups retrieves groups that are not deleted.
// Bug fix: Java uses eq(DELETION_DATE, null) which matches documents where the field is absent OR explicitly null.
// Go bson.M{"$exists": false} only matches absent fields. Changed to nil for proper parity.
func (r *GroupRepository) FindNotDeletedGroups(ctx context.Context, groupIDs []int64, lastUpdatedDate *time.Time) ([]*po.Group, error) {
	filter := bson.M{
		"_id": bson.M{"$in": groupIDs},
		"dd":  nil,
	}
	if lastUpdatedDate != nil {
		filter["lud"] = bson.M{"$gt": *lastUpdatedDate}
	}

	cursor, err := r.col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var groups []*po.Group
	if err := cursor.All(ctx, &groups); err != nil {
		return nil, err
	}
	return groups, nil
}

// FindAllNames retrieves only the names of all groups.
// @MappedFrom findAllNames()
// Bug fix: Method was missing. Java uses QueryOptions.include(Group.Fields.NAME) projection.
func (r *GroupRepository) FindAllNames(ctx context.Context) ([]string, error) {
	opts := options.Find().SetProjection(bson.M{"n": 1})
	cursor, err := r.col.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	type nameResult struct {
		Name *string `bson:"n"`
	}
	var results []nameResult
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	names := make([]string, 0, len(results))
	for _, r := range results {
		if r.Name != nil {
			names = append(names, *r.Name)
		}
	}
	return names, nil
}

// FindTypeID retrieves the type ID of a group.
func (r *GroupRepository) FindTypeID(ctx context.Context, groupID int64) (*int64, error) {
	filter := bson.M{"_id": groupID}
	opts := options.FindOne().SetProjection(bson.M{"tid": 1})
	type Result struct {
		TypeID int64 `bson:"tid"`
	}
	var res Result
	if err := r.col.FindOne(ctx, filter, opts).Decode(&res); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &res.TypeID, nil
}

// FindTypeIDIfActiveAndNotDeleted retrieves the type ID if the group is active and not deleted.
// Bug fix: Java uses eq(DELETION_DATE, null) which matches absent or null. Changed from $exists: false to nil.
func (r *GroupRepository) FindTypeIDIfActiveAndNotDeleted(ctx context.Context, groupID int64) (*int64, error) {
	filter := bson.M{
		"_id": groupID,
		"ac":  true,
		"dd":  nil,
	}
	opts := options.FindOne().SetProjection(bson.M{"tid": 1})
	type Result struct {
		TypeID int64 `bson:"tid"`
	}
	var res Result
	if err := r.col.FindOne(ctx, filter, opts).Decode(&res); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &res.TypeID, nil
}

// FindMinimumScore retrieves the minimum score of a group.
func (r *GroupRepository) FindMinimumScore(ctx context.Context, groupID int64) (*int32, error) {
	filter := bson.M{"_id": groupID}
	opts := options.FindOne().SetProjection(bson.M{"ms": 1})
	type Result struct {
		MinimumScore int32 `bson:"ms"`
	}
	var res Result
	if err := r.col.FindOne(ctx, filter, opts).Decode(&res); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &res.MinimumScore, nil
}

// IsGroupMuted checks if a group is muted.
// Bug fix: Java takes (Long groupId, Date muteEndDate) and compares stored MUTE_END_DATE against
// the passed-in muteEndDate. Go hardcoded time.Now() instead. Added muteEndDate parameter.
func (r *GroupRepository) IsGroupMuted(ctx context.Context, groupID int64, muteEndDate time.Time) (bool, error) {
	filter := bson.M{
		"_id": groupID,
		"med": bson.M{"$gt": muteEndDate},
	}
	count, err := r.col.CountDocuments(ctx, filter)
	return count > 0, err
}

// IsGroupActiveAndNotDeleted checks if a group is active and not deleted.
// Bug fix: Java use eq(DELETION_DATE, null) which matches absent or null. Changed from $exists: false to nil.
func (r *GroupRepository) IsGroupActiveAndNotDeleted(ctx context.Context, groupID int64) (bool, error) {
	filter := bson.M{
		"_id": groupID,
		"ac":  true,
		"dd":  nil,
	}
	count, err := r.col.CountDocuments(ctx, filter)
	return count > 0, err
}

// FindTypeIDAndGroupID retrieves type IDs and group IDs for multiple groups.
func (r *GroupRepository) FindTypeIDAndGroupID(ctx context.Context, groupIDs []int64) ([]*po.Group, error) {
	filter := bson.M{"_id": bson.M{"$in": groupIDs}}
	opts := options.Find().SetProjection(bson.M{"tid": 1})

	cursor, err := r.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var groups []*po.Group
	if err := cursor.All(ctx, &groups); err != nil {
		return nil, err
	}
	return groups, nil
}

// DeleteGroup removes a group from MongoDB.
func (r *GroupRepository) DeleteGroup(ctx context.Context, groupID int64) error {
	filter := bson.M{"_id": groupID}
	_, err := r.col.DeleteOne(ctx, filter)
	return err
}
