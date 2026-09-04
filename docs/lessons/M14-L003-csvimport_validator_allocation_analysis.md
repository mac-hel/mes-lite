### Lesson 14.3 Scope

Use allocation profiles and Go escape analysis to understand why the CSV validation benchmark allocates, then make only small allocation improvements that profiles and benchmarks prove worthwhile.

#### Business Context

Large imports create many short-lived objects. Even when CPU time is acceptable, unnecessary allocations increase garbage collector work and can hurt latency under load.

#### Problem

The memory profile shows allocation pressure, but it does not explain which source values escape to the heap or whether any allocation is avoidable without damaging readability.

#### Design Discussion

L14.3 should combine `pprof` allocation output with compiler escape-analysis output from `go test -gcflags=-m`. The lesson should distinguish unavoidable allocations from avoidable ones.

The likely candidates are result-slice growth in collect-all validation and benchmark input construction. CSV parser allocations are probably expected. Any optimization should be validated by comparing benchmark results before and after the change.

#### Go Concepts

- escape analysis with `-gcflags=-m`
- stack versus heap allocation
- allocation profiles versus compiler escape reports
- slice growth and capacity
- garbage collector pressure from short-lived objects

#### Architecture Concepts

- allocation reduction without API distortion
- preserving streaming production paths over optimizing test helpers
- benchmark comparison before and after changes

#### Tests

- Run escape analysis for `internal/csvimport`.
- Inspect allocation profile source lines with `pprof`.
- If code changes are made, compare benchmark results before and after.
- Keep correctness tests, build, vet and lint passing.

#### Exercises

- Explain why returning a slice usually means its backing array escapes.
- Preallocate a result slice in a controlled benchmark and compare `B/op`.
- Identify one allocation that is worth accepting for readability.

#### Interview Questions

- What is escape analysis?
- What makes a value escape to the heap?
- Why are fewer allocations not automatically better if the code becomes harder to maintain?
- How does slice capacity affect allocations?

### Lesson 14.3 Completion Notes

#### Business Context

MES Lite's CSV import reader now avoids one allocation class per CSV row while keeping the same streaming API and import behavior.

#### Problem

The L14.2 memory profile showed allocation pressure in CSV validation. The next step was to distinguish unavoidable allocations from avoidable ones and verify any change with benchmarks instead of assuming it helped.

#### Design Discussion

The allocation profile showed that `ValidateProductionEntries` allocates heavily when appending every valid record into a result slice. That is expected for this collect-all helper. The production import service already avoids this shape by validating and saving records in bounded batches.

The more useful finding was in `encoding/csv.(*Reader).readRecord`, which dominated allocated objects. Go's standard CSV reader supports `ReuseRecord`, which reuses the returned `[]string` backing array between reads. That is safe here because `ProductionEntryReader.Read` immediately copies fields into a `ProductionEntryRow` struct and never exposes the raw CSV record slice to callers.

Escape analysis confirmed the expected heap allocations: returned readers, result slices, error formatting paths and append-backed collections escape. Those are normal for this API shape. The chosen change only removes avoidable per-row CSV record-slice allocation.

#### Implementation

- Enabled `csv.Reader.ReuseRecord` in `NewProductionEntryReader`.
- Kept `ProductionEntryReader.Read` API unchanged.
- Kept validation and import service behavior unchanged.
- Left collect-all `ValidateProductionEntries` unchanged because the production service already uses bounded streaming batches.

#### Benchmark Results

Before `ReuseRecord` on `10000_rows`:

- `13507771 ns/op`, `6126466 B/op`, `20038 allocs/op`.

After `ReuseRecord`:

- `100_rows`: `112674 ns/op`, `43640 B/op`, `127 allocs/op`.
- `1000_rows`: `1071355 ns/op`, `306041 B/op`, `1030 allocs/op`.
- `10000_rows`: `12707274 ns/op`, `5166480 B/op`, `10038 allocs/op`.

The important improvement is allocation-related: 10,000-row validation now performs about 10,000 fewer allocations and allocates about 960 KB less per operation. Runtime changed within benchmark noise, so this lesson does not claim a proven CPU-speed improvement.

#### Tests

- Inspected allocation source lines with `go tool pprof -list ValidateProductionEntries /tmp/opencode/csvimport_l143_mem_before.out`.
- Ran package escape analysis with `go build -gcflags='github.com/mac-hel/mes-lite/internal/csvimport=-m=2' ./internal/csvimport`.
- Compared benchmarks with `go test ./internal/csvimport -run '^$' -bench '^BenchmarkValidateProductionEntries/10000_rows$' -benchmem -memprofile /tmp/opencode/csvimport_l143_mem_before.out`.
- Compared after-change memory profile with `go test ./internal/csvimport -run '^$' -bench '^BenchmarkValidateProductionEntries/10000_rows$' -benchmem -memprofile /tmp/opencode/csvimport_l143_mem_after.out`.
- Verified full benchmark suite with `go test ./internal/csvimport -run '^$' -bench '^BenchmarkValidateProductionEntries$' -benchmem`.
- Verified with `go test ./internal/csvimport -count=1`.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `go vet ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

No broad refactor was needed. The single-line `ReuseRecord` change is preferable to introducing a new validation API or preallocation heuristic because it reduces allocations without changing package contracts.

#### Code Review

An experienced Go engineer would approve this optimization because it is profile-guided, small and uses a standard-library feature exactly where ownership makes it safe.

The review caveat is important: `ReuseRecord` would be unsafe if callers received and retained the raw `[]string` from `encoding/csv.Reader.Read`. MES Lite does not expose that slice; it maps fields into a struct before returning.

#### Exercises

- Disable `ReuseRecord` and rerun the benchmark to confirm the allocation difference.
- Add a benchmark for mixed valid and invalid rows and compare error-allocation behavior.
- Explain why preallocating `ValidationResult.Records` is not straightforward when reading from a stream.

#### Interview Questions

- What does it mean for a value to escape to the heap?
- Why did `ValidateProductionEntries` allocate even though records are local variables in the loop?
- Why is `csv.Reader.ReuseRecord` safe in this package but potentially unsafe in another API?
- Why should allocation improvements be benchmarked instead of inferred from code inspection?

#### Roadmap Update

- Lesson 14.3 completed.
- Current lesson moved to Lesson 14.4.
- Runtime `Escape Analysis` and Performance `allocations` marked complete in the Knowledge Matrix.
