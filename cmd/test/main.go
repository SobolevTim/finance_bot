package main

import (
	"log"
	"log/slog"

	"github.com/SobolevTim/finance_bot/config"
	"github.com/SobolevTim/finance_bot/internal/pkg/logger"
)

func main() {
	// Подключаем конфигурацию
	config, err := config.LoadConfig("config")
	if err != nil {
		log.Fatalln("Failed to load config:", err)
	}

	// Подключаем логгер
	logger.InitLogger(config.App.Env)

	slog.Info("Test no module logger")

	// Получаем логгер для модуля
	testLogger := logger.GetLogger("test")
	testLogger.Info("Test info logger")
	testLogger.Debug("Test debug logger")
	testLogger.Error("Test error logger")
	testLogger.Warn("Test warn logger")
}
