package telegram

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/SobolevTim/finance_bot/internal/service"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

const daysPerPage = 7

// handleInline обрабатывает инлайн-кнопки
func (b *Bot) inlinehandlers(update telego.Update) {
	cb := update.CallbackQuery

	msg := cb.Message
	data := cb.Data
	chatID := msg.GetChat().ID

	b.logger.Debug("Получено инлайн-событие", "callbackData", data, "tgID", chatID, "messageID", msg.GetMessageID())

	switch {
	case strings.HasPrefix(data, "expenses_page_"):
		var page int
		if _, err := fmt.Sscanf(data, "expenses_page_%d", &page); err == nil {
			b.handleExpenseCommand(update, page)
		} else {
			b.logger.Warn("Не удалось распарсить номер страницы", "data", data)
		}
	case strings.HasPrefix(data, "add_"):
		b.HandleAddExpenseCallback(msg.GetMessageID(), chatID, data)
	default:
		_ = b.Client.AnswerCallbackQuery(context.Background(), &telego.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			Text:            "Обновление данных...",
			ShowAlert:       false,
		})
		return
	}
}

// handleExpenseCommand обрабатывает команду /expense
func (b *Bot) handleExpenseCommand(update telego.Update, page int) {
	var chatID int64
	var messageID int

	if update.CallbackQuery != nil {
		chatID = update.CallbackQuery.Message.GetChat().ID
		messageID = update.CallbackQuery.Message.GetMessageID()
	} else if update.Message != nil {
		chatID = update.Message.Chat.ID
	} else {
		b.logger.Warn("handleExpenseCommand: update не содержит Message или CallbackQuery")
		return
	}

	message, keyboard, err := b.buildExpenseSummary(chatID, page)
	if err != nil {
		b.SendErrorMessage(chatID, "Ошибка при получении данных о расходах")
		return
	}

	if update.CallbackQuery != nil {
		// Ответ на callback
		_ = b.Client.AnswerCallbackQuery(context.Background(), &telego.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			Text:            "Обновление данных...",
			ShowAlert:       false,
		})

		// Обновляем сообщение
		b.EditMessageWithKeyboard(chatID, messageID, message, keyboard)
	} else {
		// Отправляем новое сообщение
		b.SendMessageWithKeyboard(chatID, message, keyboard)
	}
}

// buildExpenseSummary строит сводку расходов за неделю
func (b *Bot) buildExpenseSummary(chatID int64, page int) (string, *telego.InlineKeyboardMarkup, error) {
	if page < -10 {
		page = -10
	}

	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	startOfWeek := now.AddDate(0, 0, -weekday+1-(page*7)).Truncate(24 * time.Hour)
	endOfWeek := startOfWeek.AddDate(0, 0, 6)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	expenses, err := b.Service.GetExpenses(ctx, chatID, startOfWeek, endOfWeek)
	if err != nil {
		return "", nil, err
	}

	total, avg, max, maxDate := calculateSummary(expenses)
	var message string
	if len(expenses) == 0 || (total == 0 && avg == 0 && max == 0 && maxDate.IsZero()) {
		message = fmt.Sprintf("*Обзор расходов за неделю (%s – %s):*\n"+
			"— Нет расходов за неделю\n\n",
			startOfWeek.Format("02.01.2006"), endOfWeek.Format("02.01.2006"),
		)
	} else {
		message = fmt.Sprintf(
			"*Обзор расходов за неделю (%s – %s):*\n"+
				"— Общая сумма: %.2f₽\n"+
				"— Средний расход: %.2f₽\n"+
				"— Макс. расход: %.2f₽ (%s)\n\n",
			startOfWeek.Format("02.01.2006"), endOfWeek.Format("02.01.2006"),
			total, avg, max, maxDate.Format("02.01.2006"),
		)
	}

	for _, e := range expenses {
		message += fmt.Sprintf("📅 %s: %.2f₽\n", e.Date.Format("02.01.2006"), e.Amount)
	}

	keyboard := tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⬅ Пред. неделя").WithCallbackData(fmt.Sprintf("expenses_page_%d", page+1)),
			tu.InlineKeyboardButton("След. неделя ➡").WithCallbackData(fmt.Sprintf("expenses_page_%d", page-1)),
		),
	)

	return message, keyboard, nil
}

// calculateSummary вычисляет сводку расходов
func calculateSummary(expenses []*service.ExpenseDTO) (total float64, avg float64, max float64, maxDate time.Time) {
	if len(expenses) == 0 {
		return 0, 0, 0, time.Time{}
	}
	for _, exp := range expenses {
		total += exp.Amount
		if exp.Amount > max {
			max = exp.Amount
			maxDate = exp.Date
		}
	}
	avg = total / float64(len(expenses))
	return total, avg, max, maxDate
}

// HandleAddExpenseCallback обрабатывает inline-кнопки для записи расхода.
func (b *Bot) HandleAddExpenseCallback(id int, chatID int64, callbackData string) {
	b.logger.Debug("Обработка инлайн-кнопки для записи расхода", "callbackData", callbackData, "chatID", chatID)
	// Обработка выбора даты в инлайн календаре
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	entry, err := b.Service.GetExpenseStatus(ctx, chatID)
	if err != nil || entry == nil {
		b.logger.Error("Ошибка получения статуса записи расхода", "error", err)
		b.SendErrorMessage(chatID, "Произошла ошибка. Попробуйте еще раз")
		return
	}

	switch {
	// Кнопка с датой
	case strings.HasPrefix(callbackData, "add_date_"):
		callbackData = strings.TrimPrefix(callbackData, "add_date_")

		// Регулярное выражение для извлечения компонентов даты
		re := regexp.MustCompile(`^(\d{4})_(\d{1,2})_(\d{1,2})$`)
		matches := re.FindStringSubmatch(callbackData)
		if len(matches) != 4 {
			b.logger.Error("Некорректный формат даты",
				"input", callbackData,
				"expected_format", "YYYY_MM_DD")
			b.SendErrorMessage(chatID, "❌ Неверный формат даты. Используйте ГГГГ_ММ_ДД")
			return
		}

		// Парсинг компонентов
		year, err := strconv.Atoi(matches[1])
		if err != nil || year < 2000 || year > 2100 {
			b.logger.Error("Некорректный год", "year", matches[1])
			b.SendErrorMessage(chatID, "❌ Недопустимый год (2000-2100)")
			return
		}

		month, err := strconv.Atoi(matches[2])
		if err != nil || month < 1 || month > 12 {
			b.logger.Error("Некорректный месяц", "month", matches[2])
			b.SendErrorMessage(chatID, "❌ Недопустимый месяц (1-12)")
			return
		}

		day, err := strconv.Atoi(matches[3])
		if err != nil || day < 1 || day > 31 {
			b.logger.Error("Некорректный день", "day", matches[3])
			b.SendErrorMessage(chatID, "❌ Недопустимый день (1-31)")
			return
		}

		// Создание объекта времени с валидацией
		date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
		if date.Day() != day || date.Month() != time.Month(month) || date.Year() != year {
			b.logger.Error("Некорректная дата",
				"input", callbackData,
				"parsed", date.Format("2006-01-02"))
			b.SendErrorMessage(chatID, "❌ Недопустимая дата (проверьте календарь)")
			return
		}

		if date.After(time.Now().AddDate(0, 1, 0)) { // +1 месяца
			b.SendErrorMessage(chatID, "❌ Нельзя добавлять расходы более чем на 1 месяц вперед")
			return
		}

		b.AddExpenseDate(id, chatID, date)
	//case strings.HasPrefix(callbackData, "add_amount_"):
	case strings.HasPrefix(callbackData, "add_cancel"):
		err := b.Service.DeleteStatus(ctx, chatID)
		if err != nil {
			b.logger.Error("Ошибка удаления статус", "error", err)
			b.SendErrorMessage(chatID, "Произошла ошибка. Попробуйте еще раз.")
		}
		b.SendMessage(chatID, "Ввод траты отменен")
	// Обработка выбора категории. Ожидается формат "add_category_<Category>"
	case strings.HasPrefix(callbackData, "add_category_"):
		cat := strings.TrimPrefix(callbackData, "add_category_")
		entry.Category = cat
		entry.Step = "note"
		entry.ID = id
		err := b.Service.SetExpenseStatus(ctx, chatID, entry)
		if err != nil {
			b.logger.Error("Ошибка установки статуса записи расхода", "error", err)
			b.SendErrorMessage(chatID, "Произошла ошибка. Попробуйте еще раз")
			return
		}
		// Предлагаем добавить примечание или пропустить
		keyboard := tu.InlineKeyboard(
			tu.InlineKeyboardRow(
				tu.InlineKeyboardButton("Добавить примечание").WithCallbackData("add_note"),
				tu.InlineKeyboardButton("Пропустить").WithCallbackData("add_skip_note"),
			),
		)

		b.EditMessageWithKeyboard(chatID, id, "Хотите добавить примечание?", keyboard)
	// Кнопка добавления примечания
	case callbackData == "add_note":
		entry.Step = "note_input"
		entry.ID = id
		err := b.Service.SetExpenseStatus(ctx, chatID, entry)
		if err != nil {
			b.logger.Error("Ошибка установки статуса записи расхода", "error", err)
			b.SendErrorMessage(chatID, "Произошла ошибка. Попробуйте еще раз")
			return
		}
		keyboard := tu.InlineKeyboard(
			tu.InlineKeyboardRow(
				tu.InlineKeyboardButton("⛔️ Отменить ввод расхода").WithCallbackData("add_cancel"),
			),
		)
		b.SendMessageWithKeyboard(chatID, "Введите примечание:", keyboard)
	// Кнопка пропуска записи примечания
	case callbackData == "add_skip_note":
		entry.Note = ""
		entry.Step = "confirm"
		entry.ID = id
		err := b.Service.SetExpenseStatus(ctx, chatID, entry)
		if err != nil {
			b.logger.Error("Ошибка установки статуса записи расхода", "error", err)
			b.SendErrorMessage(chatID, "Произошла ошибка. Попробуйте еще раз")
			return
		}
		b.editConfirmation(chatID, entry.ID, entry)
	// Подтверждение записи расхода
	case callbackData == "add_confirm":
		if err := b.Service.AddExpense(ctx, chatID, entry.Amount, entry.Date, entry.Category, entry.Note); err != nil {
			b.SendErrorMessage(chatID, "Ошибка записи расхода.")
		} else {
			b.SendMessage(chatID, "✅ Расход записан!")
		}
		b.Service.DeleteStatus(ctx, chatID)
	}
}

// addDateToExpense обрабатывает выбор даты в инлайн-календаре.
// В данном случае мы просто обновляем дату в статусе записи расхода.
func (b *Bot) AddExpenseDate(id int, chatID int64, date time.Time) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	entry, err := b.Service.GetExpenseStatus(ctx, chatID)
	if err != nil || entry == nil {
		b.logger.Error("Ошибка получения статуса записи расхода", "error", err)
		b.SendErrorMessage(chatID, "Произошла ошибка. Попробуйте еще раз")
		return
	}

	entry.Step = "amount"
	entry.Date = date
	entry.ID = id
	err = b.Service.SetExpenseStatus(ctx, chatID, entry)
	if err != nil {
		b.logger.Error("Ошибка установки статуса записи расхода", "error", err)
		b.SendErrorMessage(chatID, "Произошла ошибка. Попробуйте еще раз")
		return
	}
	text := "Введите сумму расхода (можно использовать математическое выражение, например, 150+20)"
	keyboard := tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⛔️ Отменить ввод расхода").WithCallbackData("add_cancel"),
		),
	)

	b.EditMessageWithKeyboard(chatID, id, text, keyboard)
	b.logger.Debug("Выбрана дата для записи расхода", "chatID", chatID, "date", date.Format("02.01.2006"))
}
