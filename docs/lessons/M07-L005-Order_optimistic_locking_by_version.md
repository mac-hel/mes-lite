### Lesson 7.5 Scope

Review the production-order slice for consistency gaps before closing Milestone 7.

#### Business Context

Production-order planning is now writable through several workflows. Managers and leaders can assign employees and move orders through lifecycle states, so concurrent edits must not silently overwrite each other.

#### Problem

`PostgresStore.Update` already used a transaction for root state and assignment replacement, but it did not detect stale updates. Two callers could load the same order version, make different valid changes and let the later write overwrite the earlier one.

#### Design Discussion

Added optimistic locking to the production-order aggregate root, matching the existing employee/product persistence pattern. The order root now carries a `version`, persisted updates atomically require the expected version and the database increments it on success.

The `Store.Update(ctx, order) error` shape was kept. Services refetch after a successful update so API responses show the persisted version without widening repository contracts for this lesson.

#### Go Concepts

- private aggregate version field with an accessor
- sentinel conflict errors with `errors.Is`
- preserving an existing interface while improving consistency
- `sync.RWMutex`-protected stale-write checks in the in-memory adapter

#### Architecture Concepts

- optimistic locking for aggregate-root updates
- stale-write detection as part of transactional consistency
- HTTP `409 Conflict` for concurrent modification
- milestone review before moving to reporting

### Lesson 7.5 Completion Notes

#### Business Context

Production-order updates are now protected against silent lost updates. A stale assignment or status transition returns a domain conflict instead of overwriting newer aggregate state.

#### Problem

Lesson 7.4 introduced transactional update, but transaction atomicity alone did not detect concurrent edits based on old state.

#### Design Discussion

Added a `version` column to `production_orders` and threaded it through `Order`, persistence reconstruction, HTTP responses and in-memory tests. PostgreSQL updates now use `WHERE id = $1 AND version = $4` and increment `version` in the same statement.

When no row is updated, the repository distinguishes missing orders from stale versions by checking whether the order still exists. This keeps `ErrNotFound` and `ErrVersionConflict` behavior precise.

#### Implementation

- Added migration `0007_add_production_order_versions.sql`.
- Added `orders.ErrVersionConflict`.
- Added `Order.version`, `Version()` and version validation.
- Updated `RestoreOrder`, sqlc queries and generated `ordersdb` code.
- Updated `PostgresStore.Update` to use optimistic locking.
- Updated `InMemoryStore.Update` to mirror version conflict behavior.
- Updated service mutation methods to refetch after successful updates.
- Added `version` to production-order HTTP responses.

#### Tests

- Added aggregate version assertion for new orders.
- Added service assertions that mutation responses return incremented versions.
- Added handler assertions for version in create/get/mutation responses.
- Added PostgreSQL integration coverage for persisted version, incremented version, stale update conflict and failed-update rollback preserving version.
- Verified with `make sqlc`.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `golangci-lint run ./...`.

#### Milestone Review

Milestone 7 is complete. The system now has a production-order aggregate, durable multi-table persistence, product/employee reference validation, HTTP create/read/mutation endpoints, role-based access and optimistic locking for mutable order workflows.

The milestone moved the project from CRUD-focused resources toward richer business modeling. The main architectural improvement is that order lifecycle rules live in the aggregate, workflow coordination lives in `orders.Service` and persistence consistency lives in the repository adapter.

#### Follow-Ups

- Production reference foreign keys remain `NOT VALID` as known technical debt.
- Full auth-user CRUD remains postponed until there is a concrete business workflow.
- Reporting should start with read-model/query design rather than more order mutations.

#### Exercises

- Explain why optimistic locking still matters when the update is already inside a transaction.
- Add a failing test for concurrent order release/cancel and make it pass with version conflicts.
- Compare returning the refetched order with changing `Store.Update` to return the updated aggregate.

#### Interview Questions

- What problem does optimistic locking solve?
- How do you distinguish not-found from version-conflict when both return no rows?
- Why should an aggregate root own the version rather than each child table?
- What does HTTP `409 Conflict` communicate to an API client?

#### Roadmap Update

- Lesson 7.5 completed.
- Milestone 7 completed.
- Current milestone moved to Milestone 8.
- Current lesson moved to Lesson 8.1.
- Concurrency `Mutex` and `RWMutex` marked complete in the Knowledge Matrix.
