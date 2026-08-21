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

