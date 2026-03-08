-- +goose Up
-- +goose StatementBegin

-- fg: поиск по chat_id (webhook приходит с chat_id)
create index if not exists idx_fg_chat_id on fg (chat_id);

-- fg: поиск по админу ФГ
create index if not exists idx_fg_admin_tg_id on fg (admin_tg_id);

-- reports: поиск отчётов по пользователю
create index if not exists idx_reports_tg_id on reports (tg_id);

-- reports: поиск отчётов по ФГ
create index if not exists idx_reports_fg_id on reports (fg_id);

-- reports: статистика за день (по дате создания)
create index if not exists idx_reports_created_at on reports (created_at);

-- reports: составной индекс для статистики по ФГ за период
create index if not exists idx_reports_fg_id_created_at on reports (fg_id, created_at);

-- reports: составной индекс для статистики пользователя в ФГ
create index if not exists idx_reports_tg_id_fg_id on reports (tg_id, fg_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

drop index if exists idx_fg_chat_id;
drop index if exists idx_fg_admin_tg_id;
drop index if exists idx_reports_tg_id;
drop index if exists idx_reports_fg_id;
drop index if exists idx_reports_created_at;
drop index if exists idx_reports_fg_id_created_at;
drop index if exists idx_reports_tg_id_fg_id;

-- +goose StatementEnd
