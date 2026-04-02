package po

const CollectionNameGroupJoinQuestion = "groupJoinQuestion"

// GroupJoinQuestion represents the group join question entity.
// Sharded by groupId (HASH).
type GroupJoinQuestion struct {
	ID       int64    `bson:"_id"`
	GroupID  int64    `bson:"gid"` // HASH index
	Question string   `bson:"q"`
	Answers  []string `bson:"ans"` // Set of strings
	Score    int      `bson:"score"`
}
