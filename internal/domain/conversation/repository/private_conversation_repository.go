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

// Upsert updates the read date for multiple private conversations, creating them if they don't exist.
// @MappedFrom upsert(@NotNull Set<PrivateConversation.Key> keys, @NotNull Date readDate, boolean allowMoveReadDateForward)
func (r *PrivateConversationRepository) Upsert(ctx context.Context, keys []po.PrivateConversationKey, readDate time.Time, allowMoveForward bool) error {
	if len(keys) == 0 {
		return nil
	}
	filter := bson.M{
		"_id": bson.M{"$in": keys},
	}
	if !allowMoveForward {
		filter["$or"] = []bson.M{
			{"rd": bson.M{"$lt": readDate}},
			{"rd": bson.M{"$exists": false}},
		}
	}
	update := bson.M{
		"$set": bson.M{
			"rd": readDate,
		},
	}
	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateMany(ctx, filter, update, opts)
	return err
}

// FindByIds retrieves private conversations for given keys.
func (r *PrivateConversationRepository) FindByIds(ctx context.Context, keys []po.PrivateConversationKey) ([]*po.PrivateConversation, error) {
	if len(keys) == 0 {
		return nil, nil
	}
	filter := bson.M{"_id": bson.M{"$in": keys}}
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

// FindConversations retrieves private conversations for given ownerIDs.
// @MappedFrom findConversations(Collection<Long> ownerIds)
func (r *PrivateConversationRepository) FindConversations(ctx context.Context, ownerIDs []int64) ([]*po.PrivateConversation, error) {
	if len(ownerIDs) == 0 {
		return nil, nil
	}
	filter := bson.M{"_id.oid": bson.M{"$in": ownerIDs}}
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

// DeleteConversationsByOwnerIds deletes all private conversations for given ownerIDs.
// @MappedFrom deleteConversationsByOwnerIds(Set<Long> ownerIds, @Nullable ClientSession session)
func (r *PrivateConversationRepository) DeleteConversationsByOwnerIds(ctx context.Context, ownerIDs []int64) error {
	if len(ownerIDs) == 0 {
		return nil
	}
	filter := bson.M{"_id.oid": bson.M{"$in": ownerIDs}}
	_, err := r.collection.DeleteMany(ctx, filter)
	return err
}

func (r *PrivateConversationRepository) DeleteByIds(ctx context.Context, keys []po.PrivateConversationKey) error {
	if len(keys) == 0 {
		return nil
	}
	filter := bson.M{"_id": bson.M{"$in": keys}}
	_, err := r.collection.DeleteMany(ctx, filter)
	return err
}
