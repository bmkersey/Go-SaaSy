-- name: CreateOrganization :one
INSERT INTO organizations (id, name, owner_id)
VALUES ($1, $2, $3)
RETURNING id, name, owner_id;


-- name: GetOrganization :one
SELECT * FROM organizations
WHERE id = $1;


-- name: UpdateOrganizationBilling :exec
UPDATE organizations
SET
  stripe_customer_id = $2,
  stripe_subscription_id = $3,
  plan = $4,
  is_paid = $5,
  updated_at = NOW()
WHERE id = $1;