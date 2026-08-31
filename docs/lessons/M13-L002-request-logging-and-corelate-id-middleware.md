### Lesson 13.2 Scope

Add request logging middleware with correlation IDs so every HTTP request receives a stable identifier and emits one structured completion log.

#### Business Context

Operators need to diagnose API behavior across many concurrent requests. A startup log can prove the process booted, but request logs are what connect a user's API call to status code, route path and latency.

#### Problem

The application had structured startup and lifecycle logs, but individual HTTP requests were not logged. There was also no request-scoped correlation ID that could be returned to clients and propagated through future logs, metrics or traces.

#### Design Discussion

Request logging is implemented as standard `net/http` middleware in `internal/platform/logging`. It preserves an incoming `X-Request-ID` header when present, generates a UUID-shaped ID when missing, stores it in request context and returns it on the response header.

The middleware logs one completion record after the handler runs. It records method, path, status and duration in milliseconds. Query strings are intentionally not logged because they can contain sensitive values and usually create high-cardinality logs.

Fuego route registration uses `fuego.Use` rather than only global server middleware. This keeps request logging active in both runtime and existing tests that call `s.Mux.ServeHTTP` directly.

#### Go Concepts

- standard-library HTTP middleware
- response-writer wrapping for status capture
- context values for request-scoped correlation IDs
- `time.Since` for request duration
- preserving inbound headers versus generating defaults

#### Architecture Concepts

- request logging as platform middleware
- correlation ID as an observability boundary
- explicit server construction with logger injection
- avoiding metrics and traces until their dedicated lessons

### Lesson 13.2 Completion Notes

#### Business Context

MES Lite now emits one structured log record per HTTP request and returns a correlation ID to API clients.

#### Problem

Without request logs, operators could see startup failures but not normal API traffic, response statuses or request latency. Without a correlation ID, future diagnostics could not tie logs, client reports and traces together.

#### Design Discussion

Added `logging.RequestLogger` as reusable platform middleware. It accepts a `*slog.Logger`, wraps the `http.ResponseWriter` to capture status code and stores a request ID in context through an unexported key type.

The request ID comes from `X-Request-ID` when supplied by a client or gateway. If it is missing, the middleware generates one with the existing platform ID generator. The same ID is returned in the response header.

`server.NewWithLogger` registers request logging before routes are added. Existing tests can still call `server.New`, which uses a discard logger by default.

#### Implementation

- Added `logging.RequestLogger` middleware.
- Added request ID context helpers in `internal/platform/logging`.
- Added `X-Request-ID` response header propagation.
- Added status-code capture through a response-writer wrapper.
- Logged structured request fields: `request_id`, `method`, `path`, `status` and `duration_ms`.
- Added `server.NewWithLogger` for runtime logger injection before route registration.
- Updated `cmd/server` to build the HTTP server with the configured logger.

#### Tests

- Added middleware test preserving an incoming request ID.
- Added middleware test generating a request ID when missing.
- Added context propagation assertions for request IDs.
- Added structured log assertions for method, path, status and request ID.
- Added server route test proving Fuego-registered routes emit request logs and return `X-Request-ID`.
- Verified with `go test ./internal/platform/logging ./internal/server -count=1`.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `go vet ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

The previous `Server.SetLogger` hook was replaced by `NewWithLogger` because request middleware must be registered with the correct logger before routes are installed. The old `New` constructor remains as a quiet default for tests.

#### Code Review

An experienced Go engineer would approve the middleware shape because it uses `net/http`, avoids framework-specific logging code and keeps correlation IDs request-scoped rather than global.

The main caveat is that only status, path and duration are logged today. User identity, route templates, panic recovery and trace IDs are intentionally postponed until later observability work defines those boundaries.

#### Exercises

- Add a test proving missing routes also receive `X-Request-ID` once the server is run through the same handler path used in production.
- Add authenticated user ID to request logs and discuss the privacy/security trade-off.
- Decide whether `duration_ms` should be integer milliseconds or a string duration.

#### Interview Questions

- Why are correlation IDs useful in distributed systems?
- Why should request IDs be stored in context instead of a global variable?
- How does wrapping `http.ResponseWriter` let middleware observe status codes?
- Why might logging query strings be risky?

#### Roadmap Update

- Lesson 13.2 completed.
- Current lesson moved to Lesson 13.3.
