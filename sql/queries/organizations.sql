-- name: CreateOrganization :one
INSERT INTO organizations (id, name)
VALUES ($1, $2)
RETURNING *;


-- name: GetOrganization :one
SELECT * FROM organizations
WHERE id = $1;


-- name: UpdateOrganizationBilling :exec
UPDATE organizations
SET
  stripe_customer_id = $2,
  stripe_subscription_id = $3,
  plan = $4,
  updated_at = NOW()
WHERE id = $1;