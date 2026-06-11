// Package sd defines service discovery interfaces for harmoni.
//
// It provides abstractions for registering services and discovering instances:
//
//	instancer := sd.NewFixedInstancer("host1:8080", "host2:8080")
//	instances, err := instancer.Instances()
//
// Implementations for Consul, Etcd, DNS, and Eureka live in kit/sd/.
package sd

import "context"

// Registrar registers and deregisters a service instance.
type Registrar interface {
	// Register registers the service instance.
	Register(ctx context.Context) error

	// Deregister removes the service instance.
	Deregister(ctx context.Context) error
}

// Discoverer discovers service instances.
type Discoverer interface {
	// Discover returns the current set of instance addresses.
	Discover(ctx context.Context) ([]string, error)
}

// Instancer provides a push-based view of service instances.
// It combines discovery with change notifications.
type Instancer interface {
	Discoverer

	// Subscribe returns a channel that receives a value each time the set of
	// instances changes. The caller should read Instances() after receiving.
	Subscribe() <-chan struct{}

	// Instances returns the current set of instance addresses.
	Instances() ([]string, error)
}
