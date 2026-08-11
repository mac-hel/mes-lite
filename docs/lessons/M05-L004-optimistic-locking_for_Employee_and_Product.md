### Lesson 5.4 Completion Notes

#### Business Context

Administrators and team leaders may edit the same employee or product around the same time. Without conflict detection, the last write silently wins and can overwrite another user's changes.

#### Problem

Employee and product updates replaced rows without checking whether the caller edited a stale copy. This creates lost updates under concurrent usage.

#### Design Discussion

The lesson introduces optimistic locking with an integer `version` column. Clients receive the current version when reading or creating an employee/product. Update requests must submit the version they edited. The database update succeeds only when the submitted version matches the stored version, then increments the version atomically.

This keeps the solution simple and explicit. We do not hold long database locks across HTTP requests. Instead, conflicts are detected at write time and returned as `409 Conflict`.

#### Go Concepts

- optimistic concurrency with version fields
- stale-write detection with sentinel errors
- update methods returning updated structs
- concurrent integration tests with goroutines and channels

#### Architecture Concepts

- conflict detection belongs at the persistence boundary
- HTTP translates stale writes to `409 Conflict`
- database update predicates enforce atomic compare-and-swap behavior
- handlers return the incremented version after successful writes

#### Implementation

- Added `Version` to employee and product domain structs.
- Added `ErrVersionConflict` for employees and products.
- Added migration `0003_add_employee_product_versions.sql`.
- Updated sqlc create/get/list/update queries to include `version`.
- Changed employee/product store `Update` methods to return the updated entity.
- Updated in-memory stores to reject stale versions and increment versions.
- Updated PostgreSQL stores to use `WHERE id/sku = $1 AND version = $6` update predicates.
- Updated employee/product update HTTP requests to require `version`.
- Mapped stale updates to `409 Conflict`.

#### Tests

- Updated handler tests to submit versions and assert incremented response versions.
- Added handler tests for stale-version conflicts.
- Added PostgreSQL stale-version tests for employees and products.
- Added a PostgreSQL concurrent update test using two goroutines updating the same employee version; one update succeeds and one returns `ErrVersionConflict`.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

Store update contracts now return the updated entity so handlers do not guess the next version. This keeps version increments owned by persistence implementations.

#### Code Review

An experienced Go engineer would approve this optimistic-locking design for simple master-data updates. The main trade-off is API friction: clients must now send a version on updates. That is intentional because silent overwrites are worse than explicit conflicts.

#### Exercises

- Explain why the version check must happen in the SQL `WHERE` clause.
- Add a product concurrent update test mirroring the employee test.
- Design the client behavior after receiving `409 Conflict`.

#### Interview Questions

- What problem does optimistic locking solve?
- How is optimistic locking different from pessimistic locking?
- Why does `UPDATE ... WHERE version = ?` avoid lost updates?
- When would `SELECT ... FOR UPDATE` be a better choice?

#### Roadmap Update

- Lesson 5.4 completed.
- Current lesson moved to Lesson 5.5.
- Known technical debt updated: optimistic locking completed; foreign-key-backed transactional reference integrity remains pending.
