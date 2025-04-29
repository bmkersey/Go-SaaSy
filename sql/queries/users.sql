-- name: CreateUser :one
INSERT INTO users (id, email, password_hash)
VALUES (
  $1,
  $2,
  $3
)
RETURNING *;


-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;


-- name: GetUserByID :one
SELECT * FROM users
where id = $1;


-- name: UpdateUserOrg :exec
UPDATE users
SET organization_id = $2
WHERE id = $1;

-- name: GetOrgMembers :many
SELECT id, email, created_at FROM users
WHERE organization_id = $1;

-- name: UpdateUserPassword :exec
UPDATE users
SET password_hash = $1
WHERE email = $2;