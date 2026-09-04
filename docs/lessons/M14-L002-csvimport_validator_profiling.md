### Lesson 14.2 Scope

Use `pprof` CPU and memory profiles on the CSV validation benchmark to learn where time and allocations are spent before deciding whether any optimization is justified.

#### Business Context

Benchmarks show how fast a code path is. Profiles explain where that time and memory go. CSV import is a natural target because it may process large historical files.

#### Problem

The L14.1 benchmark produced numbers, but numbers alone do not identify the hot functions or allocation sources. Optimizing from benchmark totals alone would still be guesswork.

#### Design Discussion

L14.2 should generate CPU and memory profiles from the existing benchmark using Go's built-in benchmark flags. The lesson should inspect the profiles with `go tool pprof` and document the top costs before making any code changes.

Optimization remains optional. If the profile shows the current implementation is acceptable for MVP scale, the correct engineering decision may be to keep the code unchanged.

#### Go Concepts

- `runtime/pprof` through `go test` profile flags
- CPU profiles versus memory profiles
- interpreting flat and cumulative cost
- allocation profiles and garbage collector pressure
- measurement-driven optimization discipline

#### Architecture Concepts

- profiling before optimization
- separating local package hot spots from end-to-end system bottlenecks
- preserving readability unless profiling proves a change is worthwhile

#### Tests

- Generate a CPU profile from the CSV validation benchmark.
- Generate a memory profile from the CSV validation benchmark.
- Inspect top profile entries with `go tool pprof`.
- Keep correctness tests, build, vet and lint passing.

#### Exercises

- Compare CPU profile output for 1,000 and 10,000 row benchmarks.
- Inspect allocation sources and predict which ones are unavoidable because CSV parsing creates strings.
- Decide whether a code optimization is justified or whether the benchmark should simply remain as a regression guard.

#### Interview Questions

- What is the difference between benchmarking and profiling?
- What do flat and cumulative pprof costs mean?
- Why can reducing allocations improve latency even when CPU time looks acceptable?
- Why should readability usually win unless profiling proves a hot path needs optimization?

### Lesson 14.2 Completion Notes

#### Business Context

MES Lite now has CPU and memory profiles for the CSV import validation benchmark. The project can explain where benchmark time and allocations go instead of treating the benchmark result as a black box.

#### Problem

The L14.1 benchmark showed timing and allocation totals, but did not identify hot functions or allocation sources. Optimizing from totals alone would still be guesswork.

#### Design Discussion

Generated CPU and memory profiles from the `10000_rows` CSV validation benchmark and inspected them with `go tool pprof`.

The CPU profile showed most application-relevant time in CSV parsing and validation orchestration. The memory profile showed two different stories depending on the view:

- `alloc_objects`: `encoding/csv.(*Reader).readRecord` produced almost all allocated objects, which is expected because CSV parsing creates records and field strings.
- `alloc_space`: `ValidateProductionEntries` dominated allocation space because it collects all valid records into a result slice.

No production-code optimization was made. The collect-all validation helper is not the current production import path; the real CSV import service already streams records and persists bounded batches. Optimizing the helper before profiling the end-to-end async/synchronous import path would risk making test-support code more complex without improving production behavior.

#### Profile Results

- CPU profile generated with `-cpuprofile /tmp/opencode/csvimport_cpu.out`.
- Memory profile generated with `-memprofile /tmp/opencode/csvimport_mem.out`.
- Benchmark profile run: `BenchmarkValidateProductionEntries/10000_rows-8`, `6433350 ns/op`, `6126464 B/op`, `20038 allocs/op`.
- CPU top cumulative application path: `ValidateProductionEntries` at `1.69s`, `ProductionEntryReader.Read` at `1.13s`, `encoding/csv.(*Reader).readRecord` at `1.04s`.
- Allocation-space top entries: `ValidateProductionEntries` at `1156.69MB` flat, `encoding/csv.(*Reader).readRecord` at `393.03MB` flat.
- Allocation-object top entry: `encoding/csv.(*Reader).readRecord` at `5134045` objects.

#### Tests

- Verified with `go test ./internal/csvimport -run '^$' -bench '^BenchmarkValidateProductionEntries/10000_rows$' -benchmem -cpuprofile /tmp/opencode/csvimport_cpu.out -memprofile /tmp/opencode/csvimport_mem.out`.
- Inspected CPU with `go tool pprof -top /tmp/opencode/csvimport_cpu.out`.
- Inspected allocation space with `go tool pprof -top -alloc_space /tmp/opencode/csvimport_mem.out`.
- Inspected allocation objects with `go tool pprof -top -alloc_objects /tmp/opencode/csvimport_mem.out`.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `go vet ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

No code refactor was applied. The correct performance decision for this lesson was restraint: profiles identified costs, but they did not justify complicating production code.

#### Code Review

An experienced Go engineer would approve the decision not to optimize yet. The profiles show expected costs in standard-library CSV parsing and in a collect-all helper that is not the runtime import service's bounded streaming path.

The next useful review is allocation-focused: L14.3 should use escape analysis and allocation inspection to understand which values move to the heap and whether any small benchmark-side or production-side allocation reduction is both measurable and readable.

#### Exercises

- Re-run the CPU profile with `-benchtime=5s` and compare the stability of the top entries.
- Use `go tool pprof -list ValidateProductionEntries` to connect allocation space to exact source lines.
- Profile the synchronous CSV import service path and compare it with the validation-only benchmark.

#### Interview Questions

- What information does profiling add that benchmarking alone does not?
- Why can `alloc_space` and `alloc_objects` point at different functions?
- Why is it risky to optimize a helper that is not the production hot path?
- When is not changing code the correct result of a performance investigation?

#### Roadmap Update

- Lesson 14.2 completed.
- Current lesson moved to Lesson 14.3.
- Performance `pprof` marked complete in the Knowledge Matrix.
