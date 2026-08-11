### Lesson M6-pre - Test DB & PostgreSQL Error Codes

### Goal

Resolve small persistence-quality issues before adding authentication and authorization.

### Scope

- introduce a separate integration test database configuration, preferably `TEST_DATABASE_URL`
- ensure tests never use the same database as the running application or future production deployments
- keep `DATABASE_URL` for the application and migration command
- define shared string-typed PostgreSQL SQLSTATE constants
- replace raw PostgreSQL error code strings in employee, product and production PostgreSQL stores

### Business Value

Authentication will introduce more security-sensitive flows. Before that, tests should be isolated from application data and persistence error handling should be easier to read and review.

### Design Notes

Integration tests currently use `DATABASE_URL` and clean tables directly. This is acceptable for early local development but should not continue once the project moves toward production-like security work.

PostgreSQL error codes are standardized SQLSTATE strings. They should remain string values, not `iota` enums, because PostgreSQL returns strings such as `23505` and `23503`.

Preferred shape:

```go
type SQLState string

const (
    UniqueViolation     SQLState = "23505"
    ForeignKeyViolation SQLState = "23503"
    CheckViolation      SQLState = "23514"
    NotNullViolation    SQLState = "23502"
    InvalidTextValue    SQLState = "22P02"
)
```

### Definition of Done

- ✅ tests use `TEST_DATABASE_URL` or skip when it is not configured/available
- ✅ application runtime still uses `DATABASE_URL`
- ✅ PostgreSQL error code magic strings are removed from `internal/employees`, `internal/products` and `internal/production`
- ✅ shared constants live in a small internal package
- ✅ tests pass
- ✅ build passes
- ✅ lint passes
- ✅ roadmap updated

### Completion Review

Integration tests for PostgreSQL stores now read `TEST_DATABASE_URL` directly and skip when it is not configured. Application configuration remains unchanged and continues to use `DATABASE_URL`.

PostgreSQL SQLSTATE codes now live in `internal/postgres` as string-typed constants. Employee, product and production PostgreSQL stores map errors through those constants instead of raw strings.
