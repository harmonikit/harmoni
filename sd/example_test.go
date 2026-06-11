package sd_test

import (
	"fmt"

	"github.com/harmonikit/harmoni/sd"
)

func ExampleFixedInstancer() {
	fi := sd.NewFixedInstancer("server1:8080", "server2:8080")

	instances, _ := fi.Instances()
	for _, addr := range instances {
		fmt.Println(addr)
	}
	// Output:
	// server1:8080
	// server2:8080
}
