-- name: CreateEmployee :one
INSERT INTO employees (
    id,
    first_name,
    last_name,
    email,
    is_active,
    version
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING id, first_name, last_name, email, is_active, version, created_at, updated_at;

-- name: GetEmployee :one
SELECT id, first_name, last_name, email, is_active, version, created_at, updated_at
FROM employees
WHERE id = $1;

-- name: ListEmployees :many
SELECT id, first_name, last_name, email, is_active, version, created_at, updated_at
FROM employees
WHERE (@query::text = ''
    OR lower(id) LIKE '%' || lower(@query::text) || '%'
    OR lower(first_name) LIKE '%' || lower(@query::text) || '%'
    OR lower(last_name) LIKE '%' || lower(@query::text) || '%'
    OR lower(email) LIKE '%' || lower(@query::text) || '%')
  AND (@active::text = '' OR is_active = (@active::text = 'true'))
ORDER BY
    CASE WHEN @sort::text = 'id' THEN id END ASC,
    CASE WHEN @sort::text = '-id' THEN id END DESC,
    CASE WHEN @sort::text = 'name' THEN last_name END ASC,
    CASE WHEN @sort::text = 'name' THEN first_name END ASC,
    CASE WHEN @sort::text = '-name' THEN last_name END DESC,
    CASE WHEN @sort::text = '-name' THEN first_name END DESC,
    CASE WHEN @sort::text = 'email' THEN email END ASC,
    CASE WHEN @sort::text = '-email' THEN email END DESC,
    id ASC
LIMIT @limit_value OFFSET @offset_value;

-- name: UpdateEmployee :one
UPDATE employees
SET first_name = $2,
    last_name = $3,
    email = $4,
    is_active = $5,
    version = version + 1,
    updated_at = now()
WHERE id = $1 AND version = $6
RETURNING id, first_name, last_name, email, is_active, version, created_at, updated_at;
