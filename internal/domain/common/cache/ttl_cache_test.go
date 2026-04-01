package cache

import (
	"testing"
	"time"
)

func TestTTLCache(t *testing.T) {
	cache := NewTTLCache[string, string](50*time.Millisecond, 10*time.Millisecond)
	defer cache.Close()

	// Test Set and Get
	cache.Set("key1", "value1")
	val, ok := cache.Get("key1")
	if !ok || val != "value1" {
		t.Errorf("Expected to get 'value1', got '%v' (ok=%v)", val, ok)
	}

	// Test non-existent key
	_, ok = cache.Get("key2")
	if ok {
		t.Error("Expected not to find 'key2'")
	}

	// Test expiry
	time.Sleep(100 * time.Millisecond)
	_, ok = cache.Get("key1")
	if ok {
		t.Error("Expected 'key1' to be expired")
	}

	// Test Delete
	cache.Set("key3", "value3")
	cache.Delete("key3")
	_, ok = cache.Get("key3")
	if ok {
		t.Error("Expected 'key3' to be deleted")
	}
}

func BenchmarkTTLCache_Set_Get(b *testing.B) {
	cache := NewTTLCache[int, string](1*time.Minute, 1*time.Minute)
	defer cache.Close()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			cache.Set(i, "test")
			cache.Get(i)
			i++
		}
	})
}
