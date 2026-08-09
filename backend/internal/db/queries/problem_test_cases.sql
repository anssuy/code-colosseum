-- name: CreateProblemTestCase :one
INSERT INTO problem_test_cases (problem_id, input, expected_output, is_sample, ord)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: ListTestCasesForProblem :many
SELECT * FROM problem_test_cases
WHERE problem_id = $1
ORDER BY ord;

-- name: ListSampleTestCasesForProblem :many
SELECT * FROM problem_test_cases
WHERE problem_id = $1 AND is_sample = true
ORDER BY ord;

-- name: DeleteTestCase :exec
DELETE FROM problem_test_cases
WHERE id = $1;

-- name: DeleteTestCasesForProblem :exec
DELETE FROM problem_test_cases
WHERE problem_id = $1;
