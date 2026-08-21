### Lesson 9.2 Scope

Parse raw CSV rows into typed import records and collect row-level validation errors without stopping the whole import.

#### Business Context

Historical production CSV files may contain mistakes from manual spreadsheets. Managers need an import process that identifies bad rows clearly instead of failing with a vague first error.

#### Problem

Lesson 9.1 could stream raw CSV fields, but every field was still a string. The importer could not yet distinguish valid rows from rows with missing IDs, invalid quantities or malformed timestamps.

#### Design Discussion

Validation is a pipeline step after raw CSV reading and before persistence. Structural CSV errors remain fatal because the stream cannot be trusted. Row-level business-shape errors are collected so the importer can later report partial failures.

The validated type is still import-specific. It is not a `production.Entry` because production-entry IDs, reference validation and persistence behavior belong to later lessons.

#### Go Concepts

- `strconv.Atoi` for integer parsing
- `time.Parse` with `time.RFC3339`
- `errors.Is` for stream completion checks
- small error structs implementing `error`
- collecting validation errors in slices

#### Architecture Concepts

- validation pipeline step between reader and persistence
- import records separated from production domain entities
- fatal stream errors separated from row-level validation errors

### Lesson 9.2 Completion Notes

#### Business Context

MES Lite can now identify valid historical production rows and explain row-level CSV mistakes in a way that a future API can return to managers.

#### Problem

Raw CSV reading alone was not enough for import. Quantity and timestamp fields needed type parsing, and invalid rows needed structured diagnostics without blocking valid rows from being recognized.

#### Design Discussion

Added `ValidateProductionEntries`, which consumes a `ProductionEntryReader`, validates each raw row and returns a `ValidationResult` containing valid typed records plus row errors.

Malformed CSV structure is still returned as a fatal error. Missing required fields, invalid quantities and invalid timestamps are collected as `RowError` values.

#### Implementation

- Added `ProductionEntryRecord` typed import record.
- Added `RowError` with row number, field and message.
- Added `ValidationResult` containing valid records and row errors.
- Added `ValidateProductionEntries` to stream through rows and collect validation errors.
- Added quantity parsing with positive and PostgreSQL `integer` bounds checks.
- Added RFC3339 timestamp parsing and UTC normalization.
- Mapped CSV field-count failures to `ErrInvalidRecord`.

#### Tests

- Tested valid CSV rows become typed records.
- Tested timestamp offsets normalize to UTC.
- Tested mixed valid and invalid rows keep valid records and collect all row errors.
- Tested missing, non-integer, negative and too-large quantities.
- Tested fatal malformed CSV record handling.
- Tested `RowError.Error` formatting.
- Verified with `go test ./internal/csvimport -count=1`.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

The reader now maps CSV field-count parse failures to the import sentinel `ErrInvalidRecord`, making stream-shape failures easier for future API code to translate.

#### Code Review

An experienced Go engineer would approve the separation between raw reading and typed row validation. The code remains concrete, standard-library based and avoids a generic validation framework.

The main follow-up is persistence strategy: decide how validated rows should be saved in batches, where transactions start and what partial failure means for valid versus invalid rows.

#### Exercises

- Add a test proving quoted comments containing commas survive validation.
- Add a row with several invalid fields and explain why all errors are collected.
- Decide whether future import summaries should include the original raw row values.

#### Interview Questions

- Why separate fatal stream errors from row-level validation errors?
- Why should CSV import validation not directly create persisted production entries yet?
- Why is RFC3339 a good timestamp format for import boundaries?
- What trade-off appears when collecting all errors instead of stopping at the first invalid row?

#### Roadmap Update

- Lesson 9.2 completed.
- Current lesson moved to Lesson 9.3.
- L9.3 remains focused on batch persistence and transaction strategy.
