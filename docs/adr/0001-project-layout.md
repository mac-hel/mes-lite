# ADR 0001: Project Layout

**Status:** Accepted

**Deciders:** Initial project setup

**Date:** 2026-07-20

## Context

The project needs a standard, predictable layout. Every Go project should follow
conventions that make it easy for new contributors to find code, understand
boundaries, and know where to add new functionality.

The key questions were:

1. How should binaries and libraries be separated?
2. How should internal packages be protected from external consumption?
3. How should the project structure evolve as business capabilities grow?

## Decision

We use the standard Go project layout:

```
cmd/           # Application entry points (one directory per binary)
internal/      # Private packages, not importable from outside
docs/          # Documentation and ADRs
migrations/    # Database migrations
```

Key rules:

1. **`cmd/`** — each subdirectory produces one binary. The binary name matches the
   directory name. Currently only `cmd/server/` exists.
2. **`internal/`** — packages under `internal/` cannot be imported by any module
   rooted outside this project. This is enforced by the Go compiler.
3. **`docs/adr/`** — Architecture Decision Records live here.
4. **`migrations/`** — SQL migration files managed by `goose`.

## Architecture: Vertical Slices

When business packages are added (Milestone 2+), they will go directly under
`internal/` rather than under a shared `pkg/`:

```
internal/
  employees/     # Owns: handlers, logic, persistence, tests
  products/      # Owns: handlers, logic, persistence, tests
  production/    # Owns: handlers, logic, persistence, tests
```

No shared `handlers/`, `services/`, or `repositories/` directories at the top
level. Each capability owns its full vertical slice.

## Alternatives Considered

1. **Flat layout** (everything in root) — rejected because it becomes
   unmanageable as the project grows. No clear ownership boundaries.

2. **`pkg/` + `internal/` split** — rejected because `pkg/` would imply the
   packages are intended for external consumption. In this project, nothing
   is. Everything is `internal/`.

3. **Layered layout** (`handlers/`, `services/`, `repositories/`) — rejected
   because it couples code by technical layer rather than business capability.
   Changes to one feature touch four directories.

## Trade-offs

- **Pro:** Predictable. Any Go engineer recognizes this layout.
- **Pro:** Compiler-enforced encapsulation via `internal/`.
- **Pro:** Scales to multiple binaries without restructuring.
- **Con:** Slightly more nesting than a flat structure.

## Consequences

- All application code lives under `internal/`.
- Adding a new binary means creating a new directory under `cmd/`.
- No package under `internal/` should ever be imported from outside the module.
- Vertical slice structure means business logic tests live alongside handlers.
