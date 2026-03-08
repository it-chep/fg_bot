package handler

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"fg_bot/internal/modules/bot"
	"fg_bot/internal/server/middleware"
)

type Config interface {
	Token() string
}

type TgHookParser interface {
	HandleUpdate(r *http.Request) (*tgbotapi.Update, error)
}

type Handler struct {
	router    *chi.Mux
	botParser TgHookParser
	botModule *bot.Bot
}

func NewHandler(botParser TgHookParser, botModule *bot.Bot, cfg Config) *Handler {
	h := &Handler{
		router:    chi.NewRouter(),
		botParser: botParser,
		botModule: botModule,
	}

	h.setupMiddleware()
	h.setupRoutes(cfg)

	return h
}

func (h *Handler) setupMiddleware() {
	h.router.Use(middleware.LoggerMiddleware)
}

func (h *Handler) setupRoutes(cfg Config) {
	h.router.Post(fmt.Sprintf("/%s/", cfg.Token()), h.bot())
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}
