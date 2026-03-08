package internal

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"fg_bot/internal/modules/bot"
	"fg_bot/internal/pkg/logger"
	"fg_bot/internal/pkg/tg_bot"
	"fg_bot/internal/pkg/tg_bot/bot_dto"
	"fg_bot/internal/pkg/worker_pool"
	"fg_bot/internal/server"
	"fg_bot/internal/server/handler"
)

func (a *App) initDB(ctx context.Context) *App {
	pool, err := pgxpool.New(ctx, a.config.DatabaseURL())
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	a.pool = pool
	return a
}

func (a *App) initTgBot(_ context.Context) *App {
	if !a.config.BotIsActive() {
		return a
	}

	tgBot, err := tg_bot.NewTgBot(a.config)
	if err != nil {
		log.Fatal(err)
	}
	a.bot = tgBot
	return a
}

func (a *App) initErrorReporter(_ context.Context) *App {
	if !a.config.BotIsActive() || a.bot == nil {
		return a
	}

	adminID := a.config.ErrorAdminID()
	if adminID == 0 {
		return a
	}

	logger.SetErrorReporter(func(ctx context.Context, payload logger.ErrorPayload) error {
		location := "unknown"
		if payload.File != "" {
			location = fmt.Sprintf("%s:%d", filepath.Base(payload.File), payload.Line)
		}

		errText := "<nil>"
		if payload.Error != nil {
			errText = payload.Error.Error()
		}

		text := fmt.Sprintf(
			"Ошибка в боте\nМесто: %s\nФункция: %s\nСообщение: %s\nОшибка: %s",
			location,
			payload.Function,
			payload.Message,
			errText,
		)

		return a.bot.SendMessage(bot_dto.Message{
			Chat: adminID,
			Text: trimForTelegram(text),
		})
	})

	return a
}

func trimForTelegram(text string) string {
	const maxLen = 4096
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen-3] + "..."
}

func (a *App) initModules(_ context.Context) *App {
	a.modules = Modules{
		Bot: bot.New(a.pool, a.bot),
	}
	return a
}

func (a *App) initWorkers(_ context.Context) *App {
	if !a.config.BotIsActive() {
		a.workerPool = worker_pool.NewWorkerPool(nil)
		return a
	}

	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		location = time.FixedZone("MSK", 3*60*60)
	}

	workers := []worker_pool.Worker{
		//worker_pool.NewWorker(a.modules.Bot.Actions.ReportReminder, "0 21 * * *", location),
		worker_pool.NewWorker(a.modules.Bot.Actions.ReportReminder, "* * * * *", location),
	}

	a.workerPool = worker_pool.NewWorkerPool(workers)
	return a
}

func (a *App) initServer(_ context.Context) *App {
	h := handler.NewHandler(a.bot, a.modules.Bot, a.config)
	srv := server.New(h, a.config.Port())
	a.server = srv
	return a
}
