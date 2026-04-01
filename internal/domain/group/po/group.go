package po

import "time"

// Group represents a group entity in MongoDB (group collection).
// Hash Sharding by _id
type Group struct {
	ID           int64      `bson:"_id"`
	TypeID       int64      `bson:"tid"`
	CreatorID    int64      `bson:"cid"`
	OwnerID      int64      `bson:"oid"`
	Name         string     `bson:"n,omitempty"`
	Intro        string     `bson:"intro,omitempty"`
	Announcement string     `bson:"ann,omitempty"`
	MinimumScore int32      `bson:"ms"`
	CreationDate time.Time  `bson:"cd"`
	DeletionDate *time.Time `bson:"dd,omitempty"`
	MuteEndDate  *time.Time `bson:"med,omitempty"`
	IsActive     bool       `bson:"act,omitempty"`

	UserDefinedAttributes map[string]any `bson:"user_defined_attributes,omitempty"`
}
