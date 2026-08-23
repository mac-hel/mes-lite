# Various

Go tries to make **ownership, dependencies, control flow, and boundaries** visible in the code.
It has relatively few language mechanisms, but those mechanisms interact in important ways.

## Various

- variadic function parameters: `func (nums ...int)`
- `any` is an alias for: `interface{}`
    - Use it where arbitrary values are genuinely necessary, not simply to avoid designing types.
- `comparable` - A generic constraint for types supporting: `== !=`, Required for generic map keys.

## Blank imports

`import _ "some/package"`
Imports a package only for its side effects - you cannot reference the package directly.
The primary side effect is package initialization:
```go
func init() {
    // registration
}
```

Historically this pattern is common for plugins/drivers, e.g.:
```go
import (
    "database/sql"
    _ "github.com/example/database-driver"
)
sql.Open(...)   // now can do - without directly calling the driver package.
```
The driver registers itself with `database/sql` during initialization.

Blank imports can hide dependency relationships, so use them when a registration architecture intentionally requires them rather than as a general-purpose technique.

---

# Compact mental models

## Interface

```text
interface value ≈ (dynamic type, dynamic value)

(nil, nil)          -> interface == nil
(*Foo, nil)         -> interface != nil
```

## Slice

```text
slice ≈ (backing-array pointer, len, cap)

copy slice          -> copy descriptor
subslice            -> usually shared backing array
append              -> may replace backing array
```

## Map

```text
map variable -> runtime map structure

copy map value      -> still references same map data
missing lookup      -> zero value
nil map read        -> OK
nil map write       -> panic
iteration order     -> unspecified
```

## Context

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

## Interfaces and concrete types

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

## Testing

```text
domain/service
    -> fake/in-memory dependency

database adapter
    -> real database integration test

critical application flow
    -> end-to-end test
```

## Value vs pointer

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
