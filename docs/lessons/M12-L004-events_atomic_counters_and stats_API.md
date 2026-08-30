### Lesson 12.4 Scope

Add atomic machine intake counters and review the machine integration slice under the race detector.

#### Business Context

Operators need lightweight visibility into fake machine event intake while the project prepares for real observability in Milestone 13. Even before Prometheus or tracing exists, the machine slice can expose basic counters that are safe under concurrent traffic.

#### Problem

The machine service could accept, deduplicate and reject events, but it did not count those outcomes. Adding counters with a shared mutex would work, but it would couple independent numeric telemetry to the store's event-state lock.

#### Design Discussion

This lesson uses typed `sync/atomic` counters inside `machines.Service`. Each counter is independent and monotonically increasing, which makes atomics a better fit than a broader mutex-protected stats struct.

The counters track received attempts, accepted new events, duplicate retries, conflicts and invalid events. A small protected stats endpoint returns a point-in-time snapshot.

This is not the full observability milestone. There are no Prometheus metrics, logs or traces yet; those belong to Milestone 13.

#### Go Concepts

- typed atomic counters with `atomic.Uint64`
- atomic `Add` for concurrent increments
- atomic `Load` for snapshot reads
- choosing atomics only for independent scalar state
- race detector review for mixed mutex and atomic synchronization

#### Architecture Concepts

- operational counters as service-owned state
- HTTP stats endpoint as temporary visibility before formal observability
- atomics for telemetry, mutexes for compound event state
- avoiding global metrics registries before Milestone 13

### Lesson 12.4 Completion Notes

#### Business Context

MES Lite now exposes basic fake machine intake counters so concurrent event activity can be inspected without reading internal state directly.

#### Problem

Machine intake outcomes were invisible. Tests could prove behavior, but an operator or client could not see how many events were accepted, retried, conflicted or rejected as invalid.

#### Design Discussion

Added counters to `machines.Service` instead of to the in-memory store. The service owns the workflow decisions that classify an attempt as accepted, duplicate retry, conflict or invalid.

The event store still uses mutexes because it protects compound state: a slice and a map that must stay consistent. The service counters use atomics because each number can be updated independently without a larger invariant between fields.

The stats endpoint is protected with the same admin/manager RBAC as fake machine intake. Workers should not inspect integration operational counters.

#### Implementation

- Added `IntakeStats` snapshot type.
- Added typed atomic counters to `machines.Service`.
- Counted received attempts, accepted events, duplicate retries, conflicts and invalid events.
- Added `Service.Stats` using atomic loads.
- Added `IntakeStatsResponse` DTO.
- Added `GET /machines/events/stats`.
- Registered stats route with admin/manager RBAC.

#### Tests

- Added service tests for accepted, duplicate retry, conflict and invalid counters.
- Extended concurrent retry test to verify atomic counters after 50 concurrent callers.
- Added handler test for stats response mapping.
- Added server route tests proving managers can read machine stats and workers are forbidden.
- Verified with `go test ./internal/machines ./internal/server -count=1`.
- Verified with `go test ./internal/machines -race -count=5 -shuffle=on`.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `go vet ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

No global metrics package was introduced. Counters remain local to the machine service until Milestone 13 introduces formal observability.

#### Code Review

An experienced Go engineer would approve atomics here because the counters are independent scalar values. They would not be appropriate for the store's event slice and deduplication map, where multiple fields must change together.

The main caveat is that these counters are in-memory and reset on restart. That is acceptable for the learning goal and current fake integration scope.

#### Exercises

- Add a test proving stats reset when a new service is constructed.
- Add a benchmark comparing atomic counter increments with mutex-protected increments.
- Decide which machine counters should become Prometheus metrics in Milestone 13.

#### Interview Questions

- When are atomic operations preferable to a mutex?
- Why are atomics risky for compound state?
- What is the difference between a race-free counter and a consistent multi-field snapshot?
- Why does a clean race detector run not prove a design is logically correct?

#### Roadmap Update

- Lesson 12.4 completed.
- Current lesson moved to Lesson 12.5.
- Standard Library `sync/atomic` and Concurrency `Atomic` marked complete in the Knowledge Matrix.
