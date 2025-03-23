package memory

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/SobolevTim/finance_bot/internal/domain/status"
	"github.com/redis/go-redis/v9"
)

// SetStatus устанавливает статус пользователя
//
// telegramID - идентификатор пользователя
// status - статус пользователя
//
// Возвращает ошибку при возникновении проблем с Redis
func (r *MemoryRepository) SetStatus(ctx context.Context, telegramID, status string) error {
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

func (r *MemoryRepository) SetExpenseStatus(ctx context.Context, ChatID string, expenseEntry *status.ExpenseEntry) error {
	key := "expense:" + ChatID
	r.logger.Debug("Установка статуса записи расхода в Redis", "TelegramID", ChatID, "Status", expenseEntry.Step)

	data, err := json.Marshal(expenseEntry)
	if err != nil {
		r.logger.Error("Ошибка сериализации данных", "error", err)
		return err
	}

	now := time.Now()
	err = r.rdb.Set(ctx, key, data, time.Hour*1).Err()
	if err != nil {
		r.logger.Error("Ошибка записи в Redis", "error", err, "Duration", time.Since(now))
		return err
	}

	r.logger.Debug("Данные сохранены", "key", key, "Duration", time.Since(now))
	return nil
}

func (r *MemoryRepository) GetExpenseStatus(ctx context.Context, ChatID string) (*status.ExpenseEntry, error) {
	key := "expense:" + ChatID
	r.logger.Debug("Получение данных расхода из Redis", "key", key)
	now := time.Now()
	data, err := r.rdb.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			r.logger.Debug("Данные не найдены", "key", key)
			return nil, nil
		}
		r.logger.Error("Ошибка чтения из Redis", "error", err)
		return nil, err
	}

	var entry status.ExpenseEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		r.logger.Error("Ошибка десериализации", "error", err)
		return nil, err
	}
	r.logger.Debug("Данные получены", "key", key, "Duration", time.Since(now))
	return &entry, nil
}

func (r *MemoryRepository) Delete(ctx context.Context, telegramID string) error {
	key := "expense:" + telegramID
	r.logger.Debug("Удаление статуса пользователя из Redis", "key", key)
	now := time.Now()
	err := r.rdb.Del(ctx, key).Err()
	if err != nil {
		r.logger.Error("Ошибка удаления статуса пользователя из Redis", "key", key, "error", err, "Duration", time.Since(now))
		return err
	}
	r.logger.Debug("Статус пользователя удален", "key", key, "Duration", time.Since(now))
	return nil
}
