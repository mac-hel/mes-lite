-- +goose Up
CREATE INDEX production_entries_reporting_idx
    ON production_entries (occurred_at, product_sku, employee_id)
    INCLUDE (quantity);

-- +goose Down
DROP INDEX production_entries_reporting_idx;
