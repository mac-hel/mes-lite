-- name: CreateEntry :one
INSERT INTO production_entries (
    id,
    employee_id,
    product_sku,
    quantity,
    workstation,
    occurred_at,
    comment
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING id, employee_id, product_sku, quantity, workstation, occurred_at, comment, created_at;

-- name: GetEntry :one
SELECT id, employee_id, product_sku, quantity, workstation, occurred_at, comment, created_at
FROM production_entries
WHERE id = $1;

-- name: ListEntries :many
SELECT id, employee_id, product_sku, quantity, workstation, occurred_at, comment, created_at
FROM production_entries
ORDER BY occurred_at DESC, created_at DESC;
