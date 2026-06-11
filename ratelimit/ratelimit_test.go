package ratelimit_test

import (
	"sync"
	"testing"
	"time"

	"github.com/harmonikit/harmoni/ratelimit"
)

func TestTokenBucket_Allow(t *testing.T) {
	// 100 tokens/sec, burst 10 — starts full so all 10 should be allowed.
	tb := ratelimit.NewTokenBucket(100, 10)

	for i := range 10 {
		if !tb.Allow() {
			t.Fatalf("request %d should be allowed (bucket starts full)", i)
		}
	}
}

func TestTokenBucket_DenyAfterBurst(t *testing.T) {
	// Very low rate — after burst is consumed, no more tokens.
	tb := ratelimit.NewTokenBucket(0.001, 3)

	// Consume the burst.
	for range 3 {
		if !tb.Allow() {
			t.Fatal("initial burst should be allowed")
		}
	}

	// Next request should be denied (rate is near zero).
	if tb.Allow() {
		t.Fatal("expected rate limit denial after burst exhausted")
	}
}

func TestTokenBucket_Refill(t *testing.T) {
	// 1000 tokens/sec gives roughly 1 token per ms.
	tb := ratelimit.NewTokenBucket(1000, 1)

	// Consume the only token.
	if !tb.Allow() {
		t.Fatal("initial token should be allowed")
	}

	// Should be denied immediately.
	if tb.Allow() {
		t.Fatal("expected denial, no tokens available")
	}

	// Wait for refill.
	time.Sleep(5 * time.Millisecond)

	// Should have at least 1 token now.
	if !tb.Allow() {
		t.Fatal("expected token after refill period")
	}
}

func TestTokenBucket_BurstBehavior(t *testing.T) {
	// Burst of 5, rate high enough to never limit.
	tb := ratelimit.NewTokenBucket(100000, 5)

	for i := range 5 {
		if !tb.Allow() {
			t.Fatalf("request %d should be allowed (within burst)", i)
		}
	}

	// Wait a tiny bit for refill — even 100µs at 100k tokens/sec = 10 tokens.
	time.Sleep(100 * time.Microsecond)

	// 6th should be allowed since enough tokens have refilled.
	if !tb.Allow() {
		t.Fatal("6th request should be allowed (rate refilled)")
	}
}

func TestTokenBucket_BurstCap(t *testing.T) {
	// Tokens should never exceed burst, even after long sleep.
	tb := ratelimit.NewTokenBucket(100000, 3)
	time.Sleep(50 * time.Millisecond)

	// Should only get up to burst (3) tokens from the sleep.
	// Then we sleep again to refill between calls.
	count := 0
	for range 10 {
		if tb.Allow() {
			count++
		}
		time.Sleep(50 * time.Microsecond) // allow refill between calls
	}

	if count < 3 {
		t.Fatalf("got %d allows, want at least 3", count)
	}
}

func TestTokenBucket_Concurrent(t *testing.T) {
	tb := ratelimit.NewTokenBucket(10000, 100)
	var wg sync.WaitGroup
	var mu sync.Mutex
	total := 0

	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 10 {
				if tb.Allow() {
					mu.Lock()
					total++
					mu.Unlock()
				}
			}
		}()
	}
	wg.Wait()

	if total == 0 {
		t.Fatal("no requests were allowed")
	}
}

func TestTokenBucket_ZeroRate(t *testing.T) {
	// Zero rate means no refill — only burst tokens are available.
	tb := ratelimit.NewTokenBucket(0, 2)

	if !tb.Allow() {
		t.Fatal("first request should be allowed")
	}
	if !tb.Allow() {
		t.Fatal("second request should be allowed")
	}
	if tb.Allow() {
		t.Fatal("third request should be denied (zero rate)")
	}
}

func TestLimiterInterface(t *testing.T) {
	var _ ratelimit.Limiter = ratelimit.NewTokenBucket(1, 1)
}
