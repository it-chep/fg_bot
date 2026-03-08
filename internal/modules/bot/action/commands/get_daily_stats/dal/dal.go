package dal

import (
	"context"
	"fg_bot/internal/modules/bot/dal/dao"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"

	"fg_bot/internal/modules/bot/domain/participant"
)

type DAL struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *DAL {
	return &DAL{pool: pool}
}

func (d *DAL) GetParticipantsWithTodayStatus(ctx context.Context, fgID int64) ([]*participant.Participant, error) {
	sql := `
		SELECT fp.tg_id, fp.name, fp.username,
			EXISTS(SELECT 1 FROM reports r WHERE r.tg_id = fp.tg_id AND r.fg_id = $1
				AND r.created_at::date = CURRENT_DATE) as reported_today
		FROM fg_participant fp
		WHERE fp.tg_id IN (SELECT DISTINCT tg_id FROM reports WHERE fg_id = $1)
	`

	var rows []dao.ParticipantStatus
	err := pgxscan.Select(ctx, d.pool, &rows, sql, fgID)
	if err != nil {
		return nil, err
	}

	result := make([]*participant.Participant, 0, len(rows))
	for _, r := range rows {
		result = append(result, r.ToDomain())
	}
	return result, nil
}
