package database

import (
	"context"

	"github.com/SobolevTim/finance_bot/internal/domain/budget"
	"github.com/SobolevTim/finance_bot/internal/domain/user"
)

func (s *Repository) Create(ctx context.Context, u *user.User) error {
	query := `INSERT INTO users (id, telegram_id, user_name, first_name, last_name, created_at, timezone) 
              VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := s.DB.Exec(ctx, query, u.ID, u.TelegramID, u.UserName, u.FirstName, u.LastName, u.CreatedAt, u.Timezone)
	if err == nil {
		s.Logger.Debug("User created", "ID", u.ID, "TelegramID", u.TelegramID, "UserName", u.UserName, "FirstName", u.FirstName, "LastName", u.LastName, "CreatedAt", u.CreatedAt, "Timezone", u.Timezone)
	}
	return err
}

func (s *Repository) GetByTelegramID(ctx context.Context, telegramID string) (*user.User, error) {
	query := `SELECT id, telegram_id, user_name, first_name, last_name, created_at, timezone
			  FROM users
			  WHERE telegram_id = $1`

	row := s.DB.QueryRow(ctx, query, telegramID)
	u := &user.User{}
	err := row.Scan(&u.ID, &u.TelegramID, &u.UserName, &u.FirstName, &u.LastName, &u.CreatedAt, &u.Timezone)
	if err != nil {
		return nil, err
	}
	s.Logger.Debug("User found", "ID", u.ID, "TelegramID", u.TelegramID, "UserName", u.UserName, "FirstName", u.FirstName, "LastName", u.LastName, "CreatedAt", u.CreatedAt, "Timezone", u.Timezone)
	return u, nil
}

func (s *Repository) UpdateTimeZone(ctx context.Context, u *user.User) error {
	query := `UPDATE users
			  SET timezone = $1
			  WHERE id = $2`
	_, err := s.DB.Exec(ctx, query, u.Timezone, u.ID)
	if err != nil {
		return err
	}
	s.Logger.Debug("User timezone updated", "ID", u.ID, "Timezone", u.Timezone)
	return nil
}

func (s *Repository) CreateBudget(ctx context.Context, b *budget.Budget) error {
	query := `INSERT INTO budgets (id, telegram_id, amount, currency, created_at) 
			  VALUES ($1, $2, $3, $4, $5)`
	_, err := s.DB.Exec(ctx, query, b.ID, b.TelegramID, b.Amount, b.Currency, b.Date)
	if err == nil {
		s.Logger.Debug("Budget created", "ID", b.ID, "TelegramID", b.TelegramID, "Amount", b.Amount, "Currency", b.Currency, "Date", b.Date)
	}
	return err
}

func (s *Repository) GetBudgetByTelegramID(ctx context.Context, telegramID string) (*budget.Budget, error) {
	query := `SELECT id, telegram_id, amount, currency, created_at, updated_at
			  FROM budgets
			  WHERE telegram_id = $1`

	row := s.DB.QueryRow(ctx, query, telegramID)
	b := &budget.Budget{}
	err := row.Scan(&b.ID, &b.TelegramID, &b.Amount, &b.Currency, &b.Date, &b.UpdateDate)
	if err != nil {
		return nil, err
	}
	s.Logger.Debug("Budget found", "ID", b.ID, "TelegramID", b.TelegramID, "Amount", b.Amount, "Currency", b.Currency, "Date", b.Date, "UpdateDate", b.UpdateDate)
	return b, nil
}
