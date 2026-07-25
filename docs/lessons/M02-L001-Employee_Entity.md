**Lesson 2.1 — Visibility & Zero Values: Employee Entity**

**Go concepts learned:**
- **Visibility** — Exported types/functions (`Employee`, `NewEmployee`) are visible to other packages; unexported helpers (none yet) stay within the package
- **Zero values** — `var emp Employee` gives zero-valued fields: `""` for strings, `false` for `bool`; `NewEmployee` explicitly sets `IsActive: true`
    - every variable has default zero value without need to initialize it
- **Return structs** — if reasonably small, copy is inexpensive, sharing mutable state not needed
    - eliminates `nil` pointer check
    - avoids heap allocation
    - compiler optimizes copying
- **Return pointer** — if:
    - methods need to modify original instance
    - target must be shared
    - target contains synchronization properties (sharing)
    - target struct is large
- **Accept Interfaces, Return Structs**
    - accept interfaces: caller decides implementation, easier testing (mocks), lower coupling
    - return concrete types: exposes full API, avoids hiding methods

**Architecture concept:** Package cohesion — first business domain package (`internal/employees/`) that owns its type definition and construction logic

**Design decisions:**
- `Employee` is a struct (not a pointer) returned by value — follows "accept interfaces, return structs"; zero value is meaningful (all fields empty, inactive)
- No repository interface yet — "concrete first, abstract later"
- No `time.Time` fields — not yet introduced in the curriculum
- No unnecessary fields — just what's needed: identity, name, contact, active flag

**Future Improvements:**
- add invariants validation to `NewEmployee` to make sure it can't exist in invalid state
    - not validation of business rules, external systems or application workflow - just intrinsic to `Employee`
- convert `Employee's` `IsActive` field to state/status, e.g.: medical leave, vacation, business trip

**Exercises:**
1. What is the zero value of a `bool`? Of a `string`?
2. What happens if you declare `var e Employee` and access `e.ID`?
3. Why is there no validation in `NewEmployee`? (Hint: separation of concerns, single responsibility)
4. Move `employee_test.go` to a separate `employees_test` package (external test) — what breaks? Why?

**Interview questions:**
1. **Exported vs unexported** — How does Go control visibility compared to `public`/`private` in Java/PHP?
2. **Zero values** — Why is the zero value design important in Go? Give an example where it simplifies code.
3. **Return structs** — Why does `NewEmployee` return `Employee` (value) instead of `*Employee` (pointer)?
4. **Package naming** — Why is the package named `employees` and not `employee`? (Plural for domain packages, singular for utility)
