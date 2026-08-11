### Lesson 5.1 Scope

Persist employees and products in PostgreSQL so production registration no longer depends on in-memory reference data.

#### Business Context

Employees and products are master data. If they disappear after a restart, production records cannot be trusted as durable business history.

#### Problem

Production entries already persist in PostgreSQL, but employee and product stores are still in-memory. This creates mixed persistence and prevents database-backed reference validation.

#### Design Discussion

This lesson keeps the existing employee/product package APIs and replaces only the infrastructure implementation. The HTTP handlers still depend on the existing `Store` interfaces, while the composition root chooses PostgreSQL stores for the running server.

Foreign keys from production entries to employees/products are intentionally postponed to Lesson 5.4 because they require a migration strategy for existing production data and a discussion about transaction boundaries.

#### Go Concepts

- persistence adapters between domain structs and generated SQL structs
- wrapped sentinel errors for storage failures
- context propagation through repository methods

#### Architecture Concepts

- persistence implementation inside each vertical slice
- composition root selects concrete infrastructure
- domain-facing APIs remain stable while storage changes

---

### Lesson 5.1 Completion Notes

#### Business Context

Employees and products are master data required for trustworthy production history.

#### Problem

Production entries were PostgreSQL-backed, but employees and products still lived in memory. Restarting the server lost reference data and kept production registration only partially durable.

#### Design Discussion

The existing handler-facing store interfaces stayed unchanged. Each vertical slice now owns its SQL queries, generated sqlc package and PostgreSQL adapter. This keeps generated database types below the package boundary and avoids leaking infrastructure details into HTTP handlers or production registration logic.

Foreign keys were intentionally postponed to Lesson 5.4. Adding them safely needs a migration and transaction-boundary discussion, especially for existing production rows that may reference master data created before this lesson.

#### Go Concepts

- sqlc-generated code as an implementation detail
- adapter functions from database rows to domain structs
- sentinel errors wrapped with storage context
- package-local integration test setup
- PostgreSQL advisory locks to serialize concurrent test migrations

#### Architecture Concepts

- persistence adapters inside vertical slices
- composition root choosing concrete infrastructure
- stable domain-facing APIs while infrastructure changes
- explicit remaining consistency gap before foreign keys

#### Implementation

- Added `employees` and `products` PostgreSQL tables in migration `0002`.
- Added sqlc query files for employee and product create/read/list/update/search operations.
- Generated `employeesdb` and `productsdb` packages.
- Added `employees.PostgresStore` and `products.PostgresStore`.
- Wired `cmd/server` to use PostgreSQL employee and product stores.
- Preserved in-memory stores for fast handler tests.

#### Tests

- Added employee PostgreSQL store integration tests.
- Added product PostgreSQL store integration tests.
- Updated repository test setup to serialize migrations with a PostgreSQL advisory lock.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

The server composition root now uses durable stores for all currently persisted business data. No handler refactor was needed because the existing interfaces already expressed the right package boundary.

#### Code Review

An experienced Go engineer would approve the direction: SQL remains explicit, generated code is isolated, and domain errors are preserved. The main remaining issue is referential integrity: production entries still reference employee/product IDs as text without database foreign keys.

#### Exercises

- Explain why `employeesdb.Employee` should not be returned directly from HTTP handlers.
- Add a database constraint test that proves blank product names are rejected by PostgreSQL.
- Trace `POST /employees` from the handler to `pgx.QueryRow`.

#### Interview Questions

- Why keep sqlc-generated types behind repository adapters?
- What problem do advisory locks solve in integration tests?
- Why can a database constraint still be useful when Go validation already exists?
- What are the trade-offs of adding foreign keys in a later migration?

#### Roadmap Update

- Lesson 5.1 completed.
- Current lesson moved to Lesson 5.2.
- Known technical debt updated: employees/products are now PostgreSQL-backed; foreign keys and transaction consistency remain pending.
