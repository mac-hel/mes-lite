### Lesson 9.4 Scope

Expose the CSV import workflow through an authenticated API endpoint and return an import summary with partial failure details.

#### Business Context

Managers need to upload historical production CSV files and understand what happened: how many rows were accepted, how many were rejected and which rows need correction.

#### Problem

The project could stream, validate and transactionally save batches internally, but there was no API entry point and no response contract for partial failures.

#### Design Discussion

The API accepts a raw `text/csv` request body at `POST /imports/production-entries`. The handler reads from the request body, while `csvimport.Service` coordinates the workflow.

Row-level validation failures are returned in the summary and do not prevent valid rows from being persisted. Fatal CSV shape errors such as invalid headers or malformed records return `400 Bad Request`. Persistence batch failures are reported with the row number that caused the batch rollback.

#### Go Concepts

- reading raw request bodies through `io.Reader`
- service-level orchestration over reader, validator and store
- response DTOs with JSON tags
- `errors.As` for extracting wrapped `BatchError`
- in-memory adapter for fast API tests

#### Architecture Concepts

- import service as workflow coordinator
- handler owns HTTP translation only
- API summary contract for partial failures
- RBAC-protected management import endpoint

### Lesson 9.4 Completion Notes

#### Business Context

MES Lite now exposes a manager/admin CSV import endpoint for historical production entries.

#### Problem

Internal import components existed, but clients could not upload CSV data or receive a useful summary of imported versus rejected rows.

#### Design Discussion

Added `csvimport.Service` to coordinate raw CSV reading, row validation and batch persistence. The handler stays small: it passes the request body to the service and maps fatal CSV input errors to `400 Bad Request`.

The endpoint returns `ImportSummary` with total rows, valid rows, invalid rows, imported rows and row-level errors. This gives API clients enough information to correct spreadsheet mistakes without silently ignoring invalid data.

#### Implementation

- Added `ImportSummary` and `ImportError` response types.
- Added `csvimport.Service.ImportProductionEntries`.
- Added `csvimport.Handler.ImportProductionEntries`.
- Added `csvimport.InMemoryStore` for fast service/server tests.
- Wired `csvimport.PostgresStore` and handler in `cmd/server`.
- Registered `POST /imports/production-entries` with bearer auth and admin/manager RBAC.

#### Tests

- Added service test for valid rows plus validation errors.
- Added service test for fatal invalid CSV headers.
- Added service test for persistence batch errors reported as row errors.
- Added service test for unexpected store errors.
- Added server route test proving managers can import CSV.
- Added server route tests proving workers are forbidden and unauthenticated callers get `401`.
- Added server route test proving invalid CSV headers return `400`.
- Verified with `go test ./internal/csvimport -count=1`.
- Verified with `go test ./internal/server -count=1`.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

`server.New` now accepts an optional CSV import handler for tests while the production composition root wires the real PostgreSQL-backed handler explicitly.

#### Code Review

An experienced Go engineer would approve the boundary split: the handler does not validate CSV rows, the service does not know about Fuego and persistence remains behind the import store.

The main follow-up is performance review in L9.5. The current service validates into slices before saving; that is acceptable for this lesson but should be reviewed against large-file import goals.

#### Exercises

- Add a server test proving admins can import CSV.
- Add a handler/service test for a CSV containing only invalid rows.
- Decide whether import summaries should include generated production-entry IDs for successful rows.

#### Interview Questions

- Why should the handler pass an `io.Reader` instead of reading the whole body into `[]byte`?
- Why are row-level errors returned in a `200 OK` summary while malformed CSV headers return `400 Bad Request`?
- What belongs in the import service versus the import handler?
- Why should import endpoints be restricted to managers/admins rather than workers?

#### Roadmap Update

- Lesson 9.4 completed.
- Current lesson moved to Lesson 9.5.
- L9.5 remains focused on CSV import review and performance.
