# Functions

Functions enable lightweight abstraction without requiring objects or interfaces.
Functions are values.

## assignment, passing/injecting, Closures, type

```go
func add(a, b int) int {
    return a + b
}

operation := add                // can be assigned
operation(2, 3)

func calculate(a, b int, operation func(int, int) int) int {
    return operation(a, b)
}
calculate(1, 2, operation)      // passed arround

func multiplier(n int) func(int) int {
    return func(x int) int {    // Closures capture surrounding variables
        return x * n
    }
}

type Clock func() time.Time     // can make type
func NewToken(clock Clock) {    // and be injected as behavior
    now := clock()
}

fixedClock := func() time.Time {    // test can inject double
    return fixedTime
}
```

> **Be aware** that captured variables can become shared mutable state, especially when goroutines are involved.

## Function composition - orchestration, decorator

> Compose behavior - building larger operations from smaller functions.

**Orchestration** - instead of one giant function, one operation can orchestrate them:
```go
func trim(s string) string
func normalize(s string) string
func validate(s string) error

s = trim(s)
s = normalize(s)

if err := validate(s); err != nil {
    return err
}
```

**Decorator** - add behavior:
```go
type Handler func(context.Context) error

func WithLogging(next Handler) Handler {
    return func(ctx context.Context) error {
        log.Println("start")
        err := next(ctx)
        log.Println("end")
        return err
    }
}
```

## Variadic params, Functional-options pattern

**Variadic parameters** accept zero or more arguments:
```go
func sum(values ...int) int {
    total := 0
    for _, v := range values {      // `values` is a slice.
        total += v
    }
    return total
}

sum()
sum(1)
sum(1, 2, 3)

numbers := []int{1, 2, 3}
sum(numbers...)
```

**Functional-options pattern**
Configure constructor without params list - for truly optional configuration, should not hide mandatory dependencies.
- config highly optional
- some options have richer behavior then just set field
- config likely to evolve - for backwards-compatibility (to add new params without breaking existing calls)
- makes consturctor calls self-documenting
```go
type config struct {
    timeout time.Duration
    retries int
}
type Option func(*config) error
// type Option func(*Server) error     // OR modify obj directly

func WithTimeout(to time.Duration) Option {
    // return func(s *Server) error {
    return func(cfg *config) {          // receives constructed data and applies config change
        if to <= 0 {
            return fmt.Errorf("timeout must be positive")
        }
        cfg.timeout = to
        // s.timeout = to
    }
}

func NewServer(addr string, opts ...Option) *Server {   // optional config
    cfg := config{
        timeout: 30 * time.Second,
        retries: 3,
    }
    // s := &Server{
    //     addr:           addr,
    //     timeout:        30 * time.Second,       // defaults
    //     maxConnections: 10,
    // }

    for _, opt := range opts {                  // + override defaults
        // if err := opt(s); err != nil {
        if err := opt(&cfg); err != nil {
            return nil, err
        }
    }

    return &Server{
        timeout: cfg.timeout,
        retries: cfg.retries,
    }
}
srv := NewServer(
    ":8080",
    WithTimeout(5*time.Second),
    WithMaxConnections(100),
)

NewService(db Database)                                 // required deps
```

**Config struct**
Simpler, often better - config opts often set together
```go
srv := NewServer(":8080", ServerConfig{
    Timeout:        5 * time.Second,
    MaxConnections: 100,
})
```

---

# Custom Types (enums)

Gives domain meaning to primitives and compile-time type safety - can't accidentally pass a raw `int` where a `ProductCategory` is expected.

```go
type UserID string          // defined type
type UserName = string      // alias
```

Enum:
```go
type ProductCategory int                            // custom int type
const (
	CategoryVentilation ProductCategory = iota      // constant generator (1st=0, each subseqquent auto-icrements)
	CategoryFilter
	// CategoryVentilation ProductCategory = "prod_cat"     // for string (convenient around JSON/database boundaries)
)

c := ProductCategory(999)       // still can be created, so validation at boundaries may be necessary

func (pc ProductCategory) Valid() bool {            // attach behavior
	switch pc
	case CategoryVentilation, CategoryFilter:
		return true
	default:
		return false
	}
}

pc := 1
switch ProductCategory(pc) {
    case CategoryVentilation:
```
