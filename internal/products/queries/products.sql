-- name: CreateProduct :one
INSERT INTO products (
    sku,
    name,
    category,
    unit,
    is_active,
    version
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING sku, name, category, unit, is_active, version, created_at, updated_at;

-- name: GetProduct :one
SELECT sku, name, category, unit, is_active, version, created_at, updated_at
FROM products
WHERE sku = $1;

-- name: ListProducts :many
SELECT sku, name, category, unit, is_active, version, created_at, updated_at
FROM products
WHERE (@query::text = ''
    OR lower(sku) LIKE '%' || lower(@query::text) || '%'
    OR lower(name) LIKE '%' || lower(@query::text) || '%'
    OR lower(CASE category
        WHEN 0 THEN 'Ventilation'
        WHEN 1 THEN 'Filter'
        WHEN 2 THEN 'Duct'
        WHEN 3 THEN 'Mounting'
        WHEN 4 THEN 'Other'
        ELSE 'Unknown'
    END) LIKE '%' || lower(@query::text) || '%')
  AND (@active::text = '' OR is_active = (@active::text = 'true'))
ORDER BY
    CASE WHEN @sort::text = 'sku' THEN sku END ASC,
    CASE WHEN @sort::text = '-sku' THEN sku END DESC,
    CASE WHEN @sort::text = 'name' THEN name END ASC,
    CASE WHEN @sort::text = '-name' THEN name END DESC,
    CASE WHEN @sort::text = 'category' THEN category END ASC,
    CASE WHEN @sort::text = '-category' THEN category END DESC,
    sku ASC
LIMIT @limit_value OFFSET @offset_value;

-- name: UpdateProduct :one
UPDATE products
SET name = $2,
    category = $3,
    unit = $4,
    is_active = $5,
    version = version + 1,
    updated_at = now()
WHERE sku = $1 AND version = $6
RETURNING sku, name, category, unit, is_active, version, created_at, updated_at;
