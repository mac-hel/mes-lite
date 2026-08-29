# MESLite Documentation Guide

This document defines the documentation set for MESLite and how an AI development agent should use it.

## Documentation principles

- Keep documentation small, authoritative, and non-duplicative.
- Each concept should have one primary source of truth.
- Prefer concise rules, examples, diagrams, IDs, and acceptance criteria over long prose.
- Update documentation whenever product behavior, domain rules, architecture, development rules, or operations change.
- Use ADRs for important decisions and their rationale instead of rewriting historical reasoning into current-state documents.

## Structure

```text
/
├── README.md
├── AGENTS.md
└── docs/
    ├── 01-product/
    │   ├── product-context.md
    │   └── requirements.md
    ├── 02-domain/
    │   └── domain.md
    ├── 03-architecture/
    │   ├── system-architecture.md
    │   ├── software-architecture.md
    │   ├── offline-sync.md
    │   └── security.md
    ├── 04-development/
    │   ├── development-guide.md
    │   └── testing-strategy.md
    ├── 05-operations/
    │   └── operations.md
    └── adr/
        └── NNNN-*.md
```

## Documents

### `README.md`

**Purpose:** Repository entry point.

Contains:
- what MESLite is, briefly;
- repository layout;
- local setup and common commands;
- links to authoritative documentation.

**Use:** Read first when entering the repository. Keep operational details concise; do not duplicate architecture or business documentation.

---

### `AGENTS.md`

**Purpose:** Instructions for AI agents working on the repository.

Contains:
- mandatory documents to read before specific types of changes;
- architectural and coding constraints;
- rules for modifying tests, schemas, APIs, and documentation;
- required validation before completing work;
- pointers to authoritative documents.

**Use:** Always read before changing code. Treat as the execution policy for AI-assisted development.

---

### `docs/01-product/product-context.md`

**Merged from:** `product-vision.md`, `business-context.md`, `glossary.md`.

**Purpose:** Explain why MESLite exists and the manufacturing context in which it operates.

Contains:
- product vision and goals;
- target users and customers;
- goals and non-goals;
- actors and business environment;
- important manufacturing terminology.

**Use:** Read when evaluating scope, UX, product behavior, terminology, or whether a proposed feature belongs in MESLite.

---

### `docs/01-product/requirements.md`

**Purpose:** Define what the system must do.

Contains:
- functional requirements;
- non-functional requirements;
- requirement IDs;
- acceptance criteria;
- reliability, performance, security, and offline requirements.

**Use:** Read before implementing or changing behavior. Link tests and feature work to requirement IDs where practical.

---

### `docs/02-domain/domain.md`

**Merged from:** `domain-model.md`, `business-rules.md`, `workflows.md`.

**Purpose:** Define the MESLite business domain.

Contains:
- entities and value objects;
- relationships;
- terminology;
- lifecycles and state transitions;
- business workflows;
- invariants and numbered business rules.

**Use:** Read before changing business logic, persistence models, workflows, or APIs that expose domain concepts. Domain rules take precedence over implementation convenience.

---

### `docs/03-architecture/system-architecture.md`

**Purpose:** Define the high-level system and deployment architecture.

Contains:
- Core / Edge / Terminal responsibilities;
- deployment boundaries;
- communication paths;
- system-level data ownership;
- major failure assumptions;
- external integrations.

**Use:** Read before changes affecting component responsibilities, deployment topology, network communication, or cross-component behavior.

---

### `docs/03-architecture/software-architecture.md`

**Merged from:** `application-architecture.md`, `data-architecture.md`.

**Purpose:** Define how MESLite code and persistent data are structured.

Contains:
- repository/module structure;
- dependency direction and boundaries;
- domain/application/infrastructure responsibilities;
- API and persistence patterns;
- database responsibilities;
- transactions and consistency boundaries;
- data ownership within the software design;
- shared-code rules for Core, Edge, and Terminal.

**Use:** Read before creating modules, dependencies, repositories, services, schemas, APIs, or persistence logic. New code should follow existing patterns unless an ADR explicitly changes them.

---

### `docs/03-architecture/offline-sync.md`

**Purpose:** Define MESLite's offline, synchronization, and recovery guarantees.

Contains:
- Internet-loss behavior;
- Wi-Fi-loss behavior;
- terminal-local durable storage;
- Edge behavior;
- synchronization protocol;
- retries and idempotency;
- conflict handling;
- replacement/recovery scenarios;
- guarantees against acknowledged production-data loss.

**Use:** Mandatory for any change affecting production recording, synchronization, local storage, retries, connectivity, Edge, or Terminal behavior.

This document is intentionally separate because offline continuity is a core MESLite architectural requirement.

---

### `docs/03-architecture/security.md`

**Purpose:** Define cross-cutting security rules.

Contains:
- authentication;
- authorization and roles;
- secrets;
- auditability;
- data protection;
- trust boundaries;
- relevant security assumptions.

**Use:** Read when changing identity, permissions, APIs, sensitive data, deployment trust, or auditing.

---

### `docs/04-development/development-guide.md`

**Merged from:** `coding-guidelines.md`, `development-workflow.md`, `definition-of-done.md`.

**Purpose:** Define how code changes should be made.

Contains:
- coding conventions not enforced automatically;
- preferred implementation patterns;
- dependency and library rules;
- development workflow;
- migration/API-change rules;
- documentation-update rules;
- definition of done;
- required checks before completion.

**Use:** Read before implementing changes and again before declaring work complete.

---

### `docs/04-development/testing-strategy.md`

**Purpose:** Define how MESLite is verified.

Contains:
- unit, integration, API, synchronization, and end-to-end test boundaries;
- required tests for business rules;
- database-testing approach;
- resilience/failure scenarios;
- offline and recovery tests;
- test naming and organization where needed.

**Use:** Read when adding or changing behavior and when deciding what tests are required.

---

### `docs/05-operations/operations.md`

**Merged from:** `deployment.md`, `observability.md`, `backup-recovery.md`.

**Purpose:** Define how MESLite is deployed and operated reliably.

Contains:
- deployment environments;
- configuration;
- upgrades and migrations;
- logging, metrics, and diagnostics;
- backup and restore;
- disaster recovery;
- operational failure handling.

**Use:** Read when changing deployment, infrastructure, configuration, observability, persistence recovery, or upgrade behavior.

Split this document later only if operational complexity justifies it.

---

### `docs/adr/NNNN-*.md`

**Purpose:** Record important architectural or technical decisions and why they were made.

Each ADR should contain:
- context;
- decision;
- rationale;
- consequences;
- alternatives considered;
- status and date.

Examples:
- technology stack;
- monorepo choice;
- Core / Edge architecture;
- database choices;
- synchronization strategy;
- UI architecture.

**Use:** Read relevant ADRs before changing established architectural decisions. Create a new ADR when making a significant, long-lived, or difficult-to-reverse decision.

Do not edit accepted ADRs to hide historical decisions; supersede them with a new ADR when necessary.

## Relationship between documents

```text
product-context.md
    WHY / FOR WHOM
        ↓
requirements.md
    WHAT MUST HAPPEN
        ↓
domain.md
    BUSINESS MEANING + RULES
        ↓
system-architecture.md
    WHICH SYSTEM COMPONENT IS RESPONSIBLE
        ↓
software-architecture.md
    HOW CODE + DATA ARE STRUCTURED
        ↓
offline-sync.md / security.md
    CROSS-CUTTING ARCHITECTURAL CONSTRAINTS
        ↓
development-guide.md
    HOW CHANGES ARE IMPLEMENTED
        ↓
testing-strategy.md
    HOW CORRECTNESS IS VERIFIED
        ↓
operations.md
    HOW THE SYSTEM IS RUN

ADRs explain WHY major architecture choices were made.
AGENTS.md tells the AI HOW to use all of the above.
```

## AI development workflow

For every task:

1. Read `AGENTS.md`.
2. Identify the affected requirements and domain concepts.
3. Read only the architecture documents relevant to the change.
4. Check relevant ADRs before altering established design.
5. Implement according to `development-guide.md`.
6. Add tests according to `testing-strategy.md`.
7. Update affected documentation when behavior or architecture changes.
8. Validate against the definition of done before completing the task.

### Minimum reading by task type

| Change | Required documents |
|---|---|
| Product/UX behavior | `product-context.md`, `requirements.md`, relevant `domain.md` sections |
| Business logic | `requirements.md`, `domain.md`, `software-architecture.md` |
| API | `requirements.md`, `domain.md`, `software-architecture.md`, `security.md` when relevant |
| Database/schema | `domain.md`, `software-architecture.md`, relevant ADRs |
| Core/Edge/Terminal boundaries | `system-architecture.md`, `software-architecture.md`, relevant ADRs |
| Offline/sync/production persistence | `requirements.md`, `domain.md`, `system-architecture.md`, `offline-sync.md` |
| Authentication/permissions | `requirements.md`, `security.md`, `software-architecture.md` |
| Deployment/infrastructure | `system-architecture.md`, `operations.md`, relevant ADRs |
| Tests | affected specification documents + `testing-strategy.md` |

## Source-of-truth rule

If documents appear to conflict, do not silently choose one.

Use this order to determine intent:

1. explicit current requirement;
2. current domain rule;
3. current architecture document;
4. accepted ADRs explaining architecture decisions;
5. implementation.

Code is not the authoritative specification when it contradicts current documented requirements or rules. Resolve meaningful conflicts explicitly and update the affected documentation.
