# sqlc Boundaries

This document defines how MES Lite uses sqlc while preserving vertical-slice boundaries.

The goal is to keep SQL explicit and type-safe without letting generated database code become the application architecture.

## Principles

- Migrations are the shared database contract.
- sqlc query files are owned by vertical slices.
- Generated query methods belong to the slice that defines the query.
- Generated table models are implementation details, not domain models.
- HTTP handlers must not expose sqlc-generated types.
- One slice must not import another slice's generated `*db` package.
- Cross-slice reporting queries are allowed inside the reporting slice because reporting is a read-model projection layer.

## Shared Schema Input

All sqlc packages currently use the shared migration directory as schema input:

```yaml
schema: "migrations"
```

This means each generated sqlc package can see the full PostgreSQL schema.

That is intentional for now.

sqlc needs schema knowledge to type-check queries. Since migrations are shared and foreign keys connect tables across slices, using the migration directory as schema input is the simplest reliable setup.

## Expected Generated-Code Noise

Because each sqlc package sees the full schema, generated `models.go` files may contain table structs unrelated to that slice.

Example: `internal/reporting/reportingdb/models.go` may include auth, order, employee, product and production table models even though reporting only owns report queries.

This is a tooling artifact, not an architectural dependency.

The important boundary is not which table structs sqlc generated. The important boundary is which generated code the application imports and uses.

## Query Ownership

Each vertical slice owns its query directory:

```text
internal/employees/queries
internal/products/queries
internal/production/queries
internal/orders/queries
internal/auth/queries
internal/reporting/queries
```

Only queries from that directory generate methods for that slice's `*db` package.

For example:

```text
internal/reporting/queries/reports.sql
```

generates reporting query methods in:

```text
internal/reporting/reportingdb
```

Those query methods are owned by reporting.

## Allowed Usage

Slice repositories may use their own generated package.

Examples:

```go
package reporting

import "github.com/mac-hel/mes-lite/internal/reporting/reportingdb"
```

```go
package employees

import "github.com/mac-hel/mes-lite/internal/employees/employeesdb"
```

Reporting may join tables across slices because reports are read-side projections over operational data. The reporting slice owns the report query shape and maps generated rows into reporting read models.

## Disallowed Usage

Do not import another slice's generated package.

Bad examples:

```go
package reporting

import "github.com/mac-hel/mes-lite/internal/employees/employeesdb"
```

```go
package orders

import "github.com/mac-hel/mes-lite/internal/products/productsdb"
```

Do not return sqlc-generated table models from HTTP handlers.

Bad example:

```go
func (h *Handler) Get(c fuego.ContextNoBody) (reportingdb.ProductionEntry, error) {
    // ...
}
```

Do not use generated table models as domain entities.

Bad example:

```go
type Order = ordersdb.ProductionOrder
```

## Mapping Rule

Repositories map generated database rows into slice-owned types.

Example shape:

```go
func rowFromDB(row reportingdb.DailyProductionRow) DailyProductionRow {
    return DailyProductionRow{
        Day:           row.Day.Time.UTC(),
        ProductSKU:    row.ProductSku,
        TotalQuantity: int(row.TotalQuantity),
        EntryCount:    int(row.EntryCount),
    }
}
```

The exported type belongs to the slice. The generated row type stays inside persistence code.

## Reporting Exception

Reporting is allowed to read across slices.

This does not violate vertical-slice architecture because reporting is a read-model layer. It does not own employee, product or production business rules. It only projects existing data into management views.

Reporting must still follow these rules:

- own its SQL queries in `internal/reporting/queries`
- generate methods in `internal/reporting/reportingdb`
- map generated rows into reporting read models
- expose HTTP DTOs, not generated sqlc rows

## When To Revisit

The current shared-schema approach should be revisited only if the generated-code noise creates real maintenance pain.

Possible future options:

- split schema inputs by slice
- maintain sqlc-specific schema files
- reduce generated model output if sqlc supports it cleanly
- create a dedicated schema package for generation inputs

Do not do this preemptively. Splitting schema inputs can duplicate migration knowledge and make foreign-key changes harder to maintain.

## Review Checklist

When reviewing persistence code, check:

- Does the slice import only its own generated `*db` package?
- Are generated sqlc rows mapped before leaving repository code?
- Are HTTP responses slice-owned DTOs or domain/read-model types?
- Are reporting cross-table joins read-only projections?
- Is generated `models.go` noise ignored unless it leaks into application code?
