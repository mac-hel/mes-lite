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
