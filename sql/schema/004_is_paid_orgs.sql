-- +goose Up
ALTER TABLE organizations
ADD COLUMN is_paid BOOLEAN NOT NULL DEFAULT FALSE;


-- +goose Down
ALTER TABLE organizations
DROP COLUMN is_paid;