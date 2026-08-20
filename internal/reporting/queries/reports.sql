-- name: DailyProduction :many
SELECT
    date_trunc('day', occurred_at AT TIME ZONE 'UTC')::timestamptz AS day,
    product_sku,
    sum(quantity)::bigint AS total_quantity,
    count(*)::bigint AS entry_count
FROM production_entries
WHERE occurred_at >= sqlc.arg(from_time)
  AND occurred_at < sqlc.arg(to_time)
GROUP BY day, product_sku
ORDER BY day ASC, product_sku ASC;

-- name: EmployeeProductivity :many
SELECT
    pe.employee_id,
    e.first_name,
    e.last_name,
    sum(pe.quantity)::bigint AS total_quantity,
    count(*)::bigint AS entry_count
FROM production_entries pe
JOIN employees e ON e.id = pe.employee_id
WHERE pe.occurred_at >= sqlc.arg(from_time)
  AND pe.occurred_at < sqlc.arg(to_time)
GROUP BY pe.employee_id, e.first_name, e.last_name
ORDER BY total_quantity DESC, entry_count DESC, pe.employee_id ASC;

-- name: ProductStatistics :many
SELECT
    pe.product_sku,
    p.name AS product_name,
    sum(pe.quantity)::bigint AS total_quantity,
    count(*)::bigint AS entry_count,
    count(DISTINCT pe.employee_id)::bigint AS employee_count
FROM production_entries pe
JOIN products p ON p.sku = pe.product_sku
WHERE pe.occurred_at >= sqlc.arg(from_time)
  AND pe.occurred_at < sqlc.arg(to_time)
GROUP BY pe.product_sku, p.name
ORDER BY total_quantity DESC, entry_count DESC, pe.product_sku ASC;
