package idgen

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSnowflakeIdGenerator_NextIncreasingId(t *testing.T) {
	gen, err := NewSnowflakeIdGenerator(1, 1)
	require.NoError(t, err)

	id1 := gen.NextIncreasingId()
	id2 := gen.NextIncreasingId()

	assert.Greater(t, id2, id1, "IDs should strictly increase")
}

func TestSnowflakeIdGenerator_Concurrency(t *testing.T) {
	gen, err := NewSnowflakeIdGenerator(1, 1)
	require.NoError(t, err)

	idMap := sync.Map{}
	var wg sync.WaitGroup

	numGoroutines := 100
	numIdsPerGoroutine := 1000

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < numIdsPerGoroutine; j++ {
				id := gen.NextIncreasingId()
				if _, exists := idMap.LoadOrStore(id, true); exists {
					t.Errorf("Duplicate ID generated: %d", id)
				}
			}
		}()
	}

	wg.Wait()
}
