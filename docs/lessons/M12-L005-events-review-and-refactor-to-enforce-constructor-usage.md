### Lesson 12.5 Scope

Review the machine integration slice for package boundaries, synchronization choices, idempotency behavior, race safety and milestone readiness.

#### Business Context

The fake machine API is not production machine connectivity, but it establishes the learning foundation for machine-originated events. Before moving to observability, the slice should be reviewed as a complete learning milestone.

#### Problem

Milestone 12 introduced events, deduplication, synchronization and atomic counters. The remaining question was whether the exported in-memory store type should continue supporting zero-value construction or whether external consumers should be forced through the constructor.

#### Design Discussion

The constructor-boundary refactor was accepted. `InMemoryStore` was renamed to an unexported `inMemoryStore`, and `NewInMemoryStore` now returns the `Store` interface. This is a deliberate exception to the usual Go preference for returning concrete types: the goal is to hide fake infrastructure details and force consumers to receive a correctly initialized store.

Because construction is now enforced outside the package, the previous `sync.Once` initialization is no longer needed. A constructor-initialized map is simpler. The synchronization lesson remains valid historically, but the final slice favors the smaller production shape.

#### Go Concepts

- unexported concrete implementation types
- constructor-enforced initialization
- when returning an interface from a constructor can be justified
- reviewing `sync.Once` after requirements change
- milestone-level race-detector verification

#### Architecture Concepts

- package API as an enforcement boundary
- fake infrastructure hidden behind a small interface
- milestone closure through review-driven refactoring
- documented trade-off between zero-value usability and constructor enforcement

### Lesson 12.5 Completion Notes

#### Business Context

Milestone 12 is complete. MES Lite now has a fake machine-event intake path that demonstrates event modeling, idempotent event intake, synchronization and atomic counters.

#### Problem

The machine slice was functionally complete for the milestone, but the final review identified one API-boundary issue: an exported `InMemoryStore` allowed external consumers to bypass the constructor even though the store had initialization requirements.

#### Design Discussion

The final design hides the concrete in-memory store and exposes it through `NewInMemoryStore`. This answers the constructor question directly: Go cannot force constructor use for an exported struct, so the struct must be unexported if constructor use matters.

This does mean giving up zero-value usability for that concrete store outside the package. That is acceptable here because this is fake infrastructure, not a domain value type. The public API still gives callers a ready-to-use store.

The milestone also reviewed atomics versus mutexes. The store uses a mutex because it protects compound state. The service uses atomics because intake counters are independent scalar values.

#### Implementation

- Renamed exported `InMemoryStore` to unexported `inMemoryStore`.
- Changed `NewInMemoryStore` to return the `Store` interface.
- Initialized the deduplication map in the constructor.
- Removed `sync.Once` from the final in-memory store implementation.
- Kept mutex-protected save, list and lookup behavior unchanged.
- Preserved all fake machine API routes and response contracts.

#### Tests

- Updated concurrent store test to use constructor-created stores.
- Re-ran focused machine and server tests.
- Verified race safety with `go test ./internal/machines -race -count=5 -shuffle=on`.
- Verified full project behavior with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `go vet ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

The refactor removes synchronization that was only needed for zero-value initialization. This is the right simplification after choosing constructor enforcement.

No durable machine-event persistence was introduced. The milestone intentionally remains about synchronization and event intake, not database-backed integration reliability.

#### Code Review

An experienced Go engineer would approve the milestone as a learning slice. The code models events separately from production entries, handles duplicate delivery idempotently, protects shared state with a mutex, uses atomics for independent counters and has race-detector coverage for concurrent intake paths.

The main production gaps are explicit: events are in-memory, machine authentication is fake, counters reset on restart and no event-to-production-entry processing exists yet.

#### Exercises

- Try exporting `InMemoryStore` again and explain why constructor enforcement is lost.
- Design a PostgreSQL-backed machine event store with a unique `(machine_id, external_event_id)` index.
- Decide whether machine stats should be exposed through this endpoint after Prometheus metrics exist.

#### Interview Questions

- How can Go package visibility enforce constructor use?
- Why might returning an interface from a constructor be acceptable for hidden infrastructure?
- Why are atomics appropriate for counters but not for updating a map and slice together?
- What would make this fake machine integration production-ready?

#### Roadmap Update

- Lesson 12.5 completed.
- Milestone 12 completed.
- Current milestone moved to Milestone 13.
- Current lesson moved to Lesson 13.1.
- Architecture maturity, Go knowledge progress and interview readiness updated.
