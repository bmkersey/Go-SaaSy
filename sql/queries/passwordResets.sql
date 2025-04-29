-- name: CreatePasswordReset :exec
INSERT INTO password_resets (id, user_id, token, expires_at)
VALUES ($1, $2, $3, $4);


-- name: GetPasswordResetByToken :one
SELECT * FROM password_resets WHERE token = $1;


-- name: DeletePasswordReset :exec
DELETE FROM password_resets WHERE token = $1;