package get_all_stats

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/xuri/excelize/v2"

	"fg_bot/internal/modules/bot/action/commands/get_all_stats/dal"
	sharedDal "fg_bot/internal/modules/bot/dal"
	"fg_bot/internal/modules/bot/domain/participant"
	"fg_bot/internal/modules/bot/domain/report"
	"fg_bot/internal/modules/bot/dto"
	"fg_bot/internal/pkg/tg_bot"
	"fg_bot/internal/pkg/tg_bot/bot_dto"
)

type Action struct {
	bot       *tg_bot.Bot
	dal       *dal.DAL
	sharedDAL *sharedDal.DAL
}

func NewAction(pool *pgxpool.Pool, bot *tg_bot.Bot, shared *sharedDal.DAL) *Action {
	return &Action{
		bot:       bot,
		dal:       dal.New(pool),
		sharedDAL: shared,
	}
}

func (a *Action) Do(ctx context.Context, msg dto.Message) error {
	if os.Getenv("DEBUG") == "True" {
		return a.doDebug(msg)
	}

	fg, err := a.sharedDAL.GetLatestFGByAdmin(ctx, msg.User)
	if err != nil {
		return a.bot.SendMessages([]bot_dto.Message{
			{Chat: msg.ChatID, Text: "У вас нет созданных ФГ. Используйте /init_fg"},
		})
	}

	participants, err := a.dal.GetParticipantsReportCounts(ctx, fg.GetID())
	if err != nil {
		return err
	}

	if len(participants) == 0 {
		return a.bot.SendMessages([]bot_dto.Message{
			{Chat: msg.ChatID, Text: "Пока нет отчётов в этой ФГ."},
		})
	}

	f := excelize.NewFile()
	sheetName := "Отчёты"
	f.SetSheetName("Sheet1", sheetName)
	f.SetCellValue(sheetName, "A1", "Участник")
	f.SetCellValue(sheetName, "B1", "Дата")
	f.SetCellValue(sheetName, "C1", "Название отчёта")
	f.SetCellValue(sheetName, "D1", "Ссылка")
	row := 2

	var sb strings.Builder
	sb.WriteString("Отчёт по ФГ\n\n")

	for i, p := range participants {
		reports, err := a.dal.GetReportsByParticipant(ctx, p.GetTgID(), fg.GetID())
		if err != nil {
			return err
		}

		lastReport := ""
		if len(reports) > 0 {
			lastReport = reports[0].GetCreatedAt().Format("02.01.2006")
		}
		sb.WriteString(fmt.Sprintf("%d. %s — %d отч., крайний: %s\n", i+1, p.GetName(), p.GetReportCount(), lastReport))

		for _, r := range reports {
			f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), p.GetName())
			f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), r.GetCreatedAt().Format("02.01.2006"))
			f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), r.GetReportName())
			f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), r.GetReportMessageLink())
			row++
		}
	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return err
	}

	if err := a.bot.SendMessages([]bot_dto.Message{
		{Chat: msg.ChatID, Text: sb.String()},
	}); err != nil {
		return err
	}

	return a.bot.SendDocument(msg.ChatID, "отчёт_фг.xlsx", buf.Bytes(), "")
}

func (a *Action) doDebug(msg dto.Message) error {
	names := []string{
		"Иванов Иван Иванович",
		"Петрова Анна Сергеевна",
		"Сидоров Алексей Дмитриевич",
		"Козлова Мария Александровна",
		"Новиков Дмитрий Павлович",
		"Морозова Елена Викторовна",
		"Волков Артём Николаевич",
		"Лебедева Ольга Андреевна",
		"Соколов Максим Романович",
		"Егорова Татьяна Игоревна",
		"Кузнецов Сергей Олегович",
	}

	reportNames := []string{
		"Утренняя тренировка",
		"Вечерняя тренировка",
		"Кардио",
		"Силовая",
		"Растяжка",
		"Бег",
		"Йога",
		"Плавание",
		"Велосипед",
		"Функциональная",
	}

	f := excelize.NewFile()
	sheetName := "Отчёты"
	f.SetSheetName("Sheet1", sheetName)
	f.SetCellValue(sheetName, "A1", "Участник")
	f.SetCellValue(sheetName, "B1", "Дата")
	f.SetCellValue(sheetName, "C1", "Название отчёта")
	f.SetCellValue(sheetName, "D1", "Ссылка")
	row := 2

	var sb strings.Builder
	sb.WriteString("Отчёт по ФГ (DEBUG)\n\n")

	for i, name := range names {
		p := participant.New(
			participant.WithTgID(int64(100+i)),
			participant.WithName(name),
			participant.WithReportCount(21),
		)

		var lastDate time.Time
		for j := 0; j < 21; j++ {
			createdAt := time.Now().AddDate(0, 0, -rand.Intn(30))
			if lastDate.IsZero() || createdAt.After(lastDate) {
				lastDate = createdAt
			}
			r := report.New(
				report.WithReportName(reportNames[rand.Intn(len(reportNames))]),
				report.WithReportMessageLink(fmt.Sprintf("https://t.me/group/msg/%d", rand.Intn(10000))),
				report.WithCreatedAt(createdAt),
			)

			f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), p.GetName())
			f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), r.GetCreatedAt().Format("02.01.2006"))
			f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), r.GetReportName())
			f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), r.GetReportMessageLink())
			row++
		}

		sb.WriteString(fmt.Sprintf("%d. %s — %d отч., крайний: %s\n", i+1, p.GetName(), p.GetReportCount(), lastDate.Format("02.01.2006")))
	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return err
	}

	if err := a.bot.SendMessages([]bot_dto.Message{
		{Chat: msg.ChatID, Text: sb.String()},
	}); err != nil {
		return err
	}

	return a.bot.SendDocument(msg.ChatID, "отчёт_фг.xlsx", buf.Bytes(), "")
}
