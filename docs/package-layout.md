# Package Layout

This document describes where code belongs in `internal/` and which dependency
rules hold between the tiers. The rules are enforced by `depguard` in
`.golangci.yml`, not only by convention.

## Tiers

```
cmd/                    entry points; the only code allowed to import internal/server
internal/
  platform/             technical packages with no business meaning
    config/             environment configuration
    ids/                UUID-shaped identifier generation
    jobs/               background job model and in-memory queue
    postgres/           PostgreSQL SQLSTATE constants
    version/            build information
  auth/                 business slice: application users, roles, JWT, middleware
  csvimport/            business slice: historical production data import
  employees/            business slice
  orders/               business slice: production orders
  production/           business slice: production entries and corrections
  products/             business slice
  reporting/            business slice: read models for management reports
  server/               composition root: routing, middleware wiring, shutdown
```

## Rules

1. **`internal/platform/**` must not import business slices or `internal/server`.**

   A platform package has to stay importable by every slice, which means it can
   import none of them. If a platform package needs business knowledge, that is
   the signal it is not a platform package.

2. **Only `cmd/` may import `internal/server`.**

   The composition root points at everything; nothing points back at it.

Both rules are `depguard` rules in `.golangci.yml` and fail the build when
broken. In practice the compiler already catches most violations of rule 2 as an
import cycle, because `internal/server` imports every slice. The rule exists for
packages the server does not import yet.

## Deciding where a new package goes

Ask what the package would need to know:

- Knows about employees, products, production, orders or reports → business
  slice, top level under `internal/`.
- Knows only about technique — encoding, transport, identifiers, configuration,
  time, telemetry → `internal/platform/`.
- Only exists to connect the two → composition root, or `cmd/`.

A package that would need a business enum, a business error or a business type
in its signature belongs in a slice, whatever its name suggests.

## Known gaps

- **Generated sqlc packages are not enforced.** `docs/sqlc-boundaries.md` says
  generated database code stays below its slice. `internal/csvimport` currently
  imports `internal/production/productiondb` directly to reuse the production
  entry insert. Nested `internal/` directories
  (`internal/production/internal/productiondb`) would make that a compile error,
  but moving them requires updating `sqlc.yaml` output paths and deciding how
  CSV import persists batches without reaching into another slice's generated
  code. Deferred until that decision is made.
- **`jobs.TypeProductionEntryImport`** names a business workload from inside a
  platform package. It is a string constant and imports nothing, so no rule is
  broken, but the constant would sit better next to the workload it names. See
  Lesson 11.4.
