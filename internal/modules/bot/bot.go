package bot

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"fg_bot/internal/modules/bot/action"
	"fg_bot/internal/pkg/tg_bot"
)

type Bot struct {
	Actions *action.Agg
	bot     *tg_bot.Bot
}

func New(pool *pgxpool.Pool, tgBot *tg_bot.Bot) *Bot {
	return &Bot{
		Actions: action.NewAgg(pool, tgBot),
		bot:     tgBot,
	}
}
