package po

import "time"

// GroupVersion represents a group version entity in MongoDB (groupVersion collection).
// Sharded by _id (GroupID)
type GroupVersion struct {
	GroupID       int64      `bson:"_id"`
	Members       *time.Time `bson:"mbr,omitempty"`
	Blocklist     *time.Time `bson:"bl,omitempty"`
	JoinRequests  *time.Time `bson:"jr,omitempty"`
	JoinQuestions *time.Time `bson:"jq,omitempty"`
	Invitations   *time.Time `bson:"invt,omitempty"`
}
