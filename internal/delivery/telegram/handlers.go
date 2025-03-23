package telegram

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mymmrac/telego"
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

	// Получение статуса
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	status, err := b.Service.GetStatus(ctx, update.Message.Chat.ID)
	if err != nil {
		b.logger.Error("Ошибка получения статуса", "error", err)
		b.SendErrorMessage(update.Message.Chat.ID, "Произошла ошибка. Воспользуйтесь командой /start для начала работы")
		return
	}

	// Обработка статуса
	if status != "" {
		b.handlerStatus(status, update)
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
	b.logger.Debug("Запрос бюджета", "tgID", chatID, "amount", amount)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	budget, err := b.Service.UpdateBudgetByTgID(ctx, chatID, amount)
	if err != nil {
		b.logger.Error("Ошибка обновления бюджета", "error", err)
		b.SendErrorMessage(chatID, "Произошла ошибка. Попробуйте еще раз")
		return
	}
	if budget == nil {
		b.logger.Error("Ошибка обновления бюджета", "error", "budget is nil")
		b.SendErrorMessage(chatID, "Произошла ошибка. Попробуйте еще раз")
		return
	}

	text := fmt.Sprintf("Бюджет на месяц установлен: %.2f", budget.Amount.InexactFloat64())
	b.logger.Debug("Бюджет установлен", "tgID", chatID, "amount", budget.Amount.InexactFloat64())
	b.SendMessage(chatID, text)
}

func (b *Bot) handlersMessage(update telego.Update) {
	b.logger.Debug("Обработка общих сообщений", "tgID", update.Message.Chat.ID)
	text := update.Message.Text

	transaction := ParseInput(text)
	if transaction.Error != nil {
		b.logger.Error("Ошибка обработки сообщения", "error", transaction.Error)
		b.SendErrorMessage(update.Message.Chat.ID, "Видимо вы ввели что-то не так!\nНапоминаю что формат ввода: ДД.ММ СУММА ОПИСАНИЕ\nЛибо в кратком формате: СУММА ОПИСАНИЕ. Можно и без описания")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Запись транзакции в базу данных
	err := b.Service.CreateExpensByTelegramID(ctx, update.Message.Chat.ID, transaction.Result, transaction.Date, transaction.Description)
	if err != nil {
		b.logger.Error("Ошибка создания транзакции", "error", err)
		b.SendErrorMessage(update.Message.Chat.ID, "Произошла ошибка. Попробуйте еще раз")
		return
	}
}
