-- +goose Up
-- +goose StatementBegin
ALTER TABLE fg_participant
    ALTER COLUMN ping_available SET DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE fg_participant
    ALTER COLUMN ping_available SET DEFAULT TRUE;
-- +goose StatementEnd
