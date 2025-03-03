package user

import "context"

type Repository interface {
	Create(ctx context.Context, user *User) error
	GetByTelegramID(ctx context.Context, telegramID string) (*User, error)
	Update(ctx context.Context, user *User) error
}
