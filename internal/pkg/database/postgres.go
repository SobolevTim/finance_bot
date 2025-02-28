package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/SobolevTim/finance_bot/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStorage struct {
	DB     *pgxpool.Pool // Пул соединений к БД
	Logger *slog.Logger  // Логгер для модуля БД
}

func NewPostgresStorage(ctx context.Context, cfg config.Config, logger *slog.Logger) (*PostgresStorage, error) {
	config, err := pgxpool.ParseConfig(cfg.DB.URL)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга DSN: %w", err)
	}

	config.MaxConns = cfg.DB.MaxConns           // Максимальное количество соединений
	config.MinConns = cfg.DB.IdleConns          // Минимальное количество соединений
	config.HealthCheckPeriod = 30 * time.Second // Период проверки соединения с БД

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к БД: %w", err)
	}

	// Проверяем соединение
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("БД недоступна: %w", err)
	}
	logger.Info("Подключение к БД установлено")

	return &PostgresStorage{
		DB:     pool,
		Logger: logger,
	}, nil
}

func (s *PostgresStorage) Close() {
	s.DB.Close()
	s.Logger.Info("Подключение к БД закрыто")
}

func (s *PostgresStorage) Ping(ctx context.Context) error {
	return s.DB.Ping(ctx)
}
