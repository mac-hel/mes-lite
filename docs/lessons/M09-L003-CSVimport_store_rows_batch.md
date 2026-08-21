### Lesson 9.3 Scope

Persist validated CSV import records in batches and define the first transaction strategy for historical production-entry imports.

#### Business Context

Historical imports should not leave the database half-updated when one row in a valid-looking batch fails at the persistence boundary. Managers need imported history to be trustworthy.

#### Problem

Lesson 9.2 produced typed valid import records, but there was no durable save path. Saving rows one by one outside a transaction could persist earlier rows and then fail on a later foreign-key or constraint error.

#### Design Discussion

The batch persistence boundary is all-or-nothing for records passed to `SaveBatch`. If a row fails during insertion, the transaction rolls back every row in that batch.

Partial failure reporting remains a higher-level API concern for Lesson 9.4. L9.3 only answers the storage question: a batch either commits completely or leaves no production-entry rows behind.

#### Go Concepts

- explicit transaction lifecycle with `Begin`, deferred `Rollback` and `Commit`
- wrapping row-specific persistence errors with `Unwrap`
- converting validated import rows into production domain entries
- returning empty non-nil slices for empty batch results

#### Architecture Concepts

- import store as a persistence boundary for CSV import workflow
- transaction boundary around one import batch
- PostgreSQL constraints as final integrity guardrails
- no generic transaction manager before repeated need exists

### Lesson 9.3 Completion Notes

#### Business Context

MES Lite can now persist validated historical production rows in one transactional batch.

#### Problem

Validated rows were still in memory only. A future upload endpoint would have had no safe way to save rows without risking partial database writes.

#### Design Discussion

Added a CSV import store that writes validated production-entry import records to PostgreSQL. `SaveBatch` generates production-entry IDs, constructs domain `production.Entry` values and inserts them through the existing sqlc production-entry query inside one transaction.

The transaction is intentionally owned by the CSV import store because the import workflow is the business operation. A broad transaction manager was not introduced.

#### Implementation

- Added `csvimport.Store` interface with `SaveBatch`.
- Added `csvimport.PostgresStore` backed by `pgxpool` transactions.
- Added `BatchError` that records the CSV row number that failed persistence.
- Converted `ProductionEntryRecord` values to `production.Entry` values before persistence.
- Inserted batch records through existing `productiondb.CreateEntry` sqlc query.
- Mapped PostgreSQL constraint failures to production domain errors.
- Returned empty non-nil slices for empty batches.

#### Tests

- Tested empty batch behavior without needing a database.
- Added PostgreSQL integration test for successful two-row batch persistence.
- Verified generated entry IDs are returned and persisted rows can be read by the production store.
- Added rollback integration test proving a missing reference on a later row leaves zero rows committed.
- Verified row number is preserved in `BatchError`.
- Verified with `go test ./internal/csvimport -count=1`.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

No new sqlc query was needed. The import store reuses the existing production-entry insert query instead of duplicating SQL.

#### Code Review

An experienced Go engineer would approve the transaction boundary because it protects exactly one batch import operation. The main caveat is expected lesson scope: API-level import summaries and partial failure reporting are not exposed yet.

#### Exercises

- Add a test proving duplicate generated IDs would roll back the whole batch by injecting deterministic IDs.
- Explain why deferred rollback is still called after a successful commit and why that is safe.
- Decide whether future imports should use one transaction for the whole file or smaller batch transactions for very large files.

#### Interview Questions

- What should define a transaction boundary?
- Why is all-or-nothing persistence useful for a validated import batch?
- How does `errors.As` find a wrapped `BatchError`?
- When would smaller transactions be better than one large import transaction?

#### Roadmap Update

- Lesson 9.3 completed.
- Current lesson moved to Lesson 9.4.
- L9.4 remains focused on upload/API summary behavior and partial failure reporting.
