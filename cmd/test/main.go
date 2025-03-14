package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/SobolevTim/finance_bot/internal/pkg/config"
	"github.com/SobolevTim/finance_bot/internal/pkg/logger"
	"github.com/SobolevTim/finance_bot/internal/repository/database"
)

func main() {
	// Подключаем конфигурацию
	config, err := config.LoadConfig("internal/pkg/config")
	if err != nil {
		log.Fatalln("Failed to load config:", err)
	}

	// Подключаем логгер
	logger.InitLogger(config.App.Env)
	dblogger := logger.GetLogger("database")
	httplogger := logger.GetLogger("http")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	// Подключаемся к БД
	db, err := database.NewUserRepository(ctx, *config, dblogger)
	if err != nil {
		dblogger.Error("failed to connect to DB", "error", err)
		return
	}
	defer db.Close()

	router := http.NewServeMux()
	router.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
		defer cancel()
		if err := db.Ping(ctx); err != nil {
			dblogger.Error("Database ping failed", "error", err)
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	httplogger.Info("Starting server", "port", 8080)
	if err := server.ListenAndServe(); err != nil {
		httplogger.Error("server failed", "error", err)
	}
}
