package po

import "time"

// Message represents a message entity in MongoDB (message collection).
// Shard key: deliveryDate (dyd)
type Message struct {
	ID               int64      `bson:"_id"`
	ConversationID   []byte     `bson:"cid,omitempty"`
	IsGroupMessage   *bool      `bson:"gm,omitempty"`
	IsSystemMessage  *bool      `bson:"sm,omitempty"`
	DeliveryDate     time.Time  `bson:"dyd"`
	ModificationDate time.Time  `bson:"md,omitempty"`
	DeletionDate     *time.Time `bson:"dd,omitempty"`
	RecallDate       *time.Time `bson:"rd,omitempty"`
	Text             string     `bson:"txt,omitempty"`
	SenderID         int64      `bson:"sid"`
	SenderIP         *int32     `bson:"sip,omitempty"`
	SenderIPv6       []byte     `bson:"sip6,omitempty"`
	TargetID         int64      `bson:"tid"`
	Records          [][]byte   `bson:"rec,omitempty"`
	BurnAfter        *int32     `bson:"bf,omitempty"`
	ReferenceID      *int64     `bson:"rid,omitempty"`
	SequenceID       *int32     `bson:"sqid,omitempty"`
	PreMessageID     *int64     `bson:"pmid,omitempty"`

	UserDefinedAttributes map[string]any `bson:"user_defined_attributes,omitempty"`
}
