**Lesson 2.2 — Constructors, Slices & Interfaces: Creating & Listing Employees**

**Files created:**
- `internal/employees/store.go:5-8` — consumer-owned `Store` interface with `Save`/`List`
- `internal/employees/store.go:16-19` — `InMemoryStore` backed by `[]Employee` (slice)
- `internal/employees/handler.go:6-8` — `Handler` struct accepting the `Store` interface
- `internal/employees/handler.go:23-37` — `Create` handler: POST /employees
- `internal/employees/handler.go:44-55` — `List` handler: GET /employees

**Files modified:**
- `internal/server/server.go:27` — `New()` accepts `*employees.Handler`, registers routes
- `cmd/server/main.go:19-22` — composition root wires `InMemoryStore` → `Handler` → `Server`

**Go concepts learned:**
1. **Constructors** — `NewInMemoryStore()`, `NewHandler(store Store)` encapsulate initialization and enforce invariants. The composition root in `main.go` wires them together via constructor injection — no globals, no DI containers.
2. **Slices** — `[]Employee` backs the in-memory store. `append()` for Save, range-based iteration in tests, nil-vs-empty slice handling in `List` (`if emps == nil { emps = []Employee{} }` ensures JSON outputs `[]` not `null`).
3. **Interfaces** — The `Store` interface is defined in the `employees` package where it's consumed, not where it's implemented. `InMemoryStore` implicitly satisfies it — no `implements` keyword (Go's structural typing). The handler accepts `Store` (interface), `NewInMemoryStore` returns `*InMemoryStore` (concrete) — "accept interfaces, return structs."
4. **DTO separation** — `CreateEmployeeRequest` is the wire format (JSON tags with camelCase), separate from the domain `Employee` struct. `ListEmployeesResponse` wraps the slice for clean API shape.

**Architecture:**
- Package cohesion — `employees` owns its domain, store interface, HTTP handlers, DTOs, and tests
- Constructor DI — `main.go` wires dependencies explicitly, no magic
- Consumer-owned interfaces — `Store` lives with its consumer, not the implementation

**Key design decisions:**
- `List` returns `[]Employee{}` instead of `nil` to produce `[]` in JSON (better UX for clients)
- The `Store` uses `context.Context` even though the in-memory implementation ignores it — establishes the contract for future PostgreSQL-backed implementations
- No `sync.Mutex` yet — the in-memory store is not concurrency-safe. This will be addressed in Milestone 10 (Concurrency)

**Exercises:**
1. What happens to the JSON output if you remove the `if emps == nil` check in `List`?
2. Add a third employee to `TestHandler_List_AfterCreate`. What order do they appear in? Why?
3. Why does `NewHandler(store Store)` accept an interface but `NewInMemoryStore()` returns a concrete type?
4. Write an `InMemoryStore` that uses a `map[string]Employee` instead of a slice. What's the trade-off?

**Interview questions:**
1. **"Accept interfaces, return structs"** — why is this idiomatic in Go?
2. **Nil slice vs empty slice** — what's the difference in JSON marshal behavior? Why is it important for APIs?
3. **Consumer-owned interfaces** — why should the package that uses an interface define it, rather than the package that implements it?
4. **Why `context.Context` as the first parameter?** — convention, cancellation propagation, deadlines, and request-scoped values.
