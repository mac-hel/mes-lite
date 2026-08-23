# Go concepts: mental models and deeper explanations

The common theme is that Go tries to make **ownership, dependencies, control flow, and boundaries visible in the code**. It has relatively few language mechanisms, but those mechanisms interact in important ways.




---


---

# 8. Blank imports


---

# 9. Domain modeling

These next four concepts come largely from Domain-Driven Design rather than Go itself, but Go's type and package system works well with them.

---

## Value Object

### Mental model

A Value Object is defined by **what it contains**, not by an identity.

Examples:

```text
Money(100, PLN)
Email("alice@example.com")
DateRange(...)
Address(...)
```

Two Money objects with identical amount and currency represent the same value.

### Deeper explanation

A Go value object might look like:

```go
type Money struct {
    amount   int64
    currency Currency
}

func NewMoney(amount int64, currency Currency) (Money, error) {
    // validate
    return Money{
        amount:   amount,
        currency: currency,
    }, nil
}
```

Keeping fields private can ensure:

```text
If Money exists, Money is valid.
```

Ideally value objects behave immutably:

```go
func (m Money) Add(other Money) (Money, error)
```

returns another Money instead of mutating shared state.

Value semantics fit Go particularly well here.

---

## Entity

### Mental model

An Entity is defined primarily by **identity across time**.

```text
Order #123
```

is still the same order after its status changes.

### Deeper explanation

Suppose:

```go
type Order struct {
    id     OrderID
    status OrderStatus
}
```

Yesterday:

```text
Order #123 = Pending
```

Today:

```text
Order #123 = Paid
```

Its properties changed, but its identity did not.

Compare this with a Value Object:

```text
Money(10 PLN)
```

Changing it to `Money(20 PLN)` conceptually gives you a different value.

Thus:

```text
Value Object → equality by value
Entity       → equality by identity
```

---

## DTO

### Mental model

DTO means **Data Transfer Object**.

Its job is to move data across a boundary.

```go
type CreateOrderRequest struct {
    CustomerID string `json:"customer_id"`
    ProductIDs []string `json:"product_ids"`
}
```

### Deeper explanation

DTOs normally have little or no business behavior.

They exist because external representations have different needs from domain objects.

For example:

```text
HTTP JSON
       │
       ▼
CreateOrderRequest
       │ conversion
       ▼
CreateOrder command/domain values
```

A response can similarly be:

```go
type OrderResponse struct {
    ID     string `json:"id"`
    Status string `json:"status"`
}
```

Keeping DTOs separate prevents your domain model from becoming coupled to JSON, database columns, GraphQL, CSV, etc.

---

## Aggregate

### Mental model

An Aggregate is:

> A cluster of domain objects that must maintain certain invariants together.

One Entity acts as its **Aggregate Root**.

For example:

```text
Order
 ├── OrderItem
 ├── OrderItem
 └── ShippingAddress
```

`Order` might be the aggregate root.

### Deeper explanation

Suppose this invariant exists:

```text
A shipped order cannot receive new items.
```

External code should not modify individual `OrderItem`s arbitrarily.

Instead:

```go
err := order.AddItem(product, qty)
```

The `Order` controls the operation:

```go
func (o *Order) AddItem(product ProductID, qty int) error {
    if o.status == Shipped {
        return ErrAlreadyShipped
    }

    // modify internal state
    return nil
}
```

An Aggregate is often also a **transactional consistency boundary**.

You should be careful about making aggregates enormous. The idea is not “one entire business domain object graph.”

It is:

> Which things truly need to change atomically to preserve invariants?

---

# 10. Aggregate root encapsulation and public API

## Aggregate root encapsulation

### Mental model

External code should tell an aggregate **what happened**, rather than directly manipulating its internals.

Bad:

```go
order.Status = "shipped"
order.Items = append(order.Items, item)
```

Better:

```go
order.AddItem(item)
order.Ship(now)
```

### Deeper explanation

Consider:

```go
type Order struct {
    id      OrderID
    status  OrderStatus
    items   []OrderItem
    shipped *time.Time
}
```

Nothing can directly mutate these fields outside the package.

Instead, your domain exposes a controlled API:

```go
func NewOrder(...) (*Order, error)

func (o *Order) AddItem(...) error

func (o *Order) RemoveItem(...) error

func (o *Order) Pay(...) error

func (o *Order) Ship(...) error
```

This API is effectively the permitted state machine:

```text
Created
   │ Pay()
   ▼
Paid
   │ Ship()
   ▼
Shipped
```

Invalid transitions don't exist as public mutation primitives.

---

## Public API as a domain boundary

### Mental model

The exported identifiers of a Go package are a gate around the domain.

```text
outside
    │
    │ exported API only
    ▼
┌─────────────────────┐
│ order package       │
│                     │
│ internal state      │
│ invariants          │
│ business rules      │
└─────────────────────┘
```

### Deeper explanation

Go's package visibility makes this especially powerful.

Imagine:

```go
package order

type Order struct {
    id     ID
    status status
}
```

Consumers cannot do:

```go
order.status = whatever
```

They must invoke exported behavior.

That means the compiler helps enforce architectural boundaries.

Be careful when returning slices/maps:

```go
func (o Order) Items() []Item {
    return o.items
}
```

The caller may potentially mutate shared backing data.

If encapsulation matters, return a copy:

```go
func (o Order) Items() []Item {
    return slices.Clone(o.items)
}
```

Encapsulation includes indirect references, not merely field capitalization.

---

# 11. Testing

## `httptest`

### Mental model

`net/http/httptest` lets you test HTTP code using Go's normal HTTP abstractions without starting a manually managed production server.

Two especially useful tools are:

```go
httptest.NewRequest(...)
httptest.NewRecorder()
```

### Deeper explanation

For a handler:

```go
func GetUser(w http.ResponseWriter, r *http.Request) {
    // ...
}
```

you can test:

```go
req := httptest.NewRequest(
    http.MethodGet,
    "/users/123",
    nil,
)

rec := httptest.NewRecorder()

GetUser(rec, req)

res := rec.Result()
defer res.Body.Close()

if res.StatusCode != http.StatusOK {
    t.Fatalf("got %d", res.StatusCode)
}
```

This is excellent for handler-level tests.

There is also:

```go
server := httptest.NewServer(handler)
defer server.Close()
```

which creates a real local HTTP server.

Then a real HTTP client can call:

```go
http.Get(server.URL + "/users/123")
```

This exercises more of the network stack than `ResponseRecorder`.

Think of these as different test scopes:

```text
NewRecorder
   → handler behavior

NewServer
   → HTTP client/server integration
```

---

## Testcontainers

### Mental model

Testcontainers lets your integration tests start real infrastructure in containers.

Instead of mocking PostgreSQL:

```text
test
 │
 ├── start real PostgreSQL container
 ├── run migration
 ├── exercise repository/application
 └── clean container up
```

### Deeper explanation

For Go, the official `testcontainers-go` package integrates with ordinary `go test` and is primarily intended for integration/end-to-end testing with containerized dependencies. It manages container startup and cleanup, and supports readiness strategies such as waiting for a port, HTTP endpoint, log message, health check, or command. ([Testcontainers dla Go][1])

This lets you test the thing mocks cannot realistically reproduce:

```text
SQL dialect
constraints
transactions
indexes
serialization
driver behavior
migrations
```

A useful testing pyramid might therefore be:

```text
many
│
│ pure domain tests
│ application tests with fakes
│ repository integration tests + Testcontainers
│ HTTP integration tests
│
few
```

The tradeoff is speed and operational complexity. Container tests are slower than pure unit tests.

Use Testcontainers where **real infrastructure behavior is part of what you need confidence in**.

---

# 12. Dependency injection

## Dependency injection and explicit dependencies

### Mental model

Dependency injection simply means:

> Give an object/function the things it needs instead of letting it secretly find them.

Bad:

```go
func CreateOrder(...) {
    db := globalDatabase
}
```

Better:

```go
type CreateOrder struct {
    repo OrderRepository
}

func NewCreateOrder(repo OrderRepository) *CreateOrder {
    return &CreateOrder{
        repo: repo,
    }
}
```

### Deeper explanation

Suppose a use case needs:

```text
OrderRepository
PaymentGateway
Clock
```

Represent that honestly:

```go
type Service struct {
    orders   OrderRepository
    payments PaymentGateway
    clock    Clock
}

func NewService(
    orders OrderRepository,
    payments PaymentGateway,
    clock Clock,
) *Service {
    return &Service{
        orders:   orders,
        payments: payments,
        clock:    clock,
    }
}
```

Now you can understand the dependencies from the constructor alone.

The composition root, often `main`, constructs everything:

```text
main
 │
 ├── postgres repository
 ├── stripe adapter
 ├── application service
 ├── HTTP handler
 └── HTTP server
```

This is preferable to a service locator:

```go
container.Get("database")
container.Get("payments")
```

because service locators hide dependencies.

Go frequently needs no DI framework at all. Constructors and interfaces are sufficient.

---

# 13. Package and architecture boundaries

## Package boundaries

### Mental model

A Go package is both:

```text
code organization
+
visibility boundary
+
dependency unit
```

Imports describe architecture.

### Deeper explanation

Suppose:

```text
internal/
    order/
        order.go
        repository.go

    postgres/
        order_repository.go

    httpapi/
        order_handler.go
```

If `order` imports `postgres`, your domain now depends upon infrastructure.

Instead:

```text
httpapi ───────► application/order
                     ▲
                     │ implements interface
                     │
postgres ────────────┘
```

At compile time the actual imports might be:

```text
httpapi → order
postgres → order
main → httpapi + postgres + order
```

The central domain never imports PostgreSQL.

This is dependency inversion expressed through Go packages.

The `internal` directory adds another useful boundary: packages under it cannot be imported arbitrarily from outside their allowed parent tree.

---

## Vertical slice

### Mental model

A vertical slice organizes code around **business capability/use case** rather than purely around technical type.

Layer-first organization:

```text
handlers/
services/
repositories/
models/
```

Vertical organization:

```text
orders/
    create/
    cancel/
    get/

customers/
    register/
    get/
```

### Deeper explanation

Layer-first architecture tends to make one feature spread everywhere:

```text
Create Order
 │
 ├── handlers/order.go
 ├── services/order.go
 ├── repositories/order.go
 ├── dto/order.go
 └── models/order.go
```

A vertical slice tries to keep highly related change together:

```text
order/
    create/
        handler.go
        command.go
        service.go

    get/
        handler.go
        query.go

    domain/
        order.go
```

This is especially useful when different use cases evolve independently.

Vertical slices do **not** mean abandoning architecture or putting everything in one file.

Rather:

> Organize primarily by business change boundaries, and introduce technical layering inside a slice only where useful.

A small codebase may need far fewer directories than a large one.

---

## Dependency direction between layers

### Mental model

Dependencies should generally point toward business policy.

A common direction is:

```text
HTTP
 │
 ▼
Application
 │
 ▼
Domain
```

Infrastructure adapts itself to interfaces required by those layers:

```text
PostgreSQL
     │
     ▼
Application/Domain interface
```

### Deeper explanation

Imagine:

```go
// application
type OrderRepository interface {
    Get(ctx context.Context, id order.ID) (*order.Order, error)
    Save(ctx context.Context, order *order.Order) error
}
```

Postgres implements this:

```go
type OrderRepository struct {
    db *sql.DB
}
```

The application doesn't know PostgreSQL exists.

At runtime:

```text
HTTP Handler
     │
     ▼
CreateOrder use case
     │
     ▼
OrderRepository interface
     ▲
     │
Postgres implementation
```

But compile-time package dependencies look roughly like:

```text
http adapter ───► application
postgres adapter ─► application/domain
application ─────► domain
domain ──────────► almost nothing
```

`main` is special because it wires everything together:

```go
repo := postgres.NewOrderRepository(db)
service := orders.NewService(repo)
handler := httpapi.NewOrderHandler(service)
```

That is your **composition root**.

---

# 14. Simple CQRS

## CQRS: read model vs domain model

### Mental model

CQRS means:

> The model used to change state does not have to be the same model used to read state.

It does **not** inherently mean Kafka, event sourcing, asynchronous queues, or separate databases.

A perfectly valid simple CQRS architecture can be:

```text
                    PostgreSQL
                   /          \
                  /            \
Write command → Domain       SQL query → Read DTO
```

### Deeper explanation

Consider displaying an order page.

Without CQRS, you may load:

```go
Order aggregate
```

with all domain behavior just to produce:

```json
{
  "id": "...",
  "customer_name": "...",
  "total": "...",
  "status": "..."
}
```

That may require loading several domain entities only to read data.

CQRS says the read side can simply have:

```go
type OrderDetails struct {
    ID           string
    CustomerName string
    Total        string
    Status       string
}
```

and execute an optimized SQL query:

```text
orders
 JOIN customers
 JOIN ...
```

directly into that read model.

Meanwhile the write side does:

```text
CancelOrder command
       │
       ▼
load Order aggregate
       │
       ▼
order.Cancel()
       │
       ▼
save aggregate
```

The difference is intentional:

```text
WRITE SIDE
command
  │
  ▼
domain aggregate
  │
enforce invariants
  │
  ▼
database


READ SIDE
query
  │
  ▼
optimized SQL
  │
  ▼
read DTO
```

Both can use exactly the same database and transactionally consistent data.

No queue is required.

No eventual consistency is required.

No event sourcing is required.

No microservices are required.

Simple CQRS is often just:

> Rich domain model for writes, purpose-built projections/DTOs for reads.

---

# 15. How these concepts fit together

Consider a `CreateOrder` HTTP endpoint.

```text
POST /orders
      │
      ▼
HTTP handler
```

The handler decodes JSON into a DTO:

```go
type CreateOrderRequest struct {
    CustomerID string `json:"customer_id"`
}
```

It extracts the propagated context:

```go
ctx := r.Context()
```

and calls the application layer explicitly:

```go
err := createOrder.Execute(ctx, cmd)
```

The application service has explicit dependencies:

```go
type CreateOrder struct {
    orders OrderRepository
    clock  Clock
}
```

It creates domain value objects:

```text
CustomerID
OrderID
Money
```

and an aggregate:

```text
Order
 ├── ID
 ├── Items
 └── Status
```

The aggregate's unexported fields protect invariants:

```go
type Order struct {
    id     OrderID
    items  []Item
    status Status
}
```

State changes happen through its public API:

```go
order.AddItem(...)
order.Confirm(...)
```

The repository interface belongs toward the application/domain side:

```go
type OrderRepository interface {
    Save(context.Context, *Order) error
}
```

PostgreSQL implements it.

Errors travel upward:

```text
PostgreSQL error
      ↓
repository translation
      ↓
application context/wrapping
      ↓
HTTP translation
      ↓
400 / 404 / 409 / 500
```

`context.Context` travels downward:

```text
HTTP
 ↓
application
 ↓
repository
 ↓
database
```

so cancellation propagates in the opposite direction:

```text
client disconnects
        ↓
HTTP context cancelled
        ↓
repository/database operation cancelled
```

Reads need not reconstruct the aggregate:

```text
GET /orders/123
      │
      ▼
query handler
      │
      ▼
SQL optimized for display
      │
      ▼
OrderDetails DTO
      │
      ▼
JSON
```

And tests naturally split:

```text
domain rules
    → ordinary table-driven unit tests

HTTP mapping
    → httptest

PostgreSQL repository
    → Testcontainers

whole application
    → selected integration tests
```

That combination is a very idiomatic way to build a moderately complex Go application: **explicit dependencies, small interfaces, value-oriented types, package-level encapsulation, streaming interfaces at boundaries, domain models for writes, simple read models for queries, and errors/context propagated deliberately.**

[1]: https://golang.testcontainers.org/?utm_source=chatgpt.com "Testcontainers for Go"
