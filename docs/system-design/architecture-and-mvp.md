# MES Lite Architecture and MVP Design

## Overview

MES Lite should initially be built as a small number of deployable systems with clear internal boundaries.

The recommended initial shape is:

- Go backend / application server
- PostgreSQL database
- Server-rendered web UI served by the backend
- Reserved integration boundary for machines, ERP, printers, scanners, and other factory systems

The system should avoid premature separation into multiple frontend or backend applications. The production domain is still more important than deployment topology at this stage.

## System Shape

For the MVP, one application plus one database is enough:

- Go application server
- PostgreSQL database

The application server should own:

- Business logic
- HTTP routes
- Worker UI
- Manager/Admin UI
- Authentication and authorization
- Production registration
- Reporting
- Background jobs, if needed
- Integration boundaries

PostgreSQL should be the primary source of truth for production, order, user, workstation, and audit data.

```text
Worker terminal ─┐
Tablet          ─┤
Desktop         ─┼──► Go Application Server ───► PostgreSQL
Manager PC      ─┤             │
Admin PC        ─┘             ▼
                         Integration Boundary
                      ERP / PLC / machines / printers
```

## Unified Web Application

Worker and manager interfaces should be part of one unified web application. They should not be separate applications initially.

The important distinction is UX mode, not deployment unit.

Worker mode should be:

- Simple
- Fast
- Touch-friendly
- Barcode/scanner-friendly
- Low navigation
- Focused on the current workstation and task
- Optimized for repetitive production input

Manager/Admin mode should be:

- Denser
- Table-oriented
- Filterable
- Report-oriented
- Permission-controlled
- Suitable for desktop use

Example route structure:

```text
/workstation
/workstation/jobs
/workstation/operations

/manager/orders
/manager/production
/manager/reports

/admin/users
/admin/workstations
/admin/products
```

These routes can render different templates/layouts while sharing the same backend, session handling, domain model, and database.

## Context and Authorization

Permissions alone are not enough for an MES.

Access and behavior should be determined by at least three dimensions:

- User role: worker, leader, manager, admin
- Workstation or device context: station, machine, production cell, terminal
- Current production context: active order, operation, shift, batch, operator

For example, a worker logged into a terminal assigned to `Laser Cutting 3` should only see actions relevant to that station. The same user may see different options from a generic desktop or from another workstation.

Authorization should be centralized behind policy-style checks rather than scattered role comparisons.

Prefer checks such as:

```go
CanViewReports(user)
CanRegisterProduction(user, workstation)
CanCorrectProductionEntry(user)
```

Avoid spreading checks like this throughout handlers:

```go
if user.Role == "admin" || user.Role == "manager" {
    // ...
}
```

## Integration Boundary

An integration or edge layer will likely become important later, but it does not need to be a separate service in the MVP.

The system should reserve a boundary for:

- Barcode scanners
- Label printers
- Machine integrations
- PLC communication
- ERP import/export
- MQTT / OPC UA
- Local device communication

For MVP, this can be represented as internal modules/packages rather than a deployed service:

```text
internal/integrations/
internal/barcode/
internal/printing/
internal/imports/
```

The main goal is to avoid mixing external device and industrial connectivity details directly into core production logic.

## MVP Scope

The MVP should be scoped as an Excel/paper replacement.

Its goal is not to implement a full MES immediately. Its goal is to replace the simplest current manual production records with a durable, searchable, auditable system.

The first version should focus on:

- Basic users and roles
- Basic workstation definitions
- Basic product/order/operation records needed for registration
- Worker production registration
- Simple manager access to entered records
- Simple reports that replace existing Excel or paper summaries

Anything beyond replacing current manual tracking should be treated as MVP-V2 or later.

## MVP Vertical Slice

The first usable version should be a narrow vertical slice:

- Define one or more workstations
- Define the minimum product/order/operation data needed for registration
- Start an operation
- Register good quantity
- Register scrap or rework if currently tracked on paper/Excel
- Register downtime or problem if currently tracked on paper/Excel
- Finish an operation
- View entered production records
- Export or display the same summaries currently maintained manually

This keeps the first release focused on replacing real manual work rather than building a broad administration system.

## MVP-V2 Modules

After the paper/Excel replacement is working, the next version can expand toward a fuller MES.

Likely MVP-V2 modules:

1. Users and roles

   Broader permission model, correction permissions, and richer audit views.

2. Workstations

   Define stations, machines, cells, and allowed operations in more detail.

3. Products and operations

   Product definitions, routing, and operation steps beyond the minimum MVP fields.

4. Production orders

   Create and manage orders, quantities, deadlines, statuses, and planning data.

5. Worker production registration

   More complete workflow for pauses, operator changes, richer downtime classification, quality checks, and exception handling.

6. Manager views

   Operational views for active orders, workstation status, production progress, problems, downtime, and broader reports.

## MVP Reporting

MVP reporting should only replace reports currently maintained in Excel or on paper.

Examples may include:

- Daily production summary
- Produced quantities by order or operation
- Scrap/rework summary if currently tracked manually
- Downtime/problem summary if currently tracked manually
- Exportable records for manual review

Operational dashboards such as active order progress, workstation status, downtime monitoring, and broader management analytics should be treated as MVP-V2 or future expansion unless they directly replace an existing manual report.

## Worker Production Workflow

The most important MVP workflow is worker-side production registration.

The first version should support:

- Select or scan job
- Start operation
- Register good quantity
- Register scrap or rework
- Register downtime or problem
- Pause operation if currently tracked on paper/Excel
- Finish operation

The worker UI should be constrained and ergonomic. It should not feel like a reduced manager interface.

Useful worker UI properties:

- Large controls
- Clear current job and operation
- Minimal typing
- Barcode-first where possible
- Few screens
- Explicit error recovery
- Confirmation only for rare or destructive actions
- Good tablet and kiosk behavior

## Event-Oriented Production Data

MES data is naturally event-like. Important production actions should be recorded as durable append-only event records, not only reflected as final state.

Production registration should use event records as the primary write model. Current-state tables or projections may be maintained for query convenience, but they should be derived from or consistent with the durable production events.

Examples of useful production events:

- Operation started
- Quantity produced
- Scrap registered
- Rework registered
- Downtime started
- Downtime ended
- Problem reported
- Operator changed
- Order paused
- Operation completed
- Correction made

This improves:

- Traceability
- Auditability
- Reporting
- Historical reconstruction
- Correction handling
- Debugging of production discrepancies

The system may still maintain current-state tables for fast queries, but important actions should not disappear into overwritten fields.

## Corrections and Auditability

Production history should not be silently edited.

Prefer a correction model:

- Keep the original entry
- Add a correction entry
- Record who made the correction
- Record when it happened
- Record the reason

This is especially important for production quantities, scrap, downtime, and completed operations.

Normal worker flows should not allow silent correction of historical production data. Correction rights should be explicit and likely limited to leader, manager, or admin roles.

## Reliability Concerns

Even if the first version is not offline-first, the system should account for factory reliability issues.

Important concerns:

- Prevent duplicate submissions
- Make production registration idempotent where practical
- Include an idempotency key or request ID with production registration commands
- Provide clear error recovery
- Avoid losing partially submitted production data
- Handle workstation/device session timeout deliberately
- Keep durable server-side records for production actions

Offline-first behavior can be deferred unless it becomes a real operational requirement, but the design should not make it impossible later.

## Main Architectural Risk

The biggest risk is treating MES Lite as a generic CRUD admin application.

The core domain is production execution and traceability. The model should be centered around:

- Orders
- Operations
- Workstations
- Workstation state
- Production events
- Operators
- Downtime
- Scrap and rework
- Problems
- Corrections

Administrative screens are necessary, but they should support production execution rather than define the architecture.

## Recommended Next Steps

Next design work should define:

- Core database entities
- Production event model
- Initial route structure
- Worker workflows
- Manager workflows
- Permission and workstation-context model
- MVP reporting requirements

These should become the implementation roadmap for the first production-oriented version of MES Lite.
