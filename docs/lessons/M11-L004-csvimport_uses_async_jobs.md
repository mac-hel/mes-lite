### Lesson 11.4 Scope

Add an asynchronous CSV import endpoint backed by the background job queue and worker pool.

#### Business Context

Managers should be able to upload historical production CSV data, receive a job ID quickly and monitor import completion through the job status API instead of keeping one HTTP request open for the whole import.

#### Problem

The CSV import pipeline is streaming internally, but the HTTP endpoint still runs validation and persistence synchronously. Large imports keep the request open and make cancellation/status visibility impossible for clients.

#### Design Discussion

The synchronous endpoint remains available for simple MVP usage. Lesson 11.4 adds a separate async endpoint so the API contract is explicit:

- `POST /imports/production-entries/jobs` accepts CSV input and returns a queued job snapshot.
- `GET /jobs/{id}` lets clients monitor queued, running, succeeded, failed or cancelled state.
- `PUT /jobs/{id}/cancel` requests cancellation.

The job payload stores metadata pointing at a temporary upload file, not the whole CSV content. This is the smallest production-minded choice for the current in-memory job queue: queued job metadata stays small, while the CSV service can still stream from an `io.Reader` when the worker runs.

The temporary file is removed by the job handler after execution. This is not durable job storage; a process restart still loses queued/running jobs and any unprocessed temporary uploads, which remains covered by ADR 0004.

#### Go Concepts

- `os.CreateTemp` for temporary upload handoff
- `io.Copy` for streaming request body to disk
- JSON encoding for opaque job payload metadata
- worker handlers that call existing application services
- cancellation propagation from worker context into CSV import processing

#### Architecture Concepts

- async endpoint separated from synchronous endpoint
- composition root wires business work into platform workers
- queue payload carries metadata, not business dependencies
- temporary-file handoff as an implementation detail before durable queues exist

#### Implementation Plan

- Add async import service/handler behavior that writes upload data to a temporary file and enqueues a production-entry-import job.
- Add a job handler that decodes the payload, opens the temporary file and calls the existing CSV import service.
- Store import summary JSON in successful job state.
- Register `POST /imports/production-entries/jobs` with the same admin/manager RBAC as synchronous import.
- Wire the job handler into the worker pool in `cmd/server`.
- Preserve the existing synchronous import endpoint.

#### Tests

- Async handler enqueues a job and returns a job response.
- Queue-full errors become a clear HTTP response.
- Job handler imports valid CSV and records success result.
- Job handler records validation/import failure as failed job state.
- Temporary files are removed after job execution.
- Server route tests cover admin/manager access and worker rejection.
- Race detector remains clean for jobs and csvimport-focused tests.

#### Exercises

- Add a max upload size before writing the temporary file and decide which HTTP status should be returned.
- Simulate process restart after enqueue and explain why the job is lost today.
- Change the payload to store CSV bytes and compare memory behavior with temporary files.
- Add a test proving async import cancellation stops between CSV rows.

#### Interview Questions

- Why is writing an upload to a temporary file different from reading it all into memory?
- Why should the composition root connect job handlers to business services?
- What cleanup responsibilities appear when temporary files are used?
- Why keep the synchronous endpoint while adding an async endpoint?
- What would need to change for this to survive process restarts?

### Lesson 11.4 Completion Notes

#### Business Context

MES Lite can now accept a production-entry CSV import as background work. Managers get a job snapshot immediately and can use the job status endpoint from L11.3 to monitor completion.

#### Problem

The CSV import pipeline streamed rows internally, but the HTTP request still waited until validation and persistence finished. That made long imports poor API citizens and gave clients no status or cancellation handle.

#### Design Discussion

Added `POST /imports/production-entries/jobs` instead of changing the existing synchronous endpoint. This keeps the old simple workflow available and makes async behavior explicit in the route name.

The async service writes the uploaded CSV body to a temporary file and stores only the file path in the job payload. This keeps queued job payloads small and lets the worker reuse the existing streaming CSV service by opening the file as an `io.Reader`.

The worker handler lives in `csvimport`, while the worker pool remains in `internal/platform/jobs`. The composition root wires them together. This preserves the platform rule: the generic jobs package does not import business slices.

Successful jobs record the `ImportSummary` JSON as job result data. Failed jobs record the failure message through the existing worker-pool failure path. The temporary file is removed after execution in both success and failure paths.

#### Implementation

- Added `Result` to `jobs.Job` with defensive copying.
- Added `Queue.RecordResult` for running jobs.
- Added `result` to job HTTP responses.
- Added `csvimport.AsyncService`.
- Added async upload-to-temp-file handoff.
- Added `csvimport.NewProductionEntriesJobHandler` for worker execution.
- Added `Handler.ImportProductionEntriesAsync`.
- Added `Handler.AsyncEnabled` so routes are only registered when async import is wired.
- Added `POST /imports/production-entries/jobs` with admin/manager RBAC.
- Wired `cmd/server` to register the CSV import job handler with the worker pool.
- Added context cancellation checks to the CSV import streaming loop.

#### Tests

- Tested async service writes upload data to a temporary file and enqueues a job.
- Tested queue-full enqueue failures remove the temporary file.
- Tested the production-entry import job handler imports valid CSV data.
- Tested successful job execution records an import summary result.
- Tested temporary files are removed after successful execution.
- Tested invalid CSV fails the job and still removes the temporary file.
- Added server route test proving managers can enqueue async imports.
- Added server route test proving workers are forbidden.
- Verified with `go test ./internal/platform/jobs ./internal/csvimport ./internal/server -count=1`.
- Verified with `go test ./internal/platform/jobs ./internal/csvimport -race -count=3 -shuffle=on`.
- Verified with `go build ./...`.
- Verified with `go test ./... -count=1`.
- Verified with `go vet ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

The CSV import service remains the owner of row streaming, validation and persistence. Async import is a wrapper around the existing workflow rather than a second import implementation.

`server.New` now registers the async import route only when the CSV import handler is built with async support. Existing tests and simple synchronous wiring do not need to know about background jobs.

#### Code Review

An experienced Go engineer would approve the composition direction: platform jobs do not depend on `csvimport`, and business work is registered from the composition root.

The main caveat is the temporary-file durability story. If the process crashes after upload and before execution, the in-memory job is lost and the temp file may be orphaned. That is acceptable for this lesson because ADR 0004 already records the durable-queue trigger, but it is not production-complete durability.

#### Exercises

- Add a max async upload size using `http.MaxBytesReader` or an `io.LimitedReader` at the HTTP boundary.
- Add a cancellation test that stops a CSV import between rows.
- Add a status API test proving successful async import includes a JSON `result` object.
- Design how durable job storage would reference uploaded files or object storage.

#### Interview Questions

- Why does async HTTP usually return a job ID instead of the final result?
- Why should a background queue payload avoid large blobs?
- How does the composition root prevent `platform/jobs` from importing business packages?
- What cleanup failure modes are introduced by temporary files?
- Why does this still need durable storage before production-grade reliability?

#### Roadmap Update

- Lesson 11.4 completed.
- Current lesson moved to Lesson 11.5.
- Known technical debt updated for temporary-file/orphan behavior and missing job list endpoint.
