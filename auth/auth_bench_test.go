package auth_test

import (
	"context"
	"testing"

	"github.com/harmonikit/harmoni/auth"
)

func BenchmarkSetAuth(b *testing.B) {
	ctx := context.Background()
	for range b.N {
		_ = auth.SetAuth(ctx, "user")
	}
}

func BenchmarkGetAuth(b *testing.B) {
	ctx := auth.SetAuth(context.Background(), "user")
	for range b.N {
		_, _ = auth.GetAuth(ctx)
	}
}
