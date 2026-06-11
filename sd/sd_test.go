package sd_test

import (
	"context"
	"testing"
	"time"

	"github.com/harmonikit/harmoni/sd"
)

func TestFixedInstancer_Discover(t *testing.T) {
	fi := sd.NewFixedInstancer("host1:8080", "host2:8080")

	instances, err := fi.Discover(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(instances) != 2 {
		t.Fatalf("got %d instances, want 2", len(instances))
	}
	if instances[0] != "host1:8080" {
		t.Fatalf("got %q, want %q", instances[0], "host1:8080")
	}
	if instances[1] != "host2:8080" {
		t.Fatalf("got %q, want %q", instances[1], "host2:8080")
	}
}

func TestFixedInstancer_Discover_Empty(t *testing.T) {
	fi := sd.NewFixedInstancer()

	instances, err := fi.Discover(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(instances) != 0 {
		t.Fatalf("got %d instances, want 0", len(instances))
	}
}

func TestFixedInstancer_Discover_Immutability(t *testing.T) {
	fi := sd.NewFixedInstancer("host1:8080")

	instances, _ := fi.Discover(context.Background())
	// Mutate the returned slice — should not affect the FixedInstancer.
	instances[0] = "modified:8080"

	instances2, _ := fi.Discover(context.Background())
	if instances2[0] != "host1:8080" {
		t.Fatalf("got %q, want %q (mutating return value should not affect internal state)", instances2[0], "host1:8080")
	}
}

func TestFixedInstancer_Subscribe(t *testing.T) {
	fi := sd.NewFixedInstancer("host1:8080")

	ch := fi.Subscribe()
	if ch != nil {
		t.Fatal("FixedInstancer.Subscribe should return nil (never changes)")
	}
}

func TestFixedInstancer_Instances(t *testing.T) {
	fi := sd.NewFixedInstancer("a:1", "b:2")

	instances, err := fi.Instances()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(instances) != 2 {
		t.Fatalf("got %d instances, want 2", len(instances))
	}
}

func TestFixedInstancer_ContextCancelled(t *testing.T) {
	fi := sd.NewFixedInstancer("host:80")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	// FixedInstancer doesn't use context for network operations,
	// so it should still work even with a cancelled context.
	instances, err := fi.Discover(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(instances) != 1 {
		t.Fatalf("got %d instances, want 1", len(instances))
	}
}

func TestInstancerInterface(t *testing.T) {
	var _ sd.Discoverer = sd.NewFixedInstancer()
	var _ sd.Instancer = sd.NewFixedInstancer()
}
