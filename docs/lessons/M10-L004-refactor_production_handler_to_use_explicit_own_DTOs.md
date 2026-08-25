### Lesson 10.4 Scope

Review the MVP API for correctness, package boundaries, API contract quality and readiness before moving into post-MVP concurrency lessons.

#### Business Context

Milestone 10 closes the remaining backend API gaps required for the Excel/paper-replacement MVP. Before adding background jobs and concurrency, the API should be coherent, reviewable and safe enough for the current business scope.

#### Problem

The production API had completed review, idempotent registration and correction workflows, but one API-design issue remained from review: production HTTP responses returned domain types directly.

#### Design Discussion

L10.4 keeps behavior unchanged and refactors the production HTTP boundary to use explicit response DTOs. This makes the JSON contract intentional instead of depending on domain struct shape.

The milestone review accepts separate correction-history endpoints for now. An effective/current production-entry read model may be useful later, but the MVP already exposes the original entry and append-only corrections needed for auditability.

#### Go Concepts

- DTO mapping functions
- preserving API behavior during refactoring
- small interface review
- milestone-level code review discipline

#### Architecture Concepts

- HTTP contract separated from domain structs
- MVP scope closure before adding post-MVP concurrency
- correction history as explicit read data instead of silent mutation
- review-driven refactoring

### Lesson 10.4 Completion Notes

#### Business Context

The MVP API now covers production-entry review, retry-safe registration and audit-safe correction workflows.

#### Problem

Production response types were coupled to domain structs. If domain fields changed later, the HTTP API could change accidentally.

#### Design Discussion

Added explicit production response DTOs while keeping request DTOs, routes and JSON field names stable. This resolves the response-coupling issue identified during the L10.3 review.

No broad DTO refactor was applied to all older slices. L10.4 focused on the production MVP API because that is the current milestone boundary.

#### Implementation

- Added `ProductionEntryResponse`.
- Added `CorrectionResponse`.
- Changed production registration and correction handlers to return response DTOs.
- Changed production-entry and correction list responses to contain DTO slices.
- Added mapping helpers from `Entry` and `Correction` to response DTOs.
- Preserved existing route names and JSON fields.

#### Tests

- Existing production handler tests pass against the stable JSON response contract.
- Existing server route tests pass for registration, review, idempotency and corrections.
- Verified with `go test ./internal/production ./internal/server -count=1`.

#### Refactoring

The production HTTP boundary no longer exposes `Entry` or `Correction` directly in handler return types. This keeps future domain changes from accidentally becoming API changes.

#### Code Review

An experienced Go engineer would approve the MVP API shape for the current scope. The API remains small, protected by RBAC, backed by PostgreSQL constraints and tested through handler, service and repository tests.

The main caveat remains OpenAPI metadata quality for query parameters. Endpoints are generated, but explicit query-parameter documentation should still be reviewed later.

#### Exercises

- Add a contract test that asserts production-entry response JSON field names.
- Design an effective production-entry read model that combines original entry plus latest correction.
- Compare response DTOs in production with older employee/product handlers and decide whether a broad API DTO refactor is worth the churn.

#### Interview Questions

- Why should HTTP response DTOs be separate from domain structs?
- How do you decide whether a review finding needs immediate refactoring?
- Why close MVP scope before introducing background jobs?
- What risks remain when OpenAPI query parameters are not explicitly documented?

#### Roadmap Update

- Lesson 10.4 completed.
- Follow-up Lesson 10.5 added for employee/product response DTO cleanup before Milestone 10 closes.
- Current lesson moved to Lesson 10.5.
