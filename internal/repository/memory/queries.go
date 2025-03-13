package memory

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
)

func (r *MemoryRepository) SetStatus(ctx context.Context, telegramID string, status string) error {
	r.logger.Debug("Установка статуса пользователя в Redis", "TelegramID", telegramID, "Status", status)
	err := r.rdb.Set(ctx, telegramID, status, 0).Err()
	if err != nil {
		r.logger.Error("Ошибка установки статуса пользователя в Redis", "TelegramID", telegramID, "Status", status, "error", err)
		return err
	}
	r.logger.Debug("Статус пользователя установлен", "TelegramID", telegramID, "Status", status)
	return nil
}

func (r *MemoryRepository) GetStatus(ctx context.Context, telegramID string) (string, error) {
	r.logger.Debug("Получение статуса пользователя из Redis", "TelegramID", telegramID)
	status, err := r.rdb.Get(ctx, telegramID).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			r.logger.Debug("Статус пользователя не найден", "TelegramID", telegramID)
			return "", nil
		}
		r.logger.Error("Ошибка получения статуса пользователя из Redis", "TelegramID", telegramID, "error", err)
		return "", err
	}
	r.logger.Debug("Статус пользователя получен", "TelegramID", telegramID, "Status", status)
	return status, nil
}
