-- name: InsertRefreshToken :exec
INSERT INTO identity.refresh_tokens (user_id, token_hash, expires_at)
VALUES ($1, $2, $3);

-- name: GetRefreshTokenByHash :one
SELECT * FROM identity.refresh_tokens WHERE token_hash = $1;

-- name: RevokeRefreshToken :exec
UPDATE identity.refresh_tokens SET revoked_at = now()
WHERE id = $1 AND revoked_at IS NULL;