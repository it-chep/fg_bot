package dao

import (
	"time"

	"fg_bot/internal/modules/bot/domain/fg"
	"fg_bot/internal/modules/bot/domain/participant"
	"fg_bot/internal/modules/bot/domain/report"
)

type FGID struct {
	ID int64 `db:"id"`
}

type FGInfo struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

func (f *FGInfo) ToDomain() *fg.FG {
	return fg.New(
		fg.WithID(f.ID),
		fg.WithName(f.Name),
	)
}

type ParticipantStatus struct {
	TgID          int64  `db:"tg_id"`
	Name          string `db:"name"`
	Username      string `db:"username"`
	ReportedToday bool   `db:"reported_today"`
}

func (p *ParticipantStatus) ToDomain() *participant.Participant {
	return participant.New(
		participant.WithTgID(p.TgID),
		participant.WithName(p.Name),
		participant.WithUsername(p.Username),
		participant.WithReportedToday(p.ReportedToday),
	)
}

type ParticipantReportCount struct {
	TgID        int64  `db:"tg_id"`
	Name        string `db:"name"`
	Username    string `db:"username"`
	ReportCount int    `db:"report_count"`
}

func (p *ParticipantReportCount) ToDomain() *participant.Participant {
	return participant.New(
		participant.WithTgID(p.TgID),
		participant.WithName(p.Name),
		participant.WithUsername(p.Username),
		participant.WithReportCount(p.ReportCount),
	)
}

type Report struct {
	ReportMessageLink string    `db:"report_message_link"`
	ReportName        string    `db:"report_name"`
	CreatedAt         time.Time `db:"created_at"`
}

func (r *Report) ToDomain() *report.Report {
	return report.New(
		report.WithReportMessageLink(r.ReportMessageLink),
		report.WithReportName(r.ReportName),
		report.WithCreatedAt(r.CreatedAt),
	)
}
