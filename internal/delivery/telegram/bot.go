package telegram

import (
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	Client *tgbotapi.BotAPI
	logger *slog.Logger
}

func NewBot(token string, logger *slog.Logger, debug bool) (*Bot, error) {
	logger.Debug("Создание бота с токеном", "token", token)
	client, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	logger.Debug("Дебаг режим бота", "debug", debug)
	client.Debug = debug
	logger.Debug("Бот создан")
	logger.Info("Авторизация бота", "bot", client.Self.UserName)
	return &Bot{
		Client: client,
		logger: logger,
	}, nil
}

func (b *Bot) StartBot(polingType string) tgbotapi.UpdatesChannel {
	b.logger.Debug("Запуск бота", "polingType", polingType)
	var updates tgbotapi.UpdatesChannel
	if polingType == "longpolling" {
		updateConfig := tgbotapi.NewUpdate(0)
		updateConfig.Timeout = 30
		updates = b.Client.GetUpdatesChan(updateConfig)
	}
	if polingType == "webhook" {
		// webhook
	}
	return updates
}
