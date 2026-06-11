// Package transport defines transport-level interfaces for harmoni services.
//
// Server abstracts a network listener (HTTP, gRPC, etc.) with a lifecycle.
// Codec[Req, Resp] encodes and decodes request/response payloads for a
// specific transport and serialization format.
package transport

import (
	"context"
	"io"
)

// Server is a transport-level server (HTTP listener, gRPC server, etc.).
type Server interface {
	// Serve starts the server and blocks until Shutdown is called or the
	// context is cancelled.
	Serve(ctx context.Context) error

	// Shutdown gracefully stops the server, waiting for in-flight requests
	// to complete or the context to expire.
	Shutdown(ctx context.Context) error
}

// Codec encodes and decodes requests and responses for a specific transport
// and serialization format.
type Codec[Req, Resp any] interface {
	// Decode reads a request from the reader.
	Decode(ctx context.Context, r io.Reader) (Req, error)

	// Encode writes a response to the writer.
	Encode(ctx context.Context, w io.Writer, resp Resp) error
}
