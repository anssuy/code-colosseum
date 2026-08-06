-- name: CreateProblem :one
INSERT INTO problems (title, slug, difficulty, description)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetProblemByID :one
SELECT * FROM problems
WHERE id = $1;

-- name: GetProblemBySlug :one
SELECT * FROM problems
WHERE slug = $1;

-- name: ListProblems :many
SELECT * FROM problems
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListProblemsByDifficulty :many
SELECT * FROM problems
WHERE difficulty = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountProblems :one
SELECT COUNT(*) FROM problems;

-- name: CountProblemsByDifficulty :one
SELECT COUNT(*) FROM problems
WHERE difficulty = $1;

-- name: UpdateProblem :one
UPDATE problems
SET title = $2, slug = $3, difficulty = $4, description = $5
WHERE id = $1
RETURNING *;

-- name: DeleteProblem :execrows
DELETE FROM problems
WHERE id = $1;
