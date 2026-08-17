-- +goose Up
CREATE TABLE production_orders (
    id text PRIMARY KEY CHECK (btrim(id) <> ''),
    status text NOT NULL CHECK (status IN ('draft', 'released', 'in_progress', 'completed', 'cancelled')),
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL
);

CREATE TABLE production_order_lines (
    order_id text NOT NULL REFERENCES production_orders(id) ON DELETE CASCADE,
    product_sku text NOT NULL REFERENCES products(sku),
    planned_quantity integer NOT NULL CHECK (planned_quantity > 0),
    PRIMARY KEY (order_id, product_sku),
    CHECK (btrim(product_sku) <> '')
);

CREATE TABLE production_order_assignments (
    order_id text NOT NULL REFERENCES production_orders(id) ON DELETE CASCADE,
    employee_id text NOT NULL REFERENCES employees(id),
    PRIMARY KEY (order_id, employee_id),
    CHECK (btrim(employee_id) <> '')
);

-- +goose Down
DROP TABLE production_order_assignments;
DROP TABLE production_order_lines;
DROP TABLE production_orders;
