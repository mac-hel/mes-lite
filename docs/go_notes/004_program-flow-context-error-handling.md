# `defer`

`defer` schedules a call for execution when the surrounding function returns.

```go
defer fmt.Println(x)        // captures current value of `x` (args evaluated on `defer` statement)
defer func() {              // but Closure receives future value (if it change)
    fmt.Println(x)
}()
defer iWillRunFirst()       // deferred functions execute in **LIFO** order.
```

Deferred calls **run**:
* on normal return
* while unwinding a panic

They do **not** run when the process terminates directly through mechanisms:
* `os.Exit`
* `log.Fatal` (which exits the process)

Common uses:
```go
file.Close()
mu.Unlock()
cancel()
rows.Close()
```

---

# Exit

`os.Exit` / `log.Fatal` terminates immediately
does not trigger deferred functions

---

# Context

Context carries **request lifetime information** down a call tree.

Contexts form a tree.
```go
ctx
 ├── database call
 ├── HTTP call
 └── another operation
```
If the parent request is cancelled, children should eventually stop too.

## What belongs in `context.Context`?
* cancellation
* deadlines
* request-scoped values: infra metadata, tracing/spans, correlation/request IDs, and sometimes authentication identity metadata
* **NOT** business data, domain params (amount, orderID), config (db, feature flags)

Use typed, preferably unexported keys:
```go
type contextKey struct{}
ctx = context.WithValue(ctx, contextKey{}, value)
```

## context propagation
```go
ctx, cancel := context.WithTimeout(parent, 2*time.Second)
defer cancel()

db.QueryContext(ctx, query)                             // same ctx to DB
http.NewRequestWithContext(ctx, method, url, body)      // same ctx to downstream HTTP request

```
Also, do not pass `nil` context. Use: `context.Background()` or a propagated parent context.

## Context as first parameter is the established Go convention, that provides:

### 1. Immediate visibility
A reader can instantly see that the operation participates in cancellation/deadline propagation.

### 2. Consistent propagation
Call chains naturally look like:
```go
handler(ctx)
    -> service.Get(ctx)
        -> repo.Get(ctx)
            -> db.QueryContext(ctx)
```
The same cancellation/deadline propagates through API boundaries.

### 3. Cancellation
If the HTTP client disconnects:
```text
HTTP request context cancelled
    ↓
service operation cancelled
    ↓
database operation cancelled
```
assuming downstream APIs respect the context.

### 4. Deadlines
```go
ctx, cancel := context.WithTimeout(parent, 2*time.Second)
defer cancel()
```
Downstream work receives the remaining deadline.

### 5. Request-scoped values

Appropriate examples:
* trace/span IDs
* authenticated principal
* request ID

Usually inappropriate:
* application configuration
* logger dependencies solely to avoid passing them
* database handles
* ordinary function parameters

Avoid using context as an untyped parameter bag.

---

## Don't usually store context in structs

Prefer:

```go
func (s *Service) Process(ctx context.Context) error
```

rather than:

```go
type Service struct {
    ctx context.Context
}
```

The caller should control the lifetime of each operation. Passing context explicitly preserves that relationship. The official Go guidance recommends passing `Context` as an argument, generally first.

Never pass `nil` as a context. Use:

```go
context.Background()
context.TODO()
```

when necessary.

---

# Errors

Errors are values: `if err != nil { ... }`
Sentinel error: `var ErrNotFound = errors.New("not found")`
    - Useful when callers need only a category.
    - Potential downside: exported sentinel becomes API coupling
Typed error: `type MyErr struct { Field string; Err error }`
    - Useful when the caller needs structured information.
    - `func (e MyErr) Error() string { return fmt.Sprintf("f %s: %v", e.Field, e.Err) }`    - human readable string
    - `func (e MyErr) Unwrap() error { return e.Err }`                                      - unwrap underlying error

## Wrapping
**wrap** to add context at abstraction boundaries, `%w` preserves the error chain:
```go
fmt.Errorf("load user %q: %w", id, err)
```

**Inspect** (is recursive), prefer `errors.Is`/`errors.As` over `equality` or `type assertions` when errors may be wrapped:
```go
if errors.Is(err, sql.ErrNoRows) { ... }
var e *MyErr; if errors.As(err, &e) { ... }
```

## Error ownership
- avoid repeatedly wrapping with meaningless text
- do not log and return error at every layer; that often causes duplicate logs
    - usually one sufficiently high boundary owns the final logging
> Handle an error when you can recover from it or translate it meaningfully. Otherwise propagate it.

A good layered flow might look like:
```text
PostgreSQL
    │
    │ sql.ErrNoRows
    ▼
Repository adapter
    │
    │ domain.ErrOrderNotFound
    ▼
Application use case
    │
    │ "load order abc: order not found"
    ▼
HTTP handler
    │
    └── HTTP 404
```
Each boundary translates into the vocabulary of the layer above it.
