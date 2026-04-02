package po

import "time"

const CollectionNameGroupInvitation = "groupInvitation"

// GroupInvitation represents the group invitation entity.
// Sharded by inviteeId.
type GroupInvitation struct {
	ID           int64         `bson:"_id"`
	GroupID      int64         `bson:"gid,omitempty"`  // HASH index
	InviterID    int64         `bson:"irid,omitempty"` // HASH index
	InviteeID    int64         `bson:"ieid"`
	Content      string        `bson:"cnt"`
	Status       RequestStatus `bson:"stat"`
	CreationDate time.Time     `bson:"cd"`
	ResponseDate *time.Time    `bson:"rd,omitempty"`
	Reason       *string       `bson:"rsn,omitempty"`
}
