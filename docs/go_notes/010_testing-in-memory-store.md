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
