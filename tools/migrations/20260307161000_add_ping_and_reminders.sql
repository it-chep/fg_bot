-- +goose Up
-- +goose StatementBegin
ALTER TABLE fg_participant
    ADD COLUMN IF NOT EXISTS ping_available BOOLEAN NOT NULL DEFAULT TRUE;

CREATE TABLE IF NOT EXISTS fg_member
(
    fg_id      BIGINT    NOT NULL,
    tg_id      BIGINT    NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (fg_id, tg_id)
);

INSERT INTO fg_member (fg_id, tg_id)
SELECT DISTINCT fg_id, tg_id
FROM reports
ON CONFLICT (fg_id, tg_id) DO NOTHING;

CREATE TABLE IF NOT EXISTS report_reminders
(
    id          BIGSERIAL PRIMARY KEY,
    fg_id       BIGINT    NOT NULL,
    tg_id       BIGINT    NOT NULL,
    remind_date DATE      NOT NULL DEFAULT CURRENT_DATE,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE (fg_id, tg_id, remind_date)
);

CREATE INDEX IF NOT EXISTS idx_fg_member_tg_id ON fg_member (tg_id);
CREATE INDEX IF NOT EXISTS idx_report_reminders_remind_date ON report_reminders (remind_date);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS report_reminders;

DROP TABLE IF EXISTS fg_member;

ALTER TABLE fg_participant
    DROP COLUMN IF EXISTS ping_available;
-- +goose StatementEnd
