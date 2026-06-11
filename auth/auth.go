// Package auth defines authentication interfaces and context helpers for
// harmoni endpoints.
//
// Example:
//
//	type myAuth struct{}
//	func (a *myAuth) Authenticate(ctx context.Context, req *MyRequest) (context.Context, error) {
//	    ctx = auth.SetAuth(ctx, &User{ID: "123"})
//	    return ctx, nil
//	}
//	ep = auth.AuthMiddleware[MyRequest, MyResponse](&myAuth{})(ep)
package auth

import (
	"context"

	"github.com/harmonikit/harmoni/endpoint"
)

// Auth is an authentication strategy for a request type. It authenticates
// the request and returns a context with auth information embedded.
type Auth[Req any] interface {
	// Authenticate validates the request and returns a context carrying
	// authentication results (e.g., user identity, claims).
	Authenticate(ctx context.Context, req Req) (context.Context, error)
}

// AuthMiddleware returns an endpoint middleware that authenticates requests
// before passing them to the endpoint. If authentication fails, the endpoint
// is not called and the error is returned.
func AuthMiddleware[Req, Resp any](a Auth[Req]) endpoint.Middleware[Req, Resp] {
	return func(next endpoint.Endpoint[Req, Resp]) endpoint.Endpoint[Req, Resp] {
		return func(ctx context.Context, req Req) (Resp, error) {
			authCtx, err := a.Authenticate(ctx, req)
			if err != nil {
				var zero Resp
				return zero, err
			}
			return next(authCtx, req)
		}
	}
}
