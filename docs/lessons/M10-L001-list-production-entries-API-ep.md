### Lesson 10.1 Scope

Expose production-entry review over HTTP so managers and leaders can inspect entered production records without using direct database access.

#### Business Context

Team leaders and managers need to review recently entered production records, filter them by the same fields they use in manual spreadsheets and identify mistakes before correction workflows are added.

#### Problem

Workers could register production entries, and reports could aggregate production history, but there was no review endpoint for the raw entered records. Reporting summaries are not enough when a leader needs to inspect individual entries.

#### Design Discussion

The review API is intentionally a list endpoint over the existing production-entry model. It supports `employeeId`, `productSku`, `workstation`, `from`, `to`, `limit` and `offset` query parameters.

The time range is half-open: `from <= timestamp < to`. This matches the reporting API convention and avoids double-counting when clients review adjacent periods.

Workers can still register production, but only admins, managers and leaders can review historical entries. Correction workflows are postponed to L10.3.

#### Go Concepts

- query-parameter parsing with `strconv` and `time.Parse`
- option structs for explicit list filters
- nil-slice normalization for stable JSON responses
- store interface growth when the consumer has a real read use case

#### Architecture Concepts

- review API as a read workflow over production entries
- route-level RBAC for historical production review
- filtering rules owned by the production vertical slice
- no workstation aggregate before the business needs formal workstation management

### Lesson 10.1 Completion Notes

#### Business Context

MES Lite now lets authorized users review entered production records directly through the API.

#### Problem

Production registration existed, but managers and leaders could not inspect individual production entries without relying on database access or summary reports.

#### Design Discussion

Added `GET /production-entries` with review-focused filters. The endpoint returns production entries plus pagination metadata and keeps workstation as a simple text field, matching MVP scope.

The handler parses HTTP query parameters, the production service delegates the read operation and the store owns filtering. This keeps HTTP translation, application workflow and persistence responsibilities separated without adding a new abstraction.

#### Implementation

- Added `production.ListOptions` and `production.Page`.
- Added `Store.List` to the production store contract.
- Implemented filtered, paginated listing in the in-memory production store.
- Updated the production sqlc `ListEntries` query with review filters and pagination.
- Implemented filtered listing in `production.PostgresStore`.
- Added `GET /production-entries` handler response with entries and pagination.
- Registered the route with bearer auth and admin/manager/leader RBAC.

#### Tests

- Added handler test for filtered review listing.
- Added handler test for invalid query parameters.
- Added PostgreSQL store test for filtered, paginated listing.
- Added server route test proving leaders can review production entries.
- Added server route test proving workers cannot review historical production entries.
- Verified with `sqlc generate`.
- Verified with `go test ./internal/production ./internal/server -count=1`.

#### Refactoring

No correction or audit model was introduced. L10.1 stays focused on review/read behavior so L10.2 can address idempotent registration and L10.3 can address append-only corrections.

#### Code Review

An experienced Go engineer would approve the small API addition because it follows existing list/filter patterns, uses explicit options and keeps authorization at route composition. The main remaining gap is expected milestone scope: registration is not idempotent yet, and historical corrections are not audit-safe yet.

#### Exercises

- Add a test proving adjacent `from`/`to` review ranges do not duplicate boundary entries.
- Add an OpenAPI inspection test or manual check for the new query parameters.
- Decide whether production-entry review should eventually include employee/product display names or keep returning IDs only.

#### Interview Questions

- Why use an option struct instead of passing many list parameters separately?
- Why is a half-open time range safer than an inclusive `to` timestamp?
- Why should workers be allowed to create entries but not review all historical entries?
- When would this endpoint need cursor pagination instead of limit/offset?

#### Roadmap Update

- Lesson 10.1 completed.
- Current lesson moved to Lesson 10.2.
- Milestone 10 status moved to In Progress.
