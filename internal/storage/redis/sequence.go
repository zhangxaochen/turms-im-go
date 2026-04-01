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
