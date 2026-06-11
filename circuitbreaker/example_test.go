package circuitbreaker_test

import (
	"context"
	"fmt"

	"github.com/harmonikit/harmoni/circuitbreaker"
	"github.com/harmonikit/harmoni/endpoint"
)

func ExampleNopCircuitBreaker() {
	cb := circuitbreaker.NewNopCircuitBreaker[int, string]()

	ep := endpoint.Endpoint[int, string](func(ctx context.Context, req int) (string, error) {
		return fmt.Sprintf("result-%d", req), nil
	})

	resp, _ := cb.Execute(context.Background(), 42, ep)
	fmt.Printf("state=%s resp=%s\n", cb.State(), resp)
	// Output: state=closed resp=result-42
}
