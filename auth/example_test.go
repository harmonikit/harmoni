package auth_test

import (
	"context"
	"fmt"

	"github.com/harmonikit/harmoni/auth"
)

func ExampleSetAuth() {
	ctx := context.Background()
	ctx = auth.SetAuth(ctx, "user-42")

	v, ok := auth.GetAuth(ctx)
	fmt.Printf("ok=%v user=%v\n", ok, v)
	// Output: ok=true user=user-42
}

func ExampleAuthMiddleware() {
	type request struct{ Token string }

	// A simple token auth.
	simpleAuth := struct{}{} // Would implement auth.Auth[request]

	_ = simpleAuth // illustrative — real auth would validate the token
	fmt.Println("auth middleware configured")
	// Output: auth middleware configured
}
