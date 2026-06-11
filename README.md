# harmoni — core microservice interfaces for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/harmonikit/harmoni.svg)](https://pkg.go.dev/github.com/harmonikit/harmoni)
[![Go Report Card](https://goreportcard.com/badge/github.com/harmonikit/harmoni)](https://goreportcard.com/report/github.com/harmonikit/harmoni)
[![CI](https://github.com/harmonikit/harmoni/actions/workflows/ci.yml/badge.svg)](https://github.com/harmonikit/harmoni/actions/workflows/ci.yml)

Zero-dependency Go library providing the core interfaces, types, and patterns for
building type-safe microservices with Go 1.23+ generics.

## Packages

| Package | Purpose |
|---|---|
| `endpoint` | `Endpoint[Req, Resp]` — the universal RPC abstraction |
| `transport` | `Server`, `Codec[Req, Resp]` — transport-level interfaces |
| `service` | `Service` — named collection of endpoints |
| `log` | `Logger` — structured logging interface + slog adapter |
| `metrics` | `Counter`, `Gauge`, `Histogram` — metrics interfaces |
| `middleware` | `Timeout`, `Retry`, `Recovery` — endpoint middleware |
| `auth` | `Auth[Req]` — authentication interface |
| `circuitbreaker` | `CircuitBreaker[Req, Resp]` — circuit breaker interface |
| `ratelimit` | `Limiter` — rate limiting interface + token bucket |
| `sd` | `Registrar`, `Discoverer`, `Instancer` — service discovery |
| `tracing` | `Tracer[Req, Resp]`, `Span` — distributed tracing interfaces |

## Install

```bash
go get github.com/harmonikit/harmoni@latest
```

## Requirements

- Go 1.23+
- Zero external dependencies

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/harmonikit/harmoni/endpoint"
    "github.com/harmonikit/harmoni/middleware"
)

func main() {
    add := endpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
        return req + 1, nil
    })

    add = endpoint.Chain(
        middleware.Timeout[int, int](5*time.Second),
        middleware.Recovery[int, int](),
    )(add)

    resp, err := add(context.Background(), 41)
    if err != nil {
        panic(err)
    }
    fmt.Println(resp) // 42
}
```

## License

MIT
