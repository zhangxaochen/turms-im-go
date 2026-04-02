package po

import "time"

const CollectionNameGroupConversation = "groupConversation"

// GroupConversation represents the conversation for a group, tracking read dates per member.
type GroupConversation struct {
	// ID is the group ID.
	ID int64 `bson:"_id"`

	// MemberIDToReadDate tracks the read date for each member.
	// BSON key "mr" stands for "memberIdToReadDate".
	// The map keys are string representations of member int64 IDs.
	MemberIDToReadDate map[string]time.Time `bson:"mr,omitempty"`
}
