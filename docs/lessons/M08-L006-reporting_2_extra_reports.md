### Lesson 8.6 Scope

Add detailed breakdown reports that preserve the existing summary endpoints while answering product-by-employee questions.

#### Business Context

Managers need both summary and detail. A daily product total answers what was produced each day, but not who produced each product. Employee productivity answers total output per employee, but not which products each employee produced.

#### Problem

Changing `/reports/daily-production` or `/reports/employee-productivity` directly would change their aggregation level and could make API consumers double-count totals. The missing information should be exposed through explicit detail reports.

#### Design Discussion

Add two new report endpoints:

- `GET /reports/daily-employee-production`
- `GET /reports/employee-productivity/products`

`/reports/daily-employee-production` groups by day, product and employee. It answers: who produced how much of each product on each day?

`/reports/employee-productivity/products` groups by employee and product. It answers: what product mix did each employee produce during the selected period?

The existing summary reports remain unchanged.

#### Go Concepts

- adding new read models without breaking existing API contracts
- explicit DTOs for different aggregation levels
- careful naming to avoid ambiguous report semantics
- extending sqlc query packages while keeping generated rows internal

#### Architecture Concepts

- preserving API compatibility by adding endpoints instead of changing row meaning
- report granularity as part of API design
- CQRS read models for different projections over the same source data

#### Implementation Plan

- Add `DailyEmployeeProductionRow` read model.
- Add `EmployeeProductivityProductRow` read model.
- Add sqlc queries for both detailed aggregations.
- Add PostgreSQL store methods and in-memory test support.
- Add handler DTOs and routes.
- Register both endpoints with existing management-report RBAC.
- Add integration, handler and server authorization tests.

#### Tests

- Query correctness for daily employee production grouping by day/product/employee.
- Query correctness for employee product mix grouping by employee/product.
- Invalid range tests for both reports.
- Handler response mapping tests.
- RBAC tests proving management roles can read and workers cannot.

#### Exercises

- Explain why changing the existing summary report row shape would be a breaking API change.
- Add a test proving totals from the detailed report sum back to the existing summary report.
- Decide whether the detailed reports need pagination before CSV export exists.

#### Interview Questions

- Why is aggregation level part of an API contract?
- How can detailed reports cause accidental double-counting?
- When should a report use a separate endpoint instead of an optional `groupBy` parameter?
- How do CQRS read models help preserve command-model simplicity?

### Lesson 8.6 Completion Notes

#### Business Context

MES Lite now has detailed product-by-employee reporting without changing the existing summary reports.

#### Problem

The daily production report showed product totals per day, and employee productivity showed employee totals for a period. Neither report answered who produced which product at the detailed level.

#### Design Discussion

Added two explicit detail endpoints instead of changing existing row shapes. This preserves summary-report semantics and avoids accidental double-counting by API clients.

`/reports/daily-employee-production` groups by day, product and employee. `/reports/employee-productivity/products` groups by employee and product.

#### Implementation

- Added `DailyEmployeeProductionRow` read model.
- Added `EmployeeProductivityProductRow` read model.
- Added sqlc queries for both detailed aggregations.
- Generated updated `reportingdb` code.
- Added PostgreSQL store methods for both reports.
- Added in-memory store support for detailed report tests.
- Added `GET /reports/daily-employee-production`.
- Added `GET /reports/employee-productivity/products`.
- Registered both routes with existing management-report RBAC.

#### Tests

- Added PostgreSQL integration test for daily employee production grouping.
- Added PostgreSQL integration test for employee productivity by product grouping.
- Added invalid range tests for both detailed store methods.
- Added handler response mapping tests for both detailed reports.
- Added handler invalid range tests for both detailed reports.
- Added server RBAC tests proving management roles can read and workers cannot.
- Verified with `go fmt ./internal/reporting ./internal/server`.
- Verified with `sqlc generate`.
- Verified with `go test ./internal/reporting -count=1`.
- Verified with `go test ./internal/server -count=1`.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

The existing summary endpoints remain unchanged. The reporting package is still explicit rather than generic, which keeps each aggregation level reviewable.

#### Code Review

An experienced Go engineer would approve adding separate endpoints because aggregation level is part of the API contract. The main follow-up is still OpenAPI query-parameter documentation quality, which remains known technical debt.

#### Exercises

- Add a test proving detailed daily employee totals sum back to daily product totals.
- Add a test proving employee product totals sum back to employee productivity totals.
- Decide whether detailed report endpoints should support pagination before CSV export.

#### Interview Questions

- Why is changing aggregation level a breaking API change?
- How do detailed and summary reports complement each other?
- When would an optional `groupBy` parameter be better than separate endpoints?
- Why should generated sqlc row types stay out of HTTP responses?

#### Roadmap Update

- Lesson 8.6 completed.
- Milestone 8 completed again after detailed report refinement.
- Current milestone moved to Milestone 9.
- Current lesson moved to Lesson 9.1.
