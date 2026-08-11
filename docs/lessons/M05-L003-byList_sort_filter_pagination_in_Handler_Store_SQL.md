### Lesson 5.3 Completion Notes

#### Business Context

As employee and product data grows, returning every row becomes inefficient and hard to use. Team leaders and administrators need predictable list endpoints that can page, filter and sort data.

#### Problem

Employee and product repositories returned all rows. Product search existed as a separate store method, while employees had no filtering. There was no consistent pagination or sort validation.

#### Design Discussion

The lesson uses explicit per-slice `ListOptions` instead of a generic pagination abstraction. This keeps allowed filters and sort keys close to the business capability. Employees and products both support `limit`, `offset`, `sort`, `q` and `active`, but each slice owns its own valid sort fields and query matching rules.

Sorting is whitelisted instead of interpolating SQL identifiers. PostgreSQL queries use static `CASE WHEN` ordering so user input never becomes SQL syntax. In-memory stores implement the same behavior for fast tests.

#### Go Concepts

- request query parsing with `strconv`
- option structs for explicit repository APIs
- in-memory filtering with `strings`
- deterministic sorting with `sort.Slice`
- defensive validation in handlers and stores

#### Architecture Concepts

- repository API designed around query intent
- filtering/sorting rules owned by vertical slices
- SQL injection prevention through whitelisted sort values
- one list/query path instead of separate product search persistence methods

#### Implementation

- Added employee and product `ListOptions` and `Page` types.
- Added validated `limit`, `offset`, `sort`, `q` and `active` query parameters.
- Updated employee and product `Store.List` contracts to accept options.
- Removed the separate product store-level `Search` method and reused `List` with `Query`.
- Updated PostgreSQL sqlc list queries with filtering, sorting, `LIMIT` and `OFFSET`.
- Updated in-memory stores to match PostgreSQL filtering and sorting behavior.
- Added pagination metadata to list responses.

#### Tests

- Added in-memory store tests for filtering, sorting and pagination.
- Added HTTP handler tests for list query options and invalid query options.
- Updated PostgreSQL repository tests to use the new list contract.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

Product search now delegates to the same list mechanism as `GET /products`. This removes a second persistence method and keeps search behavior consistent with pagination and sorting.

#### Code Review

An experienced Go engineer would approve the main design because user-controlled sort input is whitelisted and never interpolated into SQL. The main follow-up is API documentation quality: Fuego generates the endpoints, but explicit query-parameter documentation should be reviewed later.

#### Exercises

- Add a repository test proving `offset` beyond the result size returns an empty non-nil slice.
- Add a new employee sort key and update both in-memory and PostgreSQL implementations.
- Explain why dynamic SQL string concatenation would be risky for sort fields.

#### Interview Questions

- Why is offset pagination simple but not always ideal for large datasets?
- How do you safely implement dynamic sorting in SQL?
- What is the difference between filtering in Go and filtering in SQL?
- When would cursor pagination be better than limit/offset pagination?

#### Roadmap Update

- Lesson 5.3 completed.
- Current lesson moved to Lesson 5.4.
- Known technical debt updated: pagination/filtering/sorting completed; optimistic locking and foreign-key-backed consistency remain pending.
