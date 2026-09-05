# MESLite Product Context

## Purpose

Define why MESLite exists, who it is for, the manufacturing environment it serves, its product goals and current scope boundaries, and the terminology used when discussing the product.

This document is the product-level source of truth for evaluating scope, UX, product behavior, terminology, and whether a proposed feature belongs in MESLite.

## Product vision

MESLite is a lightweight Manufacturing Execution System (MES) for small and medium-sized manufacturers that have outgrown Excel, paper, or manual production reporting but do not need or want a large, complex MES/ERP implementation.

Its core purpose is to provide simple, reliable visibility and control over production execution:

- what is being produced;
- who is working on it;
- how much has been produced;
- when work happened and how long it took;
- where delays, downtime, rejects, or scrap occur;
- what the current production status is.

MESLite should deliver the useful core of production execution and monitoring with substantially less complexity and organizational burden than enterprise-class MES/ERP systems.

The product is independent and must be useful as a stand-alone system even when a customer has no ERP.

## Target customers

Primary target customers are manufacturing companies with approximately **10–250 employees** that:

- operate their own physical production;
- need better production execution visibility and reporting;
- still rely significantly on Excel, paper, manual supervisor reporting, or incomplete ERP production functionality;
- do not have an advanced MES, or consider existing MES solutions too large, expensive, or complicated;
- want a focused production system rather than a broad enterprise platform.

The initial go-to-market geography is Poland, but MESLite is not intended to be Poland-specific.

The ventilation manufacturer used during early discovery is a pilot/customer-discovery starting point, not the definition of the target industry.

MESLite should support both discrete/batch-style manufacturing and continuous/process manufacturing. Product and domain design must therefore avoid assuming that all customers use the same order/operation workflow. Exact workflows for different manufacturing modes belong in `requirements.md` and `domain.md`.

## Product positioning

Customer-facing positioning should emphasize the operational problem rather than the MES category itself.

Core positioning:

> Simple production management and monitoring for manufacturers that have outgrown Excel and paper.

MESLite competes not only with other MES products, but also with:

- Excel;
- paper records;
- manual supervisor reporting;
- partially used ERP production modules;
- doing nothing.

The product should differentiate through simplicity, focused production functionality, fast adoption, and low organizational burden rather than by offering the largest feature set.

## Product goals

1. **Make production execution visible.** Management and production staff should be able to understand current and historical production without reconstructing it manually from spreadsheets or paper.
2. **Make shop-floor recording simple.** Operators should be able to record production work with minimal interaction and minimal typing.
3. **Provide reliable production records.** Recorded production information should be durable, auditable, and recoverable wherever technically possible.
4. **Keep production running through infrastructure failures.** Loss of Internet, local connectivity, or a Terminal should not unnecessarily stop physical production; manual recovery remains the final fallback.
5. **Remain useful without ERP.** MESLite must work as an independent stand-alone production system.
6. **Support multiple deployment needs.** The product must support cloud-only, on-premise-only, and hybrid deployments from one logical product architecture.
7. **Support multi-factory organizations from the beginning.** Multi-factory capability is part of the product scope even if early customers use only one factory.
8. **Stay lightweight.** Add complexity only when a validated production problem justifies it.
9. **Productize recurring manufacturing needs.** Customer-specific implementations should increasingly become configuration around a shared MESLite core rather than separate custom systems.

## Current functional focus

The initial product should focus on the smallest useful production-execution loop.

Expected initial capabilities include:

- products;
- production orders;
- production operations where the manufacturing process uses them;
- operators and role-based access;
- workstations;
- work entries recording production activity;
- start/stop or equivalent time recording where applicable;
- produced quantities;
- rejects;
- scrap;
- downtime/problems and reasons;
- production status and history;
- simple management dashboards and reports;
- retrospective/correction entry with auditability;
- import of existing production data where needed for migration, especially from Excel.

This list describes the current product direction, not detailed requirements or domain rules. Exact behavior, fields, workflows, and acceptance criteria belong in `requirements.md` and `domain.md`.

## Deferred product areas

MESLite is intentionally not trying to become a full ERP or enterprise platform now.

Complex business areas such as the following are **not currently planned**:

- full ERP functionality;
- accounting;
- CRM;
- purchasing;
- sophisticated warehouse management;
- other broad enterprise-management domains unrelated to focused production execution.

The following are plausible later extensions but are **not V1 requirements** unless explicitly promoted in `requirements.md`:

- inventory/material management;
- broader quality management;
- production planning and scheduling;
- predictive maintenance;
- machine connectivity and telemetry;
- simple OEE;
- AI-assisted production analysis or decision support;
- other lightweight production-related capabilities validated by real customer needs.

Nothing is permanently excluded. Future scope should be driven by validated customer problems and product coherence rather than by a generic MES feature checklist.

## Relationship with other business systems

MESLite is a stand-alone independent system and must not require an ERP.

For now, customer/sales orders are expected to live outside MESLite. MESLite primarily concerns production execution and may receive or create the production work required by the factory.

A simple customer-order capability may be introduced later if it materially helps customers that have no suitable upstream system. MESLite should not expand into a full sales/ERP system merely to support that use case.

External integrations may be added where justified, for example with ERP, warehouse software, scanners, printers, machines, or other production systems. Specific integrations are not part of the product context unless separately established as requirements.

## Actors

### Operator

Shop-floor user who performs and records production work.

Primary needs:

- see or identify the relevant work;
- start/continue/finish work where the process uses those states;
- record quantities and production outcomes;
- record rejects, scrap, downtime, or problems where applicable;
- interact with MESLite quickly and ergonomically in factory conditions.

### Leader

Front-line production supervisor/team leader.

Primary needs:

- understand current production state;
- monitor operators/workstations and pending work;
- respond to delays, downtime, recording problems, or operational exceptions;
- perform or assist with authorized corrections/reconciliation where permitted.

### Manager

Production/operations management user.

Primary needs:

- understand production progress and delays;
- review production history and performance;
- use dashboards and reports for operational decisions;
- manage production-related configuration and workflows within assigned permissions.

### Admin

Administrative/system-management user.

Primary needs:

- manage users, permissions, configuration, factories/workstations, and other administrative concerns;
- support recovery, audit, and operational administration where authorized.

Owners, CEOs, directors, and similar stakeholders use MESLite through the appropriate Manager or Admin role rather than introducing separate product roles merely because of job title.

## Business environment

MESLite runs in factory environments where:

- production may continue even when IT infrastructure is degraded;
- operators may use shared or dedicated terminals, tablets, or desktops;
- interaction should tolerate factory-floor constraints such as limited time for data entry;
- work may be organized by orders and operations, but the product must also support manufacturing models that do not fit that exact pattern;
- production information may currently originate in Excel, paper, ERP, supervisor notes, or verbal/manual processes;
- factories may be single-site or part of a multi-factory organization;
- Internet and local connectivity cannot be assumed to be continuously reliable.

A core product principle is that IT failure should not unnecessarily become a production stoppage. Detailed offline, synchronization, and recovery guarantees are defined in `../03-architecture/offline-sync.md`.

## Product principles

### Simplicity over feature count

MESLite should solve a narrow production problem clearly before adding adjacent capabilities. A simple workflow that operators actually use is more valuable than a comprehensive but burdensome MES feature set.

### Shop-floor ergonomics first

Operator interactions should minimize typing, navigation, and cognitive load. UI design should reflect real factory conditions rather than ordinary office-software assumptions.

### Reliable records over fragile convenience

Production records and corrections must remain auditable. Failures, retries, offline states, and manual recovery must not silently corrupt or erase production history.

### Configuration over customer-specific forks

Repeated customer differences should become configuration where practical. The long-term product is one shared MESLite system, not a collection of unrelated customer-specific codebases.

### Validate before expanding

New product areas should be added because repeated customer evidence justifies them, not because they appear on generic MES/Industry 4.0 checklists.

### AI follows reliable production data

AI is not part of the initial product proposition. Later AI or optimization should be used only where MESLite has sufficient trusted production data and where it improves real production decisions, for example delay risk, bottlenecks, prioritization, or abnormal-performance detection.

## Important terminology

### MES

**Manufacturing Execution System.** Software that manages and records what happens during production execution on the factory floor.

### Production order

A production instruction to manufacture a defined product/output, quantity, or production target. Exact form and lifecycle may vary by manufacturing model.

### Operation

A manufacturing step or activity within a production flow. Operations are relevant where the customer's process is naturally represented as multiple steps; not every manufacturing mode must use the same operation model.

### Operator

The canonical MESLite term for the shop-floor person performing production work.

### Leader

Front-line production supervisor/team leader.

### Manager

Production or operations management user.

### Admin

Administrative/system-management user.

### Work entry

The canonical MESLite term for a recorded unit of production work/activity. Exact fields and lifecycle belong in `domain.md`.

### Workstation

A logical or physical production location at which work is performed and to which a Terminal may be associated. A workstation is distinct from a machine.

### Machine

Production equipment used to perform manufacturing work. A machine is distinct from a workstation; their detailed relationship belongs in `domain.md`.

### Reject

A separately recorded rejected production outcome/quantity. Reject is distinct from scrap. Detailed classification rules belong in `domain.md`.

### Scrap

A separately recorded scrap outcome/quantity. Scrap is distinct from reject. Detailed classification rules belong in `domain.md`.

### Downtime

A period or event in which expected production work cannot proceed, optionally associated with a reason such as equipment failure, missing material, waiting, setup, or another production problem.

### Traceability

The ability to reconstruct relevant production history for an item, quantity, order, batch, material, operator, operation, workstation, machine, or time period where required by the product/domain.

### OEE

**Overall Equipment Effectiveness.** A manufacturing metric combining availability, performance, and quality. Only simple OEE is a possible later MESLite capability; it is not a V1 requirement.

## Current product strategy

The current target-market and positioning statements in this document are settled product direction, not merely uncommitted research hypotheses.

The implementation strategy remains evidence-driven:

1. validate workflows in real factories;
2. solve the recurring production-execution problem well;
3. identify what is common across customers;
4. convert repeated differences into configuration;
5. broaden industries and capabilities without losing the lightweight product character.

The initial ventilation customer is therefore a pilot and discovery environment, while MESLite itself is intended as a broader manufacturing product.
