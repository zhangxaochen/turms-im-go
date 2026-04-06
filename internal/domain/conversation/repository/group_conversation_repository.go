package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"im.turms/server/internal/domain/conversation/po"
	turmsmongo "im.turms/server/internal/storage/mongo"
)

type GroupConversationRepository struct {
	collection *mongo.Collection
}

func NewGroupConversationRepository(client *turmsmongo.Client) *GroupConversationRepository {
	return &GroupConversationRepository{
		collection: client.Collection(po.CollectionNameGroupConversation),
	}
}

// Upsert updates a specific member's read date in a group conversation.
func (r *GroupConversationRepository) Upsert(ctx context.Context, groupID int64, memberID int64, readDate time.Time, allowMoveForward bool) error {
	filter := bson.M{"_id": groupID}
	fieldKey := fmt.Sprintf("mr.%d", memberID)
	if !allowMoveForward {
		// Bug fix: Java parity — Java's ltOrNull produces:
		//   {$or: [{fieldKey: null}, {fieldKey: {$lt: readDate}}]}
		// where {fieldKey: null} matches both missing fields AND explicitly null values.
		// Previously Go used $exists: false which only matches missing fields.
		filter["$or"] = []bson.M{
			{fieldKey: nil},
			{fieldKey: bson.M{"$lt": readDate}},
		}
	}

	update := bson.M{
		"$set": bson.M{
			fieldKey: readDate,
		},
	}
	opts := options.Update().SetUpsert(true)

	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

// BulkUpsert updates multiple members' read dates in a group conversation.
func (r *GroupConversationRepository) BulkUpsert(ctx context.Context, groupID int64, memberIDs []int64, readDate time.Time) error {
	if len(memberIDs) == 0 {
		return nil
	}
	filter := bson.M{"_id": groupID}
	updateSet := bson.M{}
	for _, memberID := range memberIDs {
		fieldKey := fmt.Sprintf("mr.%d", memberID)
		updateSet[fieldKey] = readDate
	}
	update := bson.M{"$set": updateSet}
	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

// FindByIds retrieves the conversations for the given groupIDs.
// @MappedFrom queryGroupConversations(@NotNull Collection<Long> groupIds)
func (r *GroupConversationRepository) FindByIds(ctx context.Context, groupIDs []int64) ([]*po.GroupConversation, error) {
	if len(groupIDs) == 0 {
		return nil, nil
	}

	filter := bson.M{"_id": bson.M{"$in": groupIDs}}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*po.GroupConversation
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

// DeleteMemberConversations removes a member's read date from multiple groups.
// @MappedFrom deleteMemberConversations(Collection<Long> groupIds, Long memberId, ClientSession session)
func (r *GroupConversationRepository) DeleteMemberConversations(ctx context.Context, groupIDs []int64, memberID int64) error {
	if len(groupIDs) == 0 {
		return nil
	}
	filter := bson.M{"_id": bson.M{"$in": groupIDs}}
	fieldKey := fmt.Sprintf("mr.%d", memberID)
	update := bson.M{
		"$unset": bson.M{
			fieldKey: "",
		},
	}
	_, err := r.collection.UpdateMany(ctx, filter, update)
	return err
}

func (r *GroupConversationRepository) DeleteByIds(ctx context.Context, groupIDs []int64) (*mongo.DeleteResult, error) {
	if len(groupIDs) == 0 {
		return &mongo.DeleteResult{}, nil
	}
	filter := bson.M{"_id": bson.M{"$in": groupIDs}}
	return r.collection.DeleteMany(ctx, filter)
}
