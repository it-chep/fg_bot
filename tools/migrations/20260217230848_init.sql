-- +goose Up
-- +goose StatementBegin

create table if not exists fg
(
    id          bigserial primary key,
    name        text,                   -- название чата, то есть фг
    chat_id     bigint,                 -- id чата в телеге
    admin_tg_id bigint,                 -- id админа - ведущего ФГ
    created_at  timestamp default now() -- Дата создания ФГ
);


create table if not exists fg_admin
(
    tg_id    bigint primary key, -- id из телеги пользователя
    name     text,               -- имя пользака
    username text                -- username пользака из тг
);


create table if not exists fg_participant
(
    tg_id    bigint primary key, -- id из телеги пользователя
    name     text,               -- имя пользака
    username text                -- username пользака из тг
);

create table if not exists reports
(
    id                  bigserial primary key,
    tg_id               bigint,                 -- id пользака
    fg_id               bigint,                 -- id ФГ
    report_message_link text,                   -- ссылка на сообщение, чтобы можно было вернуться к отчету и проверить его
    created_at          timestamp default now() -- дата создания отчета
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
