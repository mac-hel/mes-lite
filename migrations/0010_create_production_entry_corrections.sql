-- +goose Up
CREATE TABLE production_entry_corrections (
    id uuid PRIMARY KEY,
    entry_id uuid NOT NULL REFERENCES production_entries(id),
    actor_user_id text NOT NULL REFERENCES auth_users(id),
    reason text NOT NULL CHECK (btrim(reason) <> ''),
    employee_id text NOT NULL REFERENCES employees(id),
    product_sku text NOT NULL REFERENCES products(sku),
    quantity integer NOT NULL CHECK (quantity > 0),
    workstation text NOT NULL CHECK (btrim(workstation) <> ''),
    occurred_at timestamptz NOT NULL,
    comment text NOT NULL DEFAULT '',
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX production_entry_corrections_entry_id_created_at_idx
    ON production_entry_corrections (entry_id, created_at DESC);

-- +goose Down
DROP TABLE production_entry_corrections;
