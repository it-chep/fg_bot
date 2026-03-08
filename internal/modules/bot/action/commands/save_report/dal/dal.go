package dal

import (
	"context"
	"fg_bot/internal/modules/bot/dal/dao"

	"github.com/georgysavva/scany/v2/pgxscan"

	"fg_bot/internal/pkg/transaction/wrapper"
)

type DAL struct {
	db wrapper.Database
}

func New(db wrapper.Database) *DAL {
	return &DAL{db: db}
}

func (d *DAL) GetFGByChatID(ctx context.Context, chatID int64) (int64, error) {
	sql := `
		select id
		from fg
		where chat_id = $1
		order by created_at desc
		limit 1
	`

	var result dao.FGID
	err := pgxscan.Get(ctx, d.db.Pool(ctx), &result, sql, chatID)
	return result.ID, err
}

func (d *DAL) UpsertParticipant(ctx context.Context, tgID int64, name, username string) error {
	sql := `
		insert into fg_participant (tg_id, name, username) values ($1, $2, $3)
		on conflict (tg_id) do update set name = $2, username = $3
	`
	args := []interface{}{tgID, name, username}

	_, err := d.db.Pool(ctx).Exec(ctx, sql, args...)
	return err
}

func (d *DAL) SaveReport(ctx context.Context, tgID, fgID int64, reportLink, reportName string) error {
	sql := `
		insert into reports (tg_id, fg_id, report_message_link, report_name) values ($1, $2, $3, $4)
	`
	args := []interface{}{tgID, fgID, reportLink, reportName}

	_, err := d.db.Pool(ctx).Exec(ctx, sql, args...)
	return err
}

func (d *DAL) LinkParticipantToFG(ctx context.Context, fgID, tgID int64) error {
	sql := `
		insert into fg_member (fg_id, tg_id) values ($1, $2)
		on conflict (fg_id, tg_id) do nothing
	`

	_, err := d.db.Pool(ctx).Exec(ctx, sql, fgID, tgID)
	return err
}
