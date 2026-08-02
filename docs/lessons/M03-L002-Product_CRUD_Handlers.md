Lesson 3.2 is complete. Let me provide the summary per the lesson template.

## L3.2 Complete — Summary

**Files created:**
- `internal/products/store.go` — `Store` interface + `InMemoryStore` with `Search` (using `strings.ToLower`/`Contains`)
- `internal/products/handler.go` — 5 handlers: Create, List, Update, Deactivate, Search (uses `c.QueryParam`)
- `internal/products/handler_test.go` — 12 tests covering create, duplicate, validation, list, update, deactivate, search by name, search by category, no results, empty query

**Files modified:**
- `internal/server/server.go` — accepts `*products.Handler`, registers 5 product routes
- `cmd/server/main.go` — wires `prodStore` + `prodHandler`
- `internal/server/server_test.go` — updated for new `New` signature
- `internal/employees/employee.go` — fixed missing doc comment

**Go concepts introduced:**
- **strings** (`strings.ToLower`, `strings.Contains`) — idiomatic string operations; the `strings` package is one of Go's most-used stdlib packages
- **QueryParam** — accessing URL query parameters via Fuego's `c.QueryParam("q")`; Go's `net/http` exposes this via `r.URL.Query().Get("q")` under the hood
- **Store interface extension** — adding `Search` to the existing CRUD contract; consumers define what they need

**Idioms demonstrated:**
- Consumer-owned interfaces stay in the feature package, not in a shared "repositories" package
- DTOs kept separate from domain model (CreateProductRequest/UpdateProductRequest vs Product)
- nil-safe slices in Search results (empty list returns `[]Product{}`, not `nil`)
- Error wrapping with `fmt.Errorf("...: %w", err)` preserving the sentinel error chain

**Common search gotcha:** The search route `/products/search` must be registered **before** `/products/{sku}` in Fuego, otherwise "search" matches the `{sku}` wildcard.
