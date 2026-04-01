package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	URI            string
	Database       string
	ConnectTimeout time.Duration
}

type Client struct {
	DB     *mongo.Database
	Client *mongo.Client
}

func NewClient(ctx context.Context, cfg Config) (*Client, error) {
	if cfg.ConnectTimeout == 0 {
		cfg.ConnectTimeout = 10 * time.Second
	}

	clientOptions := options.Client().ApplyURI(cfg.URI)

	// Context for connect and ping
	connectCtx, cancel := context.WithTimeout(ctx, cfg.ConnectTimeout)
	defer cancel()

	client, err := mongo.Connect(connectCtx, clientOptions)
	if err != nil {
		return nil, err
	}

	if err := client.Ping(connectCtx, nil); err != nil {
		return nil, err
	}

	return &Client{
		DB:     client.Database(cfg.Database),
		Client: client,
	}, nil
}

func (c *Client) Collection(name string) *mongo.Collection {
	return c.DB.Collection(name)
}

func (c *Client) Close(ctx context.Context) error {
	return c.Client.Disconnect(ctx)
}
