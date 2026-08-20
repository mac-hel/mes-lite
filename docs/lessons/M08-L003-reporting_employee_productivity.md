### Lesson 8.3 Scope

Add an employee productivity report that ranks employees by completed production quantity for a requested time range.

#### Business Context

Production managers and team leaders need to understand employee output without manually joining production entries to employee master data. The report should answer who produced how much during a period.

#### Problem

The daily production report groups by day and product, but it does not show employee-level productivity. Managers still cannot compare employee output or identify who contributed to a period's production totals.

#### Design Discussion

The reporting slice now owns a second read model: `EmployeeProductivityRow`. The SQL query joins `production_entries` with `employees` so the report can return employee IDs and names without asking clients to perform extra lookups.

The API uses the same `from`/`to` RFC3339 half-open time range as daily production. This keeps report endpoints consistent and avoids inventing per-report parameter conventions.

#### Go Concepts

- extending small interfaces when a consumer actually needs another behavior
- reusing query-parameter parsing through a helper function
- DTO conversion for read models with matching field sets
- defensive in-memory test store copies

#### Architecture Concepts

- reporting query store growing by report use case
- SQL join for read-model convenience
- consistent report API contracts
- route-level RBAC reused for management reporting

### Lesson 8.3 Completion Notes

#### Business Context

MES Lite now exposes employee productivity reporting for management users.

#### Problem

The reporting API could show production by day and product, but not by employee. This left a core management question unanswered: who produced how much in a selected time range?

#### Design Discussion

Added `EmployeeProductivity` to the reporting store instead of adding employee productivity logic to the employees or production slices. This keeps reporting as the owner of read-model aggregation while employees remains master data and production remains the write workflow.

The query orders by total quantity descending, then entry count descending, then employee ID ascending for deterministic results.

#### Implementation

- Added `EmployeeProductivityRow` read model.
- Added `EmployeeProductivity` to `reporting.Store`.
- Added sqlc query joining `production_entries` with `employees`.
- Generated updated `reportingdb` code.
- Added PostgreSQL store mapping for employee productivity rows.
- Added `GET /reports/employee-productivity`.
- Added response DTOs for employee productivity.
- Updated the reporting in-memory store for multi-report tests.
- Registered the new route with bearer security and management RBAC.

#### Tests

- Added PostgreSQL integration test for employee productivity grouping and ordering.
- Added PostgreSQL integration test for invalid employee productivity ranges.
- Added handler test for employee productivity response mapping.
- Added handler test for invalid employee productivity ranges.
- Added server route test proving leaders can read employee productivity reports.
- Added server route test proving workers cannot read employee productivity reports.
- Verified with `sqlc generate`.
- Verified with `go test ./internal/reporting -count=1`.
- Verified with `go test ./internal/server -count=1`.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

The handler now uses a shared `reportRange` helper so report endpoints parse `from` and `to` consistently. No generic report framework was introduced.

#### Code Review

An experienced Go engineer would approve the direction because each report remains explicit, SQL is readable and report-specific DTOs keep API shape separate from generated database rows. The main follow-up is product statistics in Lesson 8.4, which should reuse the reporting slice without collapsing into a generic aggregation abstraction.

#### Exercises

- Add a test proving employees with no production entries are omitted from the report.
- Add a tie-breaker test where two employees have the same total quantity.
- Discuss whether inactive employees should appear in historical productivity reports.

#### Interview Questions

- Why is it acceptable for reporting SQL to join across business slices?
- Why should read models include display names instead of forcing clients to join data?
- How do deterministic `ORDER BY` clauses improve tests and APIs?
- When would this report need pagination?

#### Roadmap Update

- Lesson 8.3 completed.
- Current lesson moved to Lesson 8.4.
