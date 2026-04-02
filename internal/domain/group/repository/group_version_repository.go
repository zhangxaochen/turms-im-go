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
func (r *GroupVersionRepository) UpdateVersion(ctx context.Context, groupID int64, field string) error {
	filter := bson.M{"_id": groupID}
	update := bson.M{"$set": bson.M{field: time.Now()}}
	opts := options.Update().SetUpsert(true)

	_, err := r.col.UpdateOne(ctx, filter, update, opts)
	return err
}

// UpdateMembersVersion updates the members version.
func (r *GroupVersionRepository) UpdateMembersVersion(ctx context.Context, groupID int64) error {
	return r.UpdateVersion(ctx, groupID, "mbr")
}

// UpdateBlocklistVersion updates the blocklist version.
func (r *GroupVersionRepository) UpdateBlocklistVersion(ctx context.Context, groupID int64) error {
	return r.UpdateVersion(ctx, groupID, "bl")
}

// UpdateJoinRequestsVersion updates the join requests version.
func (r *GroupVersionRepository) UpdateJoinRequestsVersion(ctx context.Context, groupID int64) error {
	return r.UpdateVersion(ctx, groupID, "jr")
}

// UpdateJoinQuestionsVersion updates the join questions version.
func (r *GroupVersionRepository) UpdateJoinQuestionsVersion(ctx context.Context, groupID int64) error {
	return r.UpdateVersion(ctx, groupID, "jq")
}

// UpdateInvitationsVersion updates the invitations version.
func (r *GroupVersionRepository) UpdateInvitationsVersion(ctx context.Context, groupID int64) error {
	return r.UpdateVersion(ctx, groupID, "invt")
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
