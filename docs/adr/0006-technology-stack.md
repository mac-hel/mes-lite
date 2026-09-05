# ADR 0006: Technology Stack

## Status

Accepted

## Context

MESLite needs a documented baseline for runtime, persistence, development, and
quality tooling so implementation work uses consistent defaults instead of
re-deciding tools per feature.

The stack should describe current project choices without replacing narrower
engineering documents such as package layout, sqlc boundaries, validation rules,
or background job notes.

## Decision

Use the following technology stack.

### Runtime

- **Language:** Go; version is declared in `go.mod`
- **HTTP framework:** Fuego
- **Database:** PostgreSQL; local image version is declared in `docker-compose.yml`
- **Database driver:** pgx/v5
- **SQL generation:** sqlc; config in `sqlc.yaml`
- **Migrations:** goose

### Infrastructure

- **Container image:** Docker
- **Local services:** Docker Compose
- **Development reload:** air

### Quality And Tooling

- **Linting:** golangci-lint
- **Formatting:** gofmt and goimports
- **Dependency boundaries:** golangci-lint `depguard`
- **Tests:** Go test runner, race detector, and coverage tooling

### Observability And Operations

- **Metrics:** Prometheus client
- **Tracing:** OpenTelemetry
- **Configuration:** environment variables, with `.env` loading via `godotenv`

## Alternatives Considered

- **Leave the stack undocumented:** rejected because undocumented choices invite
  duplicate decisions and inconsistent implementation patterns.
- **Put the stack only in `README.md`:** rejected because the README is an entry
  point and quick-start document, not the authoritative record for long-lived
  technical decisions.
- **Put the stack only in `development-guide.md`:** rejected because the stack is
  a project decision, while the development guide should reference decisions and
  describe day-to-day workflow without duplicating authority.

## Trade-Offs

- Pro: contributors and AI agents have one place to check current technology
  choices.
- Pro: development documentation can link to the decision instead of duplicating
  the stack.
- Con: version changes now require updating this ADR or superseding it when the
  change is significant.

## Consequences

- New code should use these tools unless a later ADR changes the decision.
- Setup and workflow documentation should reference this ADR rather than copying
  the full stack list.
- If a stack choice changes in a meaningful or difficult-to-reverse way, record
  the change with a new ADR instead of silently editing this decision history.
