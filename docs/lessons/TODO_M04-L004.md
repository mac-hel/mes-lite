### Lesson 4.4 Completion Notes

#### Business Context

Production workers must not register work for unknown or inactive employees/products. Without this rule, reports would contain production that cannot be assigned to valid business records.

#### Problem

The endpoint persisted production entries, but it only validated entry shape. It did not validate whether referenced employees/products existed or were active.

#### Design Discussion

The production slice now has an application service that coordinates business validation and persistence. The handler parses HTTP and translates errors. The service validates references and calls the store. The store persists the entry.

The transaction-boundary decision is explicit: today there is only one PostgreSQL write for production entries, while employees/products are still in-memory. A database transaction would not make this cross-resource validation atomic yet. Full transactional consistency requires moving employees/products to PostgreSQL in Milestone 5 and then validating references using database constraints or a single transaction.

#### Go Concepts

- consumer-owned interfaces for employee/product lookups
- error translation across HTTP, service and persistence boundaries
- direct struct conversion when request and command shapes match
- context propagation through handler, service and repository

#### Architecture Concepts

- application service as a business coordination boundary
- transaction boundary belongs around a complete business operation, not around arbitrary function calls
- explicit technical debt when consistency cannot yet be guaranteed by the current persistence model

#### Implementation

- Added `production.Service` and `RegisterCommand`.
- Added employee/product lookup interfaces owned by the production consumer.
- Added business errors for missing/inactive employees and products.
- Updated the handler to delegate registration to the service.
- Updated the server composition root to share employee/product stores with production validation.
- Preserved the PostgreSQL-backed production entry store.

#### Tests

- Added service tests for valid registration.
- Added service tests for missing employee, inactive employee, missing product, inactive product and invalid entry data.
- Updated handler/server tests to use seeded employee/product stores.
- Verified PostgreSQL migrations and repository tests with local Docker PostgreSQL running.

#### Refactoring

- Moved ID generation and entry construction out of the handler and into the service.
- Kept in-memory stores for employees/products until Milestone 5 rather than doing a large persistence rewrite inside this lesson.

#### Code Review

- An experienced Go engineer would approve the application-service boundary and error mapping.
- An experienced Go engineer would not consider the cross-resource consistency story complete yet because employees/products are not persisted in PostgreSQL. This is documented technical debt and is the natural input to Milestone 5.

#### Exercises

- Explain why a transaction does not help if some validated data lives outside the database.
- Add a failing test that demonstrates production registration after employee deactivation.
- Design how employee/product foreign keys would change the production schema after those tables are persisted.

#### Interview Questions

- What should define a transaction boundary?
- Why are service-level validations still useful if the database also has constraints?
- When should validation be enforced by code, by database constraints or by both?
- Why do consumer-owned interfaces reduce coupling?

#### Roadmap Update

- Lesson 4.4 completed.
- Milestone 4 completed.
- Current milestone moved to Milestone 5.
- Known technical debt updated for employee/product persistence and transactional consistency.

### Milestone 4 Review

#### Architecture Review

An experienced Go engineer would approve the milestone as a learning-oriented vertical slice: production registration has domain validation, an application service, PostgreSQL persistence, sqlc queries and integration tests.

The main architectural weakness is mixed persistence: employees/products are in-memory while production entries are PostgreSQL-backed. This is acceptable temporarily because Milestone 5 is explicitly about persistence quality, but it must not remain long-term.

#### Code Review

The code remains explicit and small. The handler does not expose sqlc types. The service owns business coordination. The repository translates infrastructure errors into domain errors.

The main improvement for the next milestone is to persist employees/products and replace application-only reference checks with stronger database-backed consistency.

#### Refactoring

No broad refactor is needed before Milestone 5. The next refactor should be persistence-focused: introduce PostgreSQL-backed repositories for employees/products and revisit foreign keys/transactions.

#### Interview Review

You should now be able to discuss why sqlc is different from an ORM, why migrations are production history, how context reaches pgx calls, why package-name stutter matters and what a transaction boundary should represent.

#### Completion Criteria

- Production entries persist in PostgreSQL.
- PostgreSQL is integrated through pgx and sqlc.
- Migrations run with goose.
- Production registration validates employee/product business references.
- Tests, build, lint and sqlc generation pass.
- Roadmap updated.
