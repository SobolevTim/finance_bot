package bot

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/SobolevTim/finance_bot/internal/database"
	"github.com/SobolevTim/finance_bot/pkg/calculator"
	"github.com/mymmrac/telego"
)

const (
	userMontlyBudget = "awaiting_total_amount" // Статус для ожидания ввода суммы на месяц для пользователя
	userNotify       = "awaiting_status"       // Статус ожидания ввода подписки на уведомления
)

var (
	// 1. Число с плавающей точкой
	floatRegex = regexp.MustCompile(`^[-+]?[0-9]*[.,]?[0-9]+([ \t]*[+-/*][ \t]*[-+]?[0-9]*[.,]?[0-9]+)*$`)
	// 2. Несколько чисел с плавающей точкой (разделенные + и -)
	multipleFloatsRegex = regexp.MustCompile(`^([-+]?[0-9]*\.?[0-9]+\s*[-+]\s*)*[-+]?[0-9]*\.?[0-9]+$`)
)

// Хранилище для состояния пользователей
var userState = make(map[int64]string)

// handleMessage обрабатывает входящие сообщения от пользователя и направляет их
// в соответствующие функции-обработчики в зависимости от содержимого сообщения и состояния пользователя.
//
// Параметры:
//   - msg: объект telego.Message, содержащий текст сообщения от пользователя.
//   - service: экземпляр database.Service, обеспечивающий доступ к операциям с базой данных.
//
// Логика обработки:
//   - Если текст начинается с "/", сообщение обрабатывается как команда и передается в handleCommand.
//   - Если текст начинается с "Дата", сообщение обрабатывается как ввод данных с определенной датой и
//     направляется в handleDataInsertAmount.
//   - Если у пользователя установлен статус userMontlyBudget, сообщение обрабатывается как ввод бюджета
//     и передается в handleAmountInput.
//   - Если у пользователя установлен статус userNotify, сообщение обрабатывается как ввод уведомления
//     и передается в handleNotifyInput.
//   - Если текст содержит одно или несколько чисел с плавающей точкой, он обрабатывается как расходы за текущий день
//     с использованием функций handleToDayAmount.
//   - По умолчанию отправляется сообщение об ошибке с инструкциями для корректного ввода данных, если текст не соответствует
//     ожидаемому формату.
func (b *Bot) handleMessage(msg *telego.Message, service *database.Service) {
	if msg.LeftChatMember != nil {
		log.Println("INFO: member hast left", msg.LeftChatMember.Username, msg.LeftChatMember.ID)
	} else if msg.NewChatMembers != nil {
		b.handleNewChat(msg)
	} else if strings.HasPrefix(msg.Text, "/") {
		b.handleCommand(msg, service)
	} else if strings.HasPrefix(msg.Text, "Дата") {
		b.handleDataInsertAmount(msg, service)
	} else if strings.HasPrefix(msg.Text, "Сколько") {
		b.handleDataGetAmount(msg, service)
	} else if state, ok := userState[msg.Chat.ID]; ok && state == userMontlyBudget {
		b.handleAmountInput(msg, service)
	} else if state, ok := userState[msg.Chat.ID]; ok && state == userNotify {
		b.handleNotifyInput(msg, service)
	} else {
		msgText := strings.TrimSpace(msg.Text)
		switch {
		case floatRegex.MatchString(msgText):
			// Обработка одного числа с плавающей точкой
			b.handleToDayAmount(msg, service)
		case multipleFloatsRegex.MatchString(msgText):
			// Обработка нескольких чисел с плавающей точкой
			b.handleToDayAmount(msg, service)
		default:
			b.sendMessage(msg.Chat.ID, "неизвестный формат сообщения. Используйте /help для получения информации о доступных командах\nДля записы трат за сегодняший день - просто напишите сумму трат.\nДля запиши трат на конкретную дату - напишите: Дата ДАТА(в формате ДД.ММ.ГГ) СУММА ТРАТ, например: Дата 01.01.24 1000")
		}
	}
}

// handleCommand обрабатывает текстовые команды, отправленные пользователем в виде сообщений,
// и вызывает соответствующие функции-обработчики для каждой команды.
//
// Параметры:
//   - msg: объект telego.Message, содержащий информацию о сообщении от пользователя.
//   - service: экземпляр database.Service, обеспечивающий доступ к операциям с базой данных.
//
// Поддерживаемые команды:
//   - /start: вызывает handleStart для приветствия нового пользователя и инициализации его данных.
//   - /help: вызывает handleHelp для отображения списка доступных команд.
//   - /cancel: вызывает handleCancel для отмены текущего действия пользователя.
//   - /update_montly_sum: вызывает handleUpdateMonthlySum для обновления месячного бюджета.
//   - /montly: вызывает handleMonthlyBudget для показа текущего бюджета пользователя.
//   - /expence: вызывает handleExpence для вывода информации о тратах за текущий месяц и среднюю сумму, оставщуюся на оставщиеся дни.
//   - /notify: вызывает handleNotify для управления напоминаниями пользователя.
//   - По умолчанию: отправляет сообщение о неизвестной команде, если текст не совпадает с известными командами.
func (b *Bot) handleCommand(msg *telego.Message, service *database.Service) {
	switch msg.Text {
	case "/start":
		b.handleStart(msg, service)
	case "/help":
		b.handleHelp(msg)
	case "/cancel":
		b.handleCancel(msg)
	case "/update_montly_sum":
		b.handleUpdateMonthlySum(msg)
	case "/montly":
		b.handleMonthlyBudget(msg, service)
	case "/expence":
		b.handleExpence(msg, service)
	case "/notify":
		b.handleNotify(msg, service)
	default:
		b.sendMessage(msg.Chat.ID, "Неизвестная команда.")
	}
}

// handleStart обрабатывает команду /start от пользователя, отправляет приветственное сообщение и
// сохраняет нового пользователя в базе данных. В случае ошибки при сохранении пользователя
// отправляет сообщение об ошибке и записывает её в лог.
//
// Параметры:
//   - msg: объект telego.Message, содержащий информацию о сообщении, отправленном пользователем.
//   - service: экземпляр database.Service, обеспечивающий доступ к операциям с базой данных.
func (b *Bot) handleStart(msg *telego.Message, service *database.Service) {
	var message string
	userState[msg.Chat.ID] = userMontlyBudget
	var newUser database.Users
	if msg.Chat.ID < 0 {
		message = fmt.Sprintf("Привет, %s!\nЯ бот для расчета финансов. Для начала введите ваш бюджет на месяц - сумму которую вы расчитываете потратить за месяц.\nДля получения списка команд, используйте /help", msg.Chat.Title)
		newUser = database.Users{
			TelegramID: msg.Chat.ID,
			Username:   msg.Chat.Title,
		}
	} else {
		message = fmt.Sprintf("Привет, %s!\nЯ бот для расчета финансов. Для начала введите ваш бюджет на месяц - сумму которую вы расчитываете потратить за месяц.\nДля получения списка команд, используйте /help", msg.From.FirstName)
		newUser = database.Users{
			TelegramID: msg.Chat.ID,
			Username:   msg.From.FirstName,
		}
	}
	if err := service.InsertStartUsers(newUser); err != nil {
		b.sendMessage(msg.Chat.ID, "Произошла ошибка при создании пользователя.")
		log.Printf("ERROR: %v", err)
		return
	}
	b.sendMessage(msg.Chat.ID, message)
}

// handleHelp отправляет пользователю справочное сообщение с перечнем доступных команд и их описанием.
//
// Параметры:
//   - msg: объект telego.Message, содержащий информацию о пользователе, отправившем запрос.
//
// Справочное сообщение включает:
//   - информацию о том, как заполнять основные траты.
//   - /cancel: отмена состояния ввода данных.
//   - /montly: отображение информации о текущем месячном бюджете.
//   - /update_montly_sum: обновление бюджета на месяц.
//   - /expence: показ оставшегося бюджета на текущий день.
//   - /notify: управление ежедневными напоминаниями о внесении трат.
func (b *Bot) handleHelp(msg *telego.Message) {
	message := "Если хотите записать траты на сегодня - просто введите сумму,\nлибо траты с операторами (\"+\", \"-\", \"*\", \"/\")\nНапример, 1000	.33 + 33 * 5 - 300.33\n\nЧтобы заполнить траты за прошедние дни\nвведи информацию в формате: Дата 01.02.24 СУММА ТРАТ\n\nЧтобы посмотреть, какую сумму вы вписали за конкретный день\nВведите Сколько 01.02.24 (Дата в формате ДД.ММ.ГГ)\n\nМои основные команды:\n/cancel - отмена состояния ввода данных.\n/montly - информация о текущем бюджете на месяц (сумма которые вы планируете тратить).\n/update_montly_sum - обновление бюджета на месяц.\n/expence - информация о тратах в этом месяце и оставшейся сумме\n/notify - начать ежедневные напоминания о внесении трат, либо отменить их"
	b.sendMessage(msg.Chat.ID, message)
}

func (b *Bot) handleCancel(msg *telego.Message) {
	delete(userState, msg.Chat.ID)
	message := "Вы отменили ввод. Для просмотра списка команд используйте /help."
	b.sendMessage(msg.Chat.ID, message)
}

func (b *Bot) handleExpence(msg *telego.Message, service *database.Service) {
	expence, err := service.GetAverageMontlyExpenses(database.Users{TelegramID: msg.Chat.ID})
	if err != nil {
		b.sendMessage(msg.Chat.ID, "Произошла ошибка при получении суммы трат за месяц.")
		log.Printf("ERROR: %v", err)
		return
	}
	user, err := service.GetMontlyBudget(database.Users{TelegramID: msg.Chat.ID})
	if err != nil {
		b.sendMessage(msg.Chat.ID, "Произошла ошибка при получении вашего бюджета на месяц.")
		log.Printf("ERROR: %v", err)
		return
	}
	now := time.Now()
	nextMonth := now.AddDate(0, 1, -now.Day()+1)
	daysRemaining := int(nextMonth.Sub(now).Hours() / 24)
	//(бюджет на месяц - сумма трат) / количество оставшихся дней.
	averageCount := (user.MonthlyBudget - expence.Amount) / daysRemaining
	message := fmt.Sprintf("В этом месяце вы уже потратили: %.2f\nНа оставшиеся %d дней средняя сумма: %.2f",
		float64(expence.Amount)/100,
		daysRemaining,
		float64(averageCount)/100)
	b.sendMessage(msg.Chat.ID, message)
}

func (b *Bot) handleNotify(msg *telego.Message, service *database.Service) {
	user, err := service.GetUserNotify(database.Users{TelegramID: msg.Chat.ID})
	if err != nil {
		b.sendMessage(msg.Chat.ID, "Произошла ошибка при получении статуса подписки на ежедневные сообщения.")
		log.Printf("ERROR: %v", err)
		return
	}
	var message string
	userState[msg.Chat.ID] = userNotify
	switch {
	case user.Notify:
		message = "Вы подписаны на получение ежедневных сообщений - напоминаний. Если хотите изменить - напишите Подписка\nЕсли хотите оставить все как есть используйте команду /cancel"
	case !user.Notify:
		message = "Вы не подписаны на получение ежедневных сообщений - напоминаний. Если хотите изменить - напишите Подписка\nЕсли хотите оставить все как есть используйте команду /cancel"
	}
	b.sendMessage(msg.Chat.ID, message)
}

func (b *Bot) handleNotifyInput(msg *telego.Message, service *database.Service) {
	text := strings.ToLower(msg.Text)
	if text != "подписка" {
		message := "Введите: \"Подписка\".\nЛибо используйте /cancel - для отмены ввода"
		b.sendMessage(msg.Chat.ID, message)
		return
	}
	user, err := service.GetUserNotify(database.Users{TelegramID: msg.Chat.ID})
	if err != nil {
		b.sendMessage(msg.Chat.ID, "Произошла ошибка при получении статуса подписки на ежедневные сообщения.")
		log.Printf("ERROR: %v", err)
		return
	}
	err = service.UpdateUserNotify(database.Users{TelegramID: msg.Chat.ID, Notify: !user.Notify})
	if err != nil {
		b.sendMessage(msg.Chat.ID, "Произошла ошибка при изменении статуса подписки на ежедневные сообщения.")
		log.Printf("ERROR: %v", err)
		return
	}
	delete(userState, msg.Chat.ID)
	message := ("Ваша подписка получение ежедневных сообщений - изменена.")
	b.sendMessage(msg.Chat.ID, message)
}

func (b *Bot) handleUpdateMonthlySum(msg *telego.Message) {
	message := "Введите новую сумму трат на месяц.\nЕсли вы передумали - используйте команду /cancel"
	userState[msg.Chat.ID] = userMontlyBudget
	b.sendMessage(msg.Chat.ID, message)
}

func (b *Bot) handleMonthlyBudget(msg *telego.Message, service *database.Service) {
	user, err := service.GetMontlyBudget(database.Users{TelegramID: msg.Chat.ID})
	if err != nil {
		b.sendMessage(msg.Chat.ID, "Произошла ошибка при получении бюджета.")
		log.Printf("ERROR: %v", err)
		return
	}
	message := fmt.Sprintf("%d.%d - Ваша сумма трат на месяц.", user.MonthlyBudget/100, user.MonthlyBudget%100)
	b.sendMessage(msg.Chat.ID, message)
}

func (b *Bot) handleAmountInput(msg *telego.Message, service *database.Service) {
	amount, err := strconv.ParseFloat(msg.Text, 64)
	if err != nil {
		message := "Введите корректную сумму (например, 123.45)."
		b.sendMessage(msg.Chat.ID, message)
		return
	}
	user := database.Users{
		TelegramID:    msg.Chat.ID,
		MonthlyBudget: int(amount * 100),
	}
	if err := service.UpdateMontlyBudget(user); err != nil {
		b.sendMessage(msg.Chat.ID, "Произошла ошибка при обновлении суммы.")
		log.Printf("ERROR: %v", err)
		return
	}
	delete(userState, msg.Chat.ID)
	message := fmt.Sprintf("%.2f - Ваша сумма трат на месяц обновлена.", amount)
	b.sendMessage(msg.Chat.ID, message)
}

func (b *Bot) handleToDayAmount(msg *telego.Message, service *database.Service) {
	amount, err := calculator.Calc(msg.Text)
	if err != nil {
		message := "Введите корректную сумму (например: 123.45 или 123.45 + 67 - 89)."
		b.sendMessage(msg.Chat.ID, message)
		return
	}
	user := database.Users{
		TelegramID: msg.Chat.ID,
	}
	expence := database.Expenses{
		Amount:      amount,
		ExpenseDate: time.Now(),
	}
	todayExp, err := service.UpdateDayExpense(user, expence)
	if err != nil {
		b.sendMessage(msg.Chat.ID, "Произошла ошибка при записи суммы трат.")
		log.Printf("ERROR: %v", err)
		return
	}
	var message string
	switch {
	case todayExp.Amount != amount:
		if amount > 0 {
			message = fmt.Sprintf("Добавил %.2f к Вашим тратам на сегодня.\nИтоговая сумма трат за сегодня: %.2f", float64(amount)/100, float64(todayExp.Amount)/100)
		} else {
			message = fmt.Sprintf("Вычел %.2f из Ваших трат за сегодня.\nИтоговая сумма трат за сегодня: %.2f", float64(amount)/100, float64(todayExp.Amount)/100)
		}
	case todayExp.Amount == amount:
		message = fmt.Sprintf("Записал %.2f к Вашим тратам на сегодня.", float64(amount)/100)
	}
	b.sendMessage(msg.Chat.ID, message)
}

func (b *Bot) handleDataInsertAmount(msg *telego.Message, service *database.Service) {
	text := strings.Split(msg.Text, " ")
	if len(text) < 3 {
		b.sendMessage(msg.Chat.ID, "Кажется Вы забыли что-то ввести!🥲\nНапоминаю, что формат ввода даных должен быть такой:\nДата 01.02.2024 ТРАТЫ\nТраты можно вводить как одним числом, так и несколько чисел с мат. операторами (сложение +; вычитание -; умножение *; деление /)")
		return
	}
	var date time.Time
	var err error

	switch {
	case len(text[1]) == 5: // Формат "01.02"
		nowYear := time.Now().Year()
		text[1] += "." + strconv.Itoa(nowYear)
		date, err = time.Parse("02.01.2006", text[1])
	case len(text[1]) == 10: // Формат "01.02.2024"
		date, err = time.Parse("02.01.2006", text[1])
	case len(text[1]) == 8: // Формат "01.02.24"
		date, err = time.Parse("02.01.06", text[1])
	default:
		err = fmt.Errorf("неизвестный формат даты")
	}
	if err != nil {
		b.sendMessage(msg.Chat.ID, "Произошла ошибка при чтении даты. Используй любой из удобных форматов:\nДД.ММ - например, 01.02; ДД.ММ.ГГ - 01.02.24; ДД.ММ.ГГГГ - 01.02.2024")
		log.Printf("ERROR: %v", err)
		return
	}

	nums := strings.Join(text[2:], " ")
	amount, err := calculator.Calc(nums)
	if err != nil {
		message := "Введите корректную сумму (например: 123.45 или 123.45 + 67 - 89)."
		b.sendMessage(msg.Chat.ID, message)
		return
	}

	user := database.Users{
		TelegramID: msg.Chat.ID,
	}
	expence := database.Expenses{
		Amount:      amount,
		ExpenseDate: date,
	}

	todayExp, err := service.UpdateDayExpense(user, expence)
	if err != nil {
		b.sendMessage(msg.Chat.ID, "Произошла ошибка при записи суммы трат.")
		log.Printf("ERROR: %v", err)
		return
	}
	var message string
	switch {
	case todayExp.Amount != amount:
		if amount > 0 {
			message = fmt.Sprintf("Добавил %.2f к Вашим тратам на дату: %v.\nИтоговая сумма трат составляет: %.2f", float64(amount)/100, date.Format("02.01.2006"), float64(todayExp.Amount)/100)
		} else {
			message = fmt.Sprintf("Вычел %.2f из Ваших трат за дату: %v.\nИтоговая сумма трат составляет: %.2f", float64(amount)/100, date.Format("02.01.2006"), float64(todayExp.Amount)/100)
		}
	case todayExp.Amount == amount:
		message = fmt.Sprintf("Записал %.2f к Вашим тратам на дату: %v.", float64(amount)/100, todayExp.ExpenseDate.Format("02.01.2006"))
	}
	b.sendMessage(msg.Chat.ID, message)
}

func (b *Bot) handleDataGetAmount(msg *telego.Message, service *database.Service) {
	text := strings.Split(msg.Text, " ")
	if len(text) < 2 {
		message := fmt.Sprintf("Кажется Вы ввели что-то не так!🥲\nНапоминаю, что формат ввода даных должен быть такой:\nДата 01.02.2024 (ДД.ММ.ГГГГ / ДД.ММ.ГГ)\n либо сокращенный формат: Дата 01.02 (ДД.ММ)\nВаше сообщение было такое: %s", msg.Text)
		b.sendMessage(msg.Chat.ID, message)
		return
	}
	var date time.Time
	var err error

	switch {
	case len(text[1]) == 5: // Формат "01.02"
		nowYear := time.Now().Year()
		text[1] += "." + strconv.Itoa(nowYear)
		date, err = time.Parse("02.01.2006", text[1])
	case len(text[1]) == 10: // Формат "01.02.2024"
		date, err = time.Parse("02.01.2006", text[1])
	case len(text[1]) == 8: // Формат "01.02.24"
		date, err = time.Parse("02.01.06", text[1])
	default:
		err = fmt.Errorf("неизвестный формат даты")
	}
	if err != nil {
		b.sendMessage(msg.Chat.ID, "Произошла ошибка при чтении даты. Используй любой из удобных форматов:\nДД.ММ - например, 01.02; ДД.ММ.ГГ - 01.02.24; ДД.ММ.ГГГГ - 01.02.2024")
		log.Printf("ERROR: %v", err)
		return
	}

	user := database.Users{
		TelegramID: msg.Chat.ID,
	}
	expense := database.Expenses{
		ExpenseDate: date,
	}
	dateExpense, err := service.GetExpenseFromDate(user, expense)
	if err != nil {
		b.sendMessage(msg.Chat.ID, "Произошла ошибка при получении суммы трат.")
		log.Printf("ERROR: %v", err)
		return
	}
	var message string
	switch {
	case dateExpense.Amount == 0:
		message = fmt.Sprintf("За %s записей о тратах нет! Если вы что-то тратили в этот день - запишите траты\nДля это воспользуйтесь конструкцией: Дата 01.02.2024 ТРАТЫ\nНапример, 01.02.2024 1000 + 500", expense.ExpenseDate.Format("02.01.2006"))
	case dateExpense.Amount < 0:
		message = fmt.Sprintf("Ого! За %s не траты а заработок! Записана сумма: %.2f", expense.ExpenseDate.Format("02.01.2006"), float64(dateExpense.Amount)/100)
	default:
		message = fmt.Sprintf("Ваши траты на дату %s составляют: %.2f", expense.ExpenseDate.Format("02.01.2006"), float64(dateExpense.Amount)/100)
	}
	b.sendMessage(msg.Chat.ID, message)
}

func (b *Bot) handleNewChat(msg *telego.Message) {
	if msg.NewChatMembers != nil {
		for _, user := range msg.NewChatMembers {
			if user.IsBot {
				message := "Если вы только что добавили меня в чат - используйте /start для начала моей работы"
				b.sendMessage(msg.Chat.ID, message)
			}
		}
	}
}
