### Lesson 12.2 Scope

Add duplicate detection and idempotent retry behavior for machine events based on the machine's external event identifier.

#### Business Context

Real machines and integration gateways retry messages when networks fail or acknowledgements are lost. MES Lite must not process the same machine signal twice just because delivery happened twice.

#### Problem

L12.1 accepted and stored fake machine events, but the same `machineId` and `externalEventId` could be submitted repeatedly and stored as separate events. That would inflate future production quantities once events start producing production entries.

#### Design Discussion

Deduplication belongs at the state boundary and idempotency belongs in the workflow. The in-memory store now rejects duplicate `(machineID, externalEventID)` keys. The service receives that storage-level duplicate signal, loads the existing event and compares business payloads.

An identical retry returns the original event. A duplicate key with different data returns `ErrEventConflict`, which the HTTP handler translates to `409 Conflict`.

This mirrors the earlier production-registration `requestId` lesson, but keeps machine events separate because the retry identity comes from the external machine system rather than from a human/API client command.

#### Go Concepts

- service-level idempotency workflow
- composite map keys for duplicate detection
- sentinel errors for duplicate versus conflict outcomes
- comparing values while ignoring generated identifiers
- `errors.Is` across handler and service boundaries

#### Architecture Concepts

- event deduplication as an integration boundary concern
- storage uniqueness as a concurrency guardrail
- idempotent retry semantics before event processing
- fake in-memory behavior documented before durable persistence

### Lesson 12.2 Completion Notes

#### Business Context

MES Lite now treats repeated machine delivery safely. A machine can retry the same event without creating duplicate machine-event records.

#### Problem

The fake machine API accepted all submissions independently. A retry with the same external event ID could create duplicate records and later duplicate production output.

#### Design Discussion

Added `machines.Service` as the intake workflow boundary. The handler now translates HTTP into a `ReceiveEventCommand`, while the service constructs the event and handles duplicate semantics.

The in-memory store owns uniqueness for `(machineID, externalEventID)` under its mutex. This matters because an application-only pre-check would be race-prone once concurrent machine submissions arrive.

Identical duplicate delivery returns the existing event. Different payload under the same external key returns `409 Conflict` so clients and operators can see that a producer reused an event identifier incorrectly.

#### Implementation

- Added `ErrDuplicateEvent`, `ErrEventConflict` and `ErrNotFound` to the machines slice.
- Added `Store.FindByExternalEventID`.
- Updated `InMemoryStore` to index events by `(machineID, externalEventID)`.
- Added `machines.Service.ReceiveEvent`.
- Added `ReceiveEventCommand`.
- Changed `machines.Handler` to depend on an `EventReceiver` interface instead of storing events directly.
- Mapped conflicting duplicate machine events to HTTP `409 Conflict`.
- Updated server and production composition wiring to use the machine service.

#### Tests

- Added store tests for duplicate external event rejection and lookup by external event ID.
- Added service tests for new event storage, identical retry returning the original event, conflicting retry and invalid events.
- Added handler tests for idempotent retry and conflicting duplicate payload.
- Existing server route authorization tests continue to pass with service-based wiring.
- Verified with `go test ./internal/machines ./internal/server -count=1`.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `go vet ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

Machine event construction and persistence moved out of the handler and into `machines.Service`. This keeps handlers focused on HTTP and gives later lessons one place to add processing, progress counters or durable idempotency.

#### Code Review

An experienced Go engineer would approve the separation of duplicate detection and idempotent semantics. The store enforces uniqueness atomically under a mutex, and the service decides whether a duplicate is safe or conflicting.

The main caveat is durability. In-memory deduplication is enough for the synchronization lesson, but a real machine integration would need a database uniqueness constraint or durable event log.

#### Exercises

- Add a concurrent test that submits the same machine event from two goroutines and proves only one record is stored.
- Add a test proving the same `externalEventId` may be reused by two different machines.
- Design the PostgreSQL unique index that would make machine-event deduplication durable.

#### Interview Questions

- What is idempotent event processing?
- Why is duplicate detection at the store boundary safer than a handler-level pre-check?
- Why should conflicting duplicate payloads return `409 Conflict`?
- How is machine event idempotency different from HTTP command idempotency?

#### Roadmap Update

- Lesson 12.2 completed.
- Current lesson moved to Lesson 12.3.
- Known technical debt updated: machine event deduplication exists but is still in-memory only.
