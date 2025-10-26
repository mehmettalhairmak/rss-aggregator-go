package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/mehmettalhairmak/rss-aggregator/internal/logger"
	"github.com/mehmettalhairmak/rss-aggregator/internal/models"
)

// TokenBucket implements a simple token bucket rate limiter
type TokenBucket struct {
	tokens     float64
	capacity   float64
	refillRate float64 // tokens per second
	lastRefill time.Time
	mu         sync.Mutex
}

// RateLimitConfig holds configuration for rate limiting
type RateLimitConfig struct {
	RequestsPerMinute int
	BurstSize         int
}

// NewTokenBucket creates a new token bucket
func NewTokenBucket(capacity float64, refillRate float64) *TokenBucket {
	return &TokenBucket{
		tokens:     capacity,
		capacity:   capacity,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

// Consume attempts to consume one token from the bucket
// Returns true if a token was available, false otherwise
func (tb *TokenBucket) Consume() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastRefill).Seconds()

	// Refill tokens based on elapsed time
	tb.tokens = tb.tokens + elapsed*tb.refillRate

	// Cap at capacity
	if tb.tokens > tb.capacity {
		tb.tokens = tb.capacity
	}

	tb.lastRefill = now

	// Try to consume a token
	if tb.tokens >= 1.0 {
		tb.tokens--
		return true
	}

	return false
}

// Global rate limiter instances
var (
	limiter *TokenBucket
)

// InitRateLimiter initializes the global rate limiter
func InitRateLimiter(config RateLimitConfig) {
	// Convert requests per minute to tokens per second
	refillRate := float64(config.RequestsPerMinute) / 60.0
	limiter = NewTokenBucket(float64(config.BurstSize), refillRate)
	logger.Infof("Rate limiter initialized: %d requests/min, burst: %d",
		config.RequestsPerMinute, config.BurstSize)
}

// RateLimit is the rate limiting middleware
func RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if limiter == nil {
			// Rate limiter not initialized, allow request
			next.ServeHTTP(w, r)
			return
		}

		// Try to consume a token
		if !limiter.Consume() {
			// Rate limit exceeded
			logger.Debug("Rate limit exceeded for client")

			models.RespondWithError(w, http.StatusTooManyRequests,
				"Rate limit exceeded. Please try again later.")
			return
		}

		// Token consumed, proceed with request
		next.ServeHTTP(w, r)
	})
}
