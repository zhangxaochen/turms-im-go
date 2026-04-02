package po

import "time"

const CollectionNameUserVersion = "userVersion"

// UserVersion represents a user's local resource version tracking.
// MongoDB Collection: userVersion
type UserVersion struct {
	UserID                       int64     `bson:"_id"` // userId
	SentFriendRequests           time.Time `bson:"sfr"` // sentFriendRequests
	ReceivedFriendRequests       time.Time `bson:"rfr"` // receivedFriendRequests
	Relationships                time.Time `bson:"r"` // relationships
	RelationshipGroups           time.Time `bson:"rg"` // relationshipGroups
	RelationshipGroupMembers     time.Time `bson:"rgm"` // relationshipGroupMembers
	GroupJoinRequests            time.Time `bson:"gjr"` // groupJoinRequests
	SentGroupInvitations         time.Time `bson:"sgi"` // sentGroupInvitations
	ReceivedGroupInvitations     time.Time `bson:"rgi"` // receivedGroupInvitations
	JoinedGroups                 time.Time `bson:"jg"` // joinedGroups
}
