package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	URI            string
	ConnectTimeout time.Duration
}

type Client struct {
	RDB redis.UniversalClient
}

// NewClient parses a standard Redis URI (or redis+sentinel) and connects.
// For full turms-cluster support, this might later be expanded to ParseClusterURL
func NewClient(ctx context.Context, cfg Config) (*Client, error) {
	if cfg.ConnectTimeout == 0 {
		cfg.ConnectTimeout = 10 * time.Second
	}

	opt, err := redis.ParseURL(cfg.URI)
	if err != nil {
		return nil, err
	}

	opt.DialTimeout = cfg.ConnectTimeout
	rdb := redis.NewClient(opt)

	pingCtx, cancel := context.WithTimeout(ctx, cfg.ConnectTimeout)
	defer cancel()

	if err := rdb.Ping(pingCtx).Err(); err != nil {
		return nil, err
	}

	return &Client{
		RDB: rdb,
	}, nil
}

func (c *Client) Close() error {
	return c.RDB.Close()
}
