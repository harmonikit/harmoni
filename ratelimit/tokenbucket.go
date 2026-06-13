package ratelimit

import (
	"sync"
	"time"
)

// TokenBucket implements the token bucket rate limiting algorithm.
// Tokens are added at a fixed rate and consumed one per Allow() call.
// It is safe for concurrent use.
type TokenBucket struct {
	mu       sync.Mutex
	rate     float64   // tokens per second
	burst    float64   // maximum tokens
	tokens   float64   // current tokens
	lastTime time.Time // last time tokens were updated
}

// NewTokenBucket creates a TokenBucket with the given rate (tokens per second)
// and burst (maximum tokens). The bucket starts full.
func NewTokenBucket(rate float64, burst float64) *TokenBucket {
	return &TokenBucket{
		mu:       sync.Mutex{},
		rate:     rate,
		burst:    burst,
		tokens:   burst,
		lastTime: time.Now(),
	}
}

// Allow reports whether a request is allowed. Each call consumes one token
// if available.
func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastTime).Seconds()
	tb.lastTime = now

	// Add tokens earned since last check, capped at burst.
	tb.tokens += elapsed * tb.rate
	if tb.tokens > tb.burst {
		tb.tokens = tb.burst
	}

	if tb.tokens >= 1 {
		tb.tokens--
		return true
	}
	return false
}
