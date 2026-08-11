# Validation Flow

This document defines how MES Lite validates data from HTTP request to PostgreSQL.

The goal is not to validate everything everywhere. The goal is to put each rule at the layer that owns it, while keeping the database as the final integrity boundary.

## Principles

- Validate as early as useful, but enforce critical integrity as late as necessary.
- Domain entities should not be constructible in invalid business state when a constructor is available.
- HTTP validation is not domain validation.
- Database constraints are not a replacement for clear domain code.
- Repositories translate storage errors; they should not invent business rules.
- Every validation error should be understandable at the API boundary.

## Validation Layers

### HTTP Handler

Handlers validate transport concerns.

Responsibilities:

- JSON shape
- required request fields
- path and query parameters
- basic format checks such as email format
- conversion from request DTOs to commands or constructor arguments
- mapping domain/application errors to HTTP errors

Handlers should not own business rules such as employee activity, product availability or transaction consistency.

### Application Service

Application services validate workflow rules.

Responsibilities:

- reference checks across slices
- active/inactive business checks
- permission checks after authentication exists
- transaction boundaries
- coordination of multiple repositories
- mapping lower-level domain errors into operation-level errors when useful

Example: production registration checks whether the employee and product exist and can be used for production.

### Domain Entity And Value Object

Domain types enforce invariants.

Responsibilities:

- required business fields
- normalized values
- supported custom types
- valid state transitions
- constructor-level validation
- mutation methods that preserve invariants

Preferred constructor shape for entities with invariants:

```go
func NewProduct(sku, name, unit string, category ProductCategory) (Product, error) {
    p := Product{
        SKU:      strings.TrimSpace(sku),
        Name:     strings.TrimSpace(name),
        Unit:     strings.TrimSpace(unit),
        Category: category,
        IsActive: true,
    }

    if err := p.Validate(); err != nil {
        return Product{}, err
    }

    return p, nil
}
```

It is acceptable for business entities to have an invalid zero value. Go encourages useful zero values when practical, but persisted business entities often require explicit construction.

### Repository

Repositories validate defensively and translate storage errors.

Responsibilities:

- calling domain validation before persistence when appropriate
- adapting domain structs to generated SQL structs
- translating database errors into package/domain errors
- preserving `context.Context` propagation
- avoiding HTTP-specific errors

Repositories should not expose sqlc-generated row types as the public package API.

### Database

PostgreSQL is the final integrity boundary.

Responsibilities:

- protect against application bugs
- protect against concurrent writes
- protect against manual SQL changes
- protect against future services or import jobs bypassing HTTP
- enforce constraints that must always be true

Database constraints should exist for rules that would corrupt stored data if violated.

## PostgreSQL Integrity Capabilities

Use these capabilities when they match the rule being enforced.

### Data Types

Choose types that make invalid data unrepresentable or harder to store.

Examples:

- `integer` for quantities instead of `text`
- `boolean` for true/false state
- `uuid` for UUID identifiers
- `timestamptz` for timestamps with time zone semantics
- `numeric` for exact decimal values when needed
- bounded `varchar(n)` when a real business maximum exists

### NOT NULL Constraints

Use `NOT NULL` when a value is required for every row.

Example:

```sql
name text NOT NULL
```

### CHECK Constraints

Use `CHECK` for row-local predicates.

Examples:

```sql
quantity integer NOT NULL CHECK (quantity > 0)
workstation text NOT NULL CHECK (btrim(workstation) <> '')
category integer NOT NULL CHECK (category BETWEEN 0 AND 4)
```

### PRIMARY KEY Constraints

Use primary keys for stable row identity.

Examples:

```sql
id uuid PRIMARY KEY
sku text PRIMARY KEY
```

### UNIQUE Constraints

Use unique constraints for business uniqueness.

Examples:

```sql
email text NOT NULL UNIQUE
sku text NOT NULL UNIQUE
```

### Foreign Keys

Use foreign keys when one row references another row that must exist.

Examples:

```sql
employee_id text NOT NULL REFERENCES employees(id)
product_sku text NOT NULL REFERENCES products(sku)
```

Foreign keys can define behavior for referenced-row changes:

- `ON DELETE RESTRICT`
- `ON DELETE CASCADE`
- `ON DELETE SET NULL`
- `ON UPDATE CASCADE`

Use destructive cascading only when it matches the business. Production history usually should not disappear because master data was deleted.

### Deferrable Constraints

Use deferrable constraints when related rows are created or changed inside one transaction and the constraint should be checked at commit time.

Example:

```sql
FOREIGN KEY (employee_id) REFERENCES employees(id) DEFERRABLE INITIALLY DEFERRED
```

This is useful for complex transactional workflows, not for ordinary CRUD by default.

### Indexes

Indexes primarily improve reads, but unique indexes also enforce integrity.

Useful forms:

- regular indexes for lookup performance
- unique indexes for uniqueness
- partial indexes for conditional uniqueness
- expression indexes for normalized uniqueness

Examples:

```sql
CREATE UNIQUE INDEX employees_email_lower_idx ON employees (lower(email));
CREATE UNIQUE INDEX active_product_names_idx ON products (lower(name)) WHERE is_active;
```

### Exclusion Constraints

Use exclusion constraints when rows must not overlap according to an operator.

Example future use case: prevent overlapping employee shifts or machine reservations.

```sql
EXCLUDE USING gist (employee_id WITH =, shift_range WITH &&)
```

### Generated Columns

Use generated columns for derived values that should always stay consistent with base columns.

Example:

```sql
normalized_email text GENERATED ALWAYS AS (lower(email)) STORED
```

### DEFAULT Values

Use defaults for database-owned values.

Examples:

```sql
created_at timestamptz NOT NULL DEFAULT now()
is_active boolean NOT NULL DEFAULT true
```

Defaults are not validation by themselves, but they reduce missing-data bugs.

### Identity Columns And Sequences

Use identity columns when the database should generate numeric identifiers.

Example:

```sql
id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY
```

MES Lite currently uses explicit business IDs and UUID-like identifiers where appropriate.

### Domains

Use PostgreSQL domains for reusable scalar constraints.

Example:

```sql
CREATE DOMAIN non_blank_text AS text CHECK (btrim(VALUE) <> '');
```

Domains can reduce duplicated checks, but they add schema abstraction. Introduce them only after repeated constraints become painful.

### Enum Types

Use PostgreSQL enum types only for stable, rarely changing sets.

Example:

```sql
CREATE TYPE product_status AS ENUM ('active', 'inactive');
```

Prefer lookup tables when users need to configure values.

### Lookup Tables

Use lookup tables for configurable or business-owned reference data.

Example future use case: user-configurable product categories.

```sql
CREATE TABLE product_categories (
    id uuid PRIMARY KEY,
    name text NOT NULL UNIQUE
);
```

### Triggers

Use triggers for cross-row or audit behavior that cannot be expressed cleanly with constraints.

Examples:

- update `updated_at`
- write audit records
- enforce complex legacy constraints

Triggers are powerful but hidden. Prefer explicit constraints and application code when possible.

### Row-Level Security

Use row-level security when the database must enforce per-user data access.

Example future use case: tenant isolation or strict authorization at the database layer.

RLS is not a substitute for application authorization, but it can provide defense-in-depth.

### Views With CHECK OPTION

Use updatable views with `WITH CHECK OPTION` when writes through a restricted view must satisfy the view predicate.

Example future use case: restrict a write path to active records only.

### Transactions

Use transactions to make multi-step business operations atomic.

Transactions enforce integrity over time, not just per row.

Examples:

- validate references and insert production entry atomically
- update an order and insert its audit event together
- rollback an import batch after validation failure

### Isolation Levels And Locks

Use isolation levels and locks when concurrent transactions can violate business rules.

Capabilities:

- `READ COMMITTED`
- `REPEATABLE READ`
- `SERIALIZABLE`
- `SELECT ... FOR UPDATE`
- advisory locks

These should be introduced only when a concrete concurrency problem exists.

## Project Convention

New code should follow this flow:

1. Handler validates transport shape.
2. Handler converts request DTO to command or constructor arguments.
3. Constructor validates entity invariants and returns `(T, error)` when invariants can fail.
4. Application service validates workflow rules and references.
5. Repository validates defensively and translates storage errors.
6. Database constraints enforce critical integrity.

When adding a new validation rule, ask:

1. Is this a transport rule?
2. Is this a domain invariant?
3. Is this a workflow rule involving multiple resources?
4. Is this a storage integrity rule that must always hold?
5. Should the same rule exist in more than one layer for defense-in-depth?
