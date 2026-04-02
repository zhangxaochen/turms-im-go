package po

import "time"

type GroupMemberRole int32

const (
	GroupMemberRole_OWNER           GroupMemberRole = 0
	GroupMemberRole_MANAGER         GroupMemberRole = 1
	GroupMemberRole_MEMBER          GroupMemberRole = 2
	GroupMemberRole_GUEST           GroupMemberRole = 3
	GroupMemberRole_ANONYMOUS_GUEST GroupMemberRole = 4
)

// GroupMemberKey is the composite primary key for a GroupMember
type GroupMemberKey struct {
	GroupID int64 `bson:"gid"`
	UserID  int64 `bson:"uid"`
}

// GroupMember represents a group member entity in MongoDB (groupMember collection).
// Shard key: _id.groupId (HASH)
type GroupMember struct {
	ID          GroupMemberKey  `bson:"_id"`
	Name        *string         `bson:"n,omitempty"`
	Role        GroupMemberRole `bson:"role"`
	JoinDate    *time.Time      `bson:"jd,omitempty"`
	MuteEndDate *time.Time      `bson:"med,omitempty"`
}
