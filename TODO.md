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

## Double id - UUID for user-facing, int for internal reference

## Sanity tests?
