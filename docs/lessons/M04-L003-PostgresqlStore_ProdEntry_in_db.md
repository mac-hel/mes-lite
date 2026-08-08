### Lesson 4.3 Completion Notes

#### Business Context

Production entries now need to survive application restarts. The production registration endpoint must write to PostgreSQL instead of only using in-memory state.

#### Problem

The project had a schema and sqlc queries, but the running application still used `production.InMemoryStore`.

#### Design Discussion

The handler continues to depend on the small `production.Store` interface. The new PostgreSQL store adapts sqlc-generated database types to the domain `production.Entry` type. This keeps generated code out of HTTP handlers and preserves the package boundary.

#### Go Concepts

- `context.Context` propagation from HTTP handler to pgx/sqlc calls
- adapter code between generated persistence models and domain models
- error translation with `errors.Is` and `errors.As`
- integer-size boundary between Go `int` and PostgreSQL `integer`

#### Architecture Concepts

- repository implementation inside the production vertical slice
- generated SQL package isolated below domain/application code
- composition root chooses concrete dependencies
- infrastructure errors translated to domain errors

#### Implementation

- Added `production.PostgresStore` backed by sqlc queries.
- Added UUID conversion helpers between string IDs and `pgtype.UUID`.
- Mapped PostgreSQL duplicate-key and constraint errors to domain errors.
- Wired `cmd/server` to use `pgxpool` and `production.NewPostgresStore`.
- Refactored `cmd/server` to use `run() int` so database pool cleanup defers run before process exit.
- Added a quantity overflow guard because PostgreSQL `integer` is 32-bit.

#### Tests

- Added PostgreSQL repository integration tests for save/find, duplicate detection, not found and invalid UUID handling.
- Tests run migrations before exercising the repository.
- Tests skip only when PostgreSQL is unavailable.
- Verified with local Docker PostgreSQL running.

#### Refactoring

- Preserved the existing in-memory store for fast handler tests.
- Kept the `Store` interface consumer-owned by the handler package slice instead of exposing sqlc directly.

#### Code Review

- An experienced Go engineer would approve this direction because the generated persistence code is isolated, domain errors are preserved and context reaches the database call.
- Remaining gap: registering production does not yet validate employee/product existence in one transaction. Lesson 4.4 will define that transactional boundary and business validation strategy.

#### Exercises

- Explain why the handler should not return `productiondb.ProductionEntry` directly.
- Add a test that proves PostgreSQL rejects invalid quantity even if application validation is bypassed.
- Trace `c.Context()` from the Fuego handler to `pgx.QueryRow`.

#### Interview Questions

- Why do repositories often translate infrastructure errors into domain errors?
- Why should generated sqlc code not become the public API of the business package?
- What happens when a request context is cancelled while pgx is waiting on a query?
- Why must Go code care that PostgreSQL `integer` is 32-bit?

#### Roadmap Update

- Lesson 4.3 completed.
- Current lesson moved to Lesson 4.4.
- Repositories marked complete in the Knowledge Matrix.
