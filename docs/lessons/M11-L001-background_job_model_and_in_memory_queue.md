### Lesson 11.1 Scope

Introduce the background job vertical slice with a job model and an in-memory, channel-backed queue. No goroutines run work yet.

#### Business Context

CSV imports of historical production data run inside the HTTP request. A large file holds the request open until every row is validated and persisted, and the client learns nothing until the whole file is finished. Managers need to hand work to the server, get an answer immediately and check on the result later.

#### Problem

Asynchronous processing needs somewhere to put accepted work before anything can process it. Starting with goroutines would mean designing the job model, the handoff and the worker lifecycle at the same time, which is three hard problems in one step.

#### Design Discussion

Lesson 11.1 builds only the data that concurrency will move around: a `jobs.Job` value type and a `jobs.Queue` that accepts and hands out work safely. Lesson 11.2 adds the workers that consume it.

The queue owns two pieces of state with different responsibilities. A map guarded by `sync.RWMutex` holds the canonical status of every tracked job. A buffered channel is the handoff between producers and consumers. The channel carries job IDs rather than job values, so a consumer always reads current state from the map instead of acting on a copy taken at enqueue time.

The queue never closes its handoff channel. There are many producers, and closing a channel other goroutines may still send on is a panic rather than a shutdown signal. Shutdown closes a dedicated `done` channel with exactly one closer.

`Enqueue` never blocks. A full queue returns `ErrQueueFull` immediately, so an HTTP caller gets backpressure instead of waiting on a channel that may stay full for minutes. The buffer capacity is the burst the application is willing to absorb.

Durability is deferred and documented in ADR 0004: jobs are lost on restart.

#### Go Concepts

- buffered channels as a bounded handoff between producers and consumers
- channel ownership: who is allowed to close a channel, and why a `done` channel is used instead
- `select` with `default` for non-blocking send and receive
- `select` picking randomly among ready cases, and what that means for draining
- `sync.RWMutex` guarding shared map state read by many goroutines
- defensive copies of slices shared across goroutines

#### Architecture Concepts

- background job queue as a producer/consumer boundary
- job status owned by one component instead of spread across workers
- in-memory infrastructure chosen deliberately, with the durability trade-off recorded in an ADR

### Lesson 11.1 Completion Notes

#### Business Context

MES Lite now has the first building block for moving long-running work out of HTTP requests.

#### Problem

The application had no representation of deferred work and no safe place to hold it between the request that accepts it and the code that will run it.

#### Design Discussion

Added `internal/jobs` with a `Job` value type and a `Queue`. `Job` carries an ID, a type, a lifecycle status, an opaque payload and lifecycle timestamps. The payload is `[]byte` so the queue stays independent of any single workload; generics remain postponed.

Status transitions are intentionally absent. A job is created queued, and nothing moves it yet, because moving a job to running belongs with the worker that does the moving in Lesson 11.2.

The queue hands out copies. `Find` and `Dequeue` return `job.clone()`, and `NewJob` copies the payload it is given, so no caller can mutate tracked state through a shared slice.

Two ordering details required care:

- `Enqueue` holds the write lock across the channel send. The send is non-blocking, so it cannot deadlock, and holding the lock means a consumer that receives the ID cannot look it up before the map entry exists.
- `Dequeue` prefers waiting work over the close signal. It tries a non-blocking receive first, and when the `done` case fires it looks one more time, because `select` picks randomly among ready cases. Since `Enqueue` checks `closed` under the same mutex that `Close` writes it under, no send can follow `close(done)`, which makes that final look exact.

#### Implementation

- Added `internal/jobs/job.go` with `Job`, `Status`, `Type`, `NewJob`, `NewJobID` and `Validate`.
- Added `Status.Valid` and `Status.Terminal` for the lifecycle states queued, running, succeeded, failed and cancelled.
- Added `TypeProductionEntryImport` as the first planned workload.
- Added `internal/jobs/queue.go` with `NewQueue`, `Enqueue`, `Dequeue`, `Find`, `Len`, `Capacity` and `Close`.
- Added `ErrInvalidJob`, `ErrQueueFull`, `ErrQueueClosed`, `ErrNotFound` and `ErrAlreadyExists`.
- Added defensive payload copies on creation and on every read.
- Added ADR `0004-introduce-background-jobs.md`.

#### Tests

- Tested job normalization, payload copying and table-driven validation.
- Tested status validity and terminal states.
- Tested FIFO enqueue/dequeue, duplicate IDs and rejection of non-queued jobs.
- Tested that a full queue returns `ErrQueueFull` without blocking and leaves no tracked job behind.
- Tested that `Dequeue` blocks on an empty queue and returns once work arrives.
- Tested context cancellation while waiting.
- Tested that `Close` drains accepted work before reporting closed, rejects new work, releases waiting consumers and is idempotent.
- Tested that `Find` returns a copy.
- Added a concurrency test with 8 producers, 200 jobs and 4 consumers asserting exactly-once delivery.
- Verified with `go build ./...`.
- Verified with `go vet ./internal/jobs`.
- Verified with `go test ./internal/jobs -count=5 -race -shuffle=on`.
- Verified with `go test ./... -count=1`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

No existing package changed. The CSV import slice keeps its synchronous path until Lesson 11.4 gives it an asynchronous one, so the current endpoint keeps working while the job machinery is built beside it.

No `Store` interface was introduced for jobs. Nothing consumes the queue yet, and interfaces belong to consumers.

#### Code Review

An experienced Go engineer would approve the scope and the channel-ownership decisions: the handoff channel is never closed by a producer, shutdown has one closer, shared state is copied on the way in and on the way out, and the race detector is clean across shuffled repeat runs.

Two caveats are deliberate. `NewJobID` duplicates the UUID-shaped generator already present in the production and orders slices; a shared helper is worth considering, but a fourth copy is a weak reason to create a cross-slice utility package. Jobs are in-memory only, which ADR 0004 records rather than hides.

#### Exercises

- Remove the final `tryReceive` inside the `done` case of `Dequeue` and write a test that shows a job can be stranded at shutdown.
- Change `Enqueue` to block with `select` on `ctx.Done()` instead of returning `ErrQueueFull`, and describe what an HTTP client would experience under load.
- Explain why holding a mutex across the channel send in `Enqueue` is safe here, and what would make it a deadlock.
- Replace the ID channel with a `chan Job` and identify what becomes stale once workers start updating status.

#### Interview Questions

- Who should close a channel, and why is closing from the receiver side wrong?
- What does a buffered channel give you that an unbuffered one does not?
- What happens when several cases in a `select` are ready at the same time?
- When is a mutex a better fit than a channel for shared state?
- Why does a queue that hands out copies avoid a whole class of data races?

#### Roadmap Update

- Lesson 11.1 completed.
- Current lesson moved to Lesson 11.2.
- Concurrency `Channels`, `select` and `Race Detector` marked complete in the Knowledge Matrix.
- Standard Library `sync` marked complete in the Knowledge Matrix.
- Known technical debt updated for in-memory job durability and duplicated ID generation.
