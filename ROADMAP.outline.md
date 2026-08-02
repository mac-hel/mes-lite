## Milestone 0 - Development Environment
- L1 — Project Foundation (Git, structure, Go module, go tooling, Makefile, golangci-lint, air, pre-commit hooks, GitHub Actions)
- L2 — Docker Environment & First ADR (docker compose)

## Milestone 1 - Bootstrap HTTP Service
(running HTTP service: helath,version, OpenAPI ep, app deployable)

## Milestone 2 - Employees
(employee management: create, update, deactivate, list)
- **L2.1** — Visibility & Zero Values: Employee Entity
- **L2.2** — Constructors & Slices: Creating & Listing Employees
- **L2.3** — Error Wrapping, Validation & Maps: Updating, Deactivating & Testing

## Milestone 3 - Products
(manage manufactured products, registered by employee: create, update, deactivate, search)
- **L3.1** — Custom Types & iota: Product Entity ✅
- **L3.2** — Stringer & strings: Product Handlers & Search ✅
- **L3.3** — Value Objects & Testing: Product Validation & Review

## Milestone 4 - Production Registration
(register production: employee+product+quantity, workstation, timestamp, comment)

## Milestone 5 - Persistence & Data Access
(production-quality database practices, persistence layer: repository, pagination, filtering, sorting, soft del, optimistic locking)

## Milestone 6 - Authentication & Authorization
(secure app: login/JWT authen, user permissions/role-based author, protected endpoints)

## Milestone 7 - Production Orders
(production planning before execution: create order, assign products and quantity, assign employees, order status)

## Milestone 8 - Reporting
(business reports - monitor production: daily production, employee productivity, product statistics + filtering, exports)

## Milestone 9 - CSV Import
(import historical Excel data: CSV upload, validation, batch import, import summary)

## Milestone Review (Milestones 0-9)
app usable by small company

## Milestone 10 - Background Jobs & Concurrency
(to not block Request - async processing: CSV imports, report generation, notifications, scheduled maintenance)
with: job queue, worker pool, progress tracking, cancellation, retry strategy

## Milestone 11 - Machine Integration & Synchronization
(concurrent events from production machines, e.g. CNC: fake machine API <- events, duplicate detection, idempotency)

## Milestone 12 - Observability
(health/readiness, issues on prod: logging, metrics, tracing)

## Milestone 13 - Performance Engineering
(understand and measure go app perfo, optimize: profiling, benchmarking, allocation analysis)

## Milestone 14 - Production Readiness
(prepare for deployment with production quality: docker images, CI/CD, graceful deploy, versioning, release process)
