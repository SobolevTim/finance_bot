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
		b.SendMessage(update.Message.Chat.ID, "Доступные команды:\n/start - регистрация\n/cancel - отмена операции\n/setbudget - установка бюджета\n/getbudget - получение бюджета\n/expense - просмотр расходов\n/add - добавление расхода")
	case "/setbudget":
		b.handlersSetBudget(update)
	case "/getbudget":
		b.handlersGetBudget(update)
	case "/expense":
		b.handleExpenseCommand(update, 0)
	case "/month":
		b.handleMonthCommand(update.Message.Chat.ID)
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

	userName := update.Message.From.Username
	firstName := update.Message.Chat.FirstName
	lastName := update.Message.Chat.LastName

	if update.Message.Chat.ID < 0 { // для группы ID = -1234567890
		userName = update.Message.Chat.Title
		firstName = update.Message.Chat.Title
		lastName = update.Message.Chat.Title
	}

	user, err := b.Service.RegisterUser(ctx, update.Message.Chat.ID, userName, firstName, lastName)

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
	b.logger.Debug("Обработка команды StartAddExpense", "chatID", chatID)
	// Создаем новое состояние записи
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	b.Service.SetExpenseStatus(ctx, chatID, &service.ExpenseEntryDTO{
		Step: "date",
	})

	now := time.Now()
	year, month, _ := now.Date()

	// Get first day of month and total days
	firstDay := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	lastDay := firstDay.AddDate(0, 1, -1)

	// Create keyboard rows
	var rows [][]telego.InlineKeyboardButton

	// Add month/year and navigation row
	monthRow := []telego.InlineKeyboardButton{
		tu.InlineKeyboardButton("←").WithCallbackData(fmt.Sprintf("add_month_%d_%d", int(month)-1, year)),
		tu.InlineKeyboardButton(fmt.Sprintf("%s %d", month.String(), year)).WithCallbackData("ignore"),
		tu.InlineKeyboardButton("→").WithCallbackData(fmt.Sprintf("add_month_%d_%d", int(month)+1, year)),
	}
	rows = append(rows, monthRow)

	// Add weekday headers (Mon-Sun)
	weekDays := []string{"Пн", "Вт", "Ср", "Чт", "Пт", "Сб", "Вс"}
	weekDayRow := make([]telego.InlineKeyboardButton, 7)
	for i, day := range weekDays {
		weekDayRow[i] = tu.InlineKeyboardButton(day).WithCallbackData("ignore")
	}
	rows = append(rows, weekDayRow)

	// Calculate first weekday (0 = Monday)
	firstWeekDay := int(firstDay.Weekday())
	if firstWeekDay == 0 {
		firstWeekDay = 7
	}
	firstWeekDay--

	// Generate calendar grid
	currentDay := 1
	for week := 0; currentDay <= lastDay.Day(); week++ {
		var row []telego.InlineKeyboardButton
		for weekDay := 0; weekDay < 7; weekDay++ {
			if (week == 0 && weekDay < firstWeekDay) || currentDay > lastDay.Day() {
				row = append(row, tu.InlineKeyboardButton(" ").WithCallbackData("ignore"))
			} else {
				dayStr := fmt.Sprintf("%d", currentDay)
				if currentDay == now.Day() {
					dayStr = "✓" + dayStr
				}
				row = append(row, tu.InlineKeyboardButton(dayStr).WithCallbackData(fmt.Sprintf("add_date_%d_%d_%d", year, month, currentDay)))
				currentDay++
			}
		}
		rows = append(rows, row)
	}

	keyboard := tu.InlineKeyboard(rows...)
	msg := tu.Message(tu.ID(chatID), "Выберите дату:").WithReplyMarkup(keyboard)
	b.Client.SendMessage(ctx, msg)
}

// handleMonthCommand обрабатывает команду /month
//
// Отправляет статистику расходов за текущий месяц
func (b *Bot) handleMonthCommand(chatID int64) {
	b.logger.Debug("Обработка команды month", "tgID", chatID)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Получение статистики расходов за текущий месяц
	expenses, sumExp, err := b.Service.GetExpensesByMonth(ctx, chatID)
	if err != nil {
		b.logger.Error("Ошибка получения расходов за месяц", "error", err)
		b.SendErrorMessage(chatID, "Произошла ошибка. Попробуйте еще раз")
		return
	}

	user, err := b.Service.GetUserByTelegramID(ctx, chatID)
	if err != nil {
		b.logger.Error("Ошибка получения пользователя", "error", err)
		b.SendErrorMessage(chatID, "Произошла ошибка. Попробуйте еще раз")
		return
	}

	budget, err := b.Service.GetCurrentBudget(ctx, user.ID)
	if err != nil {
		b.logger.Error("Ошибка получения бюджета", "error", err)
		b.SendErrorMessage(chatID, "Произошла ошибка. Попробуйте еще раз")
		return
	}

	if budget == nil {
		b.SendMessage(chatID, "Бюджет на месяц еще не установлен")
		return
	}
	userBudget := budget.Amount.InexactFloat64()

	// Формирование сообщения
	text := fmt.Sprintf("🙈 Расходы за месяц:\n")
	text += fmt.Sprintf("Всего потрачено: %.2f\n", sumExp)
	text += fmt.Sprintf("Бюджет на месяц: %.2f\n", userBudget)
	text += fmt.Sprintf("Осталось: %.2f\n", userBudget-sumExp)
	daysLeft := time.Now().AddDate(0, 0, 1).Day() - time.Now().Day()
	text += fmt.Sprintf("Среднее на день осталось: %.2f\n", (userBudget-sumExp)/float64(daysLeft))
	startDate := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, -1)
	text += fmt.Sprintf("Изнаначальное среднее: %.2f\n", userBudget/float64(endDate.Day()))

	// Отправка сообщения
	b.SendMessage(chatID, text)

	// Отправка сообщения с детализацией расходов
	var message string
	for _, exp := range expenses {
		line := fmt.Sprintf("📅 %s: %.2f₽ - %s %s", exp.Date.Format("02.01.2006"), exp.Amount, exp.CategoryIcon, exp.Category)
		if exp.Description != "" {
			line += fmt.Sprintf(" - %s", exp.Description)
		}
		message += line + "\n"
	}
	b.SendMessage(chatID, message)
}
