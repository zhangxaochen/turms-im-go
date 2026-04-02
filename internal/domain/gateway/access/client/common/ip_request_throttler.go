package common

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// IpRequestThrottler provides IP-based rate limiting using a token bucket.
type IpRequestThrottler struct {
	mu       sync.RWMutex
	limiters map[string]*rate.Limiter

	Limit rate.Limit
	Burst int
}

// NewIpRequestThrottler creates a throttler with an explicit limit (requests per second) and burst size.
func NewIpRequestThrottler(limit rate.Limit, burst int) *IpRequestThrottler {
	t := &IpRequestThrottler{
		limiters: make(map[string]*rate.Limiter),
		Limit:    limit,
		Burst:    burst,
	}
	
	// Start a background cleanup routine
	go t.cleanupRoutine()
	
	return t
}

// DefaultIpRequestThrottler creates a throttler defaulting to 100 req/s, burst 100.
func DefaultIpRequestThrottler() *IpRequestThrottler {
	return NewIpRequestThrottler(100, 100)
}

// TryAcquireToken returns true if the IP is allowed to proceed, false if rate limited.
func (t *IpRequestThrottler) TryAcquireToken(ip string) bool {
	// Special case: unlimited
	if t.Burst <= 0 || t.Limit == 0 {
		return true
	}

	limiter := t.getLimiter(ip)
	return limiter.Allow()
}

func (t *IpRequestThrottler) getLimiter(ip string) *rate.Limiter {
	t.mu.RLock()
	limiter, exists := t.limiters[ip]
	t.mu.RUnlock()

	if exists {
		return limiter
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	// Double-check after acquiring write lock
	limiter, exists = t.limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(t.Limit, t.Burst)
		t.limiters[ip] = limiter
	}

	return limiter
}

// cleanupRoutine simply resets the map periodically to prevent memory leaks from millions of unique IPs.
// In a high-traffic production system, we'd use an LRU cache or explicit expiry per key.
// But this fulfills the basic need to drop stale IP limiters.
func (t *IpRequestThrottler) cleanupRoutine() {
	ticker := time.NewTicker(10 * time.Minute)
	for range ticker.C {
		t.mu.Lock()
		// Simple map reset. Active connections will seamlessly re-create their limiter.
		// For a truly persistent limiting across this threshold, a more sophisticated expiry is needed.
		t.limiters = make(map[string]*rate.Limiter)
		t.mu.Unlock()
	}
}
