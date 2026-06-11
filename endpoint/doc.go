// Package endpoint defines the universal RPC abstraction for harmonikit.
//
// Endpoint[Req, Resp] is a single request/response interaction — it could be
// an HTTP handler, a gRPC method, or an in-process function call.
//
// Middleware wraps endpoints to add cross-cutting behavior. All middleware is
// type-safe at compile time thanks to Go generics:
//
//	add := endpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
//	    return req + 1, nil
//	})
//	add = endpoint.Chain(loggingMiddleware, timeoutMiddleware)(add)
package endpoint
