package report_reminder

import (
	"context"
	"fmt"
	"html"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"fg_bot/internal/modules/bot/action/report_reminder/dal"
	"fg_bot/internal/pkg/logger"
	"fg_bot/internal/pkg/tg_bot"
	"fg_bot/internal/pkg/tg_bot/bot_dto"
)

type Worker struct {
	bot *tg_bot.Bot
	dal *dal.DAL
}

func NewWorker(pool *pgxpool.Pool, bot *tg_bot.Bot) *Worker {
	return &Worker{
		bot: bot,
		dal: dal.New(pool),
	}
}

func (w *Worker) Do(ctx context.Context) error {
	nowUTC := time.Now().UTC()
	todayUTC := nowUTC.Truncate(24 * time.Hour)
	tomorrowUTC := todayUTC.Add(24 * time.Hour)

	participants, err := w.dal.GetParticipantsWithoutTodayReport(
		ctx,
		todayUTC,
		tomorrowUTC,
	)
	if err != nil {
		logger.Error(ctx, "[ERROR] reminder worker: failed to load participants", err)
		return err
	}

	for _, p := range participants {
		mention := fmt.Sprintf("<a href=\"tg://user?id=%d\">%s</a>", p.TgID, html.EscapeString(p.Name))
		if p.Username != "" {
			mention = "@" + html.EscapeString(p.Username)
		}
		text := fmt.Sprintf("%s, напоминание: сегодня ещё нет отчёта. Отправьте его с тегом #деньN.", mention)

		err = w.bot.SendMessage(
			bot_dto.Message{Chat: p.ChatID, Text: text},
			tg_bot.WithParseModeHTML(),
		)
		if err != nil {
			logger.Error(ctx, "[ERROR] reminder worker: failed to send message", err)
			continue
		}

		if err := w.dal.MarkReminderSent(ctx, p.FGID, p.TgID, todayUTC); err != nil {
			logger.Error(ctx, "[ERROR] reminder worker: failed to mark reminder", err)
		}
	}

	return nil
}
