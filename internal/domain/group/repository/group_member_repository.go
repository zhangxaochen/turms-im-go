package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"im.turms/server/internal/domain/group/po"
	turmsmongo "im.turms/server/internal/storage/mongo"
	"im.turms/server/pkg/protocol"
)

const GroupMemberCollectionName = "groupMember"

type GroupMemberRepository struct {
	client *turmsmongo.Client
	col    *mongo.Collection
}

func NewGroupMemberRepository(client *turmsmongo.Client) *GroupMemberRepository {
	return &GroupMemberRepository{
		client: client,
		col:    client.Collection(GroupMemberCollectionName),
	}
}

// AddGroupMember adds a member to a group or updates their role.
func (r *GroupMemberRepository) AddGroupMember(ctx context.Context, member *po.GroupMember) error {
	opts := options.Update().SetUpsert(true)
	filter := bson.M{"_id": member.ID}
	update := bson.M{"$set": member}
	
	_, err := r.col.UpdateOne(ctx, filter, update, opts)
	return err
}

// RemoveGroupMember removes a member from a group.
func (r *GroupMemberRepository) RemoveGroupMember(ctx context.Context, groupID, userID int64) error {
	filter := bson.M{
		"_id": po.GroupMemberKey{GroupID: groupID, UserID: userID},
	}
	_, err := r.col.DeleteOne(ctx, filter)
	return err
}

// FindGroupMemberRole retrieves the role of a user in a group.
func (r *GroupMemberRepository) FindGroupMemberRole(ctx context.Context, groupID, userID int64) (*protocol.GroupMemberRole, error) {
	filter := bson.M{
		"_id": po.GroupMemberKey{GroupID: groupID, UserID: userID},
	}
	opts := options.FindOne().SetProjection(bson.M{"role": 1})

	var member po.GroupMember
	if err := r.col.FindOne(ctx, filter, opts).Decode(&member); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // User is not a member of the group
		}
		return nil, err
	}
	return &member.Role, nil
}

// FindGroupMemberIDs retrieves all user IDs within a group.
func (r *GroupMemberRepository) FindGroupMemberIDs(ctx context.Context, groupID int64) ([]int64, error) {
	filter := bson.M{
		"_id.gid": groupID,
	}
	opts := options.Find().SetProjection(bson.M{"_id.uid": 1})

	cursor, err := r.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var members []po.GroupMember
	if err := cursor.All(ctx, &members); err != nil {
		return nil, err
	}

	var userIDs []int64
	for _, m := range members {
		userIDs = append(userIDs, m.ID.UserID)
	}
	return userIDs, nil
}

// IsMemberMuted checks if a specific group member is currently muted.
func (r *GroupMemberRepository) IsMemberMuted(ctx context.Context, groupID, userID int64) (bool, error) {
	filter := bson.M{
		"_id": po.GroupMemberKey{GroupID: groupID, UserID: userID},
		"med": bson.M{"$gt": time.Now()},
	}

	count, err := r.col.CountDocuments(ctx, filter, options.Count().SetLimit(1))
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// FindUserJoinedGroupIDs retrieves all group IDs that a user belongs to.
func (r *GroupMemberRepository) FindUserJoinedGroupIDs(ctx context.Context, userID int64) ([]int64, error) {
	filter := bson.M{
		"_id.uid": userID,
	}
	opts := options.Find().SetProjection(bson.M{"_id.gid": 1})

	cursor, err := r.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var members []po.GroupMember
	if err := cursor.All(ctx, &members); err != nil {
		return nil, err
	}

	var groupIDs []int64
	for _, m := range members {
		groupIDs = append(groupIDs, m.ID.GroupID)
	}
	return groupIDs, nil
}
