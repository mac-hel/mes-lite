### Lesson 7.3 Scope

Expose production-order creation and read endpoints over HTTP without adding status-transition mutation endpoints yet.

#### Business Context

Production managers need an API to create planned production orders, and team leaders need to read those orders before production starts.

#### Problem

Orders were persisted in PostgreSQL, but no HTTP route allowed clients to create or fetch them. The aggregate also has private fields, so returning the domain type directly would produce the wrong API contract.

#### Design Discussion

The orders handler uses explicit request and response DTOs. This keeps HTTP JSON shape separate from the aggregate internals and avoids exposing private-field domain types directly.

Create accepts an order ID, one or more order lines and optional assigned employee IDs. It always creates a draft order. Status-transition endpoints are postponed to Lesson 7.4 so this lesson stays focused on HTTP create/read wiring.

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

Added `orders.Handler` with explicit create/read DTOs. The handler constructs `OrderLine`, `OrderLines` and `Order`, assigns optional employees through aggregate methods and delegates persistence to `orders.Store`.

Authorization follows the existing RBAC model: admins and managers can create orders; admins, managers and leaders can read orders. Workers cannot create or read planning data through these routes.

#### Go Concepts

- DTO-to-domain conversion
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
