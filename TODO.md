## <IGNORE FOR NOW> Employee state/status
Employee.isActive -> multiple statuses (sick, vacation, out-of-work etc.)

## <POSTPONED> User-configurable product categories
Distinct entity stored in DB

## <PARTIALY> Encapsulation - hide private properties of entities (structs)
Currently all domain Entities have public fields. Would it be feasible to make fields private at current phase of project?
- <DONE> Order

## Fields as Collections
What do you think about using collection as fields that contain multiple values, e.g. `Order.lines` (`[]OrderLine`) or `Order.assignedEmployees` (`[]string`)?
