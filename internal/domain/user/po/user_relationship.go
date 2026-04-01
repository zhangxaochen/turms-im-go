package po

import "time"

const CollectionNameUserRelationship = "userRelationship"

// UserRelationshipKey corresponds to the composite _id in MongoDB.
type UserRelationshipKey struct {
	OwnerID       int64 `bson:"oid"`
	RelatedUserID int64 `bson:"rid"`
}

// UserRelationship represents the user relationship entity.
type UserRelationship struct {
	ID                UserRelationshipKey `bson:"_id"`
	Name              *string             `bson:"n,omitempty"`
	BlockDate         *time.Time          `bson:"bd,omitempty"`
	GroupIndex        *int32              `bson:"gi,omitempty"`
	EstablishmentDate *time.Time          `bson:"ed,omitempty"`
}
