package po

import "time"

const CollectionNameUserRelationshipGroup = "userRelationshipGroup"

// UserRelationshipGroup represents a relationship group.
// MongoDB Collection: userRelationshipGroup
type UserRelationshipGroup struct {
	Key          UserRelationshipGroupKey `bson:"_id"` // ownerId and index
	Name         string                   `bson:"n"`
	CreationDate time.Time                `bson:"cd"`
}

type UserRelationshipGroupKey struct {
	OwnerID int64 `bson:"oid"`
	Index   int32 `bson:"i"` // Use i as the group index
}
