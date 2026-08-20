### Lesson 8.1 Scope

Introduce reporting as a read-model/query slice before adding HTTP endpoints.

#### Business Context

Managers need production totals without manually filtering Excel sheets. Reporting reads existing operational data and summarizes it for decisions rather than changing business state.

#### Problem

The application stores production entries, employees, products and orders, but every existing slice is command/workflow oriented. Adding reports directly to the production slice would mix write-model behavior with read-specific SQL aggregation.

#### Design Discussion

Reporting starts as its own vertical slice with a PostgreSQL read store. It returns read models, not production-entry aggregates. This is a small CQRS step: commands still live in their business slices, while reporting owns query shapes optimized for management views.

The first query groups production entries by UTC day and product SKU. HTTP endpoints, response DTOs and authorization are postponed to Lesson 8.2 so this lesson can focus on SQL aggregate correctness and package boundaries.

#### Go Concepts

- read-model structs separate from domain entities
- context propagation into query methods
- time-range validation with sentinel errors
- integer conversion at SQL aggregate boundaries

#### Architecture Concepts

- CQRS as separate read models, not a new framework
- reporting vertical slice owns reporting SQL
- query store returns projections instead of aggregates
- ADR for introducing reporting read models

### Lesson 8.1 Completion Notes

#### Business Context

MES Lite now has a reporting foundation for manager-facing production summaries.

#### Problem

Operational data existed, but there was no dedicated query boundary for reports. Reusing production repositories would have forced report-shaped aggregate SQL into a write-focused slice.

#### Design Discussion

Added `internal/reporting` with a small read store. The store exposes `DailyProduction`, which returns `DailyProductionRow` values grouped by UTC day and product SKU.

This is CQRS in a pragmatic Go style: no bus, no framework and no generic query abstraction. The separation exists because reports have different data shapes and SQL needs than commands.

#### Implementation

- Added `reporting.Store` and `DailyProductionRow` read model.
- Added time-range validation with `ErrInvalidRange`.
- Added reporting sqlc configuration.
- Added `internal/reporting/queries/reports.sql`.
- Generated `internal/reporting/reportingdb`.
- Added `reporting.PostgresStore`.
- Added ADR `0003-introduce-reporting-read-models.md`.

#### Tests

- Added unit tests for report range validation.
- Added PostgreSQL integration test for daily production grouping by day and product.
- Added PostgreSQL integration test for invalid ranges.
- Verified with `sqlc generate`.
- Verified with `go test ./internal/reporting -count=1`.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

No existing write-model repository was changed. Reporting reads from existing tables through its own query package.

#### Code Review

An experienced Go engineer would approve the small CQRS introduction because it solves a real reporting shape without adding a framework or generic abstraction. The main follow-up is HTTP/API design in Lesson 8.2, including query parameters and authorization.

#### Exercises

- Add a test proving an empty date range returns an empty non-nil slice.
- Explain why `DailyProductionRow` is not a `production.Entry`.
- Add employee name to the query and discuss whether that belongs in L8.2 or L8.3.

#### Interview Questions

- What is CQRS, and what problem does it solve here?
- Why can read models differ from domain aggregates?
- Why group by UTC day instead of local server time?
- What are the risks of putting reporting SQL into command repositories?

#### Roadmap Update

- Lesson 8.1 completed.
- Current lesson moved to Lesson 8.2.
- Architecture `CQRS` marked complete in the Knowledge Matrix.
