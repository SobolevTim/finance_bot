package database

import (
	"context"
	"time"

	"github.com/SobolevTim/finance_bot/internal/domain/user"
	"github.com/google/uuid"
)

func (r *Repository) UserCreate(ctx context.Context, user *user.User) error {
	r.Logger.Debug("UserCreate", "user", user)
	query := `
		INSERT INTO users (id, telegram_id, user_name, first_name, last_name, timezone, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	now := time.Now()
	_, err := r.DB.Exec(ctx, query, user.ID, user.TelegramID, user.UserName, user.FirstName, user.LastName, user.Timezone, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		r.Logger.Debug("UserCreate", "error", err)
		return err
	}
	r.Logger.Debug("UserCreate", "success", true, "timeSince", time.Since(now))
	return nil
}

func (r *Repository) UserGetByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	r.Logger.Debug("UserGetByID", "id", id)
	query := `
		SELECT id, telegram_id, user_name, first_name, last_name, timezone, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	now := time.Now()
	row := r.DB.QueryRow(ctx, query, id)
	u := &user.User{}
	err := row.Scan(&u.ID, &u.TelegramID, &u.UserName, &u.FirstName, &u.LastName, &u.Timezone, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		r.Logger.Debug("UserGetByID", "error", err)
		return nil, err
	}
	r.Logger.Debug("UserGetByID", "success", true, "timeSince", time.Since(now))
	return u, nil
}

func (r *Repository) UserGetByTelegramID(ctx context.Context, telegramID string) (*user.User, error) {
	r.Logger.Debug("UserGetByTelegramID", "telegramID", telegramID)
	query := `
		SELECT id, telegram_id, user_name, first_name, last_name, timezone, created_at, updated_at
		FROM users
		WHERE telegram_id = $1
	`
	now := time.Now()
	row := r.DB.QueryRow(ctx, query, telegramID)
	u := &user.User{}
	err := row.Scan(&u.ID, &u.TelegramID, &u.UserName, &u.FirstName, &u.LastName, &u.Timezone, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		r.Logger.Debug("UserGetByTelegramID", "error", err)
		if err.Error() == "no rows in result set" {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}
	r.Logger.Debug("UserGetByTelegramID", "success", true, "timeSince", time.Since(now))
	return u, nil
}

func (r *Repository) UserGetByUserName(ctx context.Context, userName string) (*user.User, error) {
	r.Logger.Debug("UserGetByUserName", "userName", userName)
	query := `
		SELECT id, telegram_id, user_name, first_name, last_name, timezone, created_at, updated_at
		FROM users
		WHERE user_name = $1
	`
	now := time.Now()
	row := r.DB.QueryRow(ctx, query, userName)
	u := &user.User{}
	err := row.Scan(&u.ID, &u.TelegramID, &u.UserName, &u.FirstName, &u.LastName, &u.Timezone, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		r.Logger.Debug("UserGetByUserName", "error", err)
		if err.Error() == "no rows in result set" {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}
	r.Logger.Debug("UserGetByUserName", "success", true, "timeSince", time.Since(now))
	return u, nil
}

func (r *Repository) UserUpdate(ctx context.Context, user *user.User) error {
	r.Logger.Debug("UserUpdate", "user", user)
	query := `
		UPDATE users
		SET telegram_id = $2, user_name = $3, first_name = $4, last_name = $5, timezone = $6, updated_at = $7
		WHERE id = $1
	`
	now := time.Now()
	_, err := r.DB.Exec(ctx, query, user.ID, user.TelegramID, user.UserName, user.FirstName, user.LastName, user.Timezone, user.UpdatedAt)
	if err != nil {
		r.Logger.Debug("UserUpdate", "error", err)
	}
	r.Logger.Debug("UserUpdate", "success", true, "timeSince", time.Since(now))
	return nil
}

func (r *Repository) UserDelete(ctx context.Context, id uuid.UUID) error {
	r.Logger.Debug("UserDelete", "id", id)
	query := `
		DELETE FROM users
		WHERE id = $1
	`
	now := time.Now()
	_, err := r.DB.Exec(ctx, query, id)
	if err != nil {
		r.Logger.Debug("UserDelete", "error", err)
	}
	r.Logger.Debug("UserDelete", "success", true, "timeSince", time.Since(now))
	return nil
}
