package sd_test

import (
	"context"
	"testing"

	"github.com/harmonikit/harmoni/sd"
)

func BenchmarkFixedInstancer_Discover(b *testing.B) {
	fi := sd.NewFixedInstancer("host1:8080", "host2:8080", "host3:8080")
	ctx := context.Background()

	for range b.N {
		_, _ = fi.Discover(ctx)
	}
}

func BenchmarkFixedInstancer_Instances(b *testing.B) {
	fi := sd.NewFixedInstancer("host1:8080", "host2:8080")

	for range b.N {
		_, _ = fi.Instances()
	}
}
