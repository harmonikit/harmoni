package service_test

import (
	"testing"

	"github.com/harmonikit/harmoni/service"
)

func BenchmarkService_Name(b *testing.B) {
	svc := &testService{name: "bench-svc"}
	for range b.N {
		_ = svc.Name()
	}
}

func BenchmarkServiceMiddleware(b *testing.B) {
	identity := service.ServiceMiddleware(func(s service.Service) service.Service {
		return s
	})
	svc := &testService{name: "bench-svc"}

	for range b.N {
		_ = identity(svc)
	}
}
