-- name: CreateMatch :one
INSERT INTO matches (problem_id, player_one_id, player_two_id, status)
VALUES ($1, $2, $3, 'waiting')
RETURNING *;

-- name: GetMatch :one
SELECT * FROM matches
WHERE id = $1;

-- name: StartMatch :one
UPDATE matches
SET status = 'active', started_at = now()
WHERE id = $1 AND status = 'waiting'
RETURNING *;

-- name: FinishMatch :one
UPDATE matches
SET status = 'finished', finished_at = now(), winner_id = $2
WHERE id = $1 AND status = 'active'
RETURNING *;

-- name: AbandonMatch :one
UPDATE matches
SET status = 'abandoned', finished_at = now()
WHERE id = $1 AND status IN ('waiting', 'active')
RETURNING *;

-- name: ListMatchesForUser :many
SELECT * FROM matches
WHERE player_one_id = $1 OR player_two_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListSubmissionsForMatch :many
SELECT * FROM submissions
WHERE match_id = $1
ORDER BY created_at DESC;

-- name: CountMatchesForUser :one
SELECT COUNT(*) FROM matches
WHERE player_one_id = $1 OR player_two_id = $1;
