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

**Version:** 2.53
**Status:** IN PROGRESS
**Current milestone:** 14 - Performance Engineering
**Current lesson:** L14.1 - Benchmarking Foundations
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
- Milestone 9
- Milestone 10
- Milestone 11
- Milestone 12
- Milestone 13
**Next milestone:** 15 - Production Readiness
**Current branch:** main
**Architecture maturity:** 9.0 / 10
**Go knowledge progress:** 82%
**Interview readiness:** 78%
**Known technical debt:** Production reference foreign keys are `NOT VALID`, so PostgreSQL enforces new production entries but does not validate legacy rows created before the constraint. Query parameters are implemented for list endpoints; explicit OpenAPI query-parameter documentation should be reviewed later. Auth-user management is intentionally limited to durable bootstrap admin creation; full user-management CRUD is postponed until there is a concrete business workflow. CSV import uses bounded batch inserts with regular `INSERT` statements instead of PostgreSQL `COPY`; revisit if production-scale import performance requires it. Background jobs are in-memory only, so queued, running and completed job history is lost on restart; ADR 0004 records the durable-queue trigger. Async CSV import jobs use temporary upload files, so process crashes can leave orphaned temp files and lose queued imports until durable job storage exists. Background job status and cancellation APIs exist for individual jobs, but there is no job list endpoint yet. Machine events are currently accepted through a fake JWT-protected API and stored in memory only; deduplication and intake counters are in-memory only, events are not durable, events are not processed into production entries and machines are not authenticated with real machine credentials yet. HTTP tracing exists with no-op/stdout exporters only; OTLP collector export is not wired yet. The generated sqlc boundary is documented but not enforced: `internal/csvimport` imports `internal/production/productiondb` directly, and nested `internal/` directories would make that a compile error once CSV import stops reusing another slice's generated insert. The depguard platform rule uses a deny list naming each business slice, so adding a slice requires updating `.golangci.yml`.

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
- [x] sync
- [x] sync/atomic
- [x] database/sql (concepts)
- [x] log/slog
- [ ] runtime
- [ ] embed

**Concurrency**
- [x] Goroutines
- [x] Channels
- [x] select
- [x] WaitGroup
- [ ] errgroup
- [x] Mutex
- [x] RWMutex
- [x] Atomic
- [x] Worker Pools
- [x] Pipelines
- [x] Race Detector

**Architecture**
- [x] Vertical Slice
- [x] Dependency Injection
- [x] Package Design
- [x] Repositories
- [x] Aggregates
- [x] CQRS
- [ ] Event-driven Architecture
- [x] Observability
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
- [x] Race Detection

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

✅ Completed

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

✅ Completed

### Lessons

- **L9.1** — CSV Import Design & Streaming Reader ✅
- **L9.2** — Row Validation & Error Collection ✅
- **L9.3** — Batch Persistence & Transaction Strategy ✅
- **L9.4** — Import Summary API & Partial Failure Reporting ✅
- **L9.5** — CSV Import Review & Performance ✅

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

### Lesson 9.3 Scope

Persist validated CSV import records in batches and define the first transaction strategy for historical production-entry imports.

#### Business Context

Historical imports should not leave the database half-updated when one row in a valid-looking batch fails at the persistence boundary. Managers need imported history to be trustworthy.

#### Problem

Lesson 9.2 produced typed valid import records, but there was no durable save path. Saving rows one by one outside a transaction could persist earlier rows and then fail on a later foreign-key or constraint error.

#### Design Discussion

The batch persistence boundary is all-or-nothing for records passed to `SaveBatch`. If a row fails during insertion, the transaction rolls back every row in that batch.

Partial failure reporting remains a higher-level API concern for Lesson 9.4. L9.3 only answers the storage question: a batch either commits completely or leaves no production-entry rows behind.

#### Go Concepts

- explicit transaction lifecycle with `Begin`, deferred `Rollback` and `Commit`
- wrapping row-specific persistence errors with `Unwrap`
- converting validated import rows into production domain entries
- returning empty non-nil slices for empty batch results

#### Architecture Concepts

- import store as a persistence boundary for CSV import workflow
- transaction boundary around one import batch
- PostgreSQL constraints as final integrity guardrails
- no generic transaction manager before repeated need exists

### Lesson 9.3 Completion Notes

#### Business Context

MES Lite can now persist validated historical production rows in one transactional batch.

#### Problem

Validated rows were still in memory only. A future upload endpoint would have had no safe way to save rows without risking partial database writes.

#### Design Discussion

Added a CSV import store that writes validated production-entry import records to PostgreSQL. `SaveBatch` generates production-entry IDs, constructs domain `production.Entry` values and inserts them through the existing sqlc production-entry query inside one transaction.

The transaction is intentionally owned by the CSV import store because the import workflow is the business operation. A broad transaction manager was not introduced.

#### Implementation

- Added `csvimport.Store` interface with `SaveBatch`.
- Added `csvimport.PostgresStore` backed by `pgxpool` transactions.
- Added `BatchError` that records the CSV row number that failed persistence.
- Converted `ProductionEntryRecord` values to `production.Entry` values before persistence.
- Inserted batch records through existing `productiondb.CreateEntry` sqlc query.
- Mapped PostgreSQL constraint failures to production domain errors.
- Returned empty non-nil slices for empty batches.

#### Tests

- Tested empty batch behavior without needing a database.
- Added PostgreSQL integration test for successful two-row batch persistence.
- Verified generated entry IDs are returned and persisted rows can be read by the production store.
- Added rollback integration test proving a missing reference on a later row leaves zero rows committed.
- Verified row number is preserved in `BatchError`.
- Verified with `go test ./internal/csvimport -count=1`.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

No new sqlc query was needed. The import store reuses the existing production-entry insert query instead of duplicating SQL.

#### Code Review

An experienced Go engineer would approve the transaction boundary because it protects exactly one batch import operation. The main caveat is expected lesson scope: API-level import summaries and partial failure reporting are not exposed yet.

#### Exercises

- Add a test proving duplicate generated IDs would roll back the whole batch by injecting deterministic IDs.
- Explain why deferred rollback is still called after a successful commit and why that is safe.
- Decide whether future imports should use one transaction for the whole file or smaller batch transactions for very large files.

#### Interview Questions

- What should define a transaction boundary?
- Why is all-or-nothing persistence useful for a validated import batch?
- How does `errors.As` find a wrapped `BatchError`?
- When would smaller transactions be better than one large import transaction?

#### Roadmap Update

- Lesson 9.3 completed.
- Current lesson moved to Lesson 9.4.
- L9.4 remains focused on upload/API summary behavior and partial failure reporting.

### Lesson 9.4 Scope

Expose the CSV import workflow through an authenticated API endpoint and return an import summary with partial failure details.

#### Business Context

Managers need to upload historical production CSV files and understand what happened: how many rows were accepted, how many were rejected and which rows need correction.

#### Problem

The project could stream, validate and transactionally save batches internally, but there was no API entry point and no response contract for partial failures.

#### Design Discussion

The API accepts a raw `text/csv` request body at `POST /imports/production-entries`. The handler reads from the request body, while `csvimport.Service` coordinates the workflow.

Row-level validation failures are returned in the summary and do not prevent valid rows from being persisted. Fatal CSV shape errors such as invalid headers or malformed records return `400 Bad Request`. Persistence batch failures are reported with the row number that caused the batch rollback.

#### Go Concepts

- reading raw request bodies through `io.Reader`
- service-level orchestration over reader, validator and store
- response DTOs with JSON tags
- `errors.As` for extracting wrapped `BatchError`
- in-memory adapter for fast API tests

#### Architecture Concepts

- import service as workflow coordinator
- handler owns HTTP translation only
- API summary contract for partial failures
- RBAC-protected management import endpoint

### Lesson 9.4 Completion Notes

#### Business Context

MES Lite now exposes a manager/admin CSV import endpoint for historical production entries.

#### Problem

Internal import components existed, but clients could not upload CSV data or receive a useful summary of imported versus rejected rows.

#### Design Discussion

Added `csvimport.Service` to coordinate raw CSV reading, row validation and batch persistence. The handler stays small: it passes the request body to the service and maps fatal CSV input errors to `400 Bad Request`.

The endpoint returns `ImportSummary` with total rows, valid rows, invalid rows, imported rows and row-level errors. This gives API clients enough information to correct spreadsheet mistakes without silently ignoring invalid data.

#### Implementation

- Added `ImportSummary` and `ImportError` response types.
- Added `csvimport.Service.ImportProductionEntries`.
- Added `csvimport.Handler.ImportProductionEntries`.
- Added `csvimport.InMemoryStore` for fast service/server tests.
- Wired `csvimport.PostgresStore` and handler in `cmd/server`.
- Registered `POST /imports/production-entries` with bearer auth and admin/manager RBAC.

#### Tests

- Added service test for valid rows plus validation errors.
- Added service test for fatal invalid CSV headers.
- Added service test for persistence batch errors reported as row errors.
- Added service test for unexpected store errors.
- Added server route test proving managers can import CSV.
- Added server route tests proving workers are forbidden and unauthenticated callers get `401`.
- Added server route test proving invalid CSV headers return `400`.
- Verified with `go test ./internal/csvimport -count=1`.
- Verified with `go test ./internal/server -count=1`.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

`server.New` now accepts an optional CSV import handler for tests while the production composition root wires the real PostgreSQL-backed handler explicitly.

#### Code Review

An experienced Go engineer would approve the boundary split: the handler does not validate CSV rows, the service does not know about Fuego and persistence remains behind the import store.

The main follow-up is performance review in L9.5. The current service validates into slices before saving; that is acceptable for this lesson but should be reviewed against large-file import goals.

#### Exercises

- Add a server test proving admins can import CSV.
- Add a handler/service test for a CSV containing only invalid rows.
- Decide whether import summaries should include generated production-entry IDs for successful rows.

#### Interview Questions

- Why should the handler pass an `io.Reader` instead of reading the whole body into `[]byte`?
- Why are row-level errors returned in a `200 OK` summary while malformed CSV headers return `400 Bad Request`?
- What belongs in the import service versus the import handler?
- Why should import endpoints be restricted to managers/admins rather than workers?

#### Roadmap Update

- Lesson 9.4 completed.
- Current lesson moved to Lesson 9.5.
- L9.5 remains focused on CSV import review and performance.

### Lesson 9.5 Scope

Review the CSV import pipeline for large-file behavior before closing Milestone 9.

#### Business Context

Historical production files may be large. The import endpoint should remain predictable and should not require loading all valid rows or all validation errors into memory at once.

#### Problem

The L9.4 service streamed the HTTP body into the CSV reader, but then collected every valid row before saving. That defeated part of the milestone's streaming goal for large valid files.

#### Design Discussion

The import service now processes rows incrementally and saves valid records in bounded batches. This keeps the main valid-row memory cost proportional to the batch size rather than the file size.

Reported row errors are also capped. The summary still counts every invalid row, but only stores the first bounded set of errors. This is a production trade-off: exact counts remain available, while response size and memory usage are protected for very bad files.

If a batch fails with a row-level persistence error such as a duplicate production entry or invalid reference, the service retries that batch row by row. This keeps retry imports useful: already-imported rows can be reported as row errors while new valid rows in the same file still get imported.

#### Go Concepts

- bounded batching with slices
- exact counters separated from retained detail records
- streaming loop with `io.EOF` as normal completion
- memory safety trade-offs in API response design

#### Architecture Concepts

- import pipeline review before milestone closure
- batch size as an implementation detail
- bounded error reporting for production safety
- explicit technical debt for future PostgreSQL `COPY` optimization

### Lesson 9.5 Completion Notes

#### Business Context

MES Lite's CSV import path is now suitable for MVP historical data migration with bounded memory behavior for valid rows and bounded response/error storage for invalid rows.

#### Problem

The service accepted an `io.Reader`, but validated the entire CSV into slices before saving. Large valid imports would grow memory with file size.

#### Design Discussion

Refactored `ImportProductionEntries` into a streaming pipeline. It reads one row, validates one row and appends only valid records to a fixed-size batch. When the batch reaches the configured size, the service persists it and reuses the same slice.

The API reports at most `maxReportedErrors` row errors while continuing to count all invalid rows. The response includes `errorLimitReached` so clients know whether detailed errors were truncated.

Failed batches are isolated by retrying records individually when the failure is row-level. Unexpected infrastructure failures still terminate the import because they cannot be safely attributed to one CSV row.

#### Implementation

- Added `defaultImportBatchSize` for bounded valid-row batches.
- Added `maxReportedErrors` for bounded error detail storage.
- Refactored `ImportProductionEntries` to save batches during streaming instead of collecting all valid rows first.
- Added `errorLimitReached` to `ImportSummary`.
- Added row-by-row isolation after row-level batch persistence failures.
- Preserved exact `totalRows`, `validRows`, `invalidRows` and `importedRows` counters.
- Kept `ValidateProductionEntries` for focused validation tests, while the service uses the lower-level row validator for streaming.

#### Tests

- Added large-input test proving 1,205 valid rows are saved as batches of 500, 500 and 205.
- Added large-invalid-input test proving invalid row counts remain exact while reported errors are capped.
- Added retry/import-continuation tests proving duplicate/already-existing rows become summary errors while other valid rows still import.
- Verified no save calls happen when all rows are invalid.
- Verified with `go test ./internal/csvimport -count=1`.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

The service now owns the streaming pipeline directly instead of delegating to `ValidateProductionEntries`, because the service must interleave validation and batch persistence for large files.

#### Code Review

An experienced Go engineer would approve the milestone for MVP scale: the code is standard-library based, validates rows explicitly, persists valid rows in bounded batches, isolates row-level persistence failures and returns clear summaries.

The main performance follow-up is database insert strategy. Current batch persistence still uses one regular `INSERT` per row inside a transaction. PostgreSQL `COPY` or `pgx.CopyFrom` may be justified later for very large imports, but adding it now would complicate the lesson before real performance data exists.

#### Exercises

- Make `defaultImportBatchSize` configurable and discuss whether it belongs in environment configuration.
- Add a benchmark comparing all-at-once validation with bounded-batch streaming.
- Sketch how a future background job would stream the same import without blocking an HTTP request.

#### Interview Questions

- Why does accepting `io.Reader` not automatically guarantee low memory usage?
- Why can response error details need a cap even if row processing is streamed?
- What trade-offs exist between regular batched `INSERT` and PostgreSQL `COPY`?
- Why might CSV import become a background job in a production system?

#### Roadmap Update

- Lesson 9.5 completed.
- Milestone 9 completed.
- Current milestone moved to Milestone 10.
- Current lesson moved to Lesson 10.1.
- Concurrency `Pipelines` marked complete in the Knowledge Matrix.

### Milestone 9 Review

#### Architecture Review

An experienced Go engineer would approve the CSV import slice for MVP. The package owns raw CSV reading, typed row validation, import orchestration, batch persistence and HTTP upload handling without leaking those concerns into the production registration slice.

The design is intentionally pragmatic: CSV import converts to `production.Entry` before persistence, reuses the existing production insert query and keeps PostgreSQL constraints as the final reference-integrity guardrail.

#### Code Review

The implementation is explicit and idiomatic. It uses `io.Reader`, `encoding/csv`, sentinel errors, `errors.Is`, `errors.As`, bounded slices and transaction-backed persistence.

The main improvement for later is performance at database scale. If imports become very large, `pgx.CopyFrom`, background jobs and progress tracking should be considered.

#### Refactoring

The main L9.5 refactor changed service orchestration from collect-all validation to streaming validation plus bounded batch saves. This better matches the milestone's streaming goal.

#### Interview Review

You should now be able to explain why `io.Reader` is fundamental, why `encoding/csv.Reader` is safer than `strings.Split`, why `io.EOF` is normal stream completion, how to collect row errors safely, why transactions protect import batches and why batching is not the same as loading a full file into memory.

#### Completion Criteria

- CSV import endpoint implemented.
- CSV rows are streamed from request body.
- Row validation and structured error collection implemented.
- Valid rows are persisted in bounded transactional batches.
- Failed batches are isolated so retry imports can continue past duplicate/already-existing rows.
- Import summary reports totals and partial failures.
- Management RBAC protects import endpoint.
- Large valid input is processed in bounded batches.
- Reported error details are capped while invalid-row counts remain exact.
- Tests, build and lint pass.
- Roadmap updated.

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

✅ Completed

### Goal

Close the remaining API gaps for the Excel/paper-replacement MVP without expanding into full MES scope.

### Business Value

Managers and leaders can review entered production records, correct mistakes with auditability and trust production registration even when clients retry requests.

### Lessons

- **L10.1** — Production Entry Review API
- **L10.2** — Production Registration Idempotency
- **L10.3** — Production Entry Corrections & Audit Trail
- **L10.4** — MVP API Review
- **L10.5** — Employee/Product Response DTO Cleanup

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

### Lesson 10.2 Scope

Make production registration safe for client retries by adding an explicit request ID to the registration command.

#### Business Context

Workers may submit production from unreliable networks or simple clients. If a client times out after the server already saved an entry, retrying the same request should not create duplicate production history.

#### Problem

Production registration generated a new entry ID for every request. A retry with the same business data created another production entry, inflating quantities and making historical review unreliable.

#### Design Discussion

The API now requires `requestId` in the JSON body for `POST /production-entries`. The database stores it in `production_entries.request_id` and enforces uniqueness through a partial unique index.

The partial index applies only when `request_id <> ''`. This keeps historical CSV imports and legacy rows valid without inventing fake request IDs for data that did not originate from a retryable API command.

Retry behavior is explicit: the same `requestId` with identical production data returns the original entry. The same `requestId` with different production data returns `409 Conflict`.

#### Go Concepts

- idempotent command handling
- comparing domain values while ignoring generated IDs
- sentinel conflict errors with `errors.Is`
- unique database constraints translated into domain errors

#### Architecture Concepts

- idempotency at the command/API boundary
- database uniqueness as a concurrency guardrail
- application service resolving safe retries
- import workflow kept separate from API registration workflow

### Lesson 10.2 Completion Notes

#### Business Context

Production registration is now retry-safe for API clients that provide a stable `requestId`.

#### Problem

Repeated HTTP registrations could create duplicate production entries if a client retried after a timeout or uncertain network failure.

#### Design Discussion

Added request-ID idempotency to the registration command instead of exposing client-generated production-entry IDs. The server still owns entry identity, while the client owns retry identity.

`production.Entry` keeps `RequestID` optional so CSV import and historical rows can persist entries without request IDs. `Service.Register` requires a request ID because it represents the retryable API command path.

#### Implementation

- Added migration `0009_add_production_entry_request_ids.sql`.
- Added `request_id` to `production_entries` with a partial unique index.
- Added `RequestID` to `production.Entry`.
- Added `NewEntryWithRequestID` for API registration while preserving `NewEntry` for non-idempotent historical/import paths.
- Added `ErrRequestConflict`.
- Required `requestId` in `RegisterProductionRequest` and `Service.Register`.
- Added `Store.FindByRequestID`.
- Mapped duplicate PostgreSQL request IDs to `ErrRequestConflict`.
- Made duplicate identical retries return the original production entry.
- Made duplicate different retries return `409 Conflict`.
- Updated CSV import persistence to insert blank request IDs for historical rows.

#### Tests

- Added service test for idempotent retry returning the existing entry.
- Added service test for request-ID conflict with different production data.
- Added service test proving request ID is required for registration commands.
- Added handler tests for idempotent retry and conflict responses.
- Added PostgreSQL store test for duplicate request-ID constraint mapping and lookup.
- Added server route test for end-to-end idempotent retry behavior.
- Verified with `go test ./internal/production ./internal/server ./internal/csvimport -count=1`.

#### Refactoring

The domain constructor was extended with `NewEntryWithRequestID` instead of replacing every existing call site. This keeps CSV import and legacy test setup simple while making API registration stricter at the service boundary.

#### Code Review

An experienced Go engineer would approve the partial unique index and retry semantics. The main trade-off is putting `requestId` in the JSON body instead of an `Idempotency-Key` header. The body field is acceptable for this lesson because it keeps the command DTO explicit and OpenAPI-friendly.

#### Exercises

- Add a test that fires two concurrent registrations with the same `requestId` and proves only one row is persisted.
- Change the API to use an `Idempotency-Key` header and compare the OpenAPI and handler trade-offs.
- Explain why the database unique index is still required even though the service checks duplicates.

#### Interview Questions

- What does idempotency mean for a write endpoint?
- Why should retries return the original result instead of blindly returning `409`?
- Why is a unique database constraint safer than an application-only duplicate check?
- Why separate server-generated entry IDs from client-provided request IDs?

#### Roadmap Update

- Lesson 10.2 completed.
- Current lesson moved to Lesson 10.3.

### Lesson 10.3 Scope

Allow authorized users to correct production-entry mistakes without silently editing historical records.

#### Business Context

Manual production records contain mistakes: wrong quantity, wrong workstation, wrong product, wrong employee or wrong timestamp. Leaders and managers need to fix these mistakes while preserving who changed what and why.

#### Problem

The system stored production entries as immutable rows, but it had no correction workflow. Directly updating `production_entries` would hide the original value and weaken trust in production history.

#### Design Discussion

Corrections are append-only records in `production_entry_corrections`. A correction stores the full corrected production-entry snapshot plus `reason`, `actorUserId` and `createdAt`.

The original `production_entries` row is never updated by the correction workflow. This is a narrow audit model for production-entry corrections only, not a generic audit framework.

Admins, managers and leaders can append and review corrections. Workers can register production, but cannot correct historical records.

#### Go Concepts

- DTOs for correction workflows
- request principal extraction from context
- domain validation for append-only audit records
- error translation for correction-specific failures

#### Architecture Concepts

- append-only correction records
- auditability through actor and reason tracking
- correction workflow in the production vertical slice
- route-level RBAC for historical mutation workflows

### Lesson 10.3 Completion Notes

#### Business Context

MES Lite now supports audit-safe correction of production-entry mistakes.

#### Problem

Production records could be entered and reviewed, but mistakes could not be corrected without either leaving the mistake in place or silently changing history.

#### Design Discussion

Added a dedicated correction table and production correction model. Each correction references the original entry, stores corrected replacement values and records the authenticated actor plus a required reason.

The correction endpoint does not update the original entry. Review clients can read correction history separately and decide how to present original versus corrected values.

#### Implementation

- Added migration `0010_create_production_entry_corrections.sql`.
- Added `production.Correction` and `ErrInvalidCorrection`.
- Added `CorrectEntryCommand`.
- Added `Service.CorrectEntry` and `Service.ListCorrections`.
- Added `Store.SaveCorrection` and `Store.ListCorrections`.
- Implemented in-memory correction persistence for fast tests.
- Added sqlc queries for creating and listing corrections.
- Implemented PostgreSQL correction persistence and row mapping.
- Added `POST /production-entries/{id}/corrections`.
- Added `GET /production-entries/{id}/corrections`.
- Tracked correction actor from the authenticated request principal.
- Protected correction routes with admin/manager/leader RBAC.

#### Tests

- Added service test proving corrections append without changing the original entry.
- Added service tests for missing original entries and missing correction reasons.
- Added handler test proving actor tracking comes from request context.
- Added handler test for missing principal rejection.
- Added PostgreSQL test for saving and listing correction history.
- Added PostgreSQL test for missing original-entry correction failure.
- Added server route test proving leaders can correct entries.
- Added server route test proving workers cannot correct entries.
- Added server route test for reading correction history.
- Verified with `go test ./internal/production ./internal/server -count=1`.

#### Refactoring

No generic audit package was introduced. The correction model exists because production-entry correction is a concrete MVP workflow with actor/reason requirements.

After review, the production slice split correction behavior out of the entry-oriented interfaces. `Service` now depends on separate `EntryStore` and `CorrectionStore` interfaces, and `Handler` now depends on separate `EntryRegistrar` and `CorrectionRegistrar` interfaces. This keeps correction behavior explicit without introducing a generic audit abstraction.

#### Code Review

An experienced Go engineer would approve the append-only direction because it preserves historical facts and makes corrections explicit. The main follow-up for L10.4 is MVP-level review: decide whether the API should expose an effective/current production-entry view or keep corrections as separate history for now.

#### Exercises

- Add a test proving multiple corrections are returned newest first.
- Design a read model that returns original entry plus latest correction as an effective view.
- Explain why correction reason should be required even for simple typo fixes.

#### Interview Questions

- Why are append-only audit records safer than direct updates for historical data?
- What belongs in a correction record versus the original entry record?
- Why should actor identity come from authentication context rather than request body?
- When would a generic audit framework become justified?

#### Roadmap Update

- Lesson 10.3 completed.
- Current lesson moved to Lesson 10.4.

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

- Added `EntryResponse`.
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

Added response DTOs while preserving route behavior. Because the API is not used by external clients yet, employee and product response field names were normalized to lower camel case instead of preserving the old capitalized names that came from domain struct encoding.

#### Implementation

- Added `employees.EmployeeResponse`.
- Added `products.ProductResponse`.
- Updated employee create, list, update and deactivate handlers to return DTOs.
- Updated product create, list, update, deactivate and search handlers to return DTOs.
- Added mapping helpers for employee and product responses.
- Preserved request DTOs.
- Normalized employee/product response JSON field names to lower camel case.

#### Tests

- Updated employee handler tests to decode `EmployeeResponse`.
- Updated product handler tests to decode `ProductResponse`.
- Added response field-name assertions for employee and product create responses.
- Verified focused employee/product/server tests.

#### Refactoring

The cleanup did not change domain encapsulation or persistence. It only separated HTTP response contracts from domain structs.

#### Code Review

An experienced Go engineer would approve this as a focused boundary cleanup. Since no external client depends on the API yet, changing employee/product responses to lower camel case is acceptable and makes the API more consistent with existing request DTOs and production responses.

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

### Milestone 10 Review

#### Architecture Review

An experienced Go engineer would approve Milestone 10 as an MVP closure milestone. Production-entry registration is retry-safe, review endpoints are protected and filterable, and correction workflows preserve history through append-only records.

The architecture remains vertical-slice oriented. Production owns entry registration, review and correction behavior. Auth owns identity and principal propagation. PostgreSQL owns uniqueness and reference-integrity guardrails.

#### Code Review

The production slice now has clearer boundaries after review. Entry and correction persistence interfaces are split. Handler-facing entry and correction registrar interfaces are split. Production, employee and product HTTP responses now use DTOs instead of exposing domain structs directly.

The code avoids a generic audit framework, generic transaction manager or formal workstation model because none are required by MVP business workflows.

#### Refactoring

L10.4 refactored production HTTP responses to explicit DTOs. L10.5 extended the same response-boundary cleanup to employee and product handlers. Earlier L10.3 review split correction behavior out of entry-oriented interfaces.

#### Interview Review

You should now be able to discuss idempotent write endpoints, partial unique indexes, append-only audit records, actor tracking through request context, REST route naming for subresources and why API DTOs protect contracts from domain-model churn.

#### Completion Criteria

- Managers and leaders can list/review production entries.
- Review filters cover employee, product, workstation and time range.
- Production registration requires and enforces `requestId` idempotency.
- Retry with identical request data returns the original entry.
- Reusing a request ID for different data returns `409 Conflict`.
- Corrections are append-only and preserve original production-entry rows.
- Corrections record reason and authenticated actor user ID.
- Workers cannot review or correct historical production entries.
- Workstation remains a text field for MVP.
- Production HTTP responses use explicit DTOs.
- Employee and product HTTP responses use explicit DTOs.
- Tests, build, lint and sqlc generation pass.

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

✅ Completed

### Lessons

- **L11.1** — Background Job Model & In-Memory Queue ✅
- **L11.1a** — Shared UUID Identifier Generator ✅
- **L11.1b** — Platform Package Tier & Dependency Rules ✅
- **L11.2** — Worker Pool & Job Execution ✅
- **L11.3** — Job Progress, Cancellation & Status API ✅
- **L11.4** — Async CSV Import Job ✅
- **L11.5** — Race Detection, Shutdown & Concurrency Review ✅

### Lesson 11.1 Scope

Introduce the background job vertical slice with a job model and an in-memory, channel-backed queue. No goroutines run work yet.

#### Business Context

CSV imports of historical production data run inside the HTTP request. A large file holds the request open until every row is validated and persisted, and the client learns nothing until the whole file is finished. Managers need to hand work to the server, get an answer immediately and check on the result later.

#### Problem

Asynchronous processing needs somewhere to put accepted work before anything can process it. Starting with goroutines would mean designing the job model, the handoff and the worker lifecycle at the same time, which is three hard problems in one step.

#### Design Discussion

Lesson 11.1 builds only the data that concurrency will move around: a `jobs.Job` value type and a `jobs.Queue` that accepts and hands out work safely. Lesson 11.2 adds the workers that consume it.

The queue owns two pieces of state with different responsibilities. A map guarded by `sync.RWMutex` holds the canonical status of every tracked job. A buffered channel is the handoff between producers and consumers. The channel carries job IDs rather than job values, so a consumer always reads current state from the map instead of acting on a copy taken at enqueue time.

The queue never closes its handoff channel. There are many producers, and closing a channel other goroutines may still send on is a panic rather than a shutdown signal. Shutdown closes a dedicated `done` channel with exactly one closer.

`Enqueue` never blocks. A full queue returns `ErrQueueFull` immediately, so an HTTP caller gets backpressure instead of waiting on a channel that may stay full for minutes. The buffer capacity is the burst the application is willing to absorb.

Durability is deferred and documented in ADR 0004: jobs are lost on restart.

#### Go Concepts

- buffered channels as a bounded handoff between producers and consumers
- channel ownership: who is allowed to close a channel, and why a `done` channel is used instead
- `select` with `default` for non-blocking send and receive
- `select` picking randomly among ready cases, and what that means for draining
- `sync.RWMutex` guarding shared map state read by many goroutines
- defensive copies of slices shared across goroutines

#### Architecture Concepts

- background job queue as a producer/consumer boundary
- job status owned by one component instead of spread across workers
- in-memory infrastructure chosen deliberately, with the durability trade-off recorded in an ADR

### Lesson 11.1 Completion Notes

#### Business Context

MES Lite now has the first building block for moving long-running work out of HTTP requests.

#### Problem

The application had no representation of deferred work and no safe place to hold it between the request that accepts it and the code that will run it.

#### Design Discussion

Added `internal/jobs` with a `Job` value type and a `Queue`. `Job` carries an ID, a type, a lifecycle status, an opaque payload and lifecycle timestamps. The payload is `[]byte` so the queue stays independent of any single workload; generics remain postponed.

Status transitions are intentionally absent. A job is created queued, and nothing moves it yet, because moving a job to running belongs with the worker that does the moving in Lesson 11.2.

The queue hands out copies. `Find` and `Dequeue` return `job.clone()`, and `NewJob` copies the payload it is given, so no caller can mutate tracked state through a shared slice.

Two ordering details required care:

- `Enqueue` holds the write lock across the channel send. The send is non-blocking, so it cannot deadlock, and holding the lock means a consumer that receives the ID cannot look it up before the map entry exists.
- `Dequeue` prefers waiting work over the close signal. It tries a non-blocking receive first, and when the `done` case fires it looks one more time, because `select` picks randomly among ready cases. Since `Enqueue` checks `closed` under the same mutex that `Close` writes it under, no send can follow `close(done)`, which makes that final look exact.

#### Implementation

- Added `internal/jobs/job.go` with `Job`, `Status`, `Type`, `NewJob`, `NewJobID` and `Validate`.
- Added `Status.Valid` and `Status.Terminal` for the lifecycle states queued, running, succeeded, failed and cancelled.
- Added `TypeProductionEntryImport` as the first planned workload.
- Added `internal/jobs/queue.go` with `NewQueue`, `Enqueue`, `Dequeue`, `Find`, `Len`, `Capacity` and `Close`.
- Added `ErrInvalidJob`, `ErrQueueFull`, `ErrQueueClosed`, `ErrNotFound` and `ErrAlreadyExists`.
- Added defensive payload copies on creation and on every read.
- Added ADR `0004-introduce-background-jobs.md`.

#### Tests

- Tested job normalization, payload copying and table-driven validation.
- Tested status validity and terminal states.
- Tested FIFO enqueue/dequeue, duplicate IDs and rejection of non-queued jobs.
- Tested that a full queue returns `ErrQueueFull` without blocking and leaves no tracked job behind.
- Tested that `Dequeue` blocks on an empty queue and returns once work arrives.
- Tested context cancellation while waiting.
- Tested that `Close` drains accepted work before reporting closed, rejects new work, releases waiting consumers and is idempotent.
- Tested that `Find` returns a copy.
- Added a concurrency test with 8 producers, 200 jobs and 4 consumers asserting exactly-once delivery.
- Verified with `go build ./...`.
- Verified with `go vet ./internal/jobs`.
- Verified with `go test ./internal/jobs -count=5 -race -shuffle=on`.
- Verified with `go test ./... -count=1`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

No existing package changed. The CSV import slice keeps its synchronous path until Lesson 11.4 gives it an asynchronous one, so the current endpoint keeps working while the job machinery is built beside it.

No `Store` interface was introduced for jobs. Nothing consumes the queue yet, and interfaces belong to consumers.

#### Code Review

An experienced Go engineer would approve the scope and the channel-ownership decisions: the handoff channel is never closed by a producer, shutdown has one closer, shared state is copied on the way in and on the way out, and the race detector is clean across shuffled repeat runs.

Two caveats are deliberate. `NewJobID` duplicates the UUID-shaped generator already present in the production and orders slices; a shared helper is worth considering, but a fourth copy is a weak reason to create a cross-slice utility package. Jobs are in-memory only, which ADR 0004 records rather than hides.

#### Exercises

- Remove the final `tryReceive` inside the `done` case of `Dequeue` and write a test that shows a job can be stranded at shutdown.
- Change `Enqueue` to block with `select` on `ctx.Done()` instead of returning `ErrQueueFull`, and describe what an HTTP client would experience under load.
- Explain why holding a mutex across the channel send in `Enqueue` is safe here, and what would make it a deadlock.
- Replace the ID channel with a `chan Job` and identify what becomes stale once workers start updating status.

#### Interview Questions

- Who should close a channel, and why is closing from the receiver side wrong?
- What does a buffered channel give you that an unbuffered one does not?
- What happens when several cases in a `select` are ready at the same time?
- When is a mutex a better fit than a channel for shared state?
- Why does a queue that hands out copies avoid a whole class of data races?

#### Roadmap Update

- Lesson 11.1 completed.
- Current lesson moved to Lesson 11.2.
- Concurrency `Channels`, `select` and `Race Detector` marked complete in the Knowledge Matrix.
- Standard Library `sync` marked complete in the Knowledge Matrix.
- Known technical debt updated for in-memory job durability and duplicated ID generation.

### Lesson 11.1a Scope

Replace the three duplicated UUID-shaped identifier generators with one shared helper, and drop the error return that the standard library guarantees can never fire.

#### Business Context

None directly. Identifier generation is technical plumbing that every slice creating a record depends on, and it should behave identically everywhere.

#### Problem

`production.NewEntryID`, `orders.NewOrderID` and `jobs.NewJobID` were byte-for-byte identical apart from one word in an error message. The duplication was recorded as technical debt when the third copy was added in Lesson 11.1.

Duplication was not the only cost. The wrapper names had already started to lie: `Service.CorrectEntry` generated a *correction* ID by calling `NewEntryID()`, because that was the generator its package happened to own.

#### Design Discussion

The helper lives in `internal/ids`, following the precedent set by `internal/postgres`: a small technical package below the business slices, not a business slice of its own. It is named so call sites read `ids.New()` without stutter.

**Should the helper return an error?**

Since Go 1.24, `crypto/rand.Read` is documented to never return an error. It fills the buffer completely or crashes the program irrecoverably, because the operating system APIs it uses are documented not to fail on anything but legacy Linux systems.

That makes every `if err != nil` around the old generators unreachable code. Two options were weighed:

- `New() (string, error)` keeps the current shape and stays honest if the randomness source is ever swapped for one that can fail. The cost is that every caller keeps writing, reading and reviewing a branch the runtime can never take.
- `New() string` leans on the guarantee. It deletes roughly ten unreachable error branches across services and tests, and it matches the project's architecture principle that panics are acceptable only for unrecoverable failures. A machine with an unusable CSPRNG cannot serve requests at all.

`New() string` was chosen. "Errors are values" is not a reason to invent an error that cannot occur; it is a reason to be precise about which failures are real.

**Should the per-slice wrappers survive?**

They were removed rather than reduced to one-line delegates. Three exported functions that each forward to the same helper are still three things to name, document and keep in sync, and keeping `NewEntryID` would have preserved the misleading call in `CorrectEntry`. Call sites now say `ids.New()`, which claims nothing about which record the identifier is for.

The trade-off accepted here is churn: this touches two services, one repository and several test files, in exchange for one implementation and one test suite.

#### Go Concepts

- knowing when a standard-library contract changes, and letting the API reflect it
- not modelling impossible failures as errors
- `_, _ =` as an explicit, lint-visible discard of a documented non-error
- small technical packages below business slices
- package naming that avoids stutter: `ids.New`, not `ids.NewID`

#### Architecture Concepts

- deduplication as refactoring, not as premature abstraction
- shared plumbing kept free of business meaning
- removing an abstraction whose name had begun to lie

### Lesson 11.1a Completion Notes

#### Business Context

Identifier generation now has one implementation and one test suite for the whole application.

#### Problem

Three identical generators had accumulated across the production, orders and jobs slices, and one of them was being used for a record type its name did not describe.

#### Design Discussion

Added `internal/ids` with a single `New() string` returning a canonical version 4 UUID. Removed the three per-slice generators and pointed every call site at the helper.

The helper is deliberately ignorant of what it identifies. It has no notion of entries, orders or jobs, which is what makes it safe to share; the moment such a helper starts taking a "kind" argument, it has become a business concept in the wrong package.

`internal/auth` was left alone. It does not generate identifiers: the bootstrap administrator uses the fixed literal `"bootstrap-admin"`. When auth-user management arrives, it should call `ids.New()` rather than grow a fourth generator.

#### Implementation

- Added `internal/ids/ids.go` with `New() string`.
- Removed `production.NewEntryID`, `orders.NewOrderID` and `jobs.NewJobID`.
- Updated `production.Service.Register` and `production.Service.CorrectEntry` to call `ids.New()`.
- Updated `orders.Service` order creation to call `ids.New()`.
- Updated `csvimport` entry conversion to call `ids.New()`, removing an error return that could no longer fail.
- Removed the now-unreachable error branches at every call site.
- Dropped the `crypto/rand` and `encoding/hex` imports from three business packages.

#### Tests

- Added canonical-format tests: length, hyphen positions, version nibble, RFC 4122 variant nibble and lower-case hex.
- Added a uniqueness test over 1,000 identifiers.
- Added a concurrency test proving the generator needs no external synchronization.
- Removed the two per-slice generator tests that the shared suite replaces.
- Verified with `go build ./...`.
- Verified with `go vet ./...`.
- Verified with `gofmt -l ./cmd ./internal`.
- Verified with `go test ./... -count=1`.
- Verified with `go test ./internal/ids ./internal/jobs ./internal/production ./internal/orders -race -count=2`.
- Verified with `golangci-lint run ./...`.

PostgreSQL integration tests were compiled and vetted but skipped in this run, because no test database was available. The edits to `internal/production/postgres_store_test.go` are mechanical call-site replacements.

#### Refactoring

This lesson is the refactoring. It closes the duplication debt recorded in Lesson 11.1 and removes a naming problem that predated it.

The identifier format did not change, so no migration or data change is involved: the same UUID strings are produced by the same bit manipulation, from one place instead of three.

#### Code Review

An experienced Go engineer would approve the extraction and, more notably, the signature change. Returning an error that the standard library guarantees cannot occur is a small but real cost paid at every call site, and removing it makes the remaining error checks in those functions meaningful.

The discarded return in `_, _ = rand.Read(b[:])` deserves its comment. A silent discard would read as an oversight; the comment states the contract being relied on, which is what a reviewer needs in order to agree.

#### Exercises

- Find the last remaining place where a business slice builds an identifier without `ids.New()` and decide whether it should change.
- Replace `internal/ids` with `github.com/google/uuid` (already an indirect dependency) and argue whether the dependency earns its place.
- Explain what would have to be true about the randomness source for `New()` to need an error return again.
- Write a benchmark for `ids.New()` and predict where the allocations are before running it.

#### Interview Questions

- When is duplication cheaper than the abstraction that removes it, and how do you tell?
- Why is "errors are values" not an argument for returning an error that cannot happen?
- What does `crypto/rand.Read` do on failure since Go 1.24, and why was that changed?
- What distinguishes a legitimate shared technical package from a `utils` dumping ground?
- Why should a shared ID generator never take a "kind" or "prefix" argument?

#### Roadmap Update

- Lesson 11.1a completed.
- Current lesson moved to Lesson 11.2.
- Known technical debt updated: duplicated ID generation resolved.

### Lesson 11.1b Scope

Move technical packages into `internal/platform/` and enforce the dependency direction between tiers with `depguard`.

#### Business Context

None directly. This lesson makes the codebase communicate its own structure, which matters for every business change that follows.

#### Problem

`internal/` held thirteen packages of three different kinds with nothing distinguishing them: seven business slices, five technical packages and one composition root. The layout claimed Vertical Slice Architecture while listing `config` and `version` alongside `production` and `orders`.

The deeper problem was that the dependency direction was correct by habit rather than by rule. A graph over `internal/` showed the technical packages importing nothing internal, the slices importing each other, and `server` importing everything — but nothing prevented a future technical package from importing a business slice and inverting that.

#### Design Discussion

The directory split and the lint rules solve different halves of the problem. The directory tells a reader what a package is for; the rules tell the compiler what a package may do. The rules are the half that actually holds, and they would have been worth adding even without moving a file.

Two rules cover the direction that matters:

- `internal/platform/**` may not import a business slice or `internal/server`. A platform package must stay importable by every slice, which means importing none of them.
- Only `cmd/` may import `internal/server`. The composition root points at everything; nothing points back.

Nested `internal/` directories were considered for compiler-level enforcement and rejected for this purpose: nested `internal/` constrains *who may import a package*, not *what a package may import*, so expressing "platform must not import business" would mean burying the slices under an artificial subdirectory. It stays the right tool for the generated sqlc packages, which is recorded as a known gap.

Two classification calls were not obvious:

- **`jobs` moved to platform.** It is a generic queue over opaque `[]byte` payloads with no business imports. Putting it under the rule commits the project to keeping it that way: when Lesson 11.4 connects the queue to CSV import, the handler registration must live in the composition root rather than inside `jobs`. The constraint is the reason to move it, not a side effect.
- **`auth` stayed a business slice.** JWT signing and HTTP middleware look technical, but users, roles and login are business capabilities, and `production` imports `auth` for the correction actor. Splitting the slice to relocate its middleware would break a working boundary for a cosmetic gain.

#### Go Concepts

- import paths as the primary boundary marker in Go
- import cycles as the compiler's own architectural enforcement
- what nested `internal/` can and cannot express
- lint configuration as executable architecture documentation

#### Architecture Concepts

- tiers within a vertical-slice layout: platform, slices, composition root
- dependency direction enforced in CI rather than in review
- classification rules written down so future packages are placed consistently

### Lesson 11.1b Completion Notes

#### Business Context

The repository now states its own structure, and CI fails when that structure is broken.

#### Problem

Thirteen packages sat in one flat directory with three different roles and an unstated dependency direction.

#### Design Discussion

`config`, `ids`, `jobs`, `postgres` and `version` moved to `internal/platform/`. Business slices and the composition root did not move. No code changed beyond import paths, so the commit has no behavioral effect.

The two depguard rules were verified to fire rather than assumed to work. A rule that matches no files also reports zero issues, so each rule was probed with a deliberate violation and the failure message checked before the probe was removed.

Probing exposed something worth knowing: importing `internal/server` from any existing package is already a compile error, because `server` imports every slice and the import becomes a cycle. The second rule therefore adds little today. It was kept because it applies to packages that do not exist yet, which is where the mistake would actually happen — a new package the server does not import could import the server without any cycle.

The first rule is not redundant. `platform/jobs` importing `production` produces no cycle, so nothing but the rule prevents it.

#### Implementation

- Created `internal/platform/` and moved `config`, `ids`, `jobs`, `postgres` and `version` into it with `git mv`.
- Rewrote the import paths in every affected package and restored import grouping with `goimports`.
- Added the `depguard` linter to `.golangci.yml`.
- Added rule `platform-stays-business-free`: no business slice or composition-root imports from `internal/platform/**`.
- Added rule `composition-root-is-not-a-dependency`: `internal/server` importable only from `cmd/`.
- Added `docs/package-layout.md` describing the tiers, the rules and the test for classifying a new package.
- Added ADR `0005-platform-package-tier.md`.
- Amended ADR 0001 to record that it is extended, and ADR 0004 to record the new path for `jobs`.

#### Tests

- Probed rule one with `platform/jobs` importing `production`; confirmed the depguard failure, then removed the probe.
- Probed rule two with a throwaway package importing `internal/server`; confirmed the depguard failure, then removed the probe.
- Confirmed a clean tree afterwards with `git status`.
- Verified with `go build ./...`.
- Verified with `go vet ./...`.
- Verified with `gofmt -l` and `goimports -l` over `cmd` and `internal`.
- Verified with `go test ./... -count=1`.
- Verified with `go test ./internal/platform/... -race -count=2`.
- Verified with `golangci-lint run ./...`.

PostgreSQL integration tests were compiled and vetted but skipped, because no test database was available.

#### Refactoring

This lesson is the refactoring. No package gained or lost behavior; only paths and lint configuration changed.

Historical lesson notes still refer to the old paths. They were left alone: they record what was true when written, and rewriting them would turn a project log into a fiction.

#### Code Review

An experienced Go engineer would approve the rules more readily than the move. The move is a matter of taste that reasonable engineers split on; the rules are a guarantee the codebase did not previously have.

The reviewable weakness is the deny lists: each business slice is named explicitly, so adding a slice means editing `.golangci.yml`, and forgetting to leaves a hole. An allow-list on the platform rule would fail closed instead, at the cost of naming every permitted standard-library and third-party import. The deny list was kept for now and the maintenance burden written into ADR 0005 rather than hidden.

#### Exercises

- Add a business slice, deliberately forget the deny list, and confirm the hole exists.
- Rewrite `platform-stays-business-free` as an allow-list and judge which version you would rather maintain.
- Explain why `platform/jobs` importing `production` compiles but `production` importing `server` does not.
- Decide where an `internal/platform/telemetry` package would sit if it needed to log a production entry ID.

#### Interview Questions

- What does Go's `internal/` directory actually enforce, and what does it not?
- How does an import cycle act as an architectural constraint?
- Why is a directory structure a weaker guarantee than a lint rule?
- How would you verify that a lint rule is actually being applied?
- When does a shared technical package become a business package?

#### Roadmap Update

- Lesson 11.1b completed.
- Current lesson moved to Lesson 11.2.
- Known technical debt updated for the unenforced sqlc boundary and the depguard deny-list maintenance.

### Lesson 11.2 Scope

Add a worker pool that consumes queued jobs and executes registered handlers concurrently.

#### Business Context

The application can now accept background jobs, but accepted work never runs. A manager could enqueue a future CSV import job, yet no worker would pick it up, mark it running or record whether it succeeded.

#### Problem

The queue is only the producer/consumer handoff. It stores job state and hands out queued jobs, but it does not own execution. Running handlers inside the queue would mix scheduling, state tracking and workload behavior in one type.

#### Design Discussion

Lesson 11.2 introduces a concrete `WorkerPool` in `internal/platform/jobs`. The pool starts a fixed number of goroutines, each repeatedly dequeuing jobs and executing the handler registered for the job type.

The queue remains the owner of job status. Workers ask the queue to move a job from queued to running, then to succeeded or failed. This keeps status updates serialized behind the queue mutex instead of letting workers mutate copies.

Handlers are registered by job type. Unknown job types fail the job instead of panicking, because malformed or unsupported work should be visible as job state.

Cancellation is limited to worker lifecycle in this lesson. User-requested job cancellation and progress updates belong to Lesson 11.3.

#### Go Concepts

- goroutines for concurrent background execution
- `sync.WaitGroup` for waiting on worker shutdown
- `sync.Once` for idempotent start and stop behavior
- context cancellation for worker lifecycle
- error values recorded as failed job state

#### Architecture Concepts

- worker pool as execution infrastructure
- queue owns state; handlers own business work
- fixed concurrency limit instead of unbounded goroutine creation
- workload registry without making the queue depend on business slices

#### Implementation Plan

- Add queue status-transition methods for running, succeeded and failed states.
- Add a concrete `WorkerPool` with fixed worker count.
- Add job handler registration by `jobs.Type`.
- Execute each job exactly once after successful dequeue.
- Record handler success or failure back into the queue.
- Keep progress and external cancellation out of scope until Lesson 11.3.

#### Tests

- Status-transition tests for valid and invalid lifecycle moves.
- Worker execution success and failure tests.
- Unknown job type becomes a failed job.
- Multiple workers process multiple jobs without duplicate execution.
- Shutdown waits for accepted running work and releases workers.
- Race detector remains clean.

#### Exercises

- Add a second job type and register a second handler.
- Change the pool to launch one goroutine per job and explain why that removes backpressure.
- Add a test proving `Stop` is safe to call more than once.
- Explain why handlers receive a job copy instead of a pointer to queue state.

#### Interview Questions

- What problem does a worker pool solve?
- Why use `sync.WaitGroup` instead of sleeping until workers finish?
- Why should queues not execute business handlers directly?
- What happens if a goroutine writes shared state without synchronization?
- How does context cancellation differ from closing a work channel?

### Lesson 11.2 Completion Notes

#### Business Context

MES Lite can now execute accepted background jobs with bounded concurrency.

#### Problem

Jobs could be enqueued and dequeued, but nothing marked them running, executed work or recorded success/failure. The system had a queue, not background processing.

#### Design Discussion

Added a concrete `WorkerPool` that consumes the existing in-memory queue. The pool owns goroutine lifecycle and handler dispatch, while the queue remains the owner of job state.

This split matters: workers receive job copies and ask the queue to perform status transitions. They do not mutate shared job structs directly, which keeps state changes synchronized behind the queue mutex.

The pool uses a fixed worker count. That makes concurrency a deliberate capacity choice instead of launching one goroutine per job and hoping the runtime absorbs the load.

Unknown job types are recorded as failed jobs rather than panicking. Unsupported work is operational data, not a process-level crash.

#### Implementation

- Added `ErrInvalidStatusTransition`.
- Added `Queue.MarkRunning`, `Queue.MarkSucceeded` and `Queue.MarkFailed`.
- Added `jobs.Handler` for registered background workload functions.
- Added `WorkerPool` with fixed worker count, `Start` and `Stop`.
- Used `sync.WaitGroup` to wait for worker goroutines.
- Used `sync.Once` so starting and stopping are idempotent operations.
- Recorded handler errors into failed job state.
- Kept progress updates, external cancellation and HTTP status routes out of scope for Lesson 11.3.

#### Tests

- Added queue status-transition tests for valid and invalid moves.
- Added worker-pool configuration validation tests.
- Tested successful handler execution and lifecycle timestamp updates.
- Tested handler failure becomes failed job state.
- Tested unknown job types fail jobs instead of panicking.
- Tested multiple workers execute 100 jobs exactly once.
- Tested `Stop` waits for running work.
- Tested `Stop` is safe to call more than once.
- Verified with `go test ./internal/platform/jobs -count=1`.
- Verified with `go test ./internal/platform/jobs -race -count=3 -shuffle=on`.
- Verified with `go build ./...`.
- Verified with `go test ./... -count=1`.
- Verified with `go vet ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

The queue was extended with explicit lifecycle methods instead of exposing a generic update hook. This keeps status transitions narrow and reviewable.

No job-store interface or durable implementation was introduced. ADR 0004 already records durability as a future trigger, and adding persistence now would distract from goroutines, `WaitGroup` and worker lifecycle.

#### Code Review

An experienced Go engineer would approve the core shape: fixed worker count, queue-owned synchronized state, no unbounded goroutine creation, copied jobs passed to handlers and race-detector coverage.

The main caveat is expected lesson scope. Workers can record success/failure, but clients still cannot query job status over HTTP, observe progress or request cancellation. That is the next lesson.

#### Exercises

- Add a second job type and prove it dispatches to a different handler.
- Add a test where `Stop` times out while a handler is blocked, then explain what happens to the running goroutine.
- Change workers to launch a goroutine per dequeued job and measure how many jobs can run at once.
- Explain why status transitions should stay queue-owned instead of handler-owned.

#### Interview Questions

- What does a worker pool protect a service from?
- Why is `sync.WaitGroup` the right primitive for waiting on worker goroutines?
- Why is a fixed worker count different from a buffered queue capacity?
- Why should handler failures become job state instead of only logs?
- What data race would appear if handlers received pointers to tracked jobs?

#### Roadmap Update

- Lesson 11.2 completed.
- Current lesson moved to Lesson 11.3.
- Concurrency `Goroutines`, `WaitGroup` and `Worker Pools` marked complete in the Knowledge Matrix.
- Testing `Race Detection` marked complete in the Knowledge Matrix.
- Known technical debt updated for missing job progress, cancellation and status API.

### Lesson 11.3 Scope

Expose background job status over HTTP, add progress reporting and allow authorized users to request cancellation.

#### Business Context

Managers need visibility into accepted long-running work. Once imports move to background jobs, returning only a job ID is not enough; clients must be able to check whether work is queued, running, finished, failed or cancelled.

#### Problem

Workers can execute jobs, but the outside world cannot observe job state or ask for a job to stop. A running handler also has no safe way to report progress back to the queue.

#### Design Discussion

The queue remains the owner of canonical job state. It gains progress and cancellation state transitions, while the worker pool owns running job contexts. Cancelling a queued job marks it cancelled immediately. Cancelling a running job records the request and cancels that job's context so a cooperative handler can stop.

The HTTP API is intentionally small:

- `GET /jobs/{id}` returns one job status snapshot.
- `PUT /jobs/{id}/cancel` requests cancellation and returns the updated snapshot.

Only admins and managers can read and cancel jobs for now. Workers should not inspect system-wide background processing, and leaders do not currently own import/system-operation workflows.

#### Go Concepts

- cooperative cancellation with `context.Context`
- per-job cancel functions guarded by a mutex
- progress updates as synchronized state changes
- HTTP status mapping for asynchronous job state
- idempotent cancellation for terminal jobs

#### Architecture Concepts

- status API as an operational read boundary
- queue owns state; worker pool owns execution contexts
- cancellation request separated from immediate termination
- no durable job persistence until ADR 0004's trigger is reached

#### Implementation Plan

- Add progress and cancellation fields to `jobs.Job`.
- Add queue methods for progress reporting and cancellation state.
- Extend `WorkerPool` with per-job cancellation.
- Add a jobs HTTP handler for status and cancellation.
- Register protected job routes in the server.
- Wire an in-memory queue and worker pool in the production composition root.

#### Tests

- Queue progress validation and defensive-copy tests.
- Queued cancellation transitions directly to cancelled.
- Running cancellation marks a request and worker context cancellation marks final cancelled state.
- Job status API returns stable JSON DTOs.
- Job cancellation API maps not found and terminal states correctly.
- Server authorization tests for admin/manager access and worker rejection.
- Race detector remains clean.

#### Exercises

- Add a test handler that reports progress from 0 to 100 and is cancelled halfway.
- Decide whether leaders should be allowed to read job status once report generation jobs exist.
- Add a list endpoint for recent jobs and discuss pagination before implementing it.
- Explain why cancellation in Go is cooperative rather than forcibly killing a goroutine.

#### Interview Questions

- Why does Go use context cancellation instead of killing goroutines externally?
- What must a handler do for cancellation to be effective?
- Why should progress updates be synchronized by the queue?
- What is the difference between cancellation requested and cancelled?
- Why can terminal job cancellation be idempotent?

### Lesson 11.3 Completion Notes

#### Business Context

MES Lite now has visibility and control for individual background jobs. A client can read a job status snapshot and authorized users can request cancellation.

#### Problem

The worker pool could execute jobs, but clients had no API for observing job state, and running handlers had no synchronized way to report progress or respond to a user cancellation request.

#### Design Discussion

The queue remains the canonical state owner. It now records progress, cancellation requests and final cancelled state under the same mutex that protects job status.

The worker pool owns execution contexts. Each running job gets a child context with its own cancel function. `WorkerPool.Cancel` records the cancellation request in the queue and, if the job is running, calls that job's cancel function.

Cancellation is cooperative. The pool does not kill goroutines. A handler must observe `ctx.Done()` and return. This is idiomatic Go because forcibly stopping goroutines would risk leaving locks, transactions or files in inconsistent states.

The status API is deliberately small: `GET /jobs/{id}` and `PUT /jobs/{id}/cancel`. A list endpoint is postponed until there is real UI/API pressure to define pagination, filtering and retention rules.

#### Implementation

- Added `Progress` and `CancelRequested` to `jobs.Job`.
- Added progress validation to `Job.Validate`.
- Added `ErrInvalidProgress`.
- Added `Queue.ReportProgress`.
- Added `Queue.RequestCancellation`.
- Added `Queue.MarkCancelled`.
- Added `Queue.CancellationRequested`.
- Extended `WorkerPool` with a mutex-protected map of running job cancel functions.
- Added `WorkerPool.Cancel`.
- Changed worker execution to pass per-job contexts to handlers.
- Added `jobs.HTTPHandler` with `GET /jobs/{id}` and `PUT /jobs/{id}/cancel` behavior.
- Added `server.RegisterJobRoutes` protected by admin/manager RBAC.
- Wired an in-memory queue, worker pool and job handler in `cmd/server`.

#### Tests

- Added queue progress tests.
- Added invalid progress and invalid-state progress tests.
- Added queued and running cancellation queue tests.
- Added worker-pool queued cancellation test.
- Added worker-pool running cancellation test proving handler context cancellation becomes final cancelled state.
- Added job HTTP handler tests for status, not found and cancellation.
- Added server route tests proving managers can read job status, admins can cancel jobs and workers are forbidden.
- Verified with `go test ./internal/platform/jobs ./internal/server -count=1`.
- Verified with `go test ./internal/platform/jobs -race -count=3 -shuffle=on`.
- Verified with `go build ./...`.
- Verified with `go test ./... -count=1`.
- Verified with `go vet ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

`server.New` was not widened with another constructor argument. Instead, `server.RegisterJobRoutes` registers operational job routes when a jobs handler exists. This keeps the existing application route constructor stable while L11.4 decides how background jobs are wired into CSV import.

No durable job persistence, job listing or retention policy was introduced. Those would be separate product/API decisions rather than requirements for cancellation mechanics.

#### Code Review

An experienced Go engineer would approve the cooperative-cancellation model, queue-owned state transitions and race-detector coverage. The main implementation risk in this area is handlers that ignore their context; L11.4 must make the async CSV import handler check cancellation at natural boundaries.

The main API caveat is that only per-job status exists. That is enough once an enqueue endpoint returns a job ID, but a UI will eventually want listing and retention semantics.

#### Exercises

- Modify a test handler to ignore `ctx.Done()` and observe why cancellation request does not immediately stop work.
- Add a `GET /jobs` proposal with `status`, `limit` and `offset` filters, but do not implement it yet.
- Add a progress monotonicity rule and decide whether decreasing progress should be rejected.
- Explain why `WorkerPool.Cancel` records queue state before calling the running cancel function.

#### Interview Questions

- Why is goroutine cancellation cooperative in Go?
- What is the difference between cancelling a context and closing a channel?
- Why should a queue own status/progress updates instead of handlers mutating job pointers?
- What race would appear if the worker pool's running-job map had no mutex?
- Why might a cancelled job still show partial progress?

#### Roadmap Update

- Lesson 11.3 completed.
- Current lesson moved to Lesson 11.4.
- Known technical debt updated: job status/cancellation exists for individual jobs; job listing and async CSV wiring remain pending.

### Lesson 11.4 Scope

Add an asynchronous CSV import endpoint backed by the background job queue and worker pool.

#### Business Context

Managers should be able to upload historical production CSV data, receive a job ID quickly and monitor import completion through the job status API instead of keeping one HTTP request open for the whole import.

#### Problem

The CSV import pipeline is streaming internally, but the HTTP endpoint still runs validation and persistence synchronously. Large imports keep the request open and make cancellation/status visibility impossible for clients.

#### Design Discussion

The synchronous endpoint remains available for simple MVP usage. Lesson 11.4 adds a separate async endpoint so the API contract is explicit:

- `POST /imports/production-entries/jobs` accepts CSV input and returns a queued job snapshot.
- `GET /jobs/{id}` lets clients monitor queued, running, succeeded, failed or cancelled state.
- `PUT /jobs/{id}/cancel` requests cancellation.

The job payload stores metadata pointing at a temporary upload file, not the whole CSV content. This is the smallest production-minded choice for the current in-memory job queue: queued job metadata stays small, while the CSV service can still stream from an `io.Reader` when the worker runs.

The temporary file is removed by the job handler after execution. This is not durable job storage; a process restart still loses queued/running jobs and any unprocessed temporary uploads, which remains covered by ADR 0004.

#### Go Concepts

- `os.CreateTemp` for temporary upload handoff
- `io.Copy` for streaming request body to disk
- JSON encoding for opaque job payload metadata
- worker handlers that call existing application services
- cancellation propagation from worker context into CSV import processing

#### Architecture Concepts

- async endpoint separated from synchronous endpoint
- composition root wires business work into platform workers
- queue payload carries metadata, not business dependencies
- temporary-file handoff as an implementation detail before durable queues exist

#### Implementation Plan

- Add async import service/handler behavior that writes upload data to a temporary file and enqueues a production-entry-import job.
- Add a job handler that decodes the payload, opens the temporary file and calls the existing CSV import service.
- Store import summary JSON in successful job state.
- Register `POST /imports/production-entries/jobs` with the same admin/manager RBAC as synchronous import.
- Wire the job handler into the worker pool in `cmd/server`.
- Preserve the existing synchronous import endpoint.

#### Tests

- Async handler enqueues a job and returns a job response.
- Queue-full errors become a clear HTTP response.
- Job handler imports valid CSV and records success result.
- Job handler records validation/import failure as failed job state.
- Temporary files are removed after job execution.
- Server route tests cover admin/manager access and worker rejection.
- Race detector remains clean for jobs and csvimport-focused tests.

#### Exercises

- Add a max upload size before writing the temporary file and decide which HTTP status should be returned.
- Simulate process restart after enqueue and explain why the job is lost today.
- Change the payload to store CSV bytes and compare memory behavior with temporary files.
- Add a test proving async import cancellation stops between CSV rows.

#### Interview Questions

- Why is writing an upload to a temporary file different from reading it all into memory?
- Why should the composition root connect job handlers to business services?
- What cleanup responsibilities appear when temporary files are used?
- Why keep the synchronous endpoint while adding an async endpoint?
- What would need to change for this to survive process restarts?

### Lesson 11.4 Completion Notes

#### Business Context

MES Lite can now accept a production-entry CSV import as background work. Managers get a job snapshot immediately and can use the job status endpoint from L11.3 to monitor completion.

#### Problem

The CSV import pipeline streamed rows internally, but the HTTP request still waited until validation and persistence finished. That made long imports poor API citizens and gave clients no status or cancellation handle.

#### Design Discussion

Added `POST /imports/production-entries/jobs` instead of changing the existing synchronous endpoint. This keeps the old simple workflow available and makes async behavior explicit in the route name.

The async service writes the uploaded CSV body to a temporary file and stores only the file path in the job payload. This keeps queued job payloads small and lets the worker reuse the existing streaming CSV service by opening the file as an `io.Reader`.

The worker handler lives in `csvimport`, while the worker pool remains in `internal/platform/jobs`. The composition root wires them together. This preserves the platform rule: the generic jobs package does not import business slices.

Successful jobs record the `ImportSummary` JSON as job result data. Failed jobs record the failure message through the existing worker-pool failure path. The temporary file is removed after execution in both success and failure paths.

#### Implementation

- Added `Result` to `jobs.Job` with defensive copying.
- Added `Queue.RecordResult` for running jobs.
- Added `result` to job HTTP responses.
- Added `csvimport.AsyncService`.
- Added async upload-to-temp-file handoff.
- Added `csvimport.NewProductionEntriesJobHandler` for worker execution.
- Added `Handler.ImportProductionEntriesAsync`.
- Added `Handler.AsyncEnabled` so routes are only registered when async import is wired.
- Added `POST /imports/production-entries/jobs` with admin/manager RBAC.
- Wired `cmd/server` to register the CSV import job handler with the worker pool.
- Added context cancellation checks to the CSV import streaming loop.

#### Tests

- Tested async service writes upload data to a temporary file and enqueues a job.
- Tested queue-full enqueue failures remove the temporary file.
- Tested the production-entry import job handler imports valid CSV data.
- Tested successful job execution records an import summary result.
- Tested temporary files are removed after successful execution.
- Tested invalid CSV fails the job and still removes the temporary file.
- Added server route test proving managers can enqueue async imports.
- Added server route test proving workers are forbidden.
- Verified with `go test ./internal/platform/jobs ./internal/csvimport ./internal/server -count=1`.
- Verified with `go test ./internal/platform/jobs ./internal/csvimport -race -count=3 -shuffle=on`.
- Verified with `go build ./...`.
- Verified with `go test ./... -count=1`.
- Verified with `go vet ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

The CSV import service remains the owner of row streaming, validation and persistence. Async import is a wrapper around the existing workflow rather than a second import implementation.

`server.New` now registers the async import route only when the CSV import handler is built with async support. Existing tests and simple synchronous wiring do not need to know about background jobs.

#### Code Review

An experienced Go engineer would approve the composition direction: platform jobs do not depend on `csvimport`, and business work is registered from the composition root.

The main caveat is the temporary-file durability story. If the process crashes after upload and before execution, the in-memory job is lost and the temp file may be orphaned. That is acceptable for this lesson because ADR 0004 already records the durable-queue trigger, but it is not production-complete durability.

#### Exercises

- Add a max async upload size using `http.MaxBytesReader` or an `io.LimitedReader` at the HTTP boundary.
- Add a cancellation test that stops a CSV import between rows.
- Add a status API test proving successful async import includes a JSON `result` object.
- Design how durable job storage would reference uploaded files or object storage.

#### Interview Questions

- Why does async HTTP usually return a job ID instead of the final result?
- Why should a background queue payload avoid large blobs?
- How does the composition root prevent `platform/jobs` from importing business packages?
- What cleanup failure modes are introduced by temporary files?
- Why does this still need durable storage before production-grade reliability?

#### Roadmap Update

- Lesson 11.4 completed.
- Current lesson moved to Lesson 11.5.
- Known technical debt updated for temporary-file/orphan behavior and missing job list endpoint.

### Lesson 11.5 Scope

Review the background job implementation for race safety, shutdown behavior and milestone readiness.

#### Business Context

Background jobs now execute real CSV import work. Before moving to machine integration, the concurrency foundation should be reviewed as if it were a pull request headed toward production.

#### Problem

Concurrency bugs often hide in lifecycle edges rather than happy-path job execution. The code needs focused review around worker startup, shutdown, cancellation, queue closing, temporary-file cleanup and race-detector coverage.

#### Design Discussion

This lesson is a review and hardening lesson. It should not add a new business workflow. Any code change should be small and directly tied to shutdown safety, data-race prevention or test coverage.

The main review target is worker-pool lifecycle coordination. `Start` and `Stop` are idempotent, but lifecycle calls should be explicitly serialized so `sync.WaitGroup.Add` cannot overlap incorrectly with `Wait` and so the worker-context cancel function is read and written under synchronization.

#### Go Concepts

- `sync.WaitGroup` lifecycle rules
- race detector as a verification tool, not a proof of correctness
- serialized lifecycle state with a mutex
- graceful shutdown versus forced cancellation
- code review for concurrent state ownership

#### Architecture Concepts

- concurrency review before milestone closure
- explicit shutdown contract for infrastructure components
- known durability gaps documented rather than hidden
- milestone review discipline before starting a new domain area

#### Implementation Plan

- Add a focused lifecycle guard to `WorkerPool` if review confirms the race risk.
- Add or update tests for concurrent start/stop and stop-before-start behavior.
- Run race detector across jobs and CSV import packages.
- Run full build, tests, vet and lint.
- Complete the Milestone 11 review and advance the roadmap to Milestone 12.

#### Tests

- Worker pool stop-before-start behavior.
- Concurrent `Start` and `Stop` does not race or panic.
- Existing cancellation and async import tests remain race-clean.
- Full project verification remains green.

#### Exercises

- Remove the lifecycle mutex and run the race test repeatedly to understand what it protects.
- Explain why `WaitGroup.Add` must not race with `Wait` when the counter can be zero.
- Add an intentionally context-ignoring handler and observe shutdown timeout behavior.
- List which job guarantees are lost on process restart.

#### Interview Questions

- What does the Go race detector detect, and what does it not prove?
- Why can `sync.WaitGroup` be misused even when there is no shared map?
- What is the difference between graceful shutdown and cancellation?
- Why is lifecycle state often protected by a separate mutex?
- What would make this background job system production durable?

### Lesson 11.5 Completion Notes

#### Business Context

Milestone 11 is complete. MES Lite now has a background-job foundation that can run long-lived work outside the original HTTP request and expose job status to clients.

#### Problem

The background job system worked on the happy path, but concurrency systems fail most often at lifecycle boundaries: startup, shutdown, cancellation and shared-state access.

#### Design Discussion

The review focused on the worker-pool lifecycle. `Start` and `Stop` were idempotent, but their internal coordination was not explicit enough for `sync.WaitGroup` rules. A `WaitGroup` must not have `Add` racing with `Wait` when the counter may be zero.

The worker pool now serializes lifecycle startup/shutdown state with a small mutex. This keeps worker registration, `started` state and worker-context cancellation coordinated without changing the public API.

The lesson also validated the async CSV path under the race detector. The job queue, worker pool, CSV import and server route tests are clean under race runs.

#### Implementation

- Added lifecycle serialization to `WorkerPool`.
- Added explicit `started` state so `Stop` before `Start` returns safely.
- Moved worker-context cancellation behind a synchronized helper.
- Preserved the existing `Start`, `Stop` and `Cancel` APIs.
- Kept durable job storage, job listing and retry policy out of scope.

#### Tests

- Added `WorkerPool` stop-before-start coverage.
- Added concurrent start/stop race coverage.
- Verified with `go test ./internal/platform/jobs -race -count=5 -shuffle=on`.
- Verified with `go test ./internal/platform/jobs ./internal/csvimport ./internal/server -race -count=2 -shuffle=on`; this timed out after jobs and csvimport passed because server race tests are slower.
- Re-ran server race tests with a longer timeout: `go test ./internal/server -race -count=1 -shuffle=on`.
- Verified with `go build ./...`.
- Verified with `go test ./... -count=1`.
- Verified with `go vet ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

The fix is intentionally narrow. No generic lifecycle manager or worker framework was introduced. The worker pool owns its own lifecycle because it is currently the only component that needs this coordination.

#### Code Review

An experienced Go engineer would approve the concurrency foundation for the current project stage: shared job state is queue-owned and mutex-protected, worker concurrency is bounded, cancellation is cooperative, temp files are cleaned up on normal success/failure paths and race tests cover the key packages.

The main production-readiness caveat is durability. In-memory jobs and temp-file payloads are acceptable for this learning milestone but not for reliable production processing across restarts.

#### Exercises

- Remove the lifecycle mutex and run the concurrent start/stop test with `-race` repeatedly.
- Add an intentionally context-ignoring handler and observe `Stop` timeout behavior.
- Design a durable jobs table schema with status, payload, result and retry metadata.
- Explain which parts of the current system should change if jobs move to PostgreSQL.

#### Interview Questions

- What does the race detector find, and what can it miss?
- Why is `WaitGroup.Add` versus `Wait` ordering important?
- Why is cooperative cancellation safer than forcibly killing a goroutine?
- How do you decide whether to use channels, mutexes or both?
- What guarantees are required for durable background jobs?

#### Roadmap Update

- Lesson 11.5 completed.
- Milestone 11 completed.
- Current milestone moved to Milestone 12.
- Current lesson moved to Lesson 12.1.
- Architecture maturity, Go knowledge progress and interview readiness updated.

### Milestone 11 Review

#### Architecture Review

An experienced Go engineer would approve Milestone 11 as a learning-oriented concurrency foundation. The implementation uses channels for producer/consumer handoff, mutexes for shared state, a worker pool for bounded execution and context cancellation for cooperative shutdown.

The most important architectural boundary is that `internal/platform/jobs` remains business-free. CSV import registers business work from the composition root instead of making the platform package import a slice.

#### Code Review

The code is explicit and reviewable. Queue state transitions are named methods instead of a generic update callback. Workers receive job copies instead of pointers to canonical state. The status API exposes DTOs, and successful async imports record the import summary as job result data.

The main weaknesses are known and documented: jobs are in-memory, temp-file payloads are not durable, orphan temp files are possible after process crashes and there is no job list/retention policy yet.

#### Refactoring

The milestone included two useful refactors: shared identifier generation moved to `internal/platform/ids`, and technical packages moved to `internal/platform/` with depguard rules enforcing dependency direction.

The final L11.5 refactor serialized worker-pool lifecycle state to make startup/shutdown behavior safer under concurrent calls.

#### Interview Review

You should now be able to explain goroutines, buffered channels, channel ownership, `select`, mutex-protected shared state, `sync.WaitGroup`, worker pools, cooperative cancellation, race detector usage, shutdown trade-offs and why durable background jobs require persistence beyond in-memory queues.

#### Completion Criteria

- Background job model implemented.
- In-memory queue implemented.
- Worker pool implemented.
- Job progress, result, status and cancellation implemented.
- Async CSV import job implemented.
- Worker shutdown and lifecycle reviewed.
- Race detector passes for concurrency-focused packages.
- Tests, build, vet and lint pass.
- Roadmap updated.

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

🔄 In Progress

### Lessons

- **L12.1** — Machine Event Model & Fake Machine API
- **L12.2** — Event Deduplication & Idempotent Processing
- **L12.3** — Synchronization With Mutex, RWMutex & sync.Once
- **L12.4** — Atomic Counters & Runtime Race Review
- **L12.5** — Machine Integration Review & Milestone Closure

### Lesson 12.1 Scope

Introduce the machine-integration vertical slice with a machine event model and a fake API endpoint for manually sending events into the system.

#### Business Context

Machine integration is future scope for the business, but it is a useful learning step after the MVP API and background-job foundations. Before processing concurrent machine events, the system needs a clear event shape and an endpoint that can simulate machine input without connecting to real CNC equipment.

#### Problem

The project had no concept of a machine event. Jumping straight to concurrent processing or deduplication would mix event modeling, HTTP input, idempotency and synchronization in one lesson.

#### Design Discussion

L12.1 adds `internal/machines` as a business slice. The fake API accepts a machine ID from the path and event details from JSON, validates them, normalizes timestamps to UTC and stores the event in memory.

The endpoint is deliberately fake and protected by existing admin/manager JWT RBAC. Real machine credentials, durable storage, deduplication and processing into production entries are postponed to later lessons so this lesson stays focused on the event model and HTTP boundary.

#### Go Concepts

- custom string type for event kinds
- constructor validation with `(T, error)`
- timestamp normalization with `time.Time.UTC`
- mutex-protected in-memory state for HTTP safety
- DTO mapping at the HTTP boundary

#### Architecture Concepts

- machine integration as a vertical slice
- event model separated from production-entry registration
- fake adapter endpoint before real integration protocols
- explicit postponement of deduplication and processing

### Lesson 12.1 Completion Notes

#### Business Context

MES Lite now has a fake machine-event intake path that can be used to simulate future machine integration work.

#### Problem

The application could register manual production entries and run background jobs, but it had no event shape for machine-originated production signals.

#### Design Discussion

Added `machines.Event` as an input event, not as a production entry. A machine event may later produce a production entry, but keeping them separate prevents future machine-specific fields and deduplication state from leaking into manual registration.

The in-memory store is intentionally temporary. It gives L12.2 and L12.3 a concrete event source to evolve without adding database schema and synchronization concerns at the same time.

#### Implementation

- Added `internal/machines` business slice.
- Added `EventType` with `cycle_completed`, `state_changed` and `alarm_raised`.
- Added `machines.Event` with generated ID, machine ID, external event ID, type, occurred-at timestamp, product SKU, quantity, workstation and message.
- Added constructor and validation for machine events.
- Added mutex-protected `machines.InMemoryStore`.
- Added `machines.Handler` with `POST /machines/{machineId}/events`.
- Registered machine routes separately through `server.RegisterMachineRoutes`.
- Wired the fake machine handler in `cmd/server`.
- Updated depguard's platform deny list for the new business slice.

#### Tests

- Added machine event constructor and validation tests.
- Added in-memory store save/list and defensive snapshot tests.
- Added handler tests for valid and invalid fake machine events.
- Added server route tests proving managers can submit fake machine events, workers are forbidden and unauthenticated callers receive `401`.
- Verified with `go test ./internal/machines ./internal/server -count=1`.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `go vet ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

Machine routes were registered through a focused `RegisterMachineRoutes` function instead of widening the already large `server.New` constructor. This matches the optional job-route registration pattern and keeps the main application route constructor stable.

#### Code Review

An experienced Go engineer would approve the lesson scope because it introduces the smallest useful machine-integration boundary and avoids pretending the fake API is production machine connectivity.

The main caveats are intentional: events are not durable, deduplicated or processed yet, and the fake endpoint uses human JWT roles instead of real machine authentication.

#### Exercises

- Add a handler test for `state_changed` without product SKU or quantity.
- Add a `GET /machines/{machineId}/events` proposal and decide who should be allowed to call it.
- Explain why machine events should not be stored directly as production entries at intake time.

#### Interview Questions

- What is the difference between an event and a command?
- Why should integration input be modeled separately from internal production-entry records?
- Why is a fake API useful before a real machine protocol exists?
- What consistency problems appear when machines retry the same event?

#### Roadmap Update

- Lesson 12.1 completed.
- Current lesson moved to Lesson 12.2.
- Known technical debt updated for fake, in-memory machine event intake.

### Lesson 12.2 Scope

Add duplicate detection and idempotent retry behavior for machine events based on the machine's external event identifier.

#### Business Context

Real machines and integration gateways retry messages when networks fail or acknowledgements are lost. MES Lite must not process the same machine signal twice just because delivery happened twice.

#### Problem

L12.1 accepted and stored fake machine events, but the same `machineId` and `externalEventId` could be submitted repeatedly and stored as separate events. That would inflate future production quantities once events start producing production entries.

#### Design Discussion

Deduplication belongs at the state boundary and idempotency belongs in the workflow. The in-memory store now rejects duplicate `(machineID, externalEventID)` keys. The service receives that storage-level duplicate signal, loads the existing event and compares business payloads.

An identical retry returns the original event. A duplicate key with different data returns `ErrEventConflict`, which the HTTP handler translates to `409 Conflict`.

This mirrors the earlier production-registration `requestId` lesson, but keeps machine events separate because the retry identity comes from the external machine system rather than from a human/API client command.

#### Go Concepts

- service-level idempotency workflow
- composite map keys for duplicate detection
- sentinel errors for duplicate versus conflict outcomes
- comparing values while ignoring generated identifiers
- `errors.Is` across handler and service boundaries

#### Architecture Concepts

- event deduplication as an integration boundary concern
- storage uniqueness as a concurrency guardrail
- idempotent retry semantics before event processing
- fake in-memory behavior documented before durable persistence

### Lesson 12.2 Completion Notes

#### Business Context

MES Lite now treats repeated machine delivery safely. A machine can retry the same event without creating duplicate machine-event records.

#### Problem

The fake machine API accepted all submissions independently. A retry with the same external event ID could create duplicate records and later duplicate production output.

#### Design Discussion

Added `machines.Service` as the intake workflow boundary. The handler now translates HTTP into a `ReceiveEventCommand`, while the service constructs the event and handles duplicate semantics.

The in-memory store owns uniqueness for `(machineID, externalEventID)` under its mutex. This matters because an application-only pre-check would be race-prone once concurrent machine submissions arrive.

Identical duplicate delivery returns the existing event. Different payload under the same external key returns `409 Conflict` so clients and operators can see that a producer reused an event identifier incorrectly.

#### Implementation

- Added `ErrDuplicateEvent`, `ErrEventConflict` and `ErrNotFound` to the machines slice.
- Added `Store.FindByExternalEventID`.
- Updated `InMemoryStore` to index events by `(machineID, externalEventID)`.
- Added `machines.Service.ReceiveEvent`.
- Added `ReceiveEventCommand`.
- Changed `machines.Handler` to depend on an `EventReceiver` interface instead of storing events directly.
- Mapped conflicting duplicate machine events to HTTP `409 Conflict`.
- Updated server and production composition wiring to use the machine service.

#### Tests

- Added store tests for duplicate external event rejection and lookup by external event ID.
- Added service tests for new event storage, identical retry returning the original event, conflicting retry and invalid events.
- Added handler tests for idempotent retry and conflicting duplicate payload.
- Existing server route authorization tests continue to pass with service-based wiring.
- Verified with `go test ./internal/machines ./internal/server -count=1`.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `go vet ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

Machine event construction and persistence moved out of the handler and into `machines.Service`. This keeps handlers focused on HTTP and gives later lessons one place to add processing, progress counters or durable idempotency.

#### Code Review

An experienced Go engineer would approve the separation of duplicate detection and idempotent semantics. The store enforces uniqueness atomically under a mutex, and the service decides whether a duplicate is safe or conflicting.

The main caveat is durability. In-memory deduplication is enough for the synchronization lesson, but a real machine integration would need a database uniqueness constraint or durable event log.

#### Exercises

- Add a concurrent test that submits the same machine event from two goroutines and proves only one record is stored.
- Add a test proving the same `externalEventId` may be reused by two different machines.
- Design the PostgreSQL unique index that would make machine-event deduplication durable.

#### Interview Questions

- What is idempotent event processing?
- Why is duplicate detection at the store boundary safer than a handler-level pre-check?
- Why should conflicting duplicate payloads return `409 Conflict`?
- How is machine event idempotency different from HTTP command idempotency?

#### Roadmap Update

- Lesson 12.2 completed.
- Current lesson moved to Lesson 12.3.
- Known technical debt updated: machine event deduplication exists but is still in-memory only.

### Lesson 12.3 Scope

Review and harden the machine-event in-memory synchronization path with `sync.Mutex`, `sync.RWMutex` and `sync.Once`.

#### Business Context

Machine integrations send events concurrently. Even the fake in-memory API must behave correctly when many requests arrive at the same time, because duplicate protection is only trustworthy if the shared state is protected consistently.

#### Problem

The in-memory store used a mutex, but initialization and concurrent retry behavior were not explicitly tested under race-oriented workloads. The zero-value store also relied on ad-hoc map initialization inside `Save`.

#### Design Discussion

The store remains the synchronization boundary. Writes use `Lock`, reads use `RLock` and one-time index initialization uses `sync.Once`. `sync.Once` makes the zero value usable and prevents multiple goroutines from racing to initialize the deduplication map.

The service still owns business semantics. It receives duplicate errors from the store and turns them into either idempotent success or conflict. This keeps synchronization separate from workflow rules.

#### Go Concepts

- `sync.Once` for one-time initialization
- `sync.Mutex` for exclusive writes
- `sync.RWMutex` for concurrent-safe snapshots and lookups
- race-detector verification for synchronization assumptions
- zero-value usability for stateful Go types

#### Architecture Concepts

- synchronization owned by the in-memory infrastructure boundary
- service-level idempotency kept separate from locking mechanics
- focused concurrency hardening without adding durable storage
- tests that exercise behavior likely to fail under concurrent delivery

### Lesson 12.3 Completion Notes

#### Business Context

Machine event intake is now safer under concurrent fake-machine submissions. Repeated delivery and concurrent writes are protected by synchronized store state.

#### Problem

Duplicate detection depends on shared maps and slices. Without explicit synchronization tests, a bug could appear only under concurrent machine traffic.

#### Design Discussion

Replaced the store's ad-hoc map initialization with `sync.Once`. The store now has a clear ownership model: `sync.Once` initializes internal indexes, `Lock` protects writes and `RLock` protects read snapshots and lookups.

This lesson did not introduce `sync.Map`. The store updates a slice and a map together, so a single mutex keeps those two structures consistent. `sync.Map` is useful for specific read-mostly or append-only patterns, but it would make this small store harder to reason about.

#### Implementation

- Added `sync.Once` initialization to `machines.InMemoryStore`.
- Preserved exclusive `Lock` for event saves and duplicate-key insertion.
- Preserved `RLock` for event listing and external-event lookup.
- Made the zero-value `InMemoryStore` safe for concurrent saves.
- Kept public machine service and handler behavior unchanged.

#### Tests

- Added concurrent zero-value store save test.
- Added concurrent identical retry test proving many callers receive one stored event and the same generated event ID.
- Verified focused package behavior with `go test ./internal/machines ./internal/server -count=1`.
- Verified race safety with `go test ./internal/machines -race -count=5 -shuffle=on`.
- Verified full project behavior with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `go vet ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

The refactor is intentionally narrow. The machine store remains in memory and keeps the same public API. Only internal initialization and concurrency coverage changed.

#### Code Review

An experienced Go engineer would approve using one `RWMutex` here because the store must keep `events` and `byExternalID` consistent. Splitting locks or using `sync.Map` would add complexity without a demonstrated bottleneck.

The race detector being clean is useful evidence, not a proof. It only checks executions that the tests exercise, so the tests deliberately include concurrent duplicate delivery and concurrent first-use initialization.

#### Exercises

- Remove `sync.Once` and run the zero-value concurrent save test with `-race`.
- Replace the store map with `sync.Map` and explain why keeping the slice consistent becomes less obvious.
- Add a concurrent test where one goroutine lists events while others save events.

#### Interview Questions

- When should you use `sync.Once`?
- Why can one mutex be better than separate locks for related pieces of state?
- What is the difference between `Mutex` and `RWMutex`?
- What does a clean race-detector run prove, and what does it not prove?

#### Roadmap Update

- Lesson 12.3 completed.
- Current lesson moved to Lesson 12.4.
- L12.4 remains focused on atomic counters and runtime race review.

### Lesson 12.4 Scope

Add atomic machine intake counters and review the machine integration slice under the race detector.

#### Business Context

Operators need lightweight visibility into fake machine event intake while the project prepares for real observability in Milestone 13. Even before Prometheus or tracing exists, the machine slice can expose basic counters that are safe under concurrent traffic.

#### Problem

The machine service could accept, deduplicate and reject events, but it did not count those outcomes. Adding counters with a shared mutex would work, but it would couple independent numeric telemetry to the store's event-state lock.

#### Design Discussion

This lesson uses typed `sync/atomic` counters inside `machines.Service`. Each counter is independent and monotonically increasing, which makes atomics a better fit than a broader mutex-protected stats struct.

The counters track received attempts, accepted new events, duplicate retries, conflicts and invalid events. A small protected stats endpoint returns a point-in-time snapshot.

This is not the full observability milestone. There are no Prometheus metrics, logs or traces yet; those belong to Milestone 13.

#### Go Concepts

- typed atomic counters with `atomic.Uint64`
- atomic `Add` for concurrent increments
- atomic `Load` for snapshot reads
- choosing atomics only for independent scalar state
- race detector review for mixed mutex and atomic synchronization

#### Architecture Concepts

- operational counters as service-owned state
- HTTP stats endpoint as temporary visibility before formal observability
- atomics for telemetry, mutexes for compound event state
- avoiding global metrics registries before Milestone 13

### Lesson 12.4 Completion Notes

#### Business Context

MES Lite now exposes basic fake machine intake counters so concurrent event activity can be inspected without reading internal state directly.

#### Problem

Machine intake outcomes were invisible. Tests could prove behavior, but an operator or client could not see how many events were accepted, retried, conflicted or rejected as invalid.

#### Design Discussion

Added counters to `machines.Service` instead of to the in-memory store. The service owns the workflow decisions that classify an attempt as accepted, duplicate retry, conflict or invalid.

The event store still uses mutexes because it protects compound state: a slice and a map that must stay consistent. The service counters use atomics because each number can be updated independently without a larger invariant between fields.

The stats endpoint is protected with the same admin/manager RBAC as fake machine intake. Workers should not inspect integration operational counters.

#### Implementation

- Added `IntakeStats` snapshot type.
- Added typed atomic counters to `machines.Service`.
- Counted received attempts, accepted events, duplicate retries, conflicts and invalid events.
- Added `Service.Stats` using atomic loads.
- Added `IntakeStatsResponse` DTO.
- Added `GET /machines/events/stats`.
- Registered stats route with admin/manager RBAC.

#### Tests

- Added service tests for accepted, duplicate retry, conflict and invalid counters.
- Extended concurrent retry test to verify atomic counters after 50 concurrent callers.
- Added handler test for stats response mapping.
- Added server route tests proving managers can read machine stats and workers are forbidden.
- Verified with `go test ./internal/machines ./internal/server -count=1`.
- Verified with `go test ./internal/machines -race -count=5 -shuffle=on`.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `go vet ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

No global metrics package was introduced. Counters remain local to the machine service until Milestone 13 introduces formal observability.

#### Code Review

An experienced Go engineer would approve atomics here because the counters are independent scalar values. They would not be appropriate for the store's event slice and deduplication map, where multiple fields must change together.

The main caveat is that these counters are in-memory and reset on restart. That is acceptable for the learning goal and current fake integration scope.

#### Exercises

- Add a test proving stats reset when a new service is constructed.
- Add a benchmark comparing atomic counter increments with mutex-protected increments.
- Decide which machine counters should become Prometheus metrics in Milestone 13.

#### Interview Questions

- When are atomic operations preferable to a mutex?
- Why are atomics risky for compound state?
- What is the difference between a race-free counter and a consistent multi-field snapshot?
- Why does a clean race detector run not prove a design is logically correct?

#### Roadmap Update

- Lesson 12.4 completed.
- Current lesson moved to Lesson 12.5.
- Standard Library `sync/atomic` and Concurrency `Atomic` marked complete in the Knowledge Matrix.

### Lesson 12.5 Scope

Review the machine integration slice for package boundaries, synchronization choices, idempotency behavior, race safety and milestone readiness.

#### Business Context

The fake machine API is not production machine connectivity, but it establishes the learning foundation for machine-originated events. Before moving to observability, the slice should be reviewed as a complete learning milestone.

#### Problem

Milestone 12 introduced events, deduplication, synchronization and atomic counters. The remaining question was whether the exported in-memory store type should continue supporting zero-value construction or whether external consumers should be forced through the constructor.

#### Design Discussion

The constructor-boundary refactor was accepted. `InMemoryStore` was renamed to an unexported `inMemoryStore`, and `NewInMemoryStore` now returns the `Store` interface. This is a deliberate exception to the usual Go preference for returning concrete types: the goal is to hide fake infrastructure details and force consumers to receive a correctly initialized store.

Because construction is now enforced outside the package, the previous `sync.Once` initialization is no longer needed. A constructor-initialized map is simpler. The synchronization lesson remains valid historically, but the final slice favors the smaller production shape.

#### Go Concepts

- unexported concrete implementation types
- constructor-enforced initialization
- when returning an interface from a constructor can be justified
- reviewing `sync.Once` after requirements change
- milestone-level race-detector verification

#### Architecture Concepts

- package API as an enforcement boundary
- fake infrastructure hidden behind a small interface
- milestone closure through review-driven refactoring
- documented trade-off between zero-value usability and constructor enforcement

### Lesson 12.5 Completion Notes

#### Business Context

Milestone 12 is complete. MES Lite now has a fake machine-event intake path that demonstrates event modeling, idempotent event intake, synchronization and atomic counters.

#### Problem

The machine slice was functionally complete for the milestone, but the final review identified one API-boundary issue: an exported `InMemoryStore` allowed external consumers to bypass the constructor even though the store had initialization requirements.

#### Design Discussion

The final design hides the concrete in-memory store and exposes it through `NewInMemoryStore`. This answers the constructor question directly: Go cannot force constructor use for an exported struct, so the struct must be unexported if constructor use matters.

This does mean giving up zero-value usability for that concrete store outside the package. That is acceptable here because this is fake infrastructure, not a domain value type. The public API still gives callers a ready-to-use store.

The milestone also reviewed atomics versus mutexes. The store uses a mutex because it protects compound state. The service uses atomics because intake counters are independent scalar values.

#### Implementation

- Renamed exported `InMemoryStore` to unexported `inMemoryStore`.
- Changed `NewInMemoryStore` to return the `Store` interface.
- Initialized the deduplication map in the constructor.
- Removed `sync.Once` from the final in-memory store implementation.
- Kept mutex-protected save, list and lookup behavior unchanged.
- Preserved all fake machine API routes and response contracts.

#### Tests

- Updated concurrent store test to use constructor-created stores.
- Re-ran focused machine and server tests.
- Verified race safety with `go test ./internal/machines -race -count=5 -shuffle=on`.
- Verified full project behavior with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `go vet ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

The refactor removes synchronization that was only needed for zero-value initialization. This is the right simplification after choosing constructor enforcement.

No durable machine-event persistence was introduced. The milestone intentionally remains about synchronization and event intake, not database-backed integration reliability.

#### Code Review

An experienced Go engineer would approve the milestone as a learning slice. The code models events separately from production entries, handles duplicate delivery idempotently, protects shared state with a mutex, uses atomics for independent counters and has race-detector coverage for concurrent intake paths.

The main production gaps are explicit: events are in-memory, machine authentication is fake, counters reset on restart and no event-to-production-entry processing exists yet.

#### Exercises

- Try exporting `InMemoryStore` again and explain why constructor enforcement is lost.
- Design a PostgreSQL-backed machine event store with a unique `(machine_id, external_event_id)` index.
- Decide whether machine stats should be exposed through this endpoint after Prometheus metrics exist.

#### Interview Questions

- How can Go package visibility enforce constructor use?
- Why might returning an interface from a constructor be acceptable for hidden infrastructure?
- Why are atomics appropriate for counters but not for updating a map and slice together?
- What would make this fake machine integration production-ready?

#### Roadmap Update

- Lesson 12.5 completed.
- Milestone 12 completed.
- Current milestone moved to Milestone 13.
- Current lesson moved to Lesson 13.1.
- Architecture maturity, Go knowledge progress and interview readiness updated.

### Milestone 12 Review

#### Architecture Review

An experienced Go engineer would approve Milestone 12 as a focused learning milestone. Machine events live in their own vertical slice, are not confused with manual production entries and are accepted through a fake adapter endpoint that can evolve toward real integration later.

The architecture remains intentionally non-durable. That is acceptable because the milestone goal was synchronization and idempotent event intake, not production-grade machine connectivity.

#### Code Review

The code is explicit and small. The handler owns HTTP translation, the service owns idempotent workflow semantics, the store owns synchronized state and counters are service-local atomics.

The constructor-boundary refactor improves the public package API: external consumers can no longer construct an uninitialized in-memory store directly.

#### Refactoring

The final refactor hides fake infrastructure behind an unexported concrete type and removes `sync.Once` from the store because constructor initialization is now enforced.

#### Interview Review

You should now be able to discuss event versus command semantics, idempotent event processing, duplicate detection, mutex versus atomic trade-offs, `sync.Once`, race-detector limits and Go package visibility as an API boundary.

#### Completion Criteria

- Fake machine API implemented.
- Machine event model implemented.
- Duplicate detection implemented.
- Idempotent retry behavior implemented.
- Conflicting duplicate payloads return `409 Conflict`.
- Shared in-memory event state is mutex-protected.
- Intake counters use `sync/atomic`.
- Race detector passes for machine intake tests.
- Constructor use is enforced for external in-memory store consumers.
- Tests, build, vet and lint pass.
- Roadmap updated.

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

✅ Completed

### Lessons

- **L13.1** — Structured Logging Foundation ✅
- **L13.2** — Request Logging & Correlation IDs ✅
- **L13.3** — Prometheus Metrics ✅
- **L13.4** — OpenTelemetry Tracing ✅
- **L13.5** — Health, Readiness & Observability Review ✅

### Lesson 13.1 Scope

Introduce a small structured logging foundation with `log/slog` and centralize logger setup before adding request logging, metrics or tracing.

#### Business Context

Operators need machine-readable logs when diagnosing startup failures, migration failures and server lifecycle events. Plain text messages are hard to filter and correlate once the service runs in containers or log aggregation systems.

#### Problem

The application already used `slog` directly in command entry points and server lifecycle code, but logger setup was duplicated and not configurable. The server package also depended on the global default logger for lifecycle messages, which made tests and future request-scoped logging less explicit.

#### Design Discussion

Logging remains a platform concern, not a business slice. `internal/platform/logging` now constructs a configured `*slog.Logger` using standard-library handlers. The application supports JSON logs by default because structured JSON is the common production shape for containerized services. Text logs remain available for local debugging.

The lesson intentionally does not add request logging, correlation IDs, Prometheus metrics or OpenTelemetry. Those belong to later observability lessons so each concept remains clear.

#### Go Concepts

- `log/slog` structured logging
- `slog.Logger` as an explicit dependency
- `io.Writer` for testable log output
- configuration parsing for log level and format
- avoiding package-level globals where explicit injection is cheap

#### Architecture Concepts

- platform logging package below business slices
- command entry points configure process-wide defaults
- server lifecycle logging receives an explicit logger
- observability foundation before request correlation and metrics

### Lesson 13.1 Completion Notes

#### Business Context

MES Lite now has a configurable structured logging foundation for startup, migration and server lifecycle diagnostics.

#### Problem

Logger setup was duplicated in `cmd/server` and `cmd/migrate`, and the server package logged through the global `slog` default instead of an explicit dependency.

#### Design Discussion

Added `internal/platform/logging` with one `New` function that builds a `*slog.Logger` from an `io.Writer`, log level and format. JSON is the default because it is easy for production log collectors to parse. Text format remains available through `LOG_FORMAT=text` for local development.

`cmd/server` and `cmd/migrate` now both configure the logger through the shared package and set it as the process default. `internal/server` receives an injected logger for lifecycle logs and uses a discard logger by default so tests stay quiet.

#### Implementation

- Added `internal/platform/logging`.
- Added configurable JSON and text `slog` handlers.
- Added log-level parsing for `debug`, `info`, `warn`, `warning` and `error`.
- Added `LOG_LEVEL` and `LOG_FORMAT` configuration.
- Defaulted logs to `LOG_LEVEL=info` and `LOG_FORMAT=json`.
- Updated `cmd/server` and `cmd/migrate` to use the shared logging package.
- Added `Server.SetLogger` and changed server lifecycle logs to use the injected logger.

#### Tests

- Added logging tests for JSON structured output.
- Added logging tests for text output.
- Added logging tests for level parsing and invalid configuration.
- Updated configuration tests for logging defaults and environment values.
- Verified with `go test ./internal/platform/config ./internal/platform/logging ./internal/server -count=1`.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `go vet ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

Duplicated logger construction was removed from command entry points. Server lifecycle logging no longer relies directly on package-global logger state.

#### Code Review

An experienced Go engineer would approve this as a focused first observability step. It uses the standard library, keeps logging as platform plumbing and avoids introducing an external logging framework.

The main follow-up is request logging with correlation IDs. Startup logs are structured now, but request-level diagnostics still need a middleware boundary.

#### Exercises

- Run the server with `LOG_FORMAT=text` and compare the output with JSON logs.
- Add a test proving `LOG_LEVEL=error` suppresses info logs.
- Decide which startup values are safe to log and which could leak secrets.

#### Interview Questions

- Why are structured logs more useful than plain formatted strings in production?
- Why is `log/slog` preferable here to adding a third-party logger?
- Why should secrets such as `JWT_SECRET` never be logged?
- When is using `slog.Default()` acceptable, and when is explicit logger injection better?

#### Roadmap Update

- Lesson 13.1 completed.
- Current lesson moved to Lesson 13.2.
- Standard Library `log/slog` marked complete in the Knowledge Matrix.

### Lesson 13.2 Scope

Add request logging middleware with correlation IDs so every HTTP request receives a stable identifier and emits one structured completion log.

#### Business Context

Operators need to diagnose API behavior across many concurrent requests. A startup log can prove the process booted, but request logs are what connect a user's API call to status code, route path and latency.

#### Problem

The application had structured startup and lifecycle logs, but individual HTTP requests were not logged. There was also no request-scoped correlation ID that could be returned to clients and propagated through future logs, metrics or traces.

#### Design Discussion

Request logging is implemented as standard `net/http` middleware in `internal/platform/logging`. It preserves an incoming `X-Request-ID` header when present, generates a UUID-shaped ID when missing, stores it in request context and returns it on the response header.

The middleware logs one completion record after the handler runs. It records method, path, status and duration in milliseconds. Query strings are intentionally not logged because they can contain sensitive values and usually create high-cardinality logs.

Fuego route registration uses `fuego.Use` rather than only global server middleware. This keeps request logging active in both runtime and existing tests that call `s.Mux.ServeHTTP` directly.

#### Go Concepts

- standard-library HTTP middleware
- response-writer wrapping for status capture
- context values for request-scoped correlation IDs
- `time.Since` for request duration
- preserving inbound headers versus generating defaults

#### Architecture Concepts

- request logging as platform middleware
- correlation ID as an observability boundary
- explicit server construction with logger injection
- avoiding metrics and traces until their dedicated lessons

### Lesson 13.2 Completion Notes

#### Business Context

MES Lite now emits one structured log record per HTTP request and returns a correlation ID to API clients.

#### Problem

Without request logs, operators could see startup failures but not normal API traffic, response statuses or request latency. Without a correlation ID, future diagnostics could not tie logs, client reports and traces together.

#### Design Discussion

Added `logging.RequestLogger` as reusable platform middleware. It accepts a `*slog.Logger`, wraps the `http.ResponseWriter` to capture status code and stores a request ID in context through an unexported key type.

The request ID comes from `X-Request-ID` when supplied by a client or gateway. If it is missing, the middleware generates one with the existing platform ID generator. The same ID is returned in the response header.

`server.NewWithLogger` registers request logging before routes are added. Existing tests can still call `server.New`, which uses a discard logger by default.

#### Implementation

- Added `logging.RequestLogger` middleware.
- Added request ID context helpers in `internal/platform/logging`.
- Added `X-Request-ID` response header propagation.
- Added status-code capture through a response-writer wrapper.
- Logged structured request fields: `request_id`, `method`, `path`, `status` and `duration_ms`.
- Added `server.NewWithLogger` for runtime logger injection before route registration.
- Updated `cmd/server` to build the HTTP server with the configured logger.

#### Tests

- Added middleware test preserving an incoming request ID.
- Added middleware test generating a request ID when missing.
- Added context propagation assertions for request IDs.
- Added structured log assertions for method, path, status and request ID.
- Added server route test proving Fuego-registered routes emit request logs and return `X-Request-ID`.
- Verified with `go test ./internal/platform/logging ./internal/server -count=1`.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `go vet ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

The previous `Server.SetLogger` hook was replaced by `NewWithLogger` because request middleware must be registered with the correct logger before routes are installed. The old `New` constructor remains as a quiet default for tests.

#### Code Review

An experienced Go engineer would approve the middleware shape because it uses `net/http`, avoids framework-specific logging code and keeps correlation IDs request-scoped rather than global.

The main caveat is that only status, path and duration are logged today. User identity, route templates, panic recovery and trace IDs are intentionally postponed until later observability work defines those boundaries.

#### Exercises

- Add a test proving missing routes also receive `X-Request-ID` once the server is run through the same handler path used in production.
- Add authenticated user ID to request logs and discuss the privacy/security trade-off.
- Decide whether `duration_ms` should be integer milliseconds or a string duration.

#### Interview Questions

- Why are correlation IDs useful in distributed systems?
- Why should request IDs be stored in context instead of a global variable?
- How does wrapping `http.ResponseWriter` let middleware observe status codes?
- Why might logging query strings be risky?

#### Roadmap Update

- Lesson 13.2 completed.
- Current lesson moved to Lesson 13.3.

### Lesson 13.3 Scope

Expose Prometheus-compatible HTTP metrics and record low-cardinality request counters and durations.

#### Business Context

Operators need numeric signals that can be graphed and alerted on. Logs answer what happened for individual requests, but metrics answer aggregate questions such as how many requests are served and whether latency is increasing.

#### Problem

The application had structured request logs and correlation IDs, but no metrics endpoint. Prometheus could not scrape request volume or latency, and there was no foundation for future business or runtime metrics.

#### Design Discussion

Added `internal/platform/metrics` because metrics are technical observability plumbing, not a business slice. The package owns a Prometheus registry, HTTP request counter and request-duration histogram.

The first labels are intentionally low-cardinality: `method` and `status`. Route/path labels are postponed because raw paths can include IDs and create unbounded time series. If route-template labels are added later, they should use stable templates such as `/production-orders/{id}`, not raw URL paths.

`/metrics` is public and hidden from OpenAPI. In a real deployment, network policy or ingress rules would usually restrict it to the monitoring system. The endpoint is registered before request middleware, so metrics scrapes do not count as application traffic.

#### Go Concepts

- middleware composition around `http.Handler`
- response-writer wrapping reused for status capture
- histogram versus counter metric types
- dependency introduction when the standard library lacks the protocol tooling
- low-cardinality label design

#### Architecture Concepts

- observability platform package
- Prometheus pull model through `/metrics`
- metrics middleware separated from business handlers
- avoiding high-cardinality labels at the platform boundary

### Lesson 13.3 Completion Notes

#### Business Context

MES Lite now exposes Prometheus-compatible metrics for HTTP request volume and latency.

#### Problem

Operators could read request logs, but they could not scrape aggregate request metrics for dashboards or alerts.

#### Design Discussion

Introduced Prometheus through `github.com/prometheus/client_golang`. This dependency is justified because Prometheus exposition, registries and metric types are ecosystem standards and are not provided by the Go standard library.

The implementation uses a package-local registry instead of the global default registry. This keeps tests isolated and prevents duplicate metric registration panics when multiple servers are constructed in the same process.

HTTP metrics use only `method` and `status` labels. This avoids the common production mistake of labeling by raw path, where IDs such as `/production-orders/123` and `/production-orders/456` create separate time series.

#### Implementation

- Added `internal/platform/metrics`.
- Added `HTTPMetrics` with a Prometheus registry.
- Added `mes_lite_http_requests_total` counter.
- Added `mes_lite_http_request_duration_seconds` histogram.
- Added metrics middleware for request count and duration.
- Added `/metrics` endpoint using `promhttp.HandlerFor`.
- Registered `/metrics` before application middleware so scrapes do not count themselves.
- Hid `/metrics` from generated OpenAPI documentation.
- Added Prometheus client dependency.

#### Tests

- Added metrics middleware test for explicit status-code recording.
- Added metrics middleware test for default `200 OK` recording when handlers only write a body.
- Added server test proving `/metrics` returns request metrics after a `/health` request.
- Added server test coverage that `/metrics` scrape does not increment the application request counter.
- Verified with `go test ./internal/platform/metrics ./internal/server -count=1`.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `go vet ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

No business handlers changed. Metrics were added as platform middleware, keeping observability concerns out of the vertical slices.

#### Code Review

An experienced Go engineer would approve the small metrics foundation and isolated registry. The main positive review point is label restraint: method/status metrics are less detailed than route metrics, but they avoid unsafe cardinality until stable route-template labels are designed.

The main follow-up is tracing. Metrics can show that latency increased; traces will help explain where a slow request spent time.

#### Exercises

- Add a test proving a `404 Not Found` response increments the correct status label if routed through the same middleware path.
- Design safe route-template labels for Fuego routes without using raw URL paths.
- Decide whether `/metrics` should remain public or require network-level protection in deployment.

#### Interview Questions

- What is the difference between logs and metrics?
- Why are high-cardinality Prometheus labels dangerous?
- Why use a histogram for request duration instead of only logging duration?
- Why can a package-local Prometheus registry be better than the global default registry in tests?

#### Roadmap Update

- Lesson 13.3 completed.
- Current lesson moved to Lesson 13.4.

### Lesson 13.4 Scope

Add OpenTelemetry HTTP tracing so each request creates a server span and request logs can include a trace ID when tracing is active.

#### Business Context

Metrics can show that requests are slow or failing, but operators need traces to understand where a request spent time as the service grows. A trace ID also gives support and developers a common handle that can connect logs, client reports and distributed traces.

#### Problem

The application had structured logs, correlation IDs and Prometheus metrics, but no tracing foundation. Request logs could identify one request, but there was no OpenTelemetry span context that future database, job or integration spans could attach to.

#### Design Discussion

Tracing is added as platform middleware. Each HTTP request starts a server span with method, path, status and duration attributes. The middleware extracts W3C `traceparent` headers so upstream trace context can flow into MES Lite.

The provider supports `OTEL_TRACES_EXPORTER=none` by default and `OTEL_TRACES_EXPORTER=stdout` for local inspection. OTLP export is intentionally postponed. Adding a collector endpoint and deployment configuration is a production-integration concern, while this lesson focuses on request span creation and context propagation.

The middleware order is metrics, tracing, then request logging. This means request logs are emitted while the span context is still active, so logs can include `trace_id`.

#### Go Concepts

- OpenTelemetry tracer providers and spans
- context propagation through HTTP middleware
- W3C trace context extraction
- span attributes and status
- graceful tracer-provider shutdown

#### Architecture Concepts

- tracing as platform infrastructure
- request spans as the root for future child spans
- trace IDs connected to structured logs
- exporter choice kept in configuration

### Lesson 13.4 Completion Notes

#### Business Context

MES Lite now creates OpenTelemetry server spans for HTTP requests and can connect request logs to traces through `trace_id`.

#### Problem

Logs and metrics were useful but incomplete. There was no trace context for following one request through future internal operations.

#### Design Discussion

Added `internal/platform/tracing` with a configurable tracer provider and HTTP middleware. The middleware starts one server span per request and records low-risk request attributes. It avoids business-specific spans for now because database, job and integration tracing should be added where those boundaries are reviewed.

`cmd/server` configures the provider during startup and shuts it down with a timeout on exit. The default exporter is `none`, so local development does not emit trace payloads unexpectedly. `stdout` is available for learning and manual inspection.

Request logging now checks the active span context and includes `trace_id` when one exists.

#### Implementation

- Added `internal/platform/tracing`.
- Added `tracing.Config` and `NewProvider`.
- Added `OTEL_TRACES_EXPORTER` configuration with default `none`.
- Added `stdout` trace exporter support for local inspection.
- Installed W3C trace-context propagation.
- Added tracing HTTP middleware.
- Added span attributes for method, path, status and duration.
- Added trace IDs to request logs when an active span exists.
- Wired tracing startup and shutdown in `cmd/server`.

#### Tests

- Added invalid exporter configuration test.
- Added tracing middleware test using OpenTelemetry's span recorder.
- Added assertions for span name, method, path and response status attributes.
- Added test for missing trace ID when no span exists.
- Added request-log test proving `trace_id` is included when span context is active.
- Updated configuration tests for `OTEL_TRACES_EXPORTER`.
- Verified with `go test ./internal/platform/tracing ./internal/platform/logging ./internal/platform/config ./internal/server -count=1`.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `go vet ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

No business package was changed. Tracing was added as platform middleware and request logging gained only trace-context awareness.

#### Code Review

An experienced Go engineer would approve the foundation because tracing context is established at the HTTP boundary and exporter behavior is explicit. The main limitation is deliberate: OTLP export and child spans for database/job operations are not wired yet.

#### Exercises

- Send a request with an existing `traceparent` header and verify the server span joins the upstream trace.
- Enable `OTEL_TRACES_EXPORTER=stdout` and inspect the emitted span fields.
- Design where database query spans should be introduced without leaking sqlc details into handlers.

#### Interview Questions

- What does a trace show that logs and metrics do not?
- Why is context propagation required for distributed tracing?
- What is the difference between a trace ID and a span ID?
- Why should tracer-provider shutdown be part of graceful shutdown?

#### Roadmap Update

- Lesson 13.4 completed.
- Current lesson moved to Lesson 13.5.
- Known technical debt updated for missing OTLP exporter wiring.

### Lesson 13.5 Scope

Review health, readiness and the complete observability stack before closing Milestone 13.

#### Business Context

Operators need two different operational signals: whether the process is alive and whether it is ready to serve real business traffic. They also need confidence that logs, metrics and traces compose without leaking into business packages.

#### Problem

`/health` existed as a liveness endpoint, but `/ready` always returned success. That meant an orchestrator could send traffic to an application whose database dependency was unavailable.

#### Design Discussion

Liveness stays cheap and dependency-free: `/health` answers only whether the process can respond. Readiness now has an explicit check hook. The production composition root wires that hook to `pgxpool.Ping`, while tests can install a small fake function.

A readiness failure returns `503 Service Unavailable`, not `500 Internal Server Error`. The process is still alive; it is just not ready to receive traffic.

#### Go Concepts

- function fields for small dependency hooks
- `context.WithTimeout` for bounded readiness checks
- correct HTTP status semantics for operational endpoints
- focused tests for success and dependency failure

#### Architecture Concepts

- liveness separated from readiness
- composition root wires infrastructure health checks
- observability remains in platform packages
- milestone review before performance work begins

### Lesson 13.5 Completion Notes

#### Business Context

MES Lite now has a meaningful readiness endpoint in addition to the existing liveness endpoint.

#### Problem

The readiness endpoint returned success unconditionally. That was misleading because the server could be alive while PostgreSQL was unreachable.

#### Design Discussion

Added a readiness-check hook to `server.Server`. The server package does not import PostgreSQL or know which dependencies matter; it only calls the injected function with a bounded context. `cmd/server` wires the hook to `db.Ping`.

This keeps dependency knowledge in the composition root and preserves the server package as HTTP composition rather than database infrastructure.

#### Implementation

- Added `Server.SetReadinessCheck`.
- Changed `/ready` to return `{"status":"ready"}` when the check succeeds.
- Added a two-second timeout around readiness checks.
- Changed readiness dependency failures to return `503 Service Unavailable`.
- Logged readiness failures as structured warnings.
- Wired runtime readiness to `pgxpool.Pool.Ping` in `cmd/server`.

#### Tests

- Added server test proving `/ready` calls the readiness check and returns `200 OK` with `status=ready`.
- Added server test proving readiness-check failure returns `503 Service Unavailable`.
- Verified with `go test ./internal/server -count=1`.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `go vet ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

`NewWithLogger` now constructs `*Server` before route registration so the `/ready` route can call a method that uses server state. The public route behavior is otherwise unchanged.

#### Code Review

An experienced Go engineer would approve the liveness/readiness split. The readiness endpoint checks the database without making `internal/server` depend on pgx, and failures use the correct operational status code.

The main observability limitation remains OTLP exporter wiring. The application can create request spans and emit stdout traces for learning, but production collector export is still future work.

#### Exercises

- Add a readiness check that verifies the background worker pool has started.
- Add a test proving readiness failures are logged with `request_id` and `trace_id` when tracing is active.
- Decide whether `/health`, `/ready` and `/metrics` should be protected by network policy in deployment.

#### Interview Questions

- What is the difference between liveness and readiness?
- Why should readiness checks have timeouts?
- Why is `503 Service Unavailable` better than `500` for dependency-readiness failure?
- Why should observability middleware live outside business slices?

#### Roadmap Update

- Lesson 13.5 completed.
- Milestone 13 completed.
- Current milestone moved to Milestone 14.
- Current lesson moved to Lesson 14.1.
- Architecture `Observability` marked complete in the Knowledge Matrix.

### Milestone 13 Review

#### Architecture Review

An experienced Go engineer would approve Milestone 13 as a coherent observability foundation. Logging, metrics and tracing live in platform packages. Business slices did not gain observability dependencies, and HTTP middleware composes through standard `net/http` shapes.

The milestone intentionally keeps production export simple. Prometheus scraping is available through `/metrics`, while tracing supports no-op and stdout exporters. OTLP collector wiring remains known technical debt for a later production-readiness pass.

#### Code Review

The code is explicit and testable. Logger setup is centralized. Request logging propagates correlation IDs. Metrics use a local registry and low-cardinality labels. Tracing creates HTTP server spans and connects request logs to trace IDs. Readiness now checks PostgreSQL through an injected function.

The main improvement for later is route-template observability. Raw path labels were intentionally avoided for metrics, but future stable route-template labels would improve dashboard usefulness without creating high-cardinality series.

#### Refactoring

This milestone removed duplicated logger setup, replaced a post-construction logger setter with `NewWithLogger`, added platform packages for metrics and tracing, and changed readiness from a static route to a dependency-aware method.

#### Interview Review

You should now be able to explain structured logs, correlation IDs, Prometheus counters and histograms, label cardinality, OpenTelemetry spans, trace context propagation, liveness versus readiness and why observability belongs at platform boundaries.

#### Completion Criteria

- Structured logging implemented with `log/slog`.
- Request logging emits correlation IDs and status/latency fields.
- Prometheus-compatible `/metrics` endpoint implemented.
- HTTP request metrics avoid raw-path labels.
- OpenTelemetry request tracing implemented.
- Request logs include `trace_id` when a span is active.
- `/health` remains a cheap liveness endpoint.
- `/ready` checks PostgreSQL readiness and returns `503` when unavailable.
- Tests, build, vet and lint pass.
- Roadmap updated.

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

🔄 In Progress

### Lessons

- **L14.1** — Benchmarking Foundations
- **L14.2** — pprof CPU & Memory Profiling
- **L14.3** — Allocation Analysis & Escape Analysis
- **L14.4** — Runtime Scheduler & Garbage Collector Review
- **L14.5** — Performance Review & Optimization Discipline

### Lesson 14.1 Scope

Introduce Go benchmarks on a real code path and establish the rule that optimization starts with measurement.

#### Business Context

CSV import and reporting are likely performance-sensitive paths as data grows. Before changing implementations, the project needs repeatable measurements that can show whether a change helps or hurts.

#### Problem

The codebase has many correctness tests but no benchmarks. Without benchmarks, performance discussions rely on guesses, and future refactors could accidentally make hot paths slower or more allocation-heavy.

#### Design Discussion

Start with a focused benchmark around CSV import validation or request middleware overhead. The benchmark should measure an existing behavior without changing production code first. The goal of L14.1 is learning benchmark mechanics and avoiding misleading results, not immediately optimizing.

Benchmarks should use `testing.B`, report allocations with `b.ReportAllocs()` and keep setup outside the measured loop where possible.

#### Go Concepts

- `testing.B`
- benchmark loops and `b.N`
- `b.ReportAllocs()`
- setup versus measured work
- avoiding benchmark dead-code elimination

#### Architecture Concepts

- performance baselines before optimization
- realistic input sizes for business-relevant paths
- benchmark results as engineering evidence

#### Tests

- Add at least one benchmark for a meaningful existing path.
- Run the benchmark with allocation reporting.
- Keep all correctness tests passing.

#### Exercises

- Compare benchmark results for small and large CSV samples.
- Move setup accidentally inside the benchmark loop and explain why the result changes.
- Decide which endpoint or package should receive the next benchmark.

#### Interview Questions

- How does Go's benchmark runner choose `b.N`?
- Why should setup usually be outside the measured loop?
- What does `allocs/op` tell you?
- Why can microbenchmarks mislead production optimization?

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
