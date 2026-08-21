# Database

- `database/sql` as a stable standard-library abstraction for migration tooling
- using `pgx` through `database/sql` and `goose`

## sqlc
`sqlc` generates fully type-safe idiomatic code from SQL:

1. **Reads migration files:** parses migration/schema SQL to understand tables, columns, constraints and types.
    - `sqlc.yaml`.`schema` points sqlc to migrations directory
    - sqlc reads schema files as plain SQL (it has no knowledge of particular migration tools and ignores `goose` comments)
2. **Reads query files:** parses query SQL.
3. **Queries/Schema vidation:** checks that queries are valid against the schema.
4. **Infers Go types:** Infers go parameter and return types.
5. **Generates Go code** with typed methods.

Now application code calls generated methods.

Run `sqlc generate` when:
- migration/schema SQL changes
- query SQL changes
- sqlc config changes
- generated code is missing/stale
- before committing persistence changes

## Transactions

Keep transaction boundaries near the use case that requires atomicity.

Conceptually:

```go
err := db.WithTx(ctx, func(tx Tx) error {
    ...
})
```

Avoid letting unrelated lower-level repository methods independently choose transaction boundaries when several operations must commit atomically.
