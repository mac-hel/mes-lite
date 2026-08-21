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
