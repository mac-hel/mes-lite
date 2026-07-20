## Assistant

## Milestone 0 — Lesson 1: Project Foundation

### Business Context

Before writing a single line of business logic, we must establish the foundation that every production Go project shares. This is not "wasted time" — it's the difference between a toy project and a production service. Companies invest heavily in this setup because it pays back in developer productivity, code quality, and deployment reliability.

### Problem

We have an empty directory. We need:

1. **Version control** — track every change, enable collaboration, support review
2. **Project structure** — a consistent layout that makes code predictable
3. **Build system** — compile, test, lint, run with single commands
4. **Database** — a reproducible PostgreSQL environment for development
5. **Static analysis** — catch bugs and non-idiomatic code automatically
6. **Live reload** — edit code and see changes instantly during development
7. **CI** — every commit is verified automatically

### Design Discussion

Let's start with the first two deliverables: **Git repository** and **project structure**.

**Why Git?** Non-negotiable. Every professional project uses version control. Even for a solo project, it enables experimentation, rollback, and a history of decisions.

**Why this structure?**

```
mes-lite/
├── cmd/            # Application entry points
│   └── server/     # The HTTP server binary
├── internal/       # Private application code
├── docs/           # ADRs and documentation
├── migrations/     # Database migrations
├── .github/        # CI configuration
├── docker-compose.yml
├── Makefile
├── go.mod
└── README.md
```

The `cmd/` pattern is the Go standard. Each subdirectory is a `main` package that produces a binary. `internal/` enforces at the compiler level that nothing outside the project can import our code — this is Go's mechanism for package encapsulation.

Let me implement it.

**Tool: todowrite**

---

## Assistant

**Tool: bash**

---

## Assistant

**Tool: bash**

---

## Assistant

Now let me create the project directory structure and initialize Go.

**Tool: bash**

---

## Assistant

**Tool: write**

**Tool: bash**

---

## Assistant

Good. Now let me create the remaining Milestone 0 files. I'll start with the foundational files in parallel.

**Tool: write**

**Tool: write**

**Tool: write**

**Tool: write**

---

## Assistant

Now let me create the minimal Go entry point and CI configuration.

**Tool: write**

**Tool: write**

**Tool: write**

---

## Assistant

Let me update the main.go to use slog instead of fmt.Println to match the linter rules, and create a test file.

**Tool: write**

**Tool: write**

**Tool: write**

---

## Assistant

Now let me verify everything builds and tests pass.

**Tool: bash**

**Tool: bash**

**Tool: bash**

---

## Assistant

Build and tests pass. Let me check linter availability and set up a pre-commit hook.

**Tool: bash**

---

## Assistant

**Tool: edit**

---

## Assistant

**Tool: bash**

---

## Assistant

**Tool: bash**

---

## Assistant

**Tool: bash**

---

## Assistant

**Tool: bash**

---

## Assistant

**Tool: bash**

---

## Assistant

**Tool: read**

---

## Assistant

**Tool: edit**

---

## Assistant

**Tool: edit**

---

## Assistant

**Tool: bash**

---

## Assistant

The forbidigo linter is too aggressive. Let me fix the config.

**Tool: write**

---

## Assistant

**Tool: bash**

---

## Assistant

**Tool: edit**

---

## Assistant

**Tool: bash**

---

## Assistant

Linter is clean. Let me set up a pre-commit hook and create a `.gitignore` for the `.hooks` directory.

**Tool: write**

---

## Assistant

**Tool: bash**

---

## Assistant

Now let me make the initial commit.

**Tool: bash**

---

## Assistant

**Tool: bash**

---

## Assistant

**Tool: edit**

---

## Assistant

**Tool: bash**

---

## Assistant

Commit successful. Now let me update the ROADMAP.md state and Knowledge Matrix to reflect progress.

**Tool: read**

---

## Assistant

The ROADMAP.md was shared in the conversation but doesn't exist as a file yet. Let me create it.

**Tool: write**

---

## Assistant

**Tool: bash**

---

## Assistant

Let me do a final verification — build, test, lint, all together.

**Tool: bash**

---

## Assistant

All checks pass — build, tests (with race detector), and linter all green.

## Lesson 1 Summary — Project Foundation

### What was built

| Deliverable | Status |
|---|---|
| Git repository (`main` branch) | ✅ |
| Project structure (`cmd/`, `internal/`, `docs/`, `migrations/`) | ✅ |
| Go module initialized | ✅ |
| `Makefile` (build, test, lint, run, docker, etc.) | ✅ |
| `golangci-lint` configuration | ✅ |
| `air` live-reload configuration | ✅ |
| `Docker Compose` (PostgreSQL 17) | ✅ |
| GitHub Actions CI (build → lint → test) | ✅ |
| Pre-commit hook (`go vet`) | ✅ |
| `README.md` | ✅ |
| `ROADMAP.md` with updated state | ✅ |

### Go Concepts Learned

- **Go toolchain** — `go build`, `go test`, `go mod`, `go fmt`, `go vet`
- **Modules** — `go.mod` as dependency manifest, module path as import prefix
- **Packages** — `package main` for binaries, `package <name>` for libraries
- **`internal/`** — compiler-enforced encapsulation boundary
- **`cmd/` pattern** — standard Go project layout for multiple binaries

### Go Idioms

- **`internal/` packages** — Go's mechanism for limiting visibility without a monorepo tool
- **Slog structured logging** — standard library, production ready, no framework needed

### Architecture Concepts

- **Repository layout** — convention over configuration, predictable locations
- **Package visibility** — `internal/` prevents external imports; unexported identifiers within packages
- **Composition root** — `cmd/server/main.go` is where wiring happens

### Common Mistakes Avoided

- Not using `internal/` → every package becomes a public API
- Mixing binary and library code in the same package → split `cmd/` and `internal/`
- Ignoring linter from day one → catches non-idiomatic code immediately

### Exercises

For self-practice:
1. Run `make build` — what binary is produced and where?
2. Run `go doc cmd/server` — why does it say nothing? What would you need to change?
3. Add a new command: create `cmd/migrate/main.go` that prints "running migrations..."
4. Explain in your own words: why does Go have `internal/` but PHP/Symfony doesn't?

### Interview Questions

**Q: Why Go modules?**
A: They provide dependency management with versioning, reproducible builds via `go.sum`, and avoid the legacy `GOPATH` mode. Unlike `composer.json` / `package.json`, Go modules use Minimal Version Selection (MVS) rather than constraint solving, which is simpler and deterministic.

**Q: Why `package main`?**
A: `main` is a special package name that tells the Go toolchain "this produces an executable." The `main()` function is the entry point. Non-`main` packages produce archives that can only be imported.

**Q: Why `internal/`?**
A: Packages under `internal/` can only be imported by code rooted at the parent. This enforces encapsulation boundaries at the compiler level, preventing external consumers from depending on private implementation details.

**Q: Why `cmd/`?**
A: Standard convention for organizing multiple binaries in a single module. Each subdirectory under `cmd/` is a `main` package producing a separate binary. This avoids the anti-pattern of one giant binary with build tags.

---

Ready to proceed with **Lesson 2 of Milestone 0** (verify Docker Compose works, and perhaps write the first ADR), or move to **Milestone 1** (Bootstrap HTTP Service)?

---

