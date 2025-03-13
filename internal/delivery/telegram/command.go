package telegram

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/mymmrac/telego"
)

func (b *Bot) handlersCmd(update telego.Update) {
	b.logger.Debug("Получена команда", "command", update.Message.Text, "tgID", update.Message.Chat.ID)
	switch update.Message.Text {
	case "/start":
		b.handlersStart(update)
	case "/cancel":
		b.handkersCancel(update)
	case "/help":
		// TODO: Добавить обработку команды help
	default:
		b.logger.Debug("Неизвестная команда", "command", update.Message.Text)
		b.SendMessage(update.Message.Chat.ID, "Неизвестная команда")
	}
}

func (b *Bot) handlersStart(update telego.Update) {
	b.logger.Debug("Обработка команды start", "tgID", update.Message.Chat.ID)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Регистрация пользователя в базе данных
	tgID := strconv.FormatInt(update.Message.Chat.ID, 10)
	user, budget, err := b.UserService.RegisterUser(ctx, tgID, update.Message.From.Username, update.Message.From.FirstName, update.Message.From.LastName)

	if err != nil {
		b.logger.Error("Ошибка регистрации пользователя", "error", err)
		b.SendErrorMessage(update.Message.Chat.ID, "Ошибка регистрации пользователя")
		return
	}

	if budget.Amount == 0 {
		b.logger.Debug("Бюджет не установлен", "tgID", update.Message.Chat.ID)
		text := "Для начала работы укажите ваш бюджет на месяц"
		b.SendMessage(update.Message.Chat.ID, text)

		b.StatusMemory.SetStatus(ctx, update.Message.Chat.ID, StatusBudget)
		return
	}

	// Формирование сообщения
	text := fmt.Sprintf("Привет, %s!\nЯ бот для ведения бюджета.\nВаш бюджет на месяц %d", user.UserName, budget.Amount/100)

	// Отправка сообщения
	b.SendMessage(update.Message.Chat.ID, text)
}

func (b *Bot) handkersCancel(update telego.Update) {
	b.logger.Debug("Обработка команды cancel", "tgID", update.Message.Chat.ID)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := b.StatusMemory.SetStatus(ctx, update.Message.Chat.ID, "")
	if err != nil {
		b.logger.Error("Ошибка обновления статуса", "error", err)
		b.SendErrorMessage(update.Message.Chat.ID, "Произошла ошибка. Попробуйте еще раз")
		return
	}

	text := "Операция отменена"
	b.SendMessage(update.Message.Chat.ID, text)
}
