-- +goose Up
CREATE TABLE auth_users (
    id text PRIMARY KEY CHECK (btrim(id) <> ''),
    email text NOT NULL UNIQUE CHECK (btrim(email) <> ''),
    password_hash bytea NOT NULL CHECK (length(password_hash) > 0),
    role text NOT NULL CHECK (role IN ('admin', 'manager', 'leader', 'worker')),
    is_active boolean NOT NULL DEFAULT true,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE auth_users;
