package po

import "time"

// Group represents a group entity in MongoDB (group collection).
// Shard key: none strictly defined for single node, potentially hash on _id
type Group struct {
	ID               int64      `bson:"_id"`
	TypeID           *int64     `bson:"tid,omitempty"`
	CreatorID        *int64     `bson:"cid,omitempty"`
	OwnerID          *int64     `bson:"oid,omitempty"`
	Name             *string    `bson:"n,omitempty"`
	Intro            *string    `bson:"intro,omitempty"`
	Announcement     *string    `bson:"annt,omitempty"`
	MinimumScore     *int32     `bson:"ms,omitempty"`
	CreationDate     *time.Time `bson:"cd,omitempty"`
	DeletionDate     *time.Time `bson:"dd,omitempty"`
	LastUpdatedDate  *time.Time `bson:"lud,omitempty"`
	MuteEndDate      *time.Time `bson:"med,omitempty"`
	IsActive         *bool      `bson:"ac,omitempty"`

	UserDefinedAttributes map[string]any `bson:"uda,omitempty"`
}
