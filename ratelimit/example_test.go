package ratelimit_test

import (
	"fmt"

	"github.com/harmonikit/harmoni/ratelimit"
)

func ExampleTokenBucket() {
	tb := ratelimit.NewTokenBucket(100, 3)

	for i := range 5 {
		if tb.Allow() {
			fmt.Printf("request %d: allowed\n", i+1)
		} else {
			fmt.Printf("request %d: denied\n", i+1)
		}
	}
	// Output:
	// request 1: allowed
	// request 2: allowed
	// request 3: allowed
	// request 4: denied
	// request 5: denied
}
