-- +goose Up
ALTER TABLE organizations
ADD COLUMN owner_id UUID NOT NULL;

ALTER TABLE organizations
ADD CONSTRAINT fk_organizations_owner_id
FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE;

-- +goose Down
ALTER TABLE organizations DROP CONSTRAINT fk_organizations_owner_id;
ALTER TABLE organizations DROP COLUMN owner_id;
