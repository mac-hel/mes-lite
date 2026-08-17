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
