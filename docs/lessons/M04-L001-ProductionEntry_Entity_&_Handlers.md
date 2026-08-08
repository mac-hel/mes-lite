### Lesson 4.1 Completion Notes

#### Business Context

Production workers need to record completed work without Excel.

#### Problem

The system had employees and products, but no core production-registration capability.

#### Design Discussion

The first implementation keeps the vertical slice small: `internal/production` owns the production entry entity, validation, in-memory persistence contract and HTTP handler. PostgreSQL, sqlc and transactions are intentionally postponed to later lessons because introducing them before the business model would hide the domain rules inside infrastructure work.

#### Go Concepts

- `time.Time` for business timestamps
- `crypto/rand` and `encoding/hex` for standard-library UUID-shaped IDs
- error wrapping with sentinel errors
- context propagation from HTTP handler to store

#### Architecture Concepts

- vertical slice for production registration
- package naming that avoids stutter: `production.Entry`, not `production.ProductionEntry`
- concrete in-memory implementation before persistence abstraction grows

#### Implementation

- Added `POST /production-entries`.
- Added `production.Entry` with employee, product, quantity, workstation, timestamp and comment.
- Added production entry validation and UUID-shaped ID generation.
- Wired production handler through the server composition root.

#### Tests

- Domain validation tests added.
- In-memory store tests added.
- HTTP registration tests added.
- Server route test added.

#### Refactoring

- Renamed `ProductionEntry` to `Entry` after lint identified package-name stutter.

#### Code Review

- An experienced Go engineer would likely approve the scope for this lesson because it is small, explicit and tested.
- Remaining gap: employee/product existence is not validated yet; this belongs in the next application-service/persistence lessons where the transaction boundary can be designed properly.

#### Exercises

- Explain why `production.Entry` is a better exported name than `production.ProductionEntry`.
- Add a table test for future maximum comment length validation.
- Explain why timestamps are normalized to UTC.

#### Interview Questions

- Why is `time.Time` usually preferred over strings for timestamps in Go APIs?
- What does error wrapping give us when translating domain errors to HTTP errors?
- Why should `context.Context` be passed from the handler to persistence code?
- Why is package-name stutter considered non-idiomatic Go?

#### Roadmap Update

- Lesson 4.1 completed.
- Current lesson moved to Lesson 4.2.
- Standard Library `time` marked complete in the Knowledge Matrix.
