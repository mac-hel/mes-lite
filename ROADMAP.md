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

Proceed with the project.

## STATE

This section must always reflect the current progress.

**Version:** 1.0
**Status:** IN PROGRESS
**Current milestone:** 0 - Development Environment
**Current lesson:** L1 - Project Foundation
**Completed milestones:**
- none
**Next milestone:** 1 - Bootstrap HTTP Service
**Current branch:** main
**Architecture maturity:** 0 / 10
**Go knowledge progress:** 5%
**Interview readiness:** 0%
**Known technical debt:** None

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
- [ ] Visibility
- [ ] Variables
- [ ] Constants
- [ ] Functions
- [ ] Structs
- [ ] Methods
- [ ] Receivers
- [ ] Pointers
- [ ] Interfaces
- [ ] Embedding
- [ ] Custom Types
- [ ] iota
- [ ] Errors
- [ ] Error Wrapping
- [ ] defer
- [ ] panic
- [ ] Context
- [ ] Generics (intro)
- [ ] Reflection (overview)

**Standard Library**
- [ ] net/http
- [ ] context
- [ ] errors
- [ ] io
- [ ] encoding/json
- [ ] encoding/csv
- [ ] strings
- [ ] bytes
- [ ] time
- [ ] sync
- [ ] sync/atomic
- [ ] database/sql (concepts)
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
- [ ] Vertical Slice
- [ ] Dependency Injection
- [ ] Package Design
- [ ] Repositories
- [ ] Aggregates
- [ ] CQRS
- [ ] Event-driven Architecture
- [ ] Observability
- [ ] Production Readiness

**Persistence**
- [ ] PostgreSQL
- [ ] pgx
- [ ] sqlc
- [ ] Transactions
- [ ] Optimistic Locking
- [ ] SQL Optimization

**Testing**
- [ ] Unit Tests
- [ ] Integration Tests
- [ ] Table Tests
- [ ] httptest
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
- [ ] Server

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

- idiomatic Go
- production architecture
- concurrency
- testing
- observability
- software design
- production engineering
- reviewing Go code
- deploying production services
- contributing to mature Go codebases

Every design decision should optimize for learning.

---

# 2. Learning Goals

After completing this roadmap I should be able to:

## Go

- write idiomatic Go
- understand Go philosophy
- understand why Go is different from Java/PHP/C#
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

⬜ Not Started

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

⬜ Not Started

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

⬜ Not Started

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

⬜ Not Started

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
