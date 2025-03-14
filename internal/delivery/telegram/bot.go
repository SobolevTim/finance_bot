package telegram

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/SobolevTim/finance_bot/internal/service"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

// Bot структура бота
type Bot struct {
	Client       *telego.Bot           // Клиент телеграма
	UserService  *service.UserService  // Сервис для работы с пользователями
	StatusMemory *service.StatusMemory // Сервис для работы со статусами
	logger       *slog.Logger          // Логгер
}

// NewBot создает новый экземпляр бота
//
// token - токен бота
// userService - сервис для работы с пользователями
// statusMem - сервис для работы со статусами
// logger - логгер
// debug - режим отладки
//
// Возвращает новый экземпляр бота или ошибку
func NewBot(token string, userService *service.UserService, statusMem *service.StatusMemory, logger *slog.Logger, debug bool) (*Bot, error) {
	logger.Debug("Создание бота с токеном", "token", token)
	logger.Debug("Дебаг режим бота", "debug", debug)

	// Создаем бота
	client, err := telego.NewBot(token, telego.WithDefaultLogger(debug, true))
	if err != nil {
		return nil, err
	}
	logger.Debug("Бот создан")
	bot, err := client.GetMe(context.Background())
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении информации о боте: %w", err)
	}
	logger.Info("Авторизация бота", "bot", bot.Username, "id", bot.ID, "firstName", bot.FirstName, "lastName", bot.LastName)
	return &Bot{
		Client:       client,
		UserService:  userService,
		StatusMemory: statusMem,
		logger:       logger,
	}, nil
}

// StartBot запускает бота
//
// polingType - тип работы бота
func (b *Bot) StartBot(polingType string) {
	if polingType == "longpolling" {
		b.logger.Debug("Запуск бота", "polingType", "longpolling")
		b.StartPooling()
	}
	if polingType == "webhook" {
		b.logger.Debug("Запуск бота", "polingType", "webhook")
		// TODO webhook
	}

}

// StartPooling запускает бота с использованием longpolling
func (b *Bot) StartPooling() {
	updates, err := b.Client.UpdatesViaLongPolling(
		// TODO добавить контекст для завершения работы
		context.Background(),

		&telego.GetUpdatesParams{
			Offset:  0,
			Timeout: 10,
		},

		telego.WithLongPollingUpdateInterval(time.Second*0),
		telego.WithLongPollingRetryTimeout(time.Second*1),
		telego.WithLongPollingBuffer(100),
	)
	if err != nil {
		b.logger.Error("ошибка при получении обновлений", "error", err)
	}

	for update := range updates {
		b.logger.Debug("Получено обновление", "update", update)
		if update.Message != nil {
			b.handlers(update)
		}
	}
}

// SendErrorMessage отправляет сообщение об ошибке
//
// id - идентификатор чата
// text - текст сообщения
func (b *Bot) SendErrorMessage(id int64, text string) {
	msg := tu.Message(tu.ID(id), "❌ "+text)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	b.Client.SendMessage(ctx, msg)
	b.logger.Debug("Отправка сообщения", "message", msg.Text, "chatID", msg.ChatID)
}

// SendMessage отправляет сообщение
//
// id - идентификатор чата
// text - текст сообщения
func (b *Bot) SendMessage(id int64, text string) {
	msg := tu.Message(tu.ID(id), text)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	b.Client.SendMessage(ctx, msg)
	b.logger.Debug("Отправка сообщения", "message", msg.Text, "chatID", msg.ChatID)
}
