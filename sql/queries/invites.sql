-- name: CreateInvite :exec
INSERT INTO invites (id, org_id, token, email, expires_at)
VALUES ($1, $2, $3, $4, $5);


-- name: GetInviteByToken :one
SELECT * FROM invites
WHERE token = $1 AND used_at IS NULL AND expires_at > NOW();


-- name: MarkInviteUsed :exec
UPDATE invites
SET used_at = NOW()
WHERE id = $1;


-- name: DeleteInvite :exec
DELETE FROM invites
WHERE id = $1;