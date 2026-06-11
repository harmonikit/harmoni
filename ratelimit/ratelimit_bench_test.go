package ratelimit_test

import (
	"testing"

	"github.com/harmonikit/harmoni/ratelimit"
)

func BenchmarkTokenBucket_Allow(b *testing.B) {
	tb := ratelimit.NewTokenBucket(100000, 100000) // never limits
	for range b.N {
		tb.Allow()
	}
}

func BenchmarkTokenBucket_Allow_Contended(b *testing.B) {
	tb := ratelimit.NewTokenBucket(100000, 100000)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			tb.Allow()
		}
	})
}
