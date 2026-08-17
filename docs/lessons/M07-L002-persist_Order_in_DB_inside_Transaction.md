### Lesson 7.2 Scope

Persist production orders in PostgreSQL and enforce product/employee reference integrity at the database boundary.

#### Business Context

Production orders are planning records. Managers need planned orders to survive restarts and remain connected to valid products and assigned employees.

#### Problem

The order aggregate existed only in memory. Persisting only the root row would lose order lines and assignments, while saving multiple tables without a transaction could leave partial orders behind after a failed child insert.

#### Design Discussion

The order aggregate is stored across three tables: `production_orders`, `production_order_lines` and `production_order_assignments`. Lines reference `products(sku)` and assignments reference `employees(id)`.

Saving an order uses an explicit PostgreSQL transaction because the aggregate requires multiple writes. If any line or assignment fails, the root order row is rolled back too.

`RestoreOrder` rebuilds a persisted aggregate through one domain function instead of exposing fields or spreading unchecked struct literals through the repository.

#### Go Concepts

- explicit transaction handling with pgx
- domain reconstruction with validation
- defensive conversion between sqlc rows and domain values
- error translation from PostgreSQL constraint errors

#### Architecture Concepts

- aggregate persistence across multiple tables
- transaction boundary around one aggregate save
- database foreign keys as reference-integrity guardrails
- sqlc-generated types kept below the orders package boundary

### Lesson 7.2 Completion Notes

#### Business Context

Production orders now have durable storage. Planned order lines reference real products, and assigned employees reference real employees.

#### Problem

The orders slice had a tested aggregate but no persistence. A future API would have had nowhere durable to store manager-created orders.

#### Design Discussion

Added SQL-first persistence for the aggregate. The root row, lines and assignments are saved in one transaction. This is the first lesson where a transaction solves a concrete business consistency problem: avoiding a persisted order without all of its required child records.

Reference integrity is enforced with PostgreSQL foreign keys. Application validation still protects aggregate shape, while the database protects cross-table references.

#### Go Concepts

- `pgx` transaction lifecycle with `Begin`, `Commit` and deferred rollback
- reconstructing private-field aggregates through `RestoreOrder`
- converting `pgtype.Timestamptz` to `time.Time`
- translating SQLSTATE errors into domain errors

#### Architecture Concepts

- aggregate persistence boundary
- transaction per aggregate save
- generated SQL package hidden behind repository adapter
- database constraints as the final consistency boundary

#### Implementation

- Added migration `0006_create_production_orders.sql`.
- Added `production_orders`, `production_order_lines` and `production_order_assignments` tables.
- Added foreign keys from order lines to products and assignments to employees.
- Added orders sqlc queries and generated `internal/orders/ordersdb`.
- Added `orders.Store` and `orders.PostgresStore`.
- Added `RestoreOrder` for validated persistence reconstruction.
- Added duplicate order and invalid reference error mapping.

#### Tests

- Added orders PostgreSQL store integration tests.
- Tested save and find with multiple order lines and an assigned employee.
- Tested duplicate order save mapping to `ErrAlreadyExists`.
- Tested missing order mapping to `ErrNotFound`.
- Tested missing product reference rolls back the root order row.
- Tested missing employee reference rolls back the root order row.
- Verified with `make sqlc`.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

No generic transaction manager was introduced. The order store owns the transaction because this is currently the only multi-write aggregate persistence operation.

#### Code Review

An experienced Go engineer would approve the transaction boundary because it protects exactly one aggregate save and does not introduce framework-like transaction abstraction. The main follow-up is API wiring in Lesson 7.3.

#### Exercises

- Explain why saving an order requires a transaction but saving a product does not.
- Add a test that proves a duplicate order line cannot be inserted by bypassing the domain model.
- Draw the three order tables and identify every foreign key.

#### Interview Questions

- What should define a transaction boundary?
- Why keep foreign keys if the application already validates references?
- Why should sqlc-generated row types not become the public domain API?
- How does deferred rollback work after a successful commit?

#### Roadmap Update

- Lesson 7.2 completed.
- Current lesson moved to Lesson 7.3.
- Persistence `Transactions` marked complete in the Knowledge Matrix.
