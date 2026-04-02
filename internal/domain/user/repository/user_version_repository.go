package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"im.turms/server/internal/domain/user/po"
	turmsmongo "im.turms/server/internal/storage/mongo"
)

type UserVersionRepository interface {
	UpsertEmptyUserVersion(ctx context.Context, userID int64) error
	UpdateUserVersion(ctx context.Context, userID int64, update interface{}) error
	UpdateUserVersions(ctx context.Context, userIDs []int64, update interface{}) error
	FindUserVersion(ctx context.Context, userID int64) (*po.UserVersion, error)
	DeleteUserVersion(ctx context.Context, userID int64) error
	DeleteUserVersions(ctx context.Context, userIDs []int64) error
	UpdateSpecificVersion(ctx context.Context, userIDs []int64, field string, updateDate time.Time) error
	FindGroupJoinRequestsVersion(ctx context.Context, userID int64) (*time.Time, error)
	FindJoinedGroupVersion(ctx context.Context, userID int64) (*time.Time, error)
	FindReceivedGroupInvitationsVersion(ctx context.Context, userID int64) (*time.Time, error)
	FindRelationshipsVersion(ctx context.Context, userID int64) (*time.Time, error)
	FindRelationshipGroupsVersion(ctx context.Context, userID int64) (*time.Time, error)
	FindSentGroupInvitationsVersion(ctx context.Context, userID int64) (*time.Time, error)
	FindSentFriendRequestsVersion(ctx context.Context, userID int64) (*time.Time, error)
	FindReceivedFriendRequestsVersion(ctx context.Context, userID int64) (*time.Time, error)
}

type userVersionRepository struct {
	collection *mongo.Collection
}

func NewUserVersionRepository(mongoClient *turmsmongo.Client) UserVersionRepository {
	return &userVersionRepository{
		collection: mongoClient.Collection(po.CollectionNameUserVersion),
	}
}

// @MappedFrom upsertEmptyUserVersion(@NotNull Long userId, @NotNull Date timestamp, @Nullable ClientSession session)
func (r *userVersionRepository) UpsertEmptyUserVersion(ctx context.Context, userID int64) error {
	update := map[string]interface{}{
		"$setOnInsert": po.UserVersion{
			UserID: userID,
		},
	}
	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateOne(ctx, map[string]interface{}{"_id": userID}, update, opts)
	return err
}

func (r *userVersionRepository) UpdateUserVersion(ctx context.Context, userID int64, update interface{}) error {
	_, err := r.collection.UpdateOne(ctx, map[string]interface{}{"_id": userID}, update)
	return err
}

func (r *userVersionRepository) FindUserVersion(ctx context.Context, userID int64) (*po.UserVersion, error) {
	var version po.UserVersion
	err := r.collection.FindOne(ctx, map[string]interface{}{"_id": userID}).Decode(&version)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &version, nil
}

func (r *userVersionRepository) DeleteUserVersion(ctx context.Context, userID int64) error {
	_, err := r.collection.DeleteOne(ctx, map[string]interface{}{"_id": userID})
	return err
}

func (r *userVersionRepository) UpdateUserVersions(ctx context.Context, userIDs []int64, update interface{}) error {
	filter := map[string]interface{}{
		"_id": map[string]interface{}{
			"$in": userIDs,
		},
	}
	_, err := r.collection.UpdateMany(ctx, filter, update)
	return err
}

func (r *userVersionRepository) DeleteUserVersions(ctx context.Context, userIDs []int64) error {
	filter := map[string]interface{}{
		"_id": map[string]interface{}{
			"$in": userIDs,
		},
	}
	_, err := r.collection.DeleteMany(ctx, filter)
	return err
}

// @MappedFrom updateSpecificVersion(@NotEmpty Set<Long> userIds, @Nullable ClientSession session, @NotEmpty String... fields)
// @MappedFrom updateSpecificVersion(@Nullable Set<Long> groupIds, @NotNull String field)
// @MappedFrom updateSpecificVersion(Long userId, @Nullable ClientSession session, String... fields)
// @MappedFrom updateSpecificVersion(@NotNull Long userId, @Nullable ClientSession session, @NotEmpty String... fields)
// @MappedFrom updateSpecificVersion(Long userId, @Nullable ClientSession session, String field)
// @MappedFrom updateSpecificVersion(@NotNull Long groupId, @NotNull String field)
// @MappedFrom updateSpecificVersion(@NotNull String field)
// @MappedFrom updateSpecificVersion(@NotNull Long userId, @Nullable ClientSession session, @NotNull String field)
// @MappedFrom updateSpecificVersion(Set<Long> userIds, @Nullable ClientSession session, String... fields)
func (r *userVersionRepository) UpdateSpecificVersion(ctx context.Context, userIDs []int64, field string, updateDate time.Time) error {
	if len(userIDs) == 0 {
		return nil
	}
	filter := map[string]interface{}{
		"_id": map[string]interface{}{
			"$in": userIDs,
		},
	}
	update := map[string]interface{}{
		"$set": map[string]interface{}{
			field: updateDate,
		},
	}
	_, err := r.collection.UpdateMany(ctx, filter, update)
	return err
}

func (r *userVersionRepository) findSpecificVersion(ctx context.Context, userID int64) (*po.UserVersion, error) {
	var version po.UserVersion
	err := r.collection.FindOne(ctx, map[string]interface{}{"_id": userID}).Decode(&version)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Not found
		}
		return nil, err
	}
	return &version, nil
}

func (r *userVersionRepository) FindGroupJoinRequestsVersion(ctx context.Context, userID int64) (*time.Time, error) {
	v, err := r.findSpecificVersion(ctx, userID)
	if err != nil || v == nil {
		return nil, err
	}
	return (*time.Time)(&v.GroupJoinRequests), nil
}

func (r *userVersionRepository) FindJoinedGroupVersion(ctx context.Context, userID int64) (*time.Time, error) {
	v, err := r.findSpecificVersion(ctx, userID)
	if err != nil || v == nil {
		return nil, err
	}
	return (*time.Time)(&v.JoinedGroups), nil
}

func (r *userVersionRepository) FindReceivedGroupInvitationsVersion(ctx context.Context, userID int64) (*time.Time, error) {
	v, err := r.findSpecificVersion(ctx, userID)
	if err != nil || v == nil {
		return nil, err
	}
	return (*time.Time)(&v.ReceivedGroupInvitations), nil
}

func (r *userVersionRepository) FindRelationshipsVersion(ctx context.Context, userID int64) (*time.Time, error) {
	v, err := r.findSpecificVersion(ctx, userID)
	if err != nil || v == nil {
		return nil, err
	}
	return (*time.Time)(&v.Relationships), nil
}

func (r *userVersionRepository) FindRelationshipGroupsVersion(ctx context.Context, userID int64) (*time.Time, error) {
	v, err := r.findSpecificVersion(ctx, userID)
	if err != nil || v == nil {
		return nil, err
	}
	return (*time.Time)(&v.RelationshipGroups), nil
}

func (r *userVersionRepository) FindSentGroupInvitationsVersion(ctx context.Context, userID int64) (*time.Time, error) {
	v, err := r.findSpecificVersion(ctx, userID)
	if err != nil || v == nil {
		return nil, err
	}
	return (*time.Time)(&v.SentGroupInvitations), nil
}

func (r *userVersionRepository) FindSentFriendRequestsVersion(ctx context.Context, userID int64) (*time.Time, error) {
	v, err := r.findSpecificVersion(ctx, userID)
	if err != nil || v == nil {
		return nil, err
	}
	return (*time.Time)(&v.SentFriendRequests), nil
}

func (r *userVersionRepository) FindReceivedFriendRequestsVersion(ctx context.Context, userID int64) (*time.Time, error) {
	v, err := r.findSpecificVersion(ctx, userID)
	if err != nil || v == nil {
		return nil, err
	}
	return (*time.Time)(&v.ReceivedFriendRequests), nil
}
