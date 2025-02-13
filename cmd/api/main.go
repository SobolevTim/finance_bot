package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/SobolevTim/finance_bot/config"
	"github.com/SobolevTim/finance_bot/internal/bot"
	"github.com/SobolevTim/finance_bot/internal/database"
)

func main() {
	// Загружаем конфигурацию
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// Подключаемся к базе данных
	service, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("ERROR: Ошибка подключения к базе данных: %v", err)
	}
	defer service.DB.Close()
	defer log.Println("WARMING: Соедениние БД закрыто")

	// Выполняем миграцию
	migrationFile := "internal/database/migration.sql"
	if err := database.ApplyMigration(service, migrationFile); err != nil {
		log.Fatalf("ERROR: Ошибка применения миграции: %v", err)
	}

	// Здесь запуск веб-сервера или обработка запросов
	log.Println("INFO: Database запущена и готова к работе.")

	// Создаем и запускаем бота
	telegramBot, err := bot.NewBot(cfg.BotToken)
	if err != nil {
		log.Fatalf("ERROR: Ошибка при создании бота: %v", err)
	}

	// Запускаем бота в отдельной горутине
	go func() {
		if err := telegramBot.Start(service); err != nil {
			log.Printf("ERROR: Ошибка при запуске бота: %v", err)
		}
	}()
	log.Println("INFO: Telegram бот запущен")

	defer telegramBot.Client.StopLongPolling()
	defer log.Println("WARMING: Остановка LongPolling telegram")
	// Канал для обработки сигналов
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	// Ожидаем сигнала завершения
	<-signalChannel
}
