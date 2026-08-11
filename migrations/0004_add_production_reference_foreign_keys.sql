-- +goose Up
ALTER TABLE production_entries
    ADD CONSTRAINT production_entries_employee_id_fkey
    FOREIGN KEY (employee_id) REFERENCES employees(id)
    ON DELETE RESTRICT
    NOT VALID;

ALTER TABLE production_entries
    ADD CONSTRAINT production_entries_product_sku_fkey
    FOREIGN KEY (product_sku) REFERENCES products(sku)
    ON DELETE RESTRICT
    NOT VALID;

-- +goose Down
ALTER TABLE production_entries
    DROP CONSTRAINT production_entries_product_sku_fkey;

ALTER TABLE production_entries
    DROP CONSTRAINT production_entries_employee_id_fkey;
