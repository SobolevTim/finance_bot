package telegram

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func (b *Bot) handlersCmd(update telego.Update) {
	b.logger.Debug("Получена команда", "command", update.Message.Text, "tgID", update.Message.Chat.ID)
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
	b.logger.Debug("Обработка команды start", "tgID", update.Message.Chat.ID)

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
		b.requestBudget(update.Message.Chat.ID)
		return
	}

	// Формирование сообщения
	text := fmt.Sprintf("Привет, %s!\nЯ бот для ведения бюджета.\nВаш бюджет на месяц %d", user.UserName, budget.Amount/100)

	// Отправка сообщения
	b.SendMessage(update.Message.Chat.ID, text)
}

func (b *Bot) requestBudget(chatID int64) {
	b.logger.Debug("Запрос бюджета", "tgID", chatID)

	// Формирование сообщения
	text := "Для начала работы укажите ваш бюджет на месяц"

	// Формирование клавиатуры
	inlineKeyboard := tu.InlineKeyboard(
		tu.InlineKeyboardRow( // Row 1
			tu.InlineKeyboardButton("Прибавить на 1000"). // Column 1
									WithCallbackData("callback_1"),
			tu.InlineKeyboardButton("Убавить на 1000"). // Column 2
									WithCallbackData("callback_2"),
		),
	)
	b.SendMessageWitchKeyBoard(chatID, text, inlineKeyboard)
}
