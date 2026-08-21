# Various

## Various

- blank imports
- variadic function parameters: `func (nums ...int)`
- `any` is an alias for: `interface{}`
    - Use it where arbitrary values are genuinely necessary, not simply to avoid designing types.
- `comparable` - A generic constraint for types supporting: `== !=`, Required for generic map keys.

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
