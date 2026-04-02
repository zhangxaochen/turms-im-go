package po

import "time"

const CollectionNamePrivateConversation = "privateConversation"

// PrivateConversationKey corresponds to the _id field in MongoDB.
// MongoDB BSON marshaling order matches the struct definition order.
type PrivateConversationKey struct {
	OwnerID  int64 `bson:"oid"`
	TargetID int64 `bson:"tid"`
}

// PrivateConversation represents a private conversation between two users.
type PrivateConversation struct {
	ID       PrivateConversationKey `bson:"_id"`
	ReadDate time.Time              `bson:"rd"` // rd: read date
}
