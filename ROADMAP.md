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

**Version:** 1.7
**Status:** IN PROGRESS
**Current milestone:** 5 - Persistence & Data Access
**Current lesson:** Not started
**Completed milestones:**
- Milestone 0
- Milestone 1
- Milestone 2
- Milestone 3
- Milestone 4
**Next milestone:** 6 - Authentication & Authorization
**Current branch:** main
**Architecture maturity:** 5 / 10
**Go knowledge progress:** 30%
**Interview readiness:** 20%
**Known technical debt:** Employees and products remain in-memory until Milestone 5; production reference validation is not fully transactionally consistent yet.

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
- [ ] Aggregates
- [ ] CQRS
- [ ] Event-driven Architecture
- [ ] Observability
- [ ] Production Readiness

**Persistence**
- [x] PostgreSQL
- [x] pgx
- [x] sqlc
- [ ] Transactions
- [ ] Optimistic Locking
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
- [ ] Middleware
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

⬜ Not Started

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

## Milestone 6 - Authentication & Authorization

Status

⬜ Not Started

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

⬜ Not Started

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
