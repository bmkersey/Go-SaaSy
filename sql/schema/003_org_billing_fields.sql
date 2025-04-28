-- +goose Up
ALTER TABLE organizations
ADD COLUMN stripe_customer_id TEXT,
ADD COLUMN stripe_subscription_id TEXT,
ADD COLUMN plan TEXT DEFAULT 'free',
ADD COLUMN billing_email TEXT;

-- +goose Down
ALTER TABLE organizations
DROP COLUMN stripe_customer_id
DROP COLUMN stripe_subscription_id
DROP COLUMN plan
DROP COLUMN billing_email;
