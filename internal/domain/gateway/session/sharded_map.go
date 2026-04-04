package session

import (
	"sync"

	"im.turms/server/pkg/protocol"
)

// UserSessionsManager manages all sessions (devices) for a specific user
type UserSessionsManager struct {
	UserID   int64
	Sessions map[protocol.DeviceType]*UserSession
	// Could hold UserStatus etc.
	mu sync.RWMutex
}

func NewUserSessionsManager(userID int64) *UserSessionsManager {
	return &UserSessionsManager{
		UserID:   userID,
		Sessions: make(map[protocol.DeviceType]*UserSession),
	}
}

func (m *UserSessionsManager) AddSession(session *UserSession) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Sessions[session.DeviceType] = session
}

// @MappedFrom getSession(@NotNull DeviceType deviceType)
func (m *UserSessionsManager) GetSession(deviceType protocol.DeviceType) *UserSession {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.Sessions[deviceType]
}

func (m *UserSessionsManager) RemoveSession(deviceType protocol.DeviceType) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.Sessions, deviceType)
}

// @MappedFrom isEmpty()
func (m *UserSessionsManager) IsEmpty() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.Sessions) == 0
}

func (m *UserSessionsManager) GetAllSessions() []*UserSession {
	m.mu.RLock()
	defer m.mu.RUnlock()
	sessions := make([]*UserSession, 0, len(m.Sessions))
	for _, s := range m.Sessions {
		sessions = append(sessions, s)
	}
	return sessions
}

// ShardedUserSessionsMap is a concurrent, sharded map to reduce lock contention
type ShardedUserSessionsMap struct {
	shards []*shard
	num    uint64
}

type shard struct {
	sync.RWMutex
	m map[int64]*UserSessionsManager
}

func NewShardedUserSessionsMap(numShards int) *ShardedUserSessionsMap {
	// Must be power of 2 for fast modulo
	if numShards == 0 || (numShards&(numShards-1)) != 0 {
		numShards = 256 // default
	}
	shards := make([]*shard, numShards)
	for i := 0; i < numShards; i++ {
		shards[i] = &shard{
			m: make(map[int64]*UserSessionsManager),
		}
	}
	return &ShardedUserSessionsMap{
		shards: shards,
		num:    uint64(numShards),
	}
}

// fnv1a to string/int hashing is usually good, for int64 we can just use bit mixing
func (m *ShardedUserSessionsMap) getShard(key int64) *shard {
	// simple Wang/Jenkins hash for fast mixing 64 -> index
	key = (^key) + (key << 18)
	key = key ^ (key >> 31)
	key = key * 21
	key = key ^ (key >> 11)
	key = key + (key << 6)
	key = key ^ (key >> 22)
	idx := uint64(key) & (m.num - 1)
	return m.shards[idx]
}

func (m *ShardedUserSessionsMap) Get(userID int64) (*UserSessionsManager, bool) {
	shard := m.getShard(userID)
	shard.RLock()
	defer shard.RUnlock()
	manager, ok := shard.m[userID]
	return manager, ok
}

// GetOrAdd Returns the manager. If it didn't exist, it creates it.
func (m *ShardedUserSessionsMap) GetOrAdd(userID int64) *UserSessionsManager {
	shard := m.getShard(userID)

	shard.RLock()
	manager, ok := shard.m[userID]
	shard.RUnlock()
	if ok {
		return manager
	}

	shard.Lock()
	defer shard.Unlock()
	// Double check
	if manager, ok = shard.m[userID]; ok {
		return manager
	}
	manager = NewUserSessionsManager(userID)
	shard.m[userID] = manager
	return manager
}

func (m *ShardedUserSessionsMap) RemoveIfEmpty(userID int64) *UserSessionsManager {
	shard := m.getShard(userID)
	shard.Lock()
	defer shard.Unlock()
	if manager, ok := shard.m[userID]; ok {
		if manager.IsEmpty() {
			delete(shard.m, userID)
		}
		return manager
	}
	return nil
}

func (m *ShardedUserSessionsMap) CountOnlineUsers() int {
	count := 0
	for _, shard := range m.shards {
		shard.RLock()
		count += len(shard.m)
		shard.RUnlock()
	}
	return count
}

func (m *ShardedUserSessionsMap) Range(f func(int64, *UserSessionsManager) bool) {
	for _, shard := range m.shards {
		shard.RLock()
		managers := make([]*UserSessionsManager, 0, len(shard.m))
		for _, manager := range shard.m {
			managers = append(managers, manager)
		}
		shard.RUnlock()

		for _, manager := range managers {
			if !f(manager.UserID, manager) {
				return
			}
		}
	}
}
