# M01-L001: Bootstrap HTTP Service

## 1. Business Context

Before we can build any business features, the application needs to be a running HTTP service. This lesson creates the foundation that every future feature will build upon: a server that starts, responds to requests, and shuts down gracefully.

## 2. Problem

We need a production-quality HTTP server that:

- Starts on a configurable address
- Responds to health checks (used by load balancers, orchestrators)
- Reports its version
- Generates an OpenAPI spec automatically
- Shuts down gracefully (finishes in-flight requests when terminated)

## 3. Design Discussion

### Composition Root

In Go, the `main` function is the **composition root** — the place where dependencies are wired together. This is not a framework concept; it's explicit Go code.

```go
func main() {
    cfg := config.Load()
    srv := server.New(cfg)
    srv.Start(ctx)
}
```

Every dependency is created explicitly and passed as an argument. No global state, no service locator, no magic.

### Graceful Shutdown

When a process receives SIGTERM (e.g., from Docker/Kubernetes), it should:

1. Stop accepting new requests
2. Wait for in-flight requests to complete
3. Release resources (database connections, etc.)
4. Exit

Go's `os/signal` package provides `signal.NotifyContext` which creates a context that cancels when a signal is received.

### Configuration

Environment variables are the standard for Go services. No framework needed — `os.Getenv` and simple defaults suffice at this point.

## 4. Go Concepts

| Concept | Usage |
|---------|-------|
| `package main` | Entry point binary |
| `import` | Dependencies between packages |
| `const` / `var` | Configuration defaults |
| `struct` | Group related data (Config, response types) |
| `func` | Handlers and constructors |
| `method` | `Config.Addr()`, `Server.Start()` |
| `pointer receiver` | `func (s *Server) Start()` — mutates receiver |
| `value receiver` | `func (c Config) Addr()` — immutable |
| `error` | Return as second value from handlers |
| `defer` | Ensure cleanup (stop, cancel, close) |
| `context` | Propagation of cancellation and deadlines |

### Pointer vs Value Receivers

- Use **pointer receiver** when the method needs to modify the receiver or when the receiver is large
- Use **value receiver** for small, immutable types
- Be consistent: if one method needs a pointer receiver, use pointer for all

## 5. Architecture Concepts

### Composition Root

The `main` function is the only place where concrete types are wired together. Packages export constructors (`New`) and types, but never know about each other's dependencies.

### Dependency Injection

Dependencies are passed explicitly via function arguments (constructor injection). No global variables, no `init()` functions for wiring.

### Project Structure

```
cmd/server/main.go      — Composition root
internal/config/        — Configuration (env vars)
internal/server/        — HTTP server, routes, handlers
internal/version/       — Build information
```

`internal/` ensures these packages can never be imported by external modules.

## 6. Implementation

### Files Created

- `internal/config/config.go` — Configuration struct and loader
- `internal/server/server.go` — Server creation, route registration, graceful shutdown
- Updated `cmd/server/main.go` — Composition root

### Key Decisions

1. **Fuego for HTTP** — Provides OpenAPI generation, type-safe handlers, minimal abstraction over `net/http`. We could use raw `net/http`, but Fuego gives OpenAPI for free — a justified dependency.

2. **Config via env vars** — Simple, explicit, no framework. `os.Getenv` with defaults.

3. **Graceful shutdown in `Start()`** — Method encapsulates the signal handling and shutdown logic, keeping `main` minimal.

## 7. Tests

### Config Tests

- Default values when env vars not set
- Override from env vars
- Addr() format

### Server Tests

- Health endpoint returns 200 with `{"status":"ok"}`
- Version endpoint returns 200 with version string
- Unknown route returns 404

Tests use `httptest.NewRecorder` and call `s.Mux.ServeHTTP` directly (no real listener needed).

## 8. Refactoring

No refactoring needed at this point — the code is minimal and follows idiomatic Go.

## 9. Code Review

### What went well

- Explicit dependencies through constructor injection
- Small, cohesive packages (`config`, `server`, `version`)
- Graceful shutdown pattern matches production Go services
- Fuego provides OpenAPI generation with zero additional configuration - automatically inside `s.Run()`
- OpenAPI endpoints auto registered by `s.Run()` from `*fuego.Server`:
    - `s.SpecHandler()` registers the OpenAPI spec endpoint
    - `s.UIHandler()` registers the Swagger UI

### Future improvements

- **DONE:** Allow setting environment variables in `.env` file (git-ignored but template `.env.template` commited)
- **DONE:** Add readiness probe endpoint (for Kubernetes)
- Consider structured error responses
- Add request ID middleware

## 10. Exercises

1. Add a `/ready` endpoint that returns 200 (readiness probe)
2. Run the server, send requests with `curl`, then send SIGTERM and observe graceful shutdown
3. Read the Fuego source for `RunContext` to understand how it handles shutdown

## 11. Interview Questions

**Q: Why `net/http` instead of a framework?**
A: We use Fuego, which wraps `net/http` with minimal abstraction. The standard library's `net/http` is production-quality; Fuego adds type-safe OpenAPI generation on top.

**Q: What is a receiver?**
A: A receiver attaches a function to a type, making it a method. `func (s *Server) Start()` is a method on `*Server`.

**Q: Pointer vs value receiver?**
A: Use pointer when the method modifies the receiver or the type is large. Use value for small, immutable types. Be consistent.

**Q: Why `defer`?**
A: `defer` guarantees cleanup runs even if the function panics. Used for closing resources, unlocking mutexes, cancelling contexts.

**Q: What is graceful shutdown?**
A: When a signal is received, the server stops accepting new connections and waits for in-flight requests to complete before exiting. This prevents dropped requests during deployments.

## 12. Roadmap Update

### Go Concepts Learned
- [x] package main
- [x] imports
- [x] variables
- [x] constants
- [x] functions
- [x] structs
- [x] methods
- [x] receivers (pointer and value)
- [x] errors
- [x] defer

### Go Concepts Remaining
- Modules, Packages, Visibility (L1)
- Interfaces, Error Wrapping, Context (L2)
- Pointers (review)
