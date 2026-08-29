# ADR 0001: Project Layout

**Status:** Accepted — extended by ADR 0005, which adds an `internal/platform/`
tier for technical packages. Business slices still sit directly under
`internal/` as described below.

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

1. **`cmd/`** — each subdirectory is `main` package and produces one binary.
   The binary name matches the directory name. Currently only `cmd/server/`.
   This is standard convention to organize multiple binaries in single module.
2. **`internal/`** — packages under `internal/` can be imported only by code
   rooted at the parent, preventing external consumers (packages) from depending
   on implementation details. This is enforced by the Go compiler.
3. **`docs/adr/`** — Architecture Decision Records live here.
4. **`migrations/`** — SQL migration files managed by `goose`.
5. **`cmd/server/main.go`** - Composition root is where wiring happens

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
- No package under `internal/` can ever be imported from outside the module.
- Vertical slice structure means business logic tests live alongside handlers.

## Concepts

### Go Modules

Reference: <https://go.dev/ref/mod#modules-overview>

#### Module

A **module** is a collection of Go packages that are versioned, released, and
distributed together.

- **Module path** — unique import path (describes where to find module and what
  it does) declared in `go.mod` , for example:
  - `golang.org/x/net` — module rooted at the repository root
  - `golang.org/x/tools/gopls` — module rooted in the `gopls/` subdirectory
  - `golang.org/x/tools/gopls/v2` — major version 2 module (`/v2` suffix is
    required for versions ≥2)
- **Module root** — the directory containing `go.mod`.
- **Main module** — the module containing the current working directory when a
  `go` command is executed.

`go.mod` is the module manifest. It declares the module path, Go version, and
dependencies. Dependency resolution uses **Minimal Version Selection (MVS)**,
which is deterministic and simpler than constraint-solving approaches used by
some other package managers.

#### Package

A **package** is a collection of Go source files in the same directory that are
compiled together.

- **Package path** = `MODULE_PATH/PACKAGE_SUBDIRECTORY`, for example:
  `golang.org/x/net/html`, where `html` is the package.
- **Library package** — intended to be imported by other packages; typically
  named after its directory.
- **Executable package** — always named `main`; the Go toolchain builds it into
  an executable binary.

### Project Conventions

- **Go toolchain** — `go build`, `go test`, `go mod`, `go fmt`, `go vet`
- `internal/version/` follows the standard linker-flag pattern for embedding
  build information.
- CI runs tests with `-race` and `-shuffle=on` to detect data races and hidden
  test-order dependencies.
- The pre-commit hook runs `go vet` to catch common issues before they reach CI.

## Commands
```sh
git init --initial-branch=main --object-format=sha1
go mod init github.com/mac-hel/mes-lite"
make install-tools
chmod +x .githooks/pre-commit && git config core.hooksPath .githooks
```
