// Package ratelimit defines rate limiting interfaces and a stdlib-only token
// bucket implementation.
//
// Example:
//
//	limiter := ratelimit.NewTokenBucket(100, 10) // 100 req/s, burst 10
//	if limiter.Allow() {
//	    // handle request
//	}
package ratelimit

// Limiter is the interface for a rate limiter. Implementations must be safe
// for concurrent use.
type Limiter interface {
	// Allow reports whether a single request is allowed at this moment.
	Allow() bool
}
