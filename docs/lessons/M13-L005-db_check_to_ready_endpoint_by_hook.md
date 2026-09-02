### Lesson 13.5 Scope

Review health, readiness and the complete observability stack before closing Milestone 13.

#### Business Context

Operators need two different operational signals: whether the process is alive and whether it is ready to serve real business traffic. They also need confidence that logs, metrics and traces compose without leaking into business packages.

#### Problem

`/health` existed as a liveness endpoint, but `/ready` always returned success. That meant an orchestrator could send traffic to an application whose database dependency was unavailable.

#### Design Discussion

Liveness stays cheap and dependency-free: `/health` answers only whether the process can respond. Readiness now has an explicit check hook. The production composition root wires that hook to `pgxpool.Ping`, while tests can install a small fake function.

A readiness failure returns `503 Service Unavailable`, not `500 Internal Server Error`. The process is still alive; it is just not ready to receive traffic.

#### Go Concepts

- function fields for small dependency hooks
- `context.WithTimeout` for bounded readiness checks
- correct HTTP status semantics for operational endpoints
- focused tests for success and dependency failure

#### Architecture Concepts

- liveness separated from readiness
- composition root wires infrastructure health checks
- observability remains in platform packages
- milestone review before performance work begins

### Lesson 13.5 Completion Notes

#### Business Context

MES Lite now has a meaningful readiness endpoint in addition to the existing liveness endpoint.

#### Problem

The readiness endpoint returned success unconditionally. That was misleading because the server could be alive while PostgreSQL was unreachable.

#### Design Discussion

Added a readiness-check hook to `server.Server`. The server package does not import PostgreSQL or know which dependencies matter; it only calls the injected function with a bounded context. `cmd/server` wires the hook to `db.Ping`.

This keeps dependency knowledge in the composition root and preserves the server package as HTTP composition rather than database infrastructure.

#### Implementation

- Added `Server.SetReadinessCheck`.
- Changed `/ready` to return `{"status":"ready"}` when the check succeeds.
- Added a two-second timeout around readiness checks.
- Changed readiness dependency failures to return `503 Service Unavailable`.
- Logged readiness failures as structured warnings.
- Wired runtime readiness to `pgxpool.Pool.Ping` in `cmd/server`.

#### Tests

- Added server test proving `/ready` calls the readiness check and returns `200 OK` with `status=ready`.
- Added server test proving readiness-check failure returns `503 Service Unavailable`.
- Verified with `go test ./internal/server -count=1`.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `go vet ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

`NewWithLogger` now constructs `*Server` before route registration so the `/ready` route can call a method that uses server state. The public route behavior is otherwise unchanged.

#### Code Review

An experienced Go engineer would approve the liveness/readiness split. The readiness endpoint checks the database without making `internal/server` depend on pgx, and failures use the correct operational status code.

The main observability limitation remains OTLP exporter wiring. The application can create request spans and emit stdout traces for learning, but production collector export is still future work.

#### Exercises

- Add a readiness check that verifies the background worker pool has started.
- Add a test proving readiness failures are logged with `request_id` and `trace_id` when tracing is active.
- Decide whether `/health`, `/ready` and `/metrics` should be protected by network policy in deployment.

#### Interview Questions

- What is the difference between liveness and readiness?
- Why should readiness checks have timeouts?
- Why is `503 Service Unavailable` better than `500` for dependency-readiness failure?
- Why should observability middleware live outside business slices?

#### Roadmap Update

- Lesson 13.5 completed.
- Milestone 13 completed.
- Current milestone moved to Milestone 14.
- Current lesson moved to Lesson 14.1.
- Architecture `Observability` marked complete in the Knowledge Matrix.
