package database

import (
	"context"

	"github.com/SobolevTim/finance_bot/internal/domain/user"
)

func (s *DatabaseStore) Create(ctx context.Context, u *user.User) error {
	query := `INSERT INTO users (id, telegram_id, created_at, timezone) 
              VALUES ($1, $2, $3, $4)`

	_, err := s.DB.Exec(ctx, query, u.ID, u.TelegramID, u.CreatedAt, u.Timezone)

	return err
}

func (s *DatabaseStore) GetByTelegramID(ctx context.Context, telegramID string) (*user.User, error) {
	query := `SELECT id, telegram_id, created_at, timezone
			  FROM users
			  WHERE telegram_id = $1`

	row := s.DB.QueryRow(ctx, query, telegramID)
	u := &user.User{}
	err := row.Scan(&u.ID, &u.TelegramID, &u.CreatedAt, &u.Timezone)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s *DatabaseStore) Update(ctx context.Context, u *user.User) error {
	query := `UPDATE users
			  SET timezone = $1
			  WHERE id = $2`

	_, err := s.DB.Exec(ctx, query, u.Timezone, u.ID)
	if err != nil {
		return err
	}
	return nil
}
