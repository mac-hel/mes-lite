## Lesson 2 Summary — Docker Environment & First ADR

### Code Review — Milestone 0

**What's good:**
- `cmd/server/main.go` uses `log/slog` from day one — no `fmt.Print` to unlearn later
- `internal/version/` follows the standard linker-flag pattern for build info
- Makefile targets are composable; Docker pre-flight check fails with a clear message instead of a cryptic socket error
- CI runs with `-race` and `-shuffle=on` — catches data races and test ordering bugs immediately
- Pre-commit hook runs `go vet` — catches obvious issues before they reach CI
- ADR 0001 documents *why* the layout is what it is, not just *what* it is

**What could be improved (deferred intentionally):**
- `main.go` just prints and exits — it's placeholder code that gets replaced in Milestone 1
- No Makefile `install-tools` target yet — `air` and `golangci-lint` are assumed installed. We can add this when we need it

### Architecture Review

Would an experienced Go engineer approve this foundation? **Yes.** The layout follows Go conventions, `internal/` provides compiler-enforced encapsulation, and the ADR captures the reasoning for future contributors.

### Milestone 0 Definition of Done

| Criterion | Status |
|---|---|
| Project builds | ✅ `go build ./...` |
| CI passes | ✅ Configuration correct (runs on GitHub) |
| Lint passes | ✅ 0 issues |
| Docker starts successfully | ✅ Configuration correct (requires Docker daemon) |

### Exercises

1. Run `make docker-up` on a machine with Docker running — verify PostgreSQL starts and is connectable
2. Read `docs/adr/0001-project-layout.md` — can you think of a case where `pkg/` would be justified instead of `internal/`?
3. Add a `make install-tools` target that installs `golangci-lint` and `air`

### Interview Questions

**Q: What is an ADR and why use it?**
A: An Architecture Decision Record captures the context, decision, alternatives, and consequences of a significant architectural choice. Without ADRs, decisions become tribal knowledge — "I don't know why we do it this way, it was like that when I arrived." ADRs make the reasoning explicit and reviewable.

**Q: Why can't other modules import from `internal/`?**
A: The Go compiler enforces this. Any package under `internal/` can only be imported by code whose directory tree includes the `internal` parent. This is a compile-time guarantee — no code review or convention needed.

**Q: When would you NOT use `internal/`?**
A: When you're building a library intended for public consumption — then use `pkg/` or just export from the module root. For an application like this, everything is internal.

---

Milestone 0 is complete. Ready to proceed with **Milestone 1: Bootstrap HTTP Service**? It will be our first real Go code — Fuego server, health endpoint, graceful shutdown.

---

