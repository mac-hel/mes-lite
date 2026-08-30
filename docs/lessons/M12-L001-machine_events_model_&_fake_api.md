### Lesson 12.1 Scope

Introduce the machine-integration vertical slice with a machine event model and a fake API endpoint for manually sending events into the system.

#### Business Context

Machine integration is future scope for the business, but it is a useful learning step after the MVP API and background-job foundations. Before processing concurrent machine events, the system needs a clear event shape and an endpoint that can simulate machine input without connecting to real CNC equipment.

#### Problem

The project had no concept of a machine event. Jumping straight to concurrent processing or deduplication would mix event modeling, HTTP input, idempotency and synchronization in one lesson.

#### Design Discussion

L12.1 adds `internal/machines` as a business slice. The fake API accepts a machine ID from the path and event details from JSON, validates them, normalizes timestamps to UTC and stores the event in memory.

The endpoint is deliberately fake and protected by existing admin/manager JWT RBAC. Real machine credentials, durable storage, deduplication and processing into production entries are postponed to later lessons so this lesson stays focused on the event model and HTTP boundary.

#### Go Concepts

- custom string type for event kinds
- constructor validation with `(T, error)`
- timestamp normalization with `time.Time.UTC`
- mutex-protected in-memory state for HTTP safety
- DTO mapping at the HTTP boundary

#### Architecture Concepts

- machine integration as a vertical slice
- event model separated from production-entry registration
- fake adapter endpoint before real integration protocols
- explicit postponement of deduplication and processing

### Lesson 12.1 Completion Notes

#### Business Context

MES Lite now has a fake machine-event intake path that can be used to simulate future machine integration work.

#### Problem

The application could register manual production entries and run background jobs, but it had no event shape for machine-originated production signals.

#### Design Discussion

Added `machines.Event` as an input event, not as a production entry. A machine event may later produce a production entry, but keeping them separate prevents future machine-specific fields and deduplication state from leaking into manual registration.

The in-memory store is intentionally temporary. It gives L12.2 and L12.3 a concrete event source to evolve without adding database schema and synchronization concerns at the same time.

#### Implementation

- Added `internal/machines` business slice.
- Added `EventType` with `cycle_completed`, `state_changed` and `alarm_raised`.
- Added `machines.Event` with generated ID, machine ID, external event ID, type, occurred-at timestamp, product SKU, quantity, workstation and message.
- Added constructor and validation for machine events.
- Added mutex-protected `machines.InMemoryStore`.
- Added `machines.Handler` with `POST /machines/{machineId}/events`.
- Registered machine routes separately through `server.RegisterMachineRoutes`.
- Wired the fake machine handler in `cmd/server`.
- Updated depguard's platform deny list for the new business slice.

#### Tests

- Added machine event constructor and validation tests.
- Added in-memory store save/list and defensive snapshot tests.
- Added handler tests for valid and invalid fake machine events.
- Added server route tests proving managers can submit fake machine events, workers are forbidden and unauthenticated callers receive `401`.
- Verified with `go test ./internal/machines ./internal/server -count=1`.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `go vet ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

Machine routes were registered through a focused `RegisterMachineRoutes` function instead of widening the already large `server.New` constructor. This matches the optional job-route registration pattern and keeps the main application route constructor stable.

#### Code Review

An experienced Go engineer would approve the lesson scope because it introduces the smallest useful machine-integration boundary and avoids pretending the fake API is production machine connectivity.

The main caveats are intentional: events are not durable, deduplicated or processed yet, and the fake endpoint uses human JWT roles instead of real machine authentication.

#### Exercises

- Add a handler test for `state_changed` without product SKU or quantity.
- Add a `GET /machines/{machineId}/events` proposal and decide who should be allowed to call it.
- Explain why machine events should not be stored directly as production entries at intake time.

#### Interview Questions

- What is the difference between an event and a command?
- Why should integration input be modeled separately from internal production-entry records?
- Why is a fake API useful before a real machine protocol exists?
- What consistency problems appear when machines retry the same event?

#### Roadmap Update

- Lesson 12.1 completed.
- Current lesson moved to Lesson 12.2.
- Known technical debt updated for fake, in-memory machine event intake.
