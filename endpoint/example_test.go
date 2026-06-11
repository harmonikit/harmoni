package endpoint_test

import (
	"context"
	"fmt"
	"log"

	"github.com/harmonikit/harmoni/endpoint"
)

func ExampleEndpoint() {
	add := endpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		return req + 1, nil
	})

	resp, err := add(context.Background(), 41)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp)
	// Output: 42
}

func ExampleMiddleware() {
	// A logging middleware that logs every request.
	logMW := endpoint.Middleware[int, int](func(next endpoint.Endpoint[int, int]) endpoint.Endpoint[int, int] {
		return func(ctx context.Context, req int) (int, error) {
			fmt.Printf("request: %d\n", req)
			resp, err := next(ctx, req)
			fmt.Printf("response: %d, err: %v\n", resp, err)
			return resp, err
		}
	})

	add := endpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		return req + 1, nil
	})

	wrapped := logMW(add)
	resp, _ := wrapped(context.Background(), 41)
	fmt.Println(resp)
	// Output:
	// request: 41
	// response: 42, err: <nil>
	// 42
}

func ExampleChain() {
	doubleReq := endpoint.Middleware[int, int](func(next endpoint.Endpoint[int, int]) endpoint.Endpoint[int, int] {
		return func(ctx context.Context, req int) (int, error) {
			return next(ctx, req*2)
		}
	})

	addOne := endpoint.Middleware[int, int](func(next endpoint.Endpoint[int, int]) endpoint.Endpoint[int, int] {
		return func(ctx context.Context, req int) (int, error) {
			resp, err := next(ctx, req)
			return resp + 1, err
		}
	})

	ep := endpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		return req, nil
	})

	wrapped := endpoint.Chain(doubleReq, addOne)(ep)
	resp, _ := wrapped(context.Background(), 20)
	fmt.Println(resp)
	// Output: 41
}

func ExampleNop() {
	nop := endpoint.Nop[string, int]()
	resp, err := nop(context.Background(), "hello")
	fmt.Printf("resp=%d, err=%v\n", resp, err)
	// Output: resp=0, err=<nil>
}
