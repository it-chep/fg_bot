package save_report

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"

	"fg_bot/internal/modules/bot/action/commands/save_report/dal"
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

var reportAcceptedTemplates = []string{
	"Молодец, продолжай в том же духе!",
	"Ты умничка!",
	"Принял твой отчет!",
	"Вот это да!",
	"Отличный темп, так держать!",
	"Зачет, отчет на месте!",
	"Супер, отчет принят!",
	"Красиво! Отчет зафиксирован.",
	"Огонь! Продолжаем.",
	"Мощно, беру в учет!",
	"Все четко, принято!",
	"Прекрасно, отчет записал!",
	"Крутая работа, принято!",
	"Респект, отчет засчитан!",
	"Топ! Так держать.",
	"Ты в игре, отчет принят.",
	"Отлично сработано!",
	"Спасибо, отчет добавлен!",
	"Сильный ход, принято!",
	"Стабильно хорошо, зафиксировал отчет!",
}

func NewAction(db wrapper.Database, bot *tg_bot.Bot) *Action {
	return &Action{
		bot: bot,
		dal: dal.New(db),
	}
}

func (a *Action) Do(ctx context.Context, msg dto.Message, reportName string) error {
	user, err := a.bot.GetUser(msg.ChatID, msg.User)
	if err != nil {
		return err
	}
	if user.IsAdmin {
		return nil
	}

	fgID, err := a.dal.GetFGByChatID(ctx, msg.ChatID)
	if err != nil {
		return a.bot.SendMessages([]bot_dto.Message{
			{Chat: msg.ChatID, Text: "ФГ для этого чата не найдена. Используйте /init_fg"},
		})
	}

	// Строим ссылку на сообщение: для групп chat_id отрицательный, убираем -100 префикс
	chatIDStripped := msg.ChatID
	if chatIDStripped < 0 {
		chatIDStripped = -chatIDStripped - 1000000000000
	}
	reportLink := fmt.Sprintf("https://t.me/c/%d/%d", chatIDStripped, msg.MessageID)

	err = transaction.Exec(ctx, func(ctx context.Context) error {
		if err := a.dal.UpsertParticipant(ctx, msg.User, msg.FirstName, msg.UserName); err != nil {
			return err
		}
		if err := a.dal.LinkParticipantToFG(ctx, fgID, msg.User); err != nil {
			return err
		}

		return a.dal.SaveReport(ctx, msg.User, fgID, reportLink, reportName)
	})
	if err != nil {
		return err
	}

	return a.bot.SendMessages([]bot_dto.Message{
		{
			Chat:             msg.ChatID,
			Text:             randomAcceptedMessage(),
			ReplyToMessageID: msg.MessageID,
		},
	})
}

func randomAcceptedMessage() string {
	if len(reportAcceptedTemplates) == 0 {
		return "Отчет принят!"
	}

	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(len(reportAcceptedTemplates))))
	if err != nil {
		return "Отчет принят!"
	}

	return reportAcceptedTemplates[nBig.Int64()]
}
