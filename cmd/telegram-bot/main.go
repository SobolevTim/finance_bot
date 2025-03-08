package main

import (
	"log"
	"os"

	"github.com/SobolevTim/finance_bot/internal/delivery/telegram"
	"github.com/SobolevTim/finance_bot/internal/pkg/config"
	"github.com/SobolevTim/finance_bot/internal/pkg/logger"
)

func main() {
	// Подключаем конфигурацию
	config, err := config.LoadConfig("internal/pkg/config")
	if err != nil {
		log.Fatalln("Failed to load config:", err)
	}

	// Подключаем логгер
	logger.InitLogger(config.App.Env)
	tglogger := logger.GetLogger("telegram")

	// Создаем бота
	bot, err := telegram.NewBot(config.TG.Token, tglogger, config.TG.Debug)
	if err != nil {
		tglogger.Error("failed to create bot", "error", err)
		os.Exit(1)
	}

	// Запускаем бота
	updates := bot.StartBot("longpolling")
	for update := range updates {
		tglogger.Debug("update", "update", update)
	}
}
