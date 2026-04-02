package po

import "time"

const CollectionNameGroupBlockedUser = "groupBlockedUser"

// GroupBlockedUserKey represents the composite primary key for GroupBlockedUser.
type GroupBlockedUserKey struct {
	GroupID int64 `bson:"gid"`
	UserID  int64 `bson:"uid"`
}

// GroupBlockedUser represents the group blocked user entity.
// Sharded by _id.gid.
type GroupBlockedUser struct {
	ID          GroupBlockedUserKey `bson:"_id"`
	BlockDate   *time.Time          `bson:"bd,omitempty"`
	RequesterID int64               `bson:"rid,omitempty"` // HASH index
}
