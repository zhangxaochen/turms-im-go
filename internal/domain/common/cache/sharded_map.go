package cache

import (
	"sync"
)

// ShardedMap is a concurrent-safe map that uses sharding to reduce lock contention.
type ShardedMap[K comparable, V any] struct {
	shards []*shard[K, V]
	num    uint32
	hasher func(K) uint32
}

type shard[K comparable, V any] struct {
	sync.RWMutex
	m map[K]V
}

// NewShardedMap creates a new ShardedMap with the given number of shards and a hasher.
func NewShardedMap[K comparable, V any](numShards uint32, hasher func(K) uint32) *ShardedMap[K, V] {
	s := &ShardedMap[K, V]{
		shards: make([]*shard[K, V], numShards),
		num:    numShards,
		hasher: hasher,
	}
	for i := uint32(0); i < numShards; i++ {
		s.shards[i] = &shard[K, V]{
			m: make(map[K]V),
		}
	}
	return s
}

// NewStringShardedMap is a convenience function for string keys.
func NewStringShardedMap[V any](numShards uint32) *ShardedMap[string, V] {
	return NewShardedMap[string, V](numShards, Fnv32)
}

func (s *ShardedMap[K, V]) getShard(key K) *shard[K, V] {
	return s.shards[s.hasher(key)%s.num]
}

func (s *ShardedMap[K, V]) Set(key K, value V) {
	shard := s.getShard(key)
	shard.Lock()
	shard.m[key] = value
	shard.Unlock()
}

// @MappedFrom get(UdpNotificationType type)
// @MappedFrom get(ResponseStatusCode code, String reason)
// @MappedFrom get(ResponseStatusCode code)
func (s *ShardedMap[K, V]) Get(key K) (V, bool) {
	shard := s.getShard(key)
	shard.RLock()
	val, ok := shard.m[key]
	shard.RUnlock()
	return val, ok
}

// @MappedFrom delete(@Nullable Set<Long> groupIds, @Nullable ClientSession session)
// @MappedFrom delete(@NotEmpty Set<Long> userIds, @Nullable ClientSession session)
func (s *ShardedMap[K, V]) Delete(key K) {
	shard := s.getShard(key)
	shard.Lock()
	delete(shard.m, key)
	shard.Unlock()
}

// Fnv32 is a simple hash function for strings.
func Fnv32(key string) uint32 {
	hash := uint32(2166136261)
	const prime32 = uint32(16777619)
	for i := 0; i < len(key); i++ {
		hash *= prime32
		hash ^= uint32(key[i])
	}
	return hash
}
