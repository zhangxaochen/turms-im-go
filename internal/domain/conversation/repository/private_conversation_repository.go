package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"im.turms/server/internal/domain/conversation/po"
	turmsmongo "im.turms/server/internal/storage/mongo"
)

type PrivateConversationRepository struct {
	collection *mongo.Collection
}

func NewPrivateConversationRepository(client *turmsmongo.Client) *PrivateConversationRepository {
	return &PrivateConversationRepository{
		collection: client.Collection(po.CollectionNamePrivateConversation),
	}
}

// UpsertReadDate updates the read date for a private conversation, creating it if it doesn't exist.
func (r *PrivateConversationRepository) UpsertReadDate(ctx context.Context, ownerID int64, targetID int64, readDate time.Time) error {
	filter := po.PrivateConversationKey{
		OwnerID:  ownerID,
		TargetID: targetID,
	}
	update := bson.M{
		"$set": bson.M{
			"rd": readDate,
		},
	}
	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": filter}, update, opts)
	return err
}

// QueryPrivateConversations retrieves private conversations for given ownerIDs.
// @MappedFrom queryPrivateConversations(@NotNull Collection<Long> ownerIds, @NotNull Long targetId)
// @MappedFrom queryPrivateConversations(@NotNull Set<PrivateConversation.Key> keys)
func (r *PrivateConversationRepository) QueryPrivateConversations(ctx context.Context, ownerIDs []int64) ([]*po.PrivateConversation, error) {
	if len(ownerIDs) == 0 {
		return nil, nil
	}

	filter := bson.M{
		"_id.oid": bson.M{"$in": ownerIDs},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*po.PrivateConversation
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}
