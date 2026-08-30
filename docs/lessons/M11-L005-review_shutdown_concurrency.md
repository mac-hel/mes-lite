### Lesson 11.5 Scope

Review the background job implementation for race safety, shutdown behavior and milestone readiness.

#### Business Context

Background jobs now execute real CSV import work. Before moving to machine integration, the concurrency foundation should be reviewed as if it were a pull request headed toward production.

#### Problem

Concurrency bugs often hide in lifecycle edges rather than happy-path job execution. The code needs focused review around worker startup, shutdown, cancellation, queue closing, temporary-file cleanup and race-detector coverage.

#### Design Discussion

This lesson is a review and hardening lesson. It should not add a new business workflow. Any code change should be small and directly tied to shutdown safety, data-race prevention or test coverage.

The main review target is worker-pool lifecycle coordination. `Start` and `Stop` are idempotent, but lifecycle calls should be explicitly serialized so `sync.WaitGroup.Add` cannot overlap incorrectly with `Wait` and so the worker-context cancel function is read and written under synchronization.

#### Go Concepts

- `sync.WaitGroup` lifecycle rules
- race detector as a verification tool, not a proof of correctness
- serialized lifecycle state with a mutex
- graceful shutdown versus forced cancellation
- code review for concurrent state ownership

#### Architecture Concepts

- concurrency review before milestone closure
- explicit shutdown contract for infrastructure components
- known durability gaps documented rather than hidden
- milestone review discipline before starting a new domain area

#### Implementation Plan

- Add a focused lifecycle guard to `WorkerPool` if review confirms the race risk.
- Add or update tests for concurrent start/stop and stop-before-start behavior.
- Run race detector across jobs and CSV import packages.
- Run full build, tests, vet and lint.
- Complete the Milestone 11 review and advance the roadmap to Milestone 12.

#### Tests

- Worker pool stop-before-start behavior.
- Concurrent `Start` and `Stop` does not race or panic.
- Existing cancellation and async import tests remain race-clean.
- Full project verification remains green.

#### Exercises

- Remove the lifecycle mutex and run the race test repeatedly to understand what it protects.
- Explain why `WaitGroup.Add` must not race with `Wait` when the counter can be zero.
- Add an intentionally context-ignoring handler and observe shutdown timeout behavior.
- List which job guarantees are lost on process restart.

#### Interview Questions

- What does the Go race detector detect, and what does it not prove?
- Why can `sync.WaitGroup` be misused even when there is no shared map?
- What is the difference between graceful shutdown and cancellation?
- Why is lifecycle state often protected by a separate mutex?
- What would make this background job system production durable?

### Lesson 11.5 Completion Notes

#### Business Context

Milestone 11 is complete. MES Lite now has a background-job foundation that can run long-lived work outside the original HTTP request and expose job status to clients.

#### Problem

The background job system worked on the happy path, but concurrency systems fail most often at lifecycle boundaries: startup, shutdown, cancellation and shared-state access.

#### Design Discussion

The review focused on the worker-pool lifecycle. `Start` and `Stop` were idempotent, but their internal coordination was not explicit enough for `sync.WaitGroup` rules. A `WaitGroup` must not have `Add` racing with `Wait` when the counter may be zero.

The worker pool now serializes lifecycle startup/shutdown state with a small mutex. This keeps worker registration, `started` state and worker-context cancellation coordinated without changing the public API.

The lesson also validated the async CSV path under the race detector. The job queue, worker pool, CSV import and server route tests are clean under race runs.

#### Implementation

- Added lifecycle serialization to `WorkerPool`.
- Added explicit `started` state so `Stop` before `Start` returns safely.
- Moved worker-context cancellation behind a synchronized helper.
- Preserved the existing `Start`, `Stop` and `Cancel` APIs.
- Kept durable job storage, job listing and retry policy out of scope.

#### Tests

- Added `WorkerPool` stop-before-start coverage.
- Added concurrent start/stop race coverage.
- Verified with `go test ./internal/platform/jobs -race -count=5 -shuffle=on`.
- Verified with `go test ./internal/platform/jobs ./internal/csvimport ./internal/server -race -count=2 -shuffle=on`; this timed out after jobs and csvimport passed because server race tests are slower.
- Re-ran server race tests with a longer timeout: `go test ./internal/server -race -count=1 -shuffle=on`.
- Verified with `go build ./...`.
- Verified with `go test ./... -count=1`.
- Verified with `go vet ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

The fix is intentionally narrow. No generic lifecycle manager or worker framework was introduced. The worker pool owns its own lifecycle because it is currently the only component that needs this coordination.

#### Code Review

An experienced Go engineer would approve the concurrency foundation for the current project stage: shared job state is queue-owned and mutex-protected, worker concurrency is bounded, cancellation is cooperative, temp files are cleaned up on normal success/failure paths and race tests cover the key packages.

The main production-readiness caveat is durability. In-memory jobs and temp-file payloads are acceptable for this learning milestone but not for reliable production processing across restarts.

#### Exercises

- Remove the lifecycle mutex and run the concurrent start/stop test with `-race` repeatedly.
- Add an intentionally context-ignoring handler and observe `Stop` timeout behavior.
- Design a durable jobs table schema with status, payload, result and retry metadata.
- Explain which parts of the current system should change if jobs move to PostgreSQL.

#### Interview Questions

- What does the race detector find, and what can it miss?
- Why is `WaitGroup.Add` versus `Wait` ordering important?
- Why is cooperative cancellation safer than forcibly killing a goroutine?
- How do you decide whether to use channels, mutexes or both?
- What guarantees are required for durable background jobs?

#### Roadmap Update

- Lesson 11.5 completed.
- Milestone 11 completed.
- Current milestone moved to Milestone 12.
- Current lesson moved to Lesson 12.1.
- Architecture maturity, Go knowledge progress and interview readiness updated.

### Milestone 11 Review

#### Architecture Review

An experienced Go engineer would approve Milestone 11 as a learning-oriented concurrency foundation. The implementation uses channels for producer/consumer handoff, mutexes for shared state, a worker pool for bounded execution and context cancellation for cooperative shutdown.

The most important architectural boundary is that `internal/platform/jobs` remains business-free. CSV import registers business work from the composition root instead of making the platform package import a slice.

#### Code Review

The code is explicit and reviewable. Queue state transitions are named methods instead of a generic update callback. Workers receive job copies instead of pointers to canonical state. The status API exposes DTOs, and successful async imports record the import summary as job result data.

The main weaknesses are known and documented: jobs are in-memory, temp-file payloads are not durable, orphan temp files are possible after process crashes and there is no job list/retention policy yet.

#### Refactoring

The milestone included two useful refactors: shared identifier generation moved to `internal/platform/ids`, and technical packages moved to `internal/platform/` with depguard rules enforcing dependency direction.

The final L11.5 refactor serialized worker-pool lifecycle state to make startup/shutdown behavior safer under concurrent calls.

#### Interview Review

You should now be able to explain goroutines, buffered channels, channel ownership, `select`, mutex-protected shared state, `sync.WaitGroup`, worker pools, cooperative cancellation, race detector usage, shutdown trade-offs and why durable background jobs require persistence beyond in-memory queues.

#### Completion Criteria

- Background job model implemented.
- In-memory queue implemented.
- Worker pool implemented.
- Job progress, result, status and cancellation implemented.
- Async CSV import job implemented.
- Worker shutdown and lifecycle reviewed.
- Race detector passes for concurrency-focused packages.
- Tests, build, vet and lint pass.
- Roadmap updated.
