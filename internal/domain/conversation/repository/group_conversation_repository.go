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

// UpsertReadDate updates a specific member's read date in a group conversation.
func (r *GroupConversationRepository) UpsertReadDate(ctx context.Context, groupID int64, memberID int64, readDate time.Time) error {
	filter := bson.M{"_id": groupID}

	// MongoDB field path for the specific member's read date
	// mr stands for memberIdToReadDate
	memberField := fmt.Sprintf("mr.%d", memberID)

	update := bson.M{
		"$set": bson.M{
			memberField: readDate,
		},
	}
	opts := options.Update().SetUpsert(true)

	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

// QueryGroupConversations retrieves the conversations for the given groupIDs.
// @MappedFrom queryGroupConversations(@NotNull Collection<Long> groupIds)
func (r *GroupConversationRepository) QueryGroupConversations(ctx context.Context, groupIDs []int64) ([]*po.GroupConversation, error) {
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
