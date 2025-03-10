package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/SobolevTim/finance_bot/internal/delivery/telegram"
	"github.com/SobolevTim/finance_bot/internal/pkg/config"
	"github.com/SobolevTim/finance_bot/internal/pkg/logger"
	"github.com/SobolevTim/finance_bot/internal/repository/database"
	"github.com/SobolevTim/finance_bot/internal/service"
)

func main() {
	// Подключаем конфигурацию
	config, err := config.LoadConfig("internal/pkg/config")
	if err != nil {
		log.Fatalln("ошибка при загрузке конфигурации:", err)
	}

	// Подключаем логгеры
	logger.InitLogger(config.App.Env)
	tglogger := logger.GetLogger("telegram")
	bdlogger := logger.GetLogger("database")

	// Подключаем репозитории
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	userRepo, err := database.NewUserRepository(ctx, *config, bdlogger)
	if err != nil {
		bdlogger.Error("ошибка при создании user repository", "error", err)
		os.Exit(1)
	}

	// Подключаем сервисы
	userService := service.NewUserService(userRepo)

	// Создаем бота
	bot, err := telegram.NewBot(config.TG.Token, userService, tglogger, config.TG.Debug)
	if err != nil {
		tglogger.Error("ошибка создания бота", "error", err)
		os.Exit(1)
	}

	// Запускаем бота
	bot.StartBot(config.TG.TypePolling)
	tglogger.Info("Бот запущен")

	// Ожидаем завершения работы
	// TODO добавить завершение работы
	select {}
}
