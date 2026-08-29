# Offline, Synchronization, and Recovery

## Purpose

Define how MESLite preserves production continuity and recorded worker
actions across Internet, Wi-Fi, Edge, and Terminal failures.

This is a production-critical specification. Changes to recording, local
persistence, retries, synchronization, or recovery must preserve these
rules.

## Reliability objective

Operational objective:

> Infrastructure failure should not require production to stop. MES
> records are automatically preserved wherever technically possible,
> with an audited manual recovery process for failures where digital
> capture cannot be guaranteed.

Required behavior: - Internet loss must not stop production in
deployments with a local Core or Edge. - Wi-Fi/local connectivity loss
must not stop supported Terminal recording. - Terminal failure must
allow rapid replacement. - Previously durably transmitted/acknowledged
production data must not be lost because a later network or Terminal
failure occurs. - Catastrophic cases where data existed only on
destroyed storage are handled by explicit reconciliation/manual
recovery. - Complete MES unavailability must have a paper/manual
operational fallback.

## Durability ladder

``` text
Core authoritative storage
        ^
Edge durable journal
        ^
Terminal durable queue
        ^
Manual/paper record
```

Each layer protects continuity when the layer above cannot be reached.

The lower layers are not competing authoritative MES databases.

## Command identity and idempotency

Every production-changing command that can be queued or retried must
have a globally unique immutable `command_id`.

The same command may be transmitted multiple times because
acknowledgements can be lost.

Core processing must therefore be idempotent:

``` text
if command_id already processed:
    return previous/successful outcome
else:
    validate
    apply business effect
    record command_id/result
    commit atomically
```

A retry must never double production, scrap, downtime, or another
business effect.

Do not rely on "send exactly once" network behavior.

## Terminal recording

### Durable-first rule

When a worker records an action while the destination may be
unreachable, Terminal must persist the command in durable device storage
before treating it as locally recorded.

Durable storage must survive: - page refresh; - browser/application
restart; - device reboot; - temporary power loss, subject to the
device/storage guarantees.

An in-memory queue is insufficient.

Conceptual lifecycle:

``` text
worker action
    |
    v
persist command locally
    |
    v
show action as recorded locally
    |
    v
send/retry to Edge or Core
    |
    v
receive durable acknowledgement
    |
    v
mark/remove/archive local pending command
```

Do not discard a pending command merely because it was sent.
Discard/complete it only after the next durable layer has acknowledged
responsibility.

### Terminal status

The UI should keep synchronization understandable without requiring
workers to understand distributed systems.

At minimum, Terminal must be able to distinguish: -
connected/synchronized; - offline/degraded; - commands pending
delivery; - delivery failure requiring attention.

A supervisor/diagnostic view should expose pending counts and affected
workstation/Terminal where practical.

## Edge behavior

Edge exists in hybrid deployment to remove Internet/WAN connectivity
from the production-critical path.

### Acceptance rule

Edge must durably persist an accepted production command before
acknowledging it to Terminal.

Conceptually:

``` text
receive command
    |
validate enough to accept offline operation
    |
durably append to Edge journal
    |
ACK Terminal
    |
send/retry to Core asynchronously
```

After Edge acknowledges a command, destruction/failure of the
originating Terminal must not lose that command.

### Working cache

Edge may cache only the business/configuration data required to perform
the explicitly supported offline factory operations.

The offline working set may include concepts such as: - released/active
orders; - operations/routings needed for those orders; - products; -
workstations; - workers/permissions needed locally; - recent/current
production state.

The exact cache schema is not yet defined.

Edge must not become a freely editable second authoritative database.

### Offline capability

During Core/WAN outage, prioritize production continuity.

Expected offline capabilities include the production-critical actions
required to continue already prepared/released work, such as: -
starting/continuing operations; - registering produced quantity; -
registering scrap; - recording downtime/problems; - other worker
production events subsequently defined as critical.

Non-critical functions may be unavailable, including major
administration, cross-factory reporting, or changes whose safe semantics
require Core.

Exact supported operations must be specified as requirements as the
domain is defined.

## Synchronization

Synchronization is application-level, not arbitrary database
replication.

Preferred model:

``` text
Terminal command -> Edge journal -> Core processing
```

with explicit working/configuration data flowing toward Edge as
required.

Requirements: - retry indefinitely or until an explicit
recoverable/non-recoverable outcome is reached; - preserve command
identity across every hop; - tolerate duplicate delivery; - preserve
enough ordering/context for domain rules that require it; - expose
unresolved synchronization failures; - never silently drop a production
command.

Do not implement generic bidirectional row/table synchronization between
Cloud and Edge.

## Concurrent/offline changes

Offline systems cannot guarantee that an Edge knows about a cloud-side
change made after connectivity was lost.

Business semantics must therefore define which changes remain valid
offline.

Default architectural rule:

> Work already released/prepared for a factory may continue under the
> last durably received factory working state until synchronization
> resumes, unless a specific domain rule defines stricter behavior.

Future domain requirements must explicitly define behavior for conflicts
such as order cancellation, routing changes, permission changes, or
other safety/quality-critical updates. Do not invent conflict resolution
ad hoc in synchronization code.

## Failure scenarios

### Internet down; Edge available

``` text
Terminal -> Edge journal    X    Cloud Core
```

-   production-critical recording continues;
-   Edge accumulates pending commands;
-   cloud data may become stale;
-   Edge retries automatically;
-   commands are processed idempotently after reconnection.

### Wi-Fi/local network down; Terminal alive

``` text
Terminal durable queue    X    Edge/Core
```

-   Terminal records supported actions durably;
-   worker continues within Terminal's offline scope;
-   commands retry automatically after connectivity returns.

### Terminal powers off/restarts

Pending commands must remain in durable local storage and resume
delivery after restart.

### Terminal fails permanently after Edge/Core acknowledgement

No production data is lost; the next durable layer already owns the
command.

### Terminal fails permanently before transmission

If commands existed only on the destroyed Terminal, automatic recovery
cannot be guaranteed.

Recovery: 1. replace Terminal; 2. re-bind workstation; 3.
re-authenticate worker; 4. determine missing records from operational
evidence; 5. authorized leader/manager enters retrospective
corrections/records; 6. retain full audit trail.

Do not build complex Terminal-to-Terminal state migration as the primary
solution to this case.

## Terminal replacement

Replacement must be operationally fast and must not depend on copying
normal device state.

Preferred model: - take a spare device; - open MES Terminal; - bind it
to the workstation using a simple factory mechanism (for example
workstation QR/configuration); - authenticate the worker using the
normal identity mechanism; - continue.

Workstation identity and worker identity must therefore be recoverable
from the system/physical workstation, not stored only on a specific
Terminal.

Pending commands that never left a destroyed Terminal are handled
through reconciliation, not assumed transferable.

## Manual correction and retrospective entry

Manual recovery is a first-class audited MES capability.

Authorized users must be able to record/correct production
retrospectively for cases such as: - destroyed Terminal with
unsynchronized data; - complete MES outage; - missed worker
registration; - incorrect quantity/scrap/downtime entry; - paper
fallback reconciliation.

A retrospective/correction record should preserve, where applicable: -
business event being recorded/corrected; - effective production
time/date; - affected order/operation/workstation; - quantity/value; -
original worker if known; - person making the correction; -
reason/category; - correction entry time; - linkage to the original
record when correcting an existing entry.

Do not silently overwrite historical facts. Corrections should remain
auditable, preferably as compensating/correction records rather than
destructive edits.

## Complete-system fallback

If no usable digital path exists, for example:

``` text
Internet unavailable
+ local network unavailable
+ Edge/Core inaccessible
+ Terminals unusable/unreliable
```

production must be able to continue using a defined paper/manual
production record.

After restoration: 1. authorized staff enter the records
retrospectively; 2. reconcile totals/events against available digital
records; 3. avoid duplicate entry; 4. document the recovery reason; 5.
retain the required audit evidence according to future operational
requirements.

Paper is the final operational fallback, not part of normal
synchronization.

## Acknowledgement semantics

"Acknowledged" must mean that the receiving layer has durably accepted
responsibility.

Examples: - Terminal may show an action as locally recorded after
durable local persistence. - Terminal may remove a pending item only
after Edge/Core durable acknowledgement. - Edge may remove/complete a
journal item only after Core has durably/idempotently processed or
explicitly resolved it.

Never equate HTTP/network delivery alone with durable business
acceptance.

## Observability requirements

The system must make degraded operation visible.

At minimum support diagnostics for: - Terminal pending-command
count/state; - Edge pending-command count/state; - last successful Core
synchronization; - connectivity state; - commands that cannot be
processed automatically; - reconciliation/manual recovery activity.

Degraded connectivity should not unnecessarily block workers, but
unresolved pending data must not be hidden from supervisors/support.

## Non-goals

This architecture does not attempt to: - guarantee recovery of data
physically destroyed before it reached any second durable location; -
provide arbitrary multi-master Cloud/Edge database editing; - make every
management/administrative feature work offline; - hide catastrophic
recovery by silently modifying history.

## Open decisions

The following are not yet fixed: - Terminal durable-storage
technology; - Edge persistence technology; - exact synchronization wire
protocol; - exact cached working-set schema; - exact offline command
set; - ordering rules required by the future domain model; - retention
period for acknowledged local/Edge journal entries; - detailed
reconciliation UI and approval workflow.

AI agents must treat these as open decisions, not infer earlier examples
as requirements.
