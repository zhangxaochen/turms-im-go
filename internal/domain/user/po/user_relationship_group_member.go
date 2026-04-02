package po

import "time"

const CollectionNameUserRelationshipGroupMember = "userRelationshipGroupMember"

// UserRelationshipGroupMember represents a member in a relationship group.
// MongoDB Collection: userRelationshipGroupMember
type UserRelationshipGroupMember struct {
	Key      UserRelationshipGroupMemberKey `bson:"_id"` // ownerId, groupIndex, and relatedUserId
	JoinDate time.Time                      `bson:"jd"`
}

type UserRelationshipGroupMemberKey struct {
	OwnerID       int64 `bson:"oid"`
	GroupIndex    int32 `bson:"gi"`
	RelatedUserID int64 `bson:"rid"`
}
