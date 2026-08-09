-- name: CreateTag :one
INSERT INTO tags (slug, name)
VALUES ($1, $2)
RETURNING *;

-- name: GetTagBySlug :one
SELECT * FROM tags
WHERE slug = $1;

-- name: ListTags :many
SELECT * FROM tags
ORDER BY name;

-- name: AddTagToProblem :exec
INSERT INTO problem_tags (problem_id, tag_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: RemoveTagFromProblem :exec
DELETE FROM problem_tags
WHERE problem_id = $1 AND tag_id = $2;

-- name: ListTagsForProblem :many
SELECT t.* FROM tags t
JOIN problem_tags pt ON pt.tag_id = t.id
WHERE pt.problem_id = $1
ORDER BY t.name;

-- name: ListProblemsForTag :many
SELECT p.* FROM problems p
JOIN problem_tags pt ON pt.problem_id = p.id
WHERE pt.tag_id = $1
ORDER BY p.created_at DESC;

-- name: UpdateTag :one
UPDATE tags
SET slug = $2, name = $3
WHERE id = $1
RETURNING *;

-- name: DeleteTag :execrows
DELETE FROM tags
WHERE id = $1;
