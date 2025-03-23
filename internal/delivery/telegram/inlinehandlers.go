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

	if strings.HasPrefix(callbackData, "expenses_page_") {
		var page int
		_, err := fmt.Sscanf(callbackData, "expenses_page_%d", &page)
		if err != nil {
			return
		}
		b.handleExpenseCommand(chatID, page)
	}
}

// HandleExpenseCommand обрабатывает команду /expense
func (b *Bot) handleExpenseCommand(chatID int64, page int) {
	// Получаем расходы за текущий месяц
	year, month, _ := time.Now().Date()
	startDate := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, -1)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	expenses, err := b.Service.GetExpenses(ctx, chatID, startDate, endDate)
	if err != nil {
		b.SendErrorMessage(chatID, "Произошла ошибка при получении расходов")
		return
	}

	// Формируем сводку расходов
	totalSum, avgExpense, maxExpense, maxDate := calculateSummary(expenses)

	message := fmt.Sprintf(
		"*Обзор расходов за месяц:*\n"+
			"- Общая сумма: %.2f₽\n"+
			"- Средний расход: %.2f₽\n"+
			"- Макс. расход: %.2f₽ (%s)\n\n",
		totalSum, avgExpense, maxExpense, maxDate.Format("02.01.2006"),
	)

	// Разделяем расходы по неделям
	pagedExpenses := paginateExpenses(expenses, page)

	for _, exp := range pagedExpenses {
		message += fmt.Sprintf("📅 %s: %.2f₽\n", exp.Date.Format("02.01.2006"), exp.Amount)
	}

	// Добавляем инлайн-кнопки для листания
	inlineKeyboard := tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⬅ Пред. неделя").WithCallbackData(fmt.Sprintf("expenses_page_%d", page-1)),
			tu.InlineKeyboardButton("След. неделя ➡").WithCallbackData(fmt.Sprintf("expenses_page_%d", page+1)),
		),
	)

	msg := tu.Message(tu.ID(chatID), message).WithParseMode(telego.ModeMarkdown).WithReplyMarkup(inlineKeyboard)

	ctx, cancel = context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	b.Client.SendMessage(ctx, msg)
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

// paginateExpenses выбирает 7-дневный отрезок из списка расходов
func paginateExpenses(expenses []*service.ExpenseDTO, page int) []*service.ExpenseDTO {
	start := page * daysPerPage
	end := start + daysPerPage
	if start >= len(expenses) {
		return nil
	}
	if end > len(expenses) {
		end = len(expenses)
	}
	return expenses[start:end]
}
