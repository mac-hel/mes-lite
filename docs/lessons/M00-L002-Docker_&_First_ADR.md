## Lesson 2 Summary — Docker Environment & First ADR

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
