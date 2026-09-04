### Lesson 14.5 Scope

Close Milestone 14 with a performance-engineering review: summarize benchmark, profile, allocation and runtime findings, then define the optimization discipline MES Lite should follow before production readiness.

#### Business Context

Performance work should keep the application responsive as data grows, but it must not make MVP code harder to maintain without evidence. The project needs a clear standard for future performance changes.

#### Problem

Milestone 14 introduced benchmarks, pprof, allocation analysis and runtime review. Before moving to production readiness, those findings should be reviewed together and converted into practical engineering rules.

#### Design Discussion

L14.5 should be a milestone review and discipline lesson. It should not add a new optimization unless review finds a small, proven issue. The expected output is a documented performance baseline, known caveats and criteria for future changes such as PostgreSQL `COPY`, route-template metrics, durable jobs or database profiling.

#### Go Concepts

- comparing benchmark runs responsibly
- separating throughput, latency and allocation goals
- avoiding benchmark overfitting
- recognizing when not to optimize

#### Architecture Concepts

- performance budget thinking
- production-like data before database optimization
- code readability as a performance trade-off
- milestone review before production readiness

#### Tests

- Re-run benchmark and verification commands.
- Review whether the one L14.3 optimization remains justified.
- Keep correctness tests, build, vet and lint passing.

#### Exercises

- Write a performance-change checklist for future PRs.
- Decide when CSV import should move from regular inserts to `COPY`.
- Identify which endpoint should get an end-to-end benchmark next.

#### Interview Questions

- How do you decide whether an optimization is worth it?
- Why can benchmarks lie?
- What is the difference between improving throughput and reducing tail latency?
- Why should production-like data guide SQL performance work?

### Lesson 14.5 Completion Notes

#### Business Context

Milestone 14 is complete. MES Lite now has a repeatable performance baseline and a clear rule for future optimization: measure first, profile second, optimize only when the evidence justifies the complexity.

#### Problem

The project needed to close performance engineering as a discipline, not as a collection of isolated benchmark commands. The final review needed to decide whether the L14.3 optimization remained justified and whether more optimization should happen before production readiness.

#### Design Discussion

The CSV validation benchmark remains useful as a regression guard for a performance-sensitive import path. It is not a replacement for production-scale import testing because it does not include PostgreSQL writes, HTTP upload handling, background-job scheduling or real file I/O.

The only code optimization made in the milestone, enabling `csv.Reader.ReuseRecord`, remains justified. It is small, standard-library based, API-preserving and measurably reduces allocations. No additional optimization was made in L14.5 because no new profile evidence identified a production hot path worth complicating.

The main discipline coming out of the milestone is this checklist for future performance PRs:

- Start with a business-relevant symptom or risk.
- Add or run a benchmark that captures the relevant path.
- Use profiles to identify the real hot spot.
- Prefer small standard-library or algorithmic improvements over clever rewrites.
- Compare before and after using `benchmem`.
- Keep correctness tests, build, vet and lint green.
- Document what improved and what did not.

#### Final Benchmark Baseline

- `100_rows`: `57234 ns/op`, `43640 B/op`, `127 allocs/op`.
- `1000_rows`: `519803 ns/op`, `306041 B/op`, `1030 allocs/op`.
- `10000_rows`: `5708275 ns/op`, `5166468 B/op`, `10038 allocs/op`.

Benchmark timing varies between runs, so allocation numbers are the more stable comparison for the L14.3 change.

#### Tests

- Re-ran `go test ./internal/csvimport -run '^$' -bench '^BenchmarkValidateProductionEntries$' -benchmem`.
- Verified with `go test ./... -count=1`.
- Verified with `go build ./...`.
- Verified with `go vet ./...`.
- Verified with `golangci-lint run ./...`.

#### Refactoring

No L14.5 code refactor was applied. The milestone already made the only justified production change in L14.3.

#### Code Review

An experienced Go engineer would approve Milestone 14 because it improved performance maturity without overfitting the code to a benchmark. The codebase now has benchmark coverage, profile evidence, escape-analysis notes and one measured allocation improvement.

The main caveat is that database and end-to-end API performance are not yet benchmarked. That is acceptable because production-scale performance work should use realistic data volumes and deployment conditions.

#### Exercises

- Add an end-to-end benchmark for synchronous CSV import using an in-memory store.
- Define a production-like CSV import dataset size and expected throughput target.
- Decide what evidence would justify replacing regular batch inserts with PostgreSQL `COPY`.

#### Interview Questions

- How would you investigate a slow Go endpoint in production?
- Why should you avoid optimizing code that is not on the hot path?
- Why can reducing allocations improve latency even if throughput barely changes?
- What makes a benchmark representative enough to trust?

#### Roadmap Update

- Lesson 14.5 completed.
- Milestone 14 completed.
- Current milestone moved to Milestone 15.
- Current lesson moved to Lesson 15.1.
- Milestone 15 divided into lessons.
- Architecture maturity, Go knowledge progress and interview readiness updated.
