package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/SobolevTim/finance_bot/internal/pkg/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository структура репозитория
type Repository struct {
	DB     *pgxpool.Pool // Пул соединений к БД
	Logger *slog.Logger  // Логгер для модуля БД
}

// NewUserRepository создает новый экземпляр репозитория
//
// ctx - контекст
// cfg - конфигурация
// logger - логгер
//
// Возвращает новый экземпляр репозитория или ошибку
func NewUserRepository(ctx context.Context, cfg config.Config, logger *slog.Logger) (*Repository, error) {
	logger.Info("Подключение к БД...", "URL", cfg.DB.URL)
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

	return &Repository{
		DB:     pool,
		Logger: logger,
	}, nil
}

// Close закрывает соединение с БД
func (s *Repository) Close() {
	s.DB.Close()
	s.Logger.Info("Подключение к БД закрыто")
}

// Ping проверяет соединение с БД
func (s *Repository) Ping(ctx context.Context) error {
	s.Logger.Debug("Проверка соединения с БД...")
	return s.DB.Ping(ctx)
}
