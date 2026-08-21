# Concurrency

## Goroutine

```go
go work()
```

A goroutine is a lightweight concurrent execution unit managed by the Go runtime.

It is **not literally an OS thread**.

Many goroutines are multiplexed onto OS threads by the runtime scheduler.

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

# Synchronization and concurrency ownership

## Race vs deadlock

Different bugs:

* **data race**: unsynchronized conflicting memory accesses
* **deadlock**: goroutines cannot make further progress
* **logical race**: timing-dependent behavior even when every individual memory access is synchronized

The race detector finds data races, not every concurrency bug.

## Channel vs Mutex

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

## Mutex copying

Do not copy a value containing an actively used:

```go
sync.Mutex
sync.RWMutex
sync.Once
```

Typically use a pointer receiver/type after first use.

`go vet` can detect many accidental lock copies.
