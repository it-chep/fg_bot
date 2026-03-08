package get_daily_stats

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"fg_bot/internal/modules/bot/action/commands/get_daily_stats/dal"
	sharedDal "fg_bot/internal/modules/bot/dal"
	"fg_bot/internal/modules/bot/dto"
	"fg_bot/internal/pkg/tg_bot"
	"fg_bot/internal/pkg/tg_bot/bot_dto"
)

type Action struct {
	bot       *tg_bot.Bot
	dal       *dal.DAL
	sharedDAL *sharedDal.DAL
}

func NewAction(pool *pgxpool.Pool, bot *tg_bot.Bot, shared *sharedDal.DAL) *Action {
	return &Action{
		bot:       bot,
		dal:       dal.New(pool),
		sharedDAL: shared,
	}
}

func (a *Action) Do(ctx context.Context, msg dto.Message) error {
	fg, err := a.sharedDAL.GetLatestFGByAdmin(ctx, msg.User)
	if err != nil {
		return a.bot.SendMessages([]bot_dto.Message{
			{Chat: msg.ChatID, Text: "У вас нет созданных ФГ. Используйте /init_fg"},
		})
	}

	participants, err := a.dal.GetParticipantsWithTodayStatus(ctx, fg.GetID())
	if err != nil {
		return err
	}

	if len(participants) == 0 {
		return a.bot.SendMessages([]bot_dto.Message{
			{Chat: msg.ChatID, Text: fmt.Sprintf("Статистика ФГ «%s» за сегодня:\n\nПока нет участников с отчётами.", fg.GetName())},
		})
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Статистика ФГ «%s» за сегодня:\n\n", fg.GetName()))

	for _, p := range participants {
		username := ""
		if p.GetUsername() != "" {
			username = fmt.Sprintf(" (@%s)", p.GetUsername())
		}
		if p.GetReportedToday() {
			sb.WriteString(fmt.Sprintf("✅ %s%s — отчёт сдан\n", p.GetName(), username))
		} else {
			sb.WriteString(fmt.Sprintf("❌ %s%s — отчёт не сдан\n", p.GetName(), username))
		}
	}

	return a.bot.SendMessages([]bot_dto.Message{
		{Chat: msg.ChatID, Text: sb.String()},
	})
}
