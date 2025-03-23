package telegram

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/SobolevTim/finance_bot/internal/service"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

const daysPerPage = 7

// handleInline обрабатывает инлайн-кнопки
func (b *Bot) inlinehandlers(update telego.Update) {
	callbackData := update.CallbackQuery.Data
	chatID := update.CallbackQuery.From.ID
	b.logger.Debug("Получено инлайн-событие", "callbackData", callbackData, "tgID", chatID)

	if strings.HasPrefix(callbackData, "expenses_page_") {
		var page int
		_, err := fmt.Sscanf(callbackData, "expenses_page_%d", &page)
		if err != nil {
			return
		}
		b.handleExpenseCommand(chatID, page)
	} else if strings.HasPrefix(callbackData, "add_") {
		// Обработка inline-кнопок для записи расхода.
		b.HandleAddExpenseCallback(chatID, callbackData)
	}
}

// handleExpenseCommand обрабатывает команду /expense
func (b *Bot) handleExpenseCommand(chatID int64, page int) {
	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7 // Воскресенье -> 7
	}

	// Вычисляем начало недели (понедельник) и сдвигаем на нужное количество недель
	startOfWeek := now.AddDate(0, 0, -weekday+1-(page*7)).Truncate(24 * time.Hour)
	endOfWeek := startOfWeek.AddDate(0, 0, 6)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	expenses, err := b.Service.GetExpenses(ctx, chatID, startOfWeek, endOfWeek)
	if err != nil {
		b.SendErrorMessage(chatID, "Ошибка при получении данных о расходах")
		return
	}

	totalSum, avgExpense, maxExpense, maxDate := calculateSummary(expenses)

	message := fmt.Sprintf(
		"*Обзор расходов за неделю (%s - %s):*\n"+
			"- Общая сумма: %.2f₽\n"+
			"- Средний расход: %.2f₽\n"+
			"- Макс. расход: %.2f₽ (%s)\n\n",
		startOfWeek.Format("02.01.2006"), endOfWeek.Format("02.01.2006"),
		totalSum, avgExpense, maxExpense, maxDate.Format("02.01.2006"),
	)

	for _, exp := range expenses {
		message += fmt.Sprintf("📅 %s: %.2f₽\n", exp.Date.Format("02.01.2006"), exp.Amount)
	}

	// Инлайн-кнопки для переключения недель
	inlineKeyboard := tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⬅ Пред. неделя").WithCallbackData(fmt.Sprintf("expenses_page_%d", page+1)),
			tu.InlineKeyboardButton("След. неделя ➡").WithCallbackData(fmt.Sprintf("expenses_page_%d", page-1)),
		),
	)

	b.SendMessageWithKeyboard(chatID, message, inlineKeyboard)
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
func (b *Bot) HandleAddExpenseCallback(chatID int64, callbackData string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	entry, err := b.Service.GetExpenseStatus(ctx, chatID)
	if err != nil || entry == nil {
		b.logger.Error("Ошибка получения статуса записи расхода", "error", err)
		b.SendErrorMessage(chatID, "Произошла ошибка. Попробуйте еще раз")
		return
	}

	switch callbackData {
	case "add_date_today":
		entry.Date = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC)
		entry.Step = "amount"
		err := b.Service.SetExpenseStatus(ctx, chatID, entry)
		if err != nil {
			b.logger.Error("Ошибка установки статуса записи расхода", "error", err)
			b.SendErrorMessage(chatID, "Произошла ошибка. Попробуйте еще раз")
			return
		}
		b.sendAmountPrompt(chatID)
	case "add_date_custom":
		// Просим ввести дату текстом.
		entry.Step = "date_input"
		err := b.Service.SetExpenseStatus(ctx, chatID, entry)
		if err != nil {
			b.logger.Error("Ошибка установки статуса записи расхода", "error", err)
			b.SendErrorMessage(chatID, "Произошла ошибка. Попробуйте еще раз")
			return
		}
		b.sendTextPrompt(chatID, "Введите дату в формате ДД.ММ.ГГГГ (например, 23.03.2025):")
	default:
		// Обработка выбора категории. Ожидается формат "add_category_<Category>"
		if strings.HasPrefix(callbackData, "add_category_") {
			cat := strings.TrimPrefix(callbackData, "add_category_")
			entry.Category = cat
			entry.Step = "note"
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

			b.SendMessageWithKeyboard(chatID, "Хотите добавить примечание?", keyboard)
		} else if callbackData == "add_note" {
			entry.Step = "note_input"
			err := b.Service.SetExpenseStatus(ctx, chatID, entry)
			if err != nil {
				b.logger.Error("Ошибка установки статуса записи расхода", "error", err)
				b.SendErrorMessage(chatID, "Произошла ошибка. Попробуйте еще раз")
				return
			}
			b.sendTextPrompt(chatID, "Введите примечание:")
		} else if callbackData == "add_skip_note" {
			entry.Note = ""
			entry.Step = "confirm"
			err := b.Service.SetExpenseStatus(ctx, chatID, entry)
			if err != nil {
				b.logger.Error("Ошибка установки статуса записи расхода", "error", err)
				b.SendErrorMessage(chatID, "Произошла ошибка. Попробуйте еще раз")
				return
			}
			b.sendConfirmation(chatID, entry)
		} else if callbackData == "add_confirm" {
			// Подтверждение записи расхода
			if err := b.Service.AddExpense(ctx, chatID, entry.Amount, entry.Date, entry.Category, entry.Note); err != nil {
				b.SendErrorMessage(chatID, "Ошибка записи расхода.")
			} else {
				b.SendMessage(chatID, "✅ Расход записан!")
			}
			b.Service.DeleteStatus(ctx, chatID)
		} else if callbackData == "add_cancel" {
			b.SendErrorMessage(chatID, "Запись отменена.")
			b.Service.DeleteStatus(ctx, chatID)
		}
	}
	// Ответ на callback можно добавить при необходимости.
}
