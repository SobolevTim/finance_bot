package memory

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/SobolevTim/finance_bot/internal/pkg/config"
	"github.com/redis/go-redis/v9"
)

type MemoryRepository struct {
	rdb    *redis.Client
	logger *slog.Logger
}

func NewMemoryRepository(ctx context.Context, cfg config.Config, logger *slog.Logger) (*MemoryRepository, error) {
	logger.Info("Подключение к Redis...", "Addr", cfg.Redis.Addr)
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,     // Адрес сервера Redis
		Password: cfg.Redis.Password, // Пароль, если установлен
		DB:       cfg.Redis.DB,       // Номер базы данных (0 по умолчанию)
		PoolSize: cfg.Redis.PoolSize, // Размер пула соединений
	})

	p, err := rdb.Ping(ctx).Result()
	logger.Debug("Проверка соединения с Redis...", "Ping", p)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к Redis: %w", err)
	}
	logger.Info("Подключение к Redis установлено")
	return &MemoryRepository{
		rdb:    rdb,
		logger: logger,
	}, nil
}

func (r *MemoryRepository) Close() {
	r.logger.Debug("Закрытие подключения к Redis...")
	err := r.rdb.Close()
	if err != nil {
		r.logger.Error("Ошибка закрытия подключения к Redis", "error", err)
		return
	}
	r.logger.Info("Подключение к Redis закрыто")
}

func (r *MemoryRepository) Ping(ctx context.Context) error {
	r.logger.Debug("Проверка соединения с Redis...")
	_, err := r.rdb.Ping(ctx).Result()
	if err != nil {
		r.logger.Error("Ошибка проверки соединения с Redis", "error", err)
	}
	return err
}
