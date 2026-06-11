// Package harmoni provides the core interfaces, types, and patterns for building
// type-safe microservices in Go.
//
// harmoni is the "contract" — zero external dependencies, stdlib-only.
// It defines what a service looks like (endpoint.Endpoint), how it communicates
// (transport.Server, transport.Codec), and the cross-cutting middleware interfaces
// that wrap it (log.Logger, metrics.Counter, tracing.Tracer, etc.).
//
// For implementations, see the kit modules at github.com/harmonikit/kit.
package harmoni
