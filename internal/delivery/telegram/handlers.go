package telegram

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/mymmrac/telego"
)

func (b *Bot) handlers(update telego.Update) {
	b.logger.Debug("Получено сообщение", "message", update.Message.Text, "tgID", update.Message.Chat.ID)

	// Обработка команд
	if strings.HasPrefix(update.Message.Text, "/") {
		b.handlersCmd(update)
	}

	// Получение статуса
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	status, err := b.StatusMemory.GetStatus(ctx, update.Message.Chat.ID)
	if err != nil {
		b.logger.Error("Ошибка получения статуса", "error", err)
		b.SendErrorMessage(update.Message.Chat.ID, "Произошла ошибка. Воспользуйтесь командой /start для начала работы")
		return
	}

	// Обработка статуса
	if status != "" {
		b.handlerStatus(status, update)
	}

}

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

func (b *Bot) requestBudget(update telego.Update) {
	chatID := update.Message.Chat.ID
	b.logger.Debug("Запрос бюджета", "tgID", chatID)
	amount, err := strconv.ParseInt(update.Message.Text, 10, 64)
	if err != nil {
		b.logger.Error("Ошибка преобразования бюджета", "error", err)
		b.SendErrorMessage(chatID, "Ошибка ввода. Введите число")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = b.UserService.UpdateBudget(ctx, chatID, amount*100)
	if err != nil {
		b.logger.Error("Ошибка обновлении бюджета", "error", err)
		b.SendErrorMessage(chatID, "Произошла ошибка. Попробуйте еще раз")
		return
	}

	err = b.StatusMemory.SetStatus(ctx, chatID, "")
	if err != nil {
		b.logger.Error("Ошибка обновления статуса", "error", err)
		b.SendErrorMessage(chatID, "Произошла ошибка. Попробуйте еще раз")
		return
	}

	text := fmt.Sprintf("Бюджет на месяц установлен: %d", amount)
	b.SendMessage(chatID, text)
}
