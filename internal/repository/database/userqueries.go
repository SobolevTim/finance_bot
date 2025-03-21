package database

import (
	"context"

	"github.com/SobolevTim/finance_bot/internal/domain/user"
	"github.com/google/uuid"
)

func (r *Repository) UserCreate(ctx context.Context, user *user.User) error {
	query := `
		INSERT INTO users (id, telegram_id, user_name, first_name, last_name, timezone, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.DB.Exec(ctx, query, user.ID, user.TelegramID, user.UserName, user.FirstName, user.LastName, user.Timezone, user.CreatedAt, user.UpdatedAt)
	return err
}

func (r *Repository) UserGetByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	query := `
		SELECT id, telegram_id, user_name, first_name, last_name, timezone, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	row := r.DB.QueryRow(ctx, query, id)
	u := &user.User{}
	err := row.Scan(&u.ID, &u.TelegramID, &u.UserName, &u.FirstName, &u.LastName, &u.Timezone, &u.CreatedAt, &u.UpdatedAt)
	return u, err
}

func (r *Repository) UserGetByTelegramID(ctx context.Context, telegramID string) (*user.User, error) {
	query := `
		SELECT id, telegram_id, user_name, first_name, last_name, timezone, created_at, updated_at
		FROM users
		WHERE telegram_id = $1
	`

	row := r.DB.QueryRow(ctx, query, telegramID)
	u := &user.User{}
	err := row.Scan(&u.ID, &u.TelegramID, &u.UserName, &u.FirstName, &u.LastName, &u.Timezone, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}
	return u, nil
}

func (r *Repository) UserGetByUserName(ctx context.Context, userName string) (*user.User, error) {
	query := `
		SELECT id, telegram_id, user_name, first_name, last_name, timezone, created_at, updated_at
		FROM users
		WHERE user_name = $1
	`

	row := r.DB.QueryRow(ctx, query, userName)
	u := &user.User{}
	err := row.Scan(&u.ID, &u.TelegramID, &u.UserName, &u.FirstName, &u.LastName, &u.Timezone, &u.CreatedAt, &u.UpdatedAt)
	return u, err
}

func (r *Repository) UserUpdate(ctx context.Context, user *user.User) error {
	query := `
		UPDATE users
		SET telegram_id = $2, user_name = $3, first_name = $4, last_name = $5, timezone = $6, updated_at = $7
		WHERE id = $1
	`

	_, err := r.DB.Exec(ctx, query, user.ID, user.TelegramID, user.UserName, user.FirstName, user.LastName, user.Timezone, user.UpdatedAt)
	return err
}

func (r *Repository) UserDelete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM users
		WHERE id = $1
	`

	_, err := r.DB.Exec(ctx, query, id)
	return err
}
