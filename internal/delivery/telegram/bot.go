package telegram

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/SobolevTim/finance_bot/internal/service"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

type Bot struct {
	Client      *telego.Bot
	UserService *service.UserService
	logger      *slog.Logger
}

func NewBot(token string, userService *service.UserService, logger *slog.Logger, debug bool) (*Bot, error) {
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
		Client:      client,
		UserService: userService,
		logger:      logger,
	}, nil
}

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
		b.logger.Debug("update", "update", update)
		if update.Message != nil {
			if strings.HasPrefix(update.Message.Text, "/") {
				go b.handlersCmd(update)
			}
		}
	}
}

func (b *Bot) SendErrorMessage(id int64, text string) {
	msg := tu.Message(tu.ID(id), "❌ "+text)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	b.Client.SendMessage(ctx, msg)
	b.logger.Debug("Отправка сообщения", "message", msg.Text, "chatID", msg.ChatID)
}

func (b *Bot) SendMessage(id int64, text string) {
	msg := tu.Message(tu.ID(id), text)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	b.Client.SendMessage(ctx, msg)
	b.logger.Debug("Отправка сообщения", "message", msg.Text, "chatID", msg.ChatID)
}

func (b *Bot) SendMessageWitchKeyBoard(id int64, text string, keyboard telego.ReplyMarkup) {
	msg := tu.Message(tu.ID(id), text).WithReplyMarkup(keyboard)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	b.Client.SendMessage(ctx, msg)
	b.logger.Debug("Отправка сообщения с клавиаторой", "message", msg.Text, "chatID", msg.ChatID, "keyboard", keyboard)
}
