-- +goose Up
-- +goose StatementBegin
ALTER TABLE download_events 
    ADD COLUMN IF NOT EXISTS user_first_name TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS user_last_name TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS user_username TEXT NOT NULL DEFAULT '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE download_events 
    DROP COLUMN IF EXISTS user_first_name,
    DROP COLUMN IF EXISTS user_last_name,
    DROP COLUMN IF EXISTS user_username;
-- +goose StatementEnd
