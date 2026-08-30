### Lesson 12.3 Scope

Review and harden the machine-event in-memory synchronization path with `sync.Mutex`, `sync.RWMutex` and `sync.Once`.

#### Business Context

Machine integrations send events concurrently. Even the fake in-memory API must behave correctly when many requests arrive at the same time, because duplicate protection is only trustworthy if the shared state is protected consistently.

#### Problem

The in-memory store used a mutex, but initialization and concurrent retry behavior were not explicitly tested under race-oriented workloads. The zero-value store also relied on ad-hoc map initialization inside `Save`.

#### Design Discussion

The store remains the synchronization boundary. Writes use `Lock`, reads use `RLock` and one-time index initialization uses `sync.Once`. `sync.Once` makes the zero value usable and prevents multiple goroutines from racing to initialize the deduplication map.

The service still owns business semantics. It receives duplicate errors from the store and turns them into either idempotent success or conflict. This keeps synchronization separate from workflow rules.

#### Go Concepts

- `sync.Once` for one-time initialization
- `sync.Mutex` for exclusive writes
- `sync.RWMutex` for concurrent-safe snapshots and lookups
- race-detector verification for synchronization assumptions
- zero-value usability for stateful Go types

#### Architecture Concepts

- synchronization owned by the in-memory infrastructure boundary
- service-level idempotency kept separate from locking mechanics
- focused concurrency hardening without adding durable storage
- tests that exercise behavior likely to fail under concurrent delivery

### Lesson 12.3 Completion Notes

#### Business Context

Machine event intake is now safer under concurrent fake-machine submissions. Repeated delivery and concurrent writes are protected by synchronized store state.

#### Problem

Duplicate detection depends on shared maps and slices. Without explicit synchronization tests, a bug could appear only under concurrent machine traffic.

#### Design Discussion

Replaced the store's ad-hoc map initialization with `sync.Once`. The store now has a clear ownership model: `sync.Once` initializes internal indexes, `Lock` protects writes and `RLock` protects read snapshots and lookups.

This lesson did not introduce `sync.Map`. The store updates a slice and a map together, so a single mutex keeps those two structures consistent. `sync.Map` is useful for specific read-mostly or append-only patterns, but it would make this small store harder to reason about.

#### Implementation

- Added `sync.Once` initialization to `machines.InMemoryStore`.
- Preserved exclusive `Lock` for event saves and duplicate-key insertion.
- Preserved `RLock` for event listing and external-event lookup.
- Made the zero-value `InMemoryStore` safe for concurrent saves.
- Kept public machine service and handler behavior unchanged.

#### Tests

- Added concurrent zero-value store save test.
- Added concurrent identical retry test proving many callers receive one stored event and the same generated event ID.
- Verified focused package behavior with `go test ./internal/machines ./internal/server -count=1`.
- Verified race safety with `go test ./internal/machines -race -count=5 -shuffle=on`.
- Verified full project behavior with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `go vet ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

The refactor is intentionally narrow. The machine store remains in memory and keeps the same public API. Only internal initialization and concurrency coverage changed.

#### Code Review

An experienced Go engineer would approve using one `RWMutex` here because the store must keep `events` and `byExternalID` consistent. Splitting locks or using `sync.Map` would add complexity without a demonstrated bottleneck.

The race detector being clean is useful evidence, not a proof. It only checks executions that the tests exercise, so the tests deliberately include concurrent duplicate delivery and concurrent first-use initialization.

#### Exercises

- Remove `sync.Once` and run the zero-value concurrent save test with `-race`.
- Replace the store map with `sync.Map` and explain why keeping the slice consistent becomes less obvious.
- Add a concurrent test where one goroutine lists events while others save events.

#### Interview Questions

- When should you use `sync.Once`?
- Why can one mutex be better than separate locks for related pieces of state?
- What is the difference between `Mutex` and `RWMutex`?
- What does a clean race-detector run prove, and what does it not prove?

#### Roadmap Update

- Lesson 12.3 completed.
- Current lesson moved to Lesson 12.4.
- L12.4 remains focused on atomic counters and runtime race review.
