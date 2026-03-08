package dal

import (
	"context"
	"fg_bot/internal/modules/bot/dal/dao"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"

	"fg_bot/internal/modules/bot/domain/fg"
)

type DAL struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *DAL {
	return &DAL{pool: pool}
}

func (d *DAL) GetLatestFGByAdmin(ctx context.Context, adminTgID int64) (*fg.FG, error) {
	sql := `
		select id, name 
		from fg 
		where admin_tg_id = $1 
		order by created_at desc 
		limit 1
	`

	var row dao.FGInfo
	err := pgxscan.Get(ctx, d.pool, &row, sql, adminTgID)
	if err != nil {
		return nil, err
	}
	return row.ToDomain(), nil
}
