-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, revoked_at, user_id, expires_at)
VALUES (
    $1,
    NOW(),
    NOW(),
    NULL,
    $2,
    $3
)
RETURNING *;