-- +goose Up
CREATE TABLE organizations (
  id UUID PRIMARY KEY,
  name TEXT NOT NULL UNIQUE,
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now()
);

ALTER TABLE users
ADD CONSTRAINT fk_organization
FOREIGN KEY (organization_id) 
REFERENCES organizations(id) ON DELETE SET NULL;

-- +goose Down
ALTER TABLE users
DROP COLUMN organization_id;

DROP TABLE organizations;