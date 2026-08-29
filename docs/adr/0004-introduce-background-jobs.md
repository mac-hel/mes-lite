# ADR 0004: Introduce Background Jobs With An In-Memory Queue

## Status

Accepted — amended by ADR 0005. The package described below as
`internal/jobs` now lives at `internal/platform/jobs`, where a depguard rule
keeps it free of business imports.

## Context

Milestone 11 introduces asynchronous processing. The first concrete driver is CSV import: `POST /imports/production-entries` streams, validates and persists every row inside the HTTP request, so a large historical file holds a request open for as long as the import takes. Report generation and future notifications have the same shape.

The application needs a way to accept work, answer immediately and finish the work elsewhere. That requires a job model, a queue and, in later lessons, workers and a status API.

The queue implementation is the decision with long-term consequences. A durable queue (PostgreSQL table, or a broker such as Redis/NATS) survives restarts. An in-memory queue does not.

## Decision

Introduce `internal/jobs` as a package owning the background job model and an in-memory, channel-backed queue. (ADR 0005 later reclassified it as a platform package at `internal/platform/jobs`; it was never a business slice, since it carries opaque payloads and knows nothing about production data.)

`jobs.Queue` keeps two pieces of state:

- a map of tracked jobs guarded by a `sync.RWMutex`, which is the canonical job status,
- a buffered channel of job IDs, which is the handoff between producers and consumers.

The channel carries IDs rather than job values so a consumer always reads current state from the map instead of acting on a copy taken at enqueue time.

The queue never closes its handoff channel. There are many producers, and closing a channel other goroutines may still send on panics. Shutdown closes a dedicated `done` channel that has exactly one closer.

Enqueue never blocks. A full queue returns `ErrQueueFull` so callers get backpressure immediately; the buffer capacity is the burst the application is willing to absorb.

Durability is explicitly deferred. Jobs are lost on restart.

## Alternatives Considered

- **PostgreSQL job table with `SELECT ... FOR UPDATE SKIP LOCKED`**: rejected for now. It is the natural next step for durability, but starting there would teach SQL polling instead of goroutines, channels and `select`, which is the stated purpose of this milestone. The `Queue` API (`Enqueue`, `Dequeue`, `Find`, `Close`) is deliberately shaped so a durable implementation can replace it behind a consumer-owned interface.
- **External broker (Redis, NATS, RabbitMQ)**: rejected. It adds an operational dependency for a single-instance application serving a small company, and it would hide the concurrency concepts behind a client library.
- **Goroutine per request without a queue**: rejected. Unbounded goroutines give no backpressure, no status tracking and no way to cancel or shut down cleanly.
- **Generic `Job[T any]` payload**: rejected. Generics remain postponed until the application shows a real need. An opaque `[]byte` payload keeps the queue independent of any single workload.

## Trade-Offs

- Pro: concurrency concepts (channel ownership, `select`, buffered-channel backpressure, mutex-guarded shared state) are learned directly instead of through a client library.
- Pro: no new infrastructure dependency and no new migration.
- Pro: the queue API is small enough to swap for a durable implementation.
- Con: queued and running jobs are lost on restart or crash.
- Con: the queue is per-process, so the application cannot be scaled to several instances while jobs matter.
- Con: job payloads live in memory, which bounds how large a payload may reasonably be.

## Consequences

- Jobs accepted by the API must be treated as best-effort until durability exists. Any endpoint that enqueues work must say so in its response contract.
- Graceful shutdown must drain or explicitly abandon in-flight jobs. This is addressed in Lesson 11.5.
- If a workload becomes business-critical (an import a manager must not silently lose), that is the trigger to revisit this ADR and move the queue to PostgreSQL.
- Horizontal scaling is out of scope while this ADR stands.
