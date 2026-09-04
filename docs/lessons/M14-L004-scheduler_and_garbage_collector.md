### Lesson 14.4 Scope

Review Go runtime behavior using the CSV validation benchmark as evidence: scheduler basics, garbage collector pressure and how allocation rate affects runtime work.

#### Business Context

Performance-sensitive import paths do not run in isolation. High allocation rates can increase garbage collector work, and concurrent background imports share CPU time with HTTP request handling.

#### Problem

The project has benchmark and allocation data, but the runtime concepts behind those numbers have not been reviewed: what the scheduler does, when garbage collection runs and why allocation reduction can matter even when wall-clock time looks similar.

#### Design Discussion

L14.4 should be a review-and-measurement lesson, not a feature lesson. It should use benchmark output, `GODEBUG=gctrace=1` and runtime documentation-level analysis to connect code behavior with scheduler and garbage collector concepts.

Any code change should be avoided unless the runtime evidence reveals a simple correctness or lifecycle issue. The goal is interview readiness and production reasoning, not micro-optimization.

#### Go Concepts

- goroutine scheduling overview
- GOMAXPROCS and CPU parallelism
- garbage collector roots and heap growth
- allocation rate versus GC frequency
- interpreting `GODEBUG=gctrace=1` output

#### Architecture Concepts

- runtime behavior as part of production performance review
- avoiding unbounded background work that competes with request handling
- connecting allocation profiles to garbage collector pressure

#### Tests

- Run the CSV validation benchmark with GC trace output.
- Compare benchmark behavior with different `GOMAXPROCS` values if useful.
- Keep correctness tests, build, vet and lint passing.

#### Exercises

- Explain why a lower allocation count can reduce GC pressure even if `ns/op` is similar.
- Run the benchmark with `GOMAXPROCS=1` and compare the result with the default.
- Identify one place where unbounded goroutines would compete with request handling.

#### Interview Questions

- What does the Go scheduler schedule?
- What is `GOMAXPROCS`?
- Why is Go's garbage collector concurrent instead of stop-the-world for the whole collection?
- How can allocation-heavy code affect latency?

### Lesson 14.4 Completion Notes

#### Business Context

MES Lite now connects benchmark allocation data to runtime behavior. This matters because large CSV imports can create garbage collector work and compete for CPU with HTTP request handling and background workers.

#### Problem

Benchmarks, pprof and escape analysis showed where time and allocations happen, but they did not yet explain how the Go runtime reacts to allocation-heavy code or how scheduler parallelism changes benchmark execution.

#### Design Discussion

L14.4 used the existing CSV validation benchmark as runtime evidence instead of adding a feature. The benchmark was run with `GODEBUG=gctrace=1` to observe garbage collection and with different `GOMAXPROCS` values to observe scheduler capacity effects.

The GC trace showed frequent collections during the benchmark because each operation still allocates around 5.1 MB. The trace also showed small stop-the-world phases around concurrent GC work, which is the important production lesson: Go's collector is concurrent, but allocation-heavy code still creates runtime work and some pauses.

The `GOMAXPROCS` comparison showed that CPU parallelism can improve throughput for benchmark execution, but it does not change the allocation shape. Both runs kept `5166xxx B/op` and `10038 allocs/op`.

#### Runtime Results

- Default `GOMAXPROCS=8` with GC tracing: `12585191 ns/op`, `5166480 B/op`, `10038 allocs/op`.
- `GOMAXPROCS=1`: `14300800 ns/op`, `5166463 B/op`, `10038 allocs/op`.
- `GOMAXPROCS=8`: `12210447 ns/op`, `5166486 B/op`, `10038 allocs/op`.
- GC trace during the benchmark showed hundreds of small collections over the run, with reported GC CPU around `3%` in the main benchmark process.

#### Tests

- Ran `GODEBUG=gctrace=1 go test ./internal/csvimport -run '^$' -bench '^BenchmarkValidateProductionEntries/10000_rows$' -benchmem -benchtime=2s`.
- Ran `GOMAXPROCS=1 go test ./internal/csvimport -run '^$' -bench '^BenchmarkValidateProductionEntries/10000_rows$' -benchmem -benchtime=2s`.
- Ran `GOMAXPROCS=8 go test ./internal/csvimport -run '^$' -bench '^BenchmarkValidateProductionEntries/10000_rows$' -benchmem -benchtime=2s`.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `go vet ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

No code refactor was applied. The correct output of this lesson was runtime understanding, not another optimization.

#### Code Review

An experienced Go engineer would approve not changing code in this lesson. The runtime evidence reinforces the L14.3 allocation improvement and shows where future work should focus: reduce allocation rate only when the production import path proves it matters.

The main architecture reminder is that background imports must stay bounded. Unbounded goroutines or unbounded in-memory buffering would compete with HTTP handlers for scheduler time and increase GC pressure.

#### Exercises

- Run the same benchmark with `GOGC=off` and explain why that is not a production solution.
- Run the benchmark with `GOMAXPROCS=2` and compare it with `1` and `8`.
- Identify one background-job configuration that could overload the scheduler under many imports.

#### Interview Questions

- What is the relationship between goroutines, OS threads and `GOMAXPROCS`?
- Why can allocation-heavy code increase latency even when individual GC pauses are small?
- What do the three heap numbers in a GC trace line roughly represent?
- Why does increasing `GOMAXPROCS` not reduce `allocs/op`?

#### Roadmap Update

- Lesson 14.4 completed.
- Current lesson moved to Lesson 14.5.
- Standard Library `runtime`, Runtime `Scheduler` and Runtime `Garbage Collector` marked complete in the Knowledge Matrix.
