package handler

import (
	"fmt"
	"net/http"
	"strings"

	"fg_bot/internal/modules/bot/dto"
	"fg_bot/internal/pkg/logger"
)

func (h *Handler) bot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		event, err := h.botParser.HandleUpdate(r)
		if err != nil {
			logger.Error(r.Context(), "[ERROR] Ошибка при хендлинге ивента", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if event.SentFrom() == nil || event.FromChat() == nil {
			return
		}

		txt := ""
		if event.Message != nil {
			txt = event.Message.Text
			if strings.Contains(txt, "/start ") {
				txt = txt[len("/start "):]
			}
		} else if event.CallbackQuery != nil {
			txt = event.CallbackQuery.Data
		}

		var messageID int
		if event.Message != nil {
			messageID = event.Message.MessageID
		}

		if len(event.Message.Caption) != 0 {
			txt = fmt.Sprintf("%s (отчет с медиа)", event.Message.Caption)
		}

		msg := dto.Message{
			User:      event.SentFrom().ID,
			Text:      txt,
			ChatID:    event.FromChat().ID,
			MessageID: messageID,
			ChatTitle: event.FromChat().Title,
			UserName:  event.SentFrom().UserName,
			FirstName: event.SentFrom().FirstName,
		}

		if err = h.botModule.Route(r.Context(), msg); err != nil {
			logger.Error(
				r.Context(),
				fmt.Sprintf(
					"[ERROR] Ошибка при обработке ивента, TGID: %d, username: %s",
					event.SentFrom().ID,
					event.SentFrom().UserName,
				), err)
			w.WriteHeader(http.StatusBadGateway)
			return
		}
	}
}
