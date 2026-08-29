### Lesson 11.1b Scope

Move technical packages into `internal/platform/` and enforce the dependency direction between tiers with `depguard`.

#### Business Context

None directly. This lesson makes the codebase communicate its own structure, which matters for every business change that follows.

#### Problem

`internal/` held thirteen packages of three different kinds with nothing distinguishing them: seven business slices, five technical packages and one composition root. The layout claimed Vertical Slice Architecture while listing `config` and `version` alongside `production` and `orders`.

The deeper problem was that the dependency direction was correct by habit rather than by rule. A graph over `internal/` showed the technical packages importing nothing internal, the slices importing each other, and `server` importing everything — but nothing prevented a future technical package from importing a business slice and inverting that.

#### Design Discussion

The directory split and the lint rules solve different halves of the problem. The directory tells a reader what a package is for; the rules tell the compiler what a package may do. The rules are the half that actually holds, and they would have been worth adding even without moving a file.

Two rules cover the direction that matters:

- `internal/platform/**` may not import a business slice or `internal/server`. A platform package must stay importable by every slice, which means importing none of them.
- Only `cmd/` may import `internal/server`. The composition root points at everything; nothing points back.

Nested `internal/` directories were considered for compiler-level enforcement and rejected for this purpose: nested `internal/` constrains *who may import a package*, not *what a package may import*, so expressing "platform must not import business" would mean burying the slices under an artificial subdirectory. It stays the right tool for the generated sqlc packages, which is recorded as a known gap.

Two classification calls were not obvious:

- **`jobs` moved to platform.** It is a generic queue over opaque `[]byte` payloads with no business imports. Putting it under the rule commits the project to keeping it that way: when Lesson 11.4 connects the queue to CSV import, the handler registration must live in the composition root rather than inside `jobs`. The constraint is the reason to move it, not a side effect.
- **`auth` stayed a business slice.** JWT signing and HTTP middleware look technical, but users, roles and login are business capabilities, and `production` imports `auth` for the correction actor. Splitting the slice to relocate its middleware would break a working boundary for a cosmetic gain.

#### Go Concepts

- import paths as the primary boundary marker in Go
- import cycles as the compiler's own architectural enforcement
- what nested `internal/` can and cannot express
- lint configuration as executable architecture documentation

#### Architecture Concepts

- tiers within a vertical-slice layout: platform, slices, composition root
- dependency direction enforced in CI rather than in review
- classification rules written down so future packages are placed consistently

### Lesson 11.1b Completion Notes

#### Business Context

The repository now states its own structure, and CI fails when that structure is broken.

#### Problem

Thirteen packages sat in one flat directory with three different roles and an unstated dependency direction.

#### Design Discussion

`config`, `ids`, `jobs`, `postgres` and `version` moved to `internal/platform/`. Business slices and the composition root did not move. No code changed beyond import paths, so the commit has no behavioral effect.

The two depguard rules were verified to fire rather than assumed to work. A rule that matches no files also reports zero issues, so each rule was probed with a deliberate violation and the failure message checked before the probe was removed.

Probing exposed something worth knowing: importing `internal/server` from any existing package is already a compile error, because `server` imports every slice and the import becomes a cycle. The second rule therefore adds little today. It was kept because it applies to packages that do not exist yet, which is where the mistake would actually happen — a new package the server does not import could import the server without any cycle.

The first rule is not redundant. `platform/jobs` importing `production` produces no cycle, so nothing but the rule prevents it.

#### Implementation

- Created `internal/platform/` and moved `config`, `ids`, `jobs`, `postgres` and `version` into it with `git mv`.
- Rewrote the import paths in every affected package and restored import grouping with `goimports`.
- Added the `depguard` linter to `.golangci.yml`.
- Added rule `platform-stays-business-free`: no business slice or composition-root imports from `internal/platform/**`.
- Added rule `composition-root-is-not-a-dependency`: `internal/server` importable only from `cmd/`.
- Added `docs/package-layout.md` describing the tiers, the rules and the test for classifying a new package.
- Added ADR `0005-platform-package-tier.md`.
- Amended ADR 0001 to record that it is extended, and ADR 0004 to record the new path for `jobs`.

#### Tests

- Probed rule one with `platform/jobs` importing `production`; confirmed the depguard failure, then removed the probe.
- Probed rule two with a throwaway package importing `internal/server`; confirmed the depguard failure, then removed the probe.
- Confirmed a clean tree afterwards with `git status`.
- Verified with `go build ./...`.
- Verified with `go vet ./...`.
- Verified with `gofmt -l` and `goimports -l` over `cmd` and `internal`.
- Verified with `go test ./... -count=1`.
- Verified with `go test ./internal/platform/... -race -count=2`.
- Verified with `golangci-lint run ./...`.

PostgreSQL integration tests were compiled and vetted but skipped, because no test database was available.

#### Refactoring

This lesson is the refactoring. No package gained or lost behavior; only paths and lint configuration changed.

Historical lesson notes still refer to the old paths. They were left alone: they record what was true when written, and rewriting them would turn a project log into a fiction.

#### Code Review

An experienced Go engineer would approve the rules more readily than the move. The move is a matter of taste that reasonable engineers split on; the rules are a guarantee the codebase did not previously have.

The reviewable weakness is the deny lists: each business slice is named explicitly, so adding a slice means editing `.golangci.yml`, and forgetting to leaves a hole. An allow-list on the platform rule would fail closed instead, at the cost of naming every permitted standard-library and third-party import. The deny list was kept for now and the maintenance burden written into ADR 0005 rather than hidden.

#### Exercises

- Add a business slice, deliberately forget the deny list, and confirm the hole exists.
- Rewrite `platform-stays-business-free` as an allow-list and judge which version you would rather maintain.
- Explain why `platform/jobs` importing `production` compiles but `production` importing `server` does not.
- Decide where an `internal/platform/telemetry` package would sit if it needed to log a production entry ID.

#### Interview Questions

- What does Go's `internal/` directory actually enforce, and what does it not?
- How does an import cycle act as an architectural constraint?
- Why is a directory structure a weaker guarantee than a lint rule?
- How would you verify that a lint rule is actually being applied?
- When does a shared technical package become a business package?

#### Roadmap Update

- Lesson 11.1b completed.
- Current lesson moved to Lesson 11.2.
- Known technical debt updated for the unenforced sqlc boundary and the depguard deny-list maintenance.
