-- +goose Up
CREATE TYPE difficulty AS ENUM ('easy', 'medium', 'hard');

CREATE TABLE problems (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    difficulty difficulty NOT NULL,
    description TEXT NOT NULL,
    time_limit_ms INTEGER NOT NULL DEFAULT 2000,
    memory_limit_mb INTEGER NOT NULL DEFAULT 256,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX problems_difficulty_idx ON problems(difficulty);

-- +goose Down
DROP TABLE problems;
DROP TYPE difficulty;
