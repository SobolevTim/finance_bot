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
