-- +goose Up
CREATE TABLE production_entries (
    id uuid PRIMARY KEY,
    employee_id text NOT NULL,
    product_sku text NOT NULL,
    quantity integer NOT NULL CHECK (quantity > 0),
    workstation text NOT NULL CHECK (btrim(workstation) <> ''),
    occurred_at timestamptz NOT NULL,
    comment text NOT NULL DEFAULT '',
    created_at timestamptz NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE production_entries;
