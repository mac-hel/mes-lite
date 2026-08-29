# ADR 0005: Separate Platform Packages From Business Slices

## Status

Accepted

Extends ADR 0001 (Project Layout). The `cmd/`, `internal/`, `docs/` and
`migrations/` structure decided there is unchanged; this ADR adds one tier
inside `internal/`.

## Context

`internal/` had grown to thirteen packages mixing three different kinds of code:
seven business slices (`auth`, `csvimport`, `employees`, `orders`, `production`,
`products`, `reporting`), five technical packages (`config`, `ids`, `jobs`,
`postgres`, `version`) and one composition root (`server`).

The dependency direction between those kinds was real but unstated. A dependency
graph over `internal/` showed the technical packages had zero internal imports,
the business slices imported each other, and `server` imported everything. That
arrangement was correct by habit rather than by rule, so nothing prevented a
future technical package from importing a business slice and inverting the
direction.

Milestones 13 to 15 add observability, profiling and build-time packages, which
would roughly double the technical group and make the flat listing harder to
read.

## Decision

Introduce `internal/platform/` for technical packages with no business meaning,
and enforce the dependency direction with `depguard` rules in `.golangci.yml`:

- `internal/platform/**` may not import any business slice or `internal/server`.
- Only `cmd/` may import `internal/server`.

Business slices stay directly under `internal/`. The composition root stays at
`internal/server`.

The enforcement is the point of this ADR. The directory expresses intent; the
lint rules make the intent checkable, and they would have been worth adding even
without moving a single file.

`internal/jobs` moved into `platform/`. Today it is a generic queue over opaque
`[]byte` payloads with no business imports. Placing it under the rule commits
the project to keeping it that way: when Lesson 11.4 connects the queue to the
CSV import workload, the handler registration must live in the composition root
rather than inside `jobs`. That constraint is deliberate.

`internal/auth` stays a business slice. It holds JWT signing and HTTP middleware,
which look technical, but users, roles and login are business capabilities, and
`internal/production` imports it for the correction actor. Splitting the slice to
move its middleware would break a working boundary for a cosmetic gain.

## Alternatives Considered

- **Leave `internal/` flat and add only the depguard rules**: rejected, though it
  was close. It captures most of the value at none of the churn. The directory
  was chosen because the technical group is about to double and the flat listing
  is where a reader looks first.
- **Nested `internal/` directories to enforce the rule with the compiler**:
  rejected for this purpose. Nested `internal/` constrains who may import a
  package, not what a package may import, so expressing "platform must not
  import business" would require burying the slices under an artificial
  subdirectory. It remains the right tool for the generated sqlc packages, which
  is recorded as a known gap in `docs/package-layout.md`.
- **`internal/pkg/` or `internal/common/`**: rejected. Names that mean "things"
  attract things. `platform` states what the tier is for.
- **Move `server` under `platform/`**: rejected. A composition root is neither
  technical nor business; it is the place where both meet.

## Trade-Offs

- Pro: the dependency direction is checked by CI instead of by review attention.
- Pro: `internal/` now lists business capabilities first, which is what the
  Vertical Slice Architecture claims to optimize for.
- Pro: the rules apply to packages that do not exist yet, which is where the
  mistake is most likely.
- Con: every import of a moved package changed, producing churn and git-blame
  noise across the repository in a commit with no behavior change.
- Con: the tier boundary needs a judgment call for packages that are partly
  both, and `auth` shows those calls are not always obvious.
- Con: `depguard` deny lists name each business slice explicitly and must be
  extended when a slice is added.

## Consequences

- New packages must be classified on arrival; `docs/package-layout.md` gives the
  test to apply.
- Adding a business slice means adding it to the `platform-stays-business-free`
  deny list.
- Lesson 11.4 must register job handlers in the composition root, not inside
  `internal/platform/jobs`.
- The generated sqlc boundary remains unenforced and is the next candidate for
  compiler-level enforcement through nested `internal/`.
