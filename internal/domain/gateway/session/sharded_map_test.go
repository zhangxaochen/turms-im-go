package session

import (
	"sync"
	"testing"

	"im.turms/server/pkg/protocol"
	"github.com/stretchr/testify/assert"
)

func TestShardedUserSessionsMap_GetOrAdd(t *testing.T) {
	smap := NewShardedUserSessionsMap(16)
	userID := int64(1001)

	manager := smap.GetOrAdd(userID)
	assert.NotNil(t, manager)
	assert.Equal(t, userID, manager.UserID)

	// Add again, should return the exact same pointer
	manager2 := smap.GetOrAdd(userID)
	assert.Same(t, manager, manager2)
}

func TestShardedUserSessionsMap_CountOnlineUsers(t *testing.T) {
	smap := NewShardedUserSessionsMap(16)
	
	smap.GetOrAdd(1)
	smap.GetOrAdd(2)
	smap.GetOrAdd(3)
	
	count := smap.CountOnlineUsers()
	assert.Equal(t, 3, count)
}

func TestShardedUserSessionsMap_RemoveIfEmpty(t *testing.T) {
	smap := NewShardedUserSessionsMap(16)
	userID := int64(1001)

	manager := smap.GetOrAdd(userID)
	session := &UserSession{
		UserID:     userID,
		DeviceType: protocol.DeviceType_DESKTOP,
	}
	manager.AddSession(session)

	// Not empty, should not be removed
	smap.RemoveIfEmpty(userID)
	_, ok := smap.Get(userID)
	assert.True(t, ok)

	// Remove the session, then RemoveIfEmpty
	manager.RemoveSession(protocol.DeviceType_DESKTOP)
	smap.RemoveIfEmpty(userID)
	
	_, ok = smap.Get(userID)
	assert.False(t, ok)
}

func TestShardedUserSessionsMap_Concurrency(t *testing.T) {
	smap := NewShardedUserSessionsMap(256)
	var wg sync.WaitGroup

	numUsers := int64(1000)
	workers := 10

	// 10 workers adding 1000 users concurrently
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := int64(0); i < numUsers; i++ {
				smap.GetOrAdd(i)
			}
		}()
	}
	wg.Wait()

	assert.Equal(t, int(numUsers), smap.CountOnlineUsers())
}
