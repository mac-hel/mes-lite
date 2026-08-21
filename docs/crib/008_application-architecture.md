# Code architecture

## Ownership

For every resource ask:

```text
Who creates it?
Who owns it?
Who may mutate it?
Who closes it?
When does its lifetime end?
```

Especially important for:

* goroutines
* channels
* contexts
* DB transactions
* files
* network connections

## API compatibility

Exported identifiers are commitments.

Before exporting:

```go
type Foo ...
func Bar(...)
```

ask whether callers genuinely need them.

It is much easier to export something later than to remove it after users depend on it.

---

# DTOs and transport boundaries

Avoid letting HTTP/JSON representation accidentally define your core domain model.

Transport DTO:

```go
type createUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}
```

Domain:

```go
type User struct {
    Name  Name
    Email Email
}
```

Benefits:

* transport format can evolve separately
* validation can occur at the boundary
* domain does not accumulate JSON/database tags
* avoids over-posting fields accidentally
* API versions need not map 1:1 to domain structures

Do not create DTO layers mechanically when domain and transport structures genuinely have identical semantics and coupling is harmless.

---

# Validation

**transport**
validates only HTTP request shape, not domain rules, e.g. by using struct tags

**application Services**
validate workflow rules and references

**Domain**
- constructor `NewProduct` validates domain rules (invariants)
- mutation method `UpdateDetails` replaces mutable product fields and preserves product invariants
- `ReconstituteProduct` validates reconstructing data - should be private if possible
- both use common `validate` function but can adjust validation for own requirements

**repository**
validate before persistence as defense-in-depth

**PostgreSQL**
constraints provide the final integrity boundary

---

# Authentication / Authorization

Password hashing
```go
import "golang.org/x/crypto/bcrypt"
passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
if err != nil {
    return User{}, fmt.Errorf("hash password: %w", err)
}
var correct bool = bcrypt.CompareHashAndPassword(passwordHash, []byte(password)) == nil
```

RBAC (Role-Based Access Control) - group permissions into Roles, assign Roles to Users; instead of by-person permissions

---

# HTTP server timeouts

Production HTTP servers should deliberately consider:

```go
ReadHeaderTimeout
ReadTimeout
WriteTimeout
IdleTimeout
```

rather than blindly relying on zero/default values.

Timeout strategy must account for streaming/long-lived requests.

---

# Architecture Decision Record — ADR

An ADR records an important architecture decision.

Typical structure:

```text
Title
Status
Context
Decision
Consequences
Alternatives considered
```

Example:

```text
ADR-007: Use PostgreSQL for durable job storage
```

ADRs record **why** a decision was made, not just what the architecture currently looks like.
