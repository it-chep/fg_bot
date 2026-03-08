package dal

import (
	"context"
	"fg_bot/internal/modules/bot/dal/dao"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"

	"fg_bot/internal/modules/bot/domain/participant"
	"fg_bot/internal/modules/bot/domain/report"
)

type DAL struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *DAL {
	return &DAL{pool: pool}
}

func (d *DAL) GetParticipantsReportCounts(ctx context.Context, fgID int64) ([]*participant.Participant, error) {
	sql := `
		select fp.tg_id, fp.name, fp.username, COUNT(r.id) as report_count
		from fg_participant fp
		join reports r on r.tg_id = fp.tg_id and r.fg_id = $1
		group by fp.tg_id, fp.name, fp.username
		order by report_count desc
	`

	var rows []dao.ParticipantReportCount
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

func (d *DAL) GetReportsByParticipant(ctx context.Context, tgID, fgID int64) ([]*report.Report, error) {
	sql := `
		SELECT report_message_link, report_name, created_at FROM reports
		WHERE tg_id = $1 AND fg_id = $2
		ORDER BY created_at DESC
	`

	args := []interface{}{tgID, fgID}

	var rows []dao.Report
	err := pgxscan.Select(ctx, d.pool, &rows, sql, args...)
	if err != nil {
		return nil, err
	}

	result := make([]*report.Report, 0, len(rows))
	for _, r := range rows {
		result = append(result, r.ToDomain())
	}
	return result, nil
}
