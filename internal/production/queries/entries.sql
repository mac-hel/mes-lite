-- name: CreateEntry :one
INSERT INTO production_entries (
    id,
    request_id,
    employee_id,
    product_sku,
    quantity,
    workstation,
    occurred_at,
    comment
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING id, employee_id, product_sku, quantity, workstation, occurred_at, comment, created_at, request_id;

-- name: GetEntry :one
SELECT id, employee_id, product_sku, quantity, workstation, occurred_at, comment, created_at, request_id
FROM production_entries
WHERE id = $1;

-- name: GetEntryByRequestID :one
SELECT id, employee_id, product_sku, quantity, workstation, occurred_at, comment, created_at, request_id
FROM production_entries
WHERE request_id = $1;

-- name: ListEntries :many
SELECT id, employee_id, product_sku, quantity, workstation, occurred_at, comment, created_at, request_id
FROM production_entries
WHERE (@employee_id::text = '' OR employee_id = @employee_id::text)
  AND (@product_sku::text = '' OR product_sku = @product_sku::text)
  AND (@workstation::text = '' OR lower(workstation) LIKE '%' || lower(@workstation::text) || '%')
  AND (NOT @from_time::boolean OR occurred_at >= @from_value)
  AND (NOT @to_time::boolean OR occurred_at < @to_value)
ORDER BY occurred_at DESC, created_at DESC, id DESC
LIMIT @limit_value OFFSET @offset_value;

-- name: CreateCorrection :one
INSERT INTO production_entry_corrections (
    id,
    entry_id,
    actor_user_id,
    reason,
    employee_id,
    product_sku,
    quantity,
    workstation,
    occurred_at,
    comment,
    created_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
) RETURNING id, entry_id, actor_user_id, reason, employee_id, product_sku, quantity, workstation, occurred_at, comment, created_at;

-- name: ListCorrections :many
SELECT id, entry_id, actor_user_id, reason, employee_id, product_sku, quantity, workstation, occurred_at, comment, created_at
FROM production_entry_corrections
WHERE entry_id = $1
ORDER BY created_at DESC, id DESC;
