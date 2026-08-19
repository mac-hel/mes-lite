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
