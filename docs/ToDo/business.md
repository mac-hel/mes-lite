## <IGNORE FOR NOW> Employee state/status
Employee.isActive -> multiple statuses (sick, vacation, out-of-work etc.)

## <POSTPONED> User-configurable product categories
Distinct entity stored in DB

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
