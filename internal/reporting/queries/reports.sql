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
