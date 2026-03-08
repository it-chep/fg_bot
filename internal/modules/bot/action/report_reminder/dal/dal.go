package dal

import (
	"context"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DAL struct {
	pool *pgxpool.Pool
}

type ParticipantReminderCandidate struct {
	FGID     int64  `db:"fg_id"`
	ChatID   int64  `db:"chat_id"`
	TgID     int64  `db:"tg_id"`
	Name     string `db:"name"`
	Username string `db:"username"`
}

func New(pool *pgxpool.Pool) *DAL {
	return &DAL{pool: pool}
}

func (d *DAL) GetParticipantsWithoutTodayReport(
	ctx context.Context,
	dayStart time.Time,
	dayEnd time.Time,
) ([]ParticipantReminderCandidate, error) {
	sql := `
		SELECT
			fm.fg_id,
			f.chat_id,
			fp.tg_id,
			COALESCE(NULLIF(fp.username, ''), '') as username,
			COALESCE(NULLIF(fp.username, ''), fp.name) as name
		FROM fg_member fm
		JOIN fg f ON f.id = fm.fg_id
		JOIN fg_participant fp ON fp.tg_id = fm.tg_id
		WHERE fp.ping_available = TRUE
		  AND NOT EXISTS (
			SELECT 1
			FROM reports r
			WHERE r.fg_id = fm.fg_id
			  AND r.tg_id = fm.tg_id
			  AND r.created_at >= $1
			  AND r.created_at < $2
		  )
		  AND NOT EXISTS (
			SELECT 1
			FROM report_reminders rr
			WHERE rr.fg_id = fm.fg_id
			  AND rr.tg_id = fm.tg_id
			  AND rr.remind_date = $1::date
		  )
	`

	var rows []ParticipantReminderCandidate
	if err := pgxscan.Select(ctx, d.pool, &rows, sql, dayStart, dayEnd); err != nil {
		return nil, err
	}
	return rows, nil
}

func (d *DAL) MarkReminderSent(ctx context.Context, fgID, tgID int64, remindDate time.Time) error {
	sql := `
		INSERT INTO report_reminders (fg_id, tg_id, remind_date)
		VALUES ($1, $2, $3::date)
		ON CONFLICT (fg_id, tg_id, remind_date) DO NOTHING
	`

	_, err := d.pool.Exec(ctx, sql, fgID, tgID, remindDate)
	return err
}
