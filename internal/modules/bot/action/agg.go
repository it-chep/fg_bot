package action

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"fg_bot/internal/modules/bot/action/commands/get_all_stats"
	"fg_bot/internal/modules/bot/action/commands/get_daily_stats"
	"fg_bot/internal/modules/bot/action/commands/init_fg"
	"fg_bot/internal/modules/bot/action/commands/ping_toggle"
	"fg_bot/internal/modules/bot/action/commands/save_report"
	"fg_bot/internal/modules/bot/action/report_reminder"
	"fg_bot/internal/modules/bot/dal"
	"fg_bot/internal/pkg/tg_bot"
	"fg_bot/internal/pkg/transaction/wrapper"
)

type Agg struct {
	InitFG         *init_fg.Action
	SaveReport     *save_report.Action
	GetDailyStats  *get_daily_stats.Action
	GetAllStats    *get_all_stats.Action
	PingToggle     *ping_toggle.Action
	ReportReminder *report_reminder.Worker
}

func NewAgg(pool *pgxpool.Pool, bot *tg_bot.Bot) *Agg {
	sharedDAL := dal.New(pool)
	db := wrapper.NewDatabase(pool)
	return &Agg{
		InitFG:         init_fg.NewAction(db, bot),
		SaveReport:     save_report.NewAction(db, bot),
		GetDailyStats:  get_daily_stats.NewAction(pool, bot, sharedDAL),
		GetAllStats:    get_all_stats.NewAction(pool, bot, sharedDAL),
		PingToggle:     ping_toggle.NewAction(db, bot),
		ReportReminder: report_reminder.NewWorker(pool, bot),
	}
}
