# ROADMAP.md

> **Project:** MES Lite - Go Apprenticeship
>
> This document is the **single source of truth** for the project. It defines:
>
> - project goals
> - learning roadmap
> - AI teaching rules
> - architecture principles
> - project state
> - quality standards
> - completion criteria
>
> It should evolve together with the project.

## Instructions

Divide current Milestone into Lessons. Write it to this file.
Proceed with the next Lesson of current Milestone.

## STATE

This section must always reflect the current progress.

**Version:** 2.12
**Status:** IN PROGRESS
**Current milestone:** 7 - Production Orders
**Current lesson:** L7.2 - Order Persistence & Reference Integrity
**Completed milestones:**
- Milestone 0
- Milestone 1
- Milestone 2
- Milestone 3
- Milestone 4
- Milestone 5
- Milestone 6
**Next milestone:** 8 - Reporting
**Current branch:** main
**Architecture maturity:** 7.1 / 10
**Go knowledge progress:** 52%
**Interview readiness:** 45%
**Known technical debt:** Production reference foreign keys are `NOT VALID`, so PostgreSQL enforces new production entries but does not validate legacy rows created before the constraint. Query parameters are implemented for list endpoints; explicit OpenAPI query-parameter documentation should be reviewed later. Auth-user management is intentionally limited to durable bootstrap admin creation; full user-management CRUD is postponed until there is a concrete business workflow.

The AI must update this section at the end of every lesson and milestone.

### Value examples
Status: NOT STARTED / IN PROGRESS
Current milestone: 3 - Products
Current lesson: 3.2 Product Validation
Completed milestones:
- Milestone 0
- Milestone 1

### Knowledge Matrix

This checklist tracks learning progress across the entire project.
The AI should update it after every completed milestone.

**Language**
- [ ] Modules
- [ ] Packages
- [x] Visibility
- [x] Variables
- [x] Constants
- [x] Functions
- [x] Structs
- [x] Methods
- [x] Receivers
- [x] Pointers
- [x] Interfaces
- [ ] Embedding
- [x] Custom Types
- [x] iota
- [x] Errors
- [x] Error Wrapping
- [x] defer
- [ ] panic
- [x] Context
- [ ] Generics (intro)
- [ ] Reflection (overview)

**Standard Library**
- [x] net/http (via Fuego)
- [x] context
- [x] errors
- [ ] io
- [x] encoding/json
- [ ] encoding/csv
- [x] strings
- [ ] bytes
- [x] crypto/rand
- [x] time
- [ ] sync
- [ ] sync/atomic
- [x] database/sql (concepts)
- [ ] log/slog
- [ ] runtime
- [ ] embed

**Concurrency**
- [ ] Goroutines
- [ ] Channels
- [ ] select
- [ ] WaitGroup
- [ ] errgroup
- [ ] Mutex
- [ ] RWMutex
- [ ] Atomic
- [ ] Worker Pools
- [ ] Pipelines
- [ ] Race Detector

**Architecture**
- [x] Vertical Slice
- [x] Dependency Injection
- [x] Package Design
- [x] Repositories
- [x] Aggregates
- [ ] CQRS
- [ ] Event-driven Architecture
- [ ] Observability
- [ ] Production Readiness

**Persistence**
- [x] PostgreSQL
- [x] pgx
- [x] sqlc
- [ ] Transactions
- [x] Optimistic Locking
- [ ] SQL Optimization

**Testing**
- [x] Unit Tests
- [x] Integration Tests
- [x] Sanity Tests
- [x] Table Tests
- [x] httptest
- [ ] Testcontainers
- [ ] Benchmarks
- [ ] Fuzz Tests
- [ ] Race Detection

**Runtime**
- [ ] Scheduler
- [ ] Garbage Collector
- [ ] Escape Analysis
- [ ] Memory Layout

**HTTP**
- [x] Middleware
- [ ] RoundTripper
- [ ] Client
- [x] Server

**SQL**
- [ ] Isolation Levels
- [ ] Indexes
- [ ] Explain Analyze

**Performance**
- [ ] pprof
- [ ] trace
- [ ] benchmarks
- [ ] allocations

---

# 1. Purpose

The primary goal of this project is **not** to finish the project quickly.

The primary goal is to learn by building a real-world production application and become a **Senior Go Engineer** capable of designing, implementing, testing, deploying and maintaining production systems with confidence.

The application is only the vehicle used to learn:

- software design
- idiomatic Go
- concurrency
- observability
- reviewing Go code
- testing
- production engineering
- production architecture
- deploying production services
- contributing to mature Go codebases

Every design decision should optimize for learning.

---

# 2. Learning Goals

After completing this roadmap I should be able to:

## Go

- understand Go philosophy
- understand why Go is different from Java/PHP/C#
- write idiomatic Go
- use the standard library effectively
- know when NOT to introduce abstractions
- design cohesive packages
- write maintainable code
- understand memory semantics
- understand interfaces
- understand value vs pointer semantics

## Architecture

Design production services using:

- Vertical Slice Architecture
- domain-oriented packages
- constructor dependency injection
- explicit dependencies
- explicit error handling
- CQRS where appropriate
- event-driven components where appropriate
- clean package boundaries

## Web Development

Design and build production REST APIs using:

- Fuego
- OpenAPI
- PostgreSQL
- sqlc
- pgx
- Docker

## Concurrency

Write concurrent software safely using:

- goroutines
- channels
- errgroup
- sync
- mutexes
- atomics
- contexts
- cancellation
- worker pools
- pipelines

## Production Engineering

Understand:

- logging
- metrics
- tracing
- graceful shutdown
- configuration
- CI/CD
- profiling
- benchmarking
- performance optimization

## Testing

Become comfortable writing:

- unit tests
- integration tests
- sanity tests
- table tests
- benchmarks
- fuzz tests
- race-safe code

## Interview Readiness

Be able to confidently discuss:

- Go language internals
- scheduler
- garbage collector
- escape analysis
- interfaces
- package design
- concurrency
- architecture
- performance
- testing
- trade-offs

---

# 3. AI Contract

The AI acts as a Senior Go Engineer, mentor, reviewer and pair programmer.

Its primary objective is **teaching**, not generating code.

Whenever there is a conflict between:

- finishing features quickly
- maximizing learning

the AI must always choose learning.

---

# 4. AI Teaching Rules

These rules govern every interaction during the project. They take precedence over implementation speed or feature completeness.

## 1. Learning First

Optimize for understanding, not speed.

Never generate an entire feature immediately. Guide the implementation step by step.

## 2. Incremental Learning

Divide Milestones into Lessons as you see fit.

Each lesson introduces:

- at most one architectural concept
- at most 2–5 new Go concepts

Build on previously learned knowledge.

## 3. Explain Why

Always explain:

**Why → How → Code**

Focus on design decisions and Go philosophy before implementation.

## 4. Idiomatic Go

Always explain why a solution is idiomatic Go.

When helpful, compare it with common PHP/Symfony approaches to highlight differences in philosophy rather than syntax.

### Idioms Checklist
- Accept interfaces, return structs
- Errors are values
- Context first
- Small interfaces
- Composition over inheritance
- Zero value is useful
- Make the zero allocation path obvious
- Keep APIs explicit
- Keep packages cohesive

## 5. Standard Library First

Prefer the Go standard library.

Introduce external dependencies only when they provide clear value, and explain why they are justified.

## 6. Minimal Abstractions

Avoid unnecessary abstractions.

Do not introduce interfaces, wrappers, generic helpers or additional layers until they solve a real problem.

Generics are intentionally postponed until the application naturally requires them.

Remember:

> Concrete first. Abstract later.

## 7. Consumer-Owned Interfaces

Interfaces belong to consumers.

Prefer concrete types until multiple implementations or testing needs justify an interface.

## 8. Simplicity

Prefer:

- explicit code
- composition
- small packages
- readable APIs

Avoid cleverness and hidden magic.

## 9. Continuous Refactoring

Continuously review previously written code.

Whenever newly acquired knowledge enables a simpler or more idiomatic solution, propose a refactoring before continuing.

## 10. Continuous Code Review

Actively detect and explain:

- non-idiomatic Go
- overengineering
- hidden complexity
- unnecessary abstractions
- package boundary violations
- performance issues
- testing gaps

Every review should explain both **what** should change and **why**.

## 11. Architecture by Business

Technology serves the business.

Every feature and architectural decision should begin by explaining the business problem it solves.

Discuss alternatives, trade-offs and justify the chosen solution.

Whenever an architectural decision has long-term consequences, create or update an ADR before implementing it.

## 12. Production Quality

Never sacrifice quality for speed.

Every implementation should be production-ready.

No tutorial shortcuts.

## 13. Lesson Completion

Every lesson ends with:

- summary
- Go concepts learned
- idioms learned
- common mistakes
- exercises
- interview questions

A lesson is complete only if:

- business requirement implemented
- idiomatic code
- tests passing
- linter passing
- OpenAPI updated
- AI code review completed
- refactoring discussed
- Knowledge Matrix updated
- Project State updated

## 14. Milestone Completion

Every milestone ends with:

- architecture review
    - Would an experienced Go engineer approve this PR? If not: Why? How can it be improved?
- code review
    - Would an experienced Go engineer approve this PR? If not: Why? How can it be improved?
- refactoring
- interview review
- roadmap update

A milestone is complete only after every lesson satisfies the Lesson Completion criteria.

## 15. Interview Readiness

Continuously prepare for Senior Go interviews.

After each completed topic:

- ask interview-style questions
- discuss trade-offs
- review common mistakes
- revisit previous concepts when appropriate

## 16. Lesson Template

Every lesson follows the same structure:

1. Business Context
2. Problem
3. Design Discussion
4. Go Concepts
5. Architecture Concepts
6. Implementation
7. Tests
8. Refactoring
9. Code Review
10. Exercises
11. Interview Questions
12. Roadmap Update

---

# 5. General Development Principles

## Production Ready

Every feature should be written as if it will be deployed to production.

No tutorial shortcuts.

---

## Learn by Doing

Every Go concept should immediately be used in real code.

No isolated exercises unless required.

---

## Refactor Often

Architecture is expected to evolve.

Earlier solutions are allowed to become obsolete.

Refactoring is a required part of learning.

---

## Simplicity First

Prefer: Simple, Correct, Readable
over: Flexible, Generic, Complex

---

## Standard Library First

Whenever possible use:

- net/http
- context
- errors
- io
- sync
- encoding/json
- log/slog

before introducing third-party packages.

---

## Vertical Slice Architecture

The application should be organized around business capabilities rather than technical layers.

Avoid structures such as:

```
handlers/
services/
repositories/
models/
controllers/
```

Instead prefer:

```
employees/
products/
production/
orders/
reporting/
```

Each feature owns its:

- HTTP handlers
- business logic
- persistence
- tests
- DTOs
- validation

---

## Testing From Day One

Every milestone must include tests.

Testing is never postponed.

---

## OpenAPI First

Every endpoint should be documented.

The OpenAPI specification is part of the application.

---

## Clean Git History

Commits should be:

- small
- meaningful
- independently reviewable

---

# 6. Project Description

## Background

A small manufacturing company produces ventilation components.

Production is currently tracked using Excel spreadsheets.

Employees manually record:

- what they produced
- when they produced it
- how much they produced

Reporting is manual.

Finding historical information is slow.

Management wants a lightweight system instead of a large ERP/MES solution.

---

## Goal

Replace Excel with a modern web application.

The application should remain simple enough for a small company while being designed with production-grade engineering practices.

---

## Primary Users

### Production Worker

Registers completed work.

---

### Team Leader

Monitors current production.

Reviews completed work.

Corrects mistakes.

---

### Production Manager

Creates production orders.

Tracks progress.

Generates reports.

Measures productivity.

---

### Administrator

Maintains users.

Maintains products.

Configures the system.

---

## MVP Scope

The first version supports:

- employees
- products
- production entries
- production orders
- authentication
- reporting
- CSV import

Everything else is intentionally postponed.

---

## Future Scope

Later milestones may introduce:

- machines
- workstations
- warehouse
- traceability
- quality control
- notifications
- scheduling
- dashboards
- background jobs
- integrations
- analytics

The roadmap intentionally delays these features until they become valuable learning opportunities.

---

# 7. Non-Functional Requirements

The application should be:

- maintainable
- testable
- observable
- secure
- modular
- production ready

Performance should be considered from the beginning, but premature optimization should be avoided.

Reliability and readability are more important than micro-optimizations.

---

# 8. Technical Stack

The goal of the technology stack is **not to use every popular library**, but to
learn how production Go services are commonly built while remaining as close as
possible to the Go ecosystem.

The standard library should always be preferred unless another solution provides
clear benefits.

## Language:
- `Go` (latest stable version)

## HTTP Framework:
- `Fuego`, reasons:
    - OpenAPI-first development
    - Minimal abstraction over `net/http`
    - Type-safe handlers
    - Automatic OpenAPI generation
    - Good fit for idiomatic Go
    - Small learning surface

The AI should explain what Fuego provides and what still comes directly from the
standard library.

## Database and SQL Access
-`PostgreSQL`, reasons:
    - Production-grade relational database
    - Rich SQL features
    - Excellent Go support
    - Widely used in industry
-`pgx` (PostgreSQL driver), reasons:
    - Excellent performance
    - Idiomatic API
    - Industry standard
-`sqlc` (SQL is written manually, `sqlc` generates type-safe Go code), reasons:
    - Learn SQL instead of hiding it behind an ORM
    - Compile-time query validation
    - Better performance
    - Simpler architecture
    - Explicit queries

No ORM will be used.

Understanding SQL is considered a required Senior Go skill.

## Migrations
- `goose`, reasons:
    - Simple
    - Reliable
    - SQL-first

## Testing
- Standard library
- `Testcontainers`
- `httptest`

The testing strategy should prefer integration tests whenever they provide more
confidence than mocks.
Sanity tests only where appropriate.

## Documentation
- `OpenAPI` (automatically generated through Fuego)

Documentation is considered part of the application.

## Logging
- `log/slog`, reasons:
    - Standard library
    - Structured logging
    - Production ready

## Configuration
- Environment variables

No configuration framework unless justified.

## Observability
- OpenTelemetry
- Prometheus
- Health endpoints
- Readiness endpoints

## Development
- Docker
- Docker Compose
- Makefile
- golangci-lint
- air
- GitHub Actions

---

# 9. Architecture Principles

The project intentionally avoids copying Java or PHP architecture.

The goal is to learn how experienced Go engineers organize applications.

## Primary Architecture: Vertical Slice Architecture

Business capabilities are the primary unit of organization.

Every business capability owns its:
- handlers
- application logic
- persistence
- tests
- DTOs
- validation

Instead of:
```
handlers/
services/
repositories/
models/
```
prefer:
```
employees/
products/
production/
orders/
reporting/
```

## Dependency Direction

Dependencies always point inward.

```
HTTP
↓
Application
↓
Domain
↓
Repository Interface
↓
Infrastructure
```

Business logic must never depend on infrastructure.

## Constructor Dependency Injection
- Dependencies are passed explicitly.
- No service locator.
- No global state.
- No dependency injection container.

## Package Design

Packages should be:
- cohesive
- small
- easy to understand

Package boundaries should reflect business concepts.

## Interfaces

Interfaces belong to consumers.

Do not create interfaces "just in case".

Prefer concrete implementations until abstraction becomes necessary.

## Error Handling

Errors are values.

Never ignore returned errors.

Prefer:
- wrapping
- context
- explicit propagation

Avoid panic except for unrecoverable startup failures.

## Context

Every request receives a Context.

Contexts should:
- propagate cancellation
- propagate deadlines
- carry request-scoped values

Contexts must never contain business data.

## Explicitness

Prefer explicit code over magic.

Avoid:
- hidden dependencies
- reflection
- unnecessary frameworks
- code generation (except sqlc/OpenAPI)

## Simplicity
Prefer:
- Simple > Clever
- Readable > Generic
- Explicit > Automatic
- Maintainable > Flexible

## Standard Library

Whenever possible use standard library
Introduce additional libraries only if absolutely necessary.

---

# 10. Architecture Decision Records (ADR)

Important architectural decisions must be documented as Architecture Decision Records (ADRs).

The goal is to capture the reasoning behind decisions rather than just the final outcome.

An ADR should be created whenever a significant architectural, technological or organizational decision is made.

Examples:
- Choosing Vertical Slice Architecture
- Choosing sqlc over an ORM
- Choosing pgx
- Introducing CQRS
- Introducing background jobs
- Introducing event-driven communication
- Introducing caching
- Changing package structure

Each ADR should be stored in: `docs/adr/` directory, using the following naming convention:
```
0001-use-vertical-slice.md
0002-use-sqlc.md
0003-introduce-cqrs.md
```

Each ADR should contain:
- Status
- Context
- Decision
- Alternatives Considered
- Trade-offs
- Consequences

The AI should:
- propose creating an ADR whenever appropriate,
- discuss alternatives before making a decision,
- reference previous ADRs when they influence new decisions,
- update ADRs if a decision is intentionally superseded.

---

# 11. Recommended References

## Official

- A Tour of Go
- Effective Go
- Go Blog
- Go Language Specification
- Go Memory Model
- pkg.go.dev

## Books

Essential:

- Learning Go — Jon Bodner
- The Go Programming Language — Donovan & Kernighan
- Let's Go Further — Alex Edwards
- 100 Go Mistakes and How to Avoid Them — Teiva Harsanyi
- Concurrency in Go — Katherine Cox-Buday

Optional:

- Cloud Native Go
- Distributed Services with Go

## Blogs

- Alex Edwards
- Dave Cheney
- Ardan Labs
- ThreeDotsLabs
- Ben Johnson

## Talks

- Rob Pike
- Russ Cox
- Francesc Campoy
- GopherCon

---

# 12. Milestones

## Milestone 0 - Development Environment

Status

✅ Completed

### Lessons
- L1 — Project Foundation (Git, structure, Go module, tooling)
- L2 — Docker Environment & First ADR

### Goal

Create a professional development environment identical to what would be used in
a production Go project.

### Business Value

None.

This milestone exists entirely to establish good engineering practices.

### Deliverables

- Git repository
- Project structure
- Go modules
- Docker Compose
- PostgreSQL
- Makefile
- golangci-lint
- air
- GitHub Actions
- pre-commit hooks
- README

### Go Concepts

- Go toolchain
- modules
- packages
- gofmt
- go mod
- go test
- go build
- go run

### Standard Library

None.

### Architecture Concepts

- repository layout
- project organization
- package visibility
- project conventions

### Testing

- verify project builds
- CI executes tests

### Exercises

- create project from scratch
- explain module system
- explain package visibility

### Interview Topics

- Why Go modules?
- Why package main?
- Why internal?
- Why cmd?

### Definition of Done

- project builds
- CI passes
- lint passes
- Docker starts successfully

---

## Milestone 1 - Bootstrap HTTP Service

Status

✅ Completed

### Goal

Create the first running HTTP service.

### Business Value

The application becomes deployable.

### Features

- application startup
- health endpoint
- version endpoint
- OpenAPI endpoint

### Go Concepts

- package main
- imports
- variables
- constants
- functions
- structs
- methods
- pointers
- receivers
- errors
- defer

Maximum new concepts:
10

### Standard Library

- net/http
- context
- encoding/json
- errors

### Architecture Concepts

- composition root
- dependency injection
- graceful shutdown
- configuration

### Testing

- httptest
- first integration tests

### Exercises

- implement health endpoint
- add version endpoint
- explain graceful shutdown

### Interview Topics

- Why net/http?
- What is a receiver?
- Pointer vs value receiver?
- Why defer?

### Definition of Done

- server starts
- server shuts down gracefully
- OpenAPI generated
- tests pass
- roadmap updated

---

## Milestone 2 - Employees

Status

✅ Completed

### Lessons

- **L2.1** — Visibility & Zero Values: Employee Entity
- **L2.2** — Constructors & Slices: Creating & Listing Employees
- **L2.3** — Error Wrapping, Validation & Maps: Updating, Deactivating & Testing

### Goal

Implement employee management.

### Business Value

The company can register employees who perform production work.

### Features

- create employee
- update employee
- deactivate employee
- list employees

### Go Concepts

- visibility
- slices
- maps
- constructors
- zero values
- interfaces
- error wrapping

### Standard Library

- slices
- cmp (if applicable)
- errors

### Architecture Concepts

- Vertical Slice
- package cohesion
- repository interface
- validation
- DTO separation

### Testing

- table tests
- integration tests
- validation tests

### Exercises

- design Employee entity
- decide pointer vs value
- explain package boundary

### Interview Topics

- Why interfaces belong to consumers?
- Why slices instead of arrays?
- Zero values?
- Exported vs unexported?

### Definition of Done

- CRUD complete
- validation complete
- tests passing
- code review completed

---

## Milestone 3 - Products

Status

✅ Completed

### Lessons

- **L3.1** — Custom Types & iota: Product Entity ✅
- **L3.2** — Stringer & strings: Product Handlers & Search ✅
- **L3.3** — Value Objects & Testing: Product Validation & Review ✅

### Goal

Manage manufactured products.

### Business Value

Employees can register production for real products.

### Features

- create product
- update product
- deactivate product
- search products

### Go Concepts

- custom types
- iota
- enums
- Stringer
- value semantics

### Standard Library

- strings
- strconv
- fmt

### Architecture Concepts

- value objects
- aggregate boundaries
- package APIs

### Testing

- table-driven tests
- API integration tests

### Exercises

- design Product model
- choose custom types
- implement String()

### Interview Topics

- Why doesn't Go have enums?
- Why custom types?
- Value semantics?

### Definition of Done

- products functional
- tests passing
- OpenAPI updated
- roadmap updated

---

## Milestone 4 - Production Registration

Status

✅ Completed

### Lessons

- **L4.1** — Production Entry Domain & HTTP Registration ✅
- **L4.2** — PostgreSQL, Migrations & sqlc Setup ✅
- **L4.3** — Repository Implementation & Context Propagation ✅
- **L4.4** — Transaction Boundary & Business Validation Review ✅

### Lesson 4.1 Completion Notes

#### Business Context

Production workers need to record completed work without Excel.

#### Problem

The system had employees and products, but no core production-registration capability.

#### Design Discussion

The first implementation keeps the vertical slice small: `internal/production` owns the production entry entity, validation, in-memory persistence contract and HTTP handler. PostgreSQL, sqlc and transactions are intentionally postponed to later lessons because introducing them before the business model would hide the domain rules inside infrastructure work.

#### Go Concepts

- `time.Time` for business timestamps
- `crypto/rand` and `encoding/hex` for standard-library UUID-shaped IDs
- error wrapping with sentinel errors
- context propagation from HTTP handler to store

#### Architecture Concepts

- vertical slice for production registration
- package naming that avoids stutter: `production.Entry`, not `production.ProductionEntry`
- concrete in-memory implementation before persistence abstraction grows

#### Implementation

- Added `POST /production-entries`.
- Added `production.Entry` with employee, product, quantity, workstation, timestamp and comment.
- Added production entry validation and UUID-shaped ID generation.
- Wired production handler through the server composition root.

#### Tests

- Domain validation tests added.
- In-memory store tests added.
- HTTP registration tests added.
- Server route test added.

#### Refactoring

- Renamed `ProductionEntry` to `Entry` after lint identified package-name stutter.

#### Code Review

- An experienced Go engineer would likely approve the scope for this lesson because it is small, explicit and tested.
- Remaining gap: employee/product existence is not validated yet; this belongs in the next application-service/persistence lessons where the transaction boundary can be designed properly.

#### Exercises

- Explain why `production.Entry` is a better exported name than `production.ProductionEntry`.
- Add a table test for future maximum comment length validation.
- Explain why timestamps are normalized to UTC.

#### Interview Questions

- Why is `time.Time` usually preferred over strings for timestamps in Go APIs?
- What does error wrapping give us when translating domain errors to HTTP errors?
- Why should `context.Context` be passed from the handler to persistence code?
- Why is package-name stutter considered non-idiomatic Go?

#### Roadmap Update

- Lesson 4.1 completed.
- Current lesson moved to Lesson 4.2.
- Standard Library `time` marked complete in the Knowledge Matrix.

### Lesson 4.2 Completion Notes

#### Business Context

Production entries are core business records. Keeping them only in memory would lose production history and fail the business goal of replacing Excel.

#### Problem

The project had PostgreSQL available in Docker, but no schema migrations, no sqlc configuration and no executable migration command.

#### Design Discussion

The database contract now starts with production entries only. This avoids a large persistence rewrite and keeps the lesson focused on PostgreSQL, migrations and sqlc. The SQL schema owns non-negotiable data integrity rules such as positive quantity and non-blank workstation, while Go still validates early for better API errors.

#### Go Concepts

- `database/sql` as a stable standard-library abstraction for migration tooling
- blank imports for database driver registration
- defers and why `os.Exit` must not bypass cleanup
- generated code boundaries

#### Architecture Concepts

- SQL-first persistence design
- migration files as versioned database history
- sqlc-generated query package kept inside the production vertical slice
- dependency direction preserved: HTTP still does not depend directly on generated SQL code

#### Implementation

- Added `migrations/0001_create_production_entries.sql`.
- Added `sqlc.yaml`.
- Added `internal/production/queries/entries.sql`.
- Generated `internal/production/productiondb` with typed sqlc queries.
- Implemented `cmd/migrate` using `pgx` through `database/sql` and `goose`.
- Added `DATABASE_URL` and `MIGRATIONS_DIR` configuration.
- Added `make migrate` and `make sqlc` targets.

#### Tests

- Configuration tests cover database and migration environment variables.
- Generated sqlc code is compiled by `go test ./...` and `go build ./...`.
- Live migration execution passed against local PostgreSQL with `make migrate`.
- Verified the `production_entries` schema and goose version table through `psql`.

#### Refactoring

- Refactored `cmd/migrate` to return an exit code from `run()` instead of calling `os.Exit` before defers execute.

#### Code Review

- An experienced Go engineer would approve the direction: SQL is explicit, generated code is isolated and migrations are versioned.
- Remaining gap: the HTTP handler still uses the in-memory store. This is intentional; Lesson 4.3 will introduce a PostgreSQL-backed repository and wire context propagation through it.

#### Exercises

- Explain why `CHECK (quantity > 0)` belongs in the database even though Go also validates quantity.
- Inspect the `production_entries` table with `psql` and identify every database constraint.
- Modify `entries.sql`, rerun `make sqlc`, and inspect the generated type changes.

#### Interview Questions

- Why choose sqlc instead of an ORM?
- What is the role of migrations in production systems?
- Why does `database/sql` need a blank driver import?
- Why is it dangerous to call `os.Exit` before deferred cleanup runs?

#### Roadmap Update

- Lesson 4.2 completed.
- Current lesson moved to Lesson 4.3.
- `database/sql`, PostgreSQL, pgx and sqlc marked complete in the Knowledge Matrix.

### Lesson 4.3 Completion Notes

#### Business Context

Production entries now need to survive application restarts. The production registration endpoint must write to PostgreSQL instead of only using in-memory state.

#### Problem

The project had a schema and sqlc queries, but the running application still used `production.InMemoryStore`.

#### Design Discussion

The handler continues to depend on the small `production.Store` interface. The new PostgreSQL store adapts sqlc-generated database types to the domain `production.Entry` type. This keeps generated code out of HTTP handlers and preserves the package boundary.

#### Go Concepts

- `context.Context` propagation from HTTP handler to pgx/sqlc calls
- adapter code between generated persistence models and domain models
- error translation with `errors.Is` and `errors.As`
- integer-size boundary between Go `int` and PostgreSQL `integer`

#### Architecture Concepts

- repository implementation inside the production vertical slice
- generated SQL package isolated below domain/application code
- composition root chooses concrete dependencies
- infrastructure errors translated to domain errors

#### Implementation

- Added `production.PostgresStore` backed by sqlc queries.
- Added UUID conversion helpers between string IDs and `pgtype.UUID`.
- Mapped PostgreSQL duplicate-key and constraint errors to domain errors.
- Wired `cmd/server` to use `pgxpool` and `production.NewPostgresStore`.
- Refactored `cmd/server` to use `run() int` so database pool cleanup defers run before process exit.
- Added a quantity overflow guard because PostgreSQL `integer` is 32-bit.

#### Tests

- Added PostgreSQL repository integration tests for save/find, duplicate detection, not found and invalid UUID handling.
- Tests run migrations before exercising the repository.
- Tests skip only when PostgreSQL is unavailable.
- Verified with local Docker PostgreSQL running.

#### Refactoring

- Preserved the existing in-memory store for fast handler tests.
- Kept the `Store` interface consumer-owned by the handler package slice instead of exposing sqlc directly.

#### Code Review

- An experienced Go engineer would approve this direction because the generated persistence code is isolated, domain errors are preserved and context reaches the database call.
- Remaining gap: registering production does not yet validate employee/product existence in one transaction. Lesson 4.4 will define that transactional boundary and business validation strategy.

#### Exercises

- Explain why the handler should not return `productiondb.ProductionEntry` directly.
- Add a test that proves PostgreSQL rejects invalid quantity even if application validation is bypassed.
- Trace `c.Context()` from the Fuego handler to `pgx.QueryRow`.

#### Interview Questions

- Why do repositories often translate infrastructure errors into domain errors?
- Why should generated sqlc code not become the public API of the business package?
- What happens when a request context is cancelled while pgx is waiting on a query?
- Why must Go code care that PostgreSQL `integer` is 32-bit?

#### Roadmap Update

- Lesson 4.3 completed.
- Current lesson moved to Lesson 4.4.
- Repositories marked complete in the Knowledge Matrix.

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

### Goal

Register production performed by employees.

This is the first milestone implementing the core business capability.

### Business Value

The application begins replacing Excel.

### Features

- register production
- employee
- product
- quantity
- workstation
- timestamp
- comment

### Go Concepts

- time.Time
- UUID
- context propagation
- sqlc
- transactions
- error handling across layers

### Standard Library

- time
- context
- database/sql (concepts)
- errors

### Architecture Concepts

- application service
- repository implementation
- transactional boundary
- business validation

### Testing

- integration tests with PostgreSQL
- transaction tests
- validation tests

### Exercises

- model Production Entry
- define transaction boundary
- explain why sqlc instead of ORM

### Interview Topics

- Why sqlc?
- What belongs inside a transaction?
- Context propagation?
- Why explicit SQL?

### Definition of Done

- production entries persist
- PostgreSQL integrated
- tests passing
- OpenAPI updated
- roadmap updated

---

## Milestone 5 - Persistence & Data Access

Status

✅ Completed

### Lessons

- **L5.1** — Persist Employees & Products ✅
- **L5.2** — Validation Flow & Domain Invariants ✅
- **L5.3** — Pagination, Filtering & Sorting ✅
- **L5.4** — Optimistic Locking & Concurrent Updates ✅
- **L5.5** — Transactional Reference Integrity & Milestone Review ✅

### Lesson 5.5 Completion Notes

#### Business Context

Production entries are business records. They must not reference employees or products that do not exist in durable master data.

#### Problem

Production registration validated employee/product references in application code, but the database still stored plain text references without foreign keys. A bug, import, direct SQL write or future service could insert production entries pointing at missing master data.

#### Design Discussion

PostgreSQL now owns reference existence for new production entries through foreign keys from `production_entries.employee_id` to `employees.id` and from `production_entries.product_sku` to `products.sku`.

The constraints are created as `NOT VALID`. This is a production-safe migration pattern when existing data may be dirty: PostgreSQL enforces the constraint for new writes, but does not scan and reject legacy rows during deployment. A future cleanup migration can validate existing rows and then run `VALIDATE CONSTRAINT`.

Application validation remains valuable because it returns better business errors and checks active/inactive state. Database constraints are the final integrity boundary for reference existence.

#### Go Concepts

- database constraint errors translated into domain errors
- integration tests that prove failed writes do not leave rows behind
- distinguishing business validation from persistence integrity

#### Architecture Concepts

- referential integrity belongs in the database
- application services still own workflow rules
- database constraints protect against alternate write paths
- production-safe migration with `NOT VALID` foreign keys

#### Implementation

- Added migration `0004_add_production_reference_foreign_keys.sql`.
- Added foreign key from `production_entries.employee_id` to `employees.id`.
- Added foreign key from `production_entries.product_sku` to `products.sku`.
- Mapped PostgreSQL foreign-key violations to `production.ErrInvalidEntry` at the production persistence boundary.
- Updated production PostgreSQL integration tests to seed employee/product reference data.

#### Tests

- Added production repository test for missing employee reference.
- Verified failed foreign-key insert leaves no production row behind.
- Verified with `make sqlc`.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

No broad transaction manager was introduced. The current write is a single PostgreSQL statement, so PostgreSQL already executes it atomically. A future multi-step workflow can introduce an explicit transaction abstraction when there is a real multi-write business operation.

#### Code Review

An experienced Go engineer would approve the foreign-key enforcement and the `NOT VALID` migration choice for safe rollout. The main caveat is that legacy rows are not validated yet; this is documented technical debt rather than a hidden inconsistency.

#### Exercises

- Explain why `NOT VALID` foreign keys still protect new writes.
- Add a cleanup query that finds existing production rows with missing employee references.
- Explain why active/inactive validation stays in application code instead of a foreign key.

#### Interview Questions

- What is referential integrity?
- What does a PostgreSQL `NOT VALID` constraint do?
- Why keep application validation if the database has foreign keys?
- When should you introduce an explicit transaction boundary?

#### Roadmap Update

- Lesson 5.5 completed.
- Milestone 5 completed.
- Current milestone moved to Milestone 6.
- Known technical debt updated for `NOT VALID` reference constraints and OpenAPI query-parameter documentation.

### Milestone 5 Review

#### Architecture Review

An experienced Go engineer would approve the milestone direction. Employees, products and production entries are now PostgreSQL-backed, generated sqlc types remain below slice boundaries, HTTP handlers depend on domain-facing stores and production registration has both application-level reference validation and database-level referential integrity for new writes.

The main architectural limitation is that cross-slice transaction coordination is still implicit. This is acceptable because the current production registration persistence step is a single insert. A future multi-write operation should introduce an explicit transaction boundary rather than a generic transaction abstraction in advance.

#### Code Review

The code is explicit and idiomatic for the current maturity level. Constructors now reject invalid domain entities, repositories validate defensively, update endpoints use optimistic locking and SQL sorting avoids dynamic SQL interpolation.

The main improvement is API documentation: list query parameters are implemented and tested, but OpenAPI query parameter metadata should be reviewed later.

#### Refactoring

Store files were split into `store.go`, `store_in_memory.go` and `store_postgres.go`, improving navigation without changing exported names. Product search now reuses product listing options instead of having a separate persistence method.

#### Interview Review

You should now be able to explain sqlc vs ORM, repository adapters, constructor validation, database constraints, `NOT VALID` foreign keys, limit/offset pagination, safe dynamic sorting and optimistic locking with version columns.

#### Completion Criteria

- Employees and products persist in PostgreSQL.
- Validation flow is documented and constructors enforce invariants.
- List endpoints support pagination, filtering and sorting.
- Optimistic locking prevents stale employee/product updates.
- Production reference existence is enforced by PostgreSQL foreign keys for new writes.
- Tests, build, lint and sqlc generation pass.
- Roadmap updated.

### Lesson 5.4 Completion Notes

#### Business Context

Administrators and team leaders may edit the same employee or product around the same time. Without conflict detection, the last write silently wins and can overwrite another user's changes.

#### Problem

Employee and product updates replaced rows without checking whether the caller edited a stale copy. This creates lost updates under concurrent usage.

#### Design Discussion

The lesson introduces optimistic locking with an integer `version` column. Clients receive the current version when reading or creating an employee/product. Update requests must submit the version they edited. The database update succeeds only when the submitted version matches the stored version, then increments the version atomically.

This keeps the solution simple and explicit. We do not hold long database locks across HTTP requests. Instead, conflicts are detected at write time and returned as `409 Conflict`.

#### Go Concepts

- optimistic concurrency with version fields
- stale-write detection with sentinel errors
- update methods returning updated structs
- concurrent integration tests with goroutines and channels

#### Architecture Concepts

- conflict detection belongs at the persistence boundary
- HTTP translates stale writes to `409 Conflict`
- database update predicates enforce atomic compare-and-swap behavior
- handlers return the incremented version after successful writes

#### Implementation

- Added `Version` to employee and product domain structs.
- Added `ErrVersionConflict` for employees and products.
- Added migration `0003_add_employee_product_versions.sql`.
- Updated sqlc create/get/list/update queries to include `version`.
- Changed employee/product store `Update` methods to return the updated entity.
- Updated in-memory stores to reject stale versions and increment versions.
- Updated PostgreSQL stores to use `WHERE id/sku = $1 AND version = $6` update predicates.
- Updated employee/product update HTTP requests to require `version`.
- Mapped stale updates to `409 Conflict`.

#### Tests

- Updated handler tests to submit versions and assert incremented response versions.
- Added handler tests for stale-version conflicts.
- Added PostgreSQL stale-version tests for employees and products.
- Added a PostgreSQL concurrent update test using two goroutines updating the same employee version; one update succeeds and one returns `ErrVersionConflict`.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

Store update contracts now return the updated entity so handlers do not guess the next version. This keeps version increments owned by persistence implementations.

#### Code Review

An experienced Go engineer would approve this optimistic-locking design for simple master-data updates. The main trade-off is API friction: clients must now send a version on updates. That is intentional because silent overwrites are worse than explicit conflicts.

#### Exercises

- Explain why the version check must happen in the SQL `WHERE` clause.
- Add a product concurrent update test mirroring the employee test.
- Design the client behavior after receiving `409 Conflict`.

#### Interview Questions

- What problem does optimistic locking solve?
- How is optimistic locking different from pessimistic locking?
- Why does `UPDATE ... WHERE version = ?` avoid lost updates?
- When would `SELECT ... FOR UPDATE` be a better choice?

#### Roadmap Update

- Lesson 5.4 completed.
- Current lesson moved to Lesson 5.5.
- Known technical debt updated: optimistic locking completed; foreign-key-backed transactional reference integrity remains pending.

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

### Lesson 5.1 Scope

Persist employees and products in PostgreSQL so production registration no longer depends on in-memory reference data.

#### Business Context

Employees and products are master data. If they disappear after a restart, production records cannot be trusted as durable business history.

#### Problem

Production entries already persist in PostgreSQL, but employee and product stores are still in-memory. This creates mixed persistence and prevents database-backed reference validation.

#### Design Discussion

This lesson keeps the existing employee/product package APIs and replaces only the infrastructure implementation. The HTTP handlers still depend on the existing `Store` interfaces, while the composition root chooses PostgreSQL stores for the running server.

Foreign keys from production entries to employees/products are intentionally postponed to Lesson 5.5 because they require a migration strategy for existing production data and a discussion about transaction boundaries.

#### Go Concepts

- persistence adapters between domain structs and generated SQL structs
- wrapped sentinel errors for storage failures
- context propagation through repository methods

#### Architecture Concepts

- persistence implementation inside each vertical slice
- composition root selects concrete infrastructure
- domain-facing APIs remain stable while storage changes

### Lesson 5.1 Completion Notes

#### Business Context

Employees and products are master data required for trustworthy production history.

#### Problem

Production entries were PostgreSQL-backed, but employees and products still lived in memory. Restarting the server lost reference data and kept production registration only partially durable.

#### Design Discussion

The existing handler-facing store interfaces stayed unchanged. Each vertical slice now owns its SQL queries, generated sqlc package and PostgreSQL adapter. This keeps generated database types below the package boundary and avoids leaking infrastructure details into HTTP handlers or production registration logic.

Foreign keys were intentionally postponed to Lesson 5.5. Adding them safely needs a migration and transaction-boundary discussion, especially for existing production rows that may reference master data created before this lesson.

#### Go Concepts

- sqlc-generated code as an implementation detail
- adapter functions from database rows to domain structs
- sentinel errors wrapped with storage context
- package-local integration test setup
- PostgreSQL advisory locks to serialize concurrent test migrations

#### Architecture Concepts

- persistence adapters inside vertical slices
- composition root choosing concrete infrastructure
- stable domain-facing APIs while infrastructure changes
- explicit remaining consistency gap before foreign keys

#### Implementation

- Added `employees` and `products` PostgreSQL tables in migration `0002`.
- Added sqlc query files for employee and product create/read/list/update/search operations.
- Generated `employeesdb` and `productsdb` packages.
- Added `employees.PostgresStore` and `products.PostgresStore`.
- Wired `cmd/server` to use PostgreSQL employee and product stores.
- Preserved in-memory stores for fast handler tests.

#### Tests

- Added employee PostgreSQL store integration tests.
- Added product PostgreSQL store integration tests.
- Updated repository test setup to serialize migrations with a PostgreSQL advisory lock.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

The server composition root now uses durable stores for all currently persisted business data. No handler refactor was needed because the existing interfaces already expressed the right package boundary.

#### Code Review

An experienced Go engineer would approve the direction: SQL remains explicit, generated code is isolated, and domain errors are preserved. The main remaining issue is referential integrity: production entries still reference employee/product IDs as text without database foreign keys.

#### Exercises

- Explain why `employeesdb.Employee` should not be returned directly from HTTP handlers.
- Add a database constraint test that proves blank product names are rejected by PostgreSQL.
- Trace `POST /employees` from the handler to `pgx.QueryRow`.

#### Interview Questions

- Why keep sqlc-generated types behind repository adapters?
- What problem do advisory locks solve in integration tests?
- Why can a database constraint still be useful when Go validation already exists?
- What are the trade-offs of adding foreign keys in a later migration?

#### Roadmap Update

- Lesson 5.1 completed.
- Current lesson moved to Lesson 5.2.
- Known technical debt updated: employees/products are now PostgreSQL-backed; foreign keys and transaction consistency remain pending.

### Goal

Build a robust persistence layer and establish production-quality database
practices.

This milestone intentionally focuses on persistence quality rather than adding
new business features.

### Business Value

The application becomes reliable and maintainable as data volume grows.

### Features

- repository implementations
- pagination
- filtering
- sorting
- optimistic locking
- soft deletes (where appropriate)

### Go Concepts

- interfaces as contracts
- composition
- package boundaries
- custom errors
- sentinel vs wrapped errors

### Standard Library

- errors
- context
- database/sql concepts

### Architecture Concepts

- persistence boundaries
- repositories
- optimistic locking
- dependency direction

### Testing

- repository integration tests
- concurrent update tests
- transaction rollback tests

### Exercises

- design repository API
- compare optimistic vs pessimistic locking
- explain why repositories belong to infrastructure

### Interview Topics

- Repository pattern in Go
- Why not expose sqlc directly?
- Optimistic locking
- Error wrapping

### Definition of Done

- persistence layer complete
- integration tests pass
- concurrent updates handled safely
- roadmap updated

---

## Maintenance Before Milestone 6

Status

✅ Completed

### Lesson M6-pre - Test DB & PostgreSQL Error Codes

### Goal

Resolve small persistence-quality issues before adding authentication and authorization.

### Scope

- introduce a separate integration test database configuration, preferably `TEST_DATABASE_URL`
- ensure tests never use the same database as the running application or future production deployments
- keep `DATABASE_URL` for the application and migration command
- define shared string-typed PostgreSQL SQLSTATE constants
- replace raw PostgreSQL error code strings in employee, product and production PostgreSQL stores

### Business Value

Authentication will introduce more security-sensitive flows. Before that, tests should be isolated from application data and persistence error handling should be easier to read and review.

### Design Notes

Integration tests currently use `DATABASE_URL` and clean tables directly. This is acceptable for early local development but should not continue once the project moves toward production-like security work.

PostgreSQL error codes are standardized SQLSTATE strings. They should remain string values, not `iota` enums, because PostgreSQL returns strings such as `23505` and `23503`.

Preferred shape:

```go
type SQLState string

const (
    UniqueViolation     SQLState = "23505"
    ForeignKeyViolation SQLState = "23503"
    CheckViolation      SQLState = "23514"
    NotNullViolation    SQLState = "23502"
    InvalidTextValue    SQLState = "22P02"
)
```

### Definition of Done

- ✅ tests use `TEST_DATABASE_URL` or skip when it is not configured/available
- ✅ application runtime still uses `DATABASE_URL`
- ✅ PostgreSQL error code magic strings are removed from `internal/employees`, `internal/products` and `internal/production`
- ✅ shared constants live in a small internal package
- ✅ tests pass
- ✅ build passes
- ✅ lint passes
- ✅ roadmap updated

### Completion Review

Integration tests for PostgreSQL stores now read `TEST_DATABASE_URL` directly and skip when it is not configured. Application configuration remains unchanged and continues to use `DATABASE_URL`.

PostgreSQL SQLSTATE codes now live in `internal/postgres` as string-typed constants. Employee, product and production PostgreSQL stores map errors through those constants instead of raw strings.

---

## Milestone 6 - Authentication & Authorization

Status

🚧 In Progress

### Lessons

- **L6.1** — Identity, Passwords & Login ✅
- **L6.2** — JWT Issuing & Verification ✅
- **L6.3** — Authentication Middleware & Request Context ✅
- **L6.4** — Role-Based Authorization & Protected Endpoints ✅
- **L6.5** — Security Review, OpenAPI & Milestone Review ✅

### Lesson 6.1 Scope

Introduce the authentication vertical slice without protecting business endpoints yet.

#### Business Context

MES Lite now stores real production data. Before role-based access can exist, the system needs a way to represent application users and verify credentials.

#### Problem

The application has employees, products and production entries, but no security identity. Anyone who can reach the API can call every endpoint.

#### Design Discussion

Authentication starts with a small `auth` slice. A user is not the same as an employee: an employee represents a production worker in the business domain, while an auth user represents someone allowed to access the application.

This lesson verifies email/password credentials and returns an opaque access token. JWT signing, request middleware and authorization are intentionally postponed to keep the lesson focused on identity, password hashing and credential errors.

#### Go Concepts

- custom string role type
- password hashing with explicit error handling
- `crypto/rand` for unpredictable token bytes
- authentication errors translated to HTTP `401 Unauthorized`

#### Architecture Concepts

- authentication as its own vertical slice
- separating business employees from security users
- service boundary for credential verification
- transport response hides password hashes with `json:"-"`

### Lesson 6.1 Completion Notes

#### Business Context

MES Lite needs application users before it can protect production and master-data endpoints. Security identity is separate from employees because not every employee is necessarily an API user, and not every API user performs production work.

#### Problem

The API had no login flow. Any client could call every endpoint without proving identity.

#### Design Discussion

The first authentication step introduces a small `internal/auth` vertical slice. It verifies credentials using an auth service and returns an access token, while JWT signing, token verification, request context and endpoint protection are postponed to later lessons.

Passwords are hashed with `bcrypt` from `golang.org/x/crypto`. This is an intentional dependency because the Go standard library does not provide a production password-hashing algorithm. Raw password hashing with SHA-256 would be faster for attackers and inappropriate for real authentication.

The running server can create a bootstrap admin only when both `AUTH_BOOTSTRAP_EMAIL` and `AUTH_BOOTSTRAP_PASSWORD` are configured. There is no default admin password.

#### Go Concepts

- custom string type for roles
- constructor validation for security users
- `json:"-"` to prevent password hashes from being serialized
- `crypto/rand` plus base64 encoding for unpredictable opaque token bytes
- sentinel authentication errors translated to HTTP responses

#### Architecture Concepts

- authentication vertical slice added as `internal/auth`
- auth users kept separate from business employees
- auth handler owns HTTP parsing and error translation
- auth service owns credential verification

#### Implementation

- Added `auth.User`, `auth.Role` and role validation.
- Added bcrypt password hashing and password verification.
- Added `auth.Store` and in-memory store for the first lesson.
- Added `auth.Service.Login`.
- Added `POST /auth/login`.
- Added optional bootstrap admin configuration.
- Wired the auth handler into the server composition root.

#### Tests

- Added user password-hashing tests.
- Added login service tests for valid credentials, wrong password and inactive users.
- Added HTTP handler tests for successful and rejected login.
- Updated server route setup tests.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

No existing business endpoint was protected in this lesson. That avoids mixing password handling, JWT parsing and authorization decisions in one step.

#### Code Review

An experienced Go engineer would approve the separation between auth users and employees, the absence of plaintext password storage and the explicit bootstrap configuration. The main caveat is intentional lesson scope: returned tokens are not yet verified by middleware and users are not durable yet.

#### Exercises

- Explain why an employee and an auth user are different domain concepts.
- Add a table test for every supported role.
- Try replacing bcrypt with SHA-256 and explain why that would be weaker for password storage.

#### Interview Questions

- What is the difference between authentication and authorization?
- Why should passwords be hashed with bcrypt/argon2 instead of SHA-256?
- Why should APIs return `401 Unauthorized` for invalid login credentials?
- Why is a default admin password dangerous?

#### Roadmap Update

- Lesson 6.1 completed.
- Current lesson moved to Lesson 6.2.
- Known technical debt updated for temporary auth token and user persistence gaps.

### Lesson 6.2 Scope

Replace temporary opaque login tokens with signed JWT access tokens and add token verification inside the auth package.

#### Business Context

After login, clients need a portable credential they can send with later API requests. The server must be able to verify that credential without storing every active session in memory.

#### Problem

Lesson 6.1 returned an unpredictable token, but the application had no way to validate token claims. That was enough to teach credential verification, but not enough for middleware or protected endpoints.

#### Design Discussion

JWT is introduced through a concrete `auth.TokenManager`. The token manager signs access tokens with HMAC SHA-256 and verifies token signature, expiry and required claims. A concrete type is simpler than introducing an interface before a second implementation exists.

The JWT secret is required at server startup through `JWT_SECRET`. There is intentionally no default secret because a default signing key would make local convenience look like production security.

#### Go Concepts

- concrete dependency injection
- time-based token expiry
- wrapping security errors without leaking parser internals
- validating required claims before trusting token data

#### Architecture Concepts

- JWT issuing and verification belong to the auth slice
- token claims become the future request principal
- configuration controls secrets, not source code defaults
- middleware is postponed until token verification is independently tested

### Lesson 6.2 Completion Notes

#### Business Context

Users can now receive a signed access token after login. This prepares the API for authenticated requests without requiring server-side session storage.

#### Problem

The previous login token was opaque and not verifiable by the application. Middleware could not safely identify a caller from it.

#### Design Discussion

The auth package now owns a `TokenManager` that issues and verifies JWTs. It uses `github.com/golang-jwt/jwt/v5`, which was already present indirectly and is now a direct dependency because the project uses it directly.

JWTs are signed using `HS256`. This is simple and appropriate for a single service as long as `JWT_SECRET` is strong and private. Asymmetric signing can be revisited later if multiple services need to verify tokens without sharing the signing secret.

#### Go Concepts

- concrete types before interfaces
- `time.Time` and expiry claims
- error wrapping with a stable `ErrInvalidToken`
- table-like security tests through focused token cases

#### Architecture Concepts

- token issuing stays behind the auth service
- token verification returns a small `Principal`
- server composition root owns security configuration
- no endpoint protection until middleware has a clear boundary

#### Implementation

- Added `auth.TokenManager`.
- Added JWT issue and verify behavior.
- Added `auth.Principal` extracted from token claims.
- Changed login to return a signed JWT access token.
- Added `JWT_SECRET` configuration.
- Required `JWT_SECRET` at server startup.
- Promoted `github.com/golang-jwt/jwt/v5` to a direct dependency.

#### Tests

- Added JWT issue-and-verify tests.
- Added wrong-secret rejection test.
- Added short-secret rejection test.
- Updated login service tests to verify returned JWTs.
- Updated server and handler tests to use explicit test token managers.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

Removed temporary opaque token generation from `auth.Service`. Token generation is now centralized in `TokenManager`, making the next middleware lesson smaller.

#### Code Review

An experienced Go engineer would approve requiring an explicit JWT secret and keeping token verification inside the auth package. The main remaining gap is expected: tokens are verifiable but not yet enforced on requests.

#### Exercises

- Decode a returned JWT at jwt.io and identify which fields are claims, not secrets.
- Change the signing secret in a test and explain why verification fails.
- Explain why JWT payloads must not contain password hashes or sensitive data.

#### Interview Questions

- What problem does a JWT solve?
- What is the difference between signing and encrypting a JWT?
- Why should the server verify the signing algorithm?
- When would asymmetric JWT signing be preferable to HS256?

#### Roadmap Update

- Lesson 6.2 completed.
- Current lesson moved to Lesson 6.3.
- Known technical debt updated: JWT verification exists, middleware enforcement remains pending.

### Lesson 6.3 Scope

Add authentication middleware that verifies `Authorization: Bearer <token>` and stores the authenticated principal in request context.

#### Business Context

After login returns a signed token, business endpoints must reject anonymous callers. Workers, leaders and administrators should prove identity before touching production or master data.

#### Problem

JWTs could be issued and verified manually, but no request pipeline enforced token verification. Every business endpoint was still public.

#### Design Discussion

Authentication is implemented as route-scoped middleware using Fuego's `OptionMiddleware`. Public endpoints stay public: health, readiness, version and login. Business endpoints require a valid bearer token.

The middleware uses standard `net/http` shapes: `func(http.Handler) http.Handler`. This is idiomatic Go because middleware composes around the standard library instead of depending on framework-specific magic.

#### Go Concepts

- HTTP middleware as function composition
- request headers and `Authorization: Bearer`
- request context values for request-scoped identity
- unexported context key type to avoid collisions

#### Architecture Concepts

- authentication happens before handlers
- handlers remain focused on business work
- principal propagation is request-scoped, not global state
- role checks are postponed to a separate authorization lesson

### Lesson 6.3 Completion Notes

#### Business Context

MES Lite business endpoints now require callers to authenticate with a valid JWT before creating or reading production-related data.

#### Problem

The auth package could verify JWTs, but the server did not enforce authentication on requests.

#### Design Discussion

The solution adds a small `auth.Middleware` type. It extracts a bearer token from the `Authorization` header, verifies it with `TokenManager`, stores the resulting `Principal` in request context and calls the next handler. Missing or invalid tokens return `401 Unauthorized` before business handlers execute.

This keeps authentication separate from authorization. Lesson 6.3 answers "who are you?" Lesson 6.4 will answer "are you allowed to do this?"

#### Go Concepts

- `http.Handler` middleware chaining
- request-scoped context values
- unexported context key types
- header parsing with `strings`
- early HTTP rejection before handler execution

#### Architecture Concepts

- route-scoped middleware through Fuego `OptionMiddleware`
- public infrastructure endpoints remain unauthenticated
- protected business endpoints require JWT authentication
- authenticated principal is available without global state

#### Implementation

- Added `auth.Middleware`.
- Added `Authenticate` middleware.
- Added `ContextWithPrincipal` and `PrincipalFromContext`.
- Protected employee, product and production-entry routes.
- Kept `/auth/login`, `/health`, `/ready` and `/version` public.

#### Tests

- Added middleware success test proving the principal reaches request context.
- Added middleware tests for missing and invalid tokens.
- Updated production-entry server route test to send a bearer token.
- Added server test proving production registration requires authentication.
- Added server test proving login remains public.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

Business handlers were not changed. Authentication belongs to middleware, not to every individual handler.

#### Code Review

An experienced Go engineer would approve the standard-library middleware shape and route-scoped protection. The main remaining gap is authorization: all authenticated users currently have the same access.

#### Exercises

- Add a test proving `/employees` returns `401` without a token.
- Explain why context values are acceptable for request-scoped identity but not business data.
- Trace a request from `Authorization` header to `PrincipalFromContext`.

#### Interview Questions

- What is HTTP middleware in Go?
- Why use an unexported type as a context key?
- What data should and should not be stored in `context.Context`?
- Why separate authentication middleware from authorization checks?

#### Roadmap Update

- Lesson 6.3 completed.
- Current lesson moved to Lesson 6.4.
- HTTP middleware marked complete in the Knowledge Matrix.
- Known technical debt updated: authentication enforcement exists, role-based authorization remains pending.

### Lesson 6.4 Scope

Add role-based authorization checks to protected business routes.

#### Business Context

Different MES Lite users have different responsibilities. A production worker should register completed work, but should not administer employees. A manager should maintain products, while an administrator controls user-sensitive master data.

#### Problem

Lesson 6.3 proved user identity, but every authenticated role could call every protected endpoint. Authentication answered "who are you?" but authorization did not yet answer "are you allowed to do this?"

#### Design Discussion

RBAC starts with a simple route permission matrix:

- Admin: all business endpoints.
- Manager: product maintenance, master-data reads and production registration.
- Leader: master-data reads and production registration.
- Worker: production registration only.

This is implemented as route-scoped middleware after authentication. Missing or invalid identity returns `401 Unauthorized`; valid identity with insufficient permissions returns `403 Forbidden`.

#### Go Concepts

- variadic function parameters for allowed roles
- map set pattern with `map[Role]struct{}`
- middleware ordering
- HTTP status semantics: `401` vs `403`

#### Architecture Concepts

- authorization policy belongs near route composition for now
- authentication and authorization remain separate middleware steps
- RBAC is explicit instead of hidden in handlers
- route-level policies are easy to review during code review

### Lesson 6.4 Completion Notes

#### Business Context

MES Lite now distinguishes what authenticated users are allowed to do based on their role.

#### Problem

All authenticated users previously had the same permissions. That would let workers administer employees or products, which does not match business responsibilities.

#### Design Discussion

Authorization is implemented as `RequireRole`, a middleware factory that receives allowed roles and checks the authenticated principal already placed in context by `Authenticate`.

The policy is intentionally route-level and explicit. A more complex permission system, policy engine or database-backed permissions table would be premature. The current business rules are simple enough for direct route composition.

#### Go Concepts

- variadic parameters with `roles ...Role`
- efficient membership checks with a map set
- middleware ordering with authentication before authorization
- stable HTTP semantics for security failures

#### Architecture Concepts

- RBAC at the HTTP boundary
- explicit route permission matrix
- separation between identity proof and permission checks
- handlers remain free of security boilerplate

#### Implementation

- Added `auth.Middleware.RequireRole`.
- Added `403 Forbidden` response for insufficient permissions.
- Applied admin-only permissions to employee mutation routes.
- Allowed admins, managers and leaders to read master data.
- Allowed admins and managers to maintain products.
- Allowed admins, managers, leaders and workers to register production.

#### Tests

- Added middleware tests for matching role, missing principal and forbidden role.
- Added server test proving worker cannot create employees.
- Added server test proving leader can list products.
- Added server test proving worker can register production.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

No handler changes were required. This confirms the middleware boundary is doing the right job.

#### Code Review

An experienced Go engineer would approve the explicit route policy for this project size. The main follow-up is milestone-level review: decide whether in-memory auth users are acceptable for now or whether durable auth-user persistence must be completed before Milestone 6 closes.

#### Exercises

- Add a test proving a manager can create a product.
- Add a test proving a worker cannot list employees.
- Explain why `403 Forbidden` is different from `401 Unauthorized`.

#### Interview Questions

- What is RBAC?
- Why should authentication and authorization be separate concepts?
- When would route-level authorization become insufficient?
- How would you model permissions if roles became too coarse?

#### Roadmap Update

- Lesson 6.4 completed.
- Current lesson moved to Lesson 6.5.
- Known technical debt updated: role checks exist; durable auth-user persistence remains under review.

### Lesson 6.5 Scope

Harden authentication for milestone completion by making auth users durable, documenting OpenAPI bearer security and reviewing the milestone.

#### Business Context

Authentication is not useful if users disappear after an application restart. The first production-safe identity path must persist at least the bootstrap administrator used to access the secured API.

#### Problem

Auth users were still in-memory while the rest of the business data was PostgreSQL-backed. This meant login worked only until restart, and Milestone 6 would have ended with a major security usability gap.

#### Design Discussion

The project now persists auth users in PostgreSQL through the auth vertical slice. Full user-management CRUD is intentionally postponed. The current business requirement is secure application access, and the minimum durable path is an idempotent bootstrap admin plus database-backed login.

OpenAPI now advertises a `bearerAuth` security scheme and marks protected routes as requiring JWT bearer authentication.

#### Go Concepts

- persistence adapter reuse with sqlc
- idempotent startup behavior
- error translation for duplicate users
- direct dependency promotion when a previously indirect package is imported

#### Architecture Concepts

- durable security identity
- bootstrap workflow instead of default credentials
- OpenAPI as part of the security contract
- concrete persistence before full management workflows

### Lesson 6.5 Completion Notes

#### Business Context

MES Lite now has a restart-safe authentication path. A configured bootstrap administrator is stored in PostgreSQL and can log in after restarts.

#### Problem

In-memory auth users were acceptable for learning login mechanics but not acceptable for completing an authentication milestone.

#### Design Discussion

Auth-user persistence was added without creating full user-management CRUD. This is a deliberate minimal design: CRUD would introduce extra authorization and lifecycle rules that the business has not needed yet.

Bootstrap admin creation is idempotent. If the configured email already exists, startup leaves the existing user unchanged instead of overwriting the password on every restart.

#### Go Concepts

- sqlc-generated auth queries
- PostgreSQL `bytea` for password hashes
- idempotent bootstrap with `EXISTS`
- domain error mapping from PostgreSQL constraint errors

#### Architecture Concepts

- auth vertical slice owns auth persistence
- generated auth database types stay below the auth package boundary
- server composition root selects durable auth store
- OpenAPI security scheme documents runtime security

#### Implementation

- Added migration `0005_create_auth_users.sql`.
- Added auth sqlc queries and generated `internal/auth/authdb`.
- Added `auth.PostgresStore`.
- Added idempotent `EnsureBootstrapAdmin`.
- Switched `cmd/server` from in-memory auth store to PostgreSQL auth store.
- Added OpenAPI `bearerAuth` security scheme.
- Added bearer security requirements to protected business routes.

#### Tests

- Added auth PostgreSQL store integration tests.
- Tested save/find by email.
- Tested duplicate email maps to `ErrAlreadyExists`.
- Tested missing email maps to `ErrNotFound`.
- Tested bootstrap admin creation is idempotent and does not overwrite an existing password.
- Verified with `sqlc generate`.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

The server now uses PostgreSQL-backed auth users. In-memory auth store remains for fast tests, matching the existing employee/product testing pattern.

#### Code Review

An experienced Go engineer would approve completing durable auth persistence before closing the milestone. The remaining limitation is deliberate: there is no user-management CRUD yet. That should be introduced only when the administrator workflow is defined.

#### Exercises

- Add a test proving invalid role strings are rejected by the database.
- Explain why bootstrap should not overwrite an existing admin password.
- Inspect the OpenAPI document and identify `bearerAuth`.

#### Interview Questions

- Why is idempotent startup important?
- Why should password hashes be stored as bytes instead of plaintext?
- Why keep generated sqlc auth types out of HTTP handlers?
- What risks appear when adding user-management CRUD?

#### Roadmap Update

- Lesson 6.5 completed.
- Milestone 6 completed.
- Current milestone moved to Milestone 7.
- Known technical debt updated: durable auth-user persistence completed; full auth-user management postponed.

### Milestone 6 Review

#### Architecture Review

An experienced Go engineer would approve the milestone direction. Authentication is isolated in `internal/auth`, password hashes are never serialized, JWT issuing and verification are centralized, middleware uses standard `net/http` composition and authorization policy is explicit at route registration.

The main trade-off is that route-level RBAC is simple but not infinitely flexible. This is appropriate for the current business rules. A policy service or permission table would be premature until roles become too coarse.

#### Code Review

The code remains small and idiomatic. Handlers do not parse tokens or check roles. Generated SQL types stay below the auth package boundary. Bootstrap has no default password and requires explicit environment configuration.

One caveat: auth-user management CRUD does not exist yet. This is not hidden debt; it is intentionally postponed until an administrator workflow is designed.

#### Refactoring

Auth persistence follows the established vertical-slice repository adapter pattern. No broad auth abstraction was introduced because there is only one token implementation and one durable store implementation.

#### Interview Review

You should now be able to explain authentication vs authorization, `401` vs `403`, JWT signing vs encryption, bearer-token middleware, request-scoped context values, RBAC trade-offs and why password hashing needs a dedicated algorithm such as bcrypt.

#### Completion Criteria

- Login implemented.
- JWT issuing and verification implemented.
- Business endpoints require bearer authentication.
- Role-based route authorization implemented.
- Auth users persist in PostgreSQL.
- Bootstrap admin creation is explicit and idempotent.
- OpenAPI security scheme is configured.
- Tests, build, lint and sqlc generation pass.
- Roadmap updated.

### Goal

Secure the application.

### Business Value

Different users have different permissions.

### Features

- login
- JWT authentication
- role-based authorization
- protected endpoints

### Go Concepts

- middleware
- request context
- context values
- HTTP headers
- cryptographic randomness (overview)

### Standard Library

- context
- crypto
- net/http

### Architecture Concepts

- authentication
- authorization
- RBAC
- middleware pipeline

### Testing

- authenticated endpoints
- unauthorized access
- middleware tests

### Exercises

- implement authentication middleware
- explain Context usage
- discuss JWT trade-offs

### Interview Topics

- Authentication vs Authorization
- Why Context?
- Middleware chaining
- JWT advantages/disadvantages

### Definition of Done

- protected API
- role-based access
- tests passing
- OpenAPI updated

---

## Milestone 7 - Production Orders

Status

🚧 In Progress

### Lessons

- **L7.1** — Production Order Aggregate ✅
- **L7.2** — Order Persistence & Reference Integrity
- **L7.3** — Create & Read Production Orders API
- **L7.4** — Status Transitions & Employee Assignment
- **L7.5** — Transactional Consistency & Milestone Review

### Lesson 7.1 Scope

Model production orders as the first explicit aggregate in the system.

#### Business Context

Production managers need to plan work before employees report completed production. A production order answers which products should be produced, how many units are planned for each product and who may work on the order.

#### Problem

The system can record completed production, but it cannot represent planned production. Adding persistence or HTTP first would risk creating CRUD-shaped data structures before the business invariants are clear.

#### Design Discussion

The first lesson keeps the work inside a new `orders` vertical slice and focuses on the domain model. `Order` is the aggregate root. It owns order lines, status, assigned employees and timestamps. State changes happen through methods so business rules stay close to the data they protect.

Persistence, reference validation against products/employees and HTTP endpoints are intentionally postponed. This lets aggregate rules be tested quickly without database or framework noise.

#### Go Concepts

- composition through small structs
- custom methods that preserve invariants
- status custom type with explicit validation
- time normalization for business timestamps

#### Architecture Concepts

- aggregate root
- business invariants inside the domain model
- avoiding anemic domain models for rule-heavy workflows

### Lesson 7.1 Completion Notes

#### Business Context

Production managers need to plan work before production workers report completed output. A production order defines the products to make, the planned quantity for each product and the employees assigned to the work.

#### Problem

The application could record completed production entries, but it had no model for planned production work. Starting with database tables or HTTP handlers would have encouraged CRUD-first design before identifying the business rules.

#### Design Discussion

Added `internal/orders` as a new vertical slice and modeled `orders.Order` as the aggregate root. The aggregate owns `OrderLine` children, status transitions and employee assignment rules through methods instead of letting callers mutate fields freely.

The model was refined after review because a real production order may contain multiple product types, such as 2 shafts and 4 filters. `ProductSKU` and `PlannedQuantity` therefore belong to `OrderLine`, not directly to `Order`.

Persistence and HTTP were intentionally postponed. This keeps the first lesson focused on aggregate design and makes the invariant tests fast and clear.

#### Go Concepts

- custom `Status` type with explicit validation
- child value type through `OrderLine`
- methods that preserve invariants through copy-then-assign updates
- `time.Time` normalization to UTC
- sentinel errors for invalid state and invalid transitions

#### Architecture Concepts

- aggregate root as the consistency boundary
- rich domain behavior for rule-heavy workflows
- vertical slice package for production orders
- explicit postponement of persistence until the domain model is stable

#### Implementation

- Added `orders.Order` with order lines, status, assigned employees and timestamps.
- Added `orders.OrderLine` with product SKU and planned quantity.
- Added order statuses: draft, released, in-progress, completed and cancelled.
- Added `NewOrder`, `Validate`, `AssignEmployee`, `Release`, `Start`, `Complete` and `Cancel`.
- Added `NewOrderLine` and `OrderLine.Validate`.
- Enforced at least one line, valid line quantities, no duplicate product SKUs, required identifiers, valid status, timestamps and assignment rules.

#### Tests

- Added aggregate tests for constructor normalization and validation.
- Added order-line validation tests.
- Added duplicate-product and defensive-copy tests.
- Added status validation tests.
- Added assignment tests including idempotent assignment and closed-order rejection.
- Added status-transition tests for valid and invalid lifecycle moves.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

No existing packages required changes. The new slice is isolated so later lessons can add persistence and HTTP without leaking order rules into handlers.

#### Code Review

An experienced Go engineer would approve the small scope and invariant-focused tests. The main caveat is that fields are still exported for future serialization and persistence convenience, so code can technically bypass methods. This matches the current project style, but persistence reconstruction should be designed carefully in Lesson 7.2.

#### Exercises

- Explain why `Order` is the aggregate root instead of `AssignedEmployee`.
- Add a test for cancelling an already cancelled order and decide whether idempotency is desirable.
- Compare this model with an anemic DTO that only has public fields and no methods.

#### Interview Questions

- What is an aggregate root?
- Why should business invariants live near the data they protect?
- Why does Go prefer composition and methods over inheritance?
- When are exported fields acceptable in a domain type, and what trade-off do they create?

#### Roadmap Update

- Lesson 7.1 completed.
- Current lesson moved to Lesson 7.2.
- Architecture `Aggregates` marked complete in the Knowledge Matrix.

### Goal

Introduce production planning.

This is the first milestone centered on business rules instead of CRUD.

### Business Value

Production can be planned before execution.

### Features

- create production order
- assign products
- planned quantity
- status
- assign employees

### Go Concepts

- embedding
- composition
- custom methods
- invariants

- time
- errors

### Architecture Concepts

- aggregate root
- domain model
- business invariants
- transactional consistency

### Testing

- aggregate tests
- business rule tests
- integration tests

### Exercises

- identify aggregate root
- identify invariants
- compare anemic vs rich domain models

### Interview Topics

- Aggregate
- Invariants
- Composition vs inheritance
- Embedding

### Definition of Done

- production orders functional
- business rules enforced
- tests passing

---

## Milestone 8 - Reporting

Status

⬜ Not Started

### Goal

Provide useful business reports.

### Business Value

Managers can monitor production without Excel.

### Features

- daily production
- employee productivity
- product statistics
- filtering
- exports

### Go Concepts

- io.Writer
- streaming
- buffering
- efficient JSON generation

### Standard Library

- io
- bytes
- encoding/json

### Architecture Concepts

- CQRS
- read models
- query services

### SQL Concepts

- CTE
- GROUP BY
- aggregates
- window functions
- indexes
- execution plans (introduction)

### Testing

- report integration tests
- performance tests
- query correctness

### Exercises

- design read model
- optimize SQL query
- explain CQRS benefits

### Interview Topics

- CQRS
- Why separate reads?
- SQL optimization
- Streaming responses

### Definition of Done

- reports complete
- SQL optimized
- tests passing

---

## Milestone 9 - CSV Import

Status

⬜ Not Started

### Goal

Import historical Excel data.

This milestone introduces one of the most important Go topics:

Streaming.

### Business Value

The company can migrate historical production data.

### Features

- CSV upload
- validation
- batch import
- import summary
- partial failure reporting

### Go Concepts

- io.Reader
- bufio
- csv.Reader
- streaming
- memory efficiency
- defer revisited

### Standard Library

- io
- encoding/csv
- bufio
- os

### Architecture Concepts

- import pipeline
- validation pipeline
- error collection

### Testing

- large file tests
- malformed CSV
- import rollback
- partial import

### Exercises

- import a 1 GB CSV without loading it entirely into memory
- compare streaming vs reading the whole file
- benchmark both approaches

### Interview Topics

- Why io.Reader is fundamental in Go?
- Streaming vs buffering
- Memory usage
- Designing import pipelines

### Definition of Done

- streaming importer implemented
- large files supported
- memory usage acceptable
- tests passing

---

## Milestone Review (Milestones 0-9)

At this point the application should already be usable by a small company.

The following concepts should now be understood.

### Language

- modules
- packages
- functions
- methods
- structs
- pointers
- interfaces
- visibility
- custom types
- errors
- context
- io.Reader
- defer
- time

### Standard Library

- net/http
- context
- errors
- io
- encoding/json
- encoding/csv
- bytes
- strings
- time

### Architecture

- Vertical Slice
- dependency injection
- repositories
- aggregates
- CQRS
- package boundaries
- explicit dependencies

### Database

- PostgreSQL
- sqlc
- pgx
- transactions
- optimistic locking

### Web

- Fuego
- OpenAPI
- middleware
- authentication
- authorization

### Testing

- unit tests
- integration tests
- table tests
- httptest
- Testcontainers

### Interview Readiness Checkpoint

The AI should now conduct a comprehensive review.

The user should be able to answer questions such as:

- Why does Go avoid inheritance?
- When should an interface be introduced?
- Why is `context.Context` passed as the first argument?
- Why prefer `sqlc` over an ORM?
- Why organize code by business capability instead of technical layers?
- How does error wrapping work?
- When would you choose a pointer receiver?
- How would you structure a production Go service?

If significant gaps are found, they should be addressed before moving to Milestones 10–14.

---

## Milestone 10 - Background Jobs & Concurrency

Status

⬜ Not Started

### Goal

Learn idiomatic concurrency by introducing asynchronous processing.

Concurrency is one of Go's defining features and should be understood deeply.

### Business Value

Long-running tasks should not block HTTP requests.

Examples:

- CSV imports
- report generation
- notifications
- scheduled maintenance

### Features

- job queue
- worker pool
- progress tracking
- cancellation
- retry strategy

### Go Concepts

- goroutines
- channels
- select
- WaitGroup
- errgroup
- cancellation
- channel ownership
- channel closing

### Standard Library

- sync
- context
- time

### Architecture Concepts

- asynchronous processing
- worker pool
- producer-consumer
- pipeline

### Testing

- concurrent tests
- cancellation tests
- race detection

### Exercises

- implement worker pool
- implement cancellation
- explain channel ownership
- compare goroutines with OS threads

### Interview Topics

- What is a goroutine?
- Buffered vs unbuffered channels
- When should a channel be closed?
- WaitGroup vs errgroup
- Common concurrency mistakes

### Definition of Done

- asynchronous jobs implemented
- cancellation supported
- race detector passes
- tests passing

---

## Milestone 11 - Machine Integration & Synchronization

Status

⬜ Not Started

### Goal

Safely process concurrent events coming from production machines.

### Business Value

Prepare the system for future CNC and production machine integration.

### Features

- fake machine API
- production events
- event processing
- duplicate detection
- idempotency

### Go Concepts

- Mutex
- RWMutex
- atomic operations
- sync.Once
- sync.Map
- race detector

### Standard Library

- sync
- sync/atomic

### Architecture Concepts

- event-driven architecture
- idempotency
- synchronization
- eventual consistency (introduction)

### Testing

- race detector
- stress tests
- concurrent processing tests

### Exercises

- intentionally introduce a race condition
- fix it using Mutex
- compare Mutex vs channels
- discuss atomic operations

### Interview Topics

- Mutex vs channels
- Atomic operations
- Data races
- Event-driven systems
- Idempotency

### Definition of Done

- concurrent processing is safe
- race detector clean
- idempotent processing implemented

---

## Milestone 12 - Observability

Status

⬜ Not Started

### Goal

Make the application observable in production.

### Business Value

Operators can understand system health and diagnose problems quickly.

### Features

- structured logging
- metrics
- tracing
- health endpoint
- readiness endpoint

### Go Concepts

- slog
- structured logging
- context propagation revisited

### Standard Library

- log/slog

### Additional Technologies

- OpenTelemetry
- Prometheus

### Architecture Concepts

- observability
- correlation IDs
- diagnostics

### Testing

- log verification
- health endpoint tests
- tracing integration

### Exercises

- add structured logs
- trace one request through the application
- explain why plain printf logging is insufficient

### Interview Topics

- Logs vs Metrics vs Traces
- Structured logging
- Correlation IDs
- Observability

### Definition of Done

- structured logging implemented
- metrics available
- traces exported
- health checks complete

---

## Milestone 13 - Performance Engineering

Status

⬜ Not Started

### Goal

Understand how Go applications perform and how to optimize them.

Optimization must be based on measurement rather than assumptions.

### Business Value

The application remains responsive as usage grows.

### Features

- profiling
- benchmarking
- allocation analysis

### Go Concepts

- benchmarking
- pprof
- escape analysis
- stack vs heap
- garbage collector
- scheduler
- memory allocations

### Standard Library

- testing
- runtime
- runtime/pprof

### Architecture Concepts

- performance trade-offs
- scalability
- optimization strategy

### Testing

- benchmarks
- allocation benchmarks
- profiling sessions

### Exercises

- benchmark two implementations
- reduce allocations
- analyze escape analysis output
- profile a slow endpoint

### Interview Topics

- Escape analysis
- Garbage collector
- Scheduler
- Heap vs stack
- Benchmarking

### Definition of Done

- benchmarks written
- profiling completed
- unnecessary allocations reduced

---

## Milestone 14 - Production Readiness

Status

⬜ Not Started

### Goal

Prepare the application for deployment as if it were a real production system.

### Business Value

The project reaches production quality.

### Features

- Docker image
- CI/CD improvements
- graceful deployment
- versioning
- release process

### Go Concepts

- build flags
- build tags
- embedding
- linker flags

### Standard Library

- embed
- os
- flag

### Architecture Concepts

- deployment
- configuration management
- release strategy

### Testing

- end-to-end tests
- deployment verification
- smoke tests

### Exercises

- create production image
- explain build tags
- version the application

### Interview Topics

- How would you deploy a Go service?
- Why are Go binaries easy to deploy?
- Static linking
- Cross compilation

### Definition of Done

- production deployment possible
- CI green
- documentation updated
