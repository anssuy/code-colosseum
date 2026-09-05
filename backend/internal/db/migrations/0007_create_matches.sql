-- +goose Up
CREATE TYPE match_status AS ENUM ('waiting', 'active', 'finished', 'abandoned');

CREATE TABLE matches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    problem_id UUID NOT NULL REFERENCES problems(id) ON DELETE RESTRICT,
    player_one_id UUID REFERENCES users(id) ON DELETE SET NULL,
    player_two_id UUID REFERENCES users(id) ON DELETE SET NULL,
    winner_id UUID REFERENCES users(id) ON DELETE SET NULL,
    status match_status NOT NULL DEFAULT 'waiting',
    started_at TIMESTAMPTZ,
    finished_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX matches_player_one_id_idx ON matches(player_one_id);
CREATE INDEX matches_player_two_id_idx ON matches(player_two_id);
CREATE INDEX matches_problem_id_idx ON matches(problem_id);

ALTER TABLE submissions ADD COLUMN match_id UUID REFERENCES matches(id) ON DELETE CASCADE;
CREATE INDEX submissions_match_id_idx ON submissions(match_id);

-- +goose Down
ALTER TABLE submissions DROP COLUMN match_id;
DROP TABLE matches;
DROP TYPE match_status;
