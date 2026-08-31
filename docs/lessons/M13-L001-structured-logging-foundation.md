### Lesson 13.1 Scope

Introduce a small structured logging foundation with `log/slog` and centralize logger setup before adding request logging, metrics or tracing.

#### Business Context

Operators need machine-readable logs when diagnosing startup failures, migration failures and server lifecycle events. Plain text messages are hard to filter and correlate once the service runs in containers or log aggregation systems.

#### Problem

The application already used `slog` directly in command entry points and server lifecycle code, but logger setup was duplicated and not configurable. The server package also depended on the global default logger for lifecycle messages, which made tests and future request-scoped logging less explicit.

#### Design Discussion

Logging remains a platform concern, not a business slice. `internal/platform/logging` now constructs a configured `*slog.Logger` using standard-library handlers. The application supports JSON logs by default because structured JSON is the common production shape for containerized services. Text logs remain available for local debugging.

The lesson intentionally does not add request logging, correlation IDs, Prometheus metrics or OpenTelemetry. Those belong to later observability lessons so each concept remains clear.

#### Go Concepts

- `log/slog` structured logging
- `slog.Logger` as an explicit dependency
- `io.Writer` for testable log output
- configuration parsing for log level and format
- avoiding package-level globals where explicit injection is cheap

#### Architecture Concepts

- platform logging package below business slices
- command entry points configure process-wide defaults
- server lifecycle logging receives an explicit logger
- observability foundation before request correlation and metrics

### Lesson 13.1 Completion Notes

#### Business Context

MES Lite now has a configurable structured logging foundation for startup, migration and server lifecycle diagnostics.

#### Problem

Logger setup was duplicated in `cmd/server` and `cmd/migrate`, and the server package logged through the global `slog` default instead of an explicit dependency.

#### Design Discussion

Added `internal/platform/logging` with one `New` function that builds a `*slog.Logger` from an `io.Writer`, log level and format. JSON is the default because it is easy for production log collectors to parse. Text format remains available through `LOG_FORMAT=text` for local development.

`cmd/server` and `cmd/migrate` now both configure the logger through the shared package and set it as the process default. `internal/server` receives an injected logger for lifecycle logs and uses a discard logger by default so tests stay quiet.

#### Implementation

- Added `internal/platform/logging`.
- Added configurable JSON and text `slog` handlers.
- Added log-level parsing for `debug`, `info`, `warn`, `warning` and `error`.
- Added `LOG_LEVEL` and `LOG_FORMAT` configuration.
- Defaulted logs to `LOG_LEVEL=info` and `LOG_FORMAT=json`.
- Updated `cmd/server` and `cmd/migrate` to use the shared logging package.
- Added `Server.SetLogger` and changed server lifecycle logs to use the injected logger.

#### Tests

- Added logging tests for JSON structured output.
- Added logging tests for text output.
- Added logging tests for level parsing and invalid configuration.
- Updated configuration tests for logging defaults and environment values.
- Verified with `go test ./internal/platform/config ./internal/platform/logging ./internal/server -count=1`.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `go vet ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

Duplicated logger construction was removed from command entry points. Server lifecycle logging no longer relies directly on package-global logger state.

#### Code Review

An experienced Go engineer would approve this as a focused first observability step. It uses the standard library, keeps logging as platform plumbing and avoids introducing an external logging framework.

The main follow-up is request logging with correlation IDs. Startup logs are structured now, but request-level diagnostics still need a middleware boundary.

#### Exercises

- Run the server with `LOG_FORMAT=text` and compare the output with JSON logs.
- Add a test proving `LOG_LEVEL=error` suppresses info logs.
- Decide which startup values are safe to log and which could leak secrets.

#### Interview Questions

- Why are structured logs more useful than plain formatted strings in production?
- Why is `log/slog` preferable here to adding a third-party logger?
- Why should secrets such as `JWT_SECRET` never be logged?
- When is using `slog.Default()` acceptable, and when is explicit logger injection better?

#### Roadmap Update

- Lesson 13.1 completed.
- Current lesson moved to Lesson 13.2.
- Standard Library `log/slog` marked complete in the Knowledge Matrix.
