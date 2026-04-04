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

const GroupVersionCollectionName = "groupVersion"

type GroupVersionRepository struct {
	client *turmsmongo.Client
	col    *mongo.Collection
}

func NewGroupVersionRepository(client *turmsmongo.Client) *GroupVersionRepository {
	return &GroupVersionRepository{
		client: client,
		col:    client.Collection(GroupVersionCollectionName),
	}
}

// InsertVersion creates a new group version record.
func (r *GroupVersionRepository) InsertVersion(ctx context.Context, groupID int64) error {
	now := time.Now()
	version := &po.GroupVersion{
		GroupID:       groupID,
		Members:       &now,
		Blocklist:     &now,
		JoinRequests:  &now,
		JoinQuestions: &now,
		Invitations:   &now,
	}
	_, err := r.col.InsertOne(ctx, version)
	return err
}

// UpdateVersion updates a specific version field.
// @MappedFrom updateVersion(@NotNull Long groupId, boolean updateMembers, boolean updateBlocklist, boolean joinRequests, boolean joinQuestions)
// @MappedFrom updateVersion(Long groupId, boolean updateMembers, boolean updateBlocklist, boolean joinRequests, boolean joinQuestions)
// @MappedFrom updateVersion(Long groupId, String field)
func (r *GroupVersionRepository) UpdateVersion(ctx context.Context, groupID int64, field string) error {
	filter := bson.M{"_id": groupID}
	update := bson.M{"$set": bson.M{field: time.Now()}}
	opts := options.Update().SetUpsert(true)

	_, err := r.col.UpdateOne(ctx, filter, update, opts)
	return err
}

// UpdateMembersVersion updates the members version.
// @MappedFrom updateMembersVersion(@NotNull Long groupId)
// @MappedFrom updateMembersVersion()
// @MappedFrom updateMembersVersion(@Nullable Set<Long> groupIds)
func (r *GroupVersionRepository) UpdateMembersVersion(ctx context.Context, groupID int64) error {
	return r.UpdateVersion(ctx, groupID, "mbr")
}

// UpdateInformationVersion updates the information version.
// @MappedFrom updateInformationVersion(@NotNull Long groupId)
func (r *GroupVersionRepository) UpdateInformationVersion(ctx context.Context, groupID int64) error {
	return r.UpdateVersion(ctx, groupID, "info")
}

// UpdateBlocklistVersion updates the blocklist version.
// @MappedFrom updateBlocklistVersion(@NotNull Long groupId)
func (r *GroupVersionRepository) UpdateBlocklistVersion(ctx context.Context, groupID int64) error {
	return r.UpdateVersion(ctx, groupID, "bl")
}

// UpdateJoinRequestsVersion updates the join requests version.
// @MappedFrom updateJoinRequestsVersion(@NotNull Long groupId)
func (r *GroupVersionRepository) UpdateJoinRequestsVersion(ctx context.Context, groupID int64) error {
	return r.UpdateVersion(ctx, groupID, "jr")
}

// UpdateJoinQuestionsVersion updates the join questions version.
// @MappedFrom updateJoinQuestionsVersion(@NotNull Long groupId)
func (r *GroupVersionRepository) UpdateJoinQuestionsVersion(ctx context.Context, groupID int64) error {
	return r.UpdateVersion(ctx, groupID, "jq")
}

// UpdateInvitationsVersion updates the invitations version.
func (r *GroupVersionRepository) UpdateInvitationsVersion(ctx context.Context, groupID int64) error {
	return r.UpdateVersion(ctx, groupID, "invt")
}

// Upsert creates or updates all group version records.
// @MappedFrom upsert(Long groupId, Collection<Long> memberIds, Date readDate)
// @MappedFrom upsert(Long groupId, Long memberId, Date readDate, boolean allowMoveReadDateForward)
// @MappedFrom upsert(Set<PrivateConversation.Key> keys, Date readDate, boolean allowMoveReadDateForward)
// @MappedFrom upsert(@NotNull Long groupId, @NotNull Date timestamp)
func (r *GroupVersionRepository) Upsert(ctx context.Context, groupID int64, timestamp time.Time) error {
	filter := bson.M{"_id": groupID}
	update := bson.M{
		"$set": bson.M{
			"mbr":  timestamp,
			"bl":   timestamp,
			"jr":   timestamp,
			"jq":   timestamp,
			"invt": timestamp,
			"info": timestamp,
		},
	}
	opts := options.Update().SetUpsert(true)
	_, err := r.col.UpdateOne(ctx, filter, update, opts)
	return err
}

// DeleteByIds deletes group versions by group IDs.
func (r *GroupVersionRepository) DeleteByIds(ctx context.Context, groupIDs []int64) error {
	if len(groupIDs) == 0 {
		return nil
	}
	filter := bson.M{"_id": bson.M{"$in": groupIDs}}
	_, err := r.col.DeleteMany(ctx, filter)
	return err
}

// FindVersion retrieves the group versions.
func (r *GroupVersionRepository) FindVersion(ctx context.Context, groupID int64) (*po.GroupVersion, error) {
	filter := bson.M{"_id": groupID}
	var version po.GroupVersion
	if err := r.col.FindOne(ctx, filter).Decode(&version); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &version, nil
}

// UpdateVersions updates a specific version field for multiple groups.
func (r *GroupVersionRepository) UpdateVersions(ctx context.Context, groupIDs []int64, field string) error {
	filter := bson.M{}
	if len(groupIDs) > 0 {
		filter["_id"] = bson.M{"$in": groupIDs}
	}
	update := bson.M{"$set": bson.M{field: time.Now()}}
	_, err := r.col.UpdateMany(ctx, filter, update)
	return err
}

// FindBlocklist retrieves the blocklist version.
func (r *GroupVersionRepository) FindBlocklist(ctx context.Context, groupID int64) (*time.Time, error) {
	return r.findSpecificVersion(ctx, groupID, "bl")
}

// FindJoinRequests retrieves the join requests version.
func (r *GroupVersionRepository) FindJoinRequests(ctx context.Context, groupID int64) (*time.Time, error) {
	return r.findSpecificVersion(ctx, groupID, "jr")
}

// FindJoinQuestions retrieves the join questions version.
func (r *GroupVersionRepository) FindJoinQuestions(ctx context.Context, groupID int64) (*time.Time, error) {
	return r.findSpecificVersion(ctx, groupID, "jq")
}

// FindMembers retrieves the members version.
func (r *GroupVersionRepository) FindMembers(ctx context.Context, groupID int64) (*time.Time, error) {
	return r.findSpecificVersion(ctx, groupID, "mbr")
}

func (r *GroupVersionRepository) findSpecificVersion(ctx context.Context, groupID int64, field string) (*time.Time, error) {
	filter := bson.M{"_id": groupID}
	opts := options.FindOne().SetProjection(bson.M{field: 1})
	var version po.GroupVersion
	if err := r.col.FindOne(ctx, filter, opts).Decode(&version); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	switch field {
	case "bl":
		return version.Blocklist, nil
	case "jr":
		return version.JoinRequests, nil
	case "jq":
		return version.JoinQuestions, nil
	case "mbr":
		return version.Members, nil
	}
	return nil, nil
}

