# ROADMAP.MVP.API.md

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

**Version:** 2.27
**Status:** IN PROGRESS
**Current milestone:** 9 - CSV Import
**Current lesson:** L9.3 - Batch Persistence & Transaction Strategy
**Completed milestones:**
- Milestone 0
- Milestone 1
- Milestone 2
- Milestone 3
- Milestone 4
- Milestone 5
- Milestone 6
- Milestone 7
- Milestone 8
**Next milestone:** 10 - MVP API Completion
**Current branch:** main
**Architecture maturity:** 7.9 / 10
**Go knowledge progress:** 65%
**Interview readiness:** 59%
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
- [x] io
- [x] encoding/json
- [x] encoding/csv
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
- [x] Mutex
- [x] RWMutex
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
- [x] CQRS
- [ ] Event-driven Architecture
- [ ] Observability
- [ ] Production Readiness

**Persistence**
- [x] PostgreSQL
- [x] pgx
- [x] sqlc
- [x] Transactions
- [x] Optimistic Locking
- [x] SQL Optimization

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
- [x] Indexes
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

Reviews recently entered production records.

Reviews completed work.

Corrects mistakes.

---

### Production Manager

Creates or maintains the minimum production-order data needed for registration.

Reviews production records.

Generates reports that replace existing Excel or paper summaries.

Reviews productivity summaries where they are already tracked manually.

---

### Administrator

Maintains users.

Maintains products.

Configures the system.

---

## MVP Scope

The first version is scoped as an Excel/paper replacement. It supports only the capabilities needed to replace current manual production tracking:

- employees
- products
- production entries with a workstation text field
- minimum production-order data needed for registration
- authentication
- reports that replace existing Excel or paper summaries
- CSV import for historical manual data
- manager/leader review of entered production records
- correction of production-entry mistakes without silently editing history

Everything beyond replacing current manual tracking is intentionally postponed to MVP-V2 or later.

---

## Future Scope

MVP-V2 or later milestones may introduce:

- machines
- formal workstation management
- warehouse
- traceability
- quality control
- notifications
- scheduling
- dashboards
- formal background jobs
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

✅ Completed

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
- **L7.1a** — Aggregate Encapsulation Before Persistence ✅
- **L7.2** — Order Persistence & Reference Integrity ✅
- **L7.2a** — Order Lines Collection Refactor ✅
- **L7.3** — Create & Read Production Orders API ✅
- **L7.4** — Status Transitions & Employee Assignment ✅
- **L7.5** — Transactional Consistency & Milestone Review ✅

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

Follow-up Lesson 7.1a tightened encapsulation before persistence/API code could start depending on exported aggregate fields.

#### Code Review

An experienced Go engineer would approve the small scope and invariant-focused tests. The initial caveat was that fields were exported for future serialization and persistence convenience, so code could technically bypass methods. Lesson 7.1a resolved this for the new orders aggregate before persistence reconstruction was introduced.

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
- Current lesson moved to Lesson 7.1a.
- Architecture `Aggregates` marked complete in the Knowledge Matrix.

### Lesson 7.1a Scope

Refactor the new production-order aggregate to hide mutable fields before persistence and HTTP code depend on them.

#### Business Context

Production-order lifecycle rules only remain trustworthy if application code cannot freely mutate order state around the aggregate methods.

#### Problem

`orders.Order` and `orders.OrderLine` initially exposed fields. That matched older project entities, but it weakened the first real aggregate because callers could bypass methods such as `AssignEmployee`, `Release`, `Start`, `Complete` and `Cancel`.

#### Design Discussion

The refactor was done before Lesson 7.2 because no persistence or HTTP adapters exist yet. This is the cheapest moment to improve the API. Existing slices are intentionally left unchanged because making every older entity private would be a broad refactor with more risk and less immediate value.

#### Go Concepts

- unexported struct fields
- exported accessor methods
- defensive slice copies
- value-returning getters for immutable reads

#### Architecture Concepts

- aggregate encapsulation
- preserving invariants through package APIs
- refactoring before persistence contracts harden

### Lesson 7.1a Completion Notes

#### Business Context

Production orders now protect their lifecycle and line invariants through the aggregate API instead of exposing mutable state directly.

#### Problem

Public aggregate fields allowed external callers to mutate `status`, `lines`, timestamps or assigned employees without validation.

#### Design Discussion

`Order` and `OrderLine` now use private fields with explicit accessors. Slice accessors return copies because returning the internal slice would still allow mutation of aggregate state.

No reconstruction function was added yet. Persistence does not exist for orders, so Lesson 7.2 should introduce reconstruction only when loading database rows requires it.

#### Go Concepts

- private fields with exported methods
- defensive copying for slices
- preserving invariants with methods instead of field assignment

#### Architecture Concepts

- aggregate root encapsulation
- public API as a domain boundary
- persistence reconstruction postponed until it solves a real need

#### Implementation

- Made `orders.Order` fields private.
- Made `orders.OrderLine` fields private.
- Added accessors for order ID, lines, status, assigned employees and timestamps.
- Added accessors for order-line product SKU and planned quantity.
- Kept state changes behind aggregate methods.

#### Tests

- Updated order tests to use accessors.
- Added tests proving `Lines()` returns a defensive copy.
- Added tests proving `AssignedEmployees()` returns a defensive copy.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

Only the new orders slice was refactored. Employees, products, production entries and auth users remain public-field types until a focused future refactoring is justified.

#### Code Review

An experienced Go engineer would approve this timing because the aggregate API is now safer before persistence and HTTP contracts are created. The main follow-up is to design a careful order reconstruction function in Lesson 7.2 if sqlc rows need to rebuild persisted aggregate state.

#### Exercises

- Explain why returning a slice directly can break aggregate invariants.
- Add a test proving `OrderLine` cannot be mutated through `Lines()` from another package.
- Design a `RestoreOrder` signature for Lesson 7.2 without implementing it yet.

#### Interview Questions

- Why might an aggregate use private fields in Go?
- What is defensive copying and when is it necessary?
- Why not make every existing entity private in one large refactor?
- How can persistence reconstruct an aggregate without bypassing validation everywhere?

#### Roadmap Update

- Lesson 7.1a completed.
- Current lesson moved to Lesson 7.2.
- Lesson 7.2 remains focused on persistence and reference integrity.

### Lesson 7.2 Scope

Persist production orders in PostgreSQL and enforce product/employee reference integrity at the database boundary.

#### Business Context

Production orders are planning records. Managers need planned orders to survive restarts and remain connected to valid products and assigned employees.

#### Problem

The order aggregate existed only in memory. Persisting only the root row would lose order lines and assignments, while saving multiple tables without a transaction could leave partial orders behind after a failed child insert.

#### Design Discussion

The order aggregate is stored across three tables: `production_orders`, `production_order_lines` and `production_order_assignments`. Lines reference `products(sku)` and assignments reference `employees(id)`.

Saving an order uses an explicit PostgreSQL transaction because the aggregate requires multiple writes. If any line or assignment fails, the root order row is rolled back too.

`RestoreOrder` rebuilds a persisted aggregate through one domain function instead of exposing fields or spreading unchecked struct literals through the repository.

#### Go Concepts

- explicit transaction handling with pgx
- domain reconstruction with validation
- defensive conversion between sqlc rows and domain values
- error translation from PostgreSQL constraint errors

#### Architecture Concepts

- aggregate persistence across multiple tables
- transaction boundary around one aggregate save
- database foreign keys as reference-integrity guardrails
- sqlc-generated types kept below the orders package boundary

### Lesson 7.2 Completion Notes

#### Business Context

Production orders now have durable storage. Planned order lines reference real products, and assigned employees reference real employees.

#### Problem

The orders slice had a tested aggregate but no persistence. A future API would have had nowhere durable to store manager-created orders.

#### Design Discussion

Added SQL-first persistence for the aggregate. The root row, lines and assignments are saved in one transaction. This is the first lesson where a transaction solves a concrete business consistency problem: avoiding a persisted order without all of its required child records.

Reference integrity is enforced with PostgreSQL foreign keys. Application validation still protects aggregate shape, while the database protects cross-table references.

#### Go Concepts

- `pgx` transaction lifecycle with `Begin`, `Commit` and deferred rollback
- reconstructing private-field aggregates through `RestoreOrder`
- converting `pgtype.Timestamptz` to `time.Time`
- translating SQLSTATE errors into domain errors

#### Architecture Concepts

- aggregate persistence boundary
- transaction per aggregate save
- generated SQL package hidden behind repository adapter
- database constraints as the final consistency boundary

#### Implementation

- Added migration `0006_create_production_orders.sql`.
- Added `production_orders`, `production_order_lines` and `production_order_assignments` tables.
- Added foreign keys from order lines to products and assignments to employees.
- Added orders sqlc queries and generated `internal/orders/ordersdb`.
- Added `orders.Store` and `orders.PostgresStore`.
- Added `RestoreOrder` for validated persistence reconstruction.
- Added duplicate order and invalid reference error mapping.

#### Tests

- Added orders PostgreSQL store integration tests.
- Tested save and find with multiple order lines and an assigned employee.
- Tested duplicate order save mapping to `ErrAlreadyExists`.
- Tested missing order mapping to `ErrNotFound`.
- Tested missing product reference rolls back the root order row.
- Tested missing employee reference rolls back the root order row.
- Verified with `make sqlc`.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

No generic transaction manager was introduced. The order store owns the transaction because this is currently the only multi-write aggregate persistence operation.

#### Code Review

An experienced Go engineer would approve the transaction boundary because it protects exactly one aggregate save and does not introduce framework-like transaction abstraction. The main follow-up is API wiring in Lesson 7.3.

#### Exercises

- Explain why saving an order requires a transaction but saving a product does not.
- Add a test that proves a duplicate order line cannot be inserted by bypassing the domain model.
- Draw the three order tables and identify every foreign key.

#### Interview Questions

- What should define a transaction boundary?
- Why keep foreign keys if the application already validates references?
- Why should sqlc-generated row types not become the public domain API?
- How does deferred rollback work after a successful commit?

#### Roadmap Update

- Lesson 7.2 completed.
- Current lesson moved to Lesson 7.2a.
- Persistence `Transactions` marked complete in the Knowledge Matrix.

### Lesson 7.2a Scope

Refactor order lines from a raw slice field into a domain-specific collection type before exposing orders through HTTP.

#### Business Context

Production managers create orders with one or more planned products. The rules for those planned products belong together because they define whether an order is valid.

#### Problem

`Order` still carried line collection rules directly: at least one line, no duplicate product SKUs and defensive copying. As more persistence and API code is added, keeping these rules inside `Order.Validate` would make the aggregate root collect too many low-level details.

#### Design Discussion

Introduced `OrderLines` as a small domain collection. It owns line-specific invariants and exposes `Values()` as a defensive copy. This is not a generic collection helper; it exists because order lines have real business behavior.

Assigned employees remain a raw slice for now. Their rules are still simple and do not yet justify a separate collection type.

#### Go Concepts

- domain-specific collection structs
- variadic constructors
- defensive slice copies
- moving behavior to the type that owns the data

#### Architecture Concepts

- aggregate internals organized by business responsibility
- avoiding generic abstractions while still reducing procedural validation
- refactoring before HTTP contracts harden

### Lesson 7.2a Completion Notes

#### Business Context

Order-line rules are now modeled as their own concept. This keeps multi-product order planning explicit before API request/response contracts are introduced.

#### Problem

`Order.lines` was a raw `[]OrderLine`. The aggregate protected it, but the collection rules did not have their own name or API.

#### Design Discussion

`OrderLines` is a struct with private `values []OrderLine`. It validates that an order has at least one line, that every line is valid and that product SKUs are unique within the order.

The API intentionally returns copies through `Values()` so callers cannot mutate aggregate internals through a slice reference.

#### Go Concepts

- `OrderLines` as a domain collection type
- variadic `NewOrderLines(lines ...OrderLine)` constructor
- copy-on-read with slice values
- keeping the zero value invalid when the business requires data

#### Architecture Concepts

- collection as part of the aggregate model
- domain-specific abstraction justified by behavior
- avoiding premature abstraction for assigned employees

#### Implementation

- Added `orders.OrderLines`.
- Added `NewOrderLines`, `OrderLines.Validate`, `OrderLines.Values` and `OrderLines.Len`.
- Changed `Order.lines` from `[]OrderLine` to `OrderLines`.
- Changed `NewOrder` and `RestoreOrder` to accept `OrderLines`.
- Updated `PostgresStore` to persist and restore through `OrderLines`.

#### Tests

- Added `OrderLines` construction and validation tests.
- Moved duplicate-product tests to `OrderLines`.
- Updated order and repository tests to use `OrderLines`.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

`assignedEmployees` intentionally remains `[]string`. It has some validation, but not enough behavior yet to justify a second collection type.

#### Code Review

An experienced Go engineer would approve this refactor because `OrderLines` has real behavior and reduces aggregate validation detail without introducing a generic abstraction.

#### Exercises

- Explain why `OrderLines` is a struct instead of `type OrderLines []OrderLine`.
- Add a test proving duplicate product SKUs are rejected by `NewOrderLines`.
- Identify what future employee-assignment behavior would justify an `AssignedEmployees` type.

#### Interview Questions

- When is a collection type justified in Go?
- Why can returning a slice expose mutable internal state?
- What is the difference between domain-specific abstraction and generic abstraction?
- Why can an invalid zero value be acceptable for a business value type?

#### Roadmap Update

- Lesson 7.2a completed.
- Current lesson moved to Lesson 7.3.
- Lesson 7.3 remains focused on create/read HTTP endpoints.

### Lesson 7.3 Scope

Expose production-order creation and read endpoints over HTTP without adding status-transition mutation endpoints yet.

#### Business Context

Production managers need an API to create planned production orders, and team leaders need to read those orders before production starts.

#### Problem

Orders were persisted in PostgreSQL, but no HTTP route allowed clients to create or fetch them. The aggregate also has private fields, so returning the domain type directly would produce the wrong API contract.

#### Design Discussion

The orders handler uses explicit request and response DTOs. This keeps HTTP JSON shape separate from the aggregate internals and avoids exposing private-field domain types directly.

Create accepts one or more order lines and optional assigned employee IDs. The application generates the order ID before constructing the aggregate, and PostgreSQL enforces uniqueness. It always creates a draft order. Status-transition endpoints are postponed to Lesson 7.4 so this lesson stays focused on HTTP create/read wiring.

#### Go Concepts

- explicit DTO structs for HTTP transport
- conversion between DTOs and domain values
- optional slices in request bodies
- handler-level error translation

#### Architecture Concepts

- HTTP boundary separated from aggregate internals
- route-level RBAC for production-order planning
- composition root wiring for a new vertical slice endpoint

### Lesson 7.3 Completion Notes

#### Business Context

Managers can now create draft production orders through the API. Authorized users can read persisted production orders.

#### Problem

The order aggregate and PostgreSQL store existed, but clients had no API entry point for planned production.

#### Design Discussion

Added `orders.Handler` with explicit create/read DTOs. The handler generates a production-order ID, constructs `OrderLine`, `OrderLines` and `Order`, assigns optional employees through aggregate methods and delegates persistence to `orders.Store`.

The create request no longer accepts client-provided IDs. Client-generated IDs would be useful for idempotency or offline workflows, but those are not current requirements. Application-generated IDs keep aggregate construction explicit without exposing identity creation to API clients.

Authorization follows the existing RBAC model: admins and managers can create orders; admins, managers and leaders can read orders. Workers cannot create or read planning data through these routes.

#### Go Concepts

- DTO-to-domain conversion
- application-generated UUID-shaped IDs
- non-nil response slices for stable JSON
- error mapping to `400`, `404` and `409`
- route tests with `httptest`

#### Architecture Concepts

- vertical slice owns HTTP handlers and DTOs
- domain aggregate remains hidden behind response mapping
- server composition root wires concrete store and handler
- OpenAPI route generation through Fuego registration

#### Implementation

- Added `orders.Handler`.
- Added `POST /production-orders`.
- Added `GET /production-orders/{id}`.
- Added `CreateOrderRequest`, `CreateOrderLineRequest`, `OrderResponse` and `OrderLineResponse`.
- Added `NewOrderID` for application-generated production-order IDs.
- Added `orders.InMemoryStore` for fast handler/server tests.
- Wired `orders.PostgresStore` and `orders.Handler` in `cmd/server`.
- Registered order routes in `internal/server` with bearer security and RBAC.

#### Tests

- Added handler tests for create, validation, duplicate create, get and not found.
- Added server route tests for manager create, worker create rejection, leader read and unauthenticated rejection.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

No aggregate behavior was changed. The main refactor was transport mapping: HTTP no longer needs to know about private domain fields.

#### Code Review

An experienced Go engineer would approve keeping status transitions out of this lesson. The API now has a minimal create/read surface, and Lesson 7.4 can focus on mutation workflows without mixing route bootstrapping work.

#### Exercises

- Add a test proving a worker cannot read a production order.
- Inspect generated OpenAPI and find both production-order routes.
- Explain why `OrderResponse` is preferable to returning `Order` directly.

#### Interview Questions

- Why separate HTTP DTOs from domain types?
- What belongs in a handler versus an aggregate method?
- Why is route-level authorization acceptable for this project size?
- How does Fuego still rely on standard `net/http` concepts?

#### Roadmap Update

- Lesson 7.3 completed.
- Current lesson moved to Lesson 7.4.
- Lesson 7.4 remains focused on status-transition and assignment workflows.

### Lesson 7.4 Scope

Expose production-order employee assignment and status-transition workflows while preserving aggregate invariants.

#### Business Context

Production managers and leaders need to move planned work through its lifecycle: assign employees, release the order, start work, complete work or cancel the order.

#### Problem

The API could create and read orders, but order state could not change after creation. The persistence store also lacked an update operation for mutable aggregate state.

#### Design Discussion

Added an orders application service so workflow rules do not live in HTTP handlers. The service validates active product references on creation, validates active employees before assignment and delegates lifecycle changes to aggregate methods.

The repository update operation persists mutable order state in one transaction: root status/timestamp plus assignments. Order lines remain immutable after creation in this lesson.

#### Go Concepts

- application-service methods for workflows
- direct DTO-to-command conversion
- string trimming before lookup validation
- error mapping across service, handler and repository boundaries

#### Architecture Concepts

- service boundary for business workflows
- aggregate methods as the source of lifecycle rules
- transactional update for mutable aggregate state
- route-level RBAC for planning mutations

### Lesson 7.4 Completion Notes

#### Business Context

Production orders can now move through their lifecycle and have employees assigned through the API.

#### Problem

After Lesson 7.3, created orders stayed permanently draft unless tests modified the aggregate directly. There was no service boundary for assignment validation or status changes.

#### Design Discussion

`orders.Service` now coordinates order workflows. It validates that planned products are active during create, validates that assigned employees exist and are active, applies aggregate methods and persists the resulting order.

Status transition rules remain inside `Order`. The service decides when to load and save; the aggregate decides whether a transition is legal.

#### Go Concepts

- consumer-owned lookup interfaces for products and employees
- service commands for create and assignment workflows
- explicit error translation from dependency package errors
- `httptest` coverage for mutation routes

#### Architecture Concepts

- application service as workflow coordinator
- aggregate root as invariant owner
- repository update transaction for status and assignments
- RBAC separated from handler logic

#### Implementation

- Added `orders.Service`.
- Added product and employee lookup validation for order workflows.
- Added `Store.Update`.
- Implemented `InMemoryStore.Update`.
- Implemented transactional `PostgresStore.Update`.
- Added sqlc queries for updating order root state and replacing assignments.
- Added assignment endpoint: `POST /production-orders/{id}/assignments`.
- Added status endpoints: `PUT /production-orders/{id}/release`, `start`, `complete` and `cancel`.
- Wired orders service in `cmd/server`.
- Registered mutation routes with RBAC in `internal/server`.

#### Tests

- Added service tests for create reference validation, assignment validation and status transitions.
- Added handler tests for assignment, release/start/complete and cancel.
- Added server route tests for assignment authorization and leader release.
- Added PostgreSQL store update tests for status/assignment persistence, missing order and failed assignment rollback.
- Verified with `make sqlc`.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

Order creation moved from handler-owned domain construction to `orders.Service`. This keeps handlers responsible for HTTP translation and services responsible for business workflow coordination.

#### Code Review

An experienced Go engineer would approve the separation: handlers do not own lifecycle rules, the aggregate still owns transitions and the repository update transaction protects persisted mutable state. A future improvement could add optimistic locking for orders if concurrent planning edits become a real workflow.

#### Exercises

- Add a test proving a worker cannot read a production order.
- Add a test proving `Release` fails when no employee is assigned.
- Design how optimistic locking would apply to order status transitions.

#### Interview Questions

- What belongs in an application service versus an aggregate method?
- Why validate active employees in application code instead of with a foreign key?
- Why does replacing assignments need a transaction?
- When would optimistic locking be needed for aggregate updates?

#### Roadmap Update

- Lesson 7.4 completed.
- Current lesson moved to Lesson 7.5.
- Lesson 7.5 remains focused on transactional consistency and milestone review.

### Lesson 7.5 Scope

Review the production-order slice for consistency gaps before closing Milestone 7.

#### Business Context

Production-order planning is now writable through several workflows. Managers and leaders can assign employees and move orders through lifecycle states, so concurrent edits must not silently overwrite each other.

#### Problem

`PostgresStore.Update` already used a transaction for root state and assignment replacement, but it did not detect stale updates. Two callers could load the same order version, make different valid changes and let the later write overwrite the earlier one.

#### Design Discussion

Added optimistic locking to the production-order aggregate root, matching the existing employee/product persistence pattern. The order root now carries a `version`, persisted updates atomically require the expected version and the database increments it on success.

The `Store.Update(ctx, order) error` shape was kept. Services refetch after a successful update so API responses show the persisted version without widening repository contracts for this lesson.

#### Go Concepts

- private aggregate version field with an accessor
- sentinel conflict errors with `errors.Is`
- preserving an existing interface while improving consistency
- `sync.RWMutex`-protected stale-write checks in the in-memory adapter

#### Architecture Concepts

- optimistic locking for aggregate-root updates
- stale-write detection as part of transactional consistency
- HTTP `409 Conflict` for concurrent modification
- milestone review before moving to reporting

### Lesson 7.5 Completion Notes

#### Business Context

Production-order updates are now protected against silent lost updates. A stale assignment or status transition returns a domain conflict instead of overwriting newer aggregate state.

#### Problem

Lesson 7.4 introduced transactional update, but transaction atomicity alone did not detect concurrent edits based on old state.

#### Design Discussion

Added a `version` column to `production_orders` and threaded it through `Order`, persistence reconstruction, HTTP responses and in-memory tests. PostgreSQL updates now use `WHERE id = $1 AND version = $4` and increment `version` in the same statement.

When no row is updated, the repository distinguishes missing orders from stale versions by checking whether the order still exists. This keeps `ErrNotFound` and `ErrVersionConflict` behavior precise.

#### Implementation

- Added migration `0007_add_production_order_versions.sql`.
- Added `orders.ErrVersionConflict`.
- Added `Order.version`, `Version()` and version validation.
- Updated `RestoreOrder`, sqlc queries and generated `ordersdb` code.
- Updated `PostgresStore.Update` to use optimistic locking.
- Updated `InMemoryStore.Update` to mirror version conflict behavior.
- Updated service mutation methods to refetch after successful updates.
- Added `version` to production-order HTTP responses.

#### Tests

- Added aggregate version assertion for new orders.
- Added service assertions that mutation responses return incremented versions.
- Added handler assertions for version in create/get/mutation responses.
- Added PostgreSQL integration coverage for persisted version, incremented version, stale update conflict and failed-update rollback preserving version.
- Verified with `make sqlc`.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `golangci-lint run ./...`.

#### Milestone Review

Milestone 7 is complete. The system now has a production-order aggregate, durable multi-table persistence, product/employee reference validation, HTTP create/read/mutation endpoints, role-based access and optimistic locking for mutable order workflows.

The milestone moved the project from CRUD-focused resources toward richer business modeling. The main architectural improvement is that order lifecycle rules live in the aggregate, workflow coordination lives in `orders.Service` and persistence consistency lives in the repository adapter.

#### Follow-Ups

- Production reference foreign keys remain `NOT VALID` as known technical debt.
- Full auth-user CRUD remains postponed until there is a concrete business workflow.
- Reporting should start with read-model/query design rather than more order mutations.

#### Exercises

- Explain why optimistic locking still matters when the update is already inside a transaction.
- Add a failing test for concurrent order release/cancel and make it pass with version conflicts.
- Compare returning the refetched order with changing `Store.Update` to return the updated aggregate.

#### Interview Questions

- What problem does optimistic locking solve?
- How do you distinguish not-found from version-conflict when both return no rows?
- Why should an aggregate root own the version rather than each child table?
- What does HTTP `409 Conflict` communicate to an API client?

#### Roadmap Update

- Lesson 7.5 completed.
- Milestone 7 completed.
- Current milestone moved to Milestone 8.
- Current lesson moved to Lesson 8.1.
- Concurrency `Mutex` and `RWMutex` marked complete in the Knowledge Matrix.

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

✅ Completed

### Lessons

- **L8.1** — Reporting Query Foundations ✅
- **L8.2** — Daily Production Report API ✅
- **L8.3** — Employee Productivity Report API ✅
- **L8.4** — Product Statistics Report API ✅
- **L8.5** — Reporting Review & Query Performance ✅
- **L8.6** — Detailed Production Breakdown Reports ✅

### Lesson 8.1 Scope

Introduce reporting as a read-model/query slice before adding HTTP endpoints.

#### Business Context

Managers need production totals without manually filtering Excel sheets. Reporting reads existing operational data and produces the summaries that were previously maintained manually.

#### Problem

The application stores production entries, employees, products and orders, but every existing slice is command/workflow oriented. Adding reports directly to the production slice would mix write-model behavior with read-specific SQL aggregation.

#### Design Discussion

Reporting starts as its own vertical slice with a PostgreSQL read store. It returns read models, not production-entry aggregates. This is a small CQRS step: commands still live in their business slices, while reporting owns query shapes optimized for management views.

The first query groups production entries by UTC day and product SKU. HTTP endpoints, response DTOs and authorization are postponed to Lesson 8.2 so this lesson can focus on SQL aggregate correctness and package boundaries.

#### Go Concepts

- read-model structs separate from domain entities
- context propagation into query methods
- time-range validation with sentinel errors
- integer conversion at SQL aggregate boundaries

#### Architecture Concepts

- CQRS as separate read models, not a new framework
- reporting vertical slice owns reporting SQL
- query store returns projections instead of aggregates
- ADR for introducing reporting read models

### Lesson 8.1 Completion Notes

#### Business Context

MES Lite now has a reporting foundation for manager-facing production summaries.

#### Problem

Operational data existed, but there was no dedicated query boundary for reports. Reusing production repositories would have forced report-shaped aggregate SQL into a write-focused slice.

#### Design Discussion

Added `internal/reporting` with a small read store. The store exposes `DailyProduction`, which returns `DailyProductionRow` values grouped by UTC day and product SKU.

This is CQRS in a pragmatic Go style: no bus, no framework and no generic query abstraction. The separation exists because reports have different data shapes and SQL needs than commands.

#### Implementation

- Added `reporting.Store` and `DailyProductionRow` read model.
- Added time-range validation with `ErrInvalidRange`.
- Added reporting sqlc configuration.
- Added `internal/reporting/queries/reports.sql`.
- Generated `internal/reporting/reportingdb`.
- Added `reporting.PostgresStore`.
- Added ADR `0003-introduce-reporting-read-models.md`.

#### Tests

- Added unit tests for report range validation.
- Added PostgreSQL integration test for daily production grouping by day and product.
- Added PostgreSQL integration test for invalid ranges.
- Verified with `sqlc generate`.
- Verified with `go test ./internal/reporting -count=1`.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

No existing write-model repository was changed. Reporting reads from existing tables through its own query package.

#### Code Review

An experienced Go engineer would approve the small CQRS introduction because it solves a real reporting shape without adding a framework or generic abstraction. The main follow-up is HTTP/API design in Lesson 8.2, including query parameters and authorization.

#### Exercises

- Add a test proving an empty date range returns an empty non-nil slice.
- Explain why `DailyProductionRow` is not a `production.Entry`.
- Add employee name to the query and discuss whether that belongs in L8.2 or L8.3.

#### Interview Questions

- What is CQRS, and what problem does it solve here?
- Why can read models differ from domain aggregates?
- Why group by UTC day instead of local server time?
- What are the risks of putting reporting SQL into command repositories?

#### Roadmap Update

- Lesson 8.1 completed.
- Current lesson moved to Lesson 8.2.
- Architecture `CQRS` marked complete in the Knowledge Matrix.

### Lesson 8.2 Scope

Expose the daily production report over HTTP with query parameters and role-based access.

#### Business Context

Managers and leaders need a report they can call from a future UI or external tool to see daily production totals without exporting raw entries and aggregating them manually.

#### Problem

Lesson 8.1 created the reporting query boundary, but there was no API route. Clients could not request a report for a specific time range.

#### Design Discussion

The API accepts explicit RFC3339 `from` and `to` query parameters and treats the range as half-open: `from <= occurred_at < to`. This avoids double-counting rows when clients page through adjacent date ranges.

The endpoint returns HTTP DTOs instead of exposing the reporting store row type directly. Reports are protected with the same JWT bearer authentication and RBAC middleware as other business endpoints. Admins, managers and leaders can read reports; workers cannot.

#### Go Concepts

- parsing query parameters with `time.Parse`
- RFC3339 timestamps for API boundaries
- half-open time ranges
- DTO mapping for read models

#### Architecture Concepts

- reporting HTTP boundary over a query store
- RBAC for management read endpoints
- server composition root wiring for reporting
- OpenAPI route registration through Fuego

### Lesson 8.2 Completion Notes

#### Business Context

MES Lite now exposes a daily production report endpoint for management users.

#### Problem

The reporting query existed only as an internal store method. There was no authenticated API route for clients to request daily production totals.

#### Design Discussion

Added `GET /reports/daily-production?from=<RFC3339>&to=<RFC3339>`. The handler parses transport-level query parameters, delegates aggregation to `reporting.Store` and maps read rows to JSON response DTOs.

The endpoint is read-only and management-oriented. Route-level RBAC allows admins, managers and leaders while rejecting workers.

#### Implementation

- Added `reporting.Handler`.
- Added `DailyProductionResponse` and `DailyProductionRowResponse` DTOs.
- Added RFC3339 query-parameter parsing for `from` and `to`.
- Added `reporting.InMemoryStore` for fast handler/server tests.
- Wired `reporting.PostgresStore` and `reporting.Handler` in `cmd/server`.
- Registered `GET /reports/daily-production` in `internal/server` with bearer security and RBAC.

#### Tests

- Added handler test for successful daily production response mapping.
- Added handler tests for missing, malformed and invalid time ranges.
- Added server route test proving managers can read the report.
- Added server route test proving workers cannot read the report.
- Verified with `go fmt ./cmd/server ./internal/server ./internal/reporting`.
- Verified with `go test ./internal/reporting ./internal/server -count=1`.
- Verified with `sqlc generate`.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

`server.New` now accepts a reporting handler explicitly. This keeps constructor dependency injection consistent with the rest of the application.

#### Code Review

An experienced Go engineer would approve the narrow HTTP boundary and explicit time-range parsing. The main follow-up is to add employee-level productivity reporting in Lesson 8.3 without overloading the daily production endpoint.

#### Exercises

- Add a test proving adjacent daily ranges do not double-count boundary entries.
- Add a server test proving leaders can read daily production reports.
- Inspect the generated OpenAPI document and find the report route.

#### Interview Questions

- Why are half-open time ranges common in APIs and SQL queries?
- Why use RFC3339 for timestamps at HTTP boundaries?
- Why should reporting DTOs be separate from database row structs?
- Why should workers be forbidden from management reports even though they can register production?

#### Roadmap Update

- Lesson 8.2 completed.
- Current lesson moved to Lesson 8.3.

### Lesson 8.3 Scope

Add an employee productivity report that ranks employees by completed production quantity for a requested time range.

#### Business Context

Production managers and team leaders need to replace manually prepared employee-output summaries where those summaries are already tracked. The report should answer who produced how much during a period.

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

### Lesson 8.4 Scope

Add product-level statistics reporting for completed production over a requested time range.

#### Business Context

Production managers need product-level production summaries without manually aggregating production entries from Excel.

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

### Lesson 8.5 Scope

Review the reporting slice and add the first query-performance guardrail before closing Milestone 8.

#### Business Context

Reports will become slower as production history grows. Managers need report endpoints to remain predictable without turning every request into a full-table scan over production history.

#### Problem

The reporting queries were correct and tested, but `production_entries` had no reporting-oriented index. All three reports filter by `occurred_at` before grouping by product or employee.

#### Design Discussion

Added one targeted covering index on `production_entries` for the reporting access pattern: `(occurred_at, product_sku, employee_id) INCLUDE (quantity)`. This supports the shared time-range filter and keeps product, employee and quantity data available from the index for aggregation-heavy reports.

This is intentionally one index, not several speculative indexes. Indexes speed reads but slow writes and consume storage. More indexes should be added only after real query plans or production-like data show a need.

#### Go Concepts

- integration tests that verify migration side effects
- keeping performance behavior explicit through schema migrations
- reviewing abstractions before adding more code

#### SQL Concepts

- B-tree indexes for range predicates
- covering indexes with `INCLUDE`
- trade-offs between read speed, write cost and storage
- query-performance review before premature optimization

### Lesson 8.5 Completion Notes

#### Business Context

Milestone 8 now has useful management reports and a first database-level performance guardrail for report time-range queries.

#### Problem

Reports were implemented, but there was no index supporting their common access pattern. As production entries grow, report queries would increasingly depend on scanning the whole table.

#### Design Discussion

The three report queries share the same shape: filter production entries by `occurred_at`, then aggregate by product or employee. A single covering index is the smallest useful improvement.

The lesson did not introduce caching, materialized views or background report generation. Those are valid future tools, but they would be premature before measuring real query cost on larger data.

#### Implementation

- Added migration `0008_add_reporting_indexes.sql`.
- Added `production_entries_reporting_idx` on `(occurred_at, product_sku, employee_id) INCLUDE (quantity)`.
- Added integration test proving the reporting index exists after migrations.

#### Tests

- Verified with `go fmt ./internal/reporting`.
- Verified with `sqlc generate`.
- Verified with `go test ./internal/reporting -count=1`.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

No report abstraction was added. Three explicit report methods remain clearer than a generic reporting engine.

#### Code Review

An experienced Go engineer would approve the reporting slice for this milestone: read models are explicit, SQL is readable, generated sqlc code stays below the package boundary, routes are protected and tests cover query correctness and authorization.

The main caveat is that performance was improved with a reasoned index, not proven with production-scale `EXPLAIN ANALYZE`. That is acceptable at this project stage, but future performance work should use realistic row counts.

#### Exercises

- Run `EXPLAIN ANALYZE` for each report query before and after adding sample data.
- Compare one covering index with separate `(occurred_at, employee_id)` and `(occurred_at, product_sku)` indexes.
- Explain how additional indexes affect production-entry insert performance.

#### Interview Questions

- Why does every index have a write-time cost?
- What is a covering index?
- Why should performance optimization be measurement-driven?
- When would a materialized view be better than querying base tables directly?

#### Roadmap Update

- Lesson 8.5 completed.
- Milestone 8 completed.
- Current milestone moved to Milestone 9.
- Current lesson moved to Lesson 9.1.
- Persistence `SQL Optimization` and SQL `Indexes` marked complete in the Knowledge Matrix.

### Milestone 8 Review

#### Architecture Review

An experienced Go engineer would approve the milestone direction. Reporting is isolated in `internal/reporting`, uses CQRS-style read models and does not pollute command-oriented slices with aggregate queries.

The reporting package intentionally reads across production, employees and products. This is acceptable because reporting is a read-side projection layer. Business writes and invariants remain owned by their original slices.

#### Code Review

The code remains explicit and idiomatic. There is no generic report framework, no dynamic SQL and no leaking sqlc row types into HTTP handlers. Each endpoint has a clear DTO and route-level RBAC.

The main improvement for later is OpenAPI query-parameter documentation quality. The endpoints work and are generated, but explicit query metadata remains known technical debt.

#### Refactoring

The shared `reportRange` helper is justified because all report endpoints use the same `from`/`to` RFC3339 half-open range. No broader abstraction is needed.

#### Interview Review

You should now be able to explain CQRS read models, why reports can join across slices, why half-open time ranges avoid boundary bugs, how `GROUP BY` aggregation works, why deterministic ordering matters and what trade-offs indexes introduce.

#### Completion Criteria

- Daily production report implemented.
- Employee productivity report implemented.
- Product statistics report implemented.
- Reporting endpoints require bearer authentication and management RBAC.
- Reporting SQL is generated through sqlc.
- Reporting integration tests verify query correctness.
- Reporting performance has a first targeted index.
- Tests, build, lint and sqlc generation pass.
- Roadmap updated.

### Lesson 8.6 Scope

Add detailed breakdown reports that preserve the existing summary endpoints while answering product-by-employee questions.

#### Business Context

Managers need both summary and detail when those details are already tracked manually. A daily product total answers what was produced each day, but not who produced each product. Employee productivity answers total output per employee, but not which products each employee produced.

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

MES Lite now has detailed product-by-employee reporting for manual-summary replacement without changing the existing summary reports.

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

### Goal

Provide reports that replace existing Excel or paper summaries.

### Business Value

Managers can review production summaries without maintaining Excel reports manually.

### Features

- daily production
- employee productivity where manually tracked
- product statistics where manually tracked
- filtering

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

- manual-summary reports complete
- SQL optimized
- tests passing

---

## Milestone 9 - CSV Import

Status

🚧 In Progress

### Lessons

- **L9.1** — CSV Import Design & Streaming Reader ✅
- **L9.2** — Row Validation & Error Collection ✅
- **L9.3** — Batch Persistence & Transaction Strategy
- **L9.4** — Import Summary API & Partial Failure Reporting
- **L9.5** — CSV Import Review & Performance

### Lesson 9.1 Scope

Introduce the CSV import slice and implement a streaming reader for historical production-entry CSV files.

#### Business Context

The company has historical production data in Excel-like files. Before the system can replace manual tracking completely, it needs a safe path to read that history without loading entire files into memory.

#### Problem

CSV import can easily become memory-heavy if the implementation reads the whole file before processing. Large historical exports should be processed row by row through `io.Reader` and `encoding/csv.Reader`.

#### Design Discussion

Lesson 9.1 only defines the import boundary and raw CSV row reader. It does not parse quantities or timestamps into domain values yet, and it does not persist rows. Those responsibilities belong to later lessons so the streaming concept remains clear.

The reader validates the CSV header once, then exposes one raw production-entry row per `Read` call. This mirrors Go's standard streaming style: callers pull data until `io.EOF`.

#### Go Concepts

- `io.Reader` as the standard input abstraction
- `encoding/csv.Reader` for CSV parsing
- sentinel errors for stable import failures
- `io.EOF` as normal stream completion
- row numbers for future validation diagnostics

#### Architecture Concepts

- CSV import as its own vertical slice in `internal/csvimport`
- raw transport/import rows separated from production domain entries
- validation and persistence intentionally postponed to later pipeline steps

### Lesson 9.1 Completion Notes

#### Business Context

MES Lite now has the first building block for migrating historical production data from CSV exports.

#### Problem

The application had no import boundary. Starting with upload handlers or database writes would have mixed HTTP, parsing, validation and persistence before establishing a memory-safe reader.

#### Design Discussion

Added `internal/csvimport` with a concrete `ProductionEntryReader`. It accepts any `io.Reader`, validates the expected production-entry header and streams rows one at a time.

The expected CSV columns are `employee_id`, `product_sku`, `quantity`, `workstation`, `timestamp` and `comment`. Values remain strings in this lesson because row validation, type parsing and error collection are the focus of Lesson 9.2.

#### Implementation

- Added `ProductionEntryRow` raw import row type.
- Added `ProductionEntryReader` over `encoding/csv.Reader`.
- Added `NewProductionEntryReader` with header validation.
- Added `Read` returning one row at a time and `io.EOF` when complete.
- Added import errors `ErrInvalidHeader` and `ErrInvalidRecord`.
- Added row-number tracking for future validation diagnostics.

#### Tests

- Tested sequential row reading and trimming.
- Tested header normalization for case and whitespace.
- Tested missing and unexpected headers.
- Tested malformed CSV parse errors.
- Tested `io.EOF` stream completion.
- Tested constructor behavior only consumes the header before row reads.
- Verified with `go test ./internal/csvimport -count=1`.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

No existing production-registration code was changed. The import slice is intentionally separate because CSV rows are external input records, not validated production entries yet.

#### Code Review

An experienced Go engineer would approve the narrow scope: the code uses standard-library streaming, keeps APIs explicit and avoids premature importer abstractions. The main follow-up is to add typed row validation and collect row-level errors without stopping the whole import unnecessarily.

#### Exercises

- Add a test with 100,000 generated rows and assert rows are processed through repeated `Read` calls.
- Explain why `io.EOF` is not treated as an error condition by stream consumers.
- Extend the reader test with a CSV row containing a comma inside a quoted comment.

#### Interview Questions

- Why is `io.Reader` one of the most important interfaces in Go?
- Why should large CSV imports be streamed instead of read into memory?
- What does `encoding/csv.Reader` handle that manual `strings.Split` would get wrong?
- Why should parsing and business validation be separate from raw CSV reading?

#### Roadmap Update

- Lesson 9.1 completed.
- Current lesson moved to Lesson 9.2.
- Standard Library `io` and `encoding/csv` marked complete in the Knowledge Matrix.

### Lesson 9.2 Scope

Parse raw CSV rows into typed import records and collect row-level validation errors without stopping the whole import.

#### Business Context

Historical production CSV files may contain mistakes from manual spreadsheets. Managers need an import process that identifies bad rows clearly instead of failing with a vague first error.

#### Problem

Lesson 9.1 could stream raw CSV fields, but every field was still a string. The importer could not yet distinguish valid rows from rows with missing IDs, invalid quantities or malformed timestamps.

#### Design Discussion

Validation is a pipeline step after raw CSV reading and before persistence. Structural CSV errors remain fatal because the stream cannot be trusted. Row-level business-shape errors are collected so the importer can later report partial failures.

The validated type is still import-specific. It is not a `production.Entry` because production-entry IDs, reference validation and persistence behavior belong to later lessons.

#### Go Concepts

- `strconv.Atoi` for integer parsing
- `time.Parse` with `time.RFC3339`
- `errors.Is` for stream completion checks
- small error structs implementing `error`
- collecting validation errors in slices

#### Architecture Concepts

- validation pipeline step between reader and persistence
- import records separated from production domain entities
- fatal stream errors separated from row-level validation errors

### Lesson 9.2 Completion Notes

#### Business Context

MES Lite can now identify valid historical production rows and explain row-level CSV mistakes in a way that a future API can return to managers.

#### Problem

Raw CSV reading alone was not enough for import. Quantity and timestamp fields needed type parsing, and invalid rows needed structured diagnostics without blocking valid rows from being recognized.

#### Design Discussion

Added `ValidateProductionEntries`, which consumes a `ProductionEntryReader`, validates each raw row and returns a `ValidationResult` containing valid typed records plus row errors.

Malformed CSV structure is still returned as a fatal error. Missing required fields, invalid quantities and invalid timestamps are collected as `RowError` values.

#### Implementation

- Added `ProductionEntryRecord` typed import record.
- Added `RowError` with row number, field and message.
- Added `ValidationResult` containing valid records and row errors.
- Added `ValidateProductionEntries` to stream through rows and collect validation errors.
- Added quantity parsing with positive and PostgreSQL `integer` bounds checks.
- Added RFC3339 timestamp parsing and UTC normalization.
- Mapped CSV field-count failures to `ErrInvalidRecord`.

#### Tests

- Tested valid CSV rows become typed records.
- Tested timestamp offsets normalize to UTC.
- Tested mixed valid and invalid rows keep valid records and collect all row errors.
- Tested missing, non-integer, negative and too-large quantities.
- Tested fatal malformed CSV record handling.
- Tested `RowError.Error` formatting.
- Verified with `go test ./internal/csvimport -count=1`.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

The reader now maps CSV field-count parse failures to the import sentinel `ErrInvalidRecord`, making stream-shape failures easier for future API code to translate.

#### Code Review

An experienced Go engineer would approve the separation between raw reading and typed row validation. The code remains concrete, standard-library based and avoids a generic validation framework.

The main follow-up is persistence strategy: decide how validated rows should be saved in batches, where transactions start and what partial failure means for valid versus invalid rows.

#### Exercises

- Add a test proving quoted comments containing commas survive validation.
- Add a row with several invalid fields and explain why all errors are collected.
- Decide whether future import summaries should include the original raw row values.

#### Interview Questions

- Why separate fatal stream errors from row-level validation errors?
- Why should CSV import validation not directly create persisted production entries yet?
- Why is RFC3339 a good timestamp format for import boundaries?
- What trade-off appears when collecting all errors instead of stopping at the first invalid row?

#### Roadmap Update

- Lesson 9.2 completed.
- Current lesson moved to Lesson 9.3.
- L9.3 remains focused on batch persistence and transaction strategy.

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

At this point the application should have most backend/API foundations for a small-company Excel/paper replacement. Milestone 10 closes the remaining MVP API gaps before post-MVP learning milestones begin.

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

If significant gaps are found, they should be addressed before moving to post-MVP learning milestones 11-15.

---

## Milestone 10 - MVP API Completion

Status

⬜ Not Started

### Goal

Close the remaining API gaps for the Excel/paper-replacement MVP without expanding into full MES scope.

### Business Value

Managers and leaders can review entered production records, correct mistakes with auditability and trust production registration even when clients retry requests.

### Lessons

- **L10.1** — Production Entry Review API
- **L10.2** — Production Registration Idempotency
- **L10.3** — Production Entry Corrections & Audit Trail
- **L10.4** — MVP API Review

### Features

- list/review production entries for managers and leaders
- filters matching current manual review needs
- explicit workstation as a production-entry text field
- idempotency key or request ID for production registration
- append-only correction records for production-entry mistakes
- correction reason and actor tracking

### Non-Goals

- formal workstation management
- dashboards
- offline-first client behavior
- machine integration
- generic audit framework
- broad production-event refactor beyond the correction/idempotency needs of the current production-entry model

### Go Concepts

- idempotent command handling
- unique constraints as application guardrails
- DTOs for review and correction workflows
- domain errors for duplicate request IDs and invalid corrections

### Architecture Concepts

- append-only corrections instead of silent edits
- auditability for production history
- command idempotency at HTTP/API boundaries
- minimal MVP completion before post-MVP learning milestones

### Testing

- production-entry list integration tests
- authorization tests for manager/leader review
- duplicate request-ID tests
- correction audit tests
- worker-forbidden correction tests

### Definition of Done

- managers/leaders can list and review entered production records
- production registration accepts and enforces an idempotency key or request ID
- permitted users can correct production-entry mistakes without overwriting original records
- normal workers cannot correct historical records
- workstation remains a simple production-entry text field for MVP
- tests passing

---

## Milestone 11 - Background Jobs & Concurrency

Status

⬜ Not Started

### Goal

Learn idiomatic concurrency by introducing asynchronous processing.

Concurrency is one of Go's defining features and should be understood deeply.

### Business Value

Long-running tasks should not block HTTP requests.

Post-MVP examples:

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

## Milestone 12 - Machine Integration & Synchronization

Status

⬜ Not Started

### Goal

Safely process concurrent events coming from production machines.

### Business Value

Prepare the system for future CNC and production machine integration after the Excel/paper replacement is proven useful.

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

## Milestone 13 - Observability

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

## Milestone 14 - Performance Engineering

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

## Milestone 15 - Production Readiness

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
