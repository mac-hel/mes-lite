# Project Structure and Development Workflow

## Modules, packages, layout

**module** = versioned collection of Go packages
* defined by `go.mod`
* reproducible builds via `go.sum`
* "module path" identifies the module and allows its import
* dependency management uses **Minimal Version Selection (MVS)** to decide which module version will be compiled into build (in contrast to constraint-solving)
    * lowest version that satisfies all explicitly required versions (explicit `go get p@v1.6.0`, moudle A requires C 1.4.0, B requires 1.7.0 -> 1.7.0 is compiled)
    * no version ranges like `>=, ^, ~`, dependency always define minimal version
    * major version `example.com/lib/v2` is introduced for breaking changes - only way to import and compile different versions of same package
* `go get` changes dependencies / versions
* `go mod tidy` synchronizes `go.mod` / `go.sum` with imports
* `go work` creates a workspace containing multiple local modules

**package** = compilation / encapsulation unit
* all ordinary `.go` files in one directory normally belong to the same package
* exported identifiers start with uppercase letters
* package dependencies must form an acyclic graph

**`package main`**
* executable program
* must contain `func main()` - usually composition root
* multiple binaries commonly live under:

    ```text
    cmd/
      api/
        main.go
      worker/
        main.go
    ```

**`internal/`**
  * compiler/toolchain-enforced **import boundary**
  * code outside the parent tree of an `internal` directory cannot import it
  * useful for implementation details that must not become public APIs
  * stronger than naming conventions, but not equivalent to “module-private” in every possible layout

**`pkg/`**
  * optional community convention (**not special to the Go toolchain**)
  * often unnecessary
  * public packages can live directly at repository/module root

**organize packages** around **cohesive responsibilities**, not technical layer names by default
  * e.g. `orders`, `payments`, `users`
  * a domain package may own:
    * domain types
    * behavior/use cases
    * interfaces required by that behavior
    * tests
  * HTTP/database implementations can either live there or in adjacent packages depending on coupling and size

Go recommends short, clear, lowercase, usually single-word package names and discourages vague buckets such as `util`, `common`, `types`, and `interfaces`.

---

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

## Tooling

Build:
```bash
```

Core commands worth knowing:

```bash
go build ./cmd/server               # ~> ./server
go build -o out/mes ./cmd/server    # ~> out/mes
go fmt ./...                    # universal baseline; essentially all Go code uses it
go vet ./...
go test ./...
go test -race ./...
go test -cover ./...
go test -bench=. -benchmem
go mod tidy                     # operations on modules, also: edit, graph, vendor, why
go list ./...
go env
go doc cmd/server myFunc        # show docs for: package, type, func, var etc.
```

Common ecosystem tooling:

```text
gopls
staticcheck
golangci-lint
pprof
```

---

## Documentation

```go
// Package main is...
package main
```

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
