-- +goose Up
CREATE TABLE submissions (
   id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
   user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
   problem_id UUID NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
   language TEXT NOT NULL,
   source_code TEXT NOT NULL,
   status TEXT NOT NULL DEFAULT 'pending',
   passed_tests INT NOT NULL DEFAULT 0,
   total_tests INT NOT NULL DEFAULT 0,
   execution_time_ms BIGINT,
   created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX submissions_user_id_idx ON submissions(user_id);
CREATE INDEX submissions_problem_id_idx ON submissions(problem_id);

-- +goose Down
DROP TABLE submissions;