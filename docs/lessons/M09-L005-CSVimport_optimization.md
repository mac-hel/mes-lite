### Lesson 9.5 Scope

Review the CSV import pipeline for large-file behavior before closing Milestone 9.

#### Business Context

Historical production files may be large. The import endpoint should remain predictable and should not require loading all valid rows or all validation errors into memory at once.

#### Problem

The L9.4 service streamed the HTTP body into the CSV reader, but then collected every valid row before saving. That defeated part of the milestone's streaming goal for large valid files.

#### Design Discussion

The import service now processes rows incrementally and saves valid records in bounded batches. This keeps the main valid-row memory cost proportional to the batch size rather than the file size.

Reported row errors are also capped. The summary still counts every invalid row, but only stores the first bounded set of errors. This is a production trade-off: exact counts remain available, while response size and memory usage are protected for very bad files.

If a batch fails with a row-level persistence error such as a duplicate production entry or invalid reference, the service retries that batch row by row. This keeps retry imports useful: already-imported rows can be reported as row errors while new valid rows in the same file still get imported.

#### Go Concepts

- bounded batching with slices
- exact counters separated from retained detail records
- streaming loop with `io.EOF` as normal completion
- memory safety trade-offs in API response design

#### Architecture Concepts

- import pipeline review before milestone closure
- batch size as an implementation detail
- bounded error reporting for production safety
- explicit technical debt for future PostgreSQL `COPY` optimization

### Lesson 9.5 Completion Notes

#### Business Context

MES Lite's CSV import path is now suitable for MVP historical data migration with bounded memory behavior for valid rows and bounded response/error storage for invalid rows.

#### Problem

The service accepted an `io.Reader`, but validated the entire CSV into slices before saving. Large valid imports would grow memory with file size.

#### Design Discussion

Refactored `ImportProductionEntries` into a streaming pipeline. It reads one row, validates one row and appends only valid records to a fixed-size batch. When the batch reaches the configured size, the service persists it and reuses the same slice.

The API reports at most `maxReportedErrors` row errors while continuing to count all invalid rows. The response includes `errorLimitReached` so clients know whether detailed errors were truncated.

Failed batches are isolated by retrying records individually when the failure is row-level. Unexpected infrastructure failures still terminate the import because they cannot be safely attributed to one CSV row.

#### Implementation

- Added `defaultImportBatchSize` for bounded valid-row batches.
- Added `maxReportedErrors` for bounded error detail storage.
- Refactored `ImportProductionEntries` to save batches during streaming instead of collecting all valid rows first.
- Added `errorLimitReached` to `ImportSummary`.
- Added row-by-row isolation after row-level batch persistence failures.
- Preserved exact `totalRows`, `validRows`, `invalidRows` and `importedRows` counters.
- Kept `ValidateProductionEntries` for focused validation tests, while the service uses the lower-level row validator for streaming.

#### Tests

- Added large-input test proving 1,205 valid rows are saved as batches of 500, 500 and 205.
- Added large-invalid-input test proving invalid row counts remain exact while reported errors are capped.
- Added retry/import-continuation tests proving duplicate/already-existing rows become summary errors while other valid rows still import.
- Verified no save calls happen when all rows are invalid.
- Verified with `go test ./internal/csvimport -count=1`.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

The service now owns the streaming pipeline directly instead of delegating to `ValidateProductionEntries`, because the service must interleave validation and batch persistence for large files.

#### Code Review

An experienced Go engineer would approve the milestone for MVP scale: the code is standard-library based, validates rows explicitly, persists valid rows in bounded batches, isolates row-level persistence failures and returns clear summaries.

The main performance follow-up is database insert strategy. Current batch persistence still uses one regular `INSERT` per row inside a transaction. PostgreSQL `COPY` or `pgx.CopyFrom` may be justified later for very large imports, but adding it now would complicate the lesson before real performance data exists.

#### Exercises

- Make `defaultImportBatchSize` configurable and discuss whether it belongs in environment configuration.
- Add a benchmark comparing all-at-once validation with bounded-batch streaming.
- Sketch how a future background job would stream the same import without blocking an HTTP request.

#### Interview Questions

- Why does accepting `io.Reader` not automatically guarantee low memory usage?
- Why can response error details need a cap even if row processing is streamed?
- What trade-offs exist between regular batched `INSERT` and PostgreSQL `COPY`?
- Why might CSV import become a background job in a production system?

#### Roadmap Update

- Lesson 9.5 completed.
- Milestone 9 completed.
- Current milestone moved to Milestone 10.
- Current lesson moved to Lesson 10.1.
- Concurrency `Pipelines` marked complete in the Knowledge Matrix.
