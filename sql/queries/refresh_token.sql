-- name: CreateRefreshToken :exec
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    $3,
    NULL
);

-- name: GetRefreshToken :one
SELECT token, user_id, expires_at, revoked_at FROM refresh_tokens WHERE token = $1;

-- name: GetUserFromRefreshToken :exec
UPDATE refresh_tokens
SET updated_at = NOW(),
    token = $1,
    expires_at = $2
WHERE token = $3;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET updated_at = NOW(),
    revoked_at = NOW()
WHERE token = $1;