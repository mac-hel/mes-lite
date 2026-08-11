### Lesson 5.5 Completion Notes

#### Business Context

Production entries are business records. They must not reference employees or products that do not exist in durable master data.

#### Problem

Production registration validated employee/product references in application code, but the database still stored plain text references without foreign keys. A bug, import, direct SQL write or future service could insert production entries pointing at missing master data.

#### Design Discussion

PostgreSQL now owns reference existence for new production entries through foreign keys from `production_entries.employee_id` to `employees.id` and from `production_entries.product_sku` to `products.sku`.

The constraints are created as `NOT VALID`. This is a production-safe migration pattern when existing data may be dirty: PostgreSQL enforces the constraint for new writes, but does not scan and reject legacy rows during deployment. A future cleanup migration can validate existing rows and then run `VALIDATE CONSTRAINT`.

Application validation remains valuable because it returns better business errors and checks active/inactive state. Database constraints are the final integrity boundary for reference existence.

#### Go Concepts

- database constraint errors translated into domain errors
- integration tests that prove failed writes do not leave rows behind
- distinguishing business validation from persistence integrity

#### Architecture Concepts

- referential integrity belongs in the database
- application services still own workflow rules
- database constraints protect against alternate write paths
- production-safe migration with `NOT VALID` foreign keys

#### Implementation

- Added migration `0004_add_production_reference_foreign_keys.sql`.
- Added foreign key from `production_entries.employee_id` to `employees.id`.
- Added foreign key from `production_entries.product_sku` to `products.sku`.
- Mapped PostgreSQL foreign-key violations to `production.ErrInvalidEntry` at the production persistence boundary.
- Updated production PostgreSQL integration tests to seed employee/product reference data.

#### Tests

- Added production repository test for missing employee reference.
- Verified failed foreign-key insert leaves no production row behind.
- Verified with `make sqlc`.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

No broad transaction manager was introduced. The current write is a single PostgreSQL statement, so PostgreSQL already executes it atomically. A future multi-step workflow can introduce an explicit transaction abstraction when there is a real multi-write business operation.

#### Code Review

An experienced Go engineer would approve the foreign-key enforcement and the `NOT VALID` migration choice for safe rollout. The main caveat is that legacy rows are not validated yet; this is documented technical debt rather than a hidden inconsistency.

#### Exercises

- Explain why `NOT VALID` foreign keys still protect new writes.
- Add a cleanup query that finds existing production rows with missing employee references.
- Explain why active/inactive validation stays in application code instead of a foreign key.

#### Interview Questions

- What is referential integrity?
- What does a PostgreSQL `NOT VALID` constraint do?
- Why keep application validation if the database has foreign keys?
- When should you introduce an explicit transaction boundary?

#### Roadmap Update

- Lesson 5.5 completed.
- Milestone 5 completed.
- Current milestone moved to Milestone 6.
- Known technical debt updated for `NOT VALID` reference constraints and OpenAPI query-parameter documentation.

### Milestone 5 Review

#### Architecture Review

An experienced Go engineer would approve the milestone direction. Employees, products and production entries are now PostgreSQL-backed, generated sqlc types remain below slice boundaries, HTTP handlers depend on domain-facing stores and production registration has both application-level reference validation and database-level referential integrity for new writes.

The main architectural limitation is that cross-slice transaction coordination is still implicit. This is acceptable because the current production registration persistence step is a single insert. A future multi-write operation should introduce an explicit transaction boundary rather than a generic transaction abstraction in advance.

#### Code Review

The code is explicit and idiomatic for the current maturity level. Constructors now reject invalid domain entities, repositories validate defensively, update endpoints use optimistic locking and SQL sorting avoids dynamic SQL interpolation.

The main improvement is API documentation: list query parameters are implemented and tested, but OpenAPI query parameter metadata should be reviewed later.

#### Refactoring

Store files were split into `store.go`, `store_in_memory.go` and `store_postgres.go`, improving navigation without changing exported names. Product search now reuses product listing options instead of having a separate persistence method.

#### Interview Review

You should now be able to explain sqlc vs ORM, repository adapters, constructor validation, database constraints, `NOT VALID` foreign keys, limit/offset pagination, safe dynamic sorting and optimistic locking with version columns.

#### Completion Criteria

- Employees and products persist in PostgreSQL.
- Validation flow is documented and constructors enforce invariants.
- List endpoints support pagination, filtering and sorting.
- Optimistic locking prevents stale employee/product updates.
- Production reference existence is enforced by PostgreSQL foreign keys for new writes.
- Tests, build, lint and sqlc generation pass.
- Roadmap updated.
