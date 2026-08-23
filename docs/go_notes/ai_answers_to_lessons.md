## 1. Why does Go avoid inheritance?

Go favors **composition + interfaces** because inheritance often couples behavior to a rigid type hierarchy.

```go
type Service struct {
    Logger
}
```

You reuse implementation through composition/embedding and abstract behavior through interfaces. This reduces fragile base-class problems, deep hierarchies, implicit behavior, and “is-a” relationships that don't actually fit the domain.

The model is:

> **Reuse implementation with composition; model substitutability with interfaces.**

---

## 2. Why prefer `time.Time` over strings for timestamps?

`time.Time` gives the compiler and standard library semantic information that a string does not:

```go
CreatedAt time.Time
```

instead of:

```go
CreatedAt string
```

You get parsing once at the boundary, comparisons, arithmetic, timezone handling, validation, formatting, and methods such as `Before`, `After`, `Add`, and `Sub`.

With strings, every consumer must ask: *What format? Is it valid? Which timezone?*

Use strings for the serialized representation; use `time.Time` inside Go APIs.

---

## 3. Why RFC3339 at HTTP boundaries?

RFC3339 gives clients an interoperable, standardized and unambiguous timestamp representation:

```text
2026-08-23T17:30:00Z
2026-08-23T19:30:00+02:00
```

It includes date, time, and timezone/UTC offset and is widely supported across languages.

A good boundary is therefore:

```text
HTTP JSON          Go
RFC3339 string  ↔  time.Time
```

For APIs requiring fractional precision, RFC3339-compatible fractional seconds are also commonly used.

---

## 4. Why pass `context.Context` from handler to persistence?

Because the database operation is part of the request's lifetime.

```text
HTTP request
    ↓ ctx
application
    ↓ ctx
repository
    ↓ ctx
PostgreSQL
```

If the client disconnects, a deadline expires, or the request is cancelled, downstream work should ideally stop rather than consuming DB connections and CPU for a response nobody needs.

It also propagates request-scoped tracing/deadline information.

---

## 5. Why is `context.Context` the first argument?

It's a Go convention:

```go
func Find(ctx context.Context, id ID) (...)
```

Context controls the **lifetime of the whole operation**, rather than being business data. Putting it first makes it immediately visible and consistent across APIs and makes propagation straightforward:

```go
handler(ctx, ...)
service(ctx, ...)
repo(ctx, ...)
```

Contexts should normally be passed explicitly rather than stored in structs.

---

## 6. Why shouldn't the handler directly return application/persistence structs?

Because it couples your external HTTP contract to internal implementation.

Returning a sqlc/database struct directly can accidentally expose fields, PostgreSQL-specific types, internal state, or future schema changes.

Prefer:

```text
database row
     ↓
domain/application result
     ↓
HTTP response DTO
     ↓
JSON
```

For example:

```go
type EmployeeResponse struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}
```

Now the database schema, domain model, and API can evolve independently.

---

## 7. Why translate infrastructure errors into domain/application errors?

Callers shouldn't need to understand PostgreSQL or pgx internals.

Instead of the HTTP layer asking:

```go
errors.Is(err, pgx.ErrNoRows)
```

the repository can expose something meaningful:

```go
ErrEmployeeNotFound
```

Then:

```text
pgx.ErrNoRows
    ↓
ErrEmployeeNotFound
    ↓
HTTP 404
```

This preserves abstraction and makes infrastructure replaceable.

Important caveat: **not every DB failure becomes a domain error**. A network failure should remain an infrastructure/internal error, perhaps wrapped with context, rather than being misrepresented as a business condition.

---

# I/O

## 8. Why is `io.Reader` so important?

Because it abstracts **where bytes come from**:

```go
type Reader interface {
    Read([]byte) (int, error)
}
```

The same consumer can read from files, HTTP bodies, sockets, gzip streams, strings, buffers, etc.

```go
func Parse(r io.Reader) error
```

doesn't care whether `r` is a file or network connection.

It enables Go's powerful streaming composition:

```text
HTTP body
    ↓
gzip.Reader
    ↓
csv.Reader
    ↓
application
```

with tiny interfaces and low coupling.

---

## 9. Why isn't `io.EOF` considered a failure by stream consumers?

Technically `io.EOF` is an `error` value, but semantically it means:

> There is simply no more data.

So:

```go
for {
    n, err := r.Read(buf)

    // process n bytes

    if errors.Is(err, io.EOF) {
        break
    }
    if err != nil {
        return err
    }
}
```

End-of-stream is an expected termination condition, unlike disk failure, broken connection, malformed input, etc.

Also remember that a reader can return both `n > 0` and an error; process the bytes first.

---

# sqlc and SQL

## 10. Why choose sqlc instead of an ORM?

sqlc gives you:

> **SQL that you control + generated type-safe Go code.**

You write:

```sql
SELECT id, name
FROM employees
WHERE department_id = $1;
```

and sqlc generates the Go method and result types. ([sqlc Documentation][1])

Advantages include predictable queries, good access to PostgreSQL features, less runtime magic, easier query tuning, and compiler-visible Go types.

An ORM can still be reasonable when rapid CRUD development and object mapping matter more than precise SQL control.

---

## 11. Why explicit SQL?

Because the database executes SQL anyway.

Explicit SQL lets you directly reason about joins, indexes, locking, query plans, selected columns, round trips, aggregation, and PostgreSQL-specific features.

Compare:

```text
application intent
      ↓
explicit SQL
      ↓
PostgreSQL
```

with an ORM where you may additionally need to understand:

```text
ORM expression
      ↓
generated SQL
      ↓
actual query plan
```

Explicit SQL removes one abstraction layer from performance-critical behavior.

---

## 12. Why shouldn't generated sqlc code become the business package's public API?

Because sqlc types represent the **persistence model**, not necessarily your domain.

If business code directly exposes:

```go
db.Employee
db.CreateEmployeeParams
```

then the business layer becomes coupled to schema layout, generated names, nullable database types, pgx/sqlc details, and regeneration changes.

Prefer an adapter:

```text
business
    │ OrderRepository interface
    ▼
postgres adapter
    │
    ▼
sqlc generated code
```

sqlc is an implementation detail of persistence, not your business vocabulary.

---

# Database integrity and migrations

## 13. What is the role of migrations in production?

Migrations provide a **version-controlled history of database schema changes**:

```text
001_create_employee.sql
002_add_department.sql
003_add_index.sql
```

They let deployments reproducibly evolve tables, indexes, constraints, types, and sometimes data.

Production migrations also require thinking about backward compatibility, locking, long-running changes, rolling deployments, data backfills, and rollback/roll-forward strategy.

They are essentially deployment code for database state.

---

## 14. What happens if the request context is cancelled while pgx waits on a query?

pgx's blocking operations accept contexts. In current pgx/pgconn, the default behavior on context cancellation is for the operation to return promptly; in most circumstances the underlying connection is also closed. pgx provides configurable cancellation strategies, including PostgreSQL cancellation requests, if different behavior is desired. ([Go Packages][2])

With a connection pool, a closed/broken connection is not returned as a healthy reusable connection; the pool can establish replacements as necessary.

One related subtlety: starting a pgx transaction with a context does **not** mean cancellation automatically rolls that transaction back; pgx documents that the context on `Begin` affects the begin command itself. ([Go Packages][3])

---

## 15. Why application validation if foreign keys exist?

Because application validation provides **domain semantics and useful errors**.

Instead of letting:

```text
INSERT → FK violation
```

be your primary business validation mechanism, the application can say:

```text
"department does not exist"
"employee cannot be assigned to archived department"
```

It also lets you evaluate rules involving more than relational integrity.

But application checks alone cannot guarantee correctness because concurrency creates races.

---

## 16. Why foreign keys if the application already validates?

Because the database is the final integrity boundary.

Between:

```text
application checks department exists
```

and:

```text
application inserts employee
```

another transaction could delete that department.

And data may also be modified by bugs, scripts, admin tools, migrations, or another service.

So:

```text
Application validation
    → good semantics / domain rules

Database constraints
    → final invariant enforcement
```

The duplication is intentional defense in depth.

---

# Transactions

## 17. What should define a transaction boundary?

A **business operation that must be atomic**.

For example:

```text
PlaceOrder:
  create order
  reserve inventory
  record payment state

all succeed
OR
all roll back
```

The boundary should follow the use-case/invariant consistency requirement—not arbitrarily one repository method, table, or HTTP handler.

Keep transactions as short as practical and avoid slow external network calls inside them where possible.

---

## 18. Why does deferred rollback run after successful commit, and why is that safe?

Typical code:

```go
tx, err := db.Begin(ctx)
if err != nil {
    return err
}

defer tx.Rollback(ctx)

// operations...

if err := tx.Commit(ctx); err != nil {
    return err
}

return nil
```

`defer` executes when the function returns, so it still calls `Rollback` after `Commit`.

But the transaction is already closed. pgx explicitly documents that calling rollback on a closed transaction is safe; it returns an `ErrTxClosed`-compatible error and does nothing to the already committed transaction. This is why pgx itself demonstrates the deferred-rollback pattern. ([Go Packages][3])

Its purpose is to guarantee rollback on **every earlier error path** without duplicating cleanup code.

---

# Pagination

## 19. When is cursor pagination better than limit/offset?

Cursor pagination is better for large or frequently changing datasets.

Offset:

```sql
ORDER BY created_at
LIMIT 50 OFFSET 1000000
```

may require PostgreSQL to walk past a large number of rows, and inserts/deletes can cause duplicates or skipped records between pages.

Cursor/keyset pagination instead says roughly:

```sql
WHERE (created_at, id) > ($1, $2)
ORDER BY created_at, id
LIMIT 50
```

Advantages are better scalability and more stable traversal during concurrent changes.

Tradeoff: you normally cannot jump directly to “page 237.” You also need a deterministic, preferably indexed ordering with a unique tie-breaker.

---

# Validation

## 20. What is layered validation responsibility?

Different layers validate different concerns:

```text
HTTP/transport
  → JSON shape, required fields, syntax

Application
  → use-case prerequisites, authorization,
    orchestration rules

Domain
  → business invariants

Database
  → structural/integrity constraints
    such as UNIQUE, FK, NOT NULL
```

Example email:

```text
HTTP: value is syntactically present
Domain: Email value object must be valid
DB: unique constraint prevents duplicates
```

Some overlap is appropriate because the layers protect against different failure modes.

---

# Authentication and authorization

## 21. Why separate authentication and authorization?

They answer different questions:

**Authentication:**

> Who are you?

**Authorization:**

> Are you allowed to perform this operation on this resource?

A valid logged-in user may still not be allowed to:

```text
DELETE /employees/123
```

Keeping them separate makes policies easier to change, reason about, audit, and test.

---

## 22. What is RBAC, and what are other authorization patterns?

**RBAC — Role-Based Access Control**

```text
Alice → Manager → employee:edit
Bob   → Viewer  → employee:read
```

Permissions derive from roles.

Other common models:

* **ABAC — Attribute-Based Access Control:** decisions depend on attributes such as department, resource state, geography, time, etc.
* **ACL — Access Control List:** each resource lists users/groups and their permissions.
* **ReBAC — Relationship-Based Access Control:** permission depends on relationships, e.g. “owner of project,” “member of team that owns document.”
* **Capability-based authorization:** possession of a specific unforgeable capability grants an operation.

Real systems often combine them—for example RBAC for broad permissions plus resource ownership checks.

---

# Passwords and identities

## 23. Why store password hashes rather than plaintext? And why bytes?

The critical requirement is **never store recoverable plaintext passwords**.

Store a slow, salted password hash produced by an appropriate password-hashing algorithm such as Argon2id, bcrypt, or scrypt.

```text
password
   ↓
password KDF + salt
   ↓
stored hash
```

If the database leaks, attackers don't immediately obtain users' original passwords.

In Go, many cryptographic/password APIs naturally operate on `[]byte`, and bytes are mutable while strings are immutable. But **`[]byte` itself provides no cryptographic protection**. A password hash can validly be stored as a textual encoded hash—e.g. a PHC-style string. The security property comes from the hashing/KDF scheme and parameters, not whether the DB column is `bytea` or `text`.

---

## 24. Why are an Employee and an Auth User different domain concepts?

Because they have different identity, lifecycle and responsibilities.

```text
Employee
  name
  department
  salary
  manager
  employment status

AuthUser
  login identity
  password/credential
  MFA
  authentication state
  roles
```

An employee may have no login. An auth user might represent an administrator, contractor, customer, or machine account rather than an employee.

If you force them into one entity, authentication concerns become coupled to HR/business concerns.

They may be related:

```text
AuthUser ── optionally references ── Employee
```

without being the same concept.

---

# JWT

## 25. What problem does a JWT solve?

A JWT is a standardized token format for carrying claims between parties, commonly with cryptographic integrity protection.

For example:

```text
client
  │ Authorization: Bearer <JWT>
  ▼
API
  │ verifies signature
  ▼
claims:
  sub = user123
  role = manager
  exp = ...
```

A server can verify a signed JWT without looking up a session on every request.

### Advantages

* portable standardized format
* decentralized/offline signature verification
* useful across multiple services
* expiration and standard claims
* works well with OAuth/OIDC ecosystems

### Disadvantages

* revocation is harder than server-side sessions
* claims can become stale
* stolen bearer tokens can be replayed until invalid/expired
* token/key rotation requires care
* larger than opaque session IDs
* encourages people to put too much data into tokens
* easy to implement insecurely

JWT is a **token format**, not an authentication system by itself.

Alternatives include opaque session cookies with server-side sessions, opaque access tokens, API keys for suitable machine use cases, mTLS/client certificates, and other signed-token formats. OAuth/OIDC can also issue either JWT or opaque tokens.

---

## 26. Signing vs encrypting a JWT

A **signed JWT** usually means JWS:

```text
payload + signature
```

It guarantees that the contents haven't been modified and, assuming key trust, who signed them.

But the payload is readable:

```text
base64url(header).base64url(payload).signature
```

Base64URL is encoding, **not encryption**.

An **encrypted JWT**, normally using JWE, protects confidentiality as well as integrity when correctly constructed.

So:

```text
signed    → trustworthy, but readable
encrypted → confidential to intended recipient(s)
```

Most access-token JWTs are signed, not encrypted.

---

## 27. Why must JWT payloads not contain password hashes or sensitive data?

Because ordinary signed JWT payloads are **not secret**. Anyone possessing the token can decode the claims.

Tokens also tend to pass through clients, browser/storage mechanisms, gateways, proxies, logs, monitoring systems, and debugging tools.

Putting password hashes inside one would create an unnecessary offline cracking target and spread highly sensitive credential material through many systems.

JWT claims should follow data minimization:

```text
subject ID
issuer
audience
expiry
necessary authorization claims
```

not:

```text
password hash
private secrets
sensitive personal records
```

---

# Startup

## 28. Why is idempotent startup important?

Production processes restart frequently:

```text
deployment
crash
autoscaling
node replacement
orchestrator restart
```

Startup should therefore be safe to run repeatedly.

Bad:

```text
startup #1 → create admin
startup #2 → duplicate admin
startup #3 → crash because admin exists
```

Better:

```text
ensure required state exists
```

such that:

```text
run once   → correct state
run 20x    → same correct state
```

This matters for schema initialization, seed/reference data, resource registration, indexes, subscriptions, job setup, etc.

For distributed systems, idempotency must also survive **multiple instances starting concurrently**. Database unique constraints, transactional operations, advisory locks, or migration tooling are usually stronger than a naïve:

```text
if not exists → create
```

check.

The underlying principle is:

> **Restarting or retrying initialization should converge on the same valid state rather than accumulating side effects.**

[1]: https://docs.sqlc.dev/en/v1.26.0/tutorials/getting-started-postgresql.html?utm_source=chatgpt.com "Getting started with PostgreSQL — sqlc 1.26.0 documentation"
[2]: https://pkg.go.dev/github.com/jackc/pgx/v5/pgconn "pgconn package - github.com/jackc/pgx/v5/pgconn - Go Packages"
[3]: https://pkg.go.dev/github.com/jackc/pgx/v5 "pgx package - github.com/jackc/pgx/v5 - Go Packages"
