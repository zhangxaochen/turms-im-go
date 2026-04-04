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
		filter["gm"] = *isGroupMessage
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

// CountMessages counts messages matching the specific criteria.
func (r *MessageRepository) CountMessages(
	ctx context.Context,
	isGroupMessage *bool,
	senderIDs []int64,
	targetIDs []int64,
	deliveryDateAfter *time.Time,
	deliveryDateBefore *time.Time,
) (int64, error) {
	filter := bson.M{}

	if isGroupMessage != nil {
		filter["gm"] = *isGroupMessage
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

	return r.col.CountDocuments(ctx, filter)
}

// FindMessages finds messages matching various complex criteria.
// @MappedFrom findMessages(@Nullable Collection<Long> messageIds, @Nullable Collection<byte[]> conversationIds, @Nullable Boolean areGroupMessages, @Nullable Boolean areSystemMessages, @Nullable Set<Long> senderIds, @Nullable Set<Long> targetIds, @Nullable DateRange deliveryDateRange, @Nullable DateRange deletionDateRange, @Nullable DateRange recallDateRange, @Nullable Integer page, @Nullable Integer size, @Nullable Boolean ascending)
func (r *MessageRepository) FindMessages(
	ctx context.Context,
	messageIDs []int64,
	conversationIDs [][]byte,
	areGroupMessages *bool,
	areSystemMessages *bool,
	senderIDs []int64,
	targetIDs []int64,
	deliveryDateRange *turmsmongo.DateRange,
	deletionDateRange *turmsmongo.DateRange,
	recallDateRange *turmsmongo.DateRange,
	page *int,
	size *int,
	ascending *bool,
) ([]*po.Message, error) {
	filter := bson.M{}

	if len(messageIDs) > 0 {
		filter["_id"] = bson.M{"$in": messageIDs}
	}
	// Note: conversationIds mapping might require custom logic, but usually it maps to `cid`
	if len(conversationIDs) > 0 {
		filter["cid"] = bson.M{"$in": conversationIDs}
	}
	if areGroupMessages != nil {
		filter["gm"] = *areGroupMessages
	}
	if areSystemMessages != nil {
		filter["sm"] = *areSystemMessages
	}
	if len(senderIDs) > 0 {
		filter["sid"] = bson.M{"$in": senderIDs}
	}
	if len(targetIDs) > 0 {
		filter["tid"] = bson.M{"$in": targetIDs}
	}

	if deliveryDateRange != nil {
		dateFilter := bson.M{}
		if deliveryDateRange.Start != nil {
			dateFilter["$gt"] = *deliveryDateRange.Start
		}
		if deliveryDateRange.End != nil {
			dateFilter["$lt"] = *deliveryDateRange.End
		}
		if len(dateFilter) > 0 {
			filter["dyd"] = dateFilter
		}
	}
	if deletionDateRange != nil {
		dateFilter := bson.M{}
		if deletionDateRange.Start != nil {
			dateFilter["$gt"] = *deletionDateRange.Start
		}
		if deletionDateRange.End != nil {
			dateFilter["$lt"] = *deletionDateRange.End
		}
		if len(dateFilter) > 0 {
			filter["dd"] = dateFilter
		}
	}
	if recallDateRange != nil {
		dateFilter := bson.M{}
		if recallDateRange.Start != nil {
			dateFilter["$gt"] = *recallDateRange.Start
		}
		if recallDateRange.End != nil {
			dateFilter["$lt"] = *recallDateRange.End
		}
		if len(dateFilter) > 0 {
			filter["rd"] = dateFilter
		}
	}

	opts := options.Find()

	if page != nil && size != nil {
		skip := int64(*page * *size)
		opts.SetSkip(skip)
	}
	if size != nil {
		opts.SetLimit(int64(*size))
	}
	if ascending != nil {
		if *ascending {
			opts.SetSort(bson.D{{Key: "dyd", Value: 1}})
		} else {
			opts.SetSort(bson.D{{Key: "dyd", Value: -1}})
		}
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

// UpdateMessages updates multiple messages.
func (r *MessageRepository) UpdateMessages(
	ctx context.Context,
	messageIDs []int64,
	isSystemMessage *bool,
	senderIP *int32,
	senderIPv6 []byte,
	recallDate *time.Time,
	text *string,
	records [][]byte,
	burnAfter *int32,
) error {
	set := bson.M{}

	if isSystemMessage != nil {
		set["sm"] = *isSystemMessage
	}
	if senderIP != nil {
		set["sip"] = *senderIP
	}
	if len(senderIPv6) > 0 {
		set["sip6"] = senderIPv6
	}
	if recallDate != nil {
		set["rd"] = *recallDate
	}
	if text != nil {
		set["txt"] = *text
	}
	if records != nil {
		set["rec"] = records
	}
	if burnAfter != nil {
		set["bf"] = *burnAfter
	}

	if len(set) == 0 {
		return nil
	}

	filter := bson.M{}
	if len(messageIDs) > 0 {
		filter["_id"] = bson.M{"$in": messageIDs}
	} else {
		// If messageIDs is empty, do nothing or return nil
		return nil
	}

	update := bson.M{"$set": set}
	_, err := r.col.UpdateMany(ctx, filter, update)
	return err
}

// UpdateMessagesDeletionDate sets the deletion date for given message IDs.
func (r *MessageRepository) UpdateMessagesDeletionDate(ctx context.Context, messageIDs []int64, deletionDate *time.Time) error {
	if len(messageIDs) == 0 {
		return nil
	}
	filter := bson.M{"_id": bson.M{"$in": messageIDs}}
	set := bson.M{}
	if deletionDate == nil {
		set["dd"] = nil
	} else {
		set["dd"] = *deletionDate
	}
	update := bson.M{"$set": set}
	_, err := r.col.UpdateMany(ctx, filter, update)
	return err
}

// ExistsBySenderIdAndTargetId checks if any message exists with the given sender and target ID.
func (r *MessageRepository) ExistsBySenderIDAndTargetID(ctx context.Context, senderID int64, targetID int64) (bool, error) {
	filter := bson.M{"sid": senderID, "tid": targetID}
	err := r.col.FindOne(ctx, filter, options.FindOne().SetProjection(bson.M{"_id": 1})).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// FindDeliveryDate finds the delivery date of a message.
func (r *MessageRepository) FindDeliveryDate(ctx context.Context, messageID int64) (*time.Time, error) {
	filter := bson.M{"_id": messageID}
	var result struct {
		DeliveryDate time.Time `bson:"dyd"`
	}
	err := r.col.FindOne(ctx, filter, options.FindOne().SetProjection(bson.M{"dyd": 1})).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result.DeliveryDate, nil
}

// FindExpiredMessageIds finds IDs of messages delivered before an expiration date.
func (r *MessageRepository) FindExpiredMessageIds(ctx context.Context, expirationDate time.Time) ([]int64, error) {
	filter := bson.M{"dyd": bson.M{"$lt": expirationDate}}
	cursor, err := r.col.Find(ctx, filter, options.Find().SetProjection(bson.M{"_id": 1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var results []struct {
		ID int64 `bson:"_id"`
	}
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	var ids []int64
	for _, res := range results {
		ids = append(ids, res.ID)
	}
	return ids, nil
}

// FindMessageGroupId finds the group ID indicating the target of a group message.
func (r *MessageRepository) FindMessageGroupId(ctx context.Context, messageID int64) (*int64, error) {
	filter := bson.M{"_id": messageID}
	var result struct {
		TargetID int64 `bson:"tid"`
	}
	err := r.col.FindOne(ctx, filter, options.FindOne().SetProjection(bson.M{"tid": 1})).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result.TargetID, nil
}

// FindMessageSenderIDAndTargetIDAndIsGroupMessage returns senderId, targetId, and isGroupMessage of a message.
func (r *MessageRepository) FindMessageSenderIDAndTargetIDAndIsGroupMessage(ctx context.Context, messageID int64) (*po.Message, error) {
	filter := bson.M{"_id": messageID}
	var result po.Message
	err := r.col.FindOne(ctx, filter, options.FindOne().SetProjection(bson.M{"sid": 1, "tid": 1, "gm": 1})).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// FindIsGroupMessageAndTargetID returns isGroupMessage and targetId of a message.
func (r *MessageRepository) FindIsGroupMessageAndTargetID(ctx context.Context, messageID int64, senderID int64) (*po.Message, error) {
	filter := bson.M{"_id": messageID, "sid": senderID}
	var result po.Message
	err := r.col.FindOne(ctx, filter, options.FindOne().SetProjection(bson.M{"gm": 1, "tid": 1})).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// FindIsGroupMessageAndTargetIDAndDeliveryDate returns gm, tid, dyd of a message.
func (r *MessageRepository) FindIsGroupMessageAndTargetIDAndDeliveryDate(ctx context.Context, messageID int64, senderID int64) (*po.Message, error) {
	filter := bson.M{"_id": messageID, "sid": senderID}
	var result po.Message
	err := r.col.FindOne(ctx, filter, options.FindOne().SetProjection(bson.M{"gm": 1, "tid": 1, "dyd": 1})).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetGroupConversationID computes the conversation ID for a group message.
func (r *MessageRepository) GetGroupConversationID(groupID int64) []byte {
	b := make([]byte, 8)
	for i := 7; i >= 0; i-- {
		b[i] = byte(groupID)
		groupID >>= 8
	}
	return b
}

// GetPrivateConversationID computes the conversation ID for a private message.
func (r *MessageRepository) GetPrivateConversationID(id1 int64, id2 int64) []byte {
	minID := id1
	maxID := id2
	if id1 > id2 {
		minID = id2
		maxID = id1
	}
	b := make([]byte, 16)
	for i := 7; i >= 0; i-- {
		b[i] = byte(minID)
		minID >>= 8
	}
	for i := 15; i >= 8; i-- {
		b[i] = byte(maxID)
		maxID >>= 8
	}
	return b
}

// DeleteMessages physically deletes messages by their IDs.
func (r *MessageRepository) DeleteMessages(ctx context.Context, messageIDs []int64) error {
	if len(messageIDs) == 0 {
		return nil
	}
	filter := bson.M{"_id": bson.M{"$in": messageIDs}}
	_, err := r.col.DeleteMany(ctx, filter)
	return err
}

// CountUsersWhoSentMessage counts distinct users who sent messages.
func (r *MessageRepository) CountUsersWhoSentMessage(ctx context.Context, deliveryDateAfter *time.Time, deliveryDateBefore *time.Time, areGroupMessages *bool, areSystemMessages *bool) (int64, error) {
	filter := bson.M{}
	if areGroupMessages != nil {
		filter["gm"] = *areGroupMessages
	}
	if areSystemMessages != nil {
		filter["sm"] = *areSystemMessages
	}
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

	uniqueSenderIDs, err := r.col.Distinct(ctx, "sid", filter)
	if err != nil {
		return 0, err
	}
	return int64(len(uniqueSenderIDs)), nil
}

// CountGroupsThatSentMessages counts distinct groups that had messages sent to them.
func (r *MessageRepository) CountGroupsThatSentMessages(ctx context.Context, deliveryDateAfter *time.Time, deliveryDateBefore *time.Time) (int64, error) {
	filter := bson.M{"gm": true}
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

	uniqueTargetIDs, err := r.col.Distinct(ctx, "tid", filter)
	if err != nil {
		return 0, err
	}
	return int64(len(uniqueTargetIDs)), nil
}

// CountSentMessages counts the number of sent messages based on criteria.
func (r *MessageRepository) CountSentMessages(ctx context.Context, deliveryDateAfter *time.Time, deliveryDateBefore *time.Time, areGroupMessages *bool, areSystemMessages *bool) (int64, error) {
	filter := bson.M{}
	if areGroupMessages != nil {
		filter["gm"] = *areGroupMessages
	}
	if areSystemMessages != nil {
		filter["sm"] = *areSystemMessages
	}
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
	return r.col.CountDocuments(ctx, filter)
}
