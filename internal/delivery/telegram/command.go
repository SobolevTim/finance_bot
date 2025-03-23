package telegram

import (
	"context"
	"fmt"
	"time"

	"github.com/SobolevTim/finance_bot/internal/service"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

// handlersCmd обработка команд
//
// Обработка команд:
// /start - регистрация пользователя
// /cancel - отмена операции
// /help - получение справки
// /setbudget - установка бюджета
func (b *Bot) handlersCmd(update telego.Update) {
	b.logger.Debug("Получена команда", "command", update.Message.Text, "tgID", update.Message.Chat.ID)
	switch update.Message.Text {
	case "/start":
		b.handlersStart(update)
	case "/cancel":
		b.handlersCancel(update)
	case "/help":
		// TODO: Добавить обработку команды help
	case "/setbudget":
		b.handlersSetBudget(update)
	case "/getbudget":
		b.handlersGetBudget(update)
	case "/expense":
		b.handleExpenseCommand(update.Message.Chat.ID, 0)
	case "/add":
		b.StartAddExpense(update.Message.Chat.ID)
	default:
		b.logger.Debug("Неизвестная команда", "command", update.Message.Text)
		b.SendMessage(update.Message.Chat.ID, "Неизвестная команда")
	}
}

// handlersStart обработка команды start
//
// При получении команды регистрирует пользователя в базе данных
// и отправляет сообщение с приветствием и бюджетом
func (b *Bot) handlersStart(update telego.Update) {
	b.logger.Debug("Обработка команды start", "tgID", update.Message.Chat.ID)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	user, err := b.Service.RegisterUser(ctx, update.Message.Chat.ID, update.Message.From.Username, update.Message.Chat.FirstName, update.Message.Chat.LastName)

	if err != nil {
		b.logger.Error("Ошибка регистрации пользователя", "error", err)
		b.SendErrorMessage(update.Message.Chat.ID, "Ошибка регистрации пользователя, попробуйте еще раз")
		return
	}

	// Получение бюджета пользователя
	budget, err := b.Service.GetCurrentBudget(ctx, user.ID)
	if err != nil {
		b.logger.Error("Ошибка получения бюджета", "error", err)
		b.SendErrorMessage(update.Message.Chat.ID, "Ошибка получения бюджета, попробуйте еще раз")
		return
	}

	if budget == nil {
		b.logger.Debug("Бюджет не найден", "userID", user.ID)
		b.SendMessage(update.Message.Chat.ID, "💰Бюджет на месяц еще не установлен!\nНапишите мне сумму, которые вы закладываете на месяц")
		err := b.Service.SetStatus(ctx, update.Message.Chat.ID, StatusBudget)
		if err != nil {
			b.logger.Error("Ошибка установки статуса", "error", err)
			b.SendErrorMessage(update.Message.Chat.ID, "Что-то пошло не так. Попробуйте еще раз чуть позже")
		}
		return
	}
	// Формирование сообщения
	text := fmt.Sprintf("Привет, %s!\nЯ бот для ведения бюджета.\nВаш бюджет на месяц %.2f", user.UserName, budget.Amount.InexactFloat64())

	// Отправка сообщения
	b.SendMessage(update.Message.Chat.ID, text)
}

// handlersCancel обработка команды cancel
//
// При получении команды отменяет текущую операцию
// и отправляет сообщение об отмене
func (b *Bot) handlersCancel(update telego.Update) {
	b.logger.Debug("Обработка команды cancel", "tgID", update.Message.Chat.ID)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := b.Service.SetStatus(ctx, update.Message.Chat.ID, "")
	if err != nil {
		b.logger.Error("Ошибка обновления статуса", "error", err)
		b.SendErrorMessage(update.Message.Chat.ID, "Произошла ошибка. Попробуйте еще раз")
		return
	}

	text := "Операция отменена"
	b.SendMessage(update.Message.Chat.ID, text)
}

// handlersSetBudget обработка команды установки бюджета
//
// При получении команды устанавливает статус пользователя в StatusBudget
// и отправляет сообщение с просьбой указать бюджет
func (b *Bot) handlersSetBudget(update telego.Update) {
	b.logger.Debug("Обработка команды setbudget", "tgID", update.Message.Chat.ID)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := b.Service.SetStatus(ctx, update.Message.Chat.ID, StatusBudget)
	if err != nil {
		b.logger.Error("Ошибка обновления статуса", "error", err)
		b.SendErrorMessage(update.Message.Chat.ID, "Произошла ошибка. Попробуйте еще раз")
		return
	}

	text := "Укажите ваш бюджет на месяц"
	b.SendMessage(update.Message.Chat.ID, text)
}

// handlersGetBudget обработка команды получения бюджета
//
// При получении команды отправляет сообщение с текущим бюджетом пользователя
// Если бюджет не установлен, отправляет сообщение об этом
func (b *Bot) handlersGetBudget(update telego.Update) {
	b.logger.Debug("Обработка команды getbudget", "tgID", update.Message.Chat.ID)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	user, err := b.Service.GetUserByTelegramID(ctx, update.Message.Chat.ID)
	if err != nil {
		b.logger.Error("Ошибка получения пользователя", "error", err)
		b.SendErrorMessage(update.Message.Chat.ID, "Произошла ошибка. Попробуйте еще раз")
		return
	}

	budget, err := b.Service.GetCurrentBudget(ctx, user.ID)
	if err != nil {
		b.logger.Error("Ошибка получения бюджета", "error", err)
		b.SendErrorMessage(update.Message.Chat.ID, "Произошла ошибка. Попробуйте еще раз")
		return
	}

	if budget == nil {
		b.SendMessage(update.Message.Chat.ID, "Бюджет на месяц еще не установлен")
		return
	}

	text := fmt.Sprintf("Ваш бюджет на месяц %.2f", budget.Amount.InexactFloat64())
	b.SendMessage(update.Message.Chat.ID, text)
}

// StartAddExpense инициирует процесс записи расхода.
func (b *Bot) StartAddExpense(chatID int64) {
	// Создаем новое состояние записи
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	b.Service.SetExpenseStatus(ctx, chatID, &service.ExpenseEntryDTO{
		Step: "date",
	})

	message := "Выберите дату для записи расхода:"
	keyboard := tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("Сегодня").WithCallbackData("add_date_today"),
			tu.InlineKeyboardButton("Указать дату").WithCallbackData("add_date_custom"),
		),
	)
	msg := tu.Message(tu.ID(chatID), message).WithReplyMarkup(keyboard)
	b.Client.SendMessage(ctx, msg)
}
