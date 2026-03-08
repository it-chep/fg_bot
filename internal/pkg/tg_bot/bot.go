package tg_bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"net/http"

	"fg_bot/internal/pkg/logger"
	"fg_bot/internal/pkg/tg_bot/bot_dto"
)

type Config interface {
	WebhookURL() string
	Token() string
	UseWebhook() bool
}

type Bot struct {
	bot        *tgbotapi.BotAPI
	updates    tgbotapi.UpdatesChannel
	useWebhook bool
}

func NewTgBot(cfg Config) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(cfg.Token())
	if err != nil {
		logger.Error(context.Background(), "[ERROR] NewBotAPI", err)
		return nil, err
	}

	if cfg.UseWebhook() {
		hook, _ := tgbotapi.NewWebhook(cfg.WebhookURL() + cfg.Token() + "/")
		_, err = bot.Request(hook)
		if err != nil {
			logger.Error(context.Background(), "[ERROR] Request", err)
			return nil, err
		}

		_, err = bot.GetWebhookInfo()
		if err != nil {
			logger.Error(context.Background(), "[ERROR] GetWebhookInfo", err)
			return nil, err
		}

		return &Bot{
			bot:        bot,
			useWebhook: true,
		}, nil
	}

	_, _ = bot.Request(tgbotapi.DeleteWebhookConfig{})

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	return &Bot{
		bot:        bot,
		updates:    updates,
		useWebhook: false,
	}, nil
}

func (b *Bot) HandleUpdate(r *http.Request) (*tgbotapi.Update, error) {
	update, err := b.bot.HandleUpdate(r)
	if err != nil {
		logger.Error(r.Context(), "[ERROR] HandleUpdate", err)
		return nil, err
	}
	return update, nil
}

func (b *Bot) GetUpdates() tgbotapi.UpdatesChannel {
	return b.updates
}

func (b *Bot) GetUser(chatID int64, userID int64) (bot_dto.User, error) {
	member, err := b.bot.GetChatMember(tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
			ChatID: chatID,
			UserID: userID,
		},
	})
	if err != nil {
		logger.Error(context.Background(), "[ERROR] GetChatMember", err)
		return bot_dto.User{}, err
	}

	return bot_dto.User{
		ID:       member.User.ID,
		Name:     member.User.FirstName,
		UserName: member.User.UserName,
		IsAdmin:  member.IsCreator() || member.IsAdministrator(),
	}, nil
}

func (b *Bot) SendMessage(msg bot_dto.Message, options ...MsgOption) error {
	message := tgbotapi.NewMessage(msg.Chat, msg.Text)
	if msg.ReplyToMessageID > 0 {
		message.ReplyToMessageID = msg.ReplyToMessageID
	}
	if len(msg.Buttons) != 0 {
		rows := make([][]tgbotapi.InlineKeyboardButton, 0, len(msg.Buttons))
		for _, btn := range msg.Buttons {
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(btn.Text, btn.Text),
			))
		}
		message.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
	}

	for _, opt := range options {
		message = opt(message)
	}
	_, err := b.bot.Send(message)
	return err
}

func (b *Bot) SendMessages(messages []bot_dto.Message) error {
	for _, msg := range messages {
		if err := b.SendMessage(msg); err != nil {
			return err
		}
	}
	return nil
}

func (b *Bot) SendDocument(chatID int64, fileName string, fileBytes []byte, caption string) error {
	file := tgbotapi.FileBytes{Name: fileName, Bytes: fileBytes}
	doc := tgbotapi.NewDocument(chatID, file)
	doc.Caption = caption
	_, err := b.bot.Send(doc)
	return err
}
