package redis

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	testredis "github.com/testcontainers/testcontainers-go/modules/redis"
)

func setupRedisTestContainer(ctx context.Context, t *testing.T) (*Client, func()) {
	if testing.Short() {
		t.Skip("Skipping testcases requiring Redis in short mode")
	}

	redisContainer, err := testredis.Run(ctx,
		"redis:7.0",
	)
	require.NoError(t, err)

	uri, err := redisContainer.ConnectionString(ctx)
	require.NoError(t, err)

	client, err := NewClient(context.Background(), Config{
		URI: uri,
	})
	require.NoError(t, err)

	cleanup := func() {
		client.Close()
		if err := redisContainer.Terminate(context.Background()); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}

	return client, cleanup
}

func TestSequenceGenerator_NextPrivateMessageSequenceId(t *testing.T) {
	ctx := context.Background()
	client, cleanup := setupRedisTestContainer(ctx, t)
	defer cleanup()

	generator := NewSequenceGenerator(client)

	// Fetch 3 sequences serially
	const userID int64 = 1001
	id1, err := generator.NextPrivateMessageSequenceId(ctx, userID)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), id1)

	id2, err := generator.NextPrivateMessageSequenceId(ctx, userID)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), id2)

	id3, err := generator.NextPrivateMessageSequenceId(ctx, userID)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), id3)

	// Another user sequence should start from 1
	const anotherUserID int64 = 1002
	id4, err := generator.NextPrivateMessageSequenceId(ctx, anotherUserID)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), id4)
}

func TestSequenceGenerator_Concurrency(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, cleanup := setupRedisTestContainer(ctx, t)
	defer cleanup()

	generator := NewSequenceGenerator(client)
	const groupID int64 = 8888

	const numWorkers = 50
	const requestsPerWorker = 20

	var wg sync.WaitGroup
	errCh := make(chan error, numWorkers*requestsPerWorker)
	idCh := make(chan int64, numWorkers*requestsPerWorker)

	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < requestsPerWorker; j++ {
				id, err := generator.NextGroupMessageSequenceId(ctx, groupID)
				if err != nil {
					errCh <- err
				} else {
					idCh <- id
				}
			}
		}()
	}

	wg.Wait()
	close(errCh)
	close(idCh)

	for err := range errCh {
		assert.NoError(t, err)
	}

	// Check total ids generated
	totalIDs := numWorkers * requestsPerWorker
	assert.Equal(t, totalIDs, len(idCh))

	// Ensure all IDs are unique and between 1 and totalIDs
	seen := make(map[int64]bool)
	for id := range idCh {
		assert.False(t, seen[id], "Duplicate sequence ID detected: %d", id)
		seen[id] = true
		assert.GreaterOrEqual(t, id, int64(1))
		assert.LessOrEqual(t, id, int64(totalIDs))
	}
}
