## <IGNORE FOR NOW> Employee state/status
Employee.isActive -> multiple statuses (sick, vacation, out-of-work etc.)

## <POSTPONED> User-configurable product categories
Distinct entity stored in DB



I have few questions and ideas:

## 1. QUESTION 1
In raw sql queries user input is used in `WHERE` clause for filtering, example for `employee`:
```
WHERE (@query::text = ''
    OR lower(id) LIKE '%' || lower(@query::text) || '%'
    OR lower(first_name) LIKE '%' || lower(@query::text) || '%'
    OR lower(last_name) LIKE '%' || lower(@query::text) || '%'
    OR lower(email) LIKE '%' || lower(@query::text) || '%')
  AND (@active::text = '' OR is_active = (@active::text = 'true'))
```
Is this safe and sql-injection proof?

## 2. QUESTION 2
Should we use the same DB for production and testing? I know there is no production data yet, but maybe it would be worth to sort it out early?

## 3. QUESTION 3
How sqlc works (briefly)? Does it take into account migrations and raw sqls? Anything else? When to run `sqlc generate`?

## 4. custom types (enums) for PostgreSQL errors
In file `internal/products/store_postgres.go` method `mapPostgresError` strings are used to represent PostgreSQL errors.
The same in other store files: `internal/production/store_postgres.go` and `internal/employees/store_postgres.go`.
What you think about defining custom types (enums) for PostgreSQL errors and use them in place of strings? It would be more descriptive.

