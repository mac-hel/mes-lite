### Lesson 11.3 Scope

Expose background job status over HTTP, add progress reporting and allow authorized users to request cancellation.

#### Business Context

Managers need visibility into accepted long-running work. Once imports move to background jobs, returning only a job ID is not enough; clients must be able to check whether work is queued, running, finished, failed or cancelled.

#### Problem

Workers can execute jobs, but the outside world cannot observe job state or ask for a job to stop. A running handler also has no safe way to report progress back to the queue.

#### Design Discussion

The queue remains the owner of canonical job state. It gains progress and cancellation state transitions, while the worker pool owns running job contexts. Cancelling a queued job marks it cancelled immediately. Cancelling a running job records the request and cancels that job's context so a cooperative handler can stop.

The HTTP API is intentionally small:

- `GET /jobs/{id}` returns one job status snapshot.
- `PUT /jobs/{id}/cancel` requests cancellation and returns the updated snapshot.

Only admins and managers can read and cancel jobs for now. Workers should not inspect system-wide background processing, and leaders do not currently own import/system-operation workflows.

#### Go Concepts

- cooperative cancellation with `context.Context`
- per-job cancel functions guarded by a mutex
- progress updates as synchronized state changes
- HTTP status mapping for asynchronous job state
- idempotent cancellation for terminal jobs

#### Architecture Concepts

- status API as an operational read boundary
- queue owns state; worker pool owns execution contexts
- cancellation request separated from immediate termination
- no durable job persistence until ADR 0004's trigger is reached

#### Implementation Plan

- Add progress and cancellation fields to `jobs.Job`.
- Add queue methods for progress reporting and cancellation state.
- Extend `WorkerPool` with per-job cancellation.
- Add a jobs HTTP handler for status and cancellation.
- Register protected job routes in the server.
- Wire an in-memory queue and worker pool in the production composition root.

#### Tests

- Queue progress validation and defensive-copy tests.
- Queued cancellation transitions directly to cancelled.
- Running cancellation marks a request and worker context cancellation marks final cancelled state.
- Job status API returns stable JSON DTOs.
- Job cancellation API maps not found and terminal states correctly.
- Server authorization tests for admin/manager access and worker rejection.
- Race detector remains clean.

#### Exercises

- Add a test handler that reports progress from 0 to 100 and is cancelled halfway.
- Decide whether leaders should be allowed to read job status once report generation jobs exist.
- Add a list endpoint for recent jobs and discuss pagination before implementing it.
- Explain why cancellation in Go is cooperative rather than forcibly killing a goroutine.

#### Interview Questions

- Why does Go use context cancellation instead of killing goroutines externally?
- What must a handler do for cancellation to be effective?
- Why should progress updates be synchronized by the queue?
- What is the difference between cancellation requested and cancelled?
- Why can terminal job cancellation be idempotent?

### Lesson 11.3 Completion Notes

#### Business Context

MES Lite now has visibility and control for individual background jobs. A client can read a job status snapshot and authorized users can request cancellation.

#### Problem

The worker pool could execute jobs, but clients had no API for observing job state, and running handlers had no synchronized way to report progress or respond to a user cancellation request.

#### Design Discussion

The queue remains the canonical state owner. It now records progress, cancellation requests and final cancelled state under the same mutex that protects job status.

The worker pool owns execution contexts. Each running job gets a child context with its own cancel function. `WorkerPool.Cancel` records the cancellation request in the queue and, if the job is running, calls that job's cancel function.

Cancellation is cooperative. The pool does not kill goroutines. A handler must observe `ctx.Done()` and return. This is idiomatic Go because forcibly stopping goroutines would risk leaving locks, transactions or files in inconsistent states.

The status API is deliberately small: `GET /jobs/{id}` and `PUT /jobs/{id}/cancel`. A list endpoint is postponed until there is real UI/API pressure to define pagination, filtering and retention rules.

#### Implementation

- Added `Progress` and `CancelRequested` to `jobs.Job`.
- Added progress validation to `Job.Validate`.
- Added `ErrInvalidProgress`.
- Added `Queue.ReportProgress`.
- Added `Queue.RequestCancellation`.
- Added `Queue.MarkCancelled`.
- Added `Queue.CancellationRequested`.
- Extended `WorkerPool` with a mutex-protected map of running job cancel functions.
- Added `WorkerPool.Cancel`.
- Changed worker execution to pass per-job contexts to handlers.
- Added `jobs.HTTPHandler` with `GET /jobs/{id}` and `PUT /jobs/{id}/cancel` behavior.
- Added `server.RegisterJobRoutes` protected by admin/manager RBAC.
- Wired an in-memory queue, worker pool and job handler in `cmd/server`.

#### Tests

- Added queue progress tests.
- Added invalid progress and invalid-state progress tests.
- Added queued and running cancellation queue tests.
- Added worker-pool queued cancellation test.
- Added worker-pool running cancellation test proving handler context cancellation becomes final cancelled state.
- Added job HTTP handler tests for status, not found and cancellation.
- Added server route tests proving managers can read job status, admins can cancel jobs and workers are forbidden.
- Verified with `go test ./internal/platform/jobs ./internal/server -count=1`.
- Verified with `go test ./internal/platform/jobs -race -count=3 -shuffle=on`.
- Verified with `go build ./...`.
- Verified with `go test ./... -count=1`.
- Verified with `go vet ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

`server.New` was not widened with another constructor argument. Instead, `server.RegisterJobRoutes` registers operational job routes when a jobs handler exists. This keeps the existing application route constructor stable while L11.4 decides how background jobs are wired into CSV import.

No durable job persistence, job listing or retention policy was introduced. Those would be separate product/API decisions rather than requirements for cancellation mechanics.

#### Code Review

An experienced Go engineer would approve the cooperative-cancellation model, queue-owned state transitions and race-detector coverage. The main implementation risk in this area is handlers that ignore their context; L11.4 must make the async CSV import handler check cancellation at natural boundaries.

The main API caveat is that only per-job status exists. That is enough once an enqueue endpoint returns a job ID, but a UI will eventually want listing and retention semantics.

#### Exercises

- Modify a test handler to ignore `ctx.Done()` and observe why cancellation request does not immediately stop work.
- Add a `GET /jobs` proposal with `status`, `limit` and `offset` filters, but do not implement it yet.
- Add a progress monotonicity rule and decide whether decreasing progress should be rejected.
- Explain why `WorkerPool.Cancel` records queue state before calling the running cancel function.

#### Interview Questions

- Why is goroutine cancellation cooperative in Go?
- What is the difference between cancelling a context and closing a channel?
- Why should a queue own status/progress updates instead of handlers mutating job pointers?
- What race would appear if the worker pool's running-job map had no mutex?
- Why might a cancelled job still show partial progress?

#### Roadmap Update

- Lesson 11.3 completed.
- Current lesson moved to Lesson 11.4.
- Known technical debt updated: job status/cancellation exists for individual jobs; job listing and async CSV wiring remain pending.
