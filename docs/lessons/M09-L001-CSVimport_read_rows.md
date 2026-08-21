### Lesson 9.1 Scope

Introduce the CSV import slice and implement a streaming reader for historical production-entry CSV files.

#### Business Context

The company has historical production data in Excel-like files. Before the system can replace manual tracking completely, it needs a safe path to read that history without loading entire files into memory.

#### Problem

CSV import can easily become memory-heavy if the implementation reads the whole file before processing. Large historical exports should be processed row by row through `io.Reader` and `encoding/csv.Reader`.

#### Design Discussion

Lesson 9.1 only defines the import boundary and raw CSV row reader. It does not parse quantities or timestamps into domain values yet, and it does not persist rows. Those responsibilities belong to later lessons so the streaming concept remains clear.

The reader validates the CSV header once, then exposes one raw production-entry row per `Read` call. This mirrors Go's standard streaming style: callers pull data until `io.EOF`.

#### Go Concepts

- `io.Reader` as the standard input abstraction
- `encoding/csv.Reader` for CSV parsing
- sentinel errors for stable import failures
- `io.EOF` as normal stream completion
- row numbers for future validation diagnostics

#### Architecture Concepts

- CSV import as its own vertical slice in `internal/csvimport`
- raw transport/import rows separated from production domain entries
- validation and persistence intentionally postponed to later pipeline steps

### Lesson 9.1 Completion Notes

#### Business Context

MES Lite now has the first building block for migrating historical production data from CSV exports.

#### Problem

The application had no import boundary. Starting with upload handlers or database writes would have mixed HTTP, parsing, validation and persistence before establishing a memory-safe reader.

#### Design Discussion

Added `internal/csvimport` with a concrete `ProductionEntryReader`. It accepts any `io.Reader`, validates the expected production-entry header and streams rows one at a time.

The expected CSV columns are `employee_id`, `product_sku`, `quantity`, `workstation`, `timestamp` and `comment`. Values remain strings in this lesson because row validation, type parsing and error collection are the focus of Lesson 9.2.

#### Implementation

- Added `ProductionEntryRow` raw import row type.
- Added `ProductionEntryReader` over `encoding/csv.Reader`.
- Added `NewProductionEntryReader` with header validation.
- Added `Read` returning one row at a time and `io.EOF` when complete.
- Added import errors `ErrInvalidHeader` and `ErrInvalidRecord`.
- Added row-number tracking for future validation diagnostics.

#### Tests

- Tested sequential row reading and trimming.
- Tested header normalization for case and whitespace.
- Tested missing and unexpected headers.
- Tested malformed CSV parse errors.
- Tested `io.EOF` stream completion.
- Tested constructor behavior only consumes the header before row reads.
- Verified with `go test ./internal/csvimport -count=1`.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

No existing production-registration code was changed. The import slice is intentionally separate because CSV rows are external input records, not validated production entries yet.

#### Code Review

An experienced Go engineer would approve the narrow scope: the code uses standard-library streaming, keeps APIs explicit and avoids premature importer abstractions. The main follow-up is to add typed row validation and collect row-level errors without stopping the whole import unnecessarily.

#### Exercises

- Add a test with 100,000 generated rows and assert rows are processed through repeated `Read` calls.
- Explain why `io.EOF` is not treated as an error condition by stream consumers.
- Extend the reader test with a CSV row containing a comma inside a quoted comment.

#### Interview Questions

- Why is `io.Reader` one of the most important interfaces in Go?
- Why should large CSV imports be streamed instead of read into memory?
- What does `encoding/csv.Reader` handle that manual `strings.Split` would get wrong?
- Why should parsing and business validation be separate from raw CSV reading?

#### Roadmap Update

- Lesson 9.1 completed.
- Current lesson moved to Lesson 9.2.
- Standard Library `io` and `encoding/csv` marked complete in the Knowledge Matrix.
