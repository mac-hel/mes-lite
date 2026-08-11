### Lesson 5.2 Scope

Standardize validation across the whole request-to-database flow and refactor domain constructors so invalid entities are harder to create.

#### Business Context

MES Lite stores operational data that managers and workers rely on. Invalid employee, product or production data should be rejected consistently regardless of whether it enters through HTTP, tests, future imports, background jobs or direct repository calls.

#### Problem

Validation currently exists in multiple places but the flow is not documented. Some constructors create values first and require callers to remember a separate `Validate()` call. This makes invalid domain state possible inside the application.

#### Design Discussion

Validation should be layered instead of duplicated randomly. HTTP handlers validate transport shape. Application services validate workflow rules and references. Domain constructors and mutation methods enforce invariants. Repositories validate before persistence as defense-in-depth. PostgreSQL constraints provide the final integrity boundary.

The lesson begins by documenting the convention in `docs/validation.md`, then refactors constructors and call sites incrementally.

#### Go Concepts

- constructors returning `(T, error)`
- invalid zero values for business entities
- validation methods used internally by constructors and mutation methods
- error wrapping for invariant failures

#### Architecture Concepts

- layered validation responsibility
- domain invariants vs transport validation
- database constraints as final integrity boundary
- consistency rules for future slices and contributors

### Lesson 5.2 Completion Notes

#### Business Context

MES Lite stores production-critical data. Invalid master data or production entries should be rejected consistently regardless of whether data enters through HTTP, tests, repositories, imports or future background jobs.

#### Problem

Domain constructors could create invalid values and relied on callers to remember separate validation calls. Validation responsibilities were also implicit, making future inconsistency likely.

#### Design Discussion

Validation is now documented as a layered flow. Handlers validate transport shape. Application services validate workflows and references. Domain constructors enforce invariants. Repositories validate defensively before persistence. PostgreSQL constraints remain the final integrity boundary.

This does not make invalid state impossible in Go because fields are still exported for API serialization and test readability. It does make the normal construction path safe and documents how future code should behave.

#### Go Concepts

- constructors returning `(T, error)`
- invalid zero values for business entities
- mutation methods that preserve invariants by validating a copy before assignment
- sentinel errors wrapped with domain-specific context

#### Architecture Concepts

- documented validation ownership across layers
- domain invariants separated from HTTP validation tags
- repository validation as defense-in-depth
- database constraints as the final safety net

#### Implementation

- Added `docs/validation.md` as the project validation guideline.
- Changed `employees.NewEmployee` to return `(Employee, error)` and validate required fields.
- Added `Employee.Validate` and `Employee.UpdateDetails`.
- Changed `products.NewProduct` to return `(Product, error)`.
- Changed `Product.UpdateDetails` to return an error and preserve the previous value on invalid input.
- Changed `production.NewEntry` to return `(Entry, error)`.
- Updated handlers, services, stores and tests to handle constructor errors explicitly.
- Added defensive employee validation in in-memory and PostgreSQL stores.

#### Tests

- Added constructor rejection tests for employees, products and production entries.
- Added mutation preservation tests for employee and product updates.
- Updated repository and HTTP fixtures to fail fast on invalid test data.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

The refactor keeps public struct fields for now because existing HTTP serialization and tests rely on them. A future stricter domain model could hide fields behind accessor methods, but that would be a larger API-design change and is not needed yet.

#### Code Review

An experienced Go engineer would approve this as a clear improvement: invalid construction is harder, validation ownership is documented and repositories still provide a storage boundary. The main caveat is that exported fields still allow manual invalid mutation; this is acceptable for the current project maturity but should be revisited if domain invariants become more complex.

#### Exercises

- Explain why `NewProduct` now returns `(Product, error)` instead of `Product`.
- Add a table test for `Employee.Validate` covering every required field.
- Try to mutate a product with invalid details and explain why the old value is preserved.

#### Interview Questions

- When is it acceptable for a Go type to have an invalid zero value?
- Why should HTTP validation not replace domain validation?
- Why should repositories still validate if constructors already validate?
- What belongs in database constraints versus application code?

#### Roadmap Update

- Lesson 5.2 completed.
- Current lesson moved to Lesson 5.3.
- Known technical debt updated: constructor validation inconsistency resolved; pagination/filtering/sorting, optimistic locking and foreign-key-backed consistency remain pending.
