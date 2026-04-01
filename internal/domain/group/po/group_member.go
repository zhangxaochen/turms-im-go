package po

import "time"

const CollectionNameGroupMember = "groupMember"

// GroupMemberKey corresponds to the composite _id in MongoDB.
type GroupMemberKey struct {
	GroupID int64 `bson:"gid"`
	UserID  int64 `bson:"uid"`
}

// GroupMember represents the group member entity.
type GroupMember struct {
	ID          GroupMemberKey `bson:"_id"`
	Name        *string        `bson:"n,omitempty"`
	Role        int32          `bson:"role"`
	JoinDate    *time.Time     `bson:"jd,omitempty"`
	MuteEndDate *time.Time     `bson:"med,omitempty"`
}
