-- +goose Up
ALTER TABLE production_orders
ADD COLUMN version integer NOT NULL DEFAULT 1 CHECK (version > 0);

-- +goose Down
ALTER TABLE production_orders
DROP COLUMN version;
