### Lesson 10.5 Scope

Refactor employee and product HTTP responses to explicit DTOs without changing route behavior or JSON field names.

#### Business Context

Employees and products are part of the MVP API contract. Their response shapes should be stable for future clients before the project moves into post-MVP concurrency work.

#### Problem

The production, orders, reporting, auth and CSV import slices now use explicit response DTOs where appropriate. Employee and product handlers still return domain structs directly, so future domain changes could accidentally change the public API.

#### Design Discussion

This is a cleanup lesson, not a feature lesson. It should preserve existing behavior, authorization, persistence and JSON field names while making the HTTP boundary explicit.

The lesson should avoid broad domain encapsulation refactors. It should only introduce response DTOs and mapping helpers for employee/product handlers.

#### Go Concepts

- DTO mapping functions
- preserving API compatibility during refactoring
- avoiding package-name stutter in exported response types
- small focused cleanup commits

#### Architecture Concepts

- HTTP contracts separated from domain structs
- API boundary consistency across vertical slices
- milestone cleanup before post-MVP concurrency work

#### Implementation Plan

- Add employee response DTOs.
- Add product response DTOs.
- Update employee/product create, list, update and deactivate handlers to return DTOs.
- Keep request DTOs unchanged.
- Preserve JSON field names.
- Update handler/server tests only where they decode concrete response types.

#### Tests

- Existing employee/product handler tests should pass.
- Existing server route tests should pass.
- Add or update contract assertions only where useful.

#### Exercises

- Compare employee/product response DTOs with production `EntryResponse`.
- Explain why response DTOs are useful even when they initially match domain fields exactly.
- Decide whether future employee/product domain fields should be hidden from API responses by default.

#### Interview Questions

- Why should API contracts not depend directly on domain struct fields?
- When is a DTO refactor worth the churn?
- How do you preserve compatibility while changing handler return types?
- Why avoid refactoring domain encapsulation at the same time?

### Lesson 10.5 Completion Notes

#### Business Context

Employee and product API responses now have explicit HTTP contracts before the project moves into post-MVP concurrency work.

#### Problem

Employee and product handlers returned domain structs directly. That made the public API depend on domain field shape.

#### Design Discussion

Added response DTOs while preserving route behavior and JSON field names. The existing API used capitalized field names because the domain structs had no JSON tags, so L10.5 kept those names to avoid a breaking cleanup.

#### Implementation

- Added `employees.EmployeeResponse`.
- Added `products.ProductResponse`.
- Updated employee create, list, update and deactivate handlers to return DTOs.
- Updated product create, list, update, deactivate and search handlers to return DTOs.
- Added mapping helpers for employee and product responses.
- Preserved request DTOs and JSON field names.

#### Tests

- Updated employee handler tests to decode `EmployeeResponse`.
- Updated product handler tests to decode `ProductResponse`.
- Verified focused employee/product/server tests.

#### Refactoring

The cleanup did not change domain encapsulation or persistence. It only separated HTTP response contracts from domain structs.

#### Code Review

An experienced Go engineer would approve this as a focused boundary cleanup. The main caveat is that the preserved capitalized JSON field names are not ideal for a new API, but changing them would be a breaking API contract change and should be handled deliberately if ever needed.

#### Exercises

- Add JSON contract tests for employee and product response field names.
- Compare preserving capitalized JSON fields with migrating to lower camel case.
- Identify which older slices still expose domain types and whether changing them is worth the churn.

#### Interview Questions

- Why can a cleanup preserve an imperfect API contract instead of improving field names?
- Why are DTOs useful even when they are structurally identical to domain structs?
- How do you avoid mixing API cleanup with domain-model refactoring?
- What makes an API response change breaking?

#### Roadmap Update

- Lesson 10.5 completed.
- Milestone 10 completed.
- Current milestone moved to Milestone 11.
- Current lesson moved to Lesson 11.1.
