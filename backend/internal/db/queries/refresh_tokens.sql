-- name: CreateRefreshToken :exec
INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
VALUES ($1,$2,$3);


-- name: RotateRefreshToken :one
WITH consumed_token AS (
    DELETE
    FROM refresh_tokens AS old_token
    WHERE old_token.token_hash = $1
    AND old_token.expires_at > now() RETURNING old_token.user_id
)
INSERT
INTO refresh_tokens AS new_token (user_id, token_hash, expires_at)
SELECT consumed_token.user_id, $2, $3
FROM consumed_token RETURNING new_token.user_id;


-- name: DeleteRefreshTokenByHash :exec
DELETE
FROM refresh_tokens
WHERE token_hash = $1;


-- name: DeleteRefreshTokensByUserID :exec
DELETE
FROM refresh_tokens
WHERE user_id = $1;


-- name: DeleteExpiredRefreshTokens :exec
DELETE
FROM refresh_tokens
WHERE expires_at <= now();