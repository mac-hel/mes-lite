**Lesson 2.3 — Error Wrapping, Validation & Maps: Updating, Deactivating & Testing**

**Milestone 2 is now complete.**

**Go concepts learned:**
1. **Error Wrapping** — `fmt.Errorf("employee %q: %w", id, ErrNotFound)` wraps the sentinel error. Handlers use `errors.Is(err, ErrNotFound)` to detect the underlying cause without depending on string comparisons.
2. **Validation** — Struct tags `validate:"required,email"` declaratively define field constraints. Fuego automatically runs validation after deserializing the request body and returns 400 with descriptive messages.
3. **Maps** — `InMemoryStore` refactored from `[]Employee` to `map[string]Employee`. Key benefits: O(1) lookup by ID via comma-ok pattern (`emp, ok := s.employees[id]`). Iteration uses `for _, emp := range`. Nil maps vs `make()` — maps must be initialized before writing.
4. **Fuego HTTP errors** — `fuego.NotFoundError{Err: err}` wraps the domain error while signaling the framework to return 404 instead of 500. Demonstrates the `ErrorWithStatus` interface pattern.

**Architecture decisions:**
- Error translation layer — store returns domain errors (`ErrNotFound`), handler translates to HTTP-appropriate errors (`fuego.NotFoundError`). Keeps persistence concerns out of the domain.

**Tests added:**
- `TestHandler_Create_Validation` — table-driven test with 5 validation cases (400 for missing/invalid fields)
- `TestHandler_Update` — full update + verify changed fields
- `TestHandler_Update_NotFound` — 404 for non-existent employee
- `TestHandler_Deactivate` — soft delete sets `IsActive=false`
- `TestHandler_Deactivate_NotFound` — 404 for deactivating non-existent employee
- `TestHandler_List_AfterCreate` — fixed for non-deterministic map iteration order

**Exercises:**
1. Remove the `if errors.Is(err, ErrNotFound)` check in `Update` — what HTTP code does Fuego return now?
2. Why does `map[string]Employee` lose insertion order? How would you preserve order?
3. What happens if you call `Save` on a nil map (instead of one created by `make`)?
4. Add a `TestHandler_DuplicateCreate` — what should happen when you POST the same employee ID twice?

**Interview questions:**
1. **Error wrapping** — What does `%w` do vs `%v` in `fmt.Errorf`? How does `errors.Is` traverse the error chain?
2. **Maps vs slices** — When would you choose a map over a slice for an in-memory store?
3. **Error translation** — Why translate domain errors to HTTP errors at the handler layer rather than the store layer? What does this say about dependency direction?
4. **Validation** — Why use struct tags for validation instead of writing `if` checks in the handler? What are the trade-offs?
