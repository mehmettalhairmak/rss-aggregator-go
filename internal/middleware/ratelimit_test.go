package middleware

import (
	"sync"
	"testing"
	"time"

	"github.com/mehmettalhairmak/rss-aggregator/internal/logger"
)

func init() {
	// Initialize logger for tests
	logger.InitLogger()
}

func TestNewTokenBucket(t *testing.T) {
	capacity := 10.0
	refillRate := 1.0

	tb := NewTokenBucket(capacity, refillRate)

	if tb.tokens != capacity {
		t.Errorf("Expected tokens to be %f, got %f", capacity, tb.tokens)
	}

	if tb.capacity != capacity {
		t.Errorf("Expected capacity to be %f, got %f", capacity, tb.capacity)
	}

	if tb.refillRate != refillRate {
		t.Errorf("Expected refillRate to be %f, got %f", refillRate, tb.refillRate)
	}
}

func TestTokenBucketConsume_Success(t *testing.T) {
	tb := NewTokenBucket(10.0, 1.0)

	// First consume should succeed
	if !tb.Consume() {
		t.Error("Expected first consume to succeed")
	}

	// Should have 9 tokens left
	if tb.tokens != 9.0 {
		t.Errorf("Expected 9 tokens, got %f", tb.tokens)
	}
}

func TestTokenBucketConsume_EmptyBucket(t *testing.T) {
	tb := NewTokenBucket(1.0, 1.0)

	// First consume should succeed
	if !tb.Consume() {
		t.Error("Expected first consume to succeed")
	}

	// Second consume should fail (no tokens)
	if tb.Consume() {
		t.Error("Expected second consume to fail")
	}
}

func TestTokenBucketConsume_Refill(t *testing.T) {
	tb := NewTokenBucket(10.0, 1.0)

	// Consume all tokens
	for i := 0; i < 10; i++ {
		tb.Consume()
	}

	// Wait for refill (1 token per second)
	time.Sleep(1100 * time.Millisecond)

	// Should be able to consume now
	if !tb.Consume() {
		t.Error("Expected consume to succeed after refill")
	}
}

func TestTokenBucketConsume_Burst(t *testing.T) {
	tb := NewTokenBucket(5.0, 0.5) // Refills at 0.5 tokens/second

	// Consume all 5 tokens (burst)
	for i := 0; i < 5; i++ {
		if !tb.Consume() {
			t.Errorf("Expected consume %d to succeed", i+1)
		}
	}

	// Should have no tokens left
	if tb.tokens > 0.01 { // Allow for floating point precision
		t.Errorf("Expected 0 tokens, got %f", tb.tokens)
	}

	// Next consume should fail
	if tb.Consume() {
		t.Error("Expected consume to fail after burst")
	}
}

func TestTokenBucketConsume_CapacityCap(t *testing.T) {
	tb := NewTokenBucket(10.0, 1.0)

	// Wait long enough to exceed capacity
	time.Sleep(15000 * time.Millisecond) // 15 seconds = 15 tokens

	// Should still be capped at capacity
	if tb.tokens != 10.0 {
		t.Errorf("Expected tokens to be capped at 10, got %f", tb.tokens)
	}
}

func TestTokenBucketConsume_Concurrent(t *testing.T) {
	tb := NewTokenBucket(100.0, 10.0)

	var wg sync.WaitGroup
	successCount := 0
	var mu sync.Mutex

	// Spawn 100 goroutines trying to consume
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if tb.Consume() {
				mu.Lock()
				successCount++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	// Exactly 100 should succeed
	if successCount != 100 {
		t.Errorf("Expected 100 successful consumes, got %d", successCount)
	}
}

func TestInitRateLimiter(t *testing.T) {
	config := RateLimitConfig{
		RequestsPerMinute: 60,
		BurstSize:         10,
	}

	InitRateLimiter(config)

	if limiter == nil {
		t.Error("Expected limiter to be initialized")
	}

	// Should allow requests up to burst size
	for i := 0; i < 10; i++ {
		if !limiter.Consume() {
			t.Errorf("Expected consume %d to succeed", i+1)
		}
	}
}

func BenchmarkTokenBucketConsume(b *testing.B) {
	tb := NewTokenBucket(1000.0, 1000.0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tb.Consume()
	}
}
