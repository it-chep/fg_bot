package internal

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"fg_bot/internal/config"
	"fg_bot/internal/modules/bot"
	"fg_bot/internal/modules/bot/dto"
	"fg_bot/internal/pkg/logger"
	"fg_bot/internal/pkg/tg_bot"
	"fg_bot/internal/pkg/worker_pool"
	"fg_bot/internal/server"
)

type App struct {
	config *config.Config

	server     *server.Server
	bot        *tg_bot.Bot
	pool       *pgxpool.Pool
	workerPool worker_pool.WorkerPool

	modules Modules
}

type Modules struct {
	Bot *bot.Bot
}

func New(ctx context.Context) *App {
	cfg := config.NewConfig()

	app := &App{
		config: cfg,
	}

	app.initDB(ctx).
		initTgBot(ctx).
		initErrorReporter(ctx).
		initModules(ctx).
		initWorkers(ctx).
		initServer(ctx)

	return app
}

func (a *App) Run(ctx context.Context) {
	fmt.Printf("start server http://localhost:%s\n", a.config.Port())
	ctx = logger.ContextWithLogger(ctx, logger.New())

	if a.config.BotIsActive() {
		go a.workerPool.Run(ctx)
	}

	if !a.config.BotIsActive() || (a.config.UseWebhook() && a.config.BotIsActive()) {
		log.Fatal(a.server.ListenAndServe())
	} else {
		go func() {
			log.Fatal(a.server.ListenAndServe())
		}()
	}

	if !a.config.UseWebhook() && a.config.BotIsActive() {
		fmt.Println("Режим поллинга")
		for update := range a.bot.GetUpdates() {
			go func() {
				if update.SentFrom() == nil || update.FromChat() == nil {
					return
				}

				txt := ""
				if update.Message != nil {
					txt = update.Message.Text
					if strings.Contains(txt, "/start ") {
						txt = txt[len("/start "):]
					}
				} else if update.CallbackQuery != nil {
					txt = update.CallbackQuery.Data
				}

				var messageID int
				if update.Message != nil {
					messageID = update.Message.MessageID
				}

				if len(update.Message.Caption) != 0 {
					txt = fmt.Sprintf("%s (отчет с медиа)", update.Message.Caption)
				}

				msg := dto.Message{
					User:      update.SentFrom().ID,
					Text:      txt,
					ChatID:    update.FromChat().ID,
					MessageID: messageID,
					ChatTitle: update.FromChat().Title,
					UserName:  update.SentFrom().UserName,
					FirstName: update.SentFrom().FirstName,
				}

				if err := a.modules.Bot.Route(ctx, msg); err != nil {
					logger.Error(ctx, "Ошибка при обработке ивента", err)
				}
			}()
		}
	}
}
