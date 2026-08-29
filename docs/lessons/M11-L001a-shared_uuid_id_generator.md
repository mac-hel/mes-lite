### Lesson 11.1a Scope

Replace the three duplicated UUID-shaped identifier generators with one shared helper, and drop the error return that the standard library guarantees can never fire.

#### Business Context

None directly. Identifier generation is technical plumbing that every slice creating a record depends on, and it should behave identically everywhere.

#### Problem

`production.NewEntryID`, `orders.NewOrderID` and `jobs.NewJobID` were byte-for-byte identical apart from one word in an error message. The duplication was recorded as technical debt when the third copy was added in Lesson 11.1.

Duplication was not the only cost. The wrapper names had already started to lie: `Service.CorrectEntry` generated a *correction* ID by calling `NewEntryID()`, because that was the generator its package happened to own.

#### Design Discussion

The helper lives in `internal/ids`, following the precedent set by `internal/postgres`: a small technical package below the business slices, not a business slice of its own. It is named so call sites read `ids.New()` without stutter.

**Should the helper return an error?**

Since Go 1.24, `crypto/rand.Read` is documented to never return an error. It fills the buffer completely or crashes the program irrecoverably, because the operating system APIs it uses are documented not to fail on anything but legacy Linux systems.

That makes every `if err != nil` around the old generators unreachable code. Two options were weighed:

- `New() (string, error)` keeps the current shape and stays honest if the randomness source is ever swapped for one that can fail. The cost is that every caller keeps writing, reading and reviewing a branch the runtime can never take.
- `New() string` leans on the guarantee. It deletes roughly ten unreachable error branches across services and tests, and it matches the project's architecture principle that panics are acceptable only for unrecoverable failures. A machine with an unusable CSPRNG cannot serve requests at all.

`New() string` was chosen. "Errors are values" is not a reason to invent an error that cannot occur; it is a reason to be precise about which failures are real.

**Should the per-slice wrappers survive?**

They were removed rather than reduced to one-line delegates. Three exported functions that each forward to the same helper are still three things to name, document and keep in sync, and keeping `NewEntryID` would have preserved the misleading call in `CorrectEntry`. Call sites now say `ids.New()`, which claims nothing about which record the identifier is for.

The trade-off accepted here is churn: this touches two services, one repository and several test files, in exchange for one implementation and one test suite.

#### Go Concepts

- knowing when a standard-library contract changes, and letting the API reflect it
- not modelling impossible failures as errors
- `_, _ =` as an explicit, lint-visible discard of a documented non-error
- small technical packages below business slices
- package naming that avoids stutter: `ids.New`, not `ids.NewID`

#### Architecture Concepts

- deduplication as refactoring, not as premature abstraction
- shared plumbing kept free of business meaning
- removing an abstraction whose name had begun to lie

### Lesson 11.1a Completion Notes

#### Business Context

Identifier generation now has one implementation and one test suite for the whole application.

#### Problem

Three identical generators had accumulated across the production, orders and jobs slices, and one of them was being used for a record type its name did not describe.

#### Design Discussion

Added `internal/ids` with a single `New() string` returning a canonical version 4 UUID. Removed the three per-slice generators and pointed every call site at the helper.

The helper is deliberately ignorant of what it identifies. It has no notion of entries, orders or jobs, which is what makes it safe to share; the moment such a helper starts taking a "kind" argument, it has become a business concept in the wrong package.

`internal/auth` was left alone. It does not generate identifiers: the bootstrap administrator uses the fixed literal `"bootstrap-admin"`. When auth-user management arrives, it should call `ids.New()` rather than grow a fourth generator.

#### Implementation

- Added `internal/ids/ids.go` with `New() string`.
- Removed `production.NewEntryID`, `orders.NewOrderID` and `jobs.NewJobID`.
- Updated `production.Service.Register` and `production.Service.CorrectEntry` to call `ids.New()`.
- Updated `orders.Service` order creation to call `ids.New()`.
- Updated `csvimport` entry conversion to call `ids.New()`, removing an error return that could no longer fail.
- Removed the now-unreachable error branches at every call site.
- Dropped the `crypto/rand` and `encoding/hex` imports from three business packages.

#### Tests

- Added canonical-format tests: length, hyphen positions, version nibble, RFC 4122 variant nibble and lower-case hex.
- Added a uniqueness test over 1,000 identifiers.
- Added a concurrency test proving the generator needs no external synchronization.
- Removed the two per-slice generator tests that the shared suite replaces.
- Verified with `go build ./...`.
- Verified with `go vet ./...`.
- Verified with `gofmt -l ./cmd ./internal`.
- Verified with `go test ./... -count=1`.
- Verified with `go test ./internal/ids ./internal/jobs ./internal/production ./internal/orders -race -count=2`.
- Verified with `golangci-lint run ./...`.

PostgreSQL integration tests were compiled and vetted but skipped in this run, because no test database was available. The edits to `internal/production/postgres_store_test.go` are mechanical call-site replacements.

#### Refactoring

This lesson is the refactoring. It closes the duplication debt recorded in Lesson 11.1 and removes a naming problem that predated it.

The identifier format did not change, so no migration or data change is involved: the same UUID strings are produced by the same bit manipulation, from one place instead of three.

#### Code Review

An experienced Go engineer would approve the extraction and, more notably, the signature change. Returning an error that the standard library guarantees cannot occur is a small but real cost paid at every call site, and removing it makes the remaining error checks in those functions meaningful.

The discarded return in `_, _ = rand.Read(b[:])` deserves its comment. A silent discard would read as an oversight; the comment states the contract being relied on, which is what a reviewer needs in order to agree.

#### Exercises

- Find the last remaining place where a business slice builds an identifier without `ids.New()` and decide whether it should change.
- Replace `internal/ids` with `github.com/google/uuid` (already an indirect dependency) and argue whether the dependency earns its place.
- Explain what would have to be true about the randomness source for `New()` to need an error return again.
- Write a benchmark for `ids.New()` and predict where the allocations are before running it.

#### Interview Questions

- When is duplication cheaper than the abstraction that removes it, and how do you tell?
- Why is "errors are values" not an argument for returning an error that cannot happen?
- What does `crypto/rand.Read` do on failure since Go 1.24, and why was that changed?
- What distinguishes a legitimate shared technical package from a `utils` dumping ground?
- Why should a shared ID generator never take a "kind" or "prefix" argument?

#### Roadmap Update

- Lesson 11.1a completed.
- Current lesson moved to Lesson 11.2.
- Known technical debt updated: duplicated ID generation resolved.
