package transport_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"sync"
	"testing"

	"github.com/harmonikit/harmoni/transport"
)

// Compile-time interface assertions.
var (
	_ transport.Server              = (*mockServer)(nil)
	_ transport.Codec[string, string] = (*stringCodec)(nil)
)

// mockServer implements transport.Server for testing.
type mockServer struct {
	mu             sync.Mutex
	serveCalled    bool
	shutdownCalled bool
	serveErr       error
	shutdownErr    error
}

func (m *mockServer) Serve(ctx context.Context) error {
	m.mu.Lock()
	m.serveCalled = true
	m.mu.Unlock()
	return m.serveErr
}

func (m *mockServer) Shutdown(ctx context.Context) error {
	m.mu.Lock()
	m.shutdownCalled = true
	m.mu.Unlock()
	return m.shutdownErr
}

// stringCodec implements transport.Codec[string, string] for testing.
type stringCodec struct{}

func (c *stringCodec) Decode(ctx context.Context, r io.Reader) (string, error) {
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (c *stringCodec) Encode(ctx context.Context, w io.Writer, resp string) error {
	_, err := w.Write([]byte(resp))
	return err
}

func TestServer_Serve(t *testing.T) {
	ctx := context.Background()
	s := &mockServer{}

	err := s.Serve(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !s.serveCalled {
		t.Fatal("Serve was not called")
	}
}

func TestServer_Serve_Error(t *testing.T) {
	ctx := context.Background()
	wantErr := errors.New("listen error")
	s := &mockServer{serveErr: wantErr}

	err := s.Serve(ctx)
	if !errors.Is(err, wantErr) {
		t.Fatalf("got error %v, want %v", err, wantErr)
	}
}

func TestServer_Shutdown(t *testing.T) {
	ctx := context.Background()
	s := &mockServer{}

	err := s.Shutdown(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !s.shutdownCalled {
		t.Fatal("Shutdown was not called")
	}
}

func TestServer_Shutdown_Error(t *testing.T) {
	ctx := context.Background()
	wantErr := errors.New("shutdown timeout")
	s := &mockServer{shutdownErr: wantErr}

	err := s.Shutdown(ctx)
	if !errors.Is(err, wantErr) {
		t.Fatalf("got error %v, want %v", err, wantErr)
	}
}

func TestCodec_Decode(t *testing.T) {
	codec := &stringCodec{}
	input := strings.NewReader("hello, world")

	got, err := codec.Decode(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "hello, world" {
		t.Fatalf("got %q, want %q", got, "hello, world")
	}
}

func TestCodec_Encode(t *testing.T) {
	codec := &stringCodec{}
	var buf bytes.Buffer

	err := codec.Encode(context.Background(), &buf, "response")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.String() != "response" {
		t.Fatalf("got %q, want %q", buf.String(), "response")
	}
}

func TestCodec_Decode_Empty(t *testing.T) {
	codec := &stringCodec{}
	input := strings.NewReader("")

	got, err := codec.Decode(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "" {
		t.Fatalf("got %q, want empty string", got)
	}
}

func TestCodec_Encode_Empty(t *testing.T) {
	codec := &stringCodec{}
	var buf bytes.Buffer

	err := codec.Encode(context.Background(), &buf, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.String() != "" {
		t.Fatalf("got %q, want empty string", buf.String())
	}
}

func TestServer_Concurrent(t *testing.T) {
	s := &mockServer{}
	var wg sync.WaitGroup

	// Concurrent Serve calls should be safe.
	for range 5 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = s.Serve(context.Background())
		}()
	}
	wg.Wait()

	s.mu.Lock()
	called := s.serveCalled
	s.mu.Unlock()
	if !called {
		t.Fatal("Serve was never called")
	}
}
