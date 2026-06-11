package sd

import "context"

// FixedInstancer returns a fixed set of instances. It is useful for testing
// and for static configurations.
type FixedInstancer struct {
	instances []string
}

// NewFixedInstancer creates a FixedInstancer with the given addresses.
func NewFixedInstancer(instances ...string) *FixedInstancer {
	return &FixedInstancer{instances: instances}
}

// Discover returns the fixed set of instances.
func (f *FixedInstancer) Discover(ctx context.Context) ([]string, error) {
	result := make([]string, len(f.instances))
	copy(result, f.instances)
	return result, nil
}

// Subscribe returns a channel that never fires (instances never change).
func (f *FixedInstancer) Subscribe() <-chan struct{} {
	return nil
}

// Instances returns the current set of instance addresses.
func (f *FixedInstancer) Instances() ([]string, error) {
	return f.Discover(context.Background())
}
