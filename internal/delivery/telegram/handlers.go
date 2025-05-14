package telegram

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/SobolevTim/finance_bot/internal/pkg/calc"
	"github.com/SobolevTim/finance_bot/internal/service"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

// handlers обработка сообщений
//
// Обработка команд;
// Получение статуса;
// Обработка статуса;
// Обработка сообщения;
func (b *Bot) handlers(update telego.Update) {
	b.logger.Debug("Получено сообщение", "message", update.Message.Text, "tgID", update.Message.Chat.ID)

	// Обработка команд
	if strings.HasPrefix(update.Message.Text, "/") {
		b.handlersCmd(update)
		return
	}
	chatID := update.Message.Chat.ID

	// Получение статуса
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	status, err := b.Service.GetStatus(ctx, chatID)
	if err != nil {
		b.logger.Error("Ошибка получения статуса", "error", err)
		b.SendErrorMessage(chatID, "Произошла ошибка. Попробуйте еще раз")
		return
	}
	// Обработка статуса
	if status != "" {
		b.handlerStatus(status, update)
		return
	}

	// Получение статуса записи расхода
	statusExpense, err := b.Service.GetExpenseStatus(ctx, chatID)
	if err != nil {
		b.logger.Error("Ошибка получения статуса записи расхода", "error", err)
		b.SendErrorMessage(chatID, "Произошла ошибка. Попробуйте еще раз")
		return
	}
	// Обработка статуса записи расхода
	if statusExpense != nil {
		b.logger.Debug("Статус записи расхода получен", "tgID", chatID, "status", statusExpense.Step)
		b.HandleAddExpenseText(chatID, update.Message.Text, statusExpense)
		return
	}

	// Обработка сообщения
	b.handlersMessage(update)

}

// handlerStatus обработка статуса
//
// Обработка статуса "budget" - установка бюджета
func (b *Bot) handlerStatus(status string, update telego.Update) {
	b.logger.Debug("Обработка статуса", "status", status, "tgID", update.Message.Chat.ID)
	switch status {
	case StatusBudget:
		b.requestBudget(update)
	default:
		b.logger.Debug("Неизвестный статус", "status", status)
		b.SendErrorMessage(update.Message.Chat.ID, "Произошла ошибка. Воспользуйтесь командами:\n/start для начала работы\n/help для получения справки")
	}
}

// requestBudget запрос бюджета
//
// Запрос бюджета у пользователя и обновление в базе данных
// Отправка сообщения о результате
func (b *Bot) requestBudget(update telego.Update) {
	chatID := update.Message.Chat.ID
	amount := update.Message.Text
	b.logger.Debug("Запрос бюджета requestBudget", "tgID", chatID, "amount", amount)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	budget, err := b.Service.UpdateBudgetByTgID(ctx, chatID, amount)
	if err != nil {
		b.logger.Error("Ошибка обновления бюджета requestBudget", "error", err)
		b.SendErrorMessage(chatID, "Произошла ошибка. Попробуйте еще раз")
		return
	}
	if budget == nil {
		b.logger.Error("Ошибка обновления бюджета requestBudget", "error", "budget is nil")
		b.SendErrorMessage(chatID, "Произошла ошибка. Попробуйте еще раз")
		return
	}

	// Установка статуса
	err = b.Service.SetStatus(ctx, chatID, "")
	if err != nil {
		b.logger.Error("Ошибка установки статуса", "error", err)
		b.SendErrorMessage(chatID, "Произошла ошибка. Попробуйте еще раз")
		return
	}

	text := fmt.Sprintf("Бюджет на месяц установлен: %.2f", budget.Amount.InexactFloat64())
	b.logger.Debug("Бюджет установлен requestBudget", "tgID", chatID, "amount", budget.Amount.InexactFloat64())
	b.SendMessage(chatID, text)
}

// HandleAddExpenseText обрабатывает текстовые сообщения, поступающие на разных шагах диалога.
func (b *Bot) HandleAddExpenseText(chatID int64, text string, entry *service.ExpenseEntryDTO) {
	b.logger.Debug("Обработка текстового сообщения в HandleAddExpenseText", "tgID", chatID, "text", text, "entry", entry)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	switch entry.Step {
	case "date_input":
		t, err := time.Parse("02.01.2006", text)
		if err != nil {
			b.SendErrorMessage(chatID, "Неверный формат даты. Попробуйте еще раз.")
			return
		}
		entry.Date = t
		entry.Step = "amount"
		err = b.Service.SetExpenseStatus(ctx, chatID, entry)
		if err != nil {
			b.logger.Error("Ошибка установки статуса записи расхода", "error", err)
			b.SendErrorMessage(chatID, "Произошла ошибка. Попробуйте еще раз")
			return
		}
		b.sendAmountPrompt(chatID)
	case "amount":
		amount, err := calc.Calculate(text)
		if err != nil {
			b.SendErrorMessage(chatID, "Ошибка в вычислении суммы. Попробуйте еще раз.")
			return
		}
		entry.Amount = amount
		entry.Step = "category"
		err = b.Service.SetExpenseStatus(ctx, chatID, entry)
		if err != nil {
			b.logger.Error("Ошибка установки статуса записи расхода", "error", err)
			b.SendErrorMessage(chatID, "Произошла ошибка. Попробуйте еще раз")
			return
		}
		// Показываем кнопки для выбора категории.

		defaultCategory, err := b.Service.GetDefaultCategories(ctx)
		if err != nil {
			b.logger.Error("Ошибка получения категорий", "error", err)
			b.SendErrorMessage(chatID, "Произошла ошибка. Попробуйте еще раз.")
			return
		}

		// Создаем кнопки для категорий
		keyboards := tu.InlineKeyboard()
		for _, cat := range defaultCategory {
			keyboards.InlineKeyboard = append(keyboards.InlineKeyboard, tu.InlineKeyboardRow(
				tu.InlineKeyboardButton(cat.Icon+" "+cat.Name).WithCallbackData("add_category_"+cat.Name),
			))
		}

		b.SendMessageWithKeyboard(chatID, "Выберите категорию расхода:", keyboards)
	case "note_input":
		entry.Note = text
		entry.Step = "confirm"
		err := b.Service.SetExpenseStatus(ctx, chatID, entry)
		if err != nil {
			b.logger.Error("Ошибка установки статуса записи расхода", "error", err)
			b.SendErrorMessage(chatID, "Произошла ошибка. Попробуйте еще раз")
			return
		}
		b.sendConfirmation(chatID, entry)
	}
}

func (b *Bot) handlersMessage(update telego.Update) {
	b.logger.Debug("Обработка общих сообщений", "tgID", update.Message.Chat.ID)
	// TODO обработка сообщений
	// b.SendErrorMessage(update.Message.Chat.ID, "Произошла ошибка. Воспользуйтесь командами:\n/start для начала работы\n/help для получения справки")
}
