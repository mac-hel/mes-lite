### Lesson 10.2 Scope

Make production registration safe for client retries by adding an explicit request ID to the registration command.

#### Business Context

Workers may submit production from unreliable networks or simple clients. If a client times out after the server already saved an entry, retrying the same request should not create duplicate production history.

#### Problem

Production registration generated a new entry ID for every request. A retry with the same business data created another production entry, inflating quantities and making historical review unreliable.

#### Design Discussion

The API now requires `requestId` in the JSON body for `POST /production-entries`. The database stores it in `production_entries.request_id` and enforces uniqueness through a partial unique index.

The partial index applies only when `request_id <> ''`. This keeps historical CSV imports and legacy rows valid without inventing fake request IDs for data that did not originate from a retryable API command.

Retry behavior is explicit: the same `requestId` with identical production data returns the original entry. The same `requestId` with different production data returns `409 Conflict`.

#### Go Concepts

- idempotent command handling
- comparing domain values while ignoring generated IDs
- sentinel conflict errors with `errors.Is`
- unique database constraints translated into domain errors

#### Architecture Concepts

- idempotency at the command/API boundary
- database uniqueness as a concurrency guardrail
- application service resolving safe retries
- import workflow kept separate from API registration workflow

### Lesson 10.2 Completion Notes

#### Business Context

Production registration is now retry-safe for API clients that provide a stable `requestId`.

#### Problem

Repeated HTTP registrations could create duplicate production entries if a client retried after a timeout or uncertain network failure.

#### Design Discussion

Added request-ID idempotency to the registration command instead of exposing client-generated production-entry IDs. The server still owns entry identity, while the client owns retry identity.

`production.Entry` keeps `RequestID` optional so CSV import and historical rows can persist entries without request IDs. `Service.Register` requires a request ID because it represents the retryable API command path.

#### Implementation

- Added migration `0009_add_production_entry_request_ids.sql`.
- Added `request_id` to `production_entries` with a partial unique index.
- Added `RequestID` to `production.Entry`.
- Added `NewEntryWithRequestID` for API registration while preserving `NewEntry` for non-idempotent historical/import paths.
- Added `ErrRequestConflict`.
- Required `requestId` in `RegisterProductionRequest` and `Service.Register`.
- Added `Store.FindByRequestID`.
- Mapped duplicate PostgreSQL request IDs to `ErrRequestConflict`.
- Made duplicate identical retries return the original production entry.
- Made duplicate different retries return `409 Conflict`.
- Updated CSV import persistence to insert blank request IDs for historical rows.

#### Tests

- Added service test for idempotent retry returning the existing entry.
- Added service test for request-ID conflict with different production data.
- Added service test proving request ID is required for registration commands.
- Added handler tests for idempotent retry and conflict responses.
- Added PostgreSQL store test for duplicate request-ID constraint mapping and lookup.
- Added server route test for end-to-end idempotent retry behavior.
- Verified with `go test ./internal/production ./internal/server ./internal/csvimport -count=1`.

#### Refactoring

The domain constructor was extended with `NewEntryWithRequestID` instead of replacing every existing call site. This keeps CSV import and legacy test setup simple while making API registration stricter at the service boundary.

#### Code Review

An experienced Go engineer would approve the partial unique index and retry semantics. The main trade-off is putting `requestId` in the JSON body instead of an `Idempotency-Key` header. The body field is acceptable for this lesson because it keeps the command DTO explicit and OpenAPI-friendly.

#### Exercises

- Add a test that fires two concurrent registrations with the same `requestId` and proves only one row is persisted.
- Change the API to use an `Idempotency-Key` header and compare the OpenAPI and handler trade-offs.
- Explain why the database unique index is still required even though the service checks duplicates.

#### Interview Questions

- What does idempotency mean for a write endpoint?
- Why should retries return the original result instead of blindly returning `409`?
- Why is a unique database constraint safer than an application-only duplicate check?
- Why separate server-generated entry IDs from client-provided request IDs?

#### Roadmap Update

- Lesson 10.2 completed.
- Current lesson moved to Lesson 10.3.
