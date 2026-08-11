-- +goose Up
ALTER TABLE employees
    ADD COLUMN version integer NOT NULL DEFAULT 1 CHECK (version > 0);

ALTER TABLE products
    ADD COLUMN version integer NOT NULL DEFAULT 1 CHECK (version > 0);

-- +goose Down
ALTER TABLE products
    DROP COLUMN version;

ALTER TABLE employees
    DROP COLUMN version;
