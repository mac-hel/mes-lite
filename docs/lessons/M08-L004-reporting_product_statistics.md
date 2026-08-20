### Lesson 8.4 Scope

Add product-level statistics reporting for completed production over a requested time range.

#### Business Context

Production managers need to know which products were produced most often, how much quantity was completed and how broadly employees contributed to each product.

#### Problem

Existing reports answered daily totals and employee productivity, but not product-level performance. Managers still had to manually aggregate raw production entries to identify high-volume products.

#### Design Discussion

The reporting slice now owns a third read model: `ProductStatisticsRow`. The SQL query joins `production_entries` with `products` so the report includes product names directly in the response.

The report includes total quantity, entry count and distinct employee count. This keeps the first product statistics API useful while avoiding premature analytics such as averages, trend lines or percentages.

#### Go Concepts

- adding explicit read models instead of generic maps
- extending test fixtures without weakening type safety
- converting generated sqlc rows into API-owned DTOs
- preserving consistent query-parameter parsing across handlers

#### Architecture Concepts

- reporting as a cross-slice read model owner
- SQL aggregation with `COUNT(DISTINCT ...)`
- deterministic report ordering
- avoiding generic report abstractions before repeated pain exists

### Lesson 8.4 Completion Notes

#### Business Context

MES Lite now exposes product statistics for management users.

#### Problem

Managers could view production by day and by employee, but not by product. Product-level output is a core reporting need for understanding what the factory actually produced during a period.

#### Design Discussion

Added `ProductStatistics` to the reporting store and exposed it through `GET /reports/product-statistics`. The query aggregates completed production by product and joins product master data for display names.

The SQL orders by total quantity descending, then entry count descending, then product SKU ascending. This produces useful ranking and stable test/API output.

#### Implementation

- Added `ProductStatisticsRow` read model.
- Added `ProductStatistics` to `reporting.Store`.
- Added sqlc query joining `production_entries` with `products`.
- Generated updated `reportingdb` code.
- Added PostgreSQL store mapping for product statistics rows.
- Added `GET /reports/product-statistics`.
- Added response DTOs for product statistics.
- Updated the reporting in-memory store fixture for all report types.
- Registered the new route with bearer security and management RBAC.

#### Tests

- Added PostgreSQL integration test for product statistics grouping and ordering.
- Added PostgreSQL integration test for invalid product statistics ranges.
- Added handler test for product statistics response mapping.
- Added handler test for invalid product statistics ranges.
- Added server route test proving managers can read product statistics reports.
- Added server route test proving workers cannot read product statistics reports.
- Verified with `sqlc generate`.
- Verified with `go test ./internal/reporting -count=1`.
- Verified with `go test ./internal/server -count=1`.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

No generic report abstraction was introduced. The reporting package now has three explicit report methods, which is still clear and easier to review than a configurable aggregation layer.

#### Code Review

An experienced Go engineer would approve the explicit read-model approach and deterministic ordering. The main follow-up for Lesson 8.5 is query-performance review: decide whether indexes are needed for report time ranges and grouping.

#### Exercises

- Add a test proving products with no production entries are omitted from the report.
- Add a tie-breaker test where two products have equal total quantity and entry count.
- Discuss whether inactive products should remain visible in historical reports.

#### Interview Questions

- Why use `COUNT(DISTINCT employee_id)` in product statistics?
- Why can reporting join product master data directly?
- When would product statistics need pagination or export streaming?
- What indexes might help this report as production_entries grows?

#### Roadmap Update

- Lesson 8.4 completed.
- Current lesson moved to Lesson 8.5.
