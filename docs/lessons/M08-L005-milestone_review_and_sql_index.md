### Lesson 8.5 Scope

Review the reporting slice and add the first query-performance guardrail before closing Milestone 8.

#### Business Context

Reports will become slower as production history grows. Managers need report endpoints to remain predictable without turning every request into a full-table scan over production history.

#### Problem

The reporting queries were correct and tested, but `production_entries` had no reporting-oriented index. All three reports filter by `occurred_at` before grouping by product or employee.

#### Design Discussion

Added one targeted covering index on `production_entries` for the reporting access pattern: `(occurred_at, product_sku, employee_id) INCLUDE (quantity)`. This supports the shared time-range filter and keeps product, employee and quantity data available from the index for aggregation-heavy reports.

This is intentionally one index, not several speculative indexes. Indexes speed reads but slow writes and consume storage. More indexes should be added only after real query plans or production-like data show a need.

#### Go Concepts

- integration tests that verify migration side effects
- keeping performance behavior explicit through schema migrations
- reviewing abstractions before adding more code

#### SQL Concepts

- B-tree indexes for range predicates
- covering indexes with `INCLUDE`
- trade-offs between read speed, write cost and storage
- query-performance review before premature optimization

### Lesson 8.5 Completion Notes

#### Business Context

Milestone 8 now has useful management reports and a first database-level performance guardrail for report time-range queries.

#### Problem

Reports were implemented, but there was no index supporting their common access pattern. As production entries grow, report queries would increasingly depend on scanning the whole table.

#### Design Discussion

The three report queries share the same shape: filter production entries by `occurred_at`, then aggregate by product or employee. A single covering index is the smallest useful improvement.

The lesson did not introduce caching, materialized views or background report generation. Those are valid future tools, but they would be premature before measuring real query cost on larger data.

#### Implementation

- Added migration `0008_add_reporting_indexes.sql`.
- Added `production_entries_reporting_idx` on `(occurred_at, product_sku, employee_id) INCLUDE (quantity)`.
- Added integration test proving the reporting index exists after migrations.

#### Tests

- Verified with `go fmt ./internal/reporting`.
- Verified with `sqlc generate`.
- Verified with `go test ./internal/reporting -count=1`.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

No report abstraction was added. Three explicit report methods remain clearer than a generic reporting engine.

#### Code Review

An experienced Go engineer would approve the reporting slice for this milestone: read models are explicit, SQL is readable, generated sqlc code stays below the package boundary, routes are protected and tests cover query correctness and authorization.

The main caveat is that performance was improved with a reasoned index, not proven with production-scale `EXPLAIN ANALYZE`. That is acceptable at this project stage, but future performance work should use realistic row counts.

#### Exercises

- Run `EXPLAIN ANALYZE` for each report query before and after adding sample data.
- Compare one covering index with separate `(occurred_at, employee_id)` and `(occurred_at, product_sku)` indexes.
- Explain how additional indexes affect production-entry insert performance.

#### Interview Questions

- Why does every index have a write-time cost?
- What is a covering index?
- Why should performance optimization be measurement-driven?
- When would a materialized view be better than querying base tables directly?

#### Roadmap Update

- Lesson 8.5 completed.
- Milestone 8 completed.
- Current milestone moved to Milestone 9.
- Current lesson moved to Lesson 9.1.
- Persistence `SQL Optimization` and SQL `Indexes` marked complete in the Knowledge Matrix.
