-- name: CreateUser :one
INSERT INTO users (
    username,
    email,
    password_hash
)
VALUES (
    $1,
    $2,
    $3
)
RETURNING *;


-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = $1
LIMIT 1;


-- name: GetUserByUsername :one
SELECT *
FROM users
WHERE username = $1
LIMIT 1;


-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1
LIMIT 1;


-- name: ListUsers :many
SELECT *
FROM users
ORDER BY created_at DESC;


-- name: UpdateUserRating :one
UPDATE users
SET rating = $2
WHERE id = $1
RETURNING *;


-- name: IncrementUserWins :one
UPDATE users
SET wins = wins + 1
WHERE id = $1
RETURNING *;


-- name: IncrementUserLosses :one
UPDATE users
SET losses = losses + 1
WHERE id = $1
RETURNING *;


-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
