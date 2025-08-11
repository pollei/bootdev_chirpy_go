-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at)
    VALUES ( $1, NOW(), NOW(),  $2, $3 )
    RETURNING *;

-- name: GetRefreshTokenByToken :one
SELECT * FROM refresh_tokens
    WHERE token = $1  LIMIT 1;

-- name: RevokeRefreshToken :one
UPDATE refresh_tokens SET revoked_at = NOW() WHERE token = $1
    RETURNING *;
  -- update the revoked_at to time now

-- name: GetRefreshTokensByUserID :many
SELECT * FROM refresh_tokens
WHERE user_id = $1 ORDER BY created_at ASC LIMIT $2;

-- name: DeleteAllRefreshTokens :exec
DELETE FROM refresh_tokens;

-- name: GetAllRefreshTokens :many
SELECT * FROM refresh_tokens ORDER BY created_at ASC;
