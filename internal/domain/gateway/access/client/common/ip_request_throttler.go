package common

import (
	"sync"
	"time"
)

// TokenBucketContext provides shared configuration for all token buckets.
// @MappedFrom TokenBucketContext
type TokenBucketContext struct {
	Capacity            int32
	TokensPerPeriod     int32
	RefillIntervalNanos int64
}

// TokenBucket replicates the Java TokenBucket algorithm with discrete refills.
// @MappedFrom TokenBucket
type TokenBucket struct {
	mu                  sync.Mutex
	tokens              int32
	lastRefillTimeNanos int64
}

// tryAcquire replicates Java's token acquisition logic with a shared context.
func (b *TokenBucket) tryAcquire(ctx *TokenBucketContext, timestamp int64) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Refill phase
	if timestamp >= b.lastRefillTimeNanos+ctx.RefillIntervalNanos {
		periods := (timestamp - b.lastRefillTimeNanos) / ctx.RefillIntervalNanos
		addedTokens := int32(periods) * ctx.TokensPerPeriod

		if addedTokens < 0 { // overflow
			b.tokens = ctx.Capacity
		} else {
			b.tokens += addedTokens
			if b.tokens > ctx.Capacity {
				b.tokens = ctx.Capacity
			}
		}
		b.lastRefillTimeNanos += periods * ctx.RefillIntervalNanos
	}

	if b.tokens > 0 {
		b.tokens--
		return true
	}
	return false
}

// IpRequestThrottler provides IP-based rate limiting using a custom token bucket.
type IpRequestThrottler struct {
	mu       sync.RWMutex
	limiters map[string]*TokenBucket
	ctx      *TokenBucketContext
}

// NewIpRequestThrottler creates a throttler with the specified context configuration.
func NewIpRequestThrottler(capacity int32, tokensPerPeriod int32, refillIntervalMillis int64) *IpRequestThrottler {
	ctx := &TokenBucketContext{
		Capacity:            capacity,
		TokensPerPeriod:     tokensPerPeriod,
		RefillIntervalNanos: refillIntervalMillis * 1000000,
	}

	t := &IpRequestThrottler{
		limiters: make(map[string]*TokenBucket),
		ctx:      ctx,
	}

	go t.cleanupRoutine()

	return t
}

// DefaultIpRequestThrottler creates a default fallback throttler.
func DefaultIpRequestThrottler() *IpRequestThrottler {
	// e.g. 100 tokens, 100 tokens per 1000ms
	return NewIpRequestThrottler(100, 100, 1000)
}

// UpdateContext allows dynamic updates shared by all buckets.
func (t *IpRequestThrottler) UpdateContext(capacity int32, tokensPerPeriod int32, refillIntervalMillis int64) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.ctx = &TokenBucketContext{
		Capacity:            capacity,
		TokensPerPeriod:     tokensPerPeriod,
		RefillIntervalNanos: refillIntervalMillis * 1000000,
	}
}

// TryAcquireToken returns true if the IP is allowed to proceed, false if rate limited.
// @MappedFrom tryAcquireToken(ByteArrayWrapper ip, long timestamp)
func (t *IpRequestThrottler) TryAcquireToken(ip string, timestamp int64) bool {
	// Fast-path read of shared context
	t.mu.RLock()
	ctx := t.ctx
	t.mu.RUnlock()

	// Special case: unlimited
	if ctx.RefillIntervalNanos <= 0 || ctx.TokensPerPeriod <= 0 || ctx.Capacity <= 0 {
		// When unlimited in Java conceptually by disabling refill, it returns false when empty,
		// but unlimited means RefillInterval <= 0 usually. We'll replicate the core logic:
		// If setup invalid/unlimited, we bypass throttling and allow it.
		// Wait, Java actually says: if refillInterval <= 0 it never refills.
		// If tokensPerPeriod/capacity are 0, it never has tokens.
		// Let's assume Capacity=0 means unlimited? Nope, Java treats it literally.
		// But in Java: "bucket is 'unlimited' depends on TokenBucketContext... if refillIntervalNanos <= 0, bucket returns false when empty".
		// Actually, if we just let the bucket run, we get the exact behaviour.
	}

	bucket := t.getBucket(ip, timestamp)
	return bucket.tryAcquire(ctx, timestamp)
}

func (t *IpRequestThrottler) getBucket(ip string, timestamp int64) *TokenBucket {
	t.mu.RLock()
	bucket, exists := t.limiters[ip]
	t.mu.RUnlock()

	if exists {
		return bucket
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	// Double-check
	bucket, exists = t.limiters[ip]
	if !exists {
		bucket = &TokenBucket{
			tokens:              t.ctx.Capacity, // Start at max capacity
			lastRefillTimeNanos: timestamp,      // Current time
		}
		t.limiters[ip] = bucket
	}

	return bucket
}

// CleanupByIp provides the session closed hook for removing unneeded buckets.
func (t *IpRequestThrottler) CleanupByIp(ip string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	bucket, exists := t.limiters[ip]
	if !exists {
		return
	}

	bucket.mu.Lock()
	// Only remove if it's fully replenished
	if bucket.tokens >= t.ctx.Capacity {
		delete(t.limiters, ip)
	}
	bucket.mu.Unlock()
}

func (t *IpRequestThrottler) cleanupRoutine() {
	ticker := time.NewTicker(30 * time.Minute)
	for range ticker.C {
		now := time.Now().UnixNano()
		t.mu.Lock()
		ctx := t.ctx

		for ip, bucket := range t.limiters {
			bucket.mu.Lock()

			isIdle := now-bucket.lastRefillTimeNanos > 30*time.Minute.Nanoseconds()
			// fully replenished
			isReplenished := false
			if bucket.tokens >= ctx.Capacity {
				isReplenished = true
			} else {
				// Simulate refill to check if it WOULD be replenished
				if now >= bucket.lastRefillTimeNanos+ctx.RefillIntervalNanos {
					periods := (now - bucket.lastRefillTimeNanos) / ctx.RefillIntervalNanos
					addedTokens := int32(periods) * ctx.TokensPerPeriod
					if addedTokens < 0 || bucket.tokens+addedTokens >= ctx.Capacity {
						isReplenished = true
					}
				}
			}

			if isIdle && isReplenished {
				delete(t.limiters, ip)
			}
			bucket.mu.Unlock()
		}
		t.mu.Unlock()
	}
}
