## MESLite Configuration Synthesis

### 1. Configuration as the top-level concept

MESLite should have a unified **configuration system**.

Feature flags should not be implemented as a separate subsystem. They are one category of configuration.

Conceptual categories:

```text
Configuration
├── static deployment settings
├── dynamic technical settings
├── product/customer configuration
└── feature/release flags
```

Permanent optional product functionality may be thought of as **Capabilities**, while temporary rollout/kill-switch controls remain **Feature Flags**.

---

## 2. Configuration precedence

Configuration values should be resolved using these layers:

```text
in-code default
    ↓ overridden by
config file
    ↓ overridden by
environment variable
    ↓ overridden by
persisted runtime override
```

Conceptually:

```text
effective =
    runtime_override
    ?? environment
    ?? config_file
    ?? code_default
```

Runtime overrides are normally changed through an administrative interface/API, e.g. HTTP.

They allow configuration to change **without application redeployment or restart**.

---

## 3. In-code definitions are the source of configuration metadata

Every configuration key should be explicitly defined in code.

A definition should at minimum specify:

```text
key
type
default value
static/dynamic
validation
```

Likely useful metadata later:

```text
description
scope
runtime-mutable
secret/sensitive
restart-required
allowed range/values
```

Benefits:

* one canonical list of supported configuration;
* typo/unknown-key detection;
* typed access;
* centralized validation;
* possible automatic config documentation;
* possible automatic Admin UI/API generation.

Configuration should be accessed through typed definitions rather than arbitrary string parsing throughout the code.

---

## 4. Static vs dynamic configuration

Not every setting should be changeable while MESLite is running.

Each definition should explicitly classify this.

### Static examples

```text
deployment role
database connection
storage backend
factory identity of an Edge
network listener
fundamental encryption/security setup
```

These are normally fixed at startup.

### Dynamic examples

```text
feature flags
product capabilities
logging level
some retry/synchronization parameters
reporting limits
selected operational thresholds
```

Dynamic runtime modification should be explicitly allowed rather than being the default.

---

## 5. Runtime configuration

The highest-precedence layer should be a **persisted runtime override store**.

HTTP/API is only a mechanism for modifying this store; runtime configuration must not exist only in process memory.

Example:

```text
PUT feature.foo = false
```

should persist the override.

After application restart, the override remains effective.

Removing the override:

```text
DELETE feature.foo override
```

means:

> stop overriding this key and return control to the lower-precedence layers.

Therefore runtime storage should be **sparse**: only explicit overrides are persisted, not copies of every effective value.

---

## 6. Runtime configuration caching

The persisted runtime store is authoritative, but normal reads should use a **process-local in-memory cache**.

Conceptual read path:

```text
configuration request
        ↓
effective-value cache
        ↓ miss
runtime persisted store
        ↓
lower-precedence configuration layers
        ↓
resolved value cached
```

Normal reads should therefore be very cheap.

Caching the **resolved effective value** is preferable to caching only raw runtime overrides.

Example cache key:

```text
(company, factory, configuration-key)
```

### Write behavior

Runtime updates should follow:

```text
validate
    ↓
persist override/change
    ↓
update or invalidate affected cache entries
```

Prefer updating the local cache immediately after successful persistence where practical.

Removing an override must invalidate/update the effective-value cache so the value can fall back to environment/file/default configuration.

Cache failure must not make configuration unavailable; resolution can fall back to persisted configuration.

---

## 7. Multi-process cache invalidation

If several Core/Edge processes can observe the same configuration, changing persisted configuration must eventually invalidate their local caches.

Possible mechanisms include:

```text
database notification
pub/sub
configuration-version polling
```

Do not introduce a sophisticated distributed caching platform before it is actually needed.

Initial direction:

> process-local cache + persisted runtime store + simple invalidation/version mechanism.

---

## 8. Scope is independent from precedence

Configuration has two independent dimensions:

### Source precedence

```text
default < file < environment < runtime override
```

### Business/deployment scope

Initially:

```text
global
    ↓
company
    ↓
factory
```

More specific scope overrides broader scope.

User/workstation scope should not be introduced unless a real requirement appears.

Example:

```text
feature.machine_connectivity

default/global:
    false

Company A:
    true

Company A / Factory Kraków:
    false
```

The resolver therefore evaluates both:

```text
scope resolution
+
source precedence
=
effective configuration value
```

---

## 9. Feature flags and capabilities

Feature flags should use the general configuration infrastructure.

Main use cases:

### Permanent capabilities

Examples:

```text
quality management
machine connectivity
inventory
advanced planning
```

These represent functionality available for selected companies/factories.

### Temporary release flags

Examples:

```text
new_work_entry_ui
new_reporting_engine
```

Lifecycle:

```text
off
→ limited rollout
→ general rollout
→ old implementation removed
→ flag deleted
```

Temporary flags must not remain indefinitely.

### Kill switches

Runtime configuration provides operational kill switches.

Example:

```text
feature.production_import_v2 = false
```

allows malfunctioning functionality to be disabled immediately without rebuilding/redeploying MESLite.

Only functionality that is safe to disable should have such a switch.

Architectural invariants such as command idempotency must **not** become feature flags.

---

## 10. Configuration vs product configuration

Avoid representing every customer difference as a boolean flag.

Use:

```text
optional capability
    → capability/feature configuration

different values/options
    → ordinary product configuration

fundamentally different business semantics
    → proper domain design/module

customer-name condition
    → prohibited
```

Avoid code such as:

```go
if company == "CompanyA" {
    ...
}
```

MESLite should remain one product, with customer differences represented through configuration where practical.

---

## 11. Configuration is not authorization

These questions are separate:

```text
Is the capability enabled?
        ↓
Does this user's role permit using it?
        ↓
Can the current deployment/offline state support it?
```

Feature/configuration flags must never replace RBAC or server-side authorization.

---

## 12. Core / Edge / Terminal behavior

Runtime configuration must respect MESLite's offline architecture.

A factory must not require a live Internet connection merely to determine its effective production configuration. This follows MESLite's existing offline-continuity requirements. 

For hybrid deployments:

```text
Core
    authoritative runtime configuration
        ↓
configuration synchronization
        ↓
Edge
    durable last-known configuration snapshot/cache
        ↓
Terminal
    relevant configuration/capability projection
```

If WAN connectivity disappears, Edge continues using the **last durably received configuration** needed for supported offline factory operation.

A remote configuration change that has not yet reached an offline Edge cannot be assumed active there.

---

## 13. Version/component compatibility

A capability may require support from several MESLite components.

Example:

```text
Core v1.4      supports Feature A
Edge v1.3      does not
Terminal v1.4  supports Feature A
```

Enabling it merely because the configuration says `true` may be unsafe.

For cross-component features, effective activation may eventually require:

```text
configured enabled
AND
required Core version supports it
AND
required Edge version supports it
AND
required Terminal version supports it
```

Recommended rollout:

```text
1. Deploy new code with feature disabled.
2. Upgrade required Core/Edge/Terminal components.
3. Verify deployment.
4. Enable feature for selected company/factory.
```

This supports a common MESLite version across customers without requiring customer-specific forks.

---

## 14. Runtime configuration API semantics

Conceptual operations:

```text
Get effective configuration
Get configuration definition/metadata
Set runtime override
Remove runtime override
```

Potential application API shape:

```go
config.Get(Key)
config.Effective(Key, Scope)
config.SetRuntimeOverride(Key, Scope, Value)
config.RemoveRuntimeOverride(Key, Scope)
```

Business code should use this abstraction rather than directly querying configuration tables/environment variables.

---

## 15. Validation

Runtime values must be validated before becoming effective.

Examples:

```text
sync.batch_size
type: integer
min: 1
max: 10000
```

Unknown configuration keys should be rejected rather than silently accepted.

Validation belongs in the configuration subsystem, consistent with MESLite's general principle of enforcing rules at the layer that owns them. 

---

## 16. Auditability and security

Runtime configuration is effectively an administrative control plane.

Changes should be auditable.

Recommended audit information:

```text
configuration key
scope
previous effective value
new override/value
changed by
changed at
optional reason/comment
```

Runtime configuration changes should normally require Admin-level authorization.

Sensitive configuration/secrets may require stricter handling and should not automatically be exposed through the generic runtime API/UI.

---

## 17. Recommended core architectural model

```text
                 In-code definitions
                 + default values
                         |
        +----------------+----------------+
        |                |                |
   config file      environment      runtime overrides
                                         |
                                   persisted store
                                         |
                                  in-memory cache
                                         |
                                         v
                                  Config Resolver
                                         |
                                  effective value
                                         |
                     +-------------------+-------------------+
                     |                   |                   |
                   Core                Edge              Terminal
                                      durable            relevant
                                      snapshot           projection
```

### Main design principle

> Configuration is centrally defined, layered, typed, validated and auditable. Runtime overrides provide safe operational control, while caching keeps effective-value reads cheap. Dynamic configuration must remain compatible with MESLite's offline and distributed Core/Edge/Terminal architecture.
