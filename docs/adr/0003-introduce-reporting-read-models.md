# ADR 0003: Introduce Reporting Read Models

## Status

Accepted

## Context

Milestone 8 adds reporting over existing operational data. Reports need grouped and aggregated query shapes that do not match command-oriented domain aggregates such as production entries or production orders.

Putting report SQL into write-model repositories would couple command persistence to management views and make those repositories grow for unrelated reasons.

## Decision

Introduce `internal/reporting` as a vertical slice for reporting read models and query stores.

The reporting slice owns SQL queries optimized for reads. It returns projection structs such as `DailyProductionRow` instead of domain aggregates.

This is a pragmatic CQRS split: commands remain in existing business slices, while reporting owns read-specific models. No command bus, query bus or generic CQRS framework is introduced.

## Alternatives Considered

- Add report methods to `internal/production`: rejected because production registration is a write workflow and reports will grow different query concerns.
- Create a generic query service package: rejected because one reporting slice is enough and generic abstractions would be premature.
- Introduce materialized views immediately: rejected because current data volume and query complexity do not justify cached projections yet.

## Trade-Offs

- Pro: report SQL can evolve independently from command repositories.
- Pro: read models can match API/reporting needs without weakening domain aggregates.
- Pro: package boundaries remain business-oriented.
- Con: there is another package and sqlc output to maintain.
- Con: duplicated table knowledge appears in reporting SQL, which is acceptable because reporting intentionally reads across slices.

## Consequences

- Future report APIs should depend on the reporting slice, not on command repositories.
- Query performance review belongs in reporting milestones.
- Materialized views or indexes can be introduced later when measurement shows a need.
