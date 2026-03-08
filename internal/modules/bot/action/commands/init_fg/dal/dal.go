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

func (d *DAL) UpsertAdmin(ctx context.Context, tgID int64, name, username string) error {
	sql := `
		insert into fg_admin (tg_id, name, username) values ($1, $2, $3)
		on conflict (tg_id) do update set name = $2, username = $3
	`
	args := []interface{}{tgID, name, username}

	_, err := d.db.Pool(ctx).Exec(ctx, sql, args...)
	return err
}

func (d *DAL) CreateFG(ctx context.Context, name string, chatID, adminTgID int64) (int64, error) {
	sql := `
		insert into fg (name, chat_id, admin_tg_id) values ($1, $2, $3) returning id
	`
	args := []interface{}{name, chatID, adminTgID}

	var result dao.FGID
	err := pgxscan.Get(ctx, d.db.Pool(ctx), &result, sql, args...)
	return result.ID, err
}
