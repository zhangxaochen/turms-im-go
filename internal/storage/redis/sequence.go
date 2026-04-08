package redis

import (
	"context"
	"fmt"
)

const (
	// Prefix format for message sequences: "ps:{uid}" or "gs:{group_id}"
	PrivateMessageSequencePrefix = "ps"
	GroupMessageSequencePrefix   = "gs"
)

// SequenceGenerator encapsulates Redis sequence generation logic.
type SequenceGenerator struct {
	client *Client
}

// NewSequenceGenerator returns a distributed sequence generator.
func NewSequenceGenerator(client *Client) *SequenceGenerator {
	return &SequenceGenerator{
		client: client,
	}
}

// NextPrivateMessageSequenceId fetches the next strictly increasing Sequence ID for a specific user.
func (g *SequenceGenerator) NextPrivateMessageSequenceId(ctx context.Context, userID int64) (int64, error) {
	key := fmt.Sprintf("%s:%d", PrivateMessageSequencePrefix, userID)
	return g.client.RDB.Incr(ctx, key).Result()
}

// NextGroupMessageSequenceId fetches the next strictly increasing Sequence ID for a specific group.
func (g *SequenceGenerator) NextGroupMessageSequenceId(ctx context.Context, groupID int64) (int64, error) {
	key := fmt.Sprintf("%s:%d", GroupMessageSequencePrefix, groupID)
	return g.client.RDB.Incr(ctx, key).Result()
}

// DeleteGroupMessageSequenceIDs deletes the group message sequence ID keys for the given group IDs.
// Java equivalent: redisClientManager.execute reactiveHdel with key pattern "gs:{groupId}"
func (g *SequenceGenerator) DeleteGroupMessageSequenceIDs(ctx context.Context, groupIDs []int64) error {
	if g == nil || g.client == nil {
		return nil
	}
	if len(groupIDs) == 0 {
		return nil
	}
	keys := make([]string, len(groupIDs))
	for i, id := range groupIDs {
		keys[i] = fmt.Sprintf("%s:%d", GroupMessageSequencePrefix, id)
	}
	return g.client.RDB.Del(ctx, keys...).Err()
}

// DeletePrivateMessageSequenceIDs deletes the private message sequence ID keys for the given user IDs.
// Java equivalent: uses a Lua script (deletePrivateMessageSequenceIdScript) to atomically delete
// private message sequence IDs. Each user pair has a key like "ps:{smallerID}:{largerID}" or "ps:{uid}".
func (g *SequenceGenerator) DeletePrivateMessageSequenceIDs(ctx context.Context, userIDs []int64) error {
	if g == nil || g.client == nil {
		return nil
	}
	if len(userIDs) == 0 {
		return nil
	}
	keys := make([]string, len(userIDs))
	for i, id := range userIDs {
		keys[i] = fmt.Sprintf("%s:%d", PrivateMessageSequencePrefix, id)
	}
	return g.client.RDB.Del(ctx, keys...).Err()
}

// FetchPrivateMessageSequenceIDForPair generates the private message sequence ID using both user IDs.
// Java sorts the two user IDs and uses both as keys in the Lua script. We compute the key using
// the smaller ID first to ensure consistency.
func (g *SequenceGenerator) FetchPrivateMessageSequenceIDForPair(ctx context.Context, userID1 int64, userID2 int64) (int64, error) {
	if g == nil || g.client == nil {
		return 0, nil
	}
	// Sort user IDs to produce a consistent key regardless of parameter order
	smaller, larger := userID1, userID2
	if userID1 > userID2 {
		smaller, larger = userID2, userID1
	}
	key := fmt.Sprintf("%s:%d:%d", PrivateMessageSequencePrefix, smaller, larger)
	return g.client.RDB.Incr(ctx, key).Result()
}
