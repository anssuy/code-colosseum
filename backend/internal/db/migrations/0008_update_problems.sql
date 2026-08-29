-- +goose Up
ALTER TABLE problems
    ADD COLUMN function_name TEXT NOT NULL DEFAULT '',
    ADD COLUMN params JSONB NOT NULL DEFAULT '[]',
    ADD COLUMN return_type TEXT NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE problems
    DROP COLUMN function_name,
    DROP COLUMN params,
    DROP COLUMN return_type;
