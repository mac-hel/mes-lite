### Lesson 11.2 Scope

Add a worker pool that consumes queued jobs and executes registered handlers concurrently.

#### Business Context

The application can now accept background jobs, but accepted work never runs. A manager could enqueue a future CSV import job, yet no worker would pick it up, mark it running or record whether it succeeded.

#### Problem

The queue is only the producer/consumer handoff. It stores job state and hands out queued jobs, but it does not own execution. Running handlers inside the queue would mix scheduling, state tracking and workload behavior in one type.

#### Design Discussion

Lesson 11.2 introduces a concrete `WorkerPool` in `internal/platform/jobs`. The pool starts a fixed number of goroutines, each repeatedly dequeuing jobs and executing the handler registered for the job type.

The queue remains the owner of job status. Workers ask the queue to move a job from queued to running, then to succeeded or failed. This keeps status updates serialized behind the queue mutex instead of letting workers mutate copies.

Handlers are registered by job type. Unknown job types fail the job instead of panicking, because malformed or unsupported work should be visible as job state.

Cancellation is limited to worker lifecycle in this lesson. User-requested job cancellation and progress updates belong to Lesson 11.3.

#### Go Concepts

- goroutines for concurrent background execution
- `sync.WaitGroup` for waiting on worker shutdown
- `sync.Once` for idempotent start and stop behavior
- context cancellation for worker lifecycle
- error values recorded as failed job state

#### Architecture Concepts

- worker pool as execution infrastructure
- queue owns state; handlers own business work
- fixed concurrency limit instead of unbounded goroutine creation
- workload registry without making the queue depend on business slices

#### Implementation Plan

- Add queue status-transition methods for running, succeeded and failed states.
- Add a concrete `WorkerPool` with fixed worker count.
- Add job handler registration by `jobs.Type`.
- Execute each job exactly once after successful dequeue.
- Record handler success or failure back into the queue.
- Keep progress and external cancellation out of scope until Lesson 11.3.

#### Tests

- Status-transition tests for valid and invalid lifecycle moves.
- Worker execution success and failure tests.
- Unknown job type becomes a failed job.
- Multiple workers process multiple jobs without duplicate execution.
- Shutdown waits for accepted running work and releases workers.
- Race detector remains clean.

#### Exercises

- Add a second job type and register a second handler.
- Change the pool to launch one goroutine per job and explain why that removes backpressure.
- Add a test proving `Stop` is safe to call more than once.
- Explain why handlers receive a job copy instead of a pointer to queue state.

#### Interview Questions

- What problem does a worker pool solve?
- Why use `sync.WaitGroup` instead of sleeping until workers finish?
- Why should queues not execute business handlers directly?
- What happens if a goroutine writes shared state without synchronization?
- How does context cancellation differ from closing a work channel?

### Lesson 11.2 Completion Notes

#### Business Context

MES Lite can now execute accepted background jobs with bounded concurrency.

#### Problem

Jobs could be enqueued and dequeued, but nothing marked them running, executed work or recorded success/failure. The system had a queue, not background processing.

#### Design Discussion

Added a concrete `WorkerPool` that consumes the existing in-memory queue. The pool owns goroutine lifecycle and handler dispatch, while the queue remains the owner of job state.

This split matters: workers receive job copies and ask the queue to perform status transitions. They do not mutate shared job structs directly, which keeps state changes synchronized behind the queue mutex.

The pool uses a fixed worker count. That makes concurrency a deliberate capacity choice instead of launching one goroutine per job and hoping the runtime absorbs the load.

Unknown job types are recorded as failed jobs rather than panicking. Unsupported work is operational data, not a process-level crash.

#### Implementation

- Added `ErrInvalidStatusTransition`.
- Added `Queue.MarkRunning`, `Queue.MarkSucceeded` and `Queue.MarkFailed`.
- Added `jobs.Handler` for registered background workload functions.
- Added `WorkerPool` with fixed worker count, `Start` and `Stop`.
- Used `sync.WaitGroup` to wait for worker goroutines.
- Used `sync.Once` so starting and stopping are idempotent operations.
- Recorded handler errors into failed job state.
- Kept progress updates, external cancellation and HTTP status routes out of scope for Lesson 11.3.

#### Tests

- Added queue status-transition tests for valid and invalid moves.
- Added worker-pool configuration validation tests.
- Tested successful handler execution and lifecycle timestamp updates.
- Tested handler failure becomes failed job state.
- Tested unknown job types fail jobs instead of panicking.
- Tested multiple workers execute 100 jobs exactly once.
- Tested `Stop` waits for running work.
- Tested `Stop` is safe to call more than once.
- Verified with `go test ./internal/platform/jobs -count=1`.
- Verified with `go test ./internal/platform/jobs -race -count=3 -shuffle=on`.
- Verified with `go build ./...`.
- Verified with `go test ./... -count=1`.
- Verified with `go vet ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

The queue was extended with explicit lifecycle methods instead of exposing a generic update hook. This keeps status transitions narrow and reviewable.

No job-store interface or durable implementation was introduced. ADR 0004 already records durability as a future trigger, and adding persistence now would distract from goroutines, `WaitGroup` and worker lifecycle.

#### Code Review

An experienced Go engineer would approve the core shape: fixed worker count, queue-owned synchronized state, no unbounded goroutine creation, copied jobs passed to handlers and race-detector coverage.

The main caveat is expected lesson scope. Workers can record success/failure, but clients still cannot query job status over HTTP, observe progress or request cancellation. That is the next lesson.

#### Exercises

- Add a second job type and prove it dispatches to a different handler.
- Add a test where `Stop` times out while a handler is blocked, then explain what happens to the running goroutine.
- Change workers to launch a goroutine per dequeued job and measure how many jobs can run at once.
- Explain why status transitions should stay queue-owned instead of handler-owned.

#### Interview Questions

- What does a worker pool protect a service from?
- Why is `sync.WaitGroup` the right primitive for waiting on worker goroutines?
- Why is a fixed worker count different from a buffered queue capacity?
- Why should handler failures become job state instead of only logs?
- What data race would appear if handlers received pointers to tracked jobs?

#### Roadmap Update

- Lesson 11.2 completed.
- Current lesson moved to Lesson 11.3.
- Concurrency `Goroutines`, `WaitGroup` and `Worker Pools` marked complete in the Knowledge Matrix.
- Testing `Race Detection` marked complete in the Knowledge Matrix.
- Known technical debt updated for missing job progress, cancellation and status API.
