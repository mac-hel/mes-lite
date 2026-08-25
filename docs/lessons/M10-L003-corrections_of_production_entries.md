### Lesson 10.3 Scope

Allow authorized users to correct production-entry mistakes without silently editing historical records.

#### Business Context

Manual production records contain mistakes: wrong quantity, wrong workstation, wrong product, wrong employee or wrong timestamp. Leaders and managers need to fix these mistakes while preserving who changed what and why.

#### Problem

The system stored production entries as immutable rows, but it had no correction workflow. Directly updating `production_entries` would hide the original value and weaken trust in production history.

#### Design Discussion

Corrections are append-only records in `production_entry_corrections`. A correction stores the full corrected production-entry snapshot plus `reason`, `actorUserId` and `createdAt`.

The original `production_entries` row is never updated by the correction workflow. This is a narrow audit model for production-entry corrections only, not a generic audit framework.

Admins, managers and leaders can append and review corrections. Workers can register production, but cannot correct historical records.

#### Go Concepts

- DTOs for correction workflows
- request principal extraction from context
- domain validation for append-only audit records
- error translation for correction-specific failures

#### Architecture Concepts

- append-only correction records
- auditability through actor and reason tracking
- correction workflow in the production vertical slice
- route-level RBAC for historical mutation workflows

### Lesson 10.3 Completion Notes

#### Business Context

MES Lite now supports audit-safe correction of production-entry mistakes.

#### Problem

Production records could be entered and reviewed, but mistakes could not be corrected without either leaving the mistake in place or silently changing history.

#### Design Discussion

Added a dedicated correction table and production correction model. Each correction references the original entry, stores corrected replacement values and records the authenticated actor plus a required reason.

The correction endpoint does not update the original entry. Review clients can read correction history separately and decide how to present original versus corrected values.

#### Implementation

- Added migration `0010_create_production_entry_corrections.sql`.
- Added `production.Correction` and `ErrInvalidCorrection`.
- Added `CorrectEntryCommand`.
- Added `Service.CorrectEntry` and `Service.ListCorrections`.
- Added `Store.SaveCorrection` and `Store.ListCorrections`.
- Implemented in-memory correction persistence for fast tests.
- Added sqlc queries for creating and listing corrections.
- Implemented PostgreSQL correction persistence and row mapping.
- Added `POST /production-entries/{id}/corrections`.
- Added `GET /production-entries/{id}/corrections`.
- Tracked correction actor from the authenticated request principal.
- Protected correction routes with admin/manager/leader RBAC.

#### Tests

- Added service test proving corrections append without changing the original entry.
- Added service tests for missing original entries and missing correction reasons.
- Added handler test proving actor tracking comes from request context.
- Added handler test for missing principal rejection.
- Added PostgreSQL test for saving and listing correction history.
- Added PostgreSQL test for missing original-entry correction failure.
- Added server route test proving leaders can correct entries.
- Added server route test proving workers cannot correct entries.
- Added server route test for reading correction history.
- Verified with `go test ./internal/production ./internal/server -count=1`.

#### Refactoring

No generic audit package was introduced. The correction model exists because production-entry correction is a concrete MVP workflow with actor/reason requirements.

After review, the production slice split correction behavior out of the entry-oriented interfaces. `Service` now depends on separate `EntryStore` and `CorrectionStore` interfaces, and `Handler` now depends on separate `EntryRegistrar` and `CorrectionRegistrar` interfaces. This keeps correction behavior explicit without introducing a generic audit abstraction.

#### Code Review

An experienced Go engineer would approve the append-only direction because it preserves historical facts and makes corrections explicit. The main follow-up for L10.4 is MVP-level review: decide whether the API should expose an effective/current production-entry view or keep corrections as separate history for now.

#### Exercises

- Add a test proving multiple corrections are returned newest first.
- Design a read model that returns original entry plus latest correction as an effective view.
- Explain why correction reason should be required even for simple typo fixes.

#### Interview Questions

- Why are append-only audit records safer than direct updates for historical data?
- What belongs in a correction record versus the original entry record?
- Why should actor identity come from authentication context rather than request body?
- When would a generic audit framework become justified?

#### Roadmap Update

- Lesson 10.3 completed.
- Current lesson moved to Lesson 10.4.
