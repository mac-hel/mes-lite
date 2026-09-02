### Lesson 14.1 Scope

Introduce Go benchmarks on a real code path and establish the rule that optimization starts with measurement.

#### Business Context

CSV import and reporting are likely performance-sensitive paths as data grows. Before changing implementations, the project needs repeatable measurements that can show whether a change helps or hurts.

#### Problem

The codebase has many correctness tests but no benchmarks. Without benchmarks, performance discussions rely on guesses, and future refactors could accidentally make hot paths slower or more allocation-heavy.

#### Design Discussion

Start with a focused benchmark around CSV import validation or request middleware overhead. The benchmark should measure an existing behavior without changing production code first. The goal of L14.1 is learning benchmark mechanics and avoiding misleading results, not immediately optimizing.

Benchmarks should use `testing.B`, report allocations with `b.ReportAllocs()` and keep setup outside the measured loop where possible.

#### Go Concepts

- `testing.B`
- benchmark loops and `b.N`
- `b.ReportAllocs()`
- setup versus measured work
- avoiding benchmark dead-code elimination

#### Architecture Concepts

- performance baselines before optimization
- realistic input sizes for business-relevant paths
- benchmark results as engineering evidence

#### Tests

- Add at least one benchmark for a meaningful existing path.
- Run the benchmark with allocation reporting.
- Keep all correctness tests passing.

#### Exercises

- Compare benchmark results for small and large CSV samples.
- Move setup accidentally inside the benchmark loop and explain why the result changes.
- Decide which endpoint or package should receive the next benchmark.

#### Interview Questions

- How does Go's benchmark runner choose `b.N`?
- Why should setup usually be outside the measured loop?
- What does `allocs/op` tell you?
- Why can microbenchmarks mislead production optimization?

### Lesson 14.1 Completion Notes

#### Business Context

MES Lite now has its first repeatable performance baseline for a business-relevant path: CSV import validation.

#### Problem

The project had correctness tests for CSV import, but no benchmarks. That meant performance discussion around import behavior was based on assumptions rather than measurements.

#### Design Discussion

The first benchmark measures existing CSV validation behavior without changing production code. It uses generated CSV inputs with 100, 1,000 and 10,000 rows so small and larger import sizes can be compared.

CSV generation stays outside the measured loop. Each iteration creates a new reader because readers are consumed by validation. The benchmark stores the result in a package-level variable so the compiler cannot remove the measured work as unused.

#### Implementation

- Added `BenchmarkValidateProductionEntries` in `internal/csvimport/validation_benchmark_test.go`.
- Added benchmark cases for 100, 1,000 and 10,000 valid CSV rows.
- Used `testing.B` sub-benchmarks and `b.ReportAllocs()`.
- Added a package-level sink to avoid benchmark dead-code elimination.

#### Benchmark Results

- `100_rows`: `119705 ns/op`, `53240 B/op`, `227 allocs/op`.
- `1000_rows`: `1110757 ns/op`, `402042 B/op`, `2030 allocs/op`.
- `10000_rows`: `12646587 ns/op`, `6126470 B/op`, `20038 allocs/op`.

#### Tests

- Verified with `go test ./internal/csvimport -run '^$' -bench '^BenchmarkValidateProductionEntries$' -benchmem`.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `go vet ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

No production-code optimization was made. L14.1 intentionally establishes measurement before changing implementation.

#### Code Review

An experienced Go engineer would approve the benchmark as a useful first baseline. It measures a real import path, reports allocations and avoids obvious benchmark mistakes such as generating test input inside the measured loop.

The main caveat is that this is still a package-level benchmark, not an end-to-end import benchmark with PostgreSQL persistence. That is acceptable for the first lesson; later profiling should decide whether validation, persistence or HTTP handling is the real bottleneck.

#### Exercises

- Run the benchmark three times and compare result variance.
- Move CSV generation inside the measured loop and explain why the benchmark becomes less useful.
- Add a mixed valid/invalid CSV benchmark and compare allocations with the all-valid case.

#### Interview Questions

- Why does `b.N` exist instead of writing a fixed loop count?
- Why should benchmark input generation usually stay outside the measured loop?
- What does `B/op` measure, and why does it matter to the garbage collector?
- Why should a benchmark result lead to profiling before optimization?

#### Roadmap Update

- Lesson 14.1 completed.
- Current lesson moved to Lesson 14.2.
- Testing `Benchmarks` and Performance `benchmarks` marked complete in the Knowledge Matrix.
