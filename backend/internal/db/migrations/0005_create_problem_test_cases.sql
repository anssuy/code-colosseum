-- +goose Up
CREATE TABLE problem_test_cases (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    problem_id UUID NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
    input TEXT NOT NULL,
    expected_output TEXT NOT NULL,
    is_sample BOOLEAN NOT NULL DEFAULT false,
    ord INT NOT NULL DEFAULT 0
);

CREATE INDEX problem_test_cases_problem_id_idx ON problem_test_cases(problem_id);

-- +goose Down
DROP TABLE problem_test_cases;
