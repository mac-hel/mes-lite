-- +goose Up
ALTER TABLE production_entries
    ADD COLUMN request_id text NOT NULL DEFAULT '';

CREATE UNIQUE INDEX production_entries_request_id_key
    ON production_entries (request_id)
    WHERE request_id <> '';

-- +goose Down
DROP INDEX production_entries_request_id_key;

ALTER TABLE production_entries
    DROP COLUMN request_id;
