package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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
func (r *MessageRepository) InsertMessage(ctx context.Context, msg *po.Message) error {
	// The deliveryDate (dyd) will be used by Mongo if collection is sharded.
	_, err := r.col.InsertOne(ctx, msg)
	return err
}

// FindByID retrieves a message by its ID.
func (r *MessageRepository) FindByID(ctx context.Context, id int64) (*po.Message, error) {
	filter := map[string]interface{}{
		"_id": id,
	}
	var msg po.Message
	if err := r.col.FindOne(ctx, filter).Decode(&msg); err != nil {
		return nil, err
	}
	return &msg, nil
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

// QueryMessages supports complex querying for message pulling (offline/roaming sync).
// @MappedFrom queryMessages(@Nullable Collection<Long> messageIds, @Nullable Boolean areGroupMessages, @Nullable Boolean areSystemMessages, @Nullable Set<Long> senderIds, @Nullable Set<Long> targetIds, @Nullable DateRange deliveryDateRange, @Nullable DateRange deletionDateRange, @Nullable DateRange recallDateRange, @Nullable Integer page, @Nullable Integer size, @Nullable Boolean ascending)
// @MappedFrom queryMessages(@QueryParam(required = false)
func (r *MessageRepository) QueryMessages(
	ctx context.Context,
	isGroupMessage *bool,
	senderIDs []int64,
	targetIDs []int64,
	deliveryDateAfter *time.Time,
	deliveryDateBefore *time.Time,
	size int64,
	ascending bool,
) ([]*po.Message, error) {
	filter := bson.M{}

	if isGroupMessage != nil {
		filter["igm"] = *isGroupMessage
	}

	if len(senderIDs) > 0 {
		filter["sid"] = bson.M{"$in": senderIDs}
	}

	if len(targetIDs) > 0 {
		filter["tid"] = bson.M{"$in": targetIDs}
	}

	// Date Range
	dateFilter := bson.M{}
	if deliveryDateAfter != nil {
		dateFilter["$gt"] = *deliveryDateAfter
	}
	if deliveryDateBefore != nil {
		dateFilter["$lt"] = *deliveryDateBefore
	}
	if len(dateFilter) > 0 {
		filter["dyd"] = dateFilter
	}

	opts := options.Find().SetLimit(size)
	if ascending {
		opts.SetSort(bson.D{{Key: "dyd", Value: 1}})
	} else {
		opts.SetSort(bson.D{{Key: "dyd", Value: -1}})
	}

	cursor, err := r.col.Find(ctx, filter, opts)
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

// UpdateMessage partially updates a message's text, modificationDate, or recallDate.
func (r *MessageRepository) UpdateMessage(
	ctx context.Context,
	messageID int64,
	text *string,
	modificationDate *time.Time,
	recallDate *time.Time,
) error {
	set := bson.M{}

	if text != nil {
		set["txt"] = *text
	}
	if modificationDate != nil {
		set["md"] = *modificationDate
	}
	if recallDate != nil {
		set["rd"] = *recallDate
	}

	if len(set) == 0 {
		return nil // nothing to update
	}

	filter := bson.M{"_id": messageID}
	update := bson.M{"$set": set}

	_, err := r.col.UpdateOne(ctx, filter, update)
	return err
}
