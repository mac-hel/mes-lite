-- name: CreateOrder :one
INSERT INTO production_orders (
    id,
    status,
    version,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING id, status, version, created_at, updated_at;

-- name: CreateOrderLine :exec
INSERT INTO production_order_lines (
    order_id,
    product_sku,
    planned_quantity
) VALUES (
    $1, $2, $3
);

-- name: CreateOrderAssignment :exec
INSERT INTO production_order_assignments (
    order_id,
    employee_id
) VALUES (
    $1, $2
);

-- name: GetOrder :one
SELECT id, status, version, created_at, updated_at
FROM production_orders
WHERE id = $1;

-- name: UpdateOrder :one
UPDATE production_orders
SET status = $2,
    updated_at = $3,
    version = version + 1
WHERE id = $1 AND version = $4
RETURNING id, status, version, created_at, updated_at;

-- name: DeleteOrderAssignments :exec
DELETE FROM production_order_assignments
WHERE order_id = $1;

-- name: ListOrderLines :many
SELECT order_id, product_sku, planned_quantity
FROM production_order_lines
WHERE order_id = $1
ORDER BY product_sku ASC;

-- name: ListOrderAssignments :many
SELECT order_id, employee_id
FROM production_order_assignments
WHERE order_id = $1
ORDER BY employee_id ASC;
