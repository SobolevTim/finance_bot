package telegram

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/mymmrac/telego"
)

func (b *Bot) handlersCmd(update telego.Update) {
	b.logger.Debug("Обработка команды", "command", update.Message.Text)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	switch update.Message.Text {
	case "/start":
		b.handlersStart(ctx, update)
	case "/help":
		// TODO: Добавить обработку команды help
	default:
		b.logger.Debug("Неизвестная команда", "command", update.Message.Text)
		b.SendMessage(update.Message.Chat.ID, "Неизвестная команда")
	}
}

func (b *Bot) handlersStart(ctx context.Context, update telego.Update) {
	b.logger.Debug("Обработка команды start")

	// Регистрация пользователя в базе данных
	tgID := strconv.FormatInt(update.Message.Chat.ID, 10)
	_, err := b.UserService.RegisterUser(ctx, tgID, update.Message.From.Username, update.Message.From.FirstName, update.Message.From.LastName)
	if err != nil {
		b.logger.Error("Ошибка регистрации пользователя", "error", err)
		b.SendErrorMessage(update.Message.Chat.ID, "Ошибка регистрации пользователя")
		return
	}

	// Формирование сообщения
	text := fmt.Sprintf("Привет, %s!\nЯ бот для ведения бюджета.", update.Message.From.Username)

	// Отправка сообщения
	b.SendMessage(update.Message.Chat.ID, text)
}
