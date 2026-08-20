-- name: CreateSubmission :one
INSERT INTO submissions (
    user_id,
    problem_id,
    language,
    source_code,
    status,
    passed_tests,
    total_tests,
    execution_time_ms
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetSubmissionByID :one
SELECT * FROM submissions
WHERE id = $1
LIMIT 1;

-- name: ListSubmissionsForUser :many
SELECT * FROM submissions
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: ListSubmissionsForProblem :many
SELECT * FROM submissions
WHERE problem_id = $1
ORDER BY created_at DESC;