-- name: CreateUser :one
INSERT INTO auth_users (
    id,
    email,
    password_hash,
    role,
    is_active
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING id, email, password_hash, role, is_active, created_at, updated_at;

-- name: GetUserByEmail :one
SELECT id, email, password_hash, role, is_active, created_at, updated_at
FROM auth_users
WHERE lower(email) = lower($1);

-- name: UserExistsByEmail :one
SELECT EXISTS(
    SELECT 1
    FROM auth_users
    WHERE lower(email) = lower($1)
);
