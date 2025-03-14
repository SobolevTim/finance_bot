package memory

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

// SetStatus устанавливает статус пользователя
//
// telegramID - идентификатор пользователя
// status - статус пользователя
//
// Возвращает ошибку при возникновении проблем с Redis
func (r *MemoryRepository) SetStatus(ctx context.Context, telegramID string, status string) error {
	r.logger.Debug("Установка статуса пользователя в Redis", "TelegramID", telegramID, "Status", status)
	now := time.Now()
	err := r.rdb.Set(ctx, telegramID, status, time.Hour*24).Err() // Статус хранится 24 часа
	if err != nil {
		r.logger.Error("Ошибка установки статуса пользователя в Redis", "TelegramID", telegramID, "Status", status, "error", err, "Duration", time.Since(now))
		return err
	}
	r.logger.Debug("Статус пользователя установлен", "TelegramID", telegramID, "Status", status, "Duration", time.Since(now))
	return nil
}

// GetStatus получает статус пользователя
//
// telegramID - идентификатор пользователя
//
// Возвращает статус пользователя и ошибку при возникновении проблем с Redis
func (r *MemoryRepository) GetStatus(ctx context.Context, telegramID string) (string, error) {
	r.logger.Debug("Получение статуса пользователя из Redis", "TelegramID", telegramID)
	now := time.Now()
	status, err := r.rdb.Get(ctx, telegramID).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			r.logger.Debug("Статус пользователя не найден", "TelegramID", telegramID, "Duration", time.Since(now))
			return "", nil
		}
		r.logger.Error("Ошибка получения статуса пользователя из Redis", "TelegramID", telegramID, "error", err, "Duration", time.Since(now))
		return "", err
	}
	r.logger.Debug("Статус пользователя получен", "TelegramID", telegramID, "Status", status, "Duration", time.Since(now))
	return status, nil
}
