package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"im.turms/server/internal/domain/message/po"
	turmsmongo "im.turms/server/internal/storage/mongo"
)

const MessageCollectionName = "message"

type MessageRepository struct {
	client *turmsmongo.Client
	col    *mongo.Collection
}

func NewMessageRepository(client *turmsmongo.Client) *MessageRepository {
	return &MessageRepository{
		client: client,
		col:    client.Collection(MessageCollectionName),
	}
}

// InsertMessage inserts a single Message PO into MongoDB.
// Wait, we don't handle _id (snowflake generation) here, it should be provided by the domain service.
func (r *MessageRepository) InsertMessage(ctx context.Context, msg *po.Message) error {
	// The deliveryDate (dyd) will be used by Mongo if collection is sharded.
	_, err := r.col.InsertOne(ctx, msg)
	return err
}

// FindMessagesByTarget retrieves messages using the primary multi-key index.
func (r *MessageRepository) FindMessagesByTarget(ctx context.Context, targetID int64, opts ...*options.FindOptions) ([]*po.Message, error) {
	// Simple lookup based on targetId (tid)
	filter := map[string]interface{}{
		"tid": targetID,
	}

	cursor, err := r.col.Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var msgs []*po.Message
	if err := cursor.All(ctx, &msgs); err != nil {
		return nil, err
	}
	return msgs, nil
}
