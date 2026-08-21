## <IGNORE FOR NOW> Employee state/status
Employee.isActive -> multiple statuses (sick, vacation, out-of-work etc.)

## <POSTPONED> User-configurable product categories
Distinct entity stored in DB

## <PARTIALY> Encapsulation - hide private properties of entities (structs)
Currently all domain Entities have public fields. Would it be feasible to make fields private at current phase of project?
- <DONE> Order

## <PARTIALY> Fields as Collections
What do you think about using collections as types of fields that contain multiple values, e.g. fields like:
- <DONE> `Order.lines` (currently `[]OrderLine`, could be collection)
- <REJECTED FOR NOW> `Order.assignedEmployees` (currently `[]string`, could be collection)

## <ADDED> (as separate reports) Daily Production Report lacks num of products made by employee
Currently Daily Production Report provides:
- total number of manufactured units for product
- number of production entries for product
Shouldn't it also provide number of units manufactured by employee?
**Reports**
*Daily Production*      - daily, How much of each product was made?
*Employee Productivity* - during period, How much products did each employee made overall?
*Product Statistics*    - during period, How much of each product was made?
**Added:**
*Daily Employee Production* - daily, How much of each product did each employee made?
*Employee Productivity for Products* - during period, How much of each product did each employee made?

## <REJECTED> Each package builds all sqlc queries
Currently each package (slice) builds all sqlc queries across all packages.
Isn't this approach:
- confusing for developers
- crossing package boundaries
Shouldn't each package build only own `sqlc` queries and create only own query models?
see `docs/sqlc-boundaries.md`

## <REJECTED> Double id - UUID for user-facing, int for internal reference
What do you think about using UUID and int ID simultaneously? UUID for user-facing things, int for internal references (for speed).
ANSWER: Do not introduce int IDs just for speed right now. Keep one canonical identifier per entity. If performance becomes real later, benchmark and revisit with data.

## <ADDED> Sanity tests
Should we introduce sanity tests? Either implement new ones or convert some existing tests.
