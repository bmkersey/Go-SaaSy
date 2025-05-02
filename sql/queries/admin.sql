-- name: ListOrganizations :many
SELECT id, name, owner_id, plan, billing_email, is_paid
FROM organizations
ORDER BY name;


-- name: ListAllUsers :many
SELECT id, email, is_admin, organization_id, created_at
FROM users
ORDER BY created_at DESC;


-- name: SetUserAdmin :exec
UPDATE users
SET is_admin = true
WHERE email = $1;


-- name: ListInvitesForOrg :many
SELECT * FROM invites
WHERE org_id = $1
ORDER BY created_at DESC;
