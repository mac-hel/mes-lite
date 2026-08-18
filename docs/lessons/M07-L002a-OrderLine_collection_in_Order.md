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
