package init_fg

import (
	"context"
	"fmt"

	"fg_bot/internal/modules/bot/action/commands/init_fg/dal"
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

func (a *Action) Do(ctx context.Context, msg dto.Message) error {
	user, err := a.bot.GetUser(msg.ChatID, msg.User)
	if err != nil {
		return err
	}
	if !user.IsAdmin {
		return a.bot.SendMessages([]bot_dto.Message{
			{Chat: msg.ChatID, Text: "Только администратор чата может создать ФГ"},
		})
	}

	name := msg.ChatTitle
	if name == "" {
		name = "Без названия"
	}

	err = transaction.Exec(ctx, func(ctx context.Context) error {
		if err := a.dal.UpsertAdmin(ctx, msg.User, msg.FirstName, msg.UserName); err != nil {
			return err
		}

		_, err := a.dal.CreateFG(ctx, name, msg.ChatID, msg.User)
		return err
	})
	if err != nil {
		return err
	}

	return a.bot.SendMessages([]bot_dto.Message{
		{
			Chat: msg.ChatID,
			Text: fmt.Sprintf("ФГ «%s» создана", name),
		},
	})
}
