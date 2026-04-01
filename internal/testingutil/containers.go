package testingutil

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	testredis "github.com/testcontainers/testcontainers-go/modules/redis"
	turmsmongo "im.turms/server/internal/storage/mongo"
	turmsredis "im.turms/server/internal/storage/redis"
)

// SetupMongo starts a MongoDB container for BDD testing and returns a client and a cleanup function.
// It uses t.Skip() if the container fails to start, allowing graceful degradation in environments without Docker.
func SetupMongo(t *testing.T, dbName string) (*turmsmongo.Client, func()) {
	if testing.Short() {
		t.Skip("Skipping testcases requiring Mongo in short mode")
	}

	ctx := context.Background()
	mongodbContainer, err := mongodb.Run(ctx, "mongo:7.0")
	if err != nil {
		t.Skipf("Skipping BDD test, testcontainers failed to start MongoDB: %v", err)
	}

	uri, err := mongodbContainer.ConnectionString(ctx)
	require.NoError(t, err)

	cfg := turmsmongo.Config{
		URI:            uri,
		Database:       dbName,
		ConnectTimeout: 10 * time.Second,
	}

	client, err := turmsmongo.NewClient(ctx, cfg)
	require.NoError(t, err)

	cleanup := func() {
		client.Close(ctx)
		_ = mongodbContainer.Terminate(context.Background())
	}

	return client, cleanup
}

// SetupRedis starts a Redis container for BDD testing and returns a client and a cleanup function.
// It uses t.Skip() if the container fails to start.
func SetupRedis(t *testing.T) (*turmsredis.Client, func()) {
	if testing.Short() {
		t.Skip("Skipping testcases requiring Redis in short mode")
	}

	ctx := context.Background()
	redisContainer, err := testredis.Run(ctx, "redis:7.0")
	if err != nil {
		t.Skipf("Skipping BDD test, testcontainers failed to start Redis: %v", err)
	}

	uri, err := redisContainer.ConnectionString(ctx)
	require.NoError(t, err)

	cfg := turmsredis.Config{
		URI: uri,
	}

	client, err := turmsredis.NewClient(ctx, cfg)
	require.NoError(t, err)

	cleanup := func() {
		client.Close()
		_ = redisContainer.Terminate(context.Background())
	}

	return client, cleanup
}
