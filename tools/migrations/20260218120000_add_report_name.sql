-- +goose Up
-- +goose StatementBegin
ALTER TABLE reports ADD COLUMN report_name TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE reports DROP COLUMN report_name;
-- +goose StatementEnd
