package auth

import "context"

type contextKey struct{}

// SetAuth stores an authentication value in the context.
func SetAuth(ctx context.Context, auth any) context.Context {
	return context.WithValue(ctx, contextKey{}, auth)
}

// GetAuth retrieves the authentication value from the context.
// Returns nil, false if no auth value is present.
func GetAuth(ctx context.Context) (any, bool) {
	v := ctx.Value(contextKey{})
	if v == nil {
		return nil, false
	}
	return v, true
}

// ClearAuth returns a context with the auth value removed.
func ClearAuth(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextKey{}, nil)
}
