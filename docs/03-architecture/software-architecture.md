# Software Architecture

## Purpose

Define the code boundaries needed to implement MESLite's Core, Edge, and
Terminal architecture without duplicating the MES business domain.

This document currently defines the architectural rules established for
component/code sharing. Detailed persistence, API, transaction, and
module patterns must be added when those decisions are made; do not
invent them from this document.

## Repository strategy

Use a **single repository (monorepo)** for Core, Edge, Terminal, shared
contracts, and shared backend domain/application code.

Goals: - keep cross-component changes atomic; - keep contracts
synchronized; - maximize reuse between Core and Edge; - make
architecture easy for humans and AI agents to navigate; - avoid
duplicated implementations of MES business rules.

Exact directory names are not yet fixed. Preserve the logical boundaries
below regardless of physical layout.

## Backend domain ownership

MES business rules have **one backend implementation**.

Core and Edge may share: - domain types and rules; - value objects; -
command definitions/handlers where behavior is identical; - application
logic required in both modes; - validation and business invariants; -
synchronization-related contracts.

Do not create separate `core-domain` and `edge-domain` implementations
of the same MES rules.

Core remains authoritative at runtime in deployments where Edge is
present. Code sharing does not imply equal data authority.

## Terminal boundary

Terminal is a client application and must **not duplicate or import the
backend MES domain implementation**.

Terminal should share stable external contracts needed to communicate
with Core/Edge, such as: - command/API schemas; - identifiers and
wire-format enums where appropriate; - request/response DTO definitions
or generated equivalents; - protocol/version information.

Backend internals must not leak into Terminal merely for code reuse.

In particular, Terminal must not become responsible for authoritative: -
order/routing invariants; - production consistency rules; - database
transaction rules; - synchronization reconciliation rules; - server
authorization decisions.

These rules must be enforced by the server side even when Terminal
performs UX validation.

## Terminal-local model

Terminal may have its own small client-side model for concerns that
genuinely belong to the device/UI:

-   current operator session;
-   workstation binding;
-   pending commands;
-   command delivery state;
-   connectivity state;
-   offline UX;
-   durable local persistence;
-   cached display data required for supported offline operation.

This is a **client/offline model**, not a second MES domain.

## Core and Edge implementation

Core and Edge are server-side roles. Prefer reuse over separate
products.

A desirable direction is one shared server codebase that can be
composed/configured for deployment roles such as:

``` text
cloud/core
edge
standalone/on-premise
```

However, whether these roles are delivered as one executable with modes
or as multiple executables is **not yet decided**. Do not assume either
approach without an explicit decision.

Regardless of packaging:

-   shared behavior must come from shared modules, not copy/paste;
-   Edge-specific buffering/cache/sync code must remain isolated from
    authoritative Core persistence concerns;
-   cloud-specific multi-site/reporting concerns must not be required
    for factory-local production continuity;
-   standalone/on-premise deployment must not require cloud
    connectivity.

## Dependency rules

Conceptual dependency direction:

``` text
                    shared MES domain/application
                         /               \
                    Core role          Edge role
                         \               /
                         external contracts
                                |
                             Terminal
                                |
                       terminal-local model
```

Rules:

1.  Domain/business logic must not depend on Terminal/UI code.
2.  Domain/business logic must not depend on deployment topology.
3.  Terminal may depend on external contracts, never backend
    implementation internals.
4.  Edge-specific synchronization infrastructure must not become part of
    the MES domain model.
5.  Core must validate authoritative business rules regardless of
    validation already performed by Edge or Terminal.
6.  Shared contracts must remain explicit and versionable; do not use
    internal persistence models as API contracts.
7.  Infrastructure concerns must not redefine domain semantics.

## Command-oriented production changes

Production-changing actions should cross unreliable boundaries as
explicit commands/events with stable identifiers rather than as
replicated database state.

Example concept:

``` text
RegisterProduction
- command_id
- operator_id
- workstation_id
- order/operation identity
- quantity
- timestamp/context required by the contract
```

Exact schemas are not defined here.

This supports: - durable Terminal queues; - durable Edge journals; -
retries; - idempotent Core processing; - auditability; - simpler
synchronization than multi-master database replication.

See `offline-sync.md`.

## Persistence roles

Technology choices are intentionally not fixed by this document.

Required logical roles are:

-   **Core persistence:** authoritative MES data.
-   **Edge persistence:** durable pending-command journal plus bounded
    cached working data.
-   **Terminal persistence:** durable pending-command queue plus
    necessary client/offline state.

Do not select or couple these roles to PostgreSQL, SQLite, IndexedDB, or
another technology unless separately decided. Earlier architectural
discussion identified these as candidates, not final choices.

## Source-of-truth rules for implementation

When implementing a feature:

1.  Determine whether the rule is an MES business rule, Edge
    synchronization concern, or Terminal UX/offline concern.
2.  Put MES business rules in the shared backend domain/application
    boundary.
3.  Put network/durability forwarding behavior in Edge/Terminal
    infrastructure as appropriate.
4.  Expose only explicit contracts across component boundaries.
5.  Never solve synchronization by allowing multiple components to
    silently own overlapping authoritative state.

## Decisions still intentionally open

The conversation has **not** yet fully defined: - concrete
repository/directory layout; - programming languages/frameworks as
binding architecture; - database technologies; - one server executable
vs separate Core/Edge executables; - detailed API style; -
persistence/repository patterns; - transaction boundaries; -
schema/migration strategy; - dependency injection/composition
conventions.

AI agents must not infer these as settled architecture. Record them here
and/or in ADRs when decided.
