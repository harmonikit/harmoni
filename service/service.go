// Package service defines the Service interface for harmoni.
//
// A Service is a named collection of endpoints that together implement a
// business capability. Middleware wraps services to add behavior like
// logging, metrics, or initialization — but cross-cutting RPC concerns
// (auth, rate limiting, tracing) should wrap endpoints, not services.
package service

// Service is a named collection of endpoints.
type Service interface {
	// Name returns the service name, used for logging and metrics labels.
	Name() string
}

// Middleware wraps a Service to add cross-cutting behavior at the
// service level. For RPC-level middleware (auth, rate limiting, tracing),
// use endpoint.Middleware instead.
type Middleware func(Service) Service
