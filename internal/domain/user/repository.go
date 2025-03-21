package user

import (
	"context"

	"github.com/google/uuid"
)

// Repository определяет контракт для работы с хранилищем пользователей
type Repository interface {
	UserCreate(ctx context.Context, user *User) error
	UserGetByID(ctx context.Context, id uuid.UUID) (*User, error)
	UserGetByTelegramID(ctx context.Context, telegramID string) (*User, error)
	UserGetByUserName(ctx context.Context, userName string) (*User, error) // Новый метод
	UserUpdate(ctx context.Context, user *User) error
	UserDelete(ctx context.Context, id uuid.UUID) error
}
