-- +goose Up
CREATE TABLE employees (
    id text PRIMARY KEY CHECK (btrim(id) <> ''),
    first_name text NOT NULL CHECK (btrim(first_name) <> ''),
    last_name text NOT NULL CHECK (btrim(last_name) <> ''),
    email text NOT NULL CHECK (btrim(email) <> ''),
    is_active boolean NOT NULL DEFAULT true,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE products (
    sku text PRIMARY KEY CHECK (btrim(sku) <> ''),
    name text NOT NULL CHECK (btrim(name) <> ''),
    category integer NOT NULL CHECK (category BETWEEN 0 AND 4),
    unit text NOT NULL CHECK (btrim(unit) <> ''),
    is_active boolean NOT NULL DEFAULT true,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE products;
DROP TABLE employees;
