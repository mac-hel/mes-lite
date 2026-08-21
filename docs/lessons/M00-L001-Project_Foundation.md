## Lesson 1 Summary — Project Foundation

### Exercises

- create project from scratch
- explain module system
- explain package visibility

For self-practice:
1. Run `make build` — what binary is produced and where?
    - in cwd or in dir pointed by -o argument
2. Run `go doc cmd/server` — why does it say nothing? What would you need to change?
    - package-level doc comment must be placed before `package main`
3. Add a new command: create `cmd/migrate/main.go` that prints "running migrations..."
4. Explain in your own words: why does Go have `internal/` but PHP/Symfony doesn't?
    - go packages are compilation/import boundaries and internal restricts them to sibling dirs.
    - php has namespaces but they are not module/package boundaries, visibility is handled by classes, conventions (@internal), static analysis, framework architecture

### Interview Questions

**Q: Why Go modules?**
A: They provide dependency management with versioning, reproducible builds via `go.sum`, and avoid the legacy `GOPATH` mode, also allow code import via module path.
Unlike `composer.json` / `package.json`, Go modules use Minimal Version Selection (MVS) rather than constraint solving, which is simpler and deterministic.

**Q: Why `package main`?**
A: `main` is a special package name that tells the Go toolchain "this produces an executable." The `main()` function is the entry point. Non-`main` packages produce archives that can only be imported.

**Q: Why `internal/`?**
A: Packages under `internal/` can only be imported by code rooted at the parent. This enforces encapsulation boundaries at the compiler level, preventing external consumers from depending on private implementation details.

**Q: Why `cmd/`?**
A: Standard convention for organizing multiple binaries in a single module. Each subdirectory under `cmd/` is a `main` package producing a separate binary. This avoids the anti-pattern of one giant binary with build tags.
