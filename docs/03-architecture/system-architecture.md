# System Architecture

## Purpose

Define MESLite's high-level components, deployment models,
responsibilities, communication paths, and failure assumptions.

## Architectural principles

1.  **Production must not depend on Internet availability.**
2.  **Loss of Wi-Fi or a worker terminal must not require production to
    stop.**
3.  MESLite uses **one logical system and business model** across
    deployment variants; deployment topology changes, not product
    semantics.
4.  Prefer a **single authoritative MES Core** over bidirectional
    database ownership.
5.  Edge and Terminal provide progressively local durability and
    continuity; they do not become independent competing systems of
    record.
6.  When digital capture is impossible, production continues using an
    explicit manual/paper recovery process.
7.  Components and deployment modes must remain as simple as reliability
    requirements allow.

## Components

### MES Core

The authoritative MES service.

Responsibilities: - authoritative business logic and validation; -
authoritative persistent production and configuration data; - APIs used
by clients or Edge; - production, order, administration, reporting, and
management functions appropriate to the deployment; - audit history; -
processing commands idempotently.

MES Core may run in the cloud or on-premise.

### Factory Edge

Optional factory-local service used when MES Core is remote.

Purpose: remove the WAN/Internet connection from the production-critical
path.

Responsibilities: - accept production-critical commands while the cloud
is unavailable; - durably journal commands before acknowledging them; -
cache the working data required for permitted offline operations; -
forward/retry queued commands to Core; - receive required
working/configuration data from Core; - expose
synchronization/connectivity state.

Edge is **not a second authoritative MES database** and must not
introduce arbitrary multi-master synchronization.

Offline scope should be deliberately bounded. At minimum, an Edge
deployment must support continuation of already prepared/released
production work and recording of production-critical actions.
Administrative, cross-site, and other non-critical functions may be
unavailable while Core is unreachable.

### Terminal

Worker/manager client application, normally a web application.

Responsibilities: - ergonomic production and management UI; -
identify/bind the active worker and workstation; - persist unsent worker
commands in durable device storage; - retry delivery automatically; -
show sufficient connectivity/pending-data status; - allow a replacement
terminal to be configured quickly without migrating normal application
state.

A terminal is a disposable client, **not the authoritative store of
workstation or worker state**.

## Deployment models

MESLite must support all three models from the same product
architecture.

### Cloud-only

``` text
Terminal -> Internet -> MES Core (cloud)
```

-   Core is authoritative in the cloud.
-   No Edge.
-   Terminal durable storage protects short local/WAN interruptions.
-   Suitable only where the customer's connectivity/reliability
    requirements permit it.

### On-premise-only

``` text
Terminal -> factory network -> MES Core (factory)
```

-   Core is authoritative on-premise.
-   No cloud dependency is required for production.
-   Edge is unnecessary because Core is already local.
-   Remote access/integrations depend on customer infrastructure and
    requirements.

### Hybrid

``` text
Terminal -> Factory Edge -> Internet -> MES Core (cloud)
```

-   Cloud Core is authoritative.
-   Edge provides durable local continuity during WAN outages.
-   Terminal storage protects against loss of connectivity to Edge.
-   Cloud provides centralized access and naturally supports
    multi-factory management.

Hybrid is the preferred topology when both centralized access and strong
factory continuity are required.

## Logical communication model

Normal hybrid flow:

``` text
Terminal
   |
   | production command
   v
Edge durable journal
   |
   | asynchronous/retryable delivery
   v
MES Core
   |
   v
Authoritative database
```

Every production-changing request that may be retried must carry a
globally unique command identifier. Core must safely accept duplicate
delivery without duplicating the business effect.

Do not synchronize arbitrary database rows between Core and Edge.
Synchronize explicit commands/events and the bounded
working/configuration data required by Edge.

## Failure model

### Internet/WAN failure

In hybrid deployment: - factory production-critical functions continue
through Edge; - Edge queues unsynchronized commands durably; -
synchronization resumes automatically when connectivity returns; -
cloud-only functions may be unavailable/stale.

### Wi-Fi/local network failure

-   Terminal stores newly recorded commands in durable local storage
    before considering them recorded locally.
-   Terminal retries automatically when Edge/Core becomes reachable.
-   Worker can continue within the explicitly supported offline-terminal
    scope.

### Terminal failure

-   A replacement terminal should be usable immediately.
-   Workstation identity should be recoverable by configuration/scan
    (for example a workstation QR), not by copying device state.
-   Worker identity should be re-established through normal
    authentication/badge/PIN mechanisms.
-   Data already delivered to Edge/Core is unaffected.
-   Unsynchronized data existing only on a destroyed terminal cannot be
    technically guaranteed; use the recovery workflow defined in
    `offline-sync.md`.

### Complete MES unavailability

Production must have an operational fallback independent of MES, such as
paper production records. When service is restored, authorized staff
retrospectively enter/reconcile the records with full auditability.

## Availability and durability boundary

MESLite must not claim physically impossible guarantees. If information
exists only on one terminal and that terminal's persistent storage is
destroyed before transmission, the information cannot be recovered
automatically.

Therefore the reliability model is layered:

``` text
Core authoritative storage
        ^
Edge durable journal
        ^
Terminal durable queue
        ^
Manual/paper operational record
```

Each lower layer protects production continuity when the layer above is
unavailable.

## Data authority

There is one authoritative MES Core for a deployment.

-   Core owns authoritative MES state.
-   Edge owns only its durable synchronization journal/cache state.
-   Terminal owns only transient client/session state and its durable
    pending-command queue.
-   Manual/paper records are temporary recovery evidence to be
    reconciled into Core.

Avoid designs where Cloud and Edge freely modify overlapping
authoritative database state.

## Multi-factory behavior

When Core is cloud-hosted, multiple factories connect to the same
authoritative system through their own Edge instances when required:

``` text
Factory A Edge --\
Factory B Edge ----> Cloud MES Core
Factory C Edge --/
```

This provides centralized management/reporting without making factory
production dependent on continuous Internet connectivity.

## Cross-document rules

-   Detailed durability, synchronization, retries, recovery, and
    reconciliation rules: `offline-sync.md`.
-   Repository/module structure and shared-code boundaries:
    `software-architecture.md`.
-   Functional/non-functional guarantees should be mirrored as
    requirement IDs in `../01-product/requirements.md` when that
    document is created.
-   Significant future changes to these architectural principles should
    be documented by ADR when rationale/history must be preserved.
