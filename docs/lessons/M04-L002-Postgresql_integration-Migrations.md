### Lesson 4.2 Completion Notes

#### Business Context

Production entries are core business records. Keeping them only in memory would lose production history and fail the business goal of replacing Excel.

#### Problem

The project had PostgreSQL available in Docker, but no schema migrations, no sqlc configuration and no executable migration command.

#### Design Discussion

The database contract now starts with production entries only. This avoids a large persistence rewrite and keeps the lesson focused on PostgreSQL, migrations and sqlc. The SQL schema owns non-negotiable data integrity rules such as positive quantity and non-blank workstation, while Go still validates early for better API errors.

#### Go Concepts

- `database/sql` as a stable standard-library abstraction for migration tooling
- blank imports for database driver registration
- defers and why `os.Exit` must not bypass cleanup
- generated code boundaries

#### Architecture Concepts

- SQL-first persistence design
- migration files as versioned database history
- sqlc-generated query package kept inside the production vertical slice
- dependency direction preserved: HTTP still does not depend directly on generated SQL code

#### Implementation

- Added `migrations/0001_create_production_entries.sql`.
- Added `sqlc.yaml`.
- Added `internal/production/queries/entries.sql`.
- Generated `internal/production/productiondb` with typed sqlc queries.
- Implemented `cmd/migrate` using `pgx` through `database/sql` and `goose`.
- Added `DATABASE_URL` and `MIGRATIONS_DIR` configuration.
- Added `make migrate` and `make sqlc` targets.

#### Tests

- Configuration tests cover database and migration environment variables.
- Generated sqlc code is compiled by `go test ./...` and `go build ./...`.
- Live migration execution passed against local PostgreSQL with `make migrate`.
- Verified the `production_entries` schema and goose version table through `psql`.

#### Refactoring

- Refactored `cmd/migrate` to return an exit code from `run()` instead of calling `os.Exit` before defers execute.

#### Code Review

- An experienced Go engineer would approve the direction: SQL is explicit, generated code is isolated and migrations are versioned.
- Remaining gap: the HTTP handler still uses the in-memory store. This is intentional; Lesson 4.3 will introduce a PostgreSQL-backed repository and wire context propagation through it.

#### Exercises

- Explain why `CHECK (quantity > 0)` belongs in the database even though Go also validates quantity.
- Inspect the `production_entries` table with `psql` and identify every database constraint.
- Modify `entries.sql`, rerun `make sqlc`, and inspect the generated type changes.

#### Interview Questions

- Why choose sqlc instead of an ORM?
- What is the role of migrations in production systems?
- Why does `database/sql` need a blank driver import?
- Why is it dangerous to call `os.Exit` before deferred cleanup runs?

#### Roadmap Update

- Lesson 4.2 completed.
- Current lesson moved to Lesson 4.3.
- `database/sql`, PostgreSQL, pgx and sqlc marked complete in the Knowledge Matrix.
