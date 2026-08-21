# `defer`

`defer` schedules a call for execution when the surrounding function returns.

Common uses:

```go
file.Close()
mu.Unlock()
cancel()
rows.Close()
```

Deferred calls run:

* on normal return
* while unwinding a panic

They do **not** run when the process terminates directly through mechanisms such as `os.Exit`; consequently `log.Fatal`, which exits the process, also bypasses ordinary deferred cleanup.

Arguments to a deferred call are evaluated when the `defer` statement executes:

```go
defer fmt.Println(x)
```

captures the current value of `x`.

Deferred functions execute in **LIFO** order.

---

# Exit

- `os.Exit` does not trigger deferred functions

---

# Context

`context.Context`

Context carries:
* cancellation
* deadlines
* request-scoped values

Typical API:
```go
func Fetch(ctx context.Context, id string) (User, error)
```

---

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
Typed error: `type ValidationError struct { Field string }`
    - Useful when the caller needs structured information.

## Wrapping
**wrap** to add context at abstraction boundaries, `%w` preserves the error chain:
`fmt.Errorf("load user %q: %w", id, err)`

**Inspect** (is recursive), prefer `errors.Is`/`errors.As` over `equality` or `type assertions` when errors may be wrapped:
`if errors.Is(err, sql.ErrNoRows) { ... }`
`var e *ValidationError; if errors.As(err, &e) { ... }`

## Error ownership
- avoid repeatedly wrapping with meaningless text
- do not log and return error at every layer; that often causes duplicate logs
