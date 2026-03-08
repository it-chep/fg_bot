package tg_bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type MsgOption func(tgbotapi.MessageConfig) tgbotapi.MessageConfig

func WithDisabledPreview() MsgOption {
	return func(msg tgbotapi.MessageConfig) tgbotapi.MessageConfig {
		msg.DisableWebPagePreview = true
		return msg
	}
}

func WithParseModeHTML() MsgOption {
	return func(msg tgbotapi.MessageConfig) tgbotapi.MessageConfig {
		msg.ParseMode = tgbotapi.ModeHTML
		return msg
	}
}
