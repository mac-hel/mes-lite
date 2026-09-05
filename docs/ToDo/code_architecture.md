## <PARTIALY> Encapsulation - hide private properties of entities (structs)
Currently all domain Entities have public fields. Would it be feasible to make fields private at current phase of project?
- <DONE> Order

## <PARTIALY> Fields as Collections
What do you think about using collections as types of fields that contain multiple values, e.g. fields like:
- <DONE> `Order.lines` (currently `[]OrderLine`, could be collection)
- <REJECTED FOR NOW> `Order.assignedEmployees` (currently `[]string`, could be collection)

## <REJECTED> Double id - UUID for user-facing, int for internal reference
What do you think about using UUID and int ID simultaneously? UUID for user-facing things, int for internal references (for speed).
ANSWER: Do not introduce int IDs just for speed right now. Keep one canonical identifier per entity. If performance becomes real later, benchmark and revisit with data.

## validation library
Is there a place for validation library in the project, in any package, in any layer?

## remove protections for existing data (in reality no data exist yet)
e.g. remove `request_id <> ''` from `production_entries.request_id` index
remove `NOT VALID` from migrations/0004_add_production_reference_foreign_keys.sql
check other migrations

## <REJECTED> Each package builds all sqlc queries
Currently each package (slice) builds all sqlc queries across all packages.
Isn't this approach:
- confusing for developers
- crossing package boundaries
Shouldn't each package build only own `sqlc` queries and create only own query models?
see `docs/sqlc-boundaries.md`
