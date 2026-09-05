# MESLite Domain

## Purpose

Define the current MESLite business domain: core entities, relationships, terminology, lifecycles, workflows, and business rules.

This document is the domain-level source of truth for business logic, persistence models, workflows, and APIs that expose domain concepts.

## Current domain scope

MESLite currently focuses on lightweight production execution rather than full ERP or enterprise MES modeling.

The initial domain should support:

- products or outputs to be produced;
- production orders where the customer process uses them;
- order lines with planned quantities;
- operations where production is naturally represented as manufacturing steps;
- work entries recording what actually happened on the shop floor;
- operators, workstations, quantities, rejects, scrap, downtime, statuses, and timestamps needed for production visibility and auditability.

Advanced structures such as full BOM management, sophisticated routings, scheduling, inventory reservation, quality management, subcontracting workflows, and operation dependency graphs are not current V1 domain requirements unless promoted by product requirements.

## Production work model

The initial conceptual model is:

```text
ProductionOrder
  -> OrderLine
      -> product/output
      -> planned quantity

Operation
  -> optional manufacturing step for production work

WorkEntry
  -> actual recorded execution
  -> operator
  -> workstation
  -> quantity
  -> time
  -> status
  -> rejects/scrap/downtime
```

`ProductionOrder` is a planning concept. It describes what should be made where the customer workflow uses orders.

`OrderLine` describes a product, material, output, or production target and its planned quantity within an order.

`Operation` describes a manufacturing step or activity when the production process needs step-level tracking. Operations should not be modeled as mandatory for every manufacturing mode.

`WorkEntry` records actual execution. Work entries are the primary historical evidence of who did work, where it happened, when it happened, what quantity was produced, and what production outcomes such as rejects, scrap, or downtime were observed.

The domain must not require all manufacturing work to fit a single rigid hierarchy such as production order -> product -> operation. Production work may later be grouped by work package, batch, production run, delivery, installation site, project phase, customer section, team, workstation, or another validated operational reason.

Future concepts such as work packages, batches, production runs, routing templates, operation dependencies, BOM/material structure, quality inspections, and subcontracted operations may be added only when validated customer workflows justify the added complexity.
