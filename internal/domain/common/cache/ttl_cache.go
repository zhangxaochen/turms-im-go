package cache

import (
	"sync"
	"time"
)

type item[V any] struct {
	value    V
	deadline int64
}

// TTLCache is a generic, goroutine-safe local cache with a TTL mechanism.
type TTLCache[K comparable, V any] struct {
	m    sync.RWMutex
	data map[K]item[V]
	ttl  time.Duration
	stop chan struct{}
}

// NewTTLCache creates a new TTLCache instance.
func NewTTLCache[K comparable, V any](ttl time.Duration, cleanupInterval time.Duration) *TTLCache[K, V] {
	c := &TTLCache[K, V]{
		data: make(map[K]item[V]),
		ttl:  ttl,
		stop: make(chan struct{}),
	}
	go c.cleanupLoop(cleanupInterval)
	return c
}

// Set adds or updates an item in the cache.
func (c *TTLCache[K, V]) Set(key K, value V) {
	deadline := time.Now().Add(c.ttl).UnixNano()

	c.m.Lock()
	defer c.m.Unlock()
	c.data[key] = item[V]{
		value:    value,
		deadline: deadline,
	}
}

// Get retrieves an item from the cache. Returns (value, true) if found and not expired.
func (c *TTLCache[K, V]) Get(key K) (V, bool) {
	c.m.RLock()
	i, ok := c.data[key]
	c.m.RUnlock()

	if !ok {
		var zero V
		return zero, false
	}
	if time.Now().UnixNano() > i.deadline {
		var zero V
		return zero, false
	}
	return i.value, true
}

// Delete removes an item from the cache.
func (c *TTLCache[K, V]) Delete(key K) {
	c.m.Lock()
	defer c.m.Unlock()
	delete(c.data, key)
}

// Close stops the background cleanup goroutine.
func (c *TTLCache[K, V]) Close() {
	select {
	case <-c.stop:
		// already closed
	default:
		close(c.stop)
	}
}

func (c *TTLCache[K, V]) cleanupLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			now := time.Now().UnixNano()
			c.m.Lock()
			for k, v := range c.data {
				if now > v.deadline {
					delete(c.data, k)
				}
			}
			c.m.Unlock()
		case <-c.stop:
			return
		}
	}
}
