package main

import (
	"context"
	"log"
	"time"

	"github.com/SobolevTim/finance_bot/internal/delivery/telegram"
	"github.com/SobolevTim/finance_bot/internal/pkg/config"
	"github.com/SobolevTim/finance_bot/internal/pkg/logger"
	"github.com/SobolevTim/finance_bot/internal/repository/database"
	"github.com/SobolevTim/finance_bot/internal/repository/memory"
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
	memlogger := logger.GetLogger("memorydb")

	// Подключаем репозиторий
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	repo, err := database.NewUserRepository(ctx, *config, bdlogger)
	if err != nil {
		bdlogger.Error("ошибка при создании user repository", "error", err)
		return
	}
	defer repo.Close()

	// Подключаем репозиторий Статистики
	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	StatRepo, err := memory.NewMemoryRepository(ctx, *config, memlogger)
	if err != nil {
		memlogger.Error("ошибка при создании memory repository", "error", err)
		return
	}
	defer StatRepo.Close()
	defer StatRepo.Close()

	// Подключаем сервисы
	service := service.NewService(repo, repo, StatRepo, repo, repo)

	// Создаем бота
	bot, err := telegram.NewBot(config.TG.Token, service, tglogger, config.TG.Debug)
	if err != nil {
		tglogger.Error("ошибка создания бота", "error", err)
		return
	}

	// Запускаем бота
	bot.StartBot(config.TG.TypePolling)
	tglogger.Info("Бот запущен")

	// Ожидаем завершения работы
	// TODO добавить завершение работы
	select {}
}
