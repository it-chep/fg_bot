package bot

import (
	"context"
	"regexp"

	"fg_bot/internal/modules/bot/dto"
	"fg_bot/internal/pkg/tg_bot/bot_dto"
)

var reportHashtagRe = regexp.MustCompile(`(?i)#день\d+`)

func (b *Bot) Route(ctx context.Context, msg dto.Message) error {
	if match := reportHashtagRe.FindString(msg.Text); match != "" {
		return b.Actions.SaveReport.Do(ctx, msg, match)
	}

	switch msg.Text {
	case "/start":
		if msg.ChatID == msg.User {
			return b.start(msg)
		}
	case "/init_fg":
		return b.Actions.InitFG.Do(ctx, msg)
	case "/fg_statistic":
		return b.Actions.GetDailyStats.Do(ctx, msg)
	case "/fg_stat_all":
		return b.Actions.GetAllStats.Do(ctx, msg)
	case "/ping_on":
		return b.Actions.PingToggle.Do(ctx, msg, true)
	case "/ping_off":
		return b.Actions.PingToggle.Do(ctx, msg, false)
	}
	return nil
}

func (b *Bot) start(msg dto.Message) error {
	return b.bot.SendMessages([]bot_dto.Message{
		{
			Chat: msg.ChatID,
			Text: "Привет! Я помогаю вести учёт ФГ.\n\n/init_fg — создать ФГ в этом чате\n#день1, #день2 ... — отправить отчёт\n/fg_statistic — статистика за сегодня\n/fg_stat_all — полная статистика\n/ping_on — включить напоминания\n/ping_off — выключить напоминания",
		},
	})
}
