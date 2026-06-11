package service_test

import (
	"testing"

	"github.com/harmonikit/harmoni/service"
)

// testService implements service.Service for testing.
type testService struct {
	name string
}

func (s *testService) Name() string { return s.name }

// Compile-time interface check.
var _ service.Service = (*testService)(nil)

func TestService_Name(t *testing.T) {
	svc := &testService{name: "addsvc"}

	if got := svc.Name(); got != "addsvc" {
		t.Fatalf("got %q, want %q", got, "addsvc")
	}
}

func TestServiceMiddleware_Identity(t *testing.T) {
	identity := service.ServiceMiddleware(func(s service.Service) service.Service {
		return s
	})

	svc := &testService{name: "usersvc"}
	wrapped := identity(svc)

	if got := wrapped.Name(); got != "usersvc" {
		t.Fatalf("got %q, want %q", got, "usersvc")
	}
}

func TestServiceMiddleware_Transform(t *testing.T) {
	// Middleware that prefixes the service name.
	prefixMW := service.ServiceMiddleware(func(s service.Service) service.Service {
		// Return a new service with a prefixed name.
		return &prefixedService{inner: s, prefix: "v2."}
	})

	svc := &testService{name: "profilesvc"}
	wrapped := prefixMW(svc)

	want := "v2.profilesvc"
	if got := wrapped.Name(); got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

// prefixedService wraps a Service and adds a prefix to its name.
type prefixedService struct {
	inner  service.Service
	prefix string
}

func (s *prefixedService) Name() string {
	return s.prefix + s.inner.Name()
}
