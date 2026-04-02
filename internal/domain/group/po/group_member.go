package po

import (
	"time"

	"im.turms/server/pkg/protocol"
)

// GroupMemberKey is the composite primary key for a GroupMember
type GroupMemberKey struct {
	GroupID int64 `bson:"gid"`
	UserID  int64 `bson:"uid"`
}

// GroupMember represents a group member entity in MongoDB (groupMember collection).
// Shard key: _id.groupId (HASH)
type GroupMember struct {
	ID          GroupMemberKey           `bson:"_id"`
	Name        *string                  `bson:"n,omitempty"`
	Role        protocol.GroupMemberRole `bson:"role"`
	JoinDate    *time.Time               `bson:"jd,omitempty"`
	MuteEndDate *time.Time               `bson:"med,omitempty"`
}
