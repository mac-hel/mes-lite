# Go — senior engineer crib

> last incorporated lesson: M03-L001

## Modules, packages, layout

* **module** = versioned collection of Go packages
  * defined by `go.mod`
  * module path identifies the module
  * dependencies use **Minimal Version Selection (MVS)**
  * `go get` changes dependencies / versions
  * `go mod tidy` synchronizes `go.mod` / `go.sum` with imports
  * `go work` creates a workspace containing multiple local modules

* **package** = compilation / encapsulation unit
  * all ordinary `.go` files in one directory normally belong to the same package
  * exported identifiers start with uppercase letters
  * package dependencies must form an acyclic graph

* `package main`
  * executable program
  * must contain `func main()`
  * multiple binaries commonly live under:

    ```text
    cmd/
      api/
        main.go
      worker/
        main.go
    ```

* `internal/`
  * compiler/toolchain-enforced **import boundary**
  * code outside the parent tree of an `internal` directory cannot import it
  * useful for implementation details that must not become public APIs
  * stronger than naming conventions, but not equivalent to “module-private” in every possible layout

* `pkg/`
  * optional community convention (**not special to the Go toolchain**)
  * often unnecessary
  * public packages can live directly at repository/module root

* organize packages around **cohesive responsibilities**, not technical layer names by default
  * e.g. `orders`, `payments`, `users`
  * a domain package may own:
    * domain types
    * behavior/use cases
    * interfaces required by that behavior
    * tests
  * HTTP/database implementations can either live there or in adjacent packages depending on coupling and size

Go recommends short, clear, lowercase, usually single-word package names and discourages vague buckets such as `util`, `common`, `types`, and `interfaces`.

---

## Package naming: singular vs plural

There is **no grammatical rule** saying Go packages must always be singular.

Choose the name that makes imported usage read naturally:

```go
http.Server
time.Duration
bytes.Buffer
strings.TrimSpace
users.Lookup(...)
orders.Service
```

Guidelines:

* prefer a short noun or domain concept
* singular often works well for a capability or abstraction:

  * `user`
  * `order`
  * `payment`
* plural is fine when it naturally describes operations over a collection/domain:

  * `strings`
  * `bytes`
  * `slices`
  * `maps`
* avoid awkward repetition:

  ```go
  user.UserService       // suspicious
  users.Service          // often better
  ```
* judge the API from the **caller side**:

  ```go
  package.Type
  package.Function()
  ```

The package name forms part of every exported name, so the combination should read naturally.

---

# Application architecture

## Composition root

`main()` is commonly the **composition root**
* load configuration, initialize infrastructure, wire dependencies, start the application, coordinate shutdown
* prefer explicit dependency passing over hidden global dependencies

Usually avoid:
* mutable global state
* service-locator patterns
* `init()` for dependency wiring (appropriate for rare package initialization requirements, registration mechanisms, generated code, etc., but should not hide the application's dependency graph.)
* DI framework - often unnecessary because ordinary constructors provide explicit dependency injection

---

## Operational endpoints

Typical service endpoints:

```text
/health     # is the process alive / not irrecoverably stuck?
/ready      # should this instance receive traffic?
/version    # build/revision information
/metrics    # telemetry, commonly Prometheus-compatible
```

Do not make liveness depend on every external dependency; a temporary database outage should not necessarily cause Kubernetes to restart a healthy process.

---

## Graceful shutdown

Typical order:

1. receive shutdown signal
2. stop accepting new work
3. mark service unready
4. ask HTTP/gRPC server to shut down
5. allow in-flight requests to finish within a deadline
6. cancel application/root context
7. wait for owned goroutines
8. flush telemetry/logging if needed
9. close resources
10. exit

Important distinction:

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()
```

Context cancellation is a **signal**, not a goroutine join mechanism.

Use synchronization such as:

```go
sync.WaitGroup
errgroup.Group
```

to wait for goroutines to actually terminate.

Every goroutine should have an understood lifetime:

> What starts it?
> What causes it to stop?
> Who waits for it?

---

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

# `context.Context`

Context carries:
* cancellation
* deadlines
* request-scoped values

Typical API:
```go
func Fetch(ctx context.Context, id string) (User, error)
```

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

# Various

- `import time`: var now time.Time = time.Now().UTC(), var ttl time.Duration = 15 * time.Minute, now.Add(ttl)
- `os.Exit` - does not trigger deferred functions
- `crypto/rand` and `encoding/hex` for standard-library UUID-shaped IDs
- blank imports
- variadic function parameters: `func (nums ...int)`

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

---

# Validation

**transport**
validates only HTTP request shape, not domain rules, e.g. by using struct tags

**application Services**
validate workflow rules and references

**Domain**
- constructor `NewProduct` validates domain rules (invariants)
- mutation method `UpdateDetails` replaces mutable product fields and preserves product invariants
- `ReconstituteProduct` validates reconstructing data - should be private if possible
- both use common `validate` function but can adjust validation for own requirements

**repository**
validate before persistence as defense-in-depth

**PostgreSQL**
constraints provide the final integrity boundary

---

# Interfaces

Interfaces are satisfied **implicitly**.

```go
type Reader interface {
    Read([]byte) (int, error)
}
```

No `implements Reader` declaration exists.

Prefer small interfaces focused on the behavior the consumer actually needs.

```go
type UserFinder interface {
    Find(ctx context.Context, id string) (User, error)
}
```

rather than prematurely designing large interfaces around an implementation.

---

## Why should the consumer package own the interface?

Example:

```go
package users

type UserStore interface {
    Find(ctx context.Context, id string) (User, error)
}

type Service struct {
    store UserStore
}
```

Implementation:

```go
package postgres

type UserStore struct {
    db *sql.DB
}

func (s *UserStore) Find(...) (...) { ... }
```

`postgres.UserStore` implicitly satisfies `users.UserStore`.

### Why this is useful

The **consumer knows what behavior it requires**.

`users.Service` may need only:

```go
Find(...)
```

while the database implementation may expose:

```go
Find(...)
Create(...)
Update(...)
Delete(...)
List(...)
Count(...)
BeginTx(...)
```

If the producer owns a large interface:

```go
type Repository interface {
    Find(...)
    Create(...)
    Update(...)
    Delete(...)
    List(...)
}
```

every consumer becomes coupled to capabilities it does not need.

Consumer-owned interfaces give:

* interface segregation
* smaller APIs
* easier fakes in tests
* less coupling to implementation details
* freedom for the producer to add methods to its concrete type

A producer can add:

```go
func (*Store) Vacuum()
```

without changing any consumer interface.

Go's official review guidance says interfaces generally belong in the package that **uses** them and implementations should normally return concrete types.

### Don't create interfaces only to mock something

Avoid:

```go
type UserServiceInterface interface {
    ...
}
```

merely because `UserService` exists.

Introduce an interface where there is an actual abstraction boundary / consumer need.

---

# “Accept interfaces, return structs”

More accurately:

> **Accept the narrowest useful abstraction; return concrete types when the implementation is known.**

It is a heuristic, not a language rule.

## Why accept interfaces?

Suppose you only need to write bytes:

```go
func Encode(w io.Writer, v Value) error
```

This works with:

```go
os.File
bytes.Buffer
http.ResponseWriter
gzip.Writer
network connections
```

Contrast:

```go
func Encode(w *os.File, v Value) error
```

which unnecessarily restricts callers.

Accepting a small interface expresses:

> “I care about this capability, not this implementation.”

---

## Why return concrete types?

Prefer:

```go
func NewClient(...) *Client
```

over:

```go
func NewClient(...) ClientInterface
```

when the function always creates a `*Client`.

Benefits:

### Full API remains available

The caller can use every exported method on `*Client`.

Returning an interface immediately hides methods not present in that interface.

### The producer can evolve

A new method can be added to:

```go
type Client struct { ... }
```

without changing an interface.

### Consumers define their own abstractions

A caller needing only:

```go
type Sender interface {
    Send(context.Context, Message) error
}
```

can declare that locally.

### Less premature abstraction

Returning an interface implies there is meaningful implementation polymorphism.

Often there isn't.

### Better type information

Concrete return values preserve:

* method set
* fields where exported
* documentation
* compile-time discoverability

Exceptions are common and legitimate:

```go
func Open(...) (io.ReadCloser, error)
```

when deliberately hiding implementations or when multiple implementation types are fundamental to the API.

`error` is itself an obvious interface return type.

---

# Struct return: value or pointer?

These are separate decisions:

```go
func NewConfig() Config
```

versus:

```go
func NewClient() *Client
```

## Return a value when

* the type is small
* copying has clear value semantics
* identity is irrelevant
* callers should freely copy it
* zero/value-style semantics make sense
* the type does not contain state that must not be copied

Examples:

```go
time.Time
net.IPNet // depending on semantics
Config
Coordinate
Money
```

Example:

```go
func NewPoint(x, y float64) Point {
    return Point{X: x, Y: y}
}
```

Returning a value communicates:

> “This is a value, not an identity-bearing mutable object.”

---

## Return a pointer when

* methods mutate the object
* object identity matters
* the struct is large enough that copying is undesirable
* the object contains `sync.Mutex`, `sync.Once`, atomics, or other values that should not be copied after use
* `nil` has useful semantic meaning
* construction establishes shared mutable state

Example:

```go
func NewServer(...) *Server
```

---

## Why return a value from a constructor?

A constructor does **not** have to return a pointer.

```go
func NewConfig(...) Config
```

can be preferable when the result has value semantics.

Benefits:

* communicates copyability
* avoids introducing nil as an extra state
* caller can choose whether to take its address:

  ```go
  cfg := NewConfig(...)
  use(&cfg)
  ```
* often makes APIs easier to reason about

Do **not** choose value vs pointer merely based on an assumption that “pointers avoid allocation.”

Go's compiler performs escape analysis. Whether data lives on the stack or heap is an implementation decision influenced by how the value escapes, not simply whether the source code contains `&T`.

---

# Methods and receivers

A receiver attaches a method to a defined type:

```go
func (u User) Name() string
```

## Value receiver

```go
func (p Point) Distance() float64
```

The receiver itself is copied.

Good when:

* type is small
* it has value semantics
* methods don't mutate receiver state
* copying is valid

Important:

> A value receiver does not imply deep immutability.

```go
type X struct {
    values []string
}
```

A copied `X` still has a slice referring to the same backing array.

---

## Pointer receiver

```go
func (u *User) Rename(name string)
```

Use when:

* method mutates receiver
* copying is undesirable
* type contains synchronization primitives
* identity matters
* consistency with other methods warrants it

For consistency avoid mixing pointer/value receivers without reason.

---

# Zero values

Every variable has a zero value.

Examples:

```go
int       -> 0
bool      -> false
string    -> ""
pointer   -> nil
slice     -> nil
map       -> nil
channel   -> nil
interface -> nil
```

A well-designed type often makes its zero value useful:

```go
var mu sync.Mutex
var buf bytes.Buffer
```

But not every type can or should have a fully useful zero value.
Prefer useful zero values where practical; use constructors where initialization establishes meaningful invariants or dependencies.

Constructors remain useful when:

* invariants must be established
* dependencies are required
* defaults differ from zero values
* internal maps/channels must be initialized
* validation is required

---

## Calling methods on nil pointers

This can be legal:

```go
type Sub struct {
    value string
}

func (s *Sub) Value() string {
    if s == nil {
        return ""
    }
    return s.value
}

var s *Sub
fmt.Println(s.Value())
```

Method dispatch itself does not automatically dereference `s`.

But this:

```go
func (s *Sub) Value() string {
    return s.value
}
```

will panic when `s == nil`.

Therefore:

> nil-receiver support is an explicit API design choice, **not** a general zero-value property.

---

# Nil interfaces

nil interface != interface containing nil pointer

Conceptually, an interface value contains:

```text
(dynamic type, dynamic value)
```

The zero value of an interface is:

```text
(nil, nil)
```

Example:

```go
var err error
fmt.Println(err == nil) // true
```

But:

```go
var p *MyError = nil
var err error = p
```

Conceptually:

```text
(*MyError, nil)
```

Therefore:

```go
fmt.Println(err == nil) // false
```

The interface itself is not nil because it contains a dynamic type.

## Classic bug

```go
func doSomething() error {
    var err *MyError
    return err
}
```

The caller receives an `error` whose concrete type is `*MyError`, so:

```go
err := doSomething()
fmt.Println(err == nil) // false
```

Correct:

```go
func doSomething() error {
    var err *MyError

    if err == nil {
        return nil
    }

    return err
}
```

Or structure the function so a typed nil is never converted into the interface.

### Rule of thumb

When returning an interface:

```go
return nil
```

to represent absence, rather than returning a typed nil pointer through the interface.

---

# Embedding and composition

Go has no traditional class inheritance.

Prefer composition:

```go
type Person struct {
    Name string
}

type Employee struct {
    Person
    ID string
}
```

Embedding promotes fields/methods:

```go
employee.Name
```

but `Employee` **is not** a subtype of `Person`.

Embedding is primarily:

* composition
* delegation / method promotion

not inheritance.

Be cautious when embedding exported implementation types because promoted methods become part of the embedding type's API.

---

# DTOs and transport boundaries

Avoid letting HTTP/JSON representation accidentally define your core domain model.

Transport DTO:

```go
type createUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}
```

Domain:

```go
type User struct {
    Name  Name
    Email Email
}
```

Benefits:

* transport format can evolve separately
* validation can occur at the boundary
* domain does not accumulate JSON/database tags
* avoids over-posting fields accidentally
* API versions need not map 1:1 to domain structures

Do not create DTO layers mechanically when domain and transport structures genuinely have identical semantics and coupling is harmless.

---

# Concurrency

## Goroutine

```go
go work()
```

A goroutine is a lightweight concurrent execution unit managed by the Go runtime.

It is **not literally an OS thread**.

Many goroutines are multiplexed onto OS threads by the runtime scheduler.

---

## WaitGroup

```go
var wg sync.WaitGroup

wg.Add(1)
go func() {
    defer wg.Done()
    work()
}()

wg.Wait()
```

Coordinates completion.

A `WaitGroup` does not:

* propagate errors
* cancel goroutines

For structured concurrent work, `errgroup.Group` is often useful.

---

## Channels

Typed communication/synchronization mechanism:

```go
ch := make(chan Item)
```

### Unbuffered

```go
make(chan Item)
```

Send and receive synchronize directly.

### Buffered

```go
make(chan Item, 10)
```

Send can proceed until the buffer is full.

### Closing

```go
close(ch)
```

Means:

> no more values will ever be sent.

Usually the **sender/producer** owns closing.

Receiving:

```go
v, ok := <-ch
```

When closed and drained:

```go
ok == false
```

Range:

```go
for v := range ch {
    ...
}
```

ends when the channel is closed and drained.

Do not close a channel merely to “free memory.”

### Nil channel

```go
var ch chan Item
```

Sending or receiving on a nil channel blocks forever.

This property is useful for dynamically disabling `select` cases.

### Closed channel

* receive → immediately returns remaining buffered values, then zero value
* send → panic
* close again → panic

---

# Synchronization and concurrency ownership

Know when to use:

* `sync.Mutex`
* `sync.RWMutex`
* `sync.Once`
* `sync.WaitGroup`
* atomics
* channels
* contexts

Rule of thumb:

> Use channels to communicate ownership/events; use mutexes to protect shared state.

Not every concurrent problem needs channels.

Run:

```bash
go test -race ./...
```

regularly.

---

# Structures and memory allocation

## `new`

```go
p := new(T)
```

Equivalent conceptually to:

```go
var x T
p := &x
```

Returns:

```go
*T
```

pointing to a zero-valued `T`.

Example:

```go
p := new(int)
fmt.Println(*p) // 0
```

`new` does **not** mean “heap allocation.”

The compiler determines whether storage escapes and must live on the heap.

---

## `make`

`make` applies only to:

```text
slice
map
channel
```

Examples:

```go
s := make([]int, 10, 100)
m := make(map[string]int)
ch := make(chan Job, 10)
```

Unlike `new`, `make` returns the initialized value itself:

```go
make([]T, ...) -> []T
make(map[K]V)  -> map[K]V
make(chan T)   -> chan T
```

not a pointer.

---

# Arrays

Array length is part of its type:

```go
[3]int
[4]int
```

are different types.

Arrays have value semantics:

```go
b := a
```

copies the array.

Passing an array to a function copies it unless a pointer is used.

Arrays are less common directly; slices are usually the user-facing abstraction.

---

# Slices

A slice is a small descriptor referring to an underlying array, has current length and max capacity.
`pointer → backing array + length + capacity`
Lookup: O(n)

```go
var s []int                 // nil,     len(s) == 0 , cap(s) == ?
s := []int{1, 2, 3}         // non-nil, len(s) == 3 , cap(s) == ?
s := make([]int, 4, 10)     // non-nil, len(s) == 4 , cap(s) == 10

b := s                  // copies reference; NOT independent copy: both refer to the dame backing array
b[0] = 99               // now s[0] == 99

s = append(s, 23)       // if capacity insufficient, append will: allocate new backing array, copy elements
```

## Subslice aliases memory
Slice and subslice share backing storage - changes may be visible through both slices.
```go
huge := loadHugeBuffer()
tiny := huge[2:10]              // elements 2-9; [:10] elements 0-9; [2:] elements 2-LAST
```

Copy when necessary:
```go
tiny := append([]byte(nil), huge[:10]...)
    // or
tiny := make([]byte, 10)
copy(tiny, huge[:10])
```

## Nil vs empty slice
Both are valid for operations:  `range, append, len, cap`
But they can differ under:      `== nil`, JSON encoding, reflection, some API contracts
```go
var a []int        // nil,      length 0
b := []int{}       // non-nil,  length 0
c := make([]int,0) // non-nil,  length 0
```

## sorting
    import "sort"           // package for sorting Slices and user-defined Collections
	sort.Slice(prods, func(i, j int) bool {
		a, b := prods[i], prods[j]
		switch sortKey {
		case "-sku":
			return a.SKU > b.SKU
		case "name":
			return a.Name < b.Name
		case "-name":
			return a.Name > b.Name
		default:
			return a.SKU < b.SKU
		}
	})

---

# Maps

Runtime-managed structures with reference-like behavior.
Lookup: O(1)

```go
m := make(map[string]User)                  // initialize (memory) - map must be initialized before writing
m := make(map[User]struct{}, len(users))    // Map Set pattern
m["x"] = 1                                  // write entry

m2 := m                     // copies reference; NOT independent copy: mutation of one's entries is visible through the other
```

## Nil map
Safe to use (but 0 values) except for writing.
```go
var m map[string]int    // nil map
v := m["x"]             // 0
len(m)                  // 0
delete(m,"x")           // safe
m["x"] = 1              // INVALID - panic
```

## Missing key
Zero is returned for absent values; use `ok`.
```go
v := m[key]         // returns the value type's zero value when absent
v, ok := m[key]     // to distinguish absent from a stored zero value
    if !ok { ... }
```

## Map iteration
Sort keys explicitly when deterministic output matters.
```go
for k := range m    // iteration order is unspecified
```

## Map concurrency
Ordinary maps are not generally safe for unsynchronized concurrent read/write access, typically needs:
```go
sync.Mutex
sync.RWMutex
sync.Map        // specialized; a normal map protected by a mutex is often clearer
```
---

# Slice-backed vs map-backed in-memory store

Suppose:

```go
type Store struct {
    users []User
}
```

and change it to:

```go
type Store struct {
    users map[UserID]User
}
```

This is not simply an optimization. It changes useful semantics.

## Slice-backed store

Lookup by ID:

```text
O(n)
```

Advantages:

* naturally preserves insertion/order semantics
* simple iteration
* deterministic ordering if mutation rules preserve it
* lower conceptual overhead for very small collections
* convenient when primary operation is sequential traversal

Disadvantages:

* ID lookup is linear
* delete/update requires search
* duplicate IDs require explicit prevention

---

## Map-backed store

Expected lookup/update/delete:

```text
O(1)
```

Advantages:

* efficient key lookup
* naturally expresses uniqueness by key
* simple replacement:

  ```go
  users[id] = user
  ```
* deletion:

  ```go
  delete(users, id)
  ```

Disadvantages:

* iteration order is unspecified
* usually larger memory overhead per element
* cannot naturally represent duplicate keys
* sorting is required for deterministic ordered output
* may require an additional structure if both fast lookup **and order** matter

For example:

```go
map[ID]User
[]ID
```

may be maintained together.

### Important design question

Choose based on operations and semantics, not merely asymptotic complexity.

For:

```text
10 configuration records read once
```

a slice can be perfectly reasonable.

For:

```text
100,000 users repeatedly looked up by ID
```

a map is probably a better representation.

Benchmark if the difference actually matters:

```bash
go test -bench=. -benchmem
```

---

# Structs

Struct assignment is a value copy:

```go
b := a
```

But that copy may still contain reference-like fields:

```go
type X struct {
    Values []string
    Lookup map[string]int
    Ptr    *Config
}
```

Copying `X` copies:

* the slice descriptor
* map reference-like value
* pointer

not deep copies of their referenced data.

Thus:

> struct copy ≠ deep copy.

---

# Stack vs heap

Avoid manually reasoning:

```text
pointer -> heap
value   -> stack
```

That model is wrong.

The compiler performs **escape analysis**.

A value whose lifetime cannot safely remain within its current stack frame may escape to the heap.

Example:

```go
func answer() *int {
    x := 42
    return &x
}
```

This is perfectly safe Go.

The compiler arranges the lifetime of `x` appropriately.

Useful diagnostic:

```bash
go build -gcflags="-m"
```

But optimize only when measurements justify it.

Heap allocations increase garbage-collector work, while stack allocations are generally cheaper; modern Go continues to optimize escape/allocation behavior.

---

# Strings (stdlib package)

Immutable byte data: `pointer to bytes + length`; exact representation is an implementation detail.

```go
len(s)                  // counts **bytes**, not Unicode characters.
for _, r := range s     // iterates UTF-8 decoded runes (int32 representing Unicode code point)
```
strings.Contains
strings.ToLower
strings.TrimSpace
strings.HasPrefix(s, prefix)
strings.TrimPrefix(s, prefix)

strconv.Atoi(s)/ParseInt(S,10,0)        // str -> int
strconv.ParseBool(b)                    // str -> bool

strings.Count(str, substr)              // count non-overlapping substr in str

## Stringer interface
`fmt.Stringer` interface (`String() string`) - `fmt.Print`, `%s`, `%v` all produce a readable label; useful for logging and JSON serialization.

```go
func (c ProductCategory) String() string { switch c { case CategoryVentilation: return "Ventilation"
```

---

# Custom Types (enums)

Gives domain meaning to primitives and compile-time type safety - can't accidentally pass a raw `int` where a `ProductCategory` is expected.

```go
type ProductCategory int                            // custom int type
const (
	CategoryVentilation ProductCategory = iota      // constant generator (1st=0, each subseqquent auto-icrements)
	CategoryFilter
)

func (pc) Valid() bool {
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

---

# Authentication / Authorization

Password hashing
```go
import "golang.org/x/crypto/bcrypt"
passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
if err != nil {
    return User{}, fmt.Errorf("hash password: %w", err)
}
var correct bool = bcrypt.CompareHashAndPassword(passwordHash, []byte(password)) == nil
```

RBAC (Role-Based Access Control) - group permissions into Roles, assign Roles to Users; instead of by-person permissions

---

# Testing

```bash
go test ./...
go test -race ./...
go test -cover ./...
```

Test files:

```text
foo_test.go
```

---

## Same package vs external `_test` package

Given:

```go
package users
```

tests can use:

```go
package users
```

or:

```go
package users_test
```

Go explicitly supports both models. An `_test` package is compiled as a separate package.

### Same package

```go
package users
```

Advantages:

* access to unexported identifiers
* useful for white-box testing complicated internals
* convenient for low-level implementation tests

Disadvantages:

* tests can become coupled to implementation details
* refactoring internals may break many tests even though public behavior is unchanged
* tests may accidentally validate things callers cannot observe

---

### External test package

```go
package users_test
```

Advantages:

* exercises package as a real consumer
* only exported API is visible
* encourages good API boundaries
* detects missing or awkward public abstractions
* reduces coupling to implementation

Disadvantages:

* cannot directly test unexported helpers
* test setup can occasionally be less convenient

A good practical pattern:

> Prefer external-package tests for public behavior; use same-package tests where inspecting important internals genuinely adds value.

You can use both in the same package's test suite.

---

# In-memory implementations in tests

Yes: an in-memory implementation can be completely idiomatic.

Example production dependency:

```go
type UserStore interface {
    Find(ctx context.Context, id ID) (User, error)
    Save(ctx context.Context, user User) error
}
```

Production:

```go
postgres.Store
```

Tests:

```go
memory.Store
```

This is useful for **service/domain tests** where the database itself is not what you are testing.

Example:

```go
store := memory.NewStore()
service := users.NewService(store)
```

Benefits:

* fast
* deterministic
* easy setup
* easy error injection when designed appropriately
* no external infrastructure

But there is an important trap.

## An in-memory store is not a database emulator

A map cannot faithfully reproduce:

* SQL constraints
* transaction isolation
* locking behavior
* database-generated defaults
* collation
* NULL behavior
* query semantics
* indexes
* migrations
* serialization behavior
* connection failures

Therefore tests should normally be layered.

### Domain/service tests

Use:

```text
fake / stub / in-memory implementation
```

for fast behavioral tests.

### Repository/infrastructure tests

Test the actual database adapter against the real database engine:

```text
PostgreSQL implementation
        ↓
real PostgreSQL test instance
```

Often via disposable containers or isolated test databases.

### End-to-end tests

Exercise:

```text
HTTP → service → repository → real DB
```

for a smaller number of critical flows.

### Best rule

> Fake the abstraction you own; integration-test the infrastructure semantics you depend on.

Do not write an elaborate “fake PostgreSQL” in memory. At that point you are implementing another database whose behavior may diverge from production.

---

# Test doubles

Distinguish:

### Stub

Returns predetermined data.

```go
type stubStore struct {
    user User
}
```

### Fake

Working simplified implementation.

```go
type memoryStore struct {
    users map[ID]User
}
```

### Spy

Records calls for assertions.

```go
type spyMailer struct {
    sent []Message
}
```

### Mock

Usually verifies expected interactions.

In Go, simple hand-written test doubles are often preferable to large mocking frameworks because small consumer-owned interfaces make them trivial.

---

# Table-driven tests

Idiomatic structure:

```go
tests := []struct {
    name string
    in   string
    want string
}{
    {
        name: "empty",
        in:   "",
        want: "",
    },
    {
        name: "normal",
        in:   "foo",
        want: "FOO",
    },
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        got := Transform(tt.in)
        if got != tt.want {
            t.Fatalf("Transform(%q) = %q; want %q",
                tt.in, got, tt.want)
        }
    })
}
```

Also know:

```go
t.Helper()
t.Cleanup()
t.TempDir()
t.Setenv()
```

---

# Fuzzing

Built into Go testing:

```go
func FuzzParse(f *testing.F) {
    f.Add("example")

    f.Fuzz(func(t *testing.T, s string) {
        ...
    })
}
```

Useful for:

* parsers
* serializers
* validation
* protocol handling
* functions with broad input spaces

---

# Database

- `database/sql` as a stable standard-library abstraction for migration tooling
- using `pgx` through `database/sql` and `goose`

## sqlc
`sqlc` generates fully type-safe idiomatic code from SQL:

1. **Reads migration files:** parses migration/schema SQL to understand tables, columns, constraints and types.
    - `sqlc.yaml`.`schema` points sqlc to migrations directory
    - sqlc reads schema files as plain SQL (it has no knowledge of particular migration tools and ignores `goose` comments)
2. **Reads query files:** parses query SQL.
3. **Queries/Schema vidation:** checks that queries are valid against the schema.
4. **Infers Go types:** Infers go parameter and return types.
5. **Generates Go code** with typed methods.

Now application code calls generated methods.

Run `sqlc generate` when:
- migration/schema SQL changes
- query SQL changes
- sqlc config changes
- generated code is missing/stale
- before committing persistence changes

---

# Architecture Decision Record — ADR

An ADR records an important architecture decision.

Typical structure:

```text
Title
Status
Context
Decision
Consequences
Alternatives considered
```

Example:

```text
ADR-007: Use PostgreSQL for durable job storage
```

ADRs record **why** a decision was made, not just what the architecture currently looks like.

---

# Important senior-level concepts to add

## Ownership

For every resource ask:

```text
Who creates it?
Who owns it?
Who may mutate it?
Who closes it?
When does its lifetime end?
```

Especially important for:

* goroutines
* channels
* contexts
* DB transactions
* files
* network connections

---

## API compatibility

Exported identifiers are commitments.

Before exporting:

```go
type Foo ...
func Bar(...)
```

ask whether callers genuinely need them.

It is much easier to export something later than to remove it after users depend on it.

---

## Functional options

Useful when constructors accumulate optional configuration:

```go
client := NewClient(
    WithTimeout(time.Second),
    WithLogger(logger),
)
```

Avoid using the pattern when a plain configuration struct is clearer.

---

## Generics

Use type parameters when the algorithm genuinely operates over a family of types:

```go
func Contains[T comparable](xs []T, x T) bool
```

Do not replace ordinary interfaces with generics automatically.

Rough distinction:

```text
interface -> runtime behavioral abstraction
generic   -> compile-time type abstraction
```

---

## `any`

```go
any
```

is an alias for:

```go
interface{}
```

Use it where arbitrary values are genuinely necessary, not simply to avoid designing types.

---

## `comparable`

A generic constraint for types supporting:

```go
==
!=
```

Required for generic map keys.

---

## `select`

```go
select {
case v := <-ch:
    ...
case <-ctx.Done():
    return ctx.Err()
}
```

Core primitive for coordinating channel operations and cancellation.

A `select` with no ready case blocks.

`default` makes it non-blocking:

```go
select {
case v := <-ch:
    ...
default:
}
```

Use carefully; polling loops with `default` can burn CPU.

---

## Mutex copying

Do not copy a value containing an actively used:

```go
sync.Mutex
sync.RWMutex
sync.Once
```

Typically use a pointer receiver/type after first use.

`go vet` can detect many accidental lock copies.

---

## Race vs deadlock

Different bugs:

* **data race**: unsynchronized conflicting memory accesses
* **deadlock**: goroutines cannot make further progress
* **logical race**: timing-dependent behavior even when every individual memory access is synchronized

The race detector finds data races, not every concurrency bug.

---

## HTTP server timeouts

Production HTTP servers should deliberately consider:

```go
ReadHeaderTimeout
ReadTimeout
WriteTimeout
IdleTimeout
```

rather than blindly relying on zero/default values.

Timeout strategy must account for streaming/long-lived requests.

---

## Transactions

Keep transaction boundaries near the use case that requires atomicity.

Conceptually:

```go
err := db.WithTx(ctx, func(tx Tx) error {
    ...
})
```

Avoid letting unrelated lower-level repository methods independently choose transaction boundaries when several operations must commit atomically.

---

## Observability

Three primary signals:

```text
logs
metrics
traces
```

Prefer structured logging:

```go
logger.Info("user created",
    "user_id", user.ID,
)
```

Propagate trace/request metadata through context rather than global mutable state.

---

## Tooling

Core commands worth knowing:

```bash
go fmt ./...
go vet ./...
go test ./...
go test -race ./...
go test -cover ./...
go test -bench=. -benchmem
go mod tidy
go list ./...
go env
go doc
```

Common ecosystem tooling:

```text
gopls
staticcheck
golangci-lint
pprof
```

`gofmt` is the universal baseline; Go's review guidance notes that essentially all Go code uses it.

---

# Compact mental models

### Interface

```text
interface value ≈ (dynamic type, dynamic value)

(nil, nil)          -> interface == nil
(*Foo, nil)         -> interface != nil
```

### Slice

```text
slice ≈ (backing-array pointer, len, cap)

copy slice          -> copy descriptor
subslice            -> usually shared backing array
append              -> may replace backing array
```

### Map

```text
map variable -> runtime map structure

copy map value      -> still references same map data
missing lookup      -> zero value
nil map read        -> OK
nil map write       -> panic
iteration order     -> unspecified
```

### Context

```text
caller lifetime
      ↓
context
      ↓
handler
      ↓
service
      ↓
repository
      ↓
DB/API
```

### Interfaces and concrete types

```text
consumer
   |
   | defines minimal interface it needs
   v
interface
   ^
   | implicit satisfaction
   |
concrete implementation
```

### Testing

```text
domain/service
    -> fake/in-memory dependency

database adapter
    -> real database integration test

critical application flow
    -> end-to-end test
```

### Value vs pointer

```text
value:
    small
    copyable
    identity unimportant
    immutable/value semantics

pointer:
    mutation
    identity
    synchronization
    expensive/invalid copying
    shared state
```
