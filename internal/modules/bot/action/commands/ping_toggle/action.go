package ping_toggle

import (
	"context"

	"fg_bot/internal/modules/bot/action/commands/ping_toggle/dal"
	"fg_bot/internal/modules/bot/dto"
	"fg_bot/internal/pkg/tg_bot"
	"fg_bot/internal/pkg/tg_bot/bot_dto"
	"fg_bot/internal/pkg/transaction"
	"fg_bot/internal/pkg/transaction/wrapper"
)

type Action struct {
	bot *tg_bot.Bot
	dal *dal.DAL
}

func NewAction(db wrapper.Database, bot *tg_bot.Bot) *Action {
	return &Action{
		bot: bot,
		dal: dal.New(db),
	}
}

func (a *Action) Do(ctx context.Context, msg dto.Message, enabled bool) error {
	fgID, err := a.dal.GetFGByChatID(ctx, msg.ChatID)
	if err != nil {
		return a.bot.SendMessages([]bot_dto.Message{{
			Chat: msg.ChatID,
			Text: "ФГ для этого чата не найдена. Используйте /init_fg",
		}})
	}

	err = transaction.Exec(ctx, func(ctx context.Context) error {
		if err := a.dal.UpsertParticipant(ctx, msg.User, msg.FirstName, msg.UserName); err != nil {
			return err
		}
		if err := a.dal.LinkParticipantToFG(ctx, fgID, msg.User); err != nil {
			return err
		}
		return a.dal.SetPingAvailable(ctx, msg.User, enabled)
	})
	if err != nil {
		return err
	}

	text := "Напоминания об отчётах выключены"
	if enabled {
		text = "Напоминания об отчётах включены"
	}

	return a.bot.SendMessages([]bot_dto.Message{{
		Chat: msg.ChatID,
		Text: text,
	}})
}
